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

// Package ingress provides a read-only view of Kubernetes ingress resources
// as an ingress rule configuration type store
package ingress

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/networking/v1beta1"

	"istio.io/pkg/ledger"

	ingress "k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/kubernetes"
	listerv1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/pkg/env"
	"istio.io/pkg/log"

	"istio.io/istio/pilot/pkg/model"
	kubecontroller "istio.io/istio/pilot/pkg/serviceregistry/kube/controller"
	"istio.io/istio/pkg/config/constants"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/config/schema/gvk"
	"istio.io/istio/pkg/config/schema/resource"
	"istio.io/istio/pkg/queue"
)

// In 1.0, the Gateway is defined in the namespace where the actual controller runs, and needs to be managed by
// user.
// The gateway is named by appending "-istio-autogenerated-k8s-ingress" to the name of the ingress.
//
// Currently the gateway namespace is hardcoded to istio-system (model.IstioIngressNamespace)
//
// VirtualServices are also auto-generated in the model.IstioIngressNamespace.
//
// The sync of Ingress objects to IP is done by status.go
// the 'ingress service' name is used to get the IP of the Service
// If ingress service is empty, it falls back to NodeExternalIP list, selected using the labels.
// This is using 'namespace' of pilot - but seems to be broken (never worked), since it uses Pilot's pod labels
// instead of the ingress labels.

// Follows mesh.IngressControllerMode setting to enable - OFF|STRICT|DEFAULT.
// STRICT requires "kubernetes.io/ingress.class" == mesh.IngressClass
// DEFAULT allows Ingress without explicit class.

// In 1.1:
// - K8S_INGRESS_NS - namespace of the Gateway that will act as ingress.
// - labels of the gateway set to "app=ingressgateway" for node_port, service set to 'ingressgateway' (matching default install)
//   If we need more flexibility - we can add it (but likely we'll deprecate ingress support first)
// -

var (
	schemas = collection.SchemasFor(
		collections.IstioNetworkingV1Alpha3Virtualservices,
		collections.IstioNetworkingV1Alpha3Gateways)
)

// Control needs RBAC permissions to write to Pods.

type controller struct {
	meshWatcher  mesh.Holder
	domainSuffix string

	client                 kubernetes.Interface
	queue                  queue.Instance
	ingressInformer        cache.SharedIndexInformer
	virtualServiceHandlers []func(model.Config, model.Config, model.Event)
	gatewayHandlers        []func(model.Config, model.Config, model.Event)
	// May be nil if ingress class is not supported in the cluster
	classes *v1beta1.IngressClassInformer

	serviceInformer cache.SharedIndexInformer
	serviceLister   listerv1.ServiceLister
}

var (
	// TODO: move to features ( and remove in 1.2 )
	ingressNamespace = env.RegisterStringVar("K8S_INGRESS_NS", "", "").Get()
)

var (
	errUnsupportedOp = errors.New("unsupported operation: the ingress config store is a read-only view")
)

func ingressClassSupported(client kubernetes.Interface) bool {
	_, s, _ := client.Discovery().ServerGroupsAndResources()
	// This may fail if any api service is down, but the result will still be populated, so we skip the error
	for _, res := range s {
		for _, api := range res.APIResources {
			if api.Kind == "IngressClass" && strings.HasPrefix(res.GroupVersion, "networking.k8s.io/") {
				return true
			}
		}
	}
	return false
}

// NewController creates a new Kubernetes controller
func NewController(client kubernetes.Interface, meshWatcher mesh.Holder,
	options kubecontroller.Options) model.ConfigStoreCache {

	// queue requires a time duration for a retry delay after a handler error
	q := queue.NewQueue(1 * time.Second)

	if ingressNamespace == "" {
		ingressNamespace = constants.IstioIngressNamespace
	}

	sharedInformerFactory := informers.NewSharedInformerFactoryWithOptions(client, options.ResyncPeriod,
		informers.WithNamespace(options.WatchedNamespaces))

	ingressInformer := sharedInformerFactory.Networking().V1beta1().Ingresses().Informer()
	log.Infof("Ingress controller watching namespaces %q", options.WatchedNamespaces)

	serviceInformer := sharedInformerFactory.Core().V1().Services()

	var classes *v1beta1.IngressClassInformer
	if ingressClassSupported(client) {
		sharedInformers := informers.NewSharedInformerFactory(client, options.ResyncPeriod)
		i := sharedInformers.Networking().V1beta1().IngressClasses()
		classes = &i
	} else {
		log.Infof("Skipping IngressClass, resource not supported")
	}

	c := &controller{
		meshWatcher:     meshWatcher,
		domainSuffix:    options.DomainSuffix,
		client:          client,
		queue:           q,
		ingressInformer: ingressInformer,
		classes:         classes,
		serviceInformer: serviceInformer.Informer(),
		serviceLister:   serviceInformer.Lister(),
	}

	ingressInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				q.Push(func() error {
					return c.onEvent(obj, model.EventAdd)
				})
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					q.Push(func() error {
						return c.onEvent(cur, model.EventUpdate)
					})
				}
			},
			DeleteFunc: func(obj interface{}) {
				q.Push(func() error {
					return c.onEvent(obj, model.EventDelete)
				})
			},
		})

	return c
}

func (c *controller) shouldProcessIngress(mesh *meshconfig.MeshConfig, i *ingress.Ingress) (bool, error) {
	var class *ingress.IngressClass
	if c.classes != nil && i.Spec.IngressClassName != nil {
		c, err := (*c.classes).Lister().Get(*i.Spec.IngressClassName)
		if err != nil && !kerrors.IsNotFound(err) {
			return false, fmt.Errorf("failed to get ingress class %v: %v", i.Spec.IngressClassName, err)
		}
		class = c
	}
	return shouldProcessIngressWithClass(mesh, i, class), nil
}

func (c *controller) onEvent(obj interface{}, event model.Event) error {
	if !c.HasSynced() {
		return errors.New("waiting till full synchronization")
	}

	ing, ok := obj.(*ingress.Ingress)
	process, err := c.shouldProcessIngress(c.meshWatcher.Mesh(), ing)
	if err != nil {
		return err
	}
	if !ok || !process {
		return nil
	}
	log.Infof("ingress event %s for %s/%s", event, ing.Namespace, ing.Name)

	// Trigger updates for Gateway and VirtualService
	// TODO: we could be smarter here and only trigger when real changes were found
	for _, f := range c.virtualServiceHandlers {
		f(model.Config{}, model.Config{
			ConfigMeta: model.ConfigMeta{
				GroupVersionKind: gvk.VirtualService,
			},
		}, event)
	}
	for _, f := range c.gatewayHandlers {
		f(model.Config{}, model.Config{
			ConfigMeta: model.ConfigMeta{
				GroupVersionKind: gvk.Gateway,
			},
		}, event)
	}

	return nil
}

func (c *controller) RegisterEventHandler(kind resource.GroupVersionKind, f func(model.Config, model.Config, model.Event)) {
	switch kind {
	case gvk.VirtualService:
		c.virtualServiceHandlers = append(c.virtualServiceHandlers, f)
	case gvk.Gateway:
		c.gatewayHandlers = append(c.gatewayHandlers, f)
	}
}

func (c *controller) Version() string {
	panic("implement me")
}

func (c *controller) GetResourceAtVersion(string, string) (resourceVersion string, err error) {
	panic("implement me")
}

func (c *controller) GetLedger() ledger.Ledger {
	log.Warnf("GetLedger: %s", errors.New("this operation is not supported by kube ingress controller"))
	return nil
}

func (c *controller) SetLedger(ledger.Ledger) error {
	return errors.New("this SetLedger operation is not supported by kube ingress controller")
}

func (c *controller) HasSynced() bool {
	return c.ingressInformer.HasSynced() && c.serviceInformer.HasSynced()
}

func (c *controller) Run(stop <-chan struct{}) {
	go func() {
		cache.WaitForCacheSync(stop, c.HasSynced)
		c.queue.Run(stop)
	}()
	go c.ingressInformer.Run(stop)
	go c.serviceInformer.Run(stop)
	if c.classes != nil {
		go (*c.classes).Informer().Run(stop)
	}
	<-stop
}

func (c *controller) Schemas() collection.Schemas {
	//TODO: are these two config descriptors right?
	return schemas
}

func (c *controller) Get(typ resource.GroupVersionKind, name, namespace string) *model.Config {
	return nil
}

func (c *controller) List(typ resource.GroupVersionKind, namespace string) ([]model.Config, error) {
	if typ != gvk.Gateway &&
		typ != gvk.VirtualService {
		return nil, errUnsupportedOp
	}

	out := make([]model.Config, 0)

	ingressByHost := map[string]*model.Config{}

	for _, obj := range c.ingressInformer.GetStore().List() {
		ingress := obj.(*ingress.Ingress)
		if namespace != "" && namespace != ingress.Namespace {
			continue
		}
		process, err := c.shouldProcessIngress(c.meshWatcher.Mesh(), ingress)
		if err != nil {
			return nil, err
		}
		if !process {
			continue
		}

		switch typ {
		case gvk.VirtualService:
			ConvertIngressVirtualService(*ingress, c.domainSuffix, ingressByHost, c.serviceLister)
		case gvk.Gateway:
			gateways := ConvertIngressV1alpha3(*ingress, c.meshWatcher.Mesh(), c.domainSuffix)
			out = append(out, gateways)
		}
	}

	if typ == gvk.VirtualService {
		for _, obj := range ingressByHost {
			out = append(out, *obj)
		}
	}

	return out, nil
}

func (c *controller) Create(_ model.Config) (string, error) {
	return "", errUnsupportedOp
}

func (c *controller) Update(_ model.Config) (string, error) {
	return "", errUnsupportedOp
}

func (c *controller) Delete(_ resource.GroupVersionKind, _, _ string) error {
	return errUnsupportedOp
}
