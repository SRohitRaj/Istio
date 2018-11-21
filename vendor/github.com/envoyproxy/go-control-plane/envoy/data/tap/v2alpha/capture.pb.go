// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/data/tap/v2alpha/capture.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
import types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Connection properties.
type Connection struct {
	// Global unique connection ID for Envoy session. Matches connection IDs used
	// in Envoy logs.
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// Local address.
	LocalAddress *core.Address `protobuf:"bytes,2,opt,name=local_address,json=localAddress" json:"local_address,omitempty"`
	// Remote address.
	RemoteAddress        *core.Address `protobuf:"bytes,3,opt,name=remote_address,json=remoteAddress" json:"remote_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Connection) Reset()         { *m = Connection{} }
func (m *Connection) String() string { return proto.CompactTextString(m) }
func (*Connection) ProtoMessage()    {}
func (*Connection) Descriptor() ([]byte, []int) {
	return fileDescriptor_capture_b7254821859335d7, []int{0}
}
func (m *Connection) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Connection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Connection.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Connection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Connection.Merge(dst, src)
}
func (m *Connection) XXX_Size() int {
	return m.Size()
}
func (m *Connection) XXX_DiscardUnknown() {
	xxx_messageInfo_Connection.DiscardUnknown(m)
}

var xxx_messageInfo_Connection proto.InternalMessageInfo

func (m *Connection) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Connection) GetLocalAddress() *core.Address {
	if m != nil {
		return m.LocalAddress
	}
	return nil
}

func (m *Connection) GetRemoteAddress() *core.Address {
	if m != nil {
		return m.RemoteAddress
	}
	return nil
}

// Event in a capture trace.
type Event struct {
	// Timestamp for event.
	Timestamp *types.Timestamp `protobuf:"bytes,1,opt,name=timestamp" json:"timestamp,omitempty"`
	// Read or write with content as bytes string.
	//
	// Types that are valid to be assigned to EventSelector:
	//	*Event_Read_
	//	*Event_Write_
	EventSelector        isEvent_EventSelector `protobuf_oneof:"event_selector"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_capture_b7254821859335d7, []int{1}
}
func (m *Event) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Event.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(dst, src)
}
func (m *Event) XXX_Size() int {
	return m.Size()
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

type isEvent_EventSelector interface {
	isEvent_EventSelector()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Event_Read_ struct {
	Read *Event_Read `protobuf:"bytes,2,opt,name=read,oneof"`
}
type Event_Write_ struct {
	Write *Event_Write `protobuf:"bytes,3,opt,name=write,oneof"`
}

func (*Event_Read_) isEvent_EventSelector()  {}
func (*Event_Write_) isEvent_EventSelector() {}

func (m *Event) GetEventSelector() isEvent_EventSelector {
	if m != nil {
		return m.EventSelector
	}
	return nil
}

func (m *Event) GetTimestamp() *types.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Event) GetRead() *Event_Read {
	if x, ok := m.GetEventSelector().(*Event_Read_); ok {
		return x.Read
	}
	return nil
}

func (m *Event) GetWrite() *Event_Write {
	if x, ok := m.GetEventSelector().(*Event_Write_); ok {
		return x.Write
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Event) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Event_OneofMarshaler, _Event_OneofUnmarshaler, _Event_OneofSizer, []interface{}{
		(*Event_Read_)(nil),
		(*Event_Write_)(nil),
	}
}

func _Event_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Event)
	// event_selector
	switch x := m.EventSelector.(type) {
	case *Event_Read_:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Read); err != nil {
			return err
		}
	case *Event_Write_:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Write); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Event.EventSelector has unexpected type %T", x)
	}
	return nil
}

func _Event_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Event)
	switch tag {
	case 2: // event_selector.read
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Event_Read)
		err := b.DecodeMessage(msg)
		m.EventSelector = &Event_Read_{msg}
		return true, err
	case 3: // event_selector.write
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Event_Write)
		err := b.DecodeMessage(msg)
		m.EventSelector = &Event_Write_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Event_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Event)
	// event_selector
	switch x := m.EventSelector.(type) {
	case *Event_Read_:
		s := proto.Size(x.Read)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Event_Write_:
		s := proto.Size(x.Write)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Data read by Envoy from the transport socket.
type Event_Read struct {
	// Binary data read.
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event_Read) Reset()         { *m = Event_Read{} }
func (m *Event_Read) String() string { return proto.CompactTextString(m) }
func (*Event_Read) ProtoMessage()    {}
func (*Event_Read) Descriptor() ([]byte, []int) {
	return fileDescriptor_capture_b7254821859335d7, []int{1, 0}
}
func (m *Event_Read) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Event_Read) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Event_Read.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Event_Read) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_Read.Merge(dst, src)
}
func (m *Event_Read) XXX_Size() int {
	return m.Size()
}
func (m *Event_Read) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_Read.DiscardUnknown(m)
}

var xxx_messageInfo_Event_Read proto.InternalMessageInfo

func (m *Event_Read) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// Data written by Envoy to the transport socket.
type Event_Write struct {
	// Binary data written.
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	// Stream was half closed after this write.
	EndStream            bool     `protobuf:"varint,2,opt,name=end_stream,json=endStream,proto3" json:"end_stream,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Event_Write) Reset()         { *m = Event_Write{} }
func (m *Event_Write) String() string { return proto.CompactTextString(m) }
func (*Event_Write) ProtoMessage()    {}
func (*Event_Write) Descriptor() ([]byte, []int) {
	return fileDescriptor_capture_b7254821859335d7, []int{1, 1}
}
func (m *Event_Write) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Event_Write) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Event_Write.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Event_Write) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event_Write.Merge(dst, src)
}
func (m *Event_Write) XXX_Size() int {
	return m.Size()
}
func (m *Event_Write) XXX_DiscardUnknown() {
	xxx_messageInfo_Event_Write.DiscardUnknown(m)
}

var xxx_messageInfo_Event_Write proto.InternalMessageInfo

func (m *Event_Write) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Event_Write) GetEndStream() bool {
	if m != nil {
		return m.EndStream
	}
	return false
}

// Sequence of read/write events that constitute a captured trace on a socket.
// Multiple Trace messages might be emitted for a given connection ID, with the
// sink (e.g. file set, network) responsible for later reassembly.
type Trace struct {
	// Connection properties.
	Connection *Connection `protobuf:"bytes,1,opt,name=connection" json:"connection,omitempty"`
	// Sequence of observed events.
	Events               []*Event `protobuf:"bytes,2,rep,name=events" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Trace) Reset()         { *m = Trace{} }
func (m *Trace) String() string { return proto.CompactTextString(m) }
func (*Trace) ProtoMessage()    {}
func (*Trace) Descriptor() ([]byte, []int) {
	return fileDescriptor_capture_b7254821859335d7, []int{2}
}
func (m *Trace) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Trace) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Trace.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Trace) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Trace.Merge(dst, src)
}
func (m *Trace) XXX_Size() int {
	return m.Size()
}
func (m *Trace) XXX_DiscardUnknown() {
	xxx_messageInfo_Trace.DiscardUnknown(m)
}

var xxx_messageInfo_Trace proto.InternalMessageInfo

func (m *Trace) GetConnection() *Connection {
	if m != nil {
		return m.Connection
	}
	return nil
}

func (m *Trace) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func init() {
	proto.RegisterType((*Connection)(nil), "envoy.data.tap.v2alpha.Connection")
	proto.RegisterType((*Event)(nil), "envoy.data.tap.v2alpha.Event")
	proto.RegisterType((*Event_Read)(nil), "envoy.data.tap.v2alpha.Event.Read")
	proto.RegisterType((*Event_Write)(nil), "envoy.data.tap.v2alpha.Event.Write")
	proto.RegisterType((*Trace)(nil), "envoy.data.tap.v2alpha.Trace")
}
func (m *Connection) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Connection) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.Id))
	}
	if m.LocalAddress != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.LocalAddress.Size()))
		n1, err := m.LocalAddress.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.RemoteAddress != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.RemoteAddress.Size()))
		n2, err := m.RemoteAddress.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Event) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Event) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Timestamp != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.Timestamp.Size()))
		n3, err := m.Timestamp.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.EventSelector != nil {
		nn4, err := m.EventSelector.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn4
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Event_Read_) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Read != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.Read.Size()))
		n5, err := m.Read.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	return i, nil
}
func (m *Event_Write_) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Write != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.Write.Size()))
		n6, err := m.Write.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	return i, nil
}
func (m *Event_Read) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Event_Read) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCapture(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *Event_Write) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Event_Write) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Data) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCapture(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	if m.EndStream {
		dAtA[i] = 0x10
		i++
		if m.EndStream {
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

func (m *Trace) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Trace) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Connection != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCapture(dAtA, i, uint64(m.Connection.Size()))
		n7, err := m.Connection.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	if len(m.Events) > 0 {
		for _, msg := range m.Events {
			dAtA[i] = 0x12
			i++
			i = encodeVarintCapture(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintCapture(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Connection) Size() (n int) {
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovCapture(uint64(m.Id))
	}
	if m.LocalAddress != nil {
		l = m.LocalAddress.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	if m.RemoteAddress != nil {
		l = m.RemoteAddress.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Event) Size() (n int) {
	var l int
	_ = l
	if m.Timestamp != nil {
		l = m.Timestamp.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	if m.EventSelector != nil {
		n += m.EventSelector.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Event_Read_) Size() (n int) {
	var l int
	_ = l
	if m.Read != nil {
		l = m.Read.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	return n
}
func (m *Event_Write_) Size() (n int) {
	var l int
	_ = l
	if m.Write != nil {
		l = m.Write.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	return n
}
func (m *Event_Read) Size() (n int) {
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovCapture(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Event_Write) Size() (n int) {
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovCapture(uint64(l))
	}
	if m.EndStream {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Trace) Size() (n int) {
	var l int
	_ = l
	if m.Connection != nil {
		l = m.Connection.Size()
		n += 1 + l + sovCapture(uint64(l))
	}
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovCapture(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovCapture(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCapture(x uint64) (n int) {
	return sovCapture(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Connection) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapture
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
			return fmt.Errorf("proto: Connection: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Connection: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LocalAddress", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.LocalAddress == nil {
				m.LocalAddress = &core.Address{}
			}
			if err := m.LocalAddress.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RemoteAddress", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RemoteAddress == nil {
				m.RemoteAddress = &core.Address{}
			}
			if err := m.RemoteAddress.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCapture(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCapture
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
func (m *Event) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapture
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
			return fmt.Errorf("proto: Event: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Event: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Timestamp == nil {
				m.Timestamp = &types.Timestamp{}
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Read", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &Event_Read{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.EventSelector = &Event_Read_{v}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Write", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &Event_Write{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.EventSelector = &Event_Write_{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCapture(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCapture
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
func (m *Event_Read) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapture
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
			return fmt.Errorf("proto: Read: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Read: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCapture(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCapture
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
func (m *Event_Write) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapture
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
			return fmt.Errorf("proto: Write: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Write: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndStream", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
			m.EndStream = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipCapture(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCapture
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
func (m *Trace) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapture
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
			return fmt.Errorf("proto: Trace: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Trace: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Connection", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Connection == nil {
				m.Connection = &Connection{}
			}
			if err := m.Connection.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapture
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
				return ErrInvalidLengthCapture
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, &Event{})
			if err := m.Events[len(m.Events)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCapture(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCapture
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
func skipCapture(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCapture
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
					return 0, ErrIntOverflowCapture
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
					return 0, ErrIntOverflowCapture
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
				return 0, ErrInvalidLengthCapture
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCapture
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
				next, err := skipCapture(dAtA[start:])
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
	ErrInvalidLengthCapture = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCapture   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/data/tap/v2alpha/capture.proto", fileDescriptor_capture_b7254821859335d7)
}

var fileDescriptor_capture_b7254821859335d7 = []byte{
	// 413 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xcf, 0x6a, 0x14, 0x41,
	0x10, 0xc6, 0xd3, 0x93, 0xdd, 0x60, 0x2a, 0xc9, 0x12, 0xfa, 0x20, 0x61, 0x20, 0x9b, 0xb0, 0x7a,
	0xd8, 0x53, 0x37, 0x8c, 0x08, 0x41, 0x0f, 0x92, 0x15, 0x21, 0xe7, 0x36, 0x20, 0x78, 0x59, 0x2a,
	0xd3, 0x65, 0x1c, 0x98, 0x9d, 0x6e, 0x7a, 0x3a, 0x23, 0x5e, 0x7d, 0x12, 0xf1, 0x69, 0x3c, 0xfa,
	0x08, 0xb2, 0x4f, 0x22, 0xfd, 0x27, 0x1b, 0x0f, 0x71, 0xbd, 0x4d, 0x4d, 0xff, 0xbe, 0xaf, 0xbe,
	0xaa, 0x82, 0xe7, 0xd4, 0x0d, 0xe6, 0xab, 0xd4, 0xe8, 0x51, 0x7a, 0xb4, 0x72, 0xa8, 0xb0, 0xb5,
	0x9f, 0x51, 0xd6, 0x68, 0xfd, 0x9d, 0x23, 0x61, 0x9d, 0xf1, 0x86, 0x3f, 0x8d, 0x94, 0x08, 0x94,
	0xf0, 0x68, 0x45, 0xa6, 0xca, 0xb3, 0xa4, 0x46, 0xdb, 0xc8, 0xa1, 0x92, 0xb5, 0x71, 0x24, 0x51,
	0x6b, 0x47, 0x7d, 0x9f, 0x84, 0xe5, 0xd9, 0xad, 0x31, 0xb7, 0x2d, 0xc9, 0x58, 0xdd, 0xdc, 0x7d,
	0x92, 0xbe, 0x59, 0x51, 0xef, 0x71, 0x65, 0x13, 0x30, 0xfb, 0xce, 0x00, 0xde, 0x9a, 0xae, 0xa3,
	0xda, 0x37, 0xa6, 0xe3, 0x13, 0x28, 0x1a, 0x7d, 0xc2, 0xce, 0xd9, 0x7c, 0xa4, 0x8a, 0x46, 0xf3,
	0x37, 0x70, 0xd4, 0x9a, 0x1a, 0xdb, 0x65, 0xb6, 0x3d, 0x29, 0xce, 0xd9, 0xfc, 0xa0, 0x2a, 0x45,
	0x0a, 0x84, 0xb6, 0x11, 0x43, 0x25, 0x42, 0x63, 0x71, 0x99, 0x08, 0x75, 0x18, 0x05, 0xb9, 0xe2,
	0x97, 0x30, 0x71, 0xb4, 0x32, 0x9e, 0x36, 0x0e, 0xbb, 0xff, 0x75, 0x38, 0x4a, 0x8a, 0x5c, 0xce,
	0x7e, 0x14, 0x30, 0x7e, 0x37, 0x50, 0xe7, 0xf9, 0x05, 0xec, 0x6f, 0xf2, 0xc7, 0x90, 0xc1, 0x27,
	0x4d, 0x28, 0xee, 0x27, 0x14, 0xd7, 0xf7, 0x84, 0x7a, 0x80, 0xf9, 0x05, 0x8c, 0x1c, 0xa1, 0xce,
	0xf1, 0x67, 0xe2, 0xf1, 0x7d, 0x8a, 0xd8, 0x46, 0x28, 0x42, 0x7d, 0xb5, 0xa3, 0xa2, 0x82, 0xbf,
	0x86, 0xf1, 0x17, 0xd7, 0x78, 0xca, 0xb9, 0x9f, 0x6d, 0x97, 0x7e, 0x08, 0xe8, 0xd5, 0x8e, 0x4a,
	0x9a, 0xb2, 0x84, 0x51, 0x30, 0xe3, 0x1c, 0x46, 0x41, 0x10, 0x33, 0x1f, 0xaa, 0xf8, 0x5d, 0xbe,
	0x82, 0x71, 0xa4, 0x1f, 0x7b, 0xe4, 0xa7, 0x00, 0xd4, 0xe9, 0x65, 0xef, 0x1d, 0xe1, 0x2a, 0xa6,
	0x7e, 0xa2, 0xf6, 0xa9, 0xd3, 0xef, 0xe3, 0x8f, 0xc5, 0x31, 0x4c, 0x28, 0xf4, 0x5b, 0xf6, 0xd4,
	0x52, 0xed, 0x8d, 0x9b, 0x7d, 0x63, 0x30, 0xbe, 0x76, 0x58, 0x13, 0x5f, 0x00, 0xd4, 0x9b, 0x83,
	0xe6, 0x2d, 0xfd, 0x73, 0xe0, 0x87, 0xd3, 0xab, 0xbf, 0x54, 0xfc, 0x25, 0xec, 0x45, 0xff, 0x70,
	0xef, 0xdd, 0xf9, 0x41, 0x75, 0xba, 0x75, 0x6a, 0x95, 0xe1, 0xc5, 0xf1, 0xcf, 0xf5, 0x94, 0xfd,
	0x5a, 0x4f, 0xd9, 0xef, 0xf5, 0x94, 0x7d, 0x2c, 0x86, 0xea, 0x66, 0x2f, 0x9e, 0xe5, 0xc5, 0x9f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x15, 0x47, 0x26, 0xfe, 0xe7, 0x02, 0x00, 0x00,
}
