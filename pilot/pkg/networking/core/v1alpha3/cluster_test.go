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

package v1alpha3_test

import (
	"testing"
	"time"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/onsi/gomega"

	meshconfig "istio.io/api/mesh/v1alpha1"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/model"
	core "istio.io/istio/pilot/pkg/networking/core/v1alpha3"
	"istio.io/istio/pilot/pkg/networking/core/v1alpha3/fakes"
	"istio.io/istio/pilot/pkg/networking/plugin"
)

func TestBuildGatewayClustersWithRingHashLb(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	configgen := core.NewConfigGenerator([]plugin.Plugin{})
	proxy := model.Proxy{
		ClusterID: "some-cluster-id",
		Type:      model.Router,
		IPAddress: "6.6.6.6",
		Domain:    "default.example.org",
		Metadata:  make(map[string]string),
	}

	env := buildEnv()

	clusters, err := configgen.BuildClusters(env, proxy)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	g.Expect(len(clusters)).To(gomega.Equal(2))

	cluster := clusters[0]
	g.Expect(cluster.LbPolicy).To(gomega.Equal(v2.Cluster_RING_HASH))
	g.Expect(cluster.GetRingHashLbConfig().GetMinimumRingSize().GetValue()).To(gomega.Equal(uint64(2)))
	g.Expect(cluster.Name).To(gomega.Equal("outbound|8080||*.example.org"))
	g.Expect(cluster.Type).To(gomega.Equal(v2.Cluster_EDS))
	g.Expect(cluster.ConnectTimeout).To(gomega.Equal(time.Duration(10000000001)))
}

func buildEnv() model.Environment {
	serviceDiscovery := &fakes.ServiceDiscovery{}

	serviceDiscovery.ServicesReturns([]*model.Service{
		{
			Hostname:    "*.example.org",
			Address:     "1.1.1.1",
			ClusterVIPs: make(map[string]string),
			Ports: model.PortList{
				&model.Port{
					Name:     "default",
					Port:     8080,
					Protocol: model.ProtocolHTTP,
				},
			},
		},
	}, nil)

	meshConfig := &meshconfig.MeshConfig{
		ConnectTimeout: &duration.Duration{
			Seconds: 10,
			Nanos:   1,
		},
	}

	configStore := &fakes.IstioConfigStore{}
	configStore.DestinationRuleReturns(
		&model.Config{
			ConfigMeta: model.ConfigMeta{
				Type:    model.DestinationRule.Type,
				Version: model.DestinationRule.Version,
				Name:    "acme",
			},
			Spec: &networking.DestinationRule{
				Host: "*.example.org",
				TrafficPolicy: &networking.TrafficPolicy{
					LoadBalancer: &networking.LoadBalancerSettings{
						LbPolicy: &networking.LoadBalancerSettings_ConsistentHash{
							ConsistentHash: &networking.LoadBalancerSettings_ConsistentHashLB{
								HttpHeader:      "some-header",
								MinimumRingSize: 2,
							},
						},
					},
				},
			},
		},
	)

	env := model.Environment{
		ServiceDiscovery: serviceDiscovery,
		ServiceAccounts:  &fakes.ServiceAccounts{},
		IstioConfigStore: configStore,
		Mesh:             meshConfig,
		MixerSAN:         []string{},
	}

	return env
}
