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
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"

	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
)

// Headers with special meaning in Envoy
const (
	HeaderMethod    = ":method"
	HeaderAuthority = ":authority"
	HeaderScheme    = ":scheme"
)

const (
	// UnresolvedCluster for destinations pointing to unknown clusters.
	UnresolvedCluster = "unresolved-cluster"

	// DefaultOperation is the default decorator
	DefaultOperation = "default-operation"
)

// ServiceByName claims a service entry from the registry using a host name.
type ServiceByName func(host string, contextNamespace string) *model.Service

// GuardedHost is a context-dependent virtual host entry with guarded routes.
type GuardedHost struct {
	// Port is the capture port (e.g. service port)
	Port int

	// Services are the services matching the virtual host.
	// The service host names need to be contextualized by the source.
	Services []*model.Service

	// Hosts is a list of alternative literal host names for the host.
	Hosts []string

	// Routes in the virtual host
	Routes []GuardedRoute

	// Hack: RouteConfiguration for this port
	RouteConfiguration *v2.RouteConfiguration
}

// TranslateServiceHostname matches a host against a model service.
// This cannot be externalized to core model until the registries understand namespaces.
func TranslateServiceHostname(services map[string]*model.Service, clusterDomain string) ServiceByName {
	return func(host string, contextNamespace string) *model.Service {
		if strings.Contains(host, ".") {
			return services[host]
		}

		return services[fmt.Sprintf("%s.%s.%s", host, contextNamespace, clusterDomain)]
	}
}

// TranslateVirtualHosts creates the entire routing table for Istio v1alpha3 configs.
// Services are indexed by FQDN hostnames.
// Cluster domain is used to resolve short service names (e.g. "svc.cluster.local").
func TranslateVirtualHosts(
	virtualServiceSpecs []model.Config,
	services map[string]*model.Service,
	clusterDomain string) map[int]GuardedHost {
	out := make([]GuardedHost, 0)
	serviceByName := TranslateServiceHostname(services, clusterDomain)

	// translate all virtual service configs
	for _, config := range virtualServiceSpecs {
		out = append(out, TranslateVirtualHost(config, serviceByName)...)
	}

	// compute services missing service configs
	missing := make(map[string]bool)
	for fqdn := range services {
		missing[fqdn] = true
	}
	for _, host := range out {
		for _, service := range host.Services {
			delete(missing, service.Hostname)
		}
	}

	// append default hosts for the service missing virtual services
	for fqdn := range missing {
		svc := services[fqdn]
		for _, port := range svc.Ports {
			if port.Protocol.IsHTTP() {
				cluster := model.BuildSubsetKey(model.TrafficDirectionOutbound, "", svc.Hostname, port)
				out = append(out, GuardedHost{
					Port:     port.Port,
					Services: []*model.Service{svc},
					Routes: []GuardedRoute{{
						Route: route.Route{
							Match:     TranslateHTTPRouteMatch(nil),
							Decorator: &route.Decorator{Operation: DefaultOperation},
							Action: &route.Route_Route{
								Route: &route.RouteAction{
									ClusterSpecifier: &route.RouteAction_Cluster{Cluster: cluster},
								},
							},
						},
					}},
				})
			}
		}
	}

	routePortMap := make(map[int]*GuardedHost)

	for _, guardedHost := range out {
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
			domains := generateAltVirtualHosts(svc, guardedHost.Port)
			virtualHosts = append(virtualHosts, route.VirtualHost{
				Name:    fmt.Sprintf("%s:%d", svc.Hostname, guardedHost.Port),
				Domains: domains,
				Routes:  routes,
			})
		}

		guardedHost.RouteConfiguration = &v2.RouteConfiguration{
			Name:         fmt.Sprintf("%d", guardedHost.Port),
			VirtualHosts: virtualHosts,
			ValidateClusters: &types.BoolValue{
				Value: false, // until we have rds
			},
		}
		routePortMap[guardedHost.Port] = &guardedHost
	}

	return routePortMap
}

// MatchServiceHosts splits the virtual service hosts into services and literal hosts
func MatchServiceHosts(in model.Config, serviceByName ServiceByName) ([]string, []*model.Service) {
	rule := in.Spec.(*networking.VirtualService)
	hosts := make([]string, 0)
	services := make([]*model.Service, 0)
	for _, host := range rule.Hosts {
		if svc := serviceByName(host, in.ConfigMeta.Namespace); svc != nil {
			services = append(services, svc)
		} else {
			hosts = append(hosts, host)
		}
	}
	return hosts, services
}

// TranslateVirtualHost creates virtual hosts corresponding to a virtual service.
func TranslateVirtualHost(in model.Config, serviceByName ServiceByName) []GuardedHost {
	hosts, services := MatchServiceHosts(in, serviceByName)
	serviceByPort := make(map[int][]*model.Service)
	for _, svc := range services {
		for _, port := range svc.Ports {
			if port.Protocol.IsHTTP() {
				serviceByPort[port.Port] = append(serviceByPort[port.Port], svc)
			}
		}
	}

	// if no services matched, then we have no port information -- default to 80 for now
	// TODO: use match condition ports
	if len(serviceByPort) == 0 {
		serviceByPort[80] = nil
	}

	out := make([]GuardedHost, len(serviceByPort))
	for port, services := range serviceByPort {
		clusterNaming := TranslateDestination(serviceByName, in.ConfigMeta.Namespace, port)
		routes := TranslateHTTPRoutes(in, clusterNaming)
		out = append(out, GuardedHost{
			Port:     port,
			Services: services,
			Hosts:    hosts,
			Routes:   routes,
		})
	}

	return out
}

// TranslateDestination produces a cluster naming function using the config context.
func TranslateDestination(
	serviceByName ServiceByName,
	contextNamespace string,
	defaultPort int) ClusterNaming {
	return func(destination *networking.Destination) string {
		// detect if it is a service
		svc := serviceByName(destination.Name, contextNamespace)

		// TODO: create clusters for non-service hostnames/IPs
		if svc == nil {
			return UnresolvedCluster
		}

		// default port uses port number
		svcPort, _ := svc.Ports.GetByPort(defaultPort)
		if destination.Port != nil {
			switch selector := destination.Port.Port.(type) {
			case *networking.PortSelector_Name:
				svcPort, _ = svc.Ports.Get(selector.Name)
			case *networking.PortSelector_Number:
				svcPort, _ = svc.Ports.GetByPort(int(selector.Number))
			}
		}

		if svcPort == nil {
			return UnresolvedCluster
		}

		// use subsets if it is a service
		return model.BuildSubsetKey(model.TrafficDirectionOutbound, destination.Subset, svc.Hostname, svcPort)
	}
}

// ClusterNaming specifies cluster name for a destination
type ClusterNaming func(*networking.Destination) string

// GuardedRoute are routes for a destination guarded by deployment conditions.
type GuardedRoute struct {
	Route route.Route

	// SourceLabels guarding the route
	SourceLabels map[string]string

	// Gateways pre-condition
	Gateways []string
}

// TranslateHTTPRoutes creates virtual host routes from the v1alpha3 config.
// The rule should be adapted to destination names (outbound clusters).
// Each rule is guarded by source labels.
func TranslateHTTPRoutes(in model.Config, name ClusterNaming) []GuardedRoute {
	rule, ok := in.Spec.(*networking.VirtualService)
	if !ok {
		return nil
	}

	operation := in.ConfigMeta.Name

	out := make([]GuardedRoute, 0)
	for _, http := range rule.Http {
		if len(http.Match) == 0 {
			out = append(out, TranslateHTTPRoute(http, nil, operation, name))
		} else {
			for _, match := range http.Match {
				out = append(out, TranslateHTTPRoute(http, match, operation, name))
			}
		}
	}

	return out
}

// TranslateHTTPRoute translates HTTP routes
// TODO: fault filters -- issue https://github.com/istio/api/issues/388
func TranslateHTTPRoute(in *networking.HTTPRoute,
	match *networking.HTTPMatchRequest,
	operation string,
	name ClusterNaming) GuardedRoute {
	out := route.Route{
		Match: TranslateHTTPRouteMatch(match),
		Decorator: &route.Decorator{
			Operation: operation,
		},
	}

	if redirect := in.Redirect; redirect != nil {
		out.Action = &route.Route_Redirect{
			Redirect: &route.RedirectAction{
				HostRedirect: redirect.Authority,
				PathRewriteSpecifier: &route.RedirectAction_PathRedirect{
					PathRedirect: redirect.Uri,
				},
			}}
	} else {
		action := &route.RouteAction{
			Cors:         TranslateCORSPolicy(in.CorsPolicy),
			RetryPolicy:  TranslateRetryPolicy(in.Retries),
			Timeout:      TranslateTime(in.Timeout),
			UseWebsocket: &types.BoolValue{Value: in.WebsocketUpgrade},
		}
		out.Action = &route.Route_Route{Route: action}

		if rewrite := in.Rewrite; rewrite != nil {
			action.PrefixRewrite = rewrite.Uri
			action.HostRewriteSpecifier = &route.RouteAction_HostRewrite{
				HostRewrite: rewrite.Authority,
			}
		}

		if len(in.AppendHeaders) > 0 {
			action.RequestHeadersToAdd = make([]*core.HeaderValueOption, len(in.AppendHeaders))
			for key, value := range in.AppendHeaders {
				action.RequestHeadersToAdd = append(action.RequestHeadersToAdd, &core.HeaderValueOption{
					Header: &core.HeaderValue{
						Key:   key,
						Value: value,
					},
				})
			}
		}

		if in.Mirror != nil {
			action.RequestMirrorPolicy = &route.RouteAction_RequestMirrorPolicy{Cluster: name(in.Mirror)}
		}

		weighted := make([]*route.WeightedCluster_ClusterWeight, len(in.Route))
		for _, dst := range in.Route {
			weighted = append(weighted, &route.WeightedCluster_ClusterWeight{
				Name:   name(dst.Destination),
				Weight: &types.UInt32Value{Value: uint32(dst.Weight)},
			})
		}

		// rewrite to a single cluster if there is only weighted cluster
		if len(weighted) == 1 {
			action.ClusterSpecifier = &route.RouteAction_Cluster{Cluster: weighted[0].Name}
		} else {
			action.ClusterSpecifier = &route.RouteAction_WeightedClusters{
				WeightedClusters: &route.WeightedCluster{
					Clusters: weighted,
				},
			}
		}
	}

	return GuardedRoute{
		Route:        out,
		SourceLabels: match.GetSourceLabels(),
		Gateways:     match.GetGateways(),
	}
}

// TranslateHTTPRouteMatch translates match condition
func TranslateHTTPRouteMatch(in *networking.HTTPMatchRequest) route.RouteMatch {
	out := route.RouteMatch{PathSpecifier: &route.RouteMatch_Prefix{Prefix: "/"}}
	if in == nil {
		return out
	}

	for name, stringMatch := range in.Headers {
		matcher := TranslateHeaderMatcher(name, stringMatch)
		out.Headers = append(out.Headers, &matcher)
	}

	// guarantee ordering of headers
	sort.Slice(out.Headers, func(i, j int) bool {
		if out.Headers[i].Name == out.Headers[j].Name {
			return out.Headers[i].Value < out.Headers[j].Value
		}
		return out.Headers[i].Name < out.Headers[j].Name
	})

	if in.Uri != nil {
		switch m := in.Uri.MatchType.(type) {
		case *networking.StringMatch_Exact:
			out.PathSpecifier = &route.RouteMatch_Path{Path: m.Exact}
		case *networking.StringMatch_Prefix:
			out.PathSpecifier = &route.RouteMatch_Prefix{Prefix: m.Prefix}
		case *networking.StringMatch_Regex:
			out.PathSpecifier = &route.RouteMatch_Regex{Regex: m.Regex}
		}
	}

	if in.Method != nil {
		matcher := TranslateHeaderMatcher(HeaderMethod, in.Method)
		out.Headers = append(out.Headers, &matcher)
	}

	if in.Authority != nil {
		matcher := TranslateHeaderMatcher(HeaderAuthority, in.Authority)
		out.Headers = append(out.Headers, &matcher)
	}

	if in.Scheme != nil {
		matcher := TranslateHeaderMatcher(HeaderScheme, in.Scheme)
		out.Headers = append(out.Headers, &matcher)
	}

	// TODO: match.DestinationPorts

	return out
}

// TranslateHeaderMatcher translates to HeaderMatcher
func TranslateHeaderMatcher(name string, in *networking.StringMatch) route.HeaderMatcher {
	out := route.HeaderMatcher{
		Name: name,
	}

	switch m := in.MatchType.(type) {
	case *networking.StringMatch_Exact:
		out.Value = m.Exact
	case *networking.StringMatch_Prefix:
		// Envoy regex grammar is ECMA-262 (http://en.cppreference.com/w/cpp/regex/ecmascript)
		// Golang has a slightly different regex grammar
		out.Value = fmt.Sprintf("^%s.*", regexp.QuoteMeta(m.Prefix))
		out.Regex = &types.BoolValue{Value: true}
	case *networking.StringMatch_Regex:
		out.Value = m.Regex
		out.Regex = &types.BoolValue{Value: true}
	}

	return out
}

// TranslateRetryPolicy translates retry policy
func TranslateRetryPolicy(in *networking.HTTPRetry) *route.RouteAction_RetryPolicy {
	if in != nil && in.Attempts > 0 {
		return &route.RouteAction_RetryPolicy{
			NumRetries:    &types.UInt32Value{Value: uint32(in.GetAttempts())},
			RetryOn:       "5xx,connect-failure,refused-stream",
			PerTryTimeout: TranslateTime(in.PerTryTimeout),
		}
	}
	return nil
}

// TranslateCORSPolicy translates CORS policy
func TranslateCORSPolicy(in *networking.CorsPolicy) *route.CorsPolicy {
	if in == nil {
		return nil
	}

	out := route.CorsPolicy{
		AllowOrigin: in.AllowOrigin,
		Enabled:     &types.BoolValue{Value: true},
	}
	out.AllowCredentials = in.AllowCredentials
	out.AllowHeaders = strings.Join(in.AllowHeaders, ",")
	out.AllowMethods = strings.Join(in.AllowMethods, ",")
	out.ExposeHeaders = strings.Join(in.ExposeHeaders, ",")
	if in.MaxAge != nil {
		out.MaxAge = in.MaxAge.String()
	}
	return &out
}

// TranslateTime converts time protos.
func TranslateTime(in *types.Duration) *time.Duration {
	if in == nil {
		return nil
	}
	out, err := types.DurationFromProto(in)
	if err != nil {
		log.Warnf("error converting duration %#v, using 0: %v", in, err)
	}
	return &out
}

// buildDefaultHTTPRoute builds a default route.
func buildDefaultHTTPRoute(clusterName string) *route.Route {
	return &route.Route{
		Match: TranslateHTTPRouteMatch(nil),
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
func buildInboundHTTPRouteConfig(instance *model.ServiceInstance) *v2.RouteConfiguration {
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

	return &v2.RouteConfiguration{
		Name:         clusterName,
		VirtualHosts: []route.VirtualHost{inboundVHost},
		ValidateClusters: &types.BoolValue{
			Value: false,
		},
	}
}

func last(arr []string) string {
	return arr[len(arr)-1]
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func generateAltVirtualHosts(svc *model.Service, port int) []string {
	dots := len(strings.Split(svc.Hostname, "."))
	vhosts := make([]string, 0)
	nameToSplit := reverse(svc.Hostname)
	for i := 1; i <= dots; i++ {
		variant := reverse(last(strings.SplitN(nameToSplit, ".", i)))
		variantWithPort := fmt.Sprintf("%s:%d", variant, port)
		vhosts = append(vhosts, variant)
		vhosts = append(vhosts, variantWithPort)
	}
	return vhosts
}
