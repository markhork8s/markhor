package helpers

import (
	"encoding/json"
	"testing"

	"github.com/markhork8s/markhor/pkg"
)

func TestGetLabel(t *testing.T) {
	decryptedData := make(map[string]interface{})
	labelKey := pkg.MSPARAMS_MANAGED_LABEL_KEY
	manifestKey := pkg.MARKHORPARAMS_MANIFEST_KEY

	t.Run("Label present and valid", func(t *testing.T) {
		expectedLabel := "sample label"
		paramsObj := make(map[string]interface{})
		paramsObj[labelKey] = expectedLabel
		decryptedData[manifestKey] = paramsObj
		m, err := json.Marshal(decryptedData)
		if err != nil {
			t.Fatal("Should have been able to serialize to JSON")
		}
		var um = make(map[string]interface{})
		err = json.Unmarshal([]byte(m), &um)
		if err != nil {
			t.Fatal("Should have been able to serialize to JSON")
		}

		label, ok := GetLabel(um)

		if !ok {
			t.Fatalf("Expected label to be present, but it was not")
		}
		if label != expectedLabel {
			t.Fatalf("Expected label to be %s, got %s", expectedLabel, label)
		}
	})

	t.Run("Manifest key not present", func(t *testing.T) {
		_, ok := GetLabel(make(map[string]interface{}))

		if ok {
			t.Error("Expected label to not be present as manifest key is missing")
		}
	})

	t.Run("Params object not an ordered map", func(t *testing.T) {
		decryptedData[manifestKey] = "invalid type"

		_, ok := GetLabel(decryptedData)

		if ok {
			t.Error("Expected label to not be present as params object is not an ordered map")
		}
	})

	t.Run("Label key not present in params object", func(t *testing.T) {
		paramsObj := make(map[string]interface{})
		decryptedData[manifestKey] = paramsObj

		_, ok := GetLabel(decryptedData)

		if ok {
			t.Error("Expected label to not be present as label key is missing in params object")
		}
	})

	t.Run("Label value not a string", func(t *testing.T) {
		paramsObj := make(map[string]interface{})
		paramsObj[labelKey] = 123 // setting an integer instead of a string
		decryptedData[manifestKey] = paramsObj

		_, ok := GetLabel(decryptedData)

		if ok {
			t.Error("Expected label to not be present as label value is not a string")
		}
	})
}
