// Code generated by skv2. DO NOT EDIT.

// This file contains generated Deepcopy methods for proto-based Spec and Status fields

package v1alpha1

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto for the Role.Spec
func (in *RoleSpec) DeepCopyInto(out *RoleSpec) {
	p := proto.Clone(in).(*RoleSpec)
	*out = *p
}

// DeepCopyInto for the Role.Status
func (in *RoleStatus) DeepCopyInto(out *RoleStatus) {
	p := proto.Clone(in).(*RoleStatus)
	*out = *p
}

// DeepCopyInto for the RoleBinding.Spec
func (in *RoleBindingSpec) DeepCopyInto(out *RoleBindingSpec) {
	p := proto.Clone(in).(*RoleBindingSpec)
	*out = *p
}

// DeepCopyInto for the RoleBinding.Status
func (in *RoleBindingStatus) DeepCopyInto(out *RoleBindingStatus) {
	p := proto.Clone(in).(*RoleBindingStatus)
	*out = *p
}
