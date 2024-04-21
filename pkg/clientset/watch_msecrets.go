package clientset

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/civts/markhor/pkg"
	apiV1 "github.com/civts/markhor/pkg/api/types/v1"
	v1 "github.com/civts/markhor/pkg/clientset/v1"
	"github.com/civts/markhor/pkg/handlers"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

var connectedSuccessfully = false

func WatchMarkhorSecrets(mClient *v1.MarkhorV1Client, k8sClient *kubernetes.Clientset) {
	go checkConnectTimeout(10)
	markhorSecrets, err := mClient.MarkhorSecrets("irrelevant").Watch(metav1.ListOptions{})
	if err != nil {
		e := err.Error()
		if strings.Contains(e, "the server could not find the requested resource") {
			log.Println("Kubernetes does not know what a MarkhorSecret is. Did you forget to install the CRD?")
		}
		panic(err)
	}
	connectedSuccessfully = true
	channel := markhorSecrets.ResultChan()
	log.Println("Started watching the events in the cluster")
	pkg.Healthy = true
	for event := range channel {
		markhorSecret, ok := event.Object.(*apiV1.MarkhorSecret)
		namespace := markhorSecret.ObjectMeta.Namespace
		secretName := fmt.Sprintf("%s/%s", namespace, markhorSecret.ObjectMeta.Name)
		if !ok {
			log.Println("Failed to cast the object to type MarkhorSecret")
			continue
		}
		switch event.Type {
		case watch.Added:
			log.Println("A MarkhorSecret was added:", secretName)
			handlers.HandleAddition(markhorSecret, secretName, namespace, k8sClient)
		case watch.Modified:
			log.Println("A MarkhorSecret was updated:", secretName)
			handlers.HandleAddition(markhorSecret, secretName, namespace, k8sClient)
		case watch.Deleted:
			log.Println("A MarkhorSecret was deleted:", secretName)
			handlers.HandleDeletion(markhorSecret, k8sClient)
		}
	}
	pkg.Healthy = false
	log.Println("Finished watching the events in the cluster. Most probably, the channel was closed")
}

func checkConnectTimeout(timeout int) {
	log.Print("Connecting to the k8s cluster")

	for i := 1; i <= timeout; i++ {
		time.Sleep(1 * time.Second)
		if connectedSuccessfully {
			return
		} else if i == 2 {
			log.Println("No response from the k8s cluster. Will retry until the timeout")
		}
	}

	log.Printf("Connecting to the k8s cluster timed out after %d seconds. Check the kubeconfig file and that the cluster is up.", timeout)

	os.Exit(1)
}
