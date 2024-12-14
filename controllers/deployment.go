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
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *PmnsystemReconciler) deploymentForOrc8rAccessD(_ *v1.Pmnsystem) *appsv1.Deployment {
	int32Ptr := func(i int32) *int32 { return &i }
	int64Ptr := func(i int64) *int64 { return &i }
	intstrPtr := func(i string) *intstr.IntOrString {
		val := intstr.FromString(i)
		return &val
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-accessd",
			Namespace: "pmn",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "orc8r-accessd",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
			Annotations: map[string]string{
				"app.kubernetes.io/name":       "orc8r-accessd",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":       "orc8r-accessd",
					"app.kubernetes.io/instance":   "orc8r",
					"app.kubernetes.io/managed-by": "Orc8r-Operator",
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       intstrPtr("25%"),
					MaxUnavailable: intstrPtr("25%"),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":       "orc8r-accessd",
						"app.kubernetes.io/instance":   "orc8r",
						"app.kubernetes.io/managed-by": "Orc8r-Operator",
					},
					Annotations: map[string]string{
						"app.kubernetes.io/name":       "orc8r-accessd",
						"app.kubernetes.io/instance":   "orc8r",
						"app.kubernetes.io/managed-by": "Orc8r-Operator",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "certs",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  "orc8r-controller",
									DefaultMode: int32Ptr(420),
								},
							},
						},
						{
							Name: "envdir",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  "pmn-envdir",
									DefaultMode: int32Ptr(420),
								},
							},
						},
						{
							Name: "pmn-configs-orc8r",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  "pmn-configs",
									DefaultMode: int32Ptr(420),
								},
							},
						},
					},
					DNSPolicy:                     corev1.DNSClusterFirst,
					TerminationGracePeriodSeconds: int64Ptr(30),
					RestartPolicy:                 corev1.RestartPolicyAlways,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: "artifactory",
						},
					},
					Containers: []corev1.Container{
						{
							VolumeMounts: []corev1.VolumeMount{
								{Name: "certs", MountPath: "/var/opt/magma/certs", ReadOnly: true},
								{Name: "envdir", MountPath: "/var/opt/magma/envdir", ReadOnly: true},
								{Name: "pmn-configs-orc8r", MountPath: "/var/opt/magma/configs/orc8r", ReadOnly: true},
							},
							Name:            "accessd",
							Image:           "815281572631.dkr.ecr.us-west-2.amazonaws.com/pmn/dev/controller:1.8.0-6c4579b5",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"/usr/bin/envdir",
							},
							Args: []string{
								"/var/opt/magma/envdir",
								"/var/opt/magma/bin/accessd",
								"-logtostderr=true",
								"-v=0",
							},
							Ports: []corev1.ContainerPort{
								{Name: "grpc", ContainerPort: 9091, Protocol: corev1.ProtocolTCP},
								{Name: "grpc-internal", ContainerPort: 9191, Protocol: corev1.ProtocolTCP},
							},
							Env: []corev1.EnvVar{
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
								{Name: "SQL_DRIVER", Value: "postgres"},
								{Name: "SQL_DIALECT", Value: "psql"},
								{Name: "SERVICE_HOSTNAME", ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "status.podIP",
									},
								}},
								{Name: "SERVICE_REGISTRY_MODE", Value: "k8s"},
								{Name: "HELM_RELEASE_NAME", Value: "orc8r"},
								{Name: "SERVICE_REGISTRY_NAMESPACE", Value: "pmn"},
								{Name: "HELM_VERSION_TAG", Value: "1.8.0"},
								{Name: "VERSION_TAG", Value: "1.8.0-6c4579b5"},
								{Name: "ORC8R_DOMAIN_NAME", Value: "magma.test"},
								{Name: "PUBLISHER_PORT", Value: "5442"},
								{Name: "SUBSCRIBER_PORT", Value: "443"},
								{Name: "NOTIF_PUBLISHER", Value: "notifier-internal"},
								{Name: "NOTIF_SUBSCRIBER", Value: "notifier-internal"},
								{Name: "NOTIF_CERT_CA", Value: "notifier-ca.crt"},
								{Name: "NOTIF_SERVER_CERT", Value: "notifier.crt"},
								{Name: "NOTIF_SERVER_KEY", Value: "notifier.key"},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 10,
								PeriodSeconds:       30,
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 9091,
										},
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: 10,
								PeriodSeconds:       30,
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.IntOrString{
											Type:   intstr.Int,
											IntVal: 9091,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}