// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/monitoring/v3/alert_service.proto

package monitoring

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"
import google_protobuf5 "github.com/golang/protobuf/ptypes/empty"
import google_protobuf6 "google.golang.org/genproto/protobuf/field_mask"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// The protocol for the `CreateAlertPolicy` request.
type CreateAlertPolicyRequest struct {
	// The project in which to create the alerting policy. The format is
	// `projects/[PROJECT_ID]`.
	//
	// Note that this field names the parent container in which the alerting
	// policy will be written, not the name of the created policy. The alerting
	// policy that is returned will have a name that contains a normalized
	// representation of this name as a prefix but adds a suffix of the form
	// `/alertPolicies/[POLICY_ID]`, identifying the policy in the container.
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	// The requested alerting policy. You should omit the `name` field in this
	// policy. The name will be returned in the new policy, including
	// a new [ALERT_POLICY_ID] value.
	AlertPolicy *AlertPolicy `protobuf:"bytes,2,opt,name=alert_policy,json=alertPolicy" json:"alert_policy,omitempty"`
}

func (m *CreateAlertPolicyRequest) Reset()                    { *m = CreateAlertPolicyRequest{} }
func (m *CreateAlertPolicyRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateAlertPolicyRequest) ProtoMessage()               {}
func (*CreateAlertPolicyRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *CreateAlertPolicyRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateAlertPolicyRequest) GetAlertPolicy() *AlertPolicy {
	if m != nil {
		return m.AlertPolicy
	}
	return nil
}

// The protocol for the `GetAlertPolicy` request.
type GetAlertPolicyRequest struct {
	// The alerting policy to retrieve. The format is
	//
	//     projects/[PROJECT_ID]/alertPolicies/[ALERT_POLICY_ID]
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *GetAlertPolicyRequest) Reset()                    { *m = GetAlertPolicyRequest{} }
func (m *GetAlertPolicyRequest) String() string            { return proto.CompactTextString(m) }
func (*GetAlertPolicyRequest) ProtoMessage()               {}
func (*GetAlertPolicyRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *GetAlertPolicyRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The protocol for the `ListAlertPolicies` request.
type ListAlertPoliciesRequest struct {
	// The project whose alert policies are to be listed. The format is
	//
	//     projects/[PROJECT_ID]
	//
	// Note that this field names the parent container in which the alerting
	// policies to be listed are stored. To retrieve a single alerting policy
	// by name, use the
	// [GetAlertPolicy][google.monitoring.v3.AlertPolicyService.GetAlertPolicy]
	// operation, instead.
	Name string `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	// If provided, this field specifies the criteria that must be met by
	// alert policies to be included in the response.
	//
	// For more details, see [sorting and
	// filtering](/monitoring/api/v3/sorting-and-filtering).
	Filter string `protobuf:"bytes,5,opt,name=filter" json:"filter,omitempty"`
	// A comma-separated list of fields by which to sort the result. Supports
	// the same set of field references as the `filter` field. Entries can be
	// prefixed with a minus sign to sort by the field in descending order.
	//
	// For more details, see [sorting and
	// filtering](/monitoring/api/v3/sorting-and-filtering).
	OrderBy string `protobuf:"bytes,6,opt,name=order_by,json=orderBy" json:"order_by,omitempty"`
	// The maximum number of results to return in a single response.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
	// If this field is not empty then it must contain the `nextPageToken` value
	// returned by a previous call to this method.  Using this field causes the
	// method to return more results from the previous method call.
	PageToken string `protobuf:"bytes,3,opt,name=page_token,json=pageToken" json:"page_token,omitempty"`
}

func (m *ListAlertPoliciesRequest) Reset()                    { *m = ListAlertPoliciesRequest{} }
func (m *ListAlertPoliciesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListAlertPoliciesRequest) ProtoMessage()               {}
func (*ListAlertPoliciesRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *ListAlertPoliciesRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ListAlertPoliciesRequest) GetFilter() string {
	if m != nil {
		return m.Filter
	}
	return ""
}

func (m *ListAlertPoliciesRequest) GetOrderBy() string {
	if m != nil {
		return m.OrderBy
	}
	return ""
}

func (m *ListAlertPoliciesRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListAlertPoliciesRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

// The protocol for the `ListAlertPolicies` response.
type ListAlertPoliciesResponse struct {
	// The returned alert policies.
	AlertPolicies []*AlertPolicy `protobuf:"bytes,3,rep,name=alert_policies,json=alertPolicies" json:"alert_policies,omitempty"`
	// If there might be more results than were returned, then this field is set
	// to a non-empty value. To see the additional results,
	// use that value as `pageToken` in the next call to this method.
	NextPageToken string `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken" json:"next_page_token,omitempty"`
}

func (m *ListAlertPoliciesResponse) Reset()                    { *m = ListAlertPoliciesResponse{} }
func (m *ListAlertPoliciesResponse) String() string            { return proto.CompactTextString(m) }
func (*ListAlertPoliciesResponse) ProtoMessage()               {}
func (*ListAlertPoliciesResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *ListAlertPoliciesResponse) GetAlertPolicies() []*AlertPolicy {
	if m != nil {
		return m.AlertPolicies
	}
	return nil
}

func (m *ListAlertPoliciesResponse) GetNextPageToken() string {
	if m != nil {
		return m.NextPageToken
	}
	return ""
}

// The protocol for the `UpdateAlertPolicy` request.
type UpdateAlertPolicyRequest struct {
	// Optional. A list of alerting policy field names. If this field is not
	// empty, each listed field in the existing alerting policy is set to the
	// value of the corresponding field in the supplied policy (`alert_policy`),
	// or to the field's default value if the field is not in the supplied
	// alerting policy.  Fields not listed retain their previous value.
	//
	// Examples of valid field masks include `display_name`, `documentation`,
	// `documentation.content`, `documentation.mime_type`, `user_labels`,
	// `user_label.nameofkey`, `enabled`, `conditions`, `combiner`, etc.
	//
	// If this field is empty, then the supplied alerting policy replaces the
	// existing policy. It is the same as deleting the existing policy and
	// adding the supplied policy, except for the following:
	//
	// +   The new policy will have the same `[ALERT_POLICY_ID]` as the former
	//     policy. This gives you continuity with the former policy in your
	//     notifications and incidents.
	// +   Conditions in the new policy will keep their former `[CONDITION_ID]` if
	//     the supplied condition includes the `name` field with that
	//     `[CONDITION_ID]`. If the supplied condition omits the `name` field,
	//     then a new `[CONDITION_ID]` is created.
	UpdateMask *google_protobuf6.FieldMask `protobuf:"bytes,2,opt,name=update_mask,json=updateMask" json:"update_mask,omitempty"`
	// Required. The updated alerting policy or the updated values for the
	// fields listed in `update_mask`.
	// If `update_mask` is not empty, any fields in this policy that are
	// not in `update_mask` are ignored.
	AlertPolicy *AlertPolicy `protobuf:"bytes,3,opt,name=alert_policy,json=alertPolicy" json:"alert_policy,omitempty"`
}

func (m *UpdateAlertPolicyRequest) Reset()                    { *m = UpdateAlertPolicyRequest{} }
func (m *UpdateAlertPolicyRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateAlertPolicyRequest) ProtoMessage()               {}
func (*UpdateAlertPolicyRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *UpdateAlertPolicyRequest) GetUpdateMask() *google_protobuf6.FieldMask {
	if m != nil {
		return m.UpdateMask
	}
	return nil
}

func (m *UpdateAlertPolicyRequest) GetAlertPolicy() *AlertPolicy {
	if m != nil {
		return m.AlertPolicy
	}
	return nil
}

// The protocol for the `DeleteAlertPolicy` request.
type DeleteAlertPolicyRequest struct {
	// The alerting policy to delete. The format is:
	//
	//     projects/[PROJECT_ID]/alertPolicies/[ALERT_POLICY_ID]
	//
	// For more information, see [AlertPolicy][google.monitoring.v3.AlertPolicy].
	Name string `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
}

func (m *DeleteAlertPolicyRequest) Reset()                    { *m = DeleteAlertPolicyRequest{} }
func (m *DeleteAlertPolicyRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteAlertPolicyRequest) ProtoMessage()               {}
func (*DeleteAlertPolicyRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *DeleteAlertPolicyRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateAlertPolicyRequest)(nil), "google.monitoring.v3.CreateAlertPolicyRequest")
	proto.RegisterType((*GetAlertPolicyRequest)(nil), "google.monitoring.v3.GetAlertPolicyRequest")
	proto.RegisterType((*ListAlertPoliciesRequest)(nil), "google.monitoring.v3.ListAlertPoliciesRequest")
	proto.RegisterType((*ListAlertPoliciesResponse)(nil), "google.monitoring.v3.ListAlertPoliciesResponse")
	proto.RegisterType((*UpdateAlertPolicyRequest)(nil), "google.monitoring.v3.UpdateAlertPolicyRequest")
	proto.RegisterType((*DeleteAlertPolicyRequest)(nil), "google.monitoring.v3.DeleteAlertPolicyRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AlertPolicyService service

type AlertPolicyServiceClient interface {
	// Lists the existing alerting policies for the project.
	ListAlertPolicies(ctx context.Context, in *ListAlertPoliciesRequest, opts ...grpc.CallOption) (*ListAlertPoliciesResponse, error)
	// Gets a single alerting policy.
	GetAlertPolicy(ctx context.Context, in *GetAlertPolicyRequest, opts ...grpc.CallOption) (*AlertPolicy, error)
	// Creates a new alerting policy.
	CreateAlertPolicy(ctx context.Context, in *CreateAlertPolicyRequest, opts ...grpc.CallOption) (*AlertPolicy, error)
	// Deletes an alerting policy.
	DeleteAlertPolicy(ctx context.Context, in *DeleteAlertPolicyRequest, opts ...grpc.CallOption) (*google_protobuf5.Empty, error)
	// Updates an alerting policy. You can either replace the entire policy with
	// a new one or replace only certain fields in the current alerting policy by
	// specifying the fields to be updated via `updateMask`. Returns the
	// updated alerting policy.
	UpdateAlertPolicy(ctx context.Context, in *UpdateAlertPolicyRequest, opts ...grpc.CallOption) (*AlertPolicy, error)
}

type alertPolicyServiceClient struct {
	cc *grpc.ClientConn
}

func NewAlertPolicyServiceClient(cc *grpc.ClientConn) AlertPolicyServiceClient {
	return &alertPolicyServiceClient{cc}
}

func (c *alertPolicyServiceClient) ListAlertPolicies(ctx context.Context, in *ListAlertPoliciesRequest, opts ...grpc.CallOption) (*ListAlertPoliciesResponse, error) {
	out := new(ListAlertPoliciesResponse)
	err := grpc.Invoke(ctx, "/google.monitoring.v3.AlertPolicyService/ListAlertPolicies", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertPolicyServiceClient) GetAlertPolicy(ctx context.Context, in *GetAlertPolicyRequest, opts ...grpc.CallOption) (*AlertPolicy, error) {
	out := new(AlertPolicy)
	err := grpc.Invoke(ctx, "/google.monitoring.v3.AlertPolicyService/GetAlertPolicy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertPolicyServiceClient) CreateAlertPolicy(ctx context.Context, in *CreateAlertPolicyRequest, opts ...grpc.CallOption) (*AlertPolicy, error) {
	out := new(AlertPolicy)
	err := grpc.Invoke(ctx, "/google.monitoring.v3.AlertPolicyService/CreateAlertPolicy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertPolicyServiceClient) DeleteAlertPolicy(ctx context.Context, in *DeleteAlertPolicyRequest, opts ...grpc.CallOption) (*google_protobuf5.Empty, error) {
	out := new(google_protobuf5.Empty)
	err := grpc.Invoke(ctx, "/google.monitoring.v3.AlertPolicyService/DeleteAlertPolicy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alertPolicyServiceClient) UpdateAlertPolicy(ctx context.Context, in *UpdateAlertPolicyRequest, opts ...grpc.CallOption) (*AlertPolicy, error) {
	out := new(AlertPolicy)
	err := grpc.Invoke(ctx, "/google.monitoring.v3.AlertPolicyService/UpdateAlertPolicy", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AlertPolicyService service

type AlertPolicyServiceServer interface {
	// Lists the existing alerting policies for the project.
	ListAlertPolicies(context.Context, *ListAlertPoliciesRequest) (*ListAlertPoliciesResponse, error)
	// Gets a single alerting policy.
	GetAlertPolicy(context.Context, *GetAlertPolicyRequest) (*AlertPolicy, error)
	// Creates a new alerting policy.
	CreateAlertPolicy(context.Context, *CreateAlertPolicyRequest) (*AlertPolicy, error)
	// Deletes an alerting policy.
	DeleteAlertPolicy(context.Context, *DeleteAlertPolicyRequest) (*google_protobuf5.Empty, error)
	// Updates an alerting policy. You can either replace the entire policy with
	// a new one or replace only certain fields in the current alerting policy by
	// specifying the fields to be updated via `updateMask`. Returns the
	// updated alerting policy.
	UpdateAlertPolicy(context.Context, *UpdateAlertPolicyRequest) (*AlertPolicy, error)
}

func RegisterAlertPolicyServiceServer(s *grpc.Server, srv AlertPolicyServiceServer) {
	s.RegisterService(&_AlertPolicyService_serviceDesc, srv)
}

func _AlertPolicyService_ListAlertPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAlertPoliciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertPolicyServiceServer).ListAlertPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.monitoring.v3.AlertPolicyService/ListAlertPolicies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertPolicyServiceServer).ListAlertPolicies(ctx, req.(*ListAlertPoliciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertPolicyService_GetAlertPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAlertPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertPolicyServiceServer).GetAlertPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.monitoring.v3.AlertPolicyService/GetAlertPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertPolicyServiceServer).GetAlertPolicy(ctx, req.(*GetAlertPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertPolicyService_CreateAlertPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAlertPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertPolicyServiceServer).CreateAlertPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.monitoring.v3.AlertPolicyService/CreateAlertPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertPolicyServiceServer).CreateAlertPolicy(ctx, req.(*CreateAlertPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertPolicyService_DeleteAlertPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAlertPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertPolicyServiceServer).DeleteAlertPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.monitoring.v3.AlertPolicyService/DeleteAlertPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertPolicyServiceServer).DeleteAlertPolicy(ctx, req.(*DeleteAlertPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AlertPolicyService_UpdateAlertPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAlertPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlertPolicyServiceServer).UpdateAlertPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.monitoring.v3.AlertPolicyService/UpdateAlertPolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlertPolicyServiceServer).UpdateAlertPolicy(ctx, req.(*UpdateAlertPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AlertPolicyService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.monitoring.v3.AlertPolicyService",
	HandlerType: (*AlertPolicyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListAlertPolicies",
			Handler:    _AlertPolicyService_ListAlertPolicies_Handler,
		},
		{
			MethodName: "GetAlertPolicy",
			Handler:    _AlertPolicyService_GetAlertPolicy_Handler,
		},
		{
			MethodName: "CreateAlertPolicy",
			Handler:    _AlertPolicyService_CreateAlertPolicy_Handler,
		},
		{
			MethodName: "DeleteAlertPolicy",
			Handler:    _AlertPolicyService_DeleteAlertPolicy_Handler,
		},
		{
			MethodName: "UpdateAlertPolicy",
			Handler:    _AlertPolicyService_UpdateAlertPolicy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/monitoring/v3/alert_service.proto",
}

func init() { proto.RegisterFile("google/monitoring/v3/alert_service.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 656 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x41, 0x6f, 0xd3, 0x4c,
	0x10, 0x95, 0x93, 0x36, 0x5f, 0xbb, 0xfd, 0x5a, 0x94, 0x15, 0x54, 0xae, 0x0b, 0x52, 0x30, 0x2a,
	0x54, 0xad, 0xb0, 0xa5, 0xf8, 0x04, 0x15, 0x48, 0xa4, 0x85, 0xf6, 0x40, 0xa5, 0x28, 0x85, 0x1e,
	0x50, 0xa4, 0x68, 0x93, 0x4c, 0xac, 0x25, 0x8e, 0xd7, 0x78, 0x37, 0x11, 0x29, 0xea, 0x85, 0x23,
	0x12, 0xe2, 0xc0, 0x99, 0x03, 0x47, 0x38, 0x20, 0x7e, 0x07, 0x57, 0xfe, 0x02, 0x3f, 0x04, 0x79,
	0xed, 0x34, 0x76, 0x6d, 0xab, 0x16, 0xb7, 0xcc, 0xce, 0xdb, 0x99, 0xb7, 0x6f, 0xde, 0x38, 0x68,
	0xdb, 0x66, 0xcc, 0x76, 0xc0, 0x1c, 0x31, 0x97, 0x0a, 0xe6, 0x53, 0xd7, 0x36, 0x27, 0x96, 0x49,
	0x1c, 0xf0, 0x45, 0x87, 0x83, 0x3f, 0xa1, 0x3d, 0x30, 0x3c, 0x9f, 0x09, 0x86, 0xaf, 0x87, 0x48,
	0x63, 0x8e, 0x34, 0x26, 0x96, 0x76, 0x33, 0xba, 0x4f, 0x3c, 0x6a, 0x12, 0xd7, 0x65, 0x82, 0x08,
	0xca, 0x5c, 0x1e, 0xde, 0xd1, 0x6a, 0xf9, 0xd5, 0x23, 0xc4, 0x66, 0x84, 0x90, 0x51, 0x77, 0x3c,
	0x30, 0x61, 0xe4, 0x89, 0xe9, 0xa5, 0xeb, 0x17, 0xc9, 0x01, 0x05, 0xa7, 0xdf, 0x19, 0x11, 0x3e,
	0x0c, 0x11, 0xba, 0x40, 0xea, 0xbe, 0x0f, 0x44, 0xc0, 0x93, 0xa0, 0x66, 0x93, 0x39, 0xb4, 0x37,
	0x6d, 0xc1, 0x9b, 0x31, 0x70, 0x81, 0x31, 0x5a, 0x70, 0xc9, 0x08, 0xd4, 0x72, 0x4d, 0xd9, 0x5e,
	0x6e, 0xc9, 0xdf, 0xf8, 0x00, 0xfd, 0x1f, 0xbe, 0xcd, 0x93, 0x50, 0xb5, 0x54, 0x53, 0xb6, 0x57,
	0xea, 0xb7, 0x8d, 0xac, 0xb7, 0x19, 0xf1, 0x9a, 0x2b, 0x64, 0x1e, 0xe8, 0xbb, 0xe8, 0xc6, 0x21,
	0x88, 0x62, 0x2d, 0xf5, 0x2f, 0x0a, 0x52, 0x9f, 0x53, 0x1e, 0x83, 0x53, 0xe0, 0x97, 0x2f, 0x2c,
	0xc4, 0x38, 0xae, 0xa3, 0xca, 0x80, 0x3a, 0x02, 0x7c, 0x75, 0x51, 0x9e, 0x46, 0x11, 0xde, 0x40,
	0x4b, 0xcc, 0xef, 0x83, 0xdf, 0xe9, 0x4e, 0xd5, 0x8a, 0xcc, 0xfc, 0x27, 0xe3, 0xc6, 0x14, 0x6f,
	0xa2, 0x65, 0x8f, 0xd8, 0xd0, 0xe1, 0xf4, 0x0c, 0xe4, 0x9b, 0x16, 0x5b, 0x4b, 0xc1, 0xc1, 0x09,
	0x3d, 0x03, 0x7c, 0x0b, 0x21, 0x99, 0x14, 0x6c, 0x08, 0x6e, 0x44, 0x4d, 0xc2, 0x5f, 0x04, 0x07,
	0xfa, 0x47, 0x05, 0x6d, 0x64, 0xf0, 0xe3, 0x1e, 0x73, 0x39, 0xe0, 0x23, 0xb4, 0x16, 0x13, 0x8c,
	0x02, 0x57, 0xcb, 0xb5, 0x72, 0x31, 0xc9, 0x56, 0x49, 0xbc, 0x22, 0xbe, 0x8b, 0xae, 0xb9, 0xf0,
	0x56, 0x74, 0x62, 0x5c, 0x4a, 0x92, 0xcb, 0x6a, 0x70, 0xdc, 0xbc, 0xe0, 0x13, 0xe8, 0xf5, 0xd2,
	0xeb, 0x67, 0xcf, 0x74, 0x0f, 0xad, 0x8c, 0x65, 0x4e, 0x9a, 0x20, 0x1a, 0x9f, 0x36, 0xe3, 0x32,
	0xf3, 0x89, 0xf1, 0x2c, 0xf0, 0xc9, 0x31, 0xe1, 0xc3, 0x16, 0x0a, 0xe1, 0xc1, 0xef, 0xd4, 0xf0,
	0xcb, 0xff, 0x34, 0x7c, 0x03, 0xa9, 0x07, 0xe0, 0x40, 0x51, 0xcb, 0xd5, 0x7f, 0x54, 0x10, 0x8e,
	0x41, 0x4f, 0xc2, 0xa5, 0xc2, 0x5f, 0x15, 0x54, 0x4d, 0xc9, 0x8e, 0x8d, 0x6c, 0x32, 0x79, 0xfe,
	0xd1, 0xcc, 0xc2, 0xf8, 0x70, 0x9e, 0xfa, 0xee, 0xfb, 0xdf, 0x7f, 0x3e, 0x97, 0xb6, 0xf0, 0x9d,
	0x60, 0x11, 0xdf, 0x05, 0x04, 0x1f, 0x79, 0x3e, 0x7b, 0x0d, 0x3d, 0xc1, 0xcd, 0x9d, 0x73, 0x33,
	0x39, 0xb2, 0x4f, 0x0a, 0x5a, 0x4b, 0x1a, 0x1d, 0xef, 0x66, 0x37, 0xcc, 0x5c, 0x07, 0xed, 0x6a,
	0x69, 0xf5, 0xfb, 0x92, 0xcf, 0x3d, 0xbc, 0x95, 0xc5, 0x27, 0x49, 0xc7, 0xdc, 0x39, 0x97, 0xaa,
	0xa5, 0x16, 0x3e, 0x4f, 0xb5, 0xbc, 0x2f, 0x43, 0x11, 0x5e, 0x0f, 0x24, 0x2f, 0x4b, 0x2f, 0xa2,
	0xd3, 0xc3, 0x84, 0xad, 0xf0, 0x07, 0x05, 0x55, 0x53, 0x0e, 0xc9, 0xe3, 0x98, 0x67, 0x25, 0x6d,
	0x3d, 0x65, 0xea, 0xa7, 0xc1, 0x97, 0x71, 0x26, 0xd8, 0x4e, 0x41, 0xc1, 0x7e, 0x2a, 0xa8, 0x9a,
	0xda, 0xa6, 0x3c, 0x32, 0x79, 0x6b, 0x57, 0x44, 0xb0, 0x23, 0xc9, 0xab, 0x51, 0xaf, 0x4b, 0x5e,
	0x71, 0x41, 0x8c, 0xab, 0x48, 0x26, 0xf5, 0x6b, 0x7c, 0x53, 0x90, 0xda, 0x63, 0xa3, 0xcc, 0x96,
	0x8d, 0xaa, 0xec, 0x19, 0x2d, 0x51, 0x33, 0x90, 0xa6, 0xa9, 0xbc, 0x7a, 0x1c, 0x41, 0x6d, 0xe6,
	0x10, 0xd7, 0x36, 0x98, 0x6f, 0x9b, 0x36, 0xb8, 0x52, 0x38, 0x33, 0x4c, 0x11, 0x8f, 0xf2, 0xe4,
	0xbf, 0xd0, 0xde, 0x3c, 0xfa, 0x5e, 0xd2, 0x0e, 0xc3, 0x02, 0xfb, 0x0e, 0x1b, 0xf7, 0x8d, 0xe3,
	0x79, 0xc7, 0x53, 0xeb, 0xd7, 0x2c, 0xd9, 0x96, 0xc9, 0xf6, 0x3c, 0xd9, 0x3e, 0xb5, 0xba, 0x15,
	0xd9, 0xc4, 0xfa, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x6f, 0x1f, 0xe6, 0xf0, 0x47, 0x07, 0x00, 0x00,
}
