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

package spiffe

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"gopkg.in/square/go-jose.v2"

	"istio.io/istio/pkg/config/constants"
	"istio.io/pkg/log"
)

const (
	Scheme = "spiffe"

	URIPrefix = Scheme + "://"

	// The default SPIFFE URL value for trust domain
	defaultTrustDomain = constants.DefaultKubernetesDomain
)

var (
	trustDomain      = defaultTrustDomain
	trustDomainMutex sync.RWMutex

	spiffeLog = log.RegisterScope("spiffe", "SPIFFE library logging", 0)
)

type spiffeConfig struct {
	Map map[string]string `json:"map"`
}

type bundleDoc struct {
	jose.JSONWebKeySet
	Sequence    uint64 `json:"spiffe_sequence,omitempty"`
	RefreshHint int    `json:"spiffe_refresh_hint,omitempty"`
}

func SetTrustDomain(value string) {
	// Replace special characters in spiffe
	v := strings.Replace(value, "@", ".", -1)
	trustDomainMutex.Lock()
	trustDomain = v
	trustDomainMutex.Unlock()
}

func GetTrustDomain() string {
	trustDomainMutex.RLock()
	defer trustDomainMutex.RUnlock()
	return trustDomain
}

func DetermineTrustDomain(commandLineTrustDomain string, isKubernetes bool) string {
	if len(commandLineTrustDomain) != 0 {
		return commandLineTrustDomain
	}
	if isKubernetes {
		return defaultTrustDomain
	}
	return ""
}

// GenSpiffeURI returns the formatted uri(SPIFFE format for now) for the certificate.
func GenSpiffeURI(ns, serviceAccount string) (string, error) {
	var err error
	if ns == "" || serviceAccount == "" {
		err = fmt.Errorf(
			"namespace or service account empty for SPIFFE uri ns=%v serviceAccount=%v", ns, serviceAccount)
	}
	return URIPrefix + GetTrustDomain() + "/ns/" + ns + "/sa/" + serviceAccount, err
}

// MustGenSpiffeURI returns the formatted uri(SPIFFE format for now) for the certificate and logs if there was an error.
func MustGenSpiffeURI(ns, serviceAccount string) string {
	uri, err := GenSpiffeURI(ns, serviceAccount)
	if err != nil {
		spiffeLog.Debug(err.Error())
	}
	return uri
}

// GenCustomSpiffe returns the  spiffe string that can have a custom structure
func GenCustomSpiffe(identity string) string {
	if identity == "" {
		spiffeLog.Error("spiffe identity can't be empty")
		return ""
	}

	return URIPrefix + GetTrustDomain() + "/" + identity
}

// RetrieveSpiffeBundleRootCertsFromStringInput retrieves the trusted CA certificates from a list of SPIFFE bundle endpoints.
// It can use the system cert pool and the supplied certificates to validate the endpoints.
// The input endpointTuples should be in the json format of:
//		{"map": {"foo": "URL1", "bar": "URL2"}}
func RetrieveSpiffeBundleRootCertsFromStringInput(inputString string, extraTrustedCerts []*x509.Certificate) (
	map[string][]*x509.Certificate, error) {
	var converted spiffeConfig
	err := json.Unmarshal([]byte(inputString), &converted)
	if err != nil {
		return nil, fmt.Errorf("error unmarchalling the input: %v", err)
	}
	spiffeLog.Infof("Processed SPIFFE bundle configuration: %v", converted)
	return RetrieveSpiffeBundleRootCerts(converted, extraTrustedCerts)
}

// RetrieveSpiffeBundleRootCerts retrieves the trusted CA certificates from a list of SPIFFE bundle endpoints.
// It can use the system cert pool and the supplied certificates to validate the endpoints.
func RetrieveSpiffeBundleRootCerts(config spiffeConfig, extraTrustedCerts []*x509.Certificate) (
	map[string][]*x509.Certificate, error) {
	httpClient := &http.Client{}
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to get SystemCertPool: %v", err)
	}
	for _, cert := range extraTrustedCerts {
		caCertPool.AddCert(cert)
	}

	ret := map[string][]*x509.Certificate{}
	for trustdomain, endpoint := range config.Map {
		if !strings.HasPrefix(endpoint, "https://") {
			endpoint = "https://" + endpoint
		}
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to split the SPIFFE bundle URL: %v", err)
		}

		config := &tls.Config{
			ServerName: u.Hostname(),
			RootCAs:    caCertPool,
		}

		httpClient.Transport = &http.Transport{
			TLSClientConfig: config,
		}

		resp, err := httpClient.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bundle: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b := make([]byte, 1024)
			n, _ := resp.Body.Read(b)
			return nil, fmt.Errorf("unexpected status %d fetching bundle: %s", resp.StatusCode, string(b[:n]))
		}

		doc := new(bundleDoc)
		if err := json.NewDecoder(resp.Body).Decode(doc); err != nil {
			return nil, fmt.Errorf("failed to decode bundle: %v", err)
		}

		var cert *x509.Certificate
		for i, key := range doc.Keys {
			if key.Use == "x509-svid" {
				if len(key.Certificates) != 1 {
					return nil, fmt.Errorf("expected a single certificate in x509-svid entry %d; got %d", i, len(key.Certificates))
				}
				cert = key.Certificates[0]
			}
		}
		if cert == nil {
			return nil, fmt.Errorf("trust domain [%s] at URL [%s] does not provide a X509 SVID", trustdomain, endpoint)
		}
		if certs, ok := ret[trustdomain]; ok {
			ret[trustdomain] = append(certs, cert)
		} else {
			ret[trustdomain] = []*x509.Certificate{cert}
		}
	}
	return ret, nil
}
