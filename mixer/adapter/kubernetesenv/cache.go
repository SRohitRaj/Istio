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

package kubernetesenv

import (
	"errors"
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"istio.io/istio/mixer/adapter/kubernetesenv/pkg/apis/networking/v1alpha3"
	versioned "istio.io/istio/mixer/adapter/kubernetesenv/pkg/client/clientset/versioned"
	externalversions "istio.io/istio/mixer/adapter/kubernetesenv/pkg/client/informers/externalversions"
	"istio.io/istio/mixer/pkg/adapter"
)

type (
	// internal interface used to support testing
	cacheController interface {
		Run(<-chan struct{})
		Pod(string) (*v1.Pod, bool)
		Workload(*v1.Pod) workload
		ServiceEntry(string) (*v1alpha3.ServiceEntry, bool)
		HasSynced() bool
		StopControlChannel()
	}

	controllerImpl struct {
		env      adapter.Env
		stopChan chan struct{}
		pods     cache.SharedIndexInformer
		rs       cache.SharedIndexInformer
		rc       cache.SharedIndexInformer
		se       cache.SharedIndexInformer
	}

	workload struct {
		uid, name, namespace string
		selfLinkURL          string
	}
)

func podIP(obj interface{}) ([]string, error) {
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return nil, errors.New("object is not a pod")
	}
	ip := pod.Status.PodIP
	if ip == "" {
		return nil, nil
	}
	return []string{ip}, nil
}

func endpointIP(obj interface{}) ([]string, error) {
	se, ok := obj.(*v1alpha3.ServiceEntry)
	if !ok {
		return nil, errors.New("object is not a service entry")
	}
	var ips []string
	for _, endpoint := range se.Spec.Endpoints {
		ips = append(ips, endpoint.Address)
	}
	return ips, nil
}

// Responsible for setting up the cacheController, based on the supplied client.
// It configures the index informer to list/watch k8sCache and send update events
// to a mutations channel for processing (in this case, logging).
func newCacheController(clientset kubernetes.Interface, serviceEntries versioned.Interface, refreshDuration time.Duration, env adapter.Env, stopChan chan struct{}) cacheController {
	sharedInformers := informers.NewSharedInformerFactory(clientset, refreshDuration)
	podInformer := sharedInformers.Core().V1().Pods().Informer()
	podInformer.AddIndexers(cache.Indexers{
		"ip": podIP,
	})
	versionedInformers := externalversions.NewSharedInformerFactory(serviceEntries, refreshDuration)
	seInformer := versionedInformers.Networking().V1alpha3().ServiceEntries().Informer()
	seInformer.AddIndexers(cache.Indexers{
		"endpoints": endpointIP,
	})

	return &controllerImpl{
		env:      env,
		stopChan: stopChan,
		pods:     podInformer,
		rs:       sharedInformers.Apps().V1().ReplicaSets().Informer(),
		rc:       sharedInformers.Core().V1().ReplicationControllers().Informer(),
		se:       seInformer,
	}
}

func (c *controllerImpl) StopControlChannel() {
	close(c.stopChan)
}

func (c *controllerImpl) HasSynced() bool {
	return c.pods.HasSynced() && c.rs.HasSynced() && c.rc.HasSynced() && c.se.HasSynced()
}

func (c *controllerImpl) Run(stop <-chan struct{}) {
	c.env.ScheduleDaemon(func() { c.pods.Run(stop) })
	c.env.ScheduleDaemon(func() { c.rs.Run(stop) })
	c.env.ScheduleDaemon(func() { c.rc.Run(stop) })
	c.env.ScheduleDaemon(func() { c.se.Run(stop) })
	<-stop
	// TODO: logging?
}

// Pod returns a k8s Pod object that corresponds to the supplied key, if one
// exists (and is known to the store). Keys are expected in the form of:
// namespace/name or IP address (example: "default/curl-2421989462-b2g2d.default").
func (c *controllerImpl) Pod(podKey string) (*v1.Pod, bool) {
	indexer := c.pods.GetIndexer()
	objs, err := indexer.ByIndex("ip", podKey)
	if err != nil {
		return nil, false
	}
	if len(objs) > 0 {
		pod, ok := objs[0].(*v1.Pod)
		if !ok {
			return nil, false
		}
		return pod, true
	}
	item, exists, err := indexer.GetByKey(podKey)
	if !exists || err != nil {
		return nil, false
	}
	return item.(*v1.Pod), true
}

func key(namespace, name string) string {
	return namespace + "/" + name
}

func (c *controllerImpl) Workload(pod *v1.Pod) workload {
	wl := workload{name: pod.Name, namespace: pod.Namespace, selfLinkURL: pod.SelfLink}
	if owner, found := c.rootController(&pod.ObjectMeta); found {
		wl.name = owner.Name
		wl.selfLinkURL = fmt.Sprintf("kubernetes://apis/%s/namespaces/%s/%ss/%s", owner.APIVersion, pod.Namespace, strings.ToLower(owner.Kind), owner.Name)
	}
	wl.uid = "istio://" + wl.namespace + "/workloads/" + wl.name
	return wl
}

func (c *controllerImpl) rootController(obj *metav1.ObjectMeta) (metav1.OwnerReference, bool) {
	for _, ref := range obj.OwnerReferences {
		if *ref.Controller {
			switch ref.Kind {
			case "ReplicaSet":
				indexer := c.rs.GetIndexer()
				if rs, found := c.objectMeta(indexer, key(obj.Namespace, ref.Name)); found {
					if rootRef, ok := c.rootController(rs); ok {
						return rootRef, true
					}
				}
			case "ReplicationController":
				indexer := c.rc.GetIndexer()
				if rc, found := c.objectMeta(indexer, key(obj.Namespace, ref.Name)); found {
					if rootRef, ok := c.rootController(rc); ok {
						return rootRef, true
					}
				}
			}

			return ref, true
		}
	}
	return metav1.OwnerReference{}, false
}

func (c *controllerImpl) objectMeta(keyGetter cache.KeyGetter, key string) (*metav1.ObjectMeta, bool) {
	item, exists, err := keyGetter.GetByKey(key)
	if !exists || err != nil {
		return nil, false
	}
	switch v := item.(type) {
	case *v1.ReplicationController:
		return &v.ObjectMeta, true
	case *appsv1.ReplicaSet:
		return &v.ObjectMeta, true
	}
	return nil, false
}

// ServiceEntry returns an Istio ServiceEntry object that corresponds to the supplied IP, if one
// exists (and is known to the store).
func (c *controllerImpl) ServiceEntry(ipAddr string) (*v1alpha3.ServiceEntry, bool) {
	indexer := c.se.GetIndexer()
	objs, err := indexer.ByIndex("endpoints", ipAddr)
	if err != nil || len(objs) == 0 {
		return nil, false
	}
	se, ok := objs[0].(*v1alpha3.ServiceEntry)
	if !ok {
		return nil, false
	}
	return se, true
}
