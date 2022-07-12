/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConferenceSpec defines the desired state of Conference
type ConferenceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Conference. Edit conference_types.go to remove/update
	//+optional
	ProductionTestEnabled bool `json:"production-test-enabled,omitempty"`
}

// ConferenceStatus defines the observed state of Conference
type ConferenceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+optional
	FrontendReady bool `json:"frontend-ready"`
	//+optional
	AgendaServiceReady bool `json:"agenda-service-ready"`
	//+optional
	C4pServiceReady bool `json:"c4p-service-ready"`
	//+optional
	EmailServiceReady bool `json:"email-service-ready"`

	ProdTests bool `json:"prod-tests"`

	Ready bool   `json:"ready"`
	URL   string `json:"url"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="FRONTEND",type="string",JSONPath=".status.frontend-ready"
// +kubebuilder:printcolumn:name="AGENDA",type="string",JSONPath=".status.agenda-service-ready"
// +kubebuilder:printcolumn:name="EMAIL",type="string",JSONPath=".status.email-service-ready"
// +kubebuilder:printcolumn:name="C4P",type="string",JSONPath=".status.c4p-service-ready"
// +kubebuilder:printcolumn:name="PROD TESTS",type="string",JSONPath=".status.prod-tests"
// +kubebuilder:printcolumn:name="READY",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.url"
// Conference is the Schema for the conferences API
type Conference struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConferenceSpec   `json:"spec,omitempty"`
	Status ConferenceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConferenceList contains a list of Conference
type ConferenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Conference `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Conference{}, &ConferenceList{})
}
