// Copyright 2019 Istio Authors
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

package bootstrap

import (
	"fmt"

	"istio.io/pkg/log"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/serviceregistry"
	"istio.io/istio/pilot/pkg/serviceregistry/aggregate"
	"istio.io/istio/pilot/pkg/serviceregistry/consul"
	"istio.io/istio/pilot/pkg/serviceregistry/external"
	kubecontroller "istio.io/istio/pilot/pkg/serviceregistry/kube/controller"
	"istio.io/istio/pilot/pkg/serviceregistry/memory"
	"istio.io/istio/pkg/config/host"
)

// initServiceControllers creates and initializes the service controllers
func (s *Server) initServiceControllers(args *PilotArgs) error {
	serviceControllers := aggregate.NewController()
	registered := make(map[serviceregistry.ServiceRegistry]bool)
	for _, r := range args.Service.Registries {
		serviceRegistry := serviceregistry.ServiceRegistry(r)
		if _, exists := registered[serviceRegistry]; exists {
			log.Warnf("%s registry specified multiple times.", r)
			continue
		}
		registered[serviceRegistry] = true
		log.Infof("Adding %s registry adapter", serviceRegistry)
		switch serviceRegistry {
		case serviceregistry.KubernetesRegistry:
			if err := s.initKubeRegistry(serviceControllers, args); err != nil {
				return err
			}
		case serviceregistry.MCPRegistry:
			if s.mcpDiscovery != nil {
				serviceControllers.AddRegistry(
					aggregate.Registry{
						Name:             serviceregistry.MCPRegistry,
						ServiceDiscovery: s.mcpDiscovery,
						Controller:       s.mcpDiscovery,
					})
			}
		case serviceregistry.ConsulRegistry:
			if err := s.initConsulRegistry(serviceControllers, args); err != nil {
				return err
			}
		case serviceregistry.MockRegistry:
			s.initMemoryRegistry(serviceControllers)
		default:
			return fmt.Errorf("service registry %s is not supported", r)
		}
	}

	serviceEntryStore := external.NewServiceDiscovery(s.configController, s.istioConfigStore)

	// add service entry registry to aggregator by default
	serviceEntryRegistry := aggregate.Registry{
		Name:             "ServiceEntries",
		Controller:       serviceEntryStore,
		ServiceDiscovery: serviceEntryStore,
	}
	serviceControllers.AddRegistry(serviceEntryRegistry)

	s.ServiceController = serviceControllers

	// Defer running of the service controllers.
	s.addStartFunc(func(stop <-chan struct{}) error {
		go s.ServiceController.Run(stop)
		return nil
	})

	return nil
}

// initKubeRegistry creates all the k8s service controllers under this pilot
func (s *Server) initKubeRegistry(serviceControllers *aggregate.Controller, args *PilotArgs) (err error) {
	clusterID := string(serviceregistry.KubernetesRegistry)
	log.Infof("Primary Cluster name: %s", clusterID)
	args.Config.ControllerOptions.ClusterID = clusterID
	kubectl := kubecontroller.NewController(s.kubeClient, args.Config.ControllerOptions)
	s.kubeRegistry = kubectl
	serviceControllers.AddRegistry(
		aggregate.Registry{
			Name:             serviceregistry.KubernetesRegistry,
			ClusterID:        clusterID,
			ServiceDiscovery: kubectl,
			Controller:       kubectl,
		})

	return
}

func (s *Server) initConsulRegistry(serviceControllers *aggregate.Controller, args *PilotArgs) error {
	log.Infof("Consul url: %v", args.Service.Consul.ServerURL)
	conctl, conerr := consul.NewController(
		args.Service.Consul.ServerURL)
	if conerr != nil {
		return fmt.Errorf("failed to create Consul controller: %v", conerr)
	}
	serviceControllers.AddRegistry(
		aggregate.Registry{
			Name:             serviceregistry.ConsulRegistry,
			ServiceDiscovery: conctl,
			Controller:       conctl,
		})

	return nil
}

func (s *Server) initMemoryRegistry(serviceControllers *aggregate.Controller) {
	// MemServiceDiscovery implementation
	discovery := memory.NewDiscovery(map[host.Name]*model.Service{}, 2)

	registry := aggregate.Registry{
		Name:             serviceregistry.MockRegistry,
		ServiceDiscovery: discovery,
		Controller:       &memory.MockController{},
	}

	serviceControllers.AddRegistry(registry)
}
