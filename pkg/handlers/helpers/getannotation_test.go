package helpers

import (
	"encoding/json"
	"testing"

	"github.com/civts/markhor/pkg"
)

func TestGetAnnotation(t *testing.T) {
	decryptedData := make(map[string]interface{})
	annotationKey := pkg.MSPARAMS_MANAGED_ANNOTATION_KEY
	manifestKey := pkg.MARKHORPARAMS_MANIFEST_KEY

	t.Run("Annotation present and valid", func(t *testing.T) {
		expectedAnnotation := "sample annotation"
		paramsObj := make(map[string]interface{})
		paramsObj[annotationKey] = expectedAnnotation
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

		annotation, ok := GetAnnotation(um)

		if !ok {
			t.Fatalf("Expected annotation to be present, but it was not")
		}
		if annotation != expectedAnnotation {
			t.Fatalf("Expected annotation to be %s, got %s", expectedAnnotation, annotation)
		}
	})

	t.Run("Manifest key not present", func(t *testing.T) {
		_, ok := GetAnnotation(make(map[string]interface{}))

		if ok {
			t.Error("Expected annotation to not be present as manifest key is missing")
		}
	})

	t.Run("Params object not an ordered map", func(t *testing.T) {
		decryptedData[manifestKey] = "invalid type"

		_, ok := GetAnnotation(decryptedData)

		if ok {
			t.Error("Expected annotation to not be present as params object is not an ordered map")
		}
	})

	t.Run("Annotation key not present in params object", func(t *testing.T) {
		paramsObj := make(map[string]interface{})
		decryptedData[manifestKey] = paramsObj

		_, ok := GetAnnotation(decryptedData)

		if ok {
			t.Error("Expected annotation to not be present as annotation key is missing in params object")
		}
	})

	t.Run("Annotation value not a string", func(t *testing.T) {
		paramsObj := make(map[string]interface{})
		paramsObj[annotationKey] = 123 // setting an integer instead of a string
		decryptedData[manifestKey] = paramsObj

		_, ok := GetAnnotation(decryptedData)

		if ok {
			t.Error("Expected annotation to not be present as annotation value is not a string")
		}
	})
}
