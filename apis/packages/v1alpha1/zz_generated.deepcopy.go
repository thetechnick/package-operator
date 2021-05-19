// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageDeployment) DeepCopyInto(out *PackageDeployment) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageDeployment.
func (in *PackageDeployment) DeepCopy() *PackageDeployment {
	if in == nil {
		return nil
	}
	out := new(PackageDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageDeployment) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageDeploymentList) DeepCopyInto(out *PackageDeploymentList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PackageDeployment, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageDeploymentList.
func (in *PackageDeploymentList) DeepCopy() *PackageDeploymentList {
	if in == nil {
		return nil
	}
	out := new(PackageDeploymentList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageDeploymentList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageDeploymentSpec) DeepCopyInto(out *PackageDeploymentSpec) {
	*out = *in
	if in.RevisionHistoryLimit != nil {
		in, out := &in.RevisionHistoryLimit, &out.RevisionHistoryLimit
		*out = new(int)
		**out = **in
	}
	in.Selector.DeepCopyInto(&out.Selector)
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageDeploymentSpec.
func (in *PackageDeploymentSpec) DeepCopy() *PackageDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(PackageDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageDeploymentStatus) DeepCopyInto(out *PackageDeploymentStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CollisionCount != nil {
		in, out := &in.CollisionCount, &out.CollisionCount
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageDeploymentStatus.
func (in *PackageDeploymentStatus) DeepCopy() *PackageDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(PackageDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageObject) DeepCopyInto(out *PackageObject) {
	*out = *in
	in.Object.DeepCopyInto(&out.Object)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageObject.
func (in *PackageObject) DeepCopy() *PackageObject {
	if in == nil {
		return nil
	}
	out := new(PackageObject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackagePausedObject) DeepCopyInto(out *PackagePausedObject) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackagePausedObject.
func (in *PackagePausedObject) DeepCopy() *PackagePausedObject {
	if in == nil {
		return nil
	}
	out := new(PackagePausedObject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackagePhase) DeepCopyInto(out *PackagePhase) {
	*out = *in
	if in.Objects != nil {
		in, out := &in.Objects, &out.Objects
		*out = make([]PackageObject, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackagePhase.
func (in *PackagePhase) DeepCopy() *PackagePhase {
	if in == nil {
		return nil
	}
	out := new(PackagePhase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageProbe) DeepCopyInto(out *PackageProbe) {
	*out = *in
	in.Probe.DeepCopyInto(&out.Probe)
	in.Selector.DeepCopyInto(&out.Selector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageProbe.
func (in *PackageProbe) DeepCopy() *PackageProbe {
	if in == nil {
		return nil
	}
	out := new(PackageProbe)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageProbeKindSpec) DeepCopyInto(out *PackageProbeKindSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageProbeKindSpec.
func (in *PackageProbeKindSpec) DeepCopy() *PackageProbeKindSpec {
	if in == nil {
		return nil
	}
	out := new(PackageProbeKindSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageSet) DeepCopyInto(out *PackageSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageSet.
func (in *PackageSet) DeepCopy() *PackageSet {
	if in == nil {
		return nil
	}
	out := new(PackageSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageSetList) DeepCopyInto(out *PackageSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PackageSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageSetList.
func (in *PackageSetList) DeepCopy() *PackageSetList {
	if in == nil {
		return nil
	}
	out := new(PackageSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PackageSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageSetSpec) DeepCopyInto(out *PackageSetSpec) {
	*out = *in
	if in.PausedFor != nil {
		in, out := &in.PausedFor, &out.PausedFor
		*out = make([]PackagePausedObject, len(*in))
		copy(*out, *in)
	}
	if in.Phases != nil {
		in, out := &in.Phases, &out.Phases
		*out = make([]PackagePhase, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ReadinessProbes != nil {
		in, out := &in.ReadinessProbes, &out.ReadinessProbes
		*out = make([]PackageProbe, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageSetSpec.
func (in *PackageSetSpec) DeepCopy() *PackageSetSpec {
	if in == nil {
		return nil
	}
	out := new(PackageSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageSetStatus) DeepCopyInto(out *PackageSetStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.PausedFor != nil {
		in, out := &in.PausedFor, &out.PausedFor
		*out = make([]PackagePausedObject, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageSetStatus.
func (in *PackageSetStatus) DeepCopy() *PackageSetStatus {
	if in == nil {
		return nil
	}
	out := new(PackageSetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageSetTemplate) DeepCopyInto(out *PackageSetTemplate) {
	*out = *in
	in.Metadata.DeepCopyInto(&out.Metadata)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageSetTemplate.
func (in *PackageSetTemplate) DeepCopy() *PackageSetTemplate {
	if in == nil {
		return nil
	}
	out := new(PackageSetTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PackageSetTemplateSpec) DeepCopyInto(out *PackageSetTemplateSpec) {
	*out = *in
	if in.Phases != nil {
		in, out := &in.Phases, &out.Phases
		*out = make([]PackagePhase, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ReadinessProbes != nil {
		in, out := &in.ReadinessProbes, &out.ReadinessProbes
		*out = make([]PackageProbe, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PackageSetTemplateSpec.
func (in *PackageSetTemplateSpec) DeepCopy() *PackageSetTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(PackageSetTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Probe) DeepCopyInto(out *Probe) {
	*out = *in
	if in.Condition != nil {
		in, out := &in.Condition, &out.Condition
		*out = new(ProbeConditionSpec)
		**out = **in
	}
	if in.FieldsEqual != nil {
		in, out := &in.FieldsEqual, &out.FieldsEqual
		*out = new(ProbeFieldsEqualSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Probe.
func (in *Probe) DeepCopy() *Probe {
	if in == nil {
		return nil
	}
	out := new(Probe)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProbeConditionSpec) DeepCopyInto(out *ProbeConditionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProbeConditionSpec.
func (in *ProbeConditionSpec) DeepCopy() *ProbeConditionSpec {
	if in == nil {
		return nil
	}
	out := new(ProbeConditionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProbeFieldsEqualSpec) DeepCopyInto(out *ProbeFieldsEqualSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProbeFieldsEqualSpec.
func (in *ProbeFieldsEqualSpec) DeepCopy() *ProbeFieldsEqualSpec {
	if in == nil {
		return nil
	}
	out := new(ProbeFieldsEqualSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProbeSelector) DeepCopyInto(out *ProbeSelector) {
	*out = *in
	if in.Kind != nil {
		in, out := &in.Kind, &out.Kind
		*out = new(PackageProbeKindSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProbeSelector.
func (in *ProbeSelector) DeepCopy() *ProbeSelector {
	if in == nil {
		return nil
	}
	out := new(ProbeSelector)
	in.DeepCopyInto(out)
	return out
}
