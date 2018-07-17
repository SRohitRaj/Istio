// Copyright 2017 Istio Authors
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

package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	authn "istio.io/api/authentication/v1alpha1"
	mccpb "istio.io/api/mixer/v1/config/client"
	networking "istio.io/api/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/model/test"
	"istio.io/istio/pkg/log"
)

const (
	// DefaultAuthenticationPolicyName is the name of the cluster-scoped authentication policy. Only
	// policy with this name in the cluster-scoped will be considered.
	DefaultAuthenticationPolicyName = "default"
)

// ConfigMeta is metadata attached to each configuration unit.
// The revision is optional, and if provided, identifies the
// last update operation on the object.
type ConfigMeta struct {
	// Type is a short configuration name that matches the content message type
	// (e.g. "route-rule")
	Type string `json:"type,omitempty"`

	// Group is the API group of the config.
	Group string `json:"group,omitempty"`

	// Version is the API version of the Config.
	Version string `json:"version,omitempty"`

	// Name is a unique immutable identifier in a namespace
	Name string `json:"name,omitempty"`

	// Namespace defines the space for names (optional for some types),
	// applications may choose to use namespaces for a variety of purposes
	// (security domains, fault domains, organizational domains)
	Namespace string `json:"namespace,omitempty"`

	// Domain defines the suffix of the fully qualified name past the namespace.
	// Domain is not a part of the unique key unlike name and namespace.
	Domain string `json:"domain,omitempty"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	Annotations map[string]string `json:"annotations,omitempty"`

	// ResourceVersion is an opaque identifier for tracking updates to the config registry.
	// The implementation may use a change index or a commit log for the revision.
	// The config client should not make any assumptions about revisions and rely only on
	// exact equality to implement optimistic concurrency of read-write operations.
	//
	// The lifetime of an object of a particular revision depends on the underlying data store.
	// The data store may compactify old revisions in the interest of storage optimization.
	//
	// An empty revision carries a special meaning that the associated object has
	// not been stored and assigned a revision.
	ResourceVersion string `json:"resourceVersion,omitempty"`

	// CreationTimestamp records the creation time
	CreationTimestamp meta_v1.Time `json:"resourceVersion,omitempty"`
}

// Config is a configuration unit consisting of the type of configuration, the
// key identifier that is unique per type, and the content represented as a
// protobuf message.
type Config struct {
	ConfigMeta

	// Spec holds the configuration object as a protobuf message
	Spec proto.Message
}

// ConfigStore describes a set of platform agnostic APIs that must be supported
// by the underlying platform to store and retrieve Istio configuration.
//
// Configuration key is defined to be a combination of the type, name, and
// namespace of the configuration object. The configuration key is guaranteed
// to be unique in the store.
//
// The storage interface presented here assumes that the underlying storage
// layer supports _Get_ (list), _Update_ (update), _Create_ (create) and
// _Delete_ semantics but does not guarantee any transactional semantics.
//
// _Update_, _Create_, and _Delete_ are mutator operations. These operations
// are asynchronous, and you might not see the effect immediately (e.g. _Get_
// might not return the object by key immediately after you mutate the store.)
// Intermittent errors might occur even though the operation succeeds, so you
// should always check if the object store has been modified even if the
// mutating operation returns an error.  Objects should be created with
// _Create_ operation and updated with _Update_ operation.
//
// Resource versions record the last mutation operation on each object. If a
// mutation is applied to a different revision of an object than what the
// underlying storage expects as defined by pure equality, the operation is
// blocked.  The client of this interface should not make assumptions about the
// structure or ordering of the revision identifier.
//
// Object references supplied and returned from this interface should be
// treated as read-only. Modifying them violates thread-safety.
type ConfigStore interface {
	// ConfigDescriptor exposes the configuration type schema known by the config store.
	// The type schema defines the bidrectional mapping between configuration
	// types and the protobuf encoding schema.
	ConfigDescriptor() ConfigDescriptor

	// Get retrieves a configuration element by a type and a key
	Get(typ, name, namespace string) (config *Config, exists bool)

	// List returns objects by type and namespace.
	// Use "" for the namespace to list across namespaces.
	List(typ, namespace string) ([]Config, error)

	// Create adds a new configuration object to the store. If an object with the
	// same name and namespace for the type already exists, the operation fails
	// with no side effects.
	Create(config Config) (revision string, err error)

	// Update modifies an existing configuration object in the store.  Update
	// requires that the object has been created.  Resource version prevents
	// overriding a value that has been changed between prior _Get_ and _Put_
	// operation to achieve optimistic concurrency. This method returns a new
	// revision if the operation succeeds.
	Update(config Config) (newRevision string, err error)

	// Delete removes an object from the store by key
	Delete(typ, name, namespace string) error
}

// Key function for the configuration objects
func Key(typ, name, namespace string) string {
	return fmt.Sprintf("%s/%s/%s", typ, namespace, name)
}

// Key is the unique identifier for a configuration object
func (meta *ConfigMeta) Key() string {
	return Key(meta.Type, meta.Name, meta.Namespace)
}

// ConfigStoreCache is a local fully-replicated cache of the config store.  The
// cache actively synchronizes its local state with the remote store and
// provides a notification mechanism to receive update events. As such, the
// notification handlers must be registered prior to calling _Run_, and the
// cache requires initial synchronization grace period after calling  _Run_.
//
// Update notifications require the following consistency guarantee: the view
// in the cache must be AT LEAST as fresh as the moment notification arrives, but
// MAY BE more fresh (e.g. if _Delete_ cancels an _Add_ event).
//
// Handlers execute on the single worker queue in the order they are appended.
// Handlers receive the notification event and the associated object.  Note
// that all handlers must be registered before starting the cache controller.
type ConfigStoreCache interface {
	ConfigStore

	// RegisterEventHandler adds a handler to receive config update events for a
	// configuration type
	RegisterEventHandler(typ string, handler func(Config, Event))

	// Run until a signal is received
	Run(stop <-chan struct{})

	// HasSynced returns true after initial cache synchronization is complete
	HasSynced() bool
}

// ConfigDescriptor defines the bijection between the short type name and its
// fully qualified protobuf message name
type ConfigDescriptor []ProtoSchema

// ProtoSchema provides description of the configuration schema and its key function
type ProtoSchema struct {
	// ClusterScoped is true for resource in cluster-level.
	ClusterScoped bool

	// Type is the config proto type.
	Type string

	// Plural is the type in plural.
	Plural string

	// Group is the config proto group.
	Group string

	// Version is the config proto version.
	Version string

	// MessageName refers to the protobuf message type name corresponding to the type
	MessageName string

	// Gogo is true for gogo protobuf messages
	Gogo bool

	// Validate configuration as a protobuf message assuming the object is an
	// instance of the expected message type
	Validate func(name, namespace string, config proto.Message) error
}

// Types lists all known types in the config schema
func (descriptor ConfigDescriptor) Types() []string {
	types := make([]string, 0, len(descriptor))
	for _, t := range descriptor {
		types = append(types, t.Type)
	}
	return types
}

// GetByType finds a schema by type if it is available
func (descriptor ConfigDescriptor) GetByType(name string) (ProtoSchema, bool) {
	for _, schema := range descriptor {
		if schema.Type == name {
			return schema, true
		}
	}
	return ProtoSchema{}, false
}

// IstioConfigStore is a specialized interface to access config store using
// Istio configuration types
type IstioConfigStore interface {
	ConfigStore

	// ServiceEntries lists all service entries
	ServiceEntries() []Config

	// DestinationRule returns a destination rule for a service name in a given domain.
	DestinationRule(hostname Hostname) *Config

	// VirtualServices lists all virtual services bound to the specified gateways
	VirtualServices(gateways map[string]bool) []Config

	// Gateways lists all gateways bound to the specified workload labels
	Gateways(workloadLabels LabelsCollection) []Config

	// SubsetToLabels returns the labels associated with a subset of a given service.
	SubsetToLabels(subsetName string, hostname Hostname) LabelsCollection

	// EnvoyFilter lists the envoy filter configuration bound to the specified workload labels
	EnvoyFilter(workloadLabels LabelsCollection) *Config

	// HTTPAPISpecByDestination selects Mixerclient HTTP API Specs
	// associated with destination service instances.
	HTTPAPISpecByDestination(instance *ServiceInstance) []Config

	// QuotaSpecByDestination selects Mixerclient quota specifications
	// associated with destination service instances.
	QuotaSpecByDestination(instance *ServiceInstance) []Config

	// AuthenticationPolicyByDestination selects authentication policy associated
	// with a service + port. Hostname must be FQDN.
	// If there are more than one policies at different scopes (global, namespace, service)
	// the one with the most specific scope will be selected. If there are more than
	// one with the same scope, the first one seen will be used (later, we should
	// have validation at submitting time to prevent this scenario from happening)
	AuthenticationPolicyByDestination(hostname Hostname, port *Port) *Config

	// ServiceRoles selects ServiceRoles in the specified namespace.
	ServiceRoles(namespace string) []Config

	// ServiceRoleBindings selects ServiceRoleBindings in the specified namespace.
	ServiceRoleBindings(namespace string) []Config

	// RbacConfig selects the RbacConfig of name DefaultRbacConfigName.
	RbacConfig() *Config
}

const (
	// IstioAPIGroupDomain defines API group domain of all Istio configuration resources.
	// Group domain suffix to the proto schema's group to generate the full resource group.
	IstioAPIGroupDomain = ".istio.io"

	// Default API version of an Istio config proto message.
	istioAPIVersion = "v1alpha2"

	// HeaderURI is URI HTTP header
	HeaderURI = "uri"

	// HeaderAuthority is authority HTTP header
	HeaderAuthority = "authority"

	// HeaderMethod is method HTTP header
	HeaderMethod = "method"

	// HeaderScheme is scheme HTTP header
	HeaderScheme = "scheme"

	// NamespaceAll is a designated symbol for listing across all namespaces
	NamespaceAll = ""

	// IstioMeshGateway is the built in gateway for all sidecars
	IstioMeshGateway = "mesh"
)

/*
  This conversion of CRD (== yaml files with k8s metadata) is extremely inefficient.
  The yaml is parsed (kubeyaml), converted to YAML again (FromJSONMap),
  converted to JSON (YAMLToJSON) and finally UnmarshallString in proto is called.

  The result is not cached in the model.

  In 0.7, this was the biggest factor in scalability. Moving forward we will likely
  deprecate model, and do the conversion (hopefully more efficient) only once, when
  an object is first read.
*/

var (
	// MockConfig is used purely for testing
	MockConfig = ProtoSchema{
		Type:        "mock-config",
		Plural:      "mock-configs",
		Group:       "test",
		Version:     "v1",
		MessageName: "test.MockConfig",
		Validate: func(name, namespace string, config proto.Message) error {
			if config.(*test.MockConfig).Key == "" {
				return errors.New("empty key")
			}
			return nil
		},
	}

	// VirtualService describes v1alpha3 route rules
	VirtualService = ProtoSchema{
		Type:        "virtual-service",
		Plural:      "virtual-services",
		Group:       "networking",
		Version:     "v1alpha3",
		MessageName: "istio.networking.v1alpha3.VirtualService",
		Gogo:        true,
		Validate:    ValidateVirtualService,
	}

	// Gateway describes a gateway (how a proxy is exposed on the network)
	Gateway = ProtoSchema{
		Type:        "gateway",
		Plural:      "gateways",
		Group:       "networking",
		Version:     "v1alpha3",
		MessageName: "istio.networking.v1alpha3.Gateway",
		Gogo:        true,
		Validate:    ValidateGateway,
	}

	// ServiceEntry describes service entries
	ServiceEntry = ProtoSchema{
		Type:        "service-entry",
		Plural:      "service-entries",
		Group:       "networking",
		Version:     "v1alpha3",
		MessageName: "istio.networking.v1alpha3.ServiceEntry",
		Gogo:        true,
		Validate:    ValidateServiceEntry,
	}

	// DestinationRule describes destination rules
	DestinationRule = ProtoSchema{
		Type:        "destination-rule",
		Plural:      "destination-rules",
		Group:       "networking",
		Version:     "v1alpha3",
		MessageName: "istio.networking.v1alpha3.DestinationRule",
		Validate:    ValidateDestinationRule,
	}

	// EnvoyFilter describes additional envoy filters to be inserted by Pilot
	EnvoyFilter = ProtoSchema{
		Type:        "envoy-filter",
		Plural:      "envoy-filters",
		Group:       "networking",
		Version:     "v1alpha3",
		MessageName: "istio.networking.v1alpha3.EnvoyFilter",
		Validate:    ValidateEnvoyFilter,
	}

	// HTTPAPISpec describes an HTTP API specification.
	HTTPAPISpec = ProtoSchema{
		Type:        "http-api-spec",
		Plural:      "http-api-specs",
		Group:       "config",
		Version:     istioAPIVersion,
		MessageName: "istio.mixer.v1.config.client.HTTPAPISpec",
		Validate:    ValidateHTTPAPISpec,
	}

	// HTTPAPISpecBinding describes an HTTP API specification binding.
	HTTPAPISpecBinding = ProtoSchema{
		Type:        "http-api-spec-binding",
		Plural:      "http-api-spec-bindings",
		Group:       "config",
		Version:     istioAPIVersion,
		MessageName: "istio.mixer.v1.config.client.HTTPAPISpecBinding",
		Validate:    ValidateHTTPAPISpecBinding,
	}

	// QuotaSpec describes an Quota specification.
	QuotaSpec = ProtoSchema{
		Type:        "quota-spec",
		Plural:      "quota-specs",
		Group:       "config",
		Version:     istioAPIVersion,
		MessageName: "istio.mixer.v1.config.client.QuotaSpec",
		Validate:    ValidateQuotaSpec,
	}

	// QuotaSpecBinding describes an Quota specification binding.
	QuotaSpecBinding = ProtoSchema{
		Type:        "quota-spec-binding",
		Plural:      "quota-spec-bindings",
		Group:       "config",
		Version:     istioAPIVersion,
		MessageName: "istio.mixer.v1.config.client.QuotaSpecBinding",
		Validate:    ValidateQuotaSpecBinding,
	}

	// AuthenticationPolicy describes an authentication policy.
	AuthenticationPolicy = ProtoSchema{
		Type:        "policy",
		Plural:      "policies",
		Group:       "authentication",
		Version:     "v1alpha1",
		MessageName: "istio.authentication.v1alpha1.Policy",
		Validate:    ValidateAuthenticationPolicy,
	}

	// AuthenticationMeshPolicy describes an authentication policy at mesh level.
	AuthenticationMeshPolicy = ProtoSchema{
		ClusterScoped: true,
		Type:          "mesh-policy",
		Plural:        "mesh-policies",
		Group:         "authentication",
		Version:       "v1alpha1",
		MessageName:   "istio.authentication.v1alpha1.Policy",
		Validate:      ValidateAuthenticationPolicy,
	}

	// ServiceRole describes an RBAC service role.
	ServiceRole = ProtoSchema{
		Type:        "service-role",
		Plural:      "service-roles",
		Group:       "config",
		Version:     istioAPIVersion,
		MessageName: "istio.rbac.v1alpha1.ServiceRole",
		Validate:    ValidateServiceRole,
	}

	// ServiceRoleBinding describes an RBAC service role.
	ServiceRoleBinding = ProtoSchema{
		ClusterScoped: false,
		Type:          "service-role-binding",
		Plural:        "service-role-bindings",
		Group:         "config",
		Version:       istioAPIVersion,
		MessageName:   "istio.rbac.v1alpha1.ServiceRoleBinding",
		Validate:      ValidateServiceRoleBinding,
	}

	// RbacConfig describes the mesh level RBAC config.
	RbacConfig = ProtoSchema{
		ClusterScoped: true,
		Type:          "rbac-config",
		Plural:        "rbac-configs",
		Group:         "config",
		Version:       istioAPIVersion,
		MessageName:   "istio.rbac.v1alpha1.RbacConfig",
		Validate:      ValidateRbacConfig,
	}

	// IstioConfigTypes lists all Istio config types with schemas and validation
	IstioConfigTypes = ConfigDescriptor{
		VirtualService,
		Gateway,
		ServiceEntry,
		DestinationRule,
		EnvoyFilter,
		HTTPAPISpec,
		HTTPAPISpecBinding,
		QuotaSpec,
		QuotaSpecBinding,
		AuthenticationPolicy,
		AuthenticationMeshPolicy,
		ServiceRole,
		ServiceRoleBinding,
		RbacConfig,
	}
)

// ResolveHostname produces a FQDN based on either the service or
// a concat of the namespace + domain
func ResolveHostname(meta ConfigMeta, svc *mccpb.IstioService) Hostname {
	out := svc.Name
	// if FQDN is specified, do not append domain or namespace to hostname
	// Service field has precedence over Name
	if svc.Service != "" {
		out = svc.Service
	} else {
		if svc.Namespace != "" {
			out = out + "." + svc.Namespace
		} else if meta.Namespace != "" {
			out = out + "." + meta.Namespace
		}

		if svc.Domain != "" {
			out = out + "." + svc.Domain
		} else if meta.Domain != "" {
			out = out + ".svc." + meta.Domain
		}
	}

	return Hostname(out)
}

// ResolveShortnameToFQDN uses metadata information to resolve a reference
// to shortname of the service to FQDN
func ResolveShortnameToFQDN(host string, meta ConfigMeta) Hostname {
	out := host
	// Treat the wildcard host as fully qualified. Any other variant of a wildcard hostname will contain a `.` too,
	// and skip the next if, so we only need to check for the literal wildcard itself.
	if host == "*" {
		return Hostname(out)
	}
	// if FQDN is specified, do not append domain or namespace to hostname
	if !strings.Contains(host, ".") {
		if meta.Namespace != "" {
			out = out + "." + meta.Namespace
		}

		// FIXME this is a gross hack to hardcode a service's domain name in kubernetes
		// BUG this will break non kubernetes environments if they use shortnames in the
		// rules.
		if meta.Domain != "" {
			out = out + ".svc." + meta.Domain
		}
	}

	return Hostname(out)
}

// MostSpecificHostMatch compares the elements of the stack to the needle, and returns the longest stack element
// matching the needle, or false if no element in the stack matches the needle.
func MostSpecificHostMatch(needle Hostname, stack []Hostname) (Hostname, bool) {
	sort.Sort(Hostnames(stack))
	for _, h := range stack {
		if needle.Matches(h) {
			return h, true
		}
	}
	return "", false
}

// istioConfigStore provides a simple adapter for Istio configuration types
// from the generic config registry
type istioConfigStore struct {
	ConfigStore
}

// MakeIstioStore creates a wrapper around a store.
// In pilot it is initialized with a ConfigStoreCache, tests only use
// a regular ConfigStore.
func MakeIstioStore(store ConfigStore) IstioConfigStore {
	return &istioConfigStore{store}
}

func (store *istioConfigStore) ServiceEntries() []Config {
	configs, err := store.List(ServiceEntry.Type, NamespaceAll)
	if err != nil {
		return nil
	}
	return configs
}

// TODO: move the logic to v2, read all VirtualServices once at startup and per
// change event and pass them to the config generator. Model calls to List are
// extremely expensive - and for larger number of services it doesn't make sense
// to just convert again and again, for each listener times endpoints.
func (store *istioConfigStore) VirtualServices(gateways map[string]bool) []Config {
	configs, err := store.List(VirtualService.Type, NamespaceAll)
	if err != nil {
		log.Warnf("Could not load VirtualServices. Error:\n %s \n", err)
		return nil
	}

	// TODO Sort configs by creation time
	out := make([]Config, 0)
	for _, config := range configs {
		rule := config.Spec.(*networking.VirtualService)
		if len(rule.Gateways) == 0 {
			// This rule applies only to IstioMeshGateway
			if gateways[IstioMeshGateway] {
				out = append(out, config)
			}
		} else {
			for _, g := range rule.Gateways {
				// note: Gateway names do _not_ use wildcard matching, so we do not use Hostname.Matches here
				if gateways[ResolveShortnameToFQDN(g, config.ConfigMeta).String()] {
					out = append(out, config)
					break
				} else if g == IstioMeshGateway && gateways[g] {
					// "mesh" gateway cannot be expanded into FQDN
					out = append(out, config)
					break
				}
			}
		}
	}

	// Need to parse each rule and convert the shortname to FQDN
	for _, r := range out {
		rule := r.Spec.(*networking.VirtualService)
		// resolve top level hosts
		for i, h := range rule.Hosts {
			rule.Hosts[i] = ResolveShortnameToFQDN(h, r.ConfigMeta).String()
		}
		// resolve gateways to bind to
		for i, g := range rule.Gateways {
			if g != IstioMeshGateway {
				rule.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
			}
		}
		// resolve host in http route.destination, route.mirror
		for _, d := range rule.Http {
			for _, m := range d.Match {
				for i, g := range m.Gateways {
					if g != IstioMeshGateway {
						m.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
					}
				}
			}
			for _, w := range d.Route {
				w.Destination.Host = ResolveShortnameToFQDN(w.Destination.Host, r.ConfigMeta).String()
			}
			if d.Mirror != nil {
				d.Mirror.Host = ResolveShortnameToFQDN(d.Mirror.Host, r.ConfigMeta).String()
			}
		}
		//resolve host in tcp route.destination
		for _, d := range rule.Tcp {
			for _, m := range d.Match {
				for i, g := range m.Gateways {
					if g != IstioMeshGateway {
						m.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
					}
				}
			}
			for _, w := range d.Route {
				w.Destination.Host = ResolveShortnameToFQDN(w.Destination.Host, r.ConfigMeta).String()
			}
		}
		//resolve host in tls route.destination
		for _, tls := range rule.Tls {
			for _, m := range tls.Match {
				for i, g := range m.Gateways {
					if g != IstioMeshGateway {
						m.Gateways[i] = ResolveShortnameToFQDN(g, r.ConfigMeta).String()
					}
				}
			}
			for _, w := range tls.Route {
				w.Destination.Host = ResolveShortnameToFQDN(w.Destination.Host, r.ConfigMeta).String()
			}
		}
	}

	return out
}

func (store *istioConfigStore) Gateways(workloadLabels LabelsCollection) []Config {
	configs, err := store.List(Gateway.Type, NamespaceAll)
	if err != nil {
		return nil
	}

	out := make([]Config, 0)
	for _, config := range configs {
		gateway := config.Spec.(*networking.Gateway)
		if gateway.GetSelector() == nil {
			// no selector. Applies to all workloads asking for the gateway
			out = append(out, config)
		} else {
			gatewaySelector := Labels(gateway.GetSelector())
			if workloadLabels.IsSupersetOf(gatewaySelector) {
				out = append(out, config)
			}
		}
	}
	return out
}

// NOTE: There can be only one filter for a workload. If multiple filters are defined, the behavior
// is undefined.
func (store *istioConfigStore) EnvoyFilter(workloadLabels LabelsCollection) *Config {
	configs, err := store.List(EnvoyFilter.Type, NamespaceAll)
	if err != nil {
		return nil
	}

	for _, config := range configs {
		filter := config.Spec.(*networking.EnvoyFilter)
		if filter.GetWorkloadLabels() == nil {
			// no selector. Applies to all workloads asking for the gateway
			return &config
		}
		workloadSelector := Labels(filter.GetWorkloadLabels())
		if workloadLabels.IsSupersetOf(workloadSelector) {
			return &config
		}
	}
	return nil
}

func (store *istioConfigStore) DestinationRule(hostname Hostname) *Config {
	configs, err := store.List(DestinationRule.Type, NamespaceAll)
	if err != nil {
		return nil
	}

	hosts := make([]Hostname, len(configs))
	byHosts := make(map[Hostname]*Config, len(configs))
	for i := range configs {
		rule := configs[i].Spec.(*networking.DestinationRule)
		hosts[i] = ResolveShortnameToFQDN(rule.Host, configs[i].ConfigMeta)
		byHosts[hosts[i]] = &configs[i]
	}

	if c, ok := MostSpecificHostMatch(hostname, hosts); ok {
		return byHosts[c]
	}
	return nil
}

func (store *istioConfigStore) SubsetToLabels(subsetName string, hostname Hostname) LabelsCollection {
	// empty subset
	if subsetName == "" {
		return nil
	}

	config := store.DestinationRule(hostname)
	if config == nil {
		return nil
	}

	rule := config.Spec.(*networking.DestinationRule)
	for _, subset := range rule.Subsets {
		if subset.Name == subsetName {
			return []Labels{subset.Labels}
		}
	}

	return nil
}

// HTTPAPISpecByDestination selects Mixerclient HTTP API Specs
// associated with destination service instances.
func (store *istioConfigStore) HTTPAPISpecByDestination(instance *ServiceInstance) []Config {
	bindings, err := store.List(HTTPAPISpecBinding.Type, NamespaceAll)
	if err != nil {
		return nil
	}
	specs, err := store.List(HTTPAPISpec.Type, NamespaceAll)
	if err != nil {
		return nil
	}

	// Create a set key from a reference's name and namespace.
	key := func(name, namespace string) string { return name + "/" + namespace }

	// Build the set of HTTP API spec references bound to the service instance.
	refs := make(map[string]struct{})
	for _, binding := range bindings {
		b := binding.Spec.(*mccpb.HTTPAPISpecBinding)
		for _, service := range b.Services {
			hostname := ResolveHostname(binding.ConfigMeta, service)
			if hostname == instance.Service.Hostname {
				for _, spec := range b.ApiSpecs {
					namespace := spec.Namespace
					if namespace == "" {
						namespace = binding.Namespace
					}
					refs[key(spec.Name, namespace)] = struct{}{}
				}
			}
		}
	}

	// Append any spec that is in the set of references.
	var out []Config
	for _, spec := range specs {
		if _, ok := refs[key(spec.ConfigMeta.Name, spec.ConfigMeta.Namespace)]; ok {
			out = append(out, spec)
		}
	}

	return out
}

// QuotaSpecByDestination selects Mixerclient quota specifications
// associated with destination service instances.
func (store *istioConfigStore) QuotaSpecByDestination(instance *ServiceInstance) []Config {
	bindings, err := store.List(QuotaSpecBinding.Type, NamespaceAll)
	if err != nil {
		return nil
	}
	specs, err := store.List(QuotaSpec.Type, NamespaceAll)
	if err != nil {
		return nil
	}

	// Create a set key from a reference's name and namespace.
	key := func(name, namespace string) string { return name + "/" + namespace }

	// Build the set of quota spec references bound to the service instance.
	refs := make(map[string]struct{})
	for _, binding := range bindings {
		b := binding.Spec.(*mccpb.QuotaSpecBinding)
		for _, service := range b.Services {
			hostname := ResolveHostname(binding.ConfigMeta, service)
			if hostname == instance.Service.Hostname {
				for _, spec := range b.QuotaSpecs {
					namespace := spec.Namespace
					if namespace == "" {
						namespace = binding.Namespace
					}
					refs[key(spec.Name, namespace)] = struct{}{}
				}
			}
		}
	}

	// Append any spec that is in the set of references.
	var out []Config
	for _, spec := range specs {
		if _, ok := refs[key(spec.ConfigMeta.Name, spec.ConfigMeta.Namespace)]; ok {
			out = append(out, spec)
		}
	}

	return out
}

func (store *istioConfigStore) AuthenticationPolicyByDestination(hostname Hostname, port *Port) *Config {
	// Hostname should be FQDN, so namespace can be extracted by parsing hostname.
	parts := strings.Split(string(hostname), ".")
	if len(parts) < 2 {
		// Bad hostname, return no policy.
		return nil
	}
	namespace := parts[1]
	specs, err := store.List(AuthenticationPolicy.Type, namespace)
	if err != nil {
		return nil
	}
	var out Config
	currentMatchLevel := 0
	for _, spec := range specs {
		policy := spec.Spec.(*authn.Policy)
		// Indicate if a policy matched to target destination:
		// 0 - not match.
		// 1 - global / cluster scope.
		// 2 - namespace scope.
		// 3 - workload (service).
		matchLevel := 0
		if len(policy.Targets) > 0 {
			for _, dest := range policy.Targets {
				if hostname != ResolveHostname(spec.ConfigMeta, &mccpb.IstioService{Name: dest.Name}) {
					continue
				}
				// If destination port is defined, it must match.
				if len(dest.Ports) > 0 {
					portMatched := false
					for _, portSelector := range dest.Ports {
						if port.Match(portSelector) {
							portMatched = true
							break
						}
					}
					if !portMatched {
						// Port does not match with any of port selector, skip to next target selector.
						continue
					}
				}

				matchLevel = 3
				break
			}
		} else {
			// Match on namespace level.
			matchLevel = 2
		}
		// Swap output policy that is match in more specific scope.
		if matchLevel > currentMatchLevel {
			currentMatchLevel = matchLevel
			out = spec
		}
	}
	// Non-zero currentMatchLevel implies authentication policy was found for the given host.
	if currentMatchLevel != 0 {
		return &out
	}

	// Reach here if no authentication policy found in service or namespace level; check for
	// cluster-scoped (global) policy.
	// Note: to avoid multiple global policy, we restrict that only the one with name equals to
	// `DefaultAuthenticationPolicyName` ("default") will be used. Also, targets spec should be empty.
	if specs, err := store.List(AuthenticationMeshPolicy.Type, ""); err == nil {
		for _, spec := range specs {
			if spec.Name == DefaultAuthenticationPolicyName {
				return &spec
			}
		}
	}

	return nil
}

func (store *istioConfigStore) ServiceRoles(namespace string) []Config {
	roles, err := store.List(ServiceRole.Type, namespace)
	if err != nil {
		log.Errorf("failed to get ServiceRoles in namespace %s: %v", namespace, err)
		return nil
	}

	return roles
}

func (store *istioConfigStore) ServiceRoleBindings(namespace string) []Config {
	bindings, err := store.List(ServiceRoleBinding.Type, namespace)
	if err != nil {
		log.Errorf("failed to get ServiceRoleBinding in namespace %s: %v", namespace, err)
		return nil
	}

	return bindings
}

func (store *istioConfigStore) RbacConfig() *Config {
	rbacConfigs, err := store.List(RbacConfig.Type, "")
	if err != nil {
		log.Errorf("failed to get rbacConfig: %v", err)
		return nil
	}
	if len(rbacConfigs) > 1 {
		log.Errorf("found %d RbacConfigs, expecting only 1.", len(rbacConfigs))
	}
	for _, rc := range rbacConfigs {
		if rc.Name == DefaultRbacConfigName {
			return &rc
		}
	}
	return nil
}

// SortHTTPAPISpec sorts a slice in a stable manner.
func SortHTTPAPISpec(specs []Config) {
	sort.Slice(specs, func(i, j int) bool {
		// protect against incompatible types
		irule, _ := specs[i].Spec.(*mccpb.HTTPAPISpec)
		jrule, _ := specs[j].Spec.(*mccpb.HTTPAPISpec)
		return irule == nil || jrule == nil || (specs[i].Key() < specs[j].Key())
	})
}

// SortQuotaSpec sorts a slice in a stable manner.
func SortQuotaSpec(specs []Config) {
	sort.Slice(specs, func(i, j int) bool {
		// protect against incompatible types
		irule, _ := specs[i].Spec.(*mccpb.QuotaSpec)
		jrule, _ := specs[j].Spec.(*mccpb.QuotaSpec)
		return irule == nil || jrule == nil || (specs[i].Key() < specs[j].Key())
	})
}
