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

package v1alpha3

import (
	"fmt"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	http_conn "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
)

func buildGatewayListeners(env model.Environment, node model.Proxy) ([]*xdsapi.Listener, error) {
	config := env.IstioConfigStore

	// collect workload labels
	workloadInstances, err := env.GetProxyServiceInstances(node)
	if err != nil {
		log.Errora("Failed to get gateway instances for router ", node.ID, err)
		return nil, err
	}

	var workloadLabels model.LabelsCollection
	for _, w := range workloadInstances {
		workloadLabels = append(workloadLabels, w.Labels)
	}

	gateways := config.Gateways(workloadLabels)

	if len(gateways) == 0 {
		log.Debuga("no gateways for router", node.ID)
		return []*xdsapi.Listener{}, nil
	}

	// TODO: merging makes code simpler but we lose gateway names that are needed to determine
	// the virtual services pinned to each gateway
	//gateway := &networking.Gateway{}
	//for _, spec := range gateways {
	//	err := model.MergeGateways(gateway, spec.Spec.(*networking.Gateway))
	//	if err != nil {
	//		log.Errorf("Failed to merge gateway %s for router %s", spec.Name, node.ID)
	//		return nil, fmt.Errorf("merge gateways: %s", err)
	//	}
	//}

	// HACK for the above case
	if len(gateways) > 1 {
		log.Warn("Currently, Istio cannot bind multiple gateways to the same workload")
		return []*xdsapi.Listener{}, nil
	}

	name := gateways[0].Name
	gateway := gateways[0].Spec.(*networking.Gateway)

	listeners := make([]*xdsapi.Listener, 0, len(gateway.Servers))
	listenerPortMap := make(map[uint32]bool)
	for _, server := range gateway.Servers {

		// TODO: this does not handle the case where there are two servers on same port
		if listenerPortMap[server.Port.Number] {
			log.Warnf("Multiple servers on same port is not supported yet, port %d", server.Port.Number)
			continue
		}
		switch model.Protocol(server.Port.Protocol) {
		case model.ProtocolHTTP, model.ProtocolHTTP2, model.ProtocolGRPC, model.ProtocolHTTPS:
			opts := buildListenerOpts{
				env:            env,
				proxy:          node,
				proxyInstances: nil, // only required to support deprecated mixerclient behavior
				ip:             WildcardAddress,
				port:           int(server.Port.Number),
				protocol:       model.ProtocolHTTP,
				sniHosts:       server.Hosts,
				tlsContext:     buildGatewayListenerTLSContext(server),
				bindToPort:     true,
				httpOpts: &httpListenerOpts{
					routeConfig:      buildGatewayInboundHTTPRouteConfig(env, name, server),
					rds:              "",
					useRemoteAddress: true,
					direction:        http_conn.EGRESS, // viewed as from gateway to internal
				},
			}

			l := buildListener(opts)
			listeners = append(listeners, l)
		case model.ProtocolTCP, model.ProtocolMongo:
			// TODO
			// Look at virtual service specs, and identity destinations,
			// call buildOutboundNetworkFilters.. and then construct TCPListener
			//buildTCPListener(buildOutboundNetworkFilters(clusterName, addresses, servicePort),
			//	listenAddress, uint32(servicePort.Port), servicePort.Protocol)
		}
	}

	return listeners, nil
}

func buildGatewayListenerTLSContext(server *networking.Server) *auth.DownstreamTlsContext {
	if server.Tls == nil {
		return nil
	}

	return &auth.DownstreamTlsContext{
		CommonTlsContext: &auth.CommonTlsContext{
			TlsCertificates: []*auth.TlsCertificate{
				{
					CertificateChain: &core.DataSource{
						Specifier: &core.DataSource_Filename{
							Filename: server.Tls.ServerCertificate,
						},
					},
					PrivateKey: &core.DataSource{
						Specifier: &core.DataSource_Filename{
							Filename: server.Tls.PrivateKey,
						},
					},
				},
			},
			ValidationContext: &auth.CertificateValidationContext{
				TrustedCa: &core.DataSource{
					Specifier: &core.DataSource_Filename{
						Filename: server.Tls.CaCertificates,
					},
				},
				VerifySubjectAltName: server.Tls.SubjectAltNames,
			},
			AlpnProtocols: ListenersALPNProtocols,
		},
		RequireSni: &types.BoolValue{
			Value: true, // is that OKAY?
		},
	}
}

func buildGatewayInboundHTTPRouteConfig(env model.Environment, gatewayName string,
	server *networking.Server) *xdsapi.RouteConfiguration {
	// TODO WE DO NOT SUPPORT two gateways on same workload binding to same virtual service
	virtualServices := env.VirtualServices([]string{gatewayName})

	virtualHosts := make([]route.VirtualHost, 0)
	// TODO: Need to trim output based on source label/gateway match
	for _, v := range virtualServices {
		guardedRoute := TranslateRoutes(v, nil)
		var routes []route.Route
		for _, g := range guardedRoute {
			routes = append(routes, g.Route)
		}
		domains := v.Spec.(*networking.VirtualService).Hosts

		virtualHosts = append(virtualHosts, route.VirtualHost{
			Name:    fmt.Sprintf("%s:%d", v.Name, server.Port.Number),
			Domains: domains,
			Routes:  routes,
		})
	}

	return &xdsapi.RouteConfiguration{
		Name:         fmt.Sprintf("%d", server.Port.Number),
		VirtualHosts: virtualHosts,
	}
}
