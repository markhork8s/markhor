package pkg

import (
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Sops struct {
		KeysPath string `mapstructure:"keysPath"`
	} `mapstructure:"sops"`

	Kubernetes struct {
		KubeconfigPath        string `mapstructure:"kubeconfigPath"`
		ClusterTimeoutSeconds int    `mapstructure:"clusterTimeoutSeconds"`
	} `mapstructure:"kubernetes"`

	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`

	Behavior struct {
		Fieldmanager struct {
			Name         string `mapstructure:"name"`
			ForceUpdates bool   `mapstructure:"forceUpdates"`
		} `mapstructure:"fieldmanager"`

		PruneDanglingMarkhorSecrets bool     `mapstructure:"pruneDanglingMarkhorSecrets"`
		Namespaces                  []string `mapstructure:"namespaces"`
		ExcludedNamespaces          []string `mapstructure:"excludedNamespaces"`
	} `mapstructure:"behavior"`

	MarkorSecrets struct {
		HierarchySeparator struct {
			Default        string `mapstructure:"default"`
			AllowOverride  bool   `mapstructure:"allowOverride"`
			WarnOnOverride bool   `mapstructure:"warnOnOverride"`
		} `mapstructure:"hierarchySeparator"`

		ManagedAnnotation struct {
			Default        string `mapstructure:"default"`
			AllowOverride  bool   `mapstructure:"allowOverride"`
			WarnOnOverride bool   `mapstructure:"warnOnOverride"`
		} `mapstructure:"managedAnnotation"`
	} `mapstructure:"markorSecrets"`

	Healthcheck HealthcheckConfig `mapstructure:"healthcheck"`
}

type HealthcheckConfig struct {
	Port    int  `mapstructure:"port"`
	Enabled bool `mapstructure:"enabled"`
}

// Parses the program configuration (from CLI and file).
// Will terminate the execution if it fails since it would not be illogical to
// proceed w/o a valid config
func ParseConfig() Config {
	ParseCliArgs()
	configFilePath := viper.GetString("config")

	if configFilePath != "" {
		log.Println("Reading Markhor config from user-defined path:", configFilePath)
		viper.SetConfigFile(configFilePath)
	} else {
		log.Println("Reading Markhor config from default path:", DEFAULT_CONFIG_PATH)
		viper.SetConfigFile(DEFAULT_CONFIG_PATH)
	}

	setDefaultConfigValues()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Error reading Markhor config file:", err)
		os.Exit(1)
	}

	// Define custom Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println("Error parsing Markhor config data:", err)
		os.Exit(1)
	}

	return config
}

// Defines the default Markhor config values
func setDefaultConfigValues() {
	viper.SetDefault("kubernetes.kubeconfigPath", "")
	viper.SetDefault("kubernetes.clusterTimeoutSeconds", 10)
	viper.SetDefault("sops.keysPath", "~/.config/sops/keys")
	viper.SetDefault("healthcheck.port", 8080)
	viper.SetDefault("healthcheck.enabled", true)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("behavior.fieldmanager.name", "github.com/civts/markhor")
	viper.SetDefault("behavior.namespaces", make([]string, 0))
	viper.SetDefault("behavior.excludedNamespaces", make([]string, 0))
	viper.SetDefault("behavior.fieldmanager.forceUpdates", false)
	viper.SetDefault("behavior.pruneDanglingMarkhorSecrets", true)
	viper.SetDefault("markorSecrets.hierarchySeparator.default", "/")
	viper.SetDefault("markorSecrets.hierarchySeparator.allowOverride", false)
	viper.SetDefault("markorSecrets.hierarchySeparator.warnOnOverride", true)
	viper.SetDefault("markorSecrets.managedAnnotation.default", "markhor.example.com/managed-by")
	viper.SetDefault("markorSecrets.managedAnnotation.allowOverride", false)
	viper.SetDefault("markorSecrets.managedAnnotation.warnOnOverride", true)
}

func ParseCliArgs() {
	// Define CLI flags
	pflag.StringP("config", "c", "", "Path to config file")
	helpSet := pflag.BoolP("help", "h", false, "Show this help message")
	versionSet := pflag.BoolP("version", "v", false, "Print the version of this program and exit")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Print help or version message if needed
	if *helpSet {
		pflag.PrintDefaults()
		os.Exit(0)
	} else if *versionSet {
		log.Printf("v%s\n", VERSION)
		os.Exit(0)
	}
}
