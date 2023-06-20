package controllers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"operators/PodInstanciater/api/v1alpha1"
)

func createNamespace(instance *v1alpha1.PodInstanciator) *corev1.Namespace {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.Spec.Namespace,
		},
	}
	ns.SetNamespace(instance.GetNamespace())
	return ns
}
