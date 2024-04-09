package main

import (
	"flag"
	"fmt"

	apiV1 "markhor/pkg/api/types/v1"
	cs "markhor/pkg/clientset"
	"markhor/pkg/handlers"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	apiV1.AddToScheme(scheme.Scheme)

	k8sConfig := cs.GetK8sConfig(kubeconfig)

	k8sclient := cs.GetK8sClient(k8sConfig)
	clientset, err := kubernetes.NewForConfig(k8sConfig)

	if err != nil {
		panic(err.Error())
	}

	markhorSecrets, err := k8sclient.MarkhorSecrets("irrelevant").Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to watch the events in the cluster to see when markhorSecrets are created")
	for event := range markhorSecrets.ResultChan() {
		markhorSecret, ok := event.Object.(*apiV1.MarkhorSecret)
		namespace := markhorSecret.ObjectMeta.Namespace
		secretName := fmt.Sprintf("%s/%s", namespace, markhorSecret.ObjectMeta.Name)
		if !ok {
			fmt.Println("Failed to cast the object to type MarkhorSecret")
			continue
		}
		switch event.Type {
		case watch.Added:
			fmt.Println("A MarkhorSecret was added:", secretName)
			handlers.HandleAddition(markhorSecret, secretName, namespace, clientset)
		case watch.Modified:
			fmt.Println("A MarkhorSecret was updated:", secretName)
			handlers.HandleAddition(markhorSecret, secretName, namespace, clientset)
		case watch.Deleted:
			fmt.Println("A MarkhorSecret was deleted:", secretName)
			handlers.HandleDeletion(markhorSecret, clientset)
		}
	}
	fmt.Println("Finished watching the events in the cluster. Most probably, the channel was closed")
}
