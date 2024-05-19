package helpers

import (
	"github.com/markhork8s/markhor/pkg"
)

func GetLabel(decryptedData map[string]interface{}) (string, bool) {
	params, present := decryptedData[pkg.MARKHORPARAMS_MANIFEST_KEY]
	if !present {
		return "", false
	}
	paramsObj, ok := params.(map[string]interface{})
	if !ok {
		return "", false
	}
	label, present := paramsObj[pkg.MSPARAMS_MANAGED_LABEL_KEY]
	if !present {
		return "", false
	}
	labelStr, ok := label.(string)
	if !ok {
		return "", false
	}
	return labelStr, true
}
