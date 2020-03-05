package v1alpha1

import (
	"github.com/knative/pkg/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KnativeEventingKafkaChannelSpec defines the desired state of KnativeEventingKafkaChannel
// +k8s:openapi-gen=true
type KnativeEventingKafkaChannelSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	BootstrapServers string `json:"bootstrapServers"`

	// +optional
	SetAsDefaultChannelProvisioner bool `json:"setAsDefaultChannelProvisioner,omitempty"`
}

// KnativeEventingKafkaChannelStatus defines the observed state of KnativeEventingKafkaChannel
// +k8s:openapi-gen=true
type KnativeEventingKafkaChannelStatus struct {
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

// KnativeEventingKafkaChannel is the Schema for the knativeeventingkafkachannels API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=="Ready")].status"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=="Ready")].reason"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type KnativeEventingKafkaChannel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KnativeEventingKafkaChannelSpec   `json:"spec,omitempty"`
	Status KnativeEventingKafkaChannelStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KnativeEventingKafkaChannelList contains a list of KnativeEventingKafkaChannel
type KnativeEventingKafkaChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KnativeEventingKafkaChannel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KnativeEventingKafkaChannel{}, &KnativeEventingKafkaChannelList{})
}

// check if KnativeEventingKafkaChannelStatus implements apis.ConditionsAccessor
var _ apis.ConditionsAccessor = &KnativeEventingKafkaChannelStatus{}

// GetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaChannelStatus) GetConditions() apis.Conditions {
	return s.Conditions
}

// SetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaChannelStatus) SetConditions(c apis.Conditions) {
	s.Conditions = c
}
