/*
Copyright 2024.

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

type NewPersistentVolumeClaim struct {
	OriginalPVCName   string `json:"originalpvcname"`
	OriginalNamespace string `json:"originalnamespace"`
	TargetPVCName     string `json:"targetpvcname,omitempty"`
	TargetNamespace   string `json:"targetnamespace"`
}

// SharedPVCSpec defines the desired state of SharedPVC
type SharedPVCSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SharedPVC. Edit sharedpvc_types.go to remove/update
	Image  string                   `json:"image"`
	NewPVC NewPersistentVolumeClaim `json:"newpvc"`
}

// SharedPVCStatus defines the observed state of SharedPVC
type SharedPVCStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Status         string `json:"status"`
	OriginalExists bool   `json:"originalexists"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SharedPVC is the Schema for the sharedpvcs API
type SharedPVC struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SharedPVCSpec   `json:"spec,omitempty"`
	Status SharedPVCStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SharedPVCList contains a list of SharedPVC
type SharedPVCList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SharedPVC `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SharedPVC{}, &SharedPVCList{})
}
