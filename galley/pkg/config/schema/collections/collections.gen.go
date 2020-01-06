// GENERATED FILE -- DO NOT EDIT
//

package collections

import (
	"istio.io/istio/galley/pkg/config/schema/collection"
	"istio.io/istio/galley/pkg/config/schema/resource"
	"istio.io/istio/pkg/config/validation"
)

var (

	// IstioAuthenticationV1Alpha1Meshpolicies describes the collection
	// istio/authentication/v1alpha1/meshpolicies
	IstioAuthenticationV1Alpha1Meshpolicies = collection.Builder{
		Name:     "istio/authentication/v1alpha1/meshpolicies",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "authentication.istio.io",
			Kind:          "MeshPolicy",
			Plural:        "meshpolicies",
			Version:       "v1alpha1",
			Proto:         "istio.authentication.v1alpha1.Policy",
			ProtoPackage:  "istio.io/api/authentication/v1alpha1",
			ClusterScoped: true,
			ValidateProto: validation.ValidateAuthenticationPolicy,
		}.Build(),
	}.MustBuild()

	// IstioAuthenticationV1Alpha1Policies describes the collection
	// istio/authentication/v1alpha1/policies
	IstioAuthenticationV1Alpha1Policies = collection.Builder{
		Name:     "istio/authentication/v1alpha1/policies",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "authentication.istio.io",
			Kind:          "Policy",
			Plural:        "policies",
			Version:       "v1alpha1",
			Proto:         "istio.authentication.v1alpha1.Policy",
			ProtoPackage:  "istio.io/api/authentication/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateAuthenticationPolicy,
		}.Build(),
	}.MustBuild()

	// IstioConfigV1Alpha2Adapters describes the collection
	// istio/config/v1alpha2/adapters
	IstioConfigV1Alpha2Adapters = collection.Builder{
		Name:     "istio/config/v1alpha2/adapters",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "adapter",
			Plural:        "adapters",
			Version:       "v1alpha2",
			Proto:         "google.protobuf.Struct",
			ProtoPackage:  "github.com/gogo/protobuf/types",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// IstioConfigV1Alpha2Httpapispecbindings describes the collection
	// istio/config/v1alpha2/httpapispecbindings
	IstioConfigV1Alpha2Httpapispecbindings = collection.Builder{
		Name:     "istio/config/v1alpha2/httpapispecbindings",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "HTTPAPISpecBinding",
			Plural:        "httpapispecbindings",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.HTTPAPISpecBinding",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateHTTPAPISpecBinding,
		}.Build(),
	}.MustBuild()

	// IstioConfigV1Alpha2Httpapispecs describes the collection
	// istio/config/v1alpha2/httpapispecs
	IstioConfigV1Alpha2Httpapispecs = collection.Builder{
		Name:     "istio/config/v1alpha2/httpapispecs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "HTTPAPISpec",
			Plural:        "httpapispecs",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.HTTPAPISpec",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateHTTPAPISpec,
		}.Build(),
	}.MustBuild()

	// IstioConfigV1Alpha2Templates describes the collection
	// istio/config/v1alpha2/templates
	IstioConfigV1Alpha2Templates = collection.Builder{
		Name:     "istio/config/v1alpha2/templates",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "template",
			Plural:        "templates",
			Version:       "v1alpha2",
			Proto:         "google.protobuf.Struct",
			ProtoPackage:  "github.com/gogo/protobuf/types",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// IstioMeshV1Alpha1MeshConfig describes the collection
	// istio/mesh/v1alpha1/MeshConfig
	IstioMeshV1Alpha1MeshConfig = collection.Builder{
		Name:     "istio/mesh/v1alpha1/MeshConfig",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "",
			Kind:          "MeshConfig",
			Plural:        "meshconfigs",
			Version:       "v1alpha1",
			Proto:         "istio.mesh.v1alpha1.MeshConfig",
			ProtoPackage:  "istio.io/api/mesh/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// IstioMixerV1ConfigClientQuotaspecbindings describes the collection
	// istio/mixer/v1/config/client/quotaspecbindings
	IstioMixerV1ConfigClientQuotaspecbindings = collection.Builder{
		Name:     "istio/mixer/v1/config/client/quotaspecbindings",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "QuotaSpecBinding",
			Plural:        "quotaspecbindings",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.QuotaSpecBinding",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateQuotaSpecBinding,
		}.Build(),
	}.MustBuild()

	// IstioMixerV1ConfigClientQuotaspecs describes the collection
	// istio/mixer/v1/config/client/quotaspecs
	IstioMixerV1ConfigClientQuotaspecs = collection.Builder{
		Name:     "istio/mixer/v1/config/client/quotaspecs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "QuotaSpec",
			Plural:        "quotaspecs",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.QuotaSpec",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateQuotaSpec,
		}.Build(),
	}.MustBuild()

	// IstioNetworkingV1Alpha3Destinationrules describes the collection
	// istio/networking/v1alpha3/destinationrules
	IstioNetworkingV1Alpha3Destinationrules = collection.Builder{
		Name:     "istio/networking/v1alpha3/destinationrules",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "DestinationRule",
			Plural:        "destinationrules",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.DestinationRule",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateDestinationRule,
		}.Build(),
	}.MustBuild()

	// IstioNetworkingV1Alpha3Envoyfilters describes the collection
	// istio/networking/v1alpha3/envoyfilters
	IstioNetworkingV1Alpha3Envoyfilters = collection.Builder{
		Name:     "istio/networking/v1alpha3/envoyfilters",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "EnvoyFilter",
			Plural:        "envoyfilters",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.EnvoyFilter",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateEnvoyFilter,
		}.Build(),
	}.MustBuild()

	// IstioNetworkingV1Alpha3Gateways describes the collection
	// istio/networking/v1alpha3/gateways
	IstioNetworkingV1Alpha3Gateways = collection.Builder{
		Name:     "istio/networking/v1alpha3/gateways",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "Gateway",
			Plural:        "gateways",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.Gateway",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateGateway,
		}.Build(),
	}.MustBuild()

	// IstioNetworkingV1Alpha3Serviceentries describes the collection
	// istio/networking/v1alpha3/serviceentries
	IstioNetworkingV1Alpha3Serviceentries = collection.Builder{
		Name:     "istio/networking/v1alpha3/serviceentries",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "ServiceEntry",
			Plural:        "serviceentries",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.ServiceEntry",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateServiceEntry,
		}.Build(),
	}.MustBuild()

	// IstioNetworkingV1Alpha3Sidecars describes the collection
	// istio/networking/v1alpha3/sidecars
	IstioNetworkingV1Alpha3Sidecars = collection.Builder{
		Name:     "istio/networking/v1alpha3/sidecars",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "Sidecar",
			Plural:        "sidecars",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.Sidecar",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateSidecar,
		}.Build(),
	}.MustBuild()

	// IstioNetworkingV1Alpha3SyntheticServiceentries describes the collection
	// istio/networking/v1alpha3/synthetic/serviceentries
	IstioNetworkingV1Alpha3SyntheticServiceentries = collection.Builder{
		Name:     "istio/networking/v1alpha3/synthetic/serviceentries",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "ServiceEntry",
			Plural:        "serviceentries",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.ServiceEntry",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateServiceEntry,
		}.Build(),
	}.MustBuild()

	// IstioNetworkingV1Alpha3Virtualservices describes the collection
	// istio/networking/v1alpha3/virtualservices
	IstioNetworkingV1Alpha3Virtualservices = collection.Builder{
		Name:     "istio/networking/v1alpha3/virtualservices",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "VirtualService",
			Plural:        "virtualservices",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.VirtualService",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateVirtualService,
		}.Build(),
	}.MustBuild()

	// IstioPolicyV1Beta1Attributemanifests describes the collection
	// istio/policy/v1beta1/attributemanifests
	IstioPolicyV1Beta1Attributemanifests = collection.Builder{
		Name:     "istio/policy/v1beta1/attributemanifests",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "attributemanifest",
			Plural:        "attributemanifests",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.AttributeManifest",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// IstioPolicyV1Beta1Handlers describes the collection
	// istio/policy/v1beta1/handlers
	IstioPolicyV1Beta1Handlers = collection.Builder{
		Name:     "istio/policy/v1beta1/handlers",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "handler",
			Plural:        "handlers",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.Handler",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// IstioPolicyV1Beta1Instances describes the collection
	// istio/policy/v1beta1/instances
	IstioPolicyV1Beta1Instances = collection.Builder{
		Name:     "istio/policy/v1beta1/instances",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "instance",
			Plural:        "instances",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.Instance",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// IstioPolicyV1Beta1Rules describes the collection
	// istio/policy/v1beta1/rules
	IstioPolicyV1Beta1Rules = collection.Builder{
		Name:     "istio/policy/v1beta1/rules",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "rule",
			Plural:        "rules",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.Rule",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// IstioRbacV1Alpha1Clusterrbacconfigs describes the collection
	// istio/rbac/v1alpha1/clusterrbacconfigs
	IstioRbacV1Alpha1Clusterrbacconfigs = collection.Builder{
		Name:     "istio/rbac/v1alpha1/clusterrbacconfigs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "ClusterRbacConfig",
			Plural:        "clusterrbacconfigs",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.RbacConfig",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: true,
			ValidateProto: validation.ValidateClusterRbacConfig,
		}.Build(),
	}.MustBuild()

	// IstioRbacV1Alpha1Rbacconfigs describes the collection
	// istio/rbac/v1alpha1/rbacconfigs
	IstioRbacV1Alpha1Rbacconfigs = collection.Builder{
		Name:     "istio/rbac/v1alpha1/rbacconfigs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "RbacConfig",
			Plural:        "rbacconfigs",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.RbacConfig",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateRbacConfig,
		}.Build(),
	}.MustBuild()

	// IstioRbacV1Alpha1Servicerolebindings describes the collection
	// istio/rbac/v1alpha1/servicerolebindings
	IstioRbacV1Alpha1Servicerolebindings = collection.Builder{
		Name:     "istio/rbac/v1alpha1/servicerolebindings",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "ServiceRoleBinding",
			Plural:        "servicerolebindings",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.ServiceRoleBinding",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateServiceRoleBinding,
		}.Build(),
	}.MustBuild()

	// IstioRbacV1Alpha1Serviceroles describes the collection
	// istio/rbac/v1alpha1/serviceroles
	IstioRbacV1Alpha1Serviceroles = collection.Builder{
		Name:     "istio/rbac/v1alpha1/serviceroles",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "ServiceRole",
			Plural:        "serviceroles",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.ServiceRole",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateServiceRole,
		}.Build(),
	}.MustBuild()

	// IstioSecurityV1Beta1Authorizationpolicies describes the collection
	// istio/security/v1beta1/authorizationpolicies
	IstioSecurityV1Beta1Authorizationpolicies = collection.Builder{
		Name:     "istio/security/v1beta1/authorizationpolicies",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "security.istio.io",
			Kind:          "AuthorizationPolicy",
			Plural:        "authorizationpolicies",
			Version:       "v1beta1",
			Proto:         "istio.security.v1beta1.AuthorizationPolicy",
			ProtoPackage:  "istio.io/api/security/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateAuthorizationPolicy,
		}.Build(),
	}.MustBuild()

	// IstioSecurityV1Beta1Requestauthentications describes the collection
	// istio/security/v1beta1/requestauthentications
	IstioSecurityV1Beta1Requestauthentications = collection.Builder{
		Name:     "istio/security/v1beta1/requestauthentications",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "security.istio.io",
			Kind:          "RequestAuthentication",
			Plural:        "requestauthentications",
			Version:       "v1beta1",
			Proto:         "istio.security.v1beta1.RequestAuthentication",
			ProtoPackage:  "istio.io/api/security/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateRequestAuthentication,
		}.Build(),
	}.MustBuild()

	// K8SAppsV1Deployments describes the collection k8s/apps/v1/deployments
	K8SAppsV1Deployments = collection.Builder{
		Name:     "k8s/apps/v1/deployments",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "apps",
			Kind:          "Deployment",
			Plural:        "Deployments",
			Version:       "v1",
			Proto:         "k8s.io.api.apps.v1.Deployment",
			ProtoPackage:  "k8s.io/api/apps/v1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SAuthenticationIstioIoV1Alpha1Meshpolicies describes the collection
	// k8s/authentication.istio.io/v1alpha1/meshpolicies
	K8SAuthenticationIstioIoV1Alpha1Meshpolicies = collection.Builder{
		Name:     "k8s/authentication.istio.io/v1alpha1/meshpolicies",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "authentication.istio.io",
			Kind:          "MeshPolicy",
			Plural:        "meshpolicies",
			Version:       "v1alpha1",
			Proto:         "istio.authentication.v1alpha1.Policy",
			ProtoPackage:  "istio.io/api/authentication/v1alpha1",
			ClusterScoped: true,
			ValidateProto: validation.ValidateAuthenticationPolicy,
		}.Build(),
	}.MustBuild()

	// K8SAuthenticationIstioIoV1Alpha1Policies describes the collection
	// k8s/authentication.istio.io/v1alpha1/policies
	K8SAuthenticationIstioIoV1Alpha1Policies = collection.Builder{
		Name:     "k8s/authentication.istio.io/v1alpha1/policies",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "authentication.istio.io",
			Kind:          "Policy",
			Plural:        "policies",
			Version:       "v1alpha1",
			Proto:         "istio.authentication.v1alpha1.Policy",
			ProtoPackage:  "istio.io/api/authentication/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateAuthenticationPolicy,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Adapters describes the collection
	// k8s/config.istio.io/v1alpha2/adapters
	K8SConfigIstioIoV1Alpha2Adapters = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/adapters",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "adapter",
			Plural:        "adapters",
			Version:       "v1alpha2",
			Proto:         "google.protobuf.Struct",
			ProtoPackage:  "github.com/gogo/protobuf/types",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Attributemanifests describes the collection
	// k8s/config.istio.io/v1alpha2/attributemanifests
	K8SConfigIstioIoV1Alpha2Attributemanifests = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/attributemanifests",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "attributemanifest",
			Plural:        "attributemanifests",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.AttributeManifest",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Handlers describes the collection
	// k8s/config.istio.io/v1alpha2/handlers
	K8SConfigIstioIoV1Alpha2Handlers = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/handlers",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "handler",
			Plural:        "handlers",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.Handler",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Httpapispecbindings describes the collection
	// k8s/config.istio.io/v1alpha2/httpapispecbindings
	K8SConfigIstioIoV1Alpha2Httpapispecbindings = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/httpapispecbindings",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "HTTPAPISpecBinding",
			Plural:        "httpapispecbindings",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.HTTPAPISpecBinding",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateHTTPAPISpecBinding,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Httpapispecs describes the collection
	// k8s/config.istio.io/v1alpha2/httpapispecs
	K8SConfigIstioIoV1Alpha2Httpapispecs = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/httpapispecs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "HTTPAPISpec",
			Plural:        "httpapispecs",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.HTTPAPISpec",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateHTTPAPISpec,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Instances describes the collection
	// k8s/config.istio.io/v1alpha2/instances
	K8SConfigIstioIoV1Alpha2Instances = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/instances",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "instance",
			Plural:        "instances",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.Instance",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Quotaspecbindings describes the collection
	// k8s/config.istio.io/v1alpha2/quotaspecbindings
	K8SConfigIstioIoV1Alpha2Quotaspecbindings = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/quotaspecbindings",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "QuotaSpecBinding",
			Plural:        "quotaspecbindings",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.QuotaSpecBinding",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateQuotaSpecBinding,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Quotaspecs describes the collection
	// k8s/config.istio.io/v1alpha2/quotaspecs
	K8SConfigIstioIoV1Alpha2Quotaspecs = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/quotaspecs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "QuotaSpec",
			Plural:        "quotaspecs",
			Version:       "v1alpha2",
			Proto:         "istio.mixer.v1.config.client.QuotaSpec",
			ProtoPackage:  "istio.io/api/mixer/v1/config/client",
			ClusterScoped: false,
			ValidateProto: validation.ValidateQuotaSpec,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Rules describes the collection
	// k8s/config.istio.io/v1alpha2/rules
	K8SConfigIstioIoV1Alpha2Rules = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/rules",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "rule",
			Plural:        "rules",
			Version:       "v1alpha2",
			Proto:         "istio.policy.v1beta1.Rule",
			ProtoPackage:  "istio.io/api/policy/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SConfigIstioIoV1Alpha2Templates describes the collection
	// k8s/config.istio.io/v1alpha2/templates
	K8SConfigIstioIoV1Alpha2Templates = collection.Builder{
		Name:     "k8s/config.istio.io/v1alpha2/templates",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "config.istio.io",
			Kind:          "template",
			Plural:        "templates",
			Version:       "v1alpha2",
			Proto:         "google.protobuf.Struct",
			ProtoPackage:  "github.com/gogo/protobuf/types",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SCoreV1Endpoints describes the collection k8s/core/v1/endpoints
	K8SCoreV1Endpoints = collection.Builder{
		Name:     "k8s/core/v1/endpoints",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "",
			Kind:          "Endpoints",
			Plural:        "endpoints",
			Version:       "v1",
			Proto:         "k8s.io.api.core.v1.Endpoints",
			ProtoPackage:  "k8s.io/api/core/v1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SCoreV1Namespaces describes the collection k8s/core/v1/namespaces
	K8SCoreV1Namespaces = collection.Builder{
		Name:     "k8s/core/v1/namespaces",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "",
			Kind:          "Namespace",
			Plural:        "namespaces",
			Version:       "v1",
			Proto:         "k8s.io.api.core.v1.NamespaceSpec",
			ProtoPackage:  "k8s.io/api/core/v1",
			ClusterScoped: true,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SCoreV1Nodes describes the collection k8s/core/v1/nodes
	K8SCoreV1Nodes = collection.Builder{
		Name:     "k8s/core/v1/nodes",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "",
			Kind:          "Node",
			Plural:        "nodes",
			Version:       "v1",
			Proto:         "k8s.io.api.core.v1.NodeSpec",
			ProtoPackage:  "k8s.io/api/core/v1",
			ClusterScoped: true,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SCoreV1Pods describes the collection k8s/core/v1/pods
	K8SCoreV1Pods = collection.Builder{
		Name:     "k8s/core/v1/pods",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "",
			Kind:          "Pod",
			Plural:        "pods",
			Version:       "v1",
			Proto:         "k8s.io.api.core.v1.Pod",
			ProtoPackage:  "k8s.io/api/core/v1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SCoreV1Secrets describes the collection k8s/core/v1/secrets
	K8SCoreV1Secrets = collection.Builder{
		Name:     "k8s/core/v1/secrets",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "",
			Kind:          "Secret",
			Plural:        "secrets",
			Version:       "v1",
			Proto:         "k8s.io.api.core.v1.Secret",
			ProtoPackage:  "k8s.io/api/core/v1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SCoreV1Services describes the collection k8s/core/v1/services
	K8SCoreV1Services = collection.Builder{
		Name:     "k8s/core/v1/services",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "",
			Kind:          "Service",
			Plural:        "services",
			Version:       "v1",
			Proto:         "k8s.io.api.core.v1.ServiceSpec",
			ProtoPackage:  "k8s.io/api/core/v1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SExtensionsV1Beta1Ingresses describes the collection
	// k8s/extensions/v1beta1/ingresses
	K8SExtensionsV1Beta1Ingresses = collection.Builder{
		Name:     "k8s/extensions/v1beta1/ingresses",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "extensions",
			Kind:          "Ingress",
			Plural:        "ingresses",
			Version:       "v1beta1",
			Proto:         "k8s.io.api.extensions.v1beta1.IngressSpec",
			ProtoPackage:  "k8s.io/api/extensions/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.EmptyValidate,
		}.Build(),
	}.MustBuild()

	// K8SNetworkingIstioIoV1Alpha3Destinationrules describes the collection
	// k8s/networking.istio.io/v1alpha3/destinationrules
	K8SNetworkingIstioIoV1Alpha3Destinationrules = collection.Builder{
		Name:     "k8s/networking.istio.io/v1alpha3/destinationrules",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "DestinationRule",
			Plural:        "destinationrules",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.DestinationRule",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateDestinationRule,
		}.Build(),
	}.MustBuild()

	// K8SNetworkingIstioIoV1Alpha3Envoyfilters describes the collection
	// k8s/networking.istio.io/v1alpha3/envoyfilters
	K8SNetworkingIstioIoV1Alpha3Envoyfilters = collection.Builder{
		Name:     "k8s/networking.istio.io/v1alpha3/envoyfilters",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "EnvoyFilter",
			Plural:        "envoyfilters",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.EnvoyFilter",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateEnvoyFilter,
		}.Build(),
	}.MustBuild()

	// K8SNetworkingIstioIoV1Alpha3Gateways describes the collection
	// k8s/networking.istio.io/v1alpha3/gateways
	K8SNetworkingIstioIoV1Alpha3Gateways = collection.Builder{
		Name:     "k8s/networking.istio.io/v1alpha3/gateways",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "Gateway",
			Plural:        "gateways",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.Gateway",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateGateway,
		}.Build(),
	}.MustBuild()

	// K8SNetworkingIstioIoV1Alpha3Serviceentries describes the collection
	// k8s/networking.istio.io/v1alpha3/serviceentries
	K8SNetworkingIstioIoV1Alpha3Serviceentries = collection.Builder{
		Name:     "k8s/networking.istio.io/v1alpha3/serviceentries",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "ServiceEntry",
			Plural:        "serviceentries",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.ServiceEntry",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateServiceEntry,
		}.Build(),
	}.MustBuild()

	// K8SNetworkingIstioIoV1Alpha3Sidecars describes the collection
	// k8s/networking.istio.io/v1alpha3/sidecars
	K8SNetworkingIstioIoV1Alpha3Sidecars = collection.Builder{
		Name:     "k8s/networking.istio.io/v1alpha3/sidecars",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "Sidecar",
			Plural:        "sidecars",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.Sidecar",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateSidecar,
		}.Build(),
	}.MustBuild()

	// K8SNetworkingIstioIoV1Alpha3Virtualservices describes the collection
	// k8s/networking.istio.io/v1alpha3/virtualservices
	K8SNetworkingIstioIoV1Alpha3Virtualservices = collection.Builder{
		Name:     "k8s/networking.istio.io/v1alpha3/virtualservices",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "networking.istio.io",
			Kind:          "VirtualService",
			Plural:        "virtualservices",
			Version:       "v1alpha3",
			Proto:         "istio.networking.v1alpha3.VirtualService",
			ProtoPackage:  "istio.io/api/networking/v1alpha3",
			ClusterScoped: false,
			ValidateProto: validation.ValidateVirtualService,
		}.Build(),
	}.MustBuild()

	// K8SRbacIstioIoV1Alpha1Clusterrbacconfigs describes the collection
	// k8s/rbac.istio.io/v1alpha1/clusterrbacconfigs
	K8SRbacIstioIoV1Alpha1Clusterrbacconfigs = collection.Builder{
		Name:     "k8s/rbac.istio.io/v1alpha1/clusterrbacconfigs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "ClusterRbacConfig",
			Plural:        "clusterrbacconfigs",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.RbacConfig",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: true,
			ValidateProto: validation.ValidateClusterRbacConfig,
		}.Build(),
	}.MustBuild()

	// K8SRbacIstioIoV1Alpha1Policy describes the collection
	// k8s/rbac.istio.io/v1alpha1/policy
	K8SRbacIstioIoV1Alpha1Policy = collection.Builder{
		Name:     "k8s/rbac.istio.io/v1alpha1/policy",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "ServiceRoleBinding",
			Plural:        "servicerolebindings",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.ServiceRoleBinding",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateServiceRoleBinding,
		}.Build(),
	}.MustBuild()

	// K8SRbacIstioIoV1Alpha1Rbacconfigs describes the collection
	// k8s/rbac.istio.io/v1alpha1/rbacconfigs
	K8SRbacIstioIoV1Alpha1Rbacconfigs = collection.Builder{
		Name:     "k8s/rbac.istio.io/v1alpha1/rbacconfigs",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "RbacConfig",
			Plural:        "rbacconfigs",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.RbacConfig",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateRbacConfig,
		}.Build(),
	}.MustBuild()

	// K8SRbacIstioIoV1Alpha1Serviceroles describes the collection
	// k8s/rbac.istio.io/v1alpha1/serviceroles
	K8SRbacIstioIoV1Alpha1Serviceroles = collection.Builder{
		Name:     "k8s/rbac.istio.io/v1alpha1/serviceroles",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "rbac.istio.io",
			Kind:          "ServiceRole",
			Plural:        "serviceroles",
			Version:       "v1alpha1",
			Proto:         "istio.rbac.v1alpha1.ServiceRole",
			ProtoPackage:  "istio.io/api/rbac/v1alpha1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateServiceRole,
		}.Build(),
	}.MustBuild()

	// K8SSecurityIstioIoV1Beta1Authorizationpolicies describes the collection
	// k8s/security.istio.io/v1beta1/authorizationpolicies
	K8SSecurityIstioIoV1Beta1Authorizationpolicies = collection.Builder{
		Name:     "k8s/security.istio.io/v1beta1/authorizationpolicies",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "security.istio.io",
			Kind:          "AuthorizationPolicy",
			Plural:        "authorizationpolicies",
			Version:       "v1beta1",
			Proto:         "istio.security.v1beta1.AuthorizationPolicy",
			ProtoPackage:  "istio.io/api/security/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateAuthorizationPolicy,
		}.Build(),
	}.MustBuild()

	// K8SSecurityIstioIoV1Beta1Requestauthentications describes the
	// collection k8s/security.istio.io/v1beta1/requestauthentications
	K8SSecurityIstioIoV1Beta1Requestauthentications = collection.Builder{
		Name:     "k8s/security.istio.io/v1beta1/requestauthentications",
		Disabled: false,
		Schema: resource.Builder{
			Group:         "security.istio.io",
			Kind:          "RequestAuthentication",
			Plural:        "requestauthentications",
			Version:       "v1beta1",
			Proto:         "istio.security.v1beta1.RequestAuthentication",
			ProtoPackage:  "istio.io/api/security/v1beta1",
			ClusterScoped: false,
			ValidateProto: validation.ValidateRequestAuthentication,
		}.Build(),
	}.MustBuild()

	// All contains all collections in the system.
	All = collection.NewSchemasBuilder().
		MustAdd(IstioAuthenticationV1Alpha1Meshpolicies).
		MustAdd(IstioAuthenticationV1Alpha1Policies).
		MustAdd(IstioConfigV1Alpha2Adapters).
		MustAdd(IstioConfigV1Alpha2Httpapispecbindings).
		MustAdd(IstioConfigV1Alpha2Httpapispecs).
		MustAdd(IstioConfigV1Alpha2Templates).
		MustAdd(IstioMeshV1Alpha1MeshConfig).
		MustAdd(IstioMixerV1ConfigClientQuotaspecbindings).
		MustAdd(IstioMixerV1ConfigClientQuotaspecs).
		MustAdd(IstioNetworkingV1Alpha3Destinationrules).
		MustAdd(IstioNetworkingV1Alpha3Envoyfilters).
		MustAdd(IstioNetworkingV1Alpha3Gateways).
		MustAdd(IstioNetworkingV1Alpha3Serviceentries).
		MustAdd(IstioNetworkingV1Alpha3Sidecars).
		MustAdd(IstioNetworkingV1Alpha3SyntheticServiceentries).
		MustAdd(IstioNetworkingV1Alpha3Virtualservices).
		MustAdd(IstioPolicyV1Beta1Attributemanifests).
		MustAdd(IstioPolicyV1Beta1Handlers).
		MustAdd(IstioPolicyV1Beta1Instances).
		MustAdd(IstioPolicyV1Beta1Rules).
		MustAdd(IstioRbacV1Alpha1Clusterrbacconfigs).
		MustAdd(IstioRbacV1Alpha1Rbacconfigs).
		MustAdd(IstioRbacV1Alpha1Servicerolebindings).
		MustAdd(IstioRbacV1Alpha1Serviceroles).
		MustAdd(IstioSecurityV1Beta1Authorizationpolicies).
		MustAdd(IstioSecurityV1Beta1Requestauthentications).
		MustAdd(K8SAppsV1Deployments).
		MustAdd(K8SAuthenticationIstioIoV1Alpha1Meshpolicies).
		MustAdd(K8SAuthenticationIstioIoV1Alpha1Policies).
		MustAdd(K8SConfigIstioIoV1Alpha2Adapters).
		MustAdd(K8SConfigIstioIoV1Alpha2Attributemanifests).
		MustAdd(K8SConfigIstioIoV1Alpha2Handlers).
		MustAdd(K8SConfigIstioIoV1Alpha2Httpapispecbindings).
		MustAdd(K8SConfigIstioIoV1Alpha2Httpapispecs).
		MustAdd(K8SConfigIstioIoV1Alpha2Instances).
		MustAdd(K8SConfigIstioIoV1Alpha2Quotaspecbindings).
		MustAdd(K8SConfigIstioIoV1Alpha2Quotaspecs).
		MustAdd(K8SConfigIstioIoV1Alpha2Rules).
		MustAdd(K8SConfigIstioIoV1Alpha2Templates).
		MustAdd(K8SCoreV1Endpoints).
		MustAdd(K8SCoreV1Namespaces).
		MustAdd(K8SCoreV1Nodes).
		MustAdd(K8SCoreV1Pods).
		MustAdd(K8SCoreV1Secrets).
		MustAdd(K8SCoreV1Services).
		MustAdd(K8SExtensionsV1Beta1Ingresses).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Destinationrules).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Envoyfilters).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Gateways).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Serviceentries).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Sidecars).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Virtualservices).
		MustAdd(K8SRbacIstioIoV1Alpha1Clusterrbacconfigs).
		MustAdd(K8SRbacIstioIoV1Alpha1Policy).
		MustAdd(K8SRbacIstioIoV1Alpha1Rbacconfigs).
		MustAdd(K8SRbacIstioIoV1Alpha1Serviceroles).
		MustAdd(K8SSecurityIstioIoV1Beta1Authorizationpolicies).
		MustAdd(K8SSecurityIstioIoV1Beta1Requestauthentications).
		Build()

	// Istio contains only Istio collections.
	Istio = collection.NewSchemasBuilder().
		MustAdd(IstioAuthenticationV1Alpha1Meshpolicies).
		MustAdd(IstioAuthenticationV1Alpha1Policies).
		MustAdd(IstioConfigV1Alpha2Adapters).
		MustAdd(IstioConfigV1Alpha2Httpapispecbindings).
		MustAdd(IstioConfigV1Alpha2Httpapispecs).
		MustAdd(IstioConfigV1Alpha2Templates).
		MustAdd(IstioMeshV1Alpha1MeshConfig).
		MustAdd(IstioMixerV1ConfigClientQuotaspecbindings).
		MustAdd(IstioMixerV1ConfigClientQuotaspecs).
		MustAdd(IstioNetworkingV1Alpha3Destinationrules).
		MustAdd(IstioNetworkingV1Alpha3Envoyfilters).
		MustAdd(IstioNetworkingV1Alpha3Gateways).
		MustAdd(IstioNetworkingV1Alpha3Serviceentries).
		MustAdd(IstioNetworkingV1Alpha3Sidecars).
		MustAdd(IstioNetworkingV1Alpha3SyntheticServiceentries).
		MustAdd(IstioNetworkingV1Alpha3Virtualservices).
		MustAdd(IstioPolicyV1Beta1Attributemanifests).
		MustAdd(IstioPolicyV1Beta1Handlers).
		MustAdd(IstioPolicyV1Beta1Instances).
		MustAdd(IstioPolicyV1Beta1Rules).
		MustAdd(IstioRbacV1Alpha1Clusterrbacconfigs).
		MustAdd(IstioRbacV1Alpha1Rbacconfigs).
		MustAdd(IstioRbacV1Alpha1Servicerolebindings).
		MustAdd(IstioRbacV1Alpha1Serviceroles).
		MustAdd(IstioSecurityV1Beta1Authorizationpolicies).
		MustAdd(IstioSecurityV1Beta1Requestauthentications).
		Build()

	// Kube contains only kubernetes collections.
	Kube = collection.NewSchemasBuilder().
		MustAdd(K8SAppsV1Deployments).
		MustAdd(K8SAuthenticationIstioIoV1Alpha1Meshpolicies).
		MustAdd(K8SAuthenticationIstioIoV1Alpha1Policies).
		MustAdd(K8SConfigIstioIoV1Alpha2Adapters).
		MustAdd(K8SConfigIstioIoV1Alpha2Attributemanifests).
		MustAdd(K8SConfigIstioIoV1Alpha2Handlers).
		MustAdd(K8SConfigIstioIoV1Alpha2Httpapispecbindings).
		MustAdd(K8SConfigIstioIoV1Alpha2Httpapispecs).
		MustAdd(K8SConfigIstioIoV1Alpha2Instances).
		MustAdd(K8SConfigIstioIoV1Alpha2Quotaspecbindings).
		MustAdd(K8SConfigIstioIoV1Alpha2Quotaspecs).
		MustAdd(K8SConfigIstioIoV1Alpha2Rules).
		MustAdd(K8SConfigIstioIoV1Alpha2Templates).
		MustAdd(K8SCoreV1Endpoints).
		MustAdd(K8SCoreV1Namespaces).
		MustAdd(K8SCoreV1Nodes).
		MustAdd(K8SCoreV1Pods).
		MustAdd(K8SCoreV1Secrets).
		MustAdd(K8SCoreV1Services).
		MustAdd(K8SExtensionsV1Beta1Ingresses).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Destinationrules).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Envoyfilters).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Gateways).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Serviceentries).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Sidecars).
		MustAdd(K8SNetworkingIstioIoV1Alpha3Virtualservices).
		MustAdd(K8SRbacIstioIoV1Alpha1Clusterrbacconfigs).
		MustAdd(K8SRbacIstioIoV1Alpha1Policy).
		MustAdd(K8SRbacIstioIoV1Alpha1Rbacconfigs).
		MustAdd(K8SRbacIstioIoV1Alpha1Serviceroles).
		MustAdd(K8SSecurityIstioIoV1Beta1Authorizationpolicies).
		MustAdd(K8SSecurityIstioIoV1Beta1Requestauthentications).
		Build()
)
