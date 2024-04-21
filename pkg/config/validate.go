package config

import (
	"fmt"
	"log"
)

// This function sanitizes the Markhor configuration before it is used
func ValidateConfig(c Config) {
	s := validate(c)
	if s != "" {
		log.Fatalf("Invalid configuration: %s", s)
	}
}

func validate(c Config) string {
	if c.Kubernetes.ClusterTimeoutSeconds <= 0 {
		return "Kubernetes.ClusterTimeoutSeconds must be greater than zero"
	}

	if len(c.Sops.KeysPath) == 0 {
		return "The SOPS key path can't have a length of zero"
	}

	if c.Healthcheck.Port < 1 || c.Healthcheck.Port > 65534 {
		return "The port number for the healthcheck endpoint shall be in the range [1, 65534]"
	}

	for i, f := range c.Logging.AdditionalLogFiles {
		if len(f) == 0 {
			return fmt.Sprintf("The additional log file path number %d has an invalid length of zero", i)
		}
	}
	if hasDuplicates(c.Logging.AdditionalLogFiles) {
		return "The additional log files array shall not contain duplicates"
	}

	if len(c.Behavior.Fieldmanager.Name) == 0 {
		return "The field manager name can't be an empty string"
	}

	for _, n := range c.Behavior.Namespaces {
		if len(n) == 0 {
			return "Cannot have empty namespaces in Behavior.Namespaces"
		}
	}
	if hasDuplicates(c.Behavior.Namespaces) {
		return "The namespaces array shall not contain duplicates"
	}

	for _, n := range c.Behavior.ExcludedNamespaces {
		if len(n) == 0 {
			return "Cannot have empty namespaces in Behavior.ExcludedNamespaces"
		}
	}
	if hasDuplicates(c.Behavior.ExcludedNamespaces) {
		return "The excluded namespaces array shall not contain duplicates"
	}

	if len(c.MarkorSecrets.HierarchySeparator.Default) == 0 {
		return "The default hierarchy separator can't be empty"
	}

	if len(c.MarkorSecrets.ManagedAnnotation.Default) == 0 {
		return "The default managed annotation can't be empty"
	}

	return ""
}

func hasDuplicates(slice []string) bool {
	seen := make(map[string]struct{})

	for _, value := range slice {
		if _, ok := seen[value]; ok {
			return true
		}
		seen[value] = struct{}{}
	}

	return false
}
