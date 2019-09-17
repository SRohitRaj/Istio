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

package direct

import (
	"istio.io/istio/galley/pkg/config/collection"
	"istio.io/istio/galley/pkg/config/event"
	"istio.io/istio/galley/pkg/config/processing"
	"istio.io/istio/galley/pkg/config/processor/transforms"
)

// Create a new Direct transformer.
func GetInfo(mapping map[collection.Name]collection.Name) []*transforms.Info {
	var result []*transforms.Info

	for k, v := range mapping {
		from := k
		to := v
		inputs := collection.Names{from}
		outputs := collection.Names{to}

		createFn := func(processing.ProcessorOptions) event.Transformer {

			return event.NewFnTransform(
				inputs,
				outputs,
				nil,
				nil,
				func(e event.Event, h event.Handler) {
					e = e.WithSource(to)
					h.Handle(e)
				})
		}
		result = append(result, transforms.NewInfo(inputs, outputs, createFn))
	}
	return result
}
