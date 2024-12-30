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
	"k8s.io/apimachinery/pkg/api/resource"
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
	image string,
	affinity *corev1.Affinity,
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
					Affinity: affinity,
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
	)
}
func (r *PmnsystemReconciler) orc8rAnalyticsDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-analytics",
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
	)
}
func (r *PmnsystemReconciler) orc8rBootStrapperDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-bootstrapper",
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
	)
}
func (r *PmnsystemReconciler) orc8rCertifierDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-certifier",
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
	)
}
func (r *PmnsystemReconciler) orc8rConfiguratorDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-configurator",
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

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-configurator",
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
	)
}
func (r *PmnsystemReconciler) orc8rDeviceDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-device",
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
	)
}
func (r *PmnsystemReconciler) orc8rDirectorydDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-directoryd",
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
	)
}
func (r *PmnsystemReconciler) orc8rDispatcherDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-dispatcher",
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
	)
}
func (r *PmnsystemReconciler) orc8rEventdDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-eventd",
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
	)
}
func (r *PmnsystemReconciler) orc8rmetricsdDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-metricsd",
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
	)
}
func (r *PmnsystemReconciler) orc8rNginxDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-nginx",
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
	}

	// Define volumeMounts in a separate variable
	volumeMounts := []corev1.VolumeMount{
		{Name: "certs", MountPath: "/var/opt/magma/certs", ReadOnly: true},
		{Name: "envdir", MountPath: "/var/opt/magma/envdir", ReadOnly: true},
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
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForOrc8rNginx(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "clientcert", ContainerPort: 8443, Protocol: corev1.ProtocolTCP},
		{Name: "open", ContainerPort: 8444, Protocol: corev1.ProtocolTCP},
		{Name: "api", ContainerPort: 9443, Protocol: corev1.ProtocolTCP},
		{Name: "health", ContainerPort: 80, Protocol: corev1.ProtocolTCP},
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

	// Command for the container
	command := []string{}

	args := []string{}

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

	image := cr.Spec.NginxImage

	return r.deployment(
		strategy, // Deployment strategy
		cr,
		"orc8r-nginx",
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
	)
}
func (r *PmnsystemReconciler) orc8rNotifierDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-notifier",
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
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	envVars := r.getEnvVarsForOrc8rNotifier(cr)

	// Define ports (use nil if not needed)
	ports := []corev1.ContainerPort{
		{Name: "notifier", ContainerPort: 443, Protocol: corev1.ProtocolTCP},
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
					IntVal: 5442,
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
					IntVal: 5442,
				},
			},
		},
	}

	// Command for the container
	// command := []string{
	// 	"/usr/bin/envdir",
	// }

	args := []string{
		"sh",
		"-c",
		"java -jar Orc8rNotificationService-1.0-SNAPSHOT.jar \"5442\" \"443\"",
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

	image := cr.Spec.NotifierImage

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
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
	)
}
func (r *PmnsystemReconciler) orc8rObsidianDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-obsidian",
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
	)
}
func (r *PmnsystemReconciler) orc8WorkerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-orc8r-worker",
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
	)
}
func (r *PmnsystemReconciler) orc8orchestratorDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-orchestrator",
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
	)
}
func (r *PmnsystemReconciler) orc8ServiceRegistryDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-service-registry",
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
	)
}
func (r *PmnsystemReconciler) orc8StateDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-state",
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
	)
}
func (r *PmnsystemReconciler) orc8StreamerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-streamer",
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
	)
}
func (r *PmnsystemReconciler) orc8TenantsDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-tenants",
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
	)
}
func (r *PmnsystemReconciler) orc8rHaDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-ha",
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
	)
}
func (r *PmnsystemReconciler) orc8LteDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-lte",
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
		nil,
	)
}
func (r *PmnsystemReconciler) orc8NprobeDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-nprobe",
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
	)
}
func (r *PmnsystemReconciler) orc8PolicyDbDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-policydb",
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
		nil,
	)
}
func (r *PmnsystemReconciler) orc8SmsdDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-smsd",
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
		nil,
	)
}
func (r *PmnsystemReconciler) orc8SubscriberDbCacheDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-subscriberdb-cache",
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
	)
}
func (r *PmnsystemReconciler) orc8SubscriberDbDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-subscriberdb",
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
		nil,
	)
}
func (r *PmnsystemReconciler) nmsMagmaLteDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "nms-magmalte",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "orc8r-secrets-certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "nms-certs",
					DefaultMode: int32Ptr(292),
				},
			},
		},
	}

	// Define volumeMounts in a separate variable
	volumeMounts := []corev1.VolumeMount{
		{Name: "orc8r-secrets-certs", MountPath: "/run/secrets/admin_operator.pem", SubPath: "admin_operator.pem"},
		{Name: "orc8r-secrets-certs", MountPath: "/run/secrets/admin_operator.key.pem", SubPath: "admin_operator.key.pem"},
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

	// Command for the container
	// command := []string{
	// 	"/usr/bin/envdir",
	// }

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

	image := cr.Spec.ImageMagmaLte

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
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
	)
}
func (r *PmnsystemReconciler) orc8AlertManagerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	// int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-alertmanager",
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

	// If Bevo is true, add the NET_ADMIN capability
	// if cr.Spec.Bevo {
	// 	securityContext.Capabilities = &corev1.Capabilities{
	// 		Add: []corev1.Capability{"NET_ADMIN"},
	// 	}
	// }

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	// envVars := r.getEnvVarsForDirectoryD(cr)

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
	// command := []string{
	// 	"/usr/bin/envdir",
	// }

	// args := []string{
	// 	"/var/opt/magma/envdir",
	// 	"/var/opt/magma/bin/subscriberdb",
	// 	"-run_echo_server=true",
	// 	"-logtostderr=true",
	// 	"-v=0",
	// }

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

	image := "docker.io/prom/alertmanager:v0.18.0"

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
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		affinity,                      // Affinity
	)
}
func (r *PmnsystemReconciler) orc8PrometheusCacheDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	// int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-prometheus-cache",
	}

	// Define volumes in a separate variable
	// volumes := []corev1.Volume{
	// 	{
	// 		Name: "prometheus-config",
	// 		VolumeSource: corev1.VolumeSource{
	// 			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	// 				ClaimName: "promcfg",
	// 			},
	// 		},
	// 	},
	// }

	// VolumeMounts
	// volumeMounts := []corev1.VolumeMount{
	// 	{
	// 		Name:      "prometheus-config",
	// 		MountPath: "/etc/alertmanager",
	// 		ReadOnly:  true,
	// 	},
	// }

	// Affinity
	// affinity := &corev1.Affinity{
	// 	PodAffinity: &corev1.PodAffinity{
	// 		RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
	// 			{
	// 				LabelSelector: &metav1.LabelSelector{
	// 					MatchExpressions: []metav1.LabelSelectorRequirement{
	// 						{
	// 							Key:      "app.kubernetes.io/component",
	// 							Operator: metav1.LabelSelectorOpIn,
	// 							Values:   []string{"prometheus"},
	// 						},
	// 					},
	// 				},
	// 				TopologyKey: "kubernetes.io/hostname",
	// 			},
	// 		},
	// 	},
	// }

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
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	// envVars := r.getEnvVarsForDirectoryD(cr)

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

	// Command for the container
	// command := []string{
	// 	"/usr/bin/envdir",
	// }

	args := []string{
		"-limit=500000",
		"-grpc-port=9092",
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

	image := "docker.io/facebookincubator/prometheus-edge-hub:1.1.0"

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
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
	)
}
func (r *PmnsystemReconciler) orc8rPrometheusConfigurerDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	// int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-prometheus-configurer",
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

	// If Bevo is true, add the NET_ADMIN capability
	// if cr.Spec.Bevo {
	// 	securityContext.Capabilities = &corev1.Capabilities{
	// 		Add: []corev1.Capability{"NET_ADMIN"},
	// 	}
	// }

	// Define imagePullSecrets
	imagePullSecrets := []corev1.LocalObjectReference{
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	// envVars := r.getEnvVarsForDirectoryD(cr)

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

	// Command for the container
	// command := []string{
	// 	"/usr/bin/envdir",
	// }

	args := []string{
		"-port=9100",
		"-rules-dir=/etc/configs/alert_rules/",
		"-prometheusURL=orc8r-prometheus:9090",
		"-multitenant-label=networkID",
		"-restrict-queries",
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

	image := "docker.io/facebookincubator/prometheus-configurer:1.0.4"

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
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		affinity,                      // Affinity
	)
}
func (r *PmnsystemReconciler) orc8rPrometheusKafkaAdapterDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-prometheus-kafka-adapter",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "ssl-client-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "prometheus-adapter-certs",
					DefaultMode: int32Ptr(420),
				},
			},
		},
	}

	// VolumeMounts
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "ssl-client-cert",
			MountPath: "/client_cert",
			ReadOnly:  true,
		},
	}

	// Affinity
	// affinity := &corev1.Affinity{
	// 	PodAffinity: &corev1.PodAffinity{
	// 		RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
	// 			{
	// 				LabelSelector: &metav1.LabelSelector{
	// 					MatchExpressions: []metav1.LabelSelectorRequirement{
	// 						{
	// 							Key:      "app.kubernetes.io/component",
	// 							Operator: metav1.LabelSelectorOpIn,
	// 							Values:   []string{"prometheus"},
	// 						},
	// 					},
	// 				},
	// 				TopologyKey: "kubernetes.io/hostname",
	// 			},
	// 		},
	// 	},
	// }

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
	// imagePullSecrets := []corev1.LocalObjectReference{
	// 	{Name: cr.Spec.ImagePullSecrets},
	// }

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

	// Command for the container
	// command := []string{
	// 	"/usr/bin/envdir",
	// }

	// args := []string{
	// 	"-port=9100",
	// 	"-rules-dir=/etc/configs/alert_rules/",
	// 	"-prometheusURL=orc8r-prometheus:9090",
	// 	"-multitenant-label=networkID",
	// 	"-restrict-queries",
	// }

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

	image := "telefonica/prometheus-kafka-adapter:1.9.1"

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
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
	)
}
func (r *PmnsystemReconciler) orc8rPrometheusNginxProxyDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-prometheus-nginx-proxy",
	}

	// Define volumes in a separate variable
	volumes := []corev1.Volume{
		{
			Name: "prometheus-certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  "prometheus-certs",
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
			MountPath: "/etc/nginx/conf.d/nginx_prometheus_ssl.conf",
			SubPath:   "nginx_prometheus_ssl.conf",
		},
		{
			Name:      "prometheus-certs",
			MountPath: "/etc/nginx/conf.d/prometheus.crt",
			ReadOnly:  true,
			SubPath:   "prometheus.crt",
		},
		{
			Name:      "prometheus-certs",
			MountPath: "/etc/nginx/conf.d/prometheus.key",
			SubPath:   "prometheus.key",
		},
		{
			Name:      "prometheus-certs",
			MountPath: "/etc/nginx/conf.d/prometheus-ca.crt",
			SubPath:   "prometheus-ca.crt",
			ReadOnly:  true,
		},
	}

	// Affinity
	// affinity := &corev1.Affinity{
	// 	PodAffinity: &corev1.PodAffinity{
	// 		RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
	// 			{
	// 				LabelSelector: &metav1.LabelSelector{
	// 					MatchExpressions: []metav1.LabelSelectorRequirement{
	// 						{
	// 							Key:      "app.kubernetes.io/component",
	// 							Operator: metav1.LabelSelectorOpIn,
	// 							Values:   []string{"prometheus"},
	// 						},
	// 					},
	// 				},
	// 				TopologyKey: "kubernetes.io/hostname",
	// 			},
	// 		},
	// 	},
	// }

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
	// imagePullSecrets := []corev1.LocalObjectReference{
	// 	{Name: cr.Spec.ImagePullSecrets},
	// }

	// Define environment variables if needed
	// envVars := r.getEnvVarsForPrometheusKafkaAdapter(cr)

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

	// Command for the container
	// command := []string{
	// 	"/usr/bin/envdir",
	// }

	// args := []string{
	// 	"-port=9100",
	// 	"-rules-dir=/etc/configs/alert_rules/",
	// 	"-prometheusURL=orc8r-prometheus:9090",
	// 	"-multitenant-label=networkID",
	// 	"-restrict-queries",
	// }

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

	image := "nginx:latest"

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
	)
}
func (r *PmnsystemReconciler) orc8rUserGrafanaDeployment(cr *v1.Pmnsystem) *appsv1.Deployment {
	int64Ptr := func(i int64) *int64 { return &i }
	int32Ptr := func(i int32) *int32 { return &i }

	labels := map[string]string{
		"app": "orc8r-user-grafana",
	}

	// Define volumes in a separate variable
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
					ClaimName: "grafanadatasources",
				},
			},
		},
		{
			Name: "dashboardproviders",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "grafanaproviders",
				},
			},
		},
		{
			Name: "dashboards",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "grafanadashboards",
				},
			},
		},
		{
			Name: "grafana-data",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "grafanadata",
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
			MountPath: "/var/lib/grafana/dashboards",
		},
		{
			Name:      "datasources",
			MountPath: "/etc/grafana/provisioning/datasources/",
		},
		{
			Name:      "dashboardproviders",
			MountPath: "/etc/grafana/provisioning/dashboards/",
		},
		{
			Name:      "grafana-data",
			MountPath: "/var/lib/grafana",
		},
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
		{Name: cr.Spec.ImagePullSecrets},
	}

	// Define environment variables if needed
	// envVars := r.getEnvVarsForPrometheusKafkaAdapter(cr)

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

	image := "docker.io/grafana/grafana:6.6.2"

	// Define initContainers in a separate variable
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
		corev1.PullIfNotPresent,       // Image pull policy
		resources,                     // Resources
		terminationMessagePath,        // Termination message path
		terminationMessagePolicy,      // Termination message policy
		image,                         // Image
		nil,                           // Affinity
	)
}
func (r *PmnsystemReconciler) createOrc8rPrometheusStateFullSet(cr *v1.Pmnsystem) *appsv1.StatefulSet {
	int32Ptr := func(i int32) *int32 { return &i }

	storageSize := "50Gi"

	volumeMode := "Filesystem"

	labels := map[string]string{
		"app":                          "orc8r-prometheus",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
	}

	annotations := map[string]string{
		"app":                          "orc8r-prometheus",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
	}


	// Define VolumeClaimTemplates correctly
	volumeClaimTemplates := []corev1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "prometheus-data",
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
				StorageClassName: &cr.Spec.Persistent.StorageClassName,
				VolumeMode: (*corev1.PersistentVolumeMode)(&volumeMode),
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
	// Define imagePullSecrets
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
			Name:        "orc8r-prometheus",
			Namespace:   cr.Spec.NameSpace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &cr.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: "orc8r-prometheus",
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
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
			UpdateStrategy: updateStrategy,
		},
	}
}
