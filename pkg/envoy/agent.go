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

package envoy

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"istio.io/istio/pilot/pkg/util/sets"
	"istio.io/istio/pkg/http"
	"istio.io/pkg/log"
)

var errAbort = errors.New("epoch aborted")

const errOutOfMemory = "signal: killed"

// minDrainDuration is the minimum duration for which agent waits before it actively checks
// for outstanding connections. After this period, agent terminates proxy whenever active
// connections become zero.
var minDrainDuration = 5 * time.Second

var activeConnectionCheckDelay = 1 * time.Second

var knownIstioListeners = sets.NewSet(
	"listener.0.0.0.0_15006.downstream_cx_active", "listener.0.0.0.0_15021.downstream_cx_active",
	"listener.0.0.0.0_15090.downstream_cx_active", "listener.admin.downstream_cx_active",
	"listener.admin.main_thread.downstream_cx_active",
)

// NewAgent creates a new proxy agent for the proxy start-up and clean-up functions.
func NewAgent(proxy Proxy, terminationDrainDuration time.Duration, localhost string, adminPort uint16) *Agent {
	return &Agent{
		proxy:                    proxy,
		statusCh:                 make(chan exitStatus, 1), // context might stop drainage
		drainCh:                  make(chan struct{}),
		abortCh:                  make(chan error, 1),
		terminationDrainDuration: terminationDrainDuration,
		adminPort:                adminPort,
		localhost:                localhost,
	}
}

// Proxy defines command interface for a proxy
type Proxy interface {

	// Run command for an epoch, and abort channel
	Run(int, <-chan error) error

	// Drains the current epoch.
	Drain() error

	// Cleanup command for an epoch
	Cleanup(int)

	// UpdateConfig writes a new config file
	UpdateConfig(config []byte) error
}

type Agent struct {
	// proxy commands
	proxy Proxy

	// channel for proxy exit notifications
	statusCh chan exitStatus

	drainCh chan struct{}

	abortCh chan error

	// time to allow for the proxy to drain before terminating all remaining proxy processes
	terminationDrainDuration time.Duration

	adminPort uint16
	localhost string
}

type exitStatus struct {
	epoch int
	err   error
}

// Run starts the envoy and waits until it terminates.
func (a *Agent) Run(ctx context.Context) {
	log.Info("Starting proxy agent")
	go a.runWait(0, a.abortCh)

	select {
	case status := <-a.statusCh:
		if status.err != nil {
			if status.err.Error() == errOutOfMemory {
				log.Warnf("Envoy may have been out of memory killed. Check memory usage and limits.")
			}
			log.Errorf("Epoch %d exited with error: %v", status.epoch, status.err)
		} else {
			log.Infof("Epoch %d exited normally", status.epoch)
		}

		log.Infof("No more active epochs, terminating")
	case <-ctx.Done():
		a.terminate()
		status := <-a.statusCh
		if status.err == errAbort {
			log.Infof("Epoch %d aborted normally", status.epoch)
		} else {
			log.Warnf("Epoch %d aborted abnormally", status.epoch)
		}
		log.Info("Agent has successfully terminated")
	}
}

func (a *Agent) terminate() {
	log.Infof("Agent draining Proxy")
	e := a.proxy.Drain()
	if e != nil {
		log.Warnf("Error in invoking drain listeners endpoint %v", e)
	}
	time.Sleep(minDrainDuration)
	log.Infof("Termination drain duration period is %v, checking for active connections...", a.terminationDrainDuration)
	go a.waitForDrain()
	select {
	case <-a.drainCh:
		log.Info("There are no more active connections. terminating proxy...")
		a.abortCh <- errAbort
	// TODO: remove terminationDrainDuration and rely on "terminationGracefulPeriodSeconds" of pod?
	case <-time.After(a.terminationDrainDuration):
		log.Info("Termination period complete, terminating proxy...")
		a.abortCh <- errAbort
	}
	log.Warnf("Aborted all epochs")
}

func (a *Agent) waitForDrain() {
	activeConnectionDelay := time.NewTimer(activeConnectionCheckDelay)
	for {
		select {
		case <-activeConnectionDelay.C:
			if a.activeProxyConnections() == 0 {
				a.drainCh <- struct{}{}
				return
			} else {
				log.Info("active connections are still not zero")
			}
		case <-a.abortCh:
			log.Info("Abort channel in waitForDuration")
			stopped := activeConnectionDelay.Stop()
			log.Infof("Stopped timer %v", stopped)
			if !activeConnectionDelay.Stop() {
				// if the timer has been stopped then read from the channel.
				<-activeConnectionDelay.C
			}
			return
		}
	}
}

func (a *Agent) activeProxyConnections() int {
	activeConnectionsURL := fmt.Sprintf("http://%s:%d/stats?usedonly&filter=downstream_cx_active", a.localhost, a.adminPort)
	stats, err := http.DoHTTPGet(activeConnectionsURL)
	if err != nil {
		log.Warnf("Unable to get listener stats from Envoy : %v", err)
		return -1
	}
	activeConnections := 0
	for stats.Len() > 0 {
		line, _ := stats.ReadString('\n')
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Warnf("envoy stat line is missing separator. line:%s", line)
			continue
		}
		// downstream_cx_active is accounted under "http." and "listener." for http listeners.
		// Only consider listener stats.
		if !strings.HasPrefix(parts[0], "listener.") {
			continue
		}
		// If the stat is for a known Istio listener skip it.
		if knownIstioListeners.Contains(parts[0]) {
			continue
		}
		val, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			log.Warnf("failed parsing Envoy stat %s (error: %s) line: %s", parts[0], err.Error(), line)
			continue
		}
		activeConnections += int(val)
	}
	return activeConnections
}

// runWait runs the start-up command as a go routine and waits for it to finish
func (a *Agent) runWait(epoch int, abortCh <-chan error) {
	log.Infof("Epoch %d starting", epoch)
	err := a.proxy.Run(epoch, abortCh)
	a.proxy.Cleanup(epoch)
	a.statusCh <- exitStatus{epoch: epoch, err: err}
}
