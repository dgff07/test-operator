/*
Copyright 2022.

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

	//v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demov1alpha1 "github.com/dgff07/test-operator/api/v1alpha1"
)

// TestReconciler reconciles a Test object
type TestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demo.com.example,resources=tests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.com.example,resources=tests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.com.example,resources=tests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Test object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Test Reconcile method...")
	test := demov1alpha1.Test{}
	err := r.Get(ctx, req.NamespacedName, &test)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Test resource not found. Probably it was deleted after start the reconcile")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Test crd")
		return ctrl.Result{}, err
	}
	logger.Info("Operator working well")

	// Create namespace with the same name

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      test.Name,
			Namespace: test.Name,
		},
	}

	r.Create(ctx, namespace)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1alpha1.Test{}).
		//Owns(&v1.Namespace{}).
		Complete(r)
}
