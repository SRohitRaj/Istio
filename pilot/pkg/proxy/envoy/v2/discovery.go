// Copyright 2018 Istio Authors
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

package v2

import (
	"os"
	"strconv"
	"sync"
	"time"

	ads "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"google.golang.org/grpc"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/core"
)

var (
	// Failsafe to implement periodic refresh, in case events or cache invalidation fail.
	// Disabled by default.
	periodicRefreshDuration = 0 * time.Second

	versionMutex sync.RWMutex

	// version is the timestamp of the last registry event.
	version = "0"

	periodicRefreshMetrics = 10 * time.Second
)

const (
	typePrefix = "type.googleapis.com/envoy.api.v2."

	// Constants used for XDS

	// ClusterType is used for cluster discovery. Typically first request received
	ClusterType = typePrefix + "Cluster"
	// EndpointType is used for EDS and ADS endpoint discovery. Typically second request.
	EndpointType = typePrefix + "ClusterLoadAssignment"
	// ListenerType is sent after clusters and endpoints.
	ListenerType = typePrefix + "Listener"
	// RouteType is sent after listeners.
	RouteType = typePrefix + "RouteConfiguration"
)

// DiscoveryServer is Pilot's gRPC implementation for Envoy's v2 xds APIs
type DiscoveryServer struct {
	// Env is the model environment.
	Env *model.Environment

	// MemRegistry is used for debug and load testing, allow adding services. Visible for testing.
	MemRegistry *MemServiceDiscovery

	// ConfigGenerator is responsible for generating data plane configuration using Istio networking
	// APIs and service registry info
	ConfigGenerator core.ConfigGenerator

	// ConfigController provides readiness info (if initial sync is complete)
	ConfigController model.ConfigStoreCache

	// separate rate limiter for initial connection
	initThrottle chan time.Time

	throttle chan time.Time

	// DebugConfigs controls saving snapshots of configs for /debug/adsz.
	// Defaults to false, can be enabled with PILOT_DEBUG_ADSZ_CONFIG=1
	DebugConfigs bool

	// ServiceShards for a service. This is a global (per-server) list, built from
	// incremental updates.
	EndpointShardsByService map[string]model.ServiceShards

	mutex sync.RWMutex
}

func intEnv(env string, def int) int {
	envValue := os.Getenv(env)
	if len(envValue) == 0 {
		return def
	}
	n, err := strconv.Atoi(envValue)
	if err == nil && n > 0 {
		return n
	}
	return def
}

// NewDiscoveryServer creates DiscoveryServer that sources data from Pilot's internal mesh data structures
func NewDiscoveryServer(env *model.Environment, generator core.ConfigGenerator) *DiscoveryServer {
	out := &DiscoveryServer{
		Env:                     env,
		ConfigGenerator:         generator,
		EndpointShardsByService: map[string]model.ServiceShards{},
	}
	env.PushContext = model.NewPushContext()

	go out.periodicRefresh()

	go out.periodicRefreshMetrics()

	out.DebugConfigs = os.Getenv("PILOT_DEBUG_ADSZ_CONFIG") == "1"

	pushThrottle := intEnv("PILOT_PUSH_THROTTLE", 25)
	pushBurst := intEnv("PILOT_PUSH_BURST", 100)

	adsLog.Infof("Starting ADS server with throttle=%d burst=%d", pushThrottle, pushBurst)

	// throttle rate limits the amount of `pushALL` work that is started as a result of events.
	out.throttle = initThrottle("adsPushAll", pushBurst, pushThrottle)

	// init throttle rate limits starting work on new connections from sidecars.
	out.initThrottle = initThrottle("initConnection", pushBurst*2, pushThrottle*2)

	// Note: in both cases it does not directly limit the amount of work being perform concurrently.
	// If a particular push takes a long time, it will allow more and more work, and token are being replenished
	// as work is being performed.

	return out
}

// initThrottle allocates and initializes a throttle channel with burstLimit and steady state ratePerSecond.
func initThrottle(name string, burst int, ratePerSecond int) chan time.Time {
	tick := time.NewTicker(time.Second / time.Duration(ratePerSecond))
	throttle := make(chan time.Time, burst)
	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		} // does not exit after tick.Stop()
	}()
	return throttle
}

// Register adds the ADS and EDS handles to the grpc server
func (s *DiscoveryServer) Register(rpcs *grpc.Server) {
	// EDS must remain registered for 0.8, for smooth upgrade from 0.7
	// 0.7 proxies will use this service.
	ads.RegisterAggregatedDiscoveryServiceServer(rpcs, s)
}

// Singleton, refresh the cache - may not be needed if events work properly, just a failsafe
// ( will be removed after change detection is implemented, to double check all changes are
// captured)
func (s *DiscoveryServer) periodicRefresh() {
	envOverride := os.Getenv("V2_REFRESH")
	if len(envOverride) > 0 {
		var err error
		periodicRefreshDuration, err = time.ParseDuration(envOverride)
		if err != nil {
			adsLog.Warn("Invalid value for V2_REFRESH")
		}
	}
	if periodicRefreshDuration == 0 {
		return
	}
	ticker := time.NewTicker(periodicRefreshDuration)
	defer ticker.Stop()
	for range ticker.C {
		adsLog.Infof("ADS: periodic push of envoy configs %s", versionInfo())
		s.AdsPushAll(versionInfo(), s.Env.PushContext, true, nil)
	}
}

// Push metrics are updated periodically (10s default)
func (s *DiscoveryServer) periodicRefreshMetrics() {
	envOverride := os.Getenv("V2_METRICS")
	if len(envOverride) > 0 {
		var err error
		periodicRefreshMetrics, err = time.ParseDuration(envOverride)
		if err != nil {
			adsLog.Warn("Invalid value for V2_METRICS")
		}
	}
	if periodicRefreshMetrics == 0 {
		return
	}

	ticker := time.NewTicker(periodicRefreshMetrics)
	defer ticker.Stop()
	for range ticker.C {
		push := s.Env.PushContext
		if push.End != timeZero {
			model.LastPushStatus = push
		}
		push.UpdateMetrics()
		// TODO: env to customize
		//if time.Since(push.Start) > 30*time.Second {
		// Reset the stats, some errors may still be stale.
		//s.env.PushContext = model.NewPushContext()
		//}
	}
}

// Push is called to push changes on config updates using ADS. This is set in DiscoveryService.Push,
// to avoid direct dependencies.
func (s *DiscoveryServer) Push(full bool, edsUpdates map[string]*model.ServiceShards) {
	if !full {
		adsLog.Infof("EDS Incremental Push %v", edsUpdates)
		go s.AdsPushAll(version, s.Env.PushContext, false, edsUpdates)
		return
	}
	// Reset the status during the push.
	//afterPush := true
	if s.Env.PushContext != nil {
		s.Env.PushContext.OnConfigChange()
	}
	// PushContext is reset after a config change. Previous status is
	// saved.
	t0 := time.Now()
	push := model.NewPushContext()
	err := push.InitContext(s.Env)
	if err != nil {
		adsLog.Errorf("XDS: failed to update services %v", err)
		// We can't push if we can't read the data - stick with previous version.
		// TODO: metric !!
		// TODO: metric !!
		return
	}

	if err = s.ConfigGenerator.BuildSharedPushState(s.Env, push); err != nil {
		adsLog.Errorf("XDS: Failed to rebuild share state in configgen: %v", err)
		return
	}

	err = s.updateServiceShards(push)
	if err != nil {
		return
	}

	s.Env.PushContext = push
	versionLocal := time.Now().Format(time.RFC3339)
	initContextTime := time.Since(t0)
	adsLog.Debugf("InitContext %v for push took %s", versionLocal, initContextTime)

	// TODO: propagate K8S version and use it instead
	versionMutex.Lock()
	version = versionLocal
	versionMutex.Unlock()

	go s.AdsPushAll(versionLocal, push, true, nil)
}

func nonce() string {
	return time.Now().String()
}

func versionInfo() string {
	versionMutex.RLock()
	defer versionMutex.RUnlock()
	return version
}
