package controllers

import (
	"context"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"

	"github.com/crossplane/oam-kubernetes-runtime/pkg/oam"
	"statefulsetworkload/api/v1alpha2"
)

// Reconcile error strings.
const (
	errNotStatefulSetWorkload = "object is not a statefulset workload"
)

const defaultNamespace = "default"

const labelKey = "statefulsetworkload.oam.crossplane.io"

var (
	statefulsetKind       = reflect.TypeOf(appsv1.StatefulSet{}).Name()
	statefulsetAPIVersion = appsv1.SchemeGroupVersion.String()
)

// Translator translates a StatefulSetWorkload into a StatefulSet.
// nolint:gocyclo
func Translator(ctx context.Context, w oam.Workload) ([]oam.Object, error) {
	ssw, ok := w.(*v1alpha2.StatefulSetWorkload)
	if !ok {
		return nil, errors.New(errNotStatefulSetWorkload)
	}

	s := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       statefulsetKind,
			APIVersion: statefulsetAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ssw.GetName(),
			// NOTE(hasheddan): we always create the Deployment in the default
			// namespace because there is not currently a namespace scheduling
			// mechanism in the Crossplane OAM implementation. It is likely that
			// this will be addressed in the future by adding a Scope.
			Namespace: defaultNamespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: ssw.GetName(),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					labelKey: string(ssw.GetUID()),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						labelKey: string(ssw.GetUID()),
					},
				},
			},
		},
	}

	for _, container := range ssw.Spec.Containers {
		kubernetesContainer := corev1.Container{
			Name:    container.Name,
			Image:   container.Image,
			Command: container.Command,
			Args:    container.Arguments,
		}

		for _, p := range container.Ports {
			port := corev1.ContainerPort{
				Name:          p.Name,
				ContainerPort: p.Port,
			}
			if p.Protocol != nil {
				port.Protocol = corev1.Protocol(*p.Protocol)
			}
			kubernetesContainer.Ports = append(kubernetesContainer.Ports, port)
		}

		for _, e := range container.Environment {
			if e.Value != nil {
				kubernetesContainer.Env = append(kubernetesContainer.Env, corev1.EnvVar{
					Name:  e.Name,
					Value: *e.Value,
				})
				continue
			}
			if e.FromSecret != nil {
				kubernetesContainer.Env = append(kubernetesContainer.Env, corev1.EnvVar{
					Name: e.Name,
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							Key: e.FromSecret.Key,
							LocalObjectReference: corev1.LocalObjectReference{
								Name: e.FromSecret.Name,
							},
						},
					},
				})
			}
		}

		s.Spec.Template.Spec.Containers = append(s.Spec.Template.Spec.Containers, kubernetesContainer)
	}

	return []oam.Object{s}, nil
}
