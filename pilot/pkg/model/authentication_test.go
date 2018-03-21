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

package model

import (
	"reflect"
	"testing"

	authn1 "istio.io/api/authentication/v1alpha1"
	authn "istio.io/api/authentication/v1alpha2"
	meshconfig "istio.io/api/mesh/v1alpha1"
	networking "istio.io/api/networking/v1alpha3"
)

func TestRequireTls(t *testing.T) {
	cases := []struct {
		in       authn1.Policy
		expected bool
	}{
		{
			in:       authn1.Policy{},
			expected: false,
		},
		{
			in: authn1.Policy{
				Peers: []*authn1.PeerAuthenticationMethod{{
					Params: &authn1.PeerAuthenticationMethod_Mtls{},
				}},
			},
			expected: true,
		},
		{
			in: authn1.Policy{
				Peers: []*authn1.PeerAuthenticationMethod{{
					Params: &authn1.PeerAuthenticationMethod_Jwt{},
				},
					{
						Params: &authn1.PeerAuthenticationMethod_Mtls{},
					},
				},
			},
			expected: true,
		},
	}
	for _, c := range cases {
		if got := RequireTLS(&c.in); got != c.expected {
			t.Errorf("requireTLS(%v): got(%v) != want(%v)\n", c.in, got, c.expected)
		}
	}
}

func TestLegacyAuthenticationPolicyToPolicy(t *testing.T) {
	cases := []struct {
		in       meshconfig.AuthenticationPolicy
		expected *authn1.Policy
	}{
		{
			in: meshconfig.AuthenticationPolicy_MUTUAL_TLS,
			expected: &authn1.Policy{
				Peers: []*authn1.PeerAuthenticationMethod{{
					Params: &authn1.PeerAuthenticationMethod_Mtls{},
				}},
			},
		},
		{
			in:       meshconfig.AuthenticationPolicy_NONE,
			expected: nil,
		},
	}

	for _, c := range cases {
		if got := legacyAuthenticationPolicyToPolicy(c.in); !reflect.DeepEqual(got, c.expected) {
			t.Errorf("legacyAuthenticationPolicyToPolicy(%v): got(%#v) != want(%#v)\n", c.in, got, c.expected)
		}
	}
}

func TestAuthnPortSelectorToNetworkingPortSelector(t *testing.T) {
	cases := []struct {
		in       *authn.PortSelector
		expected *networking.PortSelector
	}{
		{
			in:       nil,
			expected: nil,
		},
		{
			in:       &authn.PortSelector{},
			expected: &networking.PortSelector{},
		},
		{
			in: &authn.PortSelector{
				Port: &authn.PortSelector_Name{
					Name: "http",
				},
			},
			expected: &networking.PortSelector{
				Port: &networking.PortSelector_Name{
					Name: "http",
				},
			},
		},
		{
			in: &authn.PortSelector{
				Port: &authn.PortSelector_Number{
					Number: 8888,
				},
			},
			expected: &networking.PortSelector{
				Port: &networking.PortSelector_Number{
					Number: 8888,
				},
			},
		},
	}
	for _, c := range cases {
		if got := AuthnPortSelectorToNetworkingPortSelector(c.in); !reflect.DeepEqual(c.expected, got) {
			t.Errorf("AuthnPortSelectorToNetworkingPortSelector(%v): got(%v) != want(%v)\n", c.in, got, c.expected)
		}
	}
}
