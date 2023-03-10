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

package client

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pkg/config/schema/kubeclient"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/kubetypes"
	"istio.io/istio/pkg/ptr"
	"istio.io/pkg/log"
)

// Cached wraps a Kubernetes client providing cached read access and direct write access.
// This is based on informers, so most of the same caveats to informers apply here.
type CachedRead[T controllers.Object] interface {
	// Get looks up an object by name and namespace. If it does not exist, nil is returned
	Get(name, namespace string) T
	// List looks up an object by namespace and labels.
	// Use metav1.NamespaceAll and klabels.Everything() to select everything.
	List(namespace string, selector klabels.Selector) []T
	// ListUnfiltered is like List but ignores any *client side* filters previously configured.
	ListUnfiltered(namespace string, selector klabels.Selector) []T

	// AddEventHandler inserts a handler. The handler will be called for all Create/Update/Removals.
	// When ShutdownHandlers is called, the handler is removed.
	AddEventHandler(h cache.ResourceEventHandler)
	// HasSynced returns true when the informer is initially populated and that all handlers added
	// via AddEventHandler have been called with the initial state.
	// note: this differs from a standard informer HasSynced, which does not check handlers have been called.
	HasSynced() bool
	// ShutdownHandlers terminates all handlers added by AddEventHandler.
	// Warning: this only applies to handlers called via AddEventHandler; any handlers directly added
	// to the underlying informer are not touched
	ShutdownHandlers()
}

type Cached[T controllers.Object] interface {
	CachedRead[T]

	// Create creates a resource, returning the newly applied resource.
	Create(object T) (T, error)
	// Update updates a resource, returning the newly applied resource.
	Update(object T) (T, error)
	// UpdateStatus updates a resource's status, returning the newly applied resource.
	UpdateStatus(object T) (T, error)
	// Delete removes a resource.
	Delete(name, namespace string) error
}

type writeClient[T controllers.Object] struct {
	readClient[T]
	client kube.Client
}

type readClient[T controllers.Object] struct {
	inf                cache.SharedIndexInformer
	filter             func(t any) bool
	registeredHandlers []cache.ResourceEventHandlerRegistration
}

func (n *readClient[T]) Get(name, namespace string) T {
	obj, exists, err := n.inf.GetIndexer().GetByKey(keyFunc(name, namespace))
	if err != nil {
		return ptr.Empty[T]()
	}
	if !exists {
		return ptr.Empty[T]()
	}
	cast := obj.(T)
	if !n.applyFilter(cast) {
		return ptr.Empty[T]()
	}
	return cast
}

func (n *readClient[T]) applyFilter(t T) bool {
	if n.filter == nil {
		return true
	}
	return n.filter(t)
}

func (n *writeClient[T]) Create(object T) (T, error) {
	api := kubeclient.GetClient[T](n.client, object.GetNamespace())
	return api.Create(context.Background(), object, metav1.CreateOptions{})
}

func (n *writeClient[T]) Update(object T) (T, error) {
	api := kubeclient.GetClient[T](n.client, object.GetNamespace())
	return api.Update(context.Background(), object, metav1.UpdateOptions{})
}

func (n *writeClient[T]) UpdateStatus(object T) (T, error) {
	api, ok := kubeclient.GetClient[T](n.client, object.GetNamespace()).(kubetypes.WriteStatusAPI[T])
	if !ok {
		return ptr.Empty[T](), fmt.Errorf("%T does not support UpdateStatus", object)
	}
	return api.UpdateStatus(context.Background(), object, metav1.UpdateOptions{})
}

func (n *writeClient[T]) Delete(name, namespace string) error {
	api := kubeclient.GetClient[T](n.client, namespace)
	return api.Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (n *readClient[T]) ShutdownHandlers() {
	for _, c := range n.registeredHandlers {
		_ = n.inf.RemoveEventHandler(c)
	}
}

func (n *readClient[T]) AddEventHandler(h cache.ResourceEventHandler) {
	reg, _ := n.inf.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			if n.filter != nil && !n.filter(obj) {
				return
			}
			h.OnAdd(obj)
		},
		UpdateFunc: func(old, new any) {
			if n.filter != nil && !n.filter(new) {
				return
			}
			h.OnUpdate(old, new)
		},
		DeleteFunc: func(obj any) {
			if n.filter != nil && !n.filter(obj) {
				return
			}
			h.OnDelete(obj)
		},
	})
	n.registeredHandlers = append(n.registeredHandlers, reg)
}

func (n *readClient[T]) HasSynced() bool {
	return n.inf.HasSynced()
	/*
		TODO: client-go v0.27
		if !n.inf.HasSynced() {
			return false
		}
			for _, g := range n.registeredHandlers {
				if !g.HasSynced() {
					return false
				}
			}
		return true
	*/
}

func (n *readClient[T]) List(namespace string, selector klabels.Selector) []T {
	var res []T
	err := cache.ListAllByNamespace(n.inf.GetIndexer(), namespace, selector, func(i any) {
		cast := i.(T)
		if n.applyFilter(cast) {
			res = append(res, cast)
		}
	})

	// Should never happen
	if err != nil && features.EnableUnsafeAssertions {
		log.Fatalf("lister returned err for %v: %v", namespace, err)
	}
	return res
}

func (n *readClient[T]) ListUnfiltered(namespace string, selector klabels.Selector) []T {
	var res []T
	err := cache.ListAllByNamespace(n.inf.GetIndexer(), namespace, selector, func(i any) {
		cast := i.(T)
		res = append(res, cast)
	})

	// Should never happen
	if err != nil && features.EnableUnsafeAssertions {
		log.Fatalf("lister returned err for %v: %v", namespace, err)
	}
	return res
}

// Filter allows filtering read operations
type Filter struct {
	// A selector to restrict the list of returned objects by their labels.
	// This is a *server side* filter.
	LabelSelector string
	// A selector to restrict the list of returned objects by their fields.
	// This is a *server side* filter.
	FieldSelector string
	// ObjectFilter allows arbitrary filtering logic.
	// This is a *client side* filter. This means CPU/memory costs are still present for filtered objects.
	// Use LabelSelector or FieldSelector instead, if possible.
	ObjectFilter func(t any) bool
	// ObjectTransform allows arbitrarily modifying objects stored in the underlying cache.
	// If unset, a default transform is provided to remove ManagedFields (high cost, low value)
	ObjectTransform func(obj any) (any, error)
}

// NewCached returns a Cached for the given type.
// Internally, this uses a shared informer, so calling this multiple times will share the same internals.
func NewCached[T controllers.Object](c kube.Client) Cached[T] {
	return NewCachedFiltered[T](c, Filter{})
}

// NewCachedFiltered returns a Cached with some filter applied.
// Internally, this uses a shared informer, so calling this multiple times will share the same internals.
//
// Warning: currently, if filter.LabelSelector or filter.FieldSelector are set, the same informer will still be used
// This means there must only be one filter configuration for a given type using the same kube.Client (generally, this means the whole program).
// Use with caution.
func NewCachedFiltered[T controllers.Object](c kube.Client, filter Filter) Cached[T] {
	var inf cache.SharedIndexInformer
	if filter.LabelSelector == "" && filter.FieldSelector == "" {
		inf = kubeclient.GetInformer[T](c)
	} else {
		inf = kubeclient.GetInformerFiltered[T](c, kubetypes.InformerOptions{
			LabelSelector: filter.LabelSelector,
			FieldSelector: filter.FieldSelector,
		})
	}

	return &writeClient[T]{
		client:     c,
		readClient: newReadClient[T](c, inf, filter),
	}
}

// keyFunc is the internal API key function that returns "namespace"/"name" or
// "name" if "namespace" is empty
func keyFunc(name, namespace string) string {
	if len(namespace) == 0 {
		return name
	}
	return namespace + "/" + name
}
