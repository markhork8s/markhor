package pkg

import orderedmap "github.com/wk8/go-ordered-map/v2"

func GetAnnotation(decryptedData *orderedmap.OrderedMap[string, interface{}]) (string, bool) {
	params, present := decryptedData.Get("markhorParams")
	if !present {
		return "", false
	}
	paramsObj, ok := params.(orderedmap.OrderedMap[string, interface{}])
	if !ok {
		return "", false
	}
	annotation, present := paramsObj.Get("managedAnnotation")
	if !present {
		return "", false
	}
	annotationStr, ok := annotation.(string)
	if !ok {
		return "", false
	}
	return annotationStr, true
}
