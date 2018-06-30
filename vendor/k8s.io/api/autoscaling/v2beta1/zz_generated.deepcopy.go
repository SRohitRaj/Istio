// +build !ignore_autogenerated

/*
Copyright 2018 The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v2beta1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CrossVersionObjectReference) DeepCopyInto(out *CrossVersionObjectReference) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CrossVersionObjectReference.
func (in *CrossVersionObjectReference) DeepCopy() *CrossVersionObjectReference {
	if in == nil {
		return nil
	}
	out := new(CrossVersionObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalMetricSource) DeepCopyInto(out *ExternalMetricSource) {
	*out = *in
	if in.MetricSelector != nil {
		in, out := &in.MetricSelector, &out.MetricSelector
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.LabelSelector)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.TargetValue != nil {
		in, out := &in.TargetValue, &out.TargetValue
		if *in == nil {
			*out = nil
		} else {
			x := (*in).DeepCopy()
			*out = &x
		}
	}
	if in.TargetAverageValue != nil {
		in, out := &in.TargetAverageValue, &out.TargetAverageValue
		if *in == nil {
			*out = nil
		} else {
			x := (*in).DeepCopy()
			*out = &x
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalMetricSource.
func (in *ExternalMetricSource) DeepCopy() *ExternalMetricSource {
	if in == nil {
		return nil
	}
	out := new(ExternalMetricSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalMetricStatus) DeepCopyInto(out *ExternalMetricStatus) {
	*out = *in
	if in.MetricSelector != nil {
		in, out := &in.MetricSelector, &out.MetricSelector
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.LabelSelector)
			(*in).DeepCopyInto(*out)
		}
	}
	out.CurrentValue = in.CurrentValue.DeepCopy()
	if in.CurrentAverageValue != nil {
		in, out := &in.CurrentAverageValue, &out.CurrentAverageValue
		if *in == nil {
			*out = nil
		} else {
			x := (*in).DeepCopy()
			*out = &x
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalMetricStatus.
func (in *ExternalMetricStatus) DeepCopy() *ExternalMetricStatus {
	if in == nil {
		return nil
	}
	out := new(ExternalMetricStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalPodAutoscaler) DeepCopyInto(out *HorizontalPodAutoscaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalPodAutoscaler.
func (in *HorizontalPodAutoscaler) DeepCopy() *HorizontalPodAutoscaler {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HorizontalPodAutoscaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalPodAutoscalerCondition) DeepCopyInto(out *HorizontalPodAutoscalerCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalPodAutoscalerCondition.
func (in *HorizontalPodAutoscalerCondition) DeepCopy() *HorizontalPodAutoscalerCondition {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalPodAutoscalerList) DeepCopyInto(out *HorizontalPodAutoscalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HorizontalPodAutoscaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalPodAutoscalerList.
func (in *HorizontalPodAutoscalerList) DeepCopy() *HorizontalPodAutoscalerList {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HorizontalPodAutoscalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalPodAutoscalerSpec) DeepCopyInto(out *HorizontalPodAutoscalerSpec) {
	*out = *in
	out.ScaleTargetRef = in.ScaleTargetRef
	if in.MinReplicas != nil {
		in, out := &in.MinReplicas, &out.MinReplicas
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = make([]MetricSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalPodAutoscalerSpec.
func (in *HorizontalPodAutoscalerSpec) DeepCopy() *HorizontalPodAutoscalerSpec {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalPodAutoscalerStatus) DeepCopyInto(out *HorizontalPodAutoscalerStatus) {
	*out = *in
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		if *in == nil {
			*out = nil
		} else {
			*out = new(int64)
			**out = **in
		}
	}
	if in.LastScaleTime != nil {
		in, out := &in.LastScaleTime, &out.LastScaleTime
		if *in == nil {
			*out = nil
		} else {
			*out = (*in).DeepCopy()
		}
	}
	if in.CurrentMetrics != nil {
		in, out := &in.CurrentMetrics, &out.CurrentMetrics
		*out = make([]MetricStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]HorizontalPodAutoscalerCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalPodAutoscalerStatus.
func (in *HorizontalPodAutoscalerStatus) DeepCopy() *HorizontalPodAutoscalerStatus {
	if in == nil {
		return nil
	}
	out := new(HorizontalPodAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricSpec) DeepCopyInto(out *MetricSpec) {
	*out = *in
	if in.Object != nil {
		in, out := &in.Object, &out.Object
		if *in == nil {
			*out = nil
		} else {
			*out = new(ObjectMetricSource)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		if *in == nil {
			*out = nil
		} else {
			*out = new(PodsMetricSource)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		if *in == nil {
			*out = nil
		} else {
			*out = new(ResourceMetricSource)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.External != nil {
		in, out := &in.External, &out.External
		if *in == nil {
			*out = nil
		} else {
			*out = new(ExternalMetricSource)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricSpec.
func (in *MetricSpec) DeepCopy() *MetricSpec {
	if in == nil {
		return nil
	}
	out := new(MetricSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricStatus) DeepCopyInto(out *MetricStatus) {
	*out = *in
	if in.Object != nil {
		in, out := &in.Object, &out.Object
		if *in == nil {
			*out = nil
		} else {
			*out = new(ObjectMetricStatus)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Pods != nil {
		in, out := &in.Pods, &out.Pods
		if *in == nil {
			*out = nil
		} else {
			*out = new(PodsMetricStatus)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Resource != nil {
		in, out := &in.Resource, &out.Resource
		if *in == nil {
			*out = nil
		} else {
			*out = new(ResourceMetricStatus)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.External != nil {
		in, out := &in.External, &out.External
		if *in == nil {
			*out = nil
		} else {
			*out = new(ExternalMetricStatus)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricStatus.
func (in *MetricStatus) DeepCopy() *MetricStatus {
	if in == nil {
		return nil
	}
	out := new(MetricStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectMetricSource) DeepCopyInto(out *ObjectMetricSource) {
	*out = *in
	out.Target = in.Target
	out.TargetValue = in.TargetValue.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectMetricSource.
func (in *ObjectMetricSource) DeepCopy() *ObjectMetricSource {
	if in == nil {
		return nil
	}
	out := new(ObjectMetricSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectMetricStatus) DeepCopyInto(out *ObjectMetricStatus) {
	*out = *in
	out.Target = in.Target
	out.CurrentValue = in.CurrentValue.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectMetricStatus.
func (in *ObjectMetricStatus) DeepCopy() *ObjectMetricStatus {
	if in == nil {
		return nil
	}
	out := new(ObjectMetricStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodsMetricSource) DeepCopyInto(out *PodsMetricSource) {
	*out = *in
	out.TargetAverageValue = in.TargetAverageValue.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodsMetricSource.
func (in *PodsMetricSource) DeepCopy() *PodsMetricSource {
	if in == nil {
		return nil
	}
	out := new(PodsMetricSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodsMetricStatus) DeepCopyInto(out *PodsMetricStatus) {
	*out = *in
	out.CurrentAverageValue = in.CurrentAverageValue.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodsMetricStatus.
func (in *PodsMetricStatus) DeepCopy() *PodsMetricStatus {
	if in == nil {
		return nil
	}
	out := new(PodsMetricStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceMetricSource) DeepCopyInto(out *ResourceMetricSource) {
	*out = *in
	if in.TargetAverageUtilization != nil {
		in, out := &in.TargetAverageUtilization, &out.TargetAverageUtilization
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.TargetAverageValue != nil {
		in, out := &in.TargetAverageValue, &out.TargetAverageValue
		if *in == nil {
			*out = nil
		} else {
			x := (*in).DeepCopy()
			*out = &x
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceMetricSource.
func (in *ResourceMetricSource) DeepCopy() *ResourceMetricSource {
	if in == nil {
		return nil
	}
	out := new(ResourceMetricSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceMetricStatus) DeepCopyInto(out *ResourceMetricStatus) {
	*out = *in
	if in.CurrentAverageUtilization != nil {
		in, out := &in.CurrentAverageUtilization, &out.CurrentAverageUtilization
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	out.CurrentAverageValue = in.CurrentAverageValue.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceMetricStatus.
func (in *ResourceMetricStatus) DeepCopy() *ResourceMetricStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceMetricStatus)
	in.DeepCopyInto(out)
	return out
}
