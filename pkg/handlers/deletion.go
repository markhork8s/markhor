package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/civts/markhor/pkg"

	v1 "github.com/civts/markhor/pkg/api/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func HandleDeletion(markhorSecret *v1.MarkhorSecret, clientset *kubernetes.Clientset) {

	name := markhorSecret.ObjectMeta.Name
	namespace := markhorSecret.ObjectMeta.Namespace
	secretName := fmt.Sprintf("%s/%s", namespace, name)
	_, err := pkg.DecryptMarkhorSecret(markhorSecret)
	if err != nil {
		log.Println("Error: could not decrypt MarkhorSecret", secretName)
		return
	}

	{ // Delete the secret
		// TODO: check if it was managed by markhor
		err = clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

		if err != nil {
			log.Println("error deleting the secret:", err)
			//Apply failed with 1 conflict: conflict with>another fieldmanager has the secret
		} else {
			log.Println("secret deleted correctly", secretName)
		}
	}
}
