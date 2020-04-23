// Copyright 2020 Istio Authors
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

package util

import (
	"bytes"
	"io"
	"testing"
)

func TestProgressLog(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	testBuf := io.Writer(buf)
	testWriter = &testBuf
	expected := ""
	expect := func(e string) {
		t.Helper()
		// In buffer mode we don't overwrite old data, so we are constantly appending to the expected
		newExpected := expected + "\n" + e
		if newExpected != buf.String() {
			t.Fatalf("expected '%v', got '%v'", e, buf.String())
		}
		expected = newExpected
	}

	p := NewProgressLog()
	foo := p.NewComponent("foo")
	foo.ReportProgress()
	expect(`- Processing resources for components foo.`)

	bar := p.NewComponent("bar")
	bar.ReportProgress()
	// string buffer won't rewrite, so we append
	expect(`- Processing resources for components bar, foo.`)
	bar.ReportProgress()
	expect(`- Processing resources for components bar, foo.`)
	bar.ReportProgress()
	expect(`  Processing resources for components bar, foo.`)

	bar.ReportWaiting([]string{"deployment"})
	expect(`- Processing resources for components bar, foo. Waiting for deployment`)

	bar.ReportError("some error")
	expect(`✘ Component bar encountered an error: some error`)

	foo.ReportProgress()
	expect(`- Processing resources for components foo.`)

	foo.ReportFinished()
	expect(`✔ Component foo installed`)
}
