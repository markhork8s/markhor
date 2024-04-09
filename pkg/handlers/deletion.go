package handlers

import (
	"context"
	"fmt"

	"markhor/pkg"
	v1 "markhor/pkg/api/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func HandleDeletion(markhorSecret *v1.MarkhorSecret, clientset *kubernetes.Clientset) {

	name := markhorSecret.ObjectMeta.Name
	namespace := markhorSecret.ObjectMeta.Namespace
	secretName := fmt.Sprintf("%s/%s", namespace, name)
	_, err := pkg.DecryptMarkhorSecret(markhorSecret)
	if err != nil {
		fmt.Println("Error: something went wrong decrypting ", secretName)
		return
	}

	{ // Delete the secret
		// TODO: check if it was managed by markhor
		err = clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

		if err != nil {
			fmt.Println("error deleting the secret:", err)
			//Apply failed with 1 conflict: conflict with>another fieldmanager has the secret
		} else {
			fmt.Println("secret deleted correctly", secretName)
		}
	}
}
