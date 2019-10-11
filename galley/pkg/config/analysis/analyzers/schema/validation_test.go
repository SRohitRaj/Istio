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
package schema

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/gogo/protobuf/proto"

	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/diag"
	"istio.io/istio/galley/pkg/config/analysis/msg"
	"istio.io/istio/galley/pkg/config/meta/metadata"
	"istio.io/istio/galley/pkg/config/meta/schema/collection"
	"istio.io/istio/galley/pkg/config/resource"
	"istio.io/istio/pkg/config/schema"
)

type testContext struct {
	entries []*resource.Entry
	reports []diag.Message
}

var _ analysis.Context = &testContext{}

// Report implements analysis.Context
func (ctx *testContext) Report(c collection.Name, t diag.Message) {
	ctx.reports = append(ctx.reports, t)
}

// Find implements analysis.Context
func (ctx *testContext) Find(c collection.Name, name resource.Name) *resource.Entry { return nil }

// Exists implements analysis.Context
func (ctx *testContext) Exists(c collection.Name, name resource.Name) bool { return false }

// ForEach implements analysis.Context
func (ctx *testContext) ForEach(c collection.Name, fn analysis.IteratorFn) {
	for _, r := range ctx.entries {
		fn(r)
	}
}

// Canceled implements analysis.Context
func (ctx *testContext) Canceled() bool { return false }

func TestSchemaValidationWrapper(t *testing.T) {
	g := NewGomegaWithT(t)

	testCol := metadata.IstioNetworkingV1Alpha3Virtualservices

	m1 := &v1alpha3.VirtualService{}
	m2 := &v1alpha3.VirtualService{}

	testSchema := schema.Instance{
		Collection: testCol.String(),
		Validate: func(_, _ string, msg proto.Message) (errs error) {
			if msg == m1 {
				return nil
			}
			if msg == m2 {
				return fmt.Errorf("")
			}
			return nil
		},
	}

	a := ValidationAnalyzer{s: testSchema}
	g.Expect(a.Metadata().Inputs).To(ConsistOf(testCol))

	ctx := &testContext{
		entries: []*resource.Entry{
			{
				Item: m1,
			},
		},
	}
	a.Analyze(ctx)
	g.Expect(ctx.reports).To(BeEmpty())

	ctx = &testContext{
		entries: []*resource.Entry{
			{
				Item: m2,
			},
		},
	}
	a.Analyze(ctx)
	g.Expect(ctx.reports).To(HaveLen(1))
	g.Expect(ctx.reports[0].Type).To(Equal(msg.SchemaValidationError))
}
