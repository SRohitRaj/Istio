// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: security/proto/providers/google/istioca.proto

package google_security_istioca_v1alpha1

import (
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"
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

// Certificate request message.
type IstioCertificateRequest struct {
	// PEM-encoded certificate request.
	Csr string `protobuf:"bytes,1,opt,name=csr,proto3" json:"csr,omitempty"`
	// Optional subject ID field.
	SubjectId string `protobuf:"bytes,2,opt,name=subject_id,json=subjectId,proto3" json:"subject_id,omitempty"`
	// Optional: requested certificate validity period, in seconds.
	ValidityDuration int64 `protobuf:"varint,3,opt,name=validity_duration,json=validityDuration,proto3" json:"validity_duration,omitempty"`
}

func (m *IstioCertificateRequest) Reset()      { *m = IstioCertificateRequest{} }
func (*IstioCertificateRequest) ProtoMessage() {}
func (*IstioCertificateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_479b8b5115786bf6, []int{0}
}
func (m *IstioCertificateRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IstioCertificateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IstioCertificateRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IstioCertificateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IstioCertificateRequest.Merge(m, src)
}
func (m *IstioCertificateRequest) XXX_Size() int {
	return m.Size()
}
func (m *IstioCertificateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IstioCertificateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IstioCertificateRequest proto.InternalMessageInfo

func (m *IstioCertificateRequest) GetCsr() string {
	if m != nil {
		return m.Csr
	}
	return ""
}

func (m *IstioCertificateRequest) GetSubjectId() string {
	if m != nil {
		return m.SubjectId
	}
	return ""
}

func (m *IstioCertificateRequest) GetValidityDuration() int64 {
	if m != nil {
		return m.ValidityDuration
	}
	return 0
}

// Certificate response message.
type IstioCertificateResponse struct {
	// PEM-encoded certificate chain.
	// Leaf cert is element '0'. Root cert is element 'n'.
	CertChain []string `protobuf:"bytes,1,rep,name=cert_chain,json=certChain,proto3" json:"cert_chain,omitempty"`
}

func (m *IstioCertificateResponse) Reset()      { *m = IstioCertificateResponse{} }
func (*IstioCertificateResponse) ProtoMessage() {}
func (*IstioCertificateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_479b8b5115786bf6, []int{1}
}
func (m *IstioCertificateResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IstioCertificateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IstioCertificateResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IstioCertificateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IstioCertificateResponse.Merge(m, src)
}
func (m *IstioCertificateResponse) XXX_Size() int {
	return m.Size()
}
func (m *IstioCertificateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IstioCertificateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IstioCertificateResponse proto.InternalMessageInfo

func (m *IstioCertificateResponse) GetCertChain() []string {
	if m != nil {
		return m.CertChain
	}
	return nil
}

func init() {
	proto.RegisterType((*IstioCertificateRequest)(nil), "google.security.istioca.v1alpha1.IstioCertificateRequest")
	proto.RegisterType((*IstioCertificateResponse)(nil), "google.security.istioca.v1alpha1.IstioCertificateResponse")
}

func init() {
	proto.RegisterFile("security/proto/providers/google/istioca.proto", fileDescriptor_479b8b5115786bf6)
}

var fileDescriptor_479b8b5115786bf6 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x91, 0xbf, 0x4a, 0x7b, 0x31,
	0x14, 0xc7, 0x73, 0x7e, 0x17, 0x7e, 0xd0, 0x4c, 0xed, 0x5d, 0xbc, 0x08, 0x1e, 0x2e, 0x9d, 0x0a,
	0xe2, 0x2d, 0x55, 0x97, 0x3a, 0x5a, 0x97, 0xae, 0xf5, 0x01, 0x4a, 0x9a, 0x7b, 0x6c, 0x23, 0xa5,
	0xa9, 0x49, 0xee, 0x95, 0x6e, 0x3e, 0x80, 0x83, 0x8f, 0xd1, 0x47, 0x71, 0xec, 0xd8, 0xd1, 0xa6,
	0x8b, 0x63, 0x1f, 0x41, 0xd2, 0x3f, 0x20, 0x14, 0x11, 0x5c, 0x42, 0xf8, 0x7c, 0x38, 0x7c, 0x93,
	0xef, 0xe1, 0x17, 0x96, 0x64, 0x61, 0x94, 0x9b, 0x35, 0xa7, 0x46, 0x3b, 0x1d, 0xce, 0x52, 0xe5,
	0x64, 0x6c, 0x73, 0xa8, 0xf5, 0x70, 0x4c, 0x4d, 0x65, 0x9d, 0xd2, 0x52, 0x64, 0x5b, 0x1d, 0xa7,
	0x3b, 0x9a, 0x1d, 0xa6, 0xb2, 0x83, 0x2e, 0x5b, 0x62, 0x3c, 0x1d, 0x89, 0x56, 0xfd, 0x99, 0x9f,
	0x74, 0x03, 0xeb, 0x90, 0x71, 0xea, 0x41, 0x49, 0xe1, 0xa8, 0x47, 0x4f, 0x05, 0x59, 0x17, 0x57,
	0x79, 0x24, 0xad, 0x49, 0x20, 0x85, 0x46, 0xa5, 0x17, 0xae, 0xf1, 0x19, 0xe7, 0xb6, 0x18, 0x3c,
	0x92, 0x74, 0x7d, 0x95, 0x27, 0xff, 0xb6, 0xa2, 0xb2, 0x27, 0xdd, 0x3c, 0x3e, 0xe7, 0xb5, 0x52,
	0x8c, 0x55, 0xae, 0xdc, 0xac, 0x9f, 0x17, 0x46, 0x38, 0xa5, 0x27, 0x49, 0x94, 0x42, 0x23, 0xea,
	0x55, 0x0f, 0xe2, 0x6e, 0xcf, 0xeb, 0x6d, 0x9e, 0x1c, 0x07, 0xdb, 0xa9, 0x9e, 0x58, 0x0a, 0x39,
	0x92, 0x8c, 0xeb, 0xcb, 0x91, 0x50, 0x93, 0x04, 0xd2, 0x28, 0xe4, 0x04, 0xd2, 0x09, 0xe0, 0x72,
	0x0e, 0xc7, 0x8f, 0xbe, 0x27, 0x53, 0x2a, 0x49, 0xf1, 0x2b, 0xf0, 0x5a, 0xc7, 0x90, 0x70, 0xf4,
	0x4d, 0xc6, 0xed, 0xec, 0xb7, 0x22, 0xb2, 0x1f, 0x5a, 0x38, 0xbd, 0xf9, 0xcb, 0xe8, 0xee, 0x1f,
	0x75, 0x76, 0x7b, 0xbd, 0x58, 0x21, 0x5b, 0xae, 0x90, 0x6d, 0x56, 0x08, 0x2f, 0x1e, 0x61, 0xee,
	0x11, 0xde, 0x3d, 0xc2, 0xc2, 0x23, 0x7c, 0x78, 0x84, 0x4f, 0x8f, 0x6c, 0xe3, 0x11, 0xde, 0xd6,
	0xc8, 0x16, 0x6b, 0x64, 0xcb, 0x35, 0xb2, 0xc1, 0xff, 0xed, 0xf6, 0xae, 0xbe, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x99, 0xf6, 0x07, 0x66, 0xee, 0x01, 0x00, 0x00,
}

func (this *IstioCertificateRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IstioCertificateRequest)
	if !ok {
		that2, ok := that.(IstioCertificateRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Csr != that1.Csr {
		return false
	}
	if this.SubjectId != that1.SubjectId {
		return false
	}
	if this.ValidityDuration != that1.ValidityDuration {
		return false
	}
	return true
}
func (this *IstioCertificateResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IstioCertificateResponse)
	if !ok {
		that2, ok := that.(IstioCertificateResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.CertChain) != len(that1.CertChain) {
		return false
	}
	for i := range this.CertChain {
		if this.CertChain[i] != that1.CertChain[i] {
			return false
		}
	}
	return true
}
func (this *IstioCertificateRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&google_security_istioca_v1alpha1.IstioCertificateRequest{")
	s = append(s, "Csr: "+fmt.Sprintf("%#v", this.Csr)+",\n")
	s = append(s, "SubjectId: "+fmt.Sprintf("%#v", this.SubjectId)+",\n")
	s = append(s, "ValidityDuration: "+fmt.Sprintf("%#v", this.ValidityDuration)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *IstioCertificateResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&google_security_istioca_v1alpha1.IstioCertificateResponse{")
	s = append(s, "CertChain: "+fmt.Sprintf("%#v", this.CertChain)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringIstioca(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IstioCertificateServiceClient is the client API for IstioCertificateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IstioCertificateServiceClient interface {
	// Using provided CSR, returns a signed certificate.
	CreateCertificate(ctx context.Context, in *IstioCertificateRequest, opts ...grpc.CallOption) (*IstioCertificateResponse, error)
}

type istioCertificateServiceClient struct {
	cc *grpc.ClientConn
}

func NewIstioCertificateServiceClient(cc *grpc.ClientConn) IstioCertificateServiceClient {
	return &istioCertificateServiceClient{cc}
}

func (c *istioCertificateServiceClient) CreateCertificate(ctx context.Context, in *IstioCertificateRequest, opts ...grpc.CallOption) (*IstioCertificateResponse, error) {
	out := new(IstioCertificateResponse)
	err := c.cc.Invoke(ctx, "/google.security.istioca.v1alpha1.IstioCertificateService/CreateCertificate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IstioCertificateServiceServer is the server API for IstioCertificateService service.
type IstioCertificateServiceServer interface {
	// Using provided CSR, returns a signed certificate.
	CreateCertificate(context.Context, *IstioCertificateRequest) (*IstioCertificateResponse, error)
}

func RegisterIstioCertificateServiceServer(s *grpc.Server, srv IstioCertificateServiceServer) {
	s.RegisterService(&_IstioCertificateService_serviceDesc, srv)
}

func _IstioCertificateService_CreateCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IstioCertificateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IstioCertificateServiceServer).CreateCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.security.istioca.v1alpha1.IstioCertificateService/CreateCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IstioCertificateServiceServer).CreateCertificate(ctx, req.(*IstioCertificateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IstioCertificateService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.security.istioca.v1alpha1.IstioCertificateService",
	HandlerType: (*IstioCertificateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCertificate",
			Handler:    _IstioCertificateService_CreateCertificate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "security/proto/providers/google/istioca.proto",
}

func (m *IstioCertificateRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IstioCertificateRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Csr) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintIstioca(dAtA, i, uint64(len(m.Csr)))
		i += copy(dAtA[i:], m.Csr)
	}
	if len(m.SubjectId) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintIstioca(dAtA, i, uint64(len(m.SubjectId)))
		i += copy(dAtA[i:], m.SubjectId)
	}
	if m.ValidityDuration != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintIstioca(dAtA, i, uint64(m.ValidityDuration))
	}
	return i, nil
}

func (m *IstioCertificateResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IstioCertificateResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.CertChain) > 0 {
		for _, s := range m.CertChain {
			dAtA[i] = 0xa
			i++
			l = len(s)
			for l >= 1<<7 {
				dAtA[i] = uint8(uint64(l)&0x7f | 0x80)
				l >>= 7
				i++
			}
			dAtA[i] = uint8(l)
			i++
			i += copy(dAtA[i:], s)
		}
	}
	return i, nil
}

func encodeVarintIstioca(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *IstioCertificateRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Csr)
	if l > 0 {
		n += 1 + l + sovIstioca(uint64(l))
	}
	l = len(m.SubjectId)
	if l > 0 {
		n += 1 + l + sovIstioca(uint64(l))
	}
	if m.ValidityDuration != 0 {
		n += 1 + sovIstioca(uint64(m.ValidityDuration))
	}
	return n
}

func (m *IstioCertificateResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.CertChain) > 0 {
		for _, s := range m.CertChain {
			l = len(s)
			n += 1 + l + sovIstioca(uint64(l))
		}
	}
	return n
}

func sovIstioca(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozIstioca(x uint64) (n int) {
	return sovIstioca(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *IstioCertificateRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&IstioCertificateRequest{`,
		`Csr:` + fmt.Sprintf("%v", this.Csr) + `,`,
		`SubjectId:` + fmt.Sprintf("%v", this.SubjectId) + `,`,
		`ValidityDuration:` + fmt.Sprintf("%v", this.ValidityDuration) + `,`,
		`}`,
	}, "")
	return s
}
func (this *IstioCertificateResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&IstioCertificateResponse{`,
		`CertChain:` + fmt.Sprintf("%v", this.CertChain) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringIstioca(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *IstioCertificateRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIstioca
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
			return fmt.Errorf("proto: IstioCertificateRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IstioCertificateRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Csr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIstioca
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIstioca
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIstioca
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Csr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubjectId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIstioca
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIstioca
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIstioca
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubjectId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidityDuration", wireType)
			}
			m.ValidityDuration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIstioca
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidityDuration |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipIstioca(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIstioca
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIstioca
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
func (m *IstioCertificateResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIstioca
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
			return fmt.Errorf("proto: IstioCertificateResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IstioCertificateResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CertChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIstioca
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIstioca
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIstioca
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CertChain = append(m.CertChain, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIstioca(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIstioca
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIstioca
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
func skipIstioca(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIstioca
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
					return 0, ErrIntOverflowIstioca
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
					return 0, ErrIntOverflowIstioca
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
				return 0, ErrInvalidLengthIstioca
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthIstioca
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowIstioca
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
				next, err := skipIstioca(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthIstioca
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
	ErrInvalidLengthIstioca = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIstioca   = fmt.Errorf("proto: integer overflow")
)
