package v1

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestAddKnownTypes(t *testing.T) {
	s := runtime.NewScheme()
	err := addKnownTypes(s)

	if err != nil {
		t.Errorf("Error adding known types to scheme: %v", err)
	}

	gvk := schema.GroupVersionKind{Group: GroupName, Version: GroupVersion, Kind: "MarkhorSecret"}
	if !s.Recognizes(gvk) {
		t.Error("Scheme does not recognize the specified GroupVersionKind for MarkhorSecret")
	}

	gvkList := schema.GroupVersionKind{Group: GroupName, Version: GroupVersion, Kind: "MarkhorSecretList"}
	if !s.Recognizes(gvkList) {
		t.Error("Scheme does not recognize the specified GroupVersionKind for MarkhorSecretList")
	}
}
