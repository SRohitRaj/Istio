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

package xds

import (
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/util/protoconv"
	v3 "istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pkg/config/schema/kind"
	"istio.io/istio/pkg/util/sets"
)

type WorkloadGenerator struct {
	Server *DiscoveryServer
}

var (
	_ model.XdsResourceGenerator      = &WorkloadGenerator{}
	_ model.XdsDeltaResourceGenerator = &WorkloadGenerator{}
)

// GenerateDeltas computes Workload resources. This is design to be highly optimized to delta updates,
// and supports *on-demand* client usage. A client can subscribe with a wildcard subscription and get all
// resources (with delta updates), or on-demand and only get responses for specifically subscribed resources.
//
// Incoming requests may be for VIP or Pod IP addresses. However, all responses are Workload resources, which are pod based.
// This means subscribing to a VIP may end up pushing many resources of different name than the request.
// On-demand clients are expected to handle this (for wildcard, this is not applicable, as they don't specify any resources at all).
func (e WorkloadGenerator) GenerateDeltas(
	proxy *model.Proxy,
	req *model.PushRequest,
	w *model.WatchedResource,
) (model.Resources, model.DeletedResources, model.XdsLogDetails, bool, error) {
	updatedAddresses := model.ConfigNameOfKind(req.ConfigsUpdated, kind.Address)
	isReq := req.IsRequest()
	if len(updatedAddresses) == 0 && len(req.ConfigsUpdated) > 0 {
		// Nothing changed..
		return nil, nil, model.XdsLogDetails{}, false, nil
	}

	subs := sets.New(w.ResourceNames...)
	var addresses sets.String
	if isReq {
		// this is from request, we only send response for the subscribed address
		// At t0, a client request A, we only send A and aditional resources back to the client.
		// At t1, a client request B, we only send B and aditional resources back to the client, no A here.
		addresses = req.Delta.Subscribed
	} else {
		if w.Wildcard {
			addresses = updatedAddresses
		} else {
			// this is from the external triggers instead of request
			// send response for all the subscribed intersect with the updated
			addresses = updatedAddresses.Intersection(subs)
		}
	}

	if !w.Wildcard {
		// We only need this for on-demand. This allows us to subscribe the client to resources they
		// didn't explicitly request.
		// For wildcard, they subscribe to everything already.
		additional := e.Server.Env.ServiceDiscovery.AdditionalPodSubscriptions(proxy, addresses, subs)
		if addresses == nil {
			addresses = sets.New[string]()
		}
		addresses.Merge(additional)
	}

	// TODO: it is needlessly wasteful to do a full sync just because the rest of Istio thought it was "full"
	// The only things that can really trigger a "full" push here is trust domain or network changing, which is extremely rare
	// We do a full push for wildcard requests (initial proxy sync) or for full pushes with no ConfigsUpdates (since we don't know what changed)
	full := (isReq && w.Wildcard) || (!isReq && req.Full && len(req.ConfigsUpdated) == 0)

	// Nothing to do
	if len(addresses) == 0 && !full {
		if isReq {
			// We need to respond for requests, even if we have nothing to respond with
			return make(model.Resources, 0), nil, model.XdsLogDetails{}, false, nil
		}
		// For NOP pushes, no need
		return nil, nil, model.XdsLogDetails{}, false, nil
	}
	resources := make(model.Resources, 0)
	addrs, removed := e.Server.Env.ServiceDiscovery.AddressInformation(addresses)
	// Note: while "removed" is a weird name for a resource that never existed, this is how the spec works:
	// https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol#id2
	have := sets.New[string]()
	for _, addr := range addrs {
		// TODO(@hzxuzhonghu): calculate removed with aliases in `AddressInformation`
		aliases := addr.Aliases()
		removed.DeleteAll(aliases...)
		n := addr.ResourceName()
		have.Insert(n)
		switch w.TypeUrl {
		case v3.WorkloadType:
			if addr.GetWorkload() != nil {
				resources = append(resources, &discovery.Resource{
					Name:     n,
					Aliases:  aliases,
					Resource: protoconv.MessageToAny(addr.GetWorkload()), // TODO: pre-marshal
				})
			}
		case v3.AddressType:
			resources = append(resources, &discovery.Resource{
				Name:     n,
				Aliases:  aliases,
				Resource: protoconv.MessageToAny(addr), // TODO: pre-marshal
			})
		}
	}

	if full {
		// If it's a full push, AddressInformation won't have info to compute the full set of removals.
		// Instead, we need can see what resources are missing that we were subscribe to; those were removed.
		removed = subs.Difference(have).Merge(removed)
	}

	if !w.Wildcard {
		// For on-demand, we may have requested a VIP but gotten Pod IPs back. We need to update
		// the internal book-keeping to subscribe to the Pods, so that we push updates to those Pods.
		w.ResourceNames = subs.Merge(have).UnsortedList()
	} else {
		// For wildcard, we record all resources that have been pushed and not removed
		// It was to correctly calculate removed resources during full push alongside with specific address removed.
		w.ResourceNames = subs.Merge(have).Difference(removed).UnsortedList()
	}
	return resources, removed.UnsortedList(), model.XdsLogDetails{}, true, nil
}

func (e WorkloadGenerator) Generate(proxy *model.Proxy, w *model.WatchedResource, req *model.PushRequest) (model.Resources, model.XdsLogDetails, error) {
	resources, _, details, _, err := e.GenerateDeltas(proxy, req, w)
	return resources, details, err
}

type WorkloadRBACGenerator struct {
	Server *DiscoveryServer
}

func (e WorkloadRBACGenerator) GenerateDeltas(
	proxy *model.Proxy,
	req *model.PushRequest,
	w *model.WatchedResource,
) (model.Resources, model.DeletedResources, model.XdsLogDetails, bool, error) {
	var updatedPolicies sets.Set[model.ConfigKey]
	if len(req.ConfigsUpdated) != 0 {
		updatedPolicies = model.ConfigsOfKind(req.ConfigsUpdated, kind.AuthorizationPolicy)
		// Convert the actual Kubernetes PeerAuthentication policies to the synthetic ones
		// by adding the prefix
		//
		// This is needed because the handler that produces the ConfigUpdate blindly sends
		// the Kubernetes resource names without context of the synthetic Ambient policies
		// TODO: Split out PeerAuthentication into a separate handler in
		// https://github.com/istio/istio/blob/master/pilot/pkg/bootstrap/server.go#L882
		for p := range model.ConfigsOfKind(req.ConfigsUpdated, kind.PeerAuthentication) {
			updatedPolicies.Insert(model.ConfigKey{
				Name:      model.GetAmbientPolicyConfigName(p),
				Namespace: p.Namespace,
				Kind:      p.Kind,
			})
		}
	}
	if len(req.ConfigsUpdated) != 0 && len(updatedPolicies) == 0 {
		// This was a incremental push for a resource we don't watch... skip
		return nil, nil, model.DefaultXdsLogDetails, false, nil
	}

	policies := e.Server.Env.ServiceDiscovery.Policies(updatedPolicies)

	resources := make(model.Resources, 0)
	expected := sets.New[string]()
	if len(updatedPolicies) > 0 {
		// Partial update. Removes are ones we request but didn't get back when querying the policies
		for k := range updatedPolicies {
			expected.Insert(k.Namespace + "/" + k.Name)
		}
	} else {
		// Full update, expect everything
		expected.InsertAll(w.ResourceNames...)
	}

	removed := expected
	for _, p := range policies {
		n := p.ResourceName()
		removed.Delete(n) // We found it, so it isn't a removal
		resources = append(resources, &discovery.Resource{
			Name:     n,
			Resource: protoconv.MessageToAny(p.Authorization),
		})
	}

	return resources, sets.SortedList(removed), model.XdsLogDetails{}, true, nil
}

func (e WorkloadRBACGenerator) Generate(proxy *model.Proxy, w *model.WatchedResource, req *model.PushRequest) (model.Resources, model.XdsLogDetails, error) {
	resources, _, details, _, err := e.GenerateDeltas(proxy, req, w)
	return resources, details, err
}

var (
	_ model.XdsResourceGenerator      = &WorkloadRBACGenerator{}
	_ model.XdsDeltaResourceGenerator = &WorkloadRBACGenerator{}
)
