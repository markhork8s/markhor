package config

type Config struct {
	Kubernetes kubernetesConfig `mapstructure:"kubernetes"`

	Logging loggingConfig `mapstructure:"logging"`

	Behavior behaviorConfig `mapstructure:"behavior"`

	MarkorSecrets markhorSecretsConfig `mapstructure:"markorSecrets"`

	Healthcheck HealthcheckConfig `mapstructure:"healthcheck"`

	AdmissionController AdmissionControllerConfig `mapstructure:"admissionController"`

	Tls TlsConfig `mapstructure:"tls"`
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
	ManagedLabel       defaultOverrideStruct `mapstructure:"managedLabel"`
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

type AdmissionControllerConfig struct {
	Port    int  `mapstructure:"port"`
	Enabled bool `mapstructure:"enabled"`
}

type TlsConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertPath string `mapstructure:"certPath"`
	KeyPath  string `mapstructure:"keyPath"`
}
