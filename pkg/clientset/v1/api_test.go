package v1

import (
	"testing"

	v1 "github.com/markhork8s/markhor/pkg/api/types/v1"
	"k8s.io/client-go/rest"
)

func TestNewForConfig(t *testing.T) {
	c := &rest.Config{}
	client, err := NewForConfig(c)

	if err != nil {
		t.Fatalf("Error creating MarkhorV1Client: %v", err)
	}

	if client.restClient == nil {
		t.Fatalf("Expected REST client to be initialized, but it was nil")
	}
	group := client.restClient.APIVersion().Group
	if group != v1.GroupName {
		t.Errorf("Expected REST client to use %s as group, but it used %s", v1.GroupName, group)
	}
	version := client.restClient.APIVersion().Version
	if version != v1.GroupVersion {
		t.Errorf("Expected REST client to use %s as version, but it used %s", v1.GroupVersion, version)
	}
}

func TestMarkhorV1Client_MarkhorSecrets(t *testing.T) {
	c := &rest.Config{}
	client, _ := NewForConfig(c)

	if client == nil {
		t.Fatal("Error creating MarkhorV1Client")
	}

	secretClient := client.MarkhorSecrets()

	if secretClient == nil {
		t.Error("Expected MarkhorSecretInterface to be initialized, but it was nil")
	}
}
