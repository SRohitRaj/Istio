// Copyright 2019 Istio Authors
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

package framework

import (
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/core"
)

func TestMain(m *testing.M) {
	framework.Main("framework_test", m, framework.RequireEnvironment(core.Kube))
}

func TestBasic(t *testing.T) {
	ctx := framework.NewContext(t)
	defer ctx.Done(t)

	// Ensure that Istio can be deployed. If you're breaking this, you'll break many integration tests.
	_, err := istio.Deploy(ctx, nil)
	if err != nil {
		t.Fatalf("Istio should have deployed: %v", err)
	}
}
