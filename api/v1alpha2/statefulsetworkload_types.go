/*


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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
)

// A SecretKeySelector is a reference to a secret key in an arbitrary namespace.
type SecretKeySelector struct {
	// The name of the secret.
	Name string `json:"name"`

	// The key to select.
	Key string `json:"key"`
}

// A ContainerEnvVar specified an environment variable that should be set within
// a container.
type ContainerEnvVar struct {
	// Name of the environment variable. Must be composed of valid Unicode
	// letter and number characters, as well as _ and -.
	// +kubebuilder:validation:Pattern=^[-_a-zA-Z0-9]+$
	Name string `json:"name"`

	// Value of the environment variable.
	// +optional
	Value *string `json:"value,omitempty"`

	// FromSecret is a secret key reference which can be used to assign a value
	// to the environment variable.
	// +optional
	FromSecret *SecretKeySelector `json:"fromSecret,omitempty"`
}

// A ContainerConfigFile specifies a configuration file that should be written
// within a container.
type ContainerConfigFile struct {
	// Path within the container at which the configuration file should be
	// written.
	Path string `json:"path"`

	// Value that should be written to the configuration file.
	// +optional
	Value *string `json:"value,omitempty"`

	// FromSecret is a secret key reference which can be used to assign a value
	// to be written to the configuration file at the given path in the
	// container.
	// +optional
	FromSecret *SecretKeySelector `json:"fromSecret,omitempty"`
}

// A TransportProtocol represents a transport layer protocol.
type TransportProtocol string

// Transport protocols.
const (
	TransportProtocolTCP TransportProtocol = "TCP"
	TransportProtocolUDP TransportProtocol = "UDP"
)

// A ContainerPort specifies a port that is exposed by a container.
type ContainerPort struct {
	// Name of this port. Must be unique within its container. Must be lowercase
	// alphabetical characters.
	// +kubebuilder:validation:Pattern=^[a-z]+$
	Name string `json:"name"`

	// Port number. Must be unique within its container.
	Port int32 `json:"containerPort"`

	// Protocol used by the server listening on this port.
	// +kubebuilder:validation:Enum=TCP;UDP
	// +optional
	Protocol *TransportProtocol `json:"protocol,omitempty"`
}

//A Container represents an Open Containers Initiative (OCI) container.
type Container struct {
	//Name of this container. Must be unique within its workload
	Name string `json:"name"`

	// Image this container should run. Must be a path-like or URI-like
	// representation of an OCI image. May be prefixed with a registry address
	// and should be suffixed with a tag.
	Image string `json:"image"`

	// Command to be run by this container.
	// +optional
	Command []string `json:"command,omitempty"`

	// Arguments to be passed to the command run by this container.
	// +optional
	Arguments []string `json:"arguments,omitempty"`

	// Environment variables that should be set within this container.
	// +optional
	Environment []ContainerEnvVar `json:"env,omitempty"`

	// ConfigFiles that should be written within this container.
	// +optional
	ConfigFiles []ContainerConfigFile `json:"config,omitempty"`

	// Ports exposed by this container.
	// +optional
	Ports []ContainerPort `json:"ports,omitempty"`
}

// StatefulSetWorkloadSpec defines the desired state of StatefulSetWorkload
type StatefulSetWorkloadSpec struct {
	//Containers of which this workload consists.
	Containers []Container `json:"containers"`
}

// StatefulSetWorkloadStatus defines the observed state of StatefulSetWorkload
type StatefulSetWorkloadStatus struct {
	runtimev1alpha1.ConditionedStatus `json:",inline"`

	//Resources managed by this containerised worload
	Resources []runtimev1alpha1.TypedReference `json:"resources,omitempty"`
}

// +kubebuilder:object:root=true

// StatefulSetWorkload is the Schema for the statefulsetworkloads API
// +kubebuilder:resource:categories={crossplane,oam}
// +kubebuilder:subresource:status
type StatefulSetWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StatefulSetWorkloadSpec   `json:"spec,omitempty"`
	Status StatefulSetWorkloadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StatefulSetWorkloadList contains a list of StatefulSetWorkload
type StatefulSetWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatefulSetWorkload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatefulSetWorkload{}, &StatefulSetWorkloadList{})
}
