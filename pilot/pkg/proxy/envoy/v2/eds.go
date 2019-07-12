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
	"sync"
	"sync/atomic"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/gogo/protobuf/types"

	networkingapi "istio.io/api/networking/v1alpha3"

	"istio.io/istio/pilot/pkg/features"
	networking "istio.io/istio/pilot/pkg/networking/core/v1alpha3"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core/v1alpha3/loadbalancer"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/istio/pilot/pkg/serviceregistry"
	"istio.io/istio/pilot/pkg/serviceregistry/aggregate"
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
// TODO: if a service doesn't have split traffic - we can also skip pod and lable processing
// TODO: efficient label processing. In alpha3, the destination policies are set per service, so
// we may only need to search in a small list.

var (
	// Tracks connections, increment on each new connection.
	connectionNumber = int64(0)
)

// EdsCluster tracks eds-related info for monitored cluster. Used in 1.0, where cluster info is not source-dependent.
type EdsCluster struct {
	// mutex protects changes to this cluster
	mutex sync.Mutex

	// LoadAssignment has the pre-computed EDS response for this cluster. Any sidecar asking for the
	// cluster will get this response.
	LoadAssignment *xdsapi.ClusterLoadAssignment

	// EdsClients keeps track of all nodes monitoring the cluster.
	EdsClients map[string]*XdsConnection `json:"-"`
}

// TODO: add prom metrics !

// Return the load assignment with mutex. The field can be updated by another routine.
func loadAssignment(c *EdsCluster) *xdsapi.ClusterLoadAssignment {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.LoadAssignment
}

// buildEnvoyLbEndpoint packs the endpoint based on istio info.
func buildEnvoyLbEndpoint(uid string, family model.AddressFamily, address string, port uint32, network string, weight uint32) *endpoint.LbEndpoint {
	var addr core.Address
	switch family {
	case model.AddressFamilyTCP:
		addr = core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Address: address,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		}
	case model.AddressFamilyUnix:
		addr = core.Address{Address: &core.Address_Pipe{Pipe: &core.Pipe{Path: address}}}
	}

	epWeight := weight
	if epWeight == 0 {
		epWeight = 1
	}
	ep := &endpoint.LbEndpoint{
		LoadBalancingWeight: &types.UInt32Value{
			Value: epWeight,
		},
		HostIdentifier: &endpoint.LbEndpoint_Endpoint{
			Endpoint: &endpoint.Endpoint{
				Address: &addr,
			},
		},
	}

	// Istio telemetry depends on the metadata value being set for endpoints in the mesh.
	// Do not remove: mixerfilter depends on this logic.
	ep.Metadata = endpointMetadata(uid, network)

	return ep
}

func networkEndpointToEnvoyEndpoint(e *model.NetworkEndpoint) (*endpoint.LbEndpoint, error) {
	err := model.ValidateNetworkEndpointAddress(e)
	if err != nil {
		return nil, err
	}

	addr := util.GetNetworkEndpointAddress(e)

	epWeight := e.LbWeight
	if epWeight == 0 {
		epWeight = 1
	}
	ep := &endpoint.LbEndpoint{
		LoadBalancingWeight: &types.UInt32Value{
			Value: epWeight,
		},
		HostIdentifier: &endpoint.LbEndpoint_Endpoint{
			Endpoint: &endpoint.Endpoint{
				Address: &addr,
			},
		},
	}

	// Istio telemetry depends on the metadata value being set for endpoints in the mesh.
	// Do not remove: mixerfilter depends on this logic.
	ep.Metadata = endpointMetadata(e.UID, e.Network)

	return ep, nil
}

// Create an Istio filter metadata object with the UID and Network fields (if exist).
func endpointMetadata(uid string, network string) *core.Metadata {
	if uid == "" && network == "" {
		return nil
	}

	metadata := &core.Metadata{
		FilterMetadata: map[string]*types.Struct{
			util.IstioMetadataKey: {
				Fields: map[string]*types.Value{},
			},
		},
	}

	if uid != "" {
		metadata.FilterMetadata["istio"].Fields["uid"] = &types.Value{Kind: &types.Value_StringValue{StringValue: uid}}
	}

	if network != "" {
		metadata.FilterMetadata["istio"].Fields["network"] = &types.Value{Kind: &types.Value_StringValue{StringValue: network}}
	}

	return metadata
}

// updateClusterInc computes an envoy cluster assignment from the service shards.
// This happens when endpoints are updated.
// TODO: this code is specific for 1.0 / pre-isolation. With config scoping, two sidecars can get
// a cluster of same name but with different set of endpoints. See the
// explanation below for more details
func (s *DiscoveryServer) updateClusterInc(push *model.PushContext, clusterName string,
	edsCluster *EdsCluster) error {

	var hostname model.Hostname
	var clusterPort int
	var subsetName string
	_, subsetName, hostname, clusterPort = model.ParseSubsetKey(clusterName)

	// TODO: BUG. this code is incorrect if 1.1 isolation is used. With destination rule scoping
	// (public/private) as well as sidecar scopes allowing import of
	// specific destination rules, the destination rule for a given
	// namespace should be determined based on the sidecar scope or the
	// proxy's config namespace. As such, this code searches through all
	// destination rules, public and private and returns a completely
	// arbitrary destination rule's subset labels!
	labels := push.SubsetToLabels(nil, subsetName, hostname)

	push.Mutex.Lock()
	ports, f := push.ServicePortByHostname[hostname]
	push.Mutex.Unlock()
	if !f {
		return s.updateCluster(push, clusterName, edsCluster)
	}

	// Check that there is a matching port
	// We don't use the port though, as there could be multiple matches
	svcPort, found := ports.GetByPort(clusterPort)
	if !found {
		return s.updateCluster(push, clusterName, edsCluster)
	}

	s.mutex.RLock()
	// The service was never updated - do the full update
	se, f := s.EndpointShardsByService[string(hostname)]
	s.mutex.RUnlock()
	if !f {
		return s.updateCluster(push, clusterName, edsCluster)
	}

	locEps := buildLocalityLbEndpointsFromShards(se, svcPort, labels, clusterName, push)
	// There is a chance multiple goroutines will update the cluster at the same time.
	// This could be prevented by a lock - but because the update may be slow, it may be
	// better to accept the extra computations.
	// We still lock the access to the LoadAssignments.
	edsCluster.mutex.Lock()
	defer edsCluster.mutex.Unlock()

	edsCluster.LoadAssignment = &xdsapi.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints:   locEps,
	}
	return nil
}

// updateServiceShards will list the endpoints and create the shards.
// This is used to reconcile and to support non-k8s registries (until they migrate).
// Note that aggregated list is expensive (for large numbers) - we want to replace
// it with a model where DiscoveryServer keeps track of all endpoint registries
// directly, and calls them one by one.
func (s *DiscoveryServer) updateServiceShards(push *model.PushContext) error {

	// TODO: if ServiceDiscovery is aggregate, and all members support direct, use
	// the direct interface.
	var registries []aggregate.Registry
	var nonK8sRegistries []aggregate.Registry
	if agg, ok := s.Env.ServiceDiscovery.(*aggregate.Controller); ok {
		registries = agg.GetRegistries()
	} else {
		registries = []aggregate.Registry{
			{
				ServiceDiscovery: s.Env.ServiceDiscovery,
			},
		}
	}

	for _, registry := range registries {
		if registry.Name != serviceregistry.KubernetesRegistry {
			nonK8sRegistries = append(nonK8sRegistries, registry)
		}
	}

	// Each registry acts as a shard - we don't want to combine them because some
	// may individually update their endpoints incrementally
	for _, svc := range push.Services(nil) {
		for _, registry := range nonK8sRegistries {
			// in case this svc does not belong to the registry
			if svc, _ := registry.GetService(svc.Hostname); svc == nil {
				continue
			}

			entries := []*model.IstioEndpoint{}
			hostname := string(svc.Hostname)
			for _, port := range svc.Ports {
				if port.Protocol == model.ProtocolUDP {
					continue
				}

				// This loses track of grouping (shards)
				endpoints, err := registry.InstancesByPort(svc.Hostname, port.Port, model.LabelsCollection{})
				if err != nil {
					return err
				}

				for _, ep := range endpoints {
					entries = append(entries, &model.IstioEndpoint{
						Family:          ep.Endpoint.Family,
						Address:         ep.Endpoint.Address,
						EndpointPort:    uint32(ep.Endpoint.Port),
						ServicePortName: port.Name,
						Labels:          ep.Labels,
						UID:             ep.Endpoint.UID,
						ServiceAccount:  ep.ServiceAccount,
						Network:         ep.Endpoint.Network,
						Locality:        ep.GetLocality(),
						LbWeight:        ep.Endpoint.LbWeight,
					})
				}
			}

			s.edsUpdate(registry.ClusterID, hostname, entries, true)
		}
	}

	return nil
}

// updateCluster is called from the event (or global cache invalidation) to update
// the endpoints for the cluster.
func (s *DiscoveryServer) updateCluster(push *model.PushContext, clusterName string, edsCluster *EdsCluster) error {
	// TODO: should we lock this as well ? Once we move to event-based it may not matter.
	var locEps []endpoint.LocalityLbEndpoints
	direction, subsetName, hostname, port := model.ParseSubsetKey(clusterName)

	if direction == model.TrafficDirectionInbound ||
		direction == model.TrafficDirectionOutbound {
		labels := push.SubsetToLabels(nil, subsetName, hostname)
		instances, err := s.Env.ServiceDiscovery.InstancesByPort(hostname, port, labels)
		if err != nil {
			adsLog.Errorf("endpoints for service cluster %q returned error %v", clusterName, err)
			totalXDSInternalErrors.Increment()
			return err
		}
		if len(instances) == 0 {
			push.Add(model.ProxyStatusClusterNoInstances, clusterName, nil, "")
			adsLog.Debugf("EDS: Cluster %q (host:%s ports:%v labels:%v) has no instances", clusterName, hostname, port, labels)
		}
		edsInstances.With(clusterTag.Value(clusterName)).Record(float64(len(instances)))
		locEps = localityLbEndpointsFromInstances(instances)
	}

	for i := 0; i < len(locEps); i++ {
		var weight uint32
		for _, ep := range locEps[i].LbEndpoints {
			weight += ep.LoadBalancingWeight.GetValue()
		}
		locEps[i].LoadBalancingWeight = &types.UInt32Value{
			Value: weight,
		}
	}
	// There is a chance multiple goroutines will update the cluster at the same time.
	// This could be prevented by a lock - but because the update may be slow, it may be
	// better to accept the extra computations.
	// We still lock the access to the LoadAssignments.
	edsCluster.mutex.Lock()
	defer edsCluster.mutex.Unlock()

	// Normalize LoadBalancingWeight in range [1, 128]
	locEps = LoadBalancingWeightNormalize(locEps)

	edsCluster.LoadAssignment = &xdsapi.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints:   locEps,
	}
	return nil
}

// SvcUpdate is a callback from service discovery when service info changes.
func (s *DiscoveryServer) SvcUpdate(cluster, hostname string, ports map[string]uint32, _ map[uint32]string) {
	inboundServiceUpdates.Increment()
	pc := s.globalPushContext()
	pl := model.PortList{}
	for k, v := range ports {
		pl = append(pl, &model.Port{
			Port: int(v),
			Name: k,
		})
	}
	if cluster == string(serviceregistry.KubernetesRegistry) {
		pc.Mutex.Lock()
		pc.ServicePortByHostname[model.Hostname(hostname)] = pl
		pc.Mutex.Unlock()
	} else {
		pc.Mutex.Lock()
		ports, ok := pc.ServicePortByHostname[model.Hostname(hostname)]
		pc.Mutex.Unlock()
		if ok {
			if !reflect.DeepEqual(ports, pl) {
				adsLog.Warnf("service %s within cluster %s does not match the one in primary cluster", hostname, cluster)
			}
		} else {
			adsLog.Debugf("service %s within cluster %s occurs before primary cluster, will only take services within primary cluster", hostname, cluster)
		}
	}
}

// WorkloadUpdate is called when workload labels/annotations are updated.
func (s *DiscoveryServer) WorkloadUpdate(id string, labels map[string]string, _ map[string]string) {
	inboundWorkloadUpdates.Increment()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if labels == nil {
		// No push needed - the Endpoints object will also be triggered.
		delete(s.WorkloadsByID, id)
		return
	}
	w, f := s.WorkloadsByID[id]
	if !f {
		s.WorkloadsByID[id] = &Workload{
			Labels: labels,
		}

		fullPush := false
		adsClientsMutex.RLock()
		for _, connection := range adsClients {
			// if the workload has envoy proxy and connected to server,
			// then do a full xDS push for this proxy;
			// otherwise:
			//   case 1: the workload has no sidecar proxy, no need xDS push at all.
			//   case 2: the workload xDS connection has not been established,
			//           also no need to trigger a full push here.
			if connection.modelNode.IPAddresses[0] == id {
				fullPush = true
			}
		}
		adsClientsMutex.RUnlock()

		if fullPush {
			// First time this workload has been seen. Maybe after the first connect,
			// do a full push for this proxy in the next push epoch.
			s.proxyUpdatesMutex.Lock()
			if s.proxyUpdates == nil {
				s.proxyUpdates = make(map[string]struct{})
			}
			s.proxyUpdates[id] = struct{}{}
			s.proxyUpdatesMutex.Unlock()
		}
		return
	}
	if reflect.DeepEqual(w.Labels, labels) {
		// No label change.
		return
	}

	w.Labels = labels
	// Label changes require recomputing the config.
	// TODO: we can do a push for the affected workload only, but we need to confirm
	// no other workload can be affected. Safer option is to fallback to full push.

	adsLog.Infof("Label change, full push %s ", id)
	s.ConfigUpdate(true)
}

// EDSUpdate computes destination address membership across all clusters and networks.
// This is the main method implementing EDS.
// It replaces InstancesByPort in model - instead of iterating over all endpoints it uses
// the hostname-keyed map. And it avoids the conversion from Endpoint to ServiceEntry to envoy
// on each step: instead the conversion happens once, when an endpoint is first discovered.
func (s *DiscoveryServer) EDSUpdate(shard, serviceName string, istioEndpoints []*model.IstioEndpoint) error {
	inboundEDSUpdates.Increment()
	s.edsUpdate(shard, serviceName, istioEndpoints, false)
	return nil
}

// edsUpdate updates edsUpdates by shard, serviceName, IstioEndpoints,
// and requests a full/eds push.
func (s *DiscoveryServer) edsUpdate(shard, serviceName string,
	istioEndpoints []*model.IstioEndpoint, internal bool) {

	// edsShardUpdate replaces a subset (shard) of endpoints, as result of an incremental
	// update. The endpoint updates may be grouped by K8S clusters, other service registries
	// or by deployment. Multiple updates are debounced, to avoid too frequent pushes.
	// After debounce, the services are merged and pushed.
	s.mutex.Lock()
	defer s.mutex.Unlock()
	requireFull := false

	// To prevent memory leak.
	// Should delete the service EndpointShards, when endpoints deleted or service deleted.
	if len(istioEndpoints) == 0 {
		if s.EndpointShardsByService[serviceName] != nil {
			s.EndpointShardsByService[serviceName].mutex.Lock()
			delete(s.EndpointShardsByService[serviceName].Shards, shard)
			svcShards := len(s.EndpointShardsByService[serviceName].Shards)
			s.EndpointShardsByService[serviceName].mutex.Unlock()
			if svcShards == 0 {
				delete(s.EndpointShardsByService, serviceName)
			}
		}
		return
	}

	// Update the data structures for the service.
	// 1. Find the 'per service' data
	ep, f := s.EndpointShardsByService[serviceName]
	if !f {
		// This endpoint is for a service that was not previously loaded.
		// Return an error to force a full sync, which will also cause the
		// EndpointsShardsByService to be initialized with all services.
		ep = &EndpointShards{
			Shards:          map[string][]*model.IstioEndpoint{},
			ServiceAccounts: map[string]bool{},
		}
		s.EndpointShardsByService[serviceName] = ep
		if !internal {
			adsLog.Infof("Full push, new service %s", serviceName)
			requireFull = true
		}
	}

	// 2. Update data for the specific cluster. Each cluster gets independent
	// updates containing the full list of endpoints for the service in that cluster.
	for _, e := range istioEndpoints {
		if e.ServiceAccount != "" {
			ep.mutex.Lock()
			_, f = ep.ServiceAccounts[e.ServiceAccount]
			if !f {
				ep.ServiceAccounts[e.ServiceAccount] = true
			}
			ep.mutex.Unlock()

			if !f && !internal {
				// The entry has a service account that was not previously associated.
				// Requires a CDS push and full sync.
				adsLog.Infof("Endpoint updating service account %s %s", e.ServiceAccount, serviceName)
				requireFull = true
				break
			}
		}
	}
	ep.mutex.Lock()
	ep.Shards[shard] = istioEndpoints
	ep.mutex.Unlock()
	s.edsUpdates[serviceName] = struct{}{}

	// for internal update: this called by DiscoveryServer.Push --> updateServiceShards,
	// no need to trigger push here.
	// It is done in DiscoveryServer.Push --> AdsPushAll
	if !internal {
		if requireFull {
			s.ConfigUpdate(true)
		} else {
			s.ConfigUpdate(false)
		}
	}
}

// LocalityLbEndpointsFromInstances returns a list of Envoy v2 LocalityLbEndpoints.
// Envoy v2 Endpoints are constructed from Pilot's older data structure involving
// model.ServiceInstance objects. Envoy expects the endpoints grouped by zone, so
// a map is created - in new data structures this should be part of the model.
func localityLbEndpointsFromInstances(instances []*model.ServiceInstance) []endpoint.LocalityLbEndpoints {
	localityEpMap := make(map[string]*endpoint.LocalityLbEndpoints)
	for _, instance := range instances {
		lbEp, err := networkEndpointToEnvoyEndpoint(&instance.Endpoint)
		if err != nil {
			adsLog.Errorf("EDS: Unexpected pilot model endpoint v1 to v2 conversion: %v", err)
			totalXDSInternalErrors.Increment()
			continue
		}
		locality := instance.GetLocality()
		locLbEps, found := localityEpMap[locality]
		if !found {
			locLbEps = &endpoint.LocalityLbEndpoints{
				Locality: util.ConvertLocality(locality),
			}
			localityEpMap[locality] = locLbEps
		}
		locLbEps.LbEndpoints = append(locLbEps.LbEndpoints, *lbEp)
	}
	out := make([]endpoint.LocalityLbEndpoints, 0, len(localityEpMap))
	for _, locLbEps := range localityEpMap {
		out = append(out, *locLbEps)
	}
	return out
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
	// TODO: fail-safe, use the old implementation based on some flag.
	// Users who make sure all DestinationRules are in the right namespace and don't have override may turn it on
	// (some very-large scale customers are in this category)

	// This code is similar with the update code.
	_, subsetName, hostname, port := model.ParseSubsetKey(clusterName)

	// TODO: BUG. this code is incorrect if 1.1 isolation is used. With destination rule scoping
	// (public/private) as well as sidecar scopes allowing import of
	// specific destination rules, the destination rule for a given
	// namespace should be determined based on the sidecar scope or the
	// proxy's config namespace. As such, this code searches through all
	// destination rules, public and private and returns a completely
	// arbitrary destination rule's subset labels!
	labels := push.SubsetToLabels(proxy, subsetName, hostname)

	push.Mutex.Lock()
	portMap, f := push.ServicePortByHostname[hostname]
	push.Mutex.Unlock()
	if !f {
		// Shouldn't happen here - if it does, fallback to sending empty endpoint set
		return &xdsapi.ClusterLoadAssignment{
			ClusterName: clusterName,
		}
	}

	svcPort, f := portMap.GetByPort(port)
	if !f {
		// Shouldn't happen here - if it does, fallback to sending empty endpoint set
		return &xdsapi.ClusterLoadAssignment{
			ClusterName: clusterName,
		}
	}

	// The service was never updated - do the full update
	s.mutex.RLock()
	se, f := s.EndpointShardsByService[string(hostname)]
	s.mutex.RUnlock()
	if !f {
		// This can happen when there are no endpoints for a pod
		// In this case, push an empty set
		se = &EndpointShards{
			Shards:          map[string][]*model.IstioEndpoint{},
			ServiceAccounts: map[string]bool{},
		}
	}

	locEps := buildLocalityLbEndpointsFromShards(se, svcPort, labels, clusterName, push)

	return &xdsapi.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints:   locEps,
	}
}

// pushEds is pushing EDS updates for a single connection. Called the first time
// a client connects, for incremental updates and for full periodic updates.
func (s *DiscoveryServer) pushEds(push *model.PushContext, con *XdsConnection, version string, edsUpdatedServices map[string]struct{}) error {
	loadAssignments := []*xdsapi.ClusterLoadAssignment{}
	endpoints := 0
	empty := []string{}

	// All clusters that this endpoint is watching. For 1.0 - it's typically all clusters in the mesh.
	// For 1.1+Sidecar - it's the small set of explicitly imported clusters, using the isolated DestinationRules
	for _, clusterName := range con.Clusters {

		_, _, hostname, _ := model.ParseSubsetKey(clusterName)
		if edsUpdatedServices != nil {
			if _, ok := edsUpdatedServices[string(hostname)]; !ok {
				// Cluster was not updated, skip recomputing. This happens when we get an incremental update for a
				// specific Hostname. On connect or for full push edsUpdatedServices will be empty.
				continue
			}
		}

		l := s.loadAssignmentsForClusterIsolated(con.modelNode, push, clusterName)

		// If networks are set (by default they aren't) apply the Split Horizon
		// EDS filter on the endpoints
		if s.Env.MeshNetworks != nil && len(s.Env.MeshNetworks.Networks) > 0 {
			endpoints := EndpointsByNetworkFilter(l.Endpoints, con, s.Env)
			endpoints = LoadBalancingWeightNormalize(endpoints)
			filteredCLA := &xdsapi.ClusterLoadAssignment{
				ClusterName: l.ClusterName,
				Endpoints:   endpoints,
				Policy:      l.Policy,
			}
			l = filteredCLA
		}

		// If location prioritized load balancing is enabled, prioritize endpoints.
		if features.EnableLocalityLoadBalancing() {
			// Make a shallow copy of the cla as we are mutating the endpoints with priorities/weights relative to the calling proxy
			clonedCLA := util.CloneClusterLoadAssignment(l)
			l = &clonedCLA

			// Failover should only be enabled when there is an outlier detection, otherwise Envoy
			// will never detect the hosts are unhealthy and redirect traffic.
			enableFailover := hasOutlierDetection(push, con.modelNode, clusterName)
			loadbalancer.ApplyLocalityLBSetting(con.modelNode.Locality, l, s.Env.Mesh.LocalityLbSetting, enableFailover)
		}

		endpoints += len(l.Endpoints)
		if len(l.Endpoints) == 0 {
			empty = append(empty, clusterName)
		}
		loadAssignments = append(loadAssignments, l)
	}

	response := endpointDiscoveryResponse(loadAssignments, version)
	err := con.send(response)
	if err != nil {
		adsLog.Warnf("EDS: Send failure %s: %v", con.ConID, err)
		edsSendErrPushes.Increment()
		return err
	}
	edsPushes.Increment()

	if edsUpdatedServices == nil {
		adsLog.Infof("EDS: PUSH for node:%s clusters:%d endpoints:%d empty:%v",
			con.modelNode.ID, len(con.Clusters), endpoints, empty)
	} else {
		adsLog.Infof("EDS: PUSH INC for node:%s clusters:%d endpoints:%d empty:%v",
			con.modelNode.ID, len(con.Clusters), endpoints, empty)
	}
	return nil
}

// getDestinationRule gets the DestinationRule for a given hostname. As an optimization, this also gets the service port,
// which is needed to access the traffic policy from the destination rule.
func getDestinationRule(push *model.PushContext, proxy *model.Proxy, hostname model.Hostname, clusterPort int) (*networkingapi.DestinationRule, *model.Port) {
	for _, service := range push.Services(proxy) {
		if service.Hostname == hostname {
			config := push.DestinationRule(proxy, service)
			if config == nil {
				continue
			}
			for _, p := range service.Ports {
				if p.Port == clusterPort {
					return config.Spec.(*networkingapi.DestinationRule), p
				}
			}
		}
	}
	return nil, nil
}

func hasOutlierDetection(push *model.PushContext, proxy *model.Proxy, clusterName string) bool {
	_, subsetName, hostname, portNumber := model.ParseSubsetKey(clusterName)

	destinationRule, port := getDestinationRule(push, proxy, hostname, portNumber)
	if destinationRule == nil || port == nil {
		return false
	}

	_, outlierDetection, _, _ := networking.SelectTrafficPolicyComponents(destinationRule.TrafficPolicy, port)
	if outlierDetection != nil {
		return true
	}

	for _, subset := range destinationRule.Subsets {
		if subset.Name == subsetName {
			_, outlierDetection, _, _ := networking.SelectTrafficPolicyComponents(subset.TrafficPolicy, port)
			if outlierDetection != nil {
				return true
			}
		}
	}
	return false
}

func endpointDiscoveryResponse(loadAssignments []*xdsapi.ClusterLoadAssignment, version string) *xdsapi.DiscoveryResponse {
	out := &xdsapi.DiscoveryResponse{
		TypeUrl: EndpointType,
		// Pilot does not really care for versioning. It always supplies what's currently
		// available to it, irrespective of whether Envoy chooses to accept or reject EDS
		// responses. Pilot believes in eventual consistency and that at some point, Envoy
		// will begin seeing results it deems to be good.
		VersionInfo: version,
		Nonce:       nonce(),
	}
	for _, loadAssignment := range loadAssignments {
		resource, _ := types.MarshalAny(loadAssignment)
		out.Resources = append(out.Resources, *resource)
	}

	return out
}

// build LocalityLbEndpoints for a cluster from existing EndpointShards.
func buildLocalityLbEndpointsFromShards(
	shards *EndpointShards,
	svcPort *model.Port,
	labels model.LabelsCollection,
	clusterName string,
	push *model.PushContext) []endpoint.LocalityLbEndpoints {
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
			if !labels.HasSubsetOf(model.Labels(ep.Labels)) {
				continue
			}

			locLbEps, found := localityEpMap[ep.Locality]
			if !found {
				locLbEps = &endpoint.LocalityLbEndpoints{
					Locality: util.ConvertLocality(ep.Locality),
				}
				localityEpMap[ep.Locality] = locLbEps
			}
			if ep.EnvoyEndpoint == nil {
				ep.EnvoyEndpoint = buildEnvoyLbEndpoint(ep.UID, ep.Family, ep.Address, ep.EndpointPort, ep.Network, ep.LbWeight)
			}
			locLbEps.LbEndpoints = append(locLbEps.LbEndpoints, *ep.EnvoyEndpoint)

		}
	}
	shards.mutex.Unlock()

	locEps := make([]endpoint.LocalityLbEndpoints, 0, len(localityEpMap))
	for _, locLbEps := range localityEpMap {
		var weight uint32
		for _, ep := range locLbEps.LbEndpoints {
			weight += ep.LoadBalancingWeight.GetValue()
		}
		locLbEps.LoadBalancingWeight = &types.UInt32Value{
			Value: weight,
		}
		locEps = append(locEps, *locLbEps)
	}
	// Normalize LoadBalancingWeight in range [1, 128]
	locEps = LoadBalancingWeightNormalize(locEps)

	if len(locEps) == 0 {
		push.Add(model.ProxyStatusClusterNoInstances, clusterName, nil, "")
	}
	edsInstances.With(clusterTag.Value(clusterName)).Record(float64(len(locEps)))

	return locEps
}
