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
func (r *PmnsystemReconciler) orc8rmetricsdService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app":                          "orc8r-metricsd",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
		"orc8r.io/obsidian_handlers":   "true",
		"orc8r.io/swagger_spec":        "true",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-metricsd",
			Namespace: cr.Spec.NameSpace,
			Labels:    labels,
			Annotations: map[string]string{
				"app":                          "orc8r-metricsd",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
				"orc8r.io/obsidian_handlers_path_prefixes": "/magma/v1/networks/:network_id/alerts, /magma/v1/networks/:network_id/metrics, /magma/v1/networks/:network_id/prometheus, /magma/v1/tenants/:tenant_id/metrics, /magma/v1/tenants/targets_metadata,",
			},
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
func (r *PmnsystemReconciler) orc8rNotifierService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app":                          "orc8r-notifier",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-notifier",
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
func (r *PmnsystemReconciler) orc8rNotifierInternalService(cr *v1.Pmnsystem) *corev1.Service {
	labels := map[string]string{
		"app":                          "orc8r-notifier-internal",
		"app.kubernetes.io/instance":   "orc8r",
		"app.kubernetes.io/managed-by": "Orc8r-Operator",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-notifier-internal",
			Namespace: cr.Spec.NameSpace,
			Labels:    labels,
			Annotations: map[string]string{
				"app":                          "orc8r-notifier-internal",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
				"external-dns.alpha.kubernetes.io/hostname":                             "notifier.pmn-dev.wavelabs.in",
				"kubectl.kubernetes.io/last-applied-configuration":                      `{"apiVersion":"v1","kind":"Service","metadata":{"annotations":{"external-dns.alpha.kubernetes.io/hostname":"notifier.pmn-dev.wavelabs.in","meta.helm.sh/release-name":"orc8r","meta.helm.sh/release-namespace":"pmn","service.beta.kubernetes.io/aws-load-balancer-additional-resource-tags":"lbtype=service","service.beta.kubernetes.io/aws-load-balancer-type":"nlb"},"labels":{"app.kubernetes.io/component":"notifier","app.kubernetes.io/managed-by":"Helm"},"name":"notifier","namespace":"pmn"},"spec":{"allocateLoadBalancerNodePorts":true,"clusterIP":"172.20.227.67","clusterIPs":["172.20.227.67"],"externalTrafficPolicy":"Cluster","internalTrafficPolicy":"Cluster","ipFamilies":["IPv4"],"ipFamilyPolicy":"SingleStack","ports":[{"name":"notifier","nodePort":32001,"port":443,"protocol":"TCP","targetPort":443}],"selector":{"app.kubernetes.io/component":"notifier"},"sessionAffinity":"None","type":"LoadBalancer"},"status":{"loadBalancer":{"ingress":[{"hostname":"ae7e0c4a79e194bf8b11966c881be819-cf7367a52c2eed60.elb.us-west-2.amazonaws.com"}]}}}`,
				"service.beta.kubernetes.io/aws-load-balancer-additional-resource-tags": "lbtype=service",
				"service.beta.kubernetes.io/aws-load-balancer-type":                     "nlb",
			},
			Finalizers: []string{
				"service.kubernetes.io/load-balancer-cleanup",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeLoadBalancer,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "notifier",
					NodePort:   32001,
					Port:       4443,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(443),
				},
			},
			SessionAffinity: corev1.ServiceAffinityNone,
		},
		Status: corev1.ServiceStatus{
			LoadBalancer: corev1.LoadBalancerStatus{
				Ingress: []corev1.LoadBalancerIngress{
					{
						Hostname: "a4ca190bb09f048a19690cc67ab7038f-7f3c2eb6df85f47a.elb.us-west-2.amazonaws.com",
					},
				},
			},
		},
	}
}
