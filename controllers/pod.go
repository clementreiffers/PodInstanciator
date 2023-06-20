package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"
)

func createPod(instance *apiv1alpha1.PodInstanciator) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-pod",
			Namespace: instance.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  instance.Name + "-pod",
					Image: instance.Spec.ImageName,
					Ports: make([]corev1.ContainerPort, 0, len(instance.Spec.Ports)),
				},
			},
		},
	}

	for _, port := range instance.Spec.Ports {
		containerPort := corev1.ContainerPort{
			Name:          port.PortName,
			ContainerPort: port.PortNumber,
		}
		pod.Spec.Containers[0].Ports = append(pod.Spec.Containers[0].Ports, containerPort)
	}

	return pod
}
