//go:build integ
// +build integ

//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package sdsingressk8sca

import (
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/common/deployment"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/tests/integration/security/sds_ingress/util"
)

var (
	inst         istio.Instance
	apps         deployment.SingleNamespaceView
	echo1NS      namespace.Instance
	customConfig []echo.Config
)

func TestMain(m *testing.M) {
	// Integration test for the ingress SDS multiple Gateway flow when
	// the control plane certificate provider is k8s CA.
	// nolint: staticcheck
	framework.
		NewSuite(m).
		RequireMinVersion(20). // versions less than 1.20 doesn't have kube-root-ca.crt configmap. https://github.com/istio/istio/pull/42111
		// https://github.com/istio/istio/issues/22161. 1.22 drops support for legacy-unknown signer
		RequireMaxVersion(21).
		Setup(istio.Setup(&inst, setupConfig)).
		Setup(namespace.Setup(&echo1NS, namespace.Config{Prefix: "echo1", Inject: true})).
		Setup(func(ctx resource.Context) error {
			// Skip VM as eastwest gateway is disabled.
			s := ctx.Settings()
			s.SkipWorkloadClasses = append(s.SkipWorkloadClasses, echo.VM)
			err := util.SetupTest(ctx, &customConfig, namespace.Future(&echo1NS))
			if err != nil {
				return err
			}
			return nil
		}).
		Setup(deployment.SetupSingleNamespace(&apps, deployment.Config{
			Namespaces: []namespace.Getter{
				namespace.Future(&echo1NS),
			},
			Configs: echo.ConfigFuture(&customConfig),
		})).
		Setup(func(ctx resource.Context) error {
			return util.CreateCustomInstances(&apps)
		}).
		Run()
}

func setupConfig(_ resource.Context, cfg *istio.Config) {
	if cfg == nil {
		return
	}
	cfg.ControlPlaneValues = `
values:
  global:
    pilotCertProvider: kubernetes
`
}

func TestMtlsGatewaysK8sca(t *testing.T) {
	framework.
		NewTest(t).
		Features("security.ingress.mtls.gateway").
		Run(func(t framework.TestContext) {
			t.NewSubTest("tcp").Run(func(t framework.TestContext) {
				util.RunTestMultiMtlsGateways(t, inst, namespace.Future(&echo1NS))
			})
		})
}

func TestTlsGatewaysK8sca(t *testing.T) {
	framework.
		NewTest(t).
		Features("security.ingress.tls.gateway.K8sca").
		Run(func(t framework.TestContext) {
			t.NewSubTest("tcp").Run(func(t framework.TestContext) {
				util.RunTestMultiTLSGateways(t, inst, namespace.Future(&echo1NS))
			})
		})
}
