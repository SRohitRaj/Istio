// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mixer/adapter/dogstatsd/config/config.proto

package config

/*
	The `dogstatsd` adapter is designed to deliver Istio metric instances to a
	listening [DataDog](https://www.datadoghq.com/) agent.

	This adapter supports the [metric template](https://istio.io/docs/reference/config/policy-and-telemetry/templates/metric/).
*/

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import strconv "strconv"

import strings "strings"
import reflect "reflect"
import github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"

import encoding_binary "encoding/binary"

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

// Describes the type of metric
type Params_MetricInfo_Type int32

const (
	// Default Unknown Type
	UNKNOWN_TYPE Params_MetricInfo_Type = 0
	// Increments a DataDog counter
	COUNTER Params_MetricInfo_Type = 1
	// Sets the new value of a DataDog gauge
	GAUGE Params_MetricInfo_Type = 2
	// DISTRIBUTION is converted to a Timing Histogram for metrics with a time unit and a Histogram for all other units
	DISTRIBUTION Params_MetricInfo_Type = 3
)

var Params_MetricInfo_Type_name = map[int32]string{
	0: "UNKNOWN_TYPE",
	1: "COUNTER",
	2: "GAUGE",
	3: "DISTRIBUTION",
}
var Params_MetricInfo_Type_value = map[string]int32{
	"UNKNOWN_TYPE": 0,
	"COUNTER":      1,
	"GAUGE":        2,
	"DISTRIBUTION": 3,
}

func (Params_MetricInfo_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_config_883669f61f0765ee, []int{0, 1, 0}
}

// Configuration parameter for the DataDog adapter. These params control how Mixer telemetry is transformed and sent to a dogstatsd agent.
//
// The adapter assumes that a dogstatsd agent is running as a sidecar or at some other endpoint that the Mixer can reach.
// Any dimension that is a part of the metric is converted to a tag automatically. The configuration of the DataDog agent/daemon is outside the scope of the adapter.
type Params struct {
	// Address of the dogstatsd server.
	// Default: localhost:8125
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Prefix to prepend to all metrics handled by the adapter. Metric "bar" with prefix "foo." becomes "foo.bar" in DataDog. In order to make sure the metrics get populated into Datadog properly and avoid any billing issues, it's important to leave the metric prefix to its default value of 'istio.'
	// Default: "istio."
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"`
	// Number of individual metrics to buffer before flushing metrics to the network. When buffered, metrics are flushed every 100ms or when the buffer is filled.
	// When buffer is 0, metrics are not buffered.
	// Default: 0
	BufferLength int32 `protobuf:"varint,3,opt,name=buffer_length,json=bufferLength,proto3" json:"buffer_length,omitempty"`
	// Tags to add to every metric. "global": "tag" becomes "global:tag" in DataDog
	// Default: []
	GlobalTags map[string]string `protobuf:"bytes,4,rep,name=global_tags,json=globalTags" json:"global_tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Chance that any particular metric is sampled when emitted; can take the range [0, 1].
	// Default: 1
	SampleRate float64 `protobuf:"fixed64,5,opt,name=sample_rate,json=sampleRate,proto3" json:"sample_rate,omitempty"`
	// Map of a specific metric instance name -> info. If a metric's instance name is not in the map then the metric will not be exported to DataDog.
	Metrics              map[string]*Params_MetricInfo `protobuf:"bytes,6,rep,name=metrics" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_883669f61f0765ee, []int{0}
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

func (m *Params) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Params) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

func (m *Params) GetBufferLength() int32 {
	if m != nil {
		return m.BufferLength
	}
	return 0
}

func (m *Params) GetGlobalTags() map[string]string {
	if m != nil {
		return m.GlobalTags
	}
	return nil
}

func (m *Params) GetSampleRate() float64 {
	if m != nil {
		return m.SampleRate
	}
	return 0
}

func (m *Params) GetMetrics() map[string]*Params_MetricInfo {
	if m != nil {
		return m.Metrics
	}
	return nil
}

// Describes how to represent this metric in DataDog
type Params_MetricInfo struct {
	// Name of the metric in DataDog
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The type of metric
	Type Params_MetricInfo_Type `protobuf:"varint,2,opt,name=type,proto3,enum=adapter.dogstatsd.config.Params_MetricInfo_Type" json:"type,omitempty"`
	// Tags to add to the metric in addition to the dimensions. "tag": "val" becomes "tag:val" in DataDog
	// Default: []
	Tags                 map[string]string `protobuf:"bytes,3,rep,name=tags" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Params_MetricInfo) Reset()      { *m = Params_MetricInfo{} }
func (*Params_MetricInfo) ProtoMessage() {}
func (*Params_MetricInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_883669f61f0765ee, []int{0, 1}
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

func (m *Params_MetricInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Params_MetricInfo) GetType() Params_MetricInfo_Type {
	if m != nil {
		return m.Type
	}
	return UNKNOWN_TYPE
}

func (m *Params_MetricInfo) GetTags() map[string]string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "adapter.dogstatsd.config.Params")
	proto.RegisterMapType((map[string]string)(nil), "adapter.dogstatsd.config.Params.GlobalTagsEntry")
	proto.RegisterMapType((map[string]*Params_MetricInfo)(nil), "adapter.dogstatsd.config.Params.MetricsEntry")
	proto.RegisterType((*Params_MetricInfo)(nil), "adapter.dogstatsd.config.Params.MetricInfo")
	proto.RegisterMapType((map[string]string)(nil), "adapter.dogstatsd.config.Params.MetricInfo.TagsEntry")
	proto.RegisterEnum("adapter.dogstatsd.config.Params_MetricInfo_Type", Params_MetricInfo_Type_name, Params_MetricInfo_Type_value)
}
func (x Params_MetricInfo_Type) String() string {
	s, ok := Params_MetricInfo_Type_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
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
	if this.Address != that1.Address {
		return false
	}
	if this.Prefix != that1.Prefix {
		return false
	}
	if this.BufferLength != that1.BufferLength {
		return false
	}
	if len(this.GlobalTags) != len(that1.GlobalTags) {
		return false
	}
	for i := range this.GlobalTags {
		if this.GlobalTags[i] != that1.GlobalTags[i] {
			return false
		}
	}
	if this.SampleRate != that1.SampleRate {
		return false
	}
	if len(this.Metrics) != len(that1.Metrics) {
		return false
	}
	for i := range this.Metrics {
		if !this.Metrics[i].Equal(that1.Metrics[i]) {
			return false
		}
	}
	return true
}
func (this *Params_MetricInfo) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params_MetricInfo)
	if !ok {
		that2, ok := that.(Params_MetricInfo)
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
	if this.Name != that1.Name {
		return false
	}
	if this.Type != that1.Type {
		return false
	}
	if len(this.Tags) != len(that1.Tags) {
		return false
	}
	for i := range this.Tags {
		if this.Tags[i] != that1.Tags[i] {
			return false
		}
	}
	return true
}
func (this *Params) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 10)
	s = append(s, "&config.Params{")
	s = append(s, "Address: "+fmt.Sprintf("%#v", this.Address)+",\n")
	s = append(s, "Prefix: "+fmt.Sprintf("%#v", this.Prefix)+",\n")
	s = append(s, "BufferLength: "+fmt.Sprintf("%#v", this.BufferLength)+",\n")
	keysForGlobalTags := make([]string, 0, len(this.GlobalTags))
	for k, _ := range this.GlobalTags {
		keysForGlobalTags = append(keysForGlobalTags, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForGlobalTags)
	mapStringForGlobalTags := "map[string]string{"
	for _, k := range keysForGlobalTags {
		mapStringForGlobalTags += fmt.Sprintf("%#v: %#v,", k, this.GlobalTags[k])
	}
	mapStringForGlobalTags += "}"
	if this.GlobalTags != nil {
		s = append(s, "GlobalTags: "+mapStringForGlobalTags+",\n")
	}
	s = append(s, "SampleRate: "+fmt.Sprintf("%#v", this.SampleRate)+",\n")
	keysForMetrics := make([]string, 0, len(this.Metrics))
	for k, _ := range this.Metrics {
		keysForMetrics = append(keysForMetrics, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForMetrics)
	mapStringForMetrics := "map[string]*Params_MetricInfo{"
	for _, k := range keysForMetrics {
		mapStringForMetrics += fmt.Sprintf("%#v: %#v,", k, this.Metrics[k])
	}
	mapStringForMetrics += "}"
	if this.Metrics != nil {
		s = append(s, "Metrics: "+mapStringForMetrics+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Params_MetricInfo) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&config.Params_MetricInfo{")
	s = append(s, "Name: "+fmt.Sprintf("%#v", this.Name)+",\n")
	s = append(s, "Type: "+fmt.Sprintf("%#v", this.Type)+",\n")
	keysForTags := make([]string, 0, len(this.Tags))
	for k, _ := range this.Tags {
		keysForTags = append(keysForTags, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForTags)
	mapStringForTags := "map[string]string{"
	for _, k := range keysForTags {
		mapStringForTags += fmt.Sprintf("%#v: %#v,", k, this.Tags[k])
	}
	mapStringForTags += "}"
	if this.Tags != nil {
		s = append(s, "Tags: "+mapStringForTags+",\n")
	}
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
	if m.BufferLength != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintConfig(dAtA, i, uint64(m.BufferLength))
	}
	if len(m.GlobalTags) > 0 {
		for k, _ := range m.GlobalTags {
			dAtA[i] = 0x22
			i++
			v := m.GlobalTags[k]
			mapSize := 1 + len(k) + sovConfig(uint64(len(k))) + 1 + len(v) + sovConfig(uint64(len(v)))
			i = encodeVarintConfig(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintConfig(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintConfig(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	if m.SampleRate != 0 {
		dAtA[i] = 0x29
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.SampleRate))))
		i += 8
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
				n1, err := v.MarshalTo(dAtA[i:])
				if err != nil {
					return 0, err
				}
				i += n1
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
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintConfig(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.Type != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintConfig(dAtA, i, uint64(m.Type))
	}
	if len(m.Tags) > 0 {
		for k, _ := range m.Tags {
			dAtA[i] = 0x1a
			i++
			v := m.Tags[k]
			mapSize := 1 + len(k) + sovConfig(uint64(len(k))) + 1 + len(v) + sovConfig(uint64(len(v)))
			i = encodeVarintConfig(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintConfig(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintConfig(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
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
	if m.BufferLength != 0 {
		n += 1 + sovConfig(uint64(m.BufferLength))
	}
	if len(m.GlobalTags) > 0 {
		for k, v := range m.GlobalTags {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovConfig(uint64(len(k))) + 1 + len(v) + sovConfig(uint64(len(v)))
			n += mapEntrySize + 1 + sovConfig(uint64(mapEntrySize))
		}
	}
	if m.SampleRate != 0 {
		n += 9
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
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovConfig(uint64(m.Type))
	}
	if len(m.Tags) > 0 {
		for k, v := range m.Tags {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovConfig(uint64(len(k))) + 1 + len(v) + sovConfig(uint64(len(v)))
			n += mapEntrySize + 1 + sovConfig(uint64(mapEntrySize))
		}
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
	keysForGlobalTags := make([]string, 0, len(this.GlobalTags))
	for k, _ := range this.GlobalTags {
		keysForGlobalTags = append(keysForGlobalTags, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForGlobalTags)
	mapStringForGlobalTags := "map[string]string{"
	for _, k := range keysForGlobalTags {
		mapStringForGlobalTags += fmt.Sprintf("%v: %v,", k, this.GlobalTags[k])
	}
	mapStringForGlobalTags += "}"
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
		`BufferLength:` + fmt.Sprintf("%v", this.BufferLength) + `,`,
		`GlobalTags:` + mapStringForGlobalTags + `,`,
		`SampleRate:` + fmt.Sprintf("%v", this.SampleRate) + `,`,
		`Metrics:` + mapStringForMetrics + `,`,
		`}`,
	}, "")
	return s
}
func (this *Params_MetricInfo) String() string {
	if this == nil {
		return "nil"
	}
	keysForTags := make([]string, 0, len(this.Tags))
	for k, _ := range this.Tags {
		keysForTags = append(keysForTags, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForTags)
	mapStringForTags := "map[string]string{"
	for _, k := range keysForTags {
		mapStringForTags += fmt.Sprintf("%v: %v,", k, this.Tags[k])
	}
	mapStringForTags += "}"
	s := strings.Join([]string{`&Params_MetricInfo{`,
		`Name:` + fmt.Sprintf("%v", this.Name) + `,`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`Tags:` + mapStringForTags + `,`,
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BufferLength", wireType)
			}
			m.BufferLength = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BufferLength |= (int32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GlobalTags", wireType)
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
			if m.GlobalTags == nil {
				m.GlobalTags = make(map[string]string)
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
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthConfig
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
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
			m.GlobalTags[mapkey] = mapvalue
			iNdEx = postIndex
		case 5:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field SampleRate", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.SampleRate = float64(math.Float64frombits(v))
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
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
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tags", wireType)
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
			if m.Tags == nil {
				m.Tags = make(map[string]string)
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
						stringLenmapvalue |= (uint64(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthConfig
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
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
			m.Tags[mapkey] = mapvalue
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
	proto.RegisterFile("mixer/adapter/dogstatsd/config/config.proto", fileDescriptor_config_883669f61f0765ee)
}

var fileDescriptor_config_883669f61f0765ee = []byte{
	// 496 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0xcd, 0x8a, 0x13, 0x4d,
	0x14, 0x4d, 0xe5, 0x97, 0xdc, 0xe4, 0xfb, 0x6c, 0x8a, 0x41, 0x9a, 0x5e, 0x94, 0x61, 0xdc, 0x04,
	0x06, 0x3b, 0x32, 0x2a, 0x8a, 0xe0, 0x62, 0xe2, 0x84, 0xd0, 0xa8, 0x9d, 0xb1, 0xed, 0x20, 0xba,
	0x09, 0x95, 0xe9, 0xea, 0x32, 0xd8, 0x7f, 0x54, 0xd7, 0xc8, 0x64, 0xe7, 0x23, 0xf8, 0x18, 0xae,
	0x7c, 0x0e, 0x97, 0xb3, 0x74, 0x69, 0xda, 0x8d, 0xcb, 0xd9, 0xb8, 0x97, 0x54, 0xf5, 0xc4, 0x20,
	0x8a, 0x8e, 0xab, 0xdc, 0x7b, 0xeb, 0xdc, 0x73, 0xce, 0x3d, 0xa1, 0x61, 0x2f, 0x5e, 0x9c, 0x32,
	0x31, 0xa0, 0x01, 0xcd, 0x24, 0x13, 0x83, 0x20, 0xe5, 0xb9, 0xa4, 0x32, 0x0f, 0x06, 0xc7, 0x69,
	0x12, 0x2e, 0x78, 0xf9, 0x63, 0x67, 0x22, 0x95, 0x29, 0x36, 0x4b, 0x98, 0xbd, 0x81, 0xd9, 0xfa,
	0xdd, 0xda, 0xe1, 0x29, 0x4f, 0x15, 0x68, 0xb0, 0xae, 0x34, 0x7e, 0xf7, 0x5b, 0x03, 0x9a, 0x47,
	0x54, 0xd0, 0x38, 0xc7, 0x26, 0xb4, 0x68, 0x10, 0x08, 0x96, 0xe7, 0x26, 0xea, 0xa1, 0x7e, 0xdb,
	0xbb, 0x68, 0xf1, 0x55, 0x68, 0x66, 0x82, 0x85, 0x8b, 0x53, 0xb3, 0xaa, 0x1e, 0xca, 0x0e, 0x5f,
	0x87, 0xff, 0xe6, 0x27, 0x61, 0xc8, 0xc4, 0x2c, 0x62, 0x09, 0x97, 0xaf, 0xcc, 0x5a, 0x0f, 0xf5,
	0x1b, 0x5e, 0x57, 0x0f, 0x1f, 0xab, 0x19, 0x7e, 0x0a, 0x1d, 0x1e, 0xa5, 0x73, 0x1a, 0xcd, 0x24,
	0xe5, 0xb9, 0x59, 0xef, 0xd5, 0xfa, 0x9d, 0xfd, 0x9b, 0xf6, 0xef, 0x7c, 0xda, 0xda, 0x8d, 0x3d,
	0x56, 0x3b, 0x3e, 0xe5, 0xf9, 0x28, 0x91, 0x62, 0xe9, 0x01, 0xdf, 0x0c, 0xf0, 0x35, 0xe8, 0xe4,
	0x34, 0xce, 0x22, 0x36, 0x13, 0x54, 0x32, 0xb3, 0xd1, 0x43, 0x7d, 0xe4, 0x81, 0x1e, 0x79, 0x54,
	0x32, 0x3c, 0x86, 0x56, 0xcc, 0xa4, 0x58, 0x1c, 0xe7, 0x66, 0x53, 0xe9, 0xdd, 0xf8, 0xa3, 0xde,
	0x13, 0x8d, 0xd7, 0x62, 0x17, 0xdb, 0xd6, 0x03, 0xb8, 0xf2, 0x93, 0x11, 0x6c, 0x40, 0xed, 0x35,
	0x5b, 0x96, 0x11, 0xad, 0x4b, 0xbc, 0x03, 0x8d, 0x37, 0x34, 0x3a, 0x61, 0x65, 0x3a, 0xba, 0xb9,
	0x5f, 0xbd, 0x87, 0xac, 0x0f, 0x55, 0x00, 0x4d, 0xec, 0x24, 0x61, 0x8a, 0x31, 0xd4, 0x13, 0x1a,
	0xb3, 0x72, 0x57, 0xd5, 0xf8, 0x10, 0xea, 0x72, 0x99, 0xe9, 0xdd, 0xff, 0xff, 0x22, 0x97, 0x1f,
	0x74, 0xb6, 0xbf, 0xcc, 0x98, 0xa7, 0xb6, 0xb1, 0x03, 0x75, 0x95, 0x6e, 0x4d, 0x5d, 0x7b, 0xe7,
	0x52, 0x2c, 0x9b, 0x88, 0x15, 0x85, 0x75, 0x17, 0xda, 0xff, 0x74, 0xec, 0xee, 0x10, 0xea, 0x6b,
	0x47, 0xd8, 0x80, 0xee, 0xd4, 0x7d, 0xe4, 0x4e, 0x9e, 0xbb, 0x33, 0xff, 0xc5, 0xd1, 0xc8, 0xa8,
	0xe0, 0x0e, 0xb4, 0x1e, 0x4e, 0xa6, 0xae, 0x3f, 0xf2, 0x0c, 0x84, 0xdb, 0xd0, 0x18, 0x1f, 0x4c,
	0xc7, 0x23, 0xa3, 0xba, 0x46, 0x1e, 0x3a, 0xcf, 0x7c, 0xcf, 0x19, 0x4e, 0x7d, 0x67, 0xe2, 0x1a,
	0x35, 0x8b, 0x43, 0x77, 0xfb, 0x8f, 0xf8, 0x85, 0xfe, 0xc1, 0xb6, 0x7e, 0x67, 0x7f, 0xef, 0x12,
	0xa7, 0x6e, 0x99, 0x1d, 0xde, 0x3e, 0x5b, 0x91, 0xca, 0xa7, 0x15, 0xa9, 0x9c, 0xaf, 0x08, 0x7a,
	0x5b, 0x10, 0xf4, 0xbe, 0x20, 0xe8, 0x63, 0x41, 0xd0, 0x59, 0x41, 0xd0, 0xe7, 0x82, 0xa0, 0xaf,
	0x05, 0xa9, 0x9c, 0x17, 0x04, 0xbd, 0xfb, 0x42, 0x2a, 0x2f, 0x9b, 0x9a, 0x72, 0xde, 0x54, 0x1f,
	0xcd, 0xad, 0xef, 0x01, 0x00, 0x00, 0xff, 0xff, 0x78, 0xb4, 0xe2, 0xb6, 0x93, 0x03, 0x00, 0x00,
}
