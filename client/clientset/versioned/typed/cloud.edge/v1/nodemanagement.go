// Copyright @ 2022 Lenovo. All rights reserved
// Confidential and Restricted
// This file was generated using one of kubebuilder templated ( https://github.com/kubernetes-sigs/kubebuilder/tree/master/pkg/plugins/golang/v2/scaffolds/internal/templates) distributed under the terms of Apache-2.0 license.

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	v1 "igitlab.lenovo.com/lecp/mec/bailu/api/v1"
	scheme "igitlab.lenovo.com/lecp/mec/bailu/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// NodeManagementsGetter has a method to return a NodeManagementInterface.
// A group's client should implement this interface.
type NodeManagementsGetter interface {
	NodeManagements() NodeManagementInterface
}

// NodeManagementInterface has methods to work with NodeManagement resources.
type NodeManagementInterface interface {
	Create(ctx context.Context, nodeManagement *v1.NodeManagement, opts metav1.CreateOptions) (*v1.NodeManagement, error)
	Update(ctx context.Context, nodeManagement *v1.NodeManagement, opts metav1.UpdateOptions) (*v1.NodeManagement, error)
	UpdateStatus(ctx context.Context, nodeManagement *v1.NodeManagement, opts metav1.UpdateOptions) (*v1.NodeManagement, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.NodeManagement, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.NodeManagementList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.NodeManagement, err error)
	NodeManagementExpansion
}

// nodeManagements implements NodeManagementInterface
type nodeManagements struct {
	client rest.Interface
}

// newNodeManagements returns a NodeManagements
func newNodeManagements(c *CloudV1Client) *nodeManagements {
	return &nodeManagements{
		client: c.RESTClient(),
	}
}

// Get takes name of the nodeManagement, and returns the corresponding nodeManagement object, and an error if there is any.
func (c *nodeManagements) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.NodeManagement, err error) {
	result = &v1.NodeManagement{}
	err = c.client.Get().
		Resource("nodemanagements").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NodeManagements that match those selectors.
func (c *nodeManagements) List(ctx context.Context, opts metav1.ListOptions) (result *v1.NodeManagementList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.NodeManagementList{}
	err = c.client.Get().
		Resource("nodemanagements").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested nodeManagements.
func (c *nodeManagements) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("nodemanagements").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a nodeManagement and creates it.  Returns the server's representation of the nodeManagement, and an error, if there is any.
func (c *nodeManagements) Create(ctx context.Context, nodeManagement *v1.NodeManagement, opts metav1.CreateOptions) (result *v1.NodeManagement, err error) {
	result = &v1.NodeManagement{}
	err = c.client.Post().
		Resource("nodemanagements").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeManagement).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a nodeManagement and updates it. Returns the server's representation of the nodeManagement, and an error, if there is any.
func (c *nodeManagements) Update(ctx context.Context, nodeManagement *v1.NodeManagement, opts metav1.UpdateOptions) (result *v1.NodeManagement, err error) {
	result = &v1.NodeManagement{}
	err = c.client.Put().
		Resource("nodemanagements").
		Name(nodeManagement.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeManagement).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *nodeManagements) UpdateStatus(ctx context.Context, nodeManagement *v1.NodeManagement, opts metav1.UpdateOptions) (result *v1.NodeManagement, err error) {
	result = &v1.NodeManagement{}
	err = c.client.Put().
		Resource("nodemanagements").
		Name(nodeManagement.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(nodeManagement).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the nodeManagement and deletes it. Returns an error if one occurs.
func (c *nodeManagements) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("nodemanagements").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *nodeManagements) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("nodemanagements").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched nodeManagement.
func (c *nodeManagements) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.NodeManagement, err error) {
	result = &v1.NodeManagement{}
	err = c.client.Patch(pt).
		Resource("nodemanagements").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}