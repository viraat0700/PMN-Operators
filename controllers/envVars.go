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
	"fmt"

	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
)

func (r *PmnsystemReconciler) getEnvVarsForAccessD(cr *v1.Pmnsystem) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	// Default environment variables, always present
	envVars = append(envVars, []corev1.EnvVar{
		// Add specific environment variables here if required
		{
			Name: "DATABASE_SOURCE",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "orc8r-controller",
					},
					Key: "postgres.connstr",
				},
			},
		},
		{Name: "SERVICE_HOSTNAME", ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				APIVersion: "v1",
				FieldPath:  "status.podIP",
			},
		}},
	}...)

	for _, env := range cr.Spec.EnvVariables {
		envVars = append(envVars, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	fmt.Println("FINAL ENV VARIABLES:", envVars)
	return envVars
}
func (r *PmnsystemReconciler) getEnvVarsForDirectoryD(cr *v1.Pmnsystem) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	// Default environment variables, always present
	envVars = append(envVars, []corev1.EnvVar{
		// Add specific environment variables here if required
		{
			Name: "DATABASE_SOURCE",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "orc8r-controller",
					},
					Key: "postgres.connstr",
				},
			},
		},
		{Name: "SERVICE_HOSTNAME", ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				APIVersion: "v1",
				FieldPath:  "status.podIP",
			},
		}},
	}...)

	for _, env := range cr.Spec.EnvVariablesDirectoryD {
		envVars = append(envVars, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	fmt.Println("FINAL ENV VARIABLES:", envVars)
	return envVars
}
func (r *PmnsystemReconciler) getEnvVarsForOrc8rNginx(cr *v1.Pmnsystem) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	for _, env := range cr.Spec.EnvVariablesOrc8rNginx {
		envVars = append(envVars, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	fmt.Println("FINAL ENV VARIABLES:", envVars)
	return envVars
}
