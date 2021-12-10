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

package v1alpha1

import (
	"encoding/json"

	proto "google.golang.org/protobuf/proto"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"istio.io/pkg/log"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (m *IstioOperator) DeepCopyInto(out proto.Message) {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if err = json.Unmarshal(bytes, out); err != nil {
		log.Error(err.Error())
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IstioOperator.
func (m *IstioOperator) DeepCopy() *IstioOperator {
	if m == nil {
		return nil
	}
	out := new(IstioOperator)
	m.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (m *IstioOperator) DeepCopyObject() runtime.Object {
	if c := m.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// TODO: needs to be implemented or generated.
func (m *IstioOperator) GetObjectKind() schema.ObjectKind {
	return EmptyObjectKind
}

// IstioOperatorList contains a list of IstioOperator
type IstioOperatorList struct {
	v1.TypeMeta `json:",inline"`
	v1.ListMeta `json:"metadata,omitempty"`
	Items       []IstioOperator `json:"items"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IstioOperatorList) DeepCopyInto(out *IstioOperatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IstioOperator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IstioOperatorList.
func (in *IstioOperatorList) DeepCopy() *IstioOperatorList {
	if in == nil {
		return nil
	}
	out := new(IstioOperatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IstioOperatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// EmptyObjectKind implements the ObjectKind interface as a noop
var EmptyObjectKind = emptyObjectKind{}

type emptyObjectKind struct{}

// SetGroupVersionKind implements the ObjectKind interface
func (emptyObjectKind) SetGroupVersionKind(gvk schema.GroupVersionKind) {}

// GroupVersionKind implements the ObjectKind interface
func (emptyObjectKind) GroupVersionKind() schema.GroupVersionKind { return schema.GroupVersionKind{} }
