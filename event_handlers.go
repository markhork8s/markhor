package main

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func handleSecretUpdated(event watch.Event, clientset *kubernetes.Clientset) {
	secret, ok := event.Object.(*v1.Secret)
	if !ok {
		fmt.Println("Error decoding the secret from the event")
		return
	}

	if val, exists := secret.Annotations["sops_k8s/decryption-enabled"]; exists && val == "true" {
		fmt.Println("Decrypting the content of ", secret.Namespace, secret.Name)
		decryptedData, err := decryptSecretData(secret.Data)
		if err != nil {
			fmt.Println("Error decrypting secret data:", err)
			return
		}
		decryptedStringData, err := decryptSecretData(convertMapStringToMapByte(secret.StringData))
		if err != nil {
			fmt.Println("Error decrypting secret stringData:", err)
			return
		}

		v, ok := secret.Annotations["managed_by"]
		if !ok {
			fmt.Println("The secret ", secret.Name, " already had a value for managed_by of ", v, ". It will be replaced")
		}
		secret.Annotations["managed_by"] = "k8s_sops"
		secret.Data = decryptedData
		secret.StringData = convertMapByteToMapString(decryptedStringData)

		_, err = clientset.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
		if err != nil {
			fmt.Println("Error creating new secret:", err)
			return
		}

		fmt.Println("Decrypted secret created:", secret.Name)
	}
}

func handleSecretDeletion(clientset *kubernetes.Clientset, secret *v1.Secret) {
	if val, exists := secret.Annotations["decrypt.example.com/decryption-enabled"]; exists && val == "true" {
		decryptedSecretName := secret.Name + DECRYPTED_SUFFIX
		err := clientset.CoreV1().Secrets(secret.Namespace).Delete(context.TODO(), decryptedSecretName, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Error deleting decrypted secret %s: %v\n", decryptedSecretName, err)
		} else {
			fmt.Printf("Decrypted secret %s deleted\n", decryptedSecretName)
		}
	}
}
