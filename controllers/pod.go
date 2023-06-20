package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"
)

func createPodPorts(instance *apiv1alpha1.PodInstanciator) []corev1.ContainerPort {
	ports := make([]corev1.ContainerPort, len(instance.Spec.Ports))
	for i, port := range instance.Spec.Ports {
		ports[i] = corev1.ContainerPort{
			Name:          port.PortName,
			ContainerPort: port.PortNumber,
		}
	}
	return ports
}

func createPod(instance *apiv1alpha1.PodInstanciator) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getPodName(instance),
			Namespace: instance.Spec.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  getPodName(instance),
					Image: instance.Spec.ImageName,
					Ports: createPodPorts(instance),
				},
			},
		},
	}
	pod.SetNamespace(instance.GetNamespace())
	return pod
}
