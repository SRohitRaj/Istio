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

// Mixer filter configuration

package v1

import (
	"crypto/sha256"
	"encoding/base64"
	"net"
	"net/url"
	"sort"
	"strings"
	// TODO(nmittler): Remove this
	_ "github.com/golang/glog"

	meshconfig "istio.io/api/mesh/v1alpha1"
	mpb "istio.io/api/mixer/v1"
	mccpb "istio.io/api/mixer/v1/config/client"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
	"github.com/gogo/protobuf/proto"
)

const (
	// MixerCheckClusterName is the name of the mixer cluster used for policy checks
	MixerCheckClusterName = "mixer_check_server"

	// MixerReportClusterName is the name of the mixer cluster used for telemetry
	MixerReportClusterName = "mixer_report_server"

	// MixerFilter name and its attributes
	MixerFilter = "mixer"

	// AttrSourcePrefix all source attributes start with this prefix
	AttrSourcePrefix = "source"

	// AttrSourceIP is client source IP
	AttrSourceIP = "source.ip"

	// AttrSourceUID is platform-specific unique identifier for the client instance of the source service
	AttrSourceUID = "source.uid"

	// Labels associated with the source
	AttrSourceLabels = "source.labels"

	// AttrDestinationPrefix all destination attributes start with this prefix
	AttrDestinationPrefix = "destination"

	// AttrDestinationIP is the server source IP
	AttrDestinationIP = "destination.ip"

	// AttrDestinationUID is platform-specific unique identifier for the server instance of the target service
	AttrDestinationUID = "destination.uid"

	// Labels associated with the destination
	AttrDestinationLabels = "destination.labels"

	// AttrDestinationService is name of the target service
	AttrDestinationService = "destination.service"

	// MixerRequestCount is the quota bucket name
	MixerRequestCount = "RequestCount"

	// MixerCheck switches Check call on and off
	MixerCheck = "mixer_check"

	// MixerReport switches Report call on and off
	MixerReport = "mixer_report"

	// DisableTCPCheckCalls switches Check call on and off for tcp listeners
	DisableTCPCheckCalls = "disable_tcp_check_calls"

	// MixerForward switches attribute forwarding on and off
	MixerForward = "mixer_forward"
)

// FilterMixerConfig definition.
//
// NOTE: all fields marked as DEPRECATED are part of the original v1
// mixerclient configuration. They are deprecated and will be
// eventually removed once proxies are updated.
//
// Going forwards all mixerclient configuration should represeted by
// istio.io/api/mixer/v1/config/client/mixer_filter_config.proto and
// encoded in the `V2` field below.
//
type FilterMixerConfig struct {
	// DEPRECATED: MixerAttributes specifies the static list of attributes that are sent with
	// each request to Mixer.
	MixerAttributes map[string]string `json:"mixer_attributes,omitempty"`

	// DEPRECATED: ForwardAttributes specifies the list of attribute keys and values that
	// are forwarded as an HTTP header to the server side proxy
	ForwardAttributes map[string]string `json:"forward_attributes,omitempty"`

	// DEPRECATED: QuotaName specifies the name of the quota bucket to withdraw tokens from;
	// an empty name means no quota will be charged.
	QuotaName string `json:"quota_name,omitempty"`

	// DEPRECATED: If set to true, disables mixer check calls for TCP connections
	DisableTCPCheckCalls bool `json:"disable_tcp_check_calls,omitempty"`

	// istio.io/api/mixer/v1/config/client configuration protobuf
	// encoded as a generic map using canonical JSON encoding.
	//
	// If `V2` field is not empty, the DEPRECATED fields above should
	// be discarded.
	V2 map[string]interface{} `json:"v2,omitempty"`
}

func (*FilterMixerConfig) isNetworkFilterConfig() {}

// buildMixerCluster builds an outbound mixer cluster of a given name
func buildMixerCluster(mesh *meshconfig.MeshConfig, mixerSAN []string, server, clusterName string) *Cluster {
	cluster := buildCluster(server, clusterName, mesh.ConnectTimeout)
	cluster.CircuitBreaker = &CircuitBreaker{
		Default: DefaultCBPriority{
			MaxPendingRequests: 10000,
			MaxRequests:        10000,
		},
	}

	cluster.Features = ClusterFeatureHTTP2
	// apply auth policies
	switch mesh.DefaultConfig.ControlPlaneAuthPolicy {
	case meshconfig.AuthenticationPolicy_NONE:
		// do nothing
	case meshconfig.AuthenticationPolicy_MUTUAL_TLS:
		// apply SSL context to enable mutual TLS between Envoy proxies between app and mixer
		cluster.SSLContext = buildClusterSSLContext(model.AuthCertsPath, mixerSAN)
	}

	return cluster
}

// buildMixerClusters builds an outbound mixer cluster with configured check/report clusters
func buildMixerClusters(mesh *meshconfig.MeshConfig, role model.Node, mixerSAN []string) []*Cluster {
	mixerClusters := make([]*Cluster, 0)

	if mesh.MixerCheckServer != "" {
		mixerClusters = append(mixerClusters, buildMixerCluster(mesh, mixerSAN, mesh.MixerCheckServer, MixerCheckClusterName))
	}

	if mesh.MixerReportServer != "" {
		// if both fields point to same server, reuse the cluster
		if mesh.MixerReportServer == mesh.MixerCheckServer {
			return mixerClusters
		}
		mixerClusters = append(mixerClusters, buildMixerCluster(mesh, mixerSAN, mesh.MixerReportServer, MixerReportClusterName))
	}

	return mixerClusters
}

// buildMixerConfig build per route mixer config to be deployed at the `model.Node` workload
// with destination of Service `dest` and `destName` as the service name
func buildMixerConfig(source model.Node, destName string, dest *model.Service, config model.IstioConfigStore) map[string]string {
	var err error
	var cfg string
	var ba []byte

	sc := serviceConfig(&model.ServiceInstance{Service: dest}, config, false, false)
	addStandardNodeAttributes(sc.MixerAttributes.Attributes, AttrSourcePrefix, source, nil)
	oc := map[string]string{}
	oc[AttrDestinationService] = destName

	if ba, err = proto.Marshal(sc); err != nil {
		log.Warnf("Unable to marshal service config: %v", err)
		return oc
	}

	oc["mixer"] = base64.StdEncoding.EncodeToString(ba)
	h := sha256.New()
	_, _ = h.Write(ba)
	oc["mixer_sha"] = base64.StdEncoding.EncodeToString(h.Sum(nil))

	if cfg, err = model.ToJSON(sc); err == nil {
		oc["mixer_debug"] = base64.StdEncoding.EncodeToString([]byte(cfg))
	}
	return oc
}

func buildMixerOpaqueConfig(check, forward bool, destinationService string) map[string]string {
	keys := map[bool]string{true: "on", false: "off"}
	m := map[string]string{
		MixerReport:  "on",
		MixerCheck:   keys[check],
		MixerForward: keys[forward],
	}
	if destinationService != "" {
		m[AttrDestinationService] = destinationService
	}
	return m
}

// Mixer filter uses outbound configuration by default (forward attributes,
// but not invoke check calls)  ServiceInstances belong to the Node.
func mixerHTTPRouteConfig(mesh *meshconfig.MeshConfig, role model.Node, instances []*model.ServiceInstance, outboundRoute bool, config model.IstioConfigStore) *FilterMixerConfig { // nolint: lll
	filter := &FilterMixerConfig{
		MixerAttributes: map[string]string{
			AttrDestinationIP:  role.IPAddress,
			AttrDestinationUID: "kubernetes://" + role.ID,
		},
		ForwardAttributes: map[string]string{
			AttrSourceIP:  role.IPAddress,
			AttrSourceUID: "kubernetes://" + role.ID,
		},
		QuotaName: MixerRequestCount,
	}

	transport := &mccpb.TransportConfig{
		CheckCluster:  MixerCheckClusterName,
		ReportCluster: MixerReportClusterName,
	}
	if mesh.MixerCheckServer == mesh.MixerReportServer {
		transport.ReportCluster = transport.CheckCluster
	}

	v2 := &mccpb.HttpClientConfig{
		MixerAttributes: &mpb.Attributes{
			Attributes: map[string]*mpb.Attributes_AttributeValue{},
		},
		ServiceConfigs: map[string]*mccpb.ServiceConfig{},
		Transport:      transport,
	}

	var labels map[string]string
	if len(instances)>0 {
		labels = instances[0].Labels
	}
	addStandardNodeAttributes(v2.MixerAttributes.Attributes, AttrDestinationPrefix, role, labels)

	if role.Type == model.Sidecar && !outboundRoute {
		// Don't forward mixer attributes to the app from inbound sidecar routes
	} else {
		v2.ForwardAttributes = &mpb.Attributes{
			Attributes: map[string]*mpb.Attributes_AttributeValue{},
		}
		addStandardNodeAttributes(v2.ForwardAttributes.Attributes, AttrSourcePrefix, role, labels)
	}

	if len(instances) > 0 {
		// legacy mixerclient behavior is a comma separated list of
		// services. When can this be removed?
		var services []string
		if instances != nil {
			serviceSet := make(map[string]bool, len(instances))
			for _, instance := range instances {
				serviceSet[instance.Service.Hostname] = true
			}
			for service := range serviceSet {
				services = append(services, service)
			}
			sort.Strings(services)
		}
		filter.MixerAttributes[AttrDestinationService] = strings.Join(services, ",")

		// first service in the sorted list is the default
		v2.DefaultDestinationService = services[0]
	}

	for _, instance := range instances {
		v2.ServiceConfigs[instance.Service.Hostname] = serviceConfig(instance, config,
			outboundRoute || mesh.DisablePolicyChecks, outboundRoute)
	}

	if v2JSONMap, err := model.ToJSONMap(v2); err != nil {
		log.Warnf("Could not encode v2 HTTP mixerclient filter for node %q: %v", role, err)
	} else {
		filter.V2 = v2JSONMap
	}
	return filter
}

// addStandardNodeAttributes add standard node attributes with the given prefix
func addStandardNodeAttributes(attr map[string]*mpb.Attributes_AttributeValue, prefix string, node model.Node, labels map[string]string) {
	attr[prefix+".ip"] = &mpb.Attributes_AttributeValue{
		Value: &mpb.Attributes_AttributeValue_BytesValue{net.ParseIP(node.IPAddress)},
	}

	attr[prefix+".uid"] = &mpb.Attributes_AttributeValue{
		Value: &mpb.Attributes_AttributeValue_StringValue{"kubernetes://" + node.ID},
	}

	if len(labels) > 0 {
		attr[prefix+".labels"] = &mpb.Attributes_AttributeValue{
			Value: &mpb.Attributes_AttributeValue_StringMapValue{
				StringMapValue: &mpb.Attributes_StringMap{Entries: labels},
			},
		}
	}
}

// generate serviceConfig for a given instance
func serviceConfig(dest *model.ServiceInstance, config model.IstioConfigStore, disableCheck, disableReport bool) *mccpb.ServiceConfig{
	sc := &mccpb.ServiceConfig{
		MixerAttributes: &mpb.Attributes{
			Attributes: map[string]*mpb.Attributes_AttributeValue{
				AttrDestinationService: {
					Value: &mpb.Attributes_AttributeValue_StringValue{StringValue: dest.Service.Hostname},
				},
			},
		},
		DisableCheckCalls:  disableCheck,
		DisableReportCalls: disableReport,
	}

	if len(dest.Labels) > 0 {
		sc.MixerAttributes.Attributes[AttrDestinationLabels] = &mpb.Attributes_AttributeValue{
			Value: &mpb.Attributes_AttributeValue_StringMapValue{
				StringMapValue: &mpb.Attributes_StringMap{Entries: dest.Labels},
			},
		}
	}

	// omit API, Quota, and Auth portion of service config when
	// check is disabled.
	if disableCheck {
		return sc
	}

	apiSpecs := config.HTTPAPISpecByDestination(dest)
	model.SortHTTPAPISpec(apiSpecs)
	for _, config := range apiSpecs {
		sc.HttpApiSpec = append(sc.HttpApiSpec, config.Spec.(*mccpb.HTTPAPISpec))
	}

	quotaSpecs := config.QuotaSpecByDestination(dest)
	model.SortQuotaSpec(quotaSpecs)
	for _, config := range quotaSpecs {
		sc.QuotaSpec = append(sc.QuotaSpec, config.Spec.(*mccpb.QuotaSpec))
	}

	authSpecs := config.EndUserAuthenticationPolicySpecByDestination(dest)
	model.SortEndUserAuthenticationPolicySpec(quotaSpecs)
	if len(authSpecs) > 0 {
		spec := (authSpecs[0].Spec).(*mccpb.EndUserAuthenticationPolicySpec)

		// Update jwks_uri_envoy_cluster This cluster should be
		// created elsewhere using the same host-to-cluster naming
		// scheme, i.e. buildJWKSURIClusterNameAndAddress.
		for _, jwt := range spec.Jwts {
			if name, _, _, err := buildJWKSURIClusterNameAndAddress(jwt.JwksUri); err != nil {
				log.Warnf("Could not set jwks_uri_envoy and address for jwks_uri %q: %v",
					jwt.JwksUri, err)
			} else {
				jwt.JwksUriEnvoyCluster = name
			}
		}

		sc.EndUserAuthnSpec = spec
		if len(authSpecs) > 1 {
			// TODO - validation should catch this problem earlier at config time.
			log.Warnf("Multiple EndUserAuthenticationPolicySpec found for service %q. Selecting %v",
				dest.Service, spec)
		}
	}

	return sc
}

// Mixer TCP filter config for inbound requests.
func buildTCPMixerFilterConfig(mesh *meshconfig.MeshConfig, role model.Node, instance *model.ServiceInstance) *FilterMixerConfig {
	filter := &FilterMixerConfig{
		MixerAttributes: map[string]string{
			AttrDestinationIP:  role.IPAddress,
			AttrDestinationUID: "kubernetes://" + role.ID,
		},
	}

	transport := &mccpb.TransportConfig{
		CheckCluster:  MixerCheckClusterName,
		ReportCluster: MixerReportClusterName,
	}
	if mesh.MixerCheckServer == mesh.MixerReportServer {
		transport.ReportCluster = transport.CheckCluster
	}

	v2 := &mccpb.TcpClientConfig{

		MixerAttributes: &mpb.Attributes{
			Attributes: map[string]*mpb.Attributes_AttributeValue{
				AttrDestinationIP:      {Value: &mpb.Attributes_AttributeValue_StringValue{role.IPAddress}},
				AttrDestinationUID:     {Value: &mpb.Attributes_AttributeValue_StringValue{"kubernetes://" + role.ID}},
				AttrDestinationService: {Value: &mpb.Attributes_AttributeValue_StringValue{instance.Service.Hostname}},
			},
		},
		Transport: transport,
	}

	v2.DisableCheckCalls = mesh.DisablePolicyChecks

	if v2JSONMap, err := model.ToJSONMap(v2); err != nil {
		log.Warnf("Could not encode v2 TCP mixerclient filter for node %q: %v", role, err)
	} else {
		filter.V2 = v2JSONMap

	}
	return filter
}

const (
	// OutboundJWTURIClusterPrefix is the prefix for jwt_uri service
	// clusters external to the proxy instance
	OutboundJWTURIClusterPrefix = "jwt."
)

// buildJWKSURIClusterNameAndAddress builds the internal envoy cluster
// name and DNS address from the jwks_uri. The cluster name is used by
// the JWT auth filter to fetch public keys. The cluster name and
// address are used to build an envoy cluster that corresponds to the
// jwks_uri server.
func buildJWKSURIClusterNameAndAddress(raw string) (string, string, bool, error) {
	var useSSL bool

	u, err := url.Parse(raw)
	if err != nil {
		return "", "", useSSL, err
	}

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"

		} else {
			port = "80"
		}
	}
	address := host + ":" + port
	name := host + "|" + port

	if u.Scheme == "https" {
		useSSL = true
	}

	return truncateClusterName(OutboundJWTURIClusterPrefix + name), address, useSSL, nil
}

// buildMixerAuthFilterClusters builds the necessary clusters for the
// JWT auth filter to fetch public keys from the specified jwks_uri.
func buildMixerAuthFilterClusters(config model.IstioConfigStore, mesh *meshconfig.MeshConfig, instances []*model.ServiceInstance) Clusters {
	type authCluster struct {
		name   string
		useSSL bool
	}
	authClusters := map[string]authCluster{}
	for _, instance := range instances {
		for _, policy := range config.EndUserAuthenticationPolicySpecByDestination(instance) {
			for _, jwt := range policy.Spec.(*mccpb.EndUserAuthenticationPolicySpec).Jwts {
				if name, address, ssl, err := buildJWKSURIClusterNameAndAddress(jwt.JwksUri); err != nil {
					log.Warnf("Could not build envoy cluster and address from jwks_uri %q: %v",
						jwt.JwksUri, err)
				} else {
					authClusters[address] = authCluster{name, ssl}
				}
			}
		}
	}

	var clusters Clusters
	for address, auth := range authClusters {
		cluster := buildCluster(address, auth.name, mesh.ConnectTimeout)
		cluster.CircuitBreaker = &CircuitBreaker{
			Default: DefaultCBPriority{
				MaxPendingRequests: 10000,
				MaxRequests:        10000,
			},
		}
		if auth.useSSL {
			cluster.SSLContext = &SSLContextExternal{}
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}
