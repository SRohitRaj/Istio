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

package status

import (
	"fmt"
	"net"
	"net/http"

	"context"
	"os"
	"sync"

	"time"

	"istio.io/istio/pilot/cmd/pilot-agent/status/ready"
	"istio.io/istio/pkg/log"
)

const (
	// readyPath is for the pilot agent readiness itself.
	readyPath = "/healthz/ready"
	// appReadyPath is the path for the application after injecting.
	appReadinessPath = "/istio-take-over/readiness"
	// appHealthPath is the path for the application after injecting.
	appLivenessPath = "/istio-take-over/liveness"
)

// AppProbeInfo defines the information for Pilot agent to take over application probing.
type AppProbeInfo struct {
	Path string
	Port uint16
}

// Config for the status server.
type Config struct {
	StatusPort       uint16
	AdminPort        uint16
	ApplicationPorts []uint16
	// AppReadiness specifies how to take over Kubernetes readiness probing.
	AppReadiness *AppProbeInfo
	// AppLiveness specifies how to take over Kubernetes liveness probing.
	AppLiveness *AppProbeInfo
}

// Server provides an endpoint for handling status probes.
type Server struct {
	statusPort          uint16
	ready               *ready.Probe
	appLiveness         *AppProbeInfo
	appReadiness        *AppProbeInfo
	mutex               sync.Mutex
	lastProbeSuccessful bool
}

// NewServer creates a new status server.
func NewServer(config Config) *Server {
	return &Server{
		statusPort:   config.StatusPort,
		appLiveness:  config.AppLiveness,
		appReadiness: config.AppReadiness,
		ready: &ready.Probe{
			AdminPort:        config.AdminPort,
			ApplicationPorts: config.ApplicationPorts,
		},
	}
}

// Run opens a the status port and begins accepting probes.
func (s *Server) Run(ctx context.Context) {
	log.Infof("Opening status port %d\n", s.statusPort)

	// Add the handler for ready probes.
	http.HandleFunc(readyPath, s.handleReadyProbe)

	// TODO: how to differentiate whether app defined readiness probe or not?
	// Is it possible to not specify port in httpGet probeness?
	// Maybe use Port = -1 to differentiate and kube-injector modifies that accordingly.
	log.Infof("Pilot agent takes over readiness probe, path %v, port %v",
		s.appReadiness.Path, s.appReadiness.Port)
	http.HandleFunc(appReadinessPath, s.handleAppReadinessProbe)

	log.Infof("Pilot agent takes over liveness probe, path %v, port %v",
		s.appLiveness.Path, s.appLiveness.Port)
	http.HandleFunc(appLivenessPath, s.handleAppLivenessProbe)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.statusPort))
	if err != nil {
		log.Errorf("Error listening on status port: %v", err.Error())
		return
	}
	defer l.Close()

	go func() {
		if err := http.Serve(l, nil); err != nil {
			log.Errora(err)
			os.Exit(-1)
		}
	}()

	// Wait for the agent to be shut down.
	<-ctx.Done()
}

func (s *Server) handleReadyProbe(w http.ResponseWriter, _ *http.Request) {
	err := s.ready.Check()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)

		log.Infof("Envoy proxy is NOT ready: %s", err.Error())
		s.lastProbeSuccessful = false
	} else {
		w.WriteHeader(http.StatusOK)

		if !s.lastProbeSuccessful {
			log.Info("Envoy proxy is ready")
		}
		s.lastProbeSuccessful = true
	}
}

func (s *Server) handleAppReadinessProbe(w http.ResponseWriter, req *http.Request) {
	requestStatusCode(fmt.Sprintf("http://127.0.0.1:%d%s", s.appReadiness.Port, s.appReadiness.Path), w, req)
}

func (s *Server) handleAppLivenessProbe(w http.ResponseWriter, req *http.Request) {
	requestStatusCode(fmt.Sprintf("http://127.0.0.1:%d%s", s.appLiveness.Port, s.appLiveness.Path), w, req)
}

func requestStatusCode(appURL string, w http.ResponseWriter, req *http.Request) {
	log.Infof("Send request for probing application, url = %v", appURL)
	httpClient := &http.Client{
		// TODO: figure out the appropriate timeout?
		Timeout: 10 * time.Second,
	}

	appReq, err := http.NewRequest(req.Method, appURL, req.Body)
	for key, value := range req.Header {
		appReq.Header[key] = value
	}

	if err != nil {
		log.Errorf("Failed to copy request to probe app %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := httpClient.Do(appReq)
	if err != nil {
		log.Errorf("Request to probe app failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// We only write the status code to the response.
	w.WriteHeader(response.StatusCode)
}
