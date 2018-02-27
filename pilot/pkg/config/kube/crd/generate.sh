#!/bin/bash
#
# Copyright 2017 Istio Authors. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
################################################################################
#
# Generate structs for CRDs (see comments below)

set -o errexit
set -o nounset
set -o pipefail

cat << EOF
// Code generated by generate.sh. DO NOT EDIT!
// Sources: adapter/config/crd/config.go
// Output: adapter/config/crd/types.go

package crd

// This file contains Go definitions for Custom Resource Definition kinds
// to adhere to the idiomatic use of k8s API machinery.
// These definitions are synthesized from Istio configuration type descriptors
// as declared in the Pilot config model.

import (
    meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"

    "istio.io/istio/pilot/pkg/model"
)

type schemaType struct {
	schema     model.ProtoSchema
	object     IstioObject
	collection IstioObjectList
}

var knownTypes = map[string]schemaType{
EOF

CRDS="MockConfig RouteRule V1alpha2RouteRule IngressRule Gateway EgressRule ExternalService DestinationPolicy DestinationRule HTTPAPISpec HTTPAPISpecBinding QuotaSpec QuotaSpecBinding EndUserAuthenticationPolicySpec EndUserAuthenticationPolicySpecBinding"

for crd in $CRDS; do
cat << EOF
    model.$crd.Type: {
        schema: model.$crd,
        object: &${crd}{
            TypeMeta: meta_v1.TypeMeta{
                Kind:       "${crd}",
                APIVersion: ResourceGroup(&model.$crd) + "/" + model.$crd.Version,
            },
        },
        collection: &${crd}List{},
    },
EOF

done

cat <<EOF
}
EOF

for crd in $CRDS; do
  sed -e "1,20d;s/IstioKind/$crd/g" $1
done
