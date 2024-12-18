/*
Copyright 2024.

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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"

	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"
)

// PmnsystemReconciler reconciles a Pmnsystem object
type PmnsystemReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pmnsystems.pmnsystem.com,resources=pmnsystems,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pmnsystems.pmnsystem.com,resources=pmnsystems/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pmnsystems.pmnsystem.com,resources=pmnsystems/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PmnsystemReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("Pmnsystem", req.NamespacedName)

	r.Log.Info("Reconciling Pmnsystem")

	pmnsystem := &v1.Pmnsystem{}

	err := r.Client.Get(context.TODO(), req.NamespacedName, pmnsystem)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	var result *ctrl.Result
	// ====ensure Deployments====
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rAccessD(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rAnalyticsDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rBootStrapperDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rCertifierDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rConfiguratorDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rDeviceDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rDirectorydDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	// ====ensure PodDisruptionBudget====
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rAccessDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rAnalyticsDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rBootstrapperPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rCertifierPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rConfiguratorPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rDevicePDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rDirectorydPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	// ====ensure Service====
	svc := r.orc8rAccessDService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rAnalyticsService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rBootStrapperService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rCertifierService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rConfiguratorService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rDeviceService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rDirectoryDService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PmnsystemReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pmnsystem{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
