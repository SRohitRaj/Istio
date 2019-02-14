//  Copyright 2019 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package policybackend

import (
	"testing"

	"github.com/gogo/protobuf/proto"
	"istio.io/istio/pkg/test/framework2/components/environment"
	"istio.io/istio/pkg/test/framework2/components/environment/native"
	"istio.io/istio/pkg/test/framework2/resource"
	"istio.io/istio/pkg/test/framework2/runtime"
)

// Instance represents a deployed fake policy backend for Mixer.
type Instance interface {
	// DenyCheck indicates that the policy backend should deny all incoming check requests when deny is
	// set to true.
	DenyCheck(t testing.TB, deny bool)

	// ExpectReport checks that the backend has received the given report requests. The requests are consumed
	// after the call completes.
	ExpectReport(t testing.TB, expected ...proto.Message)

	// ExpectReportJSON checks that the backend has received the given report request.  The requests are
	// consumed after the call completes.
	ExpectReportJSON(t testing.TB, expected ...string)

	// CreateConfigSnippet for the Mixer adapter to talk to this policy backend.
	// The supplied name will be the name of the handler.
	CreateConfigSnippet(name string) string
}

func New(s resource.Context) (Instance, error) {
	switch s.Environment().Name() {
	case native.Name:
		return newNative(s, s.Environment().(*native.Environment))
	default:
		return nil, environment.UnsupportedEnvironment(s.Environment().Name())
	}
}

func NewOrFail(c *runtime.TestContext) Instance {
	i, err := New(c)
	if err != nil {
		c.T().Fatalf("Error creating PolicyBackend: %v", err)
	}

	return i
}
