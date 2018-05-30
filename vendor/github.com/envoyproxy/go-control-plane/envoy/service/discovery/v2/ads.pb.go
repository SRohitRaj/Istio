// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/service/discovery/v2/ads.proto

/*
	Package v2 is a generated protocol buffer package.

	It is generated from these files:
		envoy/service/discovery/v2/ads.proto
		envoy/service/discovery/v2/hds.proto
		envoy/service/discovery/v2/sds.proto

	It has these top-level messages:
		AdsDummy
		Capability
		HealthCheckRequest
		EndpointHealth
		EndpointHealthResponse
		HealthCheckRequestOrEndpointHealthResponse
		LocalityEndpoints
		ClusterHealthCheck
		HealthCheckSpecifier
		SdsDummy
*/
package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"

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

// [#not-implemented-hide:] Not configuration. Workaround c++ protobuf issue with importing
// services: https://github.com/google/protobuf/issues/4221
type AdsDummy struct {
}

func (m *AdsDummy) Reset()                    { *m = AdsDummy{} }
func (m *AdsDummy) String() string            { return proto.CompactTextString(m) }
func (*AdsDummy) ProtoMessage()               {}
func (*AdsDummy) Descriptor() ([]byte, []int) { return fileDescriptorAds, []int{0} }

func init() {
	proto.RegisterType((*AdsDummy)(nil), "envoy.service.discovery.v2.AdsDummy")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AggregatedDiscoveryService service

type AggregatedDiscoveryServiceClient interface {
	// This is a gRPC-only API.
	StreamAggregatedResources(ctx context.Context, opts ...grpc.CallOption) (AggregatedDiscoveryService_StreamAggregatedResourcesClient, error)
}

type aggregatedDiscoveryServiceClient struct {
	cc *grpc.ClientConn
}

func NewAggregatedDiscoveryServiceClient(cc *grpc.ClientConn) AggregatedDiscoveryServiceClient {
	return &aggregatedDiscoveryServiceClient{cc}
}

func (c *aggregatedDiscoveryServiceClient) StreamAggregatedResources(ctx context.Context, opts ...grpc.CallOption) (AggregatedDiscoveryService_StreamAggregatedResourcesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_AggregatedDiscoveryService_serviceDesc.Streams[0], c.cc, "/envoy.service.discovery.v2.AggregatedDiscoveryService/StreamAggregatedResources", opts...)
	if err != nil {
		return nil, err
	}
	x := &aggregatedDiscoveryServiceStreamAggregatedResourcesClient{stream}
	return x, nil
}

type AggregatedDiscoveryService_StreamAggregatedResourcesClient interface {
	Send(*envoy_api_v2.DiscoveryRequest) error
	Recv() (*envoy_api_v2.DiscoveryResponse, error)
	grpc.ClientStream
}

type aggregatedDiscoveryServiceStreamAggregatedResourcesClient struct {
	grpc.ClientStream
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesClient) Send(m *envoy_api_v2.DiscoveryRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesClient) Recv() (*envoy_api_v2.DiscoveryResponse, error) {
	m := new(envoy_api_v2.DiscoveryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for AggregatedDiscoveryService service

type AggregatedDiscoveryServiceServer interface {
	// This is a gRPC-only API.
	StreamAggregatedResources(AggregatedDiscoveryService_StreamAggregatedResourcesServer) error
}

func RegisterAggregatedDiscoveryServiceServer(s *grpc.Server, srv AggregatedDiscoveryServiceServer) {
	s.RegisterService(&_AggregatedDiscoveryService_serviceDesc, srv)
}

func _AggregatedDiscoveryService_StreamAggregatedResources_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AggregatedDiscoveryServiceServer).StreamAggregatedResources(&aggregatedDiscoveryServiceStreamAggregatedResourcesServer{stream})
}

type AggregatedDiscoveryService_StreamAggregatedResourcesServer interface {
	Send(*envoy_api_v2.DiscoveryResponse) error
	Recv() (*envoy_api_v2.DiscoveryRequest, error)
	grpc.ServerStream
}

type aggregatedDiscoveryServiceStreamAggregatedResourcesServer struct {
	grpc.ServerStream
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesServer) Send(m *envoy_api_v2.DiscoveryResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *aggregatedDiscoveryServiceStreamAggregatedResourcesServer) Recv() (*envoy_api_v2.DiscoveryRequest, error) {
	m := new(envoy_api_v2.DiscoveryRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AggregatedDiscoveryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "envoy.service.discovery.v2.AggregatedDiscoveryService",
	HandlerType: (*AggregatedDiscoveryServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamAggregatedResources",
			Handler:       _AggregatedDiscoveryService_StreamAggregatedResources_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "envoy/service/discovery/v2/ads.proto",
}

func (m *AdsDummy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AdsDummy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func encodeVarintAds(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *AdsDummy) Size() (n int) {
	var l int
	_ = l
	return n
}

func sovAds(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozAds(x uint64) (n int) {
	return sovAds(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AdsDummy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAds
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
			return fmt.Errorf("proto: AdsDummy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AdsDummy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipAds(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAds
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
func skipAds(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAds
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
					return 0, ErrIntOverflowAds
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
					return 0, ErrIntOverflowAds
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
				return 0, ErrInvalidLengthAds
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowAds
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
				next, err := skipAds(dAtA[start:])
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
	ErrInvalidLengthAds = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAds   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("envoy/service/discovery/v2/ads.proto", fileDescriptorAds) }

var fileDescriptorAds = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0xce, 0xb1, 0x4e, 0x86, 0x30,
	0x14, 0x05, 0x60, 0xaf, 0x83, 0x31, 0x1d, 0x71, 0xb2, 0x31, 0x35, 0xf9, 0xe3, 0xe0, 0xd4, 0x9a,
	0xfa, 0x04, 0x18, 0x9e, 0x00, 0x36, 0xb7, 0x42, 0x6f, 0x48, 0x07, 0x68, 0xed, 0x2d, 0x4d, 0xd8,
	0x1c, 0x7d, 0x34, 0x47, 0x1f, 0xc1, 0xf0, 0x24, 0x46, 0x50, 0x98, 0xfe, 0xf9, 0x7c, 0xe7, 0xe4,
	0xb0, 0x07, 0x1c, 0xb3, 0x9f, 0x15, 0x61, 0xcc, 0xae, 0x43, 0x65, 0x1d, 0x75, 0x3e, 0x63, 0x9c,
	0x55, 0xd6, 0xca, 0x58, 0x92, 0x21, 0xfa, 0xe4, 0x0b, 0xbe, 0x2a, 0xf9, 0xa7, 0xe4, 0xae, 0x64,
	0xd6, 0xfc, 0x6e, 0x5b, 0x30, 0xc1, 0xfd, 0x76, 0x8e, 0x68, 0x6d, 0x9e, 0x18, 0xbb, 0x2e, 0x2d,
	0x55, 0xd3, 0x30, 0xcc, 0xfa, 0x1d, 0x18, 0x2f, 0xfb, 0x3e, 0x62, 0x6f, 0x12, 0xda, 0xea, 0x5f,
	0x36, 0xdb, 0x6a, 0xd1, 0xb2, 0xdb, 0x26, 0x45, 0x34, 0xc3, 0x61, 0x6a, 0x24, 0x3f, 0xc5, 0x0e,
	0xa9, 0x10, 0x72, 0xbb, 0x60, 0x82, 0x93, 0x59, 0xcb, 0xbd, 0x5c, 0xe3, 0xdb, 0x84, 0x94, 0xf8,
	0xfd, 0xd9, 0x9c, 0x82, 0x1f, 0x09, 0x4f, 0x17, 0x8f, 0xf0, 0x04, 0x2f, 0x37, 0x9f, 0x8b, 0x80,
	0xaf, 0x45, 0xc0, 0xf7, 0x22, 0xe0, 0xf5, 0x32, 0xeb, 0x0f, 0x80, 0xf6, 0x6a, 0xbd, 0xfa, 0xfc,
	0x13, 0x00, 0x00, 0xff, 0xff, 0x91, 0x7b, 0xed, 0xc5, 0x0c, 0x01, 0x00, 0x00,
}
