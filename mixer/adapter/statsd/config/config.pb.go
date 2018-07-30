// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mixer/adapter/statsd/config/config.proto

package config

/*
	The `statsd` adapter enables Istio to deliver metric data to a
	[StatsD](https://github.com/etsy/statsd) monitoring backend.

	This adapter supports the [metric template](https://istio.io/docs/reference/config/policy-and-telemetry/templates/metric/).
*/

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"

import time "time"

import strconv "strconv"

import encoding_binary "encoding/binary"
import github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"

import strings "strings"
import reflect "reflect"
import github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"

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

// The type of metric.
type Params_MetricInfo_Type int32

const (
	UNKNOWN      Params_MetricInfo_Type = 0
	COUNTER      Params_MetricInfo_Type = 1
	GAUGE        Params_MetricInfo_Type = 2
	DISTRIBUTION Params_MetricInfo_Type = 3
)

var Params_MetricInfo_Type_name = map[int32]string{
	0: "UNKNOWN",
	1: "COUNTER",
	2: "GAUGE",
	3: "DISTRIBUTION",
}
var Params_MetricInfo_Type_value = map[string]int32{
	"UNKNOWN":      0,
	"COUNTER":      1,
	"GAUGE":        2,
	"DISTRIBUTION": 3,
}

func (Params_MetricInfo_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_1e600eac25cec0f9, []int{0, 0, 0}
}

// Configuration format for the `statsd` adapter.
type Params struct {
	// Address of the statsd server, e.g. localhost:8125
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Metric prefix, do not specify for no prefix
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"`
	// FlushDuration controls the maximum amount of time between sending metrics to the statsd collection server.
	// Metrics are reported when either flush_bytes is full or flush_duration time has elapsed since the last report.
	FlushDuration time.Duration `protobuf:"bytes,3,opt,name=flush_duration,json=flushDuration,stdduration" json:"flush_duration"`
	// Maximum UDP packet size to send; if not specified defaults to 512 bytes. If the statsd server is running on the
	// same (private) network 1432 bytes is recommended for better performance.
	FlushBytes int32 `protobuf:"varint,4,opt,name=flush_bytes,json=flushBytes,proto3" json:"flush_bytes,omitempty"`
	// Chance that any particular metric is sampled when incremented; can take the range [0, 1], defaults to 1 if unspecified.
	SamplingRate float32 `protobuf:"fixed32,5,opt,name=sampling_rate,json=samplingRate,proto3" json:"sampling_rate,omitempty"`
	// Map of metric name -> info. If a metric's name is not in the map then the metric will not be exported to statsd.
	Metrics              map[string]*Params_MetricInfo `protobuf:"bytes,6,rep,name=metrics" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_1e600eac25cec0f9, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(dst, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

// Describes how to represent this metric in statsd
type Params_MetricInfo struct {
	Type Params_MetricInfo_Type `protobuf:"varint,1,opt,name=type,proto3,enum=adapter.statsd.config.Params_MetricInfo_Type" json:"type,omitempty"`
	// The template will be filled with values from the metric's labels and the resulting string will be used as
	// the statsd metric name. This allows easier creation of statsd metrics like `action_name-response_code`.
	// The template strings must conform to go's text/template syntax. For the example of `action_name-response_code`,
	// we use the template:
	//    `{{.apiMethod}}-{{.responseCode}}`
	//
	// If name_template is the empty string the Istio metric name will be used for statsd metric's name.
	NameTemplate         string   `protobuf:"bytes,2,opt,name=name_template,json=nameTemplate,proto3" json:"name_template,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Params_MetricInfo) Reset()      { *m = Params_MetricInfo{} }
func (*Params_MetricInfo) ProtoMessage() {}
func (*Params_MetricInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_1e600eac25cec0f9, []int{0, 0}
}
func (m *Params_MetricInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params_MetricInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params_MetricInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *Params_MetricInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params_MetricInfo.Merge(dst, src)
}
func (m *Params_MetricInfo) XXX_Size() int {
	return m.Size()
}
func (m *Params_MetricInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_Params_MetricInfo.DiscardUnknown(m)
}

var xxx_messageInfo_Params_MetricInfo proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "adapter.statsd.config.Params")
	proto.RegisterMapType((map[string]*Params_MetricInfo)(nil), "adapter.statsd.config.Params.MetricsEntry")
	proto.RegisterType((*Params_MetricInfo)(nil), "adapter.statsd.config.Params.MetricInfo")
	proto.RegisterEnum("adapter.statsd.config.Params_MetricInfo_Type", Params_MetricInfo_Type_name, Params_MetricInfo_Type_value)
}
func (x Params_MetricInfo_Type) String() string {
	s, ok := Params_MetricInfo_Type_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Address) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintConfig(dAtA, i, uint64(len(m.Address)))
		i += copy(dAtA[i:], m.Address)
	}
	if len(m.Prefix) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintConfig(dAtA, i, uint64(len(m.Prefix)))
		i += copy(dAtA[i:], m.Prefix)
	}
	dAtA[i] = 0x1a
	i++
	i = encodeVarintConfig(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdDuration(m.FlushDuration)))
	n1, err := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.FlushDuration, dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	if m.FlushBytes != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintConfig(dAtA, i, uint64(m.FlushBytes))
	}
	if m.SamplingRate != 0 {
		dAtA[i] = 0x2d
		i++
		encoding_binary.LittleEndian.PutUint32(dAtA[i:], uint32(math.Float32bits(float32(m.SamplingRate))))
		i += 4
	}
	if len(m.Metrics) > 0 {
		for k, _ := range m.Metrics {
			dAtA[i] = 0x32
			i++
			v := m.Metrics[k]
			msgSize := 0
			if v != nil {
				msgSize = v.Size()
				msgSize += 1 + sovConfig(uint64(msgSize))
			}
			mapSize := 1 + len(k) + sovConfig(uint64(len(k))) + msgSize
			i = encodeVarintConfig(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintConfig(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			if v != nil {
				dAtA[i] = 0x12
				i++
				i = encodeVarintConfig(dAtA, i, uint64(v.Size()))
				n2, err := v.MarshalTo(dAtA[i:])
				if err != nil {
					return 0, err
				}
				i += n2
			}
		}
	}
	return i, nil
}

func (m *Params_MetricInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params_MetricInfo) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Type != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintConfig(dAtA, i, uint64(m.Type))
	}
	if len(m.NameTemplate) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintConfig(dAtA, i, uint64(len(m.NameTemplate)))
		i += copy(dAtA[i:], m.NameTemplate)
	}
	return i, nil
}

func encodeVarintConfig(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Params) Size() (n int) {
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	l = len(m.Prefix)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.FlushDuration)
	n += 1 + l + sovConfig(uint64(l))
	if m.FlushBytes != 0 {
		n += 1 + sovConfig(uint64(m.FlushBytes))
	}
	if m.SamplingRate != 0 {
		n += 5
	}
	if len(m.Metrics) > 0 {
		for k, v := range m.Metrics {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovConfig(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovConfig(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovConfig(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *Params_MetricInfo) Size() (n int) {
	var l int
	_ = l
	if m.Type != 0 {
		n += 1 + sovConfig(uint64(m.Type))
	}
	l = len(m.NameTemplate)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	return n
}

func sovConfig(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozConfig(x uint64) (n int) {
	return sovConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Params) String() string {
	if this == nil {
		return "nil"
	}
	keysForMetrics := make([]string, 0, len(this.Metrics))
	for k, _ := range this.Metrics {
		keysForMetrics = append(keysForMetrics, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForMetrics)
	mapStringForMetrics := "map[string]*Params_MetricInfo{"
	for _, k := range keysForMetrics {
		mapStringForMetrics += fmt.Sprintf("%v: %v,", k, this.Metrics[k])
	}
	mapStringForMetrics += "}"
	s := strings.Join([]string{`&Params{`,
		`Address:` + fmt.Sprintf("%v", this.Address) + `,`,
		`Prefix:` + fmt.Sprintf("%v", this.Prefix) + `,`,
		`FlushDuration:` + strings.Replace(strings.Replace(this.FlushDuration.String(), "Duration", "types.Duration", 1), `&`, ``, 1) + `,`,
		`FlushBytes:` + fmt.Sprintf("%v", this.FlushBytes) + `,`,
		`SamplingRate:` + fmt.Sprintf("%v", this.SamplingRate) + `,`,
		`Metrics:` + mapStringForMetrics + `,`,
		`}`,
	}, "")
	return s
}
func (this *Params_MetricInfo) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Params_MetricInfo{`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`NameTemplate:` + fmt.Sprintf("%v", this.NameTemplate) + `,`,
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
			wire |= (uint64(b) & 0x7F) << shift
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
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prefix", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prefix = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FlushDuration", wireType)
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
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.FlushDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FlushBytes", wireType)
			}
			m.FlushBytes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FlushBytes |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 5 {
				return fmt.Errorf("proto: wrong wireType = %d for field SamplingRate", wireType)
			}
			var v uint32
			if (iNdEx + 4) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint32(encoding_binary.LittleEndian.Uint32(dAtA[iNdEx:]))
			iNdEx += 4
			m.SamplingRate = float32(math.Float32frombits(v))
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Metrics", wireType)
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
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Metrics == nil {
				m.Metrics = make(map[string]*Params_MetricInfo)
			}
			var mapkey string
			var mapvalue *Params_MetricInfo
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
							return ErrIntOverflowConfig
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
						return ErrInvalidLengthConfig
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
							return ErrIntOverflowConfig
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
						return ErrInvalidLengthConfig
					}
					postmsgIndex := iNdEx + mapmsglen
					if mapmsglen < 0 {
						return ErrInvalidLengthConfig
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &Params_MetricInfo{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
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
			m.Metrics[mapkey] = mapvalue
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
func (m *Params_MetricInfo) Unmarshal(dAtA []byte) error {
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
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetricInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetricInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= (Params_MetricInfo_Type(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NameTemplate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
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
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NameTemplate = string(dAtA[iNdEx:postIndex])
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
			iNdEx += length
			if length < 0 {
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

func init() {
	proto.RegisterFile("mixer/adapter/statsd/config/config.proto", fileDescriptor_config_1e600eac25cec0f9)
}

var fileDescriptor_config_1e600eac25cec0f9 = []byte{
	// 494 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xbd, 0x8e, 0xd3, 0x40,
	0x10, 0xf6, 0xe6, 0xef, 0xb8, 0x4d, 0xee, 0x64, 0xad, 0x00, 0x99, 0x14, 0x9b, 0xe8, 0x68, 0x2c,
	0x24, 0xd6, 0x52, 0x68, 0x4e, 0x48, 0x20, 0x5d, 0x48, 0x74, 0x0a, 0x08, 0x07, 0x2d, 0x8e, 0x90,
	0x68, 0xa2, 0xcd, 0x79, 0x6d, 0x2c, 0xfc, 0xa7, 0xf5, 0x06, 0x9d, 0x3b, 0x1e, 0x81, 0x92, 0x47,
	0xa0, 0xa2, 0xe1, 0x25, 0x52, 0x5e, 0x49, 0x05, 0xc4, 0x34, 0x94, 0xf7, 0x08, 0xc8, 0x5e, 0x5b,
	0xa2, 0xa0, 0xb8, 0xca, 0x33, 0xdf, 0x7c, 0xdf, 0xcc, 0x37, 0xe3, 0x85, 0x66, 0x14, 0x5c, 0x72,
	0x61, 0x31, 0x97, 0xa5, 0x92, 0x0b, 0x2b, 0x93, 0x4c, 0x66, 0xae, 0x75, 0x91, 0xc4, 0x5e, 0xe0,
	0xd7, 0x1f, 0x92, 0x8a, 0x44, 0x26, 0xe8, 0x4e, 0xcd, 0x21, 0x8a, 0x43, 0x54, 0x71, 0x88, 0xfd,
	0x24, 0xf1, 0x43, 0x6e, 0x55, 0xa4, 0xcd, 0xd6, 0xb3, 0xdc, 0xad, 0x60, 0x32, 0x48, 0x62, 0x25,
	0x1b, 0xde, 0xf6, 0x13, 0x3f, 0xa9, 0x42, 0xab, 0x8c, 0x14, 0x7a, 0xf2, 0xb5, 0x03, 0x7b, 0xaf,
	0x98, 0x60, 0x51, 0x86, 0x0c, 0x78, 0xc0, 0x5c, 0x57, 0xf0, 0x2c, 0x33, 0xc0, 0x18, 0x98, 0x87,
	0xb4, 0x49, 0xd1, 0x5d, 0xd8, 0x4b, 0x05, 0xf7, 0x82, 0x4b, 0xa3, 0x55, 0x15, 0xea, 0x0c, 0x3d,
	0x87, 0xc7, 0x5e, 0xb8, 0xcd, 0xde, 0xad, 0x9b, 0x51, 0x46, 0x7b, 0x0c, 0xcc, 0xfe, 0xe4, 0x1e,
	0x51, 0x5e, 0x48, 0xe3, 0x85, 0xcc, 0x6a, 0xc2, 0xf4, 0xd6, 0xee, 0xc7, 0x48, 0xfb, 0xfc, 0x73,
	0x04, 0xe8, 0x51, 0x25, 0x6d, 0x0a, 0x68, 0x04, 0xfb, 0xaa, 0xd7, 0x26, 0x97, 0x3c, 0x33, 0x3a,
	0x63, 0x60, 0x76, 0x29, 0xac, 0xa0, 0x69, 0x89, 0xa0, 0xfb, 0xf0, 0x28, 0x63, 0x51, 0x1a, 0x06,
	0xb1, 0xbf, 0x16, 0x4c, 0x72, 0xa3, 0x3b, 0x06, 0x66, 0x8b, 0x0e, 0x1a, 0x90, 0x32, 0xc9, 0xd1,
	0x0c, 0x1e, 0x44, 0x5c, 0x8a, 0xe0, 0x22, 0x33, 0x7a, 0xe3, 0xb6, 0xd9, 0x9f, 0x3c, 0x20, 0xff,
	0xbd, 0x16, 0x51, 0x3b, 0x93, 0x97, 0x8a, 0x3c, 0x8f, 0xa5, 0xc8, 0x69, 0x23, 0x1d, 0x7e, 0x03,
	0x10, 0xaa, 0xca, 0x22, 0xf6, 0x12, 0x74, 0x06, 0x3b, 0x32, 0x4f, 0x79, 0x75, 0x95, 0xe3, 0xc9,
	0xc3, 0x9b, 0x74, 0x2c, 0x75, 0xc4, 0xc9, 0x53, 0x4e, 0x2b, 0x69, 0x69, 0x3e, 0x66, 0x11, 0x5f,
	0x4b, 0x1e, 0xa5, 0x61, 0x69, 0x5e, 0x1d, 0x72, 0x50, 0x82, 0x4e, 0x8d, 0x9d, 0x3c, 0x81, 0x9d,
	0x52, 0x82, 0xfa, 0xf0, 0x60, 0x65, 0xbf, 0xb0, 0x97, 0x6f, 0x6c, 0x5d, 0x2b, 0x93, 0x67, 0xcb,
	0x95, 0xed, 0xcc, 0xa9, 0x0e, 0xd0, 0x21, 0xec, 0x9e, 0x9f, 0xad, 0xce, 0xe7, 0x7a, 0x0b, 0xe9,
	0x70, 0x30, 0x5b, 0xbc, 0x76, 0xe8, 0x62, 0xba, 0x72, 0x16, 0x4b, 0x5b, 0x6f, 0x0f, 0x5d, 0x38,
	0xf8, 0x77, 0x1d, 0xa4, 0xc3, 0xf6, 0x7b, 0x9e, 0xd7, 0xff, 0xb2, 0x0c, 0xd1, 0x53, 0xd8, 0xfd,
	0xc0, 0xc2, 0xad, 0x9a, 0xde, 0x9f, 0x98, 0x37, 0xdd, 0x84, 0x2a, 0xd9, 0xe3, 0xd6, 0x29, 0x98,
	0x9e, 0xee, 0xf6, 0x58, 0xbb, 0xda, 0x63, 0xed, 0xfb, 0x1e, 0x6b, 0xd7, 0x7b, 0xac, 0x7d, 0x2c,
	0x30, 0xf8, 0x52, 0x60, 0x6d, 0x57, 0x60, 0x70, 0x55, 0x60, 0xf0, 0xab, 0xc0, 0xe0, 0x4f, 0x81,
	0xb5, 0xeb, 0x02, 0x83, 0x4f, 0xbf, 0xb1, 0xf6, 0xb6, 0xa7, 0xda, 0x6e, 0x7a, 0xd5, 0x6b, 0x78,
	0xf4, 0x37, 0x00, 0x00, 0xff, 0xff, 0x6a, 0x40, 0x2c, 0x67, 0xea, 0x02, 0x00, 0x00,
}
