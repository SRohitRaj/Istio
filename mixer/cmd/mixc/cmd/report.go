// Copyright 2017 Istio Authors
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

package cmd

import (
	"context"

	ot "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/spf13/cobra"

	mixerpb "istio.io/api/mixer/v1"
	"istio.io/istio/pkg/log"
)

func reportCmd(rootArgs *rootArgs) *cobra.Command {
	return &cobra.Command{
		Use:   "report",
		Short: "Invokes Mixer's Report API to generate telemetry.",
		Long: "The Report method is used to produce telemetry. Mixer\n" +
			"expects a set of attributes as input, which it uses, along with\n" +
			"its configuration, to determine which adapters to invoke and with\n" +
			"which parameters in order to output the telemetry.",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			report(rootArgs)
		}}
}

func report(rootArgs *rootArgs) {
	var attrs *mixerpb.CompressedAttributes
	var err error

	if attrs, err = parseAttributes(rootArgs); err != nil {
		log.Fatalf("%v", err)
	}

	var cs *clientState
	if cs, err = createAPIClient(rootArgs.mixerAddress, rootArgs.tracingOptions); err != nil {
		log.Fatalf("Unable to establish connection to %s: %v", rootArgs.mixerAddress, err)
	}
	defer deleteAPIClient(cs)

	span, ctx := ot.StartSpanFromContext(context.Background(), "mixc Report", ext.SpanKindRPCClient)

	for i := 0; i < rootArgs.repeat; i++ {
		request := mixerpb.ReportRequest{Attributes: []mixerpb.CompressedAttributes{*attrs}}
		_, err := cs.client.Report(ctx, &request)

		log.Infof("Report RPC returned %s", decodeError(err))
	}

	span.Finish()
}
