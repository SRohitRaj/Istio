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

// THIS FILE IS AUTOMATICALLY GENERATED.

package reportnothing

import (
	"context"

	"istio.io/istio/mixer/pkg/adapter"
)

// The `reportnothing` template represents an empty block of data, which can be useful
// in different testing scenarios.
//
// Example config:
//
// ```yaml
// apiVersion: "config.istio.io/v1alpha2"
// kind: instance
// metadata:
//   name: reportrequest
//   namespace: istio-system
// spec:
//   compiledTemplate: reportnothing
// ```

// Fully qualified name of the template
const TemplateName = "reportnothing"

// Instance is constructed by Mixer for the 'reportnothing' template.
//
// ReportNothing represents an empty block of data that is used for Report-capable
// adapters which don't require any parameters. This is primarily intended for testing
// scenarios.
type Instance struct {
	// Name of the instance as specified in configuration.
	Name string
}

// HandlerBuilder must be implemented by adapters if they want to
// process data associated with the 'reportnothing' template.
//
// Mixer uses this interface to call into the adapter at configuration time to configure
// it with adapter-specific configuration as well as all template-specific type information.
type HandlerBuilder interface {
	adapter.HandlerBuilder

	// SetReportNothingTypes is invoked by Mixer to pass the template-specific Type information for instances that an adapter
	// may receive at runtime. The type information describes the shape of the instance.
	SetReportNothingTypes(map[string]*Type /*Instance name -> Type*/)
}

// Handler must be implemented by adapter code if it wants to
// process data associated with the 'reportnothing' template.
//
// Mixer uses this interface to call into the adapter at request time in order to dispatch
// created instances to the adapter. Adapters take the incoming instances and do what they
// need to achieve their primary function.
//
// The name of each instance can be used as a key into the Type map supplied to the adapter
// at configuration time via the method 'SetReportNothingTypes'.
// These Type associated with an instance describes the shape of the instance
type Handler interface {
	adapter.Handler

	// HandleReportNothing is called by Mixer at request time to deliver instances to
	// to an adapter.
	HandleReportNothing(context.Context, []*Instance) error
}
