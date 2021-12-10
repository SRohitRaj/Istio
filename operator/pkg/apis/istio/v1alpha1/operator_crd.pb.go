// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/apis/istio/v1alpha1/operator_crd.proto

package v1alpha1

import (
	fmt "fmt"
	math "math"

	proto "google.golang.org/protobuf/proto"
	v11 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/api/operator/v1alpha1"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// IstioOperator is a CustomResourceDefinition (CRD) for an operator.
type IstioOperator struct {
	Kind                 string                      `protobuf:"bytes,5,opt,name=kind,proto3" json:"kind,omitempty"`
	ApiVersion           string                      `protobuf:"bytes,6,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Spec                 *v1alpha1.IstioOperatorSpec `protobuf:"bytes,7,opt,name=spec,proto3" json:"spec,omitempty"`
	Status				 *v1alpha1.InstallStatus     `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	v11.ObjectMeta       `json:"metadata,omitempty" protobuf:"bytes,9,opt,name=metadata"`
	v11.TypeMeta         `json:",inline"`
	Placeholder          string   `protobuf:"bytes,111,opt,name=placeholder,proto3" json:"placeholder,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IstioOperator) Reset()         { *m = IstioOperator{} }
func (m *IstioOperator) String() string { return proto.CompactTextString(m) }
func (*IstioOperator) ProtoMessage()    {}
func (*IstioOperator) Descriptor() ([]byte, []int) {
	return fileDescriptor_8eb082c28e72c148, []int{0}
}

func (m *IstioOperator) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IstioOperator.Unmarshal(m, b)
}
func (m *IstioOperator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IstioOperator.Marshal(b, m, deterministic)
}
func (m *IstioOperator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IstioOperator.Merge(m, src)
}
func (m *IstioOperator) XXX_Size() int {
	return xxx_messageInfo_IstioOperator.Size(m)
}
func (m *IstioOperator) XXX_DiscardUnknown() {
	xxx_messageInfo_IstioOperator.DiscardUnknown(m)
}

var xxx_messageInfo_IstioOperator proto.InternalMessageInfo

func (m *IstioOperator) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *IstioOperator) GetApiVersion() string {
	if m != nil {
		return m.ApiVersion
	}
	return ""
}

func (m *IstioOperator) GetPlaceholder() string {
	if m != nil {
		return m.Placeholder
	}
	return ""
}

func init() {
	proto.RegisterType((*IstioOperator)(nil), "v1alpha1.IstioOperator")
}

func init() {
	proto.RegisterFile("pkg/apis/istio/v1alpha1/operator_crd.proto", fileDescriptor_8eb082c28e72c148)
}

var fileDescriptor_8eb082c28e72c148 = []byte{
	// 145 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2a, 0xc8, 0x4e, 0xd7,
	0x4f, 0x2c, 0xc8, 0x2c, 0xd6, 0xcf, 0x2c, 0x2e, 0xc9, 0xcc, 0xd7, 0x2f, 0x33, 0x4c, 0xcc, 0x29,
	0xc8, 0x48, 0x34, 0xd4, 0xcf, 0x2f, 0x48, 0x2d, 0x4a, 0x2c, 0xc9, 0x2f, 0x8a, 0x4f, 0x2e, 0x4a,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x80, 0x49, 0x2a, 0xa5, 0x72, 0xf1, 0x7a, 0x82,
	0x94, 0xfb, 0x43, 0x15, 0x09, 0x09, 0x71, 0xb1, 0x64, 0x67, 0xe6, 0xa5, 0x48, 0xb0, 0x2a, 0x30,
	0x6a, 0x70, 0x06, 0x81, 0xd9, 0x42, 0x72, 0x5c, 0x5c, 0x89, 0x05, 0x99, 0x61, 0xa9, 0x45, 0xc5,
	0x99, 0xf9, 0x79, 0x12, 0x6c, 0x60, 0x19, 0x24, 0x11, 0x21, 0x05, 0x2e, 0xee, 0x82, 0x9c, 0xc4,
	0xe4, 0xd4, 0x8c, 0xfc, 0x9c, 0x94, 0xd4, 0x22, 0x89, 0x7c, 0xb0, 0x02, 0x64, 0xa1, 0x24, 0x36,
	0xb0, 0xbd, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xae, 0x9f, 0x05, 0xf0, 0xa5, 0x00, 0x00,
	0x00,
}
