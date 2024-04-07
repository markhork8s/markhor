package v1

import (
	"context"
	v1 "sops_k8s/api/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type SopsSecretInterface interface {
	List(opts metav1.ListOptions) (*v1.SopsSecretList, error)
	// Get(name string, options metav1.GetOptions) (*v1.SopsSecret, error)
	// Create(*v1.SopsSecret) (*v1.SopsSecret, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type sopsSecretClient struct {
	restClient rest.Interface
	ns         string
}

func (c *sopsSecretClient) List(opts metav1.ListOptions) (*v1.SopsSecretList, error) {
	result := v1.SopsSecretList{}
	err := c.restClient.
		Get().
		// Namespace(c.ns).
		Resource("sopssecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *sopsSecretClient) Get(name string, opts metav1.GetOptions) (*v1.SopsSecret, error) {
	result := v1.SopsSecret{}
	err := c.restClient.
		Get().
		//Namespace(c.ns).
		Resource("sopssecrets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *sopsSecretClient) Create(sopsSecret *v1.SopsSecret) (*v1.SopsSecret, error) {
	result := v1.SopsSecret{}
	err := c.restClient.
		Post().
		//Namespace(c.ns).
		Resource("sopssecrets").
		Body(sopsSecret).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *sopsSecretClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		//Namespace(c.ns).
		Resource("sopssecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
