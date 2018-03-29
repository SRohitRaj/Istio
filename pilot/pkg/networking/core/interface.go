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

package core

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core/v1alpha3"
)

// Networking represents the interfaces to be implemented by code that generates xDS responses
type Networking struct {
	// BuildListeners returns the list of listeners for the given proxy. This is the LDS output
	// Internally, the computation will be optimized to ensure that listeners are computed only
	// once and shared across multiple invocations of this function.
	BuildListeners func(env model.Environment, node model.Proxy) ([]*v2.Listener, error)

	// BuildClusters returns the list of clusters for the given proxy. This is the CDS output
	BuildClusters func(env model.Environment, node model.Proxy) ([]*v2.Cluster, error)

	// BuildRoutes returns the list of routes for the given proxy. This is the RDS output
	BuildRoutes func(env model.Environment, node model.Proxy, routeName string) ([]*v2.RouteConfiguration, error)
}


// NewNetworkConfiguration creates a new instance of the dataplane configuration generator
func NewNetworkConfiguration() *Networking {
	return &Networking{
		BuildListeners: v1alpha3.BuildListeners,
		BuildClusters:  v1alpha3.BuildClusters,
	}
}
