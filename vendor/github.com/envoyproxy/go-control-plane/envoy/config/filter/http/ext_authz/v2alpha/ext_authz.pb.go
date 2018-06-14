// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/http/ext_authz/v2alpha/ext_authz.proto

/*
	Package v2alpha is a generated protocol buffer package.

	It is generated from these files:
		envoy/config/filter/http/ext_authz/v2alpha/ext_authz.proto

	It has these top-level messages:
		HttpService
		ExtAuthz
*/
package v2alpha

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import envoy_api_v2_core1 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
import envoy_api_v2_core2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

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
// [#comment: The HttpService is under development and will be supported soon.]
type HttpService struct {
	// Sets the HTTP server URI which the authorization requests must be sent to.
	ServerUri *envoy_api_v2_core2.HttpUri `protobuf:"bytes,1,opt,name=server_uri,json=serverUri" json:"server_uri,omitempty"`
	// Sets an optional prefix to the value of authorization request header `path`.
	PathPrefix string `protobuf:"bytes,2,opt,name=path_prefix,json=pathPrefix,proto3" json:"path_prefix,omitempty"`
}

func (m *HttpService) Reset()                    { *m = HttpService{} }
func (m *HttpService) String() string            { return proto.CompactTextString(m) }
func (*HttpService) ProtoMessage()               {}
func (*HttpService) Descriptor() ([]byte, []int) { return fileDescriptorExtAuthz, []int{0} }

func (m *HttpService) GetServerUri() *envoy_api_v2_core2.HttpUri {
	if m != nil {
		return m.ServerUri
	}
	return nil
}

func (m *HttpService) GetPathPrefix() string {
	if m != nil {
		return m.PathPrefix
	}
	return ""
}

// External Authorization filter calls out to an external service over the
// gRPC Authorization API defined by
// :ref:`CheckRequest <envoy_api_msg_service.auth.v2alpha.CheckRequest>`.
// A failed check will cause this filter to close the HTTP request with 403(Forbidden).
type ExtAuthz struct {
	// Types that are valid to be assigned to Services:
	//	*ExtAuthz_GrpcService
	//	*ExtAuthz_HttpService
	Services isExtAuthz_Services `protobuf_oneof:"services"`
	// The filter's behaviour in case the external authorization service does
	// not respond back. When it is set to true, Envoy will also allow traffic in case of
	// communication failure between authorization service and the proxy.
	// Defaults to false.
	FailureModeAllow bool `protobuf:"varint,2,opt,name=failure_mode_allow,json=failureModeAllow,proto3" json:"failure_mode_allow,omitempty"`
}

func (m *ExtAuthz) Reset()                    { *m = ExtAuthz{} }
func (m *ExtAuthz) String() string            { return proto.CompactTextString(m) }
func (*ExtAuthz) ProtoMessage()               {}
func (*ExtAuthz) Descriptor() ([]byte, []int) { return fileDescriptorExtAuthz, []int{1} }

type isExtAuthz_Services interface {
	isExtAuthz_Services()
	MarshalTo([]byte) (int, error)
	Size() int
}

type ExtAuthz_GrpcService struct {
	GrpcService *envoy_api_v2_core1.GrpcService `protobuf:"bytes,1,opt,name=grpc_service,json=grpcService,oneof"`
}
type ExtAuthz_HttpService struct {
	HttpService *HttpService `protobuf:"bytes,3,opt,name=http_service,json=httpService,oneof"`
}

func (*ExtAuthz_GrpcService) isExtAuthz_Services() {}
func (*ExtAuthz_HttpService) isExtAuthz_Services() {}

func (m *ExtAuthz) GetServices() isExtAuthz_Services {
	if m != nil {
		return m.Services
	}
	return nil
}

func (m *ExtAuthz) GetGrpcService() *envoy_api_v2_core1.GrpcService {
	if x, ok := m.GetServices().(*ExtAuthz_GrpcService); ok {
		return x.GrpcService
	}
	return nil
}

func (m *ExtAuthz) GetHttpService() *HttpService {
	if x, ok := m.GetServices().(*ExtAuthz_HttpService); ok {
		return x.HttpService
	}
	return nil
}

func (m *ExtAuthz) GetFailureModeAllow() bool {
	if m != nil {
		return m.FailureModeAllow
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ExtAuthz) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ExtAuthz_OneofMarshaler, _ExtAuthz_OneofUnmarshaler, _ExtAuthz_OneofSizer, []interface{}{
		(*ExtAuthz_GrpcService)(nil),
		(*ExtAuthz_HttpService)(nil),
	}
}

func _ExtAuthz_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ExtAuthz)
	// services
	switch x := m.Services.(type) {
	case *ExtAuthz_GrpcService:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GrpcService); err != nil {
			return err
		}
	case *ExtAuthz_HttpService:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.HttpService); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ExtAuthz.Services has unexpected type %T", x)
	}
	return nil
}

func _ExtAuthz_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ExtAuthz)
	switch tag {
	case 1: // services.grpc_service
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(envoy_api_v2_core1.GrpcService)
		err := b.DecodeMessage(msg)
		m.Services = &ExtAuthz_GrpcService{msg}
		return true, err
	case 3: // services.http_service
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(HttpService)
		err := b.DecodeMessage(msg)
		m.Services = &ExtAuthz_HttpService{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ExtAuthz_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ExtAuthz)
	// services
	switch x := m.Services.(type) {
	case *ExtAuthz_GrpcService:
		s := proto.Size(x.GrpcService)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ExtAuthz_HttpService:
		s := proto.Size(x.HttpService)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*HttpService)(nil), "envoy.config.filter.http.ext_authz.v2alpha.HttpService")
	proto.RegisterType((*ExtAuthz)(nil), "envoy.config.filter.http.ext_authz.v2alpha.ExtAuthz")
}
func (m *HttpService) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HttpService) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.ServerUri != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintExtAuthz(dAtA, i, uint64(m.ServerUri.Size()))
		n1, err := m.ServerUri.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.PathPrefix) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintExtAuthz(dAtA, i, uint64(len(m.PathPrefix)))
		i += copy(dAtA[i:], m.PathPrefix)
	}
	return i, nil
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
	if m.Services != nil {
		nn2, err := m.Services.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn2
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

func (m *ExtAuthz_GrpcService) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.GrpcService != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintExtAuthz(dAtA, i, uint64(m.GrpcService.Size()))
		n3, err := m.GrpcService.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}
func (m *ExtAuthz_HttpService) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.HttpService != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintExtAuthz(dAtA, i, uint64(m.HttpService.Size()))
		n4, err := m.HttpService.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
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
func (m *HttpService) Size() (n int) {
	var l int
	_ = l
	if m.ServerUri != nil {
		l = m.ServerUri.Size()
		n += 1 + l + sovExtAuthz(uint64(l))
	}
	l = len(m.PathPrefix)
	if l > 0 {
		n += 1 + l + sovExtAuthz(uint64(l))
	}
	return n
}

func (m *ExtAuthz) Size() (n int) {
	var l int
	_ = l
	if m.Services != nil {
		n += m.Services.Size()
	}
	if m.FailureModeAllow {
		n += 2
	}
	return n
}

func (m *ExtAuthz_GrpcService) Size() (n int) {
	var l int
	_ = l
	if m.GrpcService != nil {
		l = m.GrpcService.Size()
		n += 1 + l + sovExtAuthz(uint64(l))
	}
	return n
}
func (m *ExtAuthz_HttpService) Size() (n int) {
	var l int
	_ = l
	if m.HttpService != nil {
		l = m.HttpService.Size()
		n += 1 + l + sovExtAuthz(uint64(l))
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
func (m *HttpService) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: HttpService: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HttpService: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServerUri", wireType)
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
			if m.ServerUri == nil {
				m.ServerUri = &envoy_api_v2_core2.HttpUri{}
			}
			if err := m.ServerUri.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PathPrefix", wireType)
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
			m.PathPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
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
			v := &envoy_api_v2_core1.GrpcService{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Services = &ExtAuthz_GrpcService{v}
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HttpService", wireType)
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
			v := &HttpService{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Services = &ExtAuthz_HttpService{v}
			iNdEx = postIndex
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
	proto.RegisterFile("envoy/config/filter/http/ext_authz/v2alpha/ext_authz.proto", fileDescriptorExtAuthz)
}

var fileDescriptorExtAuthz = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xc1, 0x4a, 0x33, 0x31,
	0x14, 0x85, 0xff, 0xfc, 0x82, 0xb6, 0x77, 0xba, 0x90, 0x80, 0x50, 0xba, 0x18, 0x4b, 0x71, 0x51,
	0x44, 0x12, 0x18, 0x17, 0xa2, 0xbb, 0x56, 0xc4, 0x6e, 0x04, 0x19, 0xe9, 0x46, 0x84, 0x21, 0x4e,
	0x33, 0x9d, 0xc0, 0xd8, 0x84, 0xdb, 0x74, 0xac, 0x3e, 0xa1, 0x4b, 0x1f, 0x41, 0xba, 0xf1, 0x35,
	0x24, 0x93, 0x91, 0x29, 0xe8, 0xc2, 0x65, 0x6e, 0xce, 0x77, 0xce, 0xb9, 0x17, 0x2e, 0xe4, 0xa2,
	0xd4, 0x2f, 0x3c, 0xd5, 0x8b, 0x4c, 0xcd, 0x79, 0xa6, 0x0a, 0x2b, 0x91, 0xe7, 0xd6, 0x1a, 0x2e,
	0xd7, 0x36, 0x11, 0x2b, 0x9b, 0xbf, 0xf2, 0x32, 0x12, 0x85, 0xc9, 0x45, 0x33, 0x61, 0x06, 0xb5,
	0xd5, 0xf4, 0xb8, 0x62, 0x99, 0x67, 0x99, 0x67, 0x99, 0x63, 0x59, 0xa3, 0xac, 0xd9, 0xde, 0x91,
	0xcf, 0x11, 0x46, 0xf1, 0x32, 0xe2, 0xa9, 0x46, 0xc9, 0xe7, 0x68, 0xd2, 0x64, 0x29, 0xb1, 0x54,
	0xa9, 0xf4, 0x8e, 0xbd, 0xfe, 0x4f, 0x95, 0xf3, 0x4b, 0x56, 0xa8, 0xbc, 0x62, 0xa0, 0x20, 0x98,
	0x58, 0x6b, 0xee, 0x3c, 0x46, 0xcf, 0x01, 0x9c, 0x83, 0x44, 0x27, 0xe9, 0x92, 0x3e, 0x19, 0x06,
	0x51, 0x8f, 0xf9, 0x5e, 0xc2, 0x28, 0x56, 0x46, 0xcc, 0xb9, 0x30, 0xc7, 0x4c, 0x51, 0xc5, 0x6d,
	0xaf, 0x9e, 0xa2, 0xa2, 0x87, 0x10, 0x18, 0x61, 0xf3, 0xc4, 0xa0, 0xcc, 0xd4, 0xba, 0xfb, 0xbf,
	0x4f, 0x86, 0xed, 0x18, 0xdc, 0xe8, 0xb6, 0x9a, 0x0c, 0x3e, 0x09, 0xb4, 0xae, 0xd6, 0x76, 0xe4,
	0xf6, 0xa0, 0x97, 0xd0, 0xd9, 0xee, 0x5b, 0x47, 0x85, 0xbf, 0x44, 0x5d, 0xa3, 0x49, 0xeb, 0x7a,
	0x93, 0x7f, 0x71, 0x30, 0x6f, 0x9e, 0xf4, 0x04, 0x68, 0x26, 0x54, 0xb1, 0x42, 0x99, 0x3c, 0xe9,
	0x99, 0x4c, 0x44, 0x51, 0xe8, 0xe7, 0x2a, 0xb9, 0x15, 0xef, 0xd7, 0x3f, 0x37, 0x7a, 0x26, 0x47,
	0x6e, 0x4e, 0x1f, 0xa0, 0x53, 0x2d, 0xff, 0x1d, 0xb9, 0x53, 0x45, 0x9e, 0xb1, 0xbf, 0x5f, 0x9d,
	0x6d, 0x9d, 0xca, 0x75, 0xc9, 0x9b, 0xe7, 0x18, 0xa0, 0x55, 0x1b, 0x2f, 0xc7, 0x07, 0x6f, 0x9b,
	0x90, 0xbc, 0x6f, 0x42, 0xf2, 0xb1, 0x09, 0xc9, 0xfd, 0x5e, 0x4d, 0x3f, 0xee, 0x56, 0x27, 0x3f,
	0xfd, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x97, 0x82, 0x1b, 0x29, 0x24, 0x02, 0x00, 0x00,
}
