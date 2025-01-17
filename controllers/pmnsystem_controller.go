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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	// "github.com/gogo/protobuf/test/required"

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
	if pmnsystem.Spec.DevEnvironment {
		r.Log.Info("DevEnvironment is true. Proceeding with PostgreSQL Deployment...")

		// Generate the deployment
		postgresDeployment := r.deploymentPostgres(pmnsystem)

		// Ensure Deployment is created only if not nil
		if postgresDeployment != nil {
			result, err = r.ensureDeployment(req, pmnsystem, postgresDeployment)
			if result != nil {
				return *result, err
			}
		} else {
			r.Log.Info("Skipping ensureDeployment as deploymentPostgres returned nil")
		}
	} else {
		r.Log.Info("DevEnvironment is false. Skipping PostgreSQL Deployment.")
	}
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
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rDispatcherDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rEventdDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rmetricsdDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rNginxDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rNotifierDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rObsidianDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8WorkerDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8orchestratorDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8ServiceRegistryDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8StateDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8StreamerDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8TenantsDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rHaDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8LteDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8NprobeDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8PolicyDbDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8SmsdDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8SubscriberDbCacheDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8SubscriberDbDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.nmsMagmaLteDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8AlertManagerDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8PrometheusCacheDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rPrometheusConfigurerDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rPrometheusKafkaAdapterDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rPrometheusNginxProxyDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8rUserGrafanaDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensureDeployment(req, pmnsystem, r.orc8AlertManagerConfigurerDeployment(pmnsystem))
	if result != nil {
		return *result, err
	}
	// result, err = r.ensureDeployment(req, pmnsystem, r.orc8rDomainProxyDeployment(pmnsystem))
	// if result != nil {
	// 	return *result, err
	// }

	// =================ensure StatefulSets=================
	result, err = r.ensureStatefulSet(req, pmnsystem, r.createOrc8rPrometheusStateFullSet(pmnsystem))
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
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rDispatcherPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rEventdPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rMetricsdPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rNginxPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rObsidianDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rWorkerDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rOrchestratorDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rServiceRegistryDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rStateDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rStreamerDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rTenantsDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rHaDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rLteDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rNprobeDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rPolicyDbDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rSmSdDbDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rSubscriberDbCachedDbDPDB(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePodDisruptionBudget(req, pmnsystem, r.orc8rSubscriberDbDbDPDB(pmnsystem))
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
	svc = r.orc8rDispatcherService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rEventdService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rmetricsdService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rNotifierService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rNotifierInternalService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rObsidianService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rWorkerService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rOrchestratorService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rServiceRegistryService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rStateService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rStreamerService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rTenantsService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rHaService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rLteService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rNprobeService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rPolicyDbService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rSmsdService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rSubscriberDbCacheService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rSubscriberDbService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.NmsMagmaLteService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rAlterManagerService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rPrometheusCacheService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rPrometheusConfigurerService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rPrometheusKafkaAdapterService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rPrometheusNginxProxyService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rUserGrafanaService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rPrometheusService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	svc = r.orc8rAlertManagerConfigurerService(pmnsystem)
	result, err = r.ensureService(pmnsystem, svc)
	if result != nil {
		return *result, err
	}
	// Check if DevEnvironment is true
	if pmnsystem.Spec.DevEnvironment {
		r.Log.Info("DevEnvironment is true. Proceeding with PostgreSQL Service...")

		// Create the service
		svc := r.servicePostgres(pmnsystem)
		if svc != nil {
			// Ensure the service exists
			result, err := r.ensureService(pmnsystem, svc)
			if result != nil {
				return *result, err
			}
		} else {
			r.Log.Info("Skipping ensureService as servicePostgres returned nil")
		}
	} else {
		r.Log.Info("DevEnvironment is false. Skipping PostgreSQL Service creation.")
	}
	//Ensure Secrets are created
	secretName := "pmn-certs"
	certDir := "/mnt/c/Users/MSI-1/Desktop/PMN-Operator-Two/pmn-operator/VMFolder"
	requiredFiles := []string{"rootCA.pem", "rootCA.key", "notifier.key", "notifier.crt", "notifier-ca.crt", "controller.key", "controller.crt", "certifier.pem", "certifier.key", "bootstrapper.key", "admin_operator.pem", "admin_operator.key.pem"}
	namespace := "pgsql"
	err = r.CreateSecretsFromCertificates(secretName, certDir, requiredFiles, namespace, pmnsystem)
	if err != nil {
		r.Log.Error(err, "Failed to create secret")
		return ctrl.Result{}, err
	}

	secretName = "prometheus-adapter-certs"
	certDir = "/mnt/c/Users/MSI-1/Desktop/PMN-Operator-Two/pmn-operator/VMFolder"
	requiredFiles = []string{"adapter-ca.crt", "adapter-client.crt", "adapter-client.key"}
	namespace = "pgsql"
	err = r.CreateSecretsFromCertificates(secretName, certDir, requiredFiles, namespace, pmnsystem)
	if err != nil {
		r.Log.Error(err, "Failed to create secret")
		return ctrl.Result{}, err
	}
	secretName = "prometheus-certs"
	certDir = "/mnt/c/Users/MSI-1/Desktop/PMN-Operator-Two/pmn-operator/VMFolder"
	requiredFiles = []string{"prometheus-ca.crt", "prometheus.crt", "prometheus.key"}
	namespace = "pgsql"
	err = r.CreateSecretsFromCertificates(secretName, certDir, requiredFiles, namespace, pmnsystem)
	if err != nil {
		r.Log.Error(err, "Failed to create secret")
		return ctrl.Result{}, err
	}
	secretName = "nms-certs"
	certDir = "/mnt/c/Users/MSI-1/Desktop/PMN-Operator-Two/pmn-operator/VMFolder"
	requiredFiles = []string{"admin_operator.key.pem", "admin_operator.pem", "controller.crt", "controller.key"}
	namespace = "pgsql"
	err = r.CreateSecretsFromCertificates(secretName, certDir, requiredFiles, namespace, pmnsystem)
	if err != nil {
		r.Log.Error(err, "Failed to create secret")
		return ctrl.Result{}, err
	}
	// ====ensure PVC====
	result, err = r.ensurePersistentVolumeClaim(pmnsystem, r.createPersistentVolumeClaim(pmnsystem))
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
		Owns(&corev1.Service{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
