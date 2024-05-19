package decrypt

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "github.com/markhork8s/markhor/pkg/api/types/v1"
	"github.com/markhork8s/markhor/pkg/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Yes, this key is used only in these tests
// The corresponding public key is: age1apq7ck5adq6dkd0c242phl42fsurvpxvt9pwk0qg7ahdex7fqppqj8pe8y
const testSk = "AGE-SECRET-KEY-1LYQ3PW2AKP02VC6WLV64944NLJPL33CS7DHJXK6GPWW8F70G0GWQ7NQMS8"

// This was encrypted with testSk. The fields and subfields have been rearranged on purpose
const cypheredData = `{
	"metadata": {
		"namespace": "default",
		"name": "sample-secret"
	},
	"apiVersion": "markhork8s.github.io/v1",
	"kind": "MarkhorSecret",
	"markhorParams": {
			"order": [
					"apiVersion",
					"kind",
					"metadata/name",
					"metadata/namespace",
					"markhorParams/order",
					"type",
					"data/session_secret",
					"stringData/another"
			]
	},
	"stringData": {
		"another": "ENC[AES256_GCM,data:EA1mlHFtMLkvI2aNUYpmIM7Mv2O+1g==,iv:9Aue7n6LHR0jKdHE3t1W1GcAciREm4cnHpxRAAKVeRo=,tag:S+/7V/y0RwdcAK3H2cVlnA==,type:str]"
	},
	"type": "Opaque",
	"sops": {
		"kms": null,
		"gcp_kms": null,
		"azure_kv": null,
		"hc_vault": null,
		"age": [
			{
				"recipient": "age1apq7ck5adq6dkd0c242phl42fsurvpxvt9pwk0qg7ahdex7fqppqj8pe8y",
				"enc": "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBxNTlLenRubTBLOEN2UEln\nMWdXejRkZHA4aUR1ZjNFR2pIS3NSdGJhRlZZClY5dytUSUlQOHdtZHI5Y3FCelJQ\nOW9PQTlVLzQrVWNwcDRpQnJIQkU4VDQKLS0tIFhXUThlT0pUaVlLcG16TWdlZkdC\naEZ5ZGNsK1hMbG5KYlJBUWJjY0xlK28KzhouuUJfFQd6RlKne8j+QyalS57xcs/h\nrPu3m9UtBsd3ChFq+eAAAkwDM3DFslII89vF8QxVrYE0svUCYfWJ4A==\n-----END AGE ENCRYPTED FILE-----\n"
			}
		],
		"lastmodified": "2024-05-19T12:03:21Z",
		"mac": "ENC[AES256_GCM,data:VY3TNEp2lteZ2tbEXLJOkj1ZIRlxSBMpbjAlaNP/4GH84XOexQUDJo8E/pAnaF5uqwHN816+a8w7SNCyl6CcpNic1Y5/zalThgfEKWxwjw0Mfk6hvmEmRlLRrfuWczXkLvGAZQT0l5V3xUegiTthN2Zxg56jqQBgfvxcDIv2ojA=,iv:jN32v++w3FADErET+7StkljdNP0n5p7nICFFpnKP804=,tag:w+XS31ng9ur6AkuqT7dUQw==,type:str]",
		"pgp": null,
		"encrypted_regex": "^(data|stringData)$",
		"version": "3.8.1"
	},
	"data": {
		"session_secret": "ENC[AES256_GCM,data:U2E3upldktVYnt10g4GL3rvMM8/5jfy+wamCsdvnFGypj51tRcKcbQ==,iv:yg5rzxaiwsWLVprFC8T6Oma8GxAYJQRiAPnUj+pEGmY=,tag:FfvHLDc3G33BYB5JExQlZw==,type:str]"
	}
}`

const expectedJSON = `{
	"apiVersion": "markhork8s.github.io/v1",
	"kind": "MarkhorSecret",
	"metadata": {
			"name": "sample-secret",
			"namespace": "default"
	},
	"markhorParams": {
			"order": [
					"apiVersion",
					"kind",
					"metadata/name",
					"metadata/namespace",
					"markhorParams/order",
					"type",
					"data/session_secret",
					"stringData/another"
			]
	},
	"type": "Opaque",
	"data": {
			"session_secret": "aHR0cHM6Ly95b3V0dS5iZS9kUXc0dzlXZ1hjUT8="
	},
	"stringData": {
			"another": "I want some pineapples"
	}
}`

var ms = config.MarkhorSecretsConfig{
	HierarchySeparator: config.DefaultOverrideStruct{
		Default:        "/",
		AllowOverride:  true,
		WarnOnOverride: false,
	},
	ManagedLabel: config.DefaultOverrideStruct{
		Default:        "unimportant",
		AllowOverride:  true,
		WarnOnOverride: false,
	},
}

func TestDecryptMarkhorSecret_Works_With_Correct_Input(t *testing.T) {
	// Yes, this key is used only in these tests
	const SOPS_KEY = "SOPS_AGE_KEY"
	prevKey := os.Getenv(SOPS_KEY)
	defer os.Setenv(SOPS_KEY, prevKey)
	os.Setenv(SOPS_KEY, testSk)

	cd := make(map[string]interface{})
	err := json.Unmarshal([]byte(cypheredData), &cd)
	if err != nil {
		t.Fatal("Failed to unmarshal encrypted JSON", err)
	}

	res, err := DecryptMarkhorSecret(cd, ms, slog.String("eid", "_"))
	if err != nil {
		t.Fatal("Failed to decrypt markhor secret:", err)
	}

	expectedObj := make(map[string]interface{})
	err = json.Unmarshal([]byte(expectedJSON), &expectedObj)
	if err != nil {
		t.Fatal("Failed to unmarshal expected JSON", err)
	}
	diff := cmp.Diff(expectedObj, res)
	if diff != "" {
		t.Fatal("The decryption did not yield the expected result", diff)
	}
}

func TestDecryptMarkhorSecret_Fails_With_Missing_Key(t *testing.T) {
	const SOPS_KEY = "SOPS_AGE_KEY"
	prevKey := os.Getenv(SOPS_KEY)
	defer os.Setenv(SOPS_KEY, prevKey)
	os.Setenv(SOPS_KEY, "")
	cd := make(map[string]interface{})
	err := json.Unmarshal([]byte(cypheredData), &cd)
	if err != nil {
		t.Fatal("Failed to unmarshal encrypted JSON", err)
	}

	_, err = DecryptMarkhorSecret(cd, ms, slog.String("eid", "_"))
	if err == nil {
		t.Fatal("The decryption should have failed since we did not provide the secret key", err)
	}
}

func TestDecryptMarkhorSecret_Fails_With_Wrong_Key(t *testing.T) {
	const SOPS_KEY = "SOPS_AGE_KEY"
	prevKey := os.Getenv(SOPS_KEY)
	defer os.Setenv(SOPS_KEY, prevKey)
	os.Setenv(SOPS_KEY, "surely not the right key")
	cd := make(map[string]interface{})
	err := json.Unmarshal([]byte(cypheredData), &cd)
	if err != nil {
		t.Fatal("Failed to unmarshal encrypted JSON", err)
	}

	_, err = DecryptMarkhorSecret(cd, ms, slog.String("eid", "_"))
	if err == nil {
		t.Fatal("The decryption should have failed since we did not provide the right key", err)
	}
}

func TestDecryptMarkhorSecret_Fails_With_Wrong_Order(t *testing.T) {
	const SOPS_KEY = "SOPS_AGE_KEY"
	prevKey := os.Getenv(SOPS_KEY)
	defer os.Setenv(SOPS_KEY, prevKey)
	os.Setenv(SOPS_KEY, testSk)
	cd := make(map[string]interface{})
	// Simulating user error in specifying the order in the markhor params
	alteredData := strings.Replace(cypheredData, "apiVersion\",\n\t\t\t\t\t\"kind", "kind\",\n\t\t\t\t\t\"apiVersion", 1)
	err := json.Unmarshal([]byte(alteredData), &cd)
	if err != nil {
		t.Fatal("Failed to unmarshal encrypted JSON", err)
	}

	_, err = DecryptMarkhorSecret(cd, ms, slog.String("eid", "_"))
	if err == nil {
		t.Fatal("The decryption should have failed since the order of the fields was not correct")
	}
}

func TestDecryptMarkhorSecret_Fails_With_Altered_File(t *testing.T) {
	const SOPS_KEY = "SOPS_AGE_KEY"
	prevKey := os.Getenv(SOPS_KEY)
	defer os.Setenv(SOPS_KEY, prevKey)
	os.Setenv(SOPS_KEY, testSk)
	cd := make(map[string]interface{})
	// Simulating an error / attacker having altered the markhor secret after the encryption
	alteredData := strings.Replace(cypheredData, `"namespace": "default"`, `"namespace": "production"`, 1)
	err := json.Unmarshal([]byte(alteredData), &cd)
	if err != nil {
		t.Fatal("Failed to unmarshal encrypted JSON", err)
	}

	_, err = DecryptMarkhorSecret(cd, ms, slog.String("eid", "_"))
	if err == nil {
		t.Fatal("The decryption should have failed since the file contents were altered")
	}
}

func TestDecryptMarkhorSecretEvent_Fails_With_Invalid_Json(t *testing.T) {
	invalidInputs := []struct {
		Reason string
		Json   string
	}{
		{
			Reason: "Empty JSON",
			Json:   "",
		},
		{
			Reason: "Invalid JSON",
			Json: `{
				notQuoted: "example"
			}`,
		},
		{
			Reason: "Invalid JSON",
			Json: `{
				notQuotedProperly": "example"
			}`,
		},
	}
	for i, input := range invalidInputs {
		t.Run(input.Reason, func(t *testing.T) {
			const SOPS_KEY = "SOPS_AGE_KEY"
			prevKey := os.Getenv(SOPS_KEY)
			defer os.Setenv(SOPS_KEY, prevKey)
			os.Setenv(SOPS_KEY, testSk)

			m := v1.MarkhorSecret{ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					LASTAPPLIEDANNOTAION_K8s: input.Json,
				},
			}}
			_, err := DecryptMarkhorSecretEvent(&m, ms, slog.String("eid", "_"))
			if err == nil {
				t.Fatal("The decryption should have failed because of the invalid JSON input number ", i, input.Json)
			}
		})
	}
}

func TestDecryptMarkhorSecretEvent_Works_With_Correct_Input(t *testing.T) {
	// Yes, this key is used only in these tests
	const SOPS_KEY = "SOPS_AGE_KEY"
	prevKey := os.Getenv(SOPS_KEY)
	defer os.Setenv(SOPS_KEY, prevKey)
	os.Setenv(SOPS_KEY, testSk)

	cd := make(map[string]interface{})
	err := json.Unmarshal([]byte(cypheredData), &cd)
	if err != nil {
		t.Fatal("Failed to unmarshal encrypted JSON", err)
	}

	m := v1.MarkhorSecret{ObjectMeta: metav1.ObjectMeta{
		Annotations: map[string]string{
			LASTAPPLIEDANNOTAION_K8s: cypheredData,
		},
	}}
	res, err := DecryptMarkhorSecretEvent(&m, ms, slog.String("eid", "_"))
	if err != nil {
		t.Fatal("Failed to decrypt markhor secret:", err)
	}
	expectedObj := make(map[string]interface{})
	err = json.Unmarshal([]byte(expectedJSON), &expectedObj)
	if err != nil {
		t.Fatal("Failed to unmarshal expected JSON", err)
	}
	diff := cmp.Diff(expectedObj, res)
	if diff != "" {
		t.Fatal("The decryption did not yield the expected result", diff)
	}
}
