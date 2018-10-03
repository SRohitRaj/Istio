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
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x31, 0x4e, 0xc3, 0x30,
	0x14, 0x86, 0x6d, 0x8a, 0x90, 0x30, 0xb4, 0x40, 0x44, 0xa5, 0xd0, 0xe1, 0x11, 0x15, 0x86, 0x88,
	0x21, 0x51, 0xe1, 0x04, 0xb4, 0x5b, 0xc7, 0x2c, 0x48, 0x2c, 0x55, 0x42, 0x4c, 0x64, 0x29, 0xed,
	0x0b, 0xb6, 0x53, 0x18, 0x39, 0x02, 0xc7, 0x60, 0xe0, 0x20, 0x19, 0x3b, 0x32, 0x21, 0x62, 0x16,
	0xc6, 0x8e, 0x8c, 0x28, 0x09, 0x2d, 0xb0, 0x59, 0xff, 0xf7, 0x3d, 0xdb, 0xef, 0x67, 0xdd, 0xa9,
	0x78, 0xe0, 0xd2, 0x9f, 0x0f, 0x7c, 0xc9, 0x33, 0x94, 0xda, 0xcb, 0x24, 0x6a, 0xb4, 0x3a, 0x42,
	0x69, 0x81, 0x5e, 0x0d, 0xbd, 0xf9, 0xa0, 0x77, 0x98, 0x60, 0x82, 0x35, 0xf2, 0xab, 0x53, 0x63,
	0xf5, 0x8e, 0xd6, 0xc3, 0xa1, 0xd6, 0x52, 0x44, 0xb9, 0xe6, 0xaa, 0x41, 0xfd, 0x17, 0xca, 0xda,
	0x41, 0x7d, 0x63, 0xc0, 0xef, 0x72, 0xae, 0xb4, 0x35, 0x66, 0xec, 0xd7, 0xb2, 0xa9, 0xd3, 0x72,
	0x77, 0xce, 0x4f, 0xbd, 0xff, 0xef, 0x78, 0x23, 0x9c, 0x66, 0x92, 0x2b, 0xc5, 0xe3, 0xcb, 0xb5,
	0x3b, 0xdc, 0x2c, 0xde, 0x8e, 0x49, 0xf0, 0x67, 0xda, 0x3a, 0x61, 0xed, 0x98, 0xdf, 0x86, 0x79,
	0xaa, 0x27, 0xf7, 0x28, 0x63, 0x65, 0x6f, 0x38, 0x2d, 0x77, 0x3b, 0xd8, 0xfd, 0x09, 0xaf, 0xaa,
	0xcc, 0x3a, 0x63, 0x07, 0x49, 0x8a, 0x51, 0x98, 0xd6, 0xce, 0xe4, 0x06, 0xf3, 0x99, 0xb6, 0x5b,
	0x0e, 0x75, 0xdb, 0xc1, 0x5e, 0x03, 0x2a, 0x6f, 0x54, 0xc5, 0xfd, 0x7d, 0xd6, 0x59, 0xfd, 0x56,
	0x65, 0x38, 0x53, 0x7c, 0x38, 0x2e, 0x4a, 0x20, 0x8b, 0x12, 0xc8, 0x6b, 0x09, 0x64, 0x59, 0x02,
	0x79, 0x34, 0x40, 0x9f, 0x0d, 0x90, 0xc2, 0x00, 0x5d, 0x18, 0xa0, 0xef, 0x06, 0xe8, 0xa7, 0x01,
	0xb2, 0x34, 0x40, 0x9f, 0x3e, 0x80, 0x5c, 0x77, 0x9b, 0x5d, 0x04, 0xfa, 0x61, 0x26, 0xfc, 0x55,
	0x35, 0x5f, 0x94, 0x46, 0x5b, 0x75, 0x27, 0x17, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x09, 0xd3,
	0xb7, 0xd0, 0x6d, 0x01, 0x00, 0x00,
}
