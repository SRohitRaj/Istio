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

package runtime

import (
	"github.com/pkg/errors"

	"time"

	"istio.io/istio/galley/pkg/metadata"
	"istio.io/istio/galley/pkg/runtime/resource"
	"istio.io/istio/pkg/log"
	"istio.io/istio/pkg/mcp/snapshot"
)

var scope = log.RegisterScope("runtime", "Galley runtime", 0)

// Processor is the main control-loop for processing incoming config events and translating them into
// component configuration
type Processor struct {
	// source interface for retrieving the events from.
	source Source

	// distributor interface for publishing config snapshots to.
	distributor Distributor

	// events channel that was obtained from source
	events chan resource.Event

	// channel that gets closed during Shutdown.
	done chan struct{}

	// channel that signals the background process as being stopped.
	stopped chan struct{}

	// The current in-memory configuration state
	state *state

	// hook that gets called after each event processing. Useful for testing.
	postProcessHook postProcessHookFn

	// pendingEvents counts the number of events awaiting publishing.
	pendingEvents int64

	// lastEventTime records the last time an event was received.
	lastEventTime time.Time

	// lastSnapshotTime records the last time a snapshot was published.
	lastSnapshotTime time.Time
}

type postProcessHookFn func()

// NewProcessor returns a new instance of a Processor
func NewProcessor(src Source, distributor Distributor, cfg *Config) *Processor {
	state := newState(metadata.Types, cfg, newPublishingStrategyWithDefaults())
	return newProcessor(state, src, distributor, nil)
}

func newProcessor(
	state *state,
	src Source,
	distributor Distributor,
	postProcessHook postProcessHookFn) *Processor {

	return &Processor{
		state:            state,
		source:           src,
		distributor:      distributor,
		postProcessHook:  postProcessHook,
		done:             make(chan struct{}),
		stopped:          make(chan struct{}),
		lastEventTime:    time.Now(),
		lastSnapshotTime: time.Now(),
	}
}

// Start the processor. This will cause processor to listen to incoming events from the provider
// and publish component configuration via the Distributor.
func (p *Processor) Start() error {
	scope.Info("Starting processor...")

	if p.events != nil {
		scope.Warn("Processor has already started")
		return errors.New("already started")
	}

	events := make(chan resource.Event, 1024)
	err := p.source.Start(func(e resource.Event) {
		events <- e
	})
	if err != nil {
		scope.Warnf("Unable to Start source: %v", err)
		return err
	}

	p.events = events

	go p.process()

	return nil
}

// Stop the processor.
func (p *Processor) Stop() {
	scope.Info("Stopping processor...")

	if p.events == nil {
		scope.Warnf("Processor has already stopped")
		return
	}

	p.source.Stop()

	close(p.done)
	<-p.stopped
	close(p.events)

	p.events = nil
	p.done = nil
	p.state = nil
}

func (p *Processor) process() {
	scope.Debugf("Starting process loop")

loop:
	for {
		select {

		// Incoming events are received through p.events
		case e := <-p.events:
			_ = p.processEvent(e)

		case <-p.state.strategy.publish:
			scope.Debug("Processor.process: publish")
			p.publish()

		// p.done signals the graceful Shutdown of the processor.
		case <-p.done:
			scope.Debug("Processor.process: done")
			break loop
		}

		if p.postProcessHook != nil {
			p.postProcessHook()
		}
	}

	p.state.close()
	close(p.stopped)
	scope.Debugf("Process.process: Exiting process loop")
}

func (p *Processor) processEvent(e resource.Event) bool {
	scope.Debugf("Incoming source event: %v", e)
	now := time.Now()
	recordProcessorEventProcessed(now.Sub(p.lastEventTime))
	p.lastEventTime = now
	p.pendingEvents++

	return p.state.apply(e)
}

func (p *Processor) publish() {
	now := time.Now()
	recordProcessorSnapshotPublished(p.pendingEvents, now.Sub(p.lastSnapshotTime))
	p.lastSnapshotTime = now
	sn := p.state.buildSnapshot()

	p.distributor.SetSnapshot(snapshot.DefaultGroup, sn)
	p.pendingEvents = 0
}
