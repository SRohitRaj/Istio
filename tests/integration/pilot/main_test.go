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

package pilot

import (
	"fmt"
	"strings"
	"testing"

	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/components/pilot"
	"istio.io/istio/tests/integration/multicluster"
)

var (
	i      istio.Instance
	pilots []pilot.Instance
)

// TestMain defines the entrypoint for pilot tests using a standard Istio installation.
// If a test requires a custom install it should go into its own package, otherwise it should go
// here to reuse a single install across tests.
func TestMain(m *testing.M) {
	framework.
		NewSuite(m).
		Setup(istio.Setup(&i, func(cfg *istio.Config) {
			// allow VMs to get discovery via ingress-gateway
			cfg.ControlPlaneValues = `
values:
  global:
    meshExpansion:
      enabled: true`
		})).
		Setup(multicluster.SetupPilots(&pilots)).
		Run()
}

func echoConfig(ns namespace.Instance, name string) echo.Config {
	return echo.Config{
		Service:   name,
		Namespace: ns,
		Ports: []echo.Port{
			{
				Name:     "http",
				Protocol: protocol.HTTP,
				// We use a port > 1024 to not require root
				InstancePort: 8090,
			},
		},
		Subsets: []echo.SubsetConfig{{}},
		Pilot:   pilots[0],
	}
}

// echoConfigForCluster builds an echo config with a setupFn for cluster-specific settings
// nameFmt should have one '%d' verb for the cluster index.
func echoConfigForCluster(ns namespace.Instance, nameFmt string) echo.Config {
	cfg := echoConfig(ns, nameFmt)
	cfg.SetupFn = setupForCluster
	return cfg
}

func setupForCluster(cfg *echo.Config) {
	if strings.Contains(cfg.Service, "%d") {
		cfg.Service = fmt.Sprintf(cfg.Service, cfg.Cluster.Index())
	}
	cfg.Pilot = pilots[cfg.Cluster.Index()]
}
