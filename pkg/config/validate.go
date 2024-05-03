package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

// This function sanitizes the Markhor configuration before it is used
func ValidateConfig(c Config) error {
	s := validate(c)
	if s != "" {
		return errors.New("Invalid configuration: " + s)
	}
	return nil
}

const maxPort = 65535

func validate(c Config) string {
	if c.Kubernetes.ClusterTimeoutSeconds <= 0 {
		return "Kubernetes.ClusterTimeoutSeconds must be greater than zero"
	}

	if c.Healthcheck.Port < 1 || c.Healthcheck.Port > maxPort {
		return fmt.Sprintf("The port number for the healthcheck endpoint shall be in the range [1, %d]", maxPort)
	}

	if c.AdmissionController.Port < 1 || c.AdmissionController.Port > maxPort {
		return fmt.Sprintf("The port number for the admission controller endpoint shall be in the range [1, %d]", maxPort)
	}

	if c.AdmissionController.Enabled && c.Healthcheck.Enabled {
		if c.AdmissionController.Port == c.Healthcheck.Port {
			return "The port number for the admission controller and healthcheck cannot be equal"
		}
	}

	if len(c.Tls.CertPath) == 0 {
		return "The TLS certificate file path can't have a length of zero"
	}

	if len(c.Tls.KeyPath) == 0 {
		return "The TLS key file path can't have a length of zero"
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

	if !contains(ValidTLSModes, c.Tls.Mode) {
		return "The TLS mode is invalid. Valid values are " + strings.Join(ValidTLSModes, ", ")
	}

	if !mapContainsKey(LoggerLevels, c.Logging.Level) {
		return "The log level is invalid. Valid values are " + strings.Join(ValidTLSModes, ", ")
	}

	return ""
}

func contains(slice []string, elem string) bool {
	for _, value := range slice {
		if value == elem {
			return true
		}
	}
	return false
}
func mapContainsKey(slice map[string]slog.Level, elem string) bool {
	for k := range slice {
		if k == elem {
			return true
		}
	}
	return false
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
