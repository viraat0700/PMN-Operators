package controllers

import (
	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *PmnsystemReconciler) orc8rAccessDService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app":                          "orc8r-accessd",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-analytics",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-bootstrapper",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-certifier",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-configurator",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-device",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-directoryd",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-dispatcher",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
		"app":                          "orc8r-eventd",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
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
			Selector: labels,
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
func (r *PmnsystemReconciler) orc8rDomainProxyService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app":                          "orc8r-domain-proxy",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-domain-proxy",
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
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       9180,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9124),
				},
				{
					Name:       "grpc-internal",
					Port:       9190,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(9224),
				},
				{
					Name:       "http",
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(10124),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
}