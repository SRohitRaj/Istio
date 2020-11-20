// +build integ
// Copyright Istio Authors. All Rights Reserved.
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

package tracing

import (
	"fmt"
	"strings"
	"testing"

	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/echoboot"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/istio/ingress"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/components/zipkin"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/tests/integration/telemetry"
)

var (
	client, server echo.Instances
	ist            istio.Instance
	ingInst        ingress.Instance
	zipkinInst     zipkin.Instance
	appNsInst      namespace.Instance
)

const (
	TraceHeader = "x-client-trace-id"
)

func GetIstioInstance() *istio.Instance {
	return &ist
}

// GetAppNamespace gets echo app namespace instance.
func GetAppNamespace() namespace.Instance {
	return appNsInst
}

func GetIngressInstance() ingress.Instance {
	return ingInst
}

func GetZipkinInstance() zipkin.Instance {
	return zipkinInst
}

func TestSetup(ctx resource.Context) (err error) {
	appNsInst, err = namespace.New(ctx, namespace.Config{
		Prefix: "echo",
		Inject: true,
	})
	if err != nil {
		return
	}
	builder := echoboot.NewBuilder(ctx)
	for _, c := range ctx.Clusters() {
		builder.
			With(nil, echo.Config{
				Service:   clientNameForCluster(c.Name()),
				Namespace: appNsInst,
				Cluster:   c,
				Ports:     nil,
				Subsets:   []echo.SubsetConfig{{}},
			}).
			With(nil, echo.Config{
				Service:   "server",
				Namespace: appNsInst,
				Cluster:   c,
				Subsets:   []echo.SubsetConfig{{}},
				Ports: []echo.Port{
					{
						Name:         "http",
						Protocol:     protocol.HTTP,
						InstancePort: 8090,
					},
					{
						Name:     "tcp",
						Protocol: protocol.TCP,
						// We use a port > 1024 to not require root
						InstancePort: 9000,
					},
				},
			}).
			Build()
	}
	echos, err := builder.Build()
	if err != nil {
		return err
	}
	client = echos.Match(echo.ServicePrefix("client"))
	server = echos.Match(echo.Service("server"))
	ingInst = ist.IngressFor(ctx.Clusters().Default())
	zipkinInst, err = zipkin.New(ctx, zipkin.Config{Cluster: ctx.Clusters().Default(), IngressAddr: ingInst.HTTPAddress()})
	if err != nil {
		return
	}

	return nil
}

func VerifyEchoTraces(t *testing.T, namespace, clName string, traces []zipkin.Trace) bool {
	wtr := WantTraceRoot(namespace, clName)
	for _, trace := range traces {
		// compare each candidate trace with the wanted trace
		for _, s := range trace.Spans {
			// find the root span of candidate trace and do recursive comparison
			if s.ParentSpanID == "" && CompareTrace(t, s, wtr) {
				return true
			}
		}
	}

	return false
}

// compareTrace recursively compares the two given spans
func CompareTrace(t *testing.T, got, want zipkin.Span) bool {
	if got.Name != want.Name || got.ServiceName != want.ServiceName {
		t.Logf("got span %+v, want span %+v", got, want)
		return false
	}
	if len(got.ChildSpans) < len(want.ChildSpans) {
		t.Logf("got %d child spans from, want %d child spans, maybe trace has not be fully reported",
			len(got.ChildSpans), len(want.ChildSpans))
		return false
	} else if len(got.ChildSpans) > len(want.ChildSpans) {
		t.Logf("got %d child spans from, want %d child spans, maybe destination rule has not became effective",
			len(got.ChildSpans), len(want.ChildSpans))
		return false
	}
	for i := range got.ChildSpans {
		if !CompareTrace(t, *got.ChildSpans[i], *want.ChildSpans[i]) {
			return false
		}
	}
	return true
}

// wantTraceRoot constructs the wanted trace and returns the root span of that trace
func WantTraceRoot(namespace, clName string) (root zipkin.Span) {
	serverSpan := zipkin.Span{
		Name:        fmt.Sprintf("server.%s.svc.cluster.local:80/*", namespace),
		ServiceName: fmt.Sprintf("server.%s", namespace),
	}

	root = zipkin.Span{
		Name:        fmt.Sprintf("server.%s.svc.cluster.local:80/*", namespace),
		ServiceName: fmt.Sprintf("%s.%s", clientNameForCluster(clName), namespace),
		ChildSpans:  []*zipkin.Span{&serverSpan},
	}
	return
}

// SendTraffic makes a client call to the "server" service on the http port.
func SendTraffic(t *testing.T, headers map[string][]string, cl resource.Cluster) error {
	t.Log("Sending Traffic...")
	for _, cltInstance := range client {
		if cltInstance.Config().Cluster != cl {
			continue
		}

		_, err := cltInstance.Call(echo.CallOptions{
			Target:   server[0],
			PortName: "http",
			Count:    telemetry.RequestCountMultipler * len(server),
			Headers:  headers,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func clientNameForCluster(clusterName string) string {
	// Convert the cluster name into a valid k8s object name.
	return "client-" + strings.Replace(clusterName, "_", "-", -1)
}
