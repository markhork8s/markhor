package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const DECRYPTED_SUFFIX string = "-decr"

func main() {
	fmt.Println("Welcome to sops_k8s")
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	watchSecrets(clientset)
}

func watchSecrets(clientset *kubernetes.Clientset) {
	watcher, err := clientset.CoreV1().Secrets("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to watch the events in the cluster to see when secrets are created")
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added, watch.Modified:
			fmt.Println("A secret was added or updated")
			handleSecretUpdated(event, clientset)
		case watch.Deleted:
			fmt.Println("A secret was deleted")
			handleSecretDeletion(clientset, event.Object.(*v1.Secret))
		}
	}
	fmt.Println("Finished watching the events in the cluster. Most probably, the channel was closed")
}
