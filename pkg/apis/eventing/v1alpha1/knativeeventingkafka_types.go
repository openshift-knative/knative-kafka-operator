package v1alpha1

import (
	"github.com/knative/pkg/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	InstallSucceeded     apis.ConditionType = "InstallSucceeded"
	DeploymentsAvailable apis.ConditionType = "DeploymentsAvailable"
)

// KnativeEventingKafkaSpec defines the desired state of KnativeEventingKafka
// +k8s:openapi-gen=true
type KnativeEventingKafkaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	BootstrapServers               string `json:"bootstrapServers"`
	// +optional
	SetAsDefaultChannelProvisioner bool   `json:"setAsDefaultChannelProvisioner,omitempty"`
}

// KnativeEventingKafkaStatus defines the observed state of KnativeEventingKafka
// +k8s:openapi-gen=true
type KnativeEventingKafkaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// The version of the installed release
	// +optional
	Version string `json:"version"`
	// The latest available observations of a resource's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions apis.Conditions `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KnativeEventingKafka is the Schema for the knativeeventingkafkas API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type KnativeEventingKafka struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KnativeEventingKafkaSpec   `json:"spec,omitempty"`
	Status KnativeEventingKafkaStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KnativeEventingKafkaList contains a list of KnativeEventingKafka
type KnativeEventingKafkaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KnativeEventingKafka `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KnativeEventingKafka{}, &KnativeEventingKafkaList{})
}
