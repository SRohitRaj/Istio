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

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/onsi/gomega"

	"istio.io/istio/pilot/pkg/model"
)

func TestPilotDefaultDomainKubernetes(t *testing.T) {
	g := gomega.NewWithT(t)
	role := &model.Proxy{}
	role.DNSDomain = ""

	domain := getDNSDomain("default", role.DNSDomain)

	g.Expect(domain).To(gomega.Equal("default.svc.cluster.local"))
}

func TestInitCustomBinary(t *testing.T) {
	g := gomega.NewWithT(t)
	origStdout := os.Stdout
	defer func() {
		os.Stdout = origStdout
	}()
	r, w, _ := os.Pipe()
	os.Stdout = w
	initCustomBinary([]string{"echo", "-n", "output"})
	w.Close()
	out, _ := ioutil.ReadAll(r)
	g.Expect(string(out)).To(gomega.Equal("output"))
}
