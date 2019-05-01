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

package vm

import (
	"fmt"

	"istio.io/istio/security/pkg/platform"
	"istio.io/istio/security/pkg/util"

	"istio.io/istio/security/pkg/caclient/protocol"
	"istio.io/istio/security/pkg/nodeagent/caclient"
)

// NodeAgent interface that should be implemented by
// various platform specific node agents.
type NodeAgent interface {
	Start() error
}

// NewNodeAgent is constructor for Node agent based on the provided Environment variable.
func NewNodeAgent(cfg *Config) (NodeAgent, error) {
	fmt.Println("new Node Agent")
	if cfg == nil {
		return nil, fmt.Errorf("nil configuration passed")
	}
	switch cfg.CAClientConfig.CAProtocol {
	case IstioCertificateService:
		return initCitadelAgentForIstioCertificateService(cfg)
	case IstioCAService:
		return initCitadelAgentForIstioCAService(cfg)
	default:
		return nil, fmt.Errorf("unrecognized CA protocol: %s. "+
			"Need to be either IstioCertificateService or IstioCAService", cfg.CAClientConfig.CAProtocol)
	}
}

// initCitadelAgentForIstioCertificateService initializes nodeAgentIntervalV2 implementation of NodeAgent interface that
// talks to IstioCertificateService gRPC service defined at https://github.com/istio/istio/blob/master/security/proto/istioca.proto
// It only supports talking to GoogleCA for now.
func initCitadelAgentForIstioCertificateService(cfg *Config) (*citadelAgent, error) {
	na := &citadelAgent{
		config:   cfg,
		certUtil: util.NewCertUtil(cfg.CAClientConfig.CSRGracePeriodPercentage),
	}

	pc, err := platform.NewClient(cfg.CAClientConfig.Env, cfg.CAClientConfig.RootCertFile, cfg.CAClientConfig.KeyFile,
		cfg.CAClientConfig.CertChainFile, cfg.CAClientConfig.CAProviderName, cfg.MapperAudience)
	if err != nil {
		return nil, err
	}
	na.pc = pc

	// For VM, only support GoogleCA with the IstioCertificateService CA protocol.
	// TODO(pitlv2109): Support Citadel, Vault, etc.
	if cfg.CAClientConfig.CAProviderName == caclient.GoogleCAName {
		googleCA, err := caclient.NewCAClient(cfg.CAClientConfig.CAAddress, cfg.CAClientConfig.CAProviderName,
			true, nil, "", "", "", "")
		if err != nil {
			return nil, err
		}
		na.caClient = googleCA
		return na, nil
	}
	return nil, fmt.Errorf("CAs other than GoogleCA is not supported using new CA protocol in VM")
}

// initCitadelAgentForIstioCAService initializes nodeAgentIntervalV implementation of NodeAgent interface that
// talks to IstioCAService gRPC service defined at https://github.com/istio/istio/blob/master/security/proto/ca_service.proto
func initCitadelAgentForIstioCAService(cfg *Config) (*nodeAgentInternal, error) {
	na := &nodeAgentInternal{
		config:   cfg,
		certUtil: util.NewCertUtil(cfg.CAClientConfig.CSRGracePeriodPercentage),
	}

	pc, err := platform.NewClient(cfg.CAClientConfig.Env, cfg.CAClientConfig.RootCertFile, cfg.CAClientConfig.KeyFile,
		cfg.CAClientConfig.CertChainFile, cfg.CAClientConfig.CAProviderName, "")
	if err != nil {
		return nil, err
	}
	na.pc = pc
	dialOpts, err := pc.GetDialOptions()
	if err != nil {
		return nil, err
	}

	grpcConn, err := protocol.NewGrpcConnection(cfg.CAClientConfig.CAAddress, dialOpts)
	if err != nil {
		return nil, err
	}
	na.caProtocol = grpcConn
	return na, nil
}
