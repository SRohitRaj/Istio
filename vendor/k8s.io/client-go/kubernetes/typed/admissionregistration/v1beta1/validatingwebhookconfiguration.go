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

package v1beta1

import (
	v1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// ValidatingWebhookConfigurationsGetter has a method to return a ValidatingWebhookConfigurationInterface.
// A group's client should implement this interface.
type ValidatingWebhookConfigurationsGetter interface {
	ValidatingWebhookConfigurations() ValidatingWebhookConfigurationInterface
}

// ValidatingWebhookConfigurationInterface has methods to work with ValidatingWebhookConfiguration resources.
type ValidatingWebhookConfigurationInterface interface {
	Create(*v1beta1.ValidatingWebhookConfiguration) (*v1beta1.ValidatingWebhookConfiguration, error)
	Update(*v1beta1.ValidatingWebhookConfiguration) (*v1beta1.ValidatingWebhookConfiguration, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.ValidatingWebhookConfiguration, error)
	List(opts v1.ListOptions) (*v1beta1.ValidatingWebhookConfigurationList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.ValidatingWebhookConfiguration, err error)
	ValidatingWebhookConfigurationExpansion
}

// validatingWebhookConfigurations implements ValidatingWebhookConfigurationInterface
type validatingWebhookConfigurations struct {
	client rest.Interface
}

// newValidatingWebhookConfigurations returns a ValidatingWebhookConfigurations
func newValidatingWebhookConfigurations(c *AdmissionregistrationV1beta1Client) *validatingWebhookConfigurations {
	return &validatingWebhookConfigurations{
		client: c.RESTClient(),
	}
}

// Get takes name of the validatingWebhookConfiguration, and returns the corresponding validatingWebhookConfiguration object, and an error if there is any.
func (c *validatingWebhookConfigurations) Get(name string, options v1.GetOptions) (result *v1beta1.ValidatingWebhookConfiguration, err error) {
	result = &v1beta1.ValidatingWebhookConfiguration{}
	err = c.client.Get().
		Resource("validatingwebhookconfigurations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ValidatingWebhookConfigurations that match those selectors.
func (c *validatingWebhookConfigurations) List(opts v1.ListOptions) (result *v1beta1.ValidatingWebhookConfigurationList, err error) {
	result = &v1beta1.ValidatingWebhookConfigurationList{}
	err = c.client.Get().
		Resource("validatingwebhookconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested validatingWebhookConfigurations.
func (c *validatingWebhookConfigurations) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("validatingwebhookconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a validatingWebhookConfiguration and creates it.  Returns the server's representation of the validatingWebhookConfiguration, and an error, if there is any.
func (c *validatingWebhookConfigurations) Create(validatingWebhookConfiguration *v1beta1.ValidatingWebhookConfiguration) (result *v1beta1.ValidatingWebhookConfiguration, err error) {
	result = &v1beta1.ValidatingWebhookConfiguration{}
	err = c.client.Post().
		Resource("validatingwebhookconfigurations").
		Body(validatingWebhookConfiguration).
		Do().
		Into(result)
	return
}

// Update takes the representation of a validatingWebhookConfiguration and updates it. Returns the server's representation of the validatingWebhookConfiguration, and an error, if there is any.
func (c *validatingWebhookConfigurations) Update(validatingWebhookConfiguration *v1beta1.ValidatingWebhookConfiguration) (result *v1beta1.ValidatingWebhookConfiguration, err error) {
	result = &v1beta1.ValidatingWebhookConfiguration{}
	err = c.client.Put().
		Resource("validatingwebhookconfigurations").
		Name(validatingWebhookConfiguration.Name).
		Body(validatingWebhookConfiguration).
		Do().
		Into(result)
	return
}

// Delete takes name of the validatingWebhookConfiguration and deletes it. Returns an error if one occurs.
func (c *validatingWebhookConfigurations) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("validatingwebhookconfigurations").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *validatingWebhookConfigurations) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("validatingwebhookconfigurations").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched validatingWebhookConfiguration.
func (c *validatingWebhookConfigurations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.ValidatingWebhookConfiguration, err error) {
	result = &v1beta1.ValidatingWebhookConfiguration{}
	err = c.client.Patch(pt).
		Resource("validatingwebhookconfigurations").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
