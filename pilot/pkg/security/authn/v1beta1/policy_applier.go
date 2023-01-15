// Copyright Istio Authors
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

package v1beta1

import (
	"fmt"
	"sort"
	"strings"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_jwt "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"

	authn_alpha "istio.io/api/authentication/v1alpha1"
	authn_filter "istio.io/api/envoy/config/filter/http/authn/v2alpha1"
	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/api/security/v1beta1"
	"istio.io/istio/pilot/pkg/extensionproviders"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking"
	"istio.io/istio/pilot/pkg/security/authn"
	authn_utils "istio.io/istio/pilot/pkg/security/authn/utils"
	authn_model "istio.io/istio/pilot/pkg/security/model"
	"istio.io/istio/pilot/pkg/util/protoconv"
	"istio.io/istio/pilot/pkg/xds/filters"
	"istio.io/istio/pkg/config"
	"istio.io/istio/pkg/config/security"
	"istio.io/istio/pkg/jwt"
	"istio.io/pkg/log"
)

var authnLog = log.RegisterScope("authn", "authn debugging", 0)

// Implementation of authn.PolicyApplier with v1beta1 API.
type v1beta1PolicyApplier struct {
	jwtPolicies []*config.Config

	peerPolices []*config.Config

	// processedJwtRules is the consolidate JWT rules from all jwtPolicies.
	processedJwtRules []*v1beta1.JWTRule

	consolidatedPeerPolicy *v1beta1.PeerAuthentication

	push *model.PushContext
}

func (a *v1beta1PolicyApplier) JwtFilter() *hcm.HttpFilter {
	if len(a.processedJwtRules) == 0 {
		return nil
	}

	filterConfigProto := convertToEnvoyJwtConfig(a.processedJwtRules, a.push)

	if filterConfigProto == nil {
		return nil
	}
	return &hcm.HttpFilter{
		Name:       authn_model.EnvoyJwtFilterName,
		ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: protoconv.MessageToAny(filterConfigProto)},
	}
}

func defaultAuthnFilter() *authn_filter.FilterConfig {
	return &authn_filter.FilterConfig{
		Policy: &authn_alpha.Policy{},
		// we can always set this field, it's no-op if mTLS is not used.
		SkipValidateTrustDomain: true,
	}
}

func (a *v1beta1PolicyApplier) setAuthnFilterForRequestAuthn(config *authn_filter.FilterConfig) *authn_filter.FilterConfig {
	if len(a.processedJwtRules) == 0 {
		// (beta) RequestAuthentication is not set for workload, do nothing.
		authnLog.Debug("AuthnFilter: RequestAuthentication (beta policy) not found, keep settings with alpha API")
		return config
	}

	if config == nil {
		config = defaultAuthnFilter()
	}

	// This is obsoleted and not needed (payload is extracted from metadata). Reset the field to remove
	// any artifacts from alpha applier.
	config.JwtOutputPayloadLocations = nil
	p := config.Policy
	// Reset origins to use with beta API
	// nolint: staticcheck
	p.Origins = []*authn_alpha.OriginAuthenticationMethod{}
	// Always set to true for beta API, as it doesn't doe rejection on missing token.
	// nolint: staticcheck
	p.OriginIsOptional = true

	// Always bind request.auth.principal from JWT origin. In v2 policy, authorization config specifies what principal to
	// choose from instead, rather than in authn config.
	// nolint: staticcheck
	p.PrincipalBinding = authn_alpha.PrincipalBinding_USE_ORIGIN
	// nolint: staticcheck
	for _, jwt := range a.processedJwtRules {
		p.Origins = append(p.Origins, &authn_alpha.OriginAuthenticationMethod{
			Jwt: &authn_alpha.Jwt{
				// used for getting the filter data, and all other fields are irrelevant.
				Issuer: jwt.GetIssuer(),
			},
		})
	}
	return config
}

// AuthNFilter returns the Istio authn filter config:
// - If RequestAuthentication is used, it overwrite the settings for request principal validation and extraction based on the new API.
// - If RequestAuthentication is used, principal binding is always set to ORIGIN.
func (a *v1beta1PolicyApplier) AuthNFilter(forSidecar bool) *hcm.HttpFilter {
	var filterConfigProto *authn_filter.FilterConfig

	// Override the config with request authentication, if applicable.
	filterConfigProto = a.setAuthnFilterForRequestAuthn(filterConfigProto)

	if filterConfigProto == nil {
		return nil
	}
	// disable clear route cache for sidecars because the JWT claim based routing is only supported on gateways.
	filterConfigProto.DisableClearRouteCache = forSidecar

	// Note: in previous Istio versions, the authn filter also handled PeerAuthentication, to extract principal.
	// This has been modified to rely on the TCP filter

	return &hcm.HttpFilter{
		Name:       filters.AuthnFilterName,
		ConfigType: &hcm.HttpFilter_TypedConfig{TypedConfig: protoconv.MessageToAny(filterConfigProto)},
	}
}

func (a *v1beta1PolicyApplier) InboundMTLSSettings(
	endpointPort uint32,
	node *model.Proxy,
	trustDomainAliases []string,
	modeOverride model.MutualTLSMode,
) authn.MTLSSettings {
	effectiveMTLSMode := modeOverride
	if effectiveMTLSMode == model.MTLSUnknown {
		effectiveMTLSMode = a.GetMutualTLSModeForPort(endpointPort)
	}
	authnLog.Debugf("InboundFilterChain: build inbound filter change for %v:%d in %s mode", node.ID, endpointPort, effectiveMTLSMode)
	var mc *meshconfig.MeshConfig
	if a.push != nil {
		mc = a.push.Mesh
	}
	// Configure TLS version based on meshconfig TLS API.
	minTLSVersion := authn_utils.GetMinTLSVersion(mc.GetMeshMTLS().GetMinProtocolVersion())
	return authn.MTLSSettings{
		Port: endpointPort,
		Mode: effectiveMTLSMode,
		TCP: authn_utils.BuildInboundTLS(effectiveMTLSMode, node, networking.ListenerProtocolTCP,
			trustDomainAliases, minTLSVersion),
		HTTP: authn_utils.BuildInboundTLS(effectiveMTLSMode, node, networking.ListenerProtocolHTTP,
			trustDomainAliases, minTLSVersion),
	}
}

// NewPolicyApplier returns new applier for v1beta1 authentication policies.
func NewPolicyApplier(rootNamespace string,
	jwtPolicies []*config.Config,
	peerPolicies []*config.Config,
	push *model.PushContext,
) authn.PolicyApplier {
	processedJwtRules := []*v1beta1.JWTRule{}

	// TODO(diemtvu) should we need to deduplicate JWT with the same issuer.
	// https://github.com/istio/istio/issues/19245
	for idx := range jwtPolicies {
		spec := jwtPolicies[idx].Spec.(*v1beta1.RequestAuthentication)
		processedJwtRules = append(processedJwtRules, spec.JwtRules...)
	}

	// Sort the jwt rules by the issuer alphabetically to make the later-on generated filter
	// config deterministic.
	sort.Slice(processedJwtRules, func(i, j int) bool {
		return strings.Compare(
			processedJwtRules[i].GetIssuer(), processedJwtRules[j].GetIssuer()) < 0
	})

	return &v1beta1PolicyApplier{
		jwtPolicies:            jwtPolicies,
		peerPolices:            peerPolicies,
		processedJwtRules:      processedJwtRules,
		consolidatedPeerPolicy: ComposePeerAuthentication(rootNamespace, peerPolicies),
		push:                   push,
	}
}

// convertToEnvoyJwtConfig converts a list of JWT rules into Envoy JWT filter config to enforce it.
// Each rule is expected corresponding to one JWT issuer (provider).
// The behavior of the filter should reject all requests with invalid token. On the other hand,
// if no token provided, the request is allowed.
func convertToEnvoyJwtConfig(jwtRules []*v1beta1.JWTRule, push *model.PushContext) *envoy_jwt.JwtAuthentication {
	if len(jwtRules) == 0 {
		return nil
	}

	providers := map[string]*envoy_jwt.JwtProvider{}
	// Each element of innerAndList is the requirement for each provider, in the form of
	// {provider OR `allow_missing`}
	// This list will be ANDed (if have more than one provider) for the final requirement.
	innerAndList := []*envoy_jwt.JwtRequirement{}

	// This is an (or) list for all providers. This will be OR with the innerAndList above so
	// it can pass the requirement in the case that providers share the same location.
	outterOrList := []*envoy_jwt.JwtRequirement{}

	for i, jwtRule := range jwtRules {
		provider := &envoy_jwt.JwtProvider{
			Issuer:               jwtRule.Issuer,
			Audiences:            jwtRule.Audiences,
			Forward:              jwtRule.ForwardOriginalToken,
			ForwardPayloadHeader: jwtRule.OutputPayloadToHeader,
			PayloadInMetadata:    jwtRule.Issuer,
		}

		for _, claimAndHeader := range jwtRule.OutputClaimToHeaders {
			provider.ClaimToHeaders = append(provider.ClaimToHeaders, &envoy_jwt.JwtClaimToHeader{
				HeaderName: claimAndHeader.Header,
				ClaimName:  claimAndHeader.Claim,
			})
		}

		for _, location := range jwtRule.FromHeaders {
			provider.FromHeaders = append(provider.FromHeaders, &envoy_jwt.JwtHeader{
				Name:        location.Name,
				ValuePrefix: location.Prefix,
			})
		}
		provider.FromParams = jwtRule.FromParams

		authnLog.Debugf("JwksFetchMode is set to: %v", features.JwksFetchMode)

		// Use remote jwks if jwksUri is not empty and JwksFetchMode not Istiod. Parse the jwksUri to get the
		// cluster name, generate the jwt filter config using remote Jwks.
		// If failed to parse the cluster name, only fallback to let istiod to fetch the jwksUri when
		// remoteJwksMode is Hybrid.
		if features.JwksFetchMode != jwt.Istiod && jwtRule.JwksUri != "" {
			jwksInfo, err := security.ParseJwksURI(jwtRule.JwksUri)
			if err != nil {
				authnLog.Errorf("Failed to parse jwt rule jwks uri %v", err)
			}
			_, cluster, err := extensionproviders.LookupCluster(push, jwksInfo.Hostname.String(), jwksInfo.Port)
			authnLog.Debugf("Look up cluster result: %v", cluster)

			if err == nil && len(cluster) > 0 {
				// This is a case of URI pointing to mesh cluster. Setup Remote Jwks and let Envoy fetch the key.
				provider.JwksSourceSpecifier = &envoy_jwt.JwtProvider_RemoteJwks{
					RemoteJwks: &envoy_jwt.RemoteJwks{
						HttpUri: &core.HttpUri{
							Uri: jwtRule.JwksUri,
							HttpUpstreamType: &core.HttpUri_Cluster{
								Cluster: cluster,
							},
							Timeout: &durationpb.Duration{Seconds: 5},
						},
						CacheDuration: &durationpb.Duration{Seconds: 5 * 60},
					},
				}
			} else if features.JwksFetchMode == jwt.Hybrid {
				provider.JwksSourceSpecifier = push.JwtKeyResolver.BuildLocalJwks(jwtRule.JwksUri, jwtRule.Issuer, "")
			} else {
				// Log error and create fake inline JWKs
				authnLog.Errorf("Failed to look up Envoy cluster %v. "+
					"Please create ServiceEntry to register external JWKs server or "+
					"set PILOT_JWT_ENABLE_REMOTE_JWKS to hybrid/istiod mode.", err)
				provider.JwksSourceSpecifier = push.JwtKeyResolver.BuildLocalJwks("", "", model.CreateFakeJwks(jwtRule.JwksUri))
			}
		} else {
			// Use inline jwks as existing flow, either jwtRule.jwks is empty or let istiod to fetch the jwtRule.jwksUri
			provider.JwksSourceSpecifier = push.JwtKeyResolver.BuildLocalJwks(jwtRule.JwksUri, jwtRule.Issuer, jwtRule.Jwks)
		}

		name := fmt.Sprintf("origins-%d", i)
		providers[name] = provider
		innerAndList = append(innerAndList, &envoy_jwt.JwtRequirement{
			RequiresType: &envoy_jwt.JwtRequirement_RequiresAny{
				RequiresAny: &envoy_jwt.JwtRequirementOrList{
					Requirements: []*envoy_jwt.JwtRequirement{
						{
							RequiresType: &envoy_jwt.JwtRequirement_ProviderName{
								ProviderName: name,
							},
						},
						{
							RequiresType: &envoy_jwt.JwtRequirement_AllowMissing{
								AllowMissing: &emptypb.Empty{},
							},
						},
					},
				},
			},
		})
		outterOrList = append(outterOrList, &envoy_jwt.JwtRequirement{
			RequiresType: &envoy_jwt.JwtRequirement_ProviderName{
				ProviderName: name,
			},
		})
	}

	// If there is only one provider, simply use an OR of {provider, `allow_missing`}.
	if len(innerAndList) == 1 {
		return &envoy_jwt.JwtAuthentication{
			Rules: []*envoy_jwt.RequirementRule{
				{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: "/",
						},
					},
					RequirementType: &envoy_jwt.RequirementRule_Requires{
						Requires: innerAndList[0],
					},
				},
			},
			Providers:           providers,
			BypassCorsPreflight: true,
		}
	}

	// If there are more than one provider, filter should OR of
	// {P1, P2 .., AND of {OR{P1, allow_missing}, OR{P2, allow_missing} ...}}
	// where the innerAnd enforce a token, if provided, must be valid, and the
	// outer OR aids the case where providers share the same location (as
	// it will always fail with the innerAND).
	outterOrList = append(outterOrList, &envoy_jwt.JwtRequirement{
		RequiresType: &envoy_jwt.JwtRequirement_RequiresAll{
			RequiresAll: &envoy_jwt.JwtRequirementAndList{
				Requirements: innerAndList,
			},
		},
	})

	return &envoy_jwt.JwtAuthentication{
		Rules: []*envoy_jwt.RequirementRule{
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				RequirementType: &envoy_jwt.RequirementRule_Requires{
					Requires: &envoy_jwt.JwtRequirement{
						RequiresType: &envoy_jwt.JwtRequirement_RequiresAny{
							RequiresAny: &envoy_jwt.JwtRequirementOrList{
								Requirements: outterOrList,
							},
						},
					},
				},
			},
		},
		Providers:           providers,
		BypassCorsPreflight: true,
	}
}

func (a *v1beta1PolicyApplier) PortLevelSetting() map[uint32]*v1beta1.PeerAuthentication_MutualTLS {
	return a.consolidatedPeerPolicy.PortLevelMtls
}

func (a *v1beta1PolicyApplier) GetMutualTLSModeForPort(endpointPort uint32) model.MutualTLSMode {
	if a.consolidatedPeerPolicy.PortLevelMtls != nil {
		if portMtls, ok := a.consolidatedPeerPolicy.PortLevelMtls[endpointPort]; ok {
			return getMutualTLSMode(portMtls)
		}
	}

	return getMutualTLSMode(a.consolidatedPeerPolicy.Mtls)
}

// getMutualTLSMode returns the MutualTLSMode enum corresponding peer MutualTLS settings.
// Input cannot be nil.
func getMutualTLSMode(mtls *v1beta1.PeerAuthentication_MutualTLS) model.MutualTLSMode {
	return model.ConvertToMutualTLSMode(mtls.Mode)
}

// ComposePeerAuthentication returns the effective PeerAuthentication given the list of applicable
// configs. This list should contains at most 1 mesh-level and 1 namespace-level configs.
// Workload-level configs should not be in root namespace (this should be guaranteed by the caller,
// though they will be safely ignored in this function). If the input config list is empty, returns
// a default policy set to a PERMISSIVE.
// If there is at least one applicable config, returns should not be nil, and is a combined policy
// based on following rules:
// - It should have the setting from the most narrow scope (i.e workload-level is preferred over
// namespace-level, which is preferred over mesh-level).
// - When there are more than one policy in the same scope (i.e workload-level), the oldest one win.
// - UNSET will be replaced with the setting from the parent. I.e UNSET port-level config will be
// replaced with config from workload-level, UNSET in workload-level config will be replaced with
// one in namespace-level and so on.
func ComposePeerAuthentication(rootNamespace string, configs []*config.Config) *v1beta1.PeerAuthentication {
	var meshCfg, namespaceCfg, workloadCfg *config.Config

	// Initial outputPolicy is set to a PERMISSIVE.
	outputPolicy := v1beta1.PeerAuthentication{
		Mtls: &v1beta1.PeerAuthentication_MutualTLS{
			Mode: v1beta1.PeerAuthentication_MutualTLS_PERMISSIVE,
		},
	}

	for _, cfg := range configs {
		spec := cfg.Spec.(*v1beta1.PeerAuthentication)
		if spec.Selector == nil || len(spec.Selector.MatchLabels) == 0 {
			// Namespace-level or mesh-level policy
			if cfg.Namespace == rootNamespace {
				if meshCfg == nil || cfg.CreationTimestamp.Before(meshCfg.CreationTimestamp) {
					authnLog.Debugf("Switch selected mesh policy to %s.%s (%v)", cfg.Name, cfg.Namespace, cfg.CreationTimestamp)
					meshCfg = cfg
				}
			} else {
				if namespaceCfg == nil || cfg.CreationTimestamp.Before(namespaceCfg.CreationTimestamp) {
					authnLog.Debugf("Switch selected namespace policy to %s.%s (%v)", cfg.Name, cfg.Namespace, cfg.CreationTimestamp)
					namespaceCfg = cfg
				}
			}
		} else if cfg.Namespace != rootNamespace {
			// Workload-level policy, aka the one with selector and not in root namespace.
			if workloadCfg == nil || cfg.CreationTimestamp.Before(workloadCfg.CreationTimestamp) {
				authnLog.Debugf("Switch selected workload policy to %s.%s (%v)", cfg.Name, cfg.Namespace, cfg.CreationTimestamp)
				workloadCfg = cfg
			}
		}
	}

	// Process in mesh, namespace, workload order to resolve inheritance (UNSET)

	if meshCfg != nil && !isMtlsModeUnset(meshCfg.Spec.(*v1beta1.PeerAuthentication).Mtls) {
		// If mesh policy is defined, update parent policy to mesh policy.
		outputPolicy.Mtls = meshCfg.Spec.(*v1beta1.PeerAuthentication).Mtls
	}

	if namespaceCfg != nil && !isMtlsModeUnset(namespaceCfg.Spec.(*v1beta1.PeerAuthentication).Mtls) {
		// If namespace policy is defined, update output policy to namespace policy. This means namespace
		// policy overwrite mesh policy.
		outputPolicy.Mtls = namespaceCfg.Spec.(*v1beta1.PeerAuthentication).Mtls
	}

	var workloadPolicy *v1beta1.PeerAuthentication
	if workloadCfg != nil {
		workloadPolicy = workloadCfg.Spec.(*v1beta1.PeerAuthentication)
	}

	if workloadPolicy != nil && !isMtlsModeUnset(workloadPolicy.Mtls) {
		// If workload policy is defined, update parent policy to workload policy.
		outputPolicy.Mtls = workloadPolicy.Mtls
	}

	if workloadPolicy != nil && workloadPolicy.PortLevelMtls != nil {
		outputPolicy.PortLevelMtls = make(map[uint32]*v1beta1.PeerAuthentication_MutualTLS, len(workloadPolicy.PortLevelMtls))
		for port, mtls := range workloadPolicy.PortLevelMtls {
			if isMtlsModeUnset(mtls) {
				// Inherit from workload level.
				outputPolicy.PortLevelMtls[port] = outputPolicy.Mtls
			} else {
				outputPolicy.PortLevelMtls[port] = mtls
			}
		}
	}

	return &outputPolicy
}

func isMtlsModeUnset(mtls *v1beta1.PeerAuthentication_MutualTLS) bool {
	return mtls == nil || mtls.Mode == v1beta1.PeerAuthentication_MutualTLS_UNSET
}
