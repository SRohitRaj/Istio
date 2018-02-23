// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: security/proto/ca_service.proto

/*
	Package istio_v1_auth is a generated protocol buffer package.

	It is generated from these files:
		security/proto/ca_service.proto

	It has these top-level messages:
		CsrRequest
		CsrResponse
*/
package istio_v1_auth

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_rpc "istio.io/gogo-genproto/googleapis/google/rpc"
import _ "github.com/gogo/protobuf/gogoproto"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

import strings "strings"
import reflect "reflect"

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

type CsrRequest struct {
	// PEM-encoded certificate signing request
	CsrPem []byte `protobuf:"bytes,1,opt,name=csr_pem,json=csrPem,proto3" json:"csr_pem,omitempty"`
	// opaque credential for node agent
	NodeAgentCredential []byte `protobuf:"bytes,2,opt,name=node_agent_credential,json=nodeAgentCredential,proto3" json:"node_agent_credential,omitempty"`
	// type of the node_agent_credential (aws/gcp/onprem/custom...)
	CredentialType string `protobuf:"bytes,3,opt,name=credential_type,json=credentialType,proto3" json:"credential_type,omitempty"`
	// the requested ttl of the certificate in minutes
	RequestedTtlMinutes int32 `protobuf:"varint,4,opt,name=requested_ttl_minutes,json=requestedTtlMinutes,proto3" json:"requested_ttl_minutes,omitempty"`
	// whether the certificate is for a CA
	ForCA bool `protobuf:"varint,5,opt,name=forCA,proto3" json:"forCA,omitempty"`
}

func (m *CsrRequest) Reset()                    { *m = CsrRequest{} }
func (*CsrRequest) ProtoMessage()               {}
func (*CsrRequest) Descriptor() ([]byte, []int) { return fileDescriptorCaService, []int{0} }

type CsrResponse struct {
	// Whether the CSR is approved.
	IsApproved bool               `protobuf:"varint,1,opt,name=is_approved,json=isApproved,proto3" json:"is_approved,omitempty"`
	Status     *google_rpc.Status `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	// The signed target cert.
	SignedCert []byte `protobuf:"bytes,3,opt,name=signed_cert,json=signedCert,proto3" json:"signed_cert,omitempty"`
	// The cert chain up to the trusted root cert. It includes all the certs between the
	// newly signed cert and the root cert.
	CertChain []byte `protobuf:"bytes,4,opt,name=cert_chain,json=certChain,proto3" json:"cert_chain,omitempty"`
}

func (m *CsrResponse) Reset()                    { *m = CsrResponse{} }
func (*CsrResponse) ProtoMessage()               {}
func (*CsrResponse) Descriptor() ([]byte, []int) { return fileDescriptorCaService, []int{1} }

func init() {
	proto.RegisterType((*CsrRequest)(nil), "istio.v1.auth.CsrRequest")
	proto.RegisterType((*CsrResponse)(nil), "istio.v1.auth.CsrResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for IstioCAService service

type IstioCAServiceClient interface {
	// A request object includes a PEM-encoded certificate signing request that
	// is generated on the Node Agent. Additionaly credential can be attached
	// within the request object for a server to authenticate the originating
	// node agent.
	HandleCSR(ctx context.Context, in *CsrRequest, opts ...grpc.CallOption) (*CsrResponse, error)
}

type istioCAServiceClient struct {
	cc *grpc.ClientConn
}

func NewIstioCAServiceClient(cc *grpc.ClientConn) IstioCAServiceClient {
	return &istioCAServiceClient{cc}
}

func (c *istioCAServiceClient) HandleCSR(ctx context.Context, in *CsrRequest, opts ...grpc.CallOption) (*CsrResponse, error) {
	out := new(CsrResponse)
	err := grpc.Invoke(ctx, "/istio.v1.auth.IstioCAService/HandleCSR", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for IstioCAService service

type IstioCAServiceServer interface {
	// A request object includes a PEM-encoded certificate signing request that
	// is generated on the Node Agent. Additionaly credential can be attached
	// within the request object for a server to authenticate the originating
	// node agent.
	HandleCSR(context.Context, *CsrRequest) (*CsrResponse, error)
}

func RegisterIstioCAServiceServer(s *grpc.Server, srv IstioCAServiceServer) {
	s.RegisterService(&_IstioCAService_serviceDesc, srv)
}

func _IstioCAService_HandleCSR_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CsrRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IstioCAServiceServer).HandleCSR(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/istio.v1.auth.IstioCAService/HandleCSR",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IstioCAServiceServer).HandleCSR(ctx, req.(*CsrRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IstioCAService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "istio.v1.auth.IstioCAService",
	HandlerType: (*IstioCAServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleCSR",
			Handler:    _IstioCAService_HandleCSR_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "security/proto/ca_service.proto",
}

func (m *CsrRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CsrRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.CsrPem) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintCaService(dAtA, i, uint64(len(m.CsrPem)))
		i += copy(dAtA[i:], m.CsrPem)
	}
	if len(m.NodeAgentCredential) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCaService(dAtA, i, uint64(len(m.NodeAgentCredential)))
		i += copy(dAtA[i:], m.NodeAgentCredential)
	}
	if len(m.CredentialType) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCaService(dAtA, i, uint64(len(m.CredentialType)))
		i += copy(dAtA[i:], m.CredentialType)
	}
	if m.RequestedTtlMinutes != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintCaService(dAtA, i, uint64(m.RequestedTtlMinutes))
	}
	if m.ForCA {
		dAtA[i] = 0x28
		i++
		if m.ForCA {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	return i, nil
}

func (m *CsrResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CsrResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.IsApproved {
		dAtA[i] = 0x8
		i++
		if m.IsApproved {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.Status != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintCaService(dAtA, i, uint64(m.Status.Size()))
		n1, err := m.Status.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.SignedCert) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintCaService(dAtA, i, uint64(len(m.SignedCert)))
		i += copy(dAtA[i:], m.SignedCert)
	}
	if len(m.CertChain) > 0 {
		dAtA[i] = 0x22
		i++
		i = encodeVarintCaService(dAtA, i, uint64(len(m.CertChain)))
		i += copy(dAtA[i:], m.CertChain)
	}
	return i, nil
}

func encodeVarintCaService(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *CsrRequest) Size() (n int) {
	var l int
	_ = l
	l = len(m.CsrPem)
	if l > 0 {
		n += 1 + l + sovCaService(uint64(l))
	}
	l = len(m.NodeAgentCredential)
	if l > 0 {
		n += 1 + l + sovCaService(uint64(l))
	}
	l = len(m.CredentialType)
	if l > 0 {
		n += 1 + l + sovCaService(uint64(l))
	}
	if m.RequestedTtlMinutes != 0 {
		n += 1 + sovCaService(uint64(m.RequestedTtlMinutes))
	}
	if m.ForCA {
		n += 2
	}
	return n
}

func (m *CsrResponse) Size() (n int) {
	var l int
	_ = l
	if m.IsApproved {
		n += 2
	}
	if m.Status != nil {
		l = m.Status.Size()
		n += 1 + l + sovCaService(uint64(l))
	}
	l = len(m.SignedCert)
	if l > 0 {
		n += 1 + l + sovCaService(uint64(l))
	}
	l = len(m.CertChain)
	if l > 0 {
		n += 1 + l + sovCaService(uint64(l))
	}
	return n
}

func sovCaService(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozCaService(x uint64) (n int) {
	return sovCaService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *CsrRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&CsrRequest{`,
		`CsrPem:` + fmt.Sprintf("%v", this.CsrPem) + `,`,
		`NodeAgentCredential:` + fmt.Sprintf("%v", this.NodeAgentCredential) + `,`,
		`CredentialType:` + fmt.Sprintf("%v", this.CredentialType) + `,`,
		`RequestedTtlMinutes:` + fmt.Sprintf("%v", this.RequestedTtlMinutes) + `,`,
		`ForCA:` + fmt.Sprintf("%v", this.ForCA) + `,`,
		`}`,
	}, "")
	return s
}
func (this *CsrResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&CsrResponse{`,
		`IsApproved:` + fmt.Sprintf("%v", this.IsApproved) + `,`,
		`Status:` + strings.Replace(fmt.Sprintf("%v", this.Status), "Status", "google_rpc.Status", 1) + `,`,
		`SignedCert:` + fmt.Sprintf("%v", this.SignedCert) + `,`,
		`CertChain:` + fmt.Sprintf("%v", this.CertChain) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringCaService(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *CsrRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCaService
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
			return fmt.Errorf("proto: CsrRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CsrRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CsrPem", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
				return ErrInvalidLengthCaService
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CsrPem = append(m.CsrPem[:0], dAtA[iNdEx:postIndex]...)
			if m.CsrPem == nil {
				m.CsrPem = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeAgentCredential", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
				return ErrInvalidLengthCaService
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeAgentCredential = append(m.NodeAgentCredential[:0], dAtA[iNdEx:postIndex]...)
			if m.NodeAgentCredential == nil {
				m.NodeAgentCredential = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CredentialType", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
				return ErrInvalidLengthCaService
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CredentialType = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequestedTtlMinutes", wireType)
			}
			m.RequestedTtlMinutes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RequestedTtlMinutes |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ForCA", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
			m.ForCA = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipCaService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCaService
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
func (m *CsrResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCaService
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
			return fmt.Errorf("proto: CsrResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CsrResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsApproved", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
			m.IsApproved = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
				return ErrInvalidLengthCaService
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Status == nil {
				m.Status = &google_rpc.Status{}
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignedCert", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
				return ErrInvalidLengthCaService
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignedCert = append(m.SignedCert[:0], dAtA[iNdEx:postIndex]...)
			if m.SignedCert == nil {
				m.SignedCert = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CertChain", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCaService
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
				return ErrInvalidLengthCaService
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CertChain = append(m.CertChain[:0], dAtA[iNdEx:postIndex]...)
			if m.CertChain == nil {
				m.CertChain = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCaService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCaService
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
func skipCaService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCaService
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
					return 0, ErrIntOverflowCaService
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
					return 0, ErrIntOverflowCaService
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
				return 0, ErrInvalidLengthCaService
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCaService
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
				next, err := skipCaService(dAtA[start:])
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
	ErrInvalidLengthCaService = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCaService   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("security/proto/ca_service.proto", fileDescriptorCaService) }

var fileDescriptorCaService = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xb1, 0x6e, 0x13, 0x31,
	0x1c, 0xc6, 0xcf, 0x40, 0x42, 0xe3, 0x84, 0x22, 0xb9, 0x45, 0x3d, 0x22, 0xe1, 0x46, 0x5d, 0x88,
	0x18, 0x2e, 0x6a, 0x78, 0x00, 0x94, 0xde, 0x02, 0x03, 0x12, 0x72, 0x32, 0xb1, 0x58, 0x87, 0xef,
	0xcf, 0xd5, 0xd2, 0xe5, 0x6c, 0x6c, 0x5f, 0xa4, 0x6c, 0x3c, 0x02, 0x33, 0x4f, 0xc0, 0xa3, 0x74,
	0x60, 0xe8, 0xc8, 0x48, 0x8e, 0x85, 0xb1, 0x8f, 0x80, 0x6c, 0x57, 0x8d, 0x90, 0xba, 0xdd, 0x7d,
	0xbf, 0x4f, 0x9f, 0xfc, 0xf9, 0x33, 0x3e, 0xb5, 0x20, 0x5a, 0x23, 0xdd, 0x76, 0xa6, 0x8d, 0x72,
	0x6a, 0x26, 0x0a, 0x6e, 0xc1, 0x6c, 0xa4, 0x80, 0x2c, 0x08, 0xe4, 0x89, 0xb4, 0x4e, 0xaa, 0x6c,
	0x73, 0x9e, 0x15, 0xad, 0xbb, 0x1c, 0x9f, 0x54, 0x4a, 0x55, 0x35, 0xcc, 0x8c, 0x16, 0x33, 0xeb,
	0x0a, 0xd7, 0xda, 0xe8, 0x1b, 0x1f, 0x57, 0xaa, 0x52, 0x31, 0xc3, 0x7f, 0x45, 0xf5, 0xec, 0x27,
	0xc2, 0x38, 0xb7, 0x86, 0xc1, 0x97, 0x16, 0xac, 0x23, 0x27, 0xf8, 0xb1, 0xb0, 0x86, 0x6b, 0x58,
	0xa7, 0x68, 0x82, 0xa6, 0x23, 0xd6, 0x17, 0xd6, 0x7c, 0x80, 0x35, 0x99, 0xe3, 0x67, 0x8d, 0x2a,
	0x81, 0x17, 0x15, 0x34, 0x8e, 0x0b, 0x03, 0x25, 0x34, 0x4e, 0x16, 0x75, 0xfa, 0x20, 0xd8, 0x8e,
	0x3c, 0x5c, 0x78, 0x96, 0xdf, 0x21, 0xf2, 0x12, 0x3f, 0xdd, 0x1b, 0xb9, 0xdb, 0x6a, 0x48, 0x1f,
	0x4e, 0xd0, 0x74, 0xc0, 0x0e, 0xf7, 0xf2, 0x6a, 0xab, 0xc1, 0x87, 0x9b, 0x78, 0x00, 0x28, 0xb9,
	0x73, 0x35, 0x5f, 0xcb, 0xa6, 0x75, 0x60, 0xd3, 0x47, 0x13, 0x34, 0xed, 0xb1, 0xa3, 0x3b, 0xb8,
	0x72, 0xf5, 0xfb, 0x88, 0xc8, 0x31, 0xee, 0x7d, 0x56, 0x26, 0x5f, 0xa4, 0xbd, 0x09, 0x9a, 0x1e,
	0xb0, 0xf8, 0x73, 0xf6, 0x1d, 0xe1, 0x61, 0xa8, 0x63, 0xb5, 0x6a, 0x2c, 0x90, 0x53, 0x3c, 0x94,
	0x96, 0x17, 0x5a, 0x1b, 0xb5, 0x81, 0x32, 0x74, 0x3a, 0x60, 0x58, 0xda, 0xc5, 0xad, 0x42, 0x5e,
	0xe1, 0x7e, 0xbc, 0xa5, 0x50, 0x64, 0x38, 0x27, 0x59, 0xbc, 0xbf, 0xcc, 0x68, 0x91, 0x2d, 0x03,
	0x61, 0xb7, 0x0e, 0x1f, 0x66, 0x65, 0xd5, 0x40, 0xc9, 0x05, 0x18, 0x17, 0xba, 0x8c, 0x18, 0x8e,
	0x52, 0x0e, 0xc6, 0x91, 0x17, 0x18, 0x7b, 0xc2, 0xc5, 0x65, 0x21, 0x9b, 0x70, 0xf8, 0x11, 0x1b,
	0x78, 0x25, 0xf7, 0xc2, 0x7c, 0x85, 0x0f, 0xdf, 0xf9, 0xad, 0xf2, 0xc5, 0x32, 0x2e, 0x48, 0x2e,
	0xf0, 0xe0, 0x6d, 0xd1, 0x94, 0x35, 0xe4, 0x4b, 0x46, 0x9e, 0x67, 0xff, 0x2d, 0x99, 0xed, 0x67,
	0x19, 0x8f, 0xef, 0x43, 0xb1, 0xe2, 0xc5, 0x9b, 0xab, 0x1d, 0x4d, 0xae, 0x77, 0x34, 0xf9, 0xb5,
	0xa3, 0xc9, 0xcd, 0x8e, 0x26, 0x5f, 0x3b, 0x8a, 0x7e, 0x74, 0x34, 0xb9, 0xea, 0x28, 0xba, 0xee,
	0x28, 0xfa, 0xdd, 0x51, 0xf4, 0xb7, 0xa3, 0xc9, 0x4d, 0x47, 0xd1, 0xb7, 0x3f, 0x34, 0xf9, 0x18,
	0x5f, 0x0c, 0xdf, 0x9c, 0x73, 0x1f, 0xf6, 0xa9, 0x1f, 0x5e, 0xc2, 0xeb, 0x7f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x51, 0xa6, 0x1d, 0xd1, 0x6a, 0x02, 0x00, 0x00,
}
