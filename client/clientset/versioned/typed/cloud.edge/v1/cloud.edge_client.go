// Copyright @ 2022 Lenovo. All rights reserved
// Confidential and Restricted
// This file was generated using one of kubebuilder templated ( https://github.com/kubernetes-sigs/kubebuilder/tree/master/pkg/plugins/golang/v2/scaffolds/internal/templates) distributed under the terms of Apache-2.0 license.

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "igitlab.lenovo.com/lecp/mec/bailu/api/v1"
	"igitlab.lenovo.com/lecp/mec/bailu/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type CloudV1Interface interface {
	RESTClient() rest.Interface
	NodeManagementsGetter
}

// CloudV1Client is used to interact with features provided by the cloud.edge group.
type CloudV1Client struct {
	restClient rest.Interface
}

func (c *CloudV1Client) NodeManagements() NodeManagementInterface {
	return newNodeManagements(c)
}

// NewForConfig creates a new CloudV1Client for the given config.
func NewForConfig(c *rest.Config) (*CloudV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &CloudV1Client{client}, nil
}

// NewForConfigOrDie creates a new CloudV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CloudV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new CloudV1Client for the given RESTClient.
func New(c rest.Interface) *CloudV1Client {
	return &CloudV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *CloudV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}