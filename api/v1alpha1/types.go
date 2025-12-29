package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Probe is the Schema for the probes API
type Probe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProbeSpec   `json:"spec,omitempty"`
	Status ProbeStatus `json:"status,omitempty"`
}

// ProbeSpec defines the desired state of Probe
type ProbeSpec struct {
	CheckType   string `json:"checkType"`
	CheckTarget string `json:"checkTarget"`
	Interval    string `json:"interval"`
	Timeout     string `json:"timeout,omitempty"`
}

// ProbeStatus defines the observed state of Probe
type ProbeStatus struct {
	Healthy       bool         `json:"healthy"`
	LastProbeTime *metav1.Time `json:"lastProbeTime,omitempty"`
	Message       string       `json:"message,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProbeList contains a list of Probe
type ProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Probe `json:"items"`
}
