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

package model

import (
	"istio.io/api/security/v1beta1"

	"istio.io/istio/pkg/config/labels"
	"istio.io/istio/pkg/config/schema/collections"
)

// MutualTLSMode is the mutule TLS mode specified by authentication policy.
type MutualTLSMode int

const (
	// MTLSUnknown is used to indicate the variable hasn't been initialized correctly (with the authentication policy).
	MTLSUnknown MutualTLSMode = iota

	// MTLSDisable if authentication policy disable mTLS.
	MTLSDisable

	// MTLSPermissive if authentication policy enable mTLS in permissive mode.
	MTLSPermissive

	// MTLSStrict if authentication policy enable mTLS in strict mode.
	MTLSStrict
)

// String converts MutualTLSMode to human readable string for debugging.
func (mode MutualTLSMode) String() string {
	switch mode {
	case MTLSDisable:
		return "DISABLE"
	case MTLSPermissive:
		return "PERMISSIVE"
	case MTLSStrict:
		return "STRICT"
	default:
		return "UNKNOWN"
	}
}

// AuthenticationPolicies organizes authentication (mTLS + JWT) policies by namespace.
type AuthenticationPolicies struct {
	// Maps from namespace to the v1beta1 authentication policies.
	requestAuthentications map[string][]Config

	peerAuthentications map[string][]Config

	rootNamespace string
}

// initAuthenticationPolicies creates a new AuthenticationPolicies struct and populates with the
// authentication policies in the mesh environment.
func initAuthenticationPolicies(env *Environment) (*AuthenticationPolicies, error) {
	policy := &AuthenticationPolicies{
		requestAuthentications: map[string][]Config{},
		peerAuthentications:    map[string][]Config{},
		rootNamespace:          env.Mesh().GetRootNamespace(),
	}

	if configs, err := env.List(
		collections.IstioSecurityV1Beta1Requestauthentications.Resource().GroupVersionKind(), NamespaceAll); err == nil {
		sortConfigByCreationTime(configs)
		policy.addRequestAuthentication(configs)
	} else {
		return nil, err
	}

	if configs, err := env.List(
		collections.IstioSecurityV1Beta1Peerauthentications.Resource().GroupVersionKind(), NamespaceAll); err == nil {
		sortConfigByCreationTime(configs)
		policy.addPeerAuthentication(configs)
	} else {
		return nil, err
	}

	return policy, nil
}

func (policy *AuthenticationPolicies) addRequestAuthentication(configs []Config) {
	for _, config := range configs {
		policy.requestAuthentications[config.Namespace] =
			append(policy.requestAuthentications[config.Namespace], config)
	}
}

func (policy *AuthenticationPolicies) addPeerAuthentication(configs []Config) {
	for _, config := range configs {
		policy.peerAuthentications[config.Namespace] =
			append(policy.peerAuthentications[config.Namespace], config)
	}
}

// GetJwtPoliciesForWorkload returns a list of JWT policies matching to labels.
func (policy *AuthenticationPolicies) GetJwtPoliciesForWorkload(namespace string,
	workloadLabels labels.Collection) []*Config {
	return getConfigsForWorkload(policy.requestAuthentications, policy.rootNamespace, namespace, workloadLabels)
}

// GetPeerAuthenticationsForWorkload returns a list of peer authentication policies matching to labels.
func (policy *AuthenticationPolicies) GetPeerAuthenticationsForWorkload(namespace string,
	workloadLabels labels.Collection) []*Config {
	return getConfigsForWorkload(policy.peerAuthentications, policy.rootNamespace, namespace, workloadLabels)
}

// GetRootNamespace return root namespace that is tracked by the policy object.
func (policy *AuthenticationPolicies) GetRootNamespace() string {
	return policy.rootNamespace
}

func getConfigsForWorkload(configsByNamespace map[string][]Config,
	rootNamespace string,
	namespace string,
	workloadLabels labels.Collection) []*Config {
	configs := make([]*Config, 0)
	lookupInNamespaces := []string{namespace}
	if namespace != rootNamespace {
		// Only check the root namespace if the (workload) namespace is not already the root namespace
		// to avoid double inclusion.
		lookupInNamespaces = append(lookupInNamespaces, rootNamespace)
	}
	for _, ns := range lookupInNamespaces {
		if nsConfig, ok := configsByNamespace[ns]; ok {
			for idx := range nsConfig {
				cfg := &nsConfig[idx]
				if ns != cfg.Namespace {
					// Should never come here. Log warning just in case.
					log.Warnf("Seeing config %s with namespace %s in map entry for %s. Ignored", cfg.Name, cfg.Namespace, ns)
					continue
				}
				var selector labels.Instance
				switch cfg.Type {
				case collections.IstioSecurityV1Beta1Requestauthentications.Resource().Kind():
					selector = labels.Instance(cfg.Spec.(*v1beta1.RequestAuthentication).GetSelector().GetMatchLabels())
				case collections.IstioSecurityV1Beta1Peerauthentications.Resource().Kind():
					selector = labels.Instance(cfg.Spec.(*v1beta1.PeerAuthentication).GetSelector().GetMatchLabels())
				default:
					log.Warnf("Not support authentication type %q", cfg.Type)
					continue
				}
				if workloadLabels.IsSupersetOf(selector) {
					configs = append(configs, cfg)
				}
			}
		}
	}

	return configs
}
