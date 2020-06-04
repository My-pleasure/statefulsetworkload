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
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// StatefulSetWorkloadSpec defines the desired state of StatefulSetWorkload
type StatefulSetWorkloadSpec struct {
	// NOTE: You can add extension of StatefulSetWorkloadSpec in the future
	// K8S native statefulsetspec
	Template v1.PodTemplateSpec `json:"template,omitempty"`

	ServiceName string `json:"serviceName,omitempty"`
}

// StatefulSetWorkloadStatus defines the observed state of StatefulSetWorkload
type StatefulSetWorkloadStatus struct {
	runtimev1alpha1.ConditionedStatus `json:",inline"`

	// Resources managed by this containerised worload
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
