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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func createOrc8rContainer(imageName, imageTag, imagePullPolicy string) corev1.Container {
	return corev1.Container{
		Name:            "orc8r-container",
		Image:           imageName + ":" + imageTag,
		ImagePullPolicy: corev1.PullPolicy(imagePullPolicy),
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "certs",
				MountPath: "/var/opt/magma/certs/certs",
				ReadOnly:  true,
			},
			{
				Name:      "envdir",
				MountPath: "/var/opt/magma/envdir",
				ReadOnly:  true,
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "http",
				ContainerPort: 8080,
			},
			{
				Name:          "grpc",
				ContainerPort: 9180,
			},
			{
				Name:          "grpc-internal",
				ContainerPort: 9190,
			},
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
			{
				Name:  "SQL_DRIVER",
				Value: "postgres",
			},
			{
				Name:  "SQL_DIALECT",
				Value: "psql",
			},
			{
				Name: "SERVICE_HOSTNAME",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
			{
				Name:  "SERVICE_REGISTRY_MODE",
				Value: "k8s",
			},
			{
				Name:  "HELM_RELEASE_NAME",
				Value: "release-name",
			},
			{
				Name:  "SERVICE_REGISTRY_NAMESPACE",
				Value: "release-namespace",
			},
			{
				Name:  "HELM_VERSION_TAG",
				Value: "\"chart-version\"",
			},
			{
				Name:  "VERSION_TAG",
				Value: "latest",
			},
			{
				Name:  "ORC8R_DOMAIN_NAME",
				Value: "domain-name",
			},
			{
				Name:  "PUBLISHER_PORT",
				Value: "port",
			},
			{
				Name:  "SUBSCRIBER_PORT",
				Value: "port",
			},
			{
				Name:  "NOTIF_PUBLISHER",
				Value: "notifier-internal",
			},
			{
				Name:  "NOTIF_SUBSCRIBER",
				Value: "notifier-internal",
			},
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{IntVal: 9180},
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				TCPSocket: &corev1.TCPSocketAction{
					Port: intstr.IntOrString{IntVal: 9180},
				},
			},
		},
	}
}
