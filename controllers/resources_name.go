package controllers

import apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"

func getPodName(instance *apiv1alpha1.PodInstanciator) string {
	return instance.Name + "-pod"
}

func getServiceName(instance *apiv1alpha1.PodInstanciator) string {
	return instance.Name + "-svc"
}

func getIngressName(instance *apiv1alpha1.PodInstanciator) string {
	return instance.Name + "-ingress"
}

func getIngressPathName(port apiv1alpha1.Port) string {
	return "/" + port.PortName
}
