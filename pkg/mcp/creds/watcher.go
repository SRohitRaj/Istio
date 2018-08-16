//  Copyright 2018 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package creds

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"sync"

	"github.com/howeyc/fsnotify"
	"github.com/spf13/cobra"

	"istio.io/istio/pkg/log"
)

var scope = log.RegisterScope("mcp-creds", "MCP Credential utilities", 0)

const (
	// DefaultCertDir is the default directory in which MCP options reside.
	DefaultCertDir = "/etc/istio/certs/"
	// DefaultCertificateFile is the default name to use for the certificate file.
	DefaultCertificateFile = "cert-chain.pem"
	// DefaultKeyFile is the default name to use for the key file.
	DefaultKeyFile = "key.pem"
	// DefaultCACertificateFile is the default name to use for the Certificate Authority's certificate file.
	DefaultCACertificateFile = "root-cert.pem"
)

// CertificateWatcher watches a x509 cert/key file and loads it up in memory as needed.
type CertificateWatcher struct {
	options Options

	stopCh <-chan struct{}

	certMutex sync.Mutex
	cert      tls.Certificate

	// Even though CA cert is not being watched, this type is still responsible for holding on to it
	// to pass into one of the create methods.
	caCertPool *x509.CertPool
}

// WatchFolder loads certificates from the given folder. It expects the
// following files:
// cert-chain.pem, key.pem: Certificate/key files for the client/server on this side.
// root-cert.pem: certificate from the CA that will be used for validating peer's certificate.
//
// Internally WatchFolder will call WatchFiles.
func WatchFolder(stop <-chan struct{}, folder string) (*CertificateWatcher, error) {
	cred := &Options{
		CertificateFile:   path.Join(folder, DefaultCertificateFile),
		KeyFile:           path.Join(folder, DefaultKeyFile),
		CACertificateFile: path.Join(folder, DefaultCACertificateFile),
	}
	return WatchFiles(stop, cred)
}

// Options defines the credential options required for MCP.
type Options struct {
	// CertificateFile to use for mTLS gRPC.
	CertificateFile string
	// KeyFile to use for mTLS gRPC.
	KeyFile string
	// CACertificateFile is the trusted root certificate authority's cert file.
	CACertificateFile string
}

// DefaultOptions returns default credential options.
func DefaultOptions() *Options {
	return &Options{
		CertificateFile:   filepath.Join(DefaultCertDir, DefaultCertificateFile),
		KeyFile:           filepath.Join(DefaultCertDir, DefaultKeyFile),
		CACertificateFile: filepath.Join(DefaultCertDir, DefaultCACertificateFile),
	}
}

// AttachCobraFlags attaches a set of Cobra flags to the given Cobra command.
//
// Cobra is the command-line processor that Istio uses. This command attaches
// the necessary set of flags to configure the MCP options.
func (c *Options) AttachCobraFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&c.CertificateFile, "certFile", "", c.CertificateFile,
		"The location of the certificate file for mutual TLS")
	cmd.PersistentFlags().StringVarP(&c.KeyFile, "keyFile", "", c.KeyFile,
		"The location of the key file for mutual TLS")
	cmd.PersistentFlags().StringVarP(&c.CACertificateFile, "caCertFile", "", c.CACertificateFile,
		"The location of the certificate file for the root certificate authority")
}

// WatchFiles loads certificate & key files from the file system. The method will start a background
// go-routine and watch for credential file changes. Callers should pass the return result to one of the
// create functions to create a transport options that can dynamically use rotated certificates.
// The supplied stop channel can be used to stop the go-routine and the watch.
func WatchFiles(stopCh <-chan struct{}, credentials *Options) (*CertificateWatcher, error) {
	w := &CertificateWatcher{
		options: *credentials,
		stopCh:  stopCh,
	}

	if err := w.start(); err != nil {
		return nil, err
	}

	return w, nil
}

// start watching and stop when the stopCh is closed. Returns an error if the initial load of the certificate
// fails.
func (c *CertificateWatcher) start() error {
	// Load CA Cert file
	caCertPool, err := loadCACert(c.options.CACertificateFile)
	if err != nil {
		return err
	}

	cert, err := loadCertPair(c.options.CertificateFile, c.options.KeyFile)
	if err != nil {
		return err
	}
	c.set(&cert)

	// TODO: https://github.com/istio/istio/issues/7877
	// It looks like fsnotify watchers have problems due to following symlinks. This needs to be handled.
	//
	certFileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	keyFileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err = certFileWatcher.Watch(c.options.CertificateFile); err != nil {
		return err
	}

	if err = keyFileWatcher.Watch(c.options.KeyFile); err != nil {
		_ = certFileWatcher.Close()
		return err
	}

	scope.Debugf("Begin watching certificate files: %s, %s: ",
		c.options.CertificateFile, c.options.KeyFile)

	// Coordinate the goroutines for orderly shutdown
	var exitSignal sync.WaitGroup
	exitSignal.Add(1)
	exitSignal.Add(1)

	go c.watch(exitSignal, certFileWatcher, keyFileWatcher)
	// Watch error events in a separate g See:
	// https://github.com/fsnotify/fsnotify#faq
	go c.watchErrors(exitSignal, certFileWatcher, keyFileWatcher)

	go closeWatchers(exitSignal, certFileWatcher, keyFileWatcher)

	c.caCertPool = caCertPool

	return nil
}

func (c *CertificateWatcher) watch(
	exitSignal sync.WaitGroup, certFileWatcher, keyFileWatcher *fsnotify.Watcher) {

	defer exitSignal.Done()

	for {
		select {
		case e := <-certFileWatcher.Event:
			if e.IsCreate() || e.IsModify() {
				cert, err := loadCertPair(c.options.CertificateFile, c.options.KeyFile)
				if err != nil {
					scope.Errorf("error loading certificates after watch event: %v", err)
				}
				c.set(&cert)
			}
			break

		case e := <-keyFileWatcher.Event:
			if e.IsCreate() || e.IsModify() {
				cert, err := loadCertPair(c.options.CertificateFile, c.options.KeyFile)
				if err != nil {
					scope.Errorf("error loading certificates after watch event: %v", err)
				}
				c.set(&cert)
			}
			break

		case <-c.stopCh:
			scope.Debug("stopping watch of certificate file changes")
			return
		}
	}
}

func (c *CertificateWatcher) watchErrors(
	exitSignal sync.WaitGroup, certFileWatcher, keyFileWatcher *fsnotify.Watcher) {

	defer exitSignal.Done()

	for {
		select {
		case e := <-keyFileWatcher.Error:
			scope.Errorf("error event while watching key file: %v", e)
			break

		case e := <-certFileWatcher.Error:
			scope.Errorf("error event while watching cert file: %v", e)
			break

		case <-c.stopCh:
			scope.Debug("stopping watch of certificate file errors")
			return
		}
	}
}

func closeWatchers(exitSignal sync.WaitGroup, certFileWatcher, keyFileWatcher *fsnotify.Watcher) {
	exitSignal.Wait()
	_ = certFileWatcher.Close()
	_ = keyFileWatcher.Close()
}

// set the certificate directly
func (c *CertificateWatcher) set(cert *tls.Certificate) {
	c.certMutex.Lock()
	defer c.certMutex.Unlock()
	c.cert = *cert
}

// get the currently loaded certificate.
func (c *CertificateWatcher) get() tls.Certificate {
	c.certMutex.Lock()
	defer c.certMutex.Unlock()
	return c.cert
}

// loadCertPair from the given set of files.
func loadCertPair(certFile, keyFile string) (tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		err = fmt.Errorf("error loading client certificate files (%s, %s): %v", certFile, keyFile, err)
	}

	return cert, err
}

// loadCACert, create a certPool and return.
func loadCACert(caCertFile string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		err = fmt.Errorf("error loading CA certificate file (%s): %v", caCertFile, err)
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err = errors.New("failed to append loaded CA certificate to the certificate pool")
		return nil, err
	}

	return certPool, nil
}
