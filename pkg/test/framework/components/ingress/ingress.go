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

package ingress

import (
	"istio.io/istio/pkg/test"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/resource"

	"time"
)

// GatewayType defines ingress gateway type
type GatewayType int

const (
	PlainText GatewayType = 0
	TLS       GatewayType = 1
	Mtls      GatewayType = 2
)

// CallOptions defines options for calling a Endpoint.
type CallOptions struct {
	// Host specifies the host to be used on the request. If not provided, an appropriate
	// default is chosen for the target Instance.
	Host string

	// Path specifies the URL path for the request.
	Path string

	// Timeout used for each individual request. Must be > 0, otherwise 1 minute is used.
	Timeout time.Duration
}

// Instance represents a deployed Ingress Gateway instance.
type Instance interface {
	resource.Resource

	// Address returns the external HTTP address of the ingress gateway (or the NodePort address,
	// when running under Minikube).
	Address() string

	//  Call makes a call through ingress.
	Call(options CallOptions) (CallResponse, error)
	CallOrFail(t test.Failer, options CallOptions) CallResponse
}

type Config struct {
	Istio istio.Instance
	// IngressType specifies the type of ingress gateway.
	IngressType GatewayType
	// CaCert is inline base64 encoded root certificate that authenticates server certificate provided
	// by ingress gateway.
	CaCert string
	// PrivateKey is inline base64 encoded private key for test client.
	PrivateKey string
	// Cert is inline base64 encoded certificate for test client.
	Cert string
}

// CallResponse is the result of a call made through Istio Ingress.
type CallResponse struct {
	// Response status code
	Code int

	// Response body
	Body string
}

// Deploy returns a new instance of echo.
func New(ctx resource.Context, cfg Config) (i Instance, err error) {
	err = resource.UnsupportedEnvironment(ctx.Environment())
	ctx.Environment().Case(environment.Kube, func() {
		i, err = newKube(ctx, cfg)
	})
	return
}

// Deploy returns a new Ingress instance or fails test
func NewOrFail(t test.Failer, ctx resource.Context, cfg Config) Instance {
	t.Helper()
	i, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("ingress.NewOrFail: %v", err)
	}

	return i
}
