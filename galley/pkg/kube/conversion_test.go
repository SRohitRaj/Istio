//  Copyright 2018 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package kube

import (
	"reflect"
	"testing"

	"github.com/gogo/protobuf/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"istio.io/istio/galley/pkg/runtime/resource"
)

func TestToProto(t *testing.T) {
	cr := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":            "foo",
				"resourceVersion": "rv",
			},
			"spec": map[string]interface{}{},
		},
	}

	s := ResourceSpec{
		Target: resource.Info{
			Kind: "google.protobuf.Empty",
		},
	}

	p, err := toProto(s, cr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var expected = &types.Empty{}
	if !reflect.DeepEqual(p, expected) {
		t.Fatalf("Mismatch\nExpected:\n%+v\nActual:\n%+v\n", expected, p)
	}
}
