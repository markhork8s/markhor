package clientset

import (
	"testing"

	"k8s.io/client-go/rest"
)

func TestGetK8sClient(t *testing.T) {
	config := &rest.Config{} // Mock config for testing
	clientSet := GetK8sClient(config)

	if clientSet == nil {
		t.Error("Expected clientSet to be initialized, but it was nil")
	}
}

func TestGetK8sConfigPanicsOutsideACluster(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	// The following is the code under test
	kubeconfig := "" // Mock kubeconfig for in-cluster configuration
	GetK8sConfig(kubeconfig)

	t.Error("Expected GetK8sConfig to panic, before getting here")
}
