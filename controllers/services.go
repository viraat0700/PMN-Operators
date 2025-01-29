package controllers

import (
	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *PmnsystemReconciler) orc8rAccessDService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabel := map[string]string{
		"app.kubernetes.io/component": "accessd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-accessd",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabel,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9091),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9191),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rAnalyticsService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "analytics",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-analytics",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9200),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9300),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rBootStrapperService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "bootstrapper",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-bootstrapper",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9088),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9188),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rCertifierService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component":  "orc8r",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/name":       "orc8r",
		"app.kubernetes.io/part-of":    "orc8r-app",
		"orc8r.io/analytics_collector": "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "certifier",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-certifier",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9086),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9186),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10089),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rConfiguratorService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "configurator",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-configurator",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9108),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9208),
				},
				{
					Name:       "moso",
					Port:       8088,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(8088),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rDeviceService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "device",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-device",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9106),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9306),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rDirectoryDService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "directoryd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-directoryd",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9100),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9102),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rDispatcherService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "dispatcher",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-dispatcher",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9096),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9196),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rEventdService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
		"orc8r.io/obsidian_handlers":  "true",
		"orc8r.io/swagger_spec":       "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "eventd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-eventd",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9121),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9221),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10121),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rmetricsdService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
		"orc8r.io/obsidian_handlers":  "true",
		"orc8r.io/swagger_spec":       "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "metricsd",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-metricsd",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9084),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9184),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10084),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rNginxProxyService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "nginx-proxy",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "nginx-proxy",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-nginx-proxy",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "health",
					Port:       80,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(80),
				},
				{
					Name:       "clientcert",
					Port:       8443,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(8443),
				},
				{
					Name:       "open",
					Port:       8444,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(8444),
				},
				{
					Name:       "api",
					Port:       443,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9443),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rNotifierService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component":  "notifier",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "notifier",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.Orc8rNotifier.ServiceSpecOrc8rNotifier.PortSpecOrc8rNotifier {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
			NodePort:   port.NodePort,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-notifier",
			Namespace: cr.Spec.NameSpace,
			Labels:    selectorLabels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.Orc8rNotifier.ServiceSpecOrc8rNotifier.Type),
			Selector:        labels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rNotifierInternalService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "notifier",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "notifier",
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-notifier-internal",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "notifier-internal",
					Port:       5442,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(5442),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rObsidianService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "obsidian",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-obsidian",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9093),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9193),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9081),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rWorkerService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "orc8r-worker",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-orc8r-worker",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9122),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9222),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rOrchestratorService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component":  "orc8r",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/name":       "orc8r",
		"app.kubernetes.io/part-of":    "orc8r-app",
		"orc8r.io/analytics_collector": "true",
		"orc8r.io/mconfig_builder":     "true",
		"orc8r.io/metrics_exporter":    "true",
		"orc8r.io/obsidian_handlers":   "true",
		"orc8r.io/state_indexer":       "true",
		"orc8r.io/stream_provider":     "true",
		"orc8r.io/swagger_spec":        "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "orchestrator",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-orchestrator",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9112),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9212),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10112),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rServiceRegistryService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "service_registry",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-service-registry",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9180),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9190),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rStateService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	matchlabels := map[string]string{
		"app.kubernetes.io/component": "state",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-state",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: matchlabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9105),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9305),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rStreamerService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "streamer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-streamer",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9082),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9182),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rTenantsService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "orc8r",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
		"orc8r.io/obsidian_handlers":  "true",
		"orc8r.io/swagger_spec":       "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "tenants",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-tenants",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9110),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9210),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10110),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rHaService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "lte-orc8r",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "ha",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-ha",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9119),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rLteService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component":  "lte-orc8r",
		"app.kubernetes.io/instance":   "lte-pmn",
		"app.kubernetes.io/name":       "lte-orc8r",
		"app.kubernetes.io/part-of":    "orc8r-app",
		"orc8r.io/analytics_collector": "true",
		"orc8r.io/mconfig_builder":     "true",
		"orc8r.io/obsidian_handlers":   "true",
		"orc8r.io/state_indexer":       "true",
		"orc8r.io/stream_provider":     "true",
		"orc8r.io/swagger_spec":        "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "lte",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-lte",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9113),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9213),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10113),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rNprobeService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "lte-orc8r",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
		"orc8r.io/obsidian_handlers":  "true",
		"orc8r.io/swagger_spec":       "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "nprobe",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-nprobe",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9666),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9766),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10088),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rPolicyDbService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "lte-orc8r",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
		"orc8r.io/obsidian_handlers":  "true",
		"orc8r.io/swagger_spec":       "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "policydb",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-policydb",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9085),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9185),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10085),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rSmsdService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "lte-orc8r",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
		"orc8r.io/obsidian_handlers":  "true",
		"orc8r.io/swagger_spec":       "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "smsd",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-smsd",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9120),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9220),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10086),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rSubscriberDbCacheService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "lte-orc8r",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "subscriberdb-cache",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-subscriberdb-cache",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9089),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10087),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rSubscriberDbService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "lte-orc8r",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
		"app.kubernetes.io/part-of":   "orc8r-app",
		"orc8r.io/obsidian_handlers":  "true",
		"orc8r.io/state_indexer":      "true",
		"orc8r.io/swagger_spec":       "true",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "subscriberdb",
		"app.kubernetes.io/instance":  "lte-pmn",
		"app.kubernetes.io/name":      "lte-orc8r",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-subscriberdb",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9083),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9183),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10083),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) NmsMagmaLteService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "magmalte",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "nms",
		"app.kubernetes.io/part-of":   "magma",
		"release_group":               "orc8r",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "magmalte",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "nms",
		"release_group":               "orc8r",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.NmsMagmaLte.ServiceSpec.PortSpec {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nms-magmalte",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.NmsMagmaLte.ServiceSpec.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rPrometheusCacheService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus-cache",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus-cache",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.PrometheusCache.ServiceSpec.PortSpec {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-prometheus-cache",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.PrometheusCache.ServiceSpec.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rPrometheusConfigurerService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.PrometheusConfigurer.ServiceSpec.PortSpec {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-prometheus-configurer",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.PrometheusConfigurer.ServiceSpec.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rPrometheusKafkaAdapterService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/instance": "orc8r",
		"app.kubernetes.io/name":     "prometheus-kafka-adapter",
		"app.kubernetes.io/version":  "1.0",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/instance": "orc8r",
		"app.kubernetes.io/name":     "prometheus-kafka-adapter",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.PrometheusKafkaAdapter.ServiceSpecPrometheusKafkaAdapter.PortSpecPrometheusKafkaAdapter {
		targetPort := intstr.FromString(port.TargetPort)
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: targetPort,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-prometheus-kafka-adapter",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.PrometheusKafkaAdapter.ServiceSpecPrometheusKafkaAdapter.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rPrometheusNginxProxyService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus-nginx",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus-nginx",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.PrometheusNginxProxy.Nginx.ServiceOrc8rSpec.PortOrc8rSpec {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
			NodePort:   port.NodePort,
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-prometheus-nginx-proxy",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.PrometheusNginxProxy.Nginx.ServiceOrc8rSpec.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rUserGrafanaService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "user-grafana",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}
	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "user-grafana",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.UserGrafana.ServiceSpec.PortSpec {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-user-grafana",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.UserGrafana.ServiceSpec.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rAlertManagerConfigurerService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "alertmanager-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}
	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "alertmanager-configurer",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.AlertmanagerConfigurer.ServiceSpec.PortSpec {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
		})
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-alertmanager-configurer",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.AlertmanagerConfigurer.ServiceSpec.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) orc8rAlterManagerService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "alertmanager",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "alertmanager",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	var servicePorts []corev1.ServicePort
	for _, port := range cr.Spec.AlertManager.ServiceSpec.PortSpec {
		servicePorts = append(servicePorts, corev1.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			Protocol:   corev1.Protocol(port.Protocol),
			TargetPort: intstr.FromInt(int(port.TargetPort)),
		})
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-alertmanager",
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
		Spec: corev1.ServiceSpec{
			Type:            corev1.ServiceType(cr.Spec.AlertManager.ServiceSpec.Type),
			Selector:        selectorLabels,
			Ports:           servicePorts,
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}
func (r *PmnsystemReconciler) servicePostgres(cr *v1.Pmnsystem) *corev1.Service {
	log := ctrl.Log.WithName("createPostgresResources")
	log.Info("Creating PostgreSQL Service...")

	postgresService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgres",
			Namespace: cr.Spec.NameSpace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": "postgres",
			},
			Ports: []corev1.ServicePort{
				{
					Port:       5432,
					TargetPort: intstr.FromInt(5432),
					NodePort:   30000,
				},
			},
		},
	}

	log.Info("PostgreSQL Service created successfully")
	return postgresService
}
func (r *PmnsystemReconciler) orc8rPrometheusService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component": "prometheus",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
		"app.kubernetes.io/version":   "1.0",
	}

	selectorLabels := map[string]string{
		"app.kubernetes.io/component": "prometheus",
		"app.kubernetes.io/instance":  "orc8r",
		"app.kubernetes.io/name":      "metrics",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-prometheus",
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
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: selectorLabels,
			Ports: []corev1.ServicePort{
				{
					Name:       "prometheus",
					Port:       9090,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9090),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
} // service for statefulset
