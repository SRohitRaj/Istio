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

package controller

import (
	"net/netip"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/proto"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"

	"istio.io/api/security/v1beta1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry/kube"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/constants"
	"istio.io/istio/pkg/config/labels"
	"istio.io/istio/pkg/config/schema/gvk"
	"istio.io/istio/pkg/config/schema/kind"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kclient"
	kubelabels "istio.io/istio/pkg/kube/labels"
	"istio.io/istio/pkg/spiffe"
	"istio.io/istio/pkg/util/sets"
	"istio.io/istio/pkg/workloadapi"
	"istio.io/istio/pkg/workloadapi/security"
)

// AmbientIndex maintains an index of ambient WorkloadInfo objects by various keys.
// These are intentionally pre-computed based on events such that lookups are efficient.
type AmbientIndex struct {
	mu sync.RWMutex
	// byService indexes by network/Service (virtual) *IP address*. A given Service may have multiple IPs, thus
	// multiple entries in the map. A given IP can have many workloads associated.
	byService map[networkAddress][]*model.WorkloadInfo
	// byPod indexes by network/podIP address.
	byPod map[string]*model.WorkloadInfo
	// services are indexed by the network/clusterIP
	services map[networkAddress]*workloadapi.Service

	// Map of ServiceAccount -> IP
	// TODO: currently, this is derived from pods. To be agnostic to the implementation,
	// we should actually be looking at Gateway.status.addresses.
	// This may be an external address (possibly even a DNS name we need to resolve), an arbitrary IP,
	// or a reference to a service.
	// If its a reference to a Service then we can find the underlying pods in that service, as an optimization.
	waypoints map[model.WaypointScope]*workloadapi.GatewayAddress

	// serviceVipIndex maintains an index of VIP -> Service
	serviceVipIndex *kclient.Index[string, *v1.Service]
}

// Lookup finds a given IP address.
func (a *AmbientIndex) Lookup(key string) []*model.AddressInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()
	// First look at pod...
	if p, f := a.byPod[key]; f {
		addr := &workloadapi.Address{
			Type: &workloadapi.Address_Workload{
				Workload: p.Workload,
			},
		}
		return []*model.AddressInfo{{Address: addr}}
	}

	// Fallback to service. Note: these IP ranges should be non-overlapping
	// cannot distunagish between Service lookup and Workloads for a Service lookup
	// since the same newtork/vip key is used.  Both Service and Workloads are
	// returned in this case
	// TODO any service to exclude in the future? for example with affinity,
	// per-node, etc - where simply load balancing will break?

	res := []*model.AddressInfo{}
	network, vip, found := strings.Cut(key, "/")
	if !found {
		log.Warnf("key(%v) did not contain the expected \"/\" character", key)
		return res
	}
	networkAddr := networkAddress{network: network, vip: vip}
	for _, wl := range a.byService[networkAddr] {
		addr := &workloadapi.Address{
			Type: &workloadapi.Address_Workload{
				Workload: wl.Workload,
			},
		}
		res = append(res, &model.AddressInfo{Address: addr})
	}
	if s, exists := a.services[networkAddr]; exists {
		addr := &workloadapi.Address{
			Type: &workloadapi.Address_Service{
				Service: s,
			},
		}
		res = append(res, &model.AddressInfo{Address: addr})
	}

	return res
}

func (a *AmbientIndex) dropWorkloadFromService(svcAddress networkAddress, workloadAddress string) {
	wls := a.byService[svcAddress]
	// TODO: this is inefficient, but basically we are trying to update a keyed element in a list
	// Probably we want a Map? But the list is nice for fast lookups
	filtered := make([]*model.WorkloadInfo, 0, len(wls))
	for _, inc := range wls {
		if inc.ResourceName() != workloadAddress {
			filtered = append(filtered, inc)
		}
	}
	a.byService[svcAddress] = filtered
}

func (a *AmbientIndex) insertWorkloadToService(svcAddress networkAddress, workload *model.WorkloadInfo) {
	// For simplicity, to insert we drop it then add it to the end.
	// TODO: optimize this
	a.dropWorkloadFromService(svcAddress, workload.ResourceName())
	a.byService[svcAddress] = append(a.byService[svcAddress], workload)
}

func (a *AmbientIndex) updateWaypoint(scope model.WaypointScope, addr *workloadapi.GatewayAddress, isDelete bool, c *Controller) map[model.ConfigKey]struct{} {
	updates := sets.New[model.ConfigKey]()
	if isDelete {
		for _, wl := range a.byPod {
			if wl.Labels[constants.ManagedGatewayLabel] == constants.ManagedGatewayMeshControllerLabel {
				continue
			}
			if wl.Namespace != scope.Namespace || (scope.ServiceAccount != "" && wl.ServiceAccount != scope.ServiceAccount) {
				continue
			}

			if wl.Waypoint != nil && proto.Equal(wl.Waypoint, addr) {
				wl.Waypoint = nil
				// If there was a change, also update the VIPs and record for a push
				updates.Insert(model.ConfigKey{Kind: kind.Address, Name: wl.ResourceName()})
			}
			updates.Merge(c.updateEndpointsOnWaypointChange(wl))
		}
	} else {
		for _, wl := range a.byPod {
			if wl.Labels[constants.ManagedGatewayLabel] == constants.ManagedGatewayMeshControllerLabel {
				continue
			}
			if wl.Namespace != scope.Namespace || (scope.ServiceAccount != "" && wl.ServiceAccount != scope.ServiceAccount) {
				continue
			}

			if wl.Waypoint == nil || !proto.Equal(wl.Waypoint, addr) {
				wl.Waypoint = addr
				// If there was a change, also update the VIPs and record for a push
				updates.Insert(model.ConfigKey{Kind: kind.Address, Name: wl.ResourceName()})
			}
			updates.Merge(c.updateEndpointsOnWaypointChange(wl))
		}
	}
	return updates
}

// All return all known workloads. Result is un-ordered
func (a *AmbientIndex) All() []*model.AddressInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()
	res := make([]*model.AddressInfo, 0, len(a.byPod)+len(a.services))
	// byPod will not have any duplicates, so we can just iterate over that.
	for _, wl := range a.byPod {
		addr := &workloadapi.Address{
			Type: &workloadapi.Address_Workload{
				Workload: wl.Workload,
			},
		}
		res = append(res, &model.AddressInfo{Address: addr})
	}
	for _, s := range a.services {
		addr := &workloadapi.Address{
			Type: &workloadapi.Address_Service{
				Service: s,
			},
		}
		res = append(res, &model.AddressInfo{Address: addr})
	}
	return res
}

func (c *Controller) WorkloadsForWaypoint(scope model.WaypointScope) []*model.WorkloadInfo {
	a := c.ambientIndex
	a.mu.RLock()
	defer a.mu.RUnlock()
	var res []*model.WorkloadInfo
	// TODO: try to precompute
	for _, w := range a.byPod {
		if a.matchesScope(scope, w) {
			res = append(res, w)
		}
	}
	return res
}

// Waypoint finds all waypoint IP addresses for a given scope.  Performs first a Namespace+ServiceAccount
// then falls back to any Namespace wide waypoints
func (c *Controller) Waypoint(scope model.WaypointScope) []netip.Addr {
	a := c.ambientIndex
	a.mu.RLock()
	defer a.mu.RUnlock()
	// TODO need to handle case where waypoints are dualstack/have multiple addresses
	if addr, f := a.waypoints[scope]; f {
		switch address := addr.Destination.(type) {
		case *workloadapi.GatewayAddress_Address:
			if ip, ok := netip.AddrFromSlice(address.Address.GetAddress()); ok {
				return []netip.Addr{ip}
			}
		case *workloadapi.GatewayAddress_Hostname:
			// TODO
		}
	}

	// Now look for namespace-wide
	scope.ServiceAccount = ""
	if addr, f := a.waypoints[scope]; f {
		switch address := addr.Destination.(type) {
		case *workloadapi.GatewayAddress_Address:
			if ip, ok := netip.AddrFromSlice(address.Address.GetAddress()); ok {
				return []netip.Addr{ip}
			}
		case *workloadapi.GatewayAddress_Hostname:
			// TODO
		}
	}

	return []netip.Addr{}
}

func (a *AmbientIndex) matchesScope(scope model.WaypointScope, w *model.WorkloadInfo) bool {
	if len(scope.ServiceAccount) == 0 {
		// We are a namespace wide waypoint. SA scope take precedence.
		// Check if there is one for this workloads service account
		if _, f := a.waypoints[model.WaypointScope{Namespace: scope.Namespace, ServiceAccount: w.ServiceAccount}]; f {
			return false
		}
	}
	if scope.ServiceAccount != "" && (w.ServiceAccount != scope.ServiceAccount) {
		return false
	}
	if w.Namespace != scope.Namespace {
		return false
	}
	// Filter out waypoints.
	if w.Labels[constants.ManagedGatewayLabel] == constants.ManagedGatewayMeshControllerLabel {
		return false
	}
	return true
}

func (c *Controller) Policies(requested sets.Set[model.ConfigKey]) []*security.Authorization {
	if !c.configCluster {
		return nil
	}
	cfgs := c.configController.List(gvk.AuthorizationPolicy, metav1.NamespaceAll)
	l := len(cfgs)
	if len(requested) > 0 {
		l = len(requested)
	}
	res := make([]*security.Authorization, 0, l)
	for _, cfg := range cfgs {
		k := model.ConfigKey{
			Kind:      kind.AuthorizationPolicy,
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
		}
		if len(requested) > 0 && !requested.Contains(k) {
			continue
		}
		pol := convertAuthorizationPolicy(c.meshWatcher.Mesh().GetRootNamespace(), cfg)
		if pol == nil {
			continue
		}
		res = append(res, pol)
	}
	return res
}

func (c *Controller) selectorAuthorizationPolicies(ns string, lbls map[string]string) []string {
	global := c.configController.List(gvk.AuthorizationPolicy, c.meshWatcher.Mesh().GetRootNamespace())
	local := c.configController.List(gvk.AuthorizationPolicy, ns)
	res := sets.New[string]()
	matches := func(c config.Config) bool {
		sel := c.Spec.(*v1beta1.AuthorizationPolicy).Selector
		if sel == nil {
			return false
		}
		return labels.Instance(sel.MatchLabels).SubsetOf(lbls)
	}

	for _, pl := range [][]config.Config{global, local} {
		for _, p := range pl {
			if matches(p) {
				res.Insert(p.Namespace + "/" + p.Name)
			}
		}
	}
	return sets.SortedList(res)
}

func (c *Controller) AuthorizationPolicyHandler(old config.Config, obj config.Config, ev model.Event) {
	getSelector := func(c config.Config) map[string]string {
		if c.Spec == nil {
			return nil
		}
		pol := c.Spec.(*v1beta1.AuthorizationPolicy)
		return pol.Selector.GetMatchLabels()
	}
	// Normal flow for AuthorizationPolicy will trigger XDS push, so we don't need to push those. But we do need
	// to update any relevant workloads and push them.
	sel := getSelector(obj)
	oldSel := getSelector(old)

	switch ev {
	case model.EventUpdate:
		if maps.Equal(sel, oldSel) {
			// Update event, but selector didn't change. No workloads to push.
			return
		}
	default:
		if sel == nil {
			// We only care about selector policies
			return
		}
	}

	pods := map[string]*v1.Pod{}
	for _, p := range c.getPodsInPolicy(obj.Namespace, sel) {
		pods[p.Status.PodIP] = p
	}
	if oldSel != nil {
		for _, p := range c.getPodsInPolicy(obj.Namespace, oldSel) {
			pods[p.Status.PodIP] = p
		}
	}

	updates := map[model.ConfigKey]struct{}{}
	for _, pod := range pods {
		newWl := c.extractWorkload(pod)
		if newWl != nil {
			// Update the pod, since it now has new VIP info
			c.ambientIndex.mu.Lock()
			c.ambientIndex.byPod[newWl.ResourceName()] = newWl
			c.ambientIndex.mu.Unlock()
			updates[model.ConfigKey{Kind: kind.Address, Name: newWl.ResourceName()}] = struct{}{}
		}
	}

	if len(updates) > 0 {
		c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
			ConfigsUpdated: updates,
			Reason:         []model.TriggerReason{model.AmbientUpdate},
		})
	}
}

func (c *Controller) getPodsInPolicy(ns string, sel map[string]string) []*v1.Pod {
	if ns == c.meshWatcher.Mesh().GetRootNamespace() {
		ns = metav1.NamespaceAll
	}
	return c.podsClient.List(ns, klabels.ValidatedSetSelector(sel))
}

func convertAuthorizationPolicy(rootns string, obj config.Config) *security.Authorization {
	pol := obj.Spec.(*v1beta1.AuthorizationPolicy)

	scope := security.Scope_WORKLOAD_SELECTOR
	if pol.Selector == nil {
		scope = security.Scope_NAMESPACE
		// TODO: TDA
		if rootns == obj.Namespace {
			scope = security.Scope_GLOBAL // TODO: global workload?
		}
	}
	action := security.Action_ALLOW
	switch pol.Action {
	case v1beta1.AuthorizationPolicy_ALLOW:
	case v1beta1.AuthorizationPolicy_DENY:
		action = security.Action_DENY
	default:
		return nil
	}
	opol := &security.Authorization{
		Name:      obj.Name,
		Namespace: obj.Namespace,
		Scope:     scope,
		Action:    action,
		Groups:    nil,
	}

	for _, rule := range pol.Rules {
		rules := handleRule(action, rule)
		if rules != nil {
			rg := &security.Group{
				Rules: rules,
			}
			opol.Groups = append(opol.Groups, rg)
		}
	}

	return opol
}

func anyNonEmpty[T any](arr ...[]T) bool {
	for _, a := range arr {
		if len(a) > 0 {
			return true
		}
	}
	return false
}

func handleRule(action security.Action, rule *v1beta1.Rule) []*security.Rules {
	toMatches := []*security.Match{}
	for _, to := range rule.To {
		op := to.Operation
		if action == security.Action_ALLOW && anyNonEmpty(op.Hosts, op.NotHosts, op.Methods, op.NotMethods, op.Paths, op.NotPaths) {
			// L7 policies never match for ALLOW
			// For DENY they will always match, so it is more restrictive
			return nil
		}
		match := &security.Match{
			DestinationPorts:    stringToPort(op.Ports),
			NotDestinationPorts: stringToPort(op.NotPorts),
		}
		// if !emptyRuleMatch(match) {
		toMatches = append(toMatches, match)
		//}
	}
	fromMatches := []*security.Match{}
	for _, from := range rule.From {
		op := from.Source
		if action == security.Action_ALLOW && anyNonEmpty(op.RemoteIpBlocks, op.NotRemoteIpBlocks, op.RequestPrincipals, op.NotRequestPrincipals) {
			// L7 policies never match for ALLOW
			// For DENY they will always match, so it is more restrictive
			return nil
		}
		match := &security.Match{
			SourceIps:     stringToIP(op.IpBlocks),
			NotSourceIps:  stringToIP(op.NotIpBlocks),
			Namespaces:    stringToMatch(op.Namespaces),
			NotNamespaces: stringToMatch(op.NotNamespaces),
			Principals:    stringToMatch(op.Principals),
			NotPrincipals: stringToMatch(op.NotPrincipals),
		}
		// if !emptyRuleMatch(match) {
		fromMatches = append(fromMatches, match)
		//}
	}

	rules := []*security.Rules{}
	if len(toMatches) > 0 {
		rules = append(rules, &security.Rules{Matches: toMatches})
	}
	if len(fromMatches) > 0 {
		rules = append(rules, &security.Rules{Matches: fromMatches})
	}
	for _, when := range rule.When {
		l4 := l4WhenAttributes.Contains(when.Key)
		if action == security.Action_ALLOW && !l4 {
			// L7 policies never match for ALLOW
			// For DENY they will always match, so it is more restrictive
			return nil
		}
		positiveMatch := &security.Match{
			Namespaces:       whenMatch("source.namespace", when, false, stringToMatch),
			Principals:       whenMatch("source.principal", when, false, stringToMatch),
			SourceIps:        whenMatch("source.ip", when, false, stringToIP),
			DestinationPorts: whenMatch("destination.port", when, false, stringToPort),
			DestinationIps:   whenMatch("destination.ip", when, false, stringToIP),

			NotNamespaces:       whenMatch("source.namespace", when, true, stringToMatch),
			NotPrincipals:       whenMatch("source.principal", when, true, stringToMatch),
			NotSourceIps:        whenMatch("source.ip", when, true, stringToIP),
			NotDestinationPorts: whenMatch("destination.port", when, true, stringToPort),
			NotDestinationIps:   whenMatch("destination.ip", when, true, stringToIP),
		}
		rules = append(rules, &security.Rules{Matches: []*security.Match{positiveMatch}})
	}
	return rules
}

var l4WhenAttributes = sets.New(
	"source.ip",
	"source.namespace",
	"source.principal",
	"destination.ip",
	"destination.port",
)

func whenMatch[T any](s string, when *v1beta1.Condition, invert bool, f func(v []string) []T) []T {
	if when.Key != s {
		return nil
	}
	if invert {
		return f(when.NotValues)
	}
	return f(when.Values)
}

func stringToMatch(rules []string) []*security.StringMatch {
	res := make([]*security.StringMatch, 0, len(rules))
	for _, v := range rules {
		var sm *security.StringMatch
		switch {
		case v == "*":
			sm = &security.StringMatch{MatchType: &security.StringMatch_Presence{}}
		case strings.HasPrefix(v, "*"):
			sm = &security.StringMatch{MatchType: &security.StringMatch_Suffix{
				Suffix: strings.TrimPrefix(v, "*"),
			}}
		case strings.HasSuffix(v, "*"):
			sm = &security.StringMatch{MatchType: &security.StringMatch_Prefix{
				Prefix: strings.TrimSuffix(v, "*"),
			}}
		default:
			sm = &security.StringMatch{MatchType: &security.StringMatch_Exact{
				Exact: v,
			}}
		}
		res = append(res, sm)
	}
	return res
}

func stringToPort(rules []string) []uint32 {
	res := make([]uint32, 0, len(rules))
	for _, m := range rules {
		p, err := strconv.ParseUint(m, 10, 32)
		if err != nil || p > 65535 {
			continue
		}
		res = append(res, uint32(p))
	}
	return res
}

func stringToIP(rules []string) []*security.Address {
	res := make([]*security.Address, 0, len(rules))
	for _, m := range rules {
		if len(m) == 0 {
			continue
		}

		var (
			ipAddr        netip.Addr
			maxCidrPrefix uint32
		)

		if strings.Contains(m, "/") {
			ipp, err := netip.ParsePrefix(m)
			if err != nil {
				continue
			}
			ipAddr = ipp.Addr()
			maxCidrPrefix = uint32(ipp.Bits())
		} else {
			ipa, err := netip.ParseAddr(m)
			if err != nil {
				continue
			}

			ipAddr = ipa
			maxCidrPrefix = uint32(ipAddr.BitLen())
		}

		res = append(res, &security.Address{
			Address: ipAddr.AsSlice(),
			Length:  maxCidrPrefix,
		})
	}
	return res
}

func (c *Controller) extractWorkload(p *v1.Pod) *model.WorkloadInfo {
	if p == nil || !IsPodRunning(p) || p.Spec.HostNetwork {
		return nil
	}
	var waypoint *workloadapi.GatewayAddress
	if p.Labels[constants.ManagedGatewayLabel] == constants.ManagedGatewayMeshControllerLabel {
		// Waypoints do not have waypoints
		waypoint = nil
	} else {
		// First check for a waypoint for our SA explicit
		// TODO: this is not robust against temporary waypoint downtime. We also need the users intent (Gateway).
		found := false
		if waypoint, found = c.ambientIndex.waypoints[model.WaypointScope{Namespace: p.Namespace, ServiceAccount: p.Spec.ServiceAccountName}]; !found {
			// if there are none, check namespace wide waypoints
			waypoint = c.ambientIndex.waypoints[model.WaypointScope{Namespace: p.Namespace}]
		}
	}

	policies := c.selectorAuthorizationPolicies(p.Namespace, p.Labels)
	wl := c.constructWorkload(p, waypoint, policies)
	if wl == nil {
		return nil
	}
	return &model.WorkloadInfo{
		Workload: wl,
		Labels:   p.Labels,
	}
}

// updateEndpointsOnWaypointChange ensures that endpoints are synced for Envoy clients. Envoy clients
// maintain information about waypoints for each destination in metadata. If the waypoint changes, we need
// to sync this metadata again (add/remove waypoint IP).
// This is only needed for waypoints, as a normal workload update will already trigger and EDS push.
func (c *Controller) updateEndpointsOnWaypointChange(wl *model.WorkloadInfo) sets.Set[model.ConfigKey] {
	updates := sets.New[model.ConfigKey]()
	// For each VIP, trigger and EDS update
	for vip := range wl.VirtualIps {
		for _, svc := range c.ambientIndex.serviceVipIndex.Lookup(vip) {
			updates.Insert(model.ConfigKey{
				Kind:      kind.ServiceEntry,
				Name:      string(kube.ServiceHostname(svc.Name, svc.Namespace, c.opts.DomainSuffix)),
				Namespace: svc.Namespace,
			})
		}
	}
	return updates
}

func (c *Controller) setupIndex() *AmbientIndex {
	idx := AmbientIndex{
		byService: map[networkAddress][]*model.WorkloadInfo{},
		byPod:     map[string]*model.WorkloadInfo{},
		waypoints: map[model.WaypointScope]*workloadapi.GatewayAddress{},
		services:  map[networkAddress]*workloadapi.Service{},
	}

	podHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			updates := idx.handlePod(nil, obj, false, c)
			if len(updates) > 0 {
				c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
					ConfigsUpdated: updates,
					Reason:         []model.TriggerReason{model.AmbientUpdate},
				})
			}
		},
		UpdateFunc: func(oldObj, newObj any) {
			updates := idx.handlePod(oldObj, newObj, false, c)
			if len(updates) > 0 {
				c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
					ConfigsUpdated: updates,
					Reason:         []model.TriggerReason{model.AmbientUpdate},
				})
			}
		},
		DeleteFunc: func(obj any) {
			updates := idx.handlePod(nil, obj, true, c)
			if len(updates) > 0 {
				c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
					ConfigsUpdated: updates,
					Reason:         []model.TriggerReason{model.AmbientUpdate},
				})
			}
		},
	}
	c.podsClient.AddEventHandler(podHandler)

	serviceHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			idx.mu.Lock()
			defer idx.mu.Unlock()
			updates := idx.handleService(obj, false, c)
			if len(updates) > 0 {
				c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
					ConfigsUpdated: updates,
					Reason:         []model.TriggerReason{model.AmbientUpdate},
				})
			}
		},
		UpdateFunc: func(oldObj, newObj any) {
			idx.mu.Lock()
			defer idx.mu.Unlock()
			updates := idx.handleService(oldObj, true, c)
			updates2 := idx.handleService(newObj, false, c)
			if updates == nil {
				updates = updates2
			} else {
				for k, v := range updates2 {
					updates[k] = v
				}
			}

			if len(updates) > 0 {
				c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
					ConfigsUpdated: updates,
					Reason:         []model.TriggerReason{model.AmbientUpdate},
				})
			}
		},
		DeleteFunc: func(obj any) {
			idx.mu.Lock()
			defer idx.mu.Unlock()
			updates := idx.handleService(obj, true, c)
			if len(updates) > 0 {
				c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
					ConfigsUpdated: updates,
					Reason:         []model.TriggerReason{model.AmbientUpdate},
				})
			}
		},
	}
	c.services.AddEventHandler(serviceHandler)
	idx.serviceVipIndex = kclient.CreateIndex[string, *v1.Service](c.services, getVIPs)
	return &idx
}

func (a *AmbientIndex) handlePod(oldObj, newObj any, isDelete bool, c *Controller) sets.Set[model.ConfigKey] {
	p := controllers.Extract[*v1.Pod](newObj)
	old := controllers.Extract[*v1.Pod](oldObj)
	if old != nil {
		// compare only labels and pod phase, which are what we care about
		if maps.Equal(old.Labels, p.Labels) &&
			maps.Equal(old.Annotations, p.Annotations) &&
			old.Status.Phase == p.Status.Phase &&
			IsPodReady(old) == IsPodReady(p) {
			return nil
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	updates := sets.New[model.ConfigKey]()

	var wl *model.WorkloadInfo
	if !isDelete {
		wl = c.extractWorkload(p)
	}
	oldWl := a.byPod[c.Network(p.Status.PodIP, p.Labels).String()+"/"+p.Status.PodIP]
	if wl == nil {
		// This is an explicit delete event, or there is no longer a Workload to create (pod NotReady, etc)

		if oldWl != nil {
			delete(a.byPod, oldWl.ResourceName())
			// If we already knew about this workload, we need to make sure we drop all VIP references as well
			for vip := range oldWl.VirtualIps {
				a.dropWorkloadFromService(networkAddress{network: oldWl.Network, vip: vip}, oldWl.ResourceName())
			}
			log.Debugf("%v: workload removed, pushing", p.Status.PodIP)
			// TODO: namespace for network?
			updates.Insert(model.ConfigKey{Kind: kind.Address, Name: oldWl.ResourceName()})
			return updates
		}
		// It was a 'delete' for a resource we didn't know yet, no need to send an event

		return updates
	}
	if oldWl != nil && proto.Equal(wl.Workload, oldWl.Workload) {
		log.Debugf("%v: no change, skipping", wl.ResourceName())

		return updates
	}
	a.byPod[wl.ResourceName()] = wl
	if oldWl != nil {
		// For updates, we will drop the VIPs and then add the new ones back. This could be optimized
		for vip := range oldWl.VirtualIps {
			a.dropWorkloadFromService(networkAddress{network: oldWl.Network, vip: vip}, oldWl.ResourceName())
		}
	}
	// Update the VIP indexes as well, as needed
	for vip := range wl.VirtualIps {
		a.insertWorkloadToService(networkAddress{network: wl.Network, vip: vip}, wl)
	}

	log.Debugf("%v: workload updated, pushing", wl.ResourceName())
	updates.Insert(model.ConfigKey{Kind: kind.Address, Name: wl.ResourceName()})

	return updates
}

func (a *AmbientIndex) handlePods(pods []*v1.Pod, c *Controller) {
	updates := sets.New[model.ConfigKey]()
	for _, p := range pods {
		updates = updates.Merge(a.handlePod(nil, p, false, c))
	}
	if len(updates) > 0 {
		c.opts.XDSUpdater.ConfigUpdate(&model.PushRequest{
			ConfigsUpdated: updates,
			Reason:         []model.TriggerReason{model.AmbientUpdate},
		})
	}
}

func (a *AmbientIndex) handleService(obj any, isDelete bool, c *Controller) sets.Set[model.ConfigKey] {
	svc := controllers.Extract[*v1.Service](obj)
	updates := sets.New[model.ConfigKey]()

	if svc.Labels[constants.ManagedGatewayLabel] == constants.ManagedGatewayMeshControllerLabel {
		scope := model.WaypointScope{Namespace: svc.Namespace, ServiceAccount: svc.Annotations[constants.WaypointServiceAccount]}

		// TODO get IP+Port from the Gateway CRD
		// https://github.com/istio/istio/issues/44230
		if svc.Spec.ClusterIP == "None" {
			// TODO handle headless Service
			return updates
		}
		svcIP := netip.MustParseAddr(svc.Spec.ClusterIP)
		addr := &workloadapi.GatewayAddress{
			Destination: &workloadapi.GatewayAddress_Address{
				Address: &workloadapi.NetworkAddress{
					Network: c.Network(svcIP.String(), make(labels.Instance, 0)).String(),
					Address: svcIP.AsSlice(),
				},
			},
			Port: uint32(svc.Spec.Ports[0].Port),
		}

		if isDelete {
			if proto.Equal(a.waypoints[scope], addr) {
				delete(a.waypoints, scope)
				updates.Merge(a.updateWaypoint(scope, addr, true, c))
			}
		} else {
			if !proto.Equal(a.waypoints[scope], addr) {
				a.waypoints[scope] = addr
				updates.Merge(a.updateWaypoint(scope, addr, false, c))
			}
		}
	}

	vips := getVIPs(svc)
	// vips are iterated over multiple times, perform their network lookup only once
	// TODO does this need to be performed for each vip?  or should every vip for a service be on the same network?
	networkAddrs := make([]networkAddress, 0, len(vips))
	for _, vip := range vips {
		networkAddrs = append(networkAddrs, networkAddress{
			vip:     vip,
			network: c.Network(vip, make(labels.Instance, 0)).String(),
		})
	}
	pods := c.getPodsInService(svc)
	var wls []*model.WorkloadInfo
	for _, p := range pods {
		// Can be nil if it's not ready, hostNetwork, etc
		wl := c.extractWorkload(p)
		if wl != nil {
			// Update the pod, since it now has new VIP info
			a.byPod[wl.ResourceName()] = wl
			wls = append(wls, wl)
		}

	}

	// We send an update for each *workload* IP address previously in the service; they may have changed
	for _, networkAddr := range networkAddrs {
		for _, wl := range a.byService[networkAddr] {
			updates.Insert(model.ConfigKey{Kind: kind.Address, Name: wl.ResourceName()})
		}
	}
	// Update indexes
	if isDelete {
		for _, networkAddr := range networkAddrs {
			delete(a.byService, networkAddr)
			delete(a.services, networkAddr)
			updates.Insert(model.ConfigKey{Kind: kind.Address, Name: networkAddr.String()})
		}
	} else {
		ports := make([]*workloadapi.Port, 0, len(svc.Spec.Ports))
		for _, p := range svc.Spec.Ports {
			ports = append(ports, &workloadapi.Port{
				ServicePort: uint32(p.Port),
				TargetPort:  uint32(p.TargetPort.IntVal),
			})
		}
		ws := &workloadapi.Service{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Hostname: string(model.ResolveShortnameToFQDN(svc.Name, config.Meta{
				Namespace: svc.Namespace,
				Domain:    "cluster.local",
			})),
			Addresses: networkVipToNetworkAddress(networkAddrs),
			Ports:     ports,
		}
		for _, networkAddr := range networkAddrs {
			a.byService[networkAddr] = wls

			a.services[networkAddr] = ws
			updates.Insert(model.ConfigKey{Kind: kind.Address, Name: networkAddr.String()})
		}
	}
	// Fetch updates again, in case it changed from adding new workloads
	for _, networkAddr := range networkAddrs {
		for _, wl := range a.byService[networkAddr] {
			updates.Insert(model.ConfigKey{Kind: kind.Address, Name: wl.ResourceName()})
		}
	}

	return updates
}

func networkVipToNetworkAddress(vips []networkAddress) []*workloadapi.NetworkAddress {
	res := make([]*workloadapi.NetworkAddress, 0, len(vips))
	for _, vip := range vips {
		res = append(res, &workloadapi.NetworkAddress{
			Network: vip.network,
			Address: netip.MustParseAddr(vip.vip).AsSlice(),
		})
	}
	return res
}

func (c *Controller) getPodsInService(svc *v1.Service) []*v1.Pod {
	if svc.Spec.Selector == nil {
		// services with nil selectors match nothing, not everything.
		return nil
	}
	return c.podsClient.List(svc.Namespace, klabels.ValidatedSetSelector(svc.Spec.Selector))
}

// AddressInformation returns all AddressInfo's in the cluster.
// This may be scoped to specific subsets by specifying a non-empty addresses field
func (c *Controller) AddressInformation(addresses sets.Set[types.NamespacedName]) ([]*model.AddressInfo, []string) {
	if len(addresses) == 0 {
		// Full update
		return c.ambientIndex.All(), nil
	}
	var wls []*model.AddressInfo
	var removed []string
	for p := range addresses {
		wname := p.Name
		// GenerateDeltas has the formatted wname from the xds request, but not sure if other callers
		// have the format enforced
		if found := strings.Count(p.Name, "/"); found == 0 {
			cNetwork := c.Network(p.Name, make(labels.Instance, 0)).String()
			wname = cNetwork + "/" + p.Name
		}
		wl := c.ambientIndex.Lookup(wname)
		if len(wl) == 0 {
			removed = append(removed, p.Name)
		} else {
			wls = append(wls, wl...)
		}
	}
	return wls, removed
}

func (c *Controller) constructWorkload(pod *v1.Pod, waypoint *workloadapi.GatewayAddress, policies []string) *workloadapi.Workload {
	vips := map[string]*workloadapi.PortList{}
	allServices := c.services.List(pod.Namespace, klabels.Everything())
	if services := getPodServices(allServices, pod); len(services) > 0 {
		for _, svc := range services {
			for _, vip := range getVIPs(svc) {
				if vips[vip] == nil {
					vips[vip] = &workloadapi.PortList{}
				}
				for _, port := range svc.Spec.Ports {
					if port.Protocol != v1.ProtocolTCP {
						continue
					}
					targetPort, err := FindPort(pod, &port)
					if err != nil {
						log.Debug(err)
						continue
					}
					vips[vip].Ports = append(vips[vip].Ports, &workloadapi.Port{
						ServicePort: uint32(port.Port),
						TargetPort:  uint32(targetPort),
					})
				}
			}
		}
	}

	wl := &workloadapi.Workload{
		Name:                  pod.Name,
		Address:               parseIP(pod.Status.PodIP),
		Network:               c.Network(pod.Status.PodIP, pod.Labels).String(),
		Namespace:             pod.Namespace,
		ServiceAccount:        pod.Spec.ServiceAccountName,
		Node:                  pod.Spec.NodeName,
		VirtualIps:            vips,
		AuthorizationPolicies: policies,
		Status:                workloadapi.WorkloadStatus_HEALTHY,
		ClusterId:             c.Cluster().String(),
	}
	if !IsPodReady(pod) {
		wl.Status = workloadapi.WorkloadStatus_UNHEALTHY
	}
	if td := spiffe.GetTrustDomain(); td != "cluster.local" {
		wl.TrustDomain = td
	}

	wl.WorkloadName, wl.WorkloadType = workloadNameAndType(pod)
	wl.CanonicalName, wl.CanonicalRevision = kubelabels.CanonicalService(pod.Labels, wl.WorkloadName)
	// If we have a remote proxy, configure it
	if waypoint != nil {
		wl.Waypoint = waypoint
	}

	if pod.Annotations[constants.AmbientRedirection] == constants.AmbientRedirectionEnabled {
		// Configured for override
		wl.TunnelProtocol = workloadapi.TunnelProtocol_HBONE
	}
	// Otherwise supports tunnel directly
	if model.SupportsTunnel(pod.Labels, model.TunnelHTTP) {
		wl.TunnelProtocol = workloadapi.TunnelProtocol_HBONE
		wl.NativeTunnel = true
	}
	return wl
}

func parseIP(ip string) []byte {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return nil
	}
	return addr.AsSlice()
}

// internal object used for
type networkAddress struct {
	network string
	vip     string
}

func (n *networkAddress) String() string {
	return n.network + "/" + n.vip
}

func getVIPs(svc *v1.Service) []string {
	res := []string{}
	if svc.Spec.ClusterIP != "" && svc.Spec.ClusterIP != "None" {
		res = append(res, svc.Spec.ClusterIP)
	}
	for _, ing := range svc.Status.LoadBalancer.Ingress {
		res = append(res, ing.IP)
	}
	return res
}

func (c *Controller) AdditionalPodSubscriptions(
	proxy *model.Proxy,
	allAddresses sets.Set[types.NamespacedName],
	currentSubs sets.Set[types.NamespacedName],
) sets.Set[types.NamespacedName] {
	shouldSubscribe := sets.New[types.NamespacedName]()

	// First, we want to handle VIP subscriptions. Example:
	// Client subscribes to VIP1. Pod1, part of VIP1, is sent.
	// The client wouldn't be explicitly subscribed to Pod1, so it would normally ignore it.
	// Since it is a part of VIP1 which we are subscribe to, add it to the subscriptions
	for s := range allAddresses {
		cNetwork := c.Network(s.Name, make(labels.Instance, 0)).String()
		for _, wl := range c.ambientIndex.Lookup(cNetwork + "/" + s.Name) {
			// We may have gotten an update for Pod, but are subscribe to a Service.
			// We need to force a subscription on the Pod as well
			switch addr := wl.Address.Type.(type) {
			case *workloadapi.Address_Workload:
				for vip := range addr.Workload.VirtualIps {
					t := types.NamespacedName{Name: vip}
					if currentSubs.Contains(t) {
						shouldSubscribe.Insert(types.NamespacedName{Name: wl.ResourceName()})
						break
					}
				}
			case *workloadapi.Address_Service:
				for _, networkAddress := range addr.Service.Addresses {
					t := types.NamespacedName{Name: string(networkAddress.Address)}
					if currentSubs.Contains(t) {
						shouldSubscribe.Insert(types.NamespacedName{Name: wl.ResourceName()})
						break
					}
				}
			}
		}
	}

	// Next, as an optimization, we will send all node-local endpoints
	if nodeName := proxy.Metadata.NodeName; nodeName != "" {
		for _, wl := range c.ambientIndex.All() {
			switch addr := wl.Address.Type.(type) {
			case *workloadapi.Address_Workload:
				if addr.Workload.Node == nodeName {
					n := types.NamespacedName{Name: wl.ResourceName()}
					if currentSubs.Contains(n) {
						continue
					}
					shouldSubscribe.Insert(n)
				}
			case *workloadapi.Address_Service:
				// Services are not constrained to a particular node
				continue
			}
		}
	}

	return shouldSubscribe
}

func workloadNameAndType(pod *v1.Pod) (string, workloadapi.WorkloadType) {
	if len(pod.GenerateName) == 0 {
		return pod.Name, workloadapi.WorkloadType_POD
	}

	// if the pod name was generated (or is scheduled for generation), we can begin an investigation into the controlling reference for the pod.
	var controllerRef metav1.OwnerReference
	controllerFound := false
	for _, ref := range pod.GetOwnerReferences() {
		if ref.Controller != nil && *ref.Controller {
			controllerRef = ref
			controllerFound = true
			break
		}
	}

	if !controllerFound {
		return pod.Name, workloadapi.WorkloadType_POD
	}

	// heuristic for deployment detection
	if controllerRef.Kind == "ReplicaSet" && strings.HasSuffix(controllerRef.Name, pod.Labels["pod-template-hash"]) {
		name := strings.TrimSuffix(controllerRef.Name, "-"+pod.Labels["pod-template-hash"])
		return name, workloadapi.WorkloadType_DEPLOYMENT
	}

	if controllerRef.Kind == "Job" {
		// figure out how to go from Job -> CronJob
		return controllerRef.Name, workloadapi.WorkloadType_JOB
	}

	if controllerRef.Kind == "CronJob" {
		// figure out how to go from Job -> CronJob
		return controllerRef.Name, workloadapi.WorkloadType_CRONJOB
	}

	return pod.Name, workloadapi.WorkloadType_POD
}
