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

package health

import (
	"time"

	"istio.io/api/networking/v1alpha3"
)

const (
	HealthInfoTypeURL string = "type.googleapis.com/istio.v1.HealthInformation"
)

type WorkloadHealthChecker struct {
	config applicationHealthCheckConfig
	prober Prober
}

// internal field purely for convenience
type applicationHealthCheckConfig struct {
	InitialDelay   time.Duration
	ProbeTimeout   time.Duration
	CheckFrequency time.Duration
	SuccessThresh  int
	FailThresh     int
}

type ProbeEvent struct {
	Healthy          bool
	UnhealthyStatus  int32
	UnhealthyMessage string
}

func NewWorkloadHealthChecker(cfg v1alpha3.ReadinessProbe, prober Prober) *WorkloadHealthChecker {
	return &WorkloadHealthChecker{
		config: applicationHealthCheckConfig{
			InitialDelay:   time.Duration(cfg.InitialDelaySeconds) * time.Second,
			ProbeTimeout:   time.Duration(cfg.TimeoutSeconds) * time.Second,
			CheckFrequency: time.Duration(cfg.PeriodSeconds) * time.Second,
			SuccessThresh:  int(cfg.SuccessThreshold),
			FailThresh:     int(cfg.FailureThreshold),
		},
		prober: prober,
	}
}

// PerformApplicationHealthCheck Performs the application-provided configuration health check.
// Instead of a heartbeat-based health checks, we only send on a health state change, and this is
// determined by the success & failure threshold provided by the user.
func (w *WorkloadHealthChecker) PerformApplicationHealthCheck(notifyHealthChange chan *ProbeEvent, quit chan struct{}) {
	defer close(notifyHealthChange)
	// delay before starting probes.
	time.Sleep(w.config.InitialDelay)

	// tracks number of success & failures after last success/failure
	numSuccess, numFail := 0, 0
	// if the last send/event was a success, this is true, by default false because we want to
	// first send a healthy message.
	lastStateHealthy := false

	periodTicker := time.NewTicker(w.config.CheckFrequency)
	for {
		select {
		case <-quit:
			return
		case <-periodTicker.C:
			// probe target
			healthy, err := w.prober.Probe(w.config.ProbeTimeout)
			if err != nil {
				healthCheckLog.Errorf("Unexpected error probing: %v", err)
			}
			if healthy {
				// we were healthy, increment success counter
				numSuccess++
				// if we reached the threshold, mark the target as healthy
				if numSuccess == w.config.SuccessThresh && !lastStateHealthy {
					notifyHealthChange <- &ProbeEvent{Healthy: true}
					numSuccess = 0
					numFail = 0
					lastStateHealthy = true
				}
			} else {
				// we were not healthy, increment fail counter
				numFail++
				// if we reached the fail threshold, mark the target as unhealthy
				if numFail == w.config.FailThresh && lastStateHealthy {
					notifyHealthChange <- &ProbeEvent{
						Healthy:          false,
						UnhealthyStatus:  500,
						UnhealthyMessage: "unhealthy",
					}
					numSuccess = 0
					numFail = 0
					lastStateHealthy = false
				}
			}
		}
	}
}

// TODO implement
func (w *WorkloadHealthChecker) PerformEnvoyHealthCheck() {

}
