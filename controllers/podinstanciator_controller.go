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

func applyResource(r *PodInstanciatorReconciler, ctx context.Context, resource client.Object, foundResource client.Object) error {
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

	ns := createNamespace(instance)
	pod := createPod(instance)
	svc := createService(instance)
	ingress := createIngress(instance)

	if err := controllerutil.SetControllerReference(instance, ns, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(instance, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	if err := controllerutil.SetControllerReference(instance, ingress, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	err = applyResource(r, ctx, ns, &corev1.Namespace{})
	if err != nil {
		logger.Error(err, "unable to create Namespace")
		return ctrl.Result{}, err
	}

	pod.SetNamespace(instance.Spec.Namespace)
	err = applyResource(r, ctx, pod, &corev1.Pod{})
	if err != nil {
		logger.Error(err, "unable to create Pod")
		return ctrl.Result{}, err
	}

	svc.SetNamespace(instance.Spec.Namespace)
	err = applyResource(r, ctx, svc, &corev1.Service{})
	if err != nil {
		logger.Error(err, "unable to create Service")
		return ctrl.Result{}, err
	}

	ingress.SetNamespace(instance.Spec.Namespace)
	err = applyResource(r, ctx, ingress, &networkingv1.Ingress{})
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
