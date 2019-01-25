// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/network/mysql_proxy/v1alpha1/mysql_proxy.proto

package v1alpha1

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
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

// [#protodoc-title: MySQL proxy]
// MySQL Proxy :ref:`configuration overview <config_network_filters_mysql_proxy>`.
type MySQLProxy struct {
	// The human readable prefix to use when emitting :ref:`statistics
	// <config_network_filters_mysql_proxy_stats>`.
	StatPrefix string `protobuf:"bytes,1,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	// [#not-implemented-hide:] The optional path to use for writing MySQL access logs.
	// If the access log field is empty, access logs will not be written.
	AccessLog            string   `protobuf:"bytes,2,opt,name=access_log,json=accessLog,proto3" json:"access_log,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MySQLProxy) Reset()         { *m = MySQLProxy{} }
func (m *MySQLProxy) String() string { return proto.CompactTextString(m) }
func (*MySQLProxy) ProtoMessage()    {}
func (*MySQLProxy) Descriptor() ([]byte, []int) {
	return fileDescriptor_mysql_proxy_ea6ac70602b155ef, []int{0}
}
func (m *MySQLProxy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MySQLProxy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MySQLProxy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *MySQLProxy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MySQLProxy.Merge(dst, src)
}
func (m *MySQLProxy) XXX_Size() int {
	return m.Size()
}
func (m *MySQLProxy) XXX_DiscardUnknown() {
	xxx_messageInfo_MySQLProxy.DiscardUnknown(m)
}

var xxx_messageInfo_MySQLProxy proto.InternalMessageInfo

func (m *MySQLProxy) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *MySQLProxy) GetAccessLog() string {
	if m != nil {
		return m.AccessLog
	}
	return ""
}

func init() {
	proto.RegisterType((*MySQLProxy)(nil), "envoy.config.filter.network.mysql_proxy.v1alpha1.MySQLProxy")
}
func (m *MySQLProxy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MySQLProxy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.StatPrefix) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMysqlProxy(dAtA, i, uint64(len(m.StatPrefix)))
		i += copy(dAtA[i:], m.StatPrefix)
	}
	if len(m.AccessLog) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintMysqlProxy(dAtA, i, uint64(len(m.AccessLog)))
		i += copy(dAtA[i:], m.AccessLog)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintMysqlProxy(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *MySQLProxy) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.StatPrefix)
	if l > 0 {
		n += 1 + l + sovMysqlProxy(uint64(l))
	}
	l = len(m.AccessLog)
	if l > 0 {
		n += 1 + l + sovMysqlProxy(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMysqlProxy(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMysqlProxy(x uint64) (n int) {
	return sovMysqlProxy(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MySQLProxy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMysqlProxy
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
			return fmt.Errorf("proto: MySQLProxy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MySQLProxy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMysqlProxy
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
				return ErrInvalidLengthMysqlProxy
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatPrefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccessLog", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMysqlProxy
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
				return ErrInvalidLengthMysqlProxy
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccessLog = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMysqlProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMysqlProxy
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
func skipMysqlProxy(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMysqlProxy
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
					return 0, ErrIntOverflowMysqlProxy
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
					return 0, ErrIntOverflowMysqlProxy
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
				return 0, ErrInvalidLengthMysqlProxy
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMysqlProxy
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
				next, err := skipMysqlProxy(dAtA[start:])
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
	ErrInvalidLengthMysqlProxy = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMysqlProxy   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/filter/network/mysql_proxy/v1alpha1/mysql_proxy.proto", fileDescriptor_mysql_proxy_ea6ac70602b155ef)
}

var fileDescriptor_mysql_proxy_ea6ac70602b155ef = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x72, 0x4a, 0xcd, 0x2b, 0xcb,
	0xaf, 0xd4, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x4f, 0xcb, 0xcc, 0x29, 0x49, 0x2d, 0xd2,
	0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0xcf, 0xad, 0x2c, 0x2e, 0xcc, 0x89, 0x2f, 0x28,
	0xca, 0xaf, 0xa8, 0xd4, 0x2f, 0x33, 0x4c, 0xcc, 0x29, 0xc8, 0x48, 0x34, 0x44, 0x16, 0xd4, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x32, 0x00, 0x9b, 0xa1, 0x07, 0x31, 0x43, 0x0f, 0x62, 0x86, 0x1e,
	0xd4, 0x0c, 0x3d, 0x64, 0xe5, 0x30, 0x33, 0xa4, 0xc4, 0xcb, 0x12, 0x73, 0x32, 0x53, 0x12, 0x4b,
	0x52, 0xf5, 0x61, 0x0c, 0x88, 0x51, 0x4a, 0xe1, 0x5c, 0x5c, 0xbe, 0x95, 0xc1, 0x81, 0x3e, 0x01,
	0x20, 0xf5, 0x42, 0x5a, 0x5c, 0xdc, 0xc5, 0x25, 0x89, 0x25, 0xf1, 0x05, 0x45, 0xa9, 0x69, 0x99,
	0x15, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x4e, 0x9c, 0xbb, 0x5e, 0x1e, 0x60, 0x66, 0x29, 0x62,
	0x52, 0x60, 0x0c, 0xe2, 0x02, 0xc9, 0x06, 0x80, 0x25, 0x85, 0x64, 0xb9, 0xb8, 0x12, 0x93, 0x93,
	0x53, 0x8b, 0x8b, 0xe3, 0x73, 0xf2, 0xd3, 0x25, 0x98, 0x40, 0x4a, 0x83, 0x38, 0x21, 0x22, 0x3e,
	0xf9, 0xe9, 0x4e, 0x11, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c,
	0x23, 0x97, 0x5d, 0x66, 0xbe, 0x1e, 0xd8, 0xd1, 0x10, 0x77, 0x91, 0xea, 0xfe, 0x00, 0xc6, 0x28,
	0x0e, 0x18, 0x3b, 0x89, 0x0d, 0xec, 0x72, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x99,
	0x1a, 0xb1, 0x4a, 0x01, 0x00, 0x00,
}
