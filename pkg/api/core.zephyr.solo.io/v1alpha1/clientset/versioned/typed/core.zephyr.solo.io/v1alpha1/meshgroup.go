/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/solo-io/mesh-projects/pkg/api/core.zephyr.solo.io/v1alpha1"
	scheme "github.com/solo-io/mesh-projects/pkg/api/core.zephyr.solo.io/v1alpha1/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MeshGroupsGetter has a method to return a MeshGroupInterface.
// A group's client should implement this interface.
type MeshGroupsGetter interface {
	MeshGroups(namespace string) MeshGroupInterface
}

// MeshGroupInterface has methods to work with MeshGroup resources.
type MeshGroupInterface interface {
	Create(*v1alpha1.MeshGroup) (*v1alpha1.MeshGroup, error)
	Update(*v1alpha1.MeshGroup) (*v1alpha1.MeshGroup, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.MeshGroup, error)
	List(opts v1.ListOptions) (*v1alpha1.MeshGroupList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MeshGroup, err error)
	MeshGroupExpansion
}

// meshGroups implements MeshGroupInterface
type meshGroups struct {
	client rest.Interface
	ns     string
}

// newMeshGroups returns a MeshGroups
func newMeshGroups(c *CoreV1alpha1Client, namespace string) *meshGroups {
	return &meshGroups{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the meshGroup, and returns the corresponding meshGroup object, and an error if there is any.
func (c *meshGroups) Get(name string, options v1.GetOptions) (result *v1alpha1.MeshGroup, err error) {
	result = &v1alpha1.MeshGroup{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("meshgroups").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MeshGroups that match those selectors.
func (c *meshGroups) List(opts v1.ListOptions) (result *v1alpha1.MeshGroupList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.MeshGroupList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("meshgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested meshGroups.
func (c *meshGroups) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("meshgroups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a meshGroup and creates it.  Returns the server's representation of the meshGroup, and an error, if there is any.
func (c *meshGroups) Create(meshGroup *v1alpha1.MeshGroup) (result *v1alpha1.MeshGroup, err error) {
	result = &v1alpha1.MeshGroup{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("meshgroups").
		Body(meshGroup).
		Do().
		Into(result)
	return
}

// Update takes the representation of a meshGroup and updates it. Returns the server's representation of the meshGroup, and an error, if there is any.
func (c *meshGroups) Update(meshGroup *v1alpha1.MeshGroup) (result *v1alpha1.MeshGroup, err error) {
	result = &v1alpha1.MeshGroup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("meshgroups").
		Name(meshGroup.Name).
		Body(meshGroup).
		Do().
		Into(result)
	return
}

// Delete takes name of the meshGroup and deletes it. Returns an error if one occurs.
func (c *meshGroups) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("meshgroups").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *meshGroups) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("meshgroups").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched meshGroup.
func (c *meshGroups) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MeshGroup, err error) {
	result = &v1alpha1.MeshGroup{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("meshgroups").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}