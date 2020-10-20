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

package crd

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// IstioKind is the generic Kubernetes API object wrapper
type IstioKind struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               map[string]interface{} `json:"spec"`
	Status             map[string]interface{} `json:"status,omitempty"`
}

// GetSpec from a wrapper
func (in *IstioKind) GetSpec() map[string]interface{} {
	return in.Spec
}

// GetStatus from a wrapper
func (in *IstioKind) GetStatus() map[string]interface{} {
	return in.Status
}

// GetObjectMeta from a wrapper
func (in *IstioKind) GetObjectMeta() meta_v1.ObjectMeta {
	return in.ObjectMeta
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IstioKind) DeepCopyInto(out *IstioKind) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IstioKind.
func (in *IstioKind) DeepCopy() *IstioKind {
	if in == nil {
		return nil
	}
	out := new(IstioKind)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IstioKind) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}

	return nil
}

// IstioObject is a k8s wrapper interface for config objects
type IstioObject interface {
	runtime.Object
	GetSpec() map[string]interface{}
	GetStatus() map[string]interface{}
	GetObjectMeta() meta_v1.ObjectMeta
}
