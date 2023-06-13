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

package cli

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"

	"istio.io/istio/istioctl/pkg/util/handlers"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/ptr"
)

type Context interface {
	CLIClientWithRevision(rev string) (kube.CLIClient, error)
	CLIClient() (kube.CLIClient, error)
	InferPodInfoFromTypedResource(name, namespace string) (pod string, ns string, err error)
	InferPodsFromTypedResource(name, namespace string) ([]string, string, error)

	// TODO(hanxiaop) entirely drop KubeConfig and KubeContext, use CLIClient instead
	KubeConfig() string
	KubeContext() string
	Namespace() string
	IstioNamespace() string
	NamespaceOrDefault(namespace string) string
	// ConfigureDefaultNamespace sets the default namespace to use for commands that don't specify a namespace.
	// This should be called before NamespaceOrDefault is called.
	ConfigureDefaultNamespace()
}

type instance struct {
	// clients are cached clients for each revision
	clients map[string]kube.CLIClient
	RootFlags
}

func newKubeClientWithRevision(kubeconfig, configContext, revision string) (kube.CLIClient, error) {
	rc, err := kube.DefaultRestConfig(kubeconfig, configContext, func(config *rest.Config) {
		// We are running a one-off command locally, so we don't need to worry too much about rate limiting
		// Bumping this up greatly decreases install time
		config.QPS = 50
		config.Burst = 100
	})
	if err != nil {
		return nil, err
	}
	return kube.NewCLIClient(kube.NewClientConfigForRestConfig(rc), revision)
}

func NewCLIContext(rootFlags RootFlags) Context {
	return &instance{
		RootFlags: rootFlags,
	}
}

func (i *instance) CLIClientWithRevision(rev string) (kube.CLIClient, error) {
	if i.clients == nil {
		i.clients = make(map[string]kube.CLIClient)
	}
	if i.clients[rev] == nil {
		client, err := newKubeClientWithRevision(i.KubeConfig(), i.KubeContext(), rev)
		if err != nil {
			return nil, err
		}
		i.clients[rev] = client
	}
	return i.clients[rev], nil
}

func (i *instance) CLIClient() (kube.CLIClient, error) {
	return i.CLIClientWithRevision("")
}

func (i *instance) InferPodInfoFromTypedResource(name, namespace string) (pod string, ns string, err error) {
	client, err := i.CLIClient()
	if err != nil {
		return "", "", err
	}
	return handlers.InferPodInfoFromTypedResource(name, i.NamespaceOrDefault(namespace), MakeKubeFactory(client))
}

func (i *instance) InferPodsFromTypedResource(name, namespace string) ([]string, string, error) {
	client, err := i.CLIClient()
	if err != nil {
		return nil, "", err
	}
	return handlers.InferPodsFromTypedResource(name, i.NamespaceOrDefault(namespace), MakeKubeFactory(client))
}

func (i *instance) NamespaceOrDefault(namespace string) string {
	return handleNamespace(namespace, i.DefaultNamespace())
}

// handleNamespace returns the defaultNamespace if the namespace is empty
func handleNamespace(ns, defaultNamespace string) string {
	if ns == corev1.NamespaceAll {
		ns = defaultNamespace
	}
	return ns
}

type fakeContext struct {
	// clients are cached clients for each revision
	clients   map[string]kube.CLIClient
	rootFlags *RootFlags
}

func (f fakeContext) CLIClientWithRevision(rev string) (kube.CLIClient, error) {
	c := kube.NewFakeClient()
	f.clients[rev] = c
	return c, nil
}

func (f fakeContext) CLIClient() (kube.CLIClient, error) {
	return f.CLIClientWithRevision("")
}

func (f fakeContext) InferPodInfoFromTypedResource(name, namespace string) (pod string, ns string, err error) {
	client, err := f.CLIClient()
	if err != nil {
		return "", "", err
	}
	return handlers.InferPodInfoFromTypedResource(name, f.NamespaceOrDefault(namespace), MakeKubeFactory(client))
}

func (f fakeContext) InferPodsFromTypedResource(name, namespace string) ([]string, string, error) {
	client, err := f.CLIClient()
	if err != nil {
		return nil, "", err
	}
	return handlers.InferPodsFromTypedResource(name, f.NamespaceOrDefault(namespace), MakeKubeFactory(client))
}

func (f fakeContext) NamespaceOrDefault(namespace string) string {
	return handleNamespace(namespace, f.rootFlags.defaultNamespace)
}

func (f fakeContext) KubeConfig() string {
	return ""
}

func (f fakeContext) KubeContext() string {
	return ""
}

func (f fakeContext) Namespace() string {
	return f.rootFlags.Namespace()
}

func (f fakeContext) IstioNamespace() string {
	return f.rootFlags.IstioNamespace()
}

func (f fakeContext) ConfigureDefaultNamespace() {
	return
}

func NewFakeContext(namespace, istioNamespace string) Context {
	ns := namespace
	ins := istioNamespace
	return &fakeContext{
		clients: map[string]kube.CLIClient{},
		rootFlags: &RootFlags{
			kubeconfig:       ptr.Of[string](""),
			configContext:    ptr.Of[string](""),
			namespace:        &ns,
			istioNamespace:   &ins,
			defaultNamespace: "",
		},
	}
}
