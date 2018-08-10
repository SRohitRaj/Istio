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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta2

import (
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// ControllerRevisions returns a ControllerRevisionInformer.
	ControllerRevisions() ControllerRevisionInformer
	// DaemonSets returns a DaemonSetInformer.
	DaemonSets() DaemonSetInformer
	// Deployments returns a DeploymentInformer.
	Deployments() DeploymentInformer
	// ReplicaSets returns a ReplicaSetInformer.
	ReplicaSets() ReplicaSetInformer
	// StatefulSets returns a StatefulSetInformer.
	StatefulSets() StatefulSetInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// ControllerRevisions returns a ControllerRevisionInformer.
func (v *version) ControllerRevisions() ControllerRevisionInformer {
	return &controllerRevisionInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// DaemonSets returns a DaemonSetInformer.
func (v *version) DaemonSets() DaemonSetInformer {
	return &daemonSetInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Deployments returns a DeploymentInformer.
func (v *version) Deployments() DeploymentInformer {
	return &deploymentInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ReplicaSets returns a ReplicaSetInformer.
func (v *version) ReplicaSets() ReplicaSetInformer {
	return &replicaSetInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// StatefulSets returns a StatefulSetInformer.
func (v *version) StatefulSets() StatefulSetInformer {
	return &statefulSetInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
