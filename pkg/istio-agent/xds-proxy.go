// Copyright 2020 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istioagent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"istio.io/pkg/log"
)

const (
	defaultClientMaxReceiveMessageSize = 1024 * 1024 * 8
	defaultInitialConnWindowSize       = 1024 * 1024 // default gRPC InitialWindowSize
	defaultInitialWindowSize           = 1024 * 1024 // default gRPC ConnWindowSize
)

type XdsProxy struct {
	stopChan             chan struct{}
	downstreamListener   net.Listener
	downstreamGrpcServer *grpc.Server
	upstreamAddress      string
	upstreamDialOptions  []grpc.DialOption
}

// XDS Proxy proxies all XDS requests from envoy to istiod, in addition to allowing subsystems inside
// the agent to also communicate with either istiod/envoy (eg dns, sds, etc).
// The goal here is to consolidate all xds related connections to istiod/envoy into a
// single tcp connection with multiple gRPC streams.
// TODO: Right now, the workloadSDS server and gatewaySDS servers are still separate
// connections. These need to be consolidated
func initXdsProxy(sa *Agent) (*XdsProxy, error) {
	var err error
	proxy := &XdsProxy{
		upstreamAddress: sa.proxyConfig.DiscoveryAddress,
	}

	if err = proxy.initDownstreamServer(); err != nil {
		return nil, err
	}

	if proxy.upstreamDialOptions, err = buildUpstreamClientDialOpts(sa); err != nil {
		return nil, err
	}

	go func() {
		_ = proxy.downstreamGrpcServer.Serve(proxy.downstreamListener)
	}()

	return proxy, nil
}

// Every time envoy makes a fresh connection to the agent, we reestablish a new connection to the upstream xds
// This ensures that a new connection between istiod and agent doesn't end up consuming pending messages from envoy
// as the new connection may not go to the same istiod. Vice versa case also applies.
func (p *XdsProxy) StreamAggregatedResources(downstream discovery.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	errChan := make(chan error)
	requestsChan := make(chan *discovery.DiscoveryRequest)
	responsesChan := make(chan *discovery.DiscoveryResponse)

	log.Infof("connecting to upstream %s", p.upstreamAddress)
	upstreamConn, err := grpc.Dial(p.upstreamAddress, p.upstreamDialOptions...)
	if err != nil {
		log.Errorf("failed to connect to upstream %s: %v", p.upstreamAddress, err)
		return err
	}

	xds := discovery.NewAggregatedDiscoveryServiceClient(upstreamConn)
	upstream, err := xds.StreamAggregatedResources(context.Background())
	if err != nil {
		log.Errorf("failed to create upstream grpc client: %v", err)
		return err
	}

	go func() {
		for {
			// from istiod
			resp, err := upstream.Recv()
			if err != nil {
				log.Errorf("upstream recv error: %v", err)
				errChan <- err
				return
			}
			responsesChan <- resp
		}
	}()
	go func() {
		for {
			// From Envoy
			req, err := downstream.Recv()
			if err != nil {
				log.Errorf("downstream recv error: %v", err)
				errChan <- err
				return
			}
			// forward to istiod
			requestsChan <- req
		}
	}()

	for {
		select {
		case err := <-errChan:
			// error receiving from downstream envoy
			// recycle connection.
			_ = upstream.CloseSend()
			// todo close downstream?
			return err
		case req := <-requestsChan:
			if err = upstream.Send(req); err != nil {
				log.Errorf("upstream send error: %v", err)
				return err
			}
		case resp := <-responsesChan:
			if err := downstream.Send(resp); err != nil {
				log.Errorf("downstream send error: %v", err)
				// we cannot return partial error and hope to restart just the downstream
				// as we are blindly proxying req/responses. For now, the best course of action
				// is to terminate upstream connection as well and restart afresh.
				return err
			}
		case <-p.stopChan:
			_ = upstream.CloseSend()
			return nil
		}
	}
}

func (p *XdsProxy) DeltaAggregatedResources(server discovery.AggregatedDiscoveryService_DeltaAggregatedResourcesServer) error {
	return errors.New("delta XDS is not implemented")
}

func (p *XdsProxy) close() {
	p.stopChan <- struct{}{}
	if p.downstreamGrpcServer != nil {
		_ = p.downstreamGrpcServer.Stop
	}
	if p.downstreamListener != nil {
		_ = p.downstreamListener.Close()
	}
}

// TODO reuse code from SDS
// TODO reuse the connection, not just code
func setUpUds(udsPath string) (net.Listener, error) {
	// Remove unix socket before use.
	if err := os.Remove(udsPath); err != nil && !os.IsNotExist(err) {
		// Anything other than "file not found" is an error.
		log.Errorf("Failed to remove unix://%s: %v", udsPath, err)
		return nil, fmt.Errorf("failed to remove unix://%s", udsPath)
	}

	// Attempt to create the folder in case it doesn't exist
	if err := os.MkdirAll(filepath.Dir(udsPath), 0750); err != nil {
		// If we cannot create it, just warn here - we will fail later if there is a real error
		log.Warnf("Failed to create directory for %v: %v", udsPath, err)
	}

	var err error
	udsListener, err := net.Listen("unix", udsPath)
	if err != nil {
		log.Errorf("Failed to listen on unix socket %q: %v", udsPath, err)
		return nil, err
	}

	// Update SDS UDS file permission so that istio-proxy has permission to access it.
	if _, err := os.Stat(udsPath); err != nil {
		log.Errorf("SDS uds file %q doesn't exist", udsPath)
		return nil, fmt.Errorf("sds uds file %q doesn't exist", udsPath)
	}
	if err := os.Chmod(udsPath, 0666); err != nil {
		log.Errorf("Failed to update %q permission", udsPath)
		return nil, fmt.Errorf("failed to update %q permission", udsPath)
	}

	return udsListener, nil
}

type fileTokenSource struct {
	path   string
	period time.Duration
}

var _ = oauth2.TokenSource(&fileTokenSource{})

func (ts *fileTokenSource) Token() (*oauth2.Token, error) {
	tokb, err := ioutil.ReadFile(ts.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file %q: %v", ts.path, err)
	}
	tok := strings.TrimSpace(string(tokb))
	if len(tok) == 0 {
		return nil, fmt.Errorf("read empty token from file %q", ts.path)
	}

	return &oauth2.Token{
		AccessToken: tok,
		Expiry:      time.Now().Add(ts.period),
	}, nil
}

func (p *XdsProxy) initDownstreamServer() error {
	l, err := setUpUds("./etc/istio/proxy/XDS")
	if err != nil {
		return err
	}
	grpcs := grpc.NewServer()
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcs, p)
	reflection.Register(grpcs)
	p.downstreamGrpcServer = grpcs
	p.downstreamListener = l
	return nil
}

func buildUpstreamClientDialOpts(sa *Agent) ([]grpc.DialOption, error) {
	tlsOpts, err := getTLSDialOption(sa)
	if err != nil {
		return nil, fmt.Errorf("failed to build TLS dial option to talk to upstream: %v", err)
	}

	keepaliveOption := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:    30 * time.Second,
		Timeout: 10 * time.Second,
	})

	initialWindowSizeOption := grpc.WithInitialWindowSize(int32(defaultInitialWindowSize))
	initialConnWindowSizeOption := grpc.WithInitialConnWindowSize(int32(defaultInitialConnWindowSize))
	msgSizeOption := grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultClientMaxReceiveMessageSize))
	dialOptions := []grpc.DialOption{
		tlsOpts,
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 5 * time.Second,
		}),
		keepaliveOption, initialWindowSizeOption, initialConnWindowSizeOption, msgSizeOption,
	}

	// TODO: This is not a valid way of detecting if we are on VM vs k8s
	// Some end users do not use Istiod for CA but run on k8s with file mounted certs
	// In these cases, while we fallback to mTLS to istiod using the provisioned certs
	// it would be ideal to keep using token plus k8s ca certs for control plane communication
	// as the intention behind provisioned certs on k8s pods is only for data plane comm.
	if sa.secOpts.ProvCert == "" {
		// only if running in k8s pod
		dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(oauth.TokenSource{&fileTokenSource{
			"./var/run/secrets/tokens/istio-token",
			time.Second * 300,
		}}))
	}
	return dialOptions, nil
}

// Returns the TLS option to use when talking to Istiod
// If provisioned cert is set, it will return a mTLS related config
// Else it will return a one-way TLS related config with the assumption
// that the consumer code will use tokens to authenticate the upstream.
func getTLSDialOption(agent *Agent) (grpc.DialOption, error) {
	var certPool *x509.CertPool
	var err error
	var rootCert []byte
	if agent.secOpts.ProvCert != "" {
		// This is most likely a VM using pre-provisioned mTLS certs (key.pem, cert-chain.pem, root-cert.pem)
		// to talk to Istiod setup mtls
		rootCert, err = ioutil.ReadFile(agent.secOpts.ProvCert + "/root-cert.pem")
	} else if agent.secOpts.PilotCertProvider != "" {
		rootCert, err = agent.loadPilotCertProviderRootCert()
	}

	if err != nil {
		return nil, err
	}

	certPool = x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(rootCert)
	if !ok {
		return nil, fmt.Errorf("failed to create TLS dial option with root certificates")
	}

	var certificate tls.Certificate
	config := tls.Config{
		Certificates: []tls.Certificate{certificate},
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			if agent.secOpts.ProvCert != "" {
				// Load the certificate from disk
				if certificate, err = tls.LoadX509KeyPair(agent.secOpts.ProvCert+"/cert-chain.pem", agent.secOpts.ProvCert+"/key.pem"); err != nil {
					return nil, err
				}
			}
			return &certificate, nil
		},
	}
	config.RootCAs = certPool
	// strip the port from the address
	parts := strings.Split(agent.proxyConfig.DiscoveryAddress, ":")
	config.ServerName = parts[0]
	config.MinVersion = tls.VersionTLS12
	transportCreds := credentials.NewTLS(&config)
	return grpc.WithTransportCredentials(transportCreds), nil
}
