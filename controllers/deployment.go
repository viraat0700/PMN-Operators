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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// func defLabels() map[string]string {
// 	return map[string]string{
// 		"app.kubernetes.io/instance":   "orc8r",
// 		"app.kubernetes.io/managed-by": "Orc8r-Operator",
// 	}
// }

// var defaultLabels = defLabels()

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
	image string,
	affinity *corev1.Affinity,
	replica *int32,
	nodeSelector map[string]string,
	tolerations []corev1.Toleration,
	matchLabels map[string]string,
) *appsv1.Deployment {
	// finalLabels := make(map[string]string)
	// for k, v := range defaultLabels {
	// 	finalLabels[k] = v
	// }
	// for k, v := range labels {
	// 	finalLabels[k] = v
	// }

	if securityContext == nil {
		securityContext = &corev1.SecurityContext{}
	}

	securityContext.Capabilities = &corev1.Capabilities{
		Add: []corev1.Capability{"NET_ADMIN"},
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cr.Spec.NameSpace,
			Labels:    labels,
			// Annotations: labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			Strategy: *strategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchLabels,
					// Annotations: labels,
				},
				Spec: corev1.PodSpec{
					Tolerations:  tolerations,
					NodeSelector: nodeSelector,
					Affinity:     affinity,
					Containers: []corev1.Container{
						{
							Name:                     name,
							Image:                    image,
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
		"app.kubernetes.io/component": "accessd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "accessd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replica
		nil,                           // Node selector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rAnalyticsDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "analytics",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9200, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9300, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9200,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9200,
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
		"/var/opt/magma/bin/analytics",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "analytics",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-analytics",
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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replicas
		nil,                           // Node selector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rBootStrapperDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "bootstrapper",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9088, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9188, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9088,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9088,
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
		"/var/opt/magma/bin/bootstrapper",
		"-cak=/var/opt/magma/certs/bootstrapper.key",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "bootstrapper",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-bootstrapper",
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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replicas
		nil,                           // Node selector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rCertifierDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "certifier",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9086, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9186, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10089, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9086,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9086,
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
		"/var/opt/magma/bin/certifier",
		"-cac=/var/opt/magma/certs/certifier.pem",
		"-cak=/var/opt/magma/certs/certifier.key",
		"-vpnc=/var/opt/magma/certs/vpn_ca.crt",
		"-vpnk=/var/opt/magma/certs/vpn_ca.key",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "certifier",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-certifier",
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
		image,                         // Image
		nil,                           // Afinity
		&replicas,                     // Replica
		nil,                           // NodeSelector
		tolerations,                   // Tolerations
		matchLabels,                   // MatchLabels
	)
}

// func (r *PmnsystemReconciler) orc8rConfiguratorDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
// 	int64Ptr := func(i int64) *int64 { return &i }
// 	int32Ptr := func(i int32) *int32 { return &i }

// 	labels := map[string]string{
// 		"app.kubernetes.io/component": "configurator",
// 		"app.kubernetes.io/instance":  "orc8r",
// 		"app.kubernetes.io/name":      "orc8r",
// 		"app.kubernetes.io/part-of":   "orc8r-app",
// 	}

// 	// Define volumes in a separate variable
// 	volumes := []corev1.Volume{
// 		{
// 			Name: "certs",
// 			VolumeSource: corev1.VolumeSource{
// 				Secret: &corev1.SecretVolumeSource{
// 					SecretName:  "pmn-certs",
// 					DefaultMode: int32Ptr(420),
// 				},
// 			},
// 		},
// 		{
// 			Name: "envdir",
// 			VolumeSource: corev1.VolumeSource{
// 				Secret: &corev1.SecretVolumeSource{
// 					SecretName:  "pmn-envdir",
// 					DefaultMode: int32Ptr(420),
// 				},
// 			},
// 		},
// 		{
// 			Name: "pmn-configs-orc8r",
// 			VolumeSource: corev1.VolumeSource{
// 				Secret: &corev1.SecretVolumeSource{
// 					SecretName:  "pmn-configs",
// 					DefaultMode: int32Ptr(420),
// 				},
// 			},
// 		},
// 	}

// 	// Define volumeMounts in a separate variable
// 	volumeMounts := []corev1.VolumeMount{
// 		{Name: "certs", MountPath: "/var/opt/magma/certs", ReadOnly: true},
// 		{Name: "envdir", MountPath: "/var/opt/magma/envdir", ReadOnly: true},
// 		{Name: "pmn-configs-orc8r", MountPath: "/var/opt/magma/configs/orc8r", ReadOnly: true},
// 	}

// 	// Define the securityContext for the container
// 	securityContext := &corev1.SecurityContext{
// 		Privileged: func(b bool) *bool { return &b }(true),
// 	}

// 	// Define imagePullSecrets
// 	imagePullSecrets := []corev1.LocalObjectReference{
// 		{Name: cr.Spec.ImagePullSecrets},
// 	}

// 	// Define environment variables if needed
// 	envVars := r.getEnvVarsForAccessD(cr)

// 	// Define ports (use nil if not needed)
// 	ports := []corev1.ContainerPort{
// 		{Name: "grpc", ContainerPort: 9108, Protocol: corev1.ProtocolTCP},
// 		{Name: "grpc-internal", ContainerPort: 9208, Protocol: corev1.ProtocolTCP},
// 		{Name: "moso", ContainerPort: 8088, Protocol: corev1.ProtocolTCP},
// 	}

// 	// Liveness and Readiness Probes
// 	livenessProbe := &corev1.Probe{
// 		FailureThreshold:    3,
// 		SuccessThreshold:    1,
// 		TimeoutSeconds:      1,
// 		InitialDelaySeconds: 10,
// 		PeriodSeconds:       30,
// 		ProbeHandler: corev1.ProbeHandler{
// 			TCPSocket: &corev1.TCPSocketAction{
// 				Port: intstr.IntOrString{
// 					Type:   intstr.Int,
// 					IntVal: 9108,
// 				},
// 			},
// 		},
// 	}

// 	readinessProbe := &corev1.Probe{
// 		FailureThreshold:    3,
// 		SuccessThreshold:    1,
// 		TimeoutSeconds:      1,
// 		InitialDelaySeconds: 10,
// 		PeriodSeconds:       30,
// 		ProbeHandler: corev1.ProbeHandler{
// 			TCPSocket: &corev1.TCPSocketAction{
// 				Port: intstr.IntOrString{
// 					Type:   intstr.Int,
// 					IntVal: 9108,
// 				},
// 			},
// 		},
// 	}

// 	// Command for the container
// 	command := []string{
// 		"/usr/bin/envdir",
// 	}

// 	args := []string{
// 		"/var/opt/magma/envdir",
// 		"/var/opt/magma/bin/configurator",
// 		"-logtostderr=true",
// 		"-v=0",
// 	}

// 	strategy := &appsv1.DeploymentStrategy{
// 		RollingUpdate: &appsv1.RollingUpdateDeployment{
// 			MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
// 			MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
// 		},
// 	}

// 	terminationGracePeriodSeconds := int64Ptr(30)

// 	resources := corev1.ResourceRequirements{}

// 	terminationMessagePath := "/dev/termination-log"

// 	terminationMessagePolicy := corev1.TerminationMessagePolicy("File")

// 	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

// 	replicas := &cr.Spec.ReplicaCount

// 	tolerations := []corev1.Toleration{}

// 	matchlabels := map[string]string{
// 		"app.kubernetes.io/component": "configurator",
// 		"app.kubernetes.io/instance":  "orc8r",
// 		"app.kubernetes.io/name":      "orc8r",
// 	}

// 	return r.deployment(
// 		strategy, // Deployment strategy
// 		cr,
// 		"orc8r-configurator",
// 		labels,                        // Labels
// 		command,                       // Command
// 		args,                          // args (nil if not needed)
// 		volumeMounts,                  // Volume mounts
// 		volumes,                       // Volumes
// 		ports,                         // Ports (empty if not needed)
// 		nil,                           // Init containers
// 		nil,                           // DNS config
// 		nil,                           // Annotations
// 		envVars,                       // Environment variables
// 		livenessProbe,                 // Liveness probe
// 		readinessProbe,                // Readiness probe
// 		securityContext,               // Security context
// 		corev1.DNSClusterFirst,        // DNS policy
// 		corev1.RestartPolicyAlways,    // Restart policy
// 		imagePullSecrets,              // Image pull secrets
// 		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
// 		corev1.PullIfNotPresent,       // Image pull policy
// 		resources,                     // Resources
// 		terminationMessagePath,        // Termination message path
// 		terminationMessagePolicy,      // Termination message policy
// 		image,                         // Image
// 		nil,                           // Affinity
// 		replicas,                      // Replicas
// 		nil,                           // Node selector
// 		tolerations,                   // Tolerations
// 		matchlabels,                   // Match labels
// 	)
// }

func (r *PmnsystemReconciler) orc8rConfiguratorDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component":  "configurator",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Helm",
		"app.kubernetes.io/name":       "orc8r",
		"app.kubernetes.io/part-of":    "orc8r-app",
		"helm.sh/chart":                "orc8r-1.8.0",
	}

	annotation := map[string]string{
		"chart-version":                     "1.8.0",
		"deployment.kubernetes.io/revision": "165",
		"meta.helm.sh/release-name":         "orc8r",
		"meta.helm.sh/release-namespace":    "pmn",
		"release-name":                      "orc8r",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9108, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9208, Protocol: corev1.ProtocolTCP},
		{Name: "moso", ContainerPort: 8088, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9108,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9108,
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
		"/var/opt/magma/bin/configurator",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchlabels := map[string]string{
		"app.kubernetes.io/component": "configurator",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-configurator",
			Namespace: cr.Spec.NameSpace,
			Labels:    labels,
			Annotations: annotation,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: matchlabels,
			},
			Strategy: *strategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchlabels,
				},
				Spec: corev1.PodSpec{
					Tolerations: tolerations,
					Containers: []corev1.Container{
						{
							Name:                     "orc8r-configurator",
							Image:                    image,
							Command:                  command,
							Args:                     args,
							Ports:                    ports,
							VolumeMounts:             volumeMounts,
							Env:                      envVars,
							SecurityContext:          securityContext,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							Resources:                resources,
							TerminationMessagePath:   terminationMessagePath,
							TerminationMessagePolicy: terminationMessagePolicy,
							ImagePullPolicy:          corev1.PullIfNotPresent,
						},
					},
					ImagePullSecrets:              imagePullSecrets,
					RestartPolicy:                 corev1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: terminationGracePeriodSeconds,
					Volumes:                       volumes,
				},
			},
		},
	}
}

func (r *PmnsystemReconciler) orc8rDeviceDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "device",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9106, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9306, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9106,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9106,
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
		"/var/opt/magma/bin/device",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "device",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-device",
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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replicas
		nil,                           // Node selector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rDirectorydDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "directoryd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9100, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9102, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9100,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9100,
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
		"/var/opt/magma/bin/directoryd",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "directoryd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-directoryd",
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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replica
		nil,                           // NodeSelector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rDispatcherDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "dispatcher",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9096, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9196, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 9080, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9096,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9096,
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
		"/var/opt/magma/bin/dispatcher",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "dispatcher",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-dispatcher",
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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replicas
		nil,                           // NodeSelector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rEventdDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "eventd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9121, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9221, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10121, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9121,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9121,
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
		"/var/opt/magma/bin/eventd",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "eventd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-eventd",
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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replicas
		nil,                           // Node selector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rmetricsdDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "metricsd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9084, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9184, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10084, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9084,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9084,
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
		"/var/opt/magma/bin/metricsd",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "metricsd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-metricsd",
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
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replicas
		nil,                           // Node selector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rNginxDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "nginx-proxy",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r",
	}

	volumes := []corev1.Volume{
		{
			Name: "certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  cr.Spec.Orc8rNginxDeployment.VolumesOrc8rNginx.SecretName[0],
					DefaultMode: int32Ptr(420),
				},
			},
		},
		{
			Name: "envdir",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  cr.Spec.Orc8rNginxDeployment.VolumesOrc8rNginx.SecretName[0],
					DefaultMode: int32Ptr(420),
				},
			},
		},
	}

	volumeMounts := []corev1.VolumeMount{
		{Name: "certs", MountPath: cr.Spec.Orc8rNginxDeployment.VolumesMountPathOrc8rNginx.MountPath[0], ReadOnly: true},
		{Name: "envdir", MountPath: cr.Spec.Orc8rNginxDeployment.VolumesMountPathOrc8rNginx.MountPath[1], ReadOnly: true},
	}

	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	envVars := r.getEnvVarsForOrc8rNginx(cr)

	ports := []corev1.ContainerPort{
		{Name: "clientcert", ContainerPort: cr.Spec.Orc8rNginxDeployment.PortOrc8rNginx.Port[0], Protocol: corev1.ProtocolTCP},
		{Name: "open", ContainerPort: cr.Spec.Orc8rNginxDeployment.PortOrc8rNginx.Port[1], Protocol: corev1.ProtocolTCP},
		{Name: "api", ContainerPort: cr.Spec.Orc8rNginxDeployment.PortOrc8rNginx.Port[2], Protocol: corev1.ProtocolTCP},
		{Name: "health", ContainerPort: cr.Spec.Orc8rNginxDeployment.PortOrc8rNginx.Port[3], Protocol: corev1.ProtocolTCP},
	}

	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromString("health"),
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromString("health"),
			},
		},
	}

	strategy := appsv1.DeploymentStrategy{
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
			MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
		},
	}

	terminationGracePeriodSeconds := int64Ptr(30)

	resources := corev1.ResourceRequirements{}

	terminationMessagePath := "/dev/termination-log"
	terminationMessagePolicy := corev1.TerminationMessagePolicy("File")

	image := cr.Spec.Orc8rNginxDeployment.ImageOrc8rNginx.Repository + ":" + cr.Spec.Orc8rNginxDeployment.ImageOrc8rNginx.Tag

	replicas := cr.Spec.Orc8rNginxDeployment.Replicas

	tolerations := cr.Spec.Orc8rNginxDeployment.Tolerations
	nodeSelector := cr.Spec.Orc8rNginxDeployment.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{}
	}

	imagePullPolicy := corev1.PullPolicy(cr.Spec.Orc8rNginxDeployment.ImageOrc8rNginx.ImagePullPolicy)

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "nginx-proxy",
	}

	labelsInTemplate := map[string]string{
		"app.kubernetes.io/component": "nginx-proxy",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-nginx",
			Namespace: cr.Spec.NameSpace,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			Strategy: strategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelsInTemplate,
				},
				Spec: corev1.PodSpec{
					NodeSelector: nodeSelector,
					Tolerations:  tolerations,
					Containers: []corev1.Container{
						{
							Name:                     "orc8r-nginx",
							Image:                    image,
							Ports:                    ports,
							VolumeMounts:             volumeMounts,
							Env:                      envVars,
							SecurityContext:          securityContext,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							Resources:                resources,
							TerminationMessagePath:   terminationMessagePath,
							TerminationMessagePolicy: terminationMessagePolicy,
							ImagePullPolicy:          imagePullPolicy,
						},
					},
					ImagePullSecrets:              imagePullSecrets,
					RestartPolicy:                 corev1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: terminationGracePeriodSeconds,
					Volumes:                       volumes,
				},
			},
		},
	}
}

func (r *PmnsystemReconciler) orc8rNotifierDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "notifier",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForOrc8rNotifier(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "notifier", ContainerPort: cr.Spec.Orc8rNotifier.PortDeployment, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: cr.Spec.Orc8rNotifier.LivenessProbePort,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: cr.Spec.Orc8rNotifier.ReadinessProbePort,
				},
			},
		},
	}

	args := cr.Spec.Orc8rNotifier.Args

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

	image := cr.Spec.Orc8rNotifier.ImageOrc8rNotifier.Repository + ":" + cr.Spec.Orc8rNotifier.ImageOrc8rNotifier.Tag

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.Orc8rNotifier.ImageOrc8rNotifier.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.Orc8rNotifier.ImageOrc8rNotifier.ImagePullPolicy)
	}

	replicas := cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchlabels := map[string]string{
		"app.kubernetes.io/component": "notifier",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-notifier",
		labels,                        // Labels
		nil,                           // Command
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
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
		&replicas,                     // Replicas
		nil,                           // Node selector
		tolerations,                   // Tolerations
		matchlabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) orc8rObsidianDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "obsidian",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9093, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9193, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 9081, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9081,
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
					IntVal: 9081,
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
		"/var/opt/magma/bin/obsidian",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchlabels := map[string]string{
		"app.kubernetes.io/component": "obsidian",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-obsidian",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchlabels,                   // match labels
	)
}

func (r *PmnsystemReconciler) orc8WorkerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r-worker",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9122, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9222, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9122,
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
					IntVal: 9122,
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
		"/var/opt/magma/bin/orc8r_worker",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "orc8r-worker",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-orc8r-worker",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8orchestratorDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "orchestrator",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9112, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9212, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10112, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9112,
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
					IntVal: 9112,
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
		"/var/opt/magma/bin/orchestrator",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "orchestrator",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-orchestrator",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8ServiceRegistryDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "service_registry",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9180, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9190, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9180,
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
					IntVal: 9180,
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
		"/var/opt/magma/bin/service_registry",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "service_registry",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-service-registry",
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
		image,                         // Image
		nil,                           // affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8StateDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "state",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9105, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9305, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9105,
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
					IntVal: 9105,
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
		"/var/opt/magma/bin/state",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "state",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-state",
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
		image,                         // Image
		nil,                           // afinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels

	)
}

func (r *PmnsystemReconciler) orc8StreamerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "streamer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9082, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9182, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9082,
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
					IntVal: 9082,
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
		"/var/opt/magma/bin/streamer",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "streamer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-streamer",
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
		image,                         // Image
		nil,                           // affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8TenantsDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "tenants",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForAccessD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9110, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9210, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10110, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9110,
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
					IntVal: 9110,
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
		"/var/opt/magma/bin/tenants",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchlabels := map[string]string{
		"app.kubernetes.io/component": "tenants",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-tenants",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchlabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8rHaDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "ha",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9119, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9119,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9119,
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
		"/var/opt/magma/bin/ha",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "ha",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-ha",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8LteDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "lte",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9113, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9213, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10113, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9113,
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
					IntVal: 9113,
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
		"/var/opt/magma/bin/lte",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "lte",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-lte",
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
		image,                         // Image
		nil,                           // affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8NprobeDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "nprobe",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9666, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9766, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10088, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9666,
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
					IntVal: 9666,
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
		"/var/opt/magma/bin/nprobe",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "nprobe",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-nprobe",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8PolicyDbDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "policydb",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9085, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9185, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10085, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9085,
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
					IntVal: 9085,
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
		"/var/opt/magma/bin/policydb",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "policydb",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-policydb",
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
		image,                         // Image
		nil,                           // affinity
		replicas,                      // replicas
		nil,                           // nodSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8SmsdDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "smsd",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9120, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9220, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10086, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9120,
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
					IntVal: 9120,
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
		"/var/opt/magma/bin/smsd",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "smsd",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-smsd",
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
		image,                         // Image
		nil,                           // affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8SubscriberDbCacheDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "subscriberdb-cache",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9089, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10087, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9089,
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
					IntVal: 9089,
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
		"/var/opt/magma/bin/subscriberdb_cache",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "subscriberdb-cache",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-subscriberdb-cache",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nil,                           // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8SubscriberDbDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "subscriberdb",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
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

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForDirectoryD(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "grpc", ContainerPort: 9083, Protocol: corev1.ProtocolTCP},
		{Name: "grpc-internal", ContainerPort: 9183, Protocol: corev1.ProtocolTCP},
		{Name: "http", ContainerPort: 10083, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 9083,
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
					IntVal: 9083,
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
		"/var/opt/magma/bin/subscriberdb",
		"-run_echo_server=true",
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

	replicas := &cr.Spec.ReplicaCount

	tolerations := []corev1.Toleration{}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "subscriberdb",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-subscriberdb",
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
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // Replicas
		nil,                           // NodeSelector
		tolerations,                   // Tolerations
		matchLabels,                   // Match labels
	)
}

func (r *PmnsystemReconciler) nmsMagmaLteDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "magmalte",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "nms",
		"app.kubernetes.io/part-of":   "magma",
		"release_group":               "orc8r",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "orc8r-secrets-certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  cr.Spec.NmsMagmaLte.VolumesNmsMagmaLte.Secretname[0],
					DefaultMode: int32Ptr(292),
				},
			},
		},
	}

	// Define volumeMounts in a separate variable
	volumeMounts := []corev1.VolumeMount{
		{Name: "orc8r-secrets-certs", MountPath: cr.Spec.NmsMagmaLte.VolumeMountNmsMagmaLte.VolumeMountPath[0], SubPath: cr.Spec.NmsMagmaLte.VolumeMountNmsMagmaLte.VolumeSubPath[0]},
		{Name: "orc8r-secrets-certs", MountPath: cr.Spec.NmsMagmaLte.VolumeMountNmsMagmaLte.VolumeMountPath[1], SubPath: cr.Spec.NmsMagmaLte.VolumeMountNmsMagmaLte.VolumeSubPath[1]},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForNmsMagmaLte(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{ContainerPort: 8081, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		FailureThreshold:    3,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/healthz",
				Scheme: corev1.URISchemeHTTP,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8081,
				},
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		FailureThreshold:    3,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/healthz",
				Scheme: corev1.URISchemeHTTP,
				Port: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8081,
				},
			},
		},
	}

	args := []string{
		"yarn",
		"run",
		"start:prod",
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

	image := cr.Spec.NmsMagmaLte.ImageMagmaLte.Repository + ":" + cr.Spec.NmsMagmaLte.ImageMagmaLte.Tag

	replicas := &cr.Spec.NmsMagmaLte.Replicas

	tolerations := cr.Spec.NmsMagmaLte.Tolerations

	nodeSelector := cr.Spec.NmsMagmaLte.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{} // Default to an empty map
	}

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.NmsMagmaLte.ImageMagmaLte.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.NmsMagmaLte.ImageMagmaLte.ImagePullPolicy)
	}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "magmalte",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "nms",
		"release_group":               "orc8r",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"nms-magmalte",
		labels,                        // Labels
		nil,                           // Command
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
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nodeSelector,                  // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // match Labels
	)
}

func (r *PmnsystemReconciler) orc8PrometheusCacheDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	// int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus-cache",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{ContainerPort: 9091, Protocol: corev1.ProtocolTCP},
		{ContainerPort: 9092, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/",
				Port:   intstr.FromInt(9091),
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/",
				Port:   intstr.FromInt(9091),
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	args := cr.Spec.PrometheusCache.Args

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

	image := cr.Spec.PrometheusCache.ImagePrometheusCache.Repository + ":" + cr.Spec.PrometheusCache.ImagePrometheusCache.Tag

	replicas := &cr.Spec.PrometheusCache.Replicas

	nodeSelector := cr.Spec.PrometheusCache.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{} // Default to an empty map
	}

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.PrometheusCache.ImagePrometheusCache.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.PrometheusCache.ImagePrometheusCache.ImagePullPolicy)
	}

	tolerations := cr.Spec.PrometheusCache.Tolerations

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus-cache",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-prometheus-cache",
		labels,                        // Labels
		nil,                           // Command
		args,                          // args (nil if not needed)
		nil,                           // Volume mounts
		nil,                           // Volumes
		ports,                         // Ports (empty if not needed)
		nil,                           // Init containers
		nil,                           // DNS config
		nil,                           // Annotations
		nil,                           // Environment variables
		livenessProbe,                 // Liveness probe
		readinessProbe,                // Readiness probe
		securityContext,               // Security context
		corev1.DNSClusterFirst,        // DNS policy
		corev1.RestartPolicyAlways,    // Restart policy
		imagePullSecrets,              // Image pull secrets
		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // Replica
		nodeSelector,                  // NodeSelector
		tolerations,                   // Tolerations
		matchLabels,                   // match labels
	)
}

func (r *PmnsystemReconciler) orc8rPrometheusConfigurerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	// int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "prometheus-config",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: cr.Spec.PrometheusConfigurer.Volume.VolumeClaimName,
				},
			},
		},
	}

	// VolumeMounts
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "prometheus-config",
			MountPath: cr.Spec.PrometheusConfigurer.Volume.VolumeMountPath,
		},
	}

	// Affinity
	affinity := &corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app.kubernetes.io/component",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"prometheus"},
							},
						},
					},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{ContainerPort: 9100, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(9100),
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(9100),
			},
		},
	}

	args := cr.Spec.PrometheusConfigurer.Args

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

	image := cr.Spec.PrometheusConfigurer.ImagePrometheusConfigurer.Repository + ":" + cr.Spec.PrometheusConfigurer.ImagePrometheusConfigurer.Tag

	replicas := &cr.Spec.PrometheusConfigurer.Replicas

	nodeSelector := cr.Spec.PrometheusConfigurer.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{} // Default to an empty map
	}

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.PrometheusConfigurer.ImagePrometheusConfigurer.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.PrometheusConfigurer.ImagePrometheusConfigurer.ImagePullPolicy)
	}

	tolerations := cr.Spec.PrometheusConfigurer.Tolerations

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-prometheus-configurer",
		labels,                        // Labels
		nil,                           // Command
		args,                          // args (nil if not needed)
		volumeMounts,                  // Volume mounts
		volumes,                       // Volumes
		ports,                         // Ports (empty if not needed)
		nil,                           // Init containers
		nil,                           // DNS config
		nil,                           // Annotations
		nil,                           // Environment variables
		livenessProbe,                 // Liveness probe
		readinessProbe,                // Readiness probe
		securityContext,               // Security context
		corev1.DNSClusterFirst,        // DNS policy
		corev1.RestartPolicyAlways,    // Restart policy
		imagePullSecrets,              // Image pull secrets
		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		affinity,                      // Affinity
		replicas,                      // replicas
		nodeSelector,                  // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // match Labels
	)
}

func (r *PmnsystemReconciler) orc8rPrometheusKafkaAdapterDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/instance": "orc8r",
		"app.kubernetes.io/name":     "prometheus-kafka-adapter",
		"app.kubernetes.io/version":  "1.0",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "ssl-client-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  cr.Spec.PrometheusKafkaAdapter.VolumeMountPathPrometheusKafkaAdapter.SecretName,
					DefaultMode: int32Ptr(420),
				},
			},
		},
	}

	// VolumeMounts
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "ssl-client-cert",
			MountPath: cr.Spec.PrometheusKafkaAdapter.VolumeMountPathPrometheusKafkaAdapter.MountPath[0],
			ReadOnly:  true,
		},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForPrometheusKafkaAdapter(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "http", ContainerPort: 8080, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/healthz",
				Port:   intstr.FromString("http"),
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/healthz",
				Port:   intstr.FromString("http"),
				Scheme: corev1.URISchemeHTTP,
			},
		},
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

	image := cr.Spec.PrometheusKafkaAdapter.ImagePrometheusKafkaAdapter.Repository + ":" + cr.Spec.PrometheusKafkaAdapter.ImagePrometheusKafkaAdapter.Tag

	nodeSelector := cr.Spec.PrometheusKafkaAdapter.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{} // Default to an empty map
	}

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.PrometheusKafkaAdapter.ImagePrometheusKafkaAdapter.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.UserGrafana.ImageUserGrafana.ImagePullPolicy)
	}

	replicas := &cr.Spec.PrometheusKafkaAdapter.Replicas

	tolerations := cr.Spec.PrometheusKafkaAdapter.Tolerations

	matchLabels := map[string]string{
		"app.kubernetes.io/instance": "orc8r",
		"app.kubernetes.io/name":     "prometheus-kafka-adapter",
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-prometheus-kafka-adapter",
		labels,                        // Labels
		nil,                           // Command
		nil,                           // args (nil if not needed)
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
		nil,                           // Image pull secrets
		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // Replica
		nodeSelector,                  // Node selector
		tolerations,                   // Tolerations
		matchLabels,                   // matchLabels
	)
}

func (r *PmnsystemReconciler) orc8rPrometheusNginxProxyDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus-nginx",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus-nginx",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "prometheus-certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  cr.Spec.PrometheusNginxProxy.Nginx.SecretName,
					DefaultMode: int32Ptr(292),
				},
			},
		},
		{
			Name: "prometheus-nginx-proxy",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "prometheus-nginx-proxy",
					},
				},
			},
		},
	}

	// VolumeMounts
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "prometheus-nginx-proxy",
			MountPath: cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.MountPath[0],
			SubPath:   cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.SubPath[0],
		},
		{
			Name:      "prometheus-certs",
			MountPath: cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.MountPath[1],
			ReadOnly:  true,
			SubPath:   cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.SubPath[1],
		},
		{
			Name:      "prometheus-certs",
			MountPath: cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.MountPath[2],
			SubPath:   cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.SubPath[2],
		},
		{
			Name:      "prometheus-certs",
			MountPath: cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.MountPath[3],
			SubPath:   cr.Spec.PrometheusNginxProxy.Nginx.VolumeMountPath.SubPath[3],
			ReadOnly:  true,
		},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{ContainerPort: 443, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(443),
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(443),
			},
		},
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

	replicas := &cr.Spec.PrometheusNginxProxy.Nginx.Replica

	image := cr.Spec.PrometheusNginxProxy.Nginx.ImagePrometheusNginxProxy.Repository + ":" + cr.Spec.PrometheusNginxProxy.Nginx.ImagePrometheusNginxProxy.Tag

	tolerations := []corev1.Toleration{}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-prometheus-nginx-proxy", // Name
		labels,                         // Labels
		nil,                            // Command
		nil,                            // args (nil if not needed)
		volumeMounts,                   // Volume mounts
		volumes,                        // Volumes
		ports,                          // Ports (empty if not needed)
		nil,                            // Init containers
		nil,                            // DNS config
		nil,                            // Annotations
		nil,                            // Environment variables
		livenessProbe,                  // Liveness probe
		readinessProbe,                 // Readiness probe
		securityContext,                // Security context
		corev1.DNSClusterFirst,         // DNS policy
		corev1.RestartPolicyAlways,     // Restart policy
		nil,                            // Image pull secrets
		terminationGracePeriodSeconds,  // terminationGracePeriodSeconds
		corev1.PullIfNotPresent,        // Image pull policy
		resources,                      // Resources
		terminationMessagePath,         // Termination message path
		terminationMessagePolicy,       // Termination message policy
		image,                          // Image
		nil,                            // Affinity
		replicas,                       // replicas
		nil,                            // Nodeselector
		tolerations,                    // toleration
		matchLabels,                    // match labels
	)
}

func (r *PmnsystemReconciler) orc8rUserGrafanaDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "user-grafana",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "user-grafana",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	claimNameGrafanaDataSource := cr.Spec.UserGrafana.VolumesUserGrafana[0].Name
	claimNameGrafanaProviders := cr.Spec.UserGrafana.VolumesUserGrafana[1].Name
	claimNameGrafanaDashboards := cr.Spec.UserGrafana.VolumesUserGrafana[2].Name
	claimNameGrafanaData := cr.Spec.UserGrafana.VolumesUserGrafana[3].Name
	mountPathGrafanaDataSource := cr.Spec.UserGrafana.VolumesUserGrafana[0].Path
	mountPathGrafanaProviders := cr.Spec.UserGrafana.VolumesUserGrafana[1].Path
	mountPathGrafanaDashBoards := cr.Spec.UserGrafana.VolumesUserGrafana[2].Path
	mountPathGrafanaData := cr.Spec.UserGrafana.VolumesUserGrafana[3].Path

	// Volumes
	volumes := []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "grafana-config-file",
					},
					DefaultMode: int32Ptr(420),
				},
			},
		},
		{
			Name: "datasources",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claimNameGrafanaDataSource,
				},
			},
		},
		{
			Name: "dashboardproviders",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claimNameGrafanaProviders,
				},
			},
		},
		{
			Name: "dashboards",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claimNameGrafanaDashboards,
				},
			},
		},
		{
			Name: "grafana-data",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claimNameGrafanaData,
				},
			},
		},
	}

	// VolumeMounts
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: "/etc/grafana",
		},
		{
			Name:      "dashboards",
			MountPath: mountPathGrafanaDashBoards,
		},
		{
			Name:      "datasources",
			MountPath: mountPathGrafanaDataSource,
		},
		{
			Name:      "dashboardproviders",
			MountPath: mountPathGrafanaProviders,
		},
		{
			Name:      "grafana-data",
			MountPath: mountPathGrafanaData,
		},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{ContainerPort: 3000, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/api/health",
				Port:   intstr.FromInt(3000),
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/api/health",
				Port:   intstr.FromInt(3000),
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	strategy := &appsv1.DeploymentStrategy{
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
			MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
		},
	}

	terminationGracePeriodSeconds := int64Ptr(30)

	resources := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
	}

	terminationMessagePath := "/dev/termination-log"

	terminationMessagePolicy := corev1.TerminationMessagePolicy("File")

	initContainers := []corev1.Container{
		{
			Name:  "volume-mount",
			Image: "busybox",
			Command: []string{
				"sh",
				"-c",
				"chmod -R 777 /grafanaData",
			},
			ImagePullPolicy: corev1.PullAlways,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "grafana-data",
					MountPath: "/grafanaData",
				},
			},
			TerminationMessagePath:   terminationMessagePath,
			TerminationMessagePolicy: terminationMessagePolicy,
			Resources:                resources,
		},
	}

	image := cr.Spec.UserGrafana.ImageUserGrafana.Repository + ":" + cr.Spec.UserGrafana.ImageUserGrafana.Tag

	replicas := &cr.Spec.UserGrafana.Replica

	// var nodeSelectorMap map[string]string
	// err := json.Unmarshal([]byte(cr.Spec.UserGrafana.NodeSelector), &nodeSelectorMap)
	// if err != nil {
	// 	r.Log.Error(err, "Failed to parse nodeSelector")
	// 	return nil
	// }

	tolerations := cr.Spec.UserGrafana.Tolerations

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.UserGrafana.ImageUserGrafana.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.UserGrafana.ImageUserGrafana.ImagePullPolicy)
	}

	nodeSelector := cr.Spec.UserGrafana.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{} // Default to an empty map
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-user-grafana",          // Name
		labels,                        // Labels
		nil,                           // Command
		nil,                           // args (nil if not needed)
		volumeMounts,                  // Volume mounts
		volumes,                       // Volumes
		ports,                         // Ports (empty if not needed)
		initContainers,                // Init containers
		nil,                           // DNS config
		nil,                           // Annotations
		nil,                           // Environment variables
		livenessProbe,                 // Liveness probe
		readinessProbe,                // Readiness probe
		securityContext,               // Security context
		corev1.DNSClusterFirst,        // DNS policy
		corev1.RestartPolicyAlways,    // Restart policy
		imagePullSecrets,              // Image pull secrets
		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
		replicas,                      // replicas
		nodeSelector,                  // nodeSelector
		tolerations,                   // tolerations
		matchLabels,                   // matchlabels
	)
}

func (r *PmnsystemReconciler) orc8AlertManagerConfigurerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "alertmanager-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}
	matchLabels := map[string]string{
		"app.kubernetes.io/component": "alertmanager-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	replica := &cr.Spec.AlertmanagerConfigurer.Replica
	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "prometheus-config",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "promcfg",
				},
			},
		},
	}

	// VolumeMounts
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "prometheus-config",
			MountPath: "/etc/configs",
			ReadOnly:  true,
		},
	}

	// Affinity
	affinity := &corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app.kubernetes.io/component",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"prometheus"},
							},
						},
					},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{ContainerPort: 9101, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(9101),
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(9101),
			},
		},
	}

	args := []string{
		fmt.Sprintf("-port=%d", cr.Spec.AlertmanagerConfigurer.AlertManagerConfigPort),
		fmt.Sprintf("-alertmanager-conf=%s", cr.Spec.AlertmanagerConfigurer.AlertManagerConfPath),
		fmt.Sprintf("-alertmanagerURL=%s", cr.Spec.AlertmanagerConfigurer.AlertmanagerURL),
		"-multitenant-label=networkID",
		"-delete-route-with-receiver=true",
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

	image := cr.Spec.AlertmanagerConfigurer.ImageAlertmanagerConfigurer.Repository + ":" + cr.Spec.AlertmanagerConfigurer.ImageAlertmanagerConfigurer.Tag

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.AlertmanagerConfigurer.ImageAlertmanagerConfigurer.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.AlertmanagerConfigurer.ImageAlertmanagerConfigurer.ImagePullPolicy)
	}

	tolerations := cr.Spec.AlertmanagerConfigurer.Tolerations

	nodeSelector := cr.Spec.AlertmanagerConfigurer.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{} // Default to an empty map
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-alertmanager-configurer",
		labels,                        // Labels
		nil,                           // Command
		args,                          // args (nil if not needed)
		volumeMounts,                  // Volume mounts
		volumes,                       // Volumes
		ports,                         // Ports (empty if not needed)
		nil,                           // Init containers
		nil,                           // DNS config
		nil,                           // Annotations
		nil,                           // Environment variables
		livenessProbe,                 // Liveness probe
		readinessProbe,                // Readiness probe
		securityContext,               // Security context
		corev1.DNSClusterFirst,        // DNS policy
		corev1.RestartPolicyAlways,    // Restart policy
		imagePullSecrets,              // Image pull secrets
		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		affinity,                      // Affinity
		replica,                       // replicas
		nodeSelector,                  // Node Selector for AlertManagerConfigurer
		tolerations,                   // tolerations
		matchLabels,                   // match Labels
	)
}

func (r *PmnsystemReconciler) orc8AlertManagerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }

	labels := map[string]string{
		"app.kubernetes.io/component": "alertmanager",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}
	matchLabels := map[string]string{
		"app.kubernetes.io/component": "alertmanager",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "prometheus-config",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "promcfg",
				},
			},
		},
	}

	// VolumeMounts
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "prometheus-config",
			MountPath: "/etc/alertmanager",
			ReadOnly:  true,
		},
	}

	// Affinity
	affinity := &corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app.kubernetes.io/component",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"prometheus"},
							},
						},
					},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		},
	}

	// Define the securityContext for the container
	securityContext := &corev1.SecurityContext{
		Privileged: func(b bool) *bool { return &b }(true),
	}

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{ContainerPort: 9093, Protocol: corev1.ProtocolTCP},
	}

	// Liveness and Readiness Probes
	livenessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/",
				Port:   intstr.FromInt(9093),
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	readinessProbe := &corev1.Probe{
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		FailureThreshold:    3,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/",
				Port:   intstr.FromInt(9093),
				Scheme: corev1.URISchemeHTTP,
			},
		},
	}

	// Command for the container
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

	image := cr.Spec.AlertManager.ImageAlertmanager.Repository + ":" + cr.Spec.AlertManager.ImageAlertmanager.Tag

	replica := &cr.Spec.Prometheus.Replicas

	tolerations := cr.Spec.AlertManager.Tolerations

	nodeSelector := cr.Spec.AlertManager.NodeSelector
	if nodeSelector == nil {
		nodeSelector = map[string]string{} // Default to an empty map
	}

	// Fetch and validate ImagePullPolicy
	imagePullPolicy := corev1.PullIfNotPresent // Default value
	switch cr.Spec.AlertManager.ImageAlertmanager.ImagePullPolicy {
	case "Always":
		imagePullPolicy = corev1.PullAlways
	case "Never":
		imagePullPolicy = corev1.PullNever
	case "IfNotPresent":
		imagePullPolicy = corev1.PullIfNotPresent
	default:
		r.Log.Info("Invalid imagePullPolicy in CR, defaulting to IfNotPresent", "imagePullPolicy", cr.Spec.AlertManager.ImageAlertmanager.ImagePullPolicy)
	}

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-alertmanager",
		labels,                        // Labels
		nil,                           // Command
		nil,                           // args (nil if not needed)
		volumeMounts,                  // Volume mounts
		volumes,                       // Volumes
		ports,                         // Ports (empty if not needed)
		nil,                           // Init containers
		nil,                           // DNS config
		nil,                           // Annotations
		nil,                           // Environment variables
		livenessProbe,                 // Liveness probe
		readinessProbe,                // Readiness probe
		securityContext,               // Security context
		corev1.DNSClusterFirst,        // DNS policy
		corev1.RestartPolicyAlways,    // Restart policy
		imagePullSecrets,              // Image pull secrets
		terminationGracePeriodSeconds, // terminationGracePeriodSeconds
		imagePullPolicy,               // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		affinity,                      // Affinity
		replica,                       // replicas
		nodeSelector,                  // nodeSelector
		tolerations,                   // toleration
		matchLabels,                   // match labels
	)
}

func (r *PmnsystemReconciler) createOrc8rPrometheusStateFullSet(cr *v1.Pmnsystem) *appsv1.StatefulSet {
	int32Ptr := func(i int32) *int32 { return &i }

	storageSize := "50Gi"

	volumeMode := "Filesystem"

	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	matchLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	// Define VolumeClaimTemplates correctly
	volumeClaimTemplates := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "prometheus-data",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(storageSize),
					},
				},
				StorageClassName: &cr.Spec.PersistentForStatefulSet.StorageClassName,
				VolumeMode:       (*corev1.PersistentVolumeMode)(&volumeMode),
			},
		},
	}

	volumes := []corev1.Volume{
		{
			Name: "prometheus-config-file",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "prometheus-config-file",
					},
					DefaultMode: int32Ptr(420),
				},
			},
		},
		{
			Name: "orc8r-alert-rules",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "orc8r-alert-rules",
					},
					DefaultMode: int32Ptr(420),
				},
			},
		},
		{
			Name: "prometheus-config",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "promcfg",
				},
			},
		},
	}

	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "prometheus-config",
			MountPath: "/etc/prometheus",
			ReadOnly:  true,
		},
		{
			Name:      "prometheus-data",
			MountPath: "/data",
		},
		{
			Name:      "prometheus-config-file",
			MountPath: "/prometheus",
		},
		{
			Name:      "orc8r-alert-rules",
			MountPath: "/etc/orc8r_alerts",
		},
	}

	livenessProbe := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/graph",
				Port:   intstr.FromInt(9090),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		FailureThreshold:    3,
		SuccessThreshold:    1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		TimeoutSeconds:      1,
	}

	readinessProbe := &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/graph",
				Port:   intstr.FromInt(9090),
				Scheme: corev1.URISchemeHTTP,
			},
		},
		FailureThreshold:    3,
		SuccessThreshold:    1,
		InitialDelaySeconds: 10,
		PeriodSeconds:       30,
		TimeoutSeconds:      1,
	}

	resource := corev1.ResourceRequirements{}

	terminationMessagePath := "/dev/termination-log"

	terminationMessagePolicy := corev1.TerminationMessageReadFile

	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}
	dnsPolicy := corev1.DNSClusterFirst

	restartPolicy := corev1.RestartPolicyAlways

	args := []string{
		"--config.file=/prometheus/prometheus.yml",
		"--storage.tsdb.path=/data",
		"--web.enable-lifecycle",
		"--enable-feature=remote-write-receiver",
		"--web.enable-admin-api",
		"--query.timeout=1m",
		"--log.level=error",
		"--storage.tsdb.no-lockfile",
		"--storage.tsdb.retention.time=6h",
	}

	// Affinity
	affinity := &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"prometheus"},
							},
						},
					},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		},
	}

	updateStrategy := appsv1.StatefulSetUpdateStrategy{
		RollingUpdate: &appsv1.RollingUpdateStatefulSetStrategy{
			Partition: int32Ptr(0),
		},
		Type: appsv1.RollingUpdateStatefulSetStrategyType,
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-prometheus",
			Namespace: cr.Spec.NameSpace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &cr.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			ServiceName: "orc8r-prometheus",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:                     "prometheus",
							Image:                    "docker.io/prom/prometheus:v2.27.1",
							Args:                     args,
							Ports:                    []corev1.ContainerPort{{ContainerPort: 9090}},
							VolumeMounts:             volumeMounts,
							LivenessProbe:            livenessProbe,
							ReadinessProbe:           readinessProbe,
							TerminationMessagePath:   terminationMessagePath,
							TerminationMessagePolicy: terminationMessagePolicy,
							Resources:                resource,
						},
					},
					ImagePullSecrets: imagePullSecrets,
					DNSPolicy:        dnsPolicy,
					Volumes:          volumes,
					Affinity:         affinity,
					RestartPolicy:    restartPolicy,
				},
			},
			VolumeClaimTemplates: volumeClaimTemplates,
			UpdateStrategy:       updateStrategy,
		},
	}
}
func (r *PmnsystemReconciler) deploymentPostgres(cr *v1.Pmnsystem) *appsv1.Deployment {
	log := ctrl.Log.WithName("createPostgresResources")

	log.Info("DevEnvironment is true. Creating PostgreSQL Deployment...")

	// Define the PostgreSQL Deployment
	postgresDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgres",
			Namespace: cr.Spec.NameSpace,
			Labels: map[string]string{
				"app": "postgres",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "postgres",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "postgres",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "postgres",
							Image: "postgres:14",
							Env: []corev1.EnvVar{
								{
									Name:  "POSTGRES_DB",
									Value: "pmndev",
								},
								{
									Name:  "POSTGRES_USER",
									Value: "pmn",
								},
								{
									Name:  "POSTGRES_PASSWORD",
									Value: "juniperprod1234",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 5432,
								},
							},
						},
					},
				},
			},
		},
	}

	log.Info("PostgreSQL Deployment created successfully")
	return postgresDeployment

}
