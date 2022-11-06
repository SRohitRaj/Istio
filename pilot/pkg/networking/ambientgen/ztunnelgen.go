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

/*
Listener ztunnel_outbound:
  Chain: <srcpodname>_<srcpodip>_to_<portname>_<hostname>_<vip>
      -> CDS: <identity>_to_<vip>_<portname>_<hostname>_outbound_internal
          transport: internal
					-> EDS:
							address: podIP
							tunnel: outbound_tunnel_lis_<identity>
  Chain: <srcpodname>_<srcpodip>_to_client_waypoint_proxy_<waypoint>
      tunneling_config: ztunnel-to-waypoint
      -> CDS: to_client_waypoint_proxy_<source identity> (EDS)
           address: WAYPOINT_IP
           transport: TLS
  Chain: <srcpodname>_<srcpodip>_to_server_waypoint_proxy_<waypoint>
      tunneling_config: ztunnel-to-waypoint
      -> CDS: <source identity>_to_server_waypoint_proxy_<server identity> (EDS)
           address: WAYPOINT_IP
           transport: TLS
  Chain: passthrough
  Chain: blackhole

Internal listener: outbound_tunnel_lis_<identity>
      tunneling_config: host.com:443
      -> CDS: outbound_tunnel_clus_<identity> (ORIG_DST)
          transport: TLS

Listener ztunnel_inbound:
  Chain: inbound_<podip>
      transport: terminate TLS
      match: CONNECT
      -> CDS: virtual_inbound (ORIG_DST)
  Chain: blackhole
*/

package ambientgen

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	accesslog "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	fileaccesslog "github.com/envoyproxy/go-control-plane/envoy/extensions/access_loggers/file/v3"
	routerfilter "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	originaldst "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/original_dst/v3"
	originalsrc "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/listener/original_src/v3"
	httpconn "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tcp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	http "github.com/envoyproxy/go-control-plane/envoy/extensions/upstreams/http/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	any "google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	"istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/ambient"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core/v1alpha3"
	"istio.io/istio/pilot/pkg/networking/core/v1alpha3/match"
	"istio.io/istio/pilot/pkg/networking/plugin/authn"
	"istio.io/istio/pilot/pkg/networking/plugin/authz"
	"istio.io/istio/pilot/pkg/networking/util"
	security "istio.io/istio/pilot/pkg/security/model"
	"istio.io/istio/pilot/pkg/serviceregistry/provider"
	"istio.io/istio/pilot/pkg/util/protoconv"
	v3 "istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pkg/config/host"
	"istio.io/istio/pkg/config/labels"
	"istio.io/istio/pkg/util/sets"
	istiolog "istio.io/pkg/log"
)

var log = istiolog.RegisterScope("ambientgen", "xDS Generator for ambient clients", 0)

var _ model.XdsResourceGenerator = &ZTunnelConfigGenerator{}

type ZTunnelConfigGenerator struct {
	EndpointIndex *model.EndpointIndex
	Workloads     ambient.Cache
}

func (g *ZTunnelConfigGenerator) Generate(
	proxy *model.Proxy,
	w *model.WatchedResource,
	req *model.PushRequest,
) (model.Resources, model.XdsLogDetails, error) {
	push := req.Push
	switch w.TypeUrl {
	case v3.ListenerType:
		return g.BuildListeners(proxy, push, w.ResourceNames), model.DefaultXdsLogDetails, nil
	case v3.ClusterType:
		return g.BuildClusters(proxy, push, w.ResourceNames), model.DefaultXdsLogDetails, nil
	case v3.EndpointType:
		return g.BuildEndpoints(proxy, push, w.ResourceNames), model.DefaultXdsLogDetails, nil
	}

	return nil, model.DefaultXdsLogDetails, nil
}

const (
	ZTunnelOutboundCapturePort         uint32 = 15001
	ZTunnelInbound2CapturePort         uint32 = 15006
	ZTunnelInboundNodeLocalCapturePort uint32 = 15088
	ZTunnelInboundCapturePort          uint32 = 15008

	// TODO: this needs to match the mark in the iptables rules.
	// And also not clash with any other mark on the host level.
	// either figure out a way to not hardcode it, or a way to not use it.
	// i think the best solution is to have this mark configurable and run the
	// iptables rules from the code, so we are sure the mark matches.
	OriginalSrcMark = 0x4d2
	OutboundMark    = 0x401
	InboundMark     = 0x402
)

// these exist on syscall package, but only on linux.
// copy these here so this file can build on any platform
const (
	SolSocket = 0x1
	SoMark    = 0x24
)

func (g *ZTunnelConfigGenerator) BuildListeners(proxy *model.Proxy, push *model.PushContext, names []string) (out model.Resources) {
	out = append(out,
		g.buildPodOutboundCaptureListener(proxy, push),
		g.buildInboundCaptureListener(proxy, push),
		g.buildInboundPlaintextCaptureListener(proxy, push),
	)
	for sa := range push.AmbientIndex.Workloads.ByIdentity {
		out = append(out, outboundTunnelListener(outboundTunnelListenerName(sa), sa))
	}

	return out
}

func (g *ZTunnelConfigGenerator) BuildClusters(proxy *model.Proxy, push *model.PushContext, names []string) model.Resources {
	var out model.Resources
	// TODO node local SAs only?
	services := proxy.SidecarScope.Services()
	workloads := push.AmbientIndex.Workloads
	for sa := range workloads.ByIdentity {
		for _, svc := range services {
			for _, port := range svc.Ports {
				c := g.serviceOutboundCluster(proxy, push, sa, svc, port.Name)
				if c == nil {
					continue
				}
				out = append(out, &discovery.Resource{Name: c.Name, Resource: protoconv.MessageToAny(c)})
			}
		}
	}

	for sa := range workloads.NodeLocalBySA(proxy.Metadata.NodeName) {
		c := outboundTunnelCluster(proxy, push, sa, sa)
		out = append(out, &discovery.Resource{Name: c.Name, Resource: protoconv.MessageToAny(c)})
	}
	for sa := range workloads.NodeLocalBySA(proxy.Metadata.NodeName) {
		c := outboundPodTunnelCluster(proxy, push, sa, sa)
		out = append(out, &discovery.Resource{Name: c.Name, Resource: protoconv.MessageToAny(c)})
	}
	for sa := range workloads.NodeLocalBySA(proxy.Metadata.NodeName) {
		c := outboundPodLocalTunnelCluster(proxy, push, sa, sa)
		out = append(out, &discovery.Resource{Name: c.Name, Resource: protoconv.MessageToAny(c)})
	}

	out = append(out, buildWaypointClusters(proxy, push)...)
	out = append(out,
		g.buildVirtualInboundCluster(),
		g.buildVirtualInboundClusterHBONE(),
		passthroughCluster(push),
		tcpPassthroughCluster(push),
		blackholeCluster(push))
	return out
}

func blackholeCluster(push *model.PushContext) *discovery.Resource {
	c := &cluster.Cluster{
		Name:                 util.BlackHoleCluster,
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_STATIC},
		ConnectTimeout:       push.Mesh.ConnectTimeout,
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
	}
	return &discovery.Resource{
		Name:     c.Name,
		Resource: protoconv.MessageToAny(c),
	}
}

func passthroughCluster(push *model.PushContext) *discovery.Resource {
	c := &cluster.Cluster{
		Name:                 util.PassthroughCluster,
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
		ConnectTimeout:       push.Mesh.ConnectTimeout,
		LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
		// TODO protocol options are copy-paste from v1alpha3 package
		TypedExtensionProtocolOptions: map[string]*any.Any{
			v3.HttpProtocolOptionsType: protoconv.MessageToAny(&http.HttpProtocolOptions{
				UpstreamProtocolOptions: &http.HttpProtocolOptions_UseDownstreamProtocolConfig{
					UseDownstreamProtocolConfig: &http.HttpProtocolOptions_UseDownstreamHttpConfig{
						HttpProtocolOptions: &core.Http1ProtocolOptions{},
						Http2ProtocolOptions: &core.Http2ProtocolOptions{
							// Envoy default value of 100 is too low for data path.
							MaxConcurrentStreams: &wrappers.UInt32Value{
								Value: 1073741824,
							},
						},
					},
				},
			}),
		},
	}
	return &discovery.Resource{Name: c.Name, Resource: protoconv.MessageToAny(c)}
}

func tcpPassthroughCluster(push *model.PushContext) *discovery.Resource {
	c := &cluster.Cluster{
		Name:                 util.PassthroughCluster + "-tcp",
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
		ConnectTimeout:       push.Mesh.ConnectTimeout,
		LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
	}
	return &discovery.Resource{Name: c.Name, Resource: protoconv.MessageToAny(c)}
}

// buildPodOutboundCaptureListener creates a single listener with a FilterChain for each combination
// of ServiceAccount from pods on the node and Service VIP in the cluster.
func (g *ZTunnelConfigGenerator) buildPodOutboundCaptureListener(proxy *model.Proxy, push *model.PushContext) *discovery.Resource {
	l := &listener.Listener{
		Name:           "ztunnel_outbound",
		UseOriginalDst: wrappers.Bool(true),
		Transparent:    wrappers.Bool(true),
		AccessLog:      accessLogString("outbound capture listener"),
		SocketOptions: []*core.SocketOption{{
			Description: "Set socket mark to packets coming back from outbound listener",
			Level:       SolSocket,
			Name:        SoMark,
			Value: &core.SocketOption_IntValue{
				IntValue: OutboundMark,
			},
			State: core.SocketOption_STATE_PREBIND,
		}},
		ListenerFilters: []*listener.ListenerFilter{
			{
				Name: wellknown.OriginalDestination,
				ConfigType: &listener.ListenerFilter_TypedConfig{
					TypedConfig: protoconv.MessageToAny(&originaldst.OriginalDst{}),
				},
			},
		},
		Address: &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address: "0.0.0.0",
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: ZTunnelOutboundCapturePort,
				},
			},
		}},
	}
	if push.Mesh.GetOutboundTrafficPolicy().GetMode() == v1alpha1.MeshConfig_OutboundTrafficPolicy_ALLOW_ANY {
		l.DefaultFilterChain = passthroughFilterChain()
	}
	// nolint: gocritic
	// if features.SidecarlessCapture == model.VariantIptables {
	l.ListenerFilters = append(l.ListenerFilters, &listener.ListenerFilter{
		Name: wellknown.OriginalSource,
		ConfigType: &listener.ListenerFilter_TypedConfig{
			TypedConfig: protoconv.MessageToAny(&originalsrc.OriginalSrc{
				Mark: OriginalSrcMark,
			}),
		},
	})
	//}

	l.ListenerFilters = append(l.ListenerFilters, &listener.ListenerFilter{
		Name: WorkloadMetadataListenerFilterName,
		ConfigType: &listener.ListenerFilter_ConfigDiscovery{
			ConfigDiscovery: &core.ExtensionConfigSource{
				ConfigSource: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{Ads: &core.AggregatedConfigSource{}},
					InitialFetchTimeout:   durationpb.New(30 * time.Second),
				},
				TypeUrls: []string{WorkloadMetadataResourcesTypeURL},
			},
		},
	})

	// match logic:
	// dest port == 15001 -> blackhole
	// source unknown -> passthrough
	// source known, has waypoint -> client waypoint
	// source known, no waypoint, dest is a VIP -> resolve VIP, use passthrough metadata from EDS for tunnel headers
	// source known, no waypoint, dest NOT a VIP -> use original src/dest for tunnel headers (headless)

	sourceMatch := match.NewSourceIP()
	sourceMatch.OnNoMatch = match.ToChain(util.PassthroughFilterChain)

	destPortMatch := match.NewDestinationPort()
	destPortMatch.OnNoMatch = match.ToMatcher(sourceMatch.Matcher)
	destPortMatch.Map[strconv.Itoa(int(l.GetAddress().GetSocketAddress().GetPortValue()))] = match.ToChain(util.BlackHoleCluster)

	services := proxy.SidecarScope.Services()
	seen := sets.New[string]()
	for _, sourceWl := range push.AmbientIndex.Workloads.NodeLocal(proxy.Metadata.NodeName) {
		sourceAndDestMatch := match.NewDestinationIP()
		// TODO: handle host network better, which has a shared IP
		sourceMatch.Map[sourceWl.PodIP] = match.ToMatcher(sourceAndDestMatch.Matcher)

		clientWaypoints := push.AmbientIndex.Waypoints.ByIdentity[sourceWl.Identity()] // TODO need to use this instead of ServiceAccountName
		clientWaypointChain := buildWaypointChain(sourceWl, clientWaypoints, "client")

		for _, svc := range services {
			// No client waypoint proxy, we build a chain per destination VIP
			vip := svc.GetAddressForProxy(proxy)

			portMatch := match.NewDestinationPort()
			sourceAndDestMatch.Map[vip] = match.ToMatcher(portMatch.Matcher)
			for _, port := range svc.Ports {
				var chain *listener.FilterChain
				serverWaypointChain := g.maybeBuildServerWaypointChain(push, sourceWl, svc)
				if serverWaypointChain != nil {
					// Has server waypoint proxy, send traffic there
					chain = serverWaypointChain
				} else if clientWaypointChain != nil {
					// Has client waypoint proxy, send traffic there
					chain = clientWaypointChain
				} else {
					// No waypoint proxy
					name := outboundServiceClusterName(sourceWl.Identity(), port.Name, svc.Hostname.String())
					chain = &listener.FilterChain{
						Name: name,
						Filters: []*listener.Filter{{
							Name: wellknown.TCPProxy,
							ConfigType: &listener.Filter_TypedConfig{TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
								AccessLog:        accessLogString("capture outbound (no waypoint proxy)"),
								StatPrefix:       name,
								ClusterSpecifier: &tcp.TcpProxy_Cluster{Cluster: name},
							},
							)},
						}},
					}
				}

				if !seen.InsertContains(chain.Name) {
					l.FilterChains = append(l.FilterChains, chain)
				}
				portMatch.Map[fmt.Sprint(port.Port)] = match.ToChain(chain.Name)
			}
		}
		// Add chain for each pod IP
		wls := push.AmbientIndex.Workloads.All()
		wls = append(wls, push.AmbientIndex.None.All()...)
		for _, wl := range wls {
			var chain *listener.FilterChain
			// Need to decide if there is a server waypoint proxy. This is somewhat problematic because a Service may span waypoint proxy and non-waypoint proxy.
			// If any workload behind the service has a waypoint proxy, we will use the waypoint proxy. In 99% of cases this is homogenous.
			serverWaypointChain := buildWaypointChain(sourceWl, push.AmbientIndex.Waypoints.ByIdentity[wl.Identity()], "server")
			if serverWaypointChain != nil {
				// Has server waypoint proxy, send traffic there
				chain = serverWaypointChain
			} else if clientWaypointChain != nil {
				// Has client waypoint proxy, send traffic there
				chain = clientWaypointChain
			} else {
				// No waypoint proxy
				// Naively, we could simply create a FC with tunnel_config here and point to an original_dst cluster.
				// This won't work for a few reasons:
				// * We need to override the port. `x-envoy-original-dst-host` cannot be used since it is an
				//   upstream header; the cluster looks for downstream headers
				// * We could add config to orig_dst cluster to override the port. This would almost work, but
				//   then we run into issues with the original_src filter. Currently, this filter is on the listener filter
				//   but it only applies for direct connections. When we go through another internal listener, the effect is lost.
				//   Ultimately that means for tunneling, we do not use the original_src but for direct calls we do. This means
				//   that we will need to go through an internal listener to "break" the original_src effect.
				// TODO: this is broken
				// If we use outboundTunnelClusterName, we get orig_dst, but x-envoy-original-dst-host is an upstream header
				// while the cluster looks for downstream headers.
				// if we make a dedicate cluster, we cannot pass the original port anymore since the context is lost.
				// We cannot create a cluster per port since it can be any port.
				// TODO2: this is still broken even with custom orig_dst. the listener sets the orig_src mark
				// If we

				tunnel := &tcp.TcpProxy_TunnelingConfig{
					Hostname: "%DOWNSTREAM_LOCAL_ADDRESS%",
					HeadersToAdd: []*core.HeaderValueOption{
						// This is for server ztunnel - not really needed for waypoint proxy
						{Header: &core.HeaderValue{Key: "x-envoy-original-dst-host", Value: "%DOWNSTREAM_LOCAL_ADDRESS%"}},
					},
				}
				// Case 1: tunnel cross node
				cluster := outboundPodTunnelClusterName(sourceWl.Identity())
				// Case 2: same node tunnel (iptables)
				if node := wl.NodeName; node != "" && node == proxy.Metadata.NodeName {
					cluster = outboundPodLocalTunnelClusterName(sourceWl.Identity())
				}
				// Case 3: direct
				if wl.Labels[ambient.LabelType] != ambient.TypeWorkload {
					cluster = util.PassthroughCluster + "-tcp"
					tunnel = nil
				}

				name := "fc-" + cluster
				chain = &listener.FilterChain{
					Name: name,
					Filters: []*listener.Filter{
						{
							Name: wellknown.TCPProxy,
							ConfigType: &listener.Filter_TypedConfig{
								TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
									AccessLog:        accessLogString("capture outbound pod (no waypoint proxy)"),
									StatPrefix:       name,
									ClusterSpecifier: &tcp.TcpProxy_Cluster{Cluster: cluster},
									TunnelingConfig:  tunnel,
								}),
							},
						},
					},
				}
			}

			if !seen.InsertContains(chain.Name) {
				l.FilterChains = append(l.FilterChains, chain)
			}
			sourceAndDestMatch.Map[wl.PodIP] = match.ToChain(chain.Name)
		}
	}

	l.FilterChainMatcher = destPortMatch.BuildMatcher()
	l.FilterChains = append(l.FilterChains, passthroughFilterChain(), blackholeFilterChain("outbound"))
	return &discovery.Resource{
		Name:     l.Name,
		Resource: protoconv.MessageToAny(l),
	}
}

// Need to decide if there is a server waypoint proxy. This is somewhat problematic because a Service may span waypoint and non-waypoint enabled proxies.
// If any workload behind the service has a waypoint, we will use the waypoint. In 99% of cases this is homogenous.
func (g *ZTunnelConfigGenerator) maybeBuildServerWaypointChain(push *model.PushContext,
	sourceWl ambient.Workload, svc *model.Service,
) *listener.FilterChain {
	var serviceWorkloads []ambient.Workload
	if svc.Attributes.ServiceRegistry == provider.External &&
		svc.Attributes.LabelSelectors == nil {
		// there are a small number of workloads specified directly by the service, check those
		shards, _ := g.EndpointIndex.ShardsForService(svc.Hostname.String(), svc.Attributes.Namespace)
		serviceWorkloads = workloadsForShards(push.AmbientIndex, shards)
	} else {
		// Find waypoints based on label selectors for any workload
		// TODO optimize this so we don't do full service selection for every service on every gen
		for _, wl := range push.AmbientIndex.Workloads.All() {
			if wl.Namespace != svc.Attributes.Namespace {
				continue
			}
			if len(svc.Attributes.LabelSelectors) == 0 || !labels.Instance(svc.Attributes.LabelSelectors).SubsetOf(wl.Labels) {
				continue
			}
			serviceWorkloads = append(serviceWorkloads, wl)
		}
	}

	// if any workload in the service has a waypoint proxy, all traffic to the service must go through it
	// TODO what happens if workloads specify multiple SAs that have waypoint proxies?
	for _, wl := range serviceWorkloads {
		if waypoints := push.AmbientIndex.Waypoints.ByIdentity[wl.Identity()]; len(waypoints) > 0 {
			return buildWaypointChain(sourceWl, waypoints, "server")
		}
	}

	return nil
}

func workloadsForShards(workloads ambient.Indexes, shards *model.EndpointShards) (out []ambient.Workload) {
	if shards == nil {
		return
	}
	shards.RLock()
	defer shards.RUnlock()

	for _, endpoints := range shards.Shards {
		for _, istioEndpoint := range endpoints {
			if w, ok := workloads.Workloads.ByIP[istioEndpoint.Address]; ok {
				out = append(out, w)
			}
		}
	}
	return out
}

func blackholeFilterChain(t string) *listener.FilterChain {
	return &listener.FilterChain{
		Name: "blackhole " + t,
		Filters: []*listener.Filter{{
			Name: wellknown.TCPProxy,
			ConfigType: &listener.Filter_TypedConfig{TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
				AccessLog:        accessLogString("blackhole " + t),
				StatPrefix:       util.BlackHoleCluster,
				ClusterSpecifier: &tcp.TcpProxy_Cluster{Cluster: "blackhole " + t},
			})},
		}},
	}
}

func buildWaypointChain(workload ambient.Workload, waypoints []ambient.Workload, t string) *listener.FilterChain {
	if len(waypoints) == 0 {
		return nil
	}

	// waypoint is just for identity (same across multiple waypoints)
	waypoint := waypoints[0]
	// For client waypoint proxy, we know the waypoint proxy and client are always the same identity which simplifies things; we can share a cluster for all
	cluster := waypointClusterName(waypoint.Identity())
	if t == "server" {
		// For server, we need to create the product of source identity x waypoint proxy
		cluster = serverWaypointClusterName(waypoint.Identity(), workload.Identity())
	}
	return &listener.FilterChain{
		Name: cluster,
		Filters: []*listener.Filter{{
			Name: wellknown.TCPProxy,
			ConfigType: &listener.Filter_TypedConfig{TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
				AccessLog:        accessLogString(fmt.Sprintf("capture outbound (to %v waypoint proxy)", t)),
				StatPrefix:       cluster,
				ClusterSpecifier: &tcp.TcpProxy_Cluster{Cluster: cluster},
				TunnelingConfig: &tcp.TcpProxy_TunnelingConfig{
					Hostname: "%DOWNSTREAM_LOCAL_ADDRESS%", // (unused, per extended connect)
					HeadersToAdd: []*core.HeaderValueOption{
						// This is for server ztunnel - not really needed for waypoint proxy
						{Header: &core.HeaderValue{Key: "x-envoy-original-dst-host", Value: "%DOWNSTREAM_LOCAL_ADDRESS%"}},

						// This is for metadata propagation
						// TODO: should we just set the baggage directly, as we have access to the Pod here (instead of using the filter)?
						{Header: &core.HeaderValue{Key: "baggage", Value: "%DYNAMIC_METADATA([\"envoy.filters.listener.workload_metadata\", \"baggage\"])%"}},
					},
				},
			},
			)},
		}},
	}
}

func passthroughFilterChain() *listener.FilterChain {
	return &listener.FilterChain{
		Name: util.PassthroughFilterChain,
		/// TODO no match – add one to make it so we only passthrough if strict mTLS to the destination is allowed
		Filters: []*listener.Filter{{
			Name: wellknown.TCPProxy,
			ConfigType: &listener.Filter_TypedConfig{TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
				AccessLog:        accessLogString("passthrough"),
				StatPrefix:       util.PassthroughCluster,
				ClusterSpecifier: &tcp.TcpProxy_Cluster{Cluster: util.PassthroughCluster},
			})},
		}},
	}
}

func (g *ZTunnelConfigGenerator) serviceOutboundCluster(
	proxy *model.Proxy, push *model.PushContext, sa string, svc *model.Service, port string,
) *cluster.Cluster {
	discoveryType := convertResolution(proxy.Type, svc)
	c := &cluster.Cluster{
		Name:                 outboundServiceClusterName(sa, port, svc.Hostname.String()),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: discoveryType},

		TransportSocketMatches: v1alpha3.InternalUpstreamSocketMatch,
	}
	switch discoveryType {
	case cluster.Cluster_STRICT_DNS, cluster.Cluster_LOGICAL_DNS:
		if proxy.SupportsIPv4() {
			c.DnsLookupFamily = cluster.Cluster_V4_ONLY
		} else {
			c.DnsLookupFamily = cluster.Cluster_V6_ONLY
		}
		dnsRate := push.Mesh.DnsRefreshRate
		c.DnsRefreshRate = dnsRate
		c.RespectDnsTtl = true
		fallthrough
	case cluster.Cluster_STATIC:
		localityLbEndpoints := g.upstreamLbEndpointsFromShards(proxy, push, sa, svc, port)
		if len(localityLbEndpoints) == 0 {
			log.Warnf("%s cluster without endpoints %s found while pushing CDS", discoveryType.String(), c.Name)
			return nil
		}
		c.LoadAssignment = &endpoint.ClusterLoadAssignment{
			ClusterName: c.Name,
			Endpoints:   localityLbEndpoints,
		}
	case cluster.Cluster_EDS:
		c.EdsClusterConfig = &cluster.Cluster_EdsClusterConfig{
			EdsConfig: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_Ads{
					Ads: &core.AggregatedConfigSource{},
				},
				InitialFetchTimeout: durationpb.New(0),
				ResourceApiVersion:  core.ApiVersion_V3,
			},
		}
	case cluster.Cluster_ORIGINAL_DST:
		c.LbPolicy = cluster.Cluster_CLUSTER_PROVIDED
	}
	return c
}

func outboundServiceClusterName(sa, port string, hostname string) string {
	return fmt.Sprintf("%s_to_%s_%s_outbound_internal", sa, port, hostname)
}

func parseServiceOutboundClusterName(clusterName string) (sa, port string, hostname string, ok bool) {
	p := strings.Split(clusterName, "_")
	if !strings.HasSuffix(clusterName, "_outbound_internal") || len(p) < 3 {
		return "", "", "", false
	}
	return p[0], p[2], p[3], true
}

func waypointClusterName(waypoint string) string {
	return fmt.Sprintf("_to_client_waypoint_proxy_%s", waypoint)
}

func serverWaypointClusterName(waypoint, workload string) string {
	return fmt.Sprintf("%s_to_server_waypoint_proxy_%s", workload, waypoint)
}

// parseWaypointClusterName parses cluster names, in the format {src}_to_{t}_waypoint_proxy_{dst} where src/dst are identities
func parseWaypointClusterName(name string) (src, dst, t string, ok bool) {
	p := strings.Split(name, "_")
	if len(p) != 6 || p[1] != "to" || p[3] != "waypoint" {
		return "", "", "", false
	}
	return p[0], p[5], p[2], true
}

func buildWaypointClusters(proxy *model.Proxy, push *model.PushContext) model.Resources {
	var clusters []*cluster.Cluster
	// Client waypoints
	for sa, waypoints := range push.AmbientIndex.Waypoints.ByIdentity {
		saWorkloads := push.AmbientIndex.Workloads.NodeLocalBySA(proxy.Metadata.NodeName)[sa]
		if len(saWorkloads) == 0 || len(waypoints) == 0 {
			// no waypoints or no workloads that use this client waypoint on the node
			continue
		}
		clusters = append(clusters, &cluster.Cluster{
			Name:                          waypointClusterName(sa),
			ClusterDiscoveryType:          &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
			LbPolicy:                      cluster.Cluster_ROUND_ROBIN,
			ConnectTimeout:                durationpb.New(2 * time.Second),
			TypedExtensionProtocolOptions: h2connectUpgrade(),
			TransportSocket: &core.TransportSocket{
				Name: "envoy.transport_sockets.tls",
				ConfigType: &core.TransportSocket_TypedConfig{TypedConfig: protoconv.MessageToAny(&tls.UpstreamTlsContext{
					CommonTlsContext: buildCommonTLSContext(proxy, sa, push, false),
				})},
			},
			EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
				EdsConfig: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{
						Ads: &core.AggregatedConfigSource{},
					},
					InitialFetchTimeout: durationpb.New(0),
					ResourceApiVersion:  core.ApiVersion_V3,
				},
			},
		})
	}
	for waypointSA, waypoints := range push.AmbientIndex.Waypoints.ByIdentity {
		for workloadSA, workloads := range push.AmbientIndex.Workloads.NodeLocalBySA(proxy.Metadata.NodeName) {
			if len(workloads) == 0 || len(waypoints) == 0 {
				// no waypoint proxies or no workloads that use this identity on the node
				continue
			}
			clusters = append(clusters, &cluster.Cluster{
				Name:                          serverWaypointClusterName(waypointSA, workloadSA),
				ClusterDiscoveryType:          &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
				LbPolicy:                      cluster.Cluster_ROUND_ROBIN,
				ConnectTimeout:                durationpb.New(2 * time.Second),
				TypedExtensionProtocolOptions: h2connectUpgrade(),
				TransportSocket: &core.TransportSocket{
					Name: "envoy.transport_sockets.tls",
					ConfigType: &core.TransportSocket_TypedConfig{TypedConfig: protoconv.MessageToAny(&tls.UpstreamTlsContext{
						CommonTlsContext: buildCommonTLSContext(proxy, workloadSA, push, false),
					})},
				},
				EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
					EdsConfig: &core.ConfigSource{
						ConfigSourceSpecifier: &core.ConfigSource_Ads{
							Ads: &core.AggregatedConfigSource{},
						},
						InitialFetchTimeout: durationpb.New(0),
						ResourceApiVersion:  core.ApiVersion_V3,
					},
				},
			})
		}
	}
	var out model.Resources
	for _, c := range clusters {
		out = append(out, &discovery.Resource{
			Name:     c.Name,
			Resource: protoconv.MessageToAny(c),
		})
	}
	return out
}

func (g *ZTunnelConfigGenerator) BuildEndpoints(proxy *model.Proxy, push *model.PushContext, names []string) model.Resources {
	out := model.Resources{}
	// ztunnel outbound to upstream
	for _, clusterName := range names {
		// sa here is already our "envoy friendly" one
		sa, port, hostname, ok := parseServiceOutboundClusterName(clusterName)
		if !ok {
			continue
		}
		svc := push.ServiceForHostname(proxy, host.Name(hostname))
		out = append(out, &discovery.Resource{
			Name: clusterName,
			Resource: protoconv.MessageToAny(&endpoint.ClusterLoadAssignment{
				ClusterName: clusterName,
				Endpoints:   g.upstreamLbEndpointsFromShards(proxy, push, sa, svc, port),
			}),
		})
	}
	// ztunnel to waypoint
	for _, clusterName := range names {
		_, dst, t, ok := parseWaypointClusterName(clusterName)
		if !ok {
			continue
		}
		out = append(out, &discovery.Resource{
			Name: clusterName,
			Resource: protoconv.MessageToAny(&endpoint.ClusterLoadAssignment{
				ClusterName: clusterName,
				Endpoints:   buildWaypointLbEndpoints(dst, t, push),
			}),
		})
	}
	return out
}

func (g *ZTunnelConfigGenerator) upstreamLbEndpointsFromShards(
	proxy *model.Proxy, push *model.PushContext, sa string, svc *model.Service, portName string,
) []*endpoint.LocalityLbEndpoints {
	if svc == nil {
		return nil
	}
	port, ok := svc.Ports.Get(portName)
	if !ok {
		return nil
	}

	var istioEndpoints []*model.IstioEndpoint
	switch svc.Resolution {
	case model.DNSLB, model.DNSRoundRobinLB:
		instances := push.ServiceInstancesByPort(svc, port.Port, nil)
		for _, instance := range instances {
			istioEndpoints = append(istioEndpoints, instance.Endpoint)
		}
	case model.ClientSideLB:
		shards, ok := g.EndpointIndex.ShardsForService(svc.Hostname.String(), svc.Attributes.Namespace)
		if !ok || shards == nil {
			log.Warnf("no endpoint shards for %s/%s", svc.Attributes.Namespace, svc.Attributes.Name)
			return nil
		}
		shards.RLock()
		for _, shard := range shards.Shards {
			istioEndpoints = append(istioEndpoints, shard...)
		}
		shards.RUnlock()
	}

	shards, ok := g.EndpointIndex.ShardsForService(svc.Hostname.String(), svc.Attributes.Namespace)
	if !ok || shards == nil {
		log.Warnf("no endpoint shards for %s/%s", svc.Attributes.Namespace, svc.Attributes.Name)
		return nil
	}
	eps := &endpoint.LocalityLbEndpoints{
		LbEndpoints: nil,
	}
	for _, istioEndpoint := range istioEndpoints {
		if portName != istioEndpoint.ServicePortName {
			continue
		}
		lbe := &endpoint.LbEndpoint{
			HostIdentifier: &endpoint.LbEndpoint_Endpoint{Endpoint: &endpoint.Endpoint{
				Address: &core.Address{
					Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Address:       istioEndpoint.Address,
							PortSpecifier: &core.SocketAddress_PortValue{PortValue: istioEndpoint.EndpointPort},
						},
					},
				},
			}},
			LoadBalancingWeight: wrappers.UInt32(1),
		}

		capturePort := ZTunnelInboundCapturePort
		// TODO passthrough for node-local upstreams without Waypoints
		if node := istioEndpoint.NodeName; node != "" && node == proxy.Metadata.NodeName {
			capturePort = ZTunnelInboundNodeLocalCapturePort
		}
		supportsTunnel := false
		if al := istioEndpoint.Labels[ambient.LabelType]; al == ambient.TypeWaypoint || al == ambient.TypeWorkload {
			// Waypointss and in-meshed workloads can do a tunnel
			supportsTunnel = true
		}
		if istioEndpoint.SupportsTunnel(model.TunnelHTTP) && istioEndpoint.EndpointPort == ZTunnelInboundCapturePort {
			// Even if it is in the mesh, if it supports tunnel directly then we should pass through the traffic if its already tunneled
			// TODO this assumes it gets captured and server ztunnel inits the tunnel
			supportsTunnel = false
		}
		if istioEndpoint.SupportsTunnel(model.TunnelHTTP) {
			// if the pod natively supports tunnel, node local doesn't change the port since we're not relying on redirection here
			capturePort = ZTunnelInboundCapturePort // TODO should this be if tunnel: h2 && !captured?
			supportsTunnel = true
		}

		if supportsTunnel {
			// TODO re-use some eds code; stable eds ordering, support multi-cluster cluster local rules, and multi-network stuff
			tunnelLis := outboundTunnelListenerName(sa)
			lbe = util.BuildInternalLbEndpoint(tunnelLis, util.BuildTunnelMetadata(
				istioEndpoint.Address,
				int(istioEndpoint.EndpointPort),
				int(capturePort)))
			lbe.Metadata.FilterMetadata[util.EnvoyTransportSocketMetadataKey] = &structpb.Struct{
				Fields: map[string]*structpb.Value{
					model.TunnelLabelShortName: {Kind: &structpb.Value_StringValue{StringValue: model.TunnelHTTP}},
				},
			}
		}
		eps.LbEndpoints = append(eps.LbEndpoints, lbe)
	}
	return []*endpoint.LocalityLbEndpoints{eps}
}

func buildWaypointLbEndpoints(waypointIdentity, t string, push *model.PushContext) []*endpoint.LocalityLbEndpoints {
	port := ZTunnelOutboundCapturePort
	if t == "server" {
		port = ZTunnelInbound2CapturePort
	}
	waypoints := push.AmbientIndex.Waypoints.ByIdentity[waypointIdentity]

	lbEndpoints := &endpoint.LocalityLbEndpoints{
		LbEndpoints: []*endpoint.LbEndpoint{},
	}
	for _, waypoint := range waypoints {
		lbEndpoints.LbEndpoints = append(lbEndpoints.LbEndpoints, &endpoint.LbEndpoint{
			HostIdentifier: &endpoint.LbEndpoint_Endpoint{Endpoint: &endpoint.Endpoint{
				Address: &core.Address{
					Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Address:       waypoint.PodIP,
							PortSpecifier: &core.SocketAddress_PortValue{PortValue: port},
						},
					},
				},
			}},
		})
	}
	return []*endpoint.LocalityLbEndpoints{lbEndpoints}
}

func outboundTunnelListenerName(sa string) string {
	return "outbound_tunnel_lis_" + sa
}

// outboundTunnelListener is built for each ServiceAccount from pods on the node.
// This listener adds the original destination headers from the dynamic EDS metadata pass through.
// We build the listener per-service account so that it can point to the corresponding cluster that presents the correct cert.
func outboundTunnelListener(name string, sa string) *discovery.Resource {
	l := &listener.Listener{
		Name:              name,
		UseOriginalDst:    wrappers.Bool(false),
		ListenerSpecifier: &listener.Listener_InternalListener{InternalListener: &listener.Listener_InternalListenerConfig{}},
		ListenerFilters:   []*listener.ListenerFilter{util.InternalListenerSetAddressFilter()},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.TCPProxy,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
						StatPrefix:       name,
						AccessLog:        accessLogString("outbound tunnel"),
						ClusterSpecifier: &tcp.TcpProxy_Cluster{Cluster: outboundTunnelClusterName(sa)},
						TunnelingConfig: &tcp.TcpProxy_TunnelingConfig{
							Hostname: "%DYNAMIC_METADATA(tunnel:destination)%",
							HeadersToAdd: []*core.HeaderValueOption{
								{Header: &core.HeaderValue{Key: "x-envoy-original-dst-host", Value: "%DYNAMIC_METADATA([\"tunnel\", \"destination\"])%"}},
							},
						},
					}),
				},
			}},
		}},
	}
	return &discovery.Resource{
		Name:     name,
		Resource: protoconv.MessageToAny(l),
	}
}

func buildCommonTLSContext(proxy *model.Proxy, identityOverride string, push *model.PushContext, inbound bool) *tls.CommonTlsContext {
	ctx := &tls.CommonTlsContext{}
	// TODO san match
	security.ApplyToCommonTLSContext(ctx, proxy, nil, authn.TrustDomainsForValidation(push.Mesh), inbound)

	if identityOverride != "" {
		ctx.TlsCertificateSdsSecretConfigs = []*tls.SdsSecretConfig{
			security.ConstructSdsSecretConfig(identityOverride),
		}
	}
	ctx.AlpnProtocols = []string{"h2"}

	ctx.TlsParams = &tls.TlsParameters{
		// Ensure TLS 1.3 is used everywhere
		TlsMaximumProtocolVersion: tls.TlsParameters_TLSv1_3,
		TlsMinimumProtocolVersion: tls.TlsParameters_TLSv1_3,
	}

	return ctx
}

// outboundTunnelCluster is per-workload SA
func outboundTunnelCluster(proxy *model.Proxy, push *model.PushContext, sa string, identityOverride string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:                 outboundTunnelClusterName(sa),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
		LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
		ConnectTimeout:       durationpb.New(2 * time.Second),
		CleanupInterval:      durationpb.New(60 * time.Second),
		LbConfig: &cluster.Cluster_OriginalDstLbConfig_{
			OriginalDstLbConfig: &cluster.Cluster_OriginalDstLbConfig{},
		},
		TypedExtensionProtocolOptions: h2connectUpgrade(),
		TransportSocket: &core.TransportSocket{
			Name: "envoy.transport_sockets.tls",
			ConfigType: &core.TransportSocket_TypedConfig{TypedConfig: protoconv.MessageToAny(&tls.UpstreamTlsContext{
				CommonTlsContext: buildCommonTLSContext(proxy, identityOverride, push, false),
			})},
		},
	}
}

// outboundTunnelCluster is per-workload SA
func outboundPodTunnelCluster(proxy *model.Proxy, push *model.PushContext, sa string, identityOverride string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:                 outboundPodTunnelClusterName(sa),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
		LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
		ConnectTimeout:       durationpb.New(2 * time.Second),
		CleanupInterval:      durationpb.New(60 * time.Second),
		LbConfig: &cluster.Cluster_OriginalDstLbConfig_{
			OriginalDstLbConfig: &cluster.Cluster_OriginalDstLbConfig{
				UpstreamPortOverride: wrappers.UInt32(ZTunnelInboundCapturePort),
			},
		},
		TypedExtensionProtocolOptions: h2connectUpgrade(),
		TransportSocket: &core.TransportSocket{
			Name: "envoy.transport_sockets.tls",
			ConfigType: &core.TransportSocket_TypedConfig{TypedConfig: protoconv.MessageToAny(&tls.UpstreamTlsContext{
				CommonTlsContext: buildCommonTLSContext(proxy, identityOverride, push, false),
			})},
		},
	}
}

// outboundTunnelCluster is per-workload SA
func outboundPodLocalTunnelCluster(proxy *model.Proxy, push *model.PushContext, sa string, identityOverride string) *cluster.Cluster {
	return &cluster.Cluster{
		Name:                 outboundPodLocalTunnelClusterName(sa),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
		LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
		ConnectTimeout:       durationpb.New(2 * time.Second),
		CleanupInterval:      durationpb.New(60 * time.Second),
		LbConfig: &cluster.Cluster_OriginalDstLbConfig_{
			OriginalDstLbConfig: &cluster.Cluster_OriginalDstLbConfig{
				UseHttpHeader:        true,
				UpstreamPortOverride: wrappers.UInt32(ZTunnelInboundNodeLocalCapturePort),
			},
		},
		TypedExtensionProtocolOptions: h2connectUpgrade(),
		TransportSocket: &core.TransportSocket{
			Name: "envoy.transport_sockets.tls",
			ConfigType: &core.TransportSocket_TypedConfig{TypedConfig: protoconv.MessageToAny(&tls.UpstreamTlsContext{
				CommonTlsContext: buildCommonTLSContext(proxy, identityOverride, push, false),
			})},
		},
	}
}

func outboundTunnelClusterName(sa string) string {
	return "outbound_tunnel_clus_" + sa
}

func outboundPodTunnelClusterName(sa string) string {
	return "outbound_pod_tunnel_clus_" + sa
}

func outboundPodLocalTunnelClusterName(sa string) string {
	return "outbound_pod_local_tunnel_clus_" + sa
}

// buildInboundCaptureListener creates a single listener with a FilterChain for each pod on the node.
func (g *ZTunnelConfigGenerator) buildInboundCaptureListener(proxy *model.Proxy, push *model.PushContext) *discovery.Resource {
	// TODO L7 stuff (deny at l4 for l7 auth if there is a waypoint proxy for the dest workload)

	l := &listener.Listener{
		Name:           "ztunnel_inbound",
		UseOriginalDst: wrappers.Bool(true),
		ListenerFilters: []*listener.ListenerFilter{
			{
				Name: wellknown.OriginalDestination,
				ConfigType: &listener.ListenerFilter_TypedConfig{
					TypedConfig: protoconv.MessageToAny(&originaldst.OriginalDst{}),
				},
			},
			{
				Name: wellknown.OriginalSource,
				ConfigType: &listener.ListenerFilter_TypedConfig{
					TypedConfig: protoconv.MessageToAny(&originalsrc.OriginalSrc{
						Mark: OriginalSrcMark,
					}),
				},
			},
		},
		Transparent: wrappers.Bool(true),
		AccessLog:   accessLogString("capture inbound listener"),
		SocketOptions: []*core.SocketOption{{
			Description: "Set socket mark to packets coming back from inbound listener",
			Level:       SolSocket,
			Name:        SoMark,
			Value: &core.SocketOption_IntValue{
				IntValue: InboundMark,
			},
			State: core.SocketOption_STATE_PREBIND,
		}},
		Address: &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				// TODO because of the port 15088 workaround, we need to use a redirect rule,
				// which means we can't bind to localhost. once we remove that workaround,
				// this can be changed back to 127.0.0.1
				Address: "0.0.0.0",
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: ZTunnelInboundCapturePort,
				},
			},
		}},
	}

	for _, workload := range push.AmbientIndex.Workloads.NodeLocal(proxy.Metadata.NodeName) {
		// Skip workloads in the host network
		if workload.HostNetwork {
			continue
		}

		if workload.Labels[model.TunnelLabel] != model.TunnelHTTP {
			dummy := &model.Proxy{
				ConfigNamespace: workload.Namespace,
				Labels:          workload.Labels,
			}
			var allowedIdentities string
			_, hasWaypoint := push.AmbientIndex.Waypoints.ByIdentity[workload.Identity()]
			if hasWaypoint {
				allowedIdentities = strings.TrimPrefix(workload.Identity(), "spiffe://")
			}
			authzBuilder := authz.NewBuilderSkipIdentity(authz.Local, push, dummy, allowedIdentities)
			tcp := authzBuilder.BuildTCP()

			var filters []*listener.Filter
			filters = append(filters, tcp...)
			filters = append(filters, &listener.Filter{
				Name: "envoy.filters.network.http_connection_manager",
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: protoconv.MessageToAny(&httpconn.HttpConnectionManager{
						AccessLog:  accessLogString("inbound hcm"),
						CodecType:  0,
						StatPrefix: "inbound_hcm",
						RouteSpecifier: &httpconn.HttpConnectionManager_RouteConfig{
							RouteConfig: &route.RouteConfiguration{
								Name: "local_route",
								VirtualHosts: []*route.VirtualHost{{
									Name:    "local_service",
									Domains: []string{"*"},
									Routes: []*route.Route{{
										Match: &route.RouteMatch{PathSpecifier: &route.RouteMatch_ConnectMatcher_{
											ConnectMatcher: &route.RouteMatch_ConnectMatcher{},
										}},
										Action: &route.Route_Route{
											Route: &route.RouteAction{
												UpgradeConfigs: []*route.RouteAction_UpgradeConfig{{
													UpgradeType:   "CONNECT",
													ConnectConfig: &route.RouteAction_UpgradeConfig_ConnectConfig{},
												}},
												ClusterSpecifier: &route.RouteAction_Cluster{
													Cluster: "virtual_inbound",
												},
											},
										},
									}},
								}},
							},
						},
						// TODO rewrite destination port to original_dest port
						HttpFilters: []*httpconn.HttpFilter{{
							Name:       "envoy.filters.http.router",
							ConfigType: &httpconn.HttpFilter_TypedConfig{TypedConfig: protoconv.MessageToAny(&routerfilter.Router{})},
						}},
						Http2ProtocolOptions: &core.Http2ProtocolOptions{
							AllowConnect: true,
						},
						UpgradeConfigs: []*httpconn.HttpConnectionManager_UpgradeConfig{{
							UpgradeType: "CONNECT",
						}},
					}),
				},
			})
			l.FilterChains = append(l.FilterChains, &listener.FilterChain{
				Name:             "inbound_" + workload.PodIP,
				FilterChainMatch: &listener.FilterChainMatch{PrefixRanges: matchIP(workload.PodIP)},
				TransportSocket: &core.TransportSocket{
					Name: "envoy.transport_sockets.tls",
					ConfigType: &core.TransportSocket_TypedConfig{TypedConfig: protoconv.MessageToAny(&tls.DownstreamTlsContext{
						CommonTlsContext: buildCommonTLSContext(proxy, workload.Identity(), push, true),
					})},
				},
				Filters: filters,
			})
		} else {
			// Pod is already handling HBONE, and this is an HBONE request. Pass it through directly.
			l.FilterChains = append(l.FilterChains, &listener.FilterChain{
				Name:             "inbound_" + workload.PodIP,
				FilterChainMatch: &listener.FilterChainMatch{PrefixRanges: matchIP(workload.PodIP)},
				Filters: []*listener.Filter{{
					Name: wellknown.TCPProxy,
					ConfigType: &listener.Filter_TypedConfig{
						TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
							StatPrefix: "virtual_inbound_hbone",
							AccessLog:  accessLogString("inbound passthrough"),
							ClusterSpecifier: &tcp.TcpProxy_Cluster{
								Cluster: "virtual_inbound_hbone",
							},
						}),
					},
				}},
			})
		}
	}
	// TODO cases where we passthrough
	l.FilterChains = append(l.FilterChains, blackholeFilterChain("inbound"))

	return &discovery.Resource{
		Name:     l.Name,
		Resource: protoconv.MessageToAny(l),
	}
}

// buildInboundCaptureListener creates a single listener with a FilterChain for each pod on the node.
func (g *ZTunnelConfigGenerator) buildInboundPlaintextCaptureListener(proxy *model.Proxy, push *model.PushContext) *discovery.Resource {
	// TODO L7 stuff (deny at l4 for l7 auth if there is a waypoint proxy for the dest workload)
	l := &listener.Listener{
		Name:           "ztunnel_inbound_plaintext",
		UseOriginalDst: wrappers.Bool(true),
		ListenerFilters: []*listener.ListenerFilter{
			{
				Name: wellknown.OriginalDestination,
				ConfigType: &listener.ListenerFilter_TypedConfig{
					TypedConfig: protoconv.MessageToAny(&originaldst.OriginalDst{}),
				},
			},
			{
				Name: wellknown.OriginalSource,
				ConfigType: &listener.ListenerFilter_TypedConfig{
					TypedConfig: protoconv.MessageToAny(&originalsrc.OriginalSrc{
						Mark: OriginalSrcMark,
					}),
				},
			},
		},
		AccessLog: accessLogString("capture inbound listener plaintext"),
		SocketOptions: []*core.SocketOption{{
			Description: "Set socket mark to packets coming back from inbound listener",
			Level:       SolSocket,
			Name:        SoMark,
			Value: &core.SocketOption_IntValue{
				IntValue: InboundMark,
			},
			State: core.SocketOption_STATE_PREBIND,
		}},
		Address: &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				// TODO because of the port 15088 workaround, we need to use a redirect rule,
				// which means we can't bind to localhost. once we remove that workaround,
				// this can be changed back to 127.0.0.1
				Address: "0.0.0.0",
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: ZTunnelInbound2CapturePort,
				},
			},
		}},
		Transparent: wrappers.Bool(true),
	}

	for _, workload := range push.AmbientIndex.Workloads.NodeLocal(proxy.Metadata.NodeName) {
		// Skip workloads in the host network
		if workload.HostNetwork {
			continue
		}

		dummy := &model.Proxy{
			ConfigNamespace: workload.Namespace,
			Labels:          workload.Labels,
		}
		authzBuilder := authz.NewBuilder(authz.Local, push, dummy)

		var filters []*listener.Filter
		filters = append(filters, authzBuilder.BuildTCP()...)
		filters = append(filters, &listener.Filter{
			Name: wellknown.TCPProxy,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: protoconv.MessageToAny(&tcp.TcpProxy{
					StatPrefix:       "virtual_inbound_plaintext",
					ClusterSpecifier: &tcp.TcpProxy_Cluster{Cluster: "virtual_inbound"},
				}),
			},
		})
		l.FilterChains = append(l.FilterChains, &listener.FilterChain{
			Name:             "inbound_" + workload.PodIP,
			FilterChainMatch: &listener.FilterChainMatch{PrefixRanges: matchIP(workload.PodIP)},
			Filters:          filters,
		})
	}
	// TODO cases where we passthrough
	l.FilterChains = append(l.FilterChains, blackholeFilterChain("inbound plaintext"))

	return &discovery.Resource{
		Name:     l.Name,
		Resource: protoconv.MessageToAny(l),
	}
}

func (g *ZTunnelConfigGenerator) buildVirtualInboundCluster() *discovery.Resource {
	c := &cluster.Cluster{
		Name:                 "virtual_inbound",
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
		LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
		LbConfig: &cluster.Cluster_OriginalDstLbConfig_{
			OriginalDstLbConfig: &cluster.Cluster_OriginalDstLbConfig{
				UseHttpHeader: true,
			},
		},
	}
	return &discovery.Resource{
		Name:     c.Name,
		Resource: protoconv.MessageToAny(c),
	}
}

// Like virtual_inbound, but always sets port to 15008. This is a huge hack to fix HBONE passhrough
// to node-local endpoints. These would send to 15088, which then gets looped back to us then
// forwarded. But we need the forwarding to go to 15008 the second iteration.
func (g *ZTunnelConfigGenerator) buildVirtualInboundClusterHBONE() *discovery.Resource {
	c := &cluster.Cluster{
		Name:                 "virtual_inbound_hbone",
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_ORIGINAL_DST},
		LbPolicy:             cluster.Cluster_CLUSTER_PROVIDED,
		LbConfig: &cluster.Cluster_OriginalDstLbConfig_{
			OriginalDstLbConfig: &cluster.Cluster_OriginalDstLbConfig{
				UseHttpHeader:        true,
				UpstreamPortOverride: wrappers.UInt32(ZTunnelInboundCapturePort),
			},
		},
	}
	return &discovery.Resource{
		Name:     c.Name,
		Resource: protoconv.MessageToAny(c),
	}
}

func matchIP(addr string) []*core.CidrRange {
	return []*core.CidrRange{{
		AddressPrefix: addr,
		PrefixLen:     wrappers.UInt32(32),
	}}
}

const EnvoyTextLogFormat = "[%START_TIME%] \"%REQ(:METHOD)% %REQ(X-ENVOY-ORIGINAL-PATH?:PATH)% " +
	"%PROTOCOL%\" %RESPONSE_CODE% %RESPONSE_FLAGS% " +
	"%RESPONSE_CODE_DETAILS% %CONNECTION_TERMINATION_DETAILS% " +
	"\"%UPSTREAM_TRANSPORT_FAILURE_REASON%\" %BYTES_RECEIVED% %BYTES_SENT% " +
	"%DURATION% %RESP(X-ENVOY-UPSTREAM-SERVICE-TIME)% \"%REQ(X-FORWARDED-FOR)%\" " +
	"\"%REQ(USER-AGENT)%\" \"%REQ(X-REQUEST-ID)%\" \"%REQ(:AUTHORITY)%\" \"%UPSTREAM_HOST%\" " +
	"%UPSTREAM_CLUSTER% %UPSTREAM_LOCAL_ADDRESS% %DOWNSTREAM_LOCAL_ADDRESS% " +
	"%DOWNSTREAM_REMOTE_ADDRESS% %REQUESTED_SERVER_NAME% %ROUTE_NAME% "

func accessLogString(prefix string) []*accesslog.AccessLog {
	inlineString := EnvoyTextLogFormat + prefix + "\n"
	return []*accesslog.AccessLog{{
		Name: "envoy.access_loggers.file",
		ConfigType: &accesslog.AccessLog_TypedConfig{TypedConfig: protoconv.MessageToAny(&fileaccesslog.FileAccessLog{
			Path: "/dev/stdout",
			AccessLogFormat: &fileaccesslog.FileAccessLog_LogFormat{LogFormat: &core.SubstitutionFormatString{
				Format: &core.SubstitutionFormatString_TextFormatSource{TextFormatSource: &core.DataSource{Specifier: &core.DataSource_InlineString{
					InlineString: inlineString,
				}}},
			}},
		})},
	}}
}

func h2connectUpgrade() map[string]*any.Any {
	return map[string]*any.Any{
		v3.HttpProtocolOptionsType: protoconv.MessageToAny(&http.HttpProtocolOptions{
			UpstreamProtocolOptions: &http.HttpProtocolOptions_ExplicitHttpConfig_{ExplicitHttpConfig: &http.HttpProtocolOptions_ExplicitHttpConfig{
				ProtocolConfig: &http.HttpProtocolOptions_ExplicitHttpConfig_Http2ProtocolOptions{
					Http2ProtocolOptions: &core.Http2ProtocolOptions{
						AllowConnect: true,
					},
				},
			}},
		}),
	}
}

func ipPortAddress(ip string, port uint32) *core.Address {
	return &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address: ip,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: port,
			},
		},
	}}
}

// TODO re-use from v1alpha3/cluster.go

func convertResolution(proxyType model.NodeType, service *model.Service) cluster.Cluster_DiscoveryType {
	switch service.Resolution {
	case model.ClientSideLB:
		return cluster.Cluster_EDS
	case model.DNSLB:
		return cluster.Cluster_STRICT_DNS
	case model.DNSRoundRobinLB:
		return cluster.Cluster_LOGICAL_DNS
	case model.Passthrough:
		// Gateways cannot use passthrough clusters. So fallback to EDS
		if proxyType == model.SidecarProxy {
			if service.Attributes.ServiceRegistry == provider.Kubernetes && features.EnableEDSForHeadless {
				return cluster.Cluster_EDS
			}

			return cluster.Cluster_ORIGINAL_DST
		}
		return cluster.Cluster_EDS
	default:
		return cluster.Cluster_EDS
	}
}
