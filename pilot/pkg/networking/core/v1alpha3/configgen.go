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

package v1alpha3

import (
//	"sync"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/plugin"
)

type ConfigGeneratorImpl struct {
	// List of plugins that modify code generated by this config generator
	Plugins []plugin.Plugin
	// List of outbound clusters keyed by configNamespace
	// Must be rebuilt for each push epoch
	PrecomputedOutboundClusters map[string][]*xdsapi.Cluster
	// TODO: add others in future
}

func NewConfigGenerator(plugins []plugin.Plugin) *ConfigGeneratorImpl {
	return &ConfigGeneratorImpl{
		Plugins: plugins,
	}
}

func (configgen *ConfigGeneratorImpl) BuildSharedPushState(env *model.Environment, push *model.PushContext) error {
	// namespaceMap := map[string]struct{}{}
	// clustersByNamespace := map[string][]*xdsapi.Cluster{}

	// services := push.Services(nil)
	// for _, svc := range services {
	// 	namespaceMap[svc.Attributes.Namespace] = struct{}{}
	// }
	// namespaceMap[""] = struct{}{}

	// // generate outbound for all namespaces in parallel.
	// wg := &sync.WaitGroup{}
	// mutex := &sync.Mutex{}
	// wg.Add(len(namespaceMap))
	// for ns := range namespaceMap {
	// 	go func(ns string) {
	// 		defer wg.Done()
	// 		dummyNode := model.Proxy{
	// 			ConfigNamespace: ns,
	// 		}
	// 		clusters := configgen.buildOutboundClusters(env, &dummyNode, push)
	// 		mutex.Lock()
	// 		clustersByNamespace[ns] = clusters
	// 		mutex.Unlock()
	// 	}(ns)
	// }
	// wg.Wait()

	// configgen.PrecomputedOutboundClusters = clustersByNamespace
	return nil
}

func (configgen *ConfigGeneratorImpl) CanUsePrecomputedCDS(proxy *model.Proxy) bool {
	return false
	// networkView := model.GetNetworkView(proxy)
	// // If we have only more than one network view for the proxy, then recompute CDS.
	// // Because, by default, we cache the CDS output for proxies in the UnnamedNetwork only.
	// if len(networkView) > 1 {
	// 	return false
	// }

	// return networkView[model.UnnamedNetwork]
}
