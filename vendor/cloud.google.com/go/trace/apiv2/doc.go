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

// Package trace is an auto-generated package for the
// Stackdriver Trace API.
//
//   NOTE: This package is in alpha. It is not stable, and is likely to change.
//
// Sends application trace data to Stackdriver Trace for viewing. Trace data
// is
// collected for all App Engine applications by default. Trace data from
// other
// applications can be provided using this API.
package trace // import "cloud.google.com/go/trace/apiv2"

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func insertMetadata(ctx context.Context, mds ...metadata.MD) context.Context {
	out, _ := metadata.FromOutgoingContext(ctx)
	out = out.Copy()
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return metadata.NewOutgoingContext(ctx, out)
}

// DefaultAuthScopes reports the default set of authentication scopes to use with this package.
func DefaultAuthScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/trace.append",
	}
}
