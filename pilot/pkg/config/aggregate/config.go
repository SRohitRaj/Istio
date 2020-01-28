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

// Package aggregate implements a read-only aggregator for config stores.
package aggregate

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/resource"
)

var errorUnsupported = errors.New("unsupported operation: the config aggregator is read-only")

// Make creates an aggregate config store from several config stores and
// unifies their descriptors
func Make(stores []model.ConfigStore) (model.ConfigStore, error) {
	union := collection.NewSchemasBuilder()
	storeTypes := make(map[resource.GroupVersionKind][]model.ConfigStore)
	for _, store := range stores {
		for _, s := range store.Schemas().All() {
			if len(storeTypes[s.Resource().GroupVersionKind()]) == 0 {
				if err := union.Add(s); err != nil {
					return nil, err
				}
			}
			storeTypes[s.Resource().GroupVersionKind()] = append(storeTypes[s.Resource().GroupVersionKind()], store)
		}
	}

	schemas := union.Build()
	if err := schemas.Validate(); err != nil {
		return nil, err
	}
	result := &store{
		schemas: schemas,
		stores:  storeTypes,
	}

	// in most cases (all cases supported by helm), pilot has only one configStore, but it is always wrapped in this
	// aggregate.  This allows us to pass through data from a single config ledger, while gracefully failing in the
	// unlikely scenario that multiple configSources are supplied.
	// in the case of multiple configSources, we could simply require that all sources share a single ledger,
	// but the ledger has to be provided at construction time since it is an implementation detail.  Perhaps we should
	// consider allowing the ledger to be set on the fly to accommodate this scenario?
	if len(stores) == 1 {
		result.getVersion = stores[0].Version
		result.getResourceAtVersion = stores[0].GetResourceAtVersion
	} else {
		result.getVersion = func() string { return "" }
		result.getResourceAtVersion = func(one, two string) (string, error) {
			return "", errors.New("config distribution status not supported on multiple configSources")
		}
	}
	return result, nil
}

// MakeCache creates an aggregate config store cache from several config store
// caches.
func MakeCache(caches []model.ConfigStoreCache) (model.ConfigStoreCache, error) {
	stores := make([]model.ConfigStore, 0, len(caches))
	for _, cache := range caches {
		stores = append(stores, cache)
	}
	store, err := Make(stores)
	if err != nil {
		return nil, err
	}
	return &storeCache{
		ConfigStore: store,
		caches:      caches,
	}, nil
}

type store struct {
	// schemas is the unified
	schemas collection.Schemas

	// stores is a mapping from config type to a store
	stores map[resource.GroupVersionKind][]model.ConfigStore

	getVersion func() string

	getResourceAtVersion func(version, key string) (resourceVersion string, err error)
}

func (cr *store) GetResourceAtVersion(version string, key string) (resourceVersion string, err error) {
	return cr.getResourceAtVersion(version, key)
}

func (cr *store) Schemas() collection.Schemas {
	return cr.schemas
}

func (cr *store) Version() string {
	return cr.getVersion()
}

// Get the first config found in the stores.
func (cr *store) Get(typ resource.GroupVersionKind, name, namespace string) *model.Config {
	for _, store := range cr.stores[typ] {
		config := store.Get(typ, name, namespace)
		if config != nil {
			return config
		}
	}
	return nil
}

// List all configs in the stores.
func (cr *store) List(typ resource.GroupVersionKind, namespace string) ([]model.Config, error) {
	if len(cr.stores[typ]) == 0 {
		return nil, fmt.Errorf("missing type %q", typ)
	}
	var errs *multierror.Error
	var configs []model.Config
	// Used to remove duplicated config
	configMap := make(map[string]struct{})

	for _, store := range cr.stores[typ] {
		storeConfigs, err := store.List(typ, namespace)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		for _, config := range storeConfigs {
			key := config.Type + config.Namespace + config.Name
			if _, exist := configMap[key]; exist {
				continue
			}
			configs = append(configs, config)
			configMap[key] = struct{}{}
		}
	}
	return configs, errs.ErrorOrNil()
}

func (cr *store) Delete(_ resource.GroupVersionKind, _, _ string) error {
	return errorUnsupported
}

func (cr *store) Create(model.Config) (string, error) {
	return "", errorUnsupported
}

func (cr *store) Update(model.Config) (string, error) {
	return "", errorUnsupported
}

type storeCache struct {
	model.ConfigStore
	caches []model.ConfigStoreCache
}

func (cr *storeCache) HasSynced() bool {
	for _, cache := range cr.caches {
		if !cache.HasSynced() {
			return false
		}
	}
	return true
}

func (cr *storeCache) RegisterEventHandler(kind resource.GroupVersionKind, handler func(model.Config, model.Config, model.Event)) {
	for _, cache := range cr.caches {
		if _, exists := cache.Schemas().FindByGroupVersionKind(kind); exists {
			cache.RegisterEventHandler(kind, handler)
		}
	}
}

func (cr *storeCache) Run(stop <-chan struct{}) {
	for _, cache := range cr.caches {
		go cache.Run(stop)
	}
	<-stop
}
