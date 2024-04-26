package decrypt

import (
	"encoding/json"
	"fmt"
	"log/slog"

	v1 "github.com/civts/markhor/pkg/api/types/v1"

	"github.com/getsops/sops/v3/decrypt"
)

// Given an encrypted MarkhorSecret, this function attempts to decrypt it using SOPS.
func DecryptMarkhorSecret(markhorSecret *v1.MarkhorSecret, eid slog.Attr) (map[string]interface{}, error) {
	jsonConfigStr := markhorSecret.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]

	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(jsonConfigStr), &jsonObj)
	if err != nil {
		slog.Error(fmt.Sprint("Error unmarshalling encrypted JSON: ", err), eid)
		return nil, err
	}

	sortedJson := sortJson(jsonObj, eid)
	encData, err := json.Marshal(sortedJson)
	if err != nil {
		slog.Error(fmt.Sprint("Error marshalling sorted encrypted JSON: ", err), eid)
		return nil, err
	}

	decryptedDataBytes, err := decrypt.Data(encData, "json")
	if err != nil {
		slog.Error(fmt.Sprint("Error decrypting JSON: ", err), eid)
		return nil, err
	}

	decryptedData := make(map[string]interface{})
	err = json.Unmarshal(decryptedDataBytes, &decryptedData)
	if err != nil {
		slog.Error(fmt.Sprint("Error parsing decrypted JSON: ", err), eid)
		return nil, err
	}

	return decryptedData, nil
}
