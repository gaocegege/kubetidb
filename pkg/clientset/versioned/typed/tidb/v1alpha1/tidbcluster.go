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

// TiDBClustersGetter has a method to return a TiDBClusterInterface.
// A group's client should implement this interface.
type TiDBClustersGetter interface {
	TiDBClusters(namespace string) TiDBClusterInterface
}

// TiDBClusterInterface has methods to work with TiDBCluster resources.
type TiDBClusterInterface interface {
	Create(*v1alpha1.TiDBCluster) (*v1alpha1.TiDBCluster, error)
	Update(*v1alpha1.TiDBCluster) (*v1alpha1.TiDBCluster, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.TiDBCluster, error)
	List(opts v1.ListOptions) (*v1alpha1.TiDBClusterList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TiDBCluster, err error)
	TiDBClusterExpansion
}

// tiDBClusters implements TiDBClusterInterface
type tiDBClusters struct {
	client rest.Interface
	ns     string
}

// newTiDBClusters returns a TiDBClusters
func newTiDBClusters(c *KubetidbV1alpha1Client, namespace string) *tiDBClusters {
	return &tiDBClusters{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the tiDBCluster, and returns the corresponding tiDBCluster object, and an error if there is any.
func (c *tiDBClusters) Get(name string, options v1.GetOptions) (result *v1alpha1.TiDBCluster, err error) {
	result = &v1alpha1.TiDBCluster{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tidbclusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TiDBClusters that match those selectors.
func (c *tiDBClusters) List(opts v1.ListOptions) (result *v1alpha1.TiDBClusterList, err error) {
	result = &v1alpha1.TiDBClusterList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tidbclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tiDBClusters.
func (c *tiDBClusters) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("tidbclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a tiDBCluster and creates it.  Returns the server's representation of the tiDBCluster, and an error, if there is any.
func (c *tiDBClusters) Create(tiDBCluster *v1alpha1.TiDBCluster) (result *v1alpha1.TiDBCluster, err error) {
	result = &v1alpha1.TiDBCluster{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("tidbclusters").
		Body(tiDBCluster).
		Do().
		Into(result)
	return
}

// Update takes the representation of a tiDBCluster and updates it. Returns the server's representation of the tiDBCluster, and an error, if there is any.
func (c *tiDBClusters) Update(tiDBCluster *v1alpha1.TiDBCluster) (result *v1alpha1.TiDBCluster, err error) {
	result = &v1alpha1.TiDBCluster{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tidbclusters").
		Name(tiDBCluster.Name).
		Body(tiDBCluster).
		Do().
		Into(result)
	return
}

// Delete takes name of the tiDBCluster and deletes it. Returns an error if one occurs.
func (c *tiDBClusters) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tidbclusters").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tiDBClusters) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tidbclusters").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched tiDBCluster.
func (c *tiDBClusters) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TiDBCluster, err error) {
	result = &v1alpha1.TiDBCluster{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("tidbclusters").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
