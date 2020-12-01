// Code generated by skv2. DO NOT EDIT.

// Generated json marshal and unmarshal functions

package v1alpha1

import (
	bytes "bytes"
	fmt "fmt"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	skv2jsonpb "github.com/solo-io/skv2/pkg/jsonpb"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var (
	marshaller   = &skv2jsonpb.Marshaler{}
	unmarshaller = &skv2jsonpb.Unmarshaler{}
)

// MarshalJSON is a custom marshaler for XdsConfigSpec
func (this *XdsConfigSpec) MarshalJSON() ([]byte, error) {
	str, err := marshaller.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for XdsConfigSpec
func (this *XdsConfigSpec) UnmarshalJSON(b []byte) error {
	return unmarshaller.Unmarshal(bytes.NewReader(b), this)
}

// MarshalJSON is a custom marshaler for XdsConfigStatus
func (this *XdsConfigStatus) MarshalJSON() ([]byte, error) {
	str, err := marshaller.MarshalToString(this)
	return []byte(str), err
}

// UnmarshalJSON is a custom unmarshaler for XdsConfigStatus
func (this *XdsConfigStatus) UnmarshalJSON(b []byte) error {
	return unmarshaller.Unmarshal(bytes.NewReader(b), this)
}
