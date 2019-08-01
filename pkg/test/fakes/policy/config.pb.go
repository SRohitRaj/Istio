// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pkg/test/fakes/policy/config.proto

// adapter config for policy backend.

package policy

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"
	types "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	math_bits "math/bits"
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

// Config for policy backend, which could be used as a fake adpater for integration test,
// supports checknothing and keyval template.
type Params struct {
	// Specify check related params.
	CheckParams *Params_CheckParams `protobuf:"bytes,1,opt,name=check_params,json=checkParams,proto3" json:"check_params,omitempty"`
	// Specify route directive related params.
	Table map[string]string `protobuf:"bytes,2,rep,name=table,proto3" json:"table,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_6befd09209a3cf70, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetCheckParams() *Params_CheckParams {
	if m != nil {
		return m.CheckParams
	}
	return nil
}

func (m *Params) GetTable() map[string]string {
	if m != nil {
		return m.Table
	}
	return nil
}

// Check params which controls check result returned by policy backend.
type Params_CheckParams struct {
	// Controls that request should be allowed or not.
	CheckAllow bool `protobuf:"varint,1,opt,name=check_allow,json=checkAllow,proto3" json:"check_allow,omitempty"`
	// Valid duration of the check result.
	ValidDuration *types.Duration `protobuf:"bytes,2,opt,name=valid_duration,json=validDuration,proto3" json:"valid_duration,omitempty"`
	// Valid request count of the check result.
	ValidCount int64 `protobuf:"varint,3,opt,name=valid_count,json=validCount,proto3" json:"valid_count,omitempty"`
}

func (m *Params_CheckParams) Reset()      { *m = Params_CheckParams{} }
func (*Params_CheckParams) ProtoMessage() {}
func (*Params_CheckParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_6befd09209a3cf70, []int{0, 0}
}
func (m *Params_CheckParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params_CheckParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params_CheckParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params_CheckParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params_CheckParams.Merge(m, src)
}
func (m *Params_CheckParams) XXX_Size() int {
	return m.Size()
}
func (m *Params_CheckParams) XXX_DiscardUnknown() {
	xxx_messageInfo_Params_CheckParams.DiscardUnknown(m)
}

var xxx_messageInfo_Params_CheckParams proto.InternalMessageInfo

func (m *Params_CheckParams) GetCheckAllow() bool {
	if m != nil {
		return m.CheckAllow
	}
	return false
}

func (m *Params_CheckParams) GetValidDuration() *types.Duration {
	if m != nil {
		return m.ValidDuration
	}
	return nil
}

func (m *Params_CheckParams) GetValidCount() int64 {
	if m != nil {
		return m.ValidCount
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "policy.Params")
	proto.RegisterMapType((map[string]string)(nil), "policy.Params.TableEntry")
	proto.RegisterType((*Params_CheckParams)(nil), "policy.Params.CheckParams")
}

func init() { proto.RegisterFile("pkg/test/fakes/policy/config.proto", fileDescriptor_6befd09209a3cf70) }

var fileDescriptor_6befd09209a3cf70 = []byte{
	// 354 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0x31, 0x4f, 0xc2, 0x40,
	0x18, 0xed, 0xb5, 0x81, 0xc8, 0x55, 0x8d, 0x69, 0x18, 0x4a, 0x87, 0x4f, 0xc2, 0xc4, 0xd4, 0x4b,
	0x70, 0x21, 0x46, 0x13, 0x15, 0xdd, 0x4d, 0xe3, 0xe4, 0x42, 0x8e, 0x52, 0x6a, 0xd3, 0xda, 0x6b,
	0xca, 0x15, 0xc3, 0xe6, 0x4f, 0xd0, 0x7f, 0xe1, 0x2f, 0x31, 0x8e, 0x8c, 0x8c, 0x72, 0x2c, 0x8e,
	0xfc, 0x04, 0x73, 0x77, 0x10, 0x8c, 0xdb, 0xf7, 0xde, 0xf7, 0xee, 0x7d, 0xef, 0xe5, 0x70, 0xa7,
	0x48, 0x63, 0xc2, 0xa3, 0x29, 0x27, 0x13, 0x9a, 0x46, 0x53, 0x52, 0xb0, 0x2c, 0x09, 0xe7, 0x24,
	0x64, 0xf9, 0x24, 0x89, 0xfd, 0xa2, 0x64, 0x9c, 0x39, 0x75, 0x4d, 0x7a, 0xcd, 0x98, 0xc5, 0x4c,
	0x51, 0x44, 0x4e, 0x7a, 0xeb, 0x41, 0xcc, 0x58, 0x9c, 0x45, 0x44, 0xa1, 0x51, 0x35, 0x21, 0xe3,
	0xaa, 0xa4, 0x3c, 0x61, 0xb9, 0xde, 0x77, 0x3e, 0x4d, 0x5c, 0xbf, 0xa7, 0x25, 0x7d, 0x9e, 0x3a,
	0x97, 0xf8, 0x30, 0x7c, 0x8a, 0xc2, 0x74, 0x58, 0x28, 0xec, 0xa2, 0x36, 0xea, 0xda, 0x3d, 0xcf,
	0xd7, 0xfe, 0xbe, 0x56, 0xf9, 0x03, 0x29, 0xd1, 0x73, 0x60, 0x87, 0x7b, 0xe0, 0x10, 0x5c, 0xe3,
	0x74, 0x94, 0x45, 0xae, 0xd9, 0xb6, 0xba, 0x76, 0xaf, 0xf5, 0xef, 0xdd, 0x83, 0xdc, 0xdd, 0xe5,
	0xbc, 0x9c, 0x07, 0x5a, 0xe7, 0xbd, 0x23, 0x6c, 0xff, 0x71, 0x73, 0x4e, 0xb1, 0xf6, 0x1b, 0xd2,
	0x2c, 0x63, 0x2f, 0xea, 0xfc, 0x41, 0x80, 0x15, 0x75, 0x2d, 0x19, 0xe7, 0x0a, 0x1f, 0xcf, 0x68,
	0x96, 0x8c, 0x87, 0xbb, 0x0e, 0xae, 0xa9, 0x22, 0xb6, 0x7c, 0x5d, 0xd2, 0xdf, 0x95, 0xf4, 0x6f,
	0xb7, 0x82, 0xe0, 0x48, 0x3d, 0xd8, 0x41, 0x79, 0x42, 0x3b, 0x84, 0xac, 0xca, 0xb9, 0x6b, 0xb5,
	0x51, 0xd7, 0x0a, 0xb0, 0xa2, 0x06, 0x92, 0xf1, 0xfa, 0x18, 0xef, 0x83, 0x3a, 0x27, 0xd8, 0x4a,
	0xa3, 0xb9, 0x4a, 0xd2, 0x08, 0xe4, 0xe8, 0x34, 0x71, 0x6d, 0x46, 0xb3, 0x2a, 0x52, 0x97, 0x1b,
	0x81, 0x06, 0xe7, 0x66, 0x1f, 0xdd, 0x5c, 0x2c, 0x56, 0x60, 0x2c, 0x57, 0x60, 0x6c, 0x56, 0x80,
	0x5e, 0x05, 0xa0, 0x0f, 0x01, 0xe8, 0x4b, 0x00, 0x5a, 0x08, 0x40, 0xdf, 0x02, 0xd0, 0x8f, 0x00,
	0x63, 0x23, 0x00, 0xbd, 0xad, 0xc1, 0x58, 0xac, 0xc1, 0x58, 0xae, 0xc1, 0x78, 0xdc, 0x7e, 0xde,
	0xa8, 0xae, 0xa2, 0x9f, 0xfd, 0x06, 0x00, 0x00, 0xff, 0xff, 0x76, 0xc6, 0x4c, 0xf6, 0xf1, 0x01,
	0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
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
	if !this.CheckParams.Equal(that1.CheckParams) {
		return false
	}
	if len(this.Table) != len(that1.Table) {
		return false
	}
	for i := range this.Table {
		if this.Table[i] != that1.Table[i] {
			return false
		}
	}
	return true
}
func (this *Params_CheckParams) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params_CheckParams)
	if !ok {
		that2, ok := that.(Params_CheckParams)
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
	if this.CheckAllow != that1.CheckAllow {
		return false
	}
	if !this.ValidDuration.Equal(that1.ValidDuration) {
		return false
	}
	if this.ValidCount != that1.ValidCount {
		return false
	}
	return true
}
func (this *Params) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&policy.Params{")
	if this.CheckParams != nil {
		s = append(s, "CheckParams: "+fmt.Sprintf("%#v", this.CheckParams)+",\n")
	}
	keysForTable := make([]string, 0, len(this.Table))
	for k, _ := range this.Table {
		keysForTable = append(keysForTable, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForTable)
	mapStringForTable := "map[string]string{"
	for _, k := range keysForTable {
		mapStringForTable += fmt.Sprintf("%#v: %#v,", k, this.Table[k])
	}
	mapStringForTable += "}"
	if this.Table != nil {
		s = append(s, "Table: "+mapStringForTable+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Params_CheckParams) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&policy.Params_CheckParams{")
	s = append(s, "CheckAllow: "+fmt.Sprintf("%#v", this.CheckAllow)+",\n")
	if this.ValidDuration != nil {
		s = append(s, "ValidDuration: "+fmt.Sprintf("%#v", this.ValidDuration)+",\n")
	}
	s = append(s, "ValidCount: "+fmt.Sprintf("%#v", this.ValidCount)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringConfig(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Table) > 0 {
		for k := range m.Table {
			v := m.Table[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintConfig(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintConfig(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintConfig(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.CheckParams != nil {
		{
			size, err := m.CheckParams.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintConfig(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Params_CheckParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params_CheckParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params_CheckParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ValidCount != 0 {
		i = encodeVarintConfig(dAtA, i, uint64(m.ValidCount))
		i--
		dAtA[i] = 0x18
	}
	if m.ValidDuration != nil {
		{
			size, err := m.ValidDuration.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintConfig(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.CheckAllow {
		i--
		if m.CheckAllow {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CheckParams != nil {
		l = m.CheckParams.Size()
		n += 1 + l + sovConfig(uint64(l))
	}
	if len(m.Table) > 0 {
		for k, v := range m.Table {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovConfig(uint64(len(k))) + 1 + len(v) + sovConfig(uint64(len(v)))
			n += mapEntrySize + 1 + sovConfig(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *Params_CheckParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CheckAllow {
		n += 2
	}
	if m.ValidDuration != nil {
		l = m.ValidDuration.Size()
		n += 1 + l + sovConfig(uint64(l))
	}
	if m.ValidCount != 0 {
		n += 1 + sovConfig(uint64(m.ValidCount))
	}
	return n
}

func sovConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozConfig(x uint64) (n int) {
	return sovConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Params) String() string {
	if this == nil {
		return "nil"
	}
	keysForTable := make([]string, 0, len(this.Table))
	for k, _ := range this.Table {
		keysForTable = append(keysForTable, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForTable)
	mapStringForTable := "map[string]string{"
	for _, k := range keysForTable {
		mapStringForTable += fmt.Sprintf("%v: %v,", k, this.Table[k])
	}
	mapStringForTable += "}"
	s := strings.Join([]string{`&Params{`,
		`CheckParams:` + strings.Replace(fmt.Sprintf("%v", this.CheckParams), "Params_CheckParams", "Params_CheckParams", 1) + `,`,
		`Table:` + mapStringForTable + `,`,
		`}`,
	}, "")
	return s
}
func (this *Params_CheckParams) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Params_CheckParams{`,
		`CheckAllow:` + fmt.Sprintf("%v", this.CheckAllow) + `,`,
		`ValidDuration:` + strings.Replace(fmt.Sprintf("%v", this.ValidDuration), "Duration", "types.Duration", 1) + `,`,
		`ValidCount:` + fmt.Sprintf("%v", this.ValidCount) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringConfig(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConfig
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CheckParams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CheckParams == nil {
				m.CheckParams = &Params_CheckParams{}
			}
			if err := m.CheckParams.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Table", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Table == nil {
				m.Table = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowConfig
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowConfig
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthConfig
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthConfig
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowConfig
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthConfig
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthConfig
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipConfig(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthConfig
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Table[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthConfig
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthConfig
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
func (m *Params_CheckParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowConfig
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
			return fmt.Errorf("proto: CheckParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CheckParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CheckAllow", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.CheckAllow = bool(v != 0)
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidDuration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ValidDuration == nil {
				m.ValidDuration = &types.Duration{}
			}
			if err := m.ValidDuration.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidCount", wireType)
			}
			m.ValidCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidCount |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthConfig
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthConfig
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
func skipConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowConfig
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
					return 0, ErrIntOverflowConfig
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
					return 0, ErrIntOverflowConfig
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
				return 0, ErrInvalidLengthConfig
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthConfig
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowConfig
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
				next, err := skipConfig(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthConfig
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
	ErrInvalidLengthConfig = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowConfig   = fmt.Errorf("proto: integer overflow")
)
