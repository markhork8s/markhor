package clientset

import (
	"log/slog"

	v1 "github.com/markhork8s/markhor/pkg/clientset/v1"

	"k8s.io/client-go/kubernetes"
)

func GetK8sClients(kubeconfig string) (*v1.MarkhorV1Client, *kubernetes.Clientset) {
	mClient, clientset := getClients(kubeconfig)
	VersionCheck(clientset)
	return mClient, clientset
}

func getClients(kubeconfig string) (*v1.MarkhorV1Client, *kubernetes.Clientset) {
	k8sConfig := GetK8sConfig(kubeconfig)

	mClient := GetK8sClient(k8sConfig)
	clientset, err := kubernetes.NewForConfig(k8sConfig)

	if err != nil {
		slog.Error("Could not get a client to communicate with the k8s cluster")
		panic(err.Error())
	}
	return mClient, clientset
}
