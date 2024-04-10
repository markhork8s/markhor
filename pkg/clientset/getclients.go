package clientset

import (
	v1 "github.com/civts/markhor/pkg/clientset/v1"

	"k8s.io/client-go/kubernetes"
)

func GetK8sClients(kubeconfig string) (*v1.MarkhorV1Client, *kubernetes.Clientset) {
	k8sConfig := GetK8sConfig(kubeconfig)

	mClient := GetK8sClient(k8sConfig)
	clientset, err := kubernetes.NewForConfig(k8sConfig)

	if err != nil {
		panic(err.Error())
	}
	return mClient, clientset
}
