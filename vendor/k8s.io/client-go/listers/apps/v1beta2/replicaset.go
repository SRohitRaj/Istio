/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1beta2

import (
	v1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ReplicaSetLister helps list ReplicaSets.
type ReplicaSetLister interface {
	// List lists all ReplicaSets in the indexer.
	List(selector labels.Selector) (ret []*v1beta2.ReplicaSet, err error)
	// ReplicaSets returns an object that can list and get ReplicaSets.
	ReplicaSets(namespace string) ReplicaSetNamespaceLister
	ReplicaSetListerExpansion
}

// replicaSetLister implements the ReplicaSetLister interface.
type replicaSetLister struct {
	indexer cache.Indexer
}

// NewReplicaSetLister returns a new ReplicaSetLister.
func NewReplicaSetLister(indexer cache.Indexer) ReplicaSetLister {
	return &replicaSetLister{indexer: indexer}
}

// List lists all ReplicaSets in the indexer.
func (s *replicaSetLister) List(selector labels.Selector) (ret []*v1beta2.ReplicaSet, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.ReplicaSet))
	})
	return ret, err
}

// ReplicaSets returns an object that can list and get ReplicaSets.
func (s *replicaSetLister) ReplicaSets(namespace string) ReplicaSetNamespaceLister {
	return replicaSetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ReplicaSetNamespaceLister helps list and get ReplicaSets.
type ReplicaSetNamespaceLister interface {
	// List lists all ReplicaSets in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta2.ReplicaSet, err error)
	// Get retrieves the ReplicaSet from the indexer for a given namespace and name.
	Get(name string) (*v1beta2.ReplicaSet, error)
	ReplicaSetNamespaceListerExpansion
}

// replicaSetNamespaceLister implements the ReplicaSetNamespaceLister
// interface.
type replicaSetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ReplicaSets in the indexer for a given namespace.
func (s replicaSetNamespaceLister) List(selector labels.Selector) (ret []*v1beta2.ReplicaSet, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.ReplicaSet))
	})
	return ret, err
}

// Get retrieves the ReplicaSet from the indexer for a given namespace and name.
func (s replicaSetNamespaceLister) Get(name string) (*v1beta2.ReplicaSet, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta2.Resource("replicaset"), name)
	}
	return obj.(*v1beta2.ReplicaSet), nil
}
