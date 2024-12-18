package controllers

import (
	v1 "github.com/viraat0700/PMN-Operator-Two/api/v1alpha1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *PmnsystemReconciler) orc8rAccessDPDB(cr *v1.Pmnsystem) *policyv1.PodDisruptionBudget {
	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy/v1",
			Kind:       "PodDisruptionBudget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-accessd",
			Namespace: cr.Spec.NameSpace,
			Annotations: map[string]string{
				"app":                          "orc8r-accessd",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
			Labels: map[string]string{
				"app":                          "orc8r-accessd",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: func(i int) *intstr.IntOrString {
				v := intstr.FromInt(i)
				return &v
			}(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":                          "orc8r-accessd",
					"app.kubernetes.io/instance":   "orc8r",
					"app.kubernetes.io/managed-by": "Orc8r-Operator",
				},
			},
		},
	}
}
func (r *PmnsystemReconciler) orc8rAnalyticsDPDB(cr *v1.Pmnsystem) *policyv1.PodDisruptionBudget {
	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy/v1",
			Kind:       "PodDisruptionBudget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-analytics",
			Namespace: cr.Spec.NameSpace,
			Annotations: map[string]string{
				"app":                          "orc8r-analytics",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
			Labels: map[string]string{
				"app":                          "orc8r-analytics",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: func(i int) *intstr.IntOrString {
				v := intstr.FromInt(i)
				return &v
			}(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":                          "orc8r-analytics",
					"app.kubernetes.io/instance":   "orc8r",
					"app.kubernetes.io/managed-by": "Orc8r-Operator",
				},
			},
		},
	}
}
func (r *PmnsystemReconciler) orc8rBootstrapperPDB(cr *v1.Pmnsystem) *policyv1.PodDisruptionBudget {
	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "policy/v1",
			Kind:       "PodDisruptionBudget",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "orc8r-bootstrapper",
			Namespace: cr.Spec.NameSpace,
			Annotations: map[string]string{
				"app":                          "orc8r-bootstrapper",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
			Labels: map[string]string{
				"app":                          "orc8r-bootstrapper",
				"app.kubernetes.io/instance":   "orc8r",
				"app.kubernetes.io/managed-by": "Orc8r-Operator",
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1.GroupVersion.Group,
					Version: v1.GroupVersion.Version,
					Kind:    "Pmnsystem",
				}),
			},
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: func(i int) *intstr.IntOrString {
				v := intstr.FromInt(i)
				return &v
			}(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":                          "orc8r-bootstrapper",
					"app.kubernetes.io/instance":   "orc8r",
					"app.kubernetes.io/managed-by": "Orc8r-Operator",
				},
			},
		},
	}
}
