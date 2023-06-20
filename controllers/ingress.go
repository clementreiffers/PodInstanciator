package controllers

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"
)

func createIngressPaths(instance *apiv1alpha1.PodInstanciator) []networkingv1.HTTPIngressPath {
	paths := make([]networkingv1.HTTPIngressPath, len(instance.Spec.Ports))
	pathType := networkingv1.PathTypePrefix
	for i, port := range instance.Spec.Ports {
		paths[i] = networkingv1.HTTPIngressPath{
			Path:     getIngressPathName(port),
			PathType: &pathType,
			Backend: networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: getServiceName(instance),
					Port: networkingv1.ServiceBackendPort{
						Number: port.PortNumber,
					},
				},
			},
		}
	}
	return paths
}

func createIngress(instance *apiv1alpha1.PodInstanciator) *networkingv1.Ingress {
	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getIngressName(instance),
			Namespace: instance.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "worker.127.0.0.1.sslip.io",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: createIngressPaths(instance),
						},
					},
				},
			},
		},
	}
}
