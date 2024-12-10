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
