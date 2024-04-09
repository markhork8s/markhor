package clientset

import (
	"log"

	clientV1 "markhor/pkg/clientset/v1"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetK8sClient(config *rest.Config) *clientV1.MarkhorV1Client {

	clientSet, err := clientV1.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientSet
}

func GetK8sConfig(kubeconfig string) *rest.Config {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("Using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("Using configuration from the command flags")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		log.Printf("Could not find a valid configuration to communicate with the k8s cluster")
		panic(err)
	}
	return config
}
