package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"
)

func createService(instance *apiv1alpha1.PodInstanciator) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getServiceName(instance),
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports:     []corev1.ServicePort{},
			Selector:  map[string]string{"app": getPodName(instance)},
			ClusterIP: "None",
		},
	}
}
