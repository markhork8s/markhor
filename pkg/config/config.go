package config

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/civts/markhor/pkg"
	"github.com/spf13/viper"
)

// Parses the program configuration (from CLI and file).
// Will terminate the execution if it fails since it would not be illogical to
// proceed w/o a valid config
func ParseConfig() (*Config, error) {
	err := parseCliArgs()
	if err != nil {
		return nil, err
	}
	configFilePath := viper.GetString("config")

	usingCustomConfigFile := configFilePath != pkg.DEFAULT_CONFIG_PATH
	if usingCustomConfigFile {
		defer slog.Info(fmt.Sprint("Reading Markhor config from user-defined path: ", configFilePath))
	} else {
		defer slog.Info(fmt.Sprint("Reading Markhor config from default path: ", pkg.DEFAULT_CONFIG_PATH))
	}

	viper.SetConfigFile(configFilePath)
	setDefaultConfigValues()

	err = viper.ReadInConfig()
	if err != nil && usingCustomConfigFile {
		return nil, errors.New("Error reading Markhor config file: " + err.Error())
	}

	// Define custom Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.New("Error parsing Markhor config data: " + err.Error())
	}

	err = ValidateConfig(config)
	if err != nil {
		return nil, err
	}

	err = SetupLogging(config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Defines the default Markhor config values
func setDefaultConfigValues() {
	viper.SetDefault("kubernetes.kubeconfigPath", DefaultKubeconfigPath)
	viper.SetDefault("kubernetes.clusterTimeoutSeconds", DefaultClusterTimeoutSeconds)
	viper.SetDefault("sops.keysPath", DefaultSopsKeysPath)
	viper.SetDefault("healthcheck.port", DefaultHealthcheckPort)
	viper.SetDefault("healthcheck.enabled", DefaultHealthcheckEnabled)
	viper.SetDefault("logging.level", DefaultLoggingLevel)
	viper.SetDefault("logging.style", DefaultLoggingStyle)
	viper.SetDefault("logging.logToStdout", DefaultLoggingLogToStdout)
	viper.SetDefault("logging.additionalLogFiles", DefaultLoggingAdditionalLogFiles)
	viper.SetDefault("behavior.fieldmanager.name", DefaultBehaviorFieldManagerName)
	viper.SetDefault("behavior.namespaces", DefaultBehaviorNamespaces)
	viper.SetDefault("behavior.excludedNamespaces", DefaultBehaviorExcludedNamespaces)
	viper.SetDefault("behavior.fieldmanager.forceUpdates", DefaultBehaviorFieldManagerForceUpdates)
	viper.SetDefault("behavior.pruneDanglingMarkhorSecrets", DefaultBehaviorPruneDanglingMarkhorSecrets)
	viper.SetDefault("markorSecrets.hierarchySeparator.default", DefaultMarkorSecretsHierarchySeparatorDefault)
	viper.SetDefault("markorSecrets.hierarchySeparator.allowOverride", DefaultMarkorSecretsHierarchySeparatorAllowOverride)
	viper.SetDefault("markorSecrets.hierarchySeparator.warnOnOverride", DefaultMarkorSecretsHierarchySeparatorWarnOnOverride)
	viper.SetDefault("markorSecrets.managedAnnotation.default", DefaultMarkorSecretsManagedAnnotationDefault)
	viper.SetDefault("markorSecrets.managedAnnotation.allowOverride", DefaultMarkorSecretsManagedAnnotationAllowOverride)
	viper.SetDefault("markorSecrets.managedAnnotation.warnOnOverride", DefaultMarkorSecretsManagedAnnotationWarnOnOverride)
}
