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

package adsc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"time"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	ads "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/conversion"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pstruct "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	mcp "istio.io/api/mcp/v1alpha1"
	"istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/util"
	v2 "istio.io/istio/pilot/pkg/proxy/envoy/v2"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/pkg/log"
)

// Config for the ADS connection.
type Config struct {
	// Namespace defaults to 'default'
	Namespace string

	// Workload defaults to 'test'
	Workload string

	// Meta includes additional metadata for the node
	Meta *pstruct.Struct

	// NodeType defaults to sidecar. "ingress" and "router" are also supported.
	NodeType string

	// IP is currently the primary key used to locate inbound configs. It is sent by client,
	// must match a known endpoint IP. Tests can use a ServiceEntry to register fake IPs.
	IP string
}

// ADSC implements a basic client for ADS, for use in stress tests and tools
// or libraries that need to connect to Istio pilot or other ADS servers.
type ADSC struct {
	// Stream is the GRPC connection stream, allowing direct GRPC send operations.
	// Set after Dial is called.
	stream ads.AggregatedDiscoveryService_StreamAggregatedResourcesClient

	conn *grpc.ClientConn

	// NodeID is the node identity sent to Pilot.
	nodeID string

	certDir string
	url     string

	watchTime time.Time

	// InitialLoad tracks the time to receive the initial configuration.
	InitialLoad time.Duration

	// httpListeners contains received listeners with a http_connection_manager filter.
	httpListeners map[string]*xdsapi.Listener

	// tcpListeners contains all listeners of type TCP (not-HTTP)
	tcpListeners map[string]*xdsapi.Listener

	// All received clusters of type eds, keyed by name
	edsClusters map[string]*xdsapi.Cluster

	// All received clusters of no-eds type, keyed by name
	clusters map[string]*xdsapi.Cluster

	// All received routes, keyed by route name
	routes map[string]*xdsapi.RouteConfiguration

	// All received endpoints, keyed by cluster name
	eds map[string]*xdsapi.ClusterLoadAssignment

	// Metadata has the node metadata to send to pilot.
	// If nil, the defaults will be used.
	Metadata *pstruct.Struct

	// Updates includes the type of the last update received from the server.
	Updates     chan *xdsapi.DiscoveryResponse
	VersionInfo map[string]string

	// Last received message, by type
	Received map[string]*xdsapi.DiscoveryResponse

	mutex sync.Mutex

	Mesh *v1alpha1.MeshConfig

	// Retrieved configurations can be stored using the common istio model interface.
	Store model.IstioConfigStore

	// Retrieved endpoints can be stored in the memory registry. This is used for CDS and EDS responses.
	Registry *v2.MemServiceDiscovery

	// LocalCacheDir is set to a base name used to save fetched resources.
	// If set, each update will be saved.
	// TODO: also load at startup - so we can support warm up in init-container, and survive
	// restarts.
	LocalCacheDir string
}

const (
	typePrefix = "type.googleapis.com/envoy.api.v2."

	// Constants used for XDS

	// ClusterType is used for cluster discovery. Typically first request received
	ClusterType = typePrefix + "Cluster"
	// EndpointType is used for EDS and ADS endpoint discovery. Typically second request.
	endpointType = typePrefix + "ClusterLoadAssignment"
	// ListenerType is sent after clusters and endpoints.
	ListenerType = typePrefix + "Listener"
	// RouteType is sent after listeners.
	routeType = typePrefix + "RouteConfiguration"
)

var (
	adscLog = log.RegisterScope("adsc", "adsc debugging", 0)
)

// Dial connects to a ADS server, with optional MTLS authentication if a cert dir is specified.
func Dial(url string, certDir string, opts *Config) (*ADSC, error) {
	adsc := &ADSC{
		Updates:     make(chan *xdsapi.DiscoveryResponse, 100),
		VersionInfo: map[string]string{},
		certDir:     certDir,
		url:         url,
		Received:    map[string]*xdsapi.DiscoveryResponse{},
	}
	if opts.Namespace == "" {
		opts.Namespace = "default"
	}
	if opts.NodeType == "" {
		opts.NodeType = "sidecar"
	}
	if opts.IP == "" {
		opts.IP = getPrivateIPIfAvailable().String()
	}
	if opts.Workload == "" {
		opts.Workload = "test-1"
	}
	adsc.Metadata = opts.Meta

	adsc.nodeID = fmt.Sprintf("%s~%s~%s.%s~%s.svc.cluster.local", opts.NodeType, opts.IP,
		opts.Workload, opts.Namespace, opts.Namespace)

	err := adsc.Run()
	return adsc, err
}

// Returns a private IP address, or unspecified IP (0.0.0.0) if no IP is available
func getPrivateIPIfAvailable() net.IP {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		default:
			continue
		}
		if !ip.IsLoopback() {
			return ip
		}
	}
	return net.IPv4zero
}

func tlsConfig(certDir string) (*tls.Config, error) {
	clientCert, err := tls.LoadX509KeyPair(certDir+"/cert-chain.pem",
		certDir+"/key.pem")
	if err != nil {
		return nil, err
	}

	serverCABytes, err := ioutil.ReadFile(certDir + "/root-cert.pem")
	if err != nil {
		return nil, err
	}
	serverCAs := x509.NewCertPool()
	if ok := serverCAs.AppendCertsFromPEM(serverCABytes); !ok {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      serverCAs,
		ServerName:   "istio-pilot.istio-system",
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return nil
		},
	}, nil
}

// Close the stream.
func (a *ADSC) Close() {
	a.mutex.Lock()
	a.conn.Close()
	a.mutex.Unlock()
}

// Run will run the ADS client.
func (a *ADSC) Run() error {

	// TODO: pass version info, nonce properly
	var err error
	if len(a.certDir) > 0 {
		tlsCfg, err := tlsConfig(a.certDir)
		if err != nil {
			return err
		}
		creds := credentials.NewTLS(tlsCfg)

		opts := []grpc.DialOption{
			// Verify Pilot cert and service account
			grpc.WithTransportCredentials(creds),
		}
		a.conn, err = grpc.Dial(a.url, opts...)
		if err != nil {
			return err
		}
	} else {
		a.conn, err = grpc.Dial(a.url, grpc.WithInsecure())
		if err != nil {
			return err
		}
	}

	xds := ads.NewAggregatedDiscoveryServiceClient(a.conn)
	edsstr, err := xds.StreamAggregatedResources(context.Background())
	if err != nil {
		return err
	}
	a.stream = edsstr
	go a.handleRecv()
	return nil
}

func (a *ADSC) handleRecv() {
	for {
		msg, err := a.stream.Recv()
		if err != nil {
			adscLog.Infof("Connection closed for node %v with err: %v", a.nodeID, err)
			a.Close()
			a.WaitClear()
			a.Updates <- nil
			return
		}

		// Group-value-kind - used for high level api generator.
		gvk := strings.SplitN(msg.TypeUrl, "/", 3)

		if msg.TypeUrl == collections.IstioMeshV1Alpha1MeshConfig.Resource().GroupVersionKind().String() &&
			len(msg.Resources) > 0 {
			rsc := msg.Resources[0]
			m := &v1alpha1.MeshConfig{}
			err = proto.Unmarshal(rsc.Value, m)
			if err != nil {
				log.Warna("Failed to unmarshal mesh config", err)
			}
			a.Mesh = m
			strResponse, err := json.MarshalIndent(m, "  ", "  ")
			if err != nil {
				continue
			}
			err = ioutil.WriteFile(a.LocalCacheDir+"_mesh.json", strResponse, 0644)
			if err != nil {
				continue
			}

			continue
		}

		listeners := []*xdsapi.Listener{}
		clusters := []*xdsapi.Cluster{}
		routes := []*xdsapi.RouteConfiguration{}
		eds := []*xdsapi.ClusterLoadAssignment{}
		log.Infoa("Update for ", gvk, len(msg.Resources))
		for _, rsc := range msg.Resources { // Any
			a.VersionInfo[rsc.TypeUrl] = msg.VersionInfo
			valBytes := rsc.Value
			switch rsc.TypeUrl {
			case "":
				// last resource of this type. IstioStore doesn't yet have a way to indicate this
				// or to support large sets - memory store has a special handling by calling Create
				// with an empty resource.
				if a.Store != nil && len(gvk) == 3 {
					a.Store.Create(model.Config{
						ConfigMeta: model.ConfigMeta{
							Group:   gvk[0],
							Version: gvk[1],
							Type:    gvk[2],
						},
					})
					log.Infoa("Sync for ", gvk, len(msg.Resources))
				}

			case ListenerType:
				{
					ll := &xdsapi.Listener{}
					_ = proto.Unmarshal(valBytes, ll)
					listeners = append(listeners, ll)
				}

			case ClusterType:
				{
					cl := &xdsapi.Cluster{}
					_ = proto.Unmarshal(valBytes, cl)
					clusters = append(clusters, cl)
				}

			case endpointType:
				{
					el := &xdsapi.ClusterLoadAssignment{}
					_ = proto.Unmarshal(valBytes, el)
					eds = append(eds, el)
				}

			case routeType:
				{
					rl := &xdsapi.RouteConfiguration{}
					_ = proto.Unmarshal(valBytes, rl)
					routes = append(routes, rl)
				}

			default:
				{
					if len(gvk) != 3 {
						continue
					}
					// Generic - fill up the store
					if a.Store != nil {
						m := &mcp.Resource{}
						types.UnmarshalAny(&types.Any{
							TypeUrl: rsc.TypeUrl,
							Value:   rsc.Value,
						}, m)
						val, err := mcpToPilot(m)
						val.Group = gvk[0]
						val.Version = gvk[1]
						val.Type = gvk[2]
						if err != nil {
							adscLog.Warna("Invalid data ", err, " ", string(valBytes))
						} else {
							cfg := a.Store.Get(val.GroupVersionKind(), val.Name, val.Namespace)
							if cfg == nil {
								a.Store.Create(*val)
							} else {
								a.Store.Update(*val)
							}
						}
						if a.LocalCacheDir != "" {
							strResponse, err := json.MarshalIndent(val, "  ", "  ")
							if err != nil {
								continue
							}
							err = ioutil.WriteFile(a.LocalCacheDir+"_res."+
								val.Type+"."+val.Namespace+"."+val.Name+".json", strResponse, 0644)
							if err != nil {
								continue
							}
						}
					}
				}
			}
		}
		// If we got no resource - still save to the store with empty name/namespace, to notify sync
		// This scheme also allows us to chunk large responses !

		a.Received[msg.TypeUrl] = msg

		// TODO: add hook to inject nacks
		a.mutex.Lock()
		a.ack(msg)
		a.mutex.Unlock()

		if len(listeners) > 0 {
			a.handleLDS(listeners)
		} else if len(clusters) > 0 {
			a.handleCDS(clusters)
		} else if len(eds) > 0 {
			a.handleEDS(eds)
		} else if len(routes) > 0 {
			a.handleRDS(routes)
		}
		select {
		case a.Updates <- msg:
		default:
		}

	}
}

func mcpToPilot(m *mcp.Resource) (*model.Config, error) {
	if m == nil || m.Metadata == nil {
		return &model.Config{}, nil
	}
	c := &model.Config{
		ConfigMeta: model.ConfigMeta{
			ResourceVersion: m.Metadata.Version,
			Labels:          m.Metadata.Labels,
			Annotations:     m.Metadata.Annotations,
		},
	}
	nsn := strings.Split(m.Metadata.Name, "/")
	if len(nsn) != 2 {
		return nil, fmt.Errorf("Invalid name %s", m.Metadata.Name)
	}
	c.Namespace = nsn[0]
	c.Name = nsn[1]
	var err error
	c.CreationTimestamp, err = types.TimestampFromProto(m.Metadata.CreateTime)
	if err != nil {
		return nil, err
	}

	pb, err := types.EmptyAny(m.Body)
	if err != nil {
		return nil, err
	}
	err = types.UnmarshalAny(m.Body, pb)
	if err != nil {
		return nil, err
	}
	c.Spec = pb
	return c, nil
}

// nolint: staticcheck
func (a *ADSC) handleLDS(ll []*xdsapi.Listener) {
	lh := map[string]*xdsapi.Listener{}
	lt := map[string]*xdsapi.Listener{}

	routes := []string{}
	ldsSize := 0

	for _, l := range ll {
		ldsSize += proto.Size(l)

		// The last filter is the actual destination for inbound listener
		if l.ApiListener != nil {
			// This is an API Listener
			// TODO: extract VIP and RDS or cluster
			continue
		}
		filter := l.FilterChains[len(l.FilterChains)-1].Filters[0]

		// The actual destination will be the next to the last if the last filter is a passthrough filter
		if l.FilterChains[len(l.FilterChains)-1].GetName() == util.PassthroughFilterChain {
			filter = l.FilterChains[len(l.FilterChains)-2].Filters[0]
		}

		if filter.Name == "mixer" {
			filter = l.FilterChains[len(l.FilterChains)-1].Filters[1]
		}
		if filter.Name == "envoy.tcp_proxy" {
			lt[l.Name] = l
			config := filter.GetConfig()
			if config == nil {
				config, _ = conversion.MessageToStruct(filter.GetTypedConfig())
			}
			c := config.Fields["cluster"].GetStringValue()
			adscLog.Debugf("TCP: %s -> %s", l.Name, c)
		} else if filter.Name == "envoy.http_connection_manager" {
			lh[l.Name] = l

			// Getting from config is too painful..
			port := l.Address.GetSocketAddress().GetPortValue()
			if port == 15002 {
				routes = append(routes, "http_proxy")
			} else {
				routes = append(routes, fmt.Sprintf("%d", port))
			}
		} else if filter.Name == "envoy.mongo_proxy" {
			// ignore for now
		} else if filter.Name == "envoy.redis_proxy" {
			// ignore for now
		} else if filter.Name == "envoy.filters.network.mysql_proxy" {
			// ignore for now
		} else {
			tm := &jsonpb.Marshaler{Indent: "  "}
			adscLog.Infof(tm.MarshalToString(l))
		}
	}

	adscLog.Infof("LDS: http=%d tcp=%d size=%d", len(lh), len(lt), ldsSize)
	if adscLog.DebugEnabled() {
		b, _ := json.MarshalIndent(ll, " ", " ")
		adscLog.Debugf(string(b))
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if len(routes) > 0 {
		a.sendRsc(routeType, routes)
	}
	a.httpListeners = lh
	a.tcpListeners = lt
}

// compact representations, for simplified debugging/testing

// TCPListener extracts the core elements from envoy Listener.
type TCPListener struct {
	// Address is the address, as expected by go Dial and Listen, including port
	Address string

	// LogFile is the access log address for the listener
	LogFile string

	// Target is the destination cluster.
	Target string
}

type Target struct {

	// Address is a go address, extracted from the mangled cluster name.
	Address string

	// Endpoints are the resolved endpoints from eds or cluster static.
	Endpoints map[string]Endpoint
}

type Endpoint struct {
	// Weight extracted from eds
	Weight int
}

// Save will save the json configs to files, using the base directory
func (a *ADSC) Save(base string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	strResponse, err := json.MarshalIndent(a.tcpListeners, "  ", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(base+"_lds_tcp.json", strResponse, 0644)
	if err != nil {
		return err
	}
	strResponse, err = json.MarshalIndent(a.httpListeners, "  ", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(base+"_lds_http.json", strResponse, 0644)
	if err != nil {
		return err
	}
	strResponse, err = json.MarshalIndent(a.routes, "  ", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(base+"_rds.json", strResponse, 0644)
	if err != nil {
		return err
	}
	strResponse, err = json.MarshalIndent(a.edsClusters, "  ", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(base+"_ecds.json", strResponse, 0644)
	if err != nil {
		return err
	}
	strResponse, err = json.MarshalIndent(a.clusters, "  ", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(base+"_cds.json", strResponse, 0644)
	if err != nil {
		return err
	}
	strResponse, err = json.MarshalIndent(a.eds, "  ", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(base+"_eds.json", strResponse, 0644)
	if err != nil {
		return err
	}

	return err
}

func (a *ADSC) handleCDS(ll []*xdsapi.Cluster) {

	cn := []string{}
	cdsSize := 0
	edscds := map[string]*xdsapi.Cluster{}
	cds := map[string]*xdsapi.Cluster{}
	for _, c := range ll {
		cdsSize += proto.Size(c)
		switch v := c.ClusterDiscoveryType.(type) {
		case *xdsapi.Cluster_Type:
			if v.Type != xdsapi.Cluster_EDS {
				cds[c.Name] = c
				continue
			}
		}
		cn = append(cn, c.Name)
		edscds[c.Name] = c
	}

	adscLog.Infof("CDS: %d size=%d", len(cn), cdsSize)

	if len(cn) > 0 {
		a.sendRsc(endpointType, cn)
	}
	if adscLog.DebugEnabled() {
		b, _ := json.MarshalIndent(ll, " ", " ")
		adscLog.Info(string(b))
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.edsClusters = edscds
	a.clusters = cds
}

func (a *ADSC) node() *core.Node {
	n := &core.Node{
		Id: a.nodeID,
	}
	if a.Metadata == nil {
		n.Metadata = &pstruct.Struct{
			Fields: map[string]*pstruct.Value{
				"ISTIO_VERSION": {Kind: &pstruct.Value_StringValue{StringValue: "65536.65536.65536"}},
			}}
	} else {
		n.Metadata = a.Metadata
		if a.Metadata.Fields["ISTIO_VERSION"] == nil {
			a.Metadata.Fields["ISTIO_VERSION"] = &pstruct.Value{Kind: &pstruct.Value_StringValue{StringValue: "65536.65536.65536"}}
		}
	}
	return n
}

func (a *ADSC) Send(req *xdsapi.DiscoveryRequest) error {
	req.Node = a.node()
	req.ResponseNonce = time.Now().String()
	return a.stream.Send(req)
}

func (a *ADSC) handleEDS(eds []*xdsapi.ClusterLoadAssignment) {
	la := map[string]*xdsapi.ClusterLoadAssignment{}
	edsSize := 0
	ep := 0
	for _, cla := range eds {
		edsSize += proto.Size(cla)
		la[cla.ClusterName] = cla
		ep += len(cla.Endpoints)
	}

	adscLog.Infof("eds: %d size=%d ep=%d", len(eds), edsSize, ep)
	if adscLog.DebugEnabled() {
		b, _ := json.MarshalIndent(eds, " ", " ")
		adscLog.Info(string(b))
	}
	if a.InitialLoad == 0 {
		// first load - Envoy loads listeners after endpoints
		_ = a.stream.Send(&xdsapi.DiscoveryRequest{
			ResponseNonce: time.Now().String(),
			Node:          a.node(),
			TypeUrl:       ListenerType,
		})
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.eds = la
}

func (a *ADSC) handleRDS(configurations []*xdsapi.RouteConfiguration) {

	vh := 0
	rcount := 0
	size := 0

	rds := map[string]*xdsapi.RouteConfiguration{}

	for _, r := range configurations {
		for _, h := range r.VirtualHosts {
			vh++
			for _, rt := range h.Routes {
				rcount++
				// Example: match:<prefix:"/" > route:<cluster:"outbound|9154||load-se-154.local" ...
				adscLog.Debugf("Handle route %v, path %v, cluster %v", h.Name, rt.Match.PathSpecifier, rt.GetRoute().GetCluster())
			}
		}
		rds[r.Name] = r
		size += proto.Size(r)
	}
	if a.InitialLoad == 0 {
		a.InitialLoad = time.Since(a.watchTime)
		adscLog.Infof("RDS: %d size=%d vhosts=%d routes=%d time=%d", len(configurations), size, vh, rcount, a.InitialLoad)
	} else {
		adscLog.Infof("RDS: %d size=%d vhosts=%d routes=%d", len(configurations), size, vh, rcount)
	}

	if adscLog.DebugEnabled() {
		b, _ := json.MarshalIndent(configurations, " ", " ")
		adscLog.Info(string(b))
	}

	a.mutex.Lock()
	a.routes = rds
	a.mutex.Unlock()
}

// WaitClear will clear the waiting events, so next call to Wait will get
// the next push type.
func (a *ADSC) WaitClear() {
	for {
		select {
		case <-a.Updates:
		default:
			return
		}
	}
}

// Wait for an updates for all the specified types
// If updates is empty, this will wait for any update
func (a *ADSC) Wait(to time.Duration, updates ...string) ([]string, error) {
	t := time.NewTimer(to)
	want := map[string]struct{}{}
	for _, update := range updates {
		switch update {
		case "lds":
			update = ListenerType
		case "cds":
			update = ClusterType
		case "eds":
			update = endpointType
		case "rds":
			update = routeType
		}
		want[update] = struct{}{}
	}
	got := make([]string, 0, len(updates))
	for {
		select {
		case t := <-a.Updates:
			if t == nil {
				return got, fmt.Errorf("Closed")
			}
			delete(want, t.TypeUrl)
			got = append(got, t.TypeUrl)
			if len(want) == 0 {
				return got, nil
			}
		case <-t.C:
			return got, fmt.Errorf("timeout, still waiting for updates: %v", want)
		}
	}
}

// WaitVersion waits for a new or updated for a typeURL.
func (a *ADSC) WaitVersion(to time.Duration, typeURL, lastVersion string) (*xdsapi.DiscoveryResponse, error) {
	t := time.NewTimer(to)
	ex := a.Received[typeURL]
	if ex != nil {
		if lastVersion == "" {
			return ex, nil
		}
		if lastVersion != ex.VersionInfo {
			return ex, nil
		}
	}

	for {
		select {
		case t := <-a.Updates:
			if t == nil {
				return nil, fmt.Errorf("Closed")
			}
			if t.TypeUrl == typeURL {
				return t, nil
			}

		case <-t.C:
			return nil, fmt.Errorf("timeout, still waiting for updates: %v", typeURL)
		}
	}
}

// EndpointsJSON returns the endpoints, formatted as JSON, for debugging.
func (a *ADSC) EndpointsJSON() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	out, _ := json.MarshalIndent(a.eds, " ", " ")
	return string(out)
}

// Watch will start watching resources, starting with LDS. Based on the LDS response
// it will start watching RDS and CDS.
func (a *ADSC) Watch() {
	a.watchTime = time.Now()
	_ = a.stream.Send(&xdsapi.DiscoveryRequest{
		ResponseNonce: time.Now().String(),
		Node:          a.node(),
		TypeUrl:       ClusterType,
	})
}

// WatchConfig will use the new experimental API watching, similar with MCP.
func (a *ADSC) WatchConfig() {
	_ = a.stream.Send(&xdsapi.DiscoveryRequest{
		ResponseNonce: time.Now().String(),
		Node:          a.node(),
		TypeUrl:       collections.IstioMeshV1Alpha1MeshConfig.Resource().GroupVersionKind().String(),
	})

	for _, sch := range collections.Pilot.All() {
		_ = a.stream.Send(&xdsapi.DiscoveryRequest{
			ResponseNonce: time.Now().String(),
			Node:          a.node(),
			TypeUrl:       sch.Resource().GroupVersionKind().String(),
		})
	}
}

func (a *ADSC) sendRsc(typeurl string, rsc []string) {
	_ = a.stream.Send(&xdsapi.DiscoveryRequest{
		ResponseNonce: "",
		Node:          a.node(),
		TypeUrl:       typeurl,
		ResourceNames: rsc,
	})
}

func (a *ADSC) ack(msg *xdsapi.DiscoveryResponse) {
	_ = a.stream.Send(&xdsapi.DiscoveryRequest{
		ResponseNonce: msg.Nonce,
		TypeUrl:       msg.TypeUrl,
		Node:          a.node(),
		VersionInfo:   msg.VersionInfo,
	})
}

// GetHTTPListeners returns all the http listeners.
func (a *ADSC) GetHTTPListeners() map[string]*xdsapi.Listener {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.httpListeners
}

// GetTCPListeners returns all the tcp listeners.
func (a *ADSC) GetTCPListeners() map[string]*xdsapi.Listener {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.tcpListeners
}

// GetEdsClusters returns all the eds type clusters.
func (a *ADSC) GetEdsClusters() map[string]*xdsapi.Cluster {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.edsClusters
}

// GetClusters returns all the non-eds type clusters.
func (a *ADSC) GetClusters() map[string]*xdsapi.Cluster {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.clusters
}

// GetRoutes returns all the routes.
func (a *ADSC) GetRoutes() map[string]*xdsapi.RouteConfiguration {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.routes
}

// GetEndpoints returns all the routes.
func (a *ADSC) GetEndpoints() map[string]*xdsapi.ClusterLoadAssignment {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.eds
}
