// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/admin/v2alpha/clusters.proto

/*
	Package envoy_admin_v2alpha is a generated protocol buffer package.

	It is generated from these files:
		envoy/admin/v2alpha/clusters.proto
		envoy/admin/v2alpha/config_dump.proto
		envoy/admin/v2alpha/metrics.proto

	It has these top-level messages:
		Clusters
		ClusterStatus
		HostStatus
		HostHealthStatus
		ConfigDump
		BootstrapConfigDump
		ListenersConfigDump
		ClustersConfigDump
		RoutesConfigDump
		SimpleMetric
*/
package envoy_admin_v2alpha

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import envoy_api_v2_core1 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
import envoy_api_v2_core2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
import envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"

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

// Admin endpoint uses this wrapper for `/clusters` to display cluster status information.
// See :ref:`/clusters <operations_admin_interface_clusters>` for more information.
type Clusters struct {
	// Mapping from cluster name to each cluster's status.
	ClusterStatuses []*ClusterStatus `protobuf:"bytes,1,rep,name=cluster_statuses,json=clusterStatuses" json:"cluster_statuses,omitempty"`
}

func (m *Clusters) Reset()                    { *m = Clusters{} }
func (m *Clusters) String() string            { return proto.CompactTextString(m) }
func (*Clusters) ProtoMessage()               {}
func (*Clusters) Descriptor() ([]byte, []int) { return fileDescriptorClusters, []int{0} }

func (m *Clusters) GetClusterStatuses() []*ClusterStatus {
	if m != nil {
		return m.ClusterStatuses
	}
	return nil
}

// Details an individual cluster's current status.
type ClusterStatus struct {
	// Name of the cluster.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Denotes whether this cluster was added via API or configured statically.
	AddedViaApi bool `protobuf:"varint,2,opt,name=added_via_api,json=addedViaApi,proto3" json:"added_via_api,omitempty"`
	// The success rate threshold used in the last interval. The threshold is used to eject hosts
	// based on their success rate. See
	// :ref:`Cluster outlier detection <arch_overview_outlier_detection>` statistics
	//
	// Note: this field may be omitted in any of the three following cases:
	// 1. There were not enough hosts with enough request volume to proceed with success rate based outlier ejection.
	// 2. The threshold is computed to be < 0 because a negative value implies that there was no threshold for that
	// interval.
	// 3. Outlier detection is not enabled for this cluster.
	SuccessRateEjectionThreshold *envoy_type.Percent `protobuf:"bytes,3,opt,name=success_rate_ejection_threshold,json=successRateEjectionThreshold" json:"success_rate_ejection_threshold,omitempty"`
	// Mapping from host address to the host's current status.
	HostStatuses []*HostStatus `protobuf:"bytes,4,rep,name=host_statuses,json=hostStatuses" json:"host_statuses,omitempty"`
}

func (m *ClusterStatus) Reset()                    { *m = ClusterStatus{} }
func (m *ClusterStatus) String() string            { return proto.CompactTextString(m) }
func (*ClusterStatus) ProtoMessage()               {}
func (*ClusterStatus) Descriptor() ([]byte, []int) { return fileDescriptorClusters, []int{1} }

func (m *ClusterStatus) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ClusterStatus) GetAddedViaApi() bool {
	if m != nil {
		return m.AddedViaApi
	}
	return false
}

func (m *ClusterStatus) GetSuccessRateEjectionThreshold() *envoy_type.Percent {
	if m != nil {
		return m.SuccessRateEjectionThreshold
	}
	return nil
}

func (m *ClusterStatus) GetHostStatuses() []*HostStatus {
	if m != nil {
		return m.HostStatuses
	}
	return nil
}

// Current state of a particular host.
type HostStatus struct {
	// Address of this host.
	Address *envoy_api_v2_core1.Address `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
	// Mapping from the name of the statistic to the current value.
	Stats map[string]*SimpleMetric `protobuf:"bytes,2,rep,name=stats" json:"stats,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	// The host's current health status.
	HealthStatus *HostHealthStatus `protobuf:"bytes,3,opt,name=health_status,json=healthStatus" json:"health_status,omitempty"`
	// Request success rate for this host over the last calculated interval.
	//
	// Note: the message will not be present if host did not have enough request volume to calculate
	// success rate or the cluster did not have enough hosts to run through success rate outlier
	// ejection.
	SuccessRate *envoy_type.Percent `protobuf:"bytes,4,opt,name=success_rate,json=successRate" json:"success_rate,omitempty"`
}

func (m *HostStatus) Reset()                    { *m = HostStatus{} }
func (m *HostStatus) String() string            { return proto.CompactTextString(m) }
func (*HostStatus) ProtoMessage()               {}
func (*HostStatus) Descriptor() ([]byte, []int) { return fileDescriptorClusters, []int{2} }

func (m *HostStatus) GetAddress() *envoy_api_v2_core1.Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *HostStatus) GetStats() map[string]*SimpleMetric {
	if m != nil {
		return m.Stats
	}
	return nil
}

func (m *HostStatus) GetHealthStatus() *HostHealthStatus {
	if m != nil {
		return m.HealthStatus
	}
	return nil
}

func (m *HostStatus) GetSuccessRate() *envoy_type.Percent {
	if m != nil {
		return m.SuccessRate
	}
	return nil
}

// Health status for a host.
type HostHealthStatus struct {
	// The host is currently failing active health checks.
	FailedActiveHealthCheck bool `protobuf:"varint,1,opt,name=failed_active_health_check,json=failedActiveHealthCheck,proto3" json:"failed_active_health_check,omitempty"`
	// The host is currently considered an outlier and has been ejected.
	FailedOutlierCheck bool `protobuf:"varint,2,opt,name=failed_outlier_check,json=failedOutlierCheck,proto3" json:"failed_outlier_check,omitempty"`
	// Health status as reported by EDS. Note: only HEALTHY and UNHEALTHY are currently supported here.
	// TODO(mrice32): pipe through remaining EDS health status possibilities.
	EdsHealthStatus envoy_api_v2_core2.HealthStatus `protobuf:"varint,3,opt,name=eds_health_status,json=edsHealthStatus,proto3,enum=envoy.api.v2.core.HealthStatus" json:"eds_health_status,omitempty"`
}

func (m *HostHealthStatus) Reset()                    { *m = HostHealthStatus{} }
func (m *HostHealthStatus) String() string            { return proto.CompactTextString(m) }
func (*HostHealthStatus) ProtoMessage()               {}
func (*HostHealthStatus) Descriptor() ([]byte, []int) { return fileDescriptorClusters, []int{3} }

func (m *HostHealthStatus) GetFailedActiveHealthCheck() bool {
	if m != nil {
		return m.FailedActiveHealthCheck
	}
	return false
}

func (m *HostHealthStatus) GetFailedOutlierCheck() bool {
	if m != nil {
		return m.FailedOutlierCheck
	}
	return false
}

func (m *HostHealthStatus) GetEdsHealthStatus() envoy_api_v2_core2.HealthStatus {
	if m != nil {
		return m.EdsHealthStatus
	}
	return envoy_api_v2_core2.HealthStatus_UNKNOWN
}

func init() {
	proto.RegisterType((*Clusters)(nil), "envoy.admin.v2alpha.Clusters")
	proto.RegisterType((*ClusterStatus)(nil), "envoy.admin.v2alpha.ClusterStatus")
	proto.RegisterType((*HostStatus)(nil), "envoy.admin.v2alpha.HostStatus")
	proto.RegisterType((*HostHealthStatus)(nil), "envoy.admin.v2alpha.HostHealthStatus")
}
func (m *Clusters) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Clusters) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ClusterStatuses) > 0 {
		for _, msg := range m.ClusterStatuses {
			dAtA[i] = 0xa
			i++
			i = encodeVarintClusters(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *ClusterStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClusterStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintClusters(dAtA, i, uint64(len(m.Name)))
		i += copy(dAtA[i:], m.Name)
	}
	if m.AddedViaApi {
		dAtA[i] = 0x10
		i++
		if m.AddedViaApi {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.SuccessRateEjectionThreshold != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintClusters(dAtA, i, uint64(m.SuccessRateEjectionThreshold.Size()))
		n1, err := m.SuccessRateEjectionThreshold.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.HostStatuses) > 0 {
		for _, msg := range m.HostStatuses {
			dAtA[i] = 0x22
			i++
			i = encodeVarintClusters(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *HostStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HostStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Address != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintClusters(dAtA, i, uint64(m.Address.Size()))
		n2, err := m.Address.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	if len(m.Stats) > 0 {
		for k, _ := range m.Stats {
			dAtA[i] = 0x12
			i++
			v := m.Stats[k]
			msgSize := 0
			if v != nil {
				msgSize = v.Size()
				msgSize += 1 + sovClusters(uint64(msgSize))
			}
			mapSize := 1 + len(k) + sovClusters(uint64(len(k))) + msgSize
			i = encodeVarintClusters(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintClusters(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			if v != nil {
				dAtA[i] = 0x12
				i++
				i = encodeVarintClusters(dAtA, i, uint64(v.Size()))
				n3, err := v.MarshalTo(dAtA[i:])
				if err != nil {
					return 0, err
				}
				i += n3
			}
		}
	}
	if m.HealthStatus != nil {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintClusters(dAtA, i, uint64(m.HealthStatus.Size()))
		n4, err := m.HealthStatus.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	if m.SuccessRate != nil {
		dAtA[i] = 0x22
		i++
		i = encodeVarintClusters(dAtA, i, uint64(m.SuccessRate.Size()))
		n5, err := m.SuccessRate.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n5
	}
	return i, nil
}

func (m *HostHealthStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HostHealthStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.FailedActiveHealthCheck {
		dAtA[i] = 0x8
		i++
		if m.FailedActiveHealthCheck {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.FailedOutlierCheck {
		dAtA[i] = 0x10
		i++
		if m.FailedOutlierCheck {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i++
	}
	if m.EdsHealthStatus != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintClusters(dAtA, i, uint64(m.EdsHealthStatus))
	}
	return i, nil
}

func encodeVarintClusters(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Clusters) Size() (n int) {
	var l int
	_ = l
	if len(m.ClusterStatuses) > 0 {
		for _, e := range m.ClusterStatuses {
			l = e.Size()
			n += 1 + l + sovClusters(uint64(l))
		}
	}
	return n
}

func (m *ClusterStatus) Size() (n int) {
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovClusters(uint64(l))
	}
	if m.AddedViaApi {
		n += 2
	}
	if m.SuccessRateEjectionThreshold != nil {
		l = m.SuccessRateEjectionThreshold.Size()
		n += 1 + l + sovClusters(uint64(l))
	}
	if len(m.HostStatuses) > 0 {
		for _, e := range m.HostStatuses {
			l = e.Size()
			n += 1 + l + sovClusters(uint64(l))
		}
	}
	return n
}

func (m *HostStatus) Size() (n int) {
	var l int
	_ = l
	if m.Address != nil {
		l = m.Address.Size()
		n += 1 + l + sovClusters(uint64(l))
	}
	if len(m.Stats) > 0 {
		for k, v := range m.Stats {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovClusters(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovClusters(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovClusters(uint64(mapEntrySize))
		}
	}
	if m.HealthStatus != nil {
		l = m.HealthStatus.Size()
		n += 1 + l + sovClusters(uint64(l))
	}
	if m.SuccessRate != nil {
		l = m.SuccessRate.Size()
		n += 1 + l + sovClusters(uint64(l))
	}
	return n
}

func (m *HostHealthStatus) Size() (n int) {
	var l int
	_ = l
	if m.FailedActiveHealthCheck {
		n += 2
	}
	if m.FailedOutlierCheck {
		n += 2
	}
	if m.EdsHealthStatus != 0 {
		n += 1 + sovClusters(uint64(m.EdsHealthStatus))
	}
	return n
}

func sovClusters(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozClusters(x uint64) (n int) {
	return sovClusters(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Clusters) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClusters
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
			return fmt.Errorf("proto: Clusters: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Clusters: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterStatuses", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterStatuses = append(m.ClusterStatuses, &ClusterStatus{})
			if err := m.ClusterStatuses[len(m.ClusterStatuses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipClusters(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthClusters
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
func (m *ClusterStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClusters
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
			return fmt.Errorf("proto: ClusterStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClusterStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddedViaApi", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.AddedViaApi = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SuccessRateEjectionThreshold", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SuccessRateEjectionThreshold == nil {
				m.SuccessRateEjectionThreshold = &envoy_type.Percent{}
			}
			if err := m.SuccessRateEjectionThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HostStatuses", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HostStatuses = append(m.HostStatuses, &HostStatus{})
			if err := m.HostStatuses[len(m.HostStatuses)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipClusters(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthClusters
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
func (m *HostStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClusters
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
			return fmt.Errorf("proto: HostStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HostStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Address == nil {
				m.Address = &envoy_api_v2_core1.Address{}
			}
			if err := m.Address.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stats", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Stats == nil {
				m.Stats = make(map[string]*SimpleMetric)
			}
			var mapkey string
			var mapvalue *SimpleMetric
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowClusters
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
							return ErrIntOverflowClusters
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
						return ErrInvalidLengthClusters
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
							return ErrIntOverflowClusters
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
						return ErrInvalidLengthClusters
					}
					postmsgIndex := iNdEx + mapmsglen
					if mapmsglen < 0 {
						return ErrInvalidLengthClusters
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &SimpleMetric{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipClusters(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthClusters
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Stats[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HealthStatus", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.HealthStatus == nil {
				m.HealthStatus = &HostHealthStatus{}
			}
			if err := m.HealthStatus.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SuccessRate", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
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
				return ErrInvalidLengthClusters
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SuccessRate == nil {
				m.SuccessRate = &envoy_type.Percent{}
			}
			if err := m.SuccessRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipClusters(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthClusters
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
func (m *HostHealthStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowClusters
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
			return fmt.Errorf("proto: HostHealthStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HostHealthStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FailedActiveHealthCheck", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FailedActiveHealthCheck = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FailedOutlierCheck", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.FailedOutlierCheck = bool(v != 0)
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EdsHealthStatus", wireType)
			}
			m.EdsHealthStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowClusters
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EdsHealthStatus |= (envoy_api_v2_core2.HealthStatus(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipClusters(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthClusters
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
func skipClusters(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowClusters
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
					return 0, ErrIntOverflowClusters
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
					return 0, ErrIntOverflowClusters
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
				return 0, ErrInvalidLengthClusters
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowClusters
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
				next, err := skipClusters(dAtA[start:])
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
	ErrInvalidLengthClusters = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowClusters   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("envoy/admin/v2alpha/clusters.proto", fileDescriptorClusters) }

var fileDescriptorClusters = []byte{
	// 549 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xd1, 0x8a, 0xd3, 0x40,
	0x14, 0x86, 0x49, 0xdb, 0xd5, 0x7a, 0xd2, 0xba, 0x75, 0x56, 0x30, 0x14, 0x69, 0xbb, 0x41, 0xa1,
	0x78, 0x91, 0x48, 0x14, 0x15, 0xbd, 0xb1, 0xae, 0x0b, 0x8b, 0xb2, 0x28, 0x59, 0x11, 0xd4, 0x8b,
	0x30, 0x4e, 0x8e, 0x64, 0xdc, 0x34, 0x09, 0x99, 0x69, 0xa0, 0x6f, 0xe8, 0x95, 0xf8, 0x08, 0x52,
	0xf0, 0xc2, 0xb7, 0x90, 0xcc, 0x4c, 0xb7, 0x59, 0xb7, 0xcb, 0x5e, 0x75, 0x32, 0xe7, 0xfb, 0xa7,
	0xff, 0x7f, 0xce, 0x0c, 0xb8, 0x98, 0x55, 0xf9, 0xd2, 0xa7, 0xf1, 0x9c, 0x67, 0x7e, 0x15, 0xd0,
	0xb4, 0x48, 0xa8, 0xcf, 0xd2, 0x85, 0x90, 0x58, 0x0a, 0xaf, 0x28, 0x73, 0x99, 0x93, 0x3d, 0xc5,
	0x78, 0x8a, 0xf1, 0x0c, 0x33, 0xdc, 0xdf, 0x26, 0x9c, 0xa3, 0x2c, 0x39, 0x33, 0xba, 0xe1, 0xd8,
	0x20, 0x05, 0xf7, 0xab, 0xc0, 0x67, 0x79, 0x89, 0x3e, 0x8d, 0xe3, 0x12, 0xc5, 0x1a, 0xb8, 0x77,
	0x11, 0x48, 0x90, 0xa6, 0x32, 0x89, 0x58, 0x82, 0xec, 0xd4, 0x50, 0x8e, 0xa6, 0xe4, 0xb2, 0x40,
	0xbf, 0xc0, 0x92, 0x61, 0x26, 0x75, 0xc5, 0xfd, 0x04, 0xdd, 0x03, 0x63, 0x95, 0x1c, 0xc3, 0xc0,
	0xd8, 0x8e, 0x84, 0xa4, 0x72, 0x21, 0x50, 0x38, 0xd6, 0xa4, 0x3d, 0xb5, 0x03, 0xd7, 0xdb, 0xe2,
	0xdf, 0x33, 0xc2, 0x13, 0xc5, 0x86, 0xbb, 0xac, 0xf9, 0x89, 0xc2, 0xfd, 0x6b, 0x41, 0xff, 0x1c,
	0x42, 0x08, 0x74, 0x32, 0x3a, 0x47, 0xc7, 0x9a, 0x58, 0xd3, 0x1b, 0xa1, 0x5a, 0x13, 0x17, 0xfa,
	0x34, 0x8e, 0x31, 0x8e, 0x2a, 0x4e, 0x23, 0x5a, 0x70, 0xa7, 0x35, 0xb1, 0xa6, 0xdd, 0xd0, 0x56,
	0x9b, 0x1f, 0x39, 0x9d, 0x15, 0x9c, 0x7c, 0x86, 0xb1, 0x58, 0x30, 0x86, 0x42, 0x44, 0x25, 0x95,
	0x18, 0xe1, 0x77, 0x64, 0x92, 0xe7, 0x59, 0x24, 0x93, 0x12, 0x45, 0x92, 0xa7, 0xb1, 0xd3, 0x9e,
	0x58, 0x53, 0x3b, 0xd8, 0x33, 0x3e, 0xeb, 0xa0, 0xde, 0x7b, 0x1d, 0x34, 0xbc, 0x6b, 0xb4, 0x21,
	0x95, 0x78, 0x68, 0x94, 0x1f, 0xd6, 0x42, 0xf2, 0x1a, 0xfa, 0x49, 0x2e, 0xe4, 0x26, 0x71, 0x47,
	0x25, 0x1e, 0x6f, 0x4d, 0x7c, 0x94, 0x0b, 0x69, 0xe2, 0xf6, 0x92, 0xb3, 0x35, 0x0a, 0xf7, 0x4f,
	0x0b, 0x60, 0x53, 0x24, 0x8f, 0xe1, 0xba, 0x19, 0x93, 0xca, 0x6a, 0x07, 0xc3, 0xf5, 0x71, 0x05,
	0xf7, 0xaa, 0xc0, 0xab, 0xe7, 0xe4, 0xcd, 0x34, 0x11, 0xae, 0x51, 0xf2, 0x12, 0x76, 0x6a, 0x17,
	0xc2, 0x69, 0x29, 0x0b, 0x0f, 0xae, 0xb0, 0xe0, 0xd5, 0x3f, 0xe2, 0x30, 0x93, 0xe5, 0x32, 0xd4,
	0x42, 0xf2, 0x06, 0xfa, 0x66, 0xfa, 0x3a, 0x8e, 0x69, 0xcb, 0xfd, 0x4b, 0x4f, 0x3a, 0x52, 0xf4,
	0x59, 0xa4, 0xc6, 0x17, 0x79, 0x02, 0xbd, 0x66, 0xd3, 0x9d, 0xce, 0xe5, 0x1d, 0xb6, 0x1b, 0x1d,
	0x1e, 0x7e, 0x01, 0xd8, 0x18, 0x23, 0x03, 0x68, 0x9f, 0xe2, 0xd2, 0x4c, 0xbc, 0x5e, 0x92, 0xa7,
	0xb0, 0x53, 0xd1, 0x74, 0x81, 0x6a, 0xd0, 0x76, 0xb0, 0xbf, 0xd5, 0xdb, 0x09, 0x9f, 0x17, 0x29,
	0x1e, 0xab, 0xb7, 0x10, 0x6a, 0xfe, 0x79, 0xeb, 0x99, 0xe5, 0xfe, 0xb4, 0x60, 0xf0, 0xbf, 0x6f,
	0xf2, 0x02, 0x86, 0xdf, 0x28, 0x4f, 0x31, 0x8e, 0x28, 0x93, 0xbc, 0xc2, 0xa8, 0xf9, 0x02, 0xd4,
	0x5f, 0x77, 0xc3, 0x3b, 0x9a, 0x98, 0x29, 0x40, 0xab, 0x0f, 0xea, 0x32, 0x79, 0x08, 0xb7, 0x8d,
	0x38, 0x5f, 0xc8, 0x94, 0x63, 0x69, 0x64, 0xfa, 0x1a, 0x12, 0x5d, 0x7b, 0xa7, 0x4b, 0x5a, 0xf1,
	0x16, 0x6e, 0x61, 0x2c, 0xa2, 0x8b, 0x8d, 0xbe, 0xb9, 0xb9, 0x35, 0x8d, 0x31, 0x9f, 0x6b, 0xf1,
	0x2e, 0xc6, 0xa2, 0xb9, 0xf1, 0xaa, 0xf7, 0x63, 0x35, 0xb2, 0x7e, 0xad, 0x46, 0xd6, 0xef, 0xd5,
	0xc8, 0xfa, 0x7a, 0x4d, 0x3d, 0xca, 0x47, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x70, 0x94, 0xfb,
	0x8a, 0x53, 0x04, 0x00, 0x00,
}
