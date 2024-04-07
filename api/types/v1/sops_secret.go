package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate controller-gen object paths=$GOFILE

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SopsSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	ApiVersion string         `json:"apiVersion"`
	Kind       string         `json:"kind"`
	Type       string         `json:"type"`
	Sops       SopsSecretSops `json:"sops"`
}

type SopsSecretSops struct {
	Kms            []string            `json:"kms"`
	GcpKms         []string            `json:"gcp_kms"`
	AzureKv        []string            `json:"azure_kv"`
	HcVault        []string            `json:"hc_vault"`
	Pgp            []string            `json:"pgp"`
	Age            []SopsSecretSopsAge `json:"age"`
	LastModified   string              `json:"lastmodified"`
	Mac            string              `json:"mac"`
	EncryptedRegex string              `json:"encrypted_regex"`
	Version        string              `json:"version"`
}

// It seems that code generation could not figure this one out
// because of the arrays
func (in *SopsSecretSops) DeepCopyInto(out *SopsSecretSops) {
	*out = *in
	out.Kms = make([]string, len(in.Kms))
	copy(out.Kms, in.Kms)
	out.GcpKms = make([]string, len(in.GcpKms))
	copy(out.GcpKms, in.GcpKms)
	out.AzureKv = make([]string, len(in.AzureKv))
	copy(out.AzureKv, in.AzureKv)
	out.HcVault = make([]string, len(in.HcVault))
	copy(out.HcVault, in.HcVault)
	out.Pgp = make([]string, len(in.Pgp))
	copy(out.Pgp, in.Pgp)
	out.Age = make([]SopsSecretSopsAge, len(in.Age))
	copy(out.Age, in.Age)
}

type SopsSecretSopsAge struct {
	Recipient string `json:"recipient"`
	Enc       string `json:"enc"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SopsSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SopsSecret `json:"items"`
}
