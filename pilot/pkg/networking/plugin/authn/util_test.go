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

package authn

import (
	"reflect"
	"testing"

	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/features"
)

func TestTrustDomainsForValidation(t *testing.T) {
	tests := []struct {
		name           string
		meshConfig     *meshconfig.MeshConfig
		skipValidation bool
		want           []string
	}{
		{
			name: "No duplicated trust domain in mesh config",
			meshConfig: &meshconfig.MeshConfig{
				TrustDomain:        "cluster.local",
				TrustDomainAliases: []string{"alias-1.domain", "some-other-alias-1.domain", "alias-2.domain"},
			},
			skipValidation: false,
			want:           []string{"cluster.local", "alias-1.domain", "some-other-alias-1.domain", "alias-2.domain"},
		},
		{
			name:           "Empty mesh config",
			meshConfig:     &meshconfig.MeshConfig{},
			skipValidation: false,
			want:           []string{},
		},
		{
			name:           "Skip Validation for mesh config",
			meshConfig:     &meshconfig.MeshConfig{},
			skipValidation: true,
			want:           []string(nil),
		},
		{
			name: "Sequential duplicated trust domains in mesh config",
			meshConfig: &meshconfig.MeshConfig{
				TrustDomain: "cluster.local",
				TrustDomainAliases: []string{
					"alias-1.domain", "alias-1.domain", "some-other-alias-1.domain", "alias-2.domain", "alias-2.domain",
				},
			},
			skipValidation: false,
			want:           []string{"cluster.local", "alias-1.domain", "some-other-alias-1.domain", "alias-2.domain"},
		},
		{
			name: "Mixed duplicated trust domains in mesh config",
			meshConfig: &meshconfig.MeshConfig{
				TrustDomain: "cluster.local",
				TrustDomainAliases: []string{
					"alias-1.domain", "cluster.local", "alias-2.domain", "some-other-alias-1.domain", "alias-2.domain", "alias-1.domain",
				},
			},
			skipValidation: false,
			want:           []string{"cluster.local", "alias-1.domain", "alias-2.domain", "some-other-alias-1.domain"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipValidation {
				features.SkipValidateTrustDomain = true
			} else {
				features.SkipValidateTrustDomain = false
			}
			if got := TrustDomainsForValidation(tt.meshConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trustDomainsForValidation() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
