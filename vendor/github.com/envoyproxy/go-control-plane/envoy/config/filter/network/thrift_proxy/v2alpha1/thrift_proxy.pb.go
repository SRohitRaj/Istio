// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/network/thrift_proxy/v2alpha1/thrift_proxy.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import types "github.com/gogo/protobuf/types"
import _ "github.com/lyft/protoc-gen-validate/validate"

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

// Thrift transport types supported by Envoy.
type TransportType int32

const (
	// For downstream connections, the Thrift proxy will attempt to determine which transport to use.
	// For upstream connections, the Thrift proxy will use same transport as the downstream
	// connection.
	AUTO_TRANSPORT TransportType = 0
	// The Thrift proxy will use the Thrift framed transport.
	FRAMED TransportType = 1
	// The Thrift proxy will use the Thrift unframed transport.
	UNFRAMED TransportType = 2
	// The Thrift proxy will assume the client is using the Thrift header transport.
	HEADER TransportType = 3
)

var TransportType_name = map[int32]string{
	0: "AUTO_TRANSPORT",
	1: "FRAMED",
	2: "UNFRAMED",
	3: "HEADER",
}
var TransportType_value = map[string]int32{
	"AUTO_TRANSPORT": 0,
	"FRAMED":         1,
	"UNFRAMED":       2,
	"HEADER":         3,
}

func (x TransportType) String() string {
	return proto.EnumName(TransportType_name, int32(x))
}
func (TransportType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_thrift_proxy_3375db5f87c43dda, []int{0}
}

// Thrift Protocol types supported by Envoy.
type ProtocolType int32

const (
	// For downstream connections, the Thrift proxy will attempt to determine which protocol to use.
	// Note that the older, non-strict (or lax) binary protocol is not included in automatic protocol
	// detection. For upstream connections, the Thrift proxy will use the same protocol as the
	// downstream connection.
	AUTO_PROTOCOL ProtocolType = 0
	// The Thrift proxy will use the Thrift binary protocol.
	BINARY ProtocolType = 1
	// The Thrift proxy will use Thrift non-strict binary protocol.
	LAX_BINARY ProtocolType = 2
	// The Thrift proxy will use the Thrift compact protocol.
	COMPACT ProtocolType = 3
	// The Thrift proxy will use the Thrift "Twitter" protocol implemented by the finagle library.
	TWITTER ProtocolType = 4
)

var ProtocolType_name = map[int32]string{
	0: "AUTO_PROTOCOL",
	1: "BINARY",
	2: "LAX_BINARY",
	3: "COMPACT",
	4: "TWITTER",
}
var ProtocolType_value = map[string]int32{
	"AUTO_PROTOCOL": 0,
	"BINARY":        1,
	"LAX_BINARY":    2,
	"COMPACT":       3,
	"TWITTER":       4,
}

func (x ProtocolType) String() string {
	return proto.EnumName(ProtocolType_name, int32(x))
}
func (ProtocolType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_thrift_proxy_3375db5f87c43dda, []int{1}
}

// [#comment:next free field: 6]
type ThriftProxy struct {
	// Supplies the type of transport that the Thrift proxy should use. Defaults to
	// :ref:`AUTO_TRANSPORT<envoy_api_enum_value_config.filter.network.thrift_proxy.v2alpha1.TransportType.AUTO_TRANSPORT>`.
	Transport TransportType `protobuf:"varint,2,opt,name=transport,proto3,enum=envoy.config.filter.network.thrift_proxy.v2alpha1.TransportType" json:"transport,omitempty"`
	// Supplies the type of protocol that the Thrift proxy should use. Defaults to
	// :ref:`AUTO_PROTOCOL<envoy_api_enum_value_config.filter.network.thrift_proxy.v2alpha1.ProtocolType.AUTO_PROTOCOL>`.
	Protocol ProtocolType `protobuf:"varint,3,opt,name=protocol,proto3,enum=envoy.config.filter.network.thrift_proxy.v2alpha1.ProtocolType" json:"protocol,omitempty"`
	// The human readable prefix to use when emitting statistics.
	StatPrefix string `protobuf:"bytes,1,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	// The route table for the connection manager is static and is specified in this property.
	RouteConfig *RouteConfiguration `protobuf:"bytes,4,opt,name=route_config,json=routeConfig,proto3" json:"route_config,omitempty"`
	// A list of individual Thrift filters that make up the filter chain for requests made to the
	// Thrift proxy. Order matters as the filters are processed sequentially. For backwards
	// compatibility, if no thrift_filters are specified, a default Thrift router filter
	// (`envoy.filters.thrift.router`) is used.
	ThriftFilters        []*ThriftFilter `protobuf:"bytes,5,rep,name=thrift_filters,json=thriftFilters,proto3" json:"thrift_filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ThriftProxy) Reset()         { *m = ThriftProxy{} }
func (m *ThriftProxy) String() string { return proto.CompactTextString(m) }
func (*ThriftProxy) ProtoMessage()    {}
func (*ThriftProxy) Descriptor() ([]byte, []int) {
	return fileDescriptor_thrift_proxy_3375db5f87c43dda, []int{0}
}
func (m *ThriftProxy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ThriftProxy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ThriftProxy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ThriftProxy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ThriftProxy.Merge(dst, src)
}
func (m *ThriftProxy) XXX_Size() int {
	return m.Size()
}
func (m *ThriftProxy) XXX_DiscardUnknown() {
	xxx_messageInfo_ThriftProxy.DiscardUnknown(m)
}

var xxx_messageInfo_ThriftProxy proto.InternalMessageInfo

func (m *ThriftProxy) GetTransport() TransportType {
	if m != nil {
		return m.Transport
	}
	return AUTO_TRANSPORT
}

func (m *ThriftProxy) GetProtocol() ProtocolType {
	if m != nil {
		return m.Protocol
	}
	return AUTO_PROTOCOL
}

func (m *ThriftProxy) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *ThriftProxy) GetRouteConfig() *RouteConfiguration {
	if m != nil {
		return m.RouteConfig
	}
	return nil
}

func (m *ThriftProxy) GetThriftFilters() []*ThriftFilter {
	if m != nil {
		return m.ThriftFilters
	}
	return nil
}

// ThriftFilter configures a Thrift filter.
// [#comment:next free field: 3]
type ThriftFilter struct {
	// The name of the filter to instantiate. The name must match a supported
	// filter. The built-in filters are:
	//
	// [#comment:TODO(zuercher): Auto generate the following list]
	// * :ref:`envoy.filters.thrift.router <config_thrift_filters_router>`
	// * :ref:`envoy.filters.thrift.rate_limit <config_thrift_filters_rate_limit>`
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Filter specific configuration which depends on the filter being instantiated. See the supported
	// filters for further documentation.
	//
	// Types that are valid to be assigned to ConfigType:
	//	*ThriftFilter_Config
	//	*ThriftFilter_TypedConfig
	ConfigType           isThriftFilter_ConfigType `protobuf_oneof:"config_type"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *ThriftFilter) Reset()         { *m = ThriftFilter{} }
func (m *ThriftFilter) String() string { return proto.CompactTextString(m) }
func (*ThriftFilter) ProtoMessage()    {}
func (*ThriftFilter) Descriptor() ([]byte, []int) {
	return fileDescriptor_thrift_proxy_3375db5f87c43dda, []int{1}
}
func (m *ThriftFilter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ThriftFilter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ThriftFilter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ThriftFilter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ThriftFilter.Merge(dst, src)
}
func (m *ThriftFilter) XXX_Size() int {
	return m.Size()
}
func (m *ThriftFilter) XXX_DiscardUnknown() {
	xxx_messageInfo_ThriftFilter.DiscardUnknown(m)
}

var xxx_messageInfo_ThriftFilter proto.InternalMessageInfo

type isThriftFilter_ConfigType interface {
	isThriftFilter_ConfigType()
	MarshalTo([]byte) (int, error)
	Size() int
}

type ThriftFilter_Config struct {
	Config *types.Struct `protobuf:"bytes,2,opt,name=config,proto3,oneof"`
}
type ThriftFilter_TypedConfig struct {
	TypedConfig *types.Any `protobuf:"bytes,3,opt,name=typed_config,json=typedConfig,proto3,oneof"`
}

func (*ThriftFilter_Config) isThriftFilter_ConfigType()      {}
func (*ThriftFilter_TypedConfig) isThriftFilter_ConfigType() {}

func (m *ThriftFilter) GetConfigType() isThriftFilter_ConfigType {
	if m != nil {
		return m.ConfigType
	}
	return nil
}

func (m *ThriftFilter) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ThriftFilter) GetConfig() *types.Struct {
	if x, ok := m.GetConfigType().(*ThriftFilter_Config); ok {
		return x.Config
	}
	return nil
}

func (m *ThriftFilter) GetTypedConfig() *types.Any {
	if x, ok := m.GetConfigType().(*ThriftFilter_TypedConfig); ok {
		return x.TypedConfig
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ThriftFilter) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ThriftFilter_OneofMarshaler, _ThriftFilter_OneofUnmarshaler, _ThriftFilter_OneofSizer, []interface{}{
		(*ThriftFilter_Config)(nil),
		(*ThriftFilter_TypedConfig)(nil),
	}
}

func _ThriftFilter_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ThriftFilter)
	// config_type
	switch x := m.ConfigType.(type) {
	case *ThriftFilter_Config:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Config); err != nil {
			return err
		}
	case *ThriftFilter_TypedConfig:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.TypedConfig); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ThriftFilter.ConfigType has unexpected type %T", x)
	}
	return nil
}

func _ThriftFilter_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ThriftFilter)
	switch tag {
	case 2: // config_type.config
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(types.Struct)
		err := b.DecodeMessage(msg)
		m.ConfigType = &ThriftFilter_Config{msg}
		return true, err
	case 3: // config_type.typed_config
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(types.Any)
		err := b.DecodeMessage(msg)
		m.ConfigType = &ThriftFilter_TypedConfig{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ThriftFilter_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ThriftFilter)
	// config_type
	switch x := m.ConfigType.(type) {
	case *ThriftFilter_Config:
		s := proto.Size(x.Config)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ThriftFilter_TypedConfig:
		s := proto.Size(x.TypedConfig)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// ThriftProtocolOptions specifies Thrift upstream protocol options. This object is used in
// in :ref:`extension_protocol_options<envoy_api_field_Cluster.extension_protocol_options>`, keyed
// by the name `envoy.filters.network.thrift_proxy`.
// [#comment:next free field: 3]
type ThriftProtocolOptions struct {
	// Supplies the type of transport that the Thrift proxy should use for upstream connections.
	// Selecting
	// :ref:`AUTO_TRANSPORT<envoy_api_enum_value_config.filter.network.thrift_proxy.v2alpha1.TransportType.AUTO_TRANSPORT>`,
	// which is the default, causes the proxy to use the same transport as the downstream connection.
	Transport TransportType `protobuf:"varint,1,opt,name=transport,proto3,enum=envoy.config.filter.network.thrift_proxy.v2alpha1.TransportType" json:"transport,omitempty"`
	// Supplies the type of protocol that the Thrift proxy should use for upstream connections.
	// Selecting
	// :ref:`AUTO_PROTOCOL<envoy_api_enum_value_config.filter.network.thrift_proxy.v2alpha1.ProtocolType.AUTO_PROTOCOL>`,
	// which is the default, causes the proxy to use the same protocol as the downstream connection.
	Protocol             ProtocolType `protobuf:"varint,2,opt,name=protocol,proto3,enum=envoy.config.filter.network.thrift_proxy.v2alpha1.ProtocolType" json:"protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ThriftProtocolOptions) Reset()         { *m = ThriftProtocolOptions{} }
func (m *ThriftProtocolOptions) String() string { return proto.CompactTextString(m) }
func (*ThriftProtocolOptions) ProtoMessage()    {}
func (*ThriftProtocolOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_thrift_proxy_3375db5f87c43dda, []int{2}
}
func (m *ThriftProtocolOptions) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ThriftProtocolOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ThriftProtocolOptions.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ThriftProtocolOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ThriftProtocolOptions.Merge(dst, src)
}
func (m *ThriftProtocolOptions) XXX_Size() int {
	return m.Size()
}
func (m *ThriftProtocolOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_ThriftProtocolOptions.DiscardUnknown(m)
}

var xxx_messageInfo_ThriftProtocolOptions proto.InternalMessageInfo

func (m *ThriftProtocolOptions) GetTransport() TransportType {
	if m != nil {
		return m.Transport
	}
	return AUTO_TRANSPORT
}

func (m *ThriftProtocolOptions) GetProtocol() ProtocolType {
	if m != nil {
		return m.Protocol
	}
	return AUTO_PROTOCOL
}

func init() {
	proto.RegisterType((*ThriftProxy)(nil), "envoy.config.filter.network.thrift_proxy.v2alpha1.ThriftProxy")
	proto.RegisterType((*ThriftFilter)(nil), "envoy.config.filter.network.thrift_proxy.v2alpha1.ThriftFilter")
	proto.RegisterType((*ThriftProtocolOptions)(nil), "envoy.config.filter.network.thrift_proxy.v2alpha1.ThriftProtocolOptions")
	proto.RegisterEnum("envoy.config.filter.network.thrift_proxy.v2alpha1.TransportType", TransportType_name, TransportType_value)
	proto.RegisterEnum("envoy.config.filter.network.thrift_proxy.v2alpha1.ProtocolType", ProtocolType_name, ProtocolType_value)
}
func (m *ThriftProxy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ThriftProxy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.StatPrefix) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(len(m.StatPrefix)))
		i += copy(dAtA[i:], m.StatPrefix)
	}
	if m.Transport != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(m.Transport))
	}
	if m.Protocol != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(m.Protocol))
	}
	if m.RouteConfig != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(m.RouteConfig.Size()))
		n1, err := m.RouteConfig.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.ThriftFilters) > 0 {
		for _, msg := range m.ThriftFilters {
			dAtA[i] = 0x2a
			i++
			i = encodeVarintThriftProxy(dAtA, i, uint64(msg.Size()))
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

func (m *ThriftFilter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ThriftFilter) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.ConfigType != nil {
		nn2, err := m.ConfigType.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn2
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *ThriftFilter_Config) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Config != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(m.Config.Size()))
		n3, err := m.Config.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}
func (m *ThriftFilter_TypedConfig) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.TypedConfig != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(m.TypedConfig.Size()))
		n4, err := m.TypedConfig.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	return i, nil
}
func (m *ThriftProtocolOptions) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ThriftProtocolOptions) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Transport != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(m.Transport))
	}
	if m.Protocol != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintThriftProxy(dAtA, i, uint64(m.Protocol))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintThriftProxy(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ThriftProxy) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.StatPrefix)
	if l > 0 {
		n += 1 + l + sovThriftProxy(uint64(l))
	}
	if m.Transport != 0 {
		n += 1 + sovThriftProxy(uint64(m.Transport))
	}
	if m.Protocol != 0 {
		n += 1 + sovThriftProxy(uint64(m.Protocol))
	}
	if m.RouteConfig != nil {
		l = m.RouteConfig.Size()
		n += 1 + l + sovThriftProxy(uint64(l))
	}
	if len(m.ThriftFilters) > 0 {
		for _, e := range m.ThriftFilters {
			l = e.Size()
			n += 1 + l + sovThriftProxy(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ThriftFilter) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovThriftProxy(uint64(l))
	}
	if m.ConfigType != nil {
		n += m.ConfigType.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ThriftFilter_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Config != nil {
		l = m.Config.Size()
		n += 1 + l + sovThriftProxy(uint64(l))
	}
	return n
}
func (m *ThriftFilter_TypedConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TypedConfig != nil {
		l = m.TypedConfig.Size()
		n += 1 + l + sovThriftProxy(uint64(l))
	}
	return n
}
func (m *ThriftProtocolOptions) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Transport != 0 {
		n += 1 + sovThriftProxy(uint64(m.Transport))
	}
	if m.Protocol != 0 {
		n += 1 + sovThriftProxy(uint64(m.Protocol))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovThriftProxy(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozThriftProxy(x uint64) (n int) {
	return sovThriftProxy(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ThriftProxy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowThriftProxy
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
			return fmt.Errorf("proto: ThriftProxy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ThriftProxy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
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
				return ErrInvalidLengthThriftProxy
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transport", wireType)
			}
			m.Transport = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Transport |= (TransportType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Protocol", wireType)
			}
			m.Protocol = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Protocol |= (ProtocolType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RouteConfig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
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
				return ErrInvalidLengthThriftProxy
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RouteConfig == nil {
				m.RouteConfig = &RouteConfiguration{}
			}
			if err := m.RouteConfig.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ThriftFilters", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
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
				return ErrInvalidLengthThriftProxy
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ThriftFilters = append(m.ThriftFilters, &ThriftFilter{})
			if err := m.ThriftFilters[len(m.ThriftFilters)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipThriftProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthThriftProxy
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
func (m *ThriftFilter) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowThriftProxy
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
			return fmt.Errorf("proto: ThriftFilter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ThriftFilter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
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
				return ErrInvalidLengthThriftProxy
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
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
				return ErrInvalidLengthThriftProxy
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &types.Struct{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.ConfigType = &ThriftFilter_Config{v}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TypedConfig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
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
				return ErrInvalidLengthThriftProxy
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &types.Any{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.ConfigType = &ThriftFilter_TypedConfig{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipThriftProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthThriftProxy
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
func (m *ThriftProtocolOptions) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowThriftProxy
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
			return fmt.Errorf("proto: ThriftProtocolOptions: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ThriftProtocolOptions: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transport", wireType)
			}
			m.Transport = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Transport |= (TransportType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Protocol", wireType)
			}
			m.Protocol = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowThriftProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Protocol |= (ProtocolType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipThriftProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthThriftProxy
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
func skipThriftProxy(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowThriftProxy
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
					return 0, ErrIntOverflowThriftProxy
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
					return 0, ErrIntOverflowThriftProxy
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
				return 0, ErrInvalidLengthThriftProxy
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowThriftProxy
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
				next, err := skipThriftProxy(dAtA[start:])
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
	ErrInvalidLengthThriftProxy = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowThriftProxy   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/filter/network/thrift_proxy/v2alpha1/thrift_proxy.proto", fileDescriptor_thrift_proxy_3375db5f87c43dda)
}

var fileDescriptor_thrift_proxy_3375db5f87c43dda = []byte{
	// 590 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x92, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xc7, 0xb3, 0x76, 0x5a, 0xda, 0x71, 0x12, 0x99, 0x55, 0x51, 0x43, 0x04, 0x51, 0xd4, 0x53,
	0xd4, 0x83, 0xad, 0x86, 0x13, 0x07, 0x04, 0x4e, 0x9a, 0xaa, 0x95, 0xda, 0xda, 0xda, 0xba, 0xe2,
	0x43, 0x42, 0x91, 0x9b, 0xda, 0x8e, 0x21, 0x78, 0xad, 0xf5, 0x26, 0x34, 0x57, 0x4e, 0xbc, 0x07,
	0x37, 0x1e, 0x81, 0x13, 0x47, 0x8e, 0x3c, 0x42, 0x95, 0x1b, 0x12, 0x0f, 0x81, 0xbc, 0xeb, 0xb4,
	0x09, 0x39, 0x35, 0x12, 0xdc, 0x76, 0x66, 0x76, 0x7e, 0xff, 0xf9, 0x82, 0x7d, 0x3f, 0x1e, 0xd3,
	0x89, 0xd9, 0xa7, 0x71, 0x10, 0x85, 0x66, 0x10, 0x0d, 0xb9, 0xcf, 0xcc, 0xd8, 0xe7, 0x1f, 0x29,
	0x7b, 0x6f, 0xf2, 0x01, 0x8b, 0x02, 0xde, 0x4b, 0x18, 0xbd, 0x9a, 0x98, 0xe3, 0x96, 0x37, 0x4c,
	0x06, 0xde, 0xde, 0x82, 0xd7, 0x48, 0x18, 0xe5, 0x14, 0xef, 0x09, 0x8a, 0x21, 0x29, 0x86, 0xa4,
	0x18, 0x39, 0xc5, 0x58, 0xf8, 0x3f, 0xa3, 0xd4, 0x9e, 0xdd, 0x5d, 0x98, 0xd1, 0x11, 0xf7, 0xa5,
	0x62, 0xed, 0x61, 0x48, 0x69, 0x38, 0xf4, 0x4d, 0x61, 0x5d, 0x8c, 0x02, 0xd3, 0x8b, 0xf3, 0x62,
	0x6a, 0x8f, 0xfe, 0x0e, 0xa5, 0x9c, 0x8d, 0xfa, 0x3c, 0x8f, 0x6e, 0x8f, 0xbd, 0x61, 0x74, 0xe9,
	0x71, 0xdf, 0x9c, 0x3d, 0xf2, 0xc0, 0x56, 0x48, 0x43, 0x2a, 0x9e, 0x66, 0xf6, 0x92, 0xde, 0x9d,
	0x6b, 0x15, 0x34, 0x57, 0x54, 0xe3, 0x64, 0xc5, 0xe0, 0x5d, 0xd0, 0x52, 0xee, 0x65, 0xa5, 0xf9,
	0x41, 0x74, 0x55, 0x45, 0x0d, 0xd4, 0xdc, 0x6c, 0x6f, 0x7e, 0xfb, 0xf5, 0x5d, 0x2d, 0x32, 0xa5,
	0x81, 0x08, 0x64, 0x51, 0x47, 0x04, 0xf1, 0x3b, 0xd8, 0xe4, 0xcc, 0x8b, 0xd3, 0x84, 0x32, 0x5e,
	0x55, 0x1a, 0xa8, 0x59, 0x69, 0xbd, 0x30, 0xee, 0x3c, 0x29, 0xc3, 0x9d, 0x31, 0xdc, 0x49, 0xe2,
	0xb7, 0x21, 0xd3, 0x5a, 0xfb, 0x84, 0x14, 0x1d, 0x91, 0x5b, 0x3c, 0x0e, 0x61, 0x43, 0x14, 0xdc,
	0xa7, 0xc3, 0xaa, 0x2a, 0xa4, 0x9e, 0xaf, 0x20, 0xe5, 0xe4, 0x88, 0x25, 0xa5, 0x1b, 0x38, 0x1e,
	0x40, 0x49, 0xec, 0xa1, 0x27, 0xb9, 0xd5, 0x62, 0x03, 0x35, 0xb5, 0x56, 0x77, 0x05, 0x31, 0x92,
	0x61, 0x3a, 0x22, 0x61, 0xc4, 0x3c, 0x1e, 0xd1, 0x98, 0x68, 0xec, 0xd6, 0x87, 0x03, 0xa8, 0xe4,
	0x89, 0x12, 0x97, 0x56, 0xd7, 0x1a, 0x6a, 0x53, 0x5b, 0xa9, 0x31, 0xb9, 0xc2, 0x03, 0xf1, 0x95,
	0x94, 0xf9, 0x9c, 0x95, 0xee, 0x7c, 0x45, 0x50, 0x9a, 0x8f, 0xe3, 0xc7, 0x50, 0x8c, 0xbd, 0x0f,
	0xfe, 0xf2, 0x72, 0x85, 0x1b, 0xef, 0xc1, 0x7a, 0xde, 0xbb, 0x22, 0x7a, 0xdf, 0x36, 0xe4, 0xc1,
	0x19, 0xb3, 0x83, 0x33, 0xce, 0xc4, 0xc1, 0x1d, 0x16, 0x48, 0xfe, 0x11, 0x3f, 0x85, 0x12, 0x9f,
	0x24, 0xfe, 0xe5, 0x6c, 0x68, 0xaa, 0x48, 0xdc, 0x5a, 0x4a, 0xb4, 0xe2, 0xc9, 0x61, 0x81, 0x68,
	0xe2, 0xaf, 0x9c, 0x42, 0xbb, 0x0c, 0x9a, 0x4c, 0xea, 0x65, 0xde, 0x9d, 0xdf, 0x08, 0x1e, 0xdc,
	0xdc, 0xa3, 0xd8, 0x88, 0x9d, 0x64, 0xa3, 0x4b, 0x17, 0xaf, 0x0d, 0xfd, 0xbf, 0x6b, 0x53, 0xfe,
	0xe1, 0xb5, 0xed, 0xda, 0x50, 0x5e, 0x28, 0x08, 0x63, 0xa8, 0x58, 0xe7, 0xae, 0xdd, 0x73, 0x89,
	0x75, 0x7a, 0xe6, 0xd8, 0xc4, 0xd5, 0x0b, 0x18, 0x60, 0xfd, 0x80, 0x58, 0x27, 0xdd, 0x7d, 0x1d,
	0xe1, 0x12, 0x6c, 0x9c, 0x9f, 0xe6, 0x96, 0x92, 0x45, 0x0e, 0xbb, 0xd6, 0x7e, 0x97, 0xe8, 0x6a,
	0xad, 0xf8, 0xf9, 0x4b, 0xbd, 0xb0, 0xfb, 0x16, 0x4a, 0xf3, 0xb2, 0xf8, 0x3e, 0x94, 0x05, 0xcf,
	0x21, 0xb6, 0x6b, 0x77, 0xec, 0x63, 0x89, 0x6b, 0x1f, 0x9d, 0x5a, 0xe4, 0xb5, 0x8e, 0x70, 0x05,
	0xe0, 0xd8, 0x7a, 0xd5, 0xcb, 0x6d, 0x05, 0x6b, 0x70, 0xaf, 0x63, 0x9f, 0x38, 0x56, 0xc7, 0xd5,
	0xd5, 0xcc, 0x70, 0x5f, 0x1e, 0xb9, 0x6e, 0x97, 0xe8, 0x45, 0x89, 0x6f, 0xeb, 0x3f, 0xa6, 0x75,
	0xf4, 0x73, 0x5a, 0x47, 0xd7, 0xd3, 0x3a, 0x7a, 0xa3, 0x8c, 0x5b, 0x17, 0xeb, 0xa2, 0x97, 0x27,
	0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x3b, 0xfc, 0xb1, 0xbc, 0x69, 0x05, 0x00, 0x00,
}
