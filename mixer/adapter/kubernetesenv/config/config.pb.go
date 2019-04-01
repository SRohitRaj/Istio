// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: mixer/adapter/kubernetesenv/config/config.proto

// The `kubernetesenv` adapter extracts information from a Kubernetes environment
// and produces attributes that can be used in downstream adapters.
//
// This adapter supports the [kubernetes template](https://istio.io/docs/reference/config/policy-and-telemetry/templates/kubernetes/).

package config

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	io "io"
	math "math"
	reflect "reflect"
	strings "strings"
	time "time"
)

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

// Configuration parameters for the kubernetes adapter. These params
// control the manner in which the kubernetes adapter discovers and
// generates values related to pod information.
//
// The adapter works by looking up pod information by UIDs (of the
// form: "kubernetes://pod.namespace"). It expects that the UIDs will be
// supplied in an input map for three distinct traffic classes (source,
// destination, and origin).
//
// For all valid UIDs supplied, this adapter generates output
// values containing information about the related pods.
type Params struct {
	// File path to discover `kubeconfig`. For in-cluster configuration,
	// this should be left unset. For local configuration, this should
	// be set to the path of a `kubeconfig` file that can be used to
	// reach a kubernetes API server.
	//
	// NOTE: The `kubernetesenv` adapter will use the value of the `KUBECONFIG` environment variable
	// in the case where it is set (overriding any value configured
	// through this proto).
	//
	// Default: "" (unset)
	KubeconfigPath string `protobuf:"bytes,1,opt,name=kubeconfig_path,json=kubeconfigPath,proto3" json:"kubeconfig_path,omitempty"`
	// Controls the resync period of the Kubernetes cluster info cache.
	// The cache will watch for events and every so often completely resync.
	// This controls how frequently the complete resync occurs.
	//
	// Default: 5 minutes
	CacheRefreshDuration time.Duration `protobuf:"bytes,2,opt,name=cache_refresh_duration,json=cacheRefreshDuration,proto3,stdduration" json:"cache_refresh_duration"`
	// Namespace of the secret created for multicluster support.
	//
	// Details on multicluster and the Kubernetes secret required to
	// access the remote cluster's credentials can be found in
	// [multicluster install](https://istio.io/docs/setup/kubernetes/install/multicluster/).
	//
	// NOTE: If `cluster_registries_namespace` is not set then the environment
	// variable `POD_NAMESPACE` is checked/used. If `POD_NAMESPACE` is not
	// set then `cluster_registries_namespace` defaults to "istio-system".
	//
	// Default: "istio-system"
	ClusterRegistriesNamespace string `protobuf:"bytes,7,opt,name=cluster_registries_namespace,json=clusterRegistriesNamespace,proto3" json:"cluster_registries_namespace,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_b321b360e762a48d, []int{0}
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
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "adapter.kubernetesenv.config.Params")
}

func init() {
	proto.RegisterFile("mixer/adapter/kubernetesenv/config/config.proto", fileDescriptor_b321b360e762a48d)
}

var fileDescriptor_b321b360e762a48d = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xcd, 0x4a, 0xeb, 0x40,
	0x18, 0x86, 0x67, 0xce, 0x49, 0xd3, 0x9c, 0x1c, 0xd0, 0x10, 0x8a, 0xd4, 0x52, 0xbe, 0x16, 0x37,
	0x76, 0x95, 0x80, 0xde, 0x40, 0x29, 0xae, 0xb2, 0x90, 0x92, 0x9d, 0x6e, 0xc2, 0x34, 0xfd, 0x9a,
	0x04, 0xdb, 0x4c, 0x98, 0x4c, 0xc4, 0xa5, 0x97, 0xe0, 0xd2, 0x4b, 0xf0, 0x52, 0xba, 0xec, 0xb2,
	0x2b, 0x35, 0x29, 0x88, 0xcb, 0x5e, 0x82, 0x34, 0x3f, 0x8a, 0xab, 0xf9, 0xf9, 0x9e, 0x79, 0xde,
	0x79, 0x75, 0x7b, 0x15, 0x3d, 0xa0, 0xb0, 0xd9, 0x9c, 0x25, 0x12, 0x85, 0x7d, 0x97, 0xcd, 0x50,
	0xc4, 0x28, 0x31, 0xc5, 0xf8, 0xde, 0xf6, 0x79, 0xbc, 0x88, 0x82, 0x7a, 0xb1, 0x12, 0xc1, 0x25,
	0x37, 0xfb, 0x35, 0x6a, 0xfd, 0x42, 0xad, 0x8a, 0xe9, 0x75, 0x02, 0x1e, 0xf0, 0x12, 0xb4, 0x0f,
	0xbb, 0xea, 0x4d, 0x0f, 0x02, 0xce, 0x83, 0x25, 0xda, 0xe5, 0x69, 0x96, 0x2d, 0xec, 0x79, 0x26,
	0x98, 0x8c, 0x78, 0x5c, 0xcd, 0xcf, 0x3e, 0xa8, 0xae, 0x4e, 0x99, 0x60, 0xab, 0xd4, 0x3c, 0xd7,
	0x8f, 0x0f, 0xe2, 0x4a, 0xe7, 0x25, 0x4c, 0x86, 0x5d, 0x3a, 0xa4, 0xa3, 0x7f, 0xee, 0xd1, 0xcf,
	0xf5, 0x94, 0xc9, 0xd0, 0xbc, 0xd1, 0x4f, 0x7c, 0xe6, 0x87, 0xe8, 0x09, 0x5c, 0x08, 0x4c, 0x43,
	0xaf, 0x71, 0x76, 0xff, 0x0c, 0xe9, 0xe8, 0xff, 0xc5, 0xa9, 0x55, 0x85, 0x5a, 0x4d, 0xa8, 0x75,
	0x55, 0x03, 0x13, 0x6d, 0xfd, 0x3a, 0x20, 0xcf, 0x6f, 0x03, 0xea, 0x76, 0x4a, 0x85, 0x5b, 0x19,
	0x9a, 0xb9, 0x39, 0xd6, 0xfb, 0xfe, 0x32, 0x4b, 0x25, 0x0a, 0x4f, 0x60, 0x10, 0xa5, 0x52, 0x44,
	0x98, 0x7a, 0x31, 0x5b, 0x61, 0x9a, 0x30, 0x1f, 0xbb, 0xed, 0xf2, 0x43, 0xbd, 0x9a, 0x71, 0xbf,
	0x91, 0xeb, 0x86, 0x70, 0x14, 0xed, 0xaf, 0xa1, 0x38, 0x8a, 0xa6, 0x18, 0x2d, 0x47, 0xd1, 0x5a,
	0x86, 0xea, 0x28, 0x9a, 0x6a, 0xb4, 0x27, 0xe3, 0x75, 0x0e, 0x64, 0x93, 0x03, 0xd9, 0xe6, 0x40,
	0xf6, 0x39, 0x90, 0xc7, 0x02, 0xe8, 0x4b, 0x01, 0x64, 0x5d, 0x00, 0xdd, 0x14, 0x40, 0xdf, 0x0b,
	0xa0, 0x9f, 0x05, 0x90, 0x7d, 0x01, 0xf4, 0x69, 0x07, 0x64, 0xb3, 0x03, 0xb2, 0xdd, 0x01, 0xb9,
	0x55, 0xab, 0xea, 0x33, 0xb5, 0xac, 0x73, 0xf9, 0x15, 0x00, 0x00, 0xff, 0xff, 0xb6, 0xfb, 0x72,
	0x08, 0xb8, 0x01, 0x00, 0x00,
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
	if len(m.KubeconfigPath) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintConfig(dAtA, i, uint64(len(m.KubeconfigPath)))
		i += copy(dAtA[i:], m.KubeconfigPath)
	}
	dAtA[i] = 0x12
	i++
	i = encodeVarintConfig(dAtA, i, uint64(github_com_gogo_protobuf_types.SizeOfStdDuration(m.CacheRefreshDuration)))
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.CacheRefreshDuration, dAtA[i:])
	if err1 != nil {
		return 0, err1
	}
	i += n1
	if len(m.ClusterRegistriesNamespace) > 0 {
		dAtA[i] = 0x3a
		i++
		i = encodeVarintConfig(dAtA, i, uint64(len(m.ClusterRegistriesNamespace)))
		i += copy(dAtA[i:], m.ClusterRegistriesNamespace)
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
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.KubeconfigPath)
	if l > 0 {
		n += 1 + l + sovConfig(uint64(l))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.CacheRefreshDuration)
	n += 1 + l + sovConfig(uint64(l))
	l = len(m.ClusterRegistriesNamespace)
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
	s := strings.Join([]string{`&Params{`,
		`KubeconfigPath:` + fmt.Sprintf("%v", this.KubeconfigPath) + `,`,
		`CacheRefreshDuration:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.CacheRefreshDuration), "Duration", "types.Duration", 1), `&`, ``, 1) + `,`,
		`ClusterRegistriesNamespace:` + fmt.Sprintf("%v", this.ClusterRegistriesNamespace) + `,`,
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
			wire |= uint64(b&0x7F) << shift
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
				return fmt.Errorf("proto: wrong wireType = %d for field KubeconfigPath", wireType)
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
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KubeconfigPath = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CacheRefreshDuration", wireType)
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
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.CacheRefreshDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClusterRegistriesNamespace", wireType)
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
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClusterRegistriesNamespace = string(dAtA[iNdEx:postIndex])
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
			if (iNdEx + skippy) < 0 {
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
			if length < 0 {
				return 0, ErrInvalidLengthConfig
			}
			iNdEx += length
			if iNdEx < 0 {
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
				if iNdEx < 0 {
					return 0, ErrInvalidLengthConfig
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
	ErrInvalidLengthConfig = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowConfig   = fmt.Errorf("proto: integer overflow")
)
