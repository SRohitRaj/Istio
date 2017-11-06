// Code generated by protoc-gen-go.
// source: broker/pkg/testing/mock/proto/fake_config.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	broker/pkg/testing/mock/proto/fake_config.proto

It has these top-level messages:
	FakeConfig
	Pair
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// Testing config resource consisting
// of a set of key-value pairs
type FakeConfig struct {
	Key   string  `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Pairs []*Pair `protobuf:"bytes,2,rep,name=pairs" json:"pairs,omitempty"`
}

func (m *FakeConfig) Reset()                    { *m = FakeConfig{} }
func (m *FakeConfig) String() string            { return proto1.CompactTextString(m) }
func (*FakeConfig) ProtoMessage()               {}
func (*FakeConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *FakeConfig) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *FakeConfig) GetPairs() []*Pair {
	if m != nil {
		return m.Pairs
	}
	return nil
}

type Pair struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Pair) Reset()                    { *m = Pair{} }
func (m *Pair) String() string            { return proto1.CompactTextString(m) }
func (*Pair) ProtoMessage()               {}
func (*Pair) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Pair) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Pair) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto1.RegisterType((*FakeConfig)(nil), "proto.FakeConfig")
	proto1.RegisterType((*Pair)(nil), "proto.Pair")
}

func init() { proto1.RegisterFile("broker/pkg/testing/mock/proto/fake_config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 154 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4f, 0x2a, 0xca, 0xcf,
	0x4e, 0x2d, 0xd2, 0x2f, 0xc8, 0x4e, 0xd7, 0x2f, 0x49, 0x2d, 0x2e, 0xc9, 0xcc, 0x4b, 0xd7, 0xcf,
	0xcd, 0x4f, 0xce, 0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x4f, 0x4b, 0xcc, 0x4e, 0x8d, 0x4f,
	0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x03, 0x8b, 0x08, 0xb1, 0x82, 0x29, 0x25, 0x47, 0x2e, 0x2e,
	0xb7, 0xc4, 0xec, 0x54, 0x67, 0xb0, 0x94, 0x90, 0x00, 0x17, 0x73, 0x76, 0x6a, 0xa5, 0x04, 0xa3,
	0x02, 0xa3, 0x06, 0x67, 0x10, 0x88, 0x29, 0xa4, 0xc8, 0xc5, 0x5a, 0x90, 0x98, 0x59, 0x54, 0x2c,
	0xc1, 0xa4, 0xc0, 0xac, 0xc1, 0x6d, 0xc4, 0x0d, 0xd1, 0xad, 0x17, 0x90, 0x98, 0x59, 0x14, 0x04,
	0x91, 0x51, 0xd2, 0xe3, 0x62, 0x01, 0x71, 0xb1, 0x68, 0x16, 0xe1, 0x62, 0x2d, 0x4b, 0xcc, 0x29,
	0x4d, 0x95, 0x60, 0x02, 0x8b, 0x41, 0x38, 0x49, 0x6c, 0x60, 0x23, 0x8c, 0x01, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x62, 0x8f, 0x63, 0xae, 0xb3, 0x00, 0x00, 0x00,
}
