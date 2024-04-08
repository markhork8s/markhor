package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate controller-gen object paths=$GOFILE

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MarkhorSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// ApiVersion string         `json:"apiVersion"`
	// Sops       SopsSecretSops `json:"sops"`
}

// type SopsSecretSops struct {
// 	Kms []string `json:"kms"`
// }

// // It seems that code generation could not figure this one out
// // because of the arrays
// func (in *SopsSecretSops) DeepCopyInto(out *SopsSecretSops) {
// 	*out = *in
// 	out.Kms = make([]string, len(in.Kms))
// 	copy(out.Kms, in.Kms)
// }

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MarkhorSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []MarkhorSecret `json:"items"`
}
