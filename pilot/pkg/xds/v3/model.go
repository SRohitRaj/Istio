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

package v3

import (
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	v2 "istio.io/istio/pilot/pkg/xds/v2"
)

const (
	ClusterType  = resource.ClusterType
	EndpointType = resource.EndpointType
	ListenerType = resource.ListenerType
	RouteType    = resource.RouteType
)

var (
	ListenerShortType = "LDS"
	RouteShortType    = "RDS"
	EndpointShortType = "EDS"
	ClusterShortType  = "CDS"
)

func GetShortType(typeURL string) string {
	switch typeURL {
	case v2.ClusterType, ClusterType:
		return ClusterShortType
	case v2.ListenerType, ListenerType:
		return ListenerShortType
	case v2.RouteType, RouteType:
		return RouteShortType
	case v2.EndpointType, EndpointType:
		return EndpointShortType
	default:
		return typeURL
	}
}
