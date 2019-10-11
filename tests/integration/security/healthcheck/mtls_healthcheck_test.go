//  Copyright 2019 Istio Authors
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

// Package healthcheck contains a test to support kubernetes app health check when mTLS is turned on.
// https://github.com/istio/istio/issues/9150.
package healthcheck

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"

	// "istio.io/istio/pkg/test/framework/components/echo/common"
	"istio.io/istio/pkg/test/framework/components/echo/echoboot"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/namespace"
)

// TestMtlsHealthCheck verifies Kubernetes HTTP health check can work when mTLS
// is enabled.
// Currently this test can only pass on Prow with a real GKE cluster, and fail
// on Minikube. For more details, see https://github.com/istio/istio/issues/12754.
func TestMtlsHealthCheck(t *testing.T) {
	framework.NewTest(t).
		RequiresEnvironment(environment.Kube).
		Run(func(ctx framework.TestContext) {
			ns := namespace.ClaimOrFail(t, ctx, "mtls-healthcheck")
			runHealthCheckDeployment(t, ctx, ns, "healthcheck", true, true)
			runHealthCheckDeployment(t, ctx, ns, "healthcheck-fail", false, false)
		})
}

func runHealthCheckDeployment(t *testing.T, ctx framework.TestContext, ns namespace.Instance,
	name string, rewrite bool, success bool) {
	t.Helper()
	policyYAML := fmt.Sprintf(`apiVersion: "authentication.istio.io/v1alpha1"
kind: "Policy"
metadata:
  name: "mtls-strict-for-%v"
spec:
  targets:
  - name: "%v"
  peers:
    - mtls:
        mode: STRICT
`, name, name)
	g.ApplyConfigOrFail(t, ns, policyYAML)
	defer g.DeleteConfigOrFail(t, ns, policyYAML)

	var healthcheck echo.Instance
	cfg := echo.Config{
		Namespace: ns,
		Service:   name,
		Pilot:     p,
		Galley:    g,
		Ports: []echo.Port{{
			Name:         "http-8080",
			Protocol:     protocol.HTTP,
			ServicePort:  8080,
			InstancePort: 8080,
		}},
		ReadinessTimeout: time.Second * 60,
	}
	cfg.Annotations = map[echo.Annotation]*echo.AnnotationValue{
		echo.SidecarRewriteAppHTTPProbers: &echo.AnnotationValue{Value: strconv.FormatBool(rewrite)},
	}
	err := echoboot.NewBuilderOrFail(t, ctx).
		With(&healthcheck, cfg).
		Build()
	gotSuccess := err == nil
	if gotSuccess != success {
		t.Errorf("health check app %v, got error %v, want success = %v", name, err, success)
	}
}
