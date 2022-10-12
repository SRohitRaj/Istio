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

package local

import (
	"strings"

	"google.golang.org/protobuf/proto"

	"istio.io/istio/pilot/pkg/config/aggregate"
	"istio.io/istio/pilot/pkg/config/memory"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/analysis/scope"
	"istio.io/istio/pkg/config/schema/collection"
)

// analyzableResourcesStore is the cache to get resources that needs to be analyzed
type analyzableResourcesStore struct {
	schemas        collection.Schemas
	resourcesStore model.ConfigStore
}

func newAnalyzableResourcesCache(initializedStore model.ConfigStoreController) *analyzableResourcesStore {
	a := &analyzableResourcesStore{
		schemas: initializedStore.Schemas(),
	}
	combinedCache, _ := aggregate.MakeWriteableCache([]model.ConfigStoreController{initializedStore}, memory.Make(a.schemas))
	a.resourcesStore = &dfCache{ConfigStore: combinedCache}
	for _, sch := range a.schemas.All() {
		initializedStore.RegisterEventHandler(sch.Resource().GroupVersionKind(), a.processConfigChanges)
	}
	return a
}

func (a *analyzableResourcesStore) getStore() model.ConfigStoreController {
	store := a.resourcesStore
	// Empty the store for the next analysis run
	a.resourcesStore = memory.Make(a.schemas)

	return &dfCache{
		ConfigStore: store,
	}
}

func (a *analyzableResourcesStore) processConfigChanges(prev config.Config, curr config.Config, event model.Event) {
	switch event {
	case model.EventAdd, model.EventUpdate:
		if !needsReAnalyze(prev, curr) {
			return
		}
		// We don't care about ResourceVersion, this is to avoid of update error of the in-memory store
		curr.ResourceVersion = ""
		_, err := a.resourcesStore.Update(curr)
		if err != nil {
			_, err = a.resourcesStore.Create(curr)
			if err != nil {
				scope.Analysis.Errorf("error create config : %v", err)
			}
		}
	case model.EventDelete:
		// Do not analyze deleted resources
	}
}

func needsReAnalyze(prev config.Config, curr config.Config) bool {
	if prev.GroupVersionKind != curr.GroupVersionKind {
		// This should never happen.
		return true
	}
	if !strings.HasSuffix(prev.GroupVersionKind.Group, "istio.io") {
		return true
	}
	prevspecProto, okProtoP := prev.Spec.(proto.Message)
	currspecProto, okProtoC := curr.Spec.(proto.Message)
	if okProtoP && okProtoC {
		return !proto.Equal(prevspecProto, currspecProto)
	}
	return true
}
