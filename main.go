package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	v1 "sops_k8s/api/types/v1"
	clientV1 "sops_k8s/clientset/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	v1.AddToScheme(scheme.Scheme)

	k8sclient := getK8sClient()

	sopsSecrets, err := k8sclient.SopsSecrets("irrelevant").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("sopsSecrets found: %+v\n", sopsSecrets)

	store := WatchResources(k8sclient)

	for {
		sopsSecretsFromStore := store.List()
		fmt.Printf("sopsSecret in store: %d\n", len(sopsSecretsFromStore))

		time.Sleep(2 * time.Second)
	}
}

func getK8sClient() *clientV1.SopsV1Client {
	config := getK8sConfig()

	clientSet, err := clientV1.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientSet
}

func getK8sConfig() *rest.Config {
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
