// Copyright 2020 Istio Authors
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

package model

import (
	"strings"

	"github.com/gogo/protobuf/jsonpb"

	networking "istio.io/api/networking/v1alpha3"

	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/util/sets"
)

func mergeVirtualServicesIfNeeded(vServices []Config) (out []Config) {
	out = make([]Config, 0, len(vServices))
	delegatesMap := map[string]Config{}
	// root virtualservices with delegate
	var rootVses []Config

	// 1. classify virtualservices
	for _, vs := range vServices {
		rule := vs.Spec.(*networking.VirtualService)
		// it is delegate, add it to the indexer cache
		if len(rule.Hosts) == 0 {
			delegatesMap[key(vs.Name, vs.Namespace)] = vs
			continue
		}

		// root vs
		if isRootVs(rule) {
			rootVses = append(rootVses, vs)
			continue
		}

		// the others are normal vs without delegate
		out = append(out, vs)
	}

	// If `PILOT_ENABLE_VIRTUAL_SERVICE_DELEGATE` feature disabled,
	// filter out invalid vs(root or delegate), this can happen after enable -> disable
	if !features.EnableVirtualServiceDelegate.Get() {
		return
	}

	// 2. merge delegates and root
	for _, root := range rootVses {
		rootVs := root.Spec.(*networking.VirtualService)
		mergedRoutes := []*networking.HTTPRoute{}
		for _, route := range rootVs.Http {
			// it is root vs with delegate
			if route.Delegate != nil {
				delegate, ok := delegatesMap[key(route.Delegate.Name, route.Delegate.Namespace)]
				if !ok {
					log.Debugf("delegate virtual service %s/%s of %s/%s not found",
						route.Delegate.Namespace, route.Delegate.Name, root.Namespace, root.Name)
					// delegate not found, ignore only the current HTTP route
					continue
				}
				// DeepCopy to prevent mutate the original delegate, it can conflict
				// when multiple routes delegate to one single VS.
				copiedDelegate := delegate.DeepCopy()
				vs := copiedDelegate.Spec.(*networking.VirtualService)
				merged := mergeHTTPRoute(route, vs.Http)
				mergedRoutes = append(mergedRoutes, merged...)
			} else {
				mergedRoutes = append(mergedRoutes, route)
			}
		}
		rootVs.Http = mergedRoutes
		if log.DebugEnabled() {
			jsonm := &jsonpb.Marshaler{Indent: "   "}
			vsString, _ := jsonm.MarshalToString(rootVs)
			log.Infof("merged virtualService: %s", vsString)
		}
		out = append(out, root)
	}

	return
}

// merge root's route with delegate's and the merged route number equals the delegate's.
// if there is a conflict with root, the route is ignored
func mergeHTTPRoute(root *networking.HTTPRoute, delegate []*networking.HTTPRoute) []*networking.HTTPRoute {
	root.Delegate = nil

	out := make([]*networking.HTTPRoute, 0, len(delegate))
	for _, subRoute := range delegate {
		// suppose there are N1 match conditions in root, N2 match conditions in delegate
		// if match condition of N2 is a subset of anyone in N1, this is a valid matching conditions
		merged, conflict := mergeHTTPMatchRequests(root.Match, subRoute.Match)
		if conflict {
			log.Debugf("HTTPMatchRequests conflict: root root %s, delegate root %s", root.Name, subRoute.Name)
			continue
		}
		subRoute.Match = merged

		if subRoute.Name == "" {
			subRoute.Name = root.Name
		}
		if subRoute.Rewrite == nil {
			subRoute.Rewrite = root.Rewrite
		}
		if subRoute.Timeout == nil {
			subRoute.Timeout = root.Timeout
		}
		if subRoute.Retries == nil {
			subRoute.Retries = root.Retries
		}
		if subRoute.Fault == nil {
			subRoute.Fault = root.Fault
		}
		if subRoute.Mirror == nil {
			subRoute.Mirror = root.Mirror
		}
		if subRoute.MirrorPercentage == nil {
			subRoute.MirrorPercentage = root.MirrorPercentage
		}
		if subRoute.CorsPolicy == nil {
			subRoute.CorsPolicy = root.CorsPolicy
		}
		if subRoute.Headers == nil {
			subRoute.Headers = root.Headers
		}

		out = append(out, subRoute)
	}
	return out
}

// return merged match conditions if not conflicts
func mergeHTTPMatchRequests(root, delegate []*networking.HTTPMatchRequest) (out []*networking.HTTPMatchRequest, conflict bool) {
	if len(root) == 0 {
		return delegate, false
	}

	if len(delegate) == 0 {
		return root, false
	}

	// each HTTPMatchRequest of delegate must find a superset in root.
	// otherwise it conflicts
	for _, subMatch := range delegate {
		foundMatch := false
		for _, rootMatch := range root {
			if hasConflict(rootMatch, subMatch) {
				continue
			}
			// merge HTTPMatchRequest
			out = append(out, mergeHTTPMatchRequest(rootMatch, subMatch))
			foundMatch = true
		}
		if !foundMatch {
			return nil, true
		}
	}
	if len(out) == 0 {
		conflict = true
	}
	return
}

func mergeHTTPMatchRequest(root, delegate *networking.HTTPMatchRequest) *networking.HTTPMatchRequest {
	out := *delegate
	if out.Name == "" {
		out.Name = root.Name
	}
	if out.Uri == nil {
		out.Uri = root.Uri
	}
	if out.Scheme == nil {
		out.Scheme = root.Scheme
	}
	if out.Method == nil {
		out.Method = root.Method
	}
	if out.Authority == nil {
		out.Authority = root.Authority
	}
	// headers
	if len(root.Headers) > 0 || len(delegate.Headers) > 0 {
		out.Headers = make(map[string]*networking.StringMatch)
	}
	for k, v := range root.Headers {
		out.Headers[k] = v
	}
	for k, v := range delegate.Headers {
		out.Headers[k] = v
	}
	// withoutheaders
	if len(root.WithoutHeaders) > 0 || len(delegate.WithoutHeaders) > 0 {
		out.WithoutHeaders = make(map[string]*networking.StringMatch)
	}
	for k, v := range root.WithoutHeaders {
		out.WithoutHeaders[k] = v
	}
	for k, v := range delegate.WithoutHeaders {
		out.WithoutHeaders[k] = v
	}
	// queryparams
	if len(root.QueryParams) > 0 || len(delegate.QueryParams) > 0 {
		out.QueryParams = make(map[string]*networking.StringMatch)
	}
	for k, v := range root.QueryParams {
		out.QueryParams[k] = v
	}
	for k, v := range delegate.QueryParams {
		out.QueryParams[k] = v
	}
	// sourcelabels
	if len(root.SourceLabels) > 0 || len(delegate.SourceLabels) > 0 {
		out.SourceLabels = make(map[string]string)
	}
	for k, v := range root.SourceLabels {
		out.SourceLabels[k] = v
	}
	for k, v := range delegate.SourceLabels {
		out.SourceLabels[k] = v
	}
	if len(out.Gateways) == 0 {
		out.Gateways = root.Gateways
	}
	if out.SourceNamespace == "" {
		out.SourceNamespace = root.SourceNamespace
	}
	if out.Port == 0 {
		out.Port = root.Port
	}
	return &out
}

func hasConflict(root, leaf *networking.HTTPMatchRequest) bool {
	roots := []*networking.StringMatch{root.Uri, root.Scheme, root.Method, root.Authority}
	leaves := []*networking.StringMatch{leaf.Uri, leaf.Scheme, leaf.Method, leaf.Authority}
	for i := range roots {
		if stringMatchConflict(roots[i], leaves[i]) {
			return true
		}
	}
	// header conflicts
	for key, leafHeader := range leaf.Headers {
		if stringMatchConflict(root.Headers[key], leafHeader) {
			return true
		}
	}

	// without headers
	for key, leafValue := range leaf.WithoutHeaders {
		if stringMatchConflict(root.Headers[key], leafValue) {
			return true
		}
	}

	// query params conflict
	for key, value := range leaf.QueryParams {
		if stringMatchConflict(root.QueryParams[key], value) {
			return true
		}
	}

	// source labels
	for key, leafValue := range leaf.SourceLabels {
		if rootValue, ok := root.SourceLabels[key]; ok && rootValue != leafValue {
			return true
		}
	}

	if len(leaf.Gateways) > 0 {
		leafGws := sets.NewSet(leaf.Gateways...)
		for _, gw := range root.Gateways {
			if leafGws.Contains(gw) {
				return true
			}
		}
	}

	if root.IgnoreUriCase != leaf.IgnoreUriCase {
		return true
	}
	if root.SourceNamespace != "" && leaf.SourceNamespace != "" && root.SourceNamespace != leaf.SourceNamespace {
		return true
	}
	if root.Port > 0 && leaf.Port > 0 && root.Port != leaf.Port {
		return true
	}

	return false
}

func stringMatchConflict(root, leaf *networking.StringMatch) bool {
	// no conflict when root or leaf is not specified
	if root == nil || leaf == nil {
		return false
	}
	// regex match is not allowed
	if root.GetRegex() != "" || leaf.GetRegex() != "" {
		return true
	}
	// root is exact match
	if exact := root.GetExact(); exact != "" {
		// leaf is prefix match, conflict
		if leaf.GetPrefix() != "" {
			return true
		}
		// both exact, but not equal
		if leaf.GetExact() != exact {
			return true
		}
		return false
	}
	// root is prefix match
	if prefix := root.GetPrefix(); prefix != "" {
		// leaf is prefix match
		if leaf.GetPrefix() != "" {
			// leaf(`/a`) is not subset of root(`/a/b`)
			return !strings.HasPrefix(leaf.GetPrefix(), prefix)
		}
		// leaf is exact match
		if leaf.GetExact() != "" {
			// leaf(`/a`) is not subset of root(`/a/b`)
			return !strings.HasPrefix(leaf.GetExact(), prefix)
		}
	}

	return true
}

func isRootVs(vs *networking.VirtualService) bool {
	for _, route := range vs.Http {
		// it is root vs with delegate
		if route.Delegate != nil {
			return true
		}
	}
	return false
}
