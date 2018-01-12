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

package v1alpha1

import (
	v1alpha1 "github.com/gaocegege/kubetidb/pkg/apis/tidb/v1alpha1"
	scheme "github.com/gaocegege/kubetidb/pkg/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TiDBsGetter has a method to return a TiDBInterface.
// A group's client should implement this interface.
type TiDBsGetter interface {
	TiDBs(namespace string) TiDBInterface
}

// TiDBInterface has methods to work with TiDB resources.
type TiDBInterface interface {
	Create(*v1alpha1.TiDB) (*v1alpha1.TiDB, error)
	Update(*v1alpha1.TiDB) (*v1alpha1.TiDB, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.TiDB, error)
	List(opts v1.ListOptions) (*v1alpha1.TiDBList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TiDB, err error)
	TiDBExpansion
}

// tiDBs implements TiDBInterface
type tiDBs struct {
	client rest.Interface
	ns     string
}

// newTiDBs returns a TiDBs
func newTiDBs(c *KubetidbV1alpha1Client, namespace string) *tiDBs {
	return &tiDBs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the tiDB, and returns the corresponding tiDB object, and an error if there is any.
func (c *tiDBs) Get(name string, options v1.GetOptions) (result *v1alpha1.TiDB, err error) {
	result = &v1alpha1.TiDB{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tidbs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TiDBs that match those selectors.
func (c *tiDBs) List(opts v1.ListOptions) (result *v1alpha1.TiDBList, err error) {
	result = &v1alpha1.TiDBList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tidbs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tiDBs.
func (c *tiDBs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("tidbs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a tiDB and creates it.  Returns the server's representation of the tiDB, and an error, if there is any.
func (c *tiDBs) Create(tiDB *v1alpha1.TiDB) (result *v1alpha1.TiDB, err error) {
	result = &v1alpha1.TiDB{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("tidbs").
		Body(tiDB).
		Do().
		Into(result)
	return
}

// Update takes the representation of a tiDB and updates it. Returns the server's representation of the tiDB, and an error, if there is any.
func (c *tiDBs) Update(tiDB *v1alpha1.TiDB) (result *v1alpha1.TiDB, err error) {
	result = &v1alpha1.TiDB{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tidbs").
		Name(tiDB.Name).
		Body(tiDB).
		Do().
		Into(result)
	return
}

// Delete takes name of the tiDB and deletes it. Returns an error if one occurs.
func (c *tiDBs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tidbs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tiDBs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tidbs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched tiDB.
func (c *tiDBs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TiDB, err error) {
	result = &v1alpha1.TiDB{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("tidbs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
