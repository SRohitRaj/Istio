// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/network/dubbo_proxy/v2alpha1/dubbo_proxy.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
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

// Configure the protocol used.
type DubboProxy_ProtocolType int32

const (
	DubboProxy_Dubbo DubboProxy_ProtocolType = 0
)

var DubboProxy_ProtocolType_name = map[int32]string{
	0: "Dubbo",
}
var DubboProxy_ProtocolType_value = map[string]int32{
	"Dubbo": 0,
}

func (x DubboProxy_ProtocolType) String() string {
	return proto.EnumName(DubboProxy_ProtocolType_name, int32(x))
}
func (DubboProxy_ProtocolType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dubbo_proxy_90d444f7a8aea9cc, []int{0, 0}
}

// Configure the serialization protocol used.
type DubboProxy_SerializationType int32

const (
	DubboProxy_Hessian2 DubboProxy_SerializationType = 0
)

var DubboProxy_SerializationType_name = map[int32]string{
	0: "Hessian2",
}
var DubboProxy_SerializationType_value = map[string]int32{
	"Hessian2": 0,
}

func (x DubboProxy_SerializationType) String() string {
	return proto.EnumName(DubboProxy_SerializationType_name, int32(x))
}
func (DubboProxy_SerializationType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_dubbo_proxy_90d444f7a8aea9cc, []int{0, 1}
}

type DubboProxy struct {
	// The human readable prefix to use when emitting statistics.
	StatPrefix           string                       `protobuf:"bytes,1,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	ProtocolType         DubboProxy_ProtocolType      `protobuf:"varint,2,opt,name=protocol_type,json=protocolType,proto3,enum=envoy.extensions.filters.network.dubbo_proxy.v2alpha1.DubboProxy_ProtocolType" json:"protocol_type,omitempty"`
	SerializationType    DubboProxy_SerializationType `protobuf:"varint,3,opt,name=serialization_type,json=serializationType,proto3,enum=envoy.extensions.filters.network.dubbo_proxy.v2alpha1.DubboProxy_SerializationType" json:"serialization_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *DubboProxy) Reset()         { *m = DubboProxy{} }
func (m *DubboProxy) String() string { return proto.CompactTextString(m) }
func (*DubboProxy) ProtoMessage()    {}
func (*DubboProxy) Descriptor() ([]byte, []int) {
	return fileDescriptor_dubbo_proxy_90d444f7a8aea9cc, []int{0}
}
func (m *DubboProxy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DubboProxy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DubboProxy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *DubboProxy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DubboProxy.Merge(dst, src)
}
func (m *DubboProxy) XXX_Size() int {
	return m.Size()
}
func (m *DubboProxy) XXX_DiscardUnknown() {
	xxx_messageInfo_DubboProxy.DiscardUnknown(m)
}

var xxx_messageInfo_DubboProxy proto.InternalMessageInfo

func (m *DubboProxy) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *DubboProxy) GetProtocolType() DubboProxy_ProtocolType {
	if m != nil {
		return m.ProtocolType
	}
	return DubboProxy_Dubbo
}

func (m *DubboProxy) GetSerializationType() DubboProxy_SerializationType {
	if m != nil {
		return m.SerializationType
	}
	return DubboProxy_Hessian2
}

func init() {
	proto.RegisterType((*DubboProxy)(nil), "envoy.extensions.filters.network.dubbo_proxy.v2alpha1.DubboProxy")
	proto.RegisterEnum("envoy.extensions.filters.network.dubbo_proxy.v2alpha1.DubboProxy_ProtocolType", DubboProxy_ProtocolType_name, DubboProxy_ProtocolType_value)
	proto.RegisterEnum("envoy.extensions.filters.network.dubbo_proxy.v2alpha1.DubboProxy_SerializationType", DubboProxy_SerializationType_name, DubboProxy_SerializationType_value)
}
func (m *DubboProxy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DubboProxy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.StatPrefix) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintDubboProxy(dAtA, i, uint64(len(m.StatPrefix)))
		i += copy(dAtA[i:], m.StatPrefix)
	}
	if m.ProtocolType != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintDubboProxy(dAtA, i, uint64(m.ProtocolType))
	}
	if m.SerializationType != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintDubboProxy(dAtA, i, uint64(m.SerializationType))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintDubboProxy(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *DubboProxy) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.StatPrefix)
	if l > 0 {
		n += 1 + l + sovDubboProxy(uint64(l))
	}
	if m.ProtocolType != 0 {
		n += 1 + sovDubboProxy(uint64(m.ProtocolType))
	}
	if m.SerializationType != 0 {
		n += 1 + sovDubboProxy(uint64(m.SerializationType))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDubboProxy(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozDubboProxy(x uint64) (n int) {
	return sovDubboProxy(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DubboProxy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDubboProxy
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
			return fmt.Errorf("proto: DubboProxy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DubboProxy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDubboProxy
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
				return ErrInvalidLengthDubboProxy
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProtocolType", wireType)
			}
			m.ProtocolType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDubboProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ProtocolType |= (DubboProxy_ProtocolType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SerializationType", wireType)
			}
			m.SerializationType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDubboProxy
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SerializationType |= (DubboProxy_SerializationType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDubboProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDubboProxy
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
func skipDubboProxy(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDubboProxy
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
					return 0, ErrIntOverflowDubboProxy
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
					return 0, ErrIntOverflowDubboProxy
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
				return 0, ErrInvalidLengthDubboProxy
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowDubboProxy
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
				next, err := skipDubboProxy(dAtA[start:])
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
	ErrInvalidLengthDubboProxy = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDubboProxy   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/filter/network/dubbo_proxy/v2alpha1/dubbo_proxy.proto", fileDescriptor_dubbo_proxy_90d444f7a8aea9cc)
}

var fileDescriptor_dubbo_proxy_90d444f7a8aea9cc = []byte{
	// 316 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x8f, 0xb1, 0x4e, 0x02, 0x31,
	0x1c, 0x87, 0xe9, 0x21, 0x46, 0xfe, 0xa2, 0x39, 0xba, 0x88, 0x0e, 0x04, 0x99, 0x88, 0x43, 0x1b,
	0xcf, 0xf8, 0x02, 0x17, 0x07, 0x27, 0x43, 0xc0, 0xc9, 0x85, 0x14, 0x28, 0xda, 0x78, 0x69, 0x9b,
	0xb6, 0x9e, 0x9c, 0x8b, 0x89, 0x93, 0x93, 0x0f, 0xe4, 0xe4, 0xe8, 0xe8, 0x23, 0x18, 0x36, 0xdf,
	0xc2, 0x5c, 0x0f, 0xe2, 0x45, 0x36, 0xdd, 0x9a, 0x7e, 0xf9, 0x7d, 0xfd, 0x0a, 0x31, 0x97, 0xa9,
	0xca, 0xe8, 0x44, 0xc9, 0x99, 0xb8, 0xa6, 0x33, 0x91, 0x38, 0x6e, 0xa8, 0xe4, 0xee, 0x5e, 0x99,
	0x5b, 0x3a, 0xbd, 0x1b, 0x8f, 0xd5, 0x48, 0x1b, 0x35, 0xcf, 0x68, 0x1a, 0xb1, 0x44, 0xdf, 0xb0,
	0xe3, 0xf2, 0x25, 0xd1, 0x46, 0x39, 0x85, 0x4f, 0xbd, 0x83, 0xf0, 0xb9, 0xe3, 0xd2, 0x0a, 0x25,
	0x2d, 0x29, 0x3c, 0x96, 0x2c, 0x45, 0xa4, 0xbc, 0x59, 0x89, 0x0e, 0xf6, 0x52, 0x96, 0x88, 0x29,
	0x73, 0x9c, 0xae, 0x0e, 0x85, 0xaf, 0xfb, 0x5c, 0x05, 0x38, 0xcb, 0x17, 0xfd, 0x7c, 0x80, 0x8f,
	0x60, 0xdb, 0x3a, 0xe6, 0x46, 0xda, 0xf0, 0x99, 0x98, 0xb7, 0x50, 0x07, 0xf5, 0xea, 0x71, 0xfd,
	0xf5, 0xeb, 0xad, 0xba, 0x61, 0x82, 0x0e, 0x1a, 0x40, 0x4e, 0xfb, 0x1e, 0xe2, 0x47, 0xd8, 0xf1,
	0x8e, 0x89, 0x4a, 0x46, 0x2e, 0xd3, 0xbc, 0x15, 0x74, 0x50, 0x6f, 0x37, 0xba, 0x20, 0x7f, 0x4a,
	0x24, 0x3f, 0x15, 0xa4, 0xbf, 0xd4, 0x5e, 0x66, 0x9a, 0xc7, 0x90, 0xbf, 0x5e, 0x7b, 0x42, 0x41,
	0x88, 0x06, 0x0d, 0x5d, 0x22, 0xf8, 0x05, 0x01, 0xb6, 0xdc, 0x08, 0x96, 0x88, 0x07, 0xe6, 0x84,
	0x92, 0x45, 0x46, 0xd5, 0x67, 0x0c, 0xff, 0x9f, 0x31, 0x2c, 0xbb, 0xd7, 0x5a, 0x9a, 0xf6, 0x37,
	0xee, 0xee, 0x43, 0xa3, 0x9c, 0x8e, 0xeb, 0x50, 0xf3, 0xba, 0xb0, 0xd2, 0x3d, 0x84, 0xe6, 0x9a,
	0x0e, 0x37, 0x60, 0xeb, 0x9c, 0x5b, 0x2b, 0x98, 0x8c, 0xc2, 0x4a, 0x1c, 0xbe, 0x2f, 0xda, 0xe8,
	0x63, 0xd1, 0x46, 0x9f, 0x8b, 0x36, 0xba, 0x0a, 0xd2, 0x68, 0xbc, 0xe9, 0xbf, 0x7b, 0xf2, 0x1d,
	0x00, 0x00, 0xff, 0xff, 0xcf, 0x2c, 0x8c, 0x17, 0x39, 0x02, 0x00, 0x00,
}
