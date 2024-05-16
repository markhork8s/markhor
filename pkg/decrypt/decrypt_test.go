package decrypt

import (
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"

	v1 "github.com/civts/markhor/pkg/api/types/v1"
	"github.com/civts/markhor/pkg/config"
	"github.com/google/go-cmp/cmp"
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
			"another": "ENC[AES256_GCM,data:m4QHjR8pkACVlLSgwoQUPQHeK+DeWA==,iv:OR7XlHKtAqIVQeBFr8d/DJC7isPHaRMzllJrxZ1/xNQ=,tag:hDCqVb3AYG3nmU03DoTJTA==,type:str]"
	},
	"type": "Opaque",
	"sops": {
		"kms": [],
		"gcp_kms": [],
		"pgp": [],
		"azure_kv": [],
		"hc_vault": [],
		"age": [{
							"recipient": "age1apq7ck5adq6dkd0c242phl42fsurvpxvt9pwk0qg7ahdex7fqppqj8pe8y",
							"enc": "-----BEGIN AGE ENCRYPTED FILE-----\nYWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBGQ24xNVIxVFVnb2F3anFC\nUkNQdFRIZDQvMFBKS2NtaXJwUkpVVDJpZlFFCmg0UHFJTnprTUwrNlQzSUE0WTMy\nT25XN1pGeU5RMXhRQmM0dTlabjBqNTQKLS0tIHBjSTUzY1lDWEtvK1BEcDV3a3ZO\nLzlwd0EzdUkyU1I5VCtxNDY5Qk11ZkkKD4o6pG1Gi+a0yfa/vX/QS2QkieuSg80O\nUr2hYHizEMW3KWH/M4UkCJOS93Zs9dYNYgM9potu+EwvhNfNwqBg4g==\n-----END AGE ENCRYPTED FILE-----\n"
					}],
		"lastmodified": "2024-04-27T12:15:51Z",
		"mac": "ENC[AES256_GCM,data:6EQMQNNwBJrNDEJaxOyj8YoOfC04Bi1w56zjqFvpwLTGv3sufuxgu5ctsaGeC/sDUi3UFPqq7PKKVXUcXFrE0HHwHs4QVj2A/Gp3O7sLKMwSfVdpYOWrWZbkgx4HpfW4Si+ooYLCFtzVNCD4/SZh01rYiyFLNzhF9RLGX+fl46Q=,iv:9qYHQXS6d6qsd7+VBp3yv4trFP2Jkj1Ek0H4mFGfAZU=,tag:LFZpTrtsSx1nVfuqCr6JJg==,type:str]",
		"encrypted_regex": "^(data|stringData)$",
		"version": "3.8.1"
	},
	"data": {
		"session_secret": "ENC[AES256_GCM,data:XF88K93P1f4eeFqQUy+WI+mmZtnVx9EE/xb0GiO89Eb/c0JXc6OnCg==,iv:y4WM97K5yOmkEHvfLoiaEjFjuoJ3/qr06EalInEEkiU=,tag:UkP4XsTC4FcjAo5SFxVNrA==,type:str]"
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
