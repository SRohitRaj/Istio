// Copyright 2017 Istio Authors
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

// Provides utility methods to generate X.509 certificates with different
// options. This implementation is Largely inspired from
// https://golang.org/src/crypto/tls/generate_cert.go.

package ca

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	_ "github.com/golang/glog"
)

// GenCSR generates a X.509 certificate sign request and private key with the given options.
func GenCSR(options CertOptions) ([]byte, []byte, error) {
	// Generates a CSR
	priv, err := rsa.GenerateKey(rand.Reader, options.RSAKeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("CSR generation fails at RSA key generation (%v)", err)
	}
	template := GenCSRTemplate(options)
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, crypto.PrivateKey(priv))
	if err != nil {
		return nil, nil, fmt.Errorf("CSR generation fails at X509 cert request generation (%v)", err)
	}

	csr, privKey := encodePem(true, csrBytes, priv)
	return csr, privKey, nil
}

// GenCSRTemplate generates a certificateRequest template with the given options.
func GenCSRTemplate(options CertOptions) x509.CertificateRequest {
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Organization: []string{options.Org},
		},
	}

	if h := options.Host; len(h) > 0 {
		s := buildSubjectAltNameExtension(h)
		template.ExtraExtensions = []pkix.Extension{*s}
	}

	return template
}
