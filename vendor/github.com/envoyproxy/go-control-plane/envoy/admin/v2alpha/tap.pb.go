// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/admin/v2alpha/tap.proto

package envoy_admin_v2alpha

import (
	fmt "fmt"
	v2alpha "github.com/envoyproxy/go-control-plane/envoy/service/tap/v2alpha"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/gogo/protobuf/proto"
	io "io"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// The /tap admin request body that is used to configure an active tap session.
type TapRequest struct {
	// The opaque configuration ID used to match the configuration to a loaded extension.
	// A tap extension configures a similar opaque ID that is used to match.
	ConfigId string `protobuf:"bytes,1,opt,name=config_id,json=configId,proto3" json:"config_id,omitempty"`
	// The tap configuration to load.
	TapConfig            *v2alpha.TapConfig `protobuf:"bytes,2,opt,name=tap_config,json=tapConfig,proto3" json:"tap_config,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *TapRequest) Reset()         { *m = TapRequest{} }
func (m *TapRequest) String() string { return proto.CompactTextString(m) }
func (*TapRequest) ProtoMessage()    {}
func (*TapRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_69ef50f493e7b8aa, []int{0}
}
func (m *TapRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TapRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TapRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TapRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TapRequest.Merge(m, src)
}
func (m *TapRequest) XXX_Size() int {
	return m.Size()
}
func (m *TapRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TapRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TapRequest proto.InternalMessageInfo

func (m *TapRequest) GetConfigId() string {
	if m != nil {
		return m.ConfigId
	}
	return ""
}

func (m *TapRequest) GetTapConfig() *v2alpha.TapConfig {
	if m != nil {
		return m.TapConfig
	}
	return nil
}

func init() {
	proto.RegisterType((*TapRequest)(nil), "envoy.admin.v2alpha.TapRequest")
}

func init() { proto.RegisterFile("envoy/admin/v2alpha/tap.proto", fileDescriptor_69ef50f493e7b8aa) }

var fileDescriptor_69ef50f493e7b8aa = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4d, 0xcd, 0x2b, 0xcb,
	0xaf, 0xd4, 0x4f, 0x4c, 0xc9, 0xcd, 0xcc, 0xd3, 0x2f, 0x33, 0x4a, 0xcc, 0x29, 0xc8, 0x48, 0xd4,
	0x2f, 0x49, 0x2c, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x06, 0x4b, 0xeb, 0x81, 0xa5,
	0xf5, 0xa0, 0xd2, 0x52, 0x6a, 0x10, 0x3d, 0xc5, 0xa9, 0x45, 0x65, 0x99, 0xc9, 0xa9, 0x20, 0xd5,
	0x70, 0x9d, 0xc9, 0xf9, 0xb9, 0xb9, 0xf9, 0x79, 0x10, 0xcd, 0x52, 0xe2, 0x65, 0x89, 0x39, 0x99,
	0x29, 0x89, 0x25, 0xa9, 0xfa, 0x30, 0x06, 0x44, 0x42, 0xa9, 0x95, 0x91, 0x8b, 0x2b, 0x24, 0xb1,
	0x20, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x48, 0x8d, 0x8b, 0x33, 0x39, 0x3f, 0x2f, 0x2d,
	0x33, 0x3d, 0x3e, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0xd3, 0x89, 0x73, 0xd7, 0xcb, 0x03,
	0xcc, 0x2c, 0x45, 0x4c, 0x0a, 0x8c, 0x41, 0x1c, 0x10, 0x39, 0xcf, 0x14, 0x21, 0x7f, 0x2e, 0xae,
	0x92, 0xc4, 0x82, 0x78, 0x08, 0x5f, 0x82, 0x49, 0x81, 0x51, 0x83, 0xdb, 0x48, 0x45, 0x0f, 0xe2,
	0x42, 0xa8, 0x63, 0xf4, 0x40, 0x4e, 0x87, 0x3a, 0x46, 0x2f, 0x24, 0xb1, 0xc0, 0x19, 0xac, 0xd6,
	0x89, 0x0b, 0x64, 0x1c, 0x6b, 0x17, 0x23, 0x93, 0x00, 0x63, 0x10, 0x67, 0x09, 0x5c, 0xd8, 0xfa,
	0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0xe4, 0x52, 0xcc, 0xcc,
	0x87, 0x18, 0x56, 0x50, 0x94, 0x5f, 0x51, 0xa9, 0x87, 0xc5, 0xe7, 0x4e, 0x1c, 0x21, 0x89, 0x05,
	0x01, 0x20, 0x2f, 0x04, 0x30, 0x26, 0xb1, 0x81, 0xfd, 0x62, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff,
	0x40, 0xb6, 0x00, 0xf2, 0x42, 0x01, 0x00, 0x00,
}

func (m *TapRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TapRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ConfigId) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintTap(dAtA, i, uint64(len(m.ConfigId)))
		i += copy(dAtA[i:], m.ConfigId)
	}
	if m.TapConfig != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintTap(dAtA, i, uint64(m.TapConfig.Size()))
		n1, err := m.TapConfig.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintTap(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *TapRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ConfigId)
	if l > 0 {
		n += 1 + l + sovTap(uint64(l))
	}
	if m.TapConfig != nil {
		l = m.TapConfig.Size()
		n += 1 + l + sovTap(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTap(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozTap(x uint64) (n int) {
	return sovTap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TapRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTap
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
			return fmt.Errorf("proto: TapRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TapRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConfigId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTap
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
				return ErrInvalidLengthTap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConfigId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TapConfig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTap
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.TapConfig == nil {
				m.TapConfig = &v2alpha.TapConfig{}
			}
			if err := m.TapConfig.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTap
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
func skipTap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTap
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
					return 0, ErrIntOverflowTap
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
					return 0, ErrIntOverflowTap
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
				return 0, ErrInvalidLengthTap
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthTap
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTap
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
				next, err := skipTap(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthTap
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
	ErrInvalidLengthTap = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTap   = fmt.Errorf("proto: integer overflow")
)
