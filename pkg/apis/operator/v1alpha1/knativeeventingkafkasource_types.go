package v1alpha1

import (
	"github.com/knative/pkg/apis"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KnativeEventingKafkaSourceSpec defines the desired state of KnativeEventingKafkaSource
// +k8s:openapi-gen=true
type KnativeEventingKafkaSourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// KnativeEventingKafkaSourceStatus defines the observed state of KnativeEventingKafkaSource
// +k8s:openapi-gen=true
type KnativeEventingKafkaSourceStatus struct {
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

// KnativeEventingKafkaSource is the Schema for the knativeeventingkafkasources API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=="Ready")].status"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=="Ready")].reason"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type KnativeEventingKafkaSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KnativeEventingKafkaSourceSpec   `json:"spec,omitempty"`
	Status KnativeEventingKafkaSourceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KnativeEventingKafkaSourceList contains a list of KnativeEventingKafkaSource
type KnativeEventingKafkaSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KnativeEventingKafkaSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KnativeEventingKafkaSource{}, &KnativeEventingKafkaSourceList{})
}

// check if KnativeEventingKafkaSourceStatus implements apis.ConditionsAccessor
var _ apis.ConditionsAccessor = &KnativeEventingKafkaSourceStatus{}

// GetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaSourceStatus) GetConditions() apis.Conditions {
	return s.Conditions
}

// SetConditions implements apis.ConditionsAccessor
func (s *KnativeEventingKafkaSourceStatus) SetConditions(c apis.Conditions) {
	s.Conditions = c
}
