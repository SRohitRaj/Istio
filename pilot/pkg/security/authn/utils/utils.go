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

package utils

import (
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"

	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking"
	"istio.io/istio/pilot/pkg/networking/util"
	authn_model "istio.io/istio/pilot/pkg/security/model"
	protovalue "istio.io/istio/pkg/proto"
)

// SupportedCiphers for server side TLS configuration.
var SupportedCiphers = []string{
	"ECDHE-ECDSA-AES256-GCM-SHA384",
	"ECDHE-RSA-AES256-GCM-SHA384",
	"ECDHE-ECDSA-AES128-GCM-SHA256",
	"ECDHE-RSA-AES128-GCM-SHA256",
	"AES256-GCM-SHA384",
	"AES128-GCM-SHA256",
}

// BuildInboundTLS returns the TLS context corresponding to the mTLS mode.
func BuildInboundTLS(mTLSMode model.MutualTLSMode, node *model.Proxy,
	protocol networking.ListenerProtocol, trustDomainAliases []string, mc *meshconfig.MeshConfig,
) *tls.DownstreamTlsContext {
	if mTLSMode == model.MTLSDisable || mTLSMode == model.MTLSUnknown {
		return nil
	}
	ctx := &tls.DownstreamTlsContext{
		CommonTlsContext:         &tls.CommonTlsContext{},
		RequireClientCertificate: protovalue.BoolTrue,
	}
	if protocol == networking.ListenerProtocolTCP {
		// For TCP with mTLS, we advertise "istio-peer-exchange" from client and
		// expect the same from server. This  is so that secure metadata exchange
		// transfer can take place between sidecars for TCP with mTLS.
		if features.MetadataExchange {
			ctx.CommonTlsContext.AlpnProtocols = util.ALPNDownstreamWithMxc
		} else {
			ctx.CommonTlsContext.AlpnProtocols = util.ALPNDownstream
		}
	} else {
		// Note that in the PERMISSIVE mode, we match filter chain on "istio" ALPN,
		// which is used to differentiate between service mesh and legacy traffic.
		//
		// Client sidecar outbound cluster's TLSContext.ALPN must include "istio".
		//
		// Server sidecar filter chain's FilterChainMatch.ApplicationProtocols must
		// include "istio" for the secure traffic, but its TLSContext.ALPN must not
		// include "istio", which would interfere with negotiation of the underlying
		// protocol, e.g. HTTP/2.
		ctx.CommonTlsContext.AlpnProtocols = util.ALPNHttp
	}

	// Set Minimum TLS version to match the default client version and allowed strong cipher suites for sidecars.
	ctx.CommonTlsContext.TlsParams = &tls.TlsParameters{
		CipherSuites: SupportedCiphers,
	}
	// Configure TLS version based on meshconfig TLS API.
	// The min TLS version is initialized to be at least TLS 1.2.
	minTLSVersion := meshconfig.MeshConfig_TLSConfig_TLSV1_2
	if mc != nil && mc.GetMeshMTLS() != nil {
		minTLSVersion = mc.GetMeshMTLS().MinProtocolVersion
	}
	configureTLSVersion(minTLSVersion, ctx.CommonTlsContext.TlsParams)

	authn_model.ApplyToCommonTLSContext(ctx.CommonTlsContext, node, []string{}, /*subjectAltNames*/
		trustDomainAliases, ctx.RequireClientCertificate.Value)
	return ctx
}

func configureTLSVersion(minTLSVer meshconfig.MeshConfig_TLSConfig_TLSProtocol, tlsParam *tls.TlsParameters) {
	// The min TLS version is at least TLS 1.2.
	tlsParam.TlsMinimumProtocolVersion = tls.TlsParameters_TLSv1_2
	// If user configures the min TLS version to be TLS 1.3, the min TLS version will be TLS 1.3.
	switch minTLSVer {
	case meshconfig.MeshConfig_TLSConfig_TLSV1_3:
		tlsParam.TlsMinimumProtocolVersion = tls.TlsParameters_TLSv1_3
	}
	// The max TLS version is configured as TLS 1.3
	tlsParam.TlsMaximumProtocolVersion = tls.TlsParameters_TLSv1_3
}
