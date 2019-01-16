// Copyright 2019 Istio Authors
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

package source

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"

	mcp "istio.io/api/mcp/v1alpha1"
	"istio.io/istio/pkg/log"
	"istio.io/istio/pkg/mcp/internal"
	"istio.io/istio/pkg/mcp/monitoring"
)

var scope = log.RegisterScope("mcp", "mcp debugging", 0)

// Request is a temporary abstraction for the MCP node request which can
// be used with the mcp.MeshConfigRequest and mcp.RequestResources. It can
// be removed once we fully cutover to mcp.RequestResources.
type Request struct {
	Collection string

	// Most recent version was that ACK/NACK'd by the sink
	VersionInfo string
	SinkNode    *mcp.SinkNode

	// hidden
	incremental bool
}

// WatchResponse contains a versioned collection of pre-serialized resources.
type WatchResponse struct {
	Collection string

	// Version of the resources in the response for the given
	// type. The node responses with this version in subsequent
	// requests as an acknowledgment.
	Version string

	// Resourced resources to be included in the response.
	Resources []*mcp.Resource

	// The original request for triggered this response
	Request *Request
}

type (
	// CancelWatchFunc allows the consumer to cancel a previous watch,
	// terminating the watch for the request.
	CancelWatchFunc func()

	// PushResponseFunc allows the consumer to push a response for the
	// corresponding watch.
	PushResponseFunc func(*WatchResponse)
)

// Watcher requests watches for configuration resources by node, last
// applied version, and type. The watch should send the responses when
// they are ready. The watch can be canceled by the consumer.
type Watcher interface {
	// Watch returns a new open watch for a non-empty request.
	//
	// Cancel is an optional function to release resources in the
	// producer. It can be called idempotently to cancel and release resources.
	Watch(*Request, PushResponseFunc) CancelWatchFunc
}

// CollectionOptions configures the per-collection updates.
type CollectionOptions struct {
	// Name of the collection, e.g. istio/networking/v1alpha3/VirtualService
	Name string
}

// CollectionOptionsFromSlice returns a slice of collection options from
// a slice of collection names.
func CollectionOptionsFromSlice(names []string) []CollectionOptions {
	options := make([]CollectionOptions, 0, len(names))
	for _, name := range names {
		options = append(options, CollectionOptions{
			Name: name,
		})
	}
	return options
}

// Options contains options for configuring MCP sources.
type Options struct {
	Watcher            Watcher
	CollectionsOptions []CollectionOptions
	Reporter           monitoring.Reporter

	// Controls the delay for re-retrying a configuration push if the previous
	// attempt was not possible, e.errgrp. the lower-level serving layer was busy. This
	// should typically be set fairly small (order of milliseconds).
	RetryPushDelay time.Duration
}

// Stream is for sending Resource messages and receiving RequestResources messages.
type Stream interface {
	Send(*mcp.Resources) error
	Recv() (*mcp.RequestResources, error)
	Context() context.Context
}

// Sources implements the resource source message exchange for MCP. It can be instantiated by client and server
// source implementations to manage the MCP message exchange.
type Source struct {
	watcher        Watcher
	collections    []CollectionOptions
	nextStreamID   int64
	reporter       monitoring.Reporter
	connections    int64
	retryPushDelay time.Duration
}

type newPushResponseState int

const (
	newPushResponseStateReady newPushResponseState = iota
	newPushResponseStateClosed
)

// watch maintains local push state of the most recent watch per-type.
type watch struct {
	// only accessed from connection goroutine
	cancel          func()
	ackedVersionMap map[string]string // resources that exist at the sink; by name and version
	pending         *mcp.Resources
	retryPushDelay  time.Duration
	incremental     bool

	// NOTE: do not hold `mu` when reading/writing to this channel.
	newPushResponseReadyChan chan newPushResponseState

	mu              sync.Mutex
	newPushResponse *WatchResponse
	timer           *time.Timer
	closed          bool
}

func (w *watch) delayedPush() {
	w.mu.Lock()
	w.timer = nil
	w.mu.Unlock()

	select {
	case w.newPushResponseReadyChan <- newPushResponseStateReady:
	default:
		time.AfterFunc(w.retryPushDelay, w.schedulePush)
	}
}

// Try to schedule pushing a response to the node. The push may
// be re-scheduled as needed. Additional care is taken to rate limit
// re-pushing responses that were previously NACK'd. This avoid flooding
// the node with responses while also allowing transient NACK'd responses
// to be retried.
func (w *watch) schedulePush() {
	w.mu.Lock()

	// close the watch
	if w.closed {
		// unlock before channel write
		w.mu.Unlock()

		select {
		case w.newPushResponseReadyChan <- newPushResponseStateClosed:
		default:
			time.AfterFunc(w.retryPushDelay, w.schedulePush)
		}
		return
	}

	// no-op if the response has already been sent
	if w.newPushResponse == nil {
		w.mu.Unlock()
		return
	}

	// Otherwise, try to schedule the response to be sent.
	if w.timer != nil {
		if !w.timer.Stop() {
			<-w.timer.C
		}
		w.timer = nil
	}
	// unlock before channel write
	w.mu.Unlock()
	select {
	case w.newPushResponseReadyChan <- newPushResponseStateReady:
	default:
		time.AfterFunc(w.retryPushDelay, w.schedulePush)
	}

}

// Save the pushed response in the newPushResponse and schedule a push. The push
// may be re-schedule as necessary but this should be transparent to the
// caller. The caller may provide a nil response to indicate that the watch
// should be closed.
func (w *watch) saveResponseAndSchedulePush(response *WatchResponse) {
	w.mu.Lock()
	w.newPushResponse = response
	if response == nil {
		w.closed = true
	}
	w.mu.Unlock()

	w.schedulePush()
}

// connection maintains per-stream connection state for a
// node. Access to the stream and watch state is serialized
// through request and response channels.
type connection struct {
	peerAddr string
	stream   Stream
	id       int64

	// unique nonce generator for req-resp pairs per xDS stream; the server
	// ignores stale nonces. nonce is only modified within send() function.
	streamNonce int64

	requestC chan *mcp.RequestResources // a channel for receiving incoming requests
	reqError error                      // holds error if request channel is closed
	watches  map[string]*watch          // per-type watches
	watcher  Watcher

	reporter monitoring.Reporter
}

const DefaultRetryPushDelay = 10 * time.Millisecond

// New creates a new resource source.
func New(options *Options) *Source {
	if options.RetryPushDelay == 0 {
		options.RetryPushDelay = DefaultRetryPushDelay
	}
	s := &Source{
		watcher:        options.Watcher,
		collections:    options.CollectionsOptions,
		reporter:       options.Reporter,
		retryPushDelay: options.RetryPushDelay,
	}
	return s
}

func (s *Source) newConnection(stream Stream) *connection {
	peerAddr := "0.0.0.0"

	peerInfo, ok := peer.FromContext(stream.Context())
	if ok {
		peerAddr = peerInfo.Addr.String()
	} else {
		scope.Warnf("No peer info found on the incoming stream.")
		peerInfo = nil
	}

	con := &connection{
		stream:   stream,
		peerAddr: peerAddr,
		requestC: make(chan *mcp.RequestResources),
		watches:  make(map[string]*watch),
		watcher:  s.watcher,
		id:       atomic.AddInt64(&s.nextStreamID, 1),
		reporter: s.reporter,
	}

	var collections []string
	for _, collection := range s.collections {
		w := &watch{
			newPushResponseReadyChan: make(chan newPushResponseState, 1),
			ackedVersionMap:          make(map[string]string),
			retryPushDelay:           s.retryPushDelay,
			incremental:              false,
		}
		con.watches[collection.Name] = w
		collections = append(collections, collection.Name)
	}

	s.reporter.SetStreamCount(atomic.AddInt64(&s.connections, 1))

	scope.Debugf("MCP: connection %v: NEW, supported collections: %#v", con, collections)

	return con
}

func (s *Source) processStream(stream Stream) error {
	con := s.newConnection(stream)

	defer s.closeConnection(con)
	go con.receive()

	// fan-in per-type response channels into single response channel for the select loop below.
	responseChan := make(chan *watch, 1)
	for _, w := range con.watches {
		go func(w *watch) {
			for state := range w.newPushResponseReadyChan {
				if state == newPushResponseStateClosed {
					break
				}
				responseChan <- w
			}

			// Any closed watch can close the overall connection. Use
			// `nil` value to indicate a closed state to the run loop
			// below instead of closing the channel to avoid closing
			// the channel multiple times.
			responseChan <- nil
		}(w)
	}

	for {
		select {
		case w, more := <-responseChan:
			if !more || w == nil {
				return status.Error(codes.Unavailable, "source canceled watch")
			}

			w.mu.Lock()
			resp := w.newPushResponse
			w.newPushResponse = nil
			w.mu.Unlock()

			// newPushResponse may have been cleared before we got to it
			if resp == nil {
				break
			}
			if err := con.pushServerResponse(w, resp); err != nil {
				return err
			}
		case req, more := <-con.requestC:
			if !more {
				return con.reqError
			}
			if err := con.processClientRequest(req); err != nil {
				return err
			}
		case <-stream.Context().Done():
			scope.Debugf("MCP: connection %v: stream done, err=%v", con, stream.Context().Err())
			return stream.Context().Err()
		}
	}
}

func (s *Source) closeConnection(con *connection) {
	con.close()
	s.reporter.SetStreamCount(atomic.AddInt64(&s.connections, -1))
}

// String implements Stringer.String.
func (con *connection) String() string {
	return fmt.Sprintf("{addr=%v id=%v}", con.peerAddr, con.id)
}

func calculateDelta(current []*mcp.Resource, acked map[string]string) (added []mcp.Resource, removed []string) {
	// TODO - consider storing desired state as a map to make this faster
	desired := make(map[string]*mcp.Resource, len(current))

	// compute diff
	for _, envelope := range current {
		prevVersion, exists := acked[envelope.Metadata.Name]
		if !exists {
			// new
			added = append(added, *envelope)
		} else if prevVersion != envelope.Metadata.Version {
			// update
			added = append(added, *envelope)
		}
		// tracking for delete
		desired[envelope.Metadata.Name] = envelope
	}

	for name := range acked {
		if _, exists := desired[name]; !exists {
			removed = append(removed, name)
		}
	}

	return added, removed
}

func (con *connection) pushServerResponse(w *watch, resp *WatchResponse) error {
	var (
		added   []mcp.Resource
		removed []string
	)

	// send an incremental update if enabled for this collection and the most
	// recent request from the sink requested it.
	var incremental bool
	if w.incremental && resp.Request.incremental {
		incremental = true
	}

	if incremental {
		added, removed = calculateDelta(resp.Resources, w.ackedVersionMap)
	} else {
		for _, resource := range resp.Resources {
			added = append(added, *resource)
		}
	}

	msg := &mcp.Resources{
		SystemVersionInfo: resp.Version,
		Collection:        resp.Collection,
		Resources:         added,
		RemovedResources:  removed,
		Incremental:       incremental,
	}

	// increment nonce
	con.streamNonce = con.streamNonce + 1
	msg.Nonce = strconv.FormatInt(con.streamNonce, 10)
	if err := con.stream.Send(msg); err != nil {
		con.reporter.RecordSendError(err, status.Code(err))

		return err
	}
	scope.Debugf("MCP: connection %v: SEND collection=%v version=%v nonce=%v",
		con, resp.Collection, resp.Version, msg.Nonce)
	w.pending = msg
	return nil
}

func (con *connection) receive() {
	defer close(con.requestC)
	for {
		req, err := con.stream.Recv()
		if err != nil {
			code := status.Code(err)
			if code == codes.Canceled || err == io.EOF {
				scope.Infof("MCP: connection %v: TERMINATED %q", con, err)
				return
			}
			con.reporter.RecordRecvError(err, code)
			scope.Errorf("MCP: connection %v: TERMINATED with errors: %v", con, err)
			// Save the stream error prior to closing the stream. The caller
			// should access the error after the channel closure.
			con.reqError = err
			return
		}
		select {
		case con.requestC <- req:
		case <-con.stream.Context().Done():
			scope.Debugf("MCP: connection %v: stream done, err=%v", con, con.stream.Context().Err())
			return
		}
	}
}

func (con *connection) close() {
	scope.Infof("MCP: connection %v: CLOSED", con)

	for _, w := range con.watches {
		if w.cancel != nil {
			w.cancel()
		}
	}
}

func (con *connection) processClientRequest(req *mcp.RequestResources) error {
	collection := req.Collection

	con.reporter.RecordRequestSize(collection, con.id, req.Size())

	w, ok := con.watches[collection]
	if !ok {
		return status.Errorf(codes.InvalidArgument, "unsupported collection %q", collection)
	}

	// nonces can be reused across streams; we verify nonce only if it initialized
	if req.ResponseNonce == "" || w.pending.GetNonce() == req.ResponseNonce {
		versionInfo := ""

		if w.pending == nil {
			scope.Infof("MCP: connection %v: WATCH for %v", con, collection)
		} else {
			versionInfo = w.pending.SystemVersionInfo
			if req.ErrorDetail != nil {
				scope.Warnf("MCP: connection %v: NACK collection=%v version=%q with nonce=%q error=%#v", // nolint: lll
					con, collection, req.ResponseNonce, versionInfo, req.ErrorDetail)
				con.reporter.RecordRequestNack(collection, con.id, codes.Code(req.ErrorDetail.Code))
			} else {
				scope.Infof("MCP: connection %v ACK collection=%v with version=%q nonce=%q",
					con, collection, versionInfo, req.ResponseNonce)
				con.reporter.RecordRequestAck(collection, con.id)

				internal.UpdateResourceVersionTracking(w.ackedVersionMap, w.pending)
			}

			// clear the pending request after we finished processing the corresponding response.
			w.pending = nil
		}

		if w.cancel != nil {
			w.cancel()
		}

		sr := &Request{
			SinkNode:    req.SinkNode,
			Collection:  collection,
			VersionInfo: versionInfo,
			incremental: req.Incremental,
		}
		w.cancel = con.watcher.Watch(sr, w.saveResponseAndSchedulePush)
	} else {
		// This error path should not happen! Skip any requests that don't match the
		// latest watch's nonce. These could be dup requests or out-of-order
		// requests from a buggy node.
		if req.ErrorDetail != nil {
			scope.Errorf("MCP: connection %v: STALE NACK collection=%v with nonce=%q (expected nonce=%q) error=%+v", // nolint: lll
				con, collection, req.ResponseNonce, w.pending.GetNonce(), req.ErrorDetail)
			con.reporter.RecordRequestNack(collection, con.id, codes.Code(req.ErrorDetail.Code))
		} else {
			scope.Errorf("MCP: connection %v: STALE ACK collection=%v with nonce=%q (expected nonce=%q)", // nolint: lll
				con, collection, req.ResponseNonce, w.pending.GetNonce())
			con.reporter.RecordRequestAck(collection, con.id)
		}
	}

	return nil
}
