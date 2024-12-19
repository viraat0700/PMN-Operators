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
	image string,
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

	image := cr.Spec.Image.Repository + ":" + cr.Spec.Image.Tag

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
	)
}