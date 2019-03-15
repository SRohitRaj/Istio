// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/api/v2/cluster/outlier_detection.proto

package cluster

import (
	bytes "bytes"
	fmt "fmt"
	io "io"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
	_ "github.com/lyft/protoc-gen-validate/validate"
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

// See the :ref:`architecture overview <arch_overview_outlier_detection>` for
// more information on outlier detection.
type OutlierDetection struct {
	// The number of consecutive 5xx responses before a consecutive 5xx ejection
	// occurs. Defaults to 5.
	Consecutive_5Xx *types.UInt32Value `protobuf:"bytes,1,opt,name=consecutive_5xx,json=consecutive5xx,proto3" json:"consecutive_5xx,omitempty"`
	// The time interval between ejection analysis sweeps. This can result in
	// both new ejections as well as hosts being returned to service. Defaults
	// to 10000ms or 10s.
	Interval *types.Duration `protobuf:"bytes,2,opt,name=interval,proto3" json:"interval,omitempty"`
	// The base time that a host is ejected for. The real time is equal to the
	// base time multiplied by the number of times the host has been ejected.
	// Defaults to 30000ms or 30s.
	BaseEjectionTime *types.Duration `protobuf:"bytes,3,opt,name=base_ejection_time,json=baseEjectionTime,proto3" json:"base_ejection_time,omitempty"`
	// The maximum % of an upstream cluster that can be ejected due to outlier
	// detection. Defaults to 10% but will eject at least one host regardless of the value.
	MaxEjectionPercent *types.UInt32Value `protobuf:"bytes,4,opt,name=max_ejection_percent,json=maxEjectionPercent,proto3" json:"max_ejection_percent,omitempty"`
	// The % chance that a host will be actually ejected when an outlier status
	// is detected through consecutive 5xx. This setting can be used to disable
	// ejection or to ramp it up slowly. Defaults to 100.
	EnforcingConsecutive_5Xx *types.UInt32Value `protobuf:"bytes,5,opt,name=enforcing_consecutive_5xx,json=enforcingConsecutive5xx,proto3" json:"enforcing_consecutive_5xx,omitempty"`
	// The % chance that a host will be actually ejected when an outlier status
	// is detected through success rate statistics. This setting can be used to
	// disable ejection or to ramp it up slowly. Defaults to 100.
	EnforcingSuccessRate *types.UInt32Value `protobuf:"bytes,6,opt,name=enforcing_success_rate,json=enforcingSuccessRate,proto3" json:"enforcing_success_rate,omitempty"`
	// The number of hosts in a cluster that must have enough request volume to
	// detect success rate outliers. If the number of hosts is less than this
	// setting, outlier detection via success rate statistics is not performed
	// for any host in the cluster. Defaults to 5.
	SuccessRateMinimumHosts *types.UInt32Value `protobuf:"bytes,7,opt,name=success_rate_minimum_hosts,json=successRateMinimumHosts,proto3" json:"success_rate_minimum_hosts,omitempty"`
	// The minimum number of total requests that must be collected in one
	// interval (as defined by the interval duration above) to include this host
	// in success rate based outlier detection. If the volume is lower than this
	// setting, outlier detection via success rate statistics is not performed
	// for that host. Defaults to 100.
	SuccessRateRequestVolume *types.UInt32Value `protobuf:"bytes,8,opt,name=success_rate_request_volume,json=successRateRequestVolume,proto3" json:"success_rate_request_volume,omitempty"`
	// This factor is used to determine the ejection threshold for success rate
	// outlier ejection. The ejection threshold is the difference between the
	// mean success rate, and the product of this factor and the standard
	// deviation of the mean success rate: mean - (stdev *
	// success_rate_stdev_factor). This factor is divided by a thousand to get a
	// double. That is, if the desired factor is 1.9, the runtime value should
	// be 1900. Defaults to 1900.
	SuccessRateStdevFactor *types.UInt32Value `protobuf:"bytes,9,opt,name=success_rate_stdev_factor,json=successRateStdevFactor,proto3" json:"success_rate_stdev_factor,omitempty"`
	// The number of consecutive gateway failures (502, 503, 504 status or
	// connection errors that are mapped to one of those status codes) before a
	// consecutive gateway failure ejection occurs. Defaults to 5.
	ConsecutiveGatewayFailure *types.UInt32Value `protobuf:"bytes,10,opt,name=consecutive_gateway_failure,json=consecutiveGatewayFailure,proto3" json:"consecutive_gateway_failure,omitempty"`
	// The % chance that a host will be actually ejected when an outlier status
	// is detected through consecutive gateway failures. This setting can be
	// used to disable ejection or to ramp it up slowly. Defaults to 0.
	EnforcingConsecutiveGatewayFailure *types.UInt32Value `protobuf:"bytes,11,opt,name=enforcing_consecutive_gateway_failure,json=enforcingConsecutiveGatewayFailure,proto3" json:"enforcing_consecutive_gateway_failure,omitempty"`
	XXX_NoUnkeyedLiteral               struct{}           `json:"-"`
	XXX_unrecognized                   []byte             `json:"-"`
	XXX_sizecache                      int32              `json:"-"`
}

func (m *OutlierDetection) Reset()         { *m = OutlierDetection{} }
func (m *OutlierDetection) String() string { return proto.CompactTextString(m) }
func (*OutlierDetection) ProtoMessage()    {}
func (*OutlierDetection) Descriptor() ([]byte, []int) {
	return fileDescriptor_56cd87362a3f00c9, []int{0}
}
func (m *OutlierDetection) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OutlierDetection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OutlierDetection.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OutlierDetection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OutlierDetection.Merge(m, src)
}
func (m *OutlierDetection) XXX_Size() int {
	return m.Size()
}
func (m *OutlierDetection) XXX_DiscardUnknown() {
	xxx_messageInfo_OutlierDetection.DiscardUnknown(m)
}

var xxx_messageInfo_OutlierDetection proto.InternalMessageInfo

func (m *OutlierDetection) GetConsecutive_5Xx() *types.UInt32Value {
	if m != nil {
		return m.Consecutive_5Xx
	}
	return nil
}

func (m *OutlierDetection) GetInterval() *types.Duration {
	if m != nil {
		return m.Interval
	}
	return nil
}

func (m *OutlierDetection) GetBaseEjectionTime() *types.Duration {
	if m != nil {
		return m.BaseEjectionTime
	}
	return nil
}

func (m *OutlierDetection) GetMaxEjectionPercent() *types.UInt32Value {
	if m != nil {
		return m.MaxEjectionPercent
	}
	return nil
}

func (m *OutlierDetection) GetEnforcingConsecutive_5Xx() *types.UInt32Value {
	if m != nil {
		return m.EnforcingConsecutive_5Xx
	}
	return nil
}

func (m *OutlierDetection) GetEnforcingSuccessRate() *types.UInt32Value {
	if m != nil {
		return m.EnforcingSuccessRate
	}
	return nil
}

func (m *OutlierDetection) GetSuccessRateMinimumHosts() *types.UInt32Value {
	if m != nil {
		return m.SuccessRateMinimumHosts
	}
	return nil
}

func (m *OutlierDetection) GetSuccessRateRequestVolume() *types.UInt32Value {
	if m != nil {
		return m.SuccessRateRequestVolume
	}
	return nil
}

func (m *OutlierDetection) GetSuccessRateStdevFactor() *types.UInt32Value {
	if m != nil {
		return m.SuccessRateStdevFactor
	}
	return nil
}

func (m *OutlierDetection) GetConsecutiveGatewayFailure() *types.UInt32Value {
	if m != nil {
		return m.ConsecutiveGatewayFailure
	}
	return nil
}

func (m *OutlierDetection) GetEnforcingConsecutiveGatewayFailure() *types.UInt32Value {
	if m != nil {
		return m.EnforcingConsecutiveGatewayFailure
	}
	return nil
}

func init() {
	proto.RegisterType((*OutlierDetection)(nil), "envoy.api.v2.cluster.OutlierDetection")
}

func init() {
	proto.RegisterFile("envoy/api/v2/cluster/outlier_detection.proto", fileDescriptor_56cd87362a3f00c9)
}

var fileDescriptor_56cd87362a3f00c9 = []byte{
	// 557 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xdd, 0x6a, 0x13, 0x41,
	0x14, 0xc7, 0xdd, 0xf4, 0x7b, 0x0a, 0x5a, 0x86, 0xd8, 0x6e, 0x5a, 0x09, 0x52, 0x10, 0xa4, 0xc8,
	0x2e, 0xa4, 0xf4, 0x01, 0x9a, 0x36, 0x55, 0x2f, 0xd4, 0x90, 0x68, 0x44, 0x54, 0x86, 0xc9, 0xe6,
	0x64, 0x1d, 0xd9, 0xdd, 0x59, 0xe7, 0x63, 0xbb, 0xf1, 0x89, 0xa4, 0x8f, 0xe0, 0x95, 0x97, 0x5e,
	0xfa, 0x08, 0x92, 0x3b, 0x9f, 0xc1, 0x1b, 0xd9, 0x9d, 0x7c, 0x6c, 0xd2, 0x80, 0xc9, 0xdd, 0xb0,
	0x73, 0x7e, 0xbf, 0xff, 0xd9, 0x99, 0xe1, 0xa0, 0x27, 0x10, 0x25, 0x7c, 0xe0, 0xd2, 0x98, 0xb9,
	0x49, 0xcd, 0xf5, 0x02, 0x2d, 0x15, 0x08, 0x97, 0x6b, 0x15, 0x30, 0x10, 0xa4, 0x07, 0x0a, 0x3c,
	0xc5, 0x78, 0xe4, 0xc4, 0x82, 0x2b, 0x8e, 0xcb, 0x79, 0xb5, 0x43, 0x63, 0xe6, 0x24, 0x35, 0x67,
	0x54, 0x7d, 0x58, 0xf5, 0x39, 0xf7, 0x03, 0x70, 0xf3, 0x9a, 0xae, 0xee, 0xbb, 0x3d, 0x2d, 0xe8,
	0x94, 0xba, 0xbd, 0x7f, 0x2d, 0x68, 0x1c, 0x83, 0x90, 0xa3, 0xfd, 0x83, 0x84, 0x06, 0xac, 0x47,
	0x15, 0xb8, 0xe3, 0xc5, 0x68, 0xa3, 0xec, 0x73, 0x9f, 0xe7, 0x4b, 0x37, 0x5b, 0x99, 0xaf, 0xc7,
	0x7f, 0xb7, 0xd0, 0xde, 0x2b, 0xd3, 0xe0, 0xe5, 0xb8, 0x3f, 0xdc, 0x40, 0xf7, 0x3c, 0x1e, 0x49,
	0xf0, 0xb4, 0x62, 0x09, 0x90, 0xb3, 0x34, 0xb5, 0xad, 0x87, 0xd6, 0xe3, 0xdd, 0xda, 0x03, 0xc7,
	0xa4, 0x3b, 0xe3, 0x74, 0xe7, 0xcd, 0xf3, 0x48, 0x9d, 0xd6, 0x3a, 0x34, 0xd0, 0xd0, 0xba, 0x5b,
	0x80, 0xce, 0xd2, 0x14, 0x9f, 0xa3, 0x6d, 0x16, 0x29, 0x10, 0x09, 0x0d, 0xec, 0x52, 0xce, 0x57,
	0x6e, 0xf1, 0x97, 0xa3, 0xbf, 0xab, 0xa3, 0xef, 0x7f, 0x7e, 0xac, 0x6d, 0xdc, 0x58, 0xa5, 0x93,
	0x3b, 0xad, 0x09, 0x86, 0xdb, 0x08, 0x77, 0xa9, 0x04, 0x02, 0x9f, 0x4d, 0x6b, 0x44, 0xb1, 0x10,
	0xec, 0xb5, 0x55, 0x64, 0x7b, 0x99, 0xa0, 0x31, 0xe2, 0x5f, 0xb3, 0x10, 0xf0, 0x3b, 0x54, 0x0e,
	0x69, 0x3a, 0x75, 0xc6, 0x20, 0x3c, 0x88, 0x94, 0xbd, 0xfe, 0xff, 0x7f, 0xac, 0xef, 0x64, 0xe6,
	0xf5, 0x93, 0x92, 0xdd, 0x6b, 0xe1, 0x90, 0xa6, 0x63, 0x6f, 0xd3, 0x28, 0xb0, 0x87, 0x2a, 0x10,
	0xf5, 0xb9, 0xf0, 0x58, 0xe4, 0x93, 0xf9, 0x33, 0xdc, 0x58, 0xcd, 0x7f, 0x30, 0x31, 0x5d, 0xcc,
	0x9e, 0xeb, 0x47, 0xb4, 0x3f, 0x0d, 0x91, 0xda, 0xf3, 0x40, 0x4a, 0x22, 0xa8, 0x02, 0x7b, 0x73,
	0xb5, 0x84, 0xf2, 0x44, 0xd3, 0x36, 0x96, 0x16, 0x55, 0xd9, 0xf1, 0x1c, 0x16, 0xa5, 0x24, 0x64,
	0x11, 0x0b, 0x75, 0x48, 0x3e, 0x71, 0xa9, 0xa4, 0xbd, 0xb5, 0xc4, 0x43, 0x38, 0x90, 0x53, 0xdd,
	0x0b, 0x43, 0x3f, 0xcb, 0x60, 0xfc, 0x1e, 0x1d, 0xcd, 0xa8, 0x05, 0x7c, 0xd1, 0x20, 0x15, 0x49,
	0x78, 0xa0, 0x43, 0xb0, 0xb7, 0x97, 0x70, 0xdb, 0x05, 0x77, 0xcb, 0xe0, 0x9d, 0x9c, 0xc6, 0x6f,
	0x51, 0x65, 0x46, 0x2e, 0x55, 0x0f, 0x12, 0xd2, 0xa7, 0x9e, 0xe2, 0xc2, 0xde, 0x59, 0x42, 0xbd,
	0x5f, 0x50, 0xb7, 0x33, 0xf8, 0x2a, 0x67, 0xf1, 0x07, 0x74, 0x54, 0xbc, 0x4a, 0x9f, 0x2a, 0xb8,
	0xa6, 0x03, 0xd2, 0xa7, 0x2c, 0xd0, 0x02, 0x6c, 0xb4, 0x84, 0xba, 0x52, 0x10, 0x3c, 0x35, 0xfc,
	0x95, 0xc1, 0xf1, 0x57, 0xf4, 0x68, 0xf1, 0x93, 0x99, 0xcf, 0xd9, 0x5d, 0xed, 0x72, 0x8f, 0x17,
	0x3d, 0x9f, 0xd9, 0xec, 0x7a, 0xff, 0xdb, 0xb0, 0x6a, 0xfd, 0x1c, 0x56, 0xad, 0x5f, 0xc3, 0xaa,
	0xf5, 0x7b, 0x58, 0xb5, 0xd0, 0x31, 0xe3, 0x4e, 0x3e, 0x97, 0x62, 0xc1, 0xd3, 0x81, 0xb3, 0x68,
	0x44, 0xd5, 0xef, 0xcf, 0x0f, 0x8c, 0x66, 0xd6, 0x4a, 0xd3, 0xba, 0x29, 0xed, 0x37, 0xf2, 0xfa,
	0xf3, 0x98, 0x39, 0x9d, 0x9a, 0x73, 0x61, 0xea, 0x5f, 0xb6, 0xbb, 0x9b, 0x79, 0xb3, 0xa7, 0xff,
	0x02, 0x00, 0x00, 0xff, 0xff, 0xde, 0xe2, 0x08, 0x28, 0x21, 0x05, 0x00, 0x00,
}

func (this *OutlierDetection) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*OutlierDetection)
	if !ok {
		that2, ok := that.(OutlierDetection)
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
	if !this.Consecutive_5Xx.Equal(that1.Consecutive_5Xx) {
		return false
	}
	if !this.Interval.Equal(that1.Interval) {
		return false
	}
	if !this.BaseEjectionTime.Equal(that1.BaseEjectionTime) {
		return false
	}
	if !this.MaxEjectionPercent.Equal(that1.MaxEjectionPercent) {
		return false
	}
	if !this.EnforcingConsecutive_5Xx.Equal(that1.EnforcingConsecutive_5Xx) {
		return false
	}
	if !this.EnforcingSuccessRate.Equal(that1.EnforcingSuccessRate) {
		return false
	}
	if !this.SuccessRateMinimumHosts.Equal(that1.SuccessRateMinimumHosts) {
		return false
	}
	if !this.SuccessRateRequestVolume.Equal(that1.SuccessRateRequestVolume) {
		return false
	}
	if !this.SuccessRateStdevFactor.Equal(that1.SuccessRateStdevFactor) {
		return false
	}
	if !this.ConsecutiveGatewayFailure.Equal(that1.ConsecutiveGatewayFailure) {
		return false
	}
	if !this.EnforcingConsecutiveGatewayFailure.Equal(that1.EnforcingConsecutiveGatewayFailure) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (m *OutlierDetection) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OutlierDetection) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Consecutive_5Xx != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.Consecutive_5Xx.Size()))
		n1, err := m.Consecutive_5Xx.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if m.Interval != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.Interval.Size()))
		n2, err := m.Interval.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if m.BaseEjectionTime != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.BaseEjectionTime.Size()))
		n3, err := m.BaseEjectionTime.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	if m.MaxEjectionPercent != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.MaxEjectionPercent.Size()))
		n4, err := m.MaxEjectionPercent.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	if m.EnforcingConsecutive_5Xx != nil {
		dAtA[i] = 0x2a
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.EnforcingConsecutive_5Xx.Size()))
		n5, err := m.EnforcingConsecutive_5Xx.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	if m.EnforcingSuccessRate != nil {
		dAtA[i] = 0x32
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.EnforcingSuccessRate.Size()))
		n6, err := m.EnforcingSuccessRate.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n6
	}
	if m.SuccessRateMinimumHosts != nil {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.SuccessRateMinimumHosts.Size()))
		n7, err := m.SuccessRateMinimumHosts.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n7
	}
	if m.SuccessRateRequestVolume != nil {
		dAtA[i] = 0x42
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.SuccessRateRequestVolume.Size()))
		n8, err := m.SuccessRateRequestVolume.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n8
	}
	if m.SuccessRateStdevFactor != nil {
		dAtA[i] = 0x4a
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.SuccessRateStdevFactor.Size()))
		n9, err := m.SuccessRateStdevFactor.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n9
	}
	if m.ConsecutiveGatewayFailure != nil {
		dAtA[i] = 0x52
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.ConsecutiveGatewayFailure.Size()))
		n10, err := m.ConsecutiveGatewayFailure.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n10
	}
	if m.EnforcingConsecutiveGatewayFailure != nil {
		dAtA[i] = 0x5a
		i++
		i = encodeVarintOutlierDetection(dAtA, i, uint64(m.EnforcingConsecutiveGatewayFailure.Size()))
		n11, err := m.EnforcingConsecutiveGatewayFailure.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n11
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintOutlierDetection(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *OutlierDetection) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Consecutive_5Xx != nil {
		l = m.Consecutive_5Xx.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.Interval != nil {
		l = m.Interval.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.BaseEjectionTime != nil {
		l = m.BaseEjectionTime.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.MaxEjectionPercent != nil {
		l = m.MaxEjectionPercent.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.EnforcingConsecutive_5Xx != nil {
		l = m.EnforcingConsecutive_5Xx.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.EnforcingSuccessRate != nil {
		l = m.EnforcingSuccessRate.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.SuccessRateMinimumHosts != nil {
		l = m.SuccessRateMinimumHosts.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.SuccessRateRequestVolume != nil {
		l = m.SuccessRateRequestVolume.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.SuccessRateStdevFactor != nil {
		l = m.SuccessRateStdevFactor.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.ConsecutiveGatewayFailure != nil {
		l = m.ConsecutiveGatewayFailure.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.EnforcingConsecutiveGatewayFailure != nil {
		l = m.EnforcingConsecutiveGatewayFailure.Size()
		n += 1 + l + sovOutlierDetection(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovOutlierDetection(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozOutlierDetection(x uint64) (n int) {
	return sovOutlierDetection(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *OutlierDetection) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOutlierDetection
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
			return fmt.Errorf("proto: OutlierDetection: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OutlierDetection: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Consecutive_5Xx", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Consecutive_5Xx == nil {
				m.Consecutive_5Xx = &types.UInt32Value{}
			}
			if err := m.Consecutive_5Xx.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Interval", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Interval == nil {
				m.Interval = &types.Duration{}
			}
			if err := m.Interval.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseEjectionTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BaseEjectionTime == nil {
				m.BaseEjectionTime = &types.Duration{}
			}
			if err := m.BaseEjectionTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxEjectionPercent", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MaxEjectionPercent == nil {
				m.MaxEjectionPercent = &types.UInt32Value{}
			}
			if err := m.MaxEjectionPercent.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnforcingConsecutive_5Xx", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.EnforcingConsecutive_5Xx == nil {
				m.EnforcingConsecutive_5Xx = &types.UInt32Value{}
			}
			if err := m.EnforcingConsecutive_5Xx.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnforcingSuccessRate", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.EnforcingSuccessRate == nil {
				m.EnforcingSuccessRate = &types.UInt32Value{}
			}
			if err := m.EnforcingSuccessRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SuccessRateMinimumHosts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SuccessRateMinimumHosts == nil {
				m.SuccessRateMinimumHosts = &types.UInt32Value{}
			}
			if err := m.SuccessRateMinimumHosts.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SuccessRateRequestVolume", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SuccessRateRequestVolume == nil {
				m.SuccessRateRequestVolume = &types.UInt32Value{}
			}
			if err := m.SuccessRateRequestVolume.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SuccessRateStdevFactor", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SuccessRateStdevFactor == nil {
				m.SuccessRateStdevFactor = &types.UInt32Value{}
			}
			if err := m.SuccessRateStdevFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsecutiveGatewayFailure", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ConsecutiveGatewayFailure == nil {
				m.ConsecutiveGatewayFailure = &types.UInt32Value{}
			}
			if err := m.ConsecutiveGatewayFailure.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnforcingConsecutiveGatewayFailure", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOutlierDetection
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
				return ErrInvalidLengthOutlierDetection
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.EnforcingConsecutiveGatewayFailure == nil {
				m.EnforcingConsecutiveGatewayFailure = &types.UInt32Value{}
			}
			if err := m.EnforcingConsecutiveGatewayFailure.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOutlierDetection(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOutlierDetection
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOutlierDetection
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
func skipOutlierDetection(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOutlierDetection
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
					return 0, ErrIntOverflowOutlierDetection
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
					return 0, ErrIntOverflowOutlierDetection
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
				return 0, ErrInvalidLengthOutlierDetection
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthOutlierDetection
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowOutlierDetection
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
				next, err := skipOutlierDetection(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthOutlierDetection
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
	ErrInvalidLengthOutlierDetection = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOutlierDetection   = fmt.Errorf("proto: integer overflow")
)
