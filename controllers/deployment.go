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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func defLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
	}
}

var defaultLabels = defLabels()

func (r *PmnsystemReconciler) deployment(
	strategy *appsv1.DeploymentStrategy,
	cr *v1.Pmnsystem,
	name string,
	labels map[string]string,
	command []string,
	args []string,
	volumeMounts []corev1.VolumeMount,
	volumes []corev1.Volume,
	ports []corev1.ContainerPort,
	initContainers []corev1.Container,
	dnsConfig *corev1.PodDNSConfig,
	_ map[string]string, // annotations map[string]string
	envVars []corev1.EnvVar,
	livenessProbe *corev1.Probe,
	readinessProbe *corev1.Probe,
	securityContext *corev1.SecurityContext,
	dnsPolicy corev1.DNSPolicy,
	restartPolicy corev1.RestartPolicy,
	imagePullSecrets []corev1.LocalObjectReference,
	terminationGracePeriodSeconds *int64,
	imagePullPolicy corev1.PullPolicy,
	resources corev1.ResourceRequirements,
	terminationMessagePath string,
	terminationMessagePolicy corev1.TerminationMessagePolicy,
) *appsv1.Deployment {
	finalLabels := make(map[string]string)
	for k, v := range defaultLabels {
		finalLabels[k] = v
	}
	for k, v := range labels {
		finalLabels[k] = v
	}

	if securityContext == nil {
		securityContext = &corev1.SecurityContext{}
	}

	securityContext.Capabilities = &corev1.Capabilities{
		Add: []corev1.Capability{"NET_ADMIN"},
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   cr.Spec.NameSpace,
			Labels:      finalLabels,
			Annotations: finalLabels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: finalLabels,
			},
			Strategy: *strategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      finalLabels,
					Annotations: finalLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:                     name,
							Image:                    "815281572631.dkr.ecr.us-west-2.amazonaws.com/pmn/dev/controller:1.8.0-6c4579b5",
							Command:                  command,
							Args:                     args,
							VolumeMounts:             volumeMounts,
							Ports:                    ports,
							Env:                      envVars,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							SecurityContext:          securityContext,
							ImagePullPolicy:          imagePullPolicy,
							Resources:                resources,
							TerminationMessagePath:   terminationMessagePath,
							TerminationMessagePolicy: terminationMessagePolicy,
						},
					},
					InitContainers:                initContainers,
					DNSConfig:                     dnsConfig,
					DNSPolicy:                     dnsPolicy,
					TerminationGracePeriodSeconds: terminationGracePeriodSeconds,
					RestartPolicy:                 restartPolicy,
					Volumes:                       volumes,
					ImagePullSecrets:              imagePullSecrets,
				},
			},
		},
	}
}

func (r *PmnsystemReconciler) orc8rAccessD(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-accessd",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "pmn-certs",
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
	}

	// Define volumeMounts in a separate variable
	volumeMounts := []corev1.VolumeMount{
		{Name: "certs", MountPath: "/var/opt/magma/certs", ReadOnly: true},
		{Name: "envdir", MountPath: "/var/opt/magma/envdir", ReadOnly: true},
		{Name: "pmn-configs-orc8r", MountPath: "/var/opt/magma/configs/orc8r", ReadOnly: true},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// If Bevo is true, add the NET_ADMIN capability
	// if cr.Spec.Bevo {
	// 	securityContext.Capabilities = &corev1.Capabilities{
	// 		Add: []corev1.Capability{"NET_ADMIN"},
	// 	}
	// }

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: "artifactory"},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9091, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9191, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
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
	}

	readinessProbe := &corev1.Probe{
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
	}

	// Command for the container
	command := []string{
		"/usr/bin/envdir",
	}

	args := []string{
		"/var/opt/magma/envdir",
		"/var/opt/magma/bin/accessd",
		"-logtostderr=true",
		"-v=0",
	}

	strategy := &appsv1.DeploymentStrategy{
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
			MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
		},
	}

	terminationGracePeriodSeconds := int64Ptr(30)

	resources := corev1.ResourceRequirements{}

	terminationMessagePath := "/dev/termination-log"

	terminationMessagePolicy := corev1.TerminationMessagePolicy("File")

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-accessd",
		labels,                        // Labels
		command,                       // Command
		args,                          // args (nil if not needed)
		volumeMounts,                  // Volume mounts
		volumes,                       // Volumes
		ports,                         // Ports (empty if not needed)
		nil,                           // Init containers
		nil,                           // DNS config
		nil,                           // Annotations
		envVars,                       // Environment variables
		livenessProbe,                 // Liveness probe
		readinessProbe,                // Readiness probe
		securityContext,               // Security context
		corev1.DNSClusterFirst,        // DNS policy
		corev1.RestartPolicyAlways,    // Restart policy
		imagePullSecrets,              // Image pull secrets
		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
	)
}
