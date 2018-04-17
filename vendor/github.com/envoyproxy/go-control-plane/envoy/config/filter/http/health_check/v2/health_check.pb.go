// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/config/filter/http/health_check/v2/health_check.proto

/*
	Package v2 is a generated protocol buffer package.

	It is generated from these files:
		envoy/config/filter/http/health_check/v2/health_check.proto

	It has these top-level messages:
		HealthCheck
*/
package v2

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/types"
import google_protobuf1 "github.com/gogo/protobuf/types"
import envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
import _ "github.com/lyft/protoc-gen-validate/validate"
import _ "github.com/gogo/protobuf/gogoproto"

import time "time"

import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type HealthCheck struct {
	// Specifies whether the filter operates in pass through mode or not.
	PassThroughMode *google_protobuf1.BoolValue `protobuf:"bytes,1,opt,name=pass_through_mode,json=passThroughMode" json:"pass_through_mode,omitempty"`
	// Specifies the incoming HTTP endpoint that should be considered the
	// health check endpoint. For example */healthcheck*.
	Endpoint string `protobuf:"bytes,2,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// If operating in pass through mode, the amount of time in milliseconds
	// that the filter should cache the upstream response.
	CacheTime *time.Duration `protobuf:"bytes,3,opt,name=cache_time,json=cacheTime,stdduration" json:"cache_time,omitempty"`
	// If operating in non-pass-through mode, specifies a set of upstream cluster
	// names and the minimum percentage of servers in each of those clusters that
	// must be healthy in order for the filter to return a 200.
	ClusterMinHealthyPercentages map[string]*envoy_type.Percent `protobuf:"bytes,4,rep,name=cluster_min_healthy_percentages,json=clusterMinHealthyPercentages" json:"cluster_min_healthy_percentages,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *HealthCheck) Reset()                    { *m = HealthCheck{} }
func (m *HealthCheck) String() string            { return proto.CompactTextString(m) }
func (*HealthCheck) ProtoMessage()               {}
func (*HealthCheck) Descriptor() ([]byte, []int) { return fileDescriptorHealthCheck, []int{0} }

func (m *HealthCheck) GetPassThroughMode() *google_protobuf1.BoolValue {
	if m != nil {
		return m.PassThroughMode
	}
	return nil
}

func (m *HealthCheck) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *HealthCheck) GetCacheTime() *time.Duration {
	if m != nil {
		return m.CacheTime
	}
	return nil
}

func (m *HealthCheck) GetClusterMinHealthyPercentages() map[string]*envoy_type.Percent {
	if m != nil {
		return m.ClusterMinHealthyPercentages
	}
	return nil
}

func init() {
	proto.RegisterType((*HealthCheck)(nil), "envoy.config.filter.http.health_check.v2.HealthCheck")
}
func (m *HealthCheck) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HealthCheck) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.PassThroughMode != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintHealthCheck(dAtA, i, uint64(m.PassThroughMode.Size()))
		n1, err := m.PassThroughMode.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.Endpoint) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintHealthCheck(dAtA, i, uint64(len(m.Endpoint)))
		i += copy(dAtA[i:], m.Endpoint)
	}
	if m.CacheTime != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintHealthCheck(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdDuration(*m.CacheTime)))
		n2, err := github_com_gogo_protobuf_types.StdDurationMarshalTo(*m.CacheTime, dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if len(m.ClusterMinHealthyPercentages) > 0 {
		for k, _ := range m.ClusterMinHealthyPercentages {
			dAtA[i] = 0x22
			i++
			v := m.ClusterMinHealthyPercentages[k]
			msgSize := 0
			if v != nil {
				msgSize = v.Size()
				msgSize += 1 + sovHealthCheck(uint64(msgSize))
			}
			mapSize := 1 + len(k) + sovHealthCheck(uint64(len(k))) + msgSize
			i = encodeVarintHealthCheck(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintHealthCheck(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			if v != nil {
				dAtA[i] = 0x12
				i++
				i = encodeVarintHealthCheck(dAtA, i, uint64(v.Size()))
				n3, err := v.MarshalTo(dAtA[i:])
				if err != nil {
					return 0, err
				}
				i += n3
			}
		}
	}
	return i, nil
}

func encodeVarintHealthCheck(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *HealthCheck) Size() (n int) {
	var l int
	_ = l
	if m.PassThroughMode != nil {
		l = m.PassThroughMode.Size()
		n += 1 + l + sovHealthCheck(uint64(l))
	}
	l = len(m.Endpoint)
	if l > 0 {
		n += 1 + l + sovHealthCheck(uint64(l))
	}
	if m.CacheTime != nil {
		l = github_com_gogo_protobuf_types.SizeOfStdDuration(*m.CacheTime)
		n += 1 + l + sovHealthCheck(uint64(l))
	}
	if len(m.ClusterMinHealthyPercentages) > 0 {
		for k, v := range m.ClusterMinHealthyPercentages {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovHealthCheck(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovHealthCheck(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovHealthCheck(uint64(mapEntrySize))
		}
	}
	return n
}

func sovHealthCheck(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozHealthCheck(x uint64) (n int) {
	return sovHealthCheck(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *HealthCheck) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowHealthCheck
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
			return fmt.Errorf("proto: HealthCheck: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HealthCheck: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PassThroughMode", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheck
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
				return ErrInvalidLengthHealthCheck
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PassThroughMode == nil {
				m.PassThroughMode = &google_protobuf1.BoolValue{}
			}
			if err := m.PassThroughMode.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Endpoint", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheck
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
				return ErrInvalidLengthHealthCheck
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Endpoint = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CacheTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheck
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
				return ErrInvalidLengthHealthCheck
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.CacheTime == nil {
				m.CacheTime = new(time.Duration)
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(m.CacheTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterMinHealthyPercentages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowHealthCheck
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
				return ErrInvalidLengthHealthCheck
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ClusterMinHealthyPercentages == nil {
				m.ClusterMinHealthyPercentages = make(map[string]*envoy_type.Percent)
			}
			var mapkey string
			var mapvalue *envoy_type.Percent
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowHealthCheck
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowHealthCheck
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthHealthCheck
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowHealthCheck
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= (int(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthHealthCheck
					}
					postmsgIndex := iNdEx + mapmsglen
					if mapmsglen < 0 {
						return ErrInvalidLengthHealthCheck
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &envoy_type.Percent{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipHealthCheck(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthHealthCheck
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.ClusterMinHealthyPercentages[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipHealthCheck(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthHealthCheck
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
func skipHealthCheck(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowHealthCheck
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
					return 0, ErrIntOverflowHealthCheck
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
					return 0, ErrIntOverflowHealthCheck
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
				return 0, ErrInvalidLengthHealthCheck
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowHealthCheck
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
				next, err := skipHealthCheck(dAtA[start:])
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
	ErrInvalidLengthHealthCheck = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowHealthCheck   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("envoy/config/filter/http/health_check/v2/health_check.proto", fileDescriptorHealthCheck)
}

var fileDescriptorHealthCheck = []byte{
	// 435 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xc1, 0x8a, 0xd4, 0x30,
	0x1c, 0xc6, 0x49, 0x3b, 0x2b, 0x4e, 0xe6, 0xe0, 0x58, 0x05, 0xeb, 0x20, 0xb3, 0xa3, 0x20, 0x8c,
	0x97, 0x04, 0xea, 0x45, 0x14, 0x3c, 0x74, 0x15, 0xbc, 0x2c, 0x48, 0x59, 0x14, 0xbc, 0x94, 0x6c,
	0xfb, 0x9f, 0x36, 0x6c, 0x9b, 0x84, 0x34, 0xad, 0xf4, 0x15, 0x7c, 0x02, 0x4f, 0x3e, 0x88, 0x27,
	0x8f, 0xde, 0xf4, 0x0d, 0x94, 0xb9, 0xf9, 0x16, 0x92, 0xa4, 0xab, 0x2e, 0x0b, 0xea, 0xed, 0x9f,
	0x7c, 0xff, 0xef, 0xcb, 0xc7, 0x2f, 0xf8, 0x09, 0x88, 0x41, 0x8e, 0xb4, 0x90, 0x62, 0xc7, 0x2b,
	0xba, 0xe3, 0x8d, 0x01, 0x4d, 0x6b, 0x63, 0x14, 0xad, 0x81, 0x35, 0xa6, 0xce, 0x8b, 0x1a, 0x8a,
	0x33, 0x3a, 0x24, 0x17, 0xce, 0x44, 0x69, 0x69, 0x64, 0xb4, 0x75, 0x66, 0xe2, 0xcd, 0xc4, 0x9b,
	0x89, 0x35, 0x93, 0x0b, 0xcb, 0x43, 0xb2, 0x5a, 0x57, 0x52, 0x56, 0x0d, 0x50, 0xe7, 0x3b, 0xed,
	0x77, 0xb4, 0xec, 0x35, 0x33, 0x5c, 0x0a, 0x9f, 0x74, 0x59, 0x7f, 0xab, 0x99, 0x52, 0xa0, 0xbb,
	0x49, 0x8f, 0x7d, 0x4d, 0x33, 0x2a, 0xa0, 0x0a, 0x74, 0x01, 0xc2, 0x4c, 0xca, 0xad, 0x81, 0x35,
	0xbc, 0x64, 0x06, 0xe8, 0xf9, 0x30, 0x09, 0x37, 0x2b, 0x59, 0x49, 0x37, 0x52, 0x3b, 0xf9, 0xdb,
	0x7b, 0x5f, 0x42, 0xbc, 0x78, 0xe1, 0xca, 0x1d, 0xd9, 0x6e, 0x51, 0x86, 0xaf, 0x2b, 0xd6, 0x75,
	0xb9, 0xa9, 0xb5, 0xec, 0xab, 0x3a, 0x6f, 0x65, 0x09, 0x31, 0xda, 0xa0, 0xed, 0x22, 0x59, 0x11,
	0x5f, 0x8a, 0x9c, 0x97, 0x22, 0xa9, 0x94, 0xcd, 0x2b, 0xd6, 0xf4, 0x90, 0xe2, 0x8f, 0x3f, 0x3e,
	0x85, 0x07, 0xef, 0x50, 0xb0, 0x44, 0xd9, 0x35, 0x1b, 0x70, 0xe2, 0xfd, 0xc7, 0xb2, 0x84, 0xe8,
	0x3e, 0xbe, 0x0a, 0xa2, 0x54, 0x92, 0x0b, 0x13, 0x07, 0x1b, 0xb4, 0x9d, 0xa7, 0x73, 0xbb, 0x3e,
	0xd3, 0xc1, 0x06, 0x65, 0xbf, 0xa4, 0xe8, 0x29, 0xc6, 0x05, 0x2b, 0x6a, 0xc8, 0x0d, 0x6f, 0x21,
	0x0e, 0xdd, 0x9b, 0xb7, 0x2f, 0xbd, 0xf9, 0x6c, 0x02, 0x95, 0xce, 0xde, 0x7f, 0x3b, 0x44, 0xd9,
	0xdc, 0x59, 0x4e, 0x78, 0x0b, 0xd1, 0x07, 0x84, 0x0f, 0x8b, 0xa6, 0xef, 0x0c, 0xe8, 0xbc, 0xe5,
	0x22, 0xf7, 0xcc, 0xc7, 0x7c, 0xe2, 0xc3, 0x2a, 0xe8, 0xe2, 0xd9, 0x26, 0xdc, 0x2e, 0x92, 0xd7,
	0xe4, 0x7f, 0x3f, 0x8a, 0xfc, 0xc1, 0x86, 0x1c, 0xf9, 0xf0, 0x63, 0x2e, 0xfc, 0xed, 0xf8, 0xf2,
	0x77, 0xf2, 0x73, 0x61, 0xf4, 0x98, 0xdd, 0x29, 0xfe, 0xb2, 0xb2, 0x2a, 0xf1, 0xdd, 0x7f, 0x46,
	0x44, 0x4b, 0x1c, 0x9e, 0xc1, 0xe8, 0x90, 0xcf, 0x33, 0x3b, 0x46, 0x0f, 0xf0, 0xc1, 0x60, 0x21,
	0x3b, 0x76, 0x8b, 0xe4, 0xc6, 0x54, 0xde, 0xfe, 0x3d, 0x99, 0xec, 0x99, 0xdf, 0x78, 0x1c, 0x3c,
	0x42, 0xe9, 0xf2, 0xf3, 0x7e, 0x8d, 0xbe, 0xee, 0xd7, 0xe8, 0xfb, 0x7e, 0x8d, 0xde, 0x04, 0x43,
	0x72, 0x7a, 0xc5, 0xc1, 0x7b, 0xf8, 0x33, 0x00, 0x00, 0xff, 0xff, 0xf6, 0x4d, 0x3e, 0x27, 0xdc,
	0x02, 0x00, 0x00,
}
