package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"
)

func createPodPorts(instance *apiv1alpha1.PodInstanciator) []corev1.ContainerPort {
	var paths []corev1.ContainerPort
	for _, port := range instance.Spec.Ports {
		containerPort := corev1.ContainerPort{
			Name:          port.PortName,
			ContainerPort: port.PortNumber,
		}
		paths = append(paths, containerPort)
	}
	return paths
}

func createPod(instance *apiv1alpha1.PodInstanciator) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-pod",
			Namespace: instance.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  instance.Name + "-pod",
					Image: instance.Spec.ImageName,
					Ports: createPodPorts(instance),
				},
			},
		},
	}
}
