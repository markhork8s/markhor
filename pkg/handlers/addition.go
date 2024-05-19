package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/markhork8s/markhor/pkg"
	"github.com/markhork8s/markhor/pkg/decrypt"
	"github.com/markhork8s/markhor/pkg/handlers/helpers"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
)

const MANAGED_BY = "Markhor"

func HandleAddition(attrs HandlerAttrs) bool {

	namespace := attrs.MarkhorSecret.ObjectMeta.Namespace
	secretName := fmt.Sprintf("%s/%s", namespace, attrs.MarkhorSecret.ObjectMeta.Name)
	decryptedData, err := decrypt.DecryptMarkhorSecretEvent(attrs.MarkhorSecret, attrs.Config.MarkorSecrets, attrs.EventId)
	if err != nil {
		slog.Error(fmt.Sprint("Could not decrypt MarkhorSecret ", secretName), attrs.EventId)
		return false
	}

	{ // Add managed-by label
		metadata, ok := decryptedData["metadata"]
		if !ok {
			slog.Error(fmt.Sprint("Missing metadata in ", secretName), attrs.EventId)
			return false
		}
		metadataObj, ok := metadata.(map[string]interface{})
		if !ok {
			slog.Error(fmt.Sprint("Missing metadata in ", secretName), attrs.EventId)
			return false
		}
		labels, ok := metadataObj["labels"]
		if !ok {
			slog.Debug(fmt.Sprint("No existing labels found in ", secretName, " will add managed-by markhor anyway"), attrs.EventId)
			labels = make(map[string]interface{})
		}
		labelsObj, ok := labels.(map[string]interface{})
		if !ok {
			slog.Error(fmt.Sprint("Labels in ", secretName, " do not appear to be a YAML object"), attrs.EventId)
			labelsObj = make(map[string]interface{})
		}

		label, present := helpers.GetLabel(decryptedData)
		useCustomLabel := false
		if present {
			if attrs.Config.MarkorSecrets.ManagedLabel.AllowOverride {
				useCustomLabel = true
				msg := fmt.Sprint("Overriding managed-by label for ", secretName)
				if attrs.Config.MarkorSecrets.ManagedLabel.WarnOnOverride {
					slog.Warn(msg)
				} else {
					slog.Debug(msg)
				}
			}
		}
		if useCustomLabel {
			labelsObj[label] = MANAGED_BY
		} else {
			labelsObj[attrs.Config.MarkorSecrets.ManagedLabel.Default] = MANAGED_BY
		}

		metadataObj["labels"] = labelsObj
		decryptedData["metadata"] = metadataObj
	}

	{ //Remove extra fields
		decryptedData[pkg.MARKHORPARAMS_MANIFEST_KEY] = nil
		decryptedData["sops"] = nil
		decryptedData["apiVersion"] = "v1"
		decryptedData["kind"] = "Secret"
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
			slog.Error(fmt.Sprint("Error creating/updating the secret: ", err), attrs.EventId)
			//Apply failed with 1 conflict: conflict with>another fieldmanager has the secret
			return false
		} else {
			return true
		}
	}
}
