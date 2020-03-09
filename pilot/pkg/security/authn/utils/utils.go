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

package utils

import (
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	ldsv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	xdsutil "github.com/envoyproxy/go-control-plane/pkg/wellknown"
	structpb "github.com/golang/protobuf/ptypes/struct"

	"istio.io/pkg/log"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking"
	"istio.io/istio/pilot/pkg/networking/util"
	protovalue "istio.io/istio/pkg/proto"
)

// BuildInboundFilterChain returns the filter chain(s) corresponding to the mTLS mode.
func BuildInboundFilterChain(mTLSMode model.MutualTLSMode, sdsUdsPath string,
	node *model.Proxy, protocol networking.ListenerProtocol) []networking.FilterChain {
	if mTLSMode == model.MTLSDisable || mTLSMode == model.MTLSUnknown {
		return nil
	}

	meta := node.Metadata
	var alpnIstioMatch *ldsv2.FilterChainMatch
	var tls *auth.DownstreamTlsContext
	alpnIstioMatch = &ldsv2.FilterChainMatch{
		ApplicationProtocols: util.ALPNMtlsTCPWithMxc,
	}
	tls = &auth.DownstreamTlsContext{
		CommonTlsContext: &auth.CommonTlsContext{
			// For TCP with mTLS, we advertise "istio-peer-exchange" from client and
			// expect the same from server. This is so that secure metadata exchange
			// transfer can take place between sidecars for TCP with mTLS.
			AlpnProtocols: util.ALPNMtlsTCPWithMxc,
		},
		RequireClientCertificate: protovalue.BoolTrue,
	}
	if protocol == networking.ListenerProtocolHTTP {
		alpnIstioMatch.ApplicationProtocols = util.ALPNMtlsHTTP
		tls.CommonTlsContext.AlpnProtocols = util.ALPNMtlsHTTP
	}

	util.ApplyToCommonTLSContext(tls.CommonTlsContext, meta, sdsUdsPath, []string{} /*subjectAltNames*/)

	if mTLSMode == model.MTLSStrict {
		log.Debug("Allow only istio mutual TLS traffic")
		return []networking.FilterChain{
			{
				TLSContext: tls,
			}}
	}
	if mTLSMode == model.MTLSPermissive {
		log.Debug("Allow both, ALPN istio and legacy traffic")
		return []networking.FilterChain{
			{
				FilterChainMatch: alpnIstioMatch,
				TLSContext:       tls,
				ListenerFilters: []*ldsv2.ListenerFilter{
					{
						Name:       xdsutil.TlsInspector,
						ConfigType: &ldsv2.ListenerFilter_Config{Config: &structpb.Struct{}},
					},
				},
			},
			{
				FilterChainMatch: &ldsv2.FilterChainMatch{},
			},
		}
	}
	return nil
}
