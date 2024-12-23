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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EnvironmentVariables struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesNMSMagmaLte struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesOrc8rNotifier struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesOrc8rNginx struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type EnvironmentVariablesDirectoryD struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ImageFluentd struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

type Image struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

type Persistent struct {
	PvcClaimName     string `json:"pvcClaimName,omitempty"`
	StorageClassName string `json:"storageClassName,omitempty"`
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PmnsystemSpec defines the desired state of Pmnsystem
type PmnsystemSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Pmnsystem. Edit pmnsystem_types.go to remove/update
	ReplicaCount              int32                               `json:"replicaCount,omitempty"`
	NginxImage                string                              `json:"nginxImage,omitempty"`
	NotifierImage                string                              `json:"notifierImage,omitempty"`
	PullPolicy                string                              `json:"pullPolicy,omitempty"`
	Persistent                Persistent                          `json:"persistent,omitempty"`
	NameSpace                 string                              `json:"nameSpace,omitempty"`
	Image                     Image                               `json:"image,omitempty"`
	ImageFluentd              ImageFluentd                         `json:"imageFluentd,omitempty"`
	EnvVariables              []EnvironmentVariables              `json:"envVariables,omitempty"`
	EnvVariablesDirectoryD    []EnvironmentVariablesDirectoryD    `json:"envVariablesDirectoryD,omitempty"`
	EnvVariablesOrc8rNginx    []EnvironmentVariablesOrc8rNginx    `json:"envVariablesOrc8rNginx,omitempty"`
	EnvVariablesOrc8rNotifier []EnvironmentVariablesOrc8rNotifier `json:"envVariablesOrc8rNotifier,omitempty"`
	EnvVariablesNMSMagmaLte   []EnvironmentVariablesNMSMagmaLte   `json:"envVariablesNMSMagmaLte,omitempty"`
	ImagePullSecrets          string                              `json:"imagePullSecrets,omitempty"`
}

// PmnsystemStatus defines the observed state of Pmnsystem
type PmnsystemStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pmnsystem is the Schema for the pmnsystems API
type Pmnsystem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PmnsystemSpec   `json:"spec,omitempty"`
	Status PmnsystemStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PmnsystemList contains a list of Pmnsystem
type PmnsystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pmnsystem `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pmnsystem{}, &PmnsystemList{})
}
