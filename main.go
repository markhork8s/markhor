package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	clientV1 "sops_k8s/clientset/v1"

	"k8s.io/api/node/v1alpha1"
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
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from the command flags")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		log.Printf("Could not find a valid configuration to communicate with the k8s cluster")
		panic(err)
	}

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	log.Printf("3")
	sopsSecrets, err := clientSet.SopsSecrets("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("sopsSecrets found: %+v\n", sopsSecrets)

	store := WatchResources(clientSet)

	for {
		sopsSecretsFromStore := store.List()
		fmt.Printf("sopsSecret in store: %d\n", len(sopsSecretsFromStore))

		time.Sleep(2 * time.Second)
	}
}
