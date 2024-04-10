package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/civts/markhor/pkg"

	v1 "github.com/civts/markhor/pkg/api/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1a "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
)

const MANAGED_BY = "github.com/civts/markhor"

func HandleAddition(markhorSecret *v1.MarkhorSecret, secretName string, namespace string, clientset *kubernetes.Clientset) {

	decryptedData, err := pkg.DecryptMarkhorSecret(markhorSecret)
	if err != nil {
		fmt.Println("Error: something went wrong decrypting ", secretName)
		return
	}

	{ // Add managed-by annotation
		annotation, present := pkg.GetAnnotation(decryptedData)
		metadata, ok := decryptedData.Get("metadata")
		if !ok {
			fmt.Println("Missing metadata in ", secretName)
			return
		}
		metadataObj, ok := metadata.(map[string]interface{})
		if !ok {
			fmt.Println("Missing metadata in ", secretName)
			return
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

		// clientset.CoreV1().Secrets("").Watch(context.TODO(), metav1.ListOptions{})
		// Apply the secret
		fieldManager := "github.com/civts/markhor"

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
}
