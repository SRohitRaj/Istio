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
	"fmt"
	"strconv"
	"testing"

	model "istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/schema/gvk"
	"istio.io/istio/pkg/spiffe"
)

func TestProxyNeedsPush(t *testing.T) {
	const (
		invalidKind = "INVALID_KIND"
		svcName     = "svc1.com"
		drName      = "dr1"
		vsName      = "vs1"
		scName      = "sc1"
		nsName      = "ns1"
		nsRoot      = "rootns"
		generalName = "name1"

		invalidNameSuffix = "invalid"
	)

	type Case struct {
		name    string
		proxy   *model.Proxy
		configs map[model.ConfigKey]struct{}
		want    bool
	}

	proxyCfg := &config.Config{Meta: config.Meta{
		Name:      generalName,
		Namespace: nsName,
	}}

	sidecar := &model.Proxy{
		Type: model.SidecarProxy, IPAddresses: []string{"127.0.0.1"}, Metadata: &model.NodeMetadata{},
		SidecarScope: &model.SidecarScope{Config: proxyCfg, RootNamespace: nsRoot}}
	gateway := &model.Proxy{Type: model.Router}

	sidecarScopeKindNames := map[config.GroupVersionKind]string{
		gvk.ServiceEntry: svcName, gvk.VirtualService: vsName, gvk.DestinationRule: drName}
	for kind, name := range sidecarScopeKindNames {
		sidecar.SidecarScope.AddConfigDependencies(model.ConfigKey{Kind: kind, Name: name, Namespace: nsName})
	}
	for kind, types := range configKindAffectedProxyTypes {
		for _, nodeType := range types {
			if nodeType == model.SidecarProxy {
				sidecar.SidecarScope.AddConfigDependencies(model.ConfigKey{
					Kind:      kind,
					Name:      generalName,
					Namespace: nsName,
				})
			}
		}
	}

	cases := []Case{
		{"no namespace or configs", sidecar, nil, true},
		{"gateway config for sidecar", sidecar, map[model.ConfigKey]struct{}{
			{
				Kind: gvk.Gateway,
				Name: generalName, Namespace: nsName}: {}}, false},
		{"gateway config for gateway", gateway, map[model.ConfigKey]struct{}{
			{
				Kind: gvk.Gateway,
				Name: generalName, Namespace: nsName}: {}}, true},
		{"sidecar config for gateway", gateway, map[model.ConfigKey]struct{}{
			{
				Kind: gvk.Sidecar,
				Name: scName, Namespace: nsName}: {}}, false},
		{"invalid config for sidecar", sidecar, map[model.ConfigKey]struct{}{
			{
				Kind: config.GroupVersionKind{Kind: invalidKind}, Name: generalName, Namespace: nsName}: {}},
			true},
		{"mixture matched and unmatched config for sidecar", sidecar, map[model.ConfigKey]struct{}{
			{Kind: gvk.DestinationRule, Name: drName, Namespace: nsName}:                   {},
			{Kind: gvk.ServiceEntry, Name: svcName + invalidNameSuffix, Namespace: nsName}: {},
		}, true},
		{"mixture unmatched and unmatched config for sidecar", sidecar, map[model.ConfigKey]struct{}{
			{Kind: gvk.DestinationRule, Name: drName + invalidNameSuffix, Namespace: nsName}: {},
			{Kind: gvk.ServiceEntry, Name: svcName + invalidNameSuffix, Namespace: nsName}:   {},
		}, false},
		{"empty configsUpdated for sidecar", sidecar, nil, true},
	}

	for kind, name := range sidecarScopeKindNames {
		cases = append(cases, Case{ // valid name
			name:    fmt.Sprintf("%s config for sidecar", kind.Kind),
			proxy:   sidecar,
			configs: map[model.ConfigKey]struct{}{{Kind: kind, Name: name, Namespace: nsName}: {}},
			want:    true,
		}, Case{ // invalid name
			name:    fmt.Sprintf("%s unmatched config for sidecar", kind.Kind),
			proxy:   sidecar,
			configs: map[model.ConfigKey]struct{}{{Kind: kind, Name: name + invalidNameSuffix, Namespace: nsName}: {}},
			want:    false,
		})
	}

	sidecarNamespaceScopeTypes := []config.GroupVersionKind{
		gvk.Sidecar, gvk.EnvoyFilter, gvk.AuthorizationPolicy, gvk.RequestAuthentication,
	}
	for _, kind := range sidecarNamespaceScopeTypes {
		cases = append(cases,
			Case{
				name:    fmt.Sprintf("%s config for sidecar in same namespace", kind.Kind),
				proxy:   sidecar,
				configs: map[model.ConfigKey]struct{}{{Kind: kind, Name: generalName, Namespace: nsName}: {}},
				want:    true,
			},
			Case{
				name:    fmt.Sprintf("%s config for sidecar in different namespace", kind.Kind),
				proxy:   sidecar,
				configs: map[model.ConfigKey]struct{}{{Kind: kind, Name: generalName, Namespace: "invalid-namespace"}: {}},
				want:    false,
			},
		)
	}

	// tests for kind-affect-proxy.
	for kind, types := range configKindAffectedProxyTypes {
		for _, nodeType := range types {
			proxy := gateway
			if nodeType == model.SidecarProxy {
				proxy = sidecar
			}
			cases = append(cases, Case{
				name:  fmt.Sprintf("kind %s affect %s", kind, nodeType),
				proxy: proxy,
				configs: map[model.ConfigKey]struct{}{
					{Kind: kind, Name: generalName + invalidNameSuffix, Namespace: nsName}: {}},
				want: true,
			})
		}
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			pushEv := &Event{pushRequest: &model.PushRequest{ConfigsUpdated: tt.configs}}
			got := ProxyNeedsPush(tt.proxy, pushEv)
			if got != tt.want {
				t.Fatalf("Got needs push = %v, expected %v", got, tt.want)
			}
		})
	}
}

func BenchmarkListEquals(b *testing.B) {
	size := 100
	var l []string
	for i := 0; i < size; i++ {
		l = append(l, strconv.Itoa(i))
	}
	var equal []string
	for i := 0; i < size; i++ {
		equal = append(equal, strconv.Itoa(i))
	}
	var notEqual []string
	for i := 0; i < size; i++ {
		notEqual = append(notEqual, strconv.Itoa(i))
	}
	notEqual[size-1] = "z"

	for n := 0; n < b.N; n++ {
		listEqualUnordered(l, equal)
		listEqualUnordered(l, notEqual)
	}
}

func TestCheckConnectionIdentity(t *testing.T) {
	cases := []struct {
		name      string
		identity  []string
		sa        string
		namespace string
		success   bool
	}{
		{
			name:      "single match",
			identity:  []string{spiffe.Identity{"cluster.local", "namespace", "serviceaccount"}.String()},
			sa:        "serviceaccount",
			namespace: "namespace",
			success:   true,
		},
		{
			name: "second match",
			identity: []string{
				spiffe.Identity{"cluster.local", "bad", "serviceaccount"}.String(),
				spiffe.Identity{"cluster.local", "namespace", "serviceaccount"}.String(),
			},
			sa:        "serviceaccount",
			namespace: "namespace",
			success:   true,
		},
		{
			name: "no match namespace",
			identity: []string{
				spiffe.Identity{"cluster.local", "bad", "serviceaccount"}.String(),
			},
			sa:        "serviceaccount",
			namespace: "namespace",
			success:   false,
		},
		{
			name: "no match service account",
			identity: []string{
				spiffe.Identity{"cluster.local", "namespace", "bad"}.String(),
			},
			sa:        "serviceaccount",
			namespace: "namespace",
			success:   false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			con := &Connection{
				proxy:      &model.Proxy{ConfigNamespace: tt.namespace, Metadata: &model.NodeMetadata{ServiceAccount: tt.sa}},
				Identities: tt.identity,
			}
			if err := checkConnectionIdentity(con); (err == nil) != tt.success {
				t.Fatalf("expected success=%v, got err=%v", tt.success, err)
			}
		})
	}
}
