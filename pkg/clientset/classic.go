package clientset

import (
	"fmt"
	"log/slog"

	clientV1 "github.com/civts/markhor/pkg/clientset/v1"

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
		slog.Info("Using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		slog.Info(fmt.Sprint("Reading k8s configuration from the specified file: ", kubeconfig))
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		slog.Error("Could not find a valid configuration to communicate with the k8s cluster")
		panic(err)
	}
	return config
}
