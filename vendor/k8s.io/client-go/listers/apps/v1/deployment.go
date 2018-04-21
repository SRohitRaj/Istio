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

// This file was automatically generated by lister-gen

package v1

import (
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DeploymentLister helps list Deployments.
type DeploymentLister interface {
	// List lists all Deployments in the indexer.
	List(selector labels.Selector) (ret []*v1.Deployment, err error)
	// Deployments returns an object that can list and get Deployments.
	Deployments(namespace string) DeploymentNamespaceLister
	DeploymentListerExpansion
}

// deploymentLister implements the DeploymentLister interface.
type deploymentLister struct {
	indexer cache.Indexer
}

// NewDeploymentLister returns a new DeploymentLister.
func NewDeploymentLister(indexer cache.Indexer) DeploymentLister {
	return &deploymentLister{indexer: indexer}
}

// List lists all Deployments in the indexer.
func (s *deploymentLister) List(selector labels.Selector) (ret []*v1.Deployment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Deployment))
	})
	return ret, err
}

// Deployments returns an object that can list and get Deployments.
func (s *deploymentLister) Deployments(namespace string) DeploymentNamespaceLister {
	return deploymentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DeploymentNamespaceLister helps list and get Deployments.
type DeploymentNamespaceLister interface {
	// List lists all Deployments in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Deployment, err error)
	// Get retrieves the Deployment from the indexer for a given namespace and name.
	Get(name string) (*v1.Deployment, error)
	DeploymentNamespaceListerExpansion
}

// deploymentNamespaceLister implements the DeploymentNamespaceLister
// interface.
type deploymentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Deployments in the indexer for a given namespace.
func (s deploymentNamespaceLister) List(selector labels.Selector) (ret []*v1.Deployment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Deployment))
	})
	return ret, err
}

// Get retrieves the Deployment from the indexer for a given namespace and name.
func (s deploymentNamespaceLister) Get(name string) (*v1.Deployment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("deployment"), name)
	}
	return obj.(*v1.Deployment), nil
}
