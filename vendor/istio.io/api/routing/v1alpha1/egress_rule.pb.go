// Code generated by protoc-gen-go. DO NOT EDIT.
// source: routing/v1alpha1/egress_rule.proto

package v1alpha1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Egress rules describe the properties of a service outside Istio. When transparent proxying
// is used, egress rules signify a white listed set of external services that microserves in the mesh
// are allowed to access. A subset of routing rules and all destination policies can be applied
// on the service targeted by an egress rule. TCP services and HTTP-based services can be expressed
// by an egress rule. The destination of an egress rule for HTTP-based services must be an IP or a domain name,
// optionally with a wildcard prefix  (e.g., *.foo.com). For TCP based services, the destination of an
// egress rule must be an IP or a block of IPs in CIDR notation.
//
// If TLS origination from the sidecar is desired, the protocol
// associated with the service port must be marked as HTTPS, and the service is expected to
// be accessed over HTTP (e.g., http://gmail.com:443). The sidecar will automatically upgrade
// the connection to TLS when initiating a connection with the external service.
//
// For example, the following egress rule describes the set of services hosted under the *.foo.com domain
//
//     kind: EgressRule
//     metadata:
//       name: foo-egress-rule
//     spec:
//       destination:
//         service: *.foo.com
//       ports:
//         - port: 80
//           protocol: http
//         - port: 443
//           protocol: https
//
// The following egress rule describes the set of services accessed by a block of IPs
//
//     kind: EgressRule
//     metadata:
//       name: bar-egress-rule
//     spec:
//       destination:
//         service: 92.198.174.192/27
//       ports:
//         - port: 111
//           protocol: tcp
//
type EgressRule struct {
	// REQUIRED: A domain name, optionally with a wildcard prefix, or an IP, or a block of IPs
	// associated with the external service.
	// ONLY the "service" field of "destination" will be taken into consideration. Name,
	// namespace, domain and labels are ignored. Routing rules and destination policies that
	// refer to these external services must have identical specification for the destination
	// as the corresponding egress rule.
	//
	// The "service" field of "destination" for HTTP-based services must be an IP or a domain name,
	// optionally with a wildcard prefix. Wildcard domain specifications must conform to format
	// allowed by Envoy's Virtual Host specification, such as “*.foo.com” or “*-bar.foo.com”.
	// The character '*' in a domain specification indicates a non-empty string. Hence, a wildcard
	// domain of form “*-bar.foo.com” will match “baz-bar.foo.com” but not “-bar.foo.com”.
	//
	// The "service" field of "destination" for TCP services must be an IP or a block of IPs in CIDR notation.
	Destination *IstioService `protobuf:"bytes,1,opt,name=destination" json:"destination,omitempty"`
	// REQUIRED: list of ports on which the external service is available.
	Ports []*EgressRule_Port `protobuf:"bytes,2,rep,name=ports" json:"ports,omitempty"`
	// Forward all the external traffic through a dedicated egress proxy. It is used in some scenarios
	// where there is a requirement that all the external traffic goes through special dedicated nodes/pods.
	// These dedicated egress nodes could then be more closely monitored for security vulnerabilities.
	//
	// The default is false, i.e. the sidecar forwards external traffic directly to the external service.
	UseEgressProxy bool `protobuf:"varint,3,opt,name=use_egress_proxy,json=useEgressProxy" json:"use_egress_proxy,omitempty"`
}

func (m *EgressRule) Reset()                    { *m = EgressRule{} }
func (m *EgressRule) String() string            { return proto.CompactTextString(m) }
func (*EgressRule) ProtoMessage()               {}
func (*EgressRule) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *EgressRule) GetDestination() *IstioService {
	if m != nil {
		return m.Destination
	}
	return nil
}

func (m *EgressRule) GetPorts() []*EgressRule_Port {
	if m != nil {
		return m.Ports
	}
	return nil
}

func (m *EgressRule) GetUseEgressProxy() bool {
	if m != nil {
		return m.UseEgressProxy
	}
	return false
}

// Port describes the properties of a specific TCP port of an external service.
type EgressRule_Port struct {
	// A valid non-negative integer port number.
	Port int32 `protobuf:"varint,1,opt,name=port" json:"port,omitempty"`
	// The protocol to communicate with the external services.
	// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|TCP|MONGO.
	Protocol string `protobuf:"bytes,2,opt,name=protocol" json:"protocol,omitempty"`
}

func (m *EgressRule_Port) Reset()                    { *m = EgressRule_Port{} }
func (m *EgressRule_Port) String() string            { return proto.CompactTextString(m) }
func (*EgressRule_Port) ProtoMessage()               {}
func (*EgressRule_Port) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0, 0} }

func (m *EgressRule_Port) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *EgressRule_Port) GetProtocol() string {
	if m != nil {
		return m.Protocol
	}
	return ""
}

func init() {
	proto.RegisterType((*EgressRule)(nil), "istio.routing.v1alpha1.EgressRule")
	proto.RegisterType((*EgressRule_Port)(nil), "istio.routing.v1alpha1.EgressRule.Port")
}

func init() { proto.RegisterFile("routing/v1alpha1/egress_rule.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xcd, 0x4b, 0x03, 0x31,
	0x14, 0xc4, 0x49, 0x3f, 0xa4, 0xbe, 0x05, 0x91, 0x1c, 0x64, 0x59, 0x10, 0xd7, 0x22, 0x98, 0x53,
	0x96, 0x56, 0xf0, 0xe6, 0x45, 0x50, 0xf0, 0x56, 0xe2, 0xcd, 0x4b, 0x59, 0xeb, 0xa3, 0x06, 0xc2,
	0xbe, 0x90, 0x8f, 0xa2, 0x7f, 0xb8, 0x77, 0x49, 0xb6, 0x6a, 0x51, 0x7b, 0x4b, 0x86, 0xf9, 0x4d,
	0x66, 0x02, 0x53, 0x47, 0x31, 0xe8, 0x6e, 0xdd, 0x6c, 0x66, 0xad, 0xb1, 0xaf, 0xed, 0xac, 0xc1,
	0xb5, 0x43, 0xef, 0x97, 0x2e, 0x1a, 0x94, 0xd6, 0x51, 0x20, 0x7e, 0xa2, 0x7d, 0xd0, 0x24, 0xb7,
	0x4e, 0xf9, 0xe5, 0xac, 0xce, 0xff, 0xb0, 0x49, 0xc0, 0x1d, 0x74, 0xfa, 0xc1, 0x00, 0xee, 0x72,
	0xa0, 0x8a, 0x06, 0xf9, 0x3d, 0x14, 0x2f, 0xe8, 0x83, 0xee, 0xda, 0xa0, 0xa9, 0x2b, 0x59, 0xcd,
	0x44, 0x31, 0xbf, 0x90, 0xff, 0xe7, 0xcb, 0x87, 0x24, 0x3f, 0xa2, 0xdb, 0xe8, 0x15, 0xaa, 0x5d,
	0x90, 0xdf, 0xc0, 0xd8, 0x92, 0x0b, 0xbe, 0x1c, 0xd4, 0x43, 0x51, 0xcc, 0x2f, 0xf7, 0x25, 0xfc,
	0x3c, 0x2d, 0x17, 0xe4, 0x82, 0xea, 0x29, 0x2e, 0xe0, 0x38, 0x7a, 0x5c, 0x6e, 0x97, 0x5a, 0x47,
	0x6f, 0xef, 0xe5, 0xb0, 0x66, 0x62, 0xa2, 0x8e, 0xa2, 0xc7, 0x1e, 0x5a, 0x24, 0xb5, 0xba, 0x86,
	0x51, 0x02, 0x39, 0x87, 0x51, 0x42, 0x73, 0xe3, 0xb1, 0xca, 0x67, 0x5e, 0xc1, 0x24, 0x8f, 0x5c,
	0x91, 0x29, 0x07, 0x35, 0x13, 0x87, 0xea, 0xfb, 0x7e, 0x7b, 0xf6, 0x74, 0xda, 0x57, 0xd2, 0xd4,
	0xb4, 0x56, 0x37, 0xbf, 0x7f, 0xea, 0xf9, 0x20, 0x5b, 0xaf, 0x3e, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x04, 0x3f, 0x9a, 0x90, 0x80, 0x01, 0x00, 0x00,
}
