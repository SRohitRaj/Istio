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

package prometheus

import (
	"testing"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prom "github.com/prometheus/common/model"

	"istio.io/istio/pkg/test/framework2/core"
)

type Instance interface {
	core.Resource

	// API Returns the core Prometheus APIs.
	API() v1.API

	// WaitForQuiesce runs the provided query periodically until the result gets stable.
	WaitForQuiesce(fmt string, args ...interface{}) (prom.Value, error)

	// WaitForOneOrMore runs the provided query and waits until one (or more for vector) values are available.
	WaitForOneOrMore(fmt string, args ...interface{}) error

	// Sum all the samples that has the given labels in the given vector value.
	Sum(val prom.Value, labels map[string]string) (float64, error)
}

// New returns a new Prometheus instance.
func New(ctx core.Context) (Instance, error) {
	switch ctx.Environment().EnvironmentName() {
	case core.Kube:
		return newKube(ctx)
	default:
		return nil, core.UnsupportedEnvironment(ctx.Environment().EnvironmentName())
	}
}

// NewOrFail returns a new Prometheus instance or fails test.
func NewOrFail(t *testing.T, ctx core.Context) Instance {
	t.Helper()
	i, err := New(ctx)
	if err != nil {
		t.Fatalf("prometheus.NewOrFail: %v", err)
	}

	return i
}
