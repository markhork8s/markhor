package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/civts/markhor/pkg"
	"github.com/google/go-cmp/cmp"
	"github.com/mohae/deepcopy"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func reset() {
	pflag.CommandLine = &pflag.FlagSet{}
	os.Args = []string{"./markhor"}
	viper.Reset()
}

// The configuration with all the default values
var defaultConfig = &Config{
	Kubernetes: kubernetesConfig{
		KubeconfigPath:        DefaultKubeconfigPath,
		ClusterTimeoutSeconds: DefaultClusterTimeoutSeconds,
	},
	Logging: loggingConfig{
		Level:              DefaultLoggingLevel,
		Style:              DefaultLoggingStyle,
		LogToStdout:        DefaultLoggingLogToStdout,
		AdditionalLogFiles: DefaultLoggingAdditionalLogFiles,
	},
	Behavior: behaviorConfig{
		Fieldmanager: fieldManagerConfig{
			Name:         DefaultBehaviorFieldManagerName,
			ForceUpdates: DefaultBehaviorFieldManagerForceUpdates,
		},
		PruneDanglingMarkhorSecrets: DefaultBehaviorPruneDanglingMarkhorSecrets,
		Namespaces:                  DefaultBehaviorNamespaces,
		ExcludedNamespaces:          DefaultBehaviorExcludedNamespaces,
	},
	MarkorSecrets: markhorSecretsConfig{
		HierarchySeparator: defaultOverrideStruct{
			Default:        DefaultMarkorSecretsHierarchySeparatorDefault,
			AllowOverride:  DefaultMarkorSecretsHierarchySeparatorAllowOverride,
			WarnOnOverride: DefaultMarkorSecretsHierarchySeparatorWarnOnOverride,
		},
		ManagedAnnotation: defaultOverrideStruct{
			Default:        DefaultMarkorSecretsManagedAnnotationDefault,
			AllowOverride:  DefaultMarkorSecretsManagedAnnotationAllowOverride,
			WarnOnOverride: DefaultMarkorSecretsManagedAnnotationWarnOnOverride,
		},
	},
	Healthcheck: HealthcheckConfig{
		Port:    DefaultHealthcheckPort,
		Enabled: DefaultHealthcheckEnabled,
	},
	AdmissionController: AdmissionControllerConfig{
		Port:    DefaultAdmissionControllerPort,
		Enabled: DefaultAdmissionControllerEnabled,
	},
	Tls: TlsConfig{
		Mode:     DefaultTlsMode,
		CertPath: DefaultTlsCertPath,
		KeyPath:  DefaultTlsKeyPath,
	},
}

// When no configuration file is provided, no error should occour in
// the configuration parsing phase
func TestValidDefaultConfig(t *testing.T) {
	reset()
	ensureDefaultConfigDoesNotExist(t)
	_, err := ParseConfig()
	if err != nil {
		t.Fatal("Parsing the config should not have failed, but instead we got this error:", err)
	}
}

// When no configuration file is provided, no error should occour in
// the configuration parsing phase AND the resulting configuration should use
// exactly the deafult values
func TestParseConfigDefaultValuesOnMissingDefaultFile(t *testing.T) {
	reset()
	ensureDefaultConfigDoesNotExist(t)

	// Parse the config without the --config flag specified
	config, err := ParseConfig()
	if err != nil {
		t.Fatal("There was an unexpected error parsing the config: ", err)
	}

	// Add assertions to validate the returned Config struct
	diff := cmp.Diff(config, defaultConfig)
	if diff != "" {
		t.Fatal("The parsed configuration is not equal to the default one:", diff)
	}
}

// Should use all the default values when an empty file is specified
func TestValidDefaultConfigOnEmptyFile(t *testing.T) {
	reset()
	// Create temporary file with the config data
	f, err := createTempFile()
	if err != nil {
		t.Fatalf("Error creating temporary config file: %v", err)
	}
	yamlName := fmt.Sprintf("%s.yaml", f.Name())
	err = os.Rename(f.Name(), yamlName)
	if err != nil {
		t.Fatalf("Error renaming temporary config file: %v", err)
	}

	defer removeTempFile(f)

	// Fake running the program with these CLI args
	// The 1st one is the program name, which is irrelevant now
	os.Args = []string{">", "--config", yamlName}

	config, err := ParseConfig()
	if err != nil {
		t.Fatal("There was an unexpected error parsing the config: ", err)
	}

	// Add assertions to validate the returned Config struct
	diff := cmp.Diff(config, defaultConfig)
	if diff != "" {
		t.Fatal("The parsed configuration is not equal to the default one:", diff)
	}
}

// When a valid config file is specified, it shall use the values specified
// therein and the default values for the other fields
func TestParseConfigValidPartialFile(t *testing.T) {
	reset()
	// Init with values that are surely different from the default ones
	k8sConfigPath := "/a/custom/config/path"
	configFileContents := "kubernetes:\n  kubeconfigPath: " + k8sConfigPath + "\n"

	// Create temporary file with the config data
	f, err := createTempFile()
	if err != nil {
		t.Fatalf("Error creating temporary config file: %v", err)
	}
	yamlName := fmt.Sprintf("%s.yaml", f.Name())
	err = os.Rename(f.Name(), yamlName)
	if err != nil {
		t.Fatalf("Error renaming temporary config file: %v", err)
	}
	_, err = f.Write([]byte(configFileContents))
	if err != nil {
		t.Fatalf("Error writing to temporary config file: %v", err)
	}

	defer removeTempFile(f)

	// Fake running the program with these CLI args
	// The 1st one is the program name, which is irrelevant now
	os.Args = []string{">", "--config", yamlName}

	// Test the ParseConfig function with the temporary config file
	config, err := ParseConfig()
	if err != nil {
		t.Fatal("There was an unexpected error parsing the config: ", err)
	}

	if config.Kubernetes.KubeconfigPath != k8sConfigPath {
		t.Fatalf("The parsed config does not use the provided value for kubernetes config path: expected %s, got %s", k8sConfigPath, config.Kubernetes.KubeconfigPath)
	}
	customDefaultConfig := deepcopy.Copy(defaultConfig).(*Config)
	customDefaultConfig.Kubernetes.KubeconfigPath = k8sConfigPath
	diff := cmp.Diff(config, customDefaultConfig)
	if diff != "" {
		t.Fatal("The two structs are not equal:", diff)
	}
}

// When a valid config file is specified, it shall use the values specified
// therein and the default values for the other fields (same as the previous
// test but w/ another key to ensure all are covered)
func TestParseConfigValidPartialFile2(t *testing.T) {
	reset()
	// Init with values that are surely different from the default ones
	name := "TheBestFieldManagerName"
	configFileContents := "behavior:\n  fieldmanager:\n    name: " + name

	// Create temporary file with the config data
	f, err := createTempFile()
	if err != nil {
		t.Fatalf("Error creating temporary config file: %v", err)
	}
	yamlName := fmt.Sprintf("%s.yaml", f.Name())
	err = os.Rename(f.Name(), yamlName)
	if err != nil {
		t.Fatalf("Error renaming temporary config file: %v", err)
	}
	_, err = f.Write([]byte(configFileContents))
	if err != nil {
		t.Fatalf("Error writing to temporary config file: %v", err)
	}

	defer removeTempFile(f)

	// Fake running the program with these CLI args
	// The 1st one is the program name, which is irrelevant now
	os.Args = []string{">", "--config", yamlName}

	// Test the ParseConfig function with the temporary config file
	config, err := ParseConfig()
	if err != nil {
		t.Fatal("There was an unexpected error parsing the config: ", err)
	}

	if config.Behavior.Fieldmanager.Name != name {
		t.Fatalf("The parsed config does not use the provided value for sops keys path: expected %s, got %s", name, config.Behavior.Fieldmanager.Name)
	}
	customDefaultConfig := deepcopy.Copy(defaultConfig).(*Config)
	customDefaultConfig.Behavior.Fieldmanager.Name = name
	diff := cmp.Diff(config, customDefaultConfig)
	if diff != "" {
		t.Fatal("The two structs are not equal:", diff)
	}
}

// When a valid configuration file is given, its values shall be taken into
// account
func TestParseConfigValidCompleteFile(t *testing.T) {
	reset()
	// Init with values that are surely different from the default ones
	healthConfig := HealthcheckConfig{
		Port:    DefaultHealthcheckPort / 2,
		Enabled: !DefaultHealthcheckEnabled,
	}
	admissionControllerConfig := AdmissionControllerConfig{
		Port:    DefaultAdmissionControllerPort / 2,
		Enabled: !DefaultAdmissionControllerEnabled,
	}
	tlsConfig := TlsConfig{
		Mode:     "file",
		CertPath: defaultConfig.Tls.CertPath + ".new",
		KeyPath:  defaultConfig.Tls.KeyPath + ".new_again",
	}
	kubernetesConfig := kubernetesConfig{
		KubeconfigPath:        DefaultKubeconfigPath + "and/more",
		ClusterTimeoutSeconds: DefaultClusterTimeoutSeconds*2 + 1,
	}
	f2, err := createTempFile()
	if err != nil {
		t.Fatalf("Error creating temporary config file: %v", err)
	}
	defer removeTempFile(f2)
	loggingConfig := loggingConfig{
		Level:              "error",
		Style:              "json",
		LogToStdout:        !DefaultLoggingLogToStdout,
		AdditionalLogFiles: append(DefaultLoggingAdditionalLogFiles, f2.Name()),
	}
	if loggingConfig.Level == DefaultLoggingLevel {
		t.Fatal("For this test, the logging level should be different from the default")
	} else if loggingConfig.Style == DefaultLoggingStyle {
		t.Fatal("For this test, the logging style should be different from the default")
	}
	fieldManagerConfig := fieldManagerConfig{
		Name:         DefaultBehaviorFieldManagerName + "but better",
		ForceUpdates: !DefaultBehaviorFieldManagerForceUpdates,
	}
	behaviorConfig := behaviorConfig{
		Fieldmanager:                fieldManagerConfig,
		PruneDanglingMarkhorSecrets: !DefaultBehaviorPruneDanglingMarkhorSecrets,
		Namespaces:                  append(DefaultBehaviorNamespaces, "example"),
		ExcludedNamespaces:          append(DefaultBehaviorExcludedNamespaces, "example2"),
	}
	hierarchySeparatorConfig := defaultOverrideStruct{
		Default:        DefaultMarkorSecretsHierarchySeparatorDefault + "-/",
		AllowOverride:  !DefaultMarkorSecretsHierarchySeparatorAllowOverride,
		WarnOnOverride: !DefaultMarkorSecretsHierarchySeparatorWarnOnOverride,
	}
	managedConfig := defaultOverrideStruct{
		Default:        DefaultMarkorSecretsManagedAnnotationDefault + "!",
		AllowOverride:  !DefaultMarkorSecretsManagedAnnotationAllowOverride,
		WarnOnOverride: !DefaultMarkorSecretsManagedAnnotationWarnOnOverride,
	}
	mSecretsConfig := markhorSecretsConfig{
		HierarchySeparator: hierarchySeparatorConfig,
		ManagedAnnotation:  managedConfig,
	}
	xConfig := Config{
		Kubernetes:          kubernetesConfig,
		Logging:             loggingConfig,
		Behavior:            behaviorConfig,
		MarkorSecrets:       mSecretsConfig,
		Healthcheck:         healthConfig,
		AdmissionController: admissionControllerConfig,
		Tls:                 tlsConfig,
	}
	yamlConf, err := yaml.Marshal(xConfig)
	if err != nil {
		t.Fatalf("Error converting cnf to YAML: %v", err)
	}

	// Create temporary file with the config data
	f, err := createTempFile()
	if err != nil {
		t.Fatalf("Error creating temporary config file: %v", err)
	}
	yamlName := fmt.Sprintf("%s.yaml", f.Name())
	err = os.Rename(f.Name(), yamlName)
	if err != nil {
		t.Fatalf("Error renaming temporary config file: %v", err)
	}
	_, err = f.Write(yamlConf)
	if err != nil {
		t.Fatalf("Error writing to temporary config file: %v", err)
	}

	defer removeTempFile(f)

	// Fake running the program with these CLI args
	// The 1st one is the program name, which is irrelevant now
	os.Args = []string{">", "--config", yamlName}

	// Test the ParseConfig function with the temporary config file
	config, err := ParseConfig()
	if err != nil {
		t.Fatal("There was an unexpected error parsing the config: ", err)
	}

	diff := cmp.Diff(config, &xConfig)
	if diff != "" {
		t.Fatal("The two structs are not equal:", diff)
	}
}

// If the file has an invalid yaml syntax, the program should exit with code 1
func TestParseConfigInvalidYamlFile(t *testing.T) {
	reset()
	configFileContents := "this is no yaml file for sure"

	// Create temporary file with the config data
	f, err := createTempFile()
	if err != nil {
		t.Fatalf("Error creating temporary config file: %v", err)
	}
	yamlName := fmt.Sprintf("%s.yaml", f.Name())
	err = os.Rename(f.Name(), yamlName)
	if err != nil {
		t.Fatalf("Error renaming temporary config file: %v", err)
	}
	_, err = f.Write([]byte(configFileContents))
	if err != nil {
		t.Fatalf("Error writing to temporary config file: %v", err)
	}

	defer removeTempFile(f)

	// Fake running the program with these CLI args
	// The 1st one is the program name, which is irrelevant now
	os.Args = []string{">", "--config", yamlName}

	// Test the ParseConfig function with the temporary config file
	_, err = ParseConfig()
	if err == nil {
		t.Fatalf("Parsing this invalid config should have failed, but it did not")
	}
}

// If the file has an invalid yaml syntax, the program should exit with code 1
func TestParseConfigInvalidConfigFile(t *testing.T) {
	reset()
	configFileContents := "kubernetes:\n  clusterTimeoutSeconds: \"oops, wrong value\""

	// Create temporary file with the config data
	f, err := createTempFile()
	if err != nil {
		t.Fatalf("Error creating temporary config file: %v", err)
	}
	yamlName := fmt.Sprintf("%s.yaml", f.Name())
	err = os.Rename(f.Name(), yamlName)
	if err != nil {
		t.Fatalf("Error renaming temporary config file: %v", err)
	}
	_, err = f.Write([]byte(configFileContents))
	if err != nil {
		t.Fatalf("Error writing to temporary config file: %v", err)
	}

	defer removeTempFile(f)

	// Fake running the program with these CLI args
	// The 1st one is the program name, which is irrelevant now
	os.Args = []string{">", "--config", yamlName}

	// Test the ParseConfig function with the temporary config file
	_, err = ParseConfig()
	if err == nil {
		t.Fatalf("Parsing this invalid config should have failed, but it did not")
	}
}

func createTempFile() (*os.File, error) {
	tmpDir := os.TempDir()
	return os.CreateTemp(tmpDir, "deleteme_")
}

func removeTempFile(f *os.File) {
	os.Remove(f.Name())
	f.Close()
}

func ensureDefaultConfigDoesNotExist(t *testing.T) {
	if _, err := os.Stat(pkg.DEFAULT_CONFIG_PATH); err == nil {
		t.Fatal("A prerequisite for this test is that this file does NOT exist, but it seems that it does. Please remove", pkg.DEFAULT_CONFIG_PATH)
	}
}
