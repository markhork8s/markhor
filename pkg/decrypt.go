package pkg

import (
	"encoding/json"
	"log"

	v1 "github.com/civts/markhor/pkg/api/types/v1"

	"github.com/getsops/sops/v3/decrypt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// Given an encrypted MarkhorSecret, this function attempts to decrypt it using SOPS.
func DecryptMarkhorSecret(markhorSecret *v1.MarkhorSecret) (*orderedmap.OrderedMap[string, interface{}], error) {
	jsonConfigStr := markhorSecret.ObjectMeta.Annotations["kubectl.kubernetes.io/last-applied-configuration"]

	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(jsonConfigStr), &jsonObj)
	if err != nil {
		log.Println("Error unmarshalling encrypted JSON:", err)
		return nil, err
	}

	sortedJson := sortJson(jsonObj)
	encData, err := json.Marshal(sortedJson)
	if err != nil {
		log.Println("Error marshalling sorted encrypted JSON:", err)
		return nil, err
	}

	decryptedDataBytes, err := decrypt.Data(encData, "json")
	if err != nil {
		log.Println("Error decrypting JSON:", err)
		return nil, err
	}

	decryptedData := orderedmap.New[string, interface{}]()
	err = json.Unmarshal(decryptedDataBytes, &decryptedData)
	if err != nil {
		log.Println("Error parsing decrypted JSON:", err)
		return nil, err
	}

	return decryptedData, nil
}
