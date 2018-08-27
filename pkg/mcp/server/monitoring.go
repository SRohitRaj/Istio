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

package server

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opencensus.io/stats/view"
)

const (
	typeURL   = "typeURL"
	errorCode = "code"
	errorStr  = "error"
	connectionID = "connectionID"
)

var (
	TypeURLTag      tag.Key
	ErrorCodeTag    tag.Key
	ErrorTag        tag.Key
	ConnectionIDTag tag.Key

	// buckets are powers of 4
	byteBuckets = []float64{1, 4, 16, 64, 256, 1024, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824}

	ClientCount = stats.Int64(
		"mcp/server/client_count",
		"The number of clients currently connected.",
		stats.UnitDimensionless)

	RequestSizesBytes = stats.Int64(
		"mcp/server/message_sizes_bytes",
		"Size of messages received from clients.",
		stats.UnitBytes)

	RequestAcks = stats.Int64(
		"mcp/server/request_acks",
		"The number of request acks received by the server.",
		stats.UnitDimensionless)

	RequestNacks = stats.Int64(
		"mcp/server/request_nacks",
		"The number of request nacks received by the server.",
		stats.UnitDimensionless)

	SendFailures = stats.Int64(
		"mcp/server/send_failures",
		"The number of send failures in the server.",
		stats.UnitDimensionless)

	RecvFailures = stats.Int64(
		"mcp/server/recv_failures",
		"The number of recv failures in the server.",
		stats.UnitDimensionless)
)

func newView(measure stats.Measure, keys []tag.Key, aggregation *view.Aggregation) *view.View {
	return &view.View{
		Name:        measure.Name(),
		Description: measure.Description(),
		Measure:     measure,
		TagKeys:     keys,
		Aggregation: aggregation,
	}
}

func init() {
	var err error
	if TypeURLTag, err = tag.NewKey(typeURL); err != nil {
		panic(err)
	}
	if ErrorCodeTag, err = tag.NewKey(errorCode); err != nil {
		panic(err)
	}
	if ErrorTag, err = tag.NewKey(errorStr); err != nil {
		panic(err)
	}
	if ConnectionIDTag, err = tag.NewKey(connectionID); err != nil {
		panic(err)
	}

	err = view.Register(
		newView(ClientCount, []tag.Key{}, view.LastValue()),
		newView(RequestSizesBytes, []tag.Key{ConnectionIDTag}, view.Distribution(byteBuckets...)),
		newView(RequestAcks, []tag.Key{TypeURLTag}, view.Count()),
		newView(RequestNacks, []tag.Key{ErrorCodeTag, TypeURLTag}, view.Count()),
		newView(SendFailures, []tag.Key{ErrorCodeTag, ErrorTag}, view.Count()),
		newView(RecvFailures, []tag.Key{ErrorCodeTag, ErrorTag}, view.Count()),
	)

	if err != nil {
		panic(err)
	}
}
