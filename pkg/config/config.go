package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/civts/markhor/pkg"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Sops sopsConfig `mapstructure:"sops"`

	Kubernetes kubernetesConfig `mapstructure:"kubernetes"`

	Logging loggingConfig `mapstructure:"logging"`

	Behavior behaviorConfig `mapstructure:"behavior"`

	MarkorSecrets markhorSecretsConfig `mapstructure:"markorSecrets"`

	Healthcheck HealthcheckConfig `mapstructure:"healthcheck"`
}

type sopsConfig struct {
	KeysPath string `mapstructure:"keysPath"`
}

type kubernetesConfig struct {
	KubeconfigPath        string `mapstructure:"kubeconfigPath"`
	ClusterTimeoutSeconds int    `mapstructure:"clusterTimeoutSeconds"`
}

type loggingConfig struct {
	Level              string   `mapstructure:"level"`
	Style              string   `mapstructure:"style"`
	LogToStdout        bool     `mapstructure:"logToStdout"`
	AdditionalLogFiles []string `mapstructure:"additionalLogFiles"`
}

type behaviorConfig struct {
	Fieldmanager fieldManagerConfig `mapstructure:"fieldmanager"`

	PruneDanglingMarkhorSecrets bool     `mapstructure:"pruneDanglingMarkhorSecrets"`
	Namespaces                  []string `mapstructure:"namespaces"`
	ExcludedNamespaces          []string `mapstructure:"excludedNamespaces"`
}

type fieldManagerConfig struct {
	Name         string `mapstructure:"name"`
	ForceUpdates bool   `mapstructure:"forceUpdates"`
}

type markhorSecretsConfig struct {
	HierarchySeparator defaultOverrideStruct `mapstructure:"hierarchySeparator"`
	ManagedAnnotation  defaultOverrideStruct `mapstructure:"managedAnnotation"`
}

type defaultOverrideStruct struct {
	Default        string `mapstructure:"default"`
	AllowOverride  bool   `mapstructure:"allowOverride"`
	WarnOnOverride bool   `mapstructure:"warnOnOverride"`
}

type HealthcheckConfig struct {
	Port    int  `mapstructure:"port"`
	Enabled bool `mapstructure:"enabled"`
}

// Parses the program configuration (from CLI and file).
// Will terminate the execution if it fails since it would not be illogical to
// proceed w/o a valid config
func ParseConfig() (*Config, error) {
	err := ParseCliArgs()
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

	SetupLogging(config)

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

func ParseCliArgs() error {
	// Define CLI flags
	pflag.StringP("config", "c", pkg.DEFAULT_CONFIG_PATH, "Path to config file")
	helpSet := pflag.BoolP("help", "h", false, "Show this help message")
	versionSet := pflag.BoolP("version", "v", false, "Print the version of this program and exit")
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return errors.New("Could not parse CLI flags: " + err.Error())
	}

	// Print help or version message if needed
	if *helpSet {
		pflag.PrintDefaults()
		os.Exit(0)
	} else if *versionSet {
		fmt.Printf("v%s\n", pkg.VERSION)
		os.Exit(0)
	}

	return nil
}
