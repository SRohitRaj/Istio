// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mixer/v1/report.proto

package v1

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Used to report telemetry after performing one or more actions.
type ReportRequest struct {
	// The attributes to use for this request.
	//
	// Each `Attributes` element represents the state of a single action. Multiple actions
	// can be provided in a single message in order to improve communication efficiency. The
	// client can accumulate a set of actions and send them all in one single message.
	//
	// Although each `Attributes` message is semantically treated as an independent
	// stand-alone entity unrelated to the other attributes within the message, this
	// message format leverages delta-encoding between attribute messages in order to
	// substantially reduce the request size and improve end-to-end efficiency. Each
	// individual set of attributes is used to modify the previous set. This eliminates
	// the need to redundantly send the same attributes multiple times over within
	// a single request.
	//
	// If a client is not sophisticated and doesn't want to use delta-encoding,
	// a degenerate case is to include all attributes in every individual message.
	Attributes []CompressedAttributes `protobuf:"bytes,1,rep,name=attributes" json:"attributes"`
	// The default message-level dictionary for all the attributes.
	// Individual attribute messages can have their own dictionaries, but if they don't
	// then this set of words, if it is provided, is used instead.
	//
	// This makes it possible to share the same dictionary for all attributes in this
	// request, which can substantially reduce the overall request size.
	DefaultWords []string `protobuf:"bytes,2,rep,name=default_words,json=defaultWords" json:"default_words,omitempty"`
	// The number of words in the global dictionary.
	// To detect global dictionary out of sync between client and server.
	GlobalWordCount uint32 `protobuf:"varint,3,opt,name=global_word_count,json=globalWordCount,proto3" json:"global_word_count,omitempty"`
}

func (m *ReportRequest) Reset()                    { *m = ReportRequest{} }
func (*ReportRequest) ProtoMessage()               {}
func (*ReportRequest) Descriptor() ([]byte, []int) { return fileDescriptorReport, []int{0} }

// Used to carry responses to telemetry reports
type ReportResponse struct {
}

func (m *ReportResponse) Reset()                    { *m = ReportResponse{} }
func (*ReportResponse) ProtoMessage()               {}
func (*ReportResponse) Descriptor() ([]byte, []int) { return fileDescriptorReport, []int{1} }

func init() {
	proto.RegisterType((*ReportRequest)(nil), "istio.mixer.v1.ReportRequest")
	proto.RegisterType((*ReportResponse)(nil), "istio.mixer.v1.ReportResponse")
}
func (m *ReportRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReportRequest) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Attributes) > 0 {
		for _, msg := range m.Attributes {
			dAtA[i] = 0xa
			i++
			i = encodeVarintReport(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.DefaultWords) > 0 {
		for _, s := range m.DefaultWords {
			dAtA[i] = 0x12
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
	if m.GlobalWordCount != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintReport(dAtA, i, uint64(m.GlobalWordCount))
	}
	return i, nil
}

func (m *ReportResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReportResponse) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	return i, nil
}

func encodeVarintReport(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *ReportRequest) Size() (n int) {
	var l int
	_ = l
	if len(m.Attributes) > 0 {
		for _, e := range m.Attributes {
			l = e.Size()
			n += 1 + l + sovReport(uint64(l))
		}
	}
	if len(m.DefaultWords) > 0 {
		for _, s := range m.DefaultWords {
			l = len(s)
			n += 1 + l + sovReport(uint64(l))
		}
	}
	if m.GlobalWordCount != 0 {
		n += 1 + sovReport(uint64(m.GlobalWordCount))
	}
	return n
}

func (m *ReportResponse) Size() (n int) {
	var l int
	_ = l
	return n
}

func sovReport(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozReport(x uint64) (n int) {
	return sovReport(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *ReportRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ReportRequest{`,
		`Attributes:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Attributes), "CompressedAttributes", "CompressedAttributes", 1), `&`, ``, 1) + `,`,
		`DefaultWords:` + fmt.Sprintf("%v", this.DefaultWords) + `,`,
		`GlobalWordCount:` + fmt.Sprintf("%v", this.GlobalWordCount) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ReportResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ReportResponse{`,
		`}`,
	}, "")
	return s
}
func valueToStringReport(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ReportRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReport
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
			return fmt.Errorf("proto: ReportRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReportRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attributes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
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
				return ErrInvalidLengthReport
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Attributes = append(m.Attributes, CompressedAttributes{})
			if err := m.Attributes[len(m.Attributes)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultWords", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
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
				return ErrInvalidLengthReport
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DefaultWords = append(m.DefaultWords, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GlobalWordCount", wireType)
			}
			m.GlobalWordCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReport
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GlobalWordCount |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipReport(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthReport
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
func (m *ReportResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReport
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
			return fmt.Errorf("proto: ReportResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReportResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipReport(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthReport
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
func skipReport(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowReport
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
					return 0, ErrIntOverflowReport
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
					return 0, ErrIntOverflowReport
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
				return 0, ErrInvalidLengthReport
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowReport
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
				next, err := skipReport(dAtA[start:])
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
	ErrInvalidLengthReport = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowReport   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("mixer/v1/report.proto", fileDescriptorReport) }

var fileDescriptorReport = []byte{
	// 294 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xcd, 0xcd, 0xac, 0x48,
	0x2d, 0xd2, 0x2f, 0x33, 0xd4, 0x2f, 0x4a, 0x2d, 0xc8, 0x2f, 0x2a, 0xd1, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0xe2, 0xcb, 0x2c, 0x2e, 0xc9, 0xcc, 0xd7, 0x03, 0x4b, 0xea, 0x95, 0x19, 0x4a, 0x89,
	0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0xa5, 0xf4, 0x41, 0x2c, 0x88, 0x2a, 0x29, 0x49, 0xb8, 0xe6, 0xc4,
	0x92, 0x92, 0xa2, 0xcc, 0xa4, 0xd2, 0x92, 0xd4, 0x62, 0x88, 0x94, 0xd2, 0x1a, 0x46, 0x2e, 0xde,
	0x20, 0xb0, 0x89, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x5e, 0x5c, 0x5c, 0x08, 0x55,
	0x12, 0x8c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0x2a, 0x7a, 0xa8, 0xf6, 0xe8, 0x39, 0xe7, 0xe7, 0x16,
	0x14, 0xa5, 0x16, 0x17, 0xa7, 0xa6, 0x38, 0xc2, 0xd5, 0x3a, 0xb1, 0x9c, 0xb8, 0x27, 0xcf, 0x10,
	0x84, 0xa4, 0x5b, 0x48, 0x99, 0x8b, 0x37, 0x25, 0x35, 0x2d, 0xb1, 0x34, 0xa7, 0x24, 0xbe, 0x3c,
	0xbf, 0x28, 0xa5, 0x58, 0x82, 0x49, 0x81, 0x59, 0x83, 0x33, 0x88, 0x07, 0x2a, 0x18, 0x0e, 0x12,
	0x13, 0xd2, 0xe2, 0x12, 0x4c, 0xcf, 0xc9, 0x4f, 0x4a, 0xcc, 0x01, 0xab, 0x89, 0x4f, 0xce, 0x2f,
	0xcd, 0x2b, 0x91, 0x60, 0x56, 0x60, 0xd4, 0xe0, 0x0d, 0xe2, 0x87, 0x48, 0x80, 0xd4, 0x39, 0x83,
	0x84, 0x95, 0x04, 0xb8, 0xf8, 0x60, 0xae, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x75, 0xf2, 0x3d,
	0xf1, 0x50, 0x8e, 0xe1, 0xc2, 0x43, 0x39, 0x86, 0x1b, 0x0f, 0xe5, 0x18, 0x3e, 0x3c, 0x94, 0x63,
	0x68, 0x78, 0x24, 0xc7, 0xb8, 0xe2, 0x91, 0x1c, 0xc3, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9,
	0x31, 0x3e, 0x78, 0x24, 0xc7, 0xf8, 0xe2, 0x91, 0x1c, 0xc3, 0x87, 0x47, 0x72, 0x8c, 0x13, 0x1e,
	0xcb, 0x31, 0x44, 0x49, 0x43, 0xfc, 0x92, 0x99, 0xaf, 0x0f, 0x66, 0xe8, 0x27, 0x16, 0x64, 0xea,
	0xc3, 0x02, 0x28, 0x89, 0x0d, 0x1c, 0x2c, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5a, 0x29,
	0x5a, 0xed, 0x70, 0x01, 0x00, 0x00,
}
