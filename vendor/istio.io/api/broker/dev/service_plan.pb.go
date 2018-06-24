// Code generated by protoc-gen-go. DO NOT EDIT.
// source: broker/dev/service_plan.proto

package dev // import "istio.io/api/broker/dev"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// $hide_from_docs
// ServicePlan defines the type of services available to Istio service
// consumers.  One or more services are included in a plan. The plan is flexible
// and subject to change along with business requirements.
type ServicePlan struct {
	// Required. Public plan information.
	Plan *CatalogPlan `protobuf:"bytes,1,opt,name=plan,proto3" json:"plan,omitempty"`
	// Required. List of the Keys of serviceclass config instance
	// that are included in the plan.
	// ServiceClass is a type of CRD resource.
	Services             []string `protobuf:"bytes,2,rep,name=services,proto3" json:"services,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServicePlan) Reset()         { *m = ServicePlan{} }
func (m *ServicePlan) String() string { return proto.CompactTextString(m) }
func (*ServicePlan) ProtoMessage()    {}
func (*ServicePlan) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_plan_7d5f4530d9bbc850, []int{0}
}
func (m *ServicePlan) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServicePlan.Unmarshal(m, b)
}
func (m *ServicePlan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServicePlan.Marshal(b, m, deterministic)
}
func (dst *ServicePlan) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServicePlan.Merge(dst, src)
}
func (m *ServicePlan) XXX_Size() int {
	return xxx_messageInfo_ServicePlan.Size(m)
}
func (m *ServicePlan) XXX_DiscardUnknown() {
	xxx_messageInfo_ServicePlan.DiscardUnknown(m)
}

var xxx_messageInfo_ServicePlan proto.InternalMessageInfo

func (m *ServicePlan) GetPlan() *CatalogPlan {
	if m != nil {
		return m.Plan
	}
	return nil
}

func (m *ServicePlan) GetServices() []string {
	if m != nil {
		return m.Services
	}
	return nil
}

// $hide_from_docs
// CatalogPlan defines listing information for this service plan within the
// exposed catalog.  The message is a subset of OSBI plan fields defined in
// https://github.com/openservicebrokerapi
type CatalogPlan struct {
	// Required. Public service plan name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Required. Public unique service plan guid.
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// Required. Public short service plan description.
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CatalogPlan) Reset()         { *m = CatalogPlan{} }
func (m *CatalogPlan) String() string { return proto.CompactTextString(m) }
func (*CatalogPlan) ProtoMessage()    {}
func (*CatalogPlan) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_plan_7d5f4530d9bbc850, []int{1}
}
func (m *CatalogPlan) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CatalogPlan.Unmarshal(m, b)
}
func (m *CatalogPlan) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CatalogPlan.Marshal(b, m, deterministic)
}
func (dst *CatalogPlan) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CatalogPlan.Merge(dst, src)
}
func (m *CatalogPlan) XXX_Size() int {
	return xxx_messageInfo_CatalogPlan.Size(m)
}
func (m *CatalogPlan) XXX_DiscardUnknown() {
	xxx_messageInfo_CatalogPlan.DiscardUnknown(m)
}

var xxx_messageInfo_CatalogPlan proto.InternalMessageInfo

func (m *CatalogPlan) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CatalogPlan) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *CatalogPlan) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*ServicePlan)(nil), "istio.broker.dev.ServicePlan")
	proto.RegisterType((*CatalogPlan)(nil), "istio.broker.dev.CatalogPlan")
}

func init() {
	proto.RegisterFile("broker/dev/service_plan.proto", fileDescriptor_service_plan_7d5f4530d9bbc850)
}

var fileDescriptor_service_plan_7d5f4530d9bbc850 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0x69, 0x5a, 0xc4, 0x4e, 0x40, 0x24, 0x17, 0xa3, 0x50, 0x28, 0x3d, 0xf5, 0x94, 0xa0,
	0xbe, 0x81, 0xbe, 0x80, 0xb4, 0x37, 0x11, 0x24, 0x6d, 0x82, 0x0c, 0xd6, 0x24, 0x24, 0xa1, 0xcf,
	0x2f, 0x4d, 0x97, 0xdd, 0xb2, 0xb7, 0x99, 0x7f, 0xbe, 0x99, 0xf9, 0x7f, 0x68, 0xa6, 0xe0, 0x7e,
	0x4d, 0x90, 0xda, 0xac, 0x32, 0x9a, 0xb0, 0xe2, 0x6c, 0xbe, 0xfd, 0xa2, 0xac, 0xf0, 0xc1, 0x25,
	0xc7, 0xee, 0x31, 0x26, 0x74, 0x62, 0x87, 0x84, 0x36, 0x6b, 0xf7, 0x05, 0x74, 0xdc, 0xb9, 0x8f,
	0x45, 0x59, 0xf6, 0x0c, 0xd5, 0x86, 0xf3, 0xa2, 0x2d, 0x7a, 0xfa, 0xd2, 0x88, 0x6b, 0x5e, 0xbc,
	0xab, 0xa4, 0x16, 0xf7, 0xb3, 0xc1, 0x43, 0x46, 0xd9, 0x13, 0xdc, 0x9e, 0x3e, 0x45, 0x4e, 0xda,
	0xb2, 0xaf, 0x87, 0x73, 0xdf, 0x8d, 0x40, 0x0f, 0x0b, 0x8c, 0x41, 0x65, 0xd5, 0x9f, 0xc9, 0xd7,
	0xeb, 0x21, 0xd7, 0xec, 0x0e, 0x08, 0x6a, 0x4e, 0xb2, 0x42, 0x50, 0xb3, 0x16, 0xa8, 0x36, 0x71,
	0x0e, 0xe8, 0x13, 0x3a, 0xcb, 0xcb, 0x3c, 0x38, 0x4a, 0x6f, 0x8f, 0x9f, 0x0f, 0xbb, 0x2d, 0x74,
	0x52, 0x79, 0x94, 0x97, 0xc8, 0xd3, 0x4d, 0x8e, 0xf9, 0xfa, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x9c,
	0x34, 0x35, 0x2a, 0x07, 0x01, 0x00, 0x00,
}
