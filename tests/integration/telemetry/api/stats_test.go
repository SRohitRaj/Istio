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

package api

import (
	"testing"

	"istio.io/istio/tests/integration/telemetry/common"
)

// TestTelemetryAPIStats verifies the stats filter could emit expected client and server side
// metrics when configured with the Telemetry API (with EnvoyFilters disabled)
// This test focuses on stats filter and metadata exchange filter could work coherently with
// proxy bootstrap config with Wasm runtime. To avoid flake, it does not verify correctness
// of metrics, which should be covered by integration test in proxy repo.
func TestTelemetryAPIStats(t *testing.T) {
	common.TestStatsFilter(t, "observability.telemetry.stats.prometheus.http.nullvm", common.DefaultBucketCount)
}

func TestTelemetryAPITCPStats(t *testing.T) { // nolint:interfacer
	common.TestStatsTCPFilter(t, "observability.telemetry.stats.prometheus.tcp")
}
