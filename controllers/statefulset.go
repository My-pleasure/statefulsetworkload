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

const labelKey = "statefulsetworkload.oam.crossplane.io"

var (
	statefulsetKind       = reflect.TypeOf(appsv1.StatefulSet{}).Name()
	statefulsetAPIVersion = appsv1.SchemeGroupVersion.String()
)

// Translator translates a StatefulSetWorkload into a StatefulSet.
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
			Name:      ssw.GetName(),
			Namespace: ssw.GetNamespace(),
		},
		Spec: appsv1.StatefulSetSpec{
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
	// NOTE: If you don't specify ServiceName it will be the defult
	if ssw.Spec.ServiceName != "" {
		s.Spec.ServiceName = ssw.Spec.ServiceName
	} else {
		s.Spec.ServiceName = ssw.GetName()
	}

	if len(ssw.Spec.Template.Spec.ImagePullSecrets) != 0 {
		s.Spec.Template.Spec.ImagePullSecrets = ssw.Spec.Template.Spec.ImagePullSecrets
	}

	for _, container := range ssw.Spec.Template.Spec.Containers {
		kubernetesContainer := corev1.Container{
			Name:           container.Name,
			Image:          container.Image,
			Command:        container.Command,
			Args:           container.Args,
			EnvFrom:        container.EnvFrom,
			Env:            container.Env,
			Resources:      container.Resources,
			VolumeMounts:   container.VolumeMounts,
			LivenessProbe:  container.LivenessProbe,
			ReadinessProbe: container.ReadinessProbe,
		}

		for _, p := range container.Ports {
			port := corev1.ContainerPort{
				Name:          p.Name,
				ContainerPort: p.ContainerPort,
				Protocol:      p.Protocol,
			}
			kubernetesContainer.Ports = append(kubernetesContainer.Ports, port)
		}

		s.Spec.Template.Spec.Containers = append(s.Spec.Template.Spec.Containers, kubernetesContainer)
	}

	return []oam.Object{s}, nil
}
