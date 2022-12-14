//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package metricsconfig

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsConfig) DeepCopyInto(out *MetricsConfig) {
	*out = *in
	if in.MetricsWithDuration != nil {
		in, out := &in.MetricsWithDuration, &out.MetricsWithDuration
		*out = make([]MetricsWithDuration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MetricsWithClusterTypeLabels != nil {
		in, out := &in.MetricsWithClusterTypeLabels, &out.MetricsWithClusterTypeLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsConfig.
func (in *MetricsConfig) DeepCopy() *MetricsConfig {
	if in == nil {
		return nil
	}
	out := new(MetricsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsWithDuration) DeepCopyInto(out *MetricsWithDuration) {
	*out = *in
	if in.Duration != nil {
		in, out := &in.Duration, &out.Duration
		*out = new(v1.Duration)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsWithDuration.
func (in *MetricsWithDuration) DeepCopy() *MetricsWithDuration {
	if in == nil {
		return nil
	}
	out := new(MetricsWithDuration)
	in.DeepCopyInto(out)
	return out
}
