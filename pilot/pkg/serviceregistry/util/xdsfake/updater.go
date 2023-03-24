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

package xdsfake

import (
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/cluster"
	"istio.io/istio/pkg/test"
	"istio.io/pkg/log"
)

// NewFakeXDS creates a XdsUpdater reporting events via a channel.
func NewFakeXDS() *Updater {
	return &Updater{
		Events: make(chan Event, 100),
	}
}

// NewWithDelegate creates a XdsUpdater reporting events via a channel.
func NewWithDelegate(delegate model.XDSUpdater) *Updater {
	return &Updater{
		Events:   make(chan Event, 100),
		Delegate: delegate,
	}
}

// Updater is used to test the registry.
type Updater struct {
	// Events tracks notifications received by the updater
	Events   chan Event
	Delegate model.XDSUpdater
}

var _ model.XDSUpdater = &Updater{}

func (fx *Updater) ConfigUpdate(req *model.PushRequest) {
	names := []string{}
	if req != nil && len(req.ConfigsUpdated) > 0 {
		for key := range req.ConfigsUpdated {
			names = append(names, key.Name)
		}
	}
	sort.Strings(names)
	id := strings.Join(names, ",")
	event := "xds"
	if req.Full {
		event += " full"
	}
	select {
	case fx.Events <- Event{Type: event, ID: id}:
	default:
	}
	if fx.Delegate != nil {
		fx.Delegate.ConfigUpdate(req)
	}
}

func (fx *Updater) ProxyUpdate(c cluster.ID, ip string) {
	select {
	case fx.Events <- Event{Type: "proxy", ID: ip}:
	default:
	}
	if fx.Delegate != nil {
		fx.Delegate.ProxyUpdate(c, ip)
	}
}

// Event is used to watch XdsEvents
type Event struct {
	// Type of the event
	Type string

	// The id of the event
	ID string

	Namespace string

	// The endpoints associated with an EDS push if any
	Endpoints []*model.IstioEndpoint

	// EndpointCount, used in matches only
	EndpointCount int
}

func (fx *Updater) EDSUpdate(c model.ShardKey, hostname string, ns string, entry []*model.IstioEndpoint) {
	select {
	case fx.Events <- Event{Type: "eds", ID: hostname, Endpoints: entry, Namespace: ns}:
	default:
	}
	if fx.Delegate != nil {
		fx.Delegate.EDSUpdate(c, hostname, ns, entry)
	}
}

func (fx *Updater) EDSCacheUpdate(c model.ShardKey, hostname, ns string, entry []*model.IstioEndpoint) {
	select {
	case fx.Events <- Event{Type: "eds cache", ID: hostname, Endpoints: entry, Namespace: ns}:
	default:
	}
	if fx.Delegate != nil {
		fx.Delegate.EDSCacheUpdate(c, hostname, ns, entry)
	}
}

// SvcUpdate is called when a service port mapping definition is updated.
// This interface is WIP - labels, annotations and other changes to service may be
// updated to force a EDS and CDS recomputation and incremental push, as it doesn't affect
// LDS/RDS.
func (fx *Updater) SvcUpdate(c model.ShardKey, hostname string, ns string, ev model.Event) {
	select {
	case fx.Events <- Event{Type: "service", ID: hostname, Namespace: ns}:
	default:
	}
	if fx.Delegate != nil {
		fx.Delegate.SvcUpdate(c, hostname, ns, ev)
	}
}

func (fx *Updater) RemoveShard(shardKey model.ShardKey) {
	select {
	case fx.Events <- Event{Type: "removeShard", ID: shardKey.String()}:
	default:
	}
	if fx.Delegate != nil {
		fx.Delegate.RemoveShard(shardKey)
	}
}

func (fx *Updater) WaitOrFail(t test.Failer, et string) *Event {
	t.Helper()
	for {
		select {
		case e := <-fx.Events:
			if e.Type == et {
				return &e
			}
			log.Infof("skipping event %q want %q", e.Type, et)
			continue
		case <-time.After(time.Second * 5):
			t.Fatalf("timed out waiting for %v", et)
		}
	}
}

func (fx *Updater) MatchOrFail(t test.Failer, events ...Event) {
	t.Helper()

	for {
		if len(events) == 0 {
			return
		}
		select {
		case e := <-fx.Events:
			found := false
			for i, want := range events {
				if e.Type == want.Type &&
					(want.ID == "" || e.ID == want.ID) &&
					(want.Namespace == "" || want.Namespace == e.Namespace) &&
					(want.EndpointCount == 0 || want.EndpointCount == len(e.Endpoints)) {
					// Matched - delete event from desired
					events = slices.Delete(events, i, i+1)
					found = true
					break
				}
			}
			if !found {
				log.Infof("skipping event %q/%v", e.Type, e.ID)
			}
			continue
		case <-time.After(time.Second * 5):
			t.Fatalf("timed out waiting for %v", events)
		}
	}
}

// Clear any pending event
func (fx *Updater) Clear() {
	wait := true
	for wait {
		select {
		case <-fx.Events:
		default:
			wait = false
		}
	}
}

// AssertEmpty ensures there are no events in the channel
func (fx *Updater) AssertEmpty(t test.Failer, dur time.Duration) {
	if dur == 0 {
		select {
		case e := <-fx.Events:
			t.Fatalf("got unexpected event %+v", e)
		default:
		}
	} else {
		select {
		case e := <-fx.Events:
			t.Fatalf("got unexpected event %+v", e)
		case <-time.After(dur):
		}
	}
}
