// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mesh/v1alpha1/config.proto

/*
Package v1alpha1 is a generated protocol buffer package.

It is generated from these files:
	mesh/v1alpha1/config.proto

It has these top-level messages:
	ServerAddress
	ProxyConfig
	MeshConfig
*/
package v1alpha1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/duration"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// AuthenticationPolicy defines authentication policy. It can be set for
// different scopes (mesh, service …), and the most narrow scope with
// non-INHERIT value will be used.
// Mesh policy cannot be INHERIT.
type AuthenticationPolicy int32

const (
	// Do not encrypt Envoy to Envoy traffic.
	AuthenticationPolicy_NONE AuthenticationPolicy = 0
	// Envoy to Envoy traffic is wrapped into mutual TLS connections.
	AuthenticationPolicy_MUTUAL_TLS AuthenticationPolicy = 1
	// Use the policy defined by the parent scope. Should not be used for mesh
	// policy.
	AuthenticationPolicy_INHERIT AuthenticationPolicy = 1000
)

var AuthenticationPolicy_name = map[int32]string{
	0:    "NONE",
	1:    "MUTUAL_TLS",
	1000: "INHERIT",
}
var AuthenticationPolicy_value = map[string]int32{
	"NONE":       0,
	"MUTUAL_TLS": 1,
	"INHERIT":    1000,
}

func (x AuthenticationPolicy) String() string {
	return proto.EnumName(AuthenticationPolicy_name, int32(x))
}
func (AuthenticationPolicy) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type MeshConfig_IngressControllerMode int32

const (
	// Disables Istio ingress controller.
	MeshConfig_OFF MeshConfig_IngressControllerMode = 0
	// Istio ingress controller will act on ingress resources that do not
	// contain any annotation or whose annotations match the value
	// specified in the ingress_class parameter described earlier. Use this
	// mode if Istio ingress controller will be the default ingress
	// controller for the entire kubernetes cluster.
	MeshConfig_DEFAULT MeshConfig_IngressControllerMode = 1
	// Istio ingress controller will only act on ingress resources whose
	// annotations match the value specified in the ingress_class parameter
	// described earlier. Use this mode if Istio ingress controller will be
	// a secondary ingress controller (e.g., in addition to a
	// cloud-provided ingress controller).
	MeshConfig_STRICT MeshConfig_IngressControllerMode = 2
)

var MeshConfig_IngressControllerMode_name = map[int32]string{
	0: "OFF",
	1: "DEFAULT",
	2: "STRICT",
}
var MeshConfig_IngressControllerMode_value = map[string]int32{
	"OFF":     0,
	"DEFAULT": 1,
	"STRICT":  2,
}

func (x MeshConfig_IngressControllerMode) String() string {
	return proto.EnumName(MeshConfig_IngressControllerMode_name, int32(x))
}
func (MeshConfig_IngressControllerMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2, 0}
}

// TODO AuthPolicy needs to be removed and merged with AuthPolicy defined above
type MeshConfig_AuthPolicy int32

const (
	// Do not encrypt Envoy to Envoy traffic.
	MeshConfig_NONE MeshConfig_AuthPolicy = 0
	// Envoy to Envoy traffic is wrapped into mutual TLS connections.
	MeshConfig_MUTUAL_TLS MeshConfig_AuthPolicy = 1
)

var MeshConfig_AuthPolicy_name = map[int32]string{
	0: "NONE",
	1: "MUTUAL_TLS",
}
var MeshConfig_AuthPolicy_value = map[string]int32{
	"NONE":       0,
	"MUTUAL_TLS": 1,
}

func (x MeshConfig_AuthPolicy) String() string {
	return proto.EnumName(MeshConfig_AuthPolicy_name, int32(x))
}
func (MeshConfig_AuthPolicy) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 1} }

// ServerAddress specifies the address of Istio components like mixer, pilot, etc.
// At least one of the field needs to be specified.
type ServerAddress struct {
	// The address for mTLS server, e.g., (_istio-pilot:15003_)
	MutualTls string `protobuf:"bytes,1,opt,name=mutual_tls,json=mutualTls" json:"mutual_tls,omitempty"`
	// The address for plain text server, e.g., (_istio-pilot:15005_)
	PlainText string `protobuf:"bytes,2,opt,name=plain_text,json=plainText" json:"plain_text,omitempty"`
}

func (m *ServerAddress) Reset()                    { *m = ServerAddress{} }
func (m *ServerAddress) String() string            { return proto.CompactTextString(m) }
func (*ServerAddress) ProtoMessage()               {}
func (*ServerAddress) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServerAddress) GetMutualTls() string {
	if m != nil {
		return m.MutualTls
	}
	return ""
}

func (m *ServerAddress) GetPlainText() string {
	if m != nil {
		return m.PlainText
	}
	return ""
}

// ProxyConfig defines variables for individual Envoy instances.
type ProxyConfig struct {
	// Path to the generated configuration file directory.
	// Proxy agent generates the actual configuration and stores it in this directory.
	ConfigPath string `protobuf:"bytes,1,opt,name=config_path,json=configPath" json:"config_path,omitempty"`
	// Path to the proxy binary
	BinaryPath string `protobuf:"bytes,2,opt,name=binary_path,json=binaryPath" json:"binary_path,omitempty"`
	// Service cluster defines the name for the service_cluster that is
	// shared by all Envoy instances. This setting corresponds to
	// _--service-cluster_ flag in Envoy.  In a typical Envoy deployment, the
	// _service-cluster_ flag is used to identify the caller, for
	// source-based routing scenarios.
	//
	// Since Istio does not assign a local service/service version to each
	// Envoy instance, the name is same for all of them.  However, the
	// source/caller's identity (e.g., IP address) is encoded in the
	// _--service-node_ flag when launching Envoy.  When the RDS service
	// receives API calls from Envoy, it uses the value of the _service-node_
	// flag to compute routes that are relative to the service instances
	// located at that IP address.
	ServiceCluster string `protobuf:"bytes,3,opt,name=service_cluster,json=serviceCluster" json:"service_cluster,omitempty"`
	// The time in seconds that Envoy will drain connections during a hot
	// restart. MUST be >=1s (e.g., _1s/1m/1h_)
	DrainDuration *google_protobuf.Duration `protobuf:"bytes,4,opt,name=drain_duration,json=drainDuration" json:"drain_duration,omitempty"`
	// The time in seconds that Envoy will wait before shutting down the
	// parent process during a hot restart. MUST be >=1s (e.g., _1s/1m/1h_).
	// MUST BE greater than _drain_duration_ parameter.
	ParentShutdownDuration *google_protobuf.Duration `protobuf:"bytes,5,opt,name=parent_shutdown_duration,json=parentShutdownDuration" json:"parent_shutdown_duration,omitempty"`
	// Deprecated, use server_address instead.
	DiscoveryAddress string `protobuf:"bytes,6,opt,name=discovery_address,json=discoveryAddress" json:"discovery_address,omitempty"`
	// Polling interval for service discovery (used by EDS, CDS, LDS, but not RDS). (MUST BE >=1ms)
	DiscoveryRefreshDelay *google_protobuf.Duration `protobuf:"bytes,7,opt,name=discovery_refresh_delay,json=discoveryRefreshDelay" json:"discovery_refresh_delay,omitempty"`
	// Address of the Zipkin service (e.g. _zipkin:9411_).
	ZipkinAddress string `protobuf:"bytes,8,opt,name=zipkin_address,json=zipkinAddress" json:"zipkin_address,omitempty"`
	// Connection timeout used by Envoy for supporting services. (MUST BE >=1ms)
	ConnectTimeout *google_protobuf.Duration `protobuf:"bytes,9,opt,name=connect_timeout,json=connectTimeout" json:"connect_timeout,omitempty"`
	// IP Address and Port of a statsd UDP listener (e.g. _10.75.241.127:9125_).
	StatsdUdpAddress string `protobuf:"bytes,10,opt,name=statsd_udp_address,json=statsdUdpAddress" json:"statsd_udp_address,omitempty"`
	// Port on which Envoy should listen for administrative commands.
	ProxyAdminPort int32 `protobuf:"varint,11,opt,name=proxy_admin_port,json=proxyAdminPort" json:"proxy_admin_port,omitempty"`
	// The availability zone where this Envoy instance is running. When running
	// Envoy as a sidecar in Kubernetes, this flag must be one of the availability
	// zones assigned to a node using failure-domain.beta.kubernetes.io/zone annotation.
	AvailabilityZone string `protobuf:"bytes,12,opt,name=availability_zone,json=availabilityZone" json:"availability_zone,omitempty"`
	// Authentication policy defines the global switch to control authentication
	// for Envoy-to-Envoy communication for istio components Mixer and Pilot.
	ControlPlaneAuthPolicy AuthenticationPolicy `protobuf:"varint,13,opt,name=control_plane_auth_policy,json=controlPlaneAuthPolicy,enum=istio.mesh.v1alpha1.AuthenticationPolicy" json:"control_plane_auth_policy,omitempty"`
	// File path of custom proxy configuration, currently used by proxies
	// in front of Mixer and Pilot.
	CustomConfigFile string `protobuf:"bytes,14,opt,name=custom_config_file,json=customConfigFile" json:"custom_config_file,omitempty"`
	// Maximum length of name field in Envoy's metrics. The length of the name field
	// is determined by the length of a name field in a service and the set of labels that
	// comprise a particular version of the service. The default value is set to 189 characters.
	// Envoy's internal metrics take up 67 characters, for a total of 256 character name per metric.
	// Increase the value of this field if you find that the metrics from Envoys are truncated.
	StatNameLength int32 `protobuf:"varint,15,opt,name=stat_name_length,json=statNameLength" json:"stat_name_length,omitempty"`
	// The number of worker threads to run. Default value is number of cores on the machine.
	Concurrency int32 `protobuf:"varint,16,opt,name=concurrency" json:"concurrency,omitempty"`
	// Address of the discovery service exposing xDS.
	Pilot *ServerAddress `protobuf:"bytes,17,opt,name=pilot" json:"pilot,omitempty"`
}

func (m *ProxyConfig) Reset()                    { *m = ProxyConfig{} }
func (m *ProxyConfig) String() string            { return proto.CompactTextString(m) }
func (*ProxyConfig) ProtoMessage()               {}
func (*ProxyConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProxyConfig) GetConfigPath() string {
	if m != nil {
		return m.ConfigPath
	}
	return ""
}

func (m *ProxyConfig) GetBinaryPath() string {
	if m != nil {
		return m.BinaryPath
	}
	return ""
}

func (m *ProxyConfig) GetServiceCluster() string {
	if m != nil {
		return m.ServiceCluster
	}
	return ""
}

func (m *ProxyConfig) GetDrainDuration() *google_protobuf.Duration {
	if m != nil {
		return m.DrainDuration
	}
	return nil
}

func (m *ProxyConfig) GetParentShutdownDuration() *google_protobuf.Duration {
	if m != nil {
		return m.ParentShutdownDuration
	}
	return nil
}

func (m *ProxyConfig) GetDiscoveryAddress() string {
	if m != nil {
		return m.DiscoveryAddress
	}
	return ""
}

func (m *ProxyConfig) GetDiscoveryRefreshDelay() *google_protobuf.Duration {
	if m != nil {
		return m.DiscoveryRefreshDelay
	}
	return nil
}

func (m *ProxyConfig) GetZipkinAddress() string {
	if m != nil {
		return m.ZipkinAddress
	}
	return ""
}

func (m *ProxyConfig) GetConnectTimeout() *google_protobuf.Duration {
	if m != nil {
		return m.ConnectTimeout
	}
	return nil
}

func (m *ProxyConfig) GetStatsdUdpAddress() string {
	if m != nil {
		return m.StatsdUdpAddress
	}
	return ""
}

func (m *ProxyConfig) GetProxyAdminPort() int32 {
	if m != nil {
		return m.ProxyAdminPort
	}
	return 0
}

func (m *ProxyConfig) GetAvailabilityZone() string {
	if m != nil {
		return m.AvailabilityZone
	}
	return ""
}

func (m *ProxyConfig) GetControlPlaneAuthPolicy() AuthenticationPolicy {
	if m != nil {
		return m.ControlPlaneAuthPolicy
	}
	return AuthenticationPolicy_NONE
}

func (m *ProxyConfig) GetCustomConfigFile() string {
	if m != nil {
		return m.CustomConfigFile
	}
	return ""
}

func (m *ProxyConfig) GetStatNameLength() int32 {
	if m != nil {
		return m.StatNameLength
	}
	return 0
}

func (m *ProxyConfig) GetConcurrency() int32 {
	if m != nil {
		return m.Concurrency
	}
	return 0
}

func (m *ProxyConfig) GetPilot() *ServerAddress {
	if m != nil {
		return m.Pilot
	}
	return nil
}

// MeshConfig defines mesh-wide variables shared by all Envoy instances in the
// Istio service mesh.
type MeshConfig struct {
	// Deprecated, use mixer_check instead.
	MixerCheckServer string `protobuf:"bytes,1,opt,name=mixer_check_server,json=mixerCheckServer" json:"mixer_check_server,omitempty"`
	// Deprecated, use mixer_report instead.
	MixerReportServer string `protobuf:"bytes,2,opt,name=mixer_report_server,json=mixerReportServer" json:"mixer_report_server,omitempty"`
	// Disable policy checks by the mixer service. Default
	// is false, i.e. mixer policy check is enabled by default.
	DisablePolicyChecks bool `protobuf:"varint,3,opt,name=disable_policy_checks,json=disablePolicyChecks" json:"disable_policy_checks,omitempty"`
	// Port on which Envoy should listen for incoming connections from
	// other services.
	ProxyListenPort int32 `protobuf:"varint,4,opt,name=proxy_listen_port,json=proxyListenPort" json:"proxy_listen_port,omitempty"`
	// Port on which Envoy should listen for HTTP PROXY requests if set.
	ProxyHttpPort int32 `protobuf:"varint,5,opt,name=proxy_http_port,json=proxyHttpPort" json:"proxy_http_port,omitempty"`
	// Connection timeout used by Envoy. (MUST BE >=1ms)
	ConnectTimeout *google_protobuf.Duration `protobuf:"bytes,6,opt,name=connect_timeout,json=connectTimeout" json:"connect_timeout,omitempty"`
	// Class of ingress resources to be processed by Istio ingress
	// controller.  This corresponds to the value of
	// "kubernetes.io/ingress.class" annotation.
	IngressClass string `protobuf:"bytes,7,opt,name=ingress_class,json=ingressClass" json:"ingress_class,omitempty"`
	// Name of the kubernetes service used for the istio ingress controller.
	IngressService string `protobuf:"bytes,8,opt,name=ingress_service,json=ingressService" json:"ingress_service,omitempty"`
	// Defines whether to use Istio ingress controller for annotated or all ingress resources.
	IngressControllerMode MeshConfig_IngressControllerMode `protobuf:"varint,9,opt,name=ingress_controller_mode,json=ingressControllerMode,enum=istio.mesh.v1alpha1.MeshConfig_IngressControllerMode" json:"ingress_controller_mode,omitempty"`
	// Authentication policy defines the global switch to control authentication
	// for Envoy-to-Envoy communication.
	// Use authentication_policy instead.
	AuthPolicy MeshConfig_AuthPolicy `protobuf:"varint,10,opt,name=auth_policy,json=authPolicy,enum=istio.mesh.v1alpha1.MeshConfig_AuthPolicy" json:"auth_policy,omitempty"`
	// Polling interval for RDS (MUST BE >=1ms)
	RdsRefreshDelay *google_protobuf.Duration `protobuf:"bytes,11,opt,name=rds_refresh_delay,json=rdsRefreshDelay" json:"rds_refresh_delay,omitempty"`
	// Flag to control generation of trace spans and request IDs.
	// Requires a trace span collector defined in the proxy configuration.
	EnableTracing bool `protobuf:"varint,12,opt,name=enable_tracing,json=enableTracing" json:"enable_tracing,omitempty"`
	// File address for the proxy access log (e.g. /dev/stdout).
	// Empty value disables access logging.
	AccessLogFile string `protobuf:"bytes,13,opt,name=access_log_file,json=accessLogFile" json:"access_log_file,omitempty"`
	// Default proxy config used by the proxy injection mechanism operating in the mesh
	// (e.g. Kubernetes admission controller)
	// In case of Kubernetes, the proxy config is applied once during the injection process,
	// and remain constant for the duration of the pod. The rest of the mesh config can be changed
	// at runtime and config gets distributed dynamically.
	DefaultConfig *ProxyConfig `protobuf:"bytes,14,opt,name=default_config,json=defaultConfig" json:"default_config,omitempty"`
	// List of remote services for which mTLS authentication should not be expected by Istio .
	// Typically, these are control services (e.g kubernetes API server) that don't have Istio sidecar
	// to transparently terminate mTLS authentication.
	// It has no effect if the authentication policy is already 'NONE'.
	// DO NOT use this setting for services that are managed by Istio (i.e. using Istio sidecar).
	// Instead, use service-level annotations to overwrite the authentication policy.
	MtlsExcludedServices []string `protobuf:"bytes,15,rep,name=mtls_excluded_services,json=mtlsExcludedServices" json:"mtls_excluded_services,omitempty"`
	// DEPRECATED. Mixer address. This option will be removed soon. Please
	// use mixer_check and mixer_report.
	MixerAddress string `protobuf:"bytes,16,opt,name=mixer_address,json=mixerAddress" json:"mixer_address,omitempty"`
	// Address of the server that will be used by the proxies for policy
	// check calls. By using different names for mixerCheck and mixerReport, it
	// is possible to have one set of mixer servers handle policy check calls,
	// while another set of mixer servers handle telemetry calls.
	//
	// NOTE: Omitting mixerCheck while specifying mixerReport is
	// equivalent to setting disablePolicyChecks to true.
	MixerCheck *ServerAddress `protobuf:"bytes,17,opt,name=mixer_check,json=mixerCheck" json:"mixer_check,omitempty"`
	// Address of the server that will be used by the proxies for policy report
	// calls.
	MixerReport *ServerAddress `protobuf:"bytes,18,opt,name=mixer_report,json=mixerReport" json:"mixer_report,omitempty"`
}

func (m *MeshConfig) Reset()                    { *m = MeshConfig{} }
func (m *MeshConfig) String() string            { return proto.CompactTextString(m) }
func (*MeshConfig) ProtoMessage()               {}
func (*MeshConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MeshConfig) GetMixerCheckServer() string {
	if m != nil {
		return m.MixerCheckServer
	}
	return ""
}

func (m *MeshConfig) GetMixerReportServer() string {
	if m != nil {
		return m.MixerReportServer
	}
	return ""
}

func (m *MeshConfig) GetDisablePolicyChecks() bool {
	if m != nil {
		return m.DisablePolicyChecks
	}
	return false
}

func (m *MeshConfig) GetProxyListenPort() int32 {
	if m != nil {
		return m.ProxyListenPort
	}
	return 0
}

func (m *MeshConfig) GetProxyHttpPort() int32 {
	if m != nil {
		return m.ProxyHttpPort
	}
	return 0
}

func (m *MeshConfig) GetConnectTimeout() *google_protobuf.Duration {
	if m != nil {
		return m.ConnectTimeout
	}
	return nil
}

func (m *MeshConfig) GetIngressClass() string {
	if m != nil {
		return m.IngressClass
	}
	return ""
}

func (m *MeshConfig) GetIngressService() string {
	if m != nil {
		return m.IngressService
	}
	return ""
}

func (m *MeshConfig) GetIngressControllerMode() MeshConfig_IngressControllerMode {
	if m != nil {
		return m.IngressControllerMode
	}
	return MeshConfig_OFF
}

func (m *MeshConfig) GetAuthPolicy() MeshConfig_AuthPolicy {
	if m != nil {
		return m.AuthPolicy
	}
	return MeshConfig_NONE
}

func (m *MeshConfig) GetRdsRefreshDelay() *google_protobuf.Duration {
	if m != nil {
		return m.RdsRefreshDelay
	}
	return nil
}

func (m *MeshConfig) GetEnableTracing() bool {
	if m != nil {
		return m.EnableTracing
	}
	return false
}

func (m *MeshConfig) GetAccessLogFile() string {
	if m != nil {
		return m.AccessLogFile
	}
	return ""
}

func (m *MeshConfig) GetDefaultConfig() *ProxyConfig {
	if m != nil {
		return m.DefaultConfig
	}
	return nil
}

func (m *MeshConfig) GetMtlsExcludedServices() []string {
	if m != nil {
		return m.MtlsExcludedServices
	}
	return nil
}

func (m *MeshConfig) GetMixerAddress() string {
	if m != nil {
		return m.MixerAddress
	}
	return ""
}

func (m *MeshConfig) GetMixerCheck() *ServerAddress {
	if m != nil {
		return m.MixerCheck
	}
	return nil
}

func (m *MeshConfig) GetMixerReport() *ServerAddress {
	if m != nil {
		return m.MixerReport
	}
	return nil
}

func init() {
	proto.RegisterType((*ServerAddress)(nil), "istio.mesh.v1alpha1.ServerAddress")
	proto.RegisterType((*ProxyConfig)(nil), "istio.mesh.v1alpha1.ProxyConfig")
	proto.RegisterType((*MeshConfig)(nil), "istio.mesh.v1alpha1.MeshConfig")
	proto.RegisterEnum("istio.mesh.v1alpha1.AuthenticationPolicy", AuthenticationPolicy_name, AuthenticationPolicy_value)
	proto.RegisterEnum("istio.mesh.v1alpha1.MeshConfig_IngressControllerMode", MeshConfig_IngressControllerMode_name, MeshConfig_IngressControllerMode_value)
	proto.RegisterEnum("istio.mesh.v1alpha1.MeshConfig_AuthPolicy", MeshConfig_AuthPolicy_name, MeshConfig_AuthPolicy_value)
}

func init() { proto.RegisterFile("mesh/v1alpha1/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1043 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x56, 0x5b, 0x4f, 0x23, 0x47,
	0x13, 0x5d, 0x73, 0x33, 0x94, 0xf1, 0x85, 0x66, 0x61, 0x67, 0xd1, 0xf7, 0x25, 0x96, 0xa3, 0x6c,
	0x1c, 0xb2, 0xb2, 0xb5, 0x24, 0x91, 0x92, 0x97, 0x28, 0x60, 0x4c, 0x16, 0xc9, 0x5c, 0x32, 0x36,
	0x2f, 0xfb, 0xd2, 0x6a, 0x66, 0x1a, 0x4f, 0x6b, 0x7b, 0xa6, 0x47, 0xdd, 0x3d, 0x04, 0xef, 0x3f,
	0xc9, 0x3f, 0xcc, 0x4b, 0xfe, 0x43, 0xd4, 0x97, 0x31, 0xde, 0x95, 0x23, 0xb2, 0x8f, 0x9c, 0x3a,
	0x55, 0x35, 0x55, 0x7d, 0x4e, 0x61, 0x38, 0x48, 0xa9, 0x4a, 0xfa, 0xf7, 0x6f, 0x08, 0xcf, 0x13,
	0xf2, 0xa6, 0x1f, 0x89, 0xec, 0x8e, 0x4d, 0x7b, 0xb9, 0x14, 0x5a, 0xa0, 0x5d, 0xa6, 0x34, 0x13,
	0x3d, 0xc3, 0xe8, 0x95, 0x8c, 0x83, 0x2f, 0xa6, 0x42, 0x4c, 0x39, 0xed, 0x5b, 0xca, 0x6d, 0x71,
	0xd7, 0x8f, 0x0b, 0x49, 0x34, 0x13, 0x99, 0x4b, 0xea, 0x5c, 0x40, 0x7d, 0x4c, 0xe5, 0x3d, 0x95,
	0xc7, 0x71, 0x2c, 0xa9, 0x52, 0xe8, 0xff, 0x00, 0x69, 0xa1, 0x0b, 0xc2, 0xb1, 0xe6, 0x2a, 0xa8,
	0xb4, 0x2b, 0xdd, 0xad, 0x70, 0xcb, 0x21, 0x13, 0x6e, 0xc3, 0x39, 0x27, 0x2c, 0xc3, 0x9a, 0x3e,
	0xe8, 0x60, 0xc5, 0x85, 0x2d, 0x32, 0xa1, 0x0f, 0xba, 0xf3, 0x67, 0x15, 0x6a, 0xd7, 0x52, 0x3c,
	0xcc, 0x06, 0xf6, 0xcb, 0xd0, 0x97, 0x50, 0x73, 0xdf, 0x88, 0x73, 0xa2, 0x13, 0x5f, 0x0e, 0x1c,
	0x74, 0x4d, 0x74, 0x62, 0x08, 0xb7, 0x2c, 0x23, 0x72, 0xe6, 0x08, 0xae, 0x20, 0x38, 0xc8, 0x12,
	0xbe, 0x81, 0xa6, 0xa2, 0xf2, 0x9e, 0x45, 0x14, 0x47, 0xbc, 0x50, 0x9a, 0xca, 0x60, 0xd5, 0x92,
	0x1a, 0x1e, 0x1e, 0x38, 0x14, 0xfd, 0x0a, 0x8d, 0x58, 0x9a, 0x2f, 0x2b, 0x27, 0x0c, 0xd6, 0xda,
	0x95, 0x6e, 0xed, 0xe8, 0x65, 0xcf, 0xad, 0xa0, 0x57, 0xae, 0xa0, 0x77, 0xea, 0x09, 0x61, 0xdd,
	0x26, 0x94, 0x7f, 0xa2, 0x31, 0x04, 0x39, 0x91, 0x34, 0xd3, 0x58, 0x25, 0x85, 0x8e, 0xc5, 0x1f,
	0x0b, 0xb5, 0xd6, 0x9f, 0xaa, 0xb5, 0xef, 0x52, 0xc7, 0x3e, 0x73, 0x5e, 0xf4, 0x3b, 0xd8, 0x89,
	0x99, 0x8a, 0xc4, 0x3d, 0x95, 0x33, 0x4c, 0xdc, 0x92, 0x83, 0x0d, 0x3b, 0x41, 0x6b, 0x1e, 0x28,
	0x97, 0xff, 0x3b, 0xbc, 0x78, 0x24, 0x4b, 0x7a, 0x27, 0xa9, 0x4a, 0x70, 0x4c, 0x39, 0x99, 0x05,
	0xd5, 0xa7, 0x3e, 0x60, 0x6f, 0x9e, 0x19, 0xba, 0xc4, 0x53, 0x93, 0x87, 0xbe, 0x86, 0xc6, 0x07,
	0x96, 0xbf, 0x67, 0xd9, 0xbc, 0xf9, 0xa6, 0x6d, 0x5e, 0x77, 0x68, 0xd9, 0xf9, 0x04, 0x9a, 0x91,
	0xc8, 0x32, 0x1a, 0x69, 0xac, 0x59, 0x4a, 0x45, 0xa1, 0x83, 0xad, 0xa7, 0x3a, 0x36, 0x7c, 0xc6,
	0xc4, 0x25, 0xa0, 0xd7, 0x80, 0x94, 0x26, 0x5a, 0xc5, 0xb8, 0x88, 0xf3, 0x79, 0x3b, 0x70, 0xb3,
	0xba, 0xc8, 0x4d, 0x9c, 0x97, 0x1d, 0xbb, 0xd0, 0xca, 0x8d, 0x52, 0x30, 0x89, 0x53, 0x96, 0xe1,
	0x5c, 0x48, 0x1d, 0xd4, 0xda, 0x95, 0xee, 0x7a, 0xd8, 0xb0, 0xf8, 0xb1, 0x81, 0xaf, 0x85, 0xd4,
	0x66, 0x85, 0xe4, 0x9e, 0x30, 0x4e, 0x6e, 0x19, 0x67, 0x7a, 0x86, 0x3f, 0x88, 0x8c, 0x06, 0xdb,
	0xae, 0xec, 0x62, 0xe0, 0x9d, 0xc8, 0x28, 0x8a, 0xe1, 0x65, 0x24, 0x32, 0x2d, 0x05, 0xc7, 0x39,
	0x27, 0x19, 0xc5, 0xa4, 0xd0, 0x09, 0xce, 0x05, 0x67, 0xd1, 0x2c, 0xa8, 0xb7, 0x2b, 0xdd, 0xc6,
	0xd1, 0xb7, 0xbd, 0x25, 0x4e, 0xe9, 0x1d, 0x17, 0x3a, 0xa1, 0x99, 0x66, 0x91, 0x1d, 0xee, 0xda,
	0x26, 0x84, 0xfb, 0xbe, 0xd6, 0xb5, 0x29, 0x65, 0x18, 0x0e, 0x37, 0xa3, 0x46, 0x85, 0xd2, 0x22,
	0xc5, 0x5e, 0xde, 0x77, 0x8c, 0xd3, 0xa0, 0xe1, 0xbe, 0xc9, 0x45, 0x9c, 0x03, 0xce, 0x18, 0xa7,
	0x66, 0x54, 0x33, 0x3e, 0xce, 0x48, 0x4a, 0x31, 0xa7, 0xd9, 0x54, 0x27, 0x41, 0xd3, 0x8d, 0x6a,
	0xf0, 0x4b, 0x92, 0xd2, 0x91, 0x45, 0x51, 0xdb, 0xfa, 0x25, 0x2a, 0xa4, 0xa4, 0x59, 0x34, 0x0b,
	0x5a, 0x96, 0xb4, 0x08, 0xa1, 0x9f, 0x60, 0x3d, 0x67, 0x5c, 0xe8, 0x60, 0xc7, 0x3e, 0x4f, 0x67,
	0xe9, 0x2c, 0x1f, 0x59, 0x3a, 0x74, 0x09, 0x9d, 0xbf, 0x37, 0x01, 0x2e, 0xa8, 0x4a, 0xbc, 0x35,
	0x5f, 0x03, 0x4a, 0xd9, 0x03, 0x95, 0x38, 0x4a, 0x68, 0xf4, 0x1e, 0x2b, 0x9b, 0xe2, 0x1d, 0xda,
	0xb2, 0x91, 0x81, 0x09, 0xb8, 0x52, 0xa8, 0x07, 0xbb, 0x8e, 0x2d, 0xa9, 0x79, 0xa9, 0x92, 0xee,
	0xfc, 0xba, 0x63, 0x43, 0xa1, 0x8d, 0x78, 0xfe, 0x11, 0x18, 0x3d, 0x92, 0x5b, 0x4e, 0xfd, 0xee,
	0x5d, 0x1b, 0x65, 0xcd, 0xbb, 0x19, 0xee, 0xfa, 0xa0, 0x5b, 0xa7, 0x6d, 0xa4, 0xd0, 0x21, 0xec,
	0x38, 0x45, 0x70, 0xa6, 0x34, 0xf5, 0x92, 0x58, 0xb3, 0x2b, 0x68, 0xda, 0xc0, 0xc8, 0xe2, 0x56,
	0x13, 0xaf, 0xc0, 0x41, 0x38, 0xd1, 0x3a, 0x77, 0xcc, 0x75, 0xcb, 0xac, 0x5b, 0xf8, 0xad, 0xd6,
	0xb9, 0xe5, 0x2d, 0xd1, 0xf5, 0xc6, 0xe7, 0xea, 0xfa, 0x2b, 0xa8, 0xb3, 0x6c, 0x6a, 0x56, 0x89,
	0x23, 0x4e, 0x94, 0xb2, 0x5e, 0xdc, 0x0a, 0xb7, 0x3d, 0x38, 0x30, 0x98, 0xb9, 0x53, 0x25, 0xc9,
	0x1f, 0x26, 0x6f, 0xb4, 0x86, 0x87, 0xc7, 0x0e, 0x45, 0x29, 0xbc, 0x98, 0x57, 0x73, 0xe2, 0xe2,
	0x54, 0xe2, 0x54, 0xc4, 0xd4, 0x3a, 0xae, 0x71, 0xf4, 0xe3, 0xd2, 0x27, 0x7d, 0x7c, 0xb9, 0xde,
	0xb9, 0xef, 0x3b, 0xcf, 0xbe, 0x10, 0x31, 0x0d, 0xf7, 0xd8, 0x32, 0x18, 0x5d, 0x41, 0x6d, 0xd1,
	0x01, 0x60, 0x5b, 0x1c, 0x3e, 0xd5, 0xe2, 0x51, 0xea, 0x27, 0x2b, 0x41, 0x25, 0x04, 0xf2, 0x28,
	0xfd, 0x21, 0xec, 0xc8, 0x58, 0x7d, 0x72, 0x9d, 0x6a, 0x4f, 0xed, 0xb4, 0x29, 0x63, 0xf5, 0xe9,
	0x5d, 0xa2, 0x99, 0xd5, 0x87, 0x96, 0x24, 0x62, 0xd9, 0xd4, 0x3a, 0x7a, 0x33, 0xac, 0x3b, 0x74,
	0xe2, 0x40, 0xf3, 0xce, 0x24, 0x8a, 0xcc, 0xb2, 0xb8, 0xf0, 0x2e, 0xab, 0xbb, 0xfb, 0xe5, 0xe0,
	0x91, 0x70, 0x16, 0xfb, 0x0d, 0x1a, 0x31, 0xbd, 0x23, 0x05, 0xd7, 0xde, 0x91, 0xd6, 0x8c, 0xb5,
	0xa3, 0xf6, 0xd2, 0x49, 0x17, 0xfe, 0x45, 0x85, 0x75, 0x9f, 0xe7, 0x6d, 0xf1, 0x03, 0xec, 0xa7,
	0x9a, 0x2b, 0x4c, 0x1f, 0x22, 0x5e, 0xc4, 0x34, 0x2e, 0x5f, 0x53, 0x05, 0xcd, 0xf6, 0x6a, 0x77,
	0x2b, 0x7c, 0x6e, 0xa2, 0x43, 0x1f, 0xf4, 0x6f, 0xaa, 0x8c, 0x44, 0x9c, 0x3d, 0xca, 0xab, 0xd7,
	0x72, 0x12, 0xb1, 0x60, 0x79, 0xf1, 0x06, 0x50, 0x5b, 0x70, 0xdc, 0x67, 0x18, 0x18, 0x1e, 0xed,
	0x88, 0x86, 0xb0, 0xbd, 0x68, 0xc4, 0x00, 0xfd, 0xe7, 0x2a, 0xb5, 0x05, 0x97, 0x76, 0x7e, 0x86,
	0xbd, 0xa5, 0x32, 0x42, 0x55, 0x58, 0xbd, 0x3a, 0x3b, 0x6b, 0x3d, 0x43, 0x35, 0xa8, 0x9e, 0x0e,
	0xcf, 0x8e, 0x6f, 0x46, 0x93, 0x56, 0x05, 0x01, 0x6c, 0x8c, 0x27, 0xe1, 0xf9, 0x60, 0xd2, 0x5a,
	0xe9, 0xbc, 0x02, 0x58, 0xb8, 0x84, 0x9b, 0xb0, 0x76, 0x79, 0x75, 0x39, 0x6c, 0x3d, 0x43, 0x0d,
	0x80, 0x8b, 0x9b, 0xc9, 0xcd, 0xf1, 0x08, 0x4f, 0x46, 0xe3, 0x56, 0xe5, 0xf0, 0x17, 0x78, 0xbe,
	0xec, 0xa6, 0xfe, 0x7b, 0x06, 0xda, 0x86, 0xea, 0xf9, 0xe5, 0xdb, 0x61, 0x78, 0x3e, 0x69, 0xfd,
	0x55, 0x3d, 0xf9, 0xdf, 0xbb, 0x03, 0x37, 0x14, 0x13, 0x7d, 0x92, 0xb3, 0xfe, 0x47, 0x3f, 0x7d,
	0x6e, 0x37, 0xac, 0xc6, 0xbe, 0xff, 0x27, 0x00, 0x00, 0xff, 0xff, 0x74, 0xbd, 0x2a, 0x4d, 0x12,
	0x09, 0x00, 0x00,
}
