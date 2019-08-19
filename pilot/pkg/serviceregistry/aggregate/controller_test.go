// Copyright 2017 Istio Authors
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

package aggregate_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry"
	"istio.io/istio/pilot/pkg/serviceregistry/aggregate"
	"istio.io/istio/pilot/pkg/serviceregistry/aggregate/mock"
	"istio.io/istio/pilot/pkg/serviceregistry/memory"
	"istio.io/istio/pkg/config/host"
	"istio.io/istio/pkg/config/labels"
	"istio.io/istio/pkg/config/protocol"
)

// MockController specifies a mock Controller for testing
type MockController struct{}

func (c *MockController) AppendServiceHandler(f func(*model.Service, model.Event)) error {
	return nil
}

func (c *MockController) AppendInstanceHandler(f func(*model.ServiceInstance, model.Event)) error {
	return nil
}

func (c *MockController) Run(<-chan struct{}) {}

var discovery1 *memory.ServiceDiscovery
var discovery2 *memory.ServiceDiscovery

func buildMockController() *aggregate.Controller {
	discovery1 = memory.NewDiscovery(
		map[host.Name]*model.Service{
			memory.HelloService.Hostname:   memory.HelloService,
			memory.ExtHTTPService.Hostname: memory.ExtHTTPService,
		}, 2)

	discovery2 = memory.NewDiscovery(
		map[host.Name]*model.Service{
			memory.WorldService.Hostname:    memory.WorldService,
			memory.ExtHTTPSService.Hostname: memory.ExtHTTPSService,
		}, 2)

	registry1 := aggregate.Registry{
		Name:             serviceregistry.ServiceRegistry("mockAdapter1"),
		ServiceDiscovery: discovery1,
		Controller:       &MockController{},
	}

	registry2 := aggregate.Registry{
		Name:             serviceregistry.ServiceRegistry("mockAdapter2"),
		ServiceDiscovery: discovery2,
		Controller:       &MockController{},
	}

	ctls := aggregate.NewController()
	ctls.AddRegistry(registry1)
	ctls.AddRegistry(registry2)

	return ctls
}

func buildMockControllerForMultiCluster() *aggregate.Controller {
	discovery1 = memory.NewDiscovery(
		map[host.Name]*model.Service{
			memory.HelloService.Hostname: memory.MakeService("hello.default.svc.cluster.local", "10.1.1.0"),
		}, 2)

	discovery2 = memory.NewDiscovery(
		map[host.Name]*model.Service{
			memory.HelloService.Hostname: memory.MakeService("hello.default.svc.cluster.local", "10.1.2.0"),
			memory.WorldService.Hostname: memory.WorldService,
		}, 2)

	registry1 := aggregate.Registry{
		Name:             serviceregistry.ServiceRegistry("mockAdapter1"),
		ClusterID:        "cluster-1",
		ServiceDiscovery: discovery1,
		Controller:       &MockController{},
	}

	registry2 := aggregate.Registry{
		Name:             serviceregistry.ServiceRegistry("mockAdapter2"),
		ClusterID:        "cluster-2",
		ServiceDiscovery: discovery2,
		Controller:       &MockController{},
	}

	ctls := aggregate.NewController()
	ctls.AddRegistry(registry1)
	ctls.AddRegistry(registry2)

	return ctls
}

func TestServicesError(t *testing.T) {
	aggregateCtl := buildMockController()

	discovery1.ServicesError = errors.New("mock Services() error")

	// List Services from aggregate controller
	_, err := aggregateCtl.Services()
	if err == nil {
		t.Fatal("Aggregate controller should return error if one discovery client experience error")
	}
}
func TestServicesForMultiCluster(t *testing.T) {
	aggregateCtl := buildMockControllerForMultiCluster()
	// List Services from aggregate controller
	services, err := aggregateCtl.Services()
	if err != nil {
		t.Fatalf("Services() encountered unexpected error: %v", err)
	}

	// Set up ground truth hostname values
	serviceMap := map[host.Name]bool{
		memory.HelloService.Hostname: false,
		memory.WorldService.Hostname: false,
	}

	svcCount := 0
	// Compare return value to ground truth
	for _, svc := range services {
		if counted, existed := serviceMap[svc.Hostname]; existed && !counted {
			svcCount++
			serviceMap[svc.Hostname] = true
		}
	}

	if svcCount != len(serviceMap) {
		t.Fatalf("Service map expected size %d, actual %v", svcCount, serviceMap)
	}

	//Now verify ClusterVIPs for each service
	ClusterVIPs := map[host.Name]map[string]string{
		memory.HelloService.Hostname: {
			"cluster-1": "10.1.1.0",
			"cluster-2": "10.1.2.0",
		},
		memory.WorldService.Hostname: {
			"cluster-2": "10.2.0.0",
		},
	}
	for _, svc := range services {
		if !reflect.DeepEqual(svc.ClusterVIPs, ClusterVIPs[svc.Hostname]) {
			t.Fatalf("Service %s ClusterVIPs actual %v, expected %v", svc.Hostname, svc.ClusterVIPs, ClusterVIPs[svc.Hostname])
		}
	}
	t.Logf("Return service ClusterVIPs match ground truth")
}

func TestServices(t *testing.T) {
	aggregateCtl := buildMockController()
	// List Services from aggregate controller
	services, err := aggregateCtl.Services()

	// Set up ground truth hostname values
	serviceMap := map[host.Name]bool{
		memory.HelloService.Hostname:    false,
		memory.ExtHTTPService.Hostname:  false,
		memory.WorldService.Hostname:    false,
		memory.ExtHTTPSService.Hostname: false,
	}

	if err != nil {
		t.Fatalf("Services() encountered unexpected error: %v", err)
	}

	svcCount := 0
	// Compare return value to ground truth
	for _, svc := range services {
		if counted, existed := serviceMap[svc.Hostname]; existed && !counted {
			svcCount++
			serviceMap[svc.Hostname] = true
		}
	}

	if svcCount != len(serviceMap) {
		t.Fatal("Return services does not match ground truth")
	}
}

func TestGetService(t *testing.T) {
	aggregateCtl := buildMockController()

	// Get service from mockAdapter1
	svc, err := aggregateCtl.GetService(memory.HelloService.Hostname)
	if err != nil {
		t.Fatalf("GetService() encountered unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("Fail to get service")
	}
	if svc.Hostname != memory.HelloService.Hostname {
		t.Fatal("Returned service is incorrect")
	}

	// Get service from mockAdapter2
	svc, err = aggregateCtl.GetService(memory.WorldService.Hostname)
	if err != nil {
		t.Fatalf("GetService() encountered unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("Fail to get service")
	}
	if svc.Hostname != memory.WorldService.Hostname {
		t.Fatal("Returned service is incorrect")
	}
}

func TestGetServiceError(t *testing.T) {
	aggregateCtl := buildMockController()

	discovery1.GetServiceError = errors.New("mock GetService() error")

	// Get service from client with error
	svc, err := aggregateCtl.GetService(memory.HelloService.Hostname)
	if err == nil {
		fmt.Println(svc)
		t.Fatal("Aggregate controller should return error if one discovery client experiences " +
			"error and no service is found")
	}
	if svc != nil {
		t.Fatal("GetService() should return nil if no service found")
	}

	// Get service from client without error
	svc, err = aggregateCtl.GetService(memory.WorldService.Hostname)
	if err != nil {
		t.Fatal("Aggregate controller should not return error if service is found")
	}
	if svc == nil {
		t.Fatal("Fail to get service")
	}
	if svc.Hostname != memory.WorldService.Hostname {
		t.Fatal("Returned service is incorrect")
	}
}

func TestGetProxyServiceInstances(t *testing.T) {
	aggregateCtl := buildMockController()

	// Get Instances from mockAdapter1
	instances, err := aggregateCtl.GetProxyServiceInstances(&model.Proxy{IPAddresses: []string{memory.HelloInstanceV0}})
	if err != nil {
		t.Fatalf("GetProxyServiceInstances() encountered unexpected error: %v", err)
	}
	if len(instances) != 6 {
		t.Fatalf("Returned GetProxyServiceInstances' amount %d is not correct", len(instances))
	}
	for _, inst := range instances {
		if inst.Service.Hostname != memory.HelloService.Hostname {
			t.Fatal("Returned Instance is incorrect")
		}
	}

	// Get Instances from mockAdapter2
	instances, err = aggregateCtl.GetProxyServiceInstances(&model.Proxy{IPAddresses: []string{memory.MakeIP(memory.WorldService, 1)}})
	if err != nil {
		t.Fatalf("GetProxyServiceInstances() encountered unexpected error: %v", err)
	}
	if len(instances) != 6 {
		t.Fatalf("Returned GetProxyServiceInstances' amount %d is not correct", len(instances))
	}
	for _, inst := range instances {
		if inst.Service.Hostname != memory.WorldService.Hostname {
			t.Fatal("Returned Instance is incorrect")
		}
	}
}

func TestGetProxyServiceInstancesError(t *testing.T) {
	aggregateCtl := buildMockController()

	discovery1.GetProxyServiceInstancesError = errors.New("mock GetProxyServiceInstances() error")

	// Get Instances from client with error
	instances, err := aggregateCtl.GetProxyServiceInstances(&model.Proxy{IPAddresses: []string{memory.HelloInstanceV0}})
	if err == nil {
		t.Fatal("Aggregate controller should return error if one discovery client experiences " +
			"error and no instances are found")
	}
	if len(instances) != 0 {
		t.Fatal("GetProxyServiceInstances() should return no instances is client experiences error")
	}

	// Get Instances from client without error
	instances, err = aggregateCtl.GetProxyServiceInstances(&model.Proxy{IPAddresses: []string{memory.MakeIP(memory.WorldService, 1)}})
	if err != nil {
		t.Fatal("Aggregate controller should not return error if instances are found")
	}
	if len(instances) != 6 {
		t.Fatalf("Returned GetProxyServiceInstances' amount %d is not correct", len(instances))
	}
	for _, inst := range instances {
		if inst.Service.Hostname != memory.WorldService.Hostname {
			t.Fatal("Returned Instance is incorrect")
		}
	}
}

func TestInstances(t *testing.T) {
	aggregateCtl := buildMockController()

	// Get Instances from mockAdapter1
	instances, err := aggregateCtl.InstancesByPort(memory.HelloService,
		80,
		labels.Collection{})
	if err != nil {
		t.Fatalf("Instances() encountered unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatal("Returned wrong number of instances from controller")
	}
	for _, instance := range instances {
		if instance.Service.Hostname != memory.HelloService.Hostname {
			t.Fatal("Returned instance's hostname does not match desired value")
		}
		if _, ok := instance.Service.Ports.Get(memory.PortHTTPName); !ok {
			t.Fatal("Returned instance does not contain desired port")
		}
	}

	// Get Instances from mockAdapter2
	instances, err = aggregateCtl.InstancesByPort(memory.WorldService,
		80,
		labels.Collection{})
	if err != nil {
		t.Fatalf("Instances() encountered unexpected error: %v", err)
	}
	if len(instances) != 2 {
		t.Fatal("Returned wrong number of instances from controller")
	}
	for _, instance := range instances {
		if instance.Service.Hostname != memory.WorldService.Hostname {
			t.Fatal("Returned instance's hostname does not match desired value")
		}
		if _, ok := instance.Service.Ports.Get(memory.PortHTTPName); !ok {
			t.Fatal("Returned instance does not contain desired port")
		}
	}
}

func TestInstancesError(t *testing.T) {
	aggregateCtl := buildMockController()

	discovery1.InstancesError = errors.New("mock Instances() error")

	// Get Instances from client with error
	instances, err := aggregateCtl.InstancesByPort(memory.HelloService,
		80,
		labels.Collection{})
	if err == nil {
		t.Fatal("Aggregate controller should return error if one discovery client experiences " +
			"error and no instances are found")
	}
	if len(instances) != 0 {
		t.Fatal("Returned wrong number of instances from controller")
	}

	// Get Instances from client without error
	instances, err = aggregateCtl.InstancesByPort(memory.WorldService,
		80,
		labels.Collection{})
	if err != nil {
		t.Fatalf("Instances() should not return error is instances are found: %v", err)
	}
	if len(instances) != 2 {
		t.Fatal("Returned wrong number of instances from controller")
	}
	for _, instance := range instances {
		if instance.Service.Hostname != memory.WorldService.Hostname {
			t.Fatal("Returned instance's hostname does not match desired value")
		}
		if _, ok := instance.Service.Ports.Get(memory.PortHTTPName); !ok {
			t.Fatal("Returned instance does not contain desired port")
		}
	}
}

func TestGetIstioServiceAccounts(t *testing.T) {
	aggregateCtl := buildMockController()

	// Get accounts from mockAdapter1
	accounts := aggregateCtl.GetIstioServiceAccounts(memory.HelloService, []int{})
	expected := make([]string, 0)

	if len(accounts) != len(expected) {
		t.Fatal("Incorrect account result returned")
	}

	for i := 0; i < len(accounts); i++ {
		if accounts[i] != expected[i] {
			t.Fatal("Returned account result does not match expected one")
		}
	}

	// Get accounts from mockAdapter2
	accounts = aggregateCtl.GetIstioServiceAccounts(memory.WorldService, []int{})
	expected = []string{
		"spiffe://cluster.local/ns/default/sa/serviceaccount1",
		"spiffe://cluster.local/ns/default/sa/serviceaccount2",
	}

	if len(accounts) != len(expected) {
		t.Fatal("Incorrect account result returned")
	}

	for i := 0; i < len(accounts); i++ {
		if accounts[i] != expected[i] {
			t.Fatal("Returned account result does not match expected one", accounts[i], expected[i])
		}
	}
}

func TestManagementPorts(t *testing.T) {
	aggregateCtl := buildMockController()
	expected := model.PortList{{
		Name:     "http",
		Port:     3333,
		Protocol: protocol.HTTP,
	}, {
		Name:     "custom",
		Port:     9999,
		Protocol: protocol.TCP,
	}}

	// Get management ports from mockAdapter1
	ports := aggregateCtl.ManagementPorts(memory.HelloInstanceV0)
	if len(ports) != 2 {
		t.Fatal("Returned wrong number of ports from controller")
	}
	for i := 0; i < len(ports); i++ {
		if ports[i].Name != expected[i].Name || ports[i].Port != expected[i].Port ||
			ports[i].Protocol != expected[i].Protocol {
			t.Fatal("Returned management ports result does not match expected one")
		}
	}

	// Get management ports from mockAdapter2
	ports = aggregateCtl.ManagementPorts(memory.MakeIP(memory.WorldService, 0))
	if len(ports) != len(expected) {
		t.Fatal("Returned wrong number of ports from controller")
	}
	for i := 0; i < len(ports); i++ {
		if ports[i].Name != expected[i].Name || ports[i].Port != expected[i].Port ||
			ports[i].Protocol != expected[i].Protocol {
			t.Fatal("Returned management ports result does not match expected one")
		}
	}
}

func TestAddRegistry(t *testing.T) {

	registries := []aggregate.Registry{
		{
			Name:      "registry1",
			ClusterID: "cluster1",
		},
		{
			Name:      "registry2",
			ClusterID: "cluster2",
		},
	}
	ctrl := aggregate.NewController()
	for _, r := range registries {
		ctrl.AddRegistry(r)
	}
	if l := len(ctrl.GetRegistries()); l != 2 {
		t.Fatalf("Expected length of the registries slice should be 2, got %d", l)
	}
}

func TestDeleteRegistry(t *testing.T) {
	registries := []aggregate.Registry{
		{
			Name:      "registry1",
			ClusterID: "cluster1",
		},
		{
			Name:      "registry2",
			ClusterID: "cluster2",
		},
	}
	ctrl := aggregate.NewController()
	for _, r := range registries {
		ctrl.AddRegistry(r)
	}
	ctrl.DeleteRegistry(registries[0].ClusterID)
	if l := len(ctrl.GetRegistries()); l != 1 {
		t.Fatalf("Expected length of the registries slice should be 1, got %d", l)
	}
}

func TestGetRegistries(t *testing.T) {
	registries := []aggregate.Registry{
		{
			Name:      "registry1",
			ClusterID: "cluster1",
		},
		{
			Name:      "registry2",
			ClusterID: "cluster2",
		},
	}
	ctrl := aggregate.NewController()
	for _, r := range registries {
		ctrl.AddRegistry(r)
	}
	result := ctrl.GetRegistries()
	if len(ctrl.GetRegistries()) != len(result) {
		t.Fatal("Length of the original registries slice does not match to returned by GetRegistries.")
	}

	for i, registry := range ctrl.GetRegistries() {
		if !reflect.DeepEqual(result[i], registry) {
			t.Fatal("The original registries slice and resulting slice supposed to be identical.")
		}
	}
}

func TestGetProxyWorkloadLabels(t *testing.T) {
	aggregateCtl, cancel := mock.NewFakeAggregateControllerForMultiCluster()
	defer cancel()

	ip := "192.168.0.10"

	testCases := []struct {
		name      string
		proxyIP   string
		clusterID string
		expect    labels.Collection
	}{
		{
			name:      "should get labels from its own cluster workload",
			proxyIP:   ip,
			clusterID: "cluster1",
			expect:    labels.Collection{{"app": "cluster1"}},
		},
		{
			name:      "should get labels from its own cluster workload",
			proxyIP:   ip,
			clusterID: "cluster2",
			expect:    labels.Collection{{"app": "cluster2"}},
		},
		{
			name:      "can not get a workload from its cluster",
			proxyIP:   ip,
			clusterID: "cluster3",
			expect:    labels.Collection{},
		},
		{
			name:      "can not get a workload from its cluster",
			proxyIP:   "not exist ip",
			clusterID: "cluster1",
			expect:    labels.Collection{},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			workloadlabels, err := aggregateCtl.GetProxyWorkloadLabels(&model.Proxy{IPAddresses: []string{test.proxyIP}, ClusterID: test.clusterID})
			if err != nil {
				t.Fatalf("Failed get proxy workloadLabels: %v", err)
			}

			if !reflect.DeepEqual(workloadlabels, test.expect) {
				t.Errorf("Expect labels %v != %v", test.expect, workloadlabels)
			}
		})
	}

}
