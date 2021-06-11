// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	packagesv1alpha1 "github.com/thetechnick/package-operator/apis/packages/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Handover) DeepCopyInto(out *Handover) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Handover.
func (in *Handover) DeepCopy() *Handover {
	if in == nil {
		return nil
	}
	out := new(Handover)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Handover) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverList) DeepCopyInto(out *HandoverList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Handover, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverList.
func (in *HandoverList) DeepCopy() *HandoverList {
	if in == nil {
		return nil
	}
	out := new(HandoverList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HandoverList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverRef) DeepCopyInto(out *HandoverRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverRef.
func (in *HandoverRef) DeepCopy() *HandoverRef {
	if in == nil {
		return nil
	}
	out := new(HandoverRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverSpec) DeepCopyInto(out *HandoverSpec) {
	*out = *in
	in.Strategy.DeepCopyInto(&out.Strategy)
	out.Target = in.Target
	if in.Probes != nil {
		in, out := &in.Probes, &out.Probes
		*out = make([]packagesv1alpha1.Probe, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverSpec.
func (in *HandoverSpec) DeepCopy() *HandoverSpec {
	if in == nil {
		return nil
	}
	out := new(HandoverSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverStatus) DeepCopyInto(out *HandoverStatus) {
	*out = *in
	if in.Processing != nil {
		in, out := &in.Processing, &out.Processing
		*out = make([]HandoverRef, len(*in))
		copy(*out, *in)
	}
	out.Stats = in.Stats
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverStatus.
func (in *HandoverStatus) DeepCopy() *HandoverStatus {
	if in == nil {
		return nil
	}
	out := new(HandoverStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverStatusStats) DeepCopyInto(out *HandoverStatusStats) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverStatusStats.
func (in *HandoverStatusStats) DeepCopy() *HandoverStatusStats {
	if in == nil {
		return nil
	}
	out := new(HandoverStatusStats)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverStrategy) DeepCopyInto(out *HandoverStrategy) {
	*out = *in
	if in.Relabel != nil {
		in, out := &in.Relabel, &out.Relabel
		*out = new(HandoverStrategyRelabelSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverStrategy.
func (in *HandoverStrategy) DeepCopy() *HandoverStrategy {
	if in == nil {
		return nil
	}
	out := new(HandoverStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverStrategyRelabelSpec) DeepCopyInto(out *HandoverStrategyRelabelSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverStrategyRelabelSpec.
func (in *HandoverStrategyRelabelSpec) DeepCopy() *HandoverStrategyRelabelSpec {
	if in == nil {
		return nil
	}
	out := new(HandoverStrategyRelabelSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HandoverTarget) DeepCopyInto(out *HandoverTarget) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HandoverTarget.
func (in *HandoverTarget) DeepCopy() *HandoverTarget {
	if in == nil {
		return nil
	}
	out := new(HandoverTarget)
	in.DeepCopyInto(out)
	return out
}
