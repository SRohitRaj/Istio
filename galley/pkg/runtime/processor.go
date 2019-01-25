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
	"time"

	"github.com/pkg/errors"

	"istio.io/istio/galley/pkg/metadata"
	"istio.io/istio/galley/pkg/runtime/monitoring"
	"istio.io/istio/galley/pkg/runtime/processing"
	"istio.io/istio/galley/pkg/runtime/publish"
	"istio.io/istio/galley/pkg/runtime/resource"
	"istio.io/istio/galley/pkg/runtime/states/defaultstate"
	"istio.io/istio/pkg/log"
	sn "istio.io/istio/pkg/mcp/snapshot"
)

var scope = log.RegisterScope("runtime", "Galley runtime", 0)

// Processor is the main control-loop for processing incoming config events and translating them into
// component configuration
type Processor struct {
	// source interface for retrieving the events from.
	source Source

	// events channel that was obtained from source
	events chan resource.Event

	// handler for events.
	handler processing.Handler

	// channel that gets closed during Shutdown.
	done chan struct{}

	// channel that signals the background process as being stopped.
	stopped chan struct{}

	// The current in-memory configuration State
	state *defaultstate.State

	// hook that gets called after each event processing. Useful for testing.
	postProcessHook postProcessHookFn

	// lastEventTime records the last time an event was received.
	lastEventTime time.Time
}

type postProcessHookFn func()

// NewProcessor returns a new instance of a Processor
func NewProcessor(src Source, distributor publish.Distributor, cfg *Config) *Processor {
	state := defaultstate.New(sn.DefaultGroup, cfg.DomainSuffix, metadata.Types, publish.NewStrategyWithDefaults(), distributor)
	return newProcessor(state, src, nil)
}

func newProcessor(
	state *defaultstate.State,
	src Source,
	postProcessHook postProcessHookFn) *Processor {

	now := time.Now()
	return &Processor{
		handler:         buildDispatcher(state),
		state:           state,
		source:          src,
		postProcessHook: postProcessHook,
		done:            make(chan struct{}),
		stopped:         make(chan struct{}),
		lastEventTime:   now,
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
}

func (p *Processor) process() {
	scope.Debugf("Starting process loop")

loop:
	for {
		select {

		// Incoming events are received through p.events
		case e := <-p.events:
			p.processEvent(e)

		case <-p.state.PublishChan():
			scope.Debug("Processor.process: publish")
			p.state.Publish()

		// p.done signals the graceful Shutdown of the processor.
		case <-p.done:
			scope.Debug("Processor.process: done")
			break loop
		}

		if p.postProcessHook != nil {
			p.postProcessHook()
		}
	}

	p.state.Close()
	close(p.stopped)
	scope.Debugf("Process.process: Exiting process loop")
}

func (p *Processor) processEvent(e resource.Event) {
	scope.Debugf("Incoming source event: %v", e)
	p.recordEvent()

	if e.Kind == resource.FullSync {
		scope.Infof("Synchronization is complete, starting distribution.")
		p.state.OnFullSync()
		return
	}

	p.handler.Handle(e)
}

func (p *Processor) recordEvent() {
	now := time.Now()
	monitoring.RecordProcessorEventProcessed(now.Sub(p.lastEventTime))
	p.lastEventTime = now
}

type schemaBasedHandler interface {
	processing.Handler
	GetSchema() *resource.Schema
}

func buildDispatcher(handlers ...schemaBasedHandler) *processing.Dispatcher {
	b := processing.NewDispatcherBuilder()
	for _, h := range handlers {
		for _, spec := range h.GetSchema().All() {
			b.Add(spec.Collection, h)
		}
	}
	return b.Build()
}
