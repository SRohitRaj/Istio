// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/type/http_status.proto

package envoy_type

import (
	fmt "fmt"
	io "io"
	math "math"

	proto "github.com/gogo/protobuf/proto"
	_ "github.com/lyft/protoc-gen-validate/validate"
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

// HTTP response codes supported in Envoy.
// For more details: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
type StatusCode int32

const (
	// Empty - This code not part of the HTTP status code specification, but it is needed for proto
	// `enum` type.
	StatusCode_Empty                         StatusCode = 0
	StatusCode_Continue                      StatusCode = 100
	StatusCode_OK                            StatusCode = 200
	StatusCode_Created                       StatusCode = 201
	StatusCode_Accepted                      StatusCode = 202
	StatusCode_NonAuthoritativeInformation   StatusCode = 203
	StatusCode_NoContent                     StatusCode = 204
	StatusCode_ResetContent                  StatusCode = 205
	StatusCode_PartialContent                StatusCode = 206
	StatusCode_MultiStatus                   StatusCode = 207
	StatusCode_AlreadyReported               StatusCode = 208
	StatusCode_IMUsed                        StatusCode = 226
	StatusCode_MultipleChoices               StatusCode = 300
	StatusCode_MovedPermanently              StatusCode = 301
	StatusCode_Found                         StatusCode = 302
	StatusCode_SeeOther                      StatusCode = 303
	StatusCode_NotModified                   StatusCode = 304
	StatusCode_UseProxy                      StatusCode = 305
	StatusCode_TemporaryRedirect             StatusCode = 307
	StatusCode_PermanentRedirect             StatusCode = 308
	StatusCode_BadRequest                    StatusCode = 400
	StatusCode_Unauthorized                  StatusCode = 401
	StatusCode_PaymentRequired               StatusCode = 402
	StatusCode_Forbidden                     StatusCode = 403
	StatusCode_NotFound                      StatusCode = 404
	StatusCode_MethodNotAllowed              StatusCode = 405
	StatusCode_NotAcceptable                 StatusCode = 406
	StatusCode_ProxyAuthenticationRequired   StatusCode = 407
	StatusCode_RequestTimeout                StatusCode = 408
	StatusCode_Conflict                      StatusCode = 409
	StatusCode_Gone                          StatusCode = 410
	StatusCode_LengthRequired                StatusCode = 411
	StatusCode_PreconditionFailed            StatusCode = 412
	StatusCode_PayloadTooLarge               StatusCode = 413
	StatusCode_URITooLong                    StatusCode = 414
	StatusCode_UnsupportedMediaType          StatusCode = 415
	StatusCode_RangeNotSatisfiable           StatusCode = 416
	StatusCode_ExpectationFailed             StatusCode = 417
	StatusCode_MisdirectedRequest            StatusCode = 421
	StatusCode_UnprocessableEntity           StatusCode = 422
	StatusCode_Locked                        StatusCode = 423
	StatusCode_FailedDependency              StatusCode = 424
	StatusCode_UpgradeRequired               StatusCode = 426
	StatusCode_PreconditionRequired          StatusCode = 428
	StatusCode_TooManyRequests               StatusCode = 429
	StatusCode_RequestHeaderFieldsTooLarge   StatusCode = 431
	StatusCode_InternalServerError           StatusCode = 500
	StatusCode_NotImplemented                StatusCode = 501
	StatusCode_BadGateway                    StatusCode = 502
	StatusCode_ServiceUnavailable            StatusCode = 503
	StatusCode_GatewayTimeout                StatusCode = 504
	StatusCode_HTTPVersionNotSupported       StatusCode = 505
	StatusCode_VariantAlsoNegotiates         StatusCode = 506
	StatusCode_InsufficientStorage           StatusCode = 507
	StatusCode_LoopDetected                  StatusCode = 508
	StatusCode_NotExtended                   StatusCode = 510
	StatusCode_NetworkAuthenticationRequired StatusCode = 511
)

var StatusCode_name = map[int32]string{
	0:   "Empty",
	100: "Continue",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "NonAuthoritativeInformation",
	204: "NoContent",
	205: "ResetContent",
	206: "PartialContent",
	207: "MultiStatus",
	208: "AlreadyReported",
	226: "IMUsed",
	300: "MultipleChoices",
	301: "MovedPermanently",
	302: "Found",
	303: "SeeOther",
	304: "NotModified",
	305: "UseProxy",
	307: "TemporaryRedirect",
	308: "PermanentRedirect",
	400: "BadRequest",
	401: "Unauthorized",
	402: "PaymentRequired",
	403: "Forbidden",
	404: "NotFound",
	405: "MethodNotAllowed",
	406: "NotAcceptable",
	407: "ProxyAuthenticationRequired",
	408: "RequestTimeout",
	409: "Conflict",
	410: "Gone",
	411: "LengthRequired",
	412: "PreconditionFailed",
	413: "PayloadTooLarge",
	414: "URITooLong",
	415: "UnsupportedMediaType",
	416: "RangeNotSatisfiable",
	417: "ExpectationFailed",
	421: "MisdirectedRequest",
	422: "UnprocessableEntity",
	423: "Locked",
	424: "FailedDependency",
	426: "UpgradeRequired",
	428: "PreconditionRequired",
	429: "TooManyRequests",
	431: "RequestHeaderFieldsTooLarge",
	500: "InternalServerError",
	501: "NotImplemented",
	502: "BadGateway",
	503: "ServiceUnavailable",
	504: "GatewayTimeout",
	505: "HTTPVersionNotSupported",
	506: "VariantAlsoNegotiates",
	507: "InsufficientStorage",
	508: "LoopDetected",
	510: "NotExtended",
	511: "NetworkAuthenticationRequired",
}

var StatusCode_value = map[string]int32{
	"Empty":                         0,
	"Continue":                      100,
	"OK":                            200,
	"Created":                       201,
	"Accepted":                      202,
	"NonAuthoritativeInformation":   203,
	"NoContent":                     204,
	"ResetContent":                  205,
	"PartialContent":                206,
	"MultiStatus":                   207,
	"AlreadyReported":               208,
	"IMUsed":                        226,
	"MultipleChoices":               300,
	"MovedPermanently":              301,
	"Found":                         302,
	"SeeOther":                      303,
	"NotModified":                   304,
	"UseProxy":                      305,
	"TemporaryRedirect":             307,
	"PermanentRedirect":             308,
	"BadRequest":                    400,
	"Unauthorized":                  401,
	"PaymentRequired":               402,
	"Forbidden":                     403,
	"NotFound":                      404,
	"MethodNotAllowed":              405,
	"NotAcceptable":                 406,
	"ProxyAuthenticationRequired":   407,
	"RequestTimeout":                408,
	"Conflict":                      409,
	"Gone":                          410,
	"LengthRequired":                411,
	"PreconditionFailed":            412,
	"PayloadTooLarge":               413,
	"URITooLong":                    414,
	"UnsupportedMediaType":          415,
	"RangeNotSatisfiable":           416,
	"ExpectationFailed":             417,
	"MisdirectedRequest":            421,
	"UnprocessableEntity":           422,
	"Locked":                        423,
	"FailedDependency":              424,
	"UpgradeRequired":               426,
	"PreconditionRequired":          428,
	"TooManyRequests":               429,
	"RequestHeaderFieldsTooLarge":   431,
	"InternalServerError":           500,
	"NotImplemented":                501,
	"BadGateway":                    502,
	"ServiceUnavailable":            503,
	"GatewayTimeout":                504,
	"HTTPVersionNotSupported":       505,
	"VariantAlsoNegotiates":         506,
	"InsufficientStorage":           507,
	"LoopDetected":                  508,
	"NotExtended":                   510,
	"NetworkAuthenticationRequired": 511,
}

func (x StatusCode) String() string {
	return proto.EnumName(StatusCode_name, int32(x))
}

func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7544d7adacd3389b, []int{0}
}

// HTTP status.
type HttpStatus struct {
	// Supplies HTTP response code.
	Code                 StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=envoy.type.StatusCode" json:"code,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *HttpStatus) Reset()         { *m = HttpStatus{} }
func (m *HttpStatus) String() string { return proto.CompactTextString(m) }
func (*HttpStatus) ProtoMessage()    {}
func (*HttpStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_7544d7adacd3389b, []int{0}
}
func (m *HttpStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HttpStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HttpStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HttpStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HttpStatus.Merge(m, src)
}
func (m *HttpStatus) XXX_Size() int {
	return m.Size()
}
func (m *HttpStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_HttpStatus.DiscardUnknown(m)
}

var xxx_messageInfo_HttpStatus proto.InternalMessageInfo

func (m *HttpStatus) GetCode() StatusCode {
	if m != nil {
		return m.Code
	}
	return StatusCode_Empty
}

func init() {
	proto.RegisterEnum("envoy.type.StatusCode", StatusCode_name, StatusCode_value)
	proto.RegisterType((*HttpStatus)(nil), "envoy.type.HttpStatus")
}

func init() { proto.RegisterFile("envoy/type/http_status.proto", fileDescriptor_7544d7adacd3389b) }

var fileDescriptor_7544d7adacd3389b = []byte{
	// 931 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0x49, 0x73, 0x5b, 0x45,
	0x10, 0xce, 0xd3, 0x64, 0xf3, 0x64, 0xeb, 0x4c, 0x9c, 0xd8, 0x84, 0xe0, 0x72, 0xe5, 0x44, 0x71,
	0xb0, 0xab, 0xe0, 0xc4, 0x0d, 0xdb, 0xb1, 0x63, 0x17, 0x96, 0xa2, 0x92, 0xa5, 0x5c, 0xa9, 0xf1,
	0x9b, 0x96, 0x34, 0x95, 0xa7, 0xe9, 0x97, 0x79, 0x2d, 0xd9, 0x8f, 0x23, 0xbf, 0x80, 0x7d, 0x5f,
	0x0f, 0x2c, 0x45, 0x25, 0x04, 0x0a, 0xb8, 0x70, 0xe2, 0x18, 0xf6, 0xfc, 0x04, 0xca, 0x37, 0xee,
	0xec, 0x50, 0x40, 0xcd, 0x68, 0xb1, 0x2f, 0xdc, 0xf4, 0x7a, 0x7a, 0xf9, 0x96, 0x56, 0xcb, 0x4b,
	0xe8, 0x06, 0x54, 0x2e, 0x72, 0x99, 0xe3, 0x62, 0x97, 0x39, 0x7f, 0xa2, 0x60, 0xcd, 0xfd, 0x62,
	0x21, 0xf7, 0xc4, 0xa4, 0x64, 0x7c, 0x5d, 0x08, 0xaf, 0x17, 0x67, 0x06, 0x3a, 0xb3, 0x46, 0x33,
	0x2e, 0x8e, 0x7f, 0x0c, 0x93, 0x2e, 0xd7, 0xa4, 0x5c, 0x67, 0xce, 0xb7, 0x62, 0xa1, 0x7a, 0x4c,
	0x1e, 0x4e, 0xc9, 0xe0, 0x6c, 0x32, 0x9f, 0x3c, 0x78, 0xfa, 0xe1, 0x0b, 0x0b, 0xfb, 0x1d, 0x16,
	0x86, 0x19, 0x2b, 0x64, 0x70, 0x79, 0xfa, 0x8b, 0x9f, 0xbe, 0x14, 0x47, 0x9e, 0x4a, 0x2a, 0xf3,
	0x87, 0xc6, 0xbf, 0x20, 0x69, 0xc4, 0xca, 0x87, 0x3e, 0x9f, 0x92, 0x72, 0x3f, 0x55, 0x4d, 0xc9,
	0x23, 0xab, 0xbd, 0x9c, 0x4b, 0x38, 0xa4, 0x4e, 0xca, 0xe3, 0x2b, 0xe4, 0xd8, 0xba, 0x3e, 0x82,
	0x51, 0xc7, 0x64, 0xe5, 0xda, 0xe3, 0x70, 0x37, 0x51, 0x27, 0xe5, 0xb1, 0x15, 0x8f, 0x9a, 0xd1,
	0xc0, 0x57, 0x89, 0x3a, 0x25, 0x8f, 0x2f, 0xa5, 0x29, 0xe6, 0xe1, 0xf3, 0xeb, 0x44, 0xcd, 0xcb,
	0xfb, 0x6b, 0xe4, 0x96, 0xfa, 0xdc, 0x25, 0x6f, 0x59, 0xb3, 0x1d, 0xe0, 0x86, 0x6b, 0x93, 0xef,
	0x69, 0xb6, 0xe4, 0xe0, 0x9b, 0x44, 0x9d, 0x96, 0x53, 0x35, 0x0a, 0x7d, 0xd1, 0x31, 0x7c, 0x9b,
	0xa8, 0xb3, 0xf2, 0x64, 0x03, 0x0b, 0xe4, 0x71, 0xe8, 0xbb, 0x44, 0x9d, 0x93, 0xa7, 0xeb, 0xda,
	0xb3, 0xd5, 0xd9, 0x38, 0xf8, 0x7d, 0xa2, 0x40, 0x9e, 0xa8, 0xf6, 0x33, 0xb6, 0x43, 0xac, 0xf0,
	0x43, 0xa2, 0xa6, 0xe5, 0x99, 0xa5, 0xcc, 0xa3, 0x36, 0x65, 0x03, 0x73, 0xf2, 0x01, 0xc1, 0xbd,
	0x44, 0x9d, 0x90, 0x47, 0x37, 0xaa, 0xad, 0x02, 0x0d, 0xec, 0xc5, 0x94, 0x58, 0x94, 0x67, 0xb8,
	0xd2, 0x25, 0x9b, 0x62, 0x01, 0xb7, 0x2a, 0xea, 0xbc, 0x84, 0x2a, 0x0d, 0xd0, 0xd4, 0xd1, 0xf7,
	0xb4, 0x43, 0xc7, 0x59, 0x09, 0xb7, 0x2b, 0x4a, 0xca, 0x23, 0x6b, 0xd4, 0x77, 0x06, 0x3e, 0xaa,
	0x04, 0x5a, 0x5b, 0x88, 0xd7, 0xb8, 0x8b, 0x1e, 0xee, 0x54, 0xc2, 0xf0, 0x1a, 0x71, 0x95, 0x8c,
	0x6d, 0x5b, 0x34, 0xf0, 0x71, 0x4c, 0x68, 0x15, 0x58, 0xf7, 0xb4, 0x5b, 0xc2, 0x27, 0x15, 0x75,
	0x41, 0x9e, 0x6d, 0x62, 0x2f, 0x27, 0xaf, 0x7d, 0xd9, 0x40, 0x63, 0x3d, 0xa6, 0x0c, 0x9f, 0xc6,
	0xf8, 0x64, 0xca, 0x24, 0xfe, 0x59, 0x45, 0x9d, 0x91, 0x72, 0x59, 0x9b, 0x06, 0xde, 0xec, 0x63,
	0xc1, 0xf0, 0xb4, 0x08, 0x32, 0xb4, 0x9c, 0x1e, 0xea, 0xf6, 0x24, 0x1a, 0x78, 0x46, 0x04, 0xf0,
	0x75, 0x5d, 0xf6, 0x62, 0xe5, 0xcd, 0xbe, 0xf5, 0x68, 0xe0, 0x59, 0x11, 0xf4, 0x5b, 0x23, 0xbf,
	0x6d, 0x8d, 0x41, 0x07, 0xcf, 0x89, 0x00, 0xa4, 0x46, 0x3c, 0x04, 0xfe, 0xbc, 0x88, 0xdc, 0x90,
	0xbb, 0x64, 0x6a, 0xc4, 0x4b, 0x59, 0x46, 0x3b, 0x68, 0xe0, 0x05, 0xa1, 0x94, 0x3c, 0x15, 0x02,
	0xd1, 0x29, 0xbd, 0x9d, 0x21, 0xbc, 0x28, 0x82, 0x57, 0x11, 0x7f, 0x70, 0x0b, 0x1d, 0xdb, 0x34,
	0x7a, 0x34, 0x99, 0xf5, 0x92, 0x08, 0x46, 0x8c, 0x20, 0x36, 0x6d, 0x0f, 0xa9, 0xcf, 0xf0, 0x72,
	0x1c, 0xb8, 0x42, 0xae, 0x9d, 0xd9, 0x94, 0xe1, 0x15, 0xa1, 0xa6, 0xe4, 0xe1, 0xab, 0xe4, 0x10,
	0x5e, 0x8d, 0xe9, 0x9b, 0xe8, 0x3a, 0xdc, 0x9d, 0xf4, 0x78, 0x4d, 0xa8, 0x19, 0xa9, 0xea, 0x1e,
	0x53, 0x72, 0xc6, 0x86, 0xf6, 0x6b, 0xda, 0x66, 0x68, 0xe0, 0xf5, 0x31, 0xbd, 0x8c, 0xb4, 0x69,
	0x12, 0x6d, 0x6a, 0xdf, 0x41, 0x78, 0x43, 0x04, 0x61, 0x5a, 0x8d, 0x8d, 0x10, 0x21, 0xd7, 0x81,
	0x37, 0x85, 0xba, 0x4f, 0x4e, 0xb7, 0x5c, 0xd1, 0xcf, 0x87, 0x0e, 0x57, 0xd1, 0x58, 0xdd, 0x2c,
	0x73, 0x84, 0xb7, 0x84, 0x9a, 0x95, 0xe7, 0x1a, 0xda, 0x75, 0xb0, 0x46, 0xbc, 0xa5, 0xd9, 0x16,
	0x6d, 0x1b, 0xa9, 0xbd, 0x2d, 0x82, 0xec, 0xab, 0xbb, 0x39, 0xa6, 0xac, 0x0f, 0xcc, 0x7c, 0x27,
	0x82, 0xa9, 0xda, 0x62, 0x68, 0x03, 0x4e, 0xe4, 0x7f, 0x37, 0xb6, 0x6a, 0xb9, 0xdc, 0x53, 0x8a,
	0x45, 0x11, 0x9a, 0xac, 0x3a, 0xb6, 0x5c, 0xc2, 0x7b, 0x22, 0xec, 0xd3, 0x26, 0xa5, 0x37, 0xd0,
	0xc0, 0xfb, 0x51, 0xdd, 0x61, 0xb3, 0x2b, 0x98, 0xa3, 0x33, 0xe8, 0xd2, 0x12, 0x3e, 0x88, 0x54,
	0x5a, 0x79, 0xc7, 0x6b, 0x83, 0x13, 0xe6, 0x1f, 0x46, 0xe4, 0x07, 0x99, 0x4f, 0x9e, 0x6e, 0xc5,
	0x82, 0x26, 0x51, 0x55, 0xbb, 0x72, 0x84, 0xa1, 0x80, 0xdb, 0xd1, 0x90, 0xd1, 0xe7, 0x3a, 0x6a,
	0x83, 0x7e, 0xcd, 0x62, 0x66, 0x8a, 0x89, 0x3a, 0x77, 0x22, 0xcc, 0x0d, 0xc7, 0xe8, 0x9d, 0xce,
	0xb6, 0xd0, 0x0f, 0xd0, 0xaf, 0x7a, 0x4f, 0x1e, 0x7e, 0x8e, 0xda, 0xd7, 0x88, 0x37, 0x7a, 0x79,
	0x86, 0x61, 0x63, 0xd0, 0xc0, 0x2f, 0x62, 0xb4, 0x65, 0x57, 0x35, 0xe3, 0x8e, 0x2e, 0xe1, 0xd7,
	0xc8, 0x3f, 0xd4, 0xd9, 0x14, 0x5b, 0x4e, 0x0f, 0xb4, 0xcd, 0xa2, 0x60, 0xbf, 0xc5, 0xf2, 0x51,
	0xda, 0xd8, 0xe9, 0xdf, 0x85, 0xba, 0x24, 0x67, 0xd6, 0x9b, 0xcd, 0xfa, 0x75, 0xf4, 0x85, 0x25,
	0x17, 0x54, 0x1e, 0xdb, 0x00, 0x7f, 0x08, 0x75, 0x51, 0x9e, 0xbf, 0xae, 0xbd, 0xd5, 0x8e, 0x97,
	0xb2, 0x82, 0x6a, 0xd8, 0x21, 0xb6, 0x9a, 0xb1, 0x80, 0x3f, 0x47, 0x38, 0x8b, 0x7e, 0xbb, 0x6d,
	0x53, 0x8b, 0x8e, 0xb7, 0x98, 0xbc, 0xee, 0x20, 0xfc, 0x15, 0xf7, 0x7c, 0x93, 0x28, 0xbf, 0x82,
	0x1c, 0x2d, 0x80, 0xbf, 0xc5, 0xe8, 0xcf, 0xb5, 0xba, 0xcb, 0x41, 0x51, 0x03, 0xff, 0x08, 0x75,
	0x59, 0x3e, 0x50, 0x43, 0xde, 0x21, 0x7f, 0xe3, 0x7f, 0x76, 0xf3, 0x5f, 0xb1, 0xfc, 0xe8, 0xdd,
	0xbd, 0xb9, 0xe4, 0xde, 0xde, 0x5c, 0xf2, 0xe3, 0xde, 0x5c, 0x22, 0x67, 0x2d, 0x0d, 0x6f, 0x5f,
	0x1e, 0x36, 0xfa, 0xc0, 0x19, 0x5c, 0x3e, 0xb3, 0x7f, 0x2d, 0xeb, 0xe1, 0x80, 0xd6, 0x93, 0xed,
	0xa3, 0xf1, 0x92, 0x3e, 0xf2, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x6d, 0xca, 0xec, 0x2e, 0x8e,
	0x05, 0x00, 0x00,
}

func (m *HttpStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HttpStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Code != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintHttpStatus(dAtA, i, uint64(m.Code))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintHttpStatus(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *HttpStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Code != 0 {
		n += 1 + sovHttpStatus(uint64(m.Code))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovHttpStatus(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozHttpStatus(x uint64) (n int) {
	return sovHttpStatus(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HttpStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHttpStatus
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
			return fmt.Errorf("proto: HttpStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HttpStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHttpStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= StatusCode(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipHttpStatus(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHttpStatus
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthHttpStatus
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
func skipHttpStatus(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHttpStatus
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
					return 0, ErrIntOverflowHttpStatus
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
					return 0, ErrIntOverflowHttpStatus
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
				return 0, ErrInvalidLengthHttpStatus
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthHttpStatus
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowHttpStatus
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
				next, err := skipHttpStatus(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthHttpStatus
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
	ErrInvalidLengthHttpStatus = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHttpStatus   = fmt.Errorf("proto: integer overflow")
)
