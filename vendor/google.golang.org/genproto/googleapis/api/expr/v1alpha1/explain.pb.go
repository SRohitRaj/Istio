// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/expr/v1alpha1/explain.proto

package expr

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Values of intermediate expressions produced when evaluating expression.
// Deprecated, use `EvalState` instead.
//
// Deprecated: Do not use.
type Explain struct {
	// All of the observed values.
	//
	// The field value_index is an index in the values list.
	// Separating values from steps is needed to remove redundant values.
	Values []*Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	// List of steps.
	//
	// Repeated evaluations of the same expression generate new ExprStep
	// instances. The order of such ExprStep instances matches the order of
	// elements returned by Comprehension.iter_range.
	ExprSteps            []*Explain_ExprStep `protobuf:"bytes,2,rep,name=expr_steps,json=exprSteps,proto3" json:"expr_steps,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Explain) Reset()         { *m = Explain{} }
func (m *Explain) String() string { return proto.CompactTextString(m) }
func (*Explain) ProtoMessage()    {}
func (*Explain) Descriptor() ([]byte, []int) {
	return fileDescriptor_2df9793dd8748e27, []int{0}
}

func (m *Explain) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Explain.Unmarshal(m, b)
}
func (m *Explain) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Explain.Marshal(b, m, deterministic)
}
func (m *Explain) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Explain.Merge(m, src)
}
func (m *Explain) XXX_Size() int {
	return xxx_messageInfo_Explain.Size(m)
}
func (m *Explain) XXX_DiscardUnknown() {
	xxx_messageInfo_Explain.DiscardUnknown(m)
}

var xxx_messageInfo_Explain proto.InternalMessageInfo

func (m *Explain) GetValues() []*Value {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *Explain) GetExprSteps() []*Explain_ExprStep {
	if m != nil {
		return m.ExprSteps
	}
	return nil
}

// ID and value index of one step.
type Explain_ExprStep struct {
	// ID of corresponding Expr node.
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Index of the value in the values list.
	ValueIndex           int32    `protobuf:"varint,2,opt,name=value_index,json=valueIndex,proto3" json:"value_index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Explain_ExprStep) Reset()         { *m = Explain_ExprStep{} }
func (m *Explain_ExprStep) String() string { return proto.CompactTextString(m) }
func (*Explain_ExprStep) ProtoMessage()    {}
func (*Explain_ExprStep) Descriptor() ([]byte, []int) {
	return fileDescriptor_2df9793dd8748e27, []int{0, 0}
}

func (m *Explain_ExprStep) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Explain_ExprStep.Unmarshal(m, b)
}
func (m *Explain_ExprStep) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Explain_ExprStep.Marshal(b, m, deterministic)
}
func (m *Explain_ExprStep) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Explain_ExprStep.Merge(m, src)
}
func (m *Explain_ExprStep) XXX_Size() int {
	return xxx_messageInfo_Explain_ExprStep.Size(m)
}
func (m *Explain_ExprStep) XXX_DiscardUnknown() {
	xxx_messageInfo_Explain_ExprStep.DiscardUnknown(m)
}

var xxx_messageInfo_Explain_ExprStep proto.InternalMessageInfo

func (m *Explain_ExprStep) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Explain_ExprStep) GetValueIndex() int32 {
	if m != nil {
		return m.ValueIndex
	}
	return 0
}

func init() {
	proto.RegisterType((*Explain)(nil), "google.api.expr.v1alpha1.Explain")
	proto.RegisterType((*Explain_ExprStep)(nil), "google.api.expr.v1alpha1.Explain.ExprStep")
}

func init() {
	proto.RegisterFile("google/api/expr/v1alpha1/explain.proto", fileDescriptor_2df9793dd8748e27)
}

var fileDescriptor_2df9793dd8748e27 = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xb1, 0x4b, 0x03, 0x31,
	0x14, 0xc6, 0x79, 0x29, 0x56, 0x7d, 0x15, 0x87, 0x4c, 0xa1, 0x08, 0x3d, 0x44, 0xe4, 0x70, 0x48,
	0xa8, 0x0e, 0x82, 0x75, 0x2a, 0x38, 0x74, 0x2b, 0x27, 0x38, 0xb8, 0x94, 0xe8, 0x85, 0x18, 0x88,
	0x97, 0x70, 0x39, 0x4b, 0xff, 0x4a, 0xff, 0x1e, 0x47, 0x49, 0x2e, 0x37, 0x95, 0x9b, 0xee, 0xde,
	0xfb, 0x7e, 0xdf, 0xf7, 0x91, 0x87, 0xb7, 0xda, 0x39, 0x6d, 0x95, 0x90, 0xde, 0x08, 0x75, 0xf0,
	0xad, 0xd8, 0x2f, 0xa5, 0xf5, 0x5f, 0x72, 0x19, 0x27, 0x2b, 0x4d, 0xc3, 0x7d, 0xeb, 0x3a, 0x47,
	0x59, 0xcf, 0x71, 0xe9, 0x0d, 0x8f, 0x1c, 0x1f, 0xb8, 0xf9, 0xcd, 0x68, 0xc2, 0x5e, 0xda, 0x1f,
	0xd5, 0xfb, 0xaf, 0x7f, 0x01, 0x4f, 0x5f, 0xfa, 0x44, 0xfa, 0x88, 0xd3, 0x24, 0x05, 0x06, 0xc5,
	0xa4, 0x9c, 0xdd, 0x2f, 0xf8, 0x58, 0x38, 0x7f, 0x8b, 0x5c, 0x95, 0x71, 0xba, 0x41, 0x8c, 0xf2,
	0x2e, 0x74, 0xca, 0x07, 0x46, 0x92, 0xf9, 0x6e, 0xdc, 0x9c, 0xfb, 0xe2, 0xb7, 0x7d, 0xed, 0x94,
	0xaf, 0xce, 0x55, 0xfe, 0x0b, 0xf3, 0x15, 0x9e, 0x0d, 0x6b, 0x7a, 0x89, 0xc4, 0xd4, 0x0c, 0x0a,
	0x28, 0x27, 0x15, 0x31, 0x35, 0x5d, 0xe0, 0x2c, 0x15, 0xee, 0x4c, 0x53, 0xab, 0x03, 0x23, 0x05,
	0x94, 0x27, 0x15, 0xa6, 0xd5, 0x26, 0x6e, 0x9e, 0x08, 0x83, 0xb5, 0xc3, 0xab, 0x4f, 0xf7, 0x3d,
	0x5a, 0xbe, 0xbe, 0xc8, 0xed, 0xdb, 0xf8, 0xfc, 0x2d, 0xbc, 0x3f, 0x67, 0x52, 0x3b, 0x2b, 0x1b,
	0xcd, 0x5d, 0xab, 0x85, 0x56, 0x4d, 0x3a, 0x8e, 0xe8, 0x25, 0xe9, 0x4d, 0x38, 0xbe, 0xe2, 0x2a,
	0x4e, 0x7f, 0x00, 0x1f, 0xd3, 0xc4, 0x3e, 0xfc, 0x07, 0x00, 0x00, 0xff, 0xff, 0x34, 0xf2, 0xb9,
	0x9e, 0xb2, 0x01, 0x00, 0x00,
}
