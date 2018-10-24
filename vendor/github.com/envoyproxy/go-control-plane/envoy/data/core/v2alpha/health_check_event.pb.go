// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/data/core/v2alpha/health_check_event.proto

package envoy_data_core_v2alpha

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import _ "github.com/lyft/protoc-gen-validate/validate"

import time "time"

import bytes "bytes"

import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type HealthCheckFailureType int32

const (
	HealthCheckFailureType_ACTIVE  HealthCheckFailureType = 0
	HealthCheckFailureType_PASSIVE HealthCheckFailureType = 1
	HealthCheckFailureType_NETWORK HealthCheckFailureType = 2
)

var HealthCheckFailureType_name = map[int32]string{
	0: "ACTIVE",
	1: "PASSIVE",
	2: "NETWORK",
}
var HealthCheckFailureType_value = map[string]int32{
	"ACTIVE":  0,
	"PASSIVE": 1,
	"NETWORK": 2,
}

func (x HealthCheckFailureType) String() string {
	return proto.EnumName(HealthCheckFailureType_name, int32(x))
}
func (HealthCheckFailureType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_health_check_event_993d394457bc7a4f, []int{0}
}

type HealthCheckerType int32

const (
	HealthCheckerType_HTTP  HealthCheckerType = 0
	HealthCheckerType_TCP   HealthCheckerType = 1
	HealthCheckerType_GRPC  HealthCheckerType = 2
	HealthCheckerType_REDIS HealthCheckerType = 3
)

var HealthCheckerType_name = map[int32]string{
	0: "HTTP",
	1: "TCP",
	2: "GRPC",
	3: "REDIS",
}
var HealthCheckerType_value = map[string]int32{
	"HTTP":  0,
	"TCP":   1,
	"GRPC":  2,
	"REDIS": 3,
}

func (x HealthCheckerType) String() string {
	return proto.EnumName(HealthCheckerType_name, int32(x))
}
func (HealthCheckerType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_health_check_event_993d394457bc7a4f, []int{1}
}

type HealthCheckEvent struct {
	HealthCheckerType HealthCheckerType `protobuf:"varint,1,opt,name=health_checker_type,json=healthCheckerType,proto3,enum=envoy.data.core.v2alpha.HealthCheckerType" json:"health_checker_type,omitempty"`
	Host              *core.Address     `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
	ClusterName       string            `protobuf:"bytes,3,opt,name=cluster_name,json=clusterName,proto3" json:"cluster_name,omitempty"`
	// Types that are valid to be assigned to Event:
	//	*HealthCheckEvent_EjectUnhealthyEvent
	//	*HealthCheckEvent_AddHealthyEvent
	Event isHealthCheckEvent_Event `protobuf_oneof:"event"`
	// Timestamp for event.
	Timestamp            *time.Time `protobuf:"bytes,6,opt,name=timestamp,stdtime" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *HealthCheckEvent) Reset()         { *m = HealthCheckEvent{} }
func (m *HealthCheckEvent) String() string { return proto.CompactTextString(m) }
func (*HealthCheckEvent) ProtoMessage()    {}
func (*HealthCheckEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_check_event_993d394457bc7a4f, []int{0}
}
func (m *HealthCheckEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HealthCheckEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HealthCheckEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *HealthCheckEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckEvent.Merge(dst, src)
}
func (m *HealthCheckEvent) XXX_Size() int {
	return m.Size()
}
func (m *HealthCheckEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckEvent.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckEvent proto.InternalMessageInfo

type isHealthCheckEvent_Event interface {
	isHealthCheckEvent_Event()
	Equal(interface{}) bool
	MarshalTo([]byte) (int, error)
	Size() int
}

type HealthCheckEvent_EjectUnhealthyEvent struct {
	EjectUnhealthyEvent *HealthCheckEjectUnhealthy `protobuf:"bytes,4,opt,name=eject_unhealthy_event,json=ejectUnhealthyEvent,oneof"`
}
type HealthCheckEvent_AddHealthyEvent struct {
	AddHealthyEvent *HealthCheckAddHealthy `protobuf:"bytes,5,opt,name=add_healthy_event,json=addHealthyEvent,oneof"`
}

func (*HealthCheckEvent_EjectUnhealthyEvent) isHealthCheckEvent_Event() {}
func (*HealthCheckEvent_AddHealthyEvent) isHealthCheckEvent_Event()     {}

func (m *HealthCheckEvent) GetEvent() isHealthCheckEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *HealthCheckEvent) GetHealthCheckerType() HealthCheckerType {
	if m != nil {
		return m.HealthCheckerType
	}
	return HealthCheckerType_HTTP
}

func (m *HealthCheckEvent) GetHost() *core.Address {
	if m != nil {
		return m.Host
	}
	return nil
}

func (m *HealthCheckEvent) GetClusterName() string {
	if m != nil {
		return m.ClusterName
	}
	return ""
}

func (m *HealthCheckEvent) GetEjectUnhealthyEvent() *HealthCheckEjectUnhealthy {
	if x, ok := m.GetEvent().(*HealthCheckEvent_EjectUnhealthyEvent); ok {
		return x.EjectUnhealthyEvent
	}
	return nil
}

func (m *HealthCheckEvent) GetAddHealthyEvent() *HealthCheckAddHealthy {
	if x, ok := m.GetEvent().(*HealthCheckEvent_AddHealthyEvent); ok {
		return x.AddHealthyEvent
	}
	return nil
}

func (m *HealthCheckEvent) GetTimestamp() *time.Time {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*HealthCheckEvent) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _HealthCheckEvent_OneofMarshaler, _HealthCheckEvent_OneofUnmarshaler, _HealthCheckEvent_OneofSizer, []interface{}{
		(*HealthCheckEvent_EjectUnhealthyEvent)(nil),
		(*HealthCheckEvent_AddHealthyEvent)(nil),
	}
}

func _HealthCheckEvent_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*HealthCheckEvent)
	// event
	switch x := m.Event.(type) {
	case *HealthCheckEvent_EjectUnhealthyEvent:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.EjectUnhealthyEvent); err != nil {
			return err
		}
	case *HealthCheckEvent_AddHealthyEvent:
		_ = b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.AddHealthyEvent); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("HealthCheckEvent.Event has unexpected type %T", x)
	}
	return nil
}

func _HealthCheckEvent_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*HealthCheckEvent)
	switch tag {
	case 4: // event.eject_unhealthy_event
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(HealthCheckEjectUnhealthy)
		err := b.DecodeMessage(msg)
		m.Event = &HealthCheckEvent_EjectUnhealthyEvent{msg}
		return true, err
	case 5: // event.add_healthy_event
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(HealthCheckAddHealthy)
		err := b.DecodeMessage(msg)
		m.Event = &HealthCheckEvent_AddHealthyEvent{msg}
		return true, err
	default:
		return false, nil
	}
}

func _HealthCheckEvent_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*HealthCheckEvent)
	// event
	switch x := m.Event.(type) {
	case *HealthCheckEvent_EjectUnhealthyEvent:
		s := proto.Size(x.EjectUnhealthyEvent)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *HealthCheckEvent_AddHealthyEvent:
		s := proto.Size(x.AddHealthyEvent)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type HealthCheckEjectUnhealthy struct {
	// The type of failure that caused this ejection.
	FailureType          HealthCheckFailureType `protobuf:"varint,1,opt,name=failure_type,json=failureType,proto3,enum=envoy.data.core.v2alpha.HealthCheckFailureType" json:"failure_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *HealthCheckEjectUnhealthy) Reset()         { *m = HealthCheckEjectUnhealthy{} }
func (m *HealthCheckEjectUnhealthy) String() string { return proto.CompactTextString(m) }
func (*HealthCheckEjectUnhealthy) ProtoMessage()    {}
func (*HealthCheckEjectUnhealthy) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_check_event_993d394457bc7a4f, []int{1}
}
func (m *HealthCheckEjectUnhealthy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HealthCheckEjectUnhealthy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HealthCheckEjectUnhealthy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *HealthCheckEjectUnhealthy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckEjectUnhealthy.Merge(dst, src)
}
func (m *HealthCheckEjectUnhealthy) XXX_Size() int {
	return m.Size()
}
func (m *HealthCheckEjectUnhealthy) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckEjectUnhealthy.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckEjectUnhealthy proto.InternalMessageInfo

func (m *HealthCheckEjectUnhealthy) GetFailureType() HealthCheckFailureType {
	if m != nil {
		return m.FailureType
	}
	return HealthCheckFailureType_ACTIVE
}

type HealthCheckAddHealthy struct {
	// Whether this addition is the result of the first ever health check on a host, in which case
	// the configured :ref:`healthy threshold <envoy_api_field_core.HealthCheck.healthy_threshold>`
	// is bypassed and the host is immediately added.
	FirstCheck           bool     `protobuf:"varint,1,opt,name=first_check,json=firstCheck,proto3" json:"first_check,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthCheckAddHealthy) Reset()         { *m = HealthCheckAddHealthy{} }
func (m *HealthCheckAddHealthy) String() string { return proto.CompactTextString(m) }
func (*HealthCheckAddHealthy) ProtoMessage()    {}
func (*HealthCheckAddHealthy) Descriptor() ([]byte, []int) {
	return fileDescriptor_health_check_event_993d394457bc7a4f, []int{2}
}
func (m *HealthCheckAddHealthy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HealthCheckAddHealthy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HealthCheckAddHealthy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *HealthCheckAddHealthy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckAddHealthy.Merge(dst, src)
}
func (m *HealthCheckAddHealthy) XXX_Size() int {
	return m.Size()
}
func (m *HealthCheckAddHealthy) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckAddHealthy.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckAddHealthy proto.InternalMessageInfo

func (m *HealthCheckAddHealthy) GetFirstCheck() bool {
	if m != nil {
		return m.FirstCheck
	}
	return false
}

func init() {
	proto.RegisterType((*HealthCheckEvent)(nil), "envoy.data.core.v2alpha.HealthCheckEvent")
	proto.RegisterType((*HealthCheckEjectUnhealthy)(nil), "envoy.data.core.v2alpha.HealthCheckEjectUnhealthy")
	proto.RegisterType((*HealthCheckAddHealthy)(nil), "envoy.data.core.v2alpha.HealthCheckAddHealthy")
	proto.RegisterEnum("envoy.data.core.v2alpha.HealthCheckFailureType", HealthCheckFailureType_name, HealthCheckFailureType_value)
	proto.RegisterEnum("envoy.data.core.v2alpha.HealthCheckerType", HealthCheckerType_name, HealthCheckerType_value)
}
func (this *HealthCheckEvent) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HealthCheckEvent)
	if !ok {
		that2, ok := that.(HealthCheckEvent)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.HealthCheckerType != that1.HealthCheckerType {
		return false
	}
	if !this.Host.Equal(that1.Host) {
		return false
	}
	if this.ClusterName != that1.ClusterName {
		return false
	}
	if that1.Event == nil {
		if this.Event != nil {
			return false
		}
	} else if this.Event == nil {
		return false
	} else if !this.Event.Equal(that1.Event) {
		return false
	}
	if that1.Timestamp == nil {
		if this.Timestamp != nil {
			return false
		}
	} else if !this.Timestamp.Equal(*that1.Timestamp) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *HealthCheckEvent_EjectUnhealthyEvent) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HealthCheckEvent_EjectUnhealthyEvent)
	if !ok {
		that2, ok := that.(HealthCheckEvent_EjectUnhealthyEvent)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.EjectUnhealthyEvent.Equal(that1.EjectUnhealthyEvent) {
		return false
	}
	return true
}
func (this *HealthCheckEvent_AddHealthyEvent) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HealthCheckEvent_AddHealthyEvent)
	if !ok {
		that2, ok := that.(HealthCheckEvent_AddHealthyEvent)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.AddHealthyEvent.Equal(that1.AddHealthyEvent) {
		return false
	}
	return true
}
func (this *HealthCheckEjectUnhealthy) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HealthCheckEjectUnhealthy)
	if !ok {
		that2, ok := that.(HealthCheckEjectUnhealthy)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.FailureType != that1.FailureType {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *HealthCheckAddHealthy) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*HealthCheckAddHealthy)
	if !ok {
		that2, ok := that.(HealthCheckAddHealthy)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.FirstCheck != that1.FirstCheck {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (m *HealthCheckEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HealthCheckEvent) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.HealthCheckerType != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintHealthCheckEvent(dAtA, i, uint64(m.HealthCheckerType))
	}
	if m.Host != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintHealthCheckEvent(dAtA, i, uint64(m.Host.Size()))
		n1, err := m.Host.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.ClusterName) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintHealthCheckEvent(dAtA, i, uint64(len(m.ClusterName)))
		i += copy(dAtA[i:], m.ClusterName)
	}
	if m.Event != nil {
		nn2, err := m.Event.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn2
	}
	if m.Timestamp != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintHealthCheckEvent(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdTime(*m.Timestamp)))
		n3, err := github_com_gogo_protobuf_types.StdTimeMarshalTo(*m.Timestamp, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *HealthCheckEvent_EjectUnhealthyEvent) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.EjectUnhealthyEvent != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintHealthCheckEvent(dAtA, i, uint64(m.EjectUnhealthyEvent.Size()))
		n4, err := m.EjectUnhealthyEvent.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	return i, nil
}
func (m *HealthCheckEvent_AddHealthyEvent) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.AddHealthyEvent != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintHealthCheckEvent(dAtA, i, uint64(m.AddHealthyEvent.Size()))
		n5, err := m.AddHealthyEvent.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	return i, nil
}
func (m *HealthCheckEjectUnhealthy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HealthCheckEjectUnhealthy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.FailureType != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintHealthCheckEvent(dAtA, i, uint64(m.FailureType))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *HealthCheckAddHealthy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HealthCheckAddHealthy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.FirstCheck {
		dAtA[i] = 0x8
		i++
		if m.FirstCheck {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintHealthCheckEvent(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *HealthCheckEvent) Size() (n int) {
	var l int
	_ = l
	if m.HealthCheckerType != 0 {
		n += 1 + sovHealthCheckEvent(uint64(m.HealthCheckerType))
	}
	if m.Host != nil {
		l = m.Host.Size()
		n += 1 + l + sovHealthCheckEvent(uint64(l))
	}
	l = len(m.ClusterName)
	if l > 0 {
		n += 1 + l + sovHealthCheckEvent(uint64(l))
	}
	if m.Event != nil {
		n += m.Event.Size()
	}
	if m.Timestamp != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdTime(*m.Timestamp)
		n += 1 + l + sovHealthCheckEvent(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *HealthCheckEvent_EjectUnhealthyEvent) Size() (n int) {
	var l int
	_ = l
	if m.EjectUnhealthyEvent != nil {
		l = m.EjectUnhealthyEvent.Size()
		n += 1 + l + sovHealthCheckEvent(uint64(l))
	}
	return n
}
func (m *HealthCheckEvent_AddHealthyEvent) Size() (n int) {
	var l int
	_ = l
	if m.AddHealthyEvent != nil {
		l = m.AddHealthyEvent.Size()
		n += 1 + l + sovHealthCheckEvent(uint64(l))
	}
	return n
}
func (m *HealthCheckEjectUnhealthy) Size() (n int) {
	var l int
	_ = l
	if m.FailureType != 0 {
		n += 1 + sovHealthCheckEvent(uint64(m.FailureType))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *HealthCheckAddHealthy) Size() (n int) {
	var l int
	_ = l
	if m.FirstCheck {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovHealthCheckEvent(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozHealthCheckEvent(x uint64) (n int) {
	return sovHealthCheckEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HealthCheckEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHealthCheckEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HealthCheckEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HealthCheckEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HealthCheckerType", wireType)
			}
			m.HealthCheckerType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HealthCheckerType |= (HealthCheckerType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Host", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Host == nil {
				m.Host = &core.Address{}
			}
			if err := m.Host.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EjectUnhealthyEvent", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &HealthCheckEjectUnhealthy{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Event = &HealthCheckEvent_EjectUnhealthyEvent{v}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddHealthyEvent", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &HealthCheckAddHealthy{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Event = &HealthCheckEvent_AddHealthyEvent{v}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = new(time.Time)
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(m.Timestamp, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHealthCheckEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *HealthCheckEjectUnhealthy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHealthCheckEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HealthCheckEjectUnhealthy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HealthCheckEjectUnhealthy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FailureType", wireType)
			}
			m.FailureType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FailureType |= (HealthCheckFailureType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipHealthCheckEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *HealthCheckAddHealthy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHealthCheckEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: HealthCheckAddHealthy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HealthCheckAddHealthy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FirstCheck", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FirstCheck = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipHealthCheckEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHealthCheckEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipHealthCheckEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHealthCheckEvent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowHealthCheckEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthHealthCheckEvent
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowHealthCheckEvent
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipHealthCheckEvent(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthHealthCheckEvent = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHealthCheckEvent   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/data/core/v2alpha/health_check_event.proto", fileDescriptor_health_check_event_993d394457bc7a4f)
}

var fileDescriptor_health_check_event_993d394457bc7a4f = []byte{
	// 581 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x6e, 0xd3, 0x4a,
	0x14, 0xed, 0xc4, 0x49, 0xdb, 0x4c, 0xaa, 0x3e, 0x77, 0xfa, 0x4a, 0x43, 0x85, 0x92, 0xa8, 0xab,
	0x28, 0x42, 0x36, 0x32, 0x1b, 0x24, 0xa4, 0x4a, 0x49, 0x08, 0xa4, 0x42, 0x2a, 0x91, 0x63, 0x60,
	0x83, 0xb0, 0x26, 0xf1, 0x4d, 0x6c, 0x70, 0x62, 0x6b, 0x3c, 0x31, 0x8a, 0xd8, 0xf1, 0x05, 0x7c,
	0x06, 0x9f, 0x80, 0x58, 0x75, 0xc9, 0x92, 0x3f, 0x00, 0x65, 0xd7, 0x1d, 0x9f, 0x80, 0x3c, 0xe3,
	0xa4, 0x69, 0xd3, 0x48, 0xd9, 0xf9, 0xde, 0x73, 0xcf, 0x39, 0xbe, 0x67, 0x2e, 0x7e, 0x04, 0xe3,
	0x38, 0x98, 0xea, 0x0e, 0xe5, 0x54, 0xef, 0x07, 0x0c, 0xf4, 0xd8, 0xa0, 0x7e, 0xe8, 0x52, 0xdd,
	0x05, 0xea, 0x73, 0xd7, 0xee, 0xbb, 0xd0, 0xff, 0x68, 0x43, 0x0c, 0x63, 0xae, 0x85, 0x2c, 0xe0,
	0x01, 0x39, 0x16, 0x0c, 0x2d, 0x61, 0x68, 0x09, 0x43, 0x4b, 0x19, 0x27, 0x65, 0x29, 0x45, 0x43,
	0x4f, 0x8f, 0x0d, 0x29, 0x46, 0x1d, 0x87, 0x41, 0x14, 0x49, 0xe6, 0xc9, 0x83, 0xd5, 0x81, 0x1e,
	0x8d, 0x20, 0x45, 0x4b, 0xc3, 0x20, 0x18, 0xfa, 0xa0, 0x8b, 0xaa, 0x37, 0x19, 0xe8, 0xce, 0x84,
	0x51, 0xee, 0x05, 0xe3, 0x14, 0x2f, 0xdf, 0xc6, 0xb9, 0x37, 0x82, 0x88, 0xd3, 0x51, 0xb8, 0x4e,
	0xe0, 0x13, 0xa3, 0x61, 0x08, 0x6c, 0x6e, 0x7f, 0x1c, 0x53, 0xdf, 0x73, 0x28, 0x07, 0x7d, 0xfe,
	0x91, 0x02, 0xff, 0x0f, 0x83, 0x61, 0x20, 0x3e, 0xf5, 0xe4, 0x4b, 0x76, 0x4f, 0xff, 0x2a, 0x58,
	0x6d, 0x8b, 0x10, 0x9a, 0x49, 0x06, 0xad, 0x24, 0x02, 0x32, 0xc0, 0x87, 0xcb, 0xc1, 0x00, 0xb3,
	0xf9, 0x34, 0x84, 0x22, 0xaa, 0xa0, 0xea, 0xbe, 0x51, 0xd3, 0xd6, 0x44, 0xa3, 0x2d, 0xe9, 0x00,
	0xb3, 0xa6, 0x21, 0x34, 0xf0, 0x8f, 0xab, 0x4b, 0x25, 0xf7, 0x05, 0x65, 0x54, 0x64, 0x1e, 0xb8,
	0xb7, 0x61, 0xa2, 0xe1, 0xac, 0x1b, 0x44, 0xbc, 0x98, 0xa9, 0xa0, 0x6a, 0xc1, 0x38, 0x49, 0x85,
	0x69, 0xe8, 0x69, 0xb1, 0x21, 0xa5, 0xeb, 0x32, 0x5a, 0x53, 0xcc, 0x91, 0x87, 0x78, 0xaf, 0xef,
	0x4f, 0x22, 0x0e, 0xcc, 0x1e, 0xd3, 0x11, 0x14, 0x95, 0x0a, 0xaa, 0xe6, 0x1b, 0xf9, 0xc4, 0x24,
	0xcb, 0x32, 0x15, 0x64, 0x16, 0x52, 0xf8, 0x82, 0x8e, 0x80, 0xb8, 0xf8, 0x08, 0x3e, 0x40, 0x9f,
	0xdb, 0x93, 0xb1, 0xb4, 0x9e, 0xca, 0x17, 0x2e, 0x66, 0x85, 0x9d, 0xb1, 0xc9, 0x1e, 0xad, 0x44,
	0xe0, 0xf5, 0x9c, 0xdf, 0xde, 0x32, 0x0f, 0xe1, 0x46, 0x47, 0xe6, 0xf5, 0x0e, 0x1f, 0x50, 0xc7,
	0xb1, 0x6f, 0xba, 0xe4, 0x84, 0x8b, 0xb6, 0x89, 0x4b, 0xdd, 0x71, 0xda, 0x0b, 0x87, 0xff, 0xe8,
	0xa2, 0x92, 0xea, 0x67, 0x38, 0xbf, 0x38, 0x82, 0xe2, 0x76, 0x1a, 0x95, 0xbc, 0x02, 0x6d, 0x7e,
	0x05, 0x9a, 0x35, 0x9f, 0x68, 0x64, 0xbf, 0xfe, 0x2e, 0x23, 0xf3, 0x9a, 0xd2, 0xd8, 0xc7, 0x39,
	0xf1, 0x47, 0x24, 0xf7, 0xfd, 0xea, 0x52, 0x41, 0xa7, 0x9f, 0xf1, 0xfd, 0xb5, 0x1b, 0x92, 0xf7,
	0x78, 0x6f, 0x40, 0x3d, 0x7f, 0xc2, 0x60, 0xf9, 0xcd, 0xf5, 0x4d, 0xb6, 0x78, 0x2e, 0x79, 0x2b,
	0x0f, 0x5f, 0x18, 0x5c, 0x03, 0xa7, 0x4f, 0xf0, 0xd1, 0x9d, 0x8b, 0x93, 0x32, 0x2e, 0x0c, 0x3c,
	0x16, 0x71, 0x79, 0x72, 0xc2, 0x77, 0xd7, 0xc4, 0xa2, 0x25, 0x46, 0x6b, 0x67, 0xf8, 0xde, 0xdd,
	0x66, 0x04, 0xe3, 0xed, 0x7a, 0xd3, 0x3a, 0x7f, 0xd3, 0x52, 0xb7, 0x48, 0x01, 0xef, 0x74, 0xea,
	0xdd, 0x6e, 0x52, 0xa0, 0xa4, 0xb8, 0x68, 0x59, 0x6f, 0x5f, 0x99, 0x2f, 0xd5, 0x4c, 0xed, 0x29,
	0x3e, 0x58, 0x39, 0x50, 0xb2, 0x8b, 0xb3, 0x6d, 0xcb, 0xea, 0xa8, 0x5b, 0x64, 0x07, 0x2b, 0x56,
	0xb3, 0xa3, 0xa2, 0xa4, 0xf5, 0xc2, 0xec, 0x34, 0xd5, 0x0c, 0xc9, 0xe3, 0x9c, 0xd9, 0x7a, 0x76,
	0xde, 0x55, 0x95, 0x86, 0xfa, 0x6d, 0x56, 0x42, 0x3f, 0x67, 0x25, 0xf4, 0x6b, 0x56, 0x42, 0x7f,
	0x66, 0x25, 0xd4, 0xdb, 0x16, 0xd1, 0x3f, 0xfe, 0x17, 0x00, 0x00, 0xff, 0xff, 0x45, 0x1a, 0x22,
	0x00, 0x5b, 0x04, 0x00, 0x00,
}
