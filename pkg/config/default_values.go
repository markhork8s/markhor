package config

const (
	DefaultKubeconfigPath                                = ""
	DefaultClusterTimeoutSeconds                         = 10
	DefaultSopsKeysPath                                  = "~/.config/sops/keys"
	DefaultHealthcheckPort                               = 8080
	DefaultHealthcheckEnabled                            = true
	DefaultLoggingLevel                                  = "info"
	DefaultLoggingStyle                                  = "text"
	DefaultLoggingLogToStdout                            = true
	DefaultBehaviorFieldManagerName                      = "github.com/civts/markhor"
	DefaultBehaviorFieldManagerForceUpdates              = false
	DefaultBehaviorPruneDanglingMarkhorSecrets           = true
	DefaultMarkorSecretsHierarchySeparatorDefault        = "/"
	DefaultMarkorSecretsHierarchySeparatorAllowOverride  = false
	DefaultMarkorSecretsHierarchySeparatorWarnOnOverride = true
	DefaultMarkorSecretsManagedAnnotationDefault         = "markhor.example.com/managed-by"
	DefaultMarkorSecretsManagedAnnotationAllowOverride   = false
	DefaultMarkorSecretsManagedAnnotationWarnOnOverride  = true
)

// Can't make these costants
var (
	DefaultLoggingAdditionalLogFiles  = []string{}
	DefaultBehaviorNamespaces         = []string{}
	DefaultBehaviorExcludedNamespaces = []string{}
)
