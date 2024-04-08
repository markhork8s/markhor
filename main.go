package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	v1 "sops_k8s/api/types/v1"
	clientV1 "sops_k8s/clientset/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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

	sopsSecrets, err := k8sclient.SopsSecrets("irrelevant").Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to watch the events in the cluster to see when sopsSecrets are created")
	for event := range sopsSecrets.ResultChan() {
		sopsSecret, ok := event.Object.(*v1.SopsSecret)
		if !ok {
			fmt.Println("Failed to cast the object to type SopsSecret")
			continue
		}
		switch event.Type {
		case watch.Added:
			jsonConfigStr := sopsSecret.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]
			fmt.Println("A sopsSecret was added ", sopsSecret.ObjectMeta.Namespace, "/", sopsSecret.ObjectMeta.Name)

			var jsonObj map[string]interface{}
			err := json.Unmarshal([]byte(jsonConfigStr), &jsonObj)
			if err != nil {
				fmt.Println("Error unmarshalling JSON:", err)
				return
			}

			sortedJson := sortJson(jsonObj)
			fmt.Println(sortedJson)

		case watch.Modified:
			fmt.Println("A sopsSecret was updated")
		case watch.Deleted:
			fmt.Println("A sopsSecret was deleted")
		}
	}
	fmt.Println("Finished watching the events in the cluster. Most probably, the channel was closed")
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
