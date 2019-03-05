//  Copyright 2019 Istio Authors
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

package pilot

import (
	"fmt"
	"testing"
	"time"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"istio.io/istio/pkg/test/framework2/components/environment"
	"istio.io/istio/pkg/test/framework2/components/environment/native"
	"istio.io/istio/pkg/test/framework2/components/galley"
	"istio.io/istio/pkg/test/framework2/resource"
)

// Instance of Pilot
type Instance interface {
	CallDiscovery(req *xdsapi.DiscoveryRequest) (*xdsapi.DiscoveryResponse, error)
	StartDiscovery(req *xdsapi.DiscoveryRequest) error
	WatchDiscovery(duration time.Duration, accept func(*xdsapi.DiscoveryResponse) (bool, error)) error
}

// Structured config for the Pilot component
type Config struct {
	fmt.Stringer
	// If set then pilot takes a dependency on the referenced Galley instance
	Galley galley.Instance
}

// New returns a new Galley instance.
func New(c resource.Context, config *Config) (Instance, error) {
	switch c.Environment().Name() {
	case environment.Native:
		return newNative(c, c.Environment().(*native.Environment), config)
	default:
		return nil, environment.UnsupportedEnvironment(c.Environment().Name())
	}
}

// NewOrFail returns a new Galley instance, or fails.
func NewOrFail(t *testing.T, c resource.Context, config *Config) Instance {
	i, err := New(c, config)
	if err != nil {
		t.Fatalf("Error creating Galley: %v", err)
	}
	return i
}
