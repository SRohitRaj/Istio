// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/ratelimit/v2/rls.proto

/*
	Package v2 is a generated protocol buffer package.

	It is generated from these files:
		envoy/config/ratelimit/v2/rls.proto

	It has these top-level messages:
		RateLimitServiceConfig
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

// Rate limit :ref:`configuration overview <config_rate_limit_service>`.
type RateLimitServiceConfig struct {
	// Types that are valid to be assigned to ServiceSpecifier:
	//	*RateLimitServiceConfig_ClusterName
	//	*RateLimitServiceConfig_GrpcService
	ServiceSpecifier isRateLimitServiceConfig_ServiceSpecifier `protobuf_oneof:"service_specifier"`
}

func (m *RateLimitServiceConfig) Reset()                    { *m = RateLimitServiceConfig{} }
func (m *RateLimitServiceConfig) String() string            { return proto.CompactTextString(m) }
func (*RateLimitServiceConfig) ProtoMessage()               {}
func (*RateLimitServiceConfig) Descriptor() ([]byte, []int) { return fileDescriptorRls, []int{0} }

type isRateLimitServiceConfig_ServiceSpecifier interface {
	isRateLimitServiceConfig_ServiceSpecifier()
	MarshalTo([]byte) (int, error)
	Size() int
}

type RateLimitServiceConfig_ClusterName struct {
	ClusterName string `protobuf:"bytes,1,opt,name=cluster_name,json=clusterName,proto3,oneof"`
}
type RateLimitServiceConfig_GrpcService struct {
	GrpcService *envoy_api_v2_core1.GrpcService `protobuf:"bytes,2,opt,name=grpc_service,json=grpcService,oneof"`
}

func (*RateLimitServiceConfig_ClusterName) isRateLimitServiceConfig_ServiceSpecifier() {}
func (*RateLimitServiceConfig_GrpcService) isRateLimitServiceConfig_ServiceSpecifier() {}

func (m *RateLimitServiceConfig) GetServiceSpecifier() isRateLimitServiceConfig_ServiceSpecifier {
	if m != nil {
		return m.ServiceSpecifier
	}
	return nil
}

func (m *RateLimitServiceConfig) GetClusterName() string {
	if x, ok := m.GetServiceSpecifier().(*RateLimitServiceConfig_ClusterName); ok {
		return x.ClusterName
	}
	return ""
}

func (m *RateLimitServiceConfig) GetGrpcService() *envoy_api_v2_core1.GrpcService {
	if x, ok := m.GetServiceSpecifier().(*RateLimitServiceConfig_GrpcService); ok {
		return x.GrpcService
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*RateLimitServiceConfig) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _RateLimitServiceConfig_OneofMarshaler, _RateLimitServiceConfig_OneofUnmarshaler, _RateLimitServiceConfig_OneofSizer, []interface{}{
		(*RateLimitServiceConfig_ClusterName)(nil),
		(*RateLimitServiceConfig_GrpcService)(nil),
	}
}

func _RateLimitServiceConfig_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*RateLimitServiceConfig)
	// service_specifier
	switch x := m.ServiceSpecifier.(type) {
	case *RateLimitServiceConfig_ClusterName:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.ClusterName)
	case *RateLimitServiceConfig_GrpcService:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GrpcService); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("RateLimitServiceConfig.ServiceSpecifier has unexpected type %T", x)
	}
	return nil
}

func _RateLimitServiceConfig_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*RateLimitServiceConfig)
	switch tag {
	case 1: // service_specifier.cluster_name
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.ServiceSpecifier = &RateLimitServiceConfig_ClusterName{x}
		return true, err
	case 2: // service_specifier.grpc_service
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(envoy_api_v2_core1.GrpcService)
		err := b.DecodeMessage(msg)
		m.ServiceSpecifier = &RateLimitServiceConfig_GrpcService{msg}
		return true, err
	default:
		return false, nil
	}
}

func _RateLimitServiceConfig_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*RateLimitServiceConfig)
	// service_specifier
	switch x := m.ServiceSpecifier.(type) {
	case *RateLimitServiceConfig_ClusterName:
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.ClusterName)))
		n += len(x.ClusterName)
	case *RateLimitServiceConfig_GrpcService:
		s := proto.Size(x.GrpcService)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*RateLimitServiceConfig)(nil), "envoy.config.ratelimit.v2.RateLimitServiceConfig")
}
func (m *RateLimitServiceConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RateLimitServiceConfig) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ServiceSpecifier != nil {
		nn1, err := m.ServiceSpecifier.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	return i, nil
}

func (m *RateLimitServiceConfig_ClusterName) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0xa
	i++
	i = encodeVarintRls(dAtA, i, uint64(len(m.ClusterName)))
	i += copy(dAtA[i:], m.ClusterName)
	return i, nil
}
func (m *RateLimitServiceConfig_GrpcService) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.GrpcService != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintRls(dAtA, i, uint64(m.GrpcService.Size()))
		n2, err := m.GrpcService.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}
func encodeVarintRls(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *RateLimitServiceConfig) Size() (n int) {
	var l int
	_ = l
	if m.ServiceSpecifier != nil {
		n += m.ServiceSpecifier.Size()
	}
	return n
}

func (m *RateLimitServiceConfig_ClusterName) Size() (n int) {
	var l int
	_ = l
	l = len(m.ClusterName)
	n += 1 + l + sovRls(uint64(l))
	return n
}
func (m *RateLimitServiceConfig_GrpcService) Size() (n int) {
	var l int
	_ = l
	if m.GrpcService != nil {
		l = m.GrpcService.Size()
		n += 1 + l + sovRls(uint64(l))
	}
	return n
}

func sovRls(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozRls(x uint64) (n int) {
	return sovRls(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RateLimitServiceConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRls
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
			return fmt.Errorf("proto: RateLimitServiceConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RateLimitServiceConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRls
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
				return ErrInvalidLengthRls
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ServiceSpecifier = &RateLimitServiceConfig_ClusterName{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GrpcService", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRls
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
				return ErrInvalidLengthRls
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &envoy_api_v2_core1.GrpcService{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.ServiceSpecifier = &RateLimitServiceConfig_GrpcService{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRls(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRls
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
func skipRls(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRls
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
					return 0, ErrIntOverflowRls
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
					return 0, ErrIntOverflowRls
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
				return 0, ErrInvalidLengthRls
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowRls
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
				next, err := skipRls(dAtA[start:])
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
	ErrInvalidLengthRls = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRls   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("envoy/config/ratelimit/v2/rls.proto", fileDescriptorRls) }

var fileDescriptorRls = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x3b, 0x4e, 0xc4, 0x30,
	0x14, 0x45, 0x79, 0x19, 0x40, 0x22, 0x99, 0x02, 0x52, 0x40, 0x48, 0x11, 0x45, 0x40, 0x31, 0x95,
	0x8d, 0xc2, 0x0e, 0x32, 0x05, 0x14, 0x88, 0x62, 0xe8, 0x68, 0x22, 0x63, 0xde, 0x44, 0x96, 0xf2,
	0xb1, 0x5e, 0x8c, 0x25, 0xd6, 0xc4, 0x06, 0x10, 0xd5, 0x94, 0x94, 0x2c, 0x01, 0xa5, 0x9b, 0x5d,
	0xa0, 0xc4, 0xe1, 0xd3, 0xd9, 0xf2, 0x39, 0xf7, 0x5e, 0xfb, 0xe7, 0xd8, 0xd8, 0xf6, 0x85, 0xcb,
	0xb6, 0x59, 0xab, 0x92, 0x93, 0x30, 0x58, 0xa9, 0x5a, 0x19, 0x6e, 0x33, 0x4e, 0x55, 0xc7, 0x34,
	0xb5, 0xa6, 0x0d, 0x4f, 0x47, 0x88, 0x39, 0x88, 0xfd, 0x42, 0xcc, 0x66, 0xf1, 0x85, 0xf3, 0x85,
	0x56, 0x83, 0x22, 0x5b, 0x42, 0x5e, 0x92, 0x96, 0x45, 0x87, 0x64, 0x95, 0x44, 0x17, 0x10, 0x9f,
	0x58, 0x51, 0xa9, 0x27, 0x61, 0x90, 0xff, 0x1c, 0xdc, 0xc3, 0xd9, 0x2b, 0xf8, 0xc7, 0x2b, 0x61,
	0xf0, 0x76, 0xc8, 0xbb, 0x77, 0xce, 0x72, 0xac, 0x09, 0x2f, 0xfd, 0xb9, 0xac, 0x9e, 0x3b, 0x83,
	0x54, 0x34, 0xa2, 0xc6, 0x08, 0x52, 0x58, 0x1c, 0xe4, 0xc1, 0xfb, 0x76, 0x33, 0xdb, 0x25, 0x2f,
	0x85, 0x08, 0x6e, 0x76, 0x56, 0xc1, 0x84, 0xdc, 0x89, 0x1a, 0xc3, 0xa5, 0x3f, 0xff, 0xdf, 0x1d,
	0x79, 0x29, 0x2c, 0x82, 0x2c, 0x61, 0x6e, 0xbd, 0xd0, 0x8a, 0xd9, 0x8c, 0x0d, 0x13, 0xd9, 0x35,
	0x69, 0x39, 0xb5, 0x0d, 0x21, 0xe5, 0xdf, 0x35, 0x8f, 0xfd, 0xa3, 0xc9, 0x2f, 0x3a, 0x8d, 0x52,
	0xad, 0x15, 0x52, 0xb8, 0xf7, 0xb6, 0xdd, 0xcc, 0x20, 0x3f, 0xfc, 0xe8, 0x13, 0xf8, 0xec, 0x13,
	0xf8, 0xea, 0x13, 0x78, 0xf0, 0x6c, 0xf6, 0xb8, 0x3f, 0x7e, 0xe3, 0xea, 0x3b, 0x00, 0x00, 0xff,
	0xff, 0xe3, 0xbb, 0x91, 0xf3, 0x47, 0x01, 0x00, 0x00,
}
