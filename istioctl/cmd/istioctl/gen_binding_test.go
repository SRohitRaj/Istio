// Copyright 2017 Istio Authors.
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
	"bytes"
	"strings"
	"testing"

	"istio.io/istio/pilot/test/util"
)

type genBindingTestCase struct {
	name string
	args []string

	// Typically use one of the three
	expectedOutput    string // Expected constant output
	expectedSubstring string // String output is expected to contain
	goldenFilename    string // Expected output stored in golden file

	wantException bool
}

func TestGenBinding(t *testing.T) {
	tt := []genBindingTestCase{
		{args: strings.Split("experimental gen-binding", " "),
			expectedSubstring: "Error: usage: istioctl experimental gen-binding <service:port> --cluster <ip:port> [--cluster <ip:port>]* [--vip <ip>] [--labels key1=value1,key2=value2] [--use-egress] [--egressgateway <ip:port>]", // nolint: lll
			wantException:     true,
			name:              "No args"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1.2.3.4:15443", " "),
			goldenFilename: "testdata/genbinding/reviews-1234.yaml",
			name:           "One remote, no subset"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1.2.3.4:15443 --vip 6.7.8.9", " "),
			goldenFilename: "testdata/genbinding/reviews-1234-vip.yaml",
			name:           "One remote with vip address"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1.2.3.4 --cluster 6.7.8.9", " "),
			goldenFilename: "testdata/genbinding/reviews-2remotes.yaml",
			name:           "Two remotes, no subset"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1.2.3.4:15443 --labels version=v1", " "),
			goldenFilename: "testdata/genbinding/reviews-v1.yaml",
			name:           "One remote with subset"},

		{args: strings.Split("experimental gen-binding ratings:8080 --cluster 1.2.3.4 --labels version=v1,arch=i586", " "),
			goldenFilename: "testdata/genbinding/ratings-v1-i586.yaml",
			name:           "One remote with subset"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1_2_3_4:15443", " "),
			expectedSubstring: "Error: could not create binding: invalid Name or IP address",
			wantException:     true,
			name:              "Bad cluster hostname or IP"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1.2.3.4 --use-egress", " "),
			goldenFilename: "testdata/genbinding/use-egress.yaml",
			name:           "Use egress gateway"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1.2.3.4 --egressgateway 6.7.8.9:15443", " "),
			goldenFilename: "testdata/genbinding/egressgateway.yaml",
			name:           "Specify egress gateway"},

		{args: strings.Split("experimental gen-binding reviews:9080 --cluster 1.2.3.4 --egressgateway 6.7.8.9", " "),
			goldenFilename: "testdata/genbinding/egress-gateway-ip.yaml",
			name:           "Specify egress gateway IP address only"},

		{args: strings.Split("experimental gen-binding ratings:8080 --cluster 1.2.3.4 --labels version=v1,arch=i586 --use-egress", " "),
			goldenFilename: "testdata/genbinding/ratings-v1-i586-egress.yaml",
			name:           "One remote with subset with egress"},

		{args: strings.Split("experimental gen-binding ratings:8080 --cluster 1.2.3.4 --use-egress=false --egressgateway 1234.com", " "),
			expectedSubstring: "Error: cannot combine --use-egress=false and --egressgateway",
			wantException:     true,
			name:              "egressgateway/use-egress mismatch"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Clear, because we re-use
			remoteClusters = []string{}
			vip = ""
			addressLabels = ""
			useEgress = false
			egressGateway = "istio-egressgateway.istio-system"

			verifyGenBindingTestOutput(t, tc)
		})
	}
}

func verifyGenBindingTestOutput(t *testing.T, c genBindingTestCase) {
	t.Helper()

	// Reset the flags, so genBindingCmd.Flags() doesn't
	// include arguments from previous run
	pflags := genBindingCmd.PersistentFlags()
	genBindingCmd.ResetFlags()
	genBindingCmd.PersistentFlags().AddFlagSet(pflags)

	var out bytes.Buffer
	rootCmd.SetOutput(&out)
	rootCmd.SetArgs(c.args)

	fErr := rootCmd.Execute()
	output := out.String()

	if c.expectedOutput != "" && c.expectedOutput != output {
		t.Fatalf("Unexpected output for 'istioctl %s'\n got: %q\nwant: %q", strings.Join(c.args, " "), output, c.expectedOutput)
	}

	if c.expectedSubstring != "" && !strings.Contains(output, c.expectedSubstring) {
		t.Fatalf("Output didn't match for 'istioctl %s'\n got %v\nwant: %v", strings.Join(c.args, " "), output, c.expectedSubstring)
	}

	if c.goldenFilename != "" {
		util.CompareContent([]byte(output), c.goldenFilename, t)
	}

	if c.wantException {
		if fErr == nil {
			t.Fatalf("Wanted an exception for 'istioctl %s', didn't get one, output was %q",
				strings.Join(c.args, " "), output)
		}
	} else {
		if fErr != nil {
			t.Fatalf("Unwanted exception for 'istioctl %s': %v", strings.Join(c.args, " "), fErr)
		}
	}
}
