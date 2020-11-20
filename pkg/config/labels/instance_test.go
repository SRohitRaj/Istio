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

package labels_test

import (
	"math/rand"
	"strings"
	"testing"

	"istio.io/istio/pkg/config/labels"
)

func TestInstance(t *testing.T) {
	a := labels.Instance{"app": "a"}
	a1 := labels.Instance{"app": "a", "prod": "env"}

	if !labels.Instance(nil).SubsetOf(a) {
		t.Errorf("nil.SubsetOf({a}) => Got false")
	}

	if a.SubsetOf(nil) {
		t.Errorf("{a}.SubsetOf(nil) => Got true")
	}

	if a1.SubsetOf(a) {
		t.Errorf("%v.SubsetOf(%v) => Got true", a1, a)
	}
}

func TestInstanceValidate(t *testing.T) {
	cases := []struct {
		name  string
		tags  labels.Instance
		valid bool
	}{
		{
			name:  "empty tags",
			valid: true,
		},
		{
			name: "bad tag",
			tags: labels.Instance{"^": "^"},
		},
		{
			name:  "good tag",
			tags:  labels.Instance{"key": "value"},
			valid: true,
		},
		{
			name:  "good tag - empty value",
			tags:  labels.Instance{"key": ""},
			valid: true,
		},
		{
			name:  "good tag - DNS prefix",
			tags:  labels.Instance{"k8s.io/key": "value"},
			valid: true,
		},
		{
			name:  "good tag - subdomain DNS prefix",
			tags:  labels.Instance{"app.kubernetes.io/name": "value"},
			valid: true,
		},
		{
			name: "bad tag - empty key",
			tags: labels.Instance{"": "value"},
		},
		{
			name: "bad tag key 1",
			tags: labels.Instance{".key": "value"},
		},
		{
			name: "bad tag key 2",
			tags: labels.Instance{"key_": "value"},
		},
		{
			name: "bad tag key 3",
			tags: labels.Instance{"key$": "value"},
		},
		{
			name: "bad tag key - invalid DNS prefix",
			tags: labels.Instance{"istio./key": "value"},
		},
		{
			name: "bad tag value 1",
			tags: labels.Instance{"key": ".value"},
		},
		{
			name: "bad tag value 2",
			tags: labels.Instance{"key": "value_"},
		},
		{
			name: "bad tag value 3",
			tags: labels.Instance{"key": "value$"},
		},
	}
	for _, c := range cases {
		if got := c.tags.Validate(); (got == nil) != c.valid {
			t.Errorf("%s failed: got valid=%v but wanted valid=%v: %v", c.name, got == nil, c.valid, got)
		}
	}
}

func TestIsClusterName(t *testing.T) {
	var validClusterChars = []rune("abcdefghijklmnopqrstuvwxyz")

	clusterNameOfLength := func(n int) string {
		b := strings.Builder{}
		for i := 0; i < n; i++ {
			b.WriteRune(validClusterChars[rand.Intn(len(validClusterChars))])
		}
		return b.String()
	}

	cases := []struct {
		testName string
		in       string
		expected bool
	}{
		{
			testName: "valid",
			in:       "he.llo-wor_ld",
			expected: true,
		},
		{
			testName: "invalid first char",
			in:       "-a",
			expected: false,
		},
		{
			testName: "invalid last char",
			in:       "a-",
			expected: false,
		},
		{
			testName: "invalid middle char",
			in:       "a/b",
			expected: false,
		},
		{
			testName: "long",
			in:       clusterNameOfLength(labels.MaxClusterNameLength),
			expected: true,
		},
		{
			testName: "too long",
			in:       clusterNameOfLength(labels.MaxClusterNameLength + 1),
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			actual := labels.IsClusterName(c.in)
			if actual != c.expected {
				t.Fatalf("cluster name %s: got %t but want %t", c.in, actual, c.expected)
			}
		})
	}
}
