//  Copyright 2018 Istio Authors
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
	"os"
)

var (
	// ProxyIPAddress is the IP address of the proxy used to identify it and its
	// co-located service instances. Example: "10.60.1.6".
	ProxyIPAddress = os.Getenv("INSTANCE_IP")

	// pod name and namespace.
	PodName      = os.Getenv("POD_NAME")
	PodNamespace = os.Getenv("POD_NAMESPACE")

	// CertDir is the default location for mTLS certificates used by pilot.
	CertDir = os.Getenv("PILOT_CERT_DIR")

	// MaxConcurrentStreams indicates pilot max grpc concurrent streams.
	MaxConcurrentStreams = os.Getenv("ISTIO_GPRC_MAXSTREAMS")

	// Env var PILOT_TRACE_SAMPLING sets mesh-wide trace sampling
	// percentage, should be 0.0 - 100.0 Precision to 0.01
	TraceSampling = os.Getenv("PILOT_TRACE_SAMPLING")

	// CacheSquash is the max interval to squash a series of events.
	CacheSquash = os.Getenv("PILOT_CACHE_SQUASH")

	// Rate limit for the actual push.
	PushThrottle = os.Getenv("PILOT_PUSH_THROTTLE")
	PushBurst    = os.Getenv("PILOT_PUSH_BURST")

	// DebugConfigs controls saving snapshots of configs for /debug/adsz.
	// Defaults to false, can be enabled with PILOT_DEBUG_ADSZ_CONFIG=1
	DebugConfigs = os.Getenv("PILOT_DEBUG_ADSZ_CONFIG") == "1"

	// Duration of periodic refresh, in case events or cache invalidation fail.
	// Example: "300ms", "10s" or "2h45m".
	RefreshDuration = os.Getenv("V2_REFRESH")

	// kubeconfig file, if not specified default is `$HOME/.kube/config`
	KubeConfig = os.Getenv("KUBECONFIG")

	// DebounceAfter is the delay added to events to wait
	// after a registry/config event for debouncing.
	// This will delay the push by at least this interval, plus
	// the time getting subsequent events. If no change is
	// detected the push will happen, otherwise we'll keep
	// delaying until things settle.
	// Default is 100ms, Example: "300ms", "10s" or "2h45m".
	DebounceAfter = os.Getenv("PILOT_DEBOUNCE_AFTER")
	// DebounceMax is the maximum time to wait for events
	// while debouncing. Defaults to 10 seconds. If events keep
	// showing up with no break for this time, we'll trigger a push.
	// Default is 10s, Example: "300ms", "10s" or "2h45m".
	DebounceMax = os.Getenv("PILOT_DEBOUNCE_MAX")

	// IstioNamespace is the namespace where Istio's components are deployed.
	// Default is `istio-system`
	IstioNamespace = os.Getenv("ISTIO_NAMESPACE")

	// AzDebug indicates whether to log service registry az info.
	AzDebug = os.Getenv("VERBOSE_AZ_DEBUG") == "1"
)
