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

package securegatewaysds

import (
	"testing"

	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/istio.io/examples"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/galley"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/pilot"
	"istio.io/istio/pkg/test/framework/label"
	"istio.io/istio/pkg/test/framework/resource"
)

var (
	inst istio.Instance
	g    galley.Instance
	p    pilot.Instance
)

func TestMain(m *testing.M) {
	// Integration test for the ingress SDS Gateway flow.
	framework.
		NewSuite("sds_ingress_gateway_test", m).
		Label(label.CustomSetup).
		SetupOnEnv(environment.Kube, istio.Setup(&inst, setupConfig)).
		Setup(func(ctx resource.Context) (err error) {
			if g, err = galley.New(ctx, galley.Config{}); err != nil {
				return err
			}
			if p, err = pilot.New(ctx, pilot.Config{
				Galley: g,
			}); err != nil {
				return err
			}
			return nil
		}).
		Run()

}

func setupConfig(cfg *istio.Config) {
	if cfg == nil {
		return
	}
	// Enable SDS key/certificate provisioning for ingress gateway.
	cfg.Values["gateways.istio-ingressgateway.sds.enabled"] = "true"
	// TODO(https://github.com/istio/istio/issues/14084) remove this
	cfg.Values["pilot.env.PILOT_ENABLE_FALLTHROUGH_ROUTE"] = "0"
}

// https://preliminary.istio.io/docs/tasks/traffic-management/ingress/secure-ingress-sds/
func TestIngressSDS(t *testing.T) {
	framework.
		NewTest(t).
			Run(func(ctx framework.TestContext) {
				ex := examples.New(t, "secure-gateway-sds")

				env := ctx.Environment().(*kube.Environment)
				// Configure a TLS ingress gateway for a single host.
				ex.AddScript("", "httpbin-deployment.sh", examples.TextOutput)
				ex.AddScript("", "create-httpbin-tls-secret.sh", examples.TextOutput)
				ex.AddScript("", "create-httpbin-tls-gateway.sh", examples.TextOutput)
				ex.AddScript("", "check-envoy-sds-update-1.sh", examples.TextOutput)
				// Send an HTTPS request to access the httpbin service TLS gateway.
				if env.Settings().Minikube {
					ex.AddScript("", "curl-httpbin-tls-gateway-minikube.sh", examples.TextOutput)
				} else {
					ex.AddScript("", "curl-httpbin-tls-gateway-gke.sh", examples.TextOutput)
				}

				// Rotate secret and send HTTPS request with new credentials.
				ex.AddScript("", "rotate-httpbin-tls-secret.sh", examples.TextOutput)
				ex.AddScript("", "check-envoy-sds-update-2.sh", examples.TextOutput)
				if env.Settings().Minikube {
					ex.AddScript("", "curl-httpbin-tls-gateway-minikube-new-tls-secret.sh", examples.TextOutput)
				} else {
					ex.AddScript("", "curl-httpbin-tls-gateway-gke-new-tls-secret.sh", examples.TextOutput)
				}

				// Configure a TLS ingress gateway for multiple hosts
				ex.AddScript("", "restore-httpbin-tls-secret.sh", examples.TextOutput)
				ex.AddScript("", "helloworld-deployment.sh", examples.TextOutput)
				ex.AddScript("", "create-helloworld-tls-secret.sh", examples.TextOutput)
				ex.AddScript("", "create-helloworld-tls-gateway.sh", examples.TextOutput)
				// Send an HTTPS request to access the helloworld service TLS gateway.
				ex.AddScript("", "check-envoy-sds-update-4.sh", examples.TextOutput)
				if env.Settings().Minikube {
					ex.AddScript("", "curl-helloworld-tls-gateway-minikube.sh", examples.TextOutput)
				} else {
					ex.AddScript("", "curl-helloworld-tls-gateway-gke.sh", examples.TextOutput)
				}
				// Send an HTTPS request to access the httpbin service TLS gateway.
				if env.Settings().Minikube {
					ex.AddScript("", "curl-httpbin-tls-gateway-minikube.sh", examples.TextOutput)
				} else {
					ex.AddScript("", "curl-httpbin-tls-gateway-gke.sh", examples.TextOutput)
				}

				// Configure a mutual TLS ingress gateway
				ex.AddScript("", "rotate-httpbin-mtls-secret.sh", examples.TextOutput)
				ex.AddScript("", "create-httpbin-mtls-gateway.sh", examples.TextOutput)
				ex.AddScript("", "check-envoy-sds-update-5.sh", examples.TextOutput)
				// Send an HTTPS request to access the httpbin service mTLS gateway.
				if env.Settings().Minikube {
					ex.AddScript("", "curl-httpbin-mtls-gateway-minikube.sh", examples.TextOutput)
				} else {
					ex.AddScript("", "curl-httpbin-mtls-gateway-gke.sh", examples.TextOutput)
				}
				// Cleanup
				ex.AddScript("", "cleanup.sh", examples.TextOutput)
				ex.Run()
			})
}