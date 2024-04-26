package helpers

import (
	"github.com/civts/markhor/pkg"
)

func GetAnnotation(decryptedData map[string]interface{}) (string, bool) {
	params, present := decryptedData[pkg.MARKHORPARAMS_MANIFEST_KEY]
	if !present {
		return "", false
	}
	paramsObj, ok := params.(map[string]interface{})
	if !ok {
		return "", false
	}
	annotation, present := paramsObj[pkg.MSPARAMS_MANAGED_ANNOTATION_KEY]
	if !present {
		return "", false
	}
	annotationStr, ok := annotation.(string)
	if !ok {
		return "", false
	}
	return annotationStr, true
}
