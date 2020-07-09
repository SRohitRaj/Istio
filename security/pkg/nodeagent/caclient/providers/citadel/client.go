// Copyright Istio Authors
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

package caclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/peer"
	"istio.io/pkg/env"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	caClientInterface "istio.io/istio/security/pkg/nodeagent/caclient/interface"
	pb "istio.io/istio/security/proto"
	"istio.io/pkg/log"
)

const (
	caServerName      = "istio-citadel"
	bearerTokenPrefix = "Bearer "
)

var (
	citadelClientLog = log.RegisterScope("citadelclient", "citadel client debugging", 0)

	// ProvCert is the environment controlling the use of pre-provisioned certs, for VMs.
	// May also be used in K8S to use a Secret to bootstrap (as a 'refresh key'), but use short-lived tokens
	// with extra SAN (labels, etc) in data path.
	ProvCert = env.RegisterStringVar("PROV_CERT", "",
		"Set to a directory containing provisioned certs, for VMs").Get()

	// EMPTY_PROV_CERT needs to be set true if the ProvCert path is provided
	// but cert not exists
	EmptyProvCert = env.RegisterStringVar("EMPTY_PROV_CERT", "",
		"Set to a directory containing provisioned certs, for VMs").Get()
)

type citadelClient struct {
	caEndpoint    string
	enableTLS     bool
	caTLSRootCert []byte
	client        pb.IstioCertificateServiceClient
	clusterID     string
}

// NewCitadelClient create a CA client for Citadel.
func NewCitadelClient(endpoint string, tls bool, rootCert []byte, clusterID string) (caClientInterface.Client, error) {
	c := &citadelClient{
		caEndpoint:    endpoint,
		enableTLS:     tls,
		caTLSRootCert: rootCert,
		clusterID:     clusterID,
	}

	var opts grpc.DialOption
	var err error
	if tls {
		opts, err = c.getTLSDialOption()
		if err != nil {
			return nil, err
		}
	} else {
		opts = grpc.WithInsecure()
	}

	citadelClientLog.Infof("3333333333\n")
	citadelClientLog.Infof("%+v\n", ProvCert)
	citadelClientLog.Infof("4444444444\n")

	// TODO(JimmyCYJ): This connection is create at construction time. If conn is broken at anytime,
	//  need a way to reconnect.
	conn, err := grpc.Dial(endpoint, opts)
	if err != nil {
		citadelClientLog.Errorf("Failed to connect to endpoint %s: %v", endpoint, err)
		return nil, fmt.Errorf("failed to connect to endpoint %s", endpoint)
	}

	c.client = pb.NewIstioCertificateServiceClient(conn)
	return c, nil
}

// CSR Sign calls Citadel to sign a CSR.
func (c *citadelClient) CSRSign(ctx context.Context, reqID string, csrPEM []byte, token string,
	certValidTTLInSec int64, withToken bool) ([]string /*PEM-encoded certificate chain*/, error) {
	req := &pb.IstioCertificateRequest{
		Csr:              string(csrPEM),
		ValidityDuration: certValidTTLInSec,
	}
	citadelClientLog.Infof("CSRSign========\n")
	citadelClientLog.Infof("%+v\n",ctx)
	peer, _ := peer.FromContext(ctx)
	citadelClientLog.Infof("ssssssppppoooooo\n")
	citadelClientLog.Infof("%+v\n",ctx)
	citadelClientLog.Infof("pppppkkkkkkkkk\n")
	citadelClientLog.Infof("%+v\n",peer)
	citadelClientLog.Infof("nnnnnnnnnnnnnnnnn\n")

	// add Bearer prefix, which is required by Citadel.
	token = bearerTokenPrefix + token
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("Authorization", token, "ClusterID", c.clusterID))

	resp, err := c.client.CreateCertificate(ctx, req)
	if err != nil {
		citadelClientLog.Errorf("Failed to create certificate: %v", err)
		return nil, err
	}

	if len(resp.CertChain) <= 1 {
		citadelClientLog.Errorf("CertChain length is %d, expected more than 1", len(resp.CertChain))
		return nil, errors.New("invalid response cert chain")
	}

	return resp.CertChain, nil
}

func (c *citadelClient) getTLSDialOption() (grpc.DialOption, error) {
	// Load the TLS root certificate from the specified file.
	// Create a certificate pool
	var certPool *x509.CertPool
	var err error
	citadelClientLog.Infof("666666666\n")
	citadelClientLog.Infof("%+v", EmptyProvCert)
	citadelClientLog.Infof("777777777\n")
	citadelClientLog.Infof("%+v", ProvCert)
	if c.caTLSRootCert == nil {
		// No explicit certificate - assume the citadel-compatible server uses a public cert
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
	} else {
		certPool = x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(c.caTLSRootCert)
		if !ok {
			return nil, fmt.Errorf("failed to append certificates")
		}
	}

	var certificate tls.Certificate
	config := tls.Config{
		Certificates: []tls.Certificate{certificate},
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			if ProvCert != "" && EmptyProvCert != "true" {
				// Load the certificate from disk
				certificate, err = tls.LoadX509KeyPair(ProvCert+"/cert-chain.pem", ProvCert+"/key.pem")
				if err != nil {
					return nil, fmt.Errorf("cannot load key pair: %s", err)
				}
			}
			return &certificate, nil
		},
	}
	config.RootCAs = certPool

	// Initial implementation of citadel hardcoded the SAN to 'istio-citadel'. For backward compat, keep it.
	// TODO: remove this once istiod replaces citadel.
	// External CAs will use their normal server names.
	if strings.Contains(c.caEndpoint, "citadel") {
		config.ServerName = caServerName
	}
	// For debugging on localhost (with port forward)
	// TODO: remove once istiod is stable and we have a way to validate JWTs locally
	if strings.Contains(c.caEndpoint, "localhost") {
		config.ServerName = "istiod.istio-system.svc"
	}

	transportCreds := credentials.NewTLS(&config)
	return grpc.WithTransportCredentials(transportCreds), nil
}
