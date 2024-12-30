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
	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"

	"context"

	"k8s.io/apimachinery/pkg/api/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/api/equality"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *PmnsystemReconciler) ensureDeployment(_ ctrl.Request, instance *v1.Pmnsystem, desired *appsv1.Deployment) (*ctrl.Result, error) {
	// Step 1: Ensure the desired deployment has the correct namespace
	desired.Namespace = instance.Namespace

	// Step 2: Check if the Deployment already exists
	found := &appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: instance.Namespace, // Ensure we're checking in the correct namespace
		Name:      desired.Name,       // Name from the desired deployment
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Step 3: If the Deployment is not found, create it
		r.Log.Info("Deployment not found, creating a new one", "Deployment.Namespace", desired.Namespace, "Deployment.Name", desired.Name)

		// Set the owner reference for proper garbage collection
		if err := ctrl.SetControllerReference(instance, desired, r.Scheme); err != nil {
			r.Log.Error(err, "Failed to set owner reference on Deployment", "Deployment.Namespace", desired.Namespace, "Deployment.Name", desired.Name)
			return &ctrl.Result{}, err
		}

		// Create the Deployment
		err = r.Client.Create(context.TODO(), desired)
		if err != nil {
			// Log the error if creation fails
			r.Log.Error(err, "Failed to create Deployment", "Deployment.Namespace", desired.Namespace, "Deployment.Name", desired.Name)
			return &ctrl.Result{}, err
		}

		// Successfully created the Deployment
		r.Log.Info("Deployment created successfully", "Deployment.Namespace", desired.Namespace, "Deployment.Name", desired.Name)
		return nil, nil
	} else if err != nil {
		// Step 4: Handle other errors when fetching the Deployment
		r.Log.Error(err, "Failed to get Deployment", "Deployment.Namespace", desired.Namespace, "Deployment.Name", desired.Name)
		return &ctrl.Result{}, err
	}

	// Step 5: Compare the found Deployment's Spec with the desired Spec and update if necessary
	if !equality.Semantic.DeepEqual(found.Spec, desired.Spec) {
		r.Log.Info("Updating Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
		found.Spec = desired.Spec // Update the found deployment spec to match the desired spec

		err = r.Client.Update(context.TODO(), found)
		if err != nil {
			// Log the error if update fails
			r.Log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &ctrl.Result{}, err
		}

		// Successfully updated the Deployment
		r.Log.Info("Deployment updated successfully", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	} else {
		// No updates required, the Deployment is already up-to-date
		r.Log.Info("Deployment is up-to-date", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	}

	// Step 6: Return nil to indicate success and no requeue
	return nil, nil
}

func (r *PmnsystemReconciler) ensurePodDisruptionBudget(_ ctrl.Request, instance *v1.Pmnsystem, desired *policyv1.PodDisruptionBudget) (*ctrl.Result, error) {
	// Step 1: Ensure the desired PodDisruptionBudget has the correct namespace
	desired.Namespace = instance.Namespace

	// Step 2: Check if the PodDisruptionBudget already exists
	found := &policyv1.PodDisruptionBudget{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: instance.Namespace, // Ensure we're checking in the correct namespace
		Name:      desired.Name,       // Name from the desired PodDisruptionBudget
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Step 3: If the PodDisruptionBudget is not found, create it
		r.Log.Info("PodDisruptionBudget not found, creating a new one", "PDB.Namespace", desired.Namespace, "PDB.Name", desired.Name)

		// Set the owner reference for proper garbage collection
		if err := ctrl.SetControllerReference(instance, desired, r.Scheme); err != nil {
			r.Log.Error(err, "Failed to set owner reference on PDB", "PDB.Namespace", desired.Namespace, "PDB.Name", desired.Name)
			return &ctrl.Result{}, err
		}

		// Create the PodDisruptionBudget
		err = r.Client.Create(context.TODO(), desired)
		if err != nil {
			// Log the error if creation fails
			r.Log.Error(err, "Failed to create PodDisruptionBudget", "PDB.Namespace", desired.Namespace, "PDB.Name", desired.Name)
			return &ctrl.Result{}, err
		}

		// Successfully created the PDB
		r.Log.Info("PodDisruptionBudget created successfully", "PDB.Namespace", desired.Namespace, "PDB.Name", desired.Name)
		return nil, nil
	} else if err != nil {
		// Step 4: Handle other errors when fetching the PDB
		r.Log.Error(err, "Failed to get PodDisruptionBudget", "PDB.Namespace", desired.Namespace, "PDB.Name", desired.Name)
		return &ctrl.Result{}, err
	}

	// Step 5: Compare the found PDB's Spec with the desired Spec and update if necessary
	if !equality.Semantic.DeepEqual(found.Spec, desired.Spec) {
		r.Log.Info("Updating PodDisruptionBudget", "PDB.Namespace", found.Namespace, "PDB.Name", found.Name)
		found.Spec = desired.Spec // Update the found PDB spec to match the desired spec

		err = r.Client.Update(context.TODO(), found)
		if err != nil {
			// Log the error if update fails
			r.Log.Error(err, "Failed to update PodDisruptionBudget", "PDB.Namespace", found.Namespace, "PDB.Name", found.Name)
			return &ctrl.Result{}, err
		}

		// Successfully updated the PDB
		r.Log.Info("PodDisruptionBudget updated successfully", "PDB.Namespace", found.Namespace, "PDB.Name", found.Name)
	} else {
		// No updates required, the PDB is already up-to-date
		r.Log.Info("PodDisruptionBudget is up-to-date", "PDB.Namespace", found.Namespace, "PDB.Name", found.Name)
	}

	// Step 6: Return nil to indicate success and no requeue
	return nil, nil
}

func (r *PmnsystemReconciler) ensureService(_ *v1.Pmnsystem, desired *corev1.Service) (*ctrl.Result, error) {
	found := &corev1.Service{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: desired.Namespace,
		Name:      desired.Name,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Service not found, create it
		r.Log.Info("Service not found, creating a new one", "Service.Namespace", desired.Namespace, "Service.Name", desired.Name)
		err = r.Client.Create(context.TODO(), desired)
		if err != nil {
			r.Log.Error(err, "Failed to create Service", "Service.Namespace", desired.Namespace, "Service.Name", desired.Name)
			return &ctrl.Result{}, err
		}
		r.Log.Info("Service created successfully", "Service.Namespace", desired.Namespace, "Service.Name", desired.Name)
		return nil, nil
	} else if err != nil {
		// Failed to get the Service due to some error other than NotFound
		r.Log.Error(err, "Failed to get Service", "Service.Namespace", desired.Namespace, "Service.Name", desired.Name)
		return &ctrl.Result{}, err
	}

	// Update the Service if needed
	if !equality.Semantic.DeepEqual(found.Spec, desired.Spec) {
		r.Log.Info("Updating Service", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
		found.Spec = desired.Spec
		err = r.Client.Update(context.TODO(), found)
		if err != nil {
			r.Log.Error(err, "Failed to update Service", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
			return &ctrl.Result{}, err
		}
	}

	r.Log.Info("Service already exists and is up-to-date", "Service.Namespace", found.Namespace, "Service.Name", found.Name)
	return nil, nil
}

func (r *PmnsystemReconciler) ensurePersistentVolumeClaim(_ *v1.Pmnsystem, desired *corev1.PersistentVolumeClaim) (*ctrl.Result, error) {
	found := &corev1.PersistentVolumeClaim{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: desired.Namespace,
		Name:      desired.Name,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// PVC not found, create it
		r.Log.Info("PVC not found, creating a new one", "PVC.Namespace", desired.Namespace, "PVC.Name", desired.Name)
		err = r.Client.Create(context.TODO(), desired)
		if err != nil {
			r.Log.Error(err, "Failed to create PVC", "PVC.Namespace", desired.Namespace, "PVC.Name", desired.Name)
			return &ctrl.Result{}, err
		}
		r.Log.Info("PVC created successfully", "PVC.Namespace", desired.Namespace, "PVC.Name", desired.Name)
		return nil, nil
	} else if err != nil {
		// Error fetching PVC
		r.Log.Error(err, "Failed to get PVC", "PVC.Namespace", desired.Namespace, "PVC.Name", desired.Name)
		return &ctrl.Result{}, err
	}

	// If PVC exists, log and do nothing
	r.Log.Info("PVC already exists", "PVC.Namespace", found.Namespace, "PVC.Name", found.Name)
	return nil, nil
}

func (r *PmnsystemReconciler) ensureStatefulSet(_ ctrl.Request, instance *v1.Pmnsystem, desired *appsv1.StatefulSet) (*ctrl.Result, error) {
	// Step 1: Ensure the desired StatefulSet has the correct namespace
	desired.Namespace = instance.Namespace

	// Step 2: Check if the StatefulSet already exists
	found := &appsv1.StatefulSet{}
	err := r.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: instance.Namespace, // Ensure we're checking in the correct namespace
		Name:      desired.Name,       // Name from the desired StatefulSet
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Step 3: If the StatefulSet is not found, create it
		r.Log.Info("StatefulSet not found, creating a new one", "StatefulSet.Namespace", desired.Namespace, "StatefulSet.Name", desired.Name)

		// Set the owner reference for proper garbage collection
		if err := ctrl.SetControllerReference(instance, desired, r.Scheme); err != nil {
			r.Log.Error(err, "Failed to set owner reference on StatefulSet", "StatefulSet.Namespace", desired.Namespace, "StatefulSet.Name", desired.Name)
			return &ctrl.Result{}, err
		}

		// Create the StatefulSet
		err = r.Client.Create(context.TODO(), desired)
		if err != nil {
			// Log the error if creation fails
			r.Log.Error(err, "Failed to create StatefulSet", "StatefulSet.Namespace", desired.Namespace, "StatefulSet.Name", desired.Name)
			return &ctrl.Result{}, err
		}

		// Successfully created the StatefulSet
		r.Log.Info("StatefulSet created successfully", "StatefulSet.Namespace", desired.Namespace, "StatefulSet.Name", desired.Name)
		return nil, nil
	} else if err != nil {
		// Step 4: Handle other errors when fetching the StatefulSet
		r.Log.Error(err, "Failed to get StatefulSet", "StatefulSet.Namespace", desired.Namespace, "StatefulSet.Name", desired.Name)
		return &ctrl.Result{}, err
	}

	// Step 5: Compare the found StatefulSet's Spec with the desired Spec and update if necessary
	if !equality.Semantic.DeepEqual(found.Spec, desired.Spec) {
		r.Log.Info("Updating StatefulSet", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
		found.Spec = desired.Spec // Update the found StatefulSet spec to match the desired spec

		err = r.Client.Update(context.TODO(), found)
		if err != nil {
			// Log the error if update fails
			r.Log.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
			return &ctrl.Result{}, err
		}

		// Successfully updated the StatefulSet
		r.Log.Info("StatefulSet updated successfully", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
	} else {
		// No updates required, the StatefulSet is already up-to-date
		r.Log.Info("StatefulSet is up-to-date", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
	}

	// Step 6: Return nil to indicate success and no requeue
	return nil, nil
}
