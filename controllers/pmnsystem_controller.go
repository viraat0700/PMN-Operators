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
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

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
	var deploymentFunctions = []struct {
		Name string
		Func func(*v1.Pmnsystem) *appsv1.Deployment
	}{
		{"orc8rAccessD", r.orc8rAccessD},
		{"orc8rAnalyticsDeployment", r.orc8rAnalyticsDeployment},
		{"orc8rBootStrapperDeployment", r.orc8rBootStrapperDeployment},
		{"orc8rCertifierDeployment", r.orc8rCertifierDeployment},
		{"orc8rConfiguratorDeployment", r.orc8rConfiguratorDeployment},
		{"orc8rDeviceDeployment", r.orc8rDeviceDeployment},
		{"orc8rDirectorydDeployment", r.orc8rDirectorydDeployment},
		{"orc8rDispatcherDeployment", r.orc8rDispatcherDeployment},
		{"orc8rEventdDeployment", r.orc8rEventdDeployment},
		{"orc8rmetricsdDeployment", r.orc8rmetricsdDeployment},
		{"orc8rNginxDeployment", r.orc8rNginxDeployment},
		{"orc8rNotifierDeployment", r.orc8rNotifierDeployment},
		{"orc8rObsidianDeployment", r.orc8rObsidianDeployment},
		{"orc8WorkerDeployment", r.orc8WorkerDeployment},
		{"orc8orchestratorDeployment", r.orc8orchestratorDeployment},
		{"orc8ServiceRegistryDeployment", r.orc8ServiceRegistryDeployment},
		{"orc8StateDeployment", r.orc8StateDeployment},
		{"orc8StreamerDeployment", r.orc8StreamerDeployment},
		{"orc8TenantsDeployment", r.orc8TenantsDeployment},
		{"orc8rHaDeployment", r.orc8rHaDeployment},
		{"orc8LteDeployment", r.orc8LteDeployment},
		{"orc8NprobeDeployment", r.orc8NprobeDeployment},
		{"orc8PolicyDbDeployment", r.orc8PolicyDbDeployment},
		{"orc8SmsdDeployment", r.orc8SmsdDeployment},
		{"orc8SubscriberDbCacheDeployment", r.orc8SubscriberDbCacheDeployment},
		{"orc8SubscriberDbDeployment", r.orc8SubscriberDbDeployment},
		{"nmsMagmaLteDeployment", r.nmsMagmaLteDeployment},
		{"orc8PrometheusCacheDeployment", r.orc8PrometheusCacheDeployment},
		{"orc8rPrometheusConfigurerDeployment", r.orc8rPrometheusConfigurerDeployment},
		{"orc8rPrometheusKafkaAdapterDeployment", r.orc8rPrometheusKafkaAdapterDeployment},
		{"orc8rPrometheusNginxProxyDeployment", r.orc8rPrometheusNginxProxyDeployment},
		{"orc8rUserGrafanaDeployment", r.orc8rUserGrafanaDeployment},
		{"orc8AlertManagerConfigurerDeployment", r.orc8AlertManagerConfigurerDeployment},
		{"orc8AlertManagerDeployment", r.orc8AlertManagerDeployment},
	}
	// Iterate over deployment functions
	for _, deployment := range deploymentFunctions {
		r.Log.Info("Ensuring deployment", "DeploymentName", deployment.Name)

		// Ensure deployment
		result, err := r.ensureDeployment(req, pmnsystem, deployment.Func(pmnsystem))
		if result != nil || err != nil {
			return *result, err
		}
	}
	// =================ensure StatefulSets=================
	result, err = r.ensureStatefulSet(req, pmnsystem, r.createOrc8rPrometheusStateFullSet(pmnsystem))
	if result != nil {
		return *result, err
	}
	// ====ensure PodDisruptionBudget====
	var pdbFunctions = []struct {
		Name string
		Func func(*v1.Pmnsystem) *policyv1.PodDisruptionBudget
	}{
		{"orc8rAccessDPDB", r.orc8rAccessDPDB},
		{"orc8rAnalyticsPDB", r.orc8rAnalyticsDPDB},
		{"orc8rBootstrapperPDB", r.orc8rBootstrapperPDB},
		{"orc8rCertifierPDB", r.orc8rCertifierPDB},
		{"orc8rConfiguratorPDB", r.orc8rConfiguratorPDB},
		{"orc8rDevicePDB", r.orc8rDevicePDB},
		{"orc8rDirectorydPDB", r.orc8rDirectorydPDB},
		{"orc8rDispatcherPDB", r.orc8rDispatcherPDB},
		{"orc8rEventdPDB", r.orc8rEventdPDB},
		{"orc8rMetricsdPDB", r.orc8rMetricsdPDB},
		{"orc8rNginxPDB", r.orc8rNginxPDB},
		{"orc8rObsidianPDB", r.orc8rObsidianDPDB},
		{"orc8rWorkerPDB", r.orc8rWorkerDPDB},
		{"orc8rOrchestratorPDB", r.orc8rOrchestratorDPDB},
		{"orc8rServiceRegistryPDB", r.orc8rServiceRegistryDPDB},
		{"orc8rStatePDB", r.orc8rStateDPDB},
		{"orc8rStreamerPDB", r.orc8rStreamerDPDB},
		{"orc8rTenantsPDB", r.orc8rTenantsDPDB},
		{"orc8rHaPDB", r.orc8rHaDPDB},
		{"orc8rLtePDB", r.orc8rLteDPDB},
		{"orc8rNprobePDB", r.orc8rNprobeDPDB},
		{"orc8rPolicyDbPDB", r.orc8rPolicyDbDPDB},
		{"orc8rSmSdPDB", r.orc8rSmSdDbDPDB},
		{"orc8rSubscriberDbCachePDB", r.orc8rSubscriberDbCachedDbDPDB},
		{"orc8rSubscriberDbPDB", r.orc8rSubscriberDbDbDPDB},
	}
	// Iterate over PDB functions
	for _, pdb := range pdbFunctions {
		r.Log.Info("Ensuring PodDisruptionBudget", "PDBName", pdb.Name)

		// Ensure PodDisruptionBudget
		result, err := r.ensurePodDisruptionBudget(req, pmnsystem, pdb.Func(pmnsystem))
		if result != nil || err != nil {
			return *result, err
		}
	}
	// ====ensure Service====
	var serviceFunctions = []struct {
		Name string
		Func func(*v1.Pmnsystem) *corev1.Service
	}{
		{"orc8rAccessDService", r.orc8rAccessDService},
		{"orc8rAnalyticsService", r.orc8rAnalyticsService},
		{"orc8rBootStrapperService", r.orc8rBootStrapperService},
		{"orc8rCertifierService", r.orc8rCertifierService},
		{"orc8rConfiguratorService", r.orc8rConfiguratorService},
		{"orc8rDeviceService", r.orc8rDeviceService},
		{"orc8rDirectoryDService", r.orc8rDirectoryDService},
		{"orc8rDispatcherService", r.orc8rDispatcherService},
		{"orc8rEventdService", r.orc8rEventdService},
		{"orc8rmetricsdService", r.orc8rmetricsdService},
		{"orc8rNginxProxyService", r.orc8rNginxProxyService},
		{"orc8rNotifierService", r.orc8rNotifierService},
		{"orc8rNotifierInternalService", r.orc8rNotifierInternalService},
		{"orc8rObsidianService", r.orc8rObsidianService},
		{"orc8rWorkerService", r.orc8rWorkerService},
		{"orc8rOrchestratorService", r.orc8rOrchestratorService},
		{"orc8rServiceRegistryService", r.orc8rServiceRegistryService},
		{"orc8rStateService", r.orc8rStateService},
		{"orc8rStreamerService", r.orc8rStreamerService},
		{"orc8rTenantsService", r.orc8rTenantsService},
		{"orc8rHaService", r.orc8rHaService},
		{"orc8rLteService", r.orc8rLteService},
		{"orc8rNprobeService", r.orc8rNprobeService},
		{"orc8rPolicyDbService", r.orc8rPolicyDbService},
		{"orc8rSmsdService", r.orc8rSmsdService},
		{"orc8rSubscriberDbCacheService", r.orc8rSubscriberDbCacheService},
		{"orc8rSubscriberDbService", r.orc8rSubscriberDbService},
		{"NmsMagmaLteService", r.NmsMagmaLteService},
		{"orc8rAlterManagerService", r.orc8rAlterManagerService},
		{"orc8rPrometheusCacheService", r.orc8rPrometheusCacheService},
		{"orc8rPrometheusConfigurerService", r.orc8rPrometheusConfigurerService},
		{"orc8rPrometheusKafkaAdapterService", r.orc8rPrometheusKafkaAdapterService},
		{"orc8rPrometheusNginxProxyService", r.orc8rPrometheusNginxProxyService},
		{"orc8rUserGrafanaService", r.orc8rUserGrafanaService},
		{"orc8rPrometheusService", r.orc8rPrometheusService},
		{"orc8rAlertManagerConfigurerService", r.orc8rAlertManagerConfigurerService},
		{"orc8rBootstrapNginxService", r.orc8rBootstrapNginxService},
		{"orc8rClientcertService", r.orc8rClientcertService},
	}
	// Iterate over service functions
	for _, service := range serviceFunctions {
		r.Log.Info("Ensuring Service", "ServiceName", service.Name)

		// Call ensureService with the service function
		result, err := r.ensureService(pmnsystem, service.Func(pmnsystem))
		if result != nil || err != nil {
			return *result, err
		}
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
	for _, secretConfig := range pmnsystem.Spec.Secrets {
		secretName := secretConfig.SecretName
		certDir := pmnsystem.Spec.CertDir
		requiredFiles := secretConfig.RequiredFiles
		namespace := pmnsystem.Spec.NameSpace

		// Call CreateSecretsFromCertificates for each secret
		err = r.CreateSecretsFromCertificates(secretName, certDir, requiredFiles, namespace, pmnsystem)
		if err != nil {
			r.Log.Error(err, "Failed to create secret", "SecretName", secretName)
			return ctrl.Result{}, err
		}

		r.Log.Info("Secret created successfully", "SecretName", secretName)
	}
	// ====ensure PVC====
	result, err = r.ensurePersistentVolumeClaim(pmnsystem, r.createPersistentVolumeClaim(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePersistentVolumeClaim(pmnsystem, r.createPersistentVolumeClaimPromData(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePersistentVolumeClaim(pmnsystem, r.createPersistentVolumeClaimGrafanaProviders(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePersistentVolumeClaim(pmnsystem, r.createPersistentVolumeClaimGrafanaData(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePersistentVolumeClaim(pmnsystem, r.createPersistentVolumeClaimGrafanaDatasources(pmnsystem))
	if result != nil {
		return *result, err
	}
	result, err = r.ensurePersistentVolumeClaim(pmnsystem, r.createPersistentVolumeClaimGrafanaDashboards(pmnsystem))
	if result != nil {
		return *result, err
	}

	// Create the Job
	job := r.createOrc8rMetricsStoreConfigJob(pmnsystem)
	r.Log.Info("Creating Job", "Job.Name", job.Name)

	// Check if the Job already exists
	err = r.Client.Create(ctx, job)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			r.Log.Info("Job already exists, skipping creation", "Job.Name", job.Name)
		} else {
			r.Log.Error(err, "Failed to create Job", "Job.Name", job.Name)
			return ctrl.Result{}, err
		}
	} else {
		r.Log.Info("Job created successfully", "Job.Name", job.Name)
	}


	// Ensure Ingress exists and is recreated if deleted
    if err := r.ensureIngress(ctx, pmnsystem); err != nil {
        log.Log.Error(err, "Failed to ensure Ingress")
        return ctrl.Result{}, err
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
		Owns(&batchv1.Job{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}
