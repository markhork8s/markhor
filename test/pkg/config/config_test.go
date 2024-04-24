package config_test

import (
	"os"
	"testing"

	c "github.com/civts/markhor/pkg/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func reset() {
	pflag.CommandLine = &pflag.FlagSet{}
	os.Args = []string{"./markhor"}
	viper.Reset()
}

// When no configuration file is provided, no error should occour in
// the configuration parsing phase
func TestValidDefaultConfig(t *testing.T) {
	reset()
	ensureDefaultConfigDoesNotExist(t)
	_ = c.ParseConfig()
	//If we get here the test passed since a valid configuration was created
}

// When no configuration file is provided, no error should occour in
// the configuration parsing phase AND the resulting configuration should use
// exactly the deafult values
func TestParseConfigDefaultValuesOnMissingDefaultFile(t *testing.T) {
	reset()
	ensureDefaultConfigDoesNotExist(t)

	// Parse the config without the --config flag specified
	config := c.ParseConfig()

	// Add assertions to validate the returned Config struct
	if config.Kubernetes.KubeconfigPath != c.DefaultKubeconfigPath {
		t.Errorf("Expected KubeconfigPath to be %s, but got %s", c.DefaultKubeconfigPath, config.Kubernetes.KubeconfigPath)
	}
	if config.Kubernetes.ClusterTimeoutSeconds != c.DefaultClusterTimeoutSeconds {
		t.Errorf("Expected ClusterTimeoutSeconds to be %d, but got %d", c.DefaultClusterTimeoutSeconds, config.Kubernetes.ClusterTimeoutSeconds)
	}
	if config.Sops.KeysPath != c.DefaultSopsKeysPath {
		t.Errorf("Expected SopsKeysPath to be %s, but got %s", c.DefaultSopsKeysPath, config.Sops.KeysPath)
	}
	if config.Healthcheck.Port != c.DefaultHealthcheckPort {
		t.Errorf("Expected HealthcheckPort to be %d, but got %d", c.DefaultHealthcheckPort, config.Healthcheck.Port)
	}
	if config.Healthcheck.Enabled != c.DefaultHealthcheckEnabled {
		t.Errorf("Expected HealthcheckEnabled to be %t, but got %t", c.DefaultHealthcheckEnabled, config.Healthcheck.Enabled)
	}
	if config.Logging.Level != c.DefaultLoggingLevel {
		t.Errorf("Expected LoggingLevel to be %s, but got %s", c.DefaultLoggingLevel, config.Logging.Level)
	}
	if config.Logging.Style != c.DefaultLoggingStyle {
		t.Errorf("Expected LoggingStyle to be %s, but got %s", c.DefaultLoggingStyle, config.Logging.Style)
	}
	if config.Logging.LogToStdout != c.DefaultLoggingLogToStdout {
		t.Errorf("Expected LoggingLogToStdout to be %t, but got %t", c.DefaultLoggingLogToStdout, config.Logging.LogToStdout)
	}
	if len(config.Logging.AdditionalLogFiles) != len(c.DefaultLoggingAdditionalLogFiles) {
		t.Errorf("Expected LoggingAdditionalLogFiles to be %s, but got %s", c.DefaultLoggingAdditionalLogFiles, config.Logging.AdditionalLogFiles)
	} else {
		for i, a := range config.Logging.AdditionalLogFiles {
			b := c.DefaultLoggingAdditionalLogFiles[i]
			if b != a {
				t.Errorf("Element number %d of config.Logging.AdditionalLogFiles should be '%s' but is '%s'", i, a, b)
			}
		}
	}
	if config.Behavior.Fieldmanager.Name != c.DefaultBehaviorFieldManagerName {
		t.Errorf("Expected BehaviorFieldManagerName to be %s, but got %s", c.DefaultBehaviorFieldManagerName, config.Behavior.Fieldmanager.Name)
	}
	if len(config.Behavior.Namespaces) != len(c.DefaultBehaviorNamespaces) {
		t.Errorf("Expected BehaviorNamespaces to be %s, but got %s", c.DefaultBehaviorNamespaces, config.Behavior.Namespaces)
	} else {
		for i, a := range config.Behavior.Namespaces {
			b := c.DefaultBehaviorNamespaces[i]
			if b != a {
				t.Errorf("Element number %d of config.Behavior.Namespaces should be '%s' but is '%s'", i, a, b)
			}
		}
	}
	if len(config.Behavior.ExcludedNamespaces) != len(c.DefaultBehaviorExcludedNamespaces) {
		t.Errorf("Expected BehaviorExcludedNamespaces to be %s, but got %s", c.DefaultBehaviorExcludedNamespaces, config.Behavior.ExcludedNamespaces)
	} else {
		for i, a := range config.Behavior.ExcludedNamespaces {
			b := c.DefaultBehaviorExcludedNamespaces[i]
			if b != a {
				t.Errorf("Element number %d of config.Behavior.ExcludedNamespaces should be '%s' but is '%s'", i, a, b)
			}
		}
	}
	if config.Behavior.Fieldmanager.ForceUpdates != c.DefaultBehaviorFieldManagerForceUpdates {
		t.Errorf("Expected BehaviorFieldManagerForceUpdates to be %t, but got %t", c.DefaultBehaviorFieldManagerForceUpdates, config.Behavior.Fieldmanager.ForceUpdates)
	}
	if config.Behavior.PruneDanglingMarkhorSecrets != c.DefaultBehaviorPruneDanglingMarkhorSecrets {
		t.Errorf("Expected BehaviorPruneDanglingMarkhorSecrets to be %t, but got %t", c.DefaultBehaviorPruneDanglingMarkhorSecrets, config.Behavior.PruneDanglingMarkhorSecrets)
	}
	if config.MarkorSecrets.HierarchySeparator.Default != c.DefaultMarkorSecretsHierarchySeparatorDefault {
		t.Errorf("Expected MarkorSecretsHierarchySeparatorDefault to be %s, but got %s", c.DefaultMarkorSecretsHierarchySeparatorDefault, config.MarkorSecrets.HierarchySeparator.Default)
	}
	if config.MarkorSecrets.HierarchySeparator.AllowOverride != c.DefaultMarkorSecretsHierarchySeparatorAllowOverride {
		t.Errorf("Expected MarkorSecretsHierarchySeparatorAllowOverride to be %t, but got %t", c.DefaultMarkorSecretsHierarchySeparatorAllowOverride, config.MarkorSecrets.HierarchySeparator.AllowOverride)
	}
	if config.MarkorSecrets.HierarchySeparator.WarnOnOverride != c.DefaultMarkorSecretsHierarchySeparatorWarnOnOverride {
		t.Errorf("Expected MarkorSecretsHierarchySeparatorWarnOnOverride to be %t, but got %t", c.DefaultMarkorSecretsHierarchySeparatorWarnOnOverride, config.MarkorSecrets.HierarchySeparator.WarnOnOverride)
	}
	if config.MarkorSecrets.ManagedAnnotation.Default != c.DefaultMarkorSecretsManagedAnnotationDefault {
		t.Errorf("Expected MarkorSecretsManagedAnnotationDefault to be %s, but got %s", c.DefaultMarkorSecretsManagedAnnotationDefault, config.MarkorSecrets.ManagedAnnotation.Default)
	}
	if config.MarkorSecrets.ManagedAnnotation.AllowOverride != c.DefaultMarkorSecretsManagedAnnotationAllowOverride {
		t.Errorf("Expected MarkorSecretsManagedAnnotationAllowOverride to be %t, but got %t", c.DefaultMarkorSecretsManagedAnnotationAllowOverride, config.MarkorSecrets.ManagedAnnotation.AllowOverride)
	}
	if config.MarkorSecrets.ManagedAnnotation.WarnOnOverride != c.DefaultMarkorSecretsManagedAnnotationWarnOnOverride {
		t.Errorf("Expected MarkorSecretsManagedAnnotationWarnOnOverride to be %t, but got %t", c.DefaultMarkorSecretsManagedAnnotationWarnOnOverride, config.MarkorSecrets.ManagedAnnotation.WarnOnOverride)
	}
}
