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
package v2_test

import (
	"log"
	"testing"

	testenv "istio.io/istio/mixer/test/client/env"
	"istio.io/istio/pilot/pkg/bootstrap"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/test/env"
	"istio.io/istio/tests/util"

	"fmt"
	"time"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	ads "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	proto "github.com/gogo/protobuf/types"

	"istio.io/istio/pilot/pkg/proxy/envoy/v2"
)

type ExpectedResults struct {
	numOfNetworks  int
	localEndpoints []string
	remoteWeights  map[string]uint32
}

func TestSplitHorizonEds(t *testing.T) {
	initSplitHorizonTestEnv(t)

	verifySplitHorizonResponse(t, "network1", sidecarId("10.1.0.1", "app3"), "192.168.0.111", "80", ExpectedResults{
		numOfNetworks:  3,
		localEndpoints: []string{"10.1.0.1"},
		remoteWeights:  map[string]uint32{"192.168.0.222": 2, "192.168.0.333": 3},
	})

	verifySplitHorizonResponse(t, "network2", sidecarId("10.2.0.1", "app3"), "192.168.0.222", "80", ExpectedResults{
		numOfNetworks:  3,
		localEndpoints: []string{"10.2.0.1", "10.2.0.2"},
		remoteWeights:  map[string]uint32{"192.168.0.111": 1, "192.168.0.333": 3},
	})

	verifySplitHorizonResponse(t, "network3", sidecarId("10.3.0.1", "app3"), "192.168.0.333", "80", ExpectedResults{
		numOfNetworks:  3,
		localEndpoints: []string{"10.3.0.1", "10.3.0.2", "10.3.0.3"},
		remoteWeights:  map[string]uint32{"192.168.0.111": 1, "192.168.0.222": 2},
	})
}

func verifySplitHorizonResponse(t *testing.T, network string, sidecarId string, gatewayIp string, gatewayPort string, expected ExpectedResults) {
	edsstr, err := connectADS(util.MockPilotGrpcAddr)
	if err != nil {
		t.Fatal(err)
	}

	metadata := &proto.Struct{Fields: map[string]*proto.Value{
		"ISTIO_PROXY_VERSION":        {Kind: &proto.Value_StringValue{StringValue: "1.1"}},
		"ISTIO_NETWORK":              {Kind: &proto.Value_StringValue{StringValue: network}},
		"ISTIO_NETWORK_GATEWAY_IP":   {Kind: &proto.Value_StringValue{StringValue: gatewayIp}},
		"ISTIO_NETWORK_GATEWAY_PORT": {Kind: &proto.Value_StringValue{StringValue: gatewayPort}},
	}}

	err = sendCDSReqWithMetadata(sidecarId, metadata, edsstr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = adsReceive(edsstr, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	err = sendEDSReqWithMetadata([]string{"outbound|1080||service5.default.svc.cluster.local"}, sidecarId, metadata, edsstr)
	if err != nil {
		t.Fatal(err)
	}
	res, err := adsReceive(edsstr, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	cla, err := getLoadAssignment(res)
	if err != nil {
		t.Fatal(err)
	}
	eps := cla.Endpoints

	for i, ep := range eps {
		log.Printf("[%d] %v", i, ep)
	}

	if len(eps) != expected.numOfNetworks {
		t.Fatal(fmt.Errorf("Expecting %d endpoints but got %d", expected.numOfNetworks, len(eps)))
	}

	localEndpoints := eps[0].LbEndpoints
	if len(localEndpoints) != len(expected.localEndpoints) {
		t.Fatal(fmt.Errorf("Number of local endpoints should be %d but got %d", len(expected.localEndpoints), len(localEndpoints)))
	}

	for gwIP, weight := range expected.remoteWeights {
		found := false
		for _, ep := range eps {
			if ep.LbEndpoints[0].Endpoint.Address.GetSocketAddress().Address == gwIP && ep.LoadBalancingWeight.Value == weight {
				found = true
				break
			}
		}
		if !found {
			t.Fatal(fmt.Errorf("Couldn't find a gateway endpoint with IP %s and weight %d", gwIP, weight))
		}
	}
}

func initSplitHorizonTestEnv(t *testing.T) *bootstrap.Server {
	initMutex.Lock()
	defer initMutex.Unlock()
	testEnv = testenv.NewTestSetup(testenv.XDSTest, t)
	server := util.EnsureTestServer()
	pilotServer = server

	testEnv.Ports().PilotGrpcPort = uint16(util.MockPilotGrpcPort)
	testEnv.Ports().PilotHTTPPort = uint16(util.MockPilotHTTPPort)
	testEnv.IstioSrc = env.IstioSrc
	testEnv.IstioOut = env.IstioOut

	localIp = getLocalIP()

	networks1Lbls := map[string]string{
		"version":                    "v1",
		"ISTIO_NETWORK":              "network1",
		"ISTIO_NETWORK_GATEWAY_IP":   "192.168.0.111",
		"ISTIO_NETWORK_GATEWAY_PORT": "80"}
	networks2Lbls := map[string]string{
		"version":                    "v1",
		"ISTIO_NETWORK":              "network2",
		"ISTIO_NETWORK_GATEWAY_IP":   "192.168.0.222",
		"ISTIO_NETWORK_GATEWAY_PORT": "80"}
	networks3Lbls := map[string]string{
		"version":                    "v1",
		"ISTIO_NETWORK":              "network3",
		"ISTIO_NETWORK_GATEWAY_IP":   "192.168.0.333",
		"ISTIO_NETWORK_GATEWAY_PORT": "80"}

	// Explicit test service, in the v2 memory registry. Similar with mock.MakeService,
	// but easier to read.
	server.EnvoyXdsServer.MemRegistry.AddService("service5.default.svc.cluster.local", &model.Service{
		Hostname: "service5.default.svc.cluster.local",
		Address:  "10.10.0.1",
		Ports:    testPorts(0),
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("service5.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: "10.1.0.1",
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           networks1Lbls,
		AvailabilityZone: "az",
	})

	// Explicit test service, in the v2 memory registry. Similar with mock.MakeService,
	// but easier to read.
	server.EnvoyXdsServer.MemRegistry.AddService("service5.default.svc.cluster.local", &model.Service{
		Hostname: "service5.default.svc.cluster.local",
		Address:  "10.20.0.1",
		Ports:    testPorts(0),
	})

	server.EnvoyXdsServer.MemRegistry.AddInstance("service5.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: "10.2.0.1",
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           networks2Lbls,
		AvailabilityZone: "az",
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("service5.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: "10.2.0.2",
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           networks2Lbls,
		AvailabilityZone: "az",
	})

	// Explicit test service, in the v2 memory registry. Similar with mock.MakeService,
	// but easier to read.
	server.EnvoyXdsServer.MemRegistry.AddService("service5.default.svc.cluster.local", &model.Service{
		Hostname: "service5.default.svc.cluster.local",
		Address:  "10.30.0.1",
		Ports:    testPorts(0),
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("service5.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: "10.3.0.1",
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           networks3Lbls,
		AvailabilityZone: "az",
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("service5.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: "10.3.0.2",
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           networks3Lbls,
		AvailabilityZone: "az",
	})
	server.EnvoyXdsServer.MemRegistry.AddInstance("service5.default.svc.cluster.local", &model.ServiceInstance{
		Endpoint: model.NetworkEndpoint{
			Address: "10.3.0.3",
			Port:    2080,
			ServicePort: &model.Port{
				Name:     "http-main",
				Port:     1080,
				Protocol: model.ProtocolHTTP,
			},
		},
		Labels:           networks3Lbls,
		AvailabilityZone: "az",
	})

	// Update cache
	server.EnvoyXdsServer.ClearCacheFunc()()

	return server
}

func sendCDSReqWithMetadata(node string, metadata *proto.Struct, edsstr ads.AggregatedDiscoveryService_StreamAggregatedResourcesClient) error {
	err := edsstr.Send(&xdsapi.DiscoveryRequest{
		ResponseNonce: time.Now().String(),
		Node: &core.Node{
			Id:       node,
			Metadata: metadata,
		},
		TypeUrl: v2.ClusterType})
	if err != nil {
		return fmt.Errorf("CDS request failed: %s", err)
	}

	return nil
}

func sendEDSReqWithMetadata(clusters []string, node string, metadata *proto.Struct,
	edsstr ads.AggregatedDiscoveryService_StreamAggregatedResourcesClient) error {
	err := edsstr.Send(&xdsapi.DiscoveryRequest{
		ResponseNonce: time.Now().String(),
		Node: &core.Node{
			Id:       node,
			Metadata: metadata,
		},
		TypeUrl:       v2.EndpointType,
		ResourceNames: clusters,
	})
	if err != nil {
		return fmt.Errorf("EDS request failed: %s", err)
	}

	return nil
}
