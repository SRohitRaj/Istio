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

// Package sds implements secret discovery service in NodeAgent.
package sds

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	tls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	sds "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/istio/pilot/pkg/util/sets"
	"istio.io/istio/pilot/pkg/xds"
	v3 "istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pkg/config/schema/gvk"
	"istio.io/istio/pkg/security"
	nodeagentutil "istio.io/istio/security/pkg/nodeagent/util"
	"istio.io/pkg/log"
)

var sdsServiceLog = log.RegisterScope("sds", "SDS service debugging", 0)

type sdsservice struct {
	st security.SecretManager

	XdsServer  *xds.DiscoveryServer
	stop       chan struct{}
	rootCaPath string
}

// Assert we implement the generator interface
var _ model.XdsResourceGenerator = &sdsservice{}

func NewXdsServer(stop chan struct{}, gen model.XdsResourceGenerator) *xds.DiscoveryServer {
	s := xds.NewXDS(stop)
	s.DiscoveryServer.Generators = map[string]model.XdsResourceGenerator{
		v3.SecretType: gen,
	}
	s.DiscoveryServer.ProxyNeedsPush = func(proxy *model.Proxy, req *model.PushRequest) bool {
		// Empty changes means "all"
		if len(req.ConfigsUpdated) == 0 {
			return true
		}
		proxy.RLock()
		wr := proxy.WatchedResources[v3.SecretType]
		proxy.RUnlock()

		if wr == nil {
			return false
		}

		names := sets.NewSet(wr.ResourceNames...)
		found := false
		for name := range model.ConfigsOfKind(req.ConfigsUpdated, gvk.Secret) {
			if names.Contains(name.Name) {
				found = true
				break
			}
		}
		return found
	}
	s.DiscoveryServer.Start(stop)
	return s.DiscoveryServer
}

// newSDSService creates Secret Discovery Service which implements envoy SDS API.
func newSDSService(st security.SecretManager, options *security.Options) *sdsservice {
	ret := &sdsservice{
		st:   st,
		stop: make(chan struct{}),
	}
	ret.XdsServer = NewXdsServer(ret.stop, ret)

	ret.rootCaPath = options.CARootPath

	if options.FileMountedCerts {
		return ret
	}

	// Pre-generate workload certificates to improve startup latency and ensure that for OUTPUT_CERTS
	// case we always write a certificate. A workload can technically run without any mTLS/CA
	// configured, in which case this will fail; if it becomes noisy we should disable the entire SDS
	// server in these cases.
	go func() {
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 0
		for {
			_, err := st.GenerateSecret(security.WorkloadKeyCertResourceName, options.CARootPath)
			if err == nil {
				break
			}
			sdsServiceLog.Warnf("failed to warm certificate: %v", err)
			select {
			case <-ret.stop:
				return
			case <-time.After(b.NextBackOff()):
			}
		}
		for {
			_, err := st.GenerateSecret(security.RootCertReqResourceName, options.CARootPath)
			if err == nil {
				break
			}
			sdsServiceLog.Warnf("failed to warm root certificate: %v", err)
			select {
			case <-ret.stop:
				return
			case <-time.After(b.NextBackOff()):
			}
		}
	}()

	return ret
}

func (s *sdsservice) generate(resourceNames []string) (model.Resources, error) {
	resources := model.Resources{}
	for _, resourceName := range resourceNames {
		secret, err := s.st.GenerateSecret(resourceName, s.rootCaPath)
		if err != nil {
			// Typically, in Istiod, we do not return an error for a failure to generate a resource
			// However, here it makes sense, because we are generally streaming a single resource,
			// so sending an error will not cause a single failure to prevent the entire multiplex stream
			// of resources, and failures here are generally due to temporary networking issues to the CA
			// rather than a result of configuration issues, which trigger updates in Istiod when resolved.
			// Instead, we rely on the client to retry (with backoff) on failures.
			return nil, fmt.Errorf("failed to generate secret for %v: %v", resourceName, err)
		}

		res := util.MessageToAny(toEnvoySecret(secret, s.rootCaPath))
		resources = append(resources, &discovery.Resource{
			Name:     resourceName,
			Resource: res,
		})
	}
	return resources, nil
}

// Generate implements the XDS Generator interface. This allows the XDS server to dispatch requests
// for SecretTypeV3 to our server to generate the Envoy response.
func (s *sdsservice) Generate(_ *model.Proxy, _ *model.PushContext, w *model.WatchedResource,
	updates *model.PushRequest) (model.Resources, model.XdsLogDetails, error) {
	// updates.Full indicates we should do a complete push of all updated resources
	// In practice, all pushes should be incremental (ie, if the `default` cert changes we won't push
	// all file certs).
	if updates.Full {
		resp, err := s.generate(w.ResourceNames)
		return resp, pushLog(w.ResourceNames), err
	}
	names := []string{}
	watched := sets.NewSet(w.ResourceNames...)
	for i := range updates.ConfigsUpdated {
		if i.Kind == gvk.Secret && watched.Contains(i.Name) {
			names = append(names, i.Name)
		}
	}
	resp, err := s.generate(names)
	return resp, pushLog(names), err
}

// register adds the SDS handle to the grpc server
func (s *sdsservice) register(rpcs *grpc.Server) {
	sds.RegisterSecretDiscoveryServiceServer(rpcs, s)
}

// StreamSecrets serves SDS discovery requests and SDS push requests
func (s *sdsservice) StreamSecrets(stream sds.SecretDiscoveryService_StreamSecretsServer) error {
	return s.XdsServer.Stream(stream)
}

func (s *sdsservice) DeltaSecrets(stream sds.SecretDiscoveryService_DeltaSecretsServer) error {
	return status.Error(codes.Unimplemented, "DeltaSecrets not implemented")
}

func (s *sdsservice) FetchSecrets(ctx context.Context, discReq *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "FetchSecrets not implemented")
}

func (s *sdsservice) Close() {
	close(s.stop)
	s.XdsServer.Shutdown()
}

// toEnvoySecret converts a security.SecretItem to an Envoy tls.Secret
func toEnvoySecret(s *security.SecretItem, caRootPath string) *tls.Secret {
	secret := &tls.Secret{
		Name: s.ResourceName,
	}
	cfg := nodeagentutil.SdsCertificateConfig{}
	ok := false
	if caRootPath != "" {
		cfg, ok = nodeagentutil.SdsCertificateConfigFromResourceName(caRootPath)
	} else {
		cfg, ok = nodeagentutil.SdsCertificateConfigFromResourceName(s.ResourceName)
	}
	if s.ResourceName == security.RootCertReqResourceName || (ok && cfg.IsRootCertificate()) {
		secret.Type = &tls.Secret_ValidationContext{
			ValidationContext: &tls.CertificateValidationContext{
				TrustedCa: &core.DataSource{
					Specifier: &core.DataSource_InlineBytes{
						InlineBytes: s.RootCert,
					},
				},
			},
		}
	} else {
		secret.Type = &tls.Secret_TlsCertificate{
			TlsCertificate: &tls.TlsCertificate{
				CertificateChain: &core.DataSource{
					Specifier: &core.DataSource_InlineBytes{
						InlineBytes: s.CertificateChain,
					},
				},
				PrivateKey: &core.DataSource{
					Specifier: &core.DataSource_InlineBytes{
						InlineBytes: s.PrivateKey,
					},
				},
			},
		}
	}

	return secret
}

func pushLog(names []string) model.XdsLogDetails {
	if len(names) == 1 {
		// For common case of single resource, show which resource it was
		return model.XdsLogDetails{AdditionalInfo: "resource:" + names[0]}
	}
	return model.DefaultXdsLogDetails
}
