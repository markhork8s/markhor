package main

import (
	"context"
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
	v1a "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const MANAGED_BY = "Markhor"

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
		namespace := markhorSecret.ObjectMeta.Namespace
		secretName := fmt.Sprintf("%s/%s", namespace, markhorSecret.ObjectMeta.Name)
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

			{ // Add managed-by annotation
				annotation, present := getAnnotation(decryptedData)
				metadata, ok := decryptedData.Get("metadata")
				if !ok {
					fmt.Println("Missing metadata in ", secretName)
					continue
				}
				metadataObj, ok := metadata.(map[string]interface{})
				if !ok {
					fmt.Println("Missing metadata in ", secretName)
					continue
				}
				annotations, ok := metadataObj["annotations"]
				if !ok {
					fmt.Println("Missing annotations in ", secretName)
					annotations = make(map[string]interface{})
				}
				annotationsObj, ok := annotations.(map[string]interface{})
				if !ok {
					fmt.Println("Missing annotations in ", secretName)
					annotationsObj = make(map[string]interface{})
				}
				if present {
					annotationsObj[annotation] = MANAGED_BY
				} else {
					annotationsObj["markhor.example.com/managed-by"] = MANAGED_BY
				}
				metadataObj["annotations"] = annotationsObj
				decryptedData.Set("metadata", metadataObj)
			}

			{ //Remove extra fields
				decryptedData.Delete("markhorParams")
				decryptedData.Delete("sops")
				decryptedData.Set("apiVersion", "v1")
				decryptedData.Set("kind", "Secret")
			}

			{ //Create new secret
				// secret := &corev1.Secret{}
				secret := &v1a.SecretApplyConfiguration{}

				bytes, err := json.Marshal(decryptedData)
				if err != nil {
					fmt.Println("can't convert decrypted final to JSON:", err)
					panic(err)
				}
				if err := json.Unmarshal(bytes, secret); err != nil {
					fmt.Println("can't make secret from final JSON:", err)
					panic(err)
				}

				clientset, err := kubernetes.NewForConfig(getK8sConfig())
				if err != nil {
					panic(err.Error())
				}

				// clientset.CoreV1().Secrets("").Watch(context.TODO(), metav1.ListOptions{})
				// Apply the secret
				fieldManager := "markhor"

				_, err = clientset.CoreV1().Secrets(namespace).Apply(context.TODO(), secret, metav1.ApplyOptions{
					FieldManager: fieldManager,
				})
				if err != nil {
					fmt.Println("error creating the secret:", err)
					//Apply failed with 1 conflict: conflict with>another fieldmanager has the secret
				} else {
					fmt.Println("new secret created correctly", secretName)
				}
			}
		case watch.Modified:
			fmt.Println("A MarkhorSecret was updated:", secretName)

		case watch.Deleted:
			fmt.Println("A MarkhorSecret was deleted:", secretName)
		}
	}
	fmt.Println("Finished watching the events in the cluster. Most probably, the channel was closed")
}

func getAnnotation(decryptedData *orderedmap.OrderedMap[string, interface{}]) (string, bool) {
	params, present := decryptedData.Get("markhorParams")
	if !present {
		return "", false
	}
	paramsObj, ok := params.(orderedmap.OrderedMap[string, interface{}])
	if !ok {
		return "", false
	}
	annotation, present := paramsObj.Get("managedAnnotation")
	if !present {
		return "", false
	}
	annotationStr, ok := annotation.(string)
	if !ok {
		return "", false
	}
	return annotationStr, true
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
