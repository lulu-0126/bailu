/*
Copyright 2023.

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

type NodeActionType string

// NodeManagementSpec defines the desired state of NodeManagement
type NodeManagementSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NodeManagement. Edit nodemanagement_types.go to remove/update
	NodeAction  NodeActionType `json:"nodeAction,omitempty"`
	NodeReason  string         `json:"nodeReason,omitempty"`
	Description string         `json:"description"`
}

type NodeStatusType string

const (
	NodeStatusUnknown NodeStatusType = "unknown"
)

// NodeManagementStatus defines the observed state of NodeManagement
type NodeManagementStatus struct {
	NodeStatus NodeStatusType `json:"nodeStatus,omitempty"`
	Token     string `json:"token"`
	LastError string `json:"lastError,omitempty"`
	LastUpdate metav1.Time `json:"lastUpdate"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NodeVersion string `json:"nodeVersion"`
	NodeKind    string `json:"nodeKind"`
	NodeIP      string `json:"nodeIP"`
	NodeArc     string `json:"nodeArc"`
}

// +genclient
// +genclient:nonNamespaced
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=nodemm,scope=Cluster
//+kubebuilder:printcolumn:name="Node Status",type=string,JSONPath=`.status.nodeStatus`
//+kubebuilder:printcolumn:name="Node IP",type=string,JSONPath=`.status.nodeIP`
//+kubebuilder:printcolumn:name="Node Kind",type=string,JSONPath=`.status.nodeKind`
//+kubebuilder:printcolumn:name="Node Version",type=string,JSONPath=`.status.nodeVersion`

// NodeManagement is the Schema for the nodemanagements API
type NodeManagement struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeManagementSpec   `json:"spec,omitempty"`
	Status NodeManagementStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeManagementList contains a list of NodeManagement
type NodeManagementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeManagement `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeManagement{}, &NodeManagementList{})
}
