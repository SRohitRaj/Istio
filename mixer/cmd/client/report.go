// Copyright 2016 Google Inc.
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

package main

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	mixerpb "istio.io/mixer/api/v1"
)

func reportCmd(rootArgs *rootArgs, errorf errorFn) *cobra.Command {
	return &cobra.Command{
		Use:   "report <message>...",
		Short: "Invokes the mixer's Report API.",
		Run: func(cmd *cobra.Command, args []string) {
			report(rootArgs, args, errorf)
		}}
}

func report(rootArgs *rootArgs, args []string, errorf errorFn) {
	var attrs *mixerpb.Attributes
	var err error

	if attrs, err = parseAttributes(rootArgs); err != nil {
		errorf(err.Error())
		return
	}

	if len(args) == 0 {
		errorf("Message is missing.")
		return
	}

	var cs *clientState
	if cs, err = createAPIClient(rootArgs.mixerAddress); err != nil {
		errorf("Unable to establish connection to %s: %v", rootArgs.mixerAddress, err)
		return
	}
	defer deleteAPIClient(cs)

	var stream mixerpb.Mixer_ReportClient
	if stream, err = cs.client.Report(context.Background()); err != nil {
		errorf("Report RPC failed: %v", err)
		return
	}

	// send the request
	request := mixerpb.ReportRequest{RequestIndex: 0, AttributeUpdate: attrs}

	/* reinstate this code ASAP

	request.LogEntries = make([]*mixerpb.LogEntry, len(args))
	for i, arg := range args {
		now, _ := ptypes.TimestampProto(time.Now())
		request.LogEntries[i] = &mixerpb.LogEntry{
			Severity:  mixerpb.LogEntry_DEFAULT,
			Timestamp: now,
			Payload: &mixerpb.LogEntry_TextPayload{
				TextPayload: arg,
			},
		}
	}
	*/
	if err = stream.Send(&request); err != nil {
		errorf("Failed to send Report RPC: %v", err)
		return
	}

	var response *mixerpb.ReportResponse
	response, err = stream.Recv()
	if err == io.EOF {
		errorf("Got no response from Report RPC")
		return
	} else if err != nil {
		errorf("Failed to receive a response from Report RPC: %v", err)
		return
	}
	stream.CloseSend()

	fmt.Printf("Report RPC returned %v\n", response.Result)
}
