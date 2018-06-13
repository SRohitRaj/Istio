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

// Package sds implements secret discovery service in NodeAgent.
package sds

import (
	"context"
	"fmt"
	"io"
	"time"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	authapi "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	sds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/log"
	"istio.io/istio/security/pkg/nodeagent/cache"
)

// SecretType is used for secret discovery service to construct response.
const SecretType = "type.googleapis.com/envoy.api.v2.Secret"

type discoveryStream interface {
	Send(*xdsapi.DiscoveryResponse) error
	Recv() (*xdsapi.DiscoveryRequest, error)
	grpc.ServerStream
}

// SdsEvent represents a config or registry event that results in a push.
type SdsEvent struct {

	// If not empty, it is used to indicate the event is caused by a change in the clusters.
	// Only EDS for the listed clusters will be sent.
	clusters []string
}

type sdsConnection struct {
	// PeerAddr is the address of the client envoy, from network layer.
	PeerAddr string

	// Time of connection, for debugging.
	Connect time.Time

	// The proxy from which the connection comes from.
	modelNode *model.Proxy

	pushChannel chan *SdsEvent

	// doneChannel will be closed when the client is closed.
	doneChannel chan int

	// SDS streams implement this interface.
	stream discoveryStream
}

type sdsservice struct {
	st cache.SecretManager
	//TODO(quanlin), add below properties later:
	//1. workloadRegistry(store proxies information).
}

// newSDSService creates Secret Discovery Service which implements envoy v2 SDS API.
func newSDSService(st cache.SecretManager) *sdsservice {
	return &sdsservice{
		st: st,
	}
}

// register adds the SDS handle to the grpc server
func (s *sdsservice) register(rpcs *grpc.Server) {
	sds.RegisterSecretDiscoveryServiceServer(rpcs, s)
}

func (s *sdsservice) StreamSecrets(stream sds.SecretDiscoveryService_StreamSecretsServer) error {
	peerAddr := "Unknown peer address"
	peerInfo, ok := peer.FromContext(stream.Context())
	if ok {
		peerAddr = peerInfo.Addr.String()
	}

	var discReq *xdsapi.DiscoveryRequest
	var receiveError error
	reqChannel := make(chan *xdsapi.DiscoveryRequest, 1)

	con := newSDSConnection(peerAddr, stream)
	defer close(con.doneChannel)

	go receiveThread(con, reqChannel, &receiveError)

	for {
		// Block until a request is received.
		select {
		case discReq, ok = <-reqChannel:
			if !ok {
				// Remote side closed connection.
				return receiveError
			}

			if discReq.Node.Id == "" {
				log.Warnf("Discovery request %+v missing node id", discReq)
				continue
			}
			nt, err := model.ParseServiceNode(discReq.Node.Id)
			if err != nil {
				log.Errorf("Failed to parse service node from discovery request %+v: %v", discReq, err)
				return err
			}
			nt.Metadata = model.ParseMetadata(discReq.Node.Metadata)
			con.modelNode = &nt

			secret, err := s.st.GetSecret(discReq.Node.Id, "" /*TODO(quanlin) credential token*/)
			if err != nil {
				log.Errorf("Failed to get secret for proxy %q from secret cache: %v", discReq.Node.Id, err)
				return err
			}

			if err := s.pushSDS(secret, *con.modelNode, con); err != nil {
				log.Errorf("SDS failed to push: %v", err)
				return err
			}
		}
	}
}

func (s *sdsservice) FetchSecrets(ctx context.Context, discReq *xdsapi.DiscoveryRequest) (*xdsapi.DiscoveryResponse, error) {
	if discReq.Node.Id == "" {
		log.Warnf("SDS discovery request %+v missing node id", discReq)
		return nil, fmt.Errorf("SDS discovery request %+v missing node id", discReq)
	}

	proxy, err := model.ParseServiceNode(discReq.Node.Id)
	if err != nil {
		log.Errorf("Failed to parse service node from discovery request %+v: %v", discReq, err)
		return nil, err
	}

	proxy.Metadata = model.ParseMetadata(discReq.Node.Metadata)

	//TODO(quanlin): add proxy info in workload registry.

	secret, err := s.st.GetSecret(discReq.Node.Id, "" /*TODO(quanlin) credential token*/)
	if err != nil {
		log.Errorf("Failed to get secret for proxy %q from secret cache: %v", discReq.Node.Id, err)
		return nil, err
	}

	return sdsDiscoveryResponse(secret, proxy)
}

func (s *sdsservice) pushSDS(secret *cache.SecretItem, proxy model.Proxy, con *sdsConnection) error {
	response, err := sdsDiscoveryResponse(secret, proxy)
	if err != nil {
		log.Errorf("SDS: Failed to construct response %v", err)
		return err
	}

	if err = con.stream.Send(response); err != nil {
		log.Errorf("SDS: Send response failure %v", err)
		return err
	}

	log.Infof("SDS: push for proxy:%q addr:%q", proxy.ID, con.PeerAddr)
	return nil
}

func sdsDiscoveryResponse(s *cache.SecretItem, proxy model.Proxy) (*xdsapi.DiscoveryResponse, error) {
	//TODO(quanlin): use timestamp for versionInfo and nouce for now, may change later.
	t := time.Now().String()
	resp := &xdsapi.DiscoveryResponse{
		TypeUrl:     SecretType,
		VersionInfo: t,
		Nonce:       t,
	}

	if s == nil {
		log.Errorf("SDS: got nil secret for proxy %q", proxy.ID)
		return resp, nil
	}

	secret := &authapi.Secret{
		//TODO(quanlin): better naming.
		Name: "self-signed",
		Type: &authapi.Secret_TlsCertificate{
			TlsCertificate: &authapi.TlsCertificate{
				CertificateChain: &core.DataSource{
					Specifier: &core.DataSource_InlineBytes{
						InlineBytes: s.CertificateChain,
					},
				},
				PrivateKey: &core.DataSource{
					Specifier: &core.DataSource_InlineBytes{
						InlineBytes: s.PrivateKey,
					},
				},
			},
		},
	}

	ms, err := types.MarshalAny(secret)
	if err != nil {
		log.Errorf("Failed to mashal secret for proxy %q: %v", proxy.ID, err)
		return nil, err
	}
	resp.Resources = append(resp.Resources, *ms)

	return resp, nil
}

func newSDSConnection(peerAddr string, stream discoveryStream) *sdsConnection {
	return &sdsConnection{
		doneChannel: make(chan int, 1),
		pushChannel: make(chan *SdsEvent, 1),
		PeerAddr:    peerAddr,
		Connect:     time.Now(),
		stream:      stream,
	}
}

func receiveThread(con *sdsConnection, reqChannel chan *xdsapi.DiscoveryRequest, errP *error) {
	defer close(reqChannel) // indicates close of the remote side.
	for {
		req, err := con.stream.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled || err == io.EOF {
				log.Infof("SDS: %q terminated %v", con.PeerAddr, err)
				return
			}
			*errP = err
			log.Errorf("SDS: %q terminated with errors %v", con.PeerAddr, err)
			return
		}
		reqChannel <- req
	}
}
