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

package v2

import (
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/golang/protobuf/ptypes/wrappers"

	networkingapi "istio.io/api/networking/v1alpha3"

	"istio.io/istio/pilot/pkg/model"
	networking "istio.io/istio/pilot/pkg/networking/core/v1alpha3"
	"istio.io/istio/pilot/pkg/networking/core/v1alpha3/loadbalancer"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/istio/pilot/pkg/serviceregistry"
	"istio.io/istio/pilot/pkg/serviceregistry/aggregate"
	"istio.io/istio/pkg/config/host"
	"istio.io/istio/pkg/config/labels"
	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/config/schema/resource"
)

// EDS returns the list of endpoints (IP:port and in future labels) associated with a real
// service or a subset of a service, selected using labels.
//
// The source of info is a list of service registries.
//
// Primary event is an endpoint creation/deletion. Once the event is fired, EDS needs to
// find the list of services associated with the endpoint.
//
// In case of k8s, Endpoints event is fired when the endpoints are added to service - typically
// after readiness check. At that point we have the 'real' Service. The Endpoint includes a list
// of port numbers and names.
//
// For the subset case, the Pod referenced in the Endpoint must be looked up, and pod checked
// for labels.
//
// In addition, ExternalEndpoint includes IPs and labels directly and can be directly processed.
//
// TODO: for selector-less services (mesh expansion), skip pod processing
// TODO: optimize the code path for ExternalEndpoint, no additional processing needed
// TODO: if a service doesn't have split traffic - we can also skip pod and label processing
// TODO: efficient label processing. In alpha3, the destination policies are set per service, so
// we may only need to search in a small list.

var (
	// Tracks connections, increment on each new connection.
	connectionNumber = int64(0)
)

// TODO: add prom metrics !

// buildEnvoyLbEndpoint packs the endpoint based on istio info.
func buildEnvoyLbEndpoint(e *model.IstioEndpoint, push *model.PushContext) *endpoint.LbEndpoint {
	addr := util.BuildAddress(e.Address, e.EndpointPort)

	epWeight := e.LbWeight
	if epWeight == 0 {
		epWeight = 1
	}
	ep := &endpoint.LbEndpoint{
		LoadBalancingWeight: &wrappers.UInt32Value{
			Value: epWeight,
		},
		HostIdentifier: &endpoint.LbEndpoint_Endpoint{
			Endpoint: &endpoint.Endpoint{
				Address: addr,
			},
		},
	}

	// Istio telemetry depends on the metadata value being set for endpoints in the mesh.
	// Istio endpoint level tls transport socket configuration depends on this logic
	// Do not remove
	ep.Metadata = util.BuildLbEndpointMetadata(e.UID, e.Network, e.TLSMode, push)

	return ep
}

// updateServiceShards will list the endpoints and create the shards.
// This is used to reconcile and to support non-k8s registries (until they migrate).
// Note that aggregated list is expensive (for large numbers) - we want to replace
// it with a model where DiscoveryServer keeps track of all endpoint registries
// directly, and calls them one by one.
func (s *DiscoveryServer) updateServiceShards(push *model.PushContext) error {
	var registries []serviceregistry.Instance
	var nonK8sRegistries []serviceregistry.Instance
	if agg, ok := s.Env.ServiceDiscovery.(*aggregate.Controller); ok {
		registries = agg.GetRegistries()
	} else {
		registries = []serviceregistry.Instance{
			serviceregistry.Simple{
				ServiceDiscovery: s.Env.ServiceDiscovery,
			},
		}
	}

	for _, registry := range registries {
		if registry.Provider() != serviceregistry.Kubernetes {
			nonK8sRegistries = append(nonK8sRegistries, registry)
		}
	}

	// Each registry acts as a shard - we don't want to combine them because some
	// may individually update their endpoints incrementally
	for _, svc := range push.Services(nil) {
		for _, registry := range nonK8sRegistries {
			// skip the service in case this svc does not belong to the registry.
			if svc.Attributes.ServiceRegistry != string(registry.Provider()) {
				continue
			}
			endpoints := make([]*model.IstioEndpoint, 0)
			for _, port := range svc.Ports {
				if port.Protocol == protocol.UDP {
					continue
				}

				// This loses track of grouping (shards)
				instances, err := registry.InstancesByPort(svc, port.Port, labels.Collection{})
				if err != nil {
					return err
				}

				for _, inst := range instances {
					endpoints = append(endpoints, inst.Endpoint)
				}
			}

			// TODO(nmittler): Should we get the cluster from the endpoints instead? May require organizing endpoints by cluster first.
			s.edsUpdate(registry.Cluster(), string(svc.Hostname), svc.Attributes.Namespace, endpoints, true)
		}
	}

	return nil
}

// SvcUpdate is a callback from service discovery when service info changes.
func (s *DiscoveryServer) SvcUpdate(cluster, hostname string, namespace string, event model.Event) {
	// When a service deleted, we should cleanup the endpoint shards and also remove keys from EndpointShardsByService to
	// prevent memory leaks.
	if event == model.EventDelete {
		inboundServiceDeletes.Increment()
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.deleteService(cluster, hostname, namespace)
	} else {
		inboundServiceUpdates.Increment()
	}
}

// Update clusters for an incremental EDS push, and initiate the push.
// Only clusters that changed are updated/pushed.
func (s *DiscoveryServer) edsIncremental(version string, req *model.PushRequest) {
	adsLog.Infof("XDS:EDSInc Pushing:%s Services:%v ConnectedEndpoints:%d",
		version, req.EdsUpdates, s.adsClientCount())
	s.startPush(req)
}

// EDSUpdate computes destination address membership across all clusters and networks.
// This is the main method implementing EDS.
// It replaces InstancesByPort in model - instead of iterating over all endpoints it uses
// the hostname-keyed map. And it avoids the conversion from Endpoint to ServiceEntry to envoy
// on each step: instead the conversion happens once, when an endpoint is first discovered.
func (s *DiscoveryServer) EDSUpdate(clusterID, serviceName string, namespace string,
	istioEndpoints []*model.IstioEndpoint) error {
	inboundEDSUpdates.Increment()
	s.edsUpdate(clusterID, serviceName, namespace, istioEndpoints, false)
	return nil
}

// edsUpdate updates edsUpdates by clusterID, serviceName, IstioEndpoints,
// and requests a full/eds push.
func (s *DiscoveryServer) edsUpdate(clusterID, serviceName string, namespace string,
	istioEndpoints []*model.IstioEndpoint, internal bool) {
	// edsShardUpdate replaces a subset (shard) of endpoints, as result of an incremental
	// update. The endpoint updates may be grouped by K8S clusters, other service registries
	// or by deployment. Multiple updates are debounced, to avoid too frequent pushes.
	// After debounce, the services are merged and pushed.
	s.mutex.Lock()
	defer s.mutex.Unlock()
	requireFull := false

	// Should delete the service EndpointShards when endpoints become zero to prevent memory leak,
	// but we should not do not delete the keys from EndpointShardsByService map - that will trigger
	// unnecessary full push which can become a real problem if a pod is in crashloop and thus endpoints
	// flip flopping between 1 and 0.
	if len(istioEndpoints) == 0 {
		if s.EndpointShardsByService[serviceName][namespace] != nil {
			s.deleteEndpointShards(clusterID, serviceName, namespace)
			adsLog.Infof("Incremental push, service %s has no endpoints", serviceName)
			s.ConfigUpdate(&model.PushRequest{
				Full:              false,
				NamespacesUpdated: map[string]struct{}{namespace: {}},
				EdsUpdates:        map[string]struct{}{serviceName: {}},
				Reason:            []model.TriggerReason{model.EndpointUpdate},
			})
		}
		return
	}

	// Update the data structures for the service.
	// 1. Find the 'per service' data
	if _, f := s.EndpointShardsByService[serviceName]; !f {
		s.EndpointShardsByService[serviceName] = map[string]*EndpointShards{}
	}
	ep, f := s.EndpointShardsByService[serviceName][namespace]
	if !f {
		// This endpoint is for a service that was not previously loaded.
		// Return an error to force a full sync, which will also cause the
		// EndpointsShardsByService to be initialized with all services.
		ep = &EndpointShards{
			Shards:          map[string][]*model.IstioEndpoint{},
			ServiceAccounts: map[string]bool{},
		}
		s.EndpointShardsByService[serviceName][namespace] = ep
		if !internal {
			adsLog.Infof("Full push, new service %s", serviceName)
			requireFull = true
		}
	}

	// 2. Update data for the specific cluster. Each cluster gets independent
	// updates containing the full list of endpoints for the service in that cluster.
	serviceAccounts := map[string]bool{}
	for _, e := range istioEndpoints {
		if e.ServiceAccount != "" {
			serviceAccounts[e.ServiceAccount] = true
		}
	}

	if !reflect.DeepEqual(serviceAccounts, ep.ServiceAccounts) {
		adsLog.Debugf("Updating service accounts now, svc %v, before service account %v, after %v",
			serviceName, ep.ServiceAccounts, serviceAccounts)
		if !internal {
			requireFull = true
			adsLog.Infof("Full push, service accounts changed, %v", serviceName)
		}
	}

	ep.mutex.Lock()
	ep.Shards[clusterID] = istioEndpoints
	ep.ServiceAccounts = serviceAccounts
	ep.mutex.Unlock()

	// for internal update: this called by DiscoveryServer.Push --> updateServiceShards,
	// no need to trigger push here.
	// It is done in DiscoveryServer.Push --> AdsPushAll
	if !internal {
		var edsUpdates map[string]struct{}
		if !requireFull {
			edsUpdates = map[string]struct{}{serviceName: {}}
		}
		s.ConfigUpdate(&model.PushRequest{
			Full:               requireFull,
			NamespacesUpdated:  map[string]struct{}{namespace: {}},
			ConfigTypesUpdated: map[resource.GroupVersionKind]struct{}{collections.IstioNetworkingV1Alpha3Serviceentries.Resource().GroupVersionKind(): {}},
			EdsUpdates:         edsUpdates,
			Reason:             []model.TriggerReason{model.EndpointUpdate},
		})
	}
}

// deleteEndpointShards deletes matching endpoint shards from EndpointShardsByService map. This is called when
// endpoints are deleted.
func (s *DiscoveryServer) deleteEndpointShards(cluster, serviceName, namespace string) {
	if s.EndpointShardsByService[serviceName][namespace] != nil {
		s.EndpointShardsByService[serviceName][namespace].mutex.Lock()
		delete(s.EndpointShardsByService[serviceName][namespace].Shards, cluster)
		s.EndpointShardsByService[serviceName][namespace].mutex.Unlock()
	}
}

// deleteService deletes all service related references from EndpointShardsByService. This is called
// when a service is deleted.
func (s *DiscoveryServer) deleteService(cluster, serviceName, namespace string) {
	if s.EndpointShardsByService[serviceName][namespace] != nil {
		s.EndpointShardsByService[serviceName][namespace].mutex.Lock()
		delete(s.EndpointShardsByService[serviceName][namespace].Shards, cluster)
		svcShards := len(s.EndpointShardsByService[serviceName][namespace].Shards)
		s.EndpointShardsByService[serviceName][namespace].mutex.Unlock()
		if svcShards == 0 {
			delete(s.EndpointShardsByService[serviceName], namespace)
		}
		if len(s.EndpointShardsByService[serviceName]) == 0 {
			delete(s.EndpointShardsByService, serviceName)
		}
	}
}

func connectionID(node string) string {
	id := atomic.AddInt64(&connectionNumber, 1)
	return node + "-" + strconv.FormatInt(id, 10)
}

// loadAssignmentsForClusterIsolated return the endpoints for a proxy in an isolated namespace
// Initial implementation is computing the endpoints on the flight - caching will be added as needed, based on
// perf tests. The logic to compute is based on the current UpdateClusterInc
func (s *DiscoveryServer) loadAssignmentsForClusterIsolated(proxy *model.Proxy, push *model.PushContext,
	clusterName string) *xdsapi.ClusterLoadAssignment {
	_, subsetName, hostname, port := model.ParseSubsetKey(clusterName)

	// TODO: BUG. this code is incorrect if 1.1 isolation is used. With destination rule scoping
	// (public/private) as well as sidecar scopes allowing import of
	// specific destination rules, the destination rule for a given
	// namespace should be determined based on the sidecar scope or the
	// proxy's config namespace. As such, this code searches through all
	// destination rules, public and private and returns a completely
	// arbitrary destination rule's subset labels!
	subsetLabels := push.SubsetToLabels(proxy, subsetName, hostname)

	push.Mutex.Lock()
	svc := proxy.SidecarScope.ServiceForHostname(hostname, push.ServiceByHostnameAndNamespace)
	push.Mutex.Unlock()
	if svc == nil {
		// Shouldn't happen here
		adsLog.Warnf("can not find the service for cluster %s", clusterName)
		return nil
	}

	// Service resolution type might have changed and Cluster may be still in the EDS cluster list of "XdsConnection.Clusters".
	// This can happen if a ServiceEntry's resolution is changed from STATIC to DNS which changes the Envoy cluster type from
	// EDS to STRICT_DNS. When pushEds is called before Envoy sends the updated cluster list via Endpoint request which in turn
	// will update "XdsConnection.Clusters", we might accidentally send EDS updates for STRICT_DNS cluster. This check gaurds
	// against such behavior and returns nil. When the updated cluster warms up in Envoy, it would update with new endpoints
	// automatically.
	// Gateways use EDS for Passthrough cluster. So we should allow Passthrough here.
	if svc.Resolution == model.DNSLB {
		adsLog.Infof("XdsConnection has %s in its eds clusters but its resolution now is updated to %v, skipping it.", clusterName, svc.Resolution)
		return nil
	}

	svcPort, f := svc.Ports.GetByPort(port)
	if !f {
		// Shouldn't happen here
		adsLog.Warnf("can not find the service port %d for cluster %s", port, clusterName)
		return nil
	}

	// The service was never updated - do the full update
	s.mutex.RLock()
	se, f := s.EndpointShardsByService[string(hostname)][svc.Attributes.Namespace]
	s.mutex.RUnlock()
	if !f {
		// Shouldn't happen here
		adsLog.Warnf("can not find the endpointShards for cluster %s", clusterName)
		return nil
	}

	locEps := buildLocalityLbEndpointsFromShards(se, svcPort, subsetLabels, clusterName, push)

	return &xdsapi.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints:   locEps,
	}
}

func (s *DiscoveryServer) generateEndpoints(
	clusterName string, proxy *model.Proxy, push *model.PushContext, edsUpdatedServices map[string]struct{},
) *xdsapi.ClusterLoadAssignment {
	_, _, hostname, _ := model.ParseSubsetKey(clusterName)
	if edsUpdatedServices != nil {
		if _, ok := edsUpdatedServices[string(hostname)]; !ok {
			// Cluster was not updated, skip recomputing. This happens when we get an incremental update for a
			// specific Hostname. On connect or for full push edsUpdatedServices will be empty.
			return nil
		}
	}

	l := s.loadAssignmentsForClusterIsolated(proxy, push, clusterName)
	if l == nil {
		return nil
	}

	// If networks are set (by default they aren't) apply the Split Horizon
	// EDS filter on the endpoints
	if push.Networks != nil && len(push.Networks.Networks) > 0 {
		endpoints := EndpointsByNetworkFilter(push, proxy.Metadata.Network, l.Endpoints)
		filteredCLA := &xdsapi.ClusterLoadAssignment{
			ClusterName: l.ClusterName,
			Endpoints:   endpoints,
			Policy:      l.Policy,
		}
		l = filteredCLA
	}

	// If locality aware routing is enabled, prioritize endpoints or set their lb weight.
	// Failover should only be enabled when there is an outlier detection, otherwise Envoy
	// will never detect the hosts are unhealthy and redirect traffic.
	enableFailover, lb := getOutlierDetectionAndLoadBalancerSettings(push, proxy, clusterName)
	lbSetting := loadbalancer.GetLocalityLbSetting(push.Mesh.GetLocalityLbSetting(), lb.GetLocalityLbSetting())
	if lbSetting != nil {
		// Make a shallow copy of the cla as we are mutating the endpoints with priorities/weights relative to the calling proxy
		clonedCLA := util.CloneClusterLoadAssignment(l)
		l = &clonedCLA
		loadbalancer.ApplyLocalityLBSetting(proxy.Locality, l, lbSetting, enableFailover)
	}
	return l
}

// pushEds is pushing EDS updates for a single connection. Called the first time
// a client connects, for incremental updates and for full periodic updates.
func (s *DiscoveryServer) pushEds(push *model.PushContext, con *XdsConnection, version string, edsUpdatedServices map[string]struct{}) error {
	pushStart := time.Now()
	loadAssignments := make([]*xdsapi.ClusterLoadAssignment, 0)
	endpoints := 0
	empty := 0

	// All clusters that this endpoint is watching. For 1.0 - it's typically all clusters in the mesh.
	// For 1.1+Sidecar - it's the small set of explicitly imported clusters, using the isolated DestinationRules
	for _, clusterName := range con.Clusters {

		l := s.generateEndpoints(clusterName, con.node, push, edsUpdatedServices)
		if l == nil {
			continue
		}

		for _, e := range l.Endpoints {
			endpoints += len(e.LbEndpoints)
		}

		if len(l.Endpoints) == 0 {
			empty++
		}
		loadAssignments = append(loadAssignments, l)
	}

	response := endpointDiscoveryResponse(loadAssignments, version, push.Version)
	err := con.send(response)
	edsPushTime.Record(time.Since(pushStart).Seconds())
	if err != nil {
		adsLog.Warnf("EDS: Send failure %s: %v", con.ConID, err)
		recordSendError(edsSendErrPushes, err)
		return err
	}
	edsPushes.Increment()

	if edsUpdatedServices == nil {
		adsLog.Infof("EDS: PUSH for node:%s clusters:%d endpoints:%d empty:%v",
			con.node.ID, len(con.Clusters), endpoints, empty)
	} else {
		adsLog.Debugf("EDS: PUSH INC for node:%s clusters:%d endpoints:%d empty:%v",
			con.node.ID, len(con.Clusters), endpoints, empty)
	}
	return nil
}

// getDestinationRule gets the DestinationRule for a given hostname. As an optimization, this also gets the service port,
// which is needed to access the traffic policy from the destination rule.
func getDestinationRule(push *model.PushContext, proxy *model.Proxy, hostname host.Name, clusterPort int) (*networkingapi.DestinationRule, *model.Port) {
	for _, service := range push.Services(proxy) {
		if service.Hostname == hostname {
			cfg := push.DestinationRule(proxy, service)
			if cfg == nil {
				continue
			}
			for _, p := range service.Ports {
				if p.Port == clusterPort {
					return cfg.Spec.(*networkingapi.DestinationRule), p
				}
			}
		}
	}
	return nil, nil
}

func getOutlierDetectionAndLoadBalancerSettings(push *model.PushContext, proxy *model.Proxy, clusterName string) (bool, *networkingapi.LoadBalancerSettings) {
	_, subsetName, hostname, portNumber := model.ParseSubsetKey(clusterName)
	var outlierDetectionEnabled = false
	var lbSettings *networkingapi.LoadBalancerSettings
	destinationRule, port := getDestinationRule(push, proxy, hostname, portNumber)
	if destinationRule == nil || port == nil {
		return false, nil
	}

	_, outlierDetection, loadBalancerSettings, _ := networking.SelectTrafficPolicyComponents(destinationRule.TrafficPolicy, port)
	lbSettings = loadBalancerSettings
	if outlierDetection != nil {
		outlierDetectionEnabled = true
	}

	for _, subset := range destinationRule.Subsets {
		if subset.Name == subsetName {
			_, outlierDetection, loadBalancerSettings, _ := networking.SelectTrafficPolicyComponents(subset.TrafficPolicy, port)
			lbSettings = loadBalancerSettings
			if outlierDetection != nil {
				outlierDetectionEnabled = true
			}
		}
	}
	return outlierDetectionEnabled, lbSettings
}

func endpointDiscoveryResponse(loadAssignments []*xdsapi.ClusterLoadAssignment, version string, noncePrefix string) *xdsapi.DiscoveryResponse {
	out := &xdsapi.DiscoveryResponse{
		TypeUrl: EndpointType,
		// Pilot does not really care for versioning. It always supplies what's currently
		// available to it, irrespective of whether Envoy chooses to accept or reject EDS
		// responses. Pilot believes in eventual consistency and that at some point, Envoy
		// will begin seeing results it deems to be good.
		VersionInfo: version,
		Nonce:       nonce(noncePrefix),
	}
	for _, loadAssignment := range loadAssignments {
		resource := util.MessageToAny(loadAssignment)
		out.Resources = append(out.Resources, resource)
	}

	return out
}

// build LocalityLbEndpoints for a cluster from existing EndpointShards.
func buildLocalityLbEndpointsFromShards(
	shards *EndpointShards,
	svcPort *model.Port,
	epLabels labels.Collection,
	clusterName string,
	push *model.PushContext) []*endpoint.LocalityLbEndpoints {
	localityEpMap := make(map[string]*endpoint.LocalityLbEndpoints)

	shards.mutex.Lock()
	// The shards are updated independently, now need to filter and merge
	// for this cluster
	for _, endpoints := range shards.Shards {
		for _, ep := range endpoints {
			if svcPort.Name != ep.ServicePortName {
				continue
			}
			// Port labels
			if !epLabels.HasSubsetOf(ep.Labels) {
				continue
			}

			locLbEps, found := localityEpMap[ep.Locality.Label]
			if !found {
				locLbEps = &endpoint.LocalityLbEndpoints{
					Locality:    util.ConvertLocality(ep.Locality.Label),
					LbEndpoints: make([]*endpoint.LbEndpoint, 0, len(endpoints)),
				}
				localityEpMap[ep.Locality.Label] = locLbEps
			}
			if ep.EnvoyEndpoint == nil {
				ep.EnvoyEndpoint = buildEnvoyLbEndpoint(ep, push)
			}
			locLbEps.LbEndpoints = append(locLbEps.LbEndpoints, ep.EnvoyEndpoint)

		}
	}
	shards.mutex.Unlock()

	locEps := make([]*endpoint.LocalityLbEndpoints, 0, len(localityEpMap))
	for _, locLbEps := range localityEpMap {
		var weight uint32
		for _, ep := range locLbEps.LbEndpoints {
			weight += ep.LoadBalancingWeight.GetValue()
		}
		locLbEps.LoadBalancingWeight = &wrappers.UInt32Value{
			Value: weight,
		}
		locEps = append(locEps, locLbEps)
	}

	if len(locEps) == 0 {
		push.AddMetric(model.ProxyStatusClusterNoInstances, clusterName, nil, "")
	}

	updateEdsStats(locEps, clusterName)

	return locEps
}

func updateEdsStats(locEps []*endpoint.LocalityLbEndpoints, cluster string) {
	edsInstances.With(clusterTag.Value(cluster)).Record(float64(len(locEps)))
	epc := 0
	for _, locLbEps := range locEps {
		epc += len(locLbEps.GetLbEndpoints())
	}
	edsAllLocalityEndpoints.With(clusterTag.Value(cluster)).Record(float64(epc))
}
