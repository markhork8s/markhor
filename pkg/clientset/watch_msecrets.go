package clientset

import (
	"fmt"

	apiV1 "github.com/civts/markhor/pkg/api/types/v1"
	v1 "github.com/civts/markhor/pkg/clientset/v1"
	"github.com/civts/markhor/pkg/handlers"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func WatchMarkhorSecrets(mClient *v1.MarkhorV1Client, k8sClient *kubernetes.Clientset) {
	markhorSecrets, err := mClient.MarkhorSecrets("irrelevant").Watch(metav1.ListOptions{})
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
			handlers.HandleAddition(markhorSecret, secretName, namespace, k8sClient)
		case watch.Modified:
			fmt.Println("A MarkhorSecret was updated:", secretName)
			handlers.HandleAddition(markhorSecret, secretName, namespace, k8sClient)
		case watch.Deleted:
			fmt.Println("A MarkhorSecret was deleted:", secretName)
			handlers.HandleDeletion(markhorSecret, k8sClient)
		}
	}
	fmt.Println("Finished watching the events in the cluster. Most probably, the channel was closed")
}
