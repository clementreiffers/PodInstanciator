/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "operators/PodInstanciater/api/v1alpha1"
)

// PodInstanciatorReconciler reconciles a PodInstanciator object
type PodInstanciatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.my.domain,resources=podinstanciators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.my.domain,resources=podinstanciators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.my.domain,resources=podinstanciators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodInstanciator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile

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

func createService(instance *apiv1alpha1.PodInstanciator) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-svc",
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports:     []corev1.ServicePort{},
			Selector:  map[string]string{"app": instance.Name + "-pod"},
			ClusterIP: "None",
		},
	}
}

func createIngress(instance *apiv1alpha1.PodInstanciator) *networkingv1.Ingress {
	pathType := networkingv1.PathTypePrefix
	ingress := &networkingv1.Ingress{
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
							Paths: []networkingv1.HTTPIngressPath{},
						},
					},
				},
			},
		},
	}
	for _, port := range instance.Spec.Ports {
		path := "/" + port.PortName
		ingress.Spec.Rules[0].HTTP.Paths = append(ingress.Spec.Rules[0].HTTP.Paths, networkingv1.HTTPIngressPath{
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
	return ingress
}

func createResource(r *PodInstanciatorReconciler, ctx context.Context, resource client.Object, foundResource client.Object) error {
	err := r.Get(ctx, types.NamespacedName{Name: resource.GetName(), Namespace: resource.GetNamespace()}, foundResource)
	if err != nil && errors.IsNotFound(err) {
		err = r.Create(ctx, resource)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (r *PodInstanciatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.Log.WithValues("PodInstanciator", req.NamespacedName)

	instance := &apiv1alpha1.PodInstanciator{}
	err := r.Get(ctx, req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	pod := createPod(instance)
	svc := createService(instance)
	ingress := createIngress(instance)

	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(instance, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(instance, ingress, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	err = createResource(r, ctx, pod, &corev1.Pod{})
	if err != nil {
		logger.Error(err, "unable to create Pod")
		return ctrl.Result{}, err
	}
	err = createResource(r, ctx, svc, &corev1.Service{})
	if err != nil {
		logger.Error(err, "unable to create Service")
		return ctrl.Result{}, err
	}
	err = createResource(r, ctx, ingress, &networkingv1.Ingress{})
	if err != nil {
		logger.Error(err, "unable to create Ingress")
		return ctrl.Result{}, err
	}

	logger.Info("all resources created!")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodInstanciatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.PodInstanciator{}).
		Complete(r)
}
