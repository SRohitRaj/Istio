// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/network/ext_authz/v2/ext_authz.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
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

// External Authorization filter calls out to an external service over the
// gRPC Authorization API defined by
// :ref:`CheckRequest <envoy_api_msg_service.auth.v2alpha.CheckRequest>`.
// A failed check will cause this filter to close the TCP connection.
type ExtAuthz struct {
	// The prefix to use when emitting statistics.
	StatPrefix string `protobuf:"bytes,1,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	// The external authorization gRPC service configuration.
	// The default timeout is set to 200ms by this filter.
	GrpcService *core.GrpcService `protobuf:"bytes,2,opt,name=grpc_service,json=grpcService,proto3" json:"grpc_service,omitempty"`
	// The filter's behaviour in case the external authorization service does
	// not respond back. When it is set to true, Envoy will also allow traffic in case of
	// communication failure between authorization service and the proxy.
	// Defaults to false.
	FailureModeAllow     bool     `protobuf:"varint,3,opt,name=failure_mode_allow,json=failureModeAllow,proto3" json:"failure_mode_allow,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExtAuthz) Reset()         { *m = ExtAuthz{} }
func (m *ExtAuthz) String() string { return proto.CompactTextString(m) }
func (*ExtAuthz) ProtoMessage()    {}
func (*ExtAuthz) Descriptor() ([]byte, []int) {
	return fileDescriptor_ext_authz_9b91fbc3b30d1404, []int{0}
}
func (m *ExtAuthz) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExtAuthz) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExtAuthz.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *ExtAuthz) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtAuthz.Merge(dst, src)
}
func (m *ExtAuthz) XXX_Size() int {
	return m.Size()
}
func (m *ExtAuthz) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtAuthz.DiscardUnknown(m)
}

var xxx_messageInfo_ExtAuthz proto.InternalMessageInfo

func (m *ExtAuthz) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *ExtAuthz) GetGrpcService() *core.GrpcService {
	if m != nil {
		return m.GrpcService
	}
	return nil
}

func (m *ExtAuthz) GetFailureModeAllow() bool {
	if m != nil {
		return m.FailureModeAllow
	}
	return false
}

func init() {
	proto.RegisterType((*ExtAuthz)(nil), "envoy.config.filter.network.ext_authz.v2.ExtAuthz")
}
func (m *ExtAuthz) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExtAuthz) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.StatPrefix) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintExtAuthz(dAtA, i, uint64(len(m.StatPrefix)))
		i += copy(dAtA[i:], m.StatPrefix)
	}
	if m.GrpcService != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExtAuthz(dAtA, i, uint64(m.GrpcService.Size()))
		n1, err := m.GrpcService.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.FailureModeAllow {
		dAtA[i] = 0x18
		i++
		if m.FailureModeAllow {
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

func encodeVarintExtAuthz(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ExtAuthz) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.StatPrefix)
	if l > 0 {
		n += 1 + l + sovExtAuthz(uint64(l))
	}
	if m.GrpcService != nil {
		l = m.GrpcService.Size()
		n += 1 + l + sovExtAuthz(uint64(l))
	}
	if m.FailureModeAllow {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovExtAuthz(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozExtAuthz(x uint64) (n int) {
	return sovExtAuthz(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ExtAuthz) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowExtAuthz
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
			return fmt.Errorf("proto: ExtAuthz: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExtAuthz: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtAuthz
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
				return ErrInvalidLengthExtAuthz
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GrpcService", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtAuthz
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
				return ErrInvalidLengthExtAuthz
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GrpcService == nil {
				m.GrpcService = &core.GrpcService{}
			}
			if err := m.GrpcService.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FailureModeAllow", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowExtAuthz
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
			m.FailureModeAllow = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipExtAuthz(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthExtAuthz
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
func skipExtAuthz(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowExtAuthz
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
					return 0, ErrIntOverflowExtAuthz
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
					return 0, ErrIntOverflowExtAuthz
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
				return 0, ErrInvalidLengthExtAuthz
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowExtAuthz
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
				next, err := skipExtAuthz(dAtA[start:])
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
	ErrInvalidLengthExtAuthz = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowExtAuthz   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/filter/network/ext_authz/v2/ext_authz.proto", fileDescriptor_ext_authz_9b91fbc3b30d1404)
}

var fileDescriptor_ext_authz_9b91fbc3b30d1404 = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xb1, 0x4a, 0xc4, 0x30,
	0x1c, 0xc6, 0x49, 0x4f, 0xe4, 0x2e, 0x75, 0x90, 0x2e, 0x96, 0x1b, 0x4a, 0x11, 0x87, 0x22, 0x92,
	0x40, 0x05, 0x71, 0xed, 0x81, 0x08, 0x82, 0x70, 0xd4, 0xcd, 0xa5, 0xc4, 0xf6, 0xdf, 0x1a, 0xac,
	0x4d, 0xc8, 0xe5, 0x72, 0x3d, 0x5f, 0xc8, 0x77, 0x70, 0x72, 0x74, 0xf4, 0x11, 0xa4, 0x9b, 0x6f,
	0x21, 0x6d, 0x4e, 0xea, 0xe8, 0x96, 0x7f, 0xbe, 0xff, 0xef, 0xfb, 0xbe, 0x04, 0x5f, 0x42, 0x63,
	0xc4, 0x96, 0xe6, 0xa2, 0x29, 0x79, 0x45, 0x4b, 0x5e, 0x6b, 0x50, 0xb4, 0x01, 0xbd, 0x11, 0xea,
	0x89, 0x42, 0xab, 0x33, 0xb6, 0xd6, 0x8f, 0x2f, 0xd4, 0xc4, 0xe3, 0x40, 0xa4, 0x12, 0x5a, 0x78,
	0xd1, 0x40, 0x12, 0x4b, 0x12, 0x4b, 0x92, 0x1d, 0x49, 0xc6, 0x65, 0x13, 0xcf, 0x4f, 0x6c, 0x06,
	0x93, 0xbc, 0xf7, 0xc9, 0x85, 0x02, 0x5a, 0x29, 0x99, 0x67, 0x2b, 0x50, 0x86, 0xe7, 0x60, 0xfd,
	0xe6, 0x47, 0x86, 0xd5, 0xbc, 0x60, 0x1a, 0xe8, 0xef, 0xc1, 0x0a, 0xc7, 0xaf, 0x08, 0x4f, 0xaf,
	0x5a, 0x9d, 0xf4, 0x76, 0xde, 0x29, 0x76, 0x57, 0x9a, 0xe9, 0x4c, 0x2a, 0x28, 0x79, 0xeb, 0xa3,
	0x10, 0x45, 0xb3, 0xc5, 0xec, 0xed, 0xfb, 0x7d, 0xb2, 0xa7, 0x9c, 0x10, 0xa5, 0xb8, 0x57, 0x97,
	0x83, 0xe8, 0x25, 0xf8, 0xe0, 0x6f, 0x8e, 0xef, 0x84, 0x28, 0x72, 0xe3, 0x80, 0xd8, 0xe2, 0x4c,
	0x72, 0x62, 0x62, 0xd2, 0xd7, 0x21, 0xd7, 0x4a, 0xe6, 0x77, 0x76, 0x2b, 0x75, 0xab, 0x71, 0xf0,
	0xce, 0xb0, 0x57, 0x32, 0x5e, 0xaf, 0x15, 0x64, 0xcf, 0xa2, 0x80, 0x8c, 0xd5, 0xb5, 0xd8, 0xf8,
	0x93, 0x10, 0x45, 0xd3, 0xf4, 0x70, 0xa7, 0xdc, 0x8a, 0x02, 0x92, 0xfe, 0x7e, 0x71, 0xf3, 0xd1,
	0x05, 0xe8, 0xb3, 0x0b, 0xd0, 0x57, 0x17, 0x20, 0x7c, 0xc1, 0x85, 0x8d, 0x92, 0x4a, 0xb4, 0x5b,
	0xf2, 0xdf, 0xef, 0x5a, 0xa2, 0x7b, 0xc7, 0xc4, 0x0f, 0xfb, 0xc3, 0xe3, 0xcf, 0x7f, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x82, 0x00, 0xd2, 0x51, 0xa1, 0x01, 0x00, 0x00,
}
