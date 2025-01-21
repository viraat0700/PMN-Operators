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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (r *PmnsystemReconciler) createPersistentVolumeClaim(cr *v1.Pmnsystem) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app.kubernetes.io/component": "promcfg",
	}

	storageclassname := cr.Spec.Persistent.StorageClassName

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Persistent.PvcClaimName,
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
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &storageclassname,
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}
func (r *PmnsystemReconciler) createPersistentVolumeClaimPromData(cr *v1.Pmnsystem) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app.kubernetes.io/component": "promdata",
	}

	storageclassname := cr.Spec.Persistent.StorageClassName

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "promdata",
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
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &storageclassname,
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}
func (r *PmnsystemReconciler) createPersistentVolumeClaimGrafanaProviders(cr *v1.Pmnsystem) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app.kubernetes.io/component": "grafanaproviders",
	}

	storageclassname := cr.Spec.Persistent.StorageClassName

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafanaproviders",
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
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &storageclassname,
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}
func (r *PmnsystemReconciler) createPersistentVolumeClaimGrafanaDatasources(cr *v1.Pmnsystem) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app.kubernetes.io/component": "grafanadatasources",
	}

	storageclassname := cr.Spec.Persistent.StorageClassName

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafanadatasources",
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
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &storageclassname,
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}
func (r *PmnsystemReconciler) createPersistentVolumeClaimGrafanaData(cr *v1.Pmnsystem) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app.kubernetes.io/component": "grafanadata",
	}

	storageclassname := cr.Spec.Persistent.StorageClassName

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafanadata",
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
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &storageclassname,
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}
func (r *PmnsystemReconciler) createPersistentVolumeClaimGrafanaDashboards(cr *v1.Pmnsystem) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app.kubernetes.io/component": "grafanadashboards",
	}

	storageclassname := cr.Spec.Persistent.StorageClassName

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafanadashboards",
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
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &storageclassname,
		},
	}

	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}
