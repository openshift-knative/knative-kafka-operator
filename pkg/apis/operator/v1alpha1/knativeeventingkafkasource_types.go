package v1alpha1

import (
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
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KnativeEventingKafkaSource is the Schema for the knativeeventingkafkasources API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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
