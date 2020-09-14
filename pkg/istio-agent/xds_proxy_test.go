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
	"net"
	"path"
	"testing"
	"time"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"istio.io/istio/pilot/pkg/xds"
	v3 "istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pkg/config/mesh"
	"istio.io/istio/pkg/security"
	"istio.io/istio/pkg/test/env"
)

func TestXdsProxyBasicFlow(t *testing.T) {
	secOpts := &security.Options{
		FileMountedCerts: true,
	}
	proxyConfig := mesh.DefaultProxyConfig()
	proxyConfig.DiscoveryAddress = "localhost:1111"

	proxyConfig.ProxyMetadata = map[string]string{
		MetadataClientCertChain: path.Join(env.IstioSrc, "tests/testdata/certs/pilot/cert-chain.pem"),
		MetadataClientCertKey:   path.Join(env.IstioSrc, "tests/testdata/certs/pilot/key.pem"),
		MetadataClientRootCert:  path.Join(env.IstioSrc, "tests/testdata/certs/pilot/root-cert.pem"),
	}
	ia := NewAgent(&proxyConfig,
		&AgentConfig{}, secOpts)
	proxy, err := initXdsProxy(ia)
	if err != nil {
		t.Fatalf("Failed to initialize xds proxy %v", err)
	}

	proxy.istiodDialOptions = []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithInsecure(),
	}
	startUpstreamServer()

	res, err := adsRequestStream("./etc/istio/proxy/XDS", &discovery.DiscoveryRequest{
		TypeUrl: v3.ClusterType,
		Node: &core.Node{
			Id: "sidecar~0.0.0.0~debug~cluster.local",
		},
	})
	if err != nil {
		t.Fatalf("Failed to receive proxied response %v", err)
	}
	if res == nil || res.TypeUrl != v3.ClusterType {
		t.Fatalf("Expected to get cluster response but got %v", res)
	}
	time.Sleep(3 * time.Minute)
	ia.Close()
}

func startUpstreamServer() {
	simple := xds.NewXDS()
	simple.StartGRPC("localhost:1111")
}

func adsRequestStream(socket string, req *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	conn, err := setupConnection(socket)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	adsClient := discovery.NewAggregatedDiscoveryServiceClient(conn)
	ctx := metadata.AppendToOutgoingContext(context.Background(), "ClusterID", "Kubernetes")
	downstream, err := adsClient.StreamAggregatedResources(ctx,
		grpc.MaxCallRecvMsgSize(defaultClientMaxReceiveMessageSize))
	if err != nil {
		return nil, err
	}
	err = downstream.Send(req)
	if err != nil {
		return nil, err
	}
	res, err := downstream.Recv()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func setupConnection(socket string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure(), grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", socket)
	}))

	conn, err := grpc.Dial(socket, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
