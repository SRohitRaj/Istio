// Copyright 2019 Istio Authors. All Rights Reserved.
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
	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	tcp_proxy "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/tcp_proxy/v2"
	xdsutil "github.com/envoyproxy/go-control-plane/pkg/util"
	google_protobuf "github.com/gogo/protobuf/types"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core/v1alpha3/envoyfilter"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/istio/pkg/proto"
	"istio.io/pkg/log"
)

var (
	// Precompute these filters as an optimization
	blackholeAnyMarshalling    = newBlackholeFilter(true)
	blackholeStructMarshalling = newBlackholeFilter(false)
)

// A stateful listener builder
type ListenerBuilder struct {
	node                   *model.Proxy
	gatewayListeners       []*xdsapi.Listener
	inboundListeners       []*xdsapi.Listener
	outboundListeners      []*xdsapi.Listener
	managementListeners    []*xdsapi.Listener
	virtualListener        *xdsapi.Listener
	virtualInboundListener *xdsapi.Listener
}

func NewListenerBuilder(node *model.Proxy) *ListenerBuilder {
	builder := &ListenerBuilder{
		node: node,
	}
	return builder
}

func (builder *ListenerBuilder) buildSidecarInboundListeners(
	configgen *ConfigGeneratorImpl,
	env *model.Environment, node *model.Proxy, push *model.PushContext,
	proxyInstances []*model.ServiceInstance) *ListenerBuilder {
	builder.inboundListeners = configgen.buildSidecarInboundListeners(env, node, push, proxyInstances)
	return builder
}

func (builder *ListenerBuilder) buildSidecarOutboundListeners(configgen *ConfigGeneratorImpl,
	env *model.Environment, node *model.Proxy, push *model.PushContext,
	proxyInstances []*model.ServiceInstance) *ListenerBuilder {
	builder.outboundListeners = configgen.buildSidecarOutboundListeners(env, node, push, proxyInstances)
	return builder
}

func (builder *ListenerBuilder) buildManagementListeners(_ *ConfigGeneratorImpl,
	env *model.Environment, node *model.Proxy, _ *model.PushContext,
	_ []*model.ServiceInstance) *ListenerBuilder {

	noneMode := node.GetInterceptionMode() == model.InterceptionNone

	// Do not generate any management port listeners if the user has specified a SidecarScope object
	// with ingress listeners. Specifying the ingress listener implies that the user wants
	// to only have those specific listeners and nothing else, in the inbound path.
	if node.SidecarScope.HasCustomIngressListeners || noneMode {
		return builder
	}
	// Let ServiceDiscovery decide which IP and Port are used for management if
	// there are multiple IPs
	mgmtListeners := make([]*xdsapi.Listener, 0)
	for _, ip := range node.IPAddresses {
		managementPorts := env.ManagementPorts(ip)
		management := buildSidecarInboundMgmtListeners(node, env, managementPorts, ip)
		mgmtListeners = append(mgmtListeners, management...)
	}
	addresses := make(map[string]*xdsapi.Listener)
	for _, listener := range builder.inboundListeners {
		if listener != nil {
			addresses[listener.Address.String()] = listener
		}
	}
	for _, listener := range builder.outboundListeners {
		if listener != nil {
			addresses[listener.Address.String()] = listener
		}
	}

	// If management listener port and service port are same, bad things happen
	// when running in kubernetes, as the probes stop responding. So, append
	// non overlapping listeners only.
	for i := range mgmtListeners {
		m := mgmtListeners[i]
		addressString := m.Address.String()
		existingListener, ok := addresses[addressString]
		if ok {
			log.Warnf("Omitting listener for management address %s due to collision with service listener (%s)",
				m.Name, existingListener.Name)
			continue
		} else {
			// dedup management listeners as well
			addresses[addressString] = m
			builder.managementListeners = append(builder.managementListeners, m)
		}

	}
	return builder
}

func (builder *ListenerBuilder) buildVirtualOutboundListener(
	configgen *ConfigGeneratorImpl,
	env *model.Environment, node *model.Proxy, push *model.PushContext,
	proxyInstances []*model.ServiceInstance) *ListenerBuilder {

	var isTransparentProxy *google_protobuf.BoolValue
	if node.GetInterceptionMode() == model.InterceptionTproxy {
		isTransparentProxy = proto.BoolTrue
	}

	tcpProxyFilter := newTCPProxyOutboundListenerFilter(env, node)

	filterChains := []listener.FilterChain{
		{
			Filters: []listener.Filter{*tcpProxyFilter},
		},
	}

	// The virtual listener will handle all traffic that does not match any other listeners, and will
	// blackhole/passthrough depending on the outbound traffic policy. When passthrough is enabled,
	// this has the risk of triggering infinite loops when requests are sent to the pod's IP, as it will
	// send requests to itself. To block this we add an additional filter chain before that will always blackhole.
	if features.RestrictPodIPTrafficLoops.Get() {
		var cidrRanges []*core.CidrRange
		for _, ip := range node.IPAddresses {
			cidrRanges = append(cidrRanges, util.ConvertAddressToCidr(ip))
		}
		blackhole := blackholeStructMarshalling
		if util.IsXDSMarshalingToAnyEnabled(node) {
			blackhole = blackholeAnyMarshalling
		}
		filterChains = append([]listener.FilterChain{{
			FilterChainMatch: &listener.FilterChainMatch{
				PrefixRanges: cidrRanges,
			},
			Filters: []listener.Filter{blackhole},
		}}, filterChains...)
	}

	actualWildcard, _ := getActualWildcardAndLocalHost(node)

	// add an extra listener that binds to the port that is the recipient of the iptables redirect
	ipTablesListener := &xdsapi.Listener{
		Name:           VirtualOutboundListenerName,
		Address:        util.BuildAddress(actualWildcard, uint32(env.Mesh.ProxyListenPort)),
		Transparent:    isTransparentProxy,
		UseOriginalDst: proto.BoolTrue,
		FilterChains:   filterChains,
	}
	configgen.onVirtualOutboundListener(env, node, push, proxyInstances,
		ipTablesListener)
	builder.virtualListener = ipTablesListener
	return builder
}

// TProxy uses only the virtual outbound listener on 15001 for both directions
// but we still ship the no-op virtual inbound listener, so that the code flow is same across REDIRECT and TPROXY.
func (builder *ListenerBuilder) buildVirtualInboundListener(env *model.Environment, node *model.Proxy) *ListenerBuilder {
	var isTransparentProxy *google_protobuf.BoolValue
	if node.GetInterceptionMode() == model.InterceptionTproxy {
		isTransparentProxy = proto.BoolTrue
	}

	actualWildcard, _ := getActualWildcardAndLocalHost(node)
	// add an extra listener that binds to the port that is the recipient of the iptables redirect
	builder.virtualInboundListener = &xdsapi.Listener{
		Name:           VirtualInboundListenerName,
		Address:        util.BuildAddress(actualWildcard, ProxyInboundListenPort),
		Transparent:    isTransparentProxy,
		UseOriginalDst: proto.BoolTrue,
		FilterChains:   newInboundPassthroughFilterChains(env, node),
	}
	return builder
}

func (builder *ListenerBuilder) patchListeners(push *model.PushContext) {
	if builder.node.Type == model.Router {
		envoyfilter.ApplyListenerPatches(networking.EnvoyFilter_GATEWAY, builder.node, push, builder.gatewayListeners, false)
		return
	}

	patchOneListener := func(listener *xdsapi.Listener) *xdsapi.Listener {
		if listener == nil {
			return nil
		}
		tempArray := []*xdsapi.Listener{listener}
		tempArray = envoyfilter.ApplyListenerPatches(networking.EnvoyFilter_SIDECAR_OUTBOUND, builder.node, push, tempArray, true)
		// temp array will either be empty [if virtual listener was removed] or will have a modified listener
		if len(tempArray) == 0 {
			return nil
		}
		return tempArray[0]
	}
	builder.virtualListener = patchOneListener(builder.virtualListener)
	builder.virtualInboundListener = patchOneListener(builder.virtualInboundListener)
	builder.managementListeners = envoyfilter.ApplyListenerPatches(networking.EnvoyFilter_SIDECAR_INBOUND, builder.node,
		push, builder.managementListeners, true)
	builder.inboundListeners = envoyfilter.ApplyListenerPatches(networking.EnvoyFilter_SIDECAR_INBOUND, builder.node,
		push, builder.inboundListeners, false)
	builder.outboundListeners = envoyfilter.ApplyListenerPatches(networking.EnvoyFilter_SIDECAR_INBOUND, builder.node,
		push, builder.outboundListeners, false)
}

func (builder *ListenerBuilder) getListeners() []*xdsapi.Listener {
	if builder.node.Type == model.SidecarProxy {
		nInbound, nOutbound, nManagement := len(builder.inboundListeners), len(builder.outboundListeners), len(builder.managementListeners)
		nVirtual, nVirtualInbound := 0, 0
		if builder.virtualListener != nil {
			nVirtual = 1
		}
		if builder.virtualInboundListener != nil {
			nVirtualInbound = 1
		}
		nListener := nInbound + nOutbound + nManagement + nVirtual + nVirtualInbound

		listeners := make([]*xdsapi.Listener, 0, nListener)
		listeners = append(listeners, builder.inboundListeners...)
		listeners = append(listeners, builder.outboundListeners...)
		listeners = append(listeners, builder.managementListeners...)
		if builder.virtualListener != nil {
			listeners = append(listeners, builder.virtualListener)
		}
		if builder.virtualInboundListener != nil {
			listeners = append(listeners, builder.virtualInboundListener)
		}

		log.Debugf("Build %d listeners for node %s including %d inbound, %d outbound, %d management, %d virtual and %d virtual inbound listeners",
			nListener,
			builder.node.ID,
			nInbound, nOutbound, nManagement,
			nVirtual, nVirtualInbound)
		return listeners
	}

	return builder.gatewayListeners
}

// Creates a new filter that will always send traffic to the blackhole cluster
func newBlackholeFilter(enableAny bool) listener.Filter {
	tcpProxy := &tcp_proxy.TcpProxy{
		StatPrefix:       util.BlackHoleCluster,
		ClusterSpecifier: &tcp_proxy.TcpProxy_Cluster{Cluster: util.BlackHoleCluster},
	}

	filter := listener.Filter{
		Name: xdsutil.TCPProxy,
	}

	if enableAny {
		filter.ConfigType = &listener.Filter_TypedConfig{TypedConfig: util.MessageToAny(tcpProxy)}
	} else {
		filter.ConfigType = &listener.Filter_Config{Config: util.MessageToStruct(tcpProxy)}
	}
	return filter
}

func newInboundPassthroughFilterChains(env *model.Environment, node *model.Proxy) []listener.FilterChain {
	// ipv4 and ipv6
	filterChains := make([]listener.FilterChain, 0, 2)
	for _, clusterName := range []string{util.InboundPassthroughClusterIpv4, util.InboundPassthroughClusterIpv6} {

		tcpProxy := &tcp_proxy.TcpProxy{
			StatPrefix:       clusterName,
			ClusterSpecifier: &tcp_proxy.TcpProxy_Cluster{Cluster: clusterName},
		}

		matchingIP := ""
		if clusterName == util.InboundPassthroughClusterIpv4 {
			matchingIP = util.InboundPassthroughBindIpv4
		} else if clusterName == util.InboundPassthroughClusterIpv6 {
			matchingIP = util.InboundPassthroughBindIpv6
		}

		filterChainMatch := listener.FilterChainMatch{
			// Port : EMPTY to match all ports
			PrefixRanges: []*core.CidrRange{
				util.ConvertAddressToCidr(matchingIP),
			},
		}
		setAccessLog(env, node, tcpProxy)
		filter := listener.Filter{
			Name: xdsutil.TCPProxy,
		}

		if util.IsXDSMarshalingToAnyEnabled(node) {
			filter.ConfigType = &listener.Filter_TypedConfig{TypedConfig: util.MessageToAny(tcpProxy)}
		} else {
			filter.ConfigType = &listener.Filter_Config{Config: util.MessageToStruct(tcpProxy)}
		}
		filterChain := listener.FilterChain{
			FilterChainMatch: &filterChainMatch,
			Filters: []listener.Filter{
				filter,
			},
		}
		filterChains = append(filterChains, filterChain)
	}

	return filterChains
}

func newTCPProxyOutboundListenerFilter(env *model.Environment, node *model.Proxy) *listener.Filter {
	tcpProxy := &tcp_proxy.TcpProxy{
		StatPrefix:       util.BlackHoleCluster,
		ClusterSpecifier: &tcp_proxy.TcpProxy_Cluster{Cluster: util.BlackHoleCluster},
	}
	if isAllowAnyOutbound(node) {
		// We need a passthrough filter to fill in the filter stack for orig_dst listener
		tcpProxy = &tcp_proxy.TcpProxy{
			StatPrefix:       util.PassthroughCluster,
			ClusterSpecifier: &tcp_proxy.TcpProxy_Cluster{Cluster: util.PassthroughCluster},
		}
		setAccessLog(env, node, tcpProxy)
	}

	filter := listener.Filter{
		Name: xdsutil.TCPProxy,
	}

	if util.IsXDSMarshalingToAnyEnabled(node) {
		filter.ConfigType = &listener.Filter_TypedConfig{TypedConfig: util.MessageToAny(tcpProxy)}
	} else {
		filter.ConfigType = &listener.Filter_Config{Config: util.MessageToStruct(tcpProxy)}
	}
	return &filter
}

func isAllowAnyOutbound(node *model.Proxy) bool {
	return node.SidecarScope.OutboundTrafficPolicy != nil && node.SidecarScope.OutboundTrafficPolicy.Mode == networking.OutboundTrafficPolicy_ALLOW_ANY
}
