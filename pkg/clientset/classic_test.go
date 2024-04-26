package clientset

import (
	"testing"

	"k8s.io/client-go/rest"
)

func TestGetK8sClient(t *testing.T) {
	config := &rest.Config{}
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

	kubeconfig := "" // Mock kubeconfig path for in-cluster configuration
	GetK8sConfig(kubeconfig)

	t.Error("Expected GetK8sConfig to panic, before getting here")
}

func TestGetK8sConfigYieldsClientWithValidConfig(t *testing.T) {
	f, err := createTempFile()
	if err != nil {
		t.Fatal("Could not write to temp file:", err)
	}
	f.WriteString(sampleKubeconfig)
	f.Sync()
	defer removeTempFile(f)

	client := GetK8sConfig(f.Name())
	if client == nil {
		t.Error("Expected GetK8sConfig to procude a valid client")
	}
}
