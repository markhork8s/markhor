package clientset

import (
	"os"
	"testing"
)

func TestGetK8sClientsPanicsOutsideACluster(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	kubeconfig := ""
	GetK8sClients(kubeconfig)

	t.Error("Expected GetK8sConfig to panic, before getting here")
}

// This kubeconfig is valid but does not give access to any real cluster
const sampleKubeconfig = `apiVersion: v1
clusters:
- cluster:
    server: https://localhost:0000
  name: nonexistent-cluster
contexts:
- context:
    cluster: nonexistent-cluster
    user: nonexistent-user
  name: nonexistent-context
current-context: nonexistent-context
kind: Config
users:
- name: nonexistent-user
  user: {}
`

func TestGetK8sClientsProducesAConfigWithRealisticKubeconfig(t *testing.T) {
	f, err := createTempFile()
	if err != nil {
		t.Fatal("Could not write to temp file:", err)
	}
	f.WriteString(sampleKubeconfig)
	f.Sync()
	defer removeTempFile(f)

	mClient, clientset := GetK8sClients(f.Name())

	if mClient == nil {
		t.Error("Expected MarkhorV1Client to be initialized, but it was nil")
	}

	if clientset == nil {
		t.Error("Expected Kubernetes clientset to be initialized, but it was nil")
	}
}

func createTempFile() (*os.File, error) {
	tmpDir := os.TempDir()
	return os.CreateTemp(tmpDir, "deleteme_")
}

func removeTempFile(f *os.File) {
	os.Remove(f.Name())
	f.Close()
}
