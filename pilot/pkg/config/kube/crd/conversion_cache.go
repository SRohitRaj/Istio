// Copyright 2018 Istio Authors.
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

package crd

import (
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/cache"
)

// TODO: figure out how we want to size the cache. We alloc size_of(cache.lruEntry) == 32 bytes per entry,
// for a (constant, startup-allocated) total of 192 KiB today. Number is derived from current cluster size
// targets of 2k services, endpoints, ingresses = 6k resources total. This only includes the cache itself,
// not the size of the objects in the cache which we're also holding refs to.
const numCacheEntries = 6000

// ObjectConverter describes a function that can convert a k8s API object into an Istio model object.
type ObjectConverter func(schema model.ProtoSchema, object IstioObject, domain string) (*model.Config, error)

// CachingConverter implements an ObjectConverter that consults a cache before performing (expensive) marshalling.
// This struct is threadsafe, but cannot be copied.
type CachingConverter struct {
	cache cache.Cache
	inner ObjectConverter
}

// keyFunc accepts an objects name and returns its key; produced by calling newKeyFunc.
type keyFunc func(string) string

// newKeyFunc curries the application of model.Key by accepting the object type and domain. The returned function accepts
// the `name` of an object of type `typ` in config domain `domain` and returns its full key.
func newKeyFunc(typ, domain string) keyFunc {
	return func(name string) string {
		return model.Key(typ, name, domain)
	}
}

// NewCachingConverter returns a CachingConverter: an ObjectConverter wrapped in a fixed-size LRU cache.
func NewCachingConverter(converter ObjectConverter) *CachingConverter {
	return &CachingConverter{
		cache: cache.SimpleLRU(numCacheEntries),
		inner: converter,
	}
}

// ConvertObject consults a cache of model.Config objects before deferring to c's underlying ObjectConverter.
func (c *CachingConverter) ConvertObject(schema model.ProtoSchema, object IstioObject, domain string) (*model.Config, error) {
	key := model.Key(schema.Type, object.GetObjectMeta().Name, domain)
	if item, exists := c.get(key); exists {
		return item, nil
	}

	item, err := c.inner(schema, object, domain)
	if err != nil {
		return nil, err
	}

	// NB: we don't cache negative results (i.e. we don't cache the result when c.inner returns an err).
	// We should evaluate if that's worthwhile, given we implement evictions and a conversion should fail consistently
	// until the resource is updated.
	c.cache.Set(key, item)
	return item, nil
}

// Evict removes an entry from the cache.
func (c *CachingConverter) Evict(key string) {
	c.cache.Remove(key)
}

// returns the cached config object and whether the object was found; if !found, the returned config is nil
func (c *CachingConverter) get(key string) (*model.Config, bool) {
	item, found := c.cache.Get(key)
	if !found {
		return nil, false
	}
	return item.(*model.Config), true
}
