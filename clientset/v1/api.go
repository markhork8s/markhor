package v1

import (
	v1 "markhor/api/types/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type MarkhorV1Interface interface {
	MarkhorSecrets(namespace string) MarkhorSecretInterface
}

type MarkhorV1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*MarkhorV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &MarkhorV1Client{restClient: client}, nil
}

func (c *MarkhorV1Client) MarkhorSecrets(namespace string) MarkhorSecretInterface {
	return &markhorSecretClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
