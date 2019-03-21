// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/transport_socket/tap/v2alpha/tap.proto

package v2

import (
	fmt "fmt"
	io "io"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	_ "github.com/lyft/protoc-gen-validate/validate"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	v2alpha "github.com/envoyproxy/go-control-plane/envoy/config/common/tap/v2alpha"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Configuration for tap transport socket. This wraps another transport socket, providing the
// ability to interpose and record in plain text any traffic that is surfaced to Envoy.
type Tap struct {
	// Common configuration for the tap transport socket.
	CommonConfig *v2alpha.CommonExtensionConfig `protobuf:"bytes,1,opt,name=common_config,json=commonConfig,proto3" json:"common_config,omitempty"`
	// The underlying transport socket being wrapped.
	TransportSocket      *core.TransportSocket `protobuf:"bytes,2,opt,name=transport_socket,json=transportSocket,proto3" json:"transport_socket,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *Tap) Reset()         { *m = Tap{} }
func (m *Tap) String() string { return proto.CompactTextString(m) }
func (*Tap) ProtoMessage()    {}
func (*Tap) Descriptor() ([]byte, []int) {
	return fileDescriptor_07cb8c0b42756e40, []int{0}
}
func (m *Tap) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Tap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Tap.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Tap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tap.Merge(m, src)
}
func (m *Tap) XXX_Size() int {
	return m.Size()
}
func (m *Tap) XXX_DiscardUnknown() {
	xxx_messageInfo_Tap.DiscardUnknown(m)
}

var xxx_messageInfo_Tap proto.InternalMessageInfo

func (m *Tap) GetCommonConfig() *v2alpha.CommonExtensionConfig {
	if m != nil {
		return m.CommonConfig
	}
	return nil
}

func (m *Tap) GetTransportSocket() *core.TransportSocket {
	if m != nil {
		return m.TransportSocket
	}
	return nil
}

func init() {
	proto.RegisterType((*Tap)(nil), "envoy.config.transport_socket.tap.v2alpha.Tap")
}

func init() {
	proto.RegisterFile("envoy/config/transport_socket/tap/v2alpha/tap.proto", fileDescriptor_07cb8c0b42756e40)
}

var fileDescriptor_07cb8c0b42756e40 = []byte{
	// 290 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x50, 0xbd, 0x4a, 0x03, 0x41,
	0x10, 0x66, 0x13, 0x14, 0x59, 0x15, 0xc3, 0x35, 0x86, 0x20, 0x87, 0xa4, 0x52, 0x90, 0x5d, 0xb8,
	0x80, 0xf6, 0x09, 0xf6, 0x41, 0xcf, 0x42, 0x9b, 0x30, 0x39, 0xd7, 0xb8, 0x98, 0xec, 0x2c, 0x77,
	0xcb, 0x92, 0xbc, 0x82, 0x8f, 0x64, 0x65, 0xa9, 0x9d, 0x8f, 0x20, 0xd7, 0xf9, 0x16, 0x72, 0x3b,
	0x17, 0xf0, 0xae, 0xb2, 0x1b, 0xe6, 0xfb, 0x9b, 0xf9, 0xf8, 0x48, 0x19, 0x8f, 0x1b, 0x99, 0xa1,
	0x79, 0xd2, 0x0b, 0xe9, 0x72, 0x30, 0x85, 0xc5, 0xdc, 0xcd, 0x0a, 0xcc, 0x5e, 0x94, 0x93, 0x0e,
	0xac, 0xf4, 0x09, 0x2c, 0xed, 0x33, 0x54, 0xb3, 0xb0, 0x39, 0x3a, 0x8c, 0xce, 0x83, 0x48, 0x90,
	0x48, 0xb4, 0x45, 0xa2, 0x22, 0xd6, 0xa2, 0xc1, 0x45, 0xc3, 0x3f, 0xc3, 0xd5, 0x0a, 0x4d, 0xc3,
	0x95, 0x56, 0x64, 0x3c, 0x38, 0x21, 0x36, 0x58, 0x2d, 0x7d, 0x22, 0x33, 0xcc, 0x95, 0x9c, 0x43,
	0xa1, 0x6a, 0xf4, 0xd8, 0xc3, 0x52, 0x3f, 0x82, 0x53, 0x72, 0x3b, 0x10, 0x30, 0xfc, 0x64, 0xbc,
	0x9b, 0x82, 0x8d, 0x16, 0xfc, 0x90, 0xec, 0x66, 0x94, 0xd7, 0x67, 0xa7, 0xec, 0x6c, 0x3f, 0xb9,
	0x14, 0x8d, 0x7b, 0xeb, 0xc4, 0x3f, 0x57, 0x8a, 0x49, 0x58, 0x5d, 0xaf, 0x9d, 0x32, 0x85, 0x46,
	0x33, 0x09, 0xc4, 0x31, 0x7f, 0xfb, 0x79, 0xef, 0xee, 0xbc, 0xb2, 0x4e, 0x8f, 0xdd, 0x1c, 0x90,
	0x8a, 0x90, 0xe8, 0x9e, 0xf7, 0xda, 0x5f, 0xf7, 0x3b, 0x21, 0x6b, 0x58, 0x67, 0x81, 0xd5, 0xc2,
	0x27, 0xa2, 0x7a, 0x41, 0xa4, 0x5b, 0xea, 0x6d, 0x60, 0x36, 0x7c, 0x8f, 0x5c, 0x0b, 0xbc, 0xfb,
	0x28, 0x63, 0xf6, 0x55, 0xc6, 0xec, 0xbb, 0x8c, 0x19, 0xbf, 0xd2, 0x48, 0x86, 0x36, 0xc7, 0xf5,
	0x46, 0xfc, 0xbb, 0xf7, 0xf1, 0x5e, 0x0a, 0x76, 0x5a, 0x95, 0x33, 0x65, 0x0f, 0x1d, 0x9f, 0xcc,
	0x77, 0x43, 0x53, 0xa3, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0x76, 0x4a, 0xf2, 0xf0, 0x01,
	0x00, 0x00,
}

func (m *Tap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Tap) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.CommonConfig != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintTap(dAtA, i, uint64(m.CommonConfig.Size()))
		n1, err := m.CommonConfig.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.TransportSocket != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintTap(dAtA, i, uint64(m.TransportSocket.Size()))
		n2, err := m.TransportSocket.MarshalTo(dAtA[i:])
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

func encodeVarintTap(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Tap) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CommonConfig != nil {
		l = m.CommonConfig.Size()
		n += 1 + l + sovTap(uint64(l))
	}
	if m.TransportSocket != nil {
		l = m.TransportSocket.Size()
		n += 1 + l + sovTap(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTap(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozTap(x uint64) (n int) {
	return sovTap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Tap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTap
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Tap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Tap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommonConfig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CommonConfig == nil {
				m.CommonConfig = &v2alpha.CommonExtensionConfig{}
			}
			if err := m.CommonConfig.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransportSocket", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TransportSocket == nil {
				m.TransportSocket = &core.TransportSocket{}
			}
			if err := m.TransportSocket.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTap
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
func skipTap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTap
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
					return 0, ErrIntOverflowTap
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
					return 0, ErrIntOverflowTap
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
			if length < 0 {
				return 0, ErrInvalidLengthTap
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthTap
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTap
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
				next, err := skipTap(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthTap
				}
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
	ErrInvalidLengthTap = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTap   = fmt.Errorf("proto: integer overflow")
)
