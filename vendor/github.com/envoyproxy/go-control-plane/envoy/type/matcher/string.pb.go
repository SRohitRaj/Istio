// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/type/matcher/string.proto

package matcher

import (
	fmt "fmt"
	io "io"
	math "math"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/gogo/protobuf/proto"
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

// Specifies the way to match a string.
type StringMatcher struct {
	// Types that are valid to be assigned to MatchPattern:
	//	*StringMatcher_Exact
	//	*StringMatcher_Prefix
	//	*StringMatcher_Suffix
	//	*StringMatcher_Regex
	MatchPattern         isStringMatcher_MatchPattern `protobuf_oneof:"match_pattern"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *StringMatcher) Reset()         { *m = StringMatcher{} }
func (m *StringMatcher) String() string { return proto.CompactTextString(m) }
func (*StringMatcher) ProtoMessage()    {}
func (*StringMatcher) Descriptor() ([]byte, []int) {
	return fileDescriptor_1dc62c75a0f154e3, []int{0}
}
func (m *StringMatcher) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StringMatcher) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StringMatcher.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StringMatcher) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringMatcher.Merge(m, src)
}
func (m *StringMatcher) XXX_Size() int {
	return m.Size()
}
func (m *StringMatcher) XXX_DiscardUnknown() {
	xxx_messageInfo_StringMatcher.DiscardUnknown(m)
}

var xxx_messageInfo_StringMatcher proto.InternalMessageInfo

type isStringMatcher_MatchPattern interface {
	isStringMatcher_MatchPattern()
	MarshalTo([]byte) (int, error)
	Size() int
}

type StringMatcher_Exact struct {
	Exact string `protobuf:"bytes,1,opt,name=exact,proto3,oneof"`
}
type StringMatcher_Prefix struct {
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3,oneof"`
}
type StringMatcher_Suffix struct {
	Suffix string `protobuf:"bytes,3,opt,name=suffix,proto3,oneof"`
}
type StringMatcher_Regex struct {
	Regex string `protobuf:"bytes,4,opt,name=regex,proto3,oneof"`
}

func (*StringMatcher_Exact) isStringMatcher_MatchPattern()  {}
func (*StringMatcher_Prefix) isStringMatcher_MatchPattern() {}
func (*StringMatcher_Suffix) isStringMatcher_MatchPattern() {}
func (*StringMatcher_Regex) isStringMatcher_MatchPattern()  {}

func (m *StringMatcher) GetMatchPattern() isStringMatcher_MatchPattern {
	if m != nil {
		return m.MatchPattern
	}
	return nil
}

func (m *StringMatcher) GetExact() string {
	if x, ok := m.GetMatchPattern().(*StringMatcher_Exact); ok {
		return x.Exact
	}
	return ""
}

func (m *StringMatcher) GetPrefix() string {
	if x, ok := m.GetMatchPattern().(*StringMatcher_Prefix); ok {
		return x.Prefix
	}
	return ""
}

func (m *StringMatcher) GetSuffix() string {
	if x, ok := m.GetMatchPattern().(*StringMatcher_Suffix); ok {
		return x.Suffix
	}
	return ""
}

func (m *StringMatcher) GetRegex() string {
	if x, ok := m.GetMatchPattern().(*StringMatcher_Regex); ok {
		return x.Regex
	}
	return ""
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*StringMatcher) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _StringMatcher_OneofMarshaler, _StringMatcher_OneofUnmarshaler, _StringMatcher_OneofSizer, []interface{}{
		(*StringMatcher_Exact)(nil),
		(*StringMatcher_Prefix)(nil),
		(*StringMatcher_Suffix)(nil),
		(*StringMatcher_Regex)(nil),
	}
}

func _StringMatcher_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*StringMatcher)
	// match_pattern
	switch x := m.MatchPattern.(type) {
	case *StringMatcher_Exact:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.Exact)
	case *StringMatcher_Prefix:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.Prefix)
	case *StringMatcher_Suffix:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.Suffix)
	case *StringMatcher_Regex:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.Regex)
	case nil:
	default:
		return fmt.Errorf("StringMatcher.MatchPattern has unexpected type %T", x)
	}
	return nil
}

func _StringMatcher_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*StringMatcher)
	switch tag {
	case 1: // match_pattern.exact
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.MatchPattern = &StringMatcher_Exact{x}
		return true, err
	case 2: // match_pattern.prefix
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.MatchPattern = &StringMatcher_Prefix{x}
		return true, err
	case 3: // match_pattern.suffix
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.MatchPattern = &StringMatcher_Suffix{x}
		return true, err
	case 4: // match_pattern.regex
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.MatchPattern = &StringMatcher_Regex{x}
		return true, err
	default:
		return false, nil
	}
}

func _StringMatcher_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*StringMatcher)
	// match_pattern
	switch x := m.MatchPattern.(type) {
	case *StringMatcher_Exact:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Exact)))
		n += len(x.Exact)
	case *StringMatcher_Prefix:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Prefix)))
		n += len(x.Prefix)
	case *StringMatcher_Suffix:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Suffix)))
		n += len(x.Suffix)
	case *StringMatcher_Regex:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Regex)))
		n += len(x.Regex)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Specifies a list of ways to match a string.
type ListStringMatcher struct {
	Patterns             []*StringMatcher `protobuf:"bytes,1,rep,name=patterns,proto3" json:"patterns,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ListStringMatcher) Reset()         { *m = ListStringMatcher{} }
func (m *ListStringMatcher) String() string { return proto.CompactTextString(m) }
func (*ListStringMatcher) ProtoMessage()    {}
func (*ListStringMatcher) Descriptor() ([]byte, []int) {
	return fileDescriptor_1dc62c75a0f154e3, []int{1}
}
func (m *ListStringMatcher) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ListStringMatcher) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ListStringMatcher.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ListStringMatcher) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListStringMatcher.Merge(m, src)
}
func (m *ListStringMatcher) XXX_Size() int {
	return m.Size()
}
func (m *ListStringMatcher) XXX_DiscardUnknown() {
	xxx_messageInfo_ListStringMatcher.DiscardUnknown(m)
}

var xxx_messageInfo_ListStringMatcher proto.InternalMessageInfo

func (m *ListStringMatcher) GetPatterns() []*StringMatcher {
	if m != nil {
		return m.Patterns
	}
	return nil
}

func init() {
	proto.RegisterType((*StringMatcher)(nil), "envoy.type.matcher.StringMatcher")
	proto.RegisterType((*ListStringMatcher)(nil), "envoy.type.matcher.ListStringMatcher")
}

func init() { proto.RegisterFile("envoy/type/matcher/string.proto", fileDescriptor_1dc62c75a0f154e3) }

var fileDescriptor_1dc62c75a0f154e3 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0xcd, 0x2b, 0xcb,
	0xaf, 0xd4, 0x2f, 0xa9, 0x2c, 0x48, 0xd5, 0xcf, 0x4d, 0x2c, 0x49, 0xce, 0x48, 0x2d, 0xd2, 0x2f,
	0x2e, 0x29, 0xca, 0xcc, 0x4b, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x02, 0x2b, 0xd0,
	0x03, 0x29, 0xd0, 0x83, 0x2a, 0x90, 0x12, 0x2f, 0x4b, 0xcc, 0xc9, 0x4c, 0x49, 0x2c, 0x49, 0xd5,
	0x87, 0x31, 0x20, 0x8a, 0x95, 0xd6, 0x32, 0x72, 0xf1, 0x06, 0x83, 0x75, 0xfb, 0x42, 0x94, 0x0a,
	0x89, 0x71, 0xb1, 0xa6, 0x56, 0x24, 0x26, 0x97, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x7a, 0x30,
	0x04, 0x41, 0xb8, 0x42, 0xca, 0x5c, 0x6c, 0x05, 0x45, 0xa9, 0x69, 0x99, 0x15, 0x12, 0x4c, 0x20,
	0x09, 0x27, 0xce, 0x5d, 0x2f, 0x0f, 0x30, 0xb3, 0x14, 0x31, 0x29, 0x30, 0x7a, 0x30, 0x04, 0x41,
	0xa5, 0x40, 0x8a, 0x8a, 0x4b, 0xd3, 0x40, 0x8a, 0x98, 0xb1, 0x28, 0x82, 0x48, 0x09, 0x29, 0x71,
	0xb1, 0x16, 0xa5, 0xa6, 0xa7, 0x56, 0x48, 0xb0, 0x80, 0xd5, 0x70, 0x81, 0xd4, 0xb0, 0x16, 0x31,
	0x6b, 0x34, 0x70, 0x80, 0x6c, 0x03, 0x4b, 0x39, 0x89, 0x71, 0xf1, 0x82, 0xdd, 0x1e, 0x5f, 0x90,
	0x58, 0x52, 0x92, 0x5a, 0x94, 0x27, 0xc4, 0xba, 0xe3, 0xe5, 0x01, 0x66, 0x46, 0xa5, 0x38, 0x2e,
	0x41, 0x9f, 0xcc, 0xe2, 0x12, 0x54, 0x27, 0x7b, 0x72, 0x71, 0x40, 0x95, 0x15, 0x4b, 0x30, 0x2a,
	0x30, 0x6b, 0x70, 0x1b, 0x29, 0xea, 0x61, 0x06, 0x82, 0x1e, 0x8a, 0x26, 0xa8, 0xb5, 0x93, 0x18,
	0x99, 0x38, 0x18, 0x83, 0xe0, 0xda, 0x9d, 0xdc, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e,
	0xf1, 0xc1, 0x23, 0x39, 0x46, 0x2e, 0x85, 0xcc, 0x7c, 0x88, 0x41, 0x05, 0x45, 0xf9, 0x15, 0x95,
	0x58, 0xcc, 0x74, 0xe2, 0x86, 0x18, 0x1a, 0x00, 0x0a, 0xcc, 0x00, 0xc6, 0x28, 0x76, 0xa8, 0x78,
	0x12, 0x1b, 0x38, 0x78, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x40, 0x1a, 0x1f, 0xf7, 0xae,
	0x01, 0x00, 0x00,
}

func (m *StringMatcher) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StringMatcher) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.MatchPattern != nil {
		nn1, err := m.MatchPattern.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn1
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *StringMatcher_Exact) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0xa
	i++
	i = encodeVarintString(dAtA, i, uint64(len(m.Exact)))
	i += copy(dAtA[i:], m.Exact)
	return i, nil
}
func (m *StringMatcher_Prefix) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0x12
	i++
	i = encodeVarintString(dAtA, i, uint64(len(m.Prefix)))
	i += copy(dAtA[i:], m.Prefix)
	return i, nil
}
func (m *StringMatcher_Suffix) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0x1a
	i++
	i = encodeVarintString(dAtA, i, uint64(len(m.Suffix)))
	i += copy(dAtA[i:], m.Suffix)
	return i, nil
}
func (m *StringMatcher_Regex) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	dAtA[i] = 0x22
	i++
	i = encodeVarintString(dAtA, i, uint64(len(m.Regex)))
	i += copy(dAtA[i:], m.Regex)
	return i, nil
}
func (m *ListStringMatcher) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ListStringMatcher) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Patterns) > 0 {
		for _, msg := range m.Patterns {
			dAtA[i] = 0xa
			i++
			i = encodeVarintString(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintString(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *StringMatcher) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MatchPattern != nil {
		n += m.MatchPattern.Size()
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *StringMatcher_Exact) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Exact)
	n += 1 + l + sovString(uint64(l))
	return n
}
func (m *StringMatcher_Prefix) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Prefix)
	n += 1 + l + sovString(uint64(l))
	return n
}
func (m *StringMatcher_Suffix) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Suffix)
	n += 1 + l + sovString(uint64(l))
	return n
}
func (m *StringMatcher_Regex) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Regex)
	n += 1 + l + sovString(uint64(l))
	return n
}
func (m *ListStringMatcher) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Patterns) > 0 {
		for _, e := range m.Patterns {
			l = e.Size()
			n += 1 + l + sovString(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovString(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozString(x uint64) (n int) {
	return sovString(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StringMatcher) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowString
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
			return fmt.Errorf("proto: StringMatcher: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StringMatcher: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exact", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowString
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
				return ErrInvalidLengthString
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthString
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MatchPattern = &StringMatcher_Exact{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowString
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
				return ErrInvalidLengthString
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthString
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MatchPattern = &StringMatcher_Prefix{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Suffix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowString
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
				return ErrInvalidLengthString
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthString
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MatchPattern = &StringMatcher_Suffix{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Regex", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowString
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
				return ErrInvalidLengthString
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthString
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MatchPattern = &StringMatcher_Regex{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipString(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthString
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthString
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
func (m *ListStringMatcher) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowString
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
			return fmt.Errorf("proto: ListStringMatcher: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ListStringMatcher: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Patterns", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowString
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
				return ErrInvalidLengthString
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthString
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Patterns = append(m.Patterns, &StringMatcher{})
			if err := m.Patterns[len(m.Patterns)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipString(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthString
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthString
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
func skipString(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowString
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
					return 0, ErrIntOverflowString
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
					return 0, ErrIntOverflowString
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
				return 0, ErrInvalidLengthString
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthString
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowString
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
				next, err := skipString(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthString
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
	ErrInvalidLengthString = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowString   = fmt.Errorf("proto: integer overflow")
)
