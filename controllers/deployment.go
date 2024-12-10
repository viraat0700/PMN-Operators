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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createDeployment creates a new Deployment for the Pmnsystem instance with hardcoded values
func (r *PmnsystemReconciler) createDeployment(_ *v1.Pmnsystem) *appsv1.Deployment {
	replicas := int32(3)
	imageName := "linuxfoundation.jfrog.io/magma-docker/magmalte"
	imageTag := "1.7.0"
	imagePullPolicy := "IfNotPresent"

	certsEnabled := false
	var volumes []corev1.Volume

	if certsEnabled {
		components := []string{"admin-operator", "bootstrapper", "controller", "certifier", "fluentd", "root", "nms"}
		for _, component := range components {
			volumes = append(volumes, corev1.Volume{
				Name: component,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "orc8r-" + component + "-tls",
					},
				},
			})
		}
	} else {
		// Static volumes when certs are not enabled
		volumes = append(volumes, corev1.Volume{
			Name: "certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "orc8r-certs",
				},
			},
		}, corev1.Volume{
			Name: "envdir",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "orc8r-envdir",
				},
			},
		})

		// Additional volumes from configs
		configs := map[string]string{
			"module1": "secret1",
			"module2": "secret2",
		} // Replace with actual data
		if len(configs) > 0 {
			for module, secretName := range configs {
				volumes = append(volumes, corev1.Volume{
					Name: secretName + "-" + module,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: secretName,
						},
					},
				})
			}
		} else {
			volumes = append(volumes, corev1.Volume{
				Name: "empty-configs",
				VolumeSource: corev1.VolumeSource{
					EmptyDir: &corev1.EmptyDirVolumeSource{},
				},
			})
		}
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "orc8r",
			Labels: map[string]string{
				"app.kubernetes.io/component":  "orc8r",
				"app.kubernetes.io/name":       "orc8r",
				"app.kubernetes.io/managed-by": "orc8r-Operator",
			},
			Annotations: map[string]string{
				"release-name-annotation": "orc8r",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/component":  "orc8r",
					"app.kubernetes.io/name":       "orc8r",
					"app.kubernetes.io/managed-by": "orc8r-Operator",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component":  "orc8r",
						"app.kubernetes.io/name":       "orc8r",
						"app.kubernetes.io/managed-by": "orc8r-Operator",
					},
					Annotations: map[string]string{
						"release-name-annotation": "orc8r",
					},
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"key": "value",
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "key",
							Operator: corev1.TolerationOpEqual,
							Value:    "value",
							Effect:   corev1.TaintEffectNoSchedule,
						},
					},
					Affinity: &corev1.Affinity{
						NodeAffinity: &corev1.NodeAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
								{
									Weight: 1,
									Preference: corev1.NodeSelectorTerm{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "key",
												Operator: corev1.NodeSelectorOpIn,
												Values:   []string{"value"},
											},
										},
									},
								},
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: "artifactory",
						},
					},
					Volumes: volumes,
					Containers: []corev1.Container{
						createOrc8rContainer(imageName, imageTag, imagePullPolicy),
					},
				},
			},
		},
	}
}
