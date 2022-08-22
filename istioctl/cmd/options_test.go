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

package cmd

import (
	"bytes"
	"regexp"
	"testing"
)

// nolint: lll
var expectedOutput = `The following options can be passed to any command:

      --context: The name of the kubeconfig context to use
  -i, --istioNamespace: Istio system namespace
  -c, --kubeconfig: Kubernetes configuration file
      --log_as_json: Whether to format output as JSON or in plain console-friendly format
      --log_caller: Comma-separated list of scopes for which to include caller information, scopes can be any of \[.*\]
      --log_output_level: Comma-separated minimum per-scope logging level of messages to output, in the form of <scope>:<level>,<scope>:<level>,... where scope can be one of \[.*\] and level can be one of \[.*\]
      --log_stacktrace_level: Comma-separated minimum per-scope logging level at which stack traces are captured, in the form of <scope>:<level>,<scope:level>,... where scope can be one of \[.*\] and level can be one of \[.*\]
      --log_target: The set of paths where to output the log. This can be any path as well as the special values stdout and stderr
  -n, --namespace: Config namespace
`

func TestLogHelp(t *testing.T) {
	var out bytes.Buffer
	rootCmd := GetRootCmd([]string{"options"})
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)

	fErr := rootCmd.Execute()
	if fErr != nil {
		t.Fatalf("options failed with %v and %q\n", fErr, out.String())
	}
	if !regexp.MustCompile(expectedOutput).Match(out.Bytes()) {
		t.Fatalf("'istioctl options' expected output\n%s\n  got\n%s",
			expectedOutput, out.String())
	}
}
