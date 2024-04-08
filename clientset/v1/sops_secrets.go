package v1

import (
	"context"
	v1 "markhor/api/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type MarkhorSecretInterface interface {
	List(opts metav1.ListOptions) (*v1.MarkhorSecretList, error)
	// Get(name string, options metav1.GetOptions) (*v1.MarkhorSecret, error)
	// Create(*v1.MarkhorSecret) (*v1.MarkhorSecret, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type markhorSecretClient struct {
	restClient rest.Interface
	ns         string
}

func (c *markhorSecretClient) List(opts metav1.ListOptions) (*v1.MarkhorSecretList, error) {
	result := v1.MarkhorSecretList{}
	err := c.restClient.
		Get().
		// Namespace(c.ns).
		Resource("markhorsecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *markhorSecretClient) Get(name string, opts metav1.GetOptions) (*v1.MarkhorSecret, error) {
	result := v1.MarkhorSecret{}
	err := c.restClient.
		Get().
		//Namespace(c.ns).
		Resource("markhorsecrets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *markhorSecretClient) Create(markhorSecret *v1.MarkhorSecret) (*v1.MarkhorSecret, error) {
	result := v1.MarkhorSecret{}
	err := c.restClient.
		Post().
		//Namespace(c.ns).
		Resource("markhorsecrets").
		Body(markhorSecret).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *markhorSecretClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		//Namespace(c.ns).
		Resource("markhorsecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
