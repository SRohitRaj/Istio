//go:build integ
// +build integ

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

package pilot

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "sigs.k8s.io/gateway-api/apis/v1beta1"

	"istio.io/istio/pilot/pkg/model/kstatus"
	"istio.io/istio/pkg/config/constants"
	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/http/headers"
	"istio.io/istio/pkg/test/echo/common/scheme"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/crd"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/check"
	"istio.io/istio/pkg/test/framework/components/istio"
	testKube "istio.io/istio/pkg/test/kube"
	"istio.io/istio/pkg/test/util/assert"
	"istio.io/istio/pkg/test/util/retry"
	ingressutil "istio.io/istio/tests/integration/security/sds_ingress/util"
)

func TestGateway(t *testing.T) {
	framework.
		NewTest(t).
		Run(func(t framework.TestContext) {
			crd.DeployGatewayAPIOrSkip(t)

			t.NewSubTest("unmanaged").Run(UnmanagedGatewayTest)
			t.NewSubTest("managed").Run(ManagedGatewayTest)
			t.NewSubTest("managed-owner").Run(ManagedOwnerGatewayTest)
			t.NewSubTest("status").Run(StatusGatewayTest)
			t.NewSubTest("managed-short-name").Run(ManagedGatewayShortNameTest)
		})
}

func ManagedOwnerGatewayTest(t framework.TestContext) {
	image := fmt.Sprintf("%s/app:%s", t.Settings().Image.Hub, strings.TrimSuffix(t.Settings().Image.Tag, "-distroless"))
	t.ConfigIstio().YAML(apps.Namespace.Name(), fmt.Sprintf(`
apiVersion: v1
kind: Service
metadata:
  name: managed-owner-istio
spec:
  ports:
  - appProtocol: http
    name: default
    port: 80
  selector:
    istio.io/gateway-name: managed-owner
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: managed-owner-istio
spec:
  selector:
    matchLabels:
      istio.io/gateway-name: managed-owner
  replicas: 1
  template:
    metadata:
      labels:
        istio.io/gateway-name: managed-owner
    spec:
      containers:
      - name: fake
        image: %s
`, image)).ApplyOrFail(t)
	cls := t.Clusters().Kube().Default()
	fetchFn := testKube.NewSinglePodFetch(cls, apps.Namespace.Name(), "istio.io/gateway-name=managed-owner")
	if _, err := testKube.WaitUntilPodsAreReady(fetchFn); err != nil {
		t.Fatal(err)
	}

	t.ConfigIstio().YAML(apps.Namespace.Name(), `
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: managed-owner
spec:
  gatewayClassName: istio
  listeners:
  - name: default
    hostname: "*.example.com"
    port: 80
    protocol: HTTP
`).ApplyOrFail(t)

	// Make sure Gateway becomes programmed..
	client := t.Clusters().Kube().Default().GatewayAPI().GatewayV1beta1().Gateways(apps.Namespace.Name())
	check := func() error {
		gw, _ := client.Get(context.Background(), "managed-owner", metav1.GetOptions{})
		if gw == nil {
			return fmt.Errorf("failed to find gateway")
		}
		cond := kstatus.GetCondition(gw.Status.Conditions, string(k8s.GatewayConditionProgrammed))
		if cond.Status != metav1.ConditionTrue {
			return fmt.Errorf("failed to find programmed condition: %+v", cond)
		}
		if cond.ObservedGeneration != gw.Generation {
			return fmt.Errorf("stale GWC generation: %+v", cond)
		}
		return nil
	}
	retry.UntilSuccessOrFail(t, check)

	// Make sure we did not overwrite our deployment or service
	dep, err := t.Clusters().Kube().Default().Kube().AppsV1().Deployments(apps.Namespace.Name()).
		Get(context.Background(), "managed-owner-istio", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, dep.Labels[constants.ManagedGatewayLabel], "")
	assert.Equal(t, dep.Spec.Template.Spec.Containers[0].Image, image)

	svc, err := t.Clusters().Kube().Default().Kube().CoreV1().Services(apps.Namespace.Name()).
		Get(context.Background(), "managed-owner-istio", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, svc.Labels[constants.ManagedGatewayLabel], "")
	assert.Equal(t, svc.Spec.Type, corev1.ServiceTypeClusterIP)
}

func ManagedGatewayTest(t framework.TestContext) {
	t.ConfigIstio().YAML(apps.Namespace.Name(), `apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: gateway
spec:
  gatewayClassName: istio
  listeners:
  - name: default
    hostname: "*.example.com"
    port: 80
    protocol: HTTP
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: http-1
spec:
  parentRefs:
  - name: gateway
  hostnames: ["bar.example.com"]
  rules:
  - backendRefs:
    - name: b
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: http-2
spec:
  parentRefs:
  - name: gateway
  hostnames: ["foo.example.com"]
  rules:
  - backendRefs:
    - name: d
      port: 80
`).ApplyOrFail(t)
	testCases := []struct {
		check echo.Checker
		from  echo.Instances
		host  string
	}{
		{
			check: check.OK(),
			from:  apps.B,
			host:  "bar.example.com",
		},
		{
			check: check.NotOK(),
			from:  apps.B,
			host:  "bar",
		},
	}
	if t.Settings().EnableDualStack {
		additionalTestCases := []struct {
			check echo.Checker
			from  echo.Instances
			host  string
		}{
			// apps.D hosts a dual-stack service,
			// apps.E hosts an ipv6 only service and
			// apps.B hosts an ipv4 only service
			{
				check: check.OK(),
				from:  apps.D,
				host:  "bar.example.com",
			},
			{
				check: check.OK(),
				from:  apps.E,
				host:  "bar.example.com",
			},
			{
				check: check.OK(),
				from:  apps.E,
				host:  "foo.example.com",
			},
			{
				check: check.OK(),
				from:  apps.D,
				host:  "foo.example.com",
			},
			{
				check: check.OK(),
				from:  apps.B,
				host:  "foo.example.com",
			},
		}
		testCases = append(testCases, additionalTestCases...)
	}
	for _, tc := range testCases {
		t.NewSubTest(fmt.Sprintf("gateway-connectivity-from-%s", tc.from[0].NamespacedName())).Run(func(t framework.TestContext) {
			tc.from[0].CallOrFail(t, echo.CallOptions{
				Port: echo.Port{
					Protocol:    protocol.HTTP,
					ServicePort: 80,
				},
				Scheme: scheme.HTTP,
				HTTP: echo.HTTP{
					Headers: headers.New().WithHost(tc.host).Build(),
				},
				Address: fmt.Sprintf("gateway-istio.%s.svc.cluster.local", apps.Namespace.Name()),
				Check:   tc.check,
			})
		})
	}
}

func ManagedGatewayShortNameTest(t framework.TestContext) {
	t.ConfigIstio().YAML(apps.Namespace.Name(), `apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: gateway
spec:
  gatewayClassName: istio
  listeners:
  - name: default
    hostname: "bar"
    port: 80
    protocol: HTTP
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: http
spec:
  parentRefs:
  - name: gateway
  rules:
  - backendRefs:
    - name: b
      port: 80
`).ApplyOrFail(t)
	apps.B[0].CallOrFail(t, echo.CallOptions{
		Port:   echo.Port{ServicePort: 80},
		Scheme: scheme.HTTP,
		HTTP: echo.HTTP{
			Headers: headers.New().WithHost("bar").Build(),
		},
		Address: fmt.Sprintf("gateway-istio.%s.svc.cluster.local", apps.Namespace.Name()),
		Check:   check.OK(),
		Retry: echo.Retry{
			Options: []retry.Option{retry.Timeout(time.Minute)},
		},
	})
	apps.B[0].CallOrFail(t, echo.CallOptions{
		Port:   echo.Port{ServicePort: 80},
		Scheme: scheme.HTTP,
		HTTP: echo.HTTP{
			Headers: headers.New().WithHost("bar.example.com").Build(),
		},
		Address: fmt.Sprintf("gateway-istio.%s.svc.cluster.local", apps.Namespace.Name()),
		Check:   check.NotOK(),
		Retry: echo.Retry{
			Options: []retry.Option{retry.Timeout(time.Minute)},
		},
	})
}

func UnmanagedGatewayTest(t framework.TestContext) {
	ingressutil.CreateIngressKubeSecret(t, "test-gateway-cert-same", ingressutil.TLS, ingressutil.IngressCredentialA,
		false, t.Clusters().Configs()...)
	ingressutil.CreateIngressKubeSecretInNamespace(t, "test-gateway-cert-cross", ingressutil.TLS, ingressutil.IngressCredentialB,
		false, apps.Namespace.Name(), t.Clusters().Configs()...)

	t.ConfigIstio().
		YAML("", `
apiVersion: gateway.networking.k8s.io/v1beta1
kind: GatewayClass
metadata:
  name: custom-istio
spec:
  controllerName: istio.io/gateway-controller
`).
		YAML("", fmt.Sprintf(`
apiVersion: gateway.networking.k8s.io/v1beta1
kind: Gateway
metadata:
  name: gateway
  namespace: istio-system
spec:
  addresses:
  - value: istio-ingressgateway
    type: Hostname
  gatewayClassName: custom-istio
  listeners:
  - name: http
    hostname: "*.domain.example"
    port: 80
    protocol: HTTP
    allowedRoutes:
      namespaces:
        from: All
  - name: tcp
    port: 31400
    protocol: TCP
    allowedRoutes:
      namespaces:
        from: All
  - name: tls-cross
    hostname: cross-namespace.domain.example
    port: 443
    protocol: HTTPS
    allowedRoutes:
      namespaces:
        from: All
    tls:
      mode: Terminate
      certificateRefs:
      - kind: Secret
        name: test-gateway-cert-cross
        namespace: "%s"
  - name: tls-same
    hostname: same-namespace.domain.example
    port: 443
    protocol: HTTPS
    allowedRoutes:
      namespaces:
        from: All
    tls:
      mode: Terminate
      certificateRefs:
      - kind: Secret
        name: test-gateway-cert-same
`, apps.Namespace.Name())).
		YAML(apps.Namespace.Name(), `
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: http
spec:
  hostnames: ["my.domain.example"]
  parentRefs:
  - name: gateway
    namespace: istio-system
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /get/
    backendRefs:
    - name: b
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: TCPRoute
metadata:
  name: tcp
spec:
  parentRefs:
  - name: gateway
    namespace: istio-system
  rules:
  - backendRefs:
    - name: b
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: b
spec:
  parentRefs:
  - group: ""
    kind: Service
    name: b
  - name: gateway
    namespace: istio-system
  hostnames: ["b"]
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /path
    filters:
    - type: RequestHeaderModifier
      requestHeaderModifier:
        add:
        - name: my-added-header
          value: added-value
    backendRefs:
    - name: b
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1alpha2
kind: GRPCRoute
metadata:
  name: grpc
spec:
  parentRefs:
  - group: ""
    kind: Service
    name: c
  - name: gateway
    namespace: istio-system
  rules:
  - matches:
    - method:
        method: Echo
    filters:
    - type: RequestHeaderModifier
      requestHeaderModifier:
        add:
        - name: my-added-header
          value: added-grpc-value
    backendRefs:
    - name: c
      port: 7070
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: tls-same
spec:
  parentRefs:
  - name: gateway
    sectionName: tls-same
    namespace: istio-system
  rules:
  - backendRefs:
    - name: b
      port: 80
---
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: tls-cross
spec:
  parentRefs:
  - name: gateway
    sectionName: tls-cross
    namespace: istio-system
  rules:
  - backendRefs:
    - name: b
      port: 80
`).YAML(apps.Namespace.Name(), fmt.Sprintf(`
apiVersion: gateway.networking.k8s.io/v1beta1
kind: ReferenceGrant
metadata:
  name: allow-gateways-to-ref-secrets
  namespace: "%s"
spec:
  from:
  - group: gateway.networking.k8s.io
    kind: Gateway
    namespace: istio-system
  to:
  - group: ""
    kind: Secret
`, apps.Namespace.Name())).
		ApplyOrFail(t)
	for _, ingr := range istio.IngressesOrFail(t, t) {
		t.NewSubTest(ingr.Cluster().StableName()).Run(func(t framework.TestContext) {
			t.NewSubTest("http").Run(func(t framework.TestContext) {
				paths := []string{"/get", "/get/", "/get/prefix"}
				for _, path := range paths {
					_ = ingr.CallOrFail(t, echo.CallOptions{
						Port: echo.Port{
							Protocol: protocol.HTTP,
						},
						HTTP: echo.HTTP{
							Path:    path,
							Headers: headers.New().WithHost("my.domain.example").Build(),
						},
						Check: check.OK(),
					})
				}
			})
			t.NewSubTest("tcp").Run(func(t framework.TestContext) {
				_ = ingr.CallOrFail(t, echo.CallOptions{
					Port: echo.Port{
						Protocol:    protocol.HTTP,
						ServicePort: 31400,
					},
					HTTP: echo.HTTP{
						Path:    "/",
						Headers: headers.New().WithHost("my.domain.example").Build(),
					},
					Check: check.OK(),
				})
			})
			t.NewSubTest("mesh").Run(func(t framework.TestContext) {
				_ = apps.A[0].CallOrFail(t, echo.CallOptions{
					To:    apps.B,
					Count: 1,
					Port: echo.Port{
						Name: "http",
					},
					HTTP: echo.HTTP{
						Path: "/path",
					},
					Check: check.And(
						check.OK(),
						check.RequestHeader("My-Added-Header", "added-value")),
				})
			})
			t.NewSubTest("mesh-grpc").Run(func(t framework.TestContext) {
				_ = apps.A[0].CallOrFail(t, echo.CallOptions{
					To:    apps.C,
					Count: 1,
					Port: echo.Port{
						Name: "grpc",
					},
					Check: check.And(
						check.OK(),
						check.RequestHeader("My-Added-Header", "added-grpc-value")),
				})
			})
			t.NewSubTest("status").Run(func(t framework.TestContext) {
				retry.UntilSuccessOrFail(t, func() error {
					gwc, err := t.Clusters().Kube().Default().GatewayAPI().GatewayV1beta1().GatewayClasses().Get(context.Background(), "istio", metav1.GetOptions{})
					if err != nil {
						return err
					}
					if s := kstatus.GetCondition(gwc.Status.Conditions, string(k8s.GatewayClassConditionStatusAccepted)).Status; s != metav1.ConditionTrue {
						return fmt.Errorf("expected status %q, got %q", metav1.ConditionTrue, s)
					}
					return nil
				})
			})
			t.NewSubTest("tls-same").Run(func(t framework.TestContext) {
				_ = ingr.CallOrFail(t, echo.CallOptions{
					Port: echo.Port{
						Protocol:    protocol.HTTPS,
						ServicePort: 443,
					},
					HTTP: echo.HTTP{
						Path:    "/",
						Headers: headers.New().WithHost("same-namespace.domain.example").Build(),
					},
					Check: check.OK(),
				})
			})
			t.NewSubTest("tls-cross").Run(func(t framework.TestContext) {
				_ = ingr.CallOrFail(t, echo.CallOptions{
					Port: echo.Port{
						Protocol:    protocol.HTTPS,
						ServicePort: 443,
					},
					HTTP: echo.HTTP{
						Path:    "/",
						Headers: headers.New().WithHost("cross-namespace.domain.example").Build(),
					},
					Check: check.OK(),
				})
			})
		})
	}
}

func StatusGatewayTest(t framework.TestContext) {
	client := t.Clusters().Kube().Default().GatewayAPI().GatewayV1beta1().GatewayClasses()

	check := func() error {
		gwc, _ := client.Get(context.Background(), "istio", metav1.GetOptions{})
		if gwc == nil {
			return fmt.Errorf("failed to find GatewayClass istio")
		}
		cond := kstatus.GetCondition(gwc.Status.Conditions, string(k8s.GatewayClassConditionStatusAccepted))
		if cond.Status != metav1.ConditionTrue {
			return fmt.Errorf("failed to find accepted condition: %+v", cond)
		}
		if cond.ObservedGeneration != gwc.Generation {
			return fmt.Errorf("stale GWC generation: %+v", cond)
		}
		return nil
	}
	retry.UntilSuccessOrFail(t, check)

	// Wipe out the status
	gwc, _ := client.Get(context.Background(), "istio", metav1.GetOptions{})
	gwc.Status.Conditions = nil
	client.Update(context.Background(), gwc, metav1.UpdateOptions{})
	// It should be added back
	retry.UntilSuccessOrFail(t, check)
}

// Verify that the envoy readiness probes are reachable at
// https://GatewayPodIP:15021/healthz/ready . This is being explicitly done
// to make sure, in dual-stack scenarios both v4 and v6 probes are reachable.
func TestGatewayReadinessProbes(t *testing.T) {
	// nolint: staticcheck
	framework.NewTest(t).
		RequiresSingleCluster().
		RequiresLocalControlPlane().
		Features("traffic.gateway.readiness").
		Run(func(t framework.TestContext) {
			c := t.Clusters().Default()
			podIPs, err := i.PodIPsFor(c, "app=istio-ingressgateway")
			if err != nil {
				t.Fatalf("error getting ingress gateway pod ips: %v", err)
			}
			for _, ip := range podIPs {
				t.NewSubTest("gateway-readiness-probe-" + ip.IP).Run(func(t framework.TestContext) {
					apps.External.All[0].CallOrFail(t, echo.CallOptions{
						Address: ip.IP,
						Port:    echo.Port{ServicePort: 15021},
						Scheme:  scheme.HTTP,
						HTTP: echo.HTTP{
							Path: "/healthz/ready",
						},
						Check: check.And(
							check.Status(200),
						),
					})
				})
			}
		})
}

// Verify that the envoy metrics endpoints are reachable at
// https://GatewayPodIP:15090/stats/prometheus . This is being explicitly done
// to make sure, in dual-stack scenarios both v4 and v6 probes are reachable.
func TestGatewayMetricsEndpoints(t *testing.T) {
	// nolint: staticcheck
	framework.NewTest(t).
		RequiresSingleCluster().
		RequiresLocalControlPlane().
		Features("traffic.gateway.metrics").
		Run(func(t framework.TestContext) {
			c := t.Clusters().Default()
			podIPs, err := i.PodIPsFor(c, "app=istio-ingressgateway")
			if err != nil {
				t.Fatalf("error getting ingress gateway pod ips: %v", err)
			}
			for _, ip := range podIPs {
				t.NewSubTest("gateway-metrics-endpoints-" + ip.IP).Run(func(t framework.TestContext) {
					apps.External.All[0].CallOrFail(t, echo.CallOptions{
						Address: ip.IP,
						Port:    echo.Port{ServicePort: 15090},
						Scheme:  scheme.HTTP,
						HTTP: echo.HTTP{
							Path: "/stats/prometheus",
						},
						Check: check.And(
							check.Status(200),
						),
					})
				})
			}
		})
}
