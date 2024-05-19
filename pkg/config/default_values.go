package config

const (
	DefaultKubeconfigPath                                = ""
	DefaultClusterTimeoutSeconds                         = 10
	DefaultHealthcheckPort                               = 8000
	DefaultHealthcheckEnabled                            = true
	DefaultAdmissionControllerPort                       = 443
	DefaultAdmissionControllerEnabled                    = true
	DefaultTlsEnabled                                    = false
	DefaultTlsCertPath                                   = "/etc/markhor/tls/server.crt"
	DefaultTlsKeyPath                                    = "/etc/markhor/tls/server.key"
	DefaultLoggingLevel                                  = "info"
	DefaultLoggingStyle                                  = "text"
	DefaultLoggingLogToStdout                            = true
	DefaultBehaviorFieldManagerName                      = "markhork8s.github.io"
	DefaultBehaviorFieldManagerForceUpdates              = false
	DefaultMarkorSecretsHierarchySeparatorDefault        = "/"
	DefaultMarkorSecretsHierarchySeparatorAllowOverride  = false
	DefaultMarkorSecretsHierarchySeparatorWarnOnOverride = true
	//#nosec G101 -- This is a false positive (gosec)
	DefaultMarkorSecretsManagedLabelDefault        = "app.kubernetes.io/managed-by"
	DefaultMarkorSecretsManagedLabelAllowOverride  = false
	DefaultMarkorSecretsManagedLabelWarnOnOverride = true
)

// Can't make these costants
var (
	DefaultLoggingAdditionalLogFiles  = []string{}
	DefaultBehaviorNamespaces         = []string{}
	DefaultBehaviorExcludedNamespaces = []string{}
)
