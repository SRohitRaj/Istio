//go:build !agent
// +build !agent

// GENERATED FILE -- DO NOT EDIT
//

package collections

import (
	"reflect"

	k8sioapiadmissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	k8sioapiappsv1 "k8s.io/api/apps/v1"
	k8sioapicorev1 "k8s.io/api/core/v1"
	k8sioapidiscoveryv1 "k8s.io/api/discovery/v1"
	k8sioapinetworkingv1 "k8s.io/api/networking/v1"
	k8sioapiextensionsapiserverpkgapisapiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	sigsk8siogatewayapiapisv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	sigsk8siogatewayapiapisv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	istioioapiextensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	istioioapimeshv1alpha1 "istio.io/api/mesh/v1alpha1"
	istioioapimetav1alpha1 "istio.io/api/meta/v1alpha1"
	istioioapinetworkingv1alpha3 "istio.io/api/networking/v1alpha3"
	istioioapinetworkingv1beta1 "istio.io/api/networking/v1beta1"
	istioioapisecurityv1beta1 "istio.io/api/security/v1beta1"
	istioioapitelemetryv1alpha1 "istio.io/api/telemetry/v1alpha1"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/resource"
	"istio.io/istio/pkg/config/validation"
)

var (
	AuthorizationPolicy = resource.Builder{
		Identifier: "AuthorizationPolicy",
		Group:      "security.istio.io",
		Kind:       "AuthorizationPolicy",
		Plural:     "authorizationpolicies",
		Version:    "v1beta1",
		VersionAliases: []string{
			"v1",
		},
		Proto: "istio.security.v1beta1.AuthorizationPolicy", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapisecurityv1beta1.AuthorizationPolicy{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/security/v1beta1", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateAuthorizationPolicy,
	}.MustBuild()

	ConfigMap = resource.Builder{
		Identifier:    "ConfigMap",
		Group:         "",
		Kind:          "ConfigMap",
		Plural:        "configmaps",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.ConfigMap",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.ConfigMap{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	CustomResourceDefinition = resource.Builder{
		Identifier:    "CustomResourceDefinition",
		Group:         "apiextensions.k8s.io",
		Kind:          "CustomResourceDefinition",
		Plural:        "customresourcedefinitions",
		Version:       "v1",
		Proto:         "k8s.io.apiextensions_apiserver.pkg.apis.apiextensions.v1.CustomResourceDefinition",
		ReflectType:   reflect.TypeOf(&k8sioapiextensionsapiserverpkgapisapiextensionsv1.CustomResourceDefinition{}).Elem(),
		ProtoPackage:  "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1",
		ClusterScoped: true,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	Deployment = resource.Builder{
		Identifier:    "Deployment",
		Group:         "apps",
		Kind:          "Deployment",
		Plural:        "deployments",
		Version:       "v1",
		Proto:         "k8s.io.api.apps.v1.DeploymentSpec",
		ReflectType:   reflect.TypeOf(&k8sioapiappsv1.DeploymentSpec{}).Elem(),
		ProtoPackage:  "k8s.io/api/apps/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	DestinationRule = resource.Builder{
		Identifier: "DestinationRule",
		Group:      "networking.istio.io",
		Kind:       "DestinationRule",
		Plural:     "destinationrules",
		Version:    "v1alpha3",
		VersionAliases: []string{
			"v1beta1",
		},
		Proto: "istio.networking.v1alpha3.DestinationRule", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.DestinationRule{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateDestinationRule,
	}.MustBuild()

	EndpointSlice = resource.Builder{
		Identifier:    "EndpointSlice",
		Group:         "",
		Kind:          "EndpointSlice",
		Plural:        "endpointslices",
		Version:       "v1",
		Proto:         "k8s.io.api.discovery.v1.EndpointSlice",
		ReflectType:   reflect.TypeOf(&k8sioapidiscoveryv1.EndpointSlice{}).Elem(),
		ProtoPackage:  "k8s.io/api/discovery/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	Endpoints = resource.Builder{
		Identifier:    "Endpoints",
		Group:         "",
		Kind:          "Endpoints",
		Plural:        "endpoints",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.Endpoints",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.Endpoints{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	EnvoyFilter = resource.Builder{
		Identifier: "EnvoyFilter",
		Group:      "networking.istio.io",
		Kind:       "EnvoyFilter",
		Plural:     "envoyfilters",
		Version:    "v1alpha3",
		Proto:      "istio.networking.v1alpha3.EnvoyFilter", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.EnvoyFilter{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateEnvoyFilter,
	}.MustBuild()

	GRPCRoute = resource.Builder{
		Identifier: "GRPCRoute",
		Group:      "gateway.networking.k8s.io",
		Kind:       "GRPCRoute",
		Plural:     "grpcroutes",
		Version:    "v1alpha2",
		Proto:      "k8s.io.gateway_api.api.v1alpha1.GRPCRouteSpec", StatusProto: "k8s.io.gateway_api.api.v1alpha1.GRPCRouteStatus",
		ReflectType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.GRPCRouteSpec{}).Elem(), StatusType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.GRPCRouteStatus{}).Elem(),
		ProtoPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2", StatusPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateGRPCRoute,
	}.MustBuild()

	Gateway = resource.Builder{
		Identifier: "Gateway",
		Group:      "networking.istio.io",
		Kind:       "Gateway",
		Plural:     "gateways",
		Version:    "v1alpha3",
		VersionAliases: []string{
			"v1beta1",
		},
		Proto: "istio.networking.v1alpha3.Gateway", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.Gateway{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateGateway,
	}.MustBuild()

	GatewayClass = resource.Builder{
		Identifier: "GatewayClass",
		Group:      "gateway.networking.k8s.io",
		Kind:       "GatewayClass",
		Plural:     "gatewayclasses",
		Version:    "v1beta1",
		VersionAliases: []string{
			"v1alpha2",
		},
		Proto: "k8s.io.gateway_api.api.v1alpha1.GatewayClassSpec", StatusProto: "k8s.io.gateway_api.api.v1alpha1.GatewayClassStatus",
		ReflectType: reflect.TypeOf(&sigsk8siogatewayapiapisv1beta1.GatewayClassSpec{}).Elem(), StatusType: reflect.TypeOf(&sigsk8siogatewayapiapisv1beta1.GatewayClassStatus{}).Elem(),
		ProtoPackage: "sigs.k8s.io/gateway-api/apis/v1beta1", StatusPackage: "sigs.k8s.io/gateway-api/apis/v1beta1",
		ClusterScoped: true,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateGatewayClass,
	}.MustBuild()

	HTTPRoute = resource.Builder{
		Identifier: "HTTPRoute",
		Group:      "gateway.networking.k8s.io",
		Kind:       "HTTPRoute",
		Plural:     "httproutes",
		Version:    "v1beta1",
		VersionAliases: []string{
			"v1alpha2",
		},
		Proto: "k8s.io.gateway_api.api.v1alpha1.HTTPRouteSpec", StatusProto: "k8s.io.gateway_api.api.v1alpha1.HTTPRouteStatus",
		ReflectType: reflect.TypeOf(&sigsk8siogatewayapiapisv1beta1.HTTPRouteSpec{}).Elem(), StatusType: reflect.TypeOf(&sigsk8siogatewayapiapisv1beta1.HTTPRouteStatus{}).Elem(),
		ProtoPackage: "sigs.k8s.io/gateway-api/apis/v1beta1", StatusPackage: "sigs.k8s.io/gateway-api/apis/v1beta1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateHTTPRoute,
	}.MustBuild()

	Ingress = resource.Builder{
		Identifier: "Ingress",
		Group:      "networking.k8s.io",
		Kind:       "Ingress",
		Plural:     "ingresses",
		Version:    "v1",
		Proto:      "k8s.io.api.networking.v1.IngressSpec", StatusProto: "k8s.io.api.networking.v1.IngressStatus",
		ReflectType: reflect.TypeOf(&k8sioapinetworkingv1.IngressSpec{}).Elem(), StatusType: reflect.TypeOf(&k8sioapinetworkingv1.IngressStatus{}).Elem(),
		ProtoPackage: "k8s.io/api/networking/v1", StatusPackage: "k8s.io/api/networking/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	IngressClass = resource.Builder{
		Identifier:    "IngressClass",
		Group:         "networking.k8s.io",
		Kind:          "IngressClass",
		Plural:        "ingressclasses",
		Version:       "v1",
		Proto:         "k8s.io.api.networking.v1.IngressClassSpec",
		ReflectType:   reflect.TypeOf(&k8sioapinetworkingv1.IngressClassSpec{}).Elem(),
		ProtoPackage:  "k8s.io/api/networking/v1",
		ClusterScoped: true,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	KubernetesGateway = resource.Builder{
		Identifier: "KubernetesGateway",
		Group:      "gateway.networking.k8s.io",
		Kind:       "Gateway",
		Plural:     "gateways",
		Version:    "v1beta1",
		VersionAliases: []string{
			"v1alpha2",
		},
		Proto: "k8s.io.gateway_api.api.v1alpha1.GatewaySpec", StatusProto: "k8s.io.gateway_api.api.v1alpha1.GatewayStatus",
		ReflectType: reflect.TypeOf(&sigsk8siogatewayapiapisv1beta1.GatewaySpec{}).Elem(), StatusType: reflect.TypeOf(&sigsk8siogatewayapiapisv1beta1.GatewayStatus{}).Elem(),
		ProtoPackage: "sigs.k8s.io/gateway-api/apis/v1beta1", StatusPackage: "sigs.k8s.io/gateway-api/apis/v1beta1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateKubernetesGateway,
	}.MustBuild()

	MeshConfig = resource.Builder{
		Identifier:    "MeshConfig",
		Group:         "",
		Kind:          "MeshConfig",
		Plural:        "meshconfigs",
		Version:       "v1alpha1",
		Proto:         "istio.mesh.v1alpha1.MeshConfig",
		ReflectType:   reflect.TypeOf(&istioioapimeshv1alpha1.MeshConfig{}).Elem(),
		ProtoPackage:  "istio.io/api/mesh/v1alpha1",
		ClusterScoped: false,
		Synthetic:     true,
		Builtin:       false,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	MeshNetworks = resource.Builder{
		Identifier:    "MeshNetworks",
		Group:         "",
		Kind:          "MeshNetworks",
		Plural:        "meshnetworks",
		Version:       "v1alpha1",
		Proto:         "istio.mesh.v1alpha1.MeshNetworks",
		ReflectType:   reflect.TypeOf(&istioioapimeshv1alpha1.MeshNetworks{}).Elem(),
		ProtoPackage:  "istio.io/api/mesh/v1alpha1",
		ClusterScoped: false,
		Synthetic:     true,
		Builtin:       false,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	MutatingWebhookConfiguration = resource.Builder{
		Identifier:    "MutatingWebhookConfiguration",
		Group:         "admissionregistration.k8s.io",
		Kind:          "MutatingWebhookConfiguration",
		Plural:        "mutatingwebhookconfigurations",
		Version:       "v1",
		Proto:         "k8s.io.api.admissionregistration.v1.MutatingWebhookConfiguration",
		ReflectType:   reflect.TypeOf(&k8sioapiadmissionregistrationv1.MutatingWebhookConfiguration{}).Elem(),
		ProtoPackage:  "k8s.io/api/admissionregistration/v1",
		ClusterScoped: true,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	Namespace = resource.Builder{
		Identifier:    "Namespace",
		Group:         "",
		Kind:          "Namespace",
		Plural:        "namespaces",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.NamespaceSpec",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.NamespaceSpec{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: true,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	Node = resource.Builder{
		Identifier:    "Node",
		Group:         "",
		Kind:          "Node",
		Plural:        "nodes",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.NodeSpec",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.NodeSpec{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: true,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	PeerAuthentication = resource.Builder{
		Identifier: "PeerAuthentication",
		Group:      "security.istio.io",
		Kind:       "PeerAuthentication",
		Plural:     "peerauthentications",
		Version:    "v1beta1",
		Proto:      "istio.security.v1beta1.PeerAuthentication", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapisecurityv1beta1.PeerAuthentication{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/security/v1beta1", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidatePeerAuthentication,
	}.MustBuild()

	Pod = resource.Builder{
		Identifier:    "Pod",
		Group:         "",
		Kind:          "Pod",
		Plural:        "pods",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.PodSpec",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.PodSpec{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	ProxyConfig = resource.Builder{
		Identifier: "ProxyConfig",
		Group:      "networking.istio.io",
		Kind:       "ProxyConfig",
		Plural:     "proxyconfigs",
		Version:    "v1beta1",
		Proto:      "istio.networking.v1beta1.ProxyConfig", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1beta1.ProxyConfig{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1beta1", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateProxyConfig,
	}.MustBuild()

	ReferenceGrant = resource.Builder{
		Identifier:    "ReferenceGrant",
		Group:         "gateway.networking.k8s.io",
		Kind:          "ReferenceGrant",
		Plural:        "referencegrants",
		Version:       "v1alpha2",
		Proto:         "k8s.io.gateway_api.api.v1alpha1.ReferenceGrantSpec",
		ReflectType:   reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.ReferenceGrantSpec{}).Elem(),
		ProtoPackage:  "sigs.k8s.io/gateway-api/apis/v1alpha2",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	RequestAuthentication = resource.Builder{
		Identifier: "RequestAuthentication",
		Group:      "security.istio.io",
		Kind:       "RequestAuthentication",
		Plural:     "requestauthentications",
		Version:    "v1beta1",
		VersionAliases: []string{
			"v1",
		},
		Proto: "istio.security.v1beta1.RequestAuthentication", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapisecurityv1beta1.RequestAuthentication{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/security/v1beta1", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateRequestAuthentication,
	}.MustBuild()

	Secret = resource.Builder{
		Identifier:    "Secret",
		Group:         "",
		Kind:          "Secret",
		Plural:        "secrets",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.Secret",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.Secret{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	Service = resource.Builder{
		Identifier:    "Service",
		Group:         "",
		Kind:          "Service",
		Plural:        "services",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.ServiceSpec",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.ServiceSpec{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	ServiceAccount = resource.Builder{
		Identifier:    "ServiceAccount",
		Group:         "",
		Kind:          "ServiceAccount",
		Plural:        "serviceaccounts",
		Version:       "v1",
		Proto:         "k8s.io.api.core.v1.ServiceAccount",
		ReflectType:   reflect.TypeOf(&k8sioapicorev1.ServiceAccount{}).Elem(),
		ProtoPackage:  "k8s.io/api/core/v1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       true,
		ValidateProto: validation.EmptyValidate,
	}.MustBuild()

	ServiceEntry = resource.Builder{
		Identifier: "ServiceEntry",
		Group:      "networking.istio.io",
		Kind:       "ServiceEntry",
		Plural:     "serviceentries",
		Version:    "v1alpha3",
		VersionAliases: []string{
			"v1beta1",
		},
		Proto: "istio.networking.v1alpha3.ServiceEntry", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.ServiceEntry{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateServiceEntry,
	}.MustBuild()

	Sidecar = resource.Builder{
		Identifier: "Sidecar",
		Group:      "networking.istio.io",
		Kind:       "Sidecar",
		Plural:     "sidecars",
		Version:    "v1alpha3",
		VersionAliases: []string{
			"v1beta1",
		},
		Proto: "istio.networking.v1alpha3.Sidecar", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.Sidecar{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateSidecar,
	}.MustBuild()

	TCPRoute = resource.Builder{
		Identifier: "TCPRoute",
		Group:      "gateway.networking.k8s.io",
		Kind:       "TCPRoute",
		Plural:     "tcproutes",
		Version:    "v1alpha2",
		Proto:      "k8s.io.gateway_api.api.v1alpha1.TCPRouteSpec", StatusProto: "k8s.io.gateway_api.api.v1alpha1.TCPRouteStatus",
		ReflectType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.TCPRouteSpec{}).Elem(), StatusType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.TCPRouteStatus{}).Elem(),
		ProtoPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2", StatusPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateTCPRoute,
	}.MustBuild()

	TLSRoute = resource.Builder{
		Identifier: "TLSRoute",
		Group:      "gateway.networking.k8s.io",
		Kind:       "TLSRoute",
		Plural:     "tlsroutes",
		Version:    "v1alpha2",
		Proto:      "k8s.io.gateway_api.api.v1alpha1.TLSRouteSpec", StatusProto: "k8s.io.gateway_api.api.v1alpha1.TLSRouteStatus",
		ReflectType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.TLSRouteSpec{}).Elem(), StatusType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.TLSRouteStatus{}).Elem(),
		ProtoPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2", StatusPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateTLSRoute,
	}.MustBuild()

	Telemetry = resource.Builder{
		Identifier: "Telemetry",
		Group:      "telemetry.istio.io",
		Kind:       "Telemetry",
		Plural:     "telemetries",
		Version:    "v1alpha1",
		Proto:      "istio.telemetry.v1alpha1.Telemetry", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapitelemetryv1alpha1.Telemetry{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/telemetry/v1alpha1", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateTelemetry,
	}.MustBuild()

	UDPRoute = resource.Builder{
		Identifier: "UDPRoute",
		Group:      "gateway.networking.k8s.io",
		Kind:       "UDPRoute",
		Plural:     "udproutes",
		Version:    "v1alpha2",
		Proto:      "k8s.io.gateway_api.api.v1alpha1.UDPRouteSpec", StatusProto: "k8s.io.gateway_api.api.v1alpha1.UDPRouteStatus",
		ReflectType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.UDPRouteSpec{}).Elem(), StatusType: reflect.TypeOf(&sigsk8siogatewayapiapisv1alpha2.UDPRouteStatus{}).Elem(),
		ProtoPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2", StatusPackage: "sigs.k8s.io/gateway-api/apis/v1alpha2",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateUDPRoute,
	}.MustBuild()

	VirtualService = resource.Builder{
		Identifier: "VirtualService",
		Group:      "networking.istio.io",
		Kind:       "VirtualService",
		Plural:     "virtualservices",
		Version:    "v1alpha3",
		VersionAliases: []string{
			"v1beta1",
		},
		Proto: "istio.networking.v1alpha3.VirtualService", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.VirtualService{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateVirtualService,
	}.MustBuild()

	WasmPlugin = resource.Builder{
		Identifier: "WasmPlugin",
		Group:      "extensions.istio.io",
		Kind:       "WasmPlugin",
		Plural:     "wasmplugins",
		Version:    "v1alpha1",
		Proto:      "istio.extensions.v1alpha1.WasmPlugin", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapiextensionsv1alpha1.WasmPlugin{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/extensions/v1alpha1", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateWasmPlugin,
	}.MustBuild()

	WorkloadEntry = resource.Builder{
		Identifier: "WorkloadEntry",
		Group:      "networking.istio.io",
		Kind:       "WorkloadEntry",
		Plural:     "workloadentries",
		Version:    "v1alpha3",
		VersionAliases: []string{
			"v1beta1",
		},
		Proto: "istio.networking.v1alpha3.WorkloadEntry", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.WorkloadEntry{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateWorkloadEntry,
	}.MustBuild()

	WorkloadGroup = resource.Builder{
		Identifier: "WorkloadGroup",
		Group:      "networking.istio.io",
		Kind:       "WorkloadGroup",
		Plural:     "workloadgroups",
		Version:    "v1alpha3",
		VersionAliases: []string{
			"v1beta1",
		},
		Proto: "istio.networking.v1alpha3.WorkloadGroup", StatusProto: "istio.meta.v1alpha1.IstioStatus",
		ReflectType: reflect.TypeOf(&istioioapinetworkingv1alpha3.WorkloadGroup{}).Elem(), StatusType: reflect.TypeOf(&istioioapimetav1alpha1.IstioStatus{}).Elem(),
		ProtoPackage: "istio.io/api/networking/v1alpha3", StatusPackage: "istio.io/api/meta/v1alpha1",
		ClusterScoped: false,
		Synthetic:     false,
		Builtin:       false,
		ValidateProto: validation.ValidateWorkloadGroup,
	}.MustBuild()

	// All contains all collections in the system.
	All = collection.NewSchemasBuilder().
		MustAdd(AuthorizationPolicy).
		MustAdd(ConfigMap).
		MustAdd(CustomResourceDefinition).
		MustAdd(Deployment).
		MustAdd(DestinationRule).
		MustAdd(EndpointSlice).
		MustAdd(Endpoints).
		MustAdd(EnvoyFilter).
		MustAdd(GRPCRoute).
		MustAdd(Gateway).
		MustAdd(GatewayClass).
		MustAdd(HTTPRoute).
		MustAdd(Ingress).
		MustAdd(IngressClass).
		MustAdd(KubernetesGateway).
		MustAdd(MeshConfig).
		MustAdd(MeshNetworks).
		MustAdd(MutatingWebhookConfiguration).
		MustAdd(Namespace).
		MustAdd(Node).
		MustAdd(PeerAuthentication).
		MustAdd(Pod).
		MustAdd(ProxyConfig).
		MustAdd(ReferenceGrant).
		MustAdd(RequestAuthentication).
		MustAdd(Secret).
		MustAdd(Service).
		MustAdd(ServiceAccount).
		MustAdd(ServiceEntry).
		MustAdd(Sidecar).
		MustAdd(TCPRoute).
		MustAdd(TLSRoute).
		MustAdd(Telemetry).
		MustAdd(UDPRoute).
		MustAdd(VirtualService).
		MustAdd(WasmPlugin).
		MustAdd(WorkloadEntry).
		MustAdd(WorkloadGroup).
		Build()

	// Kube contains only kubernetes collections.
	Kube = collection.NewSchemasBuilder().
		MustAdd(ConfigMap).
		MustAdd(CustomResourceDefinition).
		MustAdd(Deployment).
		MustAdd(EndpointSlice).
		MustAdd(Endpoints).
		MustAdd(GRPCRoute).
		MustAdd(GatewayClass).
		MustAdd(HTTPRoute).
		MustAdd(Ingress).
		MustAdd(IngressClass).
		MustAdd(KubernetesGateway).
		MustAdd(MutatingWebhookConfiguration).
		MustAdd(Namespace).
		MustAdd(Node).
		MustAdd(Pod).
		MustAdd(ReferenceGrant).
		MustAdd(Secret).
		MustAdd(Service).
		MustAdd(ServiceAccount).
		MustAdd(TCPRoute).
		MustAdd(TLSRoute).
		MustAdd(UDPRoute).
		Build()

	// Pilot contains only collections used by Pilot.
	Pilot = collection.NewSchemasBuilder().
		MustAdd(AuthorizationPolicy).
		MustAdd(DestinationRule).
		MustAdd(EnvoyFilter).
		MustAdd(Gateway).
		MustAdd(PeerAuthentication).
		MustAdd(ProxyConfig).
		MustAdd(RequestAuthentication).
		MustAdd(ServiceEntry).
		MustAdd(Sidecar).
		MustAdd(Telemetry).
		MustAdd(VirtualService).
		MustAdd(WasmPlugin).
		MustAdd(WorkloadEntry).
		MustAdd(WorkloadGroup).
		Build()

	// PilotGatewayAPI contains only collections used by Pilot, including experimental Service Api.
	PilotGatewayAPI = collection.NewSchemasBuilder().
			MustAdd(AuthorizationPolicy).
			MustAdd(DestinationRule).
			MustAdd(EnvoyFilter).
			MustAdd(GRPCRoute).
			MustAdd(Gateway).
			MustAdd(GatewayClass).
			MustAdd(HTTPRoute).
			MustAdd(KubernetesGateway).
			MustAdd(PeerAuthentication).
			MustAdd(ProxyConfig).
			MustAdd(ReferenceGrant).
			MustAdd(RequestAuthentication).
			MustAdd(ServiceEntry).
			MustAdd(Sidecar).
			MustAdd(TCPRoute).
			MustAdd(TLSRoute).
			MustAdd(Telemetry).
			MustAdd(UDPRoute).
			MustAdd(VirtualService).
			MustAdd(WasmPlugin).
			MustAdd(WorkloadEntry).
			MustAdd(WorkloadGroup).
			Build()
)
