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

	corev1 "k8s.io/api/core/v1"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (r *PmnsystemReconciler) createOrc8rMetricsStoreConfigJob(cr *v1.Pmnsystem) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-metrics-storeconfig",
			Namespace: cr.Spec.NameSpace,
			Labels: map[string]string{
				"app": "orc8r-metrics",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "orc8r-metrics-storeconfig",
					Labels: map[string]string{
						"app": "orc8r-metrics",
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:  "storeconfig",
							Image: "alpine:latest",
							Command: []string{
								"/bin/sh",
								"-c",
								`apk update && apk add --no-cache coreutils
                                 cp -n /mnt/defaults/alertmanager.yml /mnt/configs/ | true
                                 mkdir -p /mnt/configs/alert_rules && chmod +x /mnt/configs/alert_rules
                                 cp -n /mnt/defaults/*rules.yml /mnt/configs/alert_rules/ | true`,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "defaults",
									MountPath: "/mnt/defaults",
								},
								{
									Name:      "configs",
									MountPath: "/mnt/configs",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "defaults",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "orc8r-metrics-defaultconfig",
									},
								},
							},
						},
						{
							Name: "configs",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "promcfg",
								},
							},
						},
					},
				},
			},
		},
	}
}
