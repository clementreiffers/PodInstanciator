package controllers

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"
)

func createIngressPaths(instance *apiv1alpha1.PodInstanciator) []networkingv1.HTTPIngressPath {
	pathType := networkingv1.PathTypePrefix
	var paths []networkingv1.HTTPIngressPath
	for _, port := range instance.Spec.Ports {
		path := "/" + port.PortName
		paths = append(paths, networkingv1.HTTPIngressPath{
			Path:     path,
			PathType: &pathType,
			Backend: networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: instance.Name + "-svc",
					Port: networkingv1.ServiceBackendPort{
						Number: port.PortNumber,
					},
				},
			},
		})
	}
	return paths
}

func createIngress(instance *apiv1alpha1.PodInstanciator) *networkingv1.Ingress {
	return &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-ingress",
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