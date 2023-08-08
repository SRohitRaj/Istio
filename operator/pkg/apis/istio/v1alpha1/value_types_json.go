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

// Code generated by protoc-gen-golang. DO NOT EDIT.
// source: operator/v1alpha1/operator.proto

// Configuration affecting Istio control plane installation version and shape.

package v1alpha1

import (
	"encoding/json"

	github_com_golang_protobuf_jsonpb "github.com/golang/protobuf/jsonpb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var _ github_com_golang_protobuf_jsonpb.JSONPBUnmarshaler = &IntOrString{}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (this *IntOrString) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		this.Type = int64(intstr.String)
		var s string
		err := json.Unmarshal(value, &s)
		if err != nil {
			return err
		}
		this.StrVal = &wrappers.StringValue{Value: s}
		return nil
	}
	this.Type = int64(intstr.Int)
	var s int32
	err := json.Unmarshal(value, &s)
	if err != nil {
		return err
	}
	this.IntVal = &wrappers.Int32Value{Value: s}
	return nil
}

func (this *IntOrString) MarshalJSONPB(_ *github_com_golang_protobuf_jsonpb.Marshaler) ([]byte, error) {
	return this.MarshalJSON()
}

func (this *IntOrString) MarshalJSON() ([]byte, error) {
	if this.IntVal != nil {
		return json.Marshal(this.IntVal.GetValue())
	}
	return json.Marshal(this.StrVal.GetValue())
}

func (this *IntOrString) UnmarshalJSONPB(_ *github_com_golang_protobuf_jsonpb.Unmarshaler, value []byte) error {
	return this.UnmarshalJSON(value)
}

func (this *IntOrString) ToKubernetes() intstr.IntOrString {
	if this.IntVal != nil {
		return intstr.FromInt32(this.GetIntVal().GetValue())
	}
	return intstr.FromString(this.GetStrVal().GetValue())
}
