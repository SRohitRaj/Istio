// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rbac/v1alpha1/rbac.proto

package v1alpha1 // import "istio.io/api/rbac/v1alpha1"

/*
Istio RBAC (Role Based Access Control) defines ServiceRole and ServiceRoleBinding
objects.

A ServiceRole specification includes a list of rules (permissions). Each rule has
the following standard fields:
* services: a list of services.
* methods: HTTP methods. In the case of gRPC, this field is ignored because the value is always "POST".
* paths: HTTP paths or gRPC methods. Note that gRPC methods should be
  presented in the form of "packageName.serviceName/methodName".

In addition to the standard fields, operators can use custom fields in the "constraints"
section. The name of a custom field must match one of the "properties" in the "action" part
of the "authorization" template (https://github.com/istio/istio/blob/master/mixer/template/authorization/template.proto).

For example, suppose we define an instance of the "authorization" template, named "requestcontext".

```yaml
apiVersion: "config.istio.io/v1alpha1"
kind: authorization
metadata:
  name: requestcontext
  namespace: istio-system
spec:
  subject:
    user: source.user | ""
    groups: ""
    properties:
      service: source.service | ""
      namespace: source.namespace | ""
  action:
    namespace: destination.namespace | ""
    service: destination.service | ""
    method: request.method | ""
    path: request.path | ""
    properties:
      version: request.headers["version"] | ""
```

Below is an example of ServiceRole object "product-viewer", which has "read" ("GET" and "HEAD")
access to "products.svc.cluster.local" service at versions "v1" and "v2". "path" is not specified,
so it applies to any path in the service.

```yaml
apiVersion: "config.istio.io/v1alpha1"
kind: ServiceRole
metadata:
  name: products-viewer
  namespace: default
spec:
  rules:
  - services: ["products.svc.cluster.local"]
    methods: ["GET", "HEAD"]
    constraints:
    - key: "version"
      value: ["v1", "v2"]
```

A ServiceRoleBinding specification includes two parts:
* "roleRef" refers to a ServiceRole object in the same namespace.
* A list of "subjects" that are assigned the roles.

A subject is represented with a set of "properties". The name of a property must match one of
the fields ("user" or "groups" or one of the "properties") in the "subject" part of the "authorization"
template (https://github.com/istio/istio/blob/master/mixer/template/authorization/template.proto).

Below is an example of ServiceRoleBinding object "test-binding-products", which binds two subjects
to ServiceRole "product-viewer":
* User "alice@yahoo.com"
* "reviews" service in "abc" namespace.

```yaml
apiVersion: "config.istio.io/v1alpha1"
kind: ServiceRoleBinding
metadata:
  name: test-binding-products
  namespace: default
spec:
  subjects:
  - user: alice@yahoo.com
  - properties:
      service: "reviews"
      namespace: "abc"
  roleRef:
    kind: ServiceRole
    name: "products-viewer"
```
*/

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RbacConfig_Mode int32

const (
	// Disable Istio RBAC completely, any other config in RbacConfig will be ignored and Istio RBAC policies
	// will not be enforced.
	RbacConfig_OFF RbacConfig_Mode = 0
	// Enable Istio RBAC for all services and namespaces.
	RbacConfig_ON RbacConfig_Mode = 1
	// Enable Istio RBAC only for services and namespaces specified in the inclusion field. Any other
	// services and namespaces not in the inclusion field will not be enforced by Istio RBAC policies.
	RbacConfig_ON_WITH_INCLUSION RbacConfig_Mode = 2
	// Enable Istio RBAC for all services and namespaces except those specified in the exclusion field. Any other
	// services and namespaces not in the exclusion field will be enforced by Istio RBAC policies.
	RbacConfig_ON_WITH_EXCLUSION RbacConfig_Mode = 3
)

var RbacConfig_Mode_name = map[int32]string{
	0: "OFF",
	1: "ON",
	2: "ON_WITH_INCLUSION",
	3: "ON_WITH_EXCLUSION",
}
var RbacConfig_Mode_value = map[string]int32{
	"OFF":               0,
	"ON":                1,
	"ON_WITH_INCLUSION": 2,
	"ON_WITH_EXCLUSION": 3,
}

func (x RbacConfig_Mode) String() string {
	return proto.EnumName(RbacConfig_Mode_name, int32(x))
}
func (RbacConfig_Mode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{5, 0}
}

// ServiceRole specification contains a list of access rules (permissions).
// This represent the "Spec" part of the ServiceRole object. The name and namespace
// of the ServiceRole is specified in "metadata" section of the ServiceRole object.
type ServiceRole struct {
	// Required. The set of access rules (permissions) that the role has.
	Rules                []*AccessRule `protobuf:"bytes,1,rep,name=rules" json:"rules,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ServiceRole) Reset()         { *m = ServiceRole{} }
func (m *ServiceRole) String() string { return proto.CompactTextString(m) }
func (*ServiceRole) ProtoMessage()    {}
func (*ServiceRole) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{0}
}
func (m *ServiceRole) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceRole.Unmarshal(m, b)
}
func (m *ServiceRole) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceRole.Marshal(b, m, deterministic)
}
func (dst *ServiceRole) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceRole.Merge(dst, src)
}
func (m *ServiceRole) XXX_Size() int {
	return xxx_messageInfo_ServiceRole.Size(m)
}
func (m *ServiceRole) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceRole.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceRole proto.InternalMessageInfo

func (m *ServiceRole) GetRules() []*AccessRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

// AccessRule defines a permission to access a list of services.
type AccessRule struct {
	// Required. A list of service names.
	// Exact match, prefix match, and suffix match are supported for service names.
	// For example, the service name "bookstore.mtv.cluster.local" matches
	// "bookstore.mtv.cluster.local" (exact match), or "bookstore*" (prefix match),
	// or "*.mtv.cluster.local" (suffix match).
	// If set to ["*"], it refers to all services in the namespace.
	Services []string `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
	// Optional. A list of HTTP paths or gRPC methods.
	// gRPC methods must be presented as fully-qualified name in the form of
	// packageName.serviceName/methodName.
	// Exact match, prefix match, and suffix match are supported for paths.
	// For example, the path "/books/review" matches
	// "/books/review" (exact match), or "/books/*" (prefix match),
	// or "*/review" (suffix match).
	// If not specified, it applies to any path.
	Paths []string `protobuf:"bytes,2,rep,name=paths" json:"paths,omitempty"`
	// Optional. A list of HTTP methods (e.g., "GET", "POST").
	// It is ignored in gRPC case because the value is always "POST".
	// If set to ["*"] or not specified, it applies to any method.
	Methods []string `protobuf:"bytes,3,rep,name=methods" json:"methods,omitempty"`
	// Optional. Extra constraints in the ServiceRole specification.
	// The above ServiceRole examples shows an example of constraint "version".
	Constraints          []*AccessRule_Constraint `protobuf:"bytes,4,rep,name=constraints" json:"constraints,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *AccessRule) Reset()         { *m = AccessRule{} }
func (m *AccessRule) String() string { return proto.CompactTextString(m) }
func (*AccessRule) ProtoMessage()    {}
func (*AccessRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{1}
}
func (m *AccessRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccessRule.Unmarshal(m, b)
}
func (m *AccessRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccessRule.Marshal(b, m, deterministic)
}
func (dst *AccessRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessRule.Merge(dst, src)
}
func (m *AccessRule) XXX_Size() int {
	return xxx_messageInfo_AccessRule.Size(m)
}
func (m *AccessRule) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessRule.DiscardUnknown(m)
}

var xxx_messageInfo_AccessRule proto.InternalMessageInfo

func (m *AccessRule) GetServices() []string {
	if m != nil {
		return m.Services
	}
	return nil
}

func (m *AccessRule) GetPaths() []string {
	if m != nil {
		return m.Paths
	}
	return nil
}

func (m *AccessRule) GetMethods() []string {
	if m != nil {
		return m.Methods
	}
	return nil
}

func (m *AccessRule) GetConstraints() []*AccessRule_Constraint {
	if m != nil {
		return m.Constraints
	}
	return nil
}

// Definition of a custom constraint. The key of a custom constraint must match
// one of the "properties" in the "action" part of the "authorization" template
// (https://github.com/istio/istio/blob/master/mixer/template/authorization/template.proto).
type AccessRule_Constraint struct {
	// Key of the constraint.
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	// List of valid values for the constraint.
	// Exact match, prefix match, and suffix match are supported for constraint values.
	// For example, the value "v1alpha2" matches
	// "v1alpha2" (exact match), or "v1*" (prefix match),
	// or "*alpha2" (suffix match).
	Values               []string `protobuf:"bytes,2,rep,name=values" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccessRule_Constraint) Reset()         { *m = AccessRule_Constraint{} }
func (m *AccessRule_Constraint) String() string { return proto.CompactTextString(m) }
func (*AccessRule_Constraint) ProtoMessage()    {}
func (*AccessRule_Constraint) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{1, 0}
}
func (m *AccessRule_Constraint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccessRule_Constraint.Unmarshal(m, b)
}
func (m *AccessRule_Constraint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccessRule_Constraint.Marshal(b, m, deterministic)
}
func (dst *AccessRule_Constraint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccessRule_Constraint.Merge(dst, src)
}
func (m *AccessRule_Constraint) XXX_Size() int {
	return xxx_messageInfo_AccessRule_Constraint.Size(m)
}
func (m *AccessRule_Constraint) XXX_DiscardUnknown() {
	xxx_messageInfo_AccessRule_Constraint.DiscardUnknown(m)
}

var xxx_messageInfo_AccessRule_Constraint proto.InternalMessageInfo

func (m *AccessRule_Constraint) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *AccessRule_Constraint) GetValues() []string {
	if m != nil {
		return m.Values
	}
	return nil
}

// ServiceRoleBinding assigns a ServiceRole to a list of subjects.
// This represents the "Spec" part of the ServiceRoleBinding object. The name and namespace
// of the ServiceRoleBinding is specified in "metadata" section of the ServiceRoleBinding
// object.
type ServiceRoleBinding struct {
	// Required. List of subjects that are assigned the ServiceRole object.
	Subjects []*Subject `protobuf:"bytes,1,rep,name=subjects" json:"subjects,omitempty"`
	// Required. Reference to the ServiceRole object.
	RoleRef              *RoleRef `protobuf:"bytes,2,opt,name=roleRef" json:"roleRef,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceRoleBinding) Reset()         { *m = ServiceRoleBinding{} }
func (m *ServiceRoleBinding) String() string { return proto.CompactTextString(m) }
func (*ServiceRoleBinding) ProtoMessage()    {}
func (*ServiceRoleBinding) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{2}
}
func (m *ServiceRoleBinding) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceRoleBinding.Unmarshal(m, b)
}
func (m *ServiceRoleBinding) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceRoleBinding.Marshal(b, m, deterministic)
}
func (dst *ServiceRoleBinding) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceRoleBinding.Merge(dst, src)
}
func (m *ServiceRoleBinding) XXX_Size() int {
	return xxx_messageInfo_ServiceRoleBinding.Size(m)
}
func (m *ServiceRoleBinding) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceRoleBinding.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceRoleBinding proto.InternalMessageInfo

func (m *ServiceRoleBinding) GetSubjects() []*Subject {
	if m != nil {
		return m.Subjects
	}
	return nil
}

func (m *ServiceRoleBinding) GetRoleRef() *RoleRef {
	if m != nil {
		return m.RoleRef
	}
	return nil
}

// Subject defines an identity or a group of identities. The identity is either a user or
// a group or identified by a set of "properties". The name of the "properties" must match
// the "properties" in the "subject" part of the "authorization" template
// (https://github.com/istio/istio/blob/master/mixer/template/authorization/template.proto).
type Subject struct {
	// Optional. The user name/ID that the subject represents.
	User string `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
	// Optional. The group that the subject belongs to.
	Group string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
	// Optional. The set of properties that identify the subject.
	// In the above ServiceRoleBinding example, the second subject has two properties:
	//     service: "reviews"
	//     namespace: "abc"
	Properties           map[string]string `protobuf:"bytes,3,rep,name=properties" json:"properties,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Subject) Reset()         { *m = Subject{} }
func (m *Subject) String() string { return proto.CompactTextString(m) }
func (*Subject) ProtoMessage()    {}
func (*Subject) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{3}
}
func (m *Subject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Subject.Unmarshal(m, b)
}
func (m *Subject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Subject.Marshal(b, m, deterministic)
}
func (dst *Subject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subject.Merge(dst, src)
}
func (m *Subject) XXX_Size() int {
	return xxx_messageInfo_Subject.Size(m)
}
func (m *Subject) XXX_DiscardUnknown() {
	xxx_messageInfo_Subject.DiscardUnknown(m)
}

var xxx_messageInfo_Subject proto.InternalMessageInfo

func (m *Subject) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *Subject) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *Subject) GetProperties() map[string]string {
	if m != nil {
		return m.Properties
	}
	return nil
}

// RoleRef refers to a role object.
type RoleRef struct {
	// Required. The type of the role being referenced.
	// Currently, "ServiceRole" is the only supported value for "kind".
	Kind string `protobuf:"bytes,1,opt,name=kind" json:"kind,omitempty"`
	// Required. The name of the ServiceRole object being referenced.
	// The ServiceRole object must be in the same namespace as the ServiceRoleBinding
	// object.
	Name                 string   `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RoleRef) Reset()         { *m = RoleRef{} }
func (m *RoleRef) String() string { return proto.CompactTextString(m) }
func (*RoleRef) ProtoMessage()    {}
func (*RoleRef) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{4}
}
func (m *RoleRef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoleRef.Unmarshal(m, b)
}
func (m *RoleRef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoleRef.Marshal(b, m, deterministic)
}
func (dst *RoleRef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoleRef.Merge(dst, src)
}
func (m *RoleRef) XXX_Size() int {
	return xxx_messageInfo_RoleRef.Size(m)
}
func (m *RoleRef) XXX_DiscardUnknown() {
	xxx_messageInfo_RoleRef.DiscardUnknown(m)
}

var xxx_messageInfo_RoleRef proto.InternalMessageInfo

func (m *RoleRef) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *RoleRef) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// RbacConfig defines the global config to control Istio RBAC behavior.
// This Custom Resource is a singleton where only one Custom Resource should be created globally in
// the mesh and the namespace should be the same to other Istio components, which usually is istio-system.
// Note: This is enforced in both istioctl and server side, new Custom Resource will be rejected if found any
// existing one, the user should either delete the existing one or change the existing one directly.
//
// Below is an example of RbacConfig object "istio-rbac-config" which enables Istio RBAC for all
// services in the default namespace.
//
// ```yaml
// apiVersion: "config.istio.io/v1alpha1"
// kind: RbacConfig
// metadata:
//   name: istio-rbac-config
//   namespace: istio-system
// spec:
//   mode: ON_WITH_INCLUSION
//   inclusion:
//     namespaces: [ "default" ]
// ```
type RbacConfig struct {
	// Istio RBAC mode.
	Mode RbacConfig_Mode `protobuf:"varint,1,opt,name=mode,enum=istio.rbac.v1alpha1.RbacConfig_Mode" json:"mode,omitempty"`
	// A list of services or namespaces that should be enforced by Istio RBAC policies. Note: This field have
	// effect only when mode is ON_WITH_INCLUSION and will be ignored for any other modes.
	Inclusion *RbacConfig_Target `protobuf:"bytes,2,opt,name=inclusion" json:"inclusion,omitempty"`
	// A list of services or namespaces that should not be enforced by Istio RBAC policies. Note: This field have
	// effect only when mode is ON_WITH_EXCLUSION and will be ignored for any other modes.
	Exclusion            *RbacConfig_Target `protobuf:"bytes,3,opt,name=exclusion" json:"exclusion,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *RbacConfig) Reset()         { *m = RbacConfig{} }
func (m *RbacConfig) String() string { return proto.CompactTextString(m) }
func (*RbacConfig) ProtoMessage()    {}
func (*RbacConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{5}
}
func (m *RbacConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RbacConfig.Unmarshal(m, b)
}
func (m *RbacConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RbacConfig.Marshal(b, m, deterministic)
}
func (dst *RbacConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RbacConfig.Merge(dst, src)
}
func (m *RbacConfig) XXX_Size() int {
	return xxx_messageInfo_RbacConfig.Size(m)
}
func (m *RbacConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_RbacConfig.DiscardUnknown(m)
}

var xxx_messageInfo_RbacConfig proto.InternalMessageInfo

func (m *RbacConfig) GetMode() RbacConfig_Mode {
	if m != nil {
		return m.Mode
	}
	return RbacConfig_OFF
}

func (m *RbacConfig) GetInclusion() *RbacConfig_Target {
	if m != nil {
		return m.Inclusion
	}
	return nil
}

func (m *RbacConfig) GetExclusion() *RbacConfig_Target {
	if m != nil {
		return m.Exclusion
	}
	return nil
}

// Target defines a list of services or namespaces.
type RbacConfig_Target struct {
	// A list of services.
	Services []string `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
	// A list of namespaces.
	Namespaces           []string `protobuf:"bytes,2,rep,name=namespaces" json:"namespaces,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RbacConfig_Target) Reset()         { *m = RbacConfig_Target{} }
func (m *RbacConfig_Target) String() string { return proto.CompactTextString(m) }
func (*RbacConfig_Target) ProtoMessage()    {}
func (*RbacConfig_Target) Descriptor() ([]byte, []int) {
	return fileDescriptor_rbac_665e05e3b5f24d86, []int{5, 0}
}
func (m *RbacConfig_Target) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RbacConfig_Target.Unmarshal(m, b)
}
func (m *RbacConfig_Target) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RbacConfig_Target.Marshal(b, m, deterministic)
}
func (dst *RbacConfig_Target) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RbacConfig_Target.Merge(dst, src)
}
func (m *RbacConfig_Target) XXX_Size() int {
	return xxx_messageInfo_RbacConfig_Target.Size(m)
}
func (m *RbacConfig_Target) XXX_DiscardUnknown() {
	xxx_messageInfo_RbacConfig_Target.DiscardUnknown(m)
}

var xxx_messageInfo_RbacConfig_Target proto.InternalMessageInfo

func (m *RbacConfig_Target) GetServices() []string {
	if m != nil {
		return m.Services
	}
	return nil
}

func (m *RbacConfig_Target) GetNamespaces() []string {
	if m != nil {
		return m.Namespaces
	}
	return nil
}

func init() {
	proto.RegisterType((*ServiceRole)(nil), "istio.rbac.v1alpha1.ServiceRole")
	proto.RegisterType((*AccessRule)(nil), "istio.rbac.v1alpha1.AccessRule")
	proto.RegisterType((*AccessRule_Constraint)(nil), "istio.rbac.v1alpha1.AccessRule.Constraint")
	proto.RegisterType((*ServiceRoleBinding)(nil), "istio.rbac.v1alpha1.ServiceRoleBinding")
	proto.RegisterType((*Subject)(nil), "istio.rbac.v1alpha1.Subject")
	proto.RegisterMapType((map[string]string)(nil), "istio.rbac.v1alpha1.Subject.PropertiesEntry")
	proto.RegisterType((*RoleRef)(nil), "istio.rbac.v1alpha1.RoleRef")
	proto.RegisterType((*RbacConfig)(nil), "istio.rbac.v1alpha1.RbacConfig")
	proto.RegisterType((*RbacConfig_Target)(nil), "istio.rbac.v1alpha1.RbacConfig.Target")
	proto.RegisterEnum("istio.rbac.v1alpha1.RbacConfig_Mode", RbacConfig_Mode_name, RbacConfig_Mode_value)
}

func init() { proto.RegisterFile("rbac/v1alpha1/rbac.proto", fileDescriptor_rbac_665e05e3b5f24d86) }

var fileDescriptor_rbac_665e05e3b5f24d86 = []byte{
	// 530 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xd1, 0x6a, 0x13, 0x41,
	0x14, 0x75, 0x77, 0xd3, 0xa4, 0xb9, 0x01, 0x8d, 0x63, 0x95, 0x25, 0x14, 0x0d, 0x8b, 0x48, 0x11,
	0xd9, 0x90, 0x88, 0xa5, 0x08, 0x3e, 0xd8, 0x24, 0xc5, 0x40, 0x4c, 0x64, 0x52, 0x51, 0x7c, 0x29,
	0x93, 0xcd, 0x34, 0x19, 0xbb, 0xd9, 0x59, 0x66, 0x76, 0x83, 0xfd, 0x01, 0xbf, 0xcb, 0x2f, 0xd1,
	0x5f, 0x91, 0x99, 0xc9, 0x6c, 0xa2, 0xc4, 0x16, 0xdf, 0xee, 0x3d, 0xf7, 0x9e, 0xc3, 0xb9, 0x87,
	0xdd, 0x01, 0x5f, 0x4c, 0x49, 0xd4, 0x5a, 0xb5, 0x49, 0x9c, 0x2e, 0x48, 0xbb, 0xa5, 0xba, 0x30,
	0x15, 0x3c, 0xe3, 0xe8, 0x01, 0x93, 0x19, 0xe3, 0xa1, 0x46, 0xec, 0x3c, 0xe8, 0x41, 0x6d, 0x42,
	0xc5, 0x8a, 0x45, 0x14, 0xf3, 0x98, 0xa2, 0x57, 0xb0, 0x27, 0xf2, 0x98, 0x4a, 0xdf, 0x69, 0x7a,
	0x47, 0xb5, 0xce, 0x93, 0x70, 0x07, 0x27, 0x7c, 0x1b, 0x45, 0x54, 0x4a, 0x9c, 0xc7, 0x14, 0x9b,
	0xed, 0xe0, 0xa7, 0x03, 0xb0, 0x41, 0x51, 0x03, 0xf6, 0xa5, 0x11, 0x35, 0x42, 0x55, 0x5c, 0xf4,
	0xe8, 0x00, 0xf6, 0x52, 0x92, 0x2d, 0xa4, 0xef, 0xea, 0x81, 0x69, 0x90, 0x0f, 0x95, 0x25, 0xcd,
	0x16, 0x7c, 0x26, 0x7d, 0x4f, 0xe3, 0xb6, 0x45, 0x43, 0xa8, 0x45, 0x3c, 0x91, 0x99, 0x20, 0x2c,
	0xc9, 0xa4, 0x5f, 0xd2, 0xbe, 0x9e, 0xdf, 0xe2, 0x2b, 0xec, 0x16, 0x14, 0xbc, 0x4d, 0x6f, 0x1c,
	0x03, 0x6c, 0x46, 0xa8, 0x0e, 0xde, 0x15, 0xbd, 0xf6, 0x9d, 0xa6, 0x73, 0x54, 0xc5, 0xaa, 0x44,
	0x8f, 0xa0, 0xbc, 0x22, 0x71, 0x4e, 0xad, 0xbd, 0x75, 0x17, 0x7c, 0x77, 0x00, 0x6d, 0xe5, 0x74,
	0xca, 0x92, 0x19, 0x4b, 0xe6, 0xe8, 0x04, 0xf6, 0x65, 0x3e, 0xfd, 0x4a, 0xa3, 0xcc, 0x26, 0x76,
	0xb8, 0xd3, 0xd9, 0xc4, 0x2c, 0xe1, 0x62, 0x1b, 0x1d, 0x43, 0x45, 0xf0, 0x98, 0x62, 0x7a, 0xe9,
	0xbb, 0x4d, 0xe7, 0x9f, 0x44, 0x6c, 0x76, 0xb0, 0x5d, 0x0e, 0x7e, 0x38, 0x50, 0x59, 0xab, 0x21,
	0x04, 0xa5, 0x5c, 0x52, 0xb1, 0xf6, 0xaf, 0x6b, 0x15, 0xef, 0x5c, 0xf0, 0x3c, 0xd5, 0xaa, 0x55,
	0x6c, 0x1a, 0x34, 0x04, 0x48, 0x05, 0x4f, 0xa9, 0xc8, 0x18, 0x35, 0x09, 0xd7, 0x3a, 0x2f, 0x6e,
	0x72, 0x1a, 0x7e, 0x28, 0xd6, 0xfb, 0x49, 0x26, 0xae, 0xf1, 0x16, 0xbf, 0xf1, 0x06, 0xee, 0xfd,
	0x35, 0xde, 0x91, 0xe4, 0x01, 0xec, 0xe9, 0xec, 0xac, 0x11, 0xdd, 0xbc, 0x76, 0x4f, 0x9c, 0xa0,
	0x0d, 0x95, 0xf5, 0x59, 0xea, 0x82, 0x2b, 0x96, 0xcc, 0xec, 0x05, 0xaa, 0x56, 0x58, 0x42, 0x96,
	0x96, 0xa7, 0xeb, 0xe0, 0x97, 0x0b, 0x80, 0xa7, 0x24, 0xea, 0xf2, 0xe4, 0x92, 0xa9, 0xd8, 0x4b,
	0x4b, 0x3e, 0xa3, 0x9a, 0x76, 0xb7, 0xf3, 0x74, 0x77, 0x72, 0xc5, 0x7a, 0xf8, 0x9e, 0xcf, 0x28,
	0xd6, 0x0c, 0xd4, 0x83, 0x2a, 0x4b, 0xa2, 0x38, 0x97, 0x8c, 0x27, 0xeb, 0xe0, 0x9f, 0xdd, 0x46,
	0x3f, 0x27, 0x62, 0x4e, 0x33, 0xbc, 0x21, 0x2a, 0x15, 0xfa, 0xcd, 0xaa, 0x78, 0xff, 0xa7, 0x52,
	0x10, 0x1b, 0x3d, 0x28, 0x1b, 0xf0, 0xc6, 0xff, 0xe5, 0x31, 0x80, 0x8a, 0x40, 0xa6, 0x24, 0x2a,
	0xbe, 0xca, 0x2d, 0x24, 0xe8, 0x43, 0x49, 0xdd, 0x87, 0x2a, 0xe0, 0x8d, 0xcf, 0xce, 0xea, 0x77,
	0x50, 0x19, 0xdc, 0xf1, 0xa8, 0xee, 0xa0, 0x87, 0x70, 0x7f, 0x3c, 0xba, 0xf8, 0x34, 0x38, 0x7f,
	0x77, 0x31, 0x18, 0x75, 0x87, 0x1f, 0x27, 0x83, 0xf1, 0xa8, 0xee, 0x6e, 0xc3, 0xfd, 0xcf, 0x16,
	0xf6, 0x4e, 0x0f, 0xbf, 0x34, 0xcc, 0x01, 0x8c, 0xb7, 0x48, 0xca, 0x5a, 0x7f, 0xbc, 0x22, 0xd3,
	0xb2, 0x7e, 0x41, 0x5e, 0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xcc, 0xa7, 0x9a, 0xb8, 0x5d, 0x04,
	0x00, 0x00,
}
