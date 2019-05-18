// Package sdsc includes a lightweight testing client to interact with SDS.
package sdsc

import (
	"context"
	"fmt"
	"time"

	"istio.io/pkg/log"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	sds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"google.golang.org/grpc"
)

// Client is a lightweight client for testing secret discovery service server.
type Client struct {
	// adsc *adsc.ADSC
	stream     sds.SecretDiscoveryService_StreamSecretsClient
	conn       *grpc.ClientConn
	updateChan chan xdsapi.DiscoveryResponse
	nodeID     string
	udsPath    string
}

// ClientOptions contains the options for the SDS testing
type ClientOptions struct {
	ServerAddress string
}

// NewClient returns a sds client for testing.
func NewClient(options ClientOptions) (*Client, error) {
	address := fmt.Sprintf("unix://%s", options.ServerAddress)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := sds.NewSecretDiscoveryServiceClient(conn)
	stream, err := client.StreamSecrets(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{
		stream:     stream,
		conn:       conn,
		updateChan: make(chan xdsapi.DiscoveryResponse, 1),
		udsPath:    address,
	}, nil
}

// Start starts sds client to recieve the scecret updates from the server.
func (c *Client) Start() {
	go func() {
		msq, err := c.stream.Recv()
		if err != nil {
			log.Errorf("Connection closed %v", err)
			return
		}
		c.updateChan <- *msq
		log.Infof("received response from sds server %v", msq)
	}()
}

// Stop stops the sds client.
func (c *Client) Stop() {
	c.stream.CloseSend()
}

// WaitForUpdate blocks untill the error occurs or updates are pushed from the sds server.
func (c *Client) WaitForUpdate(duration time.Duration) (*xdsapi.DiscoveryResponse, error) {
	t := time.NewTimer(duration)
	for {
		select {
		case resp := <-c.updateChan:
			return &resp, nil
		case <-t.C:
			return nil, fmt.Errorf("timeout for updates")
		}
	}
}

// Send sends a request to the agent.
func (c *Client) Send() (*xdsapi.DiscoveryResponse, error) {
	err := c.stream.Send(&xdsapi.DiscoveryRequest{
		VersionInfo: "abcd",
		TypeUrl:     "type.googleapis.com/envoy.api.v2.auth.Secret",
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
