package v1

import (
	"context"

	v1 "github.com/civts/markhor/pkg/api/types/v1"

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
}

const msecretsResource = "markhorsecrets"

func (c *markhorSecretClient) List(opts metav1.ListOptions) (*v1.MarkhorSecretList, error) {
	result := v1.MarkhorSecretList{}
	err := c.restClient.
		Get().
		Resource(msecretsResource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *markhorSecretClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource(msecretsResource).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
