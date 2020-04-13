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

package v2

import (
	"net"
	"strconv"
	"strings"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_api_v2_route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	v2 "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoy_config_listener_v2 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v2"

	"istio.io/istio/pkg/config/host"

	"istio.io/istio/pilot/pkg/networking/util"
)

// Support generation of 'ApiListener' LDS responses, used for native support of gRPC.
// The same response can also be used by other apps using XDS directly.

// GRPC proposal:
// https://github.com/grpc/proposal/blob/master/A27-xds-global-load-balancing.md
//
// Note that this implementation is tested against gRPC, but it is generic - any other framework can
// use this XDS mode to get load balancing info from Istio, including MC/VM/etc.

// DNS can populate the name to cluster VIP mapping using this response.

// The corresponding RDS response is also generated - currently gRPC has special differences
// and can't understand normal Istio RDS - in particular expects "" instead of "/" as
// default prefix, and is expects just the route for one host.

func (s *DiscoveryServer) handleLDSApiType(con *XdsConnection, req *xdsapi.DiscoveryRequest) error {
	if s.handleAck(con, req) {
		return nil
	}

	push := s.globalPushContext()
	resp := &xdsapi.DiscoveryResponse{
		TypeUrl:     ListenerType,
		VersionInfo: version,
		Nonce:       nonce(push.Version),
	}
	var err error
	for _, name := range req.ResourceNames {
		portI := 0
		hostn := name
		port := ""
		if strings.Contains(name, ":") {
			hostn, port, err = net.SplitHostPort(name)
			if err != nil {
				return err
			}
			if port != "" {
				portI, err = strconv.Atoi(port)
				if err != nil {
					return err
				}
			}
		}
		ll := &xdsapi.Listener{
			Name: name,
		}
		// TODO: optimize, use a map in pre-computation of sidecar scope
		for _, el := range con.node.SidecarScope.EgressListeners {
			for _, s := range el.Services() {
				if s.Hostname.Matches(host.Name(hostn)) {
					ll.Address = &envoycore.Address{
						Address: &envoycore.Address_SocketAddress{
							SocketAddress: &envoycore.SocketAddress{
								Address: s.Address,
								PortSpecifier: &envoycore.SocketAddress_PortValue{
									PortValue: uint32(portI),
								},
							},
						},
					}
					ll.ApiListener = &envoy_config_listener_v2.ApiListener{
						ApiListener: util.MessageToAny(&v2.HttpConnectionManager{
							RouteSpecifier: &v2.HttpConnectionManager_Rds{
								Rds: &v2.Rds{
									RouteConfigName: name,
								},
							},
						}),
					}
				}
			}
		}

		lr := util.MessageToAny(ll)
		resp.Resources = append(resp.Resources, lr)
	}

	err = con.send(resp)
	if err != nil {
		adsLog.Warnf("LDS: Send failure %s: %v", con.ConID, err)
		recordSendError(ldsSendErrPushes, err)
		return err
	}

	return nil
}

// Handle a gRPC CDS request, used with the 'ApiListener' style of requests.
// The main difference is that the request includes Resources.
func (s *DiscoveryServer) handleAPICDS(con *XdsConnection, req *xdsapi.DiscoveryRequest) bool {
	if s.handleAck(con, req) {
		return true
	}

	push := s.globalPushContext()
	resp := &xdsapi.DiscoveryResponse{
		TypeUrl:     ClusterType,
		VersionInfo: version,
		Nonce:       nonce(push.Version),
	}

	// gRPC doesn't currently support any of the APIs - returning just the expected EDS result.
	// Since the code is relatively strict - we'll add info as needed.
	for _, n := range req.ResourceNames {
		hn, portn, err := net.SplitHostPort(n)
		if err != nil {
			adsLog.Warna("Failed to parse ", n, " ", err)
			continue
		}
		rc := &xdsapi.Cluster{
			Name:                 n,
			ClusterDiscoveryType: &xdsapi.Cluster_Type{Type: xdsapi.Cluster_EDS},
			EdsClusterConfig: &xdsapi.Cluster_EdsClusterConfig{
				ServiceName: "outbound|" + portn + "||" + hn,
				EdsConfig: &envoycore.ConfigSource{
					ConfigSourceSpecifier: &envoycore.ConfigSource_Ads{
						Ads: &envoycore.AggregatedConfigSource{},
					},
				},
			},
		}
		rr := util.MessageToAny(rc)
		resp.Resources = append(resp.Resources, rr)
	}

	err := con.send(resp)
	if err != nil {
		adsLog.Warnf("LDS: Send failure %s: %v", con.ConID, err)
		recordSendError(ldsSendErrPushes, err)
		return true
	}
	return true
}

// handleSplitRDS supports per-VIP routes, as used by GRPC.
// This mode is indicated by using names containing full host:port instead of just port.
// Returns true of the request is of this type.
func (s *DiscoveryServer) handleSplitRDS(con *XdsConnection, req *xdsapi.DiscoveryRequest) bool {
	for _, n := range req.ResourceNames {
		if !strings.Contains(n, ":") {
			return false // normal Istio RDS, on port
		}
	}

	if s.handleAck(con, req) {
		return true
	}
	push := s.globalPushContext()

	// Currently this mode is only used by GRPC, to extract Cluster for the default
	// route.
	// Current GRPC is also expecting the default route to be prefix=="", while we generate "/"
	// in normal response.
	// TODO: add support for full route, make sure GRPC is fixed to support both
	resp := &xdsapi.DiscoveryResponse{
		TypeUrl:     RouteType,
		VersionInfo: version,
		Nonce:       nonce(push.Version),
	}
	for _, n := range req.ResourceNames {
		hn, portn, err := net.SplitHostPort(n)
		if err != nil {
			adsLog.Warna("Failed to parse ", n, " ", err)
			continue
		}
		port, err := strconv.Atoi(portn)
		if err != nil {
			adsLog.Warna("Failed to parse port ", n, " ", err)
			continue
		}
		el := con.node.SidecarScope.GetEgressListenerForRDS(port, "")
		// TODO: use VirtualServices instead !
		// Currently gRPC doesn't support matching the path.
		svc := el.Services()
		for _, s := range svc {
			if s.Hostname.Matches(host.Name(hn)) {
				// Only generate the required route for grpc. Will need to generate more
				// as GRPC adds more features.
				rc := &xdsapi.RouteConfiguration{
					Name: n,
					VirtualHosts: []*envoy_api_v2_route.VirtualHost{
						&envoy_api_v2_route.VirtualHost{
							Name:    hn,
							Domains: []string{hn, n},

							Routes: []*envoy_api_v2_route.Route{
								&envoy_api_v2_route.Route{
									Match: &envoy_api_v2_route.RouteMatch{
										PathSpecifier: &envoy_api_v2_route.RouteMatch_Prefix{Prefix: ""},
									},
									Action: &envoy_api_v2_route.Route_Route{
										Route: &envoy_api_v2_route.RouteAction{
											ClusterSpecifier: &envoy_api_v2_route.RouteAction_Cluster{
												Cluster: n,
											},
										},
									},
								},
							},
						},
					},
				}
				rr := util.MessageToAny(rc)
				resp.Resources = append(resp.Resources, rr)
			}
		}
	}
	err := con.send(resp)
	if err != nil {
		adsLog.Warnf("LDS: Send failure %s: %v", con.ConID, err)
		recordSendError(ldsSendErrPushes, err)
		return true
	}
	return true
}
