// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/network/mongo_proxy/v2/mongo_proxy.proto

package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import v2 "github.com/envoyproxy/go-control-plane/envoy/config/filter/fault/v2"
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

type MongoProxy struct {
	// The human readable prefix to use when emitting :ref:`statistics
	// <config_network_filters_mongo_proxy_stats>`.
	StatPrefix string `protobuf:"bytes,1,opt,name=stat_prefix,json=statPrefix,proto3" json:"stat_prefix,omitempty"`
	// The optional path to use for writing Mongo access logs. If not access log
	// path is specified no access logs will be written. Note that access log is
	// also gated :ref:`runtime <config_network_filters_mongo_proxy_runtime>`.
	AccessLog string `protobuf:"bytes,2,opt,name=access_log,json=accessLog,proto3" json:"access_log,omitempty"`
	// Inject a fixed delay before proxying a Mongo operation. Delays are
	// applied to the following MongoDB operations: Query, Insert, GetMore,
	// and KillCursors. Once an active delay is in progress, all incoming
	// data up until the timer event fires will be a part of the delay.
	Delay *v2.FaultDelay `protobuf:"bytes,3,opt,name=delay" json:"delay,omitempty"`
	// Flag to specify whether :ref:`dynamic metadata
	// <config_network_filters_mongo_proxy_dynamic_metadata>` should be emitted. Defaults to false.
	EmitDynamicMetadata  bool     `protobuf:"varint,4,opt,name=emit_dynamic_metadata,json=emitDynamicMetadata,proto3" json:"emit_dynamic_metadata,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MongoProxy) Reset()         { *m = MongoProxy{} }
func (m *MongoProxy) String() string { return proto.CompactTextString(m) }
func (*MongoProxy) ProtoMessage()    {}
func (*MongoProxy) Descriptor() ([]byte, []int) {
	return fileDescriptor_mongo_proxy_c181d629f222aba8, []int{0}
}
func (m *MongoProxy) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MongoProxy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MongoProxy.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *MongoProxy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MongoProxy.Merge(dst, src)
}
func (m *MongoProxy) XXX_Size() int {
	return m.Size()
}
func (m *MongoProxy) XXX_DiscardUnknown() {
	xxx_messageInfo_MongoProxy.DiscardUnknown(m)
}

var xxx_messageInfo_MongoProxy proto.InternalMessageInfo

func (m *MongoProxy) GetStatPrefix() string {
	if m != nil {
		return m.StatPrefix
	}
	return ""
}

func (m *MongoProxy) GetAccessLog() string {
	if m != nil {
		return m.AccessLog
	}
	return ""
}

func (m *MongoProxy) GetDelay() *v2.FaultDelay {
	if m != nil {
		return m.Delay
	}
	return nil
}

func (m *MongoProxy) GetEmitDynamicMetadata() bool {
	if m != nil {
		return m.EmitDynamicMetadata
	}
	return false
}

func init() {
	proto.RegisterType((*MongoProxy)(nil), "envoy.config.filter.network.mongo_proxy.v2.MongoProxy")
}
func (m *MongoProxy) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MongoProxy) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.StatPrefix) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintMongoProxy(dAtA, i, uint64(len(m.StatPrefix)))
		i += copy(dAtA[i:], m.StatPrefix)
	}
	if len(m.AccessLog) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintMongoProxy(dAtA, i, uint64(len(m.AccessLog)))
		i += copy(dAtA[i:], m.AccessLog)
	}
	if m.Delay != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintMongoProxy(dAtA, i, uint64(m.Delay.Size()))
		n1, err := m.Delay.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.EmitDynamicMetadata {
		dAtA[i] = 0x20
		i++
		if m.EmitDynamicMetadata {
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

func encodeVarintMongoProxy(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *MongoProxy) Size() (n int) {
	var l int
	_ = l
	l = len(m.StatPrefix)
	if l > 0 {
		n += 1 + l + sovMongoProxy(uint64(l))
	}
	l = len(m.AccessLog)
	if l > 0 {
		n += 1 + l + sovMongoProxy(uint64(l))
	}
	if m.Delay != nil {
		l = m.Delay.Size()
		n += 1 + l + sovMongoProxy(uint64(l))
	}
	if m.EmitDynamicMetadata {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovMongoProxy(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozMongoProxy(x uint64) (n int) {
	return sovMongoProxy(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MongoProxy) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMongoProxy
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
			return fmt.Errorf("proto: MongoProxy: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MongoProxy: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatPrefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMongoProxy
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
				return ErrInvalidLengthMongoProxy
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
					return ErrIntOverflowMongoProxy
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
				return ErrInvalidLengthMongoProxy
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccessLog = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Delay", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMongoProxy
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
				return ErrInvalidLengthMongoProxy
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Delay == nil {
				m.Delay = &v2.FaultDelay{}
			}
			if err := m.Delay.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EmitDynamicMetadata", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMongoProxy
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
			m.EmitDynamicMetadata = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipMongoProxy(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMongoProxy
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
func skipMongoProxy(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMongoProxy
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
					return 0, ErrIntOverflowMongoProxy
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
					return 0, ErrIntOverflowMongoProxy
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
				return 0, ErrInvalidLengthMongoProxy
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowMongoProxy
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
				next, err := skipMongoProxy(dAtA[start:])
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
	ErrInvalidLengthMongoProxy = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMongoProxy   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/filter/network/mongo_proxy/v2/mongo_proxy.proto", fileDescriptor_mongo_proxy_c181d629f222aba8)
}

var fileDescriptor_mongo_proxy_c181d629f222aba8 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0x31, 0x4b, 0xc4, 0x30,
	0x18, 0x86, 0xc9, 0xdd, 0x29, 0x36, 0xb7, 0x48, 0x45, 0x2c, 0x07, 0x96, 0xe2, 0x54, 0x6e, 0x48,
	0x20, 0xae, 0xe2, 0x70, 0x1c, 0x4e, 0x1e, 0x1c, 0x1d, 0x5d, 0x4a, 0x6c, 0xd3, 0x12, 0x6c, 0x9b,
	0x92, 0x7e, 0xc6, 0xeb, 0x5f, 0x73, 0x72, 0x14, 0x5c, 0xfc, 0x09, 0xd2, 0xcd, 0x7f, 0x21, 0x69,
	0x4e, 0xb8, 0xe1, 0xb6, 0x2f, 0x79, 0xde, 0xe7, 0x4d, 0x3e, 0x7c, 0x27, 0x1a, 0xa3, 0x7a, 0x9a,
	0xa9, 0xa6, 0x90, 0x25, 0x2d, 0x64, 0x05, 0x42, 0xd3, 0x46, 0xc0, 0x9b, 0xd2, 0x2f, 0xb4, 0x56,
	0x4d, 0xa9, 0xd2, 0x56, 0xab, 0x5d, 0x4f, 0x0d, 0x3b, 0x3c, 0x92, 0x56, 0x2b, 0x50, 0xfe, 0x72,
	0xb4, 0x89, 0xb3, 0x89, 0xb3, 0xc9, 0xde, 0x26, 0x87, 0x71, 0xc3, 0x16, 0xf1, 0xb1, 0x97, 0x0a,
	0xfe, 0x5a, 0x81, 0xed, 0x1e, 0x07, 0xd7, 0xba, 0xb8, 0x32, 0xbc, 0x92, 0x39, 0x07, 0x41, 0xff,
	0x07, 0x07, 0x6e, 0xbe, 0x10, 0xc6, 0x1b, 0xdb, 0xba, 0xb5, 0xa5, 0xfe, 0x12, 0xcf, 0x3b, 0xe0,
	0x90, 0xb6, 0x5a, 0x14, 0x72, 0x17, 0xa0, 0x08, 0xc5, 0xde, 0xca, 0x7b, 0xff, 0xfd, 0x98, 0xce,
	0xf4, 0x24, 0x42, 0x09, 0xb6, 0x74, 0x3b, 0x42, 0xff, 0x1a, 0x63, 0x9e, 0x65, 0xa2, 0xeb, 0xd2,
	0x4a, 0x95, 0xc1, 0xc4, 0x46, 0x13, 0xcf, 0xdd, 0x3c, 0xaa, 0xd2, 0xbf, 0xc7, 0x27, 0xb9, 0xa8,
	0x78, 0x1f, 0x4c, 0x23, 0x14, 0xcf, 0x59, 0x4c, 0x8e, 0x2d, 0xe6, 0xfe, 0x68, 0x18, 0x79, 0xb0,
	0xc3, 0xda, 0xe6, 0x13, 0xa7, 0xf9, 0x0c, 0x5f, 0x8a, 0x5a, 0x42, 0x9a, 0xf7, 0x0d, 0xaf, 0x65,
	0x96, 0xd6, 0x02, 0x78, 0xce, 0x81, 0x07, 0xb3, 0x08, 0xc5, 0x67, 0xc9, 0x85, 0x85, 0x6b, 0xc7,
	0x36, 0x7b, 0xb4, 0x3a, 0xff, 0x1c, 0x42, 0xf4, 0x3d, 0x84, 0xe8, 0x67, 0x08, 0xd1, 0xd3, 0xc4,
	0xb0, 0xe7, 0xd3, 0x71, 0xcd, 0xdb, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xda, 0x13, 0x09, 0xf8,
	0x95, 0x01, 0x00, 0x00,
}
