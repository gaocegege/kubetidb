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
	v1alpha1 "github.com/gaocegege/kubetidb/pkg/apis/tidb/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTiDBs implements TiDBInterface
type FakeTiDBs struct {
	Fake *FakeKubetidbV1alpha1
	ns   string
}

var tidbsResource = schema.GroupVersionResource{Group: "kubetidb.gaocegege.com", Version: "v1alpha1", Resource: "tidbs"}

var tidbsKind = schema.GroupVersionKind{Group: "kubetidb.gaocegege.com", Version: "v1alpha1", Kind: "TiDB"}

// Get takes name of the tiDB, and returns the corresponding tiDB object, and an error if there is any.
func (c *FakeTiDBs) Get(name string, options v1.GetOptions) (result *v1alpha1.TiDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(tidbsResource, c.ns, name), &v1alpha1.TiDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TiDB), err
}

// List takes label and field selectors, and returns the list of TiDBs that match those selectors.
func (c *FakeTiDBs) List(opts v1.ListOptions) (result *v1alpha1.TiDBList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(tidbsResource, tidbsKind, c.ns, opts), &v1alpha1.TiDBList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TiDBList{}
	for _, item := range obj.(*v1alpha1.TiDBList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested tiDBs.
func (c *FakeTiDBs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(tidbsResource, c.ns, opts))

}

// Create takes the representation of a tiDB and creates it.  Returns the server's representation of the tiDB, and an error, if there is any.
func (c *FakeTiDBs) Create(tiDB *v1alpha1.TiDB) (result *v1alpha1.TiDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(tidbsResource, c.ns, tiDB), &v1alpha1.TiDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TiDB), err
}

// Update takes the representation of a tiDB and updates it. Returns the server's representation of the tiDB, and an error, if there is any.
func (c *FakeTiDBs) Update(tiDB *v1alpha1.TiDB) (result *v1alpha1.TiDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(tidbsResource, c.ns, tiDB), &v1alpha1.TiDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TiDB), err
}

// Delete takes name of the tiDB and deletes it. Returns an error if one occurs.
func (c *FakeTiDBs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(tidbsResource, c.ns, name), &v1alpha1.TiDB{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTiDBs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(tidbsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.TiDBList{})
	return err
}

// Patch applies the patch and returns the patched tiDB.
func (c *FakeTiDBs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.TiDB, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(tidbsResource, c.ns, name, data, subresources...), &v1alpha1.TiDB{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.TiDB), err
}
