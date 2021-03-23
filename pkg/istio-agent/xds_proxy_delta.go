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

package istioagent

import (
	"context"
	"time"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	v3 "istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pkg/istio-agent/metrics"
)

// requests from envoy
// for aditya:
// downstream -> envoy (anything "behind" xds proxy)
// upstream -> istiod (in front of xds proxy)?
func (p *XdsProxy) DeltaAggregatedResources(downstream discovery.AggregatedDiscoveryService_DeltaAggregatedResourcesServer) error {
	proxyLog.Debugf("accepted delta xds connection from envoy, forwarding to upstream")
	con := &ProxyConnection{
		upstreamError:      make(chan error, 2), // can be produced by recv and send
		downstreamError:    make(chan error, 2), // can be produced by recv and send
		deltaRequestsChan:  make(chan *discovery.DeltaDiscoveryRequest, 10),
		deltaResponsesChan: make(chan *discovery.DeltaDiscoveryResponse, 10),
		stopChan:           make(chan struct{}),
		downstreamDeltas:   downstream,
	}
	p.RegisterStream(con)
	defer p.UnregisterStream(con)

	// Handle downstream xds
	initialRequestsSent := false
	go func() {
		// Send initial request
		p.connectedMutex.RLock()
		initialRequest := p.initialDeltaRequest
		p.connectedMutex.RUnlock()

		for {
			// From Envoy
			req, err := downstream.Recv()
			if err != nil {
				con.downstreamError <- err
				return
			}
			// forward to istiod
			con.deltaRequestsChan <- req
			if !initialRequestsSent && req.TypeUrl == v3.ListenerType {
				// fire off an initial NDS request
				if _, f := p.handlers[v3.NameTableType]; f {
					con.deltaRequestsChan <- &discovery.DeltaDiscoveryRequest{
						TypeUrl: v3.NameTableType,
					}
				}
				// fire off an initial PCDS request
				if _, f := p.handlers[v3.ProxyConfigType]; f {
					con.deltaRequestsChan <- &discovery.DeltaDiscoveryRequest{
						TypeUrl: v3.ProxyConfigType,
					}
				}
				// Fire of a configured initial request, if there is one
				if initialRequest != nil {
					con.deltaRequestsChan <- initialRequest
				}
				initialRequestsSent = true
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	upstreamConn, err := grpc.DialContext(ctx, p.istiodAddress, p.istiodDialOptions...)
	if err != nil {
		proxyLog.Errorf("failed to connect to upstream %s: %v", p.istiodAddress, err)
		metrics.IstiodConnectionFailures.Increment()
		return err
	}
	defer upstreamConn.Close()

	xds := discovery.NewAggregatedDiscoveryServiceClient(upstreamConn)
	ctx = metadata.AppendToOutgoingContext(context.Background(), "ClusterID", p.clusterID)
	for k, v := range p.xdsHeaders {
		ctx = metadata.AppendToOutgoingContext(ctx, k, v)
	}
	// We must propagate upstream termination to Envoy. This ensures that we resume the full XDS sequence on new connection
	return p.HandleUpstream(ctx, con, xds)
}

func (p *XdsProxy) HandleDeltaUpstream(ctx context.Context, con *ProxyConnection, xds discovery.AggregatedDiscoveryServiceClient) error {
	deltaUpstream, err := xds.DeltaAggregatedResources(ctx, grpc.MaxCallRecvMsgSize(defaultClientMaxReceiveMessageSize))
	if err != nil {
		proxyLog.Debugf("failed to create delta upstream grpc client: %v", err)
		return err
	}
	proxyLog.Infof("connected to delta upstream XDS server: %s", p.istiodAddress)
	defer proxyLog.Debugf("disconnected from delta XDS server: %s", p.istiodAddress)

	con.upstreamDeltas = deltaUpstream

	// handle responses from istiod
	go func() {
		for {
			resp, err := deltaUpstream.Recv()
			if err != nil {
				con.upstreamError <- err
				return
			}
			con.deltaResponsesChan <- resp
		}
	}()

	go p.handleUpstreamDeltaRequest(ctx, con)
	go p.handleUpstreamDeltaResponse(con)

	// todo wasm load conversion
	for {
		select {
		case err := <-con.upstreamError:
			// error from upstream Istiod.
			if isExpectedGRPCError(err) {
				proxyLog.Debugf("upstream terminated with status %v", err)
				metrics.IstiodConnectionCancellations.Increment()
			} else {
				proxyLog.Warnf("upstream terminated with unexpected error %v", err)
				metrics.IstiodConnectionErrors.Increment()
			}
			return nil
		case err := <-con.downstreamError:
			// error from downstream Envoy.
			if isExpectedGRPCError(err) {
				proxyLog.Debugf("downstream terminated with status %v", err)
				metrics.EnvoyConnectionCancellations.Increment()
			} else {
				proxyLog.Warnf("downstream terminated with unexpected error %v", err)
				metrics.EnvoyConnectionErrors.Increment()
			}
			// On downstream error, we will return. This propagates the error to downstream envoy which will trigger reconnect
			return err
		case <-con.stopChan:
			proxyLog.Debugf("stream stopped")
			return nil
		}
	}
}

func (p *XdsProxy) handleUpstreamDeltaRequest(ctx context.Context, con *ProxyConnection) {
	defer func() {
		_ = con.upstreamDeltas.CloseSend()
	}()
	for {
		select {
		case req := <-con.deltaRequestsChan:
			proxyLog.Debugf("delta request for type url %s", req.TypeUrl)
			metrics.XdsProxyRequests.Increment()
			if req.TypeUrl == v3.ExtensionConfigurationType {
				p.ecdsLastNonce.Store(req.ResponseNonce)
			}
			if err := sendUpstreamDeltaWithTimeout(ctx, con.upstreamDeltas, req); err != nil {
				proxyLog.Errorf("upstream send error for type url %s: %v", req.TypeUrl, err)
				con.upstreamError <- err
				return
			}
		case <-con.stopChan:
			return
		}
	}
}

func (p *XdsProxy) handleUpstreamDeltaResponse(con *ProxyConnection) {
	for {
		select {
		case resp := <-con.deltaResponsesChan:
			// TODO: separate upstream response handling from requests sending, which are both time costly
			proxyLog.Debugf("response for type url %s", resp.TypeUrl)
			metrics.XdsProxyResponses.Increment()
			// TODO(howardjohn) implement handlers. We need to make handlers take in a list of
			// resources or something probably
			//if h, f := p.handlers[resp.TypeUrl]; f {
			//	err := h(resp)
			//	var errorResp *google_rpc.Status
			//	if err != nil {
			//		errorResp = &google_rpc.Status{
			//			Code:    int32(codes.Internal),
			//			Message: err.Error(),
			//		}
			//	}
			//	// Send ACK/NACK
			//	con.requestsChan <- &discovery.DiscoveryRequest{
			//		VersionInfo:   resp.VersionInfo,
			//		TypeUrl:       resp.TypeUrl,
			//		ResponseNonce: resp.Nonce,
			//		ErrorDetail:   errorResp,
			//	}
			//	continue
			//}
			switch resp.TypeUrl {
			// TODO: fix WASM
			// case v3.ExtensionConfigurationType:
			//	//if features.WasmRemoteLoadConversion {
			//	//	// If Wasm remote load conversion feature is enabled, rewrite and send.
			//	//	go p.rewriteAndForward(con, resp)
			//	//} else {
			//		// Otherwise, forward ECDS resource update directly to Envoy.
			//		forwardDeltaToEnvoy(con, resp)
			//	}
			default:
				forwardDeltaToEnvoy(con, resp)
			}
		case <-con.stopChan:
			return
		}
	}
}

func forwardDeltaToEnvoy(con *ProxyConnection, resp *discovery.DeltaDiscoveryResponse) {
	if err := sendDownstreamDeltaWithTimout(con.downstreamDeltas, resp); err != nil {
		select {
		case con.downstreamError <- err:
			proxyLog.Errorf("downstream send error: %v", err)
		default:
			proxyLog.Debugf("downstream error channel full, but get downstream send error: %v", err)
		}

		return
	}
}

func sendUpstreamDeltaWithTimeout(ctx context.Context, deltaUpstream discovery.AggregatedDiscoveryService_DeltaAggregatedResourcesClient,
	req *discovery.DeltaDiscoveryRequest) error {
	return sendWithTimeout(ctx, func(errChan chan error) {
		errChan <- deltaUpstream.Send(req)
		close(errChan)
	})
}

func sendDownstreamDeltaWithTimout(deltaUpstream discovery.AggregatedDiscoveryService_DeltaAggregatedResourcesServer,
	req *discovery.DeltaDiscoveryResponse) error {
	return sendWithTimeout(context.Background(), func(errChan chan error) {
		errChan <- deltaUpstream.Send(req)
		close(errChan)
	})
}

func (p *XdsProxy) PersistDeltaRequest(req *discovery.DeltaDiscoveryRequest) {
	var ch chan *discovery.DeltaDiscoveryRequest

	p.connectedMutex.Lock()
	if p.connected != nil {
		ch = p.connected.deltaRequestsChan
	}
	p.initialDeltaRequest = req
	p.connectedMutex.Unlock()

	// Immediately send if we are currently connect
	if ch != nil {
		ch <- req
	}
}
