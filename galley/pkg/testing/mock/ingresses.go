// Copyright 2019 Istio Authors
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

package mock

import (
	"fmt"
	"sync"

	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	extensionsv1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
)

var _ extensionsv1.IngressInterface = &ingressImpl{}

type ingressImpl struct {
	mux       sync.Mutex
	ingresses map[string]*v1beta1.Ingress
	watches   Watches
}

func newIngressInterface() extensionsv1.IngressInterface {
	return &ingressImpl{
		ingresses: make(map[string]*v1beta1.Ingress),
	}
}

func (i *ingressImpl) Create(obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.ingresses[obj.Name] = obj

	i.watches.Send(watch.Event{
		Type:   watch.Added,
		Object: obj,
	})
	return obj, nil
}

func (i *ingressImpl) Update(obj *v1beta1.Ingress) (*v1beta1.Ingress, error) {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.ingresses[obj.Name] = obj

	i.watches.Send(watch.Event{
		Type:   watch.Modified,
		Object: obj,
	})
	return obj, nil
}

func (i *ingressImpl) Delete(name string, options *v1.DeleteOptions) error {
	i.mux.Lock()
	defer i.mux.Unlock()

	obj := i.ingresses[name]
	if obj == nil {
		return fmt.Errorf("unable to delete ingress %s", name)
	}

	delete(i.ingresses, name)

	i.watches.Send(watch.Event{
		Type:   watch.Deleted,
		Object: obj,
	})
	return nil
}

func (i *ingressImpl) List(opts v1.ListOptions) (*v1beta1.IngressList, error) {
	i.mux.Lock()
	defer i.mux.Unlock()

	out := &v1beta1.IngressList{}

	for _, v := range i.ingresses {
		out.Items = append(out.Items, *v)
	}

	return out, nil
}

func (i *ingressImpl) Watch(opts v1.ListOptions) (watch.Interface, error) {
	i.mux.Lock()
	defer i.mux.Unlock()

	w := NewWatch()
	i.watches = append(i.watches, w)

	// Send add events for all current resources.
	for _, i := range i.ingresses {
		w.Send(watch.Event{
			Type:   watch.Added,
			Object: i,
		})
	}

	return w, nil
}

func (i *ingressImpl) UpdateStatus(*v1beta1.Ingress) (*v1beta1.Ingress, error) {
	panic("not implemented")

}

func (i *ingressImpl) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	panic("not implemented")

}

func (i *ingressImpl) Get(name string, options v1.GetOptions) (*v1beta1.Ingress, error) {
	panic("not implemented")

}

func (i *ingressImpl) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Ingress, err error) {
	panic("not implemented")

}
