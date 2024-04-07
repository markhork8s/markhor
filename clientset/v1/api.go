package v1

import (
	v1 "sops_k8s/api/types/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type SopsV1Interface interface {
	SopsSecrets(namespace string) SopsSecretInterface
}

type SopsV1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*SopsV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &SopsV1Client{restClient: client}, nil
}

func (c *SopsV1Client) SopsSecrets(namespace string) SopsSecretInterface {
	return &sopsSecretClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
