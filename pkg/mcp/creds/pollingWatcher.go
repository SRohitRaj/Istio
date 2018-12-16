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
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

type pollingWatcher struct {
	options      Options
	pollInterval time.Duration

	stopCh <-chan struct{}

	certMutex sync.Mutex
	cert      tls.Certificate

	// Even though CA cert is not being watched, this type is still responsible for holding on to it
	// to pass into one of the create methods.
	caCertPool *x509.CertPool

	certModTime time.Time
	keyModTime  time.Time

	// Keep the current error encountered when loading cert files while polling. This helps with testing.
	pollErr error
}

var _ CertificateWatcher = &pollingWatcher{}

func (p *pollingWatcher) certPool() *x509.CertPool {
	return p.caCertPool
}

// PollFolder loads certificates from the given folder. It expects the
// following files:
// cert-chain.pem, key.pem: Certificate/key files for the client/server on this side.
// root-cert.pem: certificate from the CA that will be used for validating peer's certificate.
//
// Internally PollFolder will call PollFiles.
func PollFolder(stop <-chan struct{}, folder string) (CertificateWatcher, error) {
	cred := &Options{
		CertificateFile:   path.Join(folder, defaultCertificateFile),
		KeyFile:           path.Join(folder, defaultKeyFile),
		CACertificateFile: path.Join(folder, defaultCACertificateFile),
	}
	return PollFiles(stop, cred)
}

// PollFiles loads certificate & key files from the file system. The method will start a background
// go-routine and watch for credential file changes. Callers should pass the return result to one of the
// create functions to create a transport options that can dynamically use rotated certificates.
// The supplied stop channel can be used to stop the go-routine and the watch.
func PollFiles(stopCh <-chan struct{}, credentials *Options) (CertificateWatcher, error) {
	return pollFiles(stopCh, credentials, time.Minute)
}

func pollFiles(stopCh <-chan struct{}, credentials *Options, interval time.Duration) (CertificateWatcher, error) {
	w := &pollingWatcher{
		options:      *credentials,
		pollInterval: interval,
		stopCh:       stopCh,
	}

	if err := w.start(); err != nil {
		return nil, err
	}

	return w, nil
}

// start watching and stop when the stopCh is closed. Returns an error if the initial load of the certificate
// fails.
func (p *pollingWatcher) start() error {
	// Load CA Cert file
	caCertPool, err := loadCACert(p.options.CACertificateFile)
	if err != nil {
		return err
	}

	if err = p.loadFiles(); err != nil {
		return err
	}

	scope.Debugf("Begin polling certificate files: %s, %s: ",
		p.options.CertificateFile, p.options.KeyFile)

	go p.poll()

	p.caCertPool = caCertPool

	return nil
}

func (p *pollingWatcher) poll() {
	t := time.NewTicker(p.pollInterval)
	for {
		select {
		case <-t.C:
			err := p.loadFiles()
			if err != nil {
				scope.Errorf("Error polling certificate files: %v", err)
			}

		case <-p.stopCh:
			t.Stop()
			scope.Debug("stopping poll of certificate file changes")
			return
		}
	}
}

func (p *pollingWatcher) loadFiles() (err error) {
	p.certMutex.Lock()
	defer p.certMutex.Unlock()

	defer func() {
		p.pollErr = err
	}()

	var newKeyModTime, newCertModTime time.Time

	var fi os.FileInfo

	// Go through files and stat.
	if fi, err = os.Stat(p.options.KeyFile); err != nil {
		err = fmt.Errorf("unable to read key file(%q): %v", p.options.KeyFile, err)
		return
	}
	newKeyModTime = fi.ModTime()

	if fi, err = os.Stat(p.options.CertificateFile); err != nil {
		err = fmt.Errorf("unable to read cert file(%q): %v", p.options.CertificateFile, err)
		return
	}
	newCertModTime = fi.ModTime()

	if !newKeyModTime.Equal(p.keyModTime) || !newCertModTime.Equal(p.certModTime) {
		var cert tls.Certificate
		cert, err = loadCertPair(p.options.CertificateFile, p.options.KeyFile)
		if err != nil {
			err = fmt.Errorf("unable load cert files as pair: %v", err)
			return
		}

		p.cert = cert
		p.keyModTime = newKeyModTime
		p.certModTime = newCertModTime
	}

	return
}

func (p *pollingWatcher) pollError() error {
	p.certMutex.Lock()
	defer p.certMutex.Unlock()
	return p.pollErr
}

// Get the currently loaded certificate.
func (p *pollingWatcher) Get() tls.Certificate {
	p.certMutex.Lock()
	defer p.certMutex.Unlock()
	return p.cert
}
