// Copyright 2019 Istio Authors
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

package mtls

import (
	"fmt"

	"istio.io/istio/galley/pkg/config/analysis/analyzers/util"
	"istio.io/istio/security/proto/authentication/v1alpha1"
)

// TargetService is a simple struct type for representing a service
// targeted by an Authentication policy.
type TargetService struct {
	fQDN string

	// Either portNumber/portName will be set, or neither. The constructors
	// prevent both from ever being set.
	portNumber uint32
	portName   string
}

// NewTargetServiceWithPortNumber creates a new TargetService using the specified
// fqdn and portNumber.
func NewTargetServiceWithPortNumber(fqdn string, portNumber uint32) TargetService {
	return TargetService{fQDN: fqdn, portNumber: portNumber}
}

// NewTargetServiceWithPortName creates a new TargetService using the specified fqdn
// and portName.
func NewTargetServiceWithPortName(fqdn, portName string) TargetService {
	return TargetService{fQDN: fqdn, portName: portName}
}

// NewTargetService creates a new TargetService using the specified fqdn. Because no
// port is specified, this implicitly represents the service bound to any port.
func NewTargetService(fqdn string) TargetService {
	return TargetService{fQDN: fqdn}
}

// FQDN is the fully-qualified domain name for the service (e.g.
// foobar.my-namespace.svc.cluster.local).
func (w TargetService) FQDN() string {
	return w.fQDN
}

// PortNumber is the port used by the service. If set (non-zero), then
// PortName must be the default value ("").
func (w TargetService) PortNumber() uint32 {
	return w.portNumber
}

// PortName is the name of the port used by the service. If set (not ""),
// then PortNumber must be the default value (0).
func (w TargetService) PortName() string {
	return w.portName
}

// PolicyChecker allows callers to add a set of v1alpha1.Policy objects in the
// mesh. Once these are loaded, you can query whether or not a specific
// TargetService will require MTLS when an incoming connection occurs using the
// IsServiceMTLSEnforced() call.
type PolicyChecker struct {
	// meshHasStrictMTLSPolicy tracks whether or not mTLS is strictly enforced on the mesh.
	meshHasStrictMTLSPolicy bool

	namespaceToMTLSMode map[string]strictMTLSMode
	serviceToMTLSMode   map[TargetService]strictMTLSMode
}

// strictMTLSMode is a helper type used to represent whether or not MTLS was
// explicitly enabled, explicitly disabled, or just not present. Useful for the
// IsServiceMTLSEnforced check to determine which policy applies.
type strictMTLSMode int

const (
	strictMTLSUnset strictMTLSMode = iota
	strictMTLSExplicitlyEnabled
	strictMTLSExplicitlyDisabled
)

// NewPolicyChecker creates a new PolicyChecker instance.
func NewPolicyChecker() *PolicyChecker {
	return &PolicyChecker{
		namespaceToMTLSMode: make(map[string]strictMTLSMode),
		serviceToMTLSMode:   make(map[TargetService]strictMTLSMode),
	}
}

// AddMeshPolicy adds a mesh-level policy to the checker. Note that there can
// only be at most one mesh level policy in effect.
func (pc *PolicyChecker) AddMeshPolicy(p *v1alpha1.Policy) {
	pc.meshHasStrictMTLSPolicy = doesPolicyEnforceMTLS(p)
}

// AddPolicy adds a new policy object to the PolicyChecker to use when later
// determining if a service is MTLS-enforced. The namespace of the policy is
// also provided as some policies can target the local namespace.
func (pc *PolicyChecker) AddPolicy(namespace string, p *v1alpha1.Policy) error {
	var policyMode strictMTLSMode
	if doesPolicyEnforceMTLS(p) {
		policyMode = strictMTLSExplicitlyEnabled
	} else {
		policyMode = strictMTLSExplicitlyDisabled
	}

	if len(p.Targets) == 0 {
		// Rule targets the namespace.
		pc.namespaceToMTLSMode[namespace] = policyMode
		return nil
	}
	// Discover the targeted service and take note. Should normalize.
	for _, target := range p.Targets {
		fqdn := util.ConvertHostToFQDN(namespace, target.Name)

		if len(target.Ports) == 0 {
			// Policy targets all ports on service
			pc.serviceToMTLSMode[NewTargetService(fqdn)] = policyMode
		}

		for _, port := range target.Ports {
			if port.GetName() != "" {
				pc.serviceToMTLSMode[NewTargetServiceWithPortName(fqdn, port.GetName())] = policyMode
			} else if port.GetNumber() != 0 {
				pc.serviceToMTLSMode[NewTargetServiceWithPortNumber(fqdn, port.GetNumber())] = policyMode
			} else {
				// Unhandled case!
				return fmt.Errorf("policy has a port with no name/number for target %s", target.Name)
			}
		}
	}

	return nil
}

// IsServiceMTLSEnforced returns true if a service requires incoming
// connections to use MTLS, or false if MTLS is not a hard-requirement (e.g.
// mode is permissive, peerIsOptional is true, etc). Only call this after adding
// all policy resources in effect via AddPolicy or AddMeshPolicy.
func (pc *PolicyChecker) IsServiceMTLSEnforced(w TargetService) (bool, error) {
	// TODO support understanding port name -> port number mappings
	switch pc.serviceToMTLSMode[w] {
	case strictMTLSExplicitlyEnabled:
		return true, nil
	case strictMTLSExplicitlyDisabled:
		return false, nil
	case strictMTLSUnset:
		// Fall through switch case
	default:
		return false, fmt.Errorf("unknown strictMTLSMode: %v", pc.serviceToMTLSMode[w])
	}

	// Try checking if its enforced on any ports
	serviceNoPort := NewTargetService(w.FQDN())
	switch pc.serviceToMTLSMode[serviceNoPort] {
	case strictMTLSExplicitlyEnabled:
		return true, nil
	case strictMTLSExplicitlyDisabled:
		return false, nil
	case strictMTLSUnset:
		// Fall through switch case
	default:
		return false, fmt.Errorf("unknown strictMTLSMode: %v", pc.serviceToMTLSMode[serviceNoPort])
	}

	// Check if enforced on namespace
	namespace, _ := util.GetResourceNameFromHost("", w.FQDN()).InterpretAsNamespaceAndName()
	if namespace == "" {
		return false, fmt.Errorf("unable to extract namespace from fqdn: %s", w.FQDN())
	}
	switch pc.namespaceToMTLSMode[namespace] {
	case strictMTLSExplicitlyEnabled:
		return true, nil
	case strictMTLSExplicitlyDisabled:
		return false, nil
	case strictMTLSUnset:
		// Fall through switch case
	default:
		return false, fmt.Errorf("unknown strictMTLSMode: %v", pc.namespaceToMTLSMode[namespace])
	}

	// Finally, defer to mesh level policy
	return pc.meshHasStrictMTLSPolicy, nil
}

// doesPolicyEnforceMTLS is a helper function to determine whether or not a
// Policy implies Strict MTLS mode.
func doesPolicyEnforceMTLS(p *v1alpha1.Policy) bool {
	if p.PeerIsOptional {
		// Connection can still occur.
		return false
	}
	for _, peer := range p.Peers {
		mtlsParams, ok := peer.Params.(*v1alpha1.PeerAuthenticationMethod_Mtls)
		if !ok {
			// Only looking for mtls methods
			continue
		}

		// The default value if no Mtls is specified on mtlsParams is strict.
		// If we do have parameters, though, ensure they do not imply permissive mode.
		return mtlsParams.Mtls == nil ||
			(!mtlsParams.Mtls.AllowTls && mtlsParams.Mtls.Mode != v1alpha1.MutualTls_PERMISSIVE)
	}

	// No MTLS configuration found
	return false
}
