// Copyright 2018 Istio Authors
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

package model

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	authn "istio.io/api/authentication/v1alpha1"
	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pkg/log"
)

const (
	// CARootCertPath is the path of ca root cert that envoy uses to validate cert got from SDS service.
	CARootCertPath = "/etc/istio/certs/ca-root-cert.pem"

	// SDSStatPrefix is the human readable prefix to use when emitting statistics for the SDS service.
	SDSStatPrefix = "sdsstat"
)

// JwtKeyResolver resolves JWT public key and JwksURI.
var JwtKeyResolver = newJwksResolver(JwtPubKeyExpireDuration, JwtPubKeyEvictionDuration, JwtPubKeyRefreshInterval)

// GetConsolidateAuthenticationPolicy returns the authentication policy for
// service specified by hostname and port, if defined.
// If not, it generates and output a policy that is equivalent to the legacy flag
// and/or service annotation. Once these legacy flags/config deprecated,
// this function can be placed by a call to store.AuthenticationPolicyByDestination
// directly.
func GetConsolidateAuthenticationPolicy(mesh *meshconfig.MeshConfig, store IstioConfigStore, hostname Hostname, port *Port) *authn.Policy {
	config := store.AuthenticationPolicyByDestination(hostname, port)
	if config != nil {
		policy := config.Spec.(*authn.Policy)
		if err := JwtKeyResolver.SetAuthenticationPolicyJwksURIs(policy); err == nil {
			return policy
		}
	}

	log.Debugf("No authentication policy found for  %s:%d. Fallback to legacy authentication mode %v\n",
		hostname, port.Port, mesh.AuthPolicy)
	return legacyAuthenticationPolicyToPolicy(mesh.AuthPolicy)
}

// ConstructSdsSecretConfig constructs SDS Sececret Configuration.
func ConstructSdsSecretConfig(serviceAccount string, refreshDuration *time.Duration, sdsUdsPath string) *auth.SdsSecretConfig {
	if serviceAccount == "" || sdsUdsPath == "" || refreshDuration == nil {
		return nil
	}

	return &auth.SdsSecretConfig{
		Name: serviceAccount,
		SdsConfig: &core.ConfigSource{
			ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
				ApiConfigSource: &core.ApiConfigSource{
					ApiType: core.ApiConfigSource_GRPC,
					GrpcServices: []*core.GrpcService{
						{
							TargetSpecifier: &core.GrpcService_GoogleGrpc_{
								GoogleGrpc: &core.GrpcService_GoogleGrpc{
									TargetUri:  sdsUdsPath,
									StatPrefix: SDSStatPrefix,
								},
							},
						},
					},
					RefreshDelay: refreshDuration,
				},
			},
		},
	}
}

// ConstructValidationContext constructs ValidationContext in CommonTlsContext.
func ConstructValidationContext(rootCAFilePath string) *auth.CommonTlsContext_ValidationContext {
	return &auth.CommonTlsContext_ValidationContext{
		ValidationContext: &auth.CertificateValidationContext{
			TrustedCa: &core.DataSource{
				Specifier: &core.DataSource_Filename{
					Filename: rootCAFilePath,
				},
			},
		},
	}
}

// If input legacy is MeshConfig_MUTUAL_TLS, return a authentication policy equivalent to it. Else,
// returns nil (implies no authentication is used)
func legacyAuthenticationPolicyToPolicy(legacy meshconfig.MeshConfig_AuthPolicy) *authn.Policy {
	if legacy == meshconfig.MeshConfig_MUTUAL_TLS {
		return &authn.Policy{
			Peers: []*authn.PeerAuthenticationMethod{{
				Params: &authn.PeerAuthenticationMethod_Mtls{
					&authn.MutualTls{},
				}}},
		}
	}
	return nil
}

// ParseJwksURI parses the input URI and returns the corresponding hostname, port, and whether SSL is used.
// URI must start with "http://" or "https://", which corresponding to "http" or "https" scheme.
// Port number is extracted from URI if available (i.e from postfix :<port>, eg. ":80"), or assigned
// to a default value based on URI scheme (80 for http and 443 for https).
// Port name is set to URI scheme value.
// Note: this is to replace [buildJWKSURIClusterNameAndAddress]
// (https://github.com/istio/istio/blob/master/pilot/pkg/proxy/envoy/v1/mixer.go#L401),
// which is used for the old EUC policy.
func ParseJwksURI(jwksURI string) (string, *Port, bool, error) {
	u, err := url.Parse(jwksURI)
	if err != nil {
		return "", nil, false, err
	}
	var useSSL bool
	var portNumber int
	switch u.Scheme {
	case "http":
		useSSL = false
		portNumber = 80
	case "https":
		useSSL = true
		portNumber = 443
	default:
		return "", nil, false, fmt.Errorf("URI scheme %q is not supported", u.Scheme)
	}

	if u.Port() != "" {
		portNumber, err = strconv.Atoi(u.Port())
		if err != nil {
			return "", nil, useSSL, err
		}
	}

	return u.Hostname(), &Port{
		Name: u.Scheme,
		Port: portNumber,
	}, useSSL, nil
}
