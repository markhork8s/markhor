package config

import (
	"testing"

	"github.com/imdario/mergo"
)

var validConfig = Config{
	Kubernetes: kubernetesConfig{
		KubeconfigPath:        "", // empty means you are running in the cluster
		ClusterTimeoutSeconds: 23,
	},
	Logging: loggingConfig{
		Level:              "info",
		Style:              "json",
		LogToStdout:        true,
		AdditionalLogFiles: []string{"file1.log", "file2.log"},
	},
	Behavior: behaviorConfig{
		Fieldmanager: fieldManagerConfig{
			Name:         "fieldManager",
			ForceUpdates: true,
		},
		Namespaces:         []string{"namespace1", "namespace2"},
		ExcludedNamespaces: []string{"excluded1", "excluded2"},
	},
	MarkorSecrets: markhorSecretsConfig{
		HierarchySeparator: defaultOverrideStruct{
			Default:        "/",
			AllowOverride:  true,
			WarnOnOverride: false,
		},
		ManagedLabel: defaultOverrideStruct{
			Default:        "managed",
			AllowOverride:  true,
			WarnOnOverride: false,
		},
	},
	Healthcheck: HealthcheckConfig{
		Port:    9091,
		Enabled: true,
	},
	AdmissionController: AdmissionControllerConfig{
		Port:    4040,
		Enabled: true,
	},
	Tls: TlsConfig{
		Enabled:  true,
		CertPath: "/some/path",
		KeyPath:  "/some/other/path",
	},
}

func TestValidateConfig_Valid(t *testing.T) {
	// Validate the configuration
	err := ValidateConfig(validConfig)

	if err != nil {
		t.Errorf("Expected valid configuration, but got error: %v", err)
	}
}

func TestValidateConfig_Invalid(t *testing.T) {
	t.Parallel()
	invalidConfigs := []struct {
		string
		Config
	}{
		{
			"Negative k8s timeout",
			Config{
				Kubernetes: kubernetesConfig{
					ClusterTimeoutSeconds: -1,
				},
			},
		},
		{
			"Healthcheck port out of range (-10)",
			Config{
				Healthcheck: HealthcheckConfig{
					Port: -10, // Less than 1
				},
			},
		},
		{
			"Healthcheck port out of range (65600)",
			Config{
				Healthcheck: HealthcheckConfig{
					Port: maxPort + 1, // Over the max
				},
			},
		},
		{
			"AdmissionController port out of range (-10)",
			Config{
				AdmissionController: AdmissionControllerConfig{
					Port: -10, // Less than 1
				},
			},
		},
		{
			"AdmissionController port out of range (65600)",
			Config{
				AdmissionController: AdmissionControllerConfig{
					Port: maxPort + 1, // Over the max
				},
			},
		},
		{
			"AdmissionController and Healthcheck on the same port",
			Config{
				AdmissionController: AdmissionControllerConfig{
					Port: 8000,
				},
				Healthcheck: HealthcheckConfig{
					Port: 8000,
				},
			},
		},
		{
			"Invalid log level",
			Config{
				Logging: loggingConfig{Level: "invalid logging level"},
			},
		},
		{
			"Invalid additional log file path length",
			Config{
				Logging: loggingConfig{
					AdditionalLogFiles: []string{"valid", ""},
				},
			},
		},
		{
			"Duplicated additional log files",
			Config{
				Logging: loggingConfig{
					AdditionalLogFiles: []string{"file1.log", "file1.log"},
				},
			},
		},
		{
			"Empty namespace in Behavior.Namespaces",
			Config{
				Behavior: behaviorConfig{
					Namespaces: []string{"valid", ""},
				},
			},
		},
		{
			"Duplicated namespaces in Behavior.Namespaces",
			Config{
				Behavior: behaviorConfig{
					Namespaces: []string{"ns1", "ns2", "ns1"},
				},
			},
		},
		{
			"Empty namespace in Behavior.ExcludedNamespaces",
			Config{
				Behavior: behaviorConfig{
					ExcludedNamespaces: []string{""},
				},
			},
		},
		{
			"Duplicated namespaces in Behavior.ExcludedNamespaces",
			Config{
				Behavior: behaviorConfig{
					ExcludedNamespaces: []string{"ns1", "ns2", "ns1"},
				},
			},
		},
	}

	for _, item := range invalidConfigs {
		c := item.Config
		t.Run("Invalid config: "+item.string, func(t *testing.T) {
			t.Parallel()
			v := validConfig
			mergo.Merge(&c, v)
			err := ValidateConfig(c)
			if err == nil {
				t.Error("Expected invalid configuration, but got no error")
			}
		})
	}

	v1 := validConfig
	v1.Healthcheck.Port = 0
	err := ValidateConfig(v1)
	if err == nil {
		t.Error("Expected invalid configuration, healthcheck port out of range (0), but got no error")
	}
	v1 = validConfig
	v1.AdmissionController.Port = 0
	err = ValidateConfig(v1)
	if err == nil {
		t.Error("Expected invalid configuration, admission controller port out of range (0), but got no error")
	}
	v1 = validConfig
	v1.Behavior.Fieldmanager.Name = ""
	err = ValidateConfig(v1)
	if err == nil {
		t.Error("Expected invalid configuration, empty field manager name, but got no error")
	}
	v1 = validConfig
	v1.Tls.CertPath = ""
	err = ValidateConfig(v1)
	if err == nil {
		t.Error("Expected invalid configuration, TLS certificate path with zero length, but got no error")
	}
	v1 = validConfig
	v1.Tls.KeyPath = ""
	err = ValidateConfig(v1)
	if err == nil {
		t.Error("Expected invalid configuration, TLS key path with zero length, but got no error")
	}
	v1 = validConfig
	v1.MarkorSecrets.HierarchySeparator.Default = ""
	err = ValidateConfig(v1)
	if err == nil {
		t.Error("Expected invalid configuration, empty default hierarchy separator, but got no error")
	}
	v1 = validConfig
	v1.MarkorSecrets.ManagedLabel.Default = ""
	err = ValidateConfig(v1)
	if err == nil {
		t.Error("Expected invalid configuration, empty default managed label, but got no error")
	}
}
