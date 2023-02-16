// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: workloadapi/workload.proto

package workloadapi

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WorkloadStatus int32

const (
	// Workload is healthy and ready to serve traffic.
	WorkloadStatus_HEALTHY WorkloadStatus = 0
	// Workload is unhealthy and NOT ready to serve traffic.
	WorkloadStatus_UNHEALTHY WorkloadStatus = 1
)

// Enum value maps for WorkloadStatus.
var (
	WorkloadStatus_name = map[int32]string{
		0: "HEALTHY",
		1: "UNHEALTHY",
	}
	WorkloadStatus_value = map[string]int32{
		"HEALTHY":   0,
		"UNHEALTHY": 1,
	}
)

func (x WorkloadStatus) Enum() *WorkloadStatus {
	p := new(WorkloadStatus)
	*p = x
	return p
}

func (x WorkloadStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WorkloadStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_workloadapi_workload_proto_enumTypes[0].Descriptor()
}

func (WorkloadStatus) Type() protoreflect.EnumType {
	return &file_workloadapi_workload_proto_enumTypes[0]
}

func (x WorkloadStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WorkloadStatus.Descriptor instead.
func (WorkloadStatus) EnumDescriptor() ([]byte, []int) {
	return file_workloadapi_workload_proto_rawDescGZIP(), []int{0}
}

type WorkloadType int32

const (
	WorkloadType_DEPLOYMENT WorkloadType = 0
	WorkloadType_CRONJOB    WorkloadType = 1
	WorkloadType_POD        WorkloadType = 2
	WorkloadType_JOB        WorkloadType = 3
)

// Enum value maps for WorkloadType.
var (
	WorkloadType_name = map[int32]string{
		0: "DEPLOYMENT",
		1: "CRONJOB",
		2: "POD",
		3: "JOB",
	}
	WorkloadType_value = map[string]int32{
		"DEPLOYMENT": 0,
		"CRONJOB":    1,
		"POD":        2,
		"JOB":        3,
	}
)

func (x WorkloadType) Enum() *WorkloadType {
	p := new(WorkloadType)
	*p = x
	return p
}

func (x WorkloadType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WorkloadType) Descriptor() protoreflect.EnumDescriptor {
	return file_workloadapi_workload_proto_enumTypes[1].Descriptor()
}

func (WorkloadType) Type() protoreflect.EnumType {
	return &file_workloadapi_workload_proto_enumTypes[1]
}

func (x WorkloadType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WorkloadType.Descriptor instead.
func (WorkloadType) EnumDescriptor() ([]byte, []int) {
	return file_workloadapi_workload_proto_rawDescGZIP(), []int{1}
}

type Protocol int32

const (
	// DIRECT means requests should be forwarded as-is.
	Protocol_DIRECT Protocol = 0
	// HTTP means requests should be tunneled over HTTP.
	// This does not dictate HTTP/1.1 vs HTTP/2; ALPN should be used for that purpose.
	Protocol_HTTP Protocol = 1
)

// Enum value maps for Protocol.
var (
	Protocol_name = map[int32]string{
		0: "DIRECT",
		1: "HTTP",
	}
	Protocol_value = map[string]int32{
		"DIRECT": 0,
		"HTTP":   1,
	}
)

func (x Protocol) Enum() *Protocol {
	p := new(Protocol)
	*p = x
	return p
}

func (x Protocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Protocol) Descriptor() protoreflect.EnumDescriptor {
	return file_workloadapi_workload_proto_enumTypes[2].Descriptor()
}

func (Protocol) Type() protoreflect.EnumType {
	return &file_workloadapi_workload_proto_enumTypes[2]
}

func (x Protocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Protocol.Descriptor instead.
func (Protocol) EnumDescriptor() ([]byte, []int) {
	return file_workloadapi_workload_proto_rawDescGZIP(), []int{2}
}

type Workload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name represents the name for the workload.
	// For Kubernetes, this is the pod name.
	// This is just for debugging and may be elided as an optimization.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Namespace represents the namespace for the workload.
	// This is just for debugging and may be elided as an optimization.
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// Address represents the IPv4/IPv6 address for the workload.
	// This should be globally unique.
	// This should not have a port number.
	// TODO: Add network as discriminator
	Address []byte `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	// Network represents the network this workload is on. This may be elided for the default network.
	// A (network,address) pair makeup a unique key for a workload *at a point in time*.
	Network string `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
	// Protocol that should be used to connect to this workload.
	Protocol Protocol `protobuf:"varint,5,opt,name=protocol,proto3,enum=istio.workload.Protocol" json:"protocol,omitempty"`
	// The SPIFFE identity of the workload. The identity is joined to form spiffe://<trust_domain>/ns/<namespace>/sa/<service_account>.
	// TrustDomain of the workload. May be elided if this is the mesh wide default (typically cluster.local)
	TrustDomain string `protobuf:"bytes,6,opt,name=trust_domain,json=trustDomain,proto3" json:"trust_domain,omitempty"`
	// ServiceAccount of the workload. May be elided if this is "default"
	ServiceAccount string `protobuf:"bytes,7,opt,name=service_account,json=serviceAccount,proto3" json:"service_account,omitempty"`
	// If present, the waypoint proxy for this workload.
	WaypointAddresses [][]byte `protobuf:"bytes,8,rep,name=waypoint_addresses,json=waypointAddresses,proto3" json:"waypoint_addresses,omitempty"`
	// Name of the node the workload runs on
	Node string `protobuf:"bytes,9,opt,name=node,proto3" json:"node,omitempty"`
	// CanonicalName for the workload. Used for telemetry.
	CanonicalName string `protobuf:"bytes,10,opt,name=canonical_name,json=canonicalName,proto3" json:"canonical_name,omitempty"`
	// CanonicalRevision for the workload. Used for telemetry.
	CanonicalRevision string `protobuf:"bytes,11,opt,name=canonical_revision,json=canonicalRevision,proto3" json:"canonical_revision,omitempty"`
	// WorkloadType represents the type of the workload. Used for telemetry.
	WorkloadType WorkloadType `protobuf:"varint,12,opt,name=workload_type,json=workloadType,proto3,enum=istio.workload.WorkloadType" json:"workload_type,omitempty"`
	// WorkloadName represents the name for the workload (of type WorkloadType). Used for telemetry.
	WorkloadName string `protobuf:"bytes,13,opt,name=workload_name,json=workloadName,proto3" json:"workload_name,omitempty"`
	// If set, indicates this workload directly speaks HBONE, and we should forward HBONE requests as-is.
	NativeHbone bool `protobuf:"varint,14,opt,name=native_hbone,json=nativeHbone,proto3" json:"native_hbone,omitempty"`
	// Virtual IPs defines a set of virtual IP addresses the workload can be reached at.
	// Typically these represent Service ClusterIPs.
	// The key is an IP address.
	VirtualIps map[string]*PortList `protobuf:"bytes,15,rep,name=virtual_ips,json=virtualIps,proto3" json:"virtual_ips,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// A list of authorization policies applicable to this workload.
	// NOTE: this *only* includes Selector based policies. Namespace and global polices
	// are returned out of band.
	AuthorizationPolicies []string       `protobuf:"bytes,16,rep,name=authorization_policies,json=authorizationPolicies,proto3" json:"authorization_policies,omitempty"`
	Status                WorkloadStatus `protobuf:"varint,17,opt,name=status,proto3,enum=istio.workload.WorkloadStatus" json:"status,omitempty"`
}

func (x *Workload) Reset() {
	*x = Workload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workloadapi_workload_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Workload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workload) ProtoMessage() {}

func (x *Workload) ProtoReflect() protoreflect.Message {
	mi := &file_workloadapi_workload_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workload.ProtoReflect.Descriptor instead.
func (*Workload) Descriptor() ([]byte, []int) {
	return file_workloadapi_workload_proto_rawDescGZIP(), []int{0}
}

func (x *Workload) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Workload) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Workload) GetAddress() []byte {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *Workload) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Workload) GetProtocol() Protocol {
	if x != nil {
		return x.Protocol
	}
	return Protocol_DIRECT
}

func (x *Workload) GetTrustDomain() string {
	if x != nil {
		return x.TrustDomain
	}
	return ""
}

func (x *Workload) GetServiceAccount() string {
	if x != nil {
		return x.ServiceAccount
	}
	return ""
}

func (x *Workload) GetWaypointAddresses() [][]byte {
	if x != nil {
		return x.WaypointAddresses
	}
	return nil
}

func (x *Workload) GetNode() string {
	if x != nil {
		return x.Node
	}
	return ""
}

func (x *Workload) GetCanonicalName() string {
	if x != nil {
		return x.CanonicalName
	}
	return ""
}

func (x *Workload) GetCanonicalRevision() string {
	if x != nil {
		return x.CanonicalRevision
	}
	return ""
}

func (x *Workload) GetWorkloadType() WorkloadType {
	if x != nil {
		return x.WorkloadType
	}
	return WorkloadType_DEPLOYMENT
}

func (x *Workload) GetWorkloadName() string {
	if x != nil {
		return x.WorkloadName
	}
	return ""
}

func (x *Workload) GetNativeHbone() bool {
	if x != nil {
		return x.NativeHbone
	}
	return false
}

func (x *Workload) GetVirtualIps() map[string]*PortList {
	if x != nil {
		return x.VirtualIps
	}
	return nil
}

func (x *Workload) GetAuthorizationPolicies() []string {
	if x != nil {
		return x.AuthorizationPolicies
	}
	return nil
}

func (x *Workload) GetStatus() WorkloadStatus {
	if x != nil {
		return x.Status
	}
	return WorkloadStatus_HEALTHY
}

// PorList represents the ports for a service
type PortList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ports []*Port `protobuf:"bytes,1,rep,name=ports,proto3" json:"ports,omitempty"`
}

func (x *PortList) Reset() {
	*x = PortList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workloadapi_workload_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PortList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PortList) ProtoMessage() {}

func (x *PortList) ProtoReflect() protoreflect.Message {
	mi := &file_workloadapi_workload_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PortList.ProtoReflect.Descriptor instead.
func (*PortList) Descriptor() ([]byte, []int) {
	return file_workloadapi_workload_proto_rawDescGZIP(), []int{1}
}

func (x *PortList) GetPorts() []*Port {
	if x != nil {
		return x.Ports
	}
	return nil
}

type Port struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Port the service is reached at (frontend).
	ServicePort uint32 `protobuf:"varint,1,opt,name=service_port,json=servicePort,proto3" json:"service_port,omitempty"`
	// Port the service forwards to (backend).
	TargetPort uint32 `protobuf:"varint,2,opt,name=target_port,json=targetPort,proto3" json:"target_port,omitempty"`
}

func (x *Port) Reset() {
	*x = Port{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workloadapi_workload_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Port) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Port) ProtoMessage() {}

func (x *Port) ProtoReflect() protoreflect.Message {
	mi := &file_workloadapi_workload_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Port.ProtoReflect.Descriptor instead.
func (*Port) Descriptor() ([]byte, []int) {
	return file_workloadapi_workload_proto_rawDescGZIP(), []int{2}
}

func (x *Port) GetServicePort() uint32 {
	if x != nil {
		return x.ServicePort
	}
	return 0
}

func (x *Port) GetTargetPort() uint32 {
	if x != nil {
		return x.TargetPort
	}
	return 0
}

var File_workloadapi_workload_proto protoreflect.FileDescriptor

var file_workloadapi_workload_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x61, 0x70, 0x69, 0x2f, 0x77, 0x6f,
	0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x69, 0x73,
	0x74, 0x69, 0x6f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0xa9, 0x06, 0x0a,
	0x08, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x12,
	0x34, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x18, 0x2e, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52, 0x08, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x74, 0x72, 0x75, 0x73, 0x74, 0x5f, 0x64,
	0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x74, 0x72, 0x75,
	0x73, 0x74, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x2d, 0x0a, 0x12, 0x77, 0x61, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x11, 0x77,
	0x61, 0x79, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x6f, 0x64, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61,
	0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x61,
	0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2d, 0x0a, 0x12, 0x63,
	0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63, 0x61, 0x6c, 0x5f, 0x72, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x63, 0x61, 0x6e, 0x6f, 0x6e, 0x69, 0x63,
	0x61, 0x6c, 0x52, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x41, 0x0a, 0x0d, 0x77, 0x6f,
	0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1c, 0x2e, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f,
	0x61, 0x64, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x68, 0x62, 0x6f,
	0x6e, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65,
	0x48, 0x62, 0x6f, 0x6e, 0x65, 0x12, 0x49, 0x0a, 0x0b, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c,
	0x5f, 0x69, 0x70, 0x73, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x73, 0x74,
	0x69, 0x6f, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x57, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x49, 0x70, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0a, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x49, 0x70, 0x73,
	0x12, 0x35, 0x0a, 0x16, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x18, 0x10, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x15, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x2e,
	0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61,
	0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a,
	0x57, 0x0a, 0x0f, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x49, 0x70, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x2e, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x2e, 0x77, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x36, 0x0a, 0x08, 0x50, 0x6f, 0x72, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x69, 0x73, 0x74, 0x69, 0x6f, 0x2e, 0x77, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x05, 0x70, 0x6f, 0x72, 0x74, 0x73,
	0x22, 0x4a, 0x0a, 0x04, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x2a, 0x2c, 0x0a, 0x0e,
	0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b,
	0x0a, 0x07, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x55,
	0x4e, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x01, 0x2a, 0x3d, 0x0a, 0x0c, 0x57, 0x6f,
	0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x44, 0x45,
	0x50, 0x4c, 0x4f, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x52,
	0x4f, 0x4e, 0x4a, 0x4f, 0x42, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x50, 0x4f, 0x44, 0x10, 0x02,
	0x12, 0x07, 0x0a, 0x03, 0x4a, 0x4f, 0x42, 0x10, 0x03, 0x2a, 0x20, 0x0a, 0x08, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x10,
	0x00, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x54, 0x54, 0x50, 0x10, 0x01, 0x42, 0x11, 0x5a, 0x0f, 0x70,
	0x6b, 0x67, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x61, 0x70, 0x69, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workloadapi_workload_proto_rawDescOnce sync.Once
	file_workloadapi_workload_proto_rawDescData = file_workloadapi_workload_proto_rawDesc
)

func file_workloadapi_workload_proto_rawDescGZIP() []byte {
	file_workloadapi_workload_proto_rawDescOnce.Do(func() {
		file_workloadapi_workload_proto_rawDescData = protoimpl.X.CompressGZIP(file_workloadapi_workload_proto_rawDescData)
	})
	return file_workloadapi_workload_proto_rawDescData
}

var file_workloadapi_workload_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_workloadapi_workload_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_workloadapi_workload_proto_goTypes = []interface{}{
	(WorkloadStatus)(0), // 0: istio.workload.WorkloadStatus
	(WorkloadType)(0),   // 1: istio.workload.WorkloadType
	(Protocol)(0),       // 2: istio.workload.Protocol
	(*Workload)(nil),    // 3: istio.workload.Workload
	(*PortList)(nil),    // 4: istio.workload.PortList
	(*Port)(nil),        // 5: istio.workload.Port
	nil,                 // 6: istio.workload.Workload.VirtualIpsEntry
}
var file_workloadapi_workload_proto_depIdxs = []int32{
	2, // 0: istio.workload.Workload.protocol:type_name -> istio.workload.Protocol
	1, // 1: istio.workload.Workload.workload_type:type_name -> istio.workload.WorkloadType
	6, // 2: istio.workload.Workload.virtual_ips:type_name -> istio.workload.Workload.VirtualIpsEntry
	0, // 3: istio.workload.Workload.status:type_name -> istio.workload.WorkloadStatus
	5, // 4: istio.workload.PortList.ports:type_name -> istio.workload.Port
	4, // 5: istio.workload.Workload.VirtualIpsEntry.value:type_name -> istio.workload.PortList
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_workloadapi_workload_proto_init() }
func file_workloadapi_workload_proto_init() {
	if File_workloadapi_workload_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_workloadapi_workload_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Workload); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workloadapi_workload_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PortList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workloadapi_workload_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Port); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_workloadapi_workload_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_workloadapi_workload_proto_goTypes,
		DependencyIndexes: file_workloadapi_workload_proto_depIdxs,
		EnumInfos:         file_workloadapi_workload_proto_enumTypes,
		MessageInfos:      file_workloadapi_workload_proto_msgTypes,
	}.Build()
	File_workloadapi_workload_proto = out.File
	file_workloadapi_workload_proto_rawDesc = nil
	file_workloadapi_workload_proto_goTypes = nil
	file_workloadapi_workload_proto_depIdxs = nil
}
