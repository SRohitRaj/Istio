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
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	listerv1 "k8s.io/client-go/listers/core/v1"

	"istio.io/pkg/log"
)

// DiscoveryNamespacesFilter tracks the set of namespaces selected for discovery, which are updated by the discovery namespace controller.
// It exposes a filter function used for filtering out objects that don't reside in namespaces selected for discovery.
type DiscoveryNamespacesFilter interface {
	// return true if the input object resides in a namespace selected for discovery
	Filter(obj interface{}) bool
	// invoked when meshConfig's discoverySelectors change, returns any newly selected namespaces and deselected namespaces
	SelectorsChanged(discoverySelectors []*metav1.LabelSelector) (selectedNamespaces []string, deselectedNamespaces []string)
	// return true if the created namespace is selected for discovery
	NamespaceCreated(ns metav1.ObjectMeta) (membershipChanged bool)
	// membershipChanged will be true if the updated namespace is newly selected or deselected for discovery
	NamespaceUpdated(oldNs, newNs metav1.ObjectMeta) (membershipChanged bool, namespaceAdded bool)
	// return true if the deleted namespace was selected for discovery
	NamespaceDeleted(ns metav1.ObjectMeta) (membershipChanged bool)
}

type discoveryNamespacesFilter struct {
	lock                sync.RWMutex
	nsLister            listerv1.NamespaceLister
	discoveryNamespaces sets.String
	discoverySelectors  []labels.Selector // nil if discovery selectors are not specified, permits all namespaces for discovery
}

func NewDiscoveryNamespacesFilter(
	nsLister listerv1.NamespaceLister,
	discoverySelectors []*metav1.LabelSelector,
) DiscoveryNamespacesFilter {
	discoveryNamespacesFilter := &discoveryNamespacesFilter{
		nsLister: nsLister,
	}

	// initialize discovery namespaces filter
	discoveryNamespacesFilter.SelectorsChanged(discoverySelectors)

	return discoveryNamespacesFilter
}

func (d *discoveryNamespacesFilter) Filter(obj interface{}) bool {
	// permit all objects if discovery selectors are not specified
	if len(d.discoverySelectors) == 0 {
		return true
	}

	// permit if object resides in a namespace labeled for discovery
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.discoveryNamespaces.Has(obj.(metav1.Object).GetNamespace())
}

// initialize the discovery filter state with the discovery selectors and selected namespaces
func (d *discoveryNamespacesFilter) SelectorsChanged(
	discoverySelectors []*metav1.LabelSelector,
) (selectedNamespaces []string, deselectedNamespaces []string) {
	oldDiscoveryNamespaces := d.discoveryNamespaces

	var selectors []labels.Selector
	newDiscoveryNamespaces := sets.NewString()

	namespaceList, err := d.nsLister.List(labels.Everything())
	if err != nil {
		log.Errorf("error initializing discovery namespaces filter, failed to list namespaces: %v", err)
		return
	}

	// convert LabelSelectors to Selectors
	for _, selector := range discoverySelectors {
		ls, err := metav1.LabelSelectorAsSelector(selector)
		if err != nil {
			log.Errorf("error initializing discovery namespaces filter, invalid discovery selector: %v", err)
			return
		}
		selectors = append(selectors, ls)
	}

	if len(selectors) != 0 {
		// range over all namespaces to get discovery namespaces
		for _, ns := range namespaceList {
			for _, selector := range selectors {
				// omitting discoverySelectors indicates discovering all namespaces
				if selector.Matches(labels.Set(ns.Labels)) {
					newDiscoveryNamespaces.Insert(ns.Name)
				}
			}
		}
	}

	selectedNamespaces = newDiscoveryNamespaces.Difference(oldDiscoveryNamespaces).List()
	deselectedNamespaces = oldDiscoveryNamespaces.Difference(newDiscoveryNamespaces).List()

	// update filter state
	d.discoveryNamespaces = newDiscoveryNamespaces
	d.discoverySelectors = selectors

	return
}

// if newly created namespace is selected, update namespace membership
func (d *discoveryNamespacesFilter) NamespaceCreated(ns metav1.ObjectMeta) (membershipChanged bool) {
	if d.isSelected(ns.Labels) {
		d.addNamespace(ns.Name)
		return true
	}
	return false
}

// if updated namespace was a member and no longer selected, or was not a member and now selected, update namespace membership
func (d *discoveryNamespacesFilter) NamespaceUpdated(oldNs, newNs metav1.ObjectMeta) (membershipChanged bool, namespaceAdded bool) {
	if d.hasNamespace(oldNs.Name) && !d.isSelected(newNs.Labels) {
		d.removeNamespace(oldNs.Name)
		return true, false
	}
	if !d.hasNamespace(oldNs.Name) && d.isSelected(newNs.Labels) {
		d.addNamespace(oldNs.Name)
		return true, true
	}
	return false, false
}

// if deleted namespace was a member, remove it
func (d *discoveryNamespacesFilter) NamespaceDeleted(ns metav1.ObjectMeta) (membershipChanged bool) {
	if d.isSelected(ns.Labels) {
		d.removeNamespace(ns.Name)
		return true
	}
	return false
}

func (d *discoveryNamespacesFilter) addNamespace(ns string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.discoveryNamespaces.Insert(ns)
}

func (d *discoveryNamespacesFilter) hasNamespace(ns string) bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.discoveryNamespaces.Has(ns)
}

func (d *discoveryNamespacesFilter) removeNamespace(ns string) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.discoveryNamespaces.Delete(ns)
}

func (d *discoveryNamespacesFilter) isSelected(labels labels.Set) bool {
	// permit all objects if discovery selectors are not specified
	if len(d.discoverySelectors) == 0 {
		return true
	}

	for _, selector := range d.discoverySelectors {
		if selector.Matches(labels) {
			return true
		}
	}

	return false
}
