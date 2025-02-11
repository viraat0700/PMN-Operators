package controllers

import (
	"context"

	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
)

func (r *PmnsystemReconciler) orc8rIngress(cr *v1.Pmnsystem) *networkingv1.Ingress {
	annotations := map[string]string{
		"nginx.ingress.kubernetes.io/backend-protocol":   "HTTPS",
		"nginx.ingress.kubernetes.io/force-ssl-redirect": "true",
		"nginx.ingress.kubernetes.io/ssl-passthrough":    "true",
		"nginx.ingress.kubernetes.io/ssl-redirect":       "true",
	}

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "orc8r-passthrough-ingress",
			Namespace:   cr.Spec.NameSpace,
			Annotations: annotations,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ptr.To("nginx"),
			Rules: []networkingv1.IngressRule{
				{
					Host: "bootstrapper-controller.operator.wavelabs.int",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: ptr.To(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "orc8r-bootstrap-nginx",
											Port: networkingv1.ServiceBackendPort{Number: 443},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "api.operator.wavelabs.int",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: ptr.To(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "orc8r-nginx-proxy",
											Port: networkingv1.ServiceBackendPort{Number: 443},
										},
									},
								},
							},
						},
					},
				},
				{
					Host: "controller.operator.wavelabs.int",
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: ptr.To(networkingv1.PathTypePrefix),
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "orc8r-clientcert-nginx",
											Port: networkingv1.ServiceBackendPort{Number: 443},
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

	return ingress
}

func (r *PmnsystemReconciler) ensureIngress(ctx context.Context, cr *v1.Pmnsystem) error {
	log := r.Log.WithValues("Ingress", cr.Name)

	// Check if Ingress already exists
	found := &networkingv1.Ingress{}
	err := r.Get(ctx, types.NamespacedName{Name: "orc8r-passthrough-ingress", Namespace: cr.Spec.NameSpace}, found)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Ingress not found, creating a new one.")
		ingress := r.orc8rIngress(cr) // Generate the Ingress manifest

		if err := r.Create(ctx, ingress); err != nil {
			log.Error(err, "Failed to create Ingress")
			return err
		}
		log.Info("Successfully created Ingress.")
		return nil
	} else if err != nil {
		log.Error(err, "Failed to get Ingress")
		return err
	}

	// Ingress exists, nothing to do
	log.Info("Ingress already exists.")
	return nil
}
