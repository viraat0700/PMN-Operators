//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariables) DeepCopyInto(out *EnvironmentVariables) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariables.
func (in *EnvironmentVariables) DeepCopy() *EnvironmentVariables {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariables)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariablesDirectoryD) DeepCopyInto(out *EnvironmentVariablesDirectoryD) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariablesDirectoryD.
func (in *EnvironmentVariablesDirectoryD) DeepCopy() *EnvironmentVariablesDirectoryD {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariablesDirectoryD)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariablesNMSMagmaLte) DeepCopyInto(out *EnvironmentVariablesNMSMagmaLte) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariablesNMSMagmaLte.
func (in *EnvironmentVariablesNMSMagmaLte) DeepCopy() *EnvironmentVariablesNMSMagmaLte {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariablesNMSMagmaLte)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariablesNMSMagmaLteDevEnv) DeepCopyInto(out *EnvironmentVariablesNMSMagmaLteDevEnv) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariablesNMSMagmaLteDevEnv.
func (in *EnvironmentVariablesNMSMagmaLteDevEnv) DeepCopy() *EnvironmentVariablesNMSMagmaLteDevEnv {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariablesNMSMagmaLteDevEnv)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariablesOrc8rNginx) DeepCopyInto(out *EnvironmentVariablesOrc8rNginx) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariablesOrc8rNginx.
func (in *EnvironmentVariablesOrc8rNginx) DeepCopy() *EnvironmentVariablesOrc8rNginx {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariablesOrc8rNginx)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariablesOrc8rNotifier) DeepCopyInto(out *EnvironmentVariablesOrc8rNotifier) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariablesOrc8rNotifier.
func (in *EnvironmentVariablesOrc8rNotifier) DeepCopy() *EnvironmentVariablesOrc8rNotifier {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariablesOrc8rNotifier)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvironmentVariablesPrometheusKafka) DeepCopyInto(out *EnvironmentVariablesPrometheusKafka) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvironmentVariablesPrometheusKafka.
func (in *EnvironmentVariablesPrometheusKafka) DeepCopy() *EnvironmentVariablesPrometheusKafka {
	if in == nil {
		return nil
	}
	out := new(EnvironmentVariablesPrometheusKafka)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Image.
func (in *Image) DeepCopy() *Image {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistenForStatefulSet) DeepCopyInto(out *PersistenForStatefulSet) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistenForStatefulSet.
func (in *PersistenForStatefulSet) DeepCopy() *PersistenForStatefulSet {
	if in == nil {
		return nil
	}
	out := new(PersistenForStatefulSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Persistent) DeepCopyInto(out *Persistent) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Persistent.
func (in *Persistent) DeepCopy() *Persistent {
	if in == nil {
		return nil
	}
	out := new(Persistent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Pmnsystem) DeepCopyInto(out *Pmnsystem) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Pmnsystem.
func (in *Pmnsystem) DeepCopy() *Pmnsystem {
	if in == nil {
		return nil
	}
	out := new(Pmnsystem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Pmnsystem) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PmnsystemList) DeepCopyInto(out *PmnsystemList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Pmnsystem, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PmnsystemList.
func (in *PmnsystemList) DeepCopy() *PmnsystemList {
	if in == nil {
		return nil
	}
	out := new(PmnsystemList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PmnsystemList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PmnsystemSpec) DeepCopyInto(out *PmnsystemSpec) {
	*out = *in
	out.Persistent = in.Persistent
	out.PersistentForStatefulSet = in.PersistentForStatefulSet
	out.Image = in.Image
	if in.EnvVariables != nil {
		in, out := &in.EnvVariables, &out.EnvVariables
		*out = make([]EnvironmentVariables, len(*in))
		copy(*out, *in)
	}
	if in.EnvVariablesDirectoryD != nil {
		in, out := &in.EnvVariablesDirectoryD, &out.EnvVariablesDirectoryD
		*out = make([]EnvironmentVariablesDirectoryD, len(*in))
		copy(*out, *in)
	}
	if in.EnvVariablesOrc8rNginx != nil {
		in, out := &in.EnvVariablesOrc8rNginx, &out.EnvVariablesOrc8rNginx
		*out = make([]EnvironmentVariablesOrc8rNginx, len(*in))
		copy(*out, *in)
	}
	if in.EnvVariablesOrc8rNotifier != nil {
		in, out := &in.EnvVariablesOrc8rNotifier, &out.EnvVariablesOrc8rNotifier
		*out = make([]EnvironmentVariablesOrc8rNotifier, len(*in))
		copy(*out, *in)
	}
	if in.EnvVariablesNMSMagmaLte != nil {
		in, out := &in.EnvVariablesNMSMagmaLte, &out.EnvVariablesNMSMagmaLte
		*out = make([]EnvironmentVariablesNMSMagmaLte, len(*in))
		copy(*out, *in)
	}
	if in.EnvVariablesNMSMagmaLteDevEnv != nil {
		in, out := &in.EnvVariablesNMSMagmaLteDevEnv, &out.EnvVariablesNMSMagmaLteDevEnv
		*out = make([]EnvironmentVariablesNMSMagmaLteDevEnv, len(*in))
		copy(*out, *in)
	}
	if in.EnvVariablesPrometheusKafka != nil {
		in, out := &in.EnvVariablesPrometheusKafka, &out.EnvVariablesPrometheusKafka
		*out = make([]EnvironmentVariablesPrometheusKafka, len(*in))
		copy(*out, *in)
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]SecretConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PmnsystemSpec.
func (in *PmnsystemSpec) DeepCopy() *PmnsystemSpec {
	if in == nil {
		return nil
	}
	out := new(PmnsystemSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PmnsystemStatus) DeepCopyInto(out *PmnsystemStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PmnsystemStatus.
func (in *PmnsystemStatus) DeepCopy() *PmnsystemStatus {
	if in == nil {
		return nil
	}
	out := new(PmnsystemStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretConfig) DeepCopyInto(out *SecretConfig) {
	*out = *in
	if in.RequiredFiles != nil {
		in, out := &in.RequiredFiles, &out.RequiredFiles
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretConfig.
func (in *SecretConfig) DeepCopy() *SecretConfig {
	if in == nil {
		return nil
	}
	out := new(SecretConfig)
	in.DeepCopyInto(out)
	return out
}
