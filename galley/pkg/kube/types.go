// GENERATED FILE -- DO NOT EDIT
//
//go:generate $GOPATH/src/istio.io/istio/galley/tools/gen-meta/gen-meta.sh kube pkg/kube/types.go
//

package kube

import "istio.io/istio/galley/pkg/kube/converter"

// Types in the schema.
var Types = Schema{}

func init() {

	Types.add(ResourceSpec{
		Kind:      "AuthenticationPolicy",
		ListKind:  "AuthenticationPolicyList",
		Singular:  "policy",
		Plural:    "policies",
		Version:   "v1alpha1",
		Group:     "authentication",
		Target:    getTargetFor("istio.authentication.v1alpha1.Policy"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "AuthenticationMeshPolicy",
		ListKind:  "AuthenticationMeshPolicyList",
		Singular:  "mesh-policy",
		Plural:    "mesh-policies",
		Version:   "v1alpha1",
		Group:     "authentication",
		Target:    getTargetFor("istio.authentication.v1alpha1.Policy"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "servicecontrol",
		ListKind:  "servicecontrolList",
		Singular:  "servicecontrol",
		Plural:    "servicecontrols",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "fluentd",
		ListKind:  "fluentdList",
		Singular:  "fluentd",
		Plural:    "fluentds",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "opa",
		ListKind:  "opaList",
		Singular:  "opa",
		Plural:    "opas",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "circonus",
		ListKind:  "circonusList",
		Singular:  "circonus",
		Plural:    "circonuses",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "statsd",
		ListKind:  "statsdList",
		Singular:  "statsd",
		Plural:    "statsds",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "bypass",
		ListKind:  "bypassList",
		Singular:  "bypass",
		Plural:    "bypasses",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "adapter",
		ListKind:  "adapterList",
		Singular:  "adapter",
		Plural:    "adapters",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "prometheus",
		ListKind:  "prometheusList",
		Singular:  "prometheus",
		Plural:    "prometheuses",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "denier",
		ListKind:  "denierList",
		Singular:  "denier",
		Plural:    "deniers",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "stdio",
		ListKind:  "stdioList",
		Singular:  "stdio",
		Plural:    "stdios",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "rbac",
		ListKind:  "rbacList",
		Singular:  "rbac",
		Plural:    "rbacs",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "kubernetesenv",
		ListKind:  "kubernetesenvList",
		Singular:  "kubernetesenv",
		Plural:    "kubernetesenvs",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "listchecker",
		ListKind:  "listcheckerList",
		Singular:  "listchecker",
		Plural:    "listcheckers",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "stackdriver",
		ListKind:  "stackdriverList",
		Singular:  "stackdriver",
		Plural:    "stackdrivers",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "solarwinds",
		ListKind:  "solarwindsList",
		Singular:  "solarwinds",
		Plural:    "solarwindses",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "signalfx",
		ListKind:  "signalfxList",
		Singular:  "signalfx",
		Plural:    "signalfxs",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "noop",
		ListKind:  "noopList",
		Singular:  "noop",
		Plural:    "noops",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "memquota",
		ListKind:  "memquotaList",
		Singular:  "memquota",
		Plural:    "memquotas",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Info"),
		Converter: converter.Get("old-mixer-adapter"),
	})

	Types.add(ResourceSpec{
		Kind:      "tracespan",
		ListKind:  "tracespanList",
		Singular:  "tracespan",
		Plural:    "tracespans",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "apikey",
		ListKind:  "apikeyList",
		Singular:  "apikey",
		Plural:    "apikeys",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "authorization",
		ListKind:  "authorizationList",
		Singular:  "authorization",
		Plural:    "authorizations",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "template",
		ListKind:  "templateList",
		Singular:  "template",
		Plural:    "templates",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "checknothing",
		ListKind:  "checknothingList",
		Singular:  "checknothing",
		Plural:    "checknothings",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "logentry",
		ListKind:  "logentryList",
		Singular:  "logentry",
		Plural:    "logentries",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "metric",
		ListKind:  "metricList",
		Singular:  "metric",
		Plural:    "metrics",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "quota",
		ListKind:  "quotaList",
		Singular:  "quota",
		Plural:    "quotas",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "reportnothing",
		ListKind:  "reportnothingList",
		Singular:  "reportnothing",
		Plural:    "reportnothings",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "servicecontrolreport",
		ListKind:  "servicecontrolreportList",
		Singular:  "servicecontrolreport",
		Plural:    "servicecontrolreports",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "kubernetes",
		ListKind:  "kubernetesList",
		Singular:  "kubernetes",
		Plural:    "kuberneteses",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "listentry",
		ListKind:  "listentryList",
		Singular:  "listentry",
		Plural:    "listentries",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.adapter.model.v1beta1.Template"),
		Converter: converter.Get("old-mixer-template"),
	})

	Types.add(ResourceSpec{
		Kind:      "HTTPAPISpec",
		ListKind:  "HTTPAPISpecList",
		Singular:  "http-api-spec",
		Plural:    "http-api-specs",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.v1.config.client.HTTPAPISpec"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "HTTPAPISpecBinding",
		ListKind:  "HTTPAPISpecBindingList",
		Singular:  "http-api-spec-binding",
		Plural:    "http-api-spec-bindings",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.v1.config.client.HTTPAPISpecBinding"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "QuotaSpec",
		ListKind:  "QuotaSpecList",
		Singular:  "quota-spec",
		Plural:    "quota-specs",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.v1.config.client.QuotaSpec"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "QuotaSpecBinding",
		ListKind:  "QuotaSpecBindingList",
		Singular:  "quota-spec-binding",
		Plural:    "quota-spec-bindings",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.mixer.v1.config.client.QuotaSpecBinding"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "DestinationRule",
		ListKind:  "DestinationRuleList",
		Singular:  "destination-rule",
		Plural:    "destination-rules",
		Version:   "v1alpha3",
		Group:     "networking",
		Target:    getTargetFor("istio.networking.v1alpha3.DestinationRule"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "EnvoyFilter",
		ListKind:  "EnvoyFilterList",
		Singular:  "envoy-filter",
		Plural:    "envoy-filters",
		Version:   "v1alpha3",
		Group:     "networking",
		Target:    getTargetFor("istio.networking.v1alpha3.EnvoyFilter"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "Gateway",
		ListKind:  "GatewayList",
		Singular:  "gateway",
		Plural:    "gateways",
		Version:   "v1alpha3",
		Group:     "networking",
		Target:    getTargetFor("istio.networking.v1alpha3.Gateway"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "ServiceEntry",
		ListKind:  "ServiceEntryList",
		Singular:  "service-entry",
		Plural:    "service-entries",
		Version:   "v1alpha3",
		Group:     "networking",
		Target:    getTargetFor("istio.networking.v1alpha3.ServiceEntry"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "VirtualService",
		ListKind:  "VirtualServiceList",
		Singular:  "virtual-service",
		Plural:    "virtual-services",
		Version:   "v1alpha3",
		Group:     "networking",
		Target:    getTargetFor("istio.networking.v1alpha3.VirtualService"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "attributemanifest",
		ListKind:  "attributemanifestList",
		Singular:  "attributemanifest",
		Plural:    "attributemanifests",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.policy.v1beta1.AttributeManifest"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "handler",
		ListKind:  "handlerList",
		Singular:  "handler",
		Plural:    "handlers",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.policy.v1beta1.Handler"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "instance",
		ListKind:  "instanceList",
		Singular:  "instance",
		Plural:    "instances",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.policy.v1beta1.Instance"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "rule",
		ListKind:  "ruleList",
		Singular:  "rule",
		Plural:    "rules",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.policy.v1beta1.Rule"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "RbacConfig",
		ListKind:  "RbacConfigList",
		Singular:  "rbac-config",
		Plural:    "rbac-configs",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.rbac.v1alpha1.RbacConfig"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "ServiceRole",
		ListKind:  "ServiceRoleList",
		Singular:  "service-role",
		Plural:    "service-roles",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.rbac.v1alpha1.ServiceRole"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "ServiceRoleBinding",
		ListKind:  "ServiceRoleBindingList",
		Singular:  "service-role-binding",
		Plural:    "service-role-bindings",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.rbac.v1alpha1.ServiceRoleBinding"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "DestinationPolicy",
		ListKind:  "DestinationPolicyList",
		Singular:  "destination-policy",
		Plural:    "destination-policies",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.routing.v1alpha1.DestinationPolicy"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "EgressRule",
		ListKind:  "EgressRuleList",
		Singular:  "egress-rule",
		Plural:    "egress-rules",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.routing.v1alpha1.EgressRule"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "IngressRule",
		ListKind:  "IngressRuleList",
		Singular:  "ingress-rule",
		Plural:    "ingress-rules",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.routing.v1alpha1.IngressRule"),
		Converter: converter.Get("identity"),
	})

	Types.add(ResourceSpec{
		Kind:      "RouteRule",
		ListKind:  "RouteRuleList",
		Singular:  "route-rule",
		Plural:    "route-rules",
		Version:   "v1alpha2",
		Group:     "config",
		Target:    getTargetFor("istio.routing.v1alpha1.RouteRule"),
		Converter: converter.Get("identity"),
	})

}
