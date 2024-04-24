package helpers

import (
	"github.com/civts/markhor/pkg"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func GetAnnotation(decryptedData *orderedmap.OrderedMap[string, interface{}]) (string, bool) {
	params, present := decryptedData.Get(pkg.MARKHORPARAMS_MANIFEST_KEY)
	if !present {
		return "", false
	}
	paramsObj, ok := params.(orderedmap.OrderedMap[string, interface{}])
	if !ok {
		return "", false
	}
	annotation, present := paramsObj.Get(pkg.MSPARAMS_MANAGED_ANNOTATION_KEY)
	if !present {
		return "", false
	}
	annotationStr, ok := annotation.(string)
	if !ok {
		return "", false
	}
	return annotationStr, true
}
