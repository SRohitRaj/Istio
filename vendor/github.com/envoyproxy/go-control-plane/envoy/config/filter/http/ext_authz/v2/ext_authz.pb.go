// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/http/ext_authz/v2/ext_authz.proto

/*
	Package v2 is a generated protocol buffer package.

	It is generated from these files:
		envoy/config/filter/http/ext_authz/v2/ext_authz.proto

	It has these top-level messages:
		ExtAuthz
*/
package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import envoy_api_v2_core1 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
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

// [#not-implemented-hide:]
// External Authorization filter calls out to an external service over the
// gRPC Authorization API defined by
// :ref:`external_auth <envoy_api_msg_auth.CheckRequest>`.
// A failed check will cause this filter to return 403 Forbidden.
type ExtAuthz struct {
	// The external authorization gRPC service configuration.
	GrpcService *envoy_api_v2_core1.GrpcService `protobuf:"bytes,1,opt,name=grpc_service,json=grpcService" json:"grpc_service,omitempty"`
	// The filter's behaviour in case the external authorization service does
	// not respond back. If set to true then in case of failure to get a
	// response back from the authorization service or getting a response that
	// is NOT denied then traffic will be permitted.
	// Defaults to false.
	FailureModeAllow bool `protobuf:"varint,2,opt,name=failure_mode_allow,json=failureModeAllow,proto3" json:"failure_mode_allow,omitempty"`
}

func (m *ExtAuthz) Reset()                    { *m = ExtAuthz{} }
func (m *ExtAuthz) String() string            { return proto.CompactTextString(m) }
func (*ExtAuthz) ProtoMessage()               {}
func (*ExtAuthz) Descriptor() ([]byte, []int) { return fileDescriptorExtAuthz, []int{0} }

func (m *ExtAuthz) GetGrpcService() *envoy_api_v2_core1.GrpcService {
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
	proto.RegisterType((*ExtAuthz)(nil), "envoy.config.filter.http.ext_authz.v2.ExtAuthz")
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
	if m.GrpcService != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintExtAuthz(dAtA, i, uint64(m.GrpcService.Size()))
		n1, err := m.GrpcService.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.FailureModeAllow {
		dAtA[i] = 0x10
		i++
		if m.FailureModeAllow {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
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
	var l int
	_ = l
	if m.GrpcService != nil {
		l = m.GrpcService.Size()
		n += 1 + l + sovExtAuthz(uint64(l))
	}
	if m.FailureModeAllow {
		n += 2
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
				m.GrpcService = &envoy_api_v2_core1.GrpcService{}
			}
			if err := m.GrpcService.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
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
	proto.RegisterFile("envoy/config/filter/http/ext_authz/v2/ext_authz.proto", fileDescriptorExtAuthz)
}

var fileDescriptorExtAuthz = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8e, 0xc1, 0x4a, 0x03, 0x31,
	0x10, 0x86, 0x49, 0x0f, 0x52, 0xb6, 0x1e, 0xca, 0x5e, 0x2c, 0x3d, 0x2c, 0x45, 0x14, 0x7a, 0x90,
	0x09, 0x44, 0x7c, 0x80, 0x15, 0xc4, 0x93, 0x97, 0x7a, 0xf3, 0x12, 0x62, 0x76, 0x76, 0x1b, 0x88,
	0x4d, 0x88, 0xd3, 0x58, 0xf5, 0x05, 0x3d, 0xfa, 0x08, 0xb2, 0x4f, 0x22, 0xd9, 0x28, 0xdb, 0x5b,
	0x26, 0xf3, 0x7f, 0xdf, 0xfc, 0xc5, 0x0d, 0xee, 0xa2, 0x7b, 0xe7, 0xda, 0xed, 0x5a, 0xd3, 0xf1,
	0xd6, 0x58, 0xc2, 0xc0, 0xb7, 0x44, 0x9e, 0xe3, 0x81, 0xa4, 0xda, 0xd3, 0xf6, 0x83, 0x47, 0x31,
	0x0e, 0xe0, 0x83, 0x23, 0x57, 0x5e, 0x0e, 0x18, 0x64, 0x0c, 0x32, 0x06, 0x09, 0x83, 0x31, 0x19,
	0xc5, 0xf2, 0x22, 0xdb, 0x95, 0x37, 0x49, 0xa2, 0x5d, 0x40, 0xde, 0x05, 0xaf, 0xe5, 0x2b, 0x86,
	0x68, 0x34, 0x66, 0xd9, 0xf2, 0x2c, 0x2a, 0x6b, 0x1a, 0x45, 0xc8, 0xff, 0x1f, 0x79, 0x71, 0xfe,
	0x59, 0x4c, 0xef, 0x0e, 0x54, 0x27, 0x5b, 0x59, 0x17, 0xa7, 0xc7, 0xe8, 0x82, 0xad, 0xd8, 0x7a,
	0x26, 0x2a, 0xc8, 0x45, 0x94, 0x37, 0x10, 0x05, 0xa4, 0x0b, 0x70, 0x1f, 0xbc, 0x7e, 0xcc, 0xa9,
	0xcd, 0xac, 0x1b, 0x87, 0xf2, 0xaa, 0x28, 0x5b, 0x65, 0xec, 0x3e, 0xa0, 0x7c, 0x71, 0x0d, 0x4a,
	0x65, 0xad, 0x7b, 0x5b, 0x4c, 0x56, 0x6c, 0x3d, 0xdd, 0xcc, 0xff, 0x36, 0x0f, 0xae, 0xc1, 0x3a,
	0xfd, 0xdf, 0xce, 0xbf, 0xfa, 0x8a, 0x7d, 0xf7, 0x15, 0xfb, 0xe9, 0x2b, 0xf6, 0x34, 0x89, 0xe2,
	0xf9, 0x64, 0x68, 0x75, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0xd6, 0x43, 0xa3, 0x68, 0x34, 0x01,
	0x00, 0x00,
}
