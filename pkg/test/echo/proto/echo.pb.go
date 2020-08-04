// Code generated by protoc-gen-go. DO NOT EDIT.
// source: echo.proto

// Generate with protoc --go_out=plugins=grpc:. echo.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EchoRequest struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoRequest) Reset()         { *m = EchoRequest{} }
func (m *EchoRequest) String() string { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()    {}
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{0}
}

func (m *EchoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoRequest.Unmarshal(m, b)
}
func (m *EchoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoRequest.Marshal(b, m, deterministic)
}
func (m *EchoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoRequest.Merge(m, src)
}
func (m *EchoRequest) XXX_Size() int {
	return xxx_messageInfo_EchoRequest.Size(m)
}
func (m *EchoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EchoRequest proto.InternalMessageInfo

func (m *EchoRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type EchoResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoResponse) Reset()         { *m = EchoResponse{} }
func (m *EchoResponse) String() string { return proto.CompactTextString(m) }
func (*EchoResponse) ProtoMessage()    {}
func (*EchoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{1}
}

func (m *EchoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoResponse.Unmarshal(m, b)
}
func (m *EchoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoResponse.Marshal(b, m, deterministic)
}
func (m *EchoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoResponse.Merge(m, src)
}
func (m *EchoResponse) XXX_Size() int {
	return xxx_messageInfo_EchoResponse.Size(m)
}
func (m *EchoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EchoResponse proto.InternalMessageInfo

func (m *EchoResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Header struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Header) Reset()         { *m = Header{} }
func (m *Header) String() string { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()    {}
func (*Header) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{2}
}

func (m *Header) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Header.Unmarshal(m, b)
}
func (m *Header) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Header.Marshal(b, m, deterministic)
}
func (m *Header) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Header.Merge(m, src)
}
func (m *Header) XXX_Size() int {
	return xxx_messageInfo_Header.Size(m)
}
func (m *Header) XXX_DiscardUnknown() {
	xxx_messageInfo_Header.DiscardUnknown(m)
}

var xxx_messageInfo_Header proto.InternalMessageInfo

func (m *Header) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Header) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type ForwardEchoRequest struct {
	Count         int32     `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	Qps           int32     `protobuf:"varint,2,opt,name=qps,proto3" json:"qps,omitempty"`
	TimeoutMicros int64     `protobuf:"varint,3,opt,name=timeout_micros,json=timeoutMicros,proto3" json:"timeout_micros,omitempty"`
	Url           string    `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	Headers       []*Header `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty"`
	Message       string    `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
	// If true, requests will be sent using h2c prior knowledge
	Http2 bool `protobuf:"varint,7,opt,name=http2,proto3" json:"http2,omitempty"`
	// If true, requests will not be sent until magic string is received
	ServerFirst          bool     `protobuf:"varint,8,opt,name=serverFirst,proto3" json:"serverFirst,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardEchoRequest) Reset()         { *m = ForwardEchoRequest{} }
func (m *ForwardEchoRequest) String() string { return proto.CompactTextString(m) }
func (*ForwardEchoRequest) ProtoMessage()    {}
func (*ForwardEchoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{3}
}

func (m *ForwardEchoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardEchoRequest.Unmarshal(m, b)
}
func (m *ForwardEchoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardEchoRequest.Marshal(b, m, deterministic)
}
func (m *ForwardEchoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardEchoRequest.Merge(m, src)
}
func (m *ForwardEchoRequest) XXX_Size() int {
	return xxx_messageInfo_ForwardEchoRequest.Size(m)
}
func (m *ForwardEchoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardEchoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardEchoRequest proto.InternalMessageInfo

func (m *ForwardEchoRequest) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

func (m *ForwardEchoRequest) GetQps() int32 {
	if m != nil {
		return m.Qps
	}
	return 0
}

func (m *ForwardEchoRequest) GetTimeoutMicros() int64 {
	if m != nil {
		return m.TimeoutMicros
	}
	return 0
}

func (m *ForwardEchoRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *ForwardEchoRequest) GetHeaders() []*Header {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *ForwardEchoRequest) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ForwardEchoRequest) GetHttp2() bool {
	if m != nil {
		return m.Http2
	}
	return false
}

func (m *ForwardEchoRequest) GetServerFirst() bool {
	if m != nil {
		return m.ServerFirst
	}
	return false
}

type ForwardEchoResponse struct {
	Output               []string `protobuf:"bytes,1,rep,name=output,proto3" json:"output,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForwardEchoResponse) Reset()         { *m = ForwardEchoResponse{} }
func (m *ForwardEchoResponse) String() string { return proto.CompactTextString(m) }
func (*ForwardEchoResponse) ProtoMessage()    {}
func (*ForwardEchoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{4}
}

func (m *ForwardEchoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForwardEchoResponse.Unmarshal(m, b)
}
func (m *ForwardEchoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForwardEchoResponse.Marshal(b, m, deterministic)
}
func (m *ForwardEchoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForwardEchoResponse.Merge(m, src)
}
func (m *ForwardEchoResponse) XXX_Size() int {
	return xxx_messageInfo_ForwardEchoResponse.Size(m)
}
func (m *ForwardEchoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ForwardEchoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ForwardEchoResponse proto.InternalMessageInfo

func (m *ForwardEchoResponse) GetOutput() []string {
	if m != nil {
		return m.Output
	}
	return nil
}

func init() {
	proto.RegisterType((*EchoRequest)(nil), "proto.EchoRequest")
	proto.RegisterType((*EchoResponse)(nil), "proto.EchoResponse")
	proto.RegisterType((*Header)(nil), "proto.Header")
	proto.RegisterType((*ForwardEchoRequest)(nil), "proto.ForwardEchoRequest")
	proto.RegisterType((*ForwardEchoResponse)(nil), "proto.ForwardEchoResponse")
}

func init() { proto.RegisterFile("echo.proto", fileDescriptor_08134aea513e0001) }

var fileDescriptor_08134aea513e0001 = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x50, 0x4d, 0x4f, 0xc2, 0x40,
	0x14, 0x4c, 0x2d, 0x2d, 0xf0, 0x2a, 0x6a, 0x16, 0x62, 0x56, 0x4e, 0x4d, 0x13, 0x43, 0x2f, 0xa2,
	0xc1, 0xbf, 0xa0, 0xc4, 0x8b, 0x97, 0xd5, 0xbb, 0xa9, 0xe5, 0xc5, 0x36, 0x02, 0x5b, 0xf6, 0x03,
	0xe3, 0x3f, 0xf0, 0xe7, 0xfa, 0x13, 0xcc, 0x7e, 0x90, 0xb4, 0xd1, 0x78, 0xea, 0x9b, 0x99, 0xd7,
	0xd9, 0x79, 0x03, 0x80, 0x65, 0xc5, 0xe7, 0x8d, 0xe0, 0x8a, 0x93, 0xc8, 0x7e, 0xb2, 0x19, 0x24,
	0xf7, 0x65, 0xc5, 0x19, 0xee, 0x34, 0x4a, 0x45, 0x28, 0xf4, 0x37, 0x28, 0x65, 0xf1, 0x86, 0x34,
	0x48, 0x83, 0x7c, 0xc8, 0x0e, 0x30, 0xcb, 0xe1, 0xd8, 0x2d, 0xca, 0x86, 0x6f, 0x25, 0xfe, 0xb3,
	0x79, 0x03, 0xf1, 0x03, 0x16, 0x2b, 0x14, 0xe4, 0x0c, 0xc2, 0x77, 0xfc, 0xf4, 0xba, 0x19, 0xc9,
	0x04, 0xa2, 0x7d, 0xb1, 0xd6, 0x48, 0x8f, 0x2c, 0xe7, 0x40, 0xf6, 0x1d, 0x00, 0x59, 0x72, 0xf1,
	0x51, 0x88, 0x55, 0x3b, 0xcc, 0x04, 0xa2, 0x92, 0xeb, 0xad, 0xb2, 0x06, 0x11, 0x73, 0xc0, 0x98,
	0xee, 0x1a, 0x69, 0x0d, 0x22, 0x66, 0x46, 0x72, 0x09, 0x27, 0xaa, 0xde, 0x20, 0xd7, 0xea, 0x65,
	0x53, 0x97, 0x82, 0x4b, 0x1a, 0xa6, 0x41, 0x1e, 0xb2, 0x91, 0x67, 0x1f, 0x2d, 0x69, 0x7e, 0xd4,
	0x62, 0x4d, 0x7b, 0x2e, 0x8d, 0x16, 0x6b, 0x32, 0x83, 0x7e, 0x65, 0x93, 0x4a, 0x1a, 0xa5, 0x61,
	0x9e, 0x2c, 0x46, 0xae, 0x9c, 0xb9, 0xcb, 0xcf, 0x0e, 0x6a, 0xfb, 0xd8, 0xb8, 0x73, 0xac, 0xc9,
	0x58, 0x29, 0xd5, 0x2c, 0x68, 0x3f, 0x0d, 0xf2, 0x01, 0x73, 0x80, 0xa4, 0x90, 0x48, 0x14, 0x7b,
	0x14, 0xcb, 0x5a, 0x48, 0x45, 0x07, 0x56, 0x6b, 0x53, 0xd9, 0x15, 0x8c, 0x3b, 0x17, 0xfb, 0x56,
	0xcf, 0x21, 0xe6, 0x5a, 0x35, 0xda, 0xdc, 0x1c, 0xe6, 0x43, 0xe6, 0xd1, 0xe2, 0x2b, 0x80, 0x53,
	0xb3, 0xf8, 0x8c, 0x52, 0x3d, 0xa1, 0xd8, 0xd7, 0x25, 0x92, 0x6b, 0xe8, 0x19, 0x8a, 0x10, 0x1f,
	0xba, 0x55, 0xdd, 0x74, 0xdc, 0xe1, 0xbc, 0xf9, 0x1d, 0x24, 0xad, 0x37, 0xc9, 0x85, 0xdf, 0xf9,
	0xdd, 0xfc, 0x74, 0xfa, 0x97, 0xe4, 0x5c, 0x5e, 0x63, 0x2b, 0xdd, 0xfe, 0x04, 0x00, 0x00, 0xff,
	0xff, 0x73, 0xa5, 0xa0, 0xe9, 0x4d, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EchoTestServiceClient is the client API for EchoTestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EchoTestServiceClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
	ForwardEcho(ctx context.Context, in *ForwardEchoRequest, opts ...grpc.CallOption) (*ForwardEchoResponse, error)
}

type echoTestServiceClient struct {
	cc *grpc.ClientConn
}

func NewEchoTestServiceClient(cc *grpc.ClientConn) EchoTestServiceClient {
	return &echoTestServiceClient{cc}
}

func (c *echoTestServiceClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, "/proto.EchoTestService/Echo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoTestServiceClient) ForwardEcho(ctx context.Context, in *ForwardEchoRequest, opts ...grpc.CallOption) (*ForwardEchoResponse, error) {
	out := new(ForwardEchoResponse)
	err := c.cc.Invoke(ctx, "/proto.EchoTestService/ForwardEcho", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoTestServiceServer is the server API for EchoTestService service.
type EchoTestServiceServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	ForwardEcho(context.Context, *ForwardEchoRequest) (*ForwardEchoResponse, error)
}

func RegisterEchoTestServiceServer(s *grpc.Server, srv EchoTestServiceServer) {
	s.RegisterService(&_EchoTestService_serviceDesc, srv)
}

func _EchoTestService_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoTestServiceServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.EchoTestService/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoTestServiceServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EchoTestService_ForwardEcho_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ForwardEchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoTestServiceServer).ForwardEcho(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.EchoTestService/ForwardEcho",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoTestServiceServer).ForwardEcho(ctx, req.(*ForwardEchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EchoTestService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.EchoTestService",
	HandlerType: (*EchoTestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _EchoTestService_Echo_Handler,
		},
		{
			MethodName: "ForwardEcho",
			Handler:    _EchoTestService_ForwardEcho_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "echo.proto",
}
