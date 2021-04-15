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

package virtualservice

import (
	"fmt"

	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/analyzers/util"
	"istio.io/istio/galley/pkg/config/analysis/msg"
	"istio.io/istio/pkg/config/host"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

// GatewayAnalyzer checks the gateways associated with each virtual service
type GatewayAnalyzer struct{}

var _ analysis.Analyzer = &GatewayAnalyzer{}

// Metadata implements Analyzer
func (s *GatewayAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "virtualservice.GatewayAnalyzer",
		Description: "Checks the gateways associated with each virtual service",
		Inputs: collection.Names{
			collections.IstioNetworkingV1Alpha3Gateways.Name(),
			collections.IstioNetworkingV1Alpha3Virtualservices.Name(),
		},
	}
}

// Analyze implements Analyzer
func (s *GatewayAnalyzer) Analyze(c analysis.Context) {
	c.ForEach(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), func(r *resource.Instance) bool {
		s.analyzeVirtualService(r, c)
		return true
	})
}

func (s *GatewayAnalyzer) analyzeVirtualService(r *resource.Instance, c analysis.Context) {
	vs := r.Message.(*v1alpha3.VirtualService)
	vsNs := r.Metadata.FullName.Namespace
	vsName := r.Metadata.FullName
	gwMap := map[string]bool{}

	for i, gwName := range vs.Gateways {
		gwMap[gwName] = true
		// This is a special-case accepted value
		if gwName == util.MeshGateway {
			continue
		}

		gwFullName := resource.NewShortOrFullName(vsNs, gwName)

		if !c.Exists(collections.IstioNetworkingV1Alpha3Gateways.Name(), gwFullName) {
			m := msg.NewReferencedResourceNotFound(r, "gateway", gwName)

			if line, ok := util.ErrorLine(r, fmt.Sprintf(util.VSGateway, i)); ok {
				m.Line = line
			}

			c.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), m)
		}

		if !vsHostInGateway(c, gwFullName, vs.Hosts) {
			m := msg.NewVirtualServiceHostNotFoundInGateway(r, vs.Hosts, vsName.String(), gwFullName.String())

			if line, ok := util.ErrorLine(r, fmt.Sprintf(util.VSGateway, i)); ok {
				m.Line = line
			}

			c.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), m)
		}
	}
	gatewayNotDeclaredReport := func(g string, protocol string, ri, mi, gi int) {
		m := msg.NewReferencedResourceNotFound(r, "gateway", g)
		if line, ok := util.ErrorLine(r, fmt.Sprintf(util.VSMatchGateway, protocol, ri, mi, gi)); ok {
			m.Line = line
		}
		c.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), m)
	}
	for i, r := range vs.Http {
		for j, m := range r.Match {
			for k, g := range m.Gateways {
				if _, ok := gwMap[g]; !ok {
					gatewayNotDeclaredReport(g, "http", i, j, k)
				}
			}
		}
	}
	for i, r := range vs.Tls {
		for j, m := range r.Match {
			for k, g := range m.Gateways {
				if _, ok := gwMap[g]; !ok {
					gatewayNotDeclaredReport(g, "tls", i, j, k)
				}
			}
		}
	}
	for i, r := range vs.Tcp {
		for j, m := range r.Match {
			for k, g := range m.Gateways {
				if _, ok := gwMap[g]; !ok {
					gatewayNotDeclaredReport(g, "tcp", i, j, k)
				}
			}
		}
	}
}

func vsHostInGateway(c analysis.Context, gateway resource.FullName, vsHosts []string) bool {
	var gatewayHosts []string

	c.ForEach(collections.IstioNetworkingV1Alpha3Gateways.Name(), func(r *resource.Instance) bool {
		if r.Metadata.FullName == gateway {
			s := r.Message.(*v1alpha3.Gateway)

			for _, v := range s.Servers {
				gatewayHosts = append(gatewayHosts, v.Hosts...)
			}
		}

		return true
	})

	for _, gh := range gatewayHosts {
		for _, vsh := range vsHosts {
			gatewayHost := host.Name(gh)
			vsHost := host.Name(vsh)

			if gatewayHost.Matches(vsHost) {
				return true
			}
		}
	}

	return false
}
