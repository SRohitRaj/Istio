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

package krt

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	acmetav1 "k8s.io/client-go/applyconfigurations/meta/v1"

	"istio.io/istio/pkg/ptr"
)

// GetKey returns the key for the provided object.
// If there is none, this will panic.
func GetKey[O any](a O) Key[O] {
	if k, ok := tryGetKey[O](a); ok {
		return k
	}
	// Allow pointer receiver as well
	if k, ok := tryGetKey[*O](&a); ok {
		return Key[O](k)
	}
	panic(fmt.Sprintf("Cannot get Key, got %T", a))
	return ""
}

// Named is a convenience struct. It is ideal to be embedded into a type that has a name and namespace,
// and will automatically implement the various interfaces to return the name, namespace, and a key based on these two.
type Named struct {
	Name, Namespace string
}

// NewNamed builds a Named object from a Kubernetes object type.
func NewNamed(o metav1.Object) Named {
	return Named{Name: o.GetName(), Namespace: o.GetNamespace()}
}

func (n Named) ResourceName() string {
	return n.Namespace + "/" + n.Name
}

func (n Named) GetName() string {
	return n.Name
}

func (n Named) GetNamespace() string {
	return n.Namespace
}

func GetApplyConfigKey[O any](a O) *Key[O] {
	val := reflect.ValueOf(a)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}

	specField := val.FieldByName("ObjectMetaApplyConfiguration")
	if !specField.IsValid() {
		return nil
	}
	meta := specField.Interface().(*acmetav1.ObjectMetaApplyConfiguration)
	if meta.Namespace != nil && len(*meta.Namespace) > 0 {
		return ptr.Of(Key[O](*meta.Namespace + "/" + *meta.Name))
	}
	return ptr.Of(Key[O](*meta.Name))
}

// keyFunc is the internal API key function that returns "namespace"/"name" or
// "name" if "namespace" is empty
func keyFunc(name, namespace string) string {
	if len(namespace) == 0 {
		return name
	}
	return namespace + "/" + name
}
