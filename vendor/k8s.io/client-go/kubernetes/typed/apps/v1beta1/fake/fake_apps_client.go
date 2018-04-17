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

package fake

import (
	v1beta1 "k8s.io/client-go/kubernetes/typed/apps/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAppsV1beta1 struct {
	*testing.Fake
}

func (c *FakeAppsV1beta1) ControllerRevisions(namespace string) v1beta1.ControllerRevisionInterface {
	return &FakeControllerRevisions{c, namespace}
}

func (c *FakeAppsV1beta1) Deployments(namespace string) v1beta1.DeploymentInterface {
	return &FakeDeployments{c, namespace}
}

func (c *FakeAppsV1beta1) Scales(namespace string) v1beta1.ScaleInterface {
	return &FakeScales{c, namespace}
}

func (c *FakeAppsV1beta1) StatefulSets(namespace string) v1beta1.StatefulSetInterface {
	return &FakeStatefulSets{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAppsV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
