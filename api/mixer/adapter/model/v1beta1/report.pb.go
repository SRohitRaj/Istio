// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mixer/adapter/model/v1beta1/report.proto

package v1beta1

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

// Expresses the result of a report call.
type ReportResult struct {
}

func (m *ReportResult) Reset()                    { *m = ReportResult{} }
func (*ReportResult) ProtoMessage()               {}
func (*ReportResult) Descriptor() ([]byte, []int) { return fileDescriptorReport, []int{0} }

func init() {
	proto.RegisterType((*ReportResult)(nil), "istio.mixer.adapter.model.v1beta1.ReportResult")
}
func (m *ReportResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReportResult) MarshalTo(dAtA []byte) (int, error) {
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
func (m *ReportResult) Size() (n int) {
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
func (this *ReportResult) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ReportResult{`,
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
func (m *ReportResult) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ReportResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReportResult: illegal tag %d (wire type %d)", fieldNum, wire)
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

func init() { proto.RegisterFile("mixer/adapter/model/v1beta1/report.proto", fileDescriptorReport) }

var fileDescriptorReport = []byte{
	// 185 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0xc8, 0xcd, 0xac, 0x48,
	0x2d, 0xd2, 0x4f, 0x4c, 0x49, 0x2c, 0x28, 0x49, 0x2d, 0xd2, 0xcf, 0xcd, 0x4f, 0x49, 0xcd, 0xd1,
	0x2f, 0x33, 0x4c, 0x4a, 0x2d, 0x49, 0x34, 0xd4, 0x2f, 0x4a, 0x2d, 0xc8, 0x2f, 0x2a, 0xd1, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x52, 0xcc, 0x2c, 0x2e, 0xc9, 0xcc, 0xd7, 0x03, 0xab, 0xd7, 0x83,
	0xaa, 0xd7, 0x03, 0xab, 0xd7, 0x83, 0xaa, 0x97, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0xab, 0xd6,
	0x07, 0xb1, 0x20, 0x1a, 0x95, 0xf8, 0xb8, 0x78, 0x82, 0xc0, 0x06, 0x05, 0xa5, 0x16, 0x97, 0xe6,
	0x94, 0x38, 0x25, 0x9c, 0x78, 0x28, 0xc7, 0x70, 0xe1, 0xa1, 0x1c, 0xc3, 0x8d, 0x87, 0x72, 0x0c,
	0x1f, 0x1e, 0xca, 0x31, 0x34, 0x3c, 0x92, 0x63, 0x5c, 0xf1, 0x48, 0x8e, 0xe1, 0xc4, 0x23, 0x39,
	0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x7c, 0xf1, 0x48, 0x8e, 0xe1, 0xc3, 0x23,
	0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0xa2, 0xf4, 0x20, 0x56, 0x67, 0xe6, 0xeb, 0x83, 0x19, 0xfa,
	0x89, 0x05, 0x99, 0xfa, 0x78, 0xdc, 0x9d, 0xc4, 0x06, 0xb6, 0xd8, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0x7d, 0x68, 0xc7, 0x30, 0xdd, 0x00, 0x00, 0x00,
}
