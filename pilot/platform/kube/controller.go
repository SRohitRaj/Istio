// Copyright 2017 Istio Authors
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

package kube

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"istio.io/istio/pilot/model"
)

const (
	// NodeRegionLabel is the well-known label for kubernetes node region
	NodeRegionLabel = "failure-domain.beta.kubernetes.io/region"
	// NodeZoneLabel is the well-known label for kubernetes node zone
	NodeZoneLabel = "failure-domain.beta.kubernetes.io/zone"
	// IstioNamespace used by default for Istio cluster-wide installation
	IstioNamespace = "istio-system"
)

// ControllerOptions stores the configurable attributes of a Controller.
type ControllerOptions struct {
	// Namespace the controller watches. If set to meta_v1.NamespaceAll (""), controller watches all namespaces
	WatchedNamespace string
	ResyncPeriod     time.Duration
	DomainSuffix     string
}

// Controller is a collection of synchronized resource watchers
// Caches are thread-safe
type Controller struct {
	domainSuffix string

	client    kubernetes.Interface
	queue     Queue
	services  cacheHandler
	endpoints cacheHandler
	nodes     cacheHandler

	pods *PodCache
}

type cacheHandler struct {
	informer cache.SharedIndexInformer
	handler  *ChainHandler
}

// NewController creates a new Kubernetes controller
func NewController(client kubernetes.Interface, options ControllerOptions) *Controller {
	glog.V(2).Infof("Service controller watching namespace %q", options.WatchedNamespace)

	// Queue requires a time duration for a retry delay after a handler error
	out := &Controller{
		domainSuffix: options.DomainSuffix,
		client:       client,
		queue:        NewQueue(1 * time.Second),
	}

	out.services = out.createInformer(&v1.Service{}, options.ResyncPeriod,
		func(opts meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Services(options.WatchedNamespace).List(opts)
		},
		func(opts meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Services(options.WatchedNamespace).Watch(opts)
		})

	out.endpoints = out.createInformer(&v1.Endpoints{}, options.ResyncPeriod,
		func(opts meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).List(opts)
		},
		func(opts meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Endpoints(options.WatchedNamespace).Watch(opts)
		})

	out.nodes = out.createInformer(&v1.Node{}, options.ResyncPeriod,
		func(opts meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Nodes().List(opts)
		},
		func(opts meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Nodes().Watch(opts)
		})

	out.pods = newPodCache(out.createInformer(&v1.Pod{}, options.ResyncPeriod,
		func(opts meta_v1.ListOptions) (runtime.Object, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).List(opts)
		},
		func(opts meta_v1.ListOptions) (watch.Interface, error) {
			return client.CoreV1().Pods(options.WatchedNamespace).Watch(opts)
		}))

	return out
}

// notify is the first handler in the handler chain.
// Returning an error causes repeated execution of the entire chain.
func (c *Controller) notify(obj interface{}, event model.Event) error {
	if !c.HasSynced() {
		return errors.New("waiting till full synchronization")
	}
	k, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		glog.V(2).Infof("Error retrieving key: %v", err)
	} else {
		glog.V(6).Infof("Event %s: key %#v", event, k)
	}
	return nil
}

func (c *Controller) createInformer(
	o runtime.Object,
	resyncPeriod time.Duration,
	lf cache.ListFunc,
	wf cache.WatchFunc) cacheHandler {
	handler := &ChainHandler{funcs: []Handler{c.notify}}

	// TODO: finer-grained index (perf)
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{ListFunc: lf, WatchFunc: wf}, o,
		resyncPeriod, cache.Indexers{})

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			// TODO: filtering functions to skip over un-referenced resources (perf)
			AddFunc: func(obj interface{}) {
				c.queue.Push(Task{handler: handler.Apply, obj: obj, event: model.EventAdd})
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					c.queue.Push(Task{handler: handler.Apply, obj: cur, event: model.EventUpdate})
				}
			},
			DeleteFunc: func(obj interface{}) {
				c.queue.Push(Task{handler: handler.Apply, obj: obj, event: model.EventDelete})
			},
		})

	return cacheHandler{informer: informer, handler: handler}
}

// HasSynced returns true after the initial state synchronization
func (c *Controller) HasSynced() bool {
	if !c.services.informer.HasSynced() ||
		!c.endpoints.informer.HasSynced() ||
		!c.pods.informer.HasSynced() {
		return false
	}

	return true
}

// Run all controllers until a signal is received
func (c *Controller) Run(stop <-chan struct{}) {
	go c.queue.Run(stop)
	go c.services.informer.Run(stop)
	go c.endpoints.informer.Run(stop)
	go c.pods.informer.Run(stop)

	<-stop
	glog.V(2).Info("Controller terminated")
}

// Services implements a service catalog operation
func (c *Controller) Services() ([]*model.Service, error) {
	list := c.services.informer.GetStore().List()
	out := make([]*model.Service, 0, len(list))

	for _, item := range list {
		if svc := convertService(*item.(*v1.Service), c.domainSuffix); svc != nil {
			out = append(out, svc)
		}
	}
	return out, nil
}

// GetService implements a service catalog operation
func (c *Controller) GetService(hostname string) (*model.Service, error) {
	name, namespace, err := parseHostname(hostname)
	if err != nil {
		glog.V(2).Infof("GetService(%s) => error %v", hostname, err)
		return nil, err
	}
	item, exists := c.serviceByKey(name, namespace)
	if !exists {
		return nil, nil
	}

	svc := convertService(*item, c.domainSuffix)
	return svc, nil
}

// serviceByKey retrieves a service by name and namespace
func (c *Controller) serviceByKey(name, namespace string) (*v1.Service, bool) {
	item, exists, err := c.services.informer.GetStore().GetByKey(KeyFunc(name, namespace))
	if err != nil {
		glog.V(2).Infof("serviceByKey(%s, %s) => error %v", name, namespace, err)
		return nil, false
	}
	if !exists {
		return nil, false
	}
	return item.(*v1.Service), true
}

// GetPodAZ retrieves the AZ for a pod.
func (c *Controller) GetPodAZ(pod *v1.Pod) (string, bool) {
	// NodeName is set by the scheduler after the pod is created
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#late-initialization
	node, exists, err := c.nodes.informer.GetStore().GetByKey(pod.Spec.NodeName)
	if !exists || err != nil {
		return "", false
	}
	region, exists := node.(*v1.Node).Labels[NodeRegionLabel]
	if !exists {
		return "", false
	}
	zone, exists := node.(*v1.Node).Labels[NodeZoneLabel]
	if !exists {
		return "", false
	}

	return fmt.Sprintf("%v/%v", region, zone), true
}

// ManagementPorts implements a service catalog operation
func (c *Controller) ManagementPorts(addr string) model.PortList {
	pod, exists := c.pods.getPodByIP(addr)
	if !exists {
		return nil
	}

	managementPorts, err := convertProbesToPorts(&pod.Spec)

	if err != nil {
		glog.V(2).Infof("Error while parsing liveliness and readiness probe ports for %s => %v", addr, err)
	}

	// We continue despite the error because healthCheckPorts could return a partial
	// list of management ports
	return managementPorts
}

// Instances implements a service catalog operation
func (c *Controller) Instances(hostname string, ports []string,
	labelsList model.LabelsCollection) ([]*model.ServiceInstance, error) {
	// Get actual service by name
	name, namespace, err := parseHostname(hostname)
	if err != nil {
		glog.V(2).Infof("parseHostname(%s) => error %v", hostname, err)
		return nil, err
	}

	item, exists := c.serviceByKey(name, namespace)
	if !exists {
		return nil, nil
	}

	// Locate all ports in the actual service
	svc := convertService(*item, c.domainSuffix)
	if svc == nil {
		return nil, nil
	}
	svcPorts := make(map[string]*model.Port)
	for _, port := range ports {
		if svcPort, exists := svc.Ports.Get(port); exists {
			svcPorts[port] = svcPort
		}
	}

	// TODO: single port service missing name
	for _, item := range c.endpoints.informer.GetStore().List() {
		ep := *item.(*v1.Endpoints)
		if ep.Name == name && ep.Namespace == namespace {
			var out []*model.ServiceInstance
			for _, ss := range ep.Subsets {
				for _, ea := range ss.Addresses {
					labels, _ := c.pods.labelsByIP(ea.IP)
					// check that one of the input labels is a subset of the labels
					if !labelsList.HasSubsetOf(labels) {
						continue
					}

					pod, exists := c.pods.getPodByIP(ea.IP)
					az, sa := "", ""
					if exists {
						az, _ = c.GetPodAZ(pod)
						sa = kubeToIstioServiceAccount(pod.Spec.ServiceAccountName, pod.GetNamespace(), c.domainSuffix)
					}

					// identify the port by name
					for _, port := range ss.Ports {
						if svcPort, exists := svcPorts[port.Name]; exists {
							out = append(out, &model.ServiceInstance{
								Endpoint: model.NetworkEndpoint{
									Address:     ea.IP,
									Port:        int(port.Port),
									ServicePort: svcPort,
								},
								Service:          svc,
								Labels:           labels,
								AvailabilityZone: az,
								ServiceAccount:   sa,
							})
						}
					}
				}
			}
			return out, nil
		}
	}
	return nil, nil
}

// HostInstances implements a service catalog operation
func (c *Controller) HostInstances(addrs map[string]bool) ([]*model.ServiceInstance, error) {
	var out []*model.ServiceInstance
	for _, item := range c.endpoints.informer.GetStore().List() {
		ep := *item.(*v1.Endpoints)
		for _, ss := range ep.Subsets {
			for _, ea := range ss.Addresses {
				if addrs[ea.IP] {
					item, exists := c.serviceByKey(ep.Name, ep.Namespace)
					if !exists {
						continue
					}
					svc := convertService(*item, c.domainSuffix)
					if svc == nil {
						continue
					}
					for _, port := range ss.Ports {
						svcPort, exists := svc.Ports.Get(port.Name)
						if !exists {
							continue
						}
						labels, _ := c.pods.labelsByIP(ea.IP)
						pod, exists := c.pods.getPodByIP(ea.IP)
						az, sa := "", ""
						if exists {
							az, _ = c.GetPodAZ(pod)
							sa = kubeToIstioServiceAccount(pod.Spec.ServiceAccountName, pod.GetNamespace(), c.domainSuffix)
						}
						out = append(out, &model.ServiceInstance{
							Endpoint: model.NetworkEndpoint{
								Address:     ea.IP,
								Port:        int(port.Port),
								ServicePort: svcPort,
							},
							Service:          svc,
							Labels:           labels,
							AvailabilityZone: az,
							ServiceAccount:   sa,
						})
					}
				}
			}
		}
	}
	return out, nil
}

func (c *Controller) WorkloadInstances(id string) ([]model.ServiceInstance, error) {
	out := make([]model.ServiceInstance, 0)

	elt, exists, err := c.pods.informer.GetStore().GetByKey(id)
	if err != nil {
		return out, err
	}
	if !exists {
		return out, nil
	}
	pod := elt.(*v1.Pod)

	for _, item := range c.endpoints.informer.GetStore().List() {
		ep := *item.(*v1.Endpoints)
		for _, ss := range ep.Subsets {
			for _, ea := range ss.Addresses {
				if ea.IP == pod.Status.PodIP {
					item, exists := c.serviceByKey(ep.Name, ep.Namespace)
					if !exists {
						continue
					}
					svc := convertService(*item, c.domainSuffix)
					if svc == nil {
						continue
					}
					labels := convertLabels(pod.ObjectMeta)
					az, _ := c.GetPodAZ(pod)
					sa := kubeToIstioServiceAccount(pod.Spec.ServiceAccountName, pod.GetNamespace(), c.domainSuffix)

					for _, port := range ss.Ports {
						svcPort, exists := svc.Ports.Get(port.Name)
						if !exists {
							continue
						}
						out = append(out, model.ServiceInstance{
							Endpoint: model.NetworkEndpoint{
								Address:     ea.IP,
								Port:        int(port.Port),
								ServicePort: svcPort,
							},
							Service:          svc,
							Labels:           labels,
							AvailabilityZone: az,
							ServiceAccount:   sa,
						})
					}
				}
			}
		}
	}
	return out, nil
}

// GetIstioServiceAccounts returns the Istio service accounts running a serivce
// hostname. Each service account is encoded according to the SPIFFE VSID spec.
// For example, a service account named "bar" in namespace "foo" is encoded as
// "spiffe://cluster.local/ns/foo/sa/bar".
func (c *Controller) GetIstioServiceAccounts(hostname string, ports []string) []string {
	saSet := make(map[string]bool)

	// Get the service accounts running service within Kubernetes. This is reflected by the pods that
	// the service is deployed on, and the service accounts of the pods.
	instances, err := c.Instances(hostname, ports, model.LabelsCollection{})
	if err != nil {
		glog.Warningf("Instances(%s) error: %v", hostname, err)
		return nil
	}
	for _, si := range instances {
		if si.ServiceAccount != "" {
			saSet[si.ServiceAccount] = true
		}
	}

	// Get the service accounts running the service, if it is deployed on VMs. This is retrieved
	// from the service annotation explicitly set by the operators.
	svc, err := c.GetService(hostname)
	if err != nil {
		glog.Warningf("GetService(%s) error: %v", hostname, err)
		return nil
	}
	if svc == nil {
		glog.V(2).Infof("GetService(%s) error: service does not exist", hostname)
		return nil
	}
	for _, serviceAccount := range svc.ServiceAccounts {
		sa := serviceAccount
		saSet[sa] = true
	}

	saArray := make([]string, 0, len(saSet))
	for sa := range saSet {
		saArray = append(saArray, sa)
	}

	return saArray
}

// AppendServiceHandler implements a service catalog operation
func (c *Controller) AppendServiceHandler(f func(*model.Service, model.Event)) error {
	c.services.handler.Append(func(obj interface{}, event model.Event) error {
		svc := *obj.(*v1.Service)

		// Do not handle "kube-system" services
		if svc.Namespace == meta_v1.NamespaceSystem {
			return nil
		}

		glog.V(2).Infof("Handle service %s in namespace %s", svc.Name, svc.Namespace)

		if svcConv := convertService(svc, c.domainSuffix); svcConv != nil {
			f(svcConv, event)
		}
		return nil
	})
	return nil
}

// AppendInstanceHandler implements a service catalog operation
func (c *Controller) AppendInstanceHandler(f func(*model.ServiceInstance, model.Event)) error {
	c.endpoints.handler.Append(func(obj interface{}, event model.Event) error {
		ep := *obj.(*v1.Endpoints)

		// Do not handle "kube-system" endpoints
		if ep.Namespace == meta_v1.NamespaceSystem {
			return nil
		}

		glog.V(2).Infof("Handle endpoint %s in namespace %s", ep.Name, ep.Namespace)
		if item, exists := c.serviceByKey(ep.Name, ep.Namespace); exists {
			if svc := convertService(*item, c.domainSuffix); svc != nil {
				// TODO: we're passing an incomplete instance to the
				// handler since endpoints is an aggregate structure
				f(&model.ServiceInstance{Service: svc}, event)
			}
		}
		return nil
	})
	return nil
}
