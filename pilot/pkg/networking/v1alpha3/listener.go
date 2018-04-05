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

package v1alpha3

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	accesslog "github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2"
	http_conn "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	tcp_proxy "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/tcp_proxy/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"

	google_protobuf "github.com/gogo/protobuf/types"

	"os"

	authn "istio.io/api/authentication/v1alpha1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
)

const (
	fileAccessLog = "envoy.file_access_log"

	envoyHTTPConnectionManager = "envoy.http_connection_manager"

	// HTTPStatPrefix indicates envoy stat prefix for http listeners
	HTTPStatPrefix = "http"

	// RDSName is the name of route-discovery-service (RDS) cluster
	RDSName = "rds"

	// RDSHttpProxy is the special name for HTTP PROXY route
	RDSHttpProxy = "http_proxy"

	// VirtualListenerName is the name for traffic capture listener
	VirtualListenerName = "virtual"

	// WildcardAddress binds to all IP addresses
	WildcardAddress = "0.0.0.0"

	// LocalhostAddress for local binding
	LocalhostAddress = "127.0.0.1"
)

var (
	// Very verbose output in the logs - full LDS response logged for each sidecar.
	// Use /debug/ldsz instead.
	verboseDebug = os.Getenv("PILOT_DUMP_ALPHA3") != ""
)

// ListenersALPNProtocols denotes the the list of ALPN protocols that the listener
// should expose
var ListenersALPNProtocols = []string{"h2", "http/1.1"}

// BuildListeners produces a list of listeners and referenced clusters for all proxies
func BuildListeners(env model.Environment, node model.Proxy) ([]*xdsapi.Listener, error) {
	switch node.Type {
	case model.Sidecar:
		return buildSidecarListeners(env, node)
	case model.Router:
		// TODO: add listeners for other protocols too
		return buildGatewayHTTPListeners(env, node)
	case model.Ingress:
		return buildLegacyIngressListeners(env, node)
	}
	return nil, nil
}

// buildSidecarListeners produces a list of listeners for sidecar proxies
func buildSidecarListeners(env model.Environment, node model.Proxy) ([]*xdsapi.Listener, error) {

	mesh := env.Mesh
	managementPorts := env.ManagementPorts(node.IPAddress)

	proxyInstances, err := env.GetProxyServiceInstances(node)
	if err != nil {
		return nil, err
	}

	services, err := env.Services()
	if err != nil {
		return nil, err
	}

	// ensure services are ordered to simplify generation logic
	sort.Slice(services, func(i, j int) bool { return services[i].Hostname < services[j].Hostname })

	listeners := make([]*xdsapi.Listener, 0)

	//if node.Type == model.Router {
	//	outbound := buildSidecarOutboundListeners(mesh, node, proxyInstances, services, config)
	//	listeners = append(listeners, outbound...)
	//} else
	if mesh.ProxyListenPort > 0 {
		inbound := buildSidecarInboundListeners(env, node, proxyInstances)
		outbound := buildSidecarOutboundListeners(env, node, proxyInstances, services)

		listeners = append(listeners, inbound...)
		listeners = append(listeners, outbound...)

		mgmtListeners := buildMgmtPortListeners(managementPorts, node.IPAddress)
		// If management listener port and service port are same, bad things happen
		// when running in kubernetes, as the probes stop responding. So, append
		// non overlapping listeners only.
		for i := range mgmtListeners {
			m := mgmtListeners[i]
			l := getByAddress(listeners, m.Address.String())
			if l != nil {
				log.Warnf("Omitting listener for management address %s (%s) due to collision with service listener %s (%s)",
					m.Name, m.Address, l.Name, l.Address)
				continue
			}
			listeners = append(listeners, m)
		}

		// We need a dummy filter to fill in the filter stack for orig_dst listener
		// TODO: Move to Listener filters and set up original dst filter there.
		dummyTCPProxy := &tcp_proxy.TcpProxy{
			StatPrefix: "Dummy",
			Cluster:    "Dummy",
		}

		// add an extra listener that binds to the port that is the recipient of the iptables redirect
		listeners = append(listeners, &xdsapi.Listener{
			Name:           VirtualListenerName,
			Address:        buildAddress(WildcardAddress, uint32(mesh.ProxyListenPort)),
			UseOriginalDst: &google_protobuf.BoolValue{true},
			FilterChains: []listener.FilterChain{
				{
					Filters: []listener.Filter{
						{
							Name:   util.TCPProxy,
							Config: messageToStruct(dummyTCPProxy),
						},
					},
				},
			},
		})
	}

	// enable HTTP PROXY port if necessary; this will add an RDS route for this port
	if mesh.ProxyHttpPort > 0 {
		useRemoteAddress := false
		traceOperation := http_conn.EGRESS
		listenAddress := LocalhostAddress

		if node.Type == model.Router {
			useRemoteAddress = true
			traceOperation = http_conn.INGRESS
			listenAddress = WildcardAddress
		}

		listeners = append(listeners, buildHTTPListener(buildHTTPListenerOpts{
			env:            env,
			proxy:          node,
			proxyInstances: proxyInstances,
			routeConfig:    buildSidecarOutboundHTTPRouteConfig(env, node, proxyInstances, services, RDSHttpProxy),
			ip:             listenAddress,
			port:           int(mesh.ProxyHttpPort),
			//rds:              RDSHttpProxy,
			useRemoteAddress: useRemoteAddress,
			direction:        traceOperation,
			authnPolicy:      nil, /* authN policy is not needed for outbound listener */
		}))
		// TODO: need inbound listeners in HTTP_PROXY case, with dedicated ingress listener.
	}

	return normalizeListeners(listeners), nil
}

// buildSidecarInboundListeners creates listeners for the server-side (inbound)
// configuration for co-located service proxyInstances.
func buildSidecarInboundListeners(env model.Environment, node model.Proxy,
	proxyInstances []*model.ServiceInstance) []*xdsapi.Listener {
	listeners := make([]*xdsapi.Listener, 0, len(proxyInstances))

	mesh := env.Mesh
	config := env.IstioConfigStore
	// inbound connections/requests are redirected to the endpoint address but appear to be sent
	// to the service address.
	for _, instance := range proxyInstances {
		endpoint := instance.Endpoint
		protocol := endpoint.ServicePort.Protocol

		var l *xdsapi.Listener
		authenticationPolicy := model.GetConsolidateAuthenticationPolicy(mesh,
			config, instance.Service.Hostname, instance.Endpoint.ServicePort)

		// Local service instances can be accessed through one of three
		// addresses: localhost, endpoint IP, and service
		// VIP. Localhost bypasses the proxy and doesn't need any TCP
		// route config. Endpoint IP is handled below and Service IP is handled
		// by outbound routes.
		// Traffic sent to our service VIP is redirected by remote
		// services' kubeproxy to our specific endpoint IP.
		switch protocol {
		case model.ProtocolHTTP, model.ProtocolHTTP2, model.ProtocolGRPC:
			l = buildHTTPListener(buildHTTPListenerOpts{
				env:              env,
				proxy:            node,
				proxyInstances:   proxyInstances,
				routeConfig:      buildInboundHTTPRouteConfig(instance),
				ip:               endpoint.Address,
				port:             endpoint.Port,
				rds:              "",
				useRemoteAddress: false,
				direction:        http_conn.INGRESS,
				authnPolicy:      authenticationPolicy,
			})

		case model.ProtocolTCP, model.ProtocolHTTPS, model.ProtocolMongo, model.ProtocolRedis:
			l = buildTCPListener(buildInboundNetworkFilters(instance), endpoint.Address, uint32(endpoint.Port), protocol)

			// TODO: set server-side mixer filter config
			//if mesh.MixerCheckServer != "" || mesh.MixerReportServer != "" {
			//	// config := v1.BuildTCPMixerFilterConfig(mesh, node, instance)
			//	l.FilterChains = append(l.FilterChains, listener.FilterChain{
			//		Filters: []listener.Filter{
			//			{
			//				// TODO(mostrowski): need proto version of mixer config.
			//				// Config: messageToStruct(&config),
			//			},
			//		},
			//	})
			//}

		default:
			log.Debugf("Unsupported inbound protocol %v for port %#v", protocol, instance.Endpoint.ServicePort)
		}

		if l != nil {
			// TODO: move to plugin
			applyInboundAuth(authenticationPolicy, l)
			listeners = append(listeners, l)
		}
	}

	return listeners
}

// buildSidecarOutboundListeners generates http and tcp listeners for outbound connections from the service instance
// TODO(github.com/istio/pilot/issues/237)
//
// Sharing tcp_proxy and http_connection_manager filters on the same port for
// different destination services doesn't work with Envoy (yet). When the
// tcp_proxy filter's route matching fails for the http service the connection
// is closed without falling back to the http_connection_manager.
//
// Temporary workaround is to add a listener for each service IP that requires
// TCP routing
//
// Connections to the ports of non-load balanced services are directed to
// the connection's original destination. This avoids costly queries of instance
// IPs and ports, but requires that ports of non-load balanced service be unique.
func buildSidecarOutboundListeners(env model.Environment, node model.Proxy,
	proxyInstances []*model.ServiceInstance, services []*model.Service) []*xdsapi.Listener {

	var tcpListeners, httpListeners []*xdsapi.Listener

	wildcardListenerPorts := make(map[int]bool)
	for _, service := range services {
		for _, servicePort := range service.Ports {
			clusterName := model.BuildSubsetKey(model.TrafficDirectionOutbound, "",
				service.Hostname, servicePort)

			var addresses []string
			var listenAddress string

			switch servicePort.Protocol {
			case model.ProtocolTCP, model.ProtocolHTTPS, model.ProtocolMongo, model.ProtocolRedis:
				if service.Resolution == model.Passthrough || node.Type == model.Router {
					// ensure only one wildcard listener is created per port if its headless service
					// or if its for a Router (where there is one wildcard TCP listener per port)
					// or if this is in environment where services don't get a dummy load balancer IP.
					if wildcardListenerPorts[servicePort.Port] {
						log.Debugf("Multiple definitions for port %d", servicePort.Port)
						continue
					}
					wildcardListenerPorts[servicePort.Port] = true
					listenAddress = WildcardAddress
					addresses = nil
				} else {
					listenAddress = service.Address
					addresses = []string{service.Address}
				}
				listener := buildTCPListener(buildOutboundNetworkFilters(clusterName, addresses, servicePort),
					listenAddress, uint32(servicePort.Port), servicePort.Protocol)
				tcpListeners = append(tcpListeners, listener)
				// TODO: Set SNI for HTTPS
			case model.ProtocolHTTP2, model.ProtocolHTTP, model.ProtocolGRPC:
				operation := http_conn.EGRESS
				useRemoteAddress := false

				if node.Type == model.Router {
					// if this is in Router mode, then use ingress style trace operation, and remote address settings
					useRemoteAddress = true
					operation = http_conn.INGRESS
				}

				httpListeners = append(httpListeners, buildHTTPListener(buildHTTPListenerOpts{
					env:            env,
					proxy:          node,
					proxyInstances: proxyInstances,
					services:       services,
					ip:             WildcardAddress,
					port:           servicePort.Port,
					//rds:              fmt.Sprintf("%d", servicePort.Port),
					routeConfig:      buildSidecarOutboundHTTPRouteConfig(env, node, proxyInstances, services, fmt.Sprintf("%d", servicePort.Port)),
					useRemoteAddress: useRemoteAddress,
					direction:        operation,
					authnPolicy:      nil, /* authn policy is not needed for outbound listener */
				}))

			}
		}
	}

	return append(tcpListeners, httpListeners...)
}

// buildMgmtPortListeners creates inbound TCP only listeners for the management ports on
// server (inbound). Management port listeners are slightly different from standard Inbound listeners
// in that, they do not have mixer filters nor do they have inbound auth.
// N.B. If a given management port is same as the service instance's endpoint port
// the pod will fail to start in Kubernetes, because the mixer service tries to
// lookup the service associated with the Pod. Since the pod is yet to be started
// and hence not bound to the service), the service lookup fails causing the mixer
// to fail the health check call. This results in a vicious cycle, where kubernetes
// restarts the unhealthy pod after successive failed health checks, and the mixer
// continues to reject the health checks as there is no service associated with
// the pod.
// So, if a user wants to use kubernetes probes with Istio, she should ensure
// that the health check ports are distinct from the service ports.
func buildMgmtPortListeners(managementPorts model.PortList, managementIP string) []*xdsapi.Listener {
	listeners := make([]*xdsapi.Listener, 0, len(managementPorts))

	if managementIP == "" {
		managementIP = "127.0.0.1"
	}

	// assumes that inbound connections/requests are sent to the endpoint address
	for _, mPort := range managementPorts {
		switch mPort.Protocol {
		case model.ProtocolHTTP, model.ProtocolHTTP2, model.ProtocolGRPC, model.ProtocolTCP,
			model.ProtocolHTTPS, model.ProtocolMongo, model.ProtocolRedis:

			instance := &model.ServiceInstance{
				Endpoint: model.NetworkEndpoint{
					Address:     managementIP,
					Port:        mPort.Port,
					ServicePort: mPort,
				},
				Service: &model.Service{
					Hostname: ManagementClusterHostname,
				},
			}

			listeners = append(listeners, buildTCPListener(buildInboundNetworkFilters(instance),
				managementIP, uint32(mPort.Port), model.ProtocolTCP))
		default:
			log.Warnf("Unsupported inbound protocol %v for management port %#v",
				mPort.Protocol, mPort)
		}
	}

	return listeners
}

// TODO: move to plugins
// applyInboundAuth adds ssl_context to the listener if the policy requires one.
func applyInboundAuth(authenticationPolicy *authn.Policy, listener *xdsapi.Listener) {

	if model.RequireTLS(authenticationPolicy) {
		// TODO(mostrowski): figure out SSL
		log.Debugf("TODO Apply authN policy %#v for %#v\n", authenticationPolicy, listener)
	}
}

// options required to build an HTTPListener
type buildHTTPListenerOpts struct { // nolint: maligned
	env            model.Environment
	proxy          model.Proxy
	proxyInstances []*model.ServiceInstance
	services       []*model.Service
	routeConfig    *xdsapi.RouteConfiguration
	ip             string
	port           int
	// bindToPort (default to false) should be set for ingress / gateway
	bindToPort       bool
	rds              string
	useRemoteAddress bool
	direction        http_conn.HttpConnectionManager_Tracing_OperationName
	authnPolicy      *authn.Policy
}

/* // Enable only to compare with RDSv1 responses
func buildDeprecatedHTTPListener(opts buildHTTPListenerOpts) *xdsapi.Listener {
	if opts.rds != "" {
		// Fetch V1 RDS response and stick it into the LDS response
		rc, _ := v1.BuildRDSRoute(opts.env.Mesh, opts.proxy, opts.rds,
			opts.env.ServiceDiscovery, opts.env.IstioConfigStore, true)
		rcBytes, _ := json.Marshal(rc)
		routeConfigDeprecated := string(rcBytes)
		return &xdsapi.Listener{
			Name:    fmt.Sprintf("http_%s_%d", opts.ip, opts.port),
			Address: buildAddress(opts.ip, uint32(opts.port)),
			FilterChains: []listener.FilterChain{
				{
					Filters: []listener.Filter{
						{
							Name: util.HTTPConnectionManager,
							DeprecatedV1: &listener.Filter_DeprecatedV1{
								Type: routeConfigDeprecated,
							},
						},
					},
				},
			},
			DeprecatedV1: &xdsapi.Listener_DeprecatedV1{
				BindToPort: &google_protobuf.BoolValue{
					Value: opts.bindToPort,
				},
			},
		}
	}
	return nil
}
*/

func buildHTTPListener(opts buildHTTPListenerOpts) *xdsapi.Listener {
	mesh := opts.env.Mesh
	filters := []*http_conn.HttpFilter{}

	filters = append(filters, &http_conn.HttpFilter{
		Name: util.CORS,
	})
	// TODO: need alphav3 fault filters.
	// filters = append(filters, buildFaultFilters(opts.config, opts.env, opts.proxy)...)
	filters = append(filters, &http_conn.HttpFilter{
		Name: util.Router,
	})

	/*	TODO(mostrowski): need to port internal build functions for mixer.
		if opts.mesh.MixerCheckServer != "" || opts.mesh.MixerReportServer != "" {
			mixerConfig := v1.BuildHTTPMixerFilterConfig(opts.mesh, opts.proxy, opts.proxyInstances, opts.outboundListener, opts.store)
		filter := &http_conn.HttpFilter{
			Name: v1.MixerFilter,
			Config:messageToStruct(mixerConfig),
		}
			filters = append([]*http_conn.HttpFilter{filter}, filters...)
		}
	*/
	refresh := time.Duration(mesh.RdsRefreshDelay.Seconds) * time.Second
	if refresh == 0 {
		// envoy crashes if 0. Will go away once we move to v2
		refresh = 5 * time.Second
	}

	if filter := buildJwtFilter(opts.authnPolicy); filter != nil {
		filters = append([]*http_conn.HttpFilter{filter}, filters...)
	}

	connectionManager := &http_conn.HttpConnectionManager{
		CodecType: http_conn.AUTO,
		AccessLog: []*accesslog.AccessLog{
			{
				Config: nil,
			},
		},
		HttpFilters:      filters,
		StatPrefix:       HTTPStatPrefix,
		UseRemoteAddress: &google_protobuf.BoolValue{opts.useRemoteAddress},
	}

	if opts.rds != "" {
		rds := &http_conn.HttpConnectionManager_Rds{
			Rds: &http_conn.Rds{
				RouteConfigName: opts.rds,
				ConfigSource: core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_ApiConfigSource{
						ApiConfigSource: &core.ApiConfigSource{
							ApiType:      core.ApiConfigSource_REST_LEGACY,
							ClusterNames: []string{RDSName},
							RefreshDelay: &refresh,
						},
					},
				},
			},
		}
		connectionManager.RouteSpecifier = rds
	} else {
		connectionManager.RouteSpecifier = &http_conn.HttpConnectionManager_RouteConfig{RouteConfig: opts.routeConfig}
	}

	if connectionManager.RouteSpecifier == nil {
		connectionManager.RouteSpecifier = &http_conn.HttpConnectionManager_RouteConfig{
			RouteConfig: opts.routeConfig,
		}
	}

	if mesh.AccessLogFile != "" {
		fl := &accesslog.FileAccessLog{
			Path: mesh.AccessLogFile,
		}

		connectionManager.AccessLog = []*accesslog.AccessLog{
			{
				Config: messageToStruct(fl),
				Name:   fileAccessLog,
			},
		}
	}

	if mesh.EnableTracing {
		connectionManager.Tracing = &http_conn.HttpConnectionManager_Tracing{
			OperationName: opts.direction,
		}
		connectionManager.GenerateRequestId = &google_protobuf.BoolValue{true}
	}

	if verboseDebug {
		connectionManagerJSON, _ := json.MarshalIndent(connectionManager, "  ", "  ")
		log.Infof("LDS: %s \n", string(connectionManagerJSON))
	}
	l := &xdsapi.Listener{
		Name:    fmt.Sprintf("http_%s_%d", opts.ip, opts.port),
		Address: buildAddress(opts.ip, uint32(opts.port)),
		FilterChains: []listener.FilterChain{
			{
				Filters: []listener.Filter{
					{
						Name:   envoyHTTPConnectionManager,
						Config: messageToStruct(connectionManager),
					},
				},
			},
		},
	}

	if !opts.bindToPort {
		l.DeprecatedV1 = &xdsapi.Listener_DeprecatedV1{
			BindToPort: &google_protobuf.BoolValue{
				Value: opts.bindToPort,
			},
		}
	}
	return l
}

// buildTCPListener constructs a listener for the TCP proxy
func buildTCPListener(filters []listener.Filter, ip string, port uint32, protocol model.Protocol) *xdsapi.Listener {
	filterChain := listener.FilterChain{
		Filters: filters,
	}
	return &xdsapi.Listener{
		Name:    fmt.Sprintf("%s_%s_%d", protocol, ip, port),
		Address: buildAddress(ip, port),
		FilterChains: []listener.FilterChain{
			filterChain,
		},
		DeprecatedV1: &xdsapi.Listener_DeprecatedV1{
			BindToPort: &google_protobuf.BoolValue{
				Value: false,
			},
		},
	}
}

// TODO: find a proper home for these http related functions
// buildDefaultHTTPRoute builds a default route.
func buildDefaultHTTPRoute(clusterName string) *route.Route {
	return &route.Route{
		Match: route.RouteMatch{PathSpecifier: &route.RouteMatch_Prefix{Prefix: "/"}},
		Decorator: &route.Decorator{
			Operation: DefaultOperation,
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{Cluster: clusterName},
			},
		},
	}
}

// buildInboundHTTPRouteConfig builds the route config with a single wildcard virtual host on the inbound path
// TODO: enable mixer configuration, websockets, trace decorators
func buildInboundHTTPRouteConfig(instance *model.ServiceInstance) *xdsapi.RouteConfiguration {
	clusterName := model.BuildSubsetKey(model.TrafficDirectionInbound, "",
		instance.Service.Hostname, instance.Endpoint.ServicePort)
	defaultRoute := buildDefaultHTTPRoute(clusterName)

	inboundVHost := route.VirtualHost{
		Name:    fmt.Sprintf("%s|http|%d", model.TrafficDirectionInbound, instance.Endpoint.ServicePort.Port),
		Domains: []string{"*"},
		Routes:  []route.Route{*defaultRoute},
	}

	// TODO: mixer disabled for now as its configuration is still in old format
	// set server-side mixer filter config for inbound HTTP routes
	//if mesh.MixerCheckServer != "" || mesh.MixerReportServer != "" {
	//	defaultRoute.OpaqueConfig = v1.BuildMixerOpaqueConfig(!mesh.DisablePolicyChecks, false, instance.Service.Hostname)
	//}

	return &xdsapi.RouteConfiguration{
		Name:         clusterName,
		VirtualHosts: []route.VirtualHost{inboundVHost},
		ValidateClusters: &google_protobuf.BoolValue{
			Value: false,
		},
	}
}

func buildSidecarOutboundHTTPRouteConfig(env model.Environment, node model.Proxy, _ []*model.ServiceInstance,
	services []*model.Service, routeName string) *xdsapi.RouteConfiguration {

	port := 0
	if routeName != RDSHttpProxy {
		var err error
		port, err = strconv.Atoi(routeName)
		if err != nil {
			return nil
		}
	}

	nameToServiceMap := make(map[string]*model.Service)
	for _, svc := range services {
		if port == 0 {
			nameToServiceMap[svc.Hostname] = svc
		} else {
			if svcPort, exists := svc.Ports.GetByPort(port); exists {
				nameToServiceMap[svc.Hostname] = &model.Service{
					Hostname: svc.Hostname,
					Address:  svc.Address,
					Ports:    []*model.Port{svcPort},
				}
			}
		}
	}

	// Get list of virtual services bound to the mesh gateway
	virtualServices := env.VirtualServices([]string{model.IstioMeshGateway})
	// TODO: Need to trim output based on source label/gateway match
	guardedHosts := TranslateVirtualHosts(virtualServices,
		nameToServiceMap, nil, node.Domain)
	vHostPortMap := make(map[int][]route.VirtualHost)

	// there should be only one guarded host in the return val since we supplied services with just one port
	for _, guardedHost := range guardedHosts {
		routes := make([]route.Route, 0)
		for _, r := range guardedHost.Routes {
			routes = append(routes, r.Route)
		}

		virtualHosts := make([]route.VirtualHost, 0)

		for _, host := range guardedHost.Hosts {
			virtualHosts = append(virtualHosts, route.VirtualHost{
				Name:    fmt.Sprintf("%s:%d", host, guardedHost.Port),
				Domains: []string{host},
				Routes:  routes,
			})
		}

		for _, svc := range guardedHost.Services {
			domains := generateAltVirtualHosts(svc.Hostname, guardedHost.Port)
			if len(svc.Address) > 0 {
				// add a vhost match for the IP (if its non CIDR)
				cidr := convertAddressToCidr(svc.Address)
				if cidr.PrefixLen.Value == 32 {
					domains = append(domains, svc.Address)
					domains = append(domains, fmt.Sprintf("%s:%d", svc.Address, guardedHost.Port))
				}
			}
			virtualHosts = append(virtualHosts, route.VirtualHost{
				Name:    fmt.Sprintf("%s:%d", svc.Hostname, guardedHost.Port),
				Domains: domains,
				Routes:  routes,
			})
		}

		vHostPortMap[guardedHost.Port] = virtualHosts
	}

	var virtualHosts []route.VirtualHost
	if routeName == RDSHttpProxy {
		virtualHosts = mergeAllVirtualHosts(vHostPortMap)
	} else {
		virtualHosts = vHostPortMap[port]
	}

	out := &xdsapi.RouteConfiguration{
		Name:         fmt.Sprintf("%d", port),
		VirtualHosts: virtualHosts,
		ValidateClusters: &google_protobuf.BoolValue{
			Value: false, // until we have rds
		},
	}
	return out

}

// Given a service, and a port, this function generates all possible HTTP Host headers.
// For example, a service of the form foo.local.campus.net on port 80 could be accessed as
// http://foo:80 within the .local network, as http://foo.local:80 (by other clients in the campus.net domain),
// as http://foo.local.campus:80, etc.
func generateAltVirtualHosts(hostname string, port int) []string {
	vhosts := []string{hostname, fmt.Sprintf("%s:%d", hostname, port)}
	for i := len(hostname) - 1; i >= 0; i-- {
		if hostname[i] == '.' {
			variant := hostname[:i]
			variantWithPort := fmt.Sprintf("%s:%d", variant, port)
			vhosts = append(vhosts, variant)
			vhosts = append(vhosts, variantWithPort)
		}
	}
	return vhosts
}

// mergeAllVirtualHosts across all ports. On routes for ports other than port 80,
// virtual hosts without an explicit port suffix (IP:PORT) should be stripped
func mergeAllVirtualHosts(vHostPortMap map[int][]route.VirtualHost) []route.VirtualHost {
	var virtualHosts []route.VirtualHost
	for p, vhosts := range vHostPortMap {
		if p == 80 {
			virtualHosts = append(virtualHosts, vhosts...)
		} else {
			for _, vhost := range vhosts {
				var newDomains []string
				for _, domain := range vhost.Domains {
					if strings.Contains(domain, ":") {
						newDomains = append(newDomains, domain)
					}
				}
				if len(newDomains) > 0 {
					vhost.Domains = newDomains
					virtualHosts = append(virtualHosts, vhost)
				}
			}
		}
	}
	return virtualHosts
}
