package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	v1 "markhor/api/types/v1"
	clientV1 "markhor/clientset/v1"

	"github.com/getsops/sops/v3/decrypt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
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

	markhorSecrets, err := k8sclient.MarkhorSecrets("irrelevant").Watch(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting to watch the events in the cluster to see when markhorSecrets are created")
	for event := range markhorSecrets.ResultChan() {
		markhorSecret, ok := event.Object.(*v1.MarkhorSecret)
		secretName := fmt.Sprintf("%s/%s", markhorSecret.ObjectMeta.Namespace, markhorSecret.ObjectMeta.Name)
		if !ok {
			fmt.Println("Failed to cast the object to type MarkhorSecret")
			continue
		}
		switch event.Type {
		case watch.Added:
			fmt.Println("A MarkhorSecret was added:", secretName)

			decryptedData, err := decryptMarkhorSecret(markhorSecret)
			if err != nil {
				fmt.Println("Error: something went wrong decrypting ", secretName)
				continue
			}

			b, err := json.Marshal(decryptedData)
			if err != nil {
				fmt.Println("Error decrypted to json ", secretName)
				continue
			}

			fmt.Println(string(b))

		case watch.Modified:
			fmt.Println("A MarkhorSecret was updated:", secretName)

		case watch.Deleted:
			fmt.Println("A MarkhorSecret was deleted:", secretName)
		}
	}
	fmt.Println("Finished watching the events in the cluster. Most probably, the channel was closed")
}

func decryptMarkhorSecret(markhorSecret *v1.MarkhorSecret) (*orderedmap.OrderedMap[string, interface{}], error) {
	jsonConfigStr := markhorSecret.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]

	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(jsonConfigStr), &jsonObj)
	if err != nil {
		fmt.Println("Error unmarshalling encrypted JSON:", err)
		return nil, err
	}

	sortedJson := sortJson(jsonObj)
	encData, err := json.Marshal(sortedJson)
	if err != nil {
		fmt.Println("Error marshalling sorted encrypted JSON:", err)
		return nil, err
	}

	decryptedDataBytes, err := decrypt.Data(encData, "json")
	if err != nil {
		fmt.Println("Error decrypting JSON:", err)
		return nil, err
	}

	decryptedData := orderedmap.New[string, interface{}]()
	err = json.Unmarshal(decryptedDataBytes, &decryptedData)
	if err != nil {
		fmt.Println("Error parsing decrypted JSON:", err)
		return nil, err
	}

	return decryptedData, nil
}

func getK8sClient() *clientV1.MarkhorV1Client {
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
