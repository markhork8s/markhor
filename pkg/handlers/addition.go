package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/civts/markhor/pkg"
	"github.com/civts/markhor/pkg/decrypt"
	"github.com/civts/markhor/pkg/handlers/helpers"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
)

const MANAGED_BY = "github.com/civts/markhor"

func HandleAddition(attrs HandlerAttrs) {

	namespace := attrs.MarkhorSecret.ObjectMeta.Namespace
	secretName := fmt.Sprintf("%s/%s", namespace, attrs.MarkhorSecret.ObjectMeta.Name)
	decryptedData, err := decrypt.DecryptMarkhorSecret(attrs.MarkhorSecret, attrs.EventId)
	if err != nil {
		slog.Error(fmt.Sprint("Could not decrypt MarkhorSecret ", secretName), attrs.EventId)
		return
	}

	{ // Add managed-by annotation
		annotation, present := helpers.GetAnnotation(decryptedData)
		metadata, ok := decryptedData.Get("metadata")
		if !ok {
			slog.Error(fmt.Sprint("Missing metadata in ", secretName), attrs.EventId)
			return
		}
		metadataObj, ok := metadata.(map[string]interface{})
		if !ok {
			slog.Error(fmt.Sprint("Missing metadata in ", secretName), attrs.EventId)
			return
		}
		annotations, ok := metadataObj["annotations"]
		if !ok {
			slog.Error(fmt.Sprint("Missing annotations in ", secretName), attrs.EventId)
			annotations = make(map[string]interface{})
		}
		annotationsObj, ok := annotations.(map[string]interface{})
		if !ok {
			slog.Error(fmt.Sprint("Missing annotations in ", secretName), attrs.EventId)
			annotationsObj = make(map[string]interface{})
		}
		if present {
			annotationsObj[annotation] = MANAGED_BY
		} else {
			annotationsObj[attrs.Config.MarkorSecrets.ManagedAnnotation.Default] = MANAGED_BY
		}
		metadataObj["annotations"] = annotationsObj
		decryptedData.Set("metadata", metadataObj)
	}

	{ //Remove extra fields
		decryptedData.Delete(pkg.MARKHORPARAMS_MANIFEST_KEY)
		decryptedData.Delete("sops")
		decryptedData.Set("apiVersion", "v1")
		decryptedData.Set("kind", "Secret")
	}

	{ //Create new secret
		secret := &v1.SecretApplyConfiguration{}

		bytes, err := json.Marshal(decryptedData)
		if err != nil {
			slog.Error(fmt.Sprint("Can't convert decrypted final to JSON: ", err), attrs.EventId)
			panic(err)
		}
		if err := json.Unmarshal(bytes, secret); err != nil {
			slog.Error(fmt.Sprint("Can't make secret from final JSON: ", err), attrs.EventId)
			panic(err)
		}

		_, err = attrs.Clientset.CoreV1().Secrets(namespace).Apply(context.TODO(), secret, metav1.ApplyOptions{
			FieldManager: attrs.Config.Behavior.Fieldmanager.Name,
			Force:        attrs.Config.Behavior.Fieldmanager.ForceUpdates,
		})
		if err != nil {
			slog.Error(fmt.Sprint("Error creating the secret: ", err), attrs.EventId)
			//Apply failed with 1 conflict: conflict with>another fieldmanager has the secret
		} else {
			slog.Info(fmt.Sprint("New secret created correctly: ", secretName), attrs.EventId)
		}
	}
}
