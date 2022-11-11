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
	"encoding/json"
	"strings"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=namespaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Namespace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	namespace := v1.Namespace{}
	err := r.Get(ctx, req.NamespacedName, &namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Namespace resource not found. Probably it was deleted after start the reconcile")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Namespace crd")
		return ctrl.Result{}, err
	}

	isNamespaceToBeDeleted := namespace.GetDeletionTimestamp() != nil
	if isNamespaceToBeDeleted {
		return ctrl.Result{}, nil
	}

	lastAppliedConfig := namespace.Annotations[v1.LastAppliedConfigAnnotation]
	if len(strings.TrimSpace(lastAppliedConfig)) == 0 {
		return ctrl.Result{}, nil
	}

	var lastNamespace v1.Namespace
	err = json.Unmarshal([]byte(lastAppliedConfig), &lastNamespace)
	if err != nil {
		logger.Error(err, "Failed to convert json to namespace struct")
		return ctrl.Result{}, err
	}

	printLabelsAdded(namespace.Labels, lastNamespace.Labels, &logger)
	printLabelsRemovedOrUpdated(namespace.Labels, lastNamespace.Labels, &logger)

	bytes, err := json.Marshal(namespace)
	if err != nil {
		return ctrl.Result{}, err
	}
	namespace.Annotations[v1.LastAppliedConfigAnnotation] = string(bytes)
	r.Update(ctx, &namespace)

	return ctrl.Result{}, nil

}

func printLabelsAdded(currentLabels map[string]string, oldLabels map[string]string, logger *logr.Logger) {

	for key, value := range currentLabels {
		oldLabelValue, keyExist := oldLabels[key]
		if value != oldLabelValue && !keyExist {
			logger.Info("NEW LABEL", "key", key, "value", value)
		}
	}
}

func printLabelsRemovedOrUpdated(currentLabels map[string]string, oldLabels map[string]string, logger *logr.Logger) {

	for key, oldValue := range oldLabels {
		currentLabelValue, keyExist := currentLabels[key]
		if oldValue != currentLabelValue {
			if keyExist {
				logger.Info("UPDATED LABEL", "key", key, "value", currentLabelValue, "oldValue", oldValue)
				continue
			}
			logger.Info("REMOVED LABEL", "key", key, "value", oldValue)
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Namespace{}).
		Complete(r)
}
