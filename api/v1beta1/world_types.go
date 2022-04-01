/*
Copyright 2022 antmoveh.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WorldSpec defines the desired state of World
type WorldSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	World string `json:"world,omitempty"`
}

// WorldStatus defines the observed state of World
type WorldStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	War string `json:"war,omitempty"`
	SyncTime metav1.Time `json:"syncTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=wd
// +kubebuilder:printcolumn:name="world",type="string",JSONPath=".spec.world"
// +kubebuilder:printcolumn:name="war",type="string",JSONPath=".status.war"
// +kubebuilder:printcolumn:name="syncTime",type="date",priority=1,JSONPath=".status.syncTime"

// World is the Schema for the worlds API
type World struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorldSpec   `json:"spec,omitempty"`
	Status WorldStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WorldList contains a list of World
type WorldList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []World `json:"items"`
}

func init() {
	SchemeBuilder.Register(&World{}, &WorldList{})
}
