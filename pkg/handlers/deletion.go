package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/civts/markhor/pkg/decrypt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HandleDeletion(args HandlerAttrs) {

	name := args.MarkhorSecret.ObjectMeta.Name
	namespace := args.MarkhorSecret.ObjectMeta.Namespace
	secretName := fmt.Sprintf("%s/%s", namespace, name)
	_, err := decrypt.DecryptMarkhorSecretEvent(args.MarkhorSecret, args.EventId)
	if err != nil {
		slog.Error(fmt.Sprint("Could not decrypt MarkhorSecret ", secretName), args.EventId)
		return
	}

	{ // Delete the secret
		// TODO: check if it was managed by markhor
		err = args.Clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})

		if err != nil {
			slog.Error(fmt.Sprint("Error deleting the secret: ", err), args.EventId)
			//Apply failed with 1 conflict: conflict with>another fieldmanager has the secret
		} else {
			slog.Info(fmt.Sprint("Secret deleted correctly ", secretName), args.EventId)
		}
	}
}
