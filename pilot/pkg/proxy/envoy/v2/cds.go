// Copyright 2018 Istio Authors
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

package v2

import (
	"fmt"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/gogo/protobuf/types"

	"istio.io/istio/pilot/pkg/model"
	authn_model "istio.io/istio/pilot/pkg/security/model"
)

// clusters aggregate a DiscoveryResponse for pushing.
func (conn *XdsConnection) clusters(response []*xdsapi.Cluster) *xdsapi.DiscoveryResponse {
	out := &xdsapi.DiscoveryResponse{
		// All resources for CDS ought to be of the type ClusterLoadAssignment
		TypeUrl: ClusterType,

		// Pilot does not really care for versioning. It always supplies what's currently
		// available to it, irrespective of whether Envoy chooses to accept or reject CDS
		// responses. Pilot believes in eventual consistency and that at some point, Envoy
		// will begin seeing results it deems to be good.
		VersionInfo: versionInfo(),
		Nonce:       nonce(),
	}

	for _, c := range response {
		cc, _ := types.MarshalAny(c)
		out.Resources = append(out.Resources, *cc)
	}

	return out
}

func (s *DiscoveryServer) pushCds(con *XdsConnection, push *model.PushContext, version string) error {
	// TODO: Modify interface to take services, and config instead of making library query registry
	rawClusters, err := s.generateRawClusters(con.modelNode, push)
	if err != nil {
		return err
	}
	if s.DebugConfigs {
		con.CDSClusters = rawClusters
	}
	response := con.clusters(rawClusters)
	err = con.send(response)
	if err != nil {
		adsLog.Warnf("CDS: Send failure %s: %v", con.ConID, err)
		cdsSendErrPushes.Increment()
		return err
	}
	cdsPushes.Increment()

	// The response can't be easily read due to 'any' marshaling.
	adsLog.Infof("CDS: PUSH for node:%s clusters:%d services:%d version:%s",
		con.modelNode.ID, len(rawClusters), len(push.Services(nil)), version)
	return nil
}

func (s *DiscoveryServer) generateRawClusters(node *model.Proxy, push *model.PushContext) ([]*xdsapi.Cluster, error) {
	rawClusters, err := s.ConfigGenerator.BuildClusters(s.Env, node, push)
	if err != nil {
		adsLog.Warnf("CDS: Failed to generate clusters for node:%s: %v", node.ID, err)
		cdsBuildErrPushes.Increment()
		return nil, err
	}

	for _, c := range rawClusters {
		if err = c.Validate(); err != nil {
			retErr := fmt.Errorf("CDS: Generated invalid cluster for node %v: %v", node, err)
			adsLog.Errorf("CDS: Generated invalid cluster for node:%s: %v, %v", node.ID, err, c)
			cdsBuildErrPushes.Increment()
			totalXDSInternalErrors.Increment()
			// Generating invalid clusters is a bug.
			// Panic instead of trying to recover from that, since we can't
			// assume anything about the state.
			panic(retErr.Error())
		}
	}
	return rawClusters, nil
}

// Set the token path for SDS if SDS_TOKEN_PATH is defined in the proxy metadata
func SetTokenPathForSdsFromProxyMetadata(c *xdsapi.Cluster, node *model.Proxy) {
	if sdsTokenPath, found := node.Metadata[model.NodeMetadataSdsTokenPath]; found && len(sdsTokenPath) > 0 {
		// Set the SDS token path in the TLS certificate config
		if c.GetTlsContext() != nil && c.GetTlsContext().GetCommonTlsContext() != nil &&
			c.GetTlsContext().GetCommonTlsContext().GetTlsCertificateSdsSecretConfigs() != nil {
			for _, sc := range c.GetTlsContext().GetCommonTlsContext().GetTlsCertificateSdsSecretConfigs() {
				if sc.GetSdsConfig() != nil && sc.GetSdsConfig().GetApiConfigSource() != nil &&
					sc.GetSdsConfig().GetApiConfigSource().GetGrpcServices() != nil {
					for _, svc := range sc.GetSdsConfig().GetApiConfigSource().GetGrpcServices() {
						// If no call-credential in the cluster, no need to set SDS token path
						if svc.GetGoogleGrpc() != nil && svc.GetGoogleGrpc().GetCallCredentials() != nil &&
							svc.GetGoogleGrpc().GetCredentialsFactoryName() == authn_model.FileBasedMetadataPlugName {
							adsLog.Debugf("Set SDS token path in TLS context based on the proxy metadata")
							svc.GetGoogleGrpc().CallCredentials =
								authn_model.ConstructgRPCCallCredentials(sdsTokenPath, authn_model.K8sSAJwtTokenHeaderKey)
						}
					}
				}
			}
		}

		// Set the SDS token path in the TLS validation context
		if c.GetTlsContext() != nil && c.GetTlsContext().GetCommonTlsContext() != nil &&
			c.GetTlsContext().GetCommonTlsContext().GetCombinedValidationContext() != nil &&
			c.GetTlsContext().GetCommonTlsContext().GetCombinedValidationContext().GetValidationContextSdsSecretConfig() != nil {
			sc := c.GetTlsContext().GetCommonTlsContext().GetCombinedValidationContext().GetValidationContextSdsSecretConfig()
			if sc.GetSdsConfig() != nil && sc.GetSdsConfig().GetApiConfigSource() != nil &&
				sc.GetSdsConfig().GetApiConfigSource().GetGrpcServices() != nil {
				for _, svc := range sc.GetSdsConfig().GetApiConfigSource().GetGrpcServices() {
					// If no call-credential in the cluster, no need to set SDS token path
					if svc.GetGoogleGrpc() != nil && svc.GetGoogleGrpc().GetCallCredentials() != nil &&
						svc.GetGoogleGrpc().GetCredentialsFactoryName() == authn_model.FileBasedMetadataPlugName {
						adsLog.Debugf("Set SDS token path in validation context based on the proxy metadata")
						svc.GetGoogleGrpc().CallCredentials =
							authn_model.ConstructgRPCCallCredentials(sdsTokenPath, authn_model.K8sSAJwtTokenHeaderKey)
					}
				}
			}
		}
	}
}
