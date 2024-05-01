package decrypt

import (
	"fmt"
	"log/slog"
	"strings"
)

func CheckSopsVersion(json map[string]interface{}) (string, bool) {
	result, ok := getVersion(json)
	if !ok {
		return result, false
	}
	return IsSupportedSopsVersion(result)
}

func getVersion(json map[string]interface{}) (string, bool) {
	sops := json["sops"]
	if sops == nil {
		return "Could not find the sops field in the provided JSON", false
	}
	sopsMap, ok := sops.(map[string]interface{})
	if !ok {
		slog.Debug(fmt.Sprint("The sops field in the provided JSON is not an object: ", sops))
		return "The sops field in the provided JSON is not an object", false
	}
	version := sopsMap["version"]
	versionStr, ok := version.(string)
	if !ok {
		slog.Debug(fmt.Sprint("The SOPS version is not a string: ", sops))
		return "The SOPS version is not a string", false
	}
	return versionStr, true
}

func IsSupportedSopsVersion(version string) (string, bool) {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 1 {
		return fmt.Sprint("The SOPS version is in an unexpected format: ", version), false
	}
	major := parts[0]
	if major != "3" {
		return fmt.Sprint("SOPS version 3 is supported, but we got ", version), false
	}
	if len(parts) < 2 {
		return fmt.Sprint("The SOPS version is in an unexpected format (missing the minor version): ", version), false
	}
	minor := parts[1]
	if minor != "8" {
		return fmt.Sprint("SOPS 3.8 is supported, but we got ", version), false
	}
	return "", true
}
