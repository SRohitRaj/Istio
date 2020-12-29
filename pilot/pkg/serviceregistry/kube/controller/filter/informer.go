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

package filter

import (
	"k8s.io/client-go/tools/cache"
)

type FilteredSharedIndexInformer interface {
	AddEventHandler(handler cache.ResourceEventHandler)
	GetIndexer() cache.Indexer
	HasSynced() bool
	Run(stopCh <-chan struct{})
}

type filteredSharedIndexInformer struct {
	filterFuncFactory func() Func
	cache.SharedIndexInformer
	filteredIndexer *filteredIndexer
}

// wrap a SharedIndexInformer's handlers and indexer with a filter predicate,
// which scopes the processed objects to only those that satisfy the predicate
func NewFilteredSharedIndexInformer(
	filterFuncFactory func() Func,
	sharedIndexInformer cache.SharedIndexInformer,
) FilteredSharedIndexInformer {
	return &filteredSharedIndexInformer{
		filterFuncFactory:   filterFuncFactory,
		SharedIndexInformer: sharedIndexInformer,
		filteredIndexer:     newFilteredIndexer(filterFuncFactory, sharedIndexInformer.GetIndexer()),
	}
}

// filter incoming objects before forwarding to event handler
func (w *filteredSharedIndexInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	w.SharedIndexInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if !w.filterFuncFactory()(obj) {
				return
			}
			handler.OnAdd(obj)
		},
		UpdateFunc: func(old, new interface{}) {
			if !w.filterFuncFactory()(new) {
				return
			}
			handler.OnUpdate(old, new)
		},
		DeleteFunc: func(obj interface{}) {
			if !w.filterFuncFactory()(obj) {
				return
			}
			handler.OnDelete(obj)
		},
	})
}

func (w *filteredSharedIndexInformer) HasSynced() bool {
	w.SharedIndexInformer.GetStore()
	return w.SharedIndexInformer.HasSynced()
}

func (w *filteredSharedIndexInformer) Run(stopCh <-chan struct{}) {
	w.SharedIndexInformer.Run(stopCh)
}

func (w *filteredSharedIndexInformer) GetIndexer() cache.Indexer {
	return w.filteredIndexer
}

type filteredIndexer struct {
	filterFuncFactory func() Func
	cache.Indexer
}

func newFilteredIndexer(
	filterFuncFactory func() Func,
	indexer cache.Indexer,
) *filteredIndexer {
	return &filteredIndexer{
		filterFuncFactory: filterFuncFactory,
		Indexer:           indexer,
	}
}

func (w filteredIndexer) List() []interface{} {
	// initialize filterFunc once to avoid listing/iterating all namespaces for each listed object
	filterFunc := w.filterFuncFactory()
	unfiltered := w.Indexer.List()
	var filtered []interface{}
	for _, obj := range unfiltered {
		if filterFunc(obj) {
			filtered = append(filtered, obj)
		}
	}
	return filtered
}

func (w filteredIndexer) GetByKey(key string) (item interface{}, exists bool, err error) {
	item, exists, err = w.Indexer.GetByKey(key)
	if !exists || err != nil {
		return nil, exists, err
	}
	if w.filterFuncFactory()(item) {
		return item, true, nil
	}
	return nil, false, nil
}
