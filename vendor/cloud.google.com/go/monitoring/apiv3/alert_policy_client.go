// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// AUTO-GENERATED CODE. DO NOT EDIT.

package monitoring

import (
	"math"
	"time"

	"cloud.google.com/go/internal/version"
	gax "github.com/googleapis/gax-go"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// AlertPolicyCallOptions contains the retry settings for each method of AlertPolicyClient.
type AlertPolicyCallOptions struct {
	ListAlertPolicies []gax.CallOption
	GetAlertPolicy    []gax.CallOption
	CreateAlertPolicy []gax.CallOption
	DeleteAlertPolicy []gax.CallOption
	UpdateAlertPolicy []gax.CallOption
}

func defaultAlertPolicyClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("monitoring.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultAlertPolicyCallOptions() *AlertPolicyCallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.3,
				})
			}),
		},
	}
	return &AlertPolicyCallOptions{
		ListAlertPolicies: retry[[2]string{"default", "idempotent"}],
		GetAlertPolicy:    retry[[2]string{"default", "idempotent"}],
		CreateAlertPolicy: retry[[2]string{"default", "non_idempotent"}],
		DeleteAlertPolicy: retry[[2]string{"default", "idempotent"}],
		UpdateAlertPolicy: retry[[2]string{"default", "non_idempotent"}],
	}
}

// AlertPolicyClient is a client for interacting with Stackdriver Monitoring API.
type AlertPolicyClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	alertPolicyClient monitoringpb.AlertPolicyServiceClient

	// The call options for this service.
	CallOptions *AlertPolicyCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewAlertPolicyClient creates a new alert policy service client.
//
// The AlertPolicyService API is used to manage (list, create, delete,
// edit) alert policies in Stackdriver Monitoring. An alerting policy is
// a description of the conditions under which some aspect of your
// system is considered to be "unhealthy" and the ways to notify
// people or services about this state. In addition to using this API, alert
// policies can also be managed through
// Stackdriver Monitoring (at https://cloud.google.com/monitoring/docs/),
// which can be reached by clicking the "Monitoring" tab in
// Cloud Console (at https://console.cloud.google.com/).
func NewAlertPolicyClient(ctx context.Context, opts ...option.ClientOption) (*AlertPolicyClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultAlertPolicyClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &AlertPolicyClient{
		conn:        conn,
		CallOptions: defaultAlertPolicyCallOptions(),

		alertPolicyClient: monitoringpb.NewAlertPolicyServiceClient(conn),
	}
	c.setGoogleClientInfo()
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *AlertPolicyClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *AlertPolicyClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *AlertPolicyClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", version.Go()}, keyval...)
	kv = append(kv, "gapic", version.Repo, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// ListAlertPolicies lists the existing alerting policies for the project.
func (c *AlertPolicyClient) ListAlertPolicies(ctx context.Context, req *monitoringpb.ListAlertPoliciesRequest, opts ...gax.CallOption) *AlertPolicyIterator {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.ListAlertPolicies[0:len(c.CallOptions.ListAlertPolicies):len(c.CallOptions.ListAlertPolicies)], opts...)
	it := &AlertPolicyIterator{}
	it.InternalFetch = func(pageSize int, pageToken string) ([]*monitoringpb.AlertPolicy, string, error) {
		var resp *monitoringpb.ListAlertPoliciesResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.alertPolicyClient.ListAlertPolicies(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.AlertPolicies, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	return it
}

// GetAlertPolicy gets a single alerting policy.
func (c *AlertPolicyClient) GetAlertPolicy(ctx context.Context, req *monitoringpb.GetAlertPolicyRequest, opts ...gax.CallOption) (*monitoringpb.AlertPolicy, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.GetAlertPolicy[0:len(c.CallOptions.GetAlertPolicy):len(c.CallOptions.GetAlertPolicy)], opts...)
	var resp *monitoringpb.AlertPolicy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.alertPolicyClient.GetAlertPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateAlertPolicy creates a new alerting policy.
func (c *AlertPolicyClient) CreateAlertPolicy(ctx context.Context, req *monitoringpb.CreateAlertPolicyRequest, opts ...gax.CallOption) (*monitoringpb.AlertPolicy, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.CreateAlertPolicy[0:len(c.CallOptions.CreateAlertPolicy):len(c.CallOptions.CreateAlertPolicy)], opts...)
	var resp *monitoringpb.AlertPolicy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.alertPolicyClient.CreateAlertPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteAlertPolicy deletes an alerting policy.
func (c *AlertPolicyClient) DeleteAlertPolicy(ctx context.Context, req *monitoringpb.DeleteAlertPolicyRequest, opts ...gax.CallOption) error {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.DeleteAlertPolicy[0:len(c.CallOptions.DeleteAlertPolicy):len(c.CallOptions.DeleteAlertPolicy)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.alertPolicyClient.DeleteAlertPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// UpdateAlertPolicy updates an alerting policy. You can either replace the entire policy with
// a new one or replace only certain fields in the current alerting policy by
// specifying the fields to be updated via updateMask. Returns the
// updated alerting policy.
func (c *AlertPolicyClient) UpdateAlertPolicy(ctx context.Context, req *monitoringpb.UpdateAlertPolicyRequest, opts ...gax.CallOption) (*monitoringpb.AlertPolicy, error) {
	ctx = insertMetadata(ctx, c.xGoogMetadata)
	opts = append(c.CallOptions.UpdateAlertPolicy[0:len(c.CallOptions.UpdateAlertPolicy):len(c.CallOptions.UpdateAlertPolicy)], opts...)
	var resp *monitoringpb.AlertPolicy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.alertPolicyClient.UpdateAlertPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// AlertPolicyIterator manages a stream of *monitoringpb.AlertPolicy.
type AlertPolicyIterator struct {
	items    []*monitoringpb.AlertPolicy
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*monitoringpb.AlertPolicy, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *AlertPolicyIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *AlertPolicyIterator) Next() (*monitoringpb.AlertPolicy, error) {
	var item *monitoringpb.AlertPolicy
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *AlertPolicyIterator) bufLen() int {
	return len(it.items)
}

func (it *AlertPolicyIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
