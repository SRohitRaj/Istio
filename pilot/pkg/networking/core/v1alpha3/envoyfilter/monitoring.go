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
package envoyfilter

import (
	"istio.io/istio/pilot/pkg/features"
	"istio.io/pkg/monitoring"
)

type ResultType string

const (
	Error   ResultType = "error"
	Skipped ResultType = "skipped"
	Applied ResultType = "applied"
)

type PatchType string

const (
	Cluster       PatchType = "cluster"
	Listener      PatchType = "listener"
	FilterChain   PatchType = "filterchain"
	NetworkFilter PatchType = "networkfilter"
	// nolint
	HttpFilter  PatchType = "httpfilter"
	Route       PatchType = "route"
	VirtualHost PatchType = "vhost"
)

var (
	patchType = monitoring.MustCreateLabel("patch")
	errorType = monitoring.MustCreateLabel("type")
	nameType  = monitoring.MustCreateLabel("name")

	totalEnvoyFilters = monitoring.NewSum(
		"pilot_total_envoy_filter",
		"Total number of Envoy filters that were applied, skipped and errored.",
		monitoring.WithLabels(nameType, patchType, errorType),
	)
)

func init() {
	if features.EnableEnvoyFilterMetrics {
		monitoring.MustRegister(totalEnvoyFilters)
	}
}

// IncrementEnvoyFilterMetric increments filter metric.
func IncrementEnvoyFilterMetric(name string, pt PatchType, applied bool) {
	if !features.EnableEnvoyFilterMetrics {
		return
	}
	resultType := Applied
	if !applied {
		resultType = Skipped
	}
	totalEnvoyFilters.With(nameType.Value(name)).With(patchType.Value(string(pt))).
		With(errorType.Value(string(resultType))).Record(1)
}

// IncrementEnvoyFilterErrorMetric increments filter metric for errors.
func IncrementEnvoyFilterErrorMetric(name string, pt PatchType) {
	if !features.EnableEnvoyFilterMetrics {
		return
	}
	totalEnvoyFilters.With(nameType.Value(name)).With(patchType.Value(string(pt))).With(errorType.Value(string(Error))).Record(1)
}

// RecordEnvoyFilterMetric increments the filter metric with the given value.
func RecordEnvoyFilterMetric(name string, pt PatchType, success bool, value float64) {
	if !features.EnableEnvoyFilterMetrics {
		return
	}
	resultType := Applied
	if !success {
		resultType = Skipped
	}
	totalEnvoyFilters.With(nameType.Value(name)).With(patchType.Value(string(pt))).With(errorType.Value(string(resultType))).Record(value)
}
