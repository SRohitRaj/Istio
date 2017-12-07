// Copyright 2017 Istio Authors
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

package istio_mixer_adapter_metricentry

import (
	"context"
	"net"
	"time"

	"istio.io/istio/mixer/pkg/adapter"
)

//
// Overview of what metric is etc..
//
// Additional overview of what metric is etc..

// Fully qualified name of the template
const TemplateName = "metricentry"

// Instance is constructed by Mixer for the 'metricentry' template.
//
// metric template is ..
// aso it is...
type Instance struct {
	// Name of the instance as specified in configuration.
	Name string

	// value is ...
	Value interface{}

	// dimensions are ...
	Dimensions map[string]interface{}

	Int64Primitive int64

	BoolPrimitive bool

	DoublePrimitive float64

	StringPrimitive string

	AnotherValueType interface{}

	DimensionsFixedInt64ValueDType map[string]int64

	TimeStamp time.Time

	Duration time.Duration

	IpAddr net.IP

	DnsName adapter.DNSName

	EmailAddr adapter.EmailAddress

	Uri adapter.URI

	Res3List []*Resource3

	Res3Map map[string]*Resource3
}

type Resource1 struct {
	Str string

	SelfRefRes1 *Resource1

	ResRef2 *Resource2
}

type Resource2 struct {
	Str string

	Res3 *Resource3

	Res3List []*Resource3

	Res3Map map[string]*Resource3
}

// resource3 comment
type Resource3 struct {

	// value is ...
	Value interface{}

	// dimensions are ...
	Dimensions map[string]interface{}

	Int64Primitive int64

	BoolPrimitive bool

	DoublePrimitive float64

	StringPrimitive string

	AnotherValueType interface{}

	DimensionsFixedInt64ValueDType map[string]int64

	TimeStamp time.Time

	Duration time.Duration
}

// HandlerBuilder must be implemented by adapters if they want to
// process data associated with the 'metricentry' template.
//
// Mixer uses this interface to call into the adapter at configuration time to configure
// it with adapter-specific configuration as well as all template-specific type information.
type HandlerBuilder interface {
	adapter.HandlerBuilder

	// SetMetricEntryTypes is invoked by Mixer to pass the template-specific Type information for instances that an adapter
	// may receive at runtime. The type information describes the shape of the instance.
	SetMetricEntryTypes(map[string]*Type /*Instance name -> Type*/)
}

// Handler must be implemented by adapter code if it wants to
// process data associated with the 'metricentry' template.
//
// Mixer uses this interface to call into the adapter at request time in order to dispatch
// created instances to the adapter. Adapters take the incoming instances and do what they
// need to achieve their primary function.
//
// The name of each instance can be used as a key into the Type map supplied to the adapter
// at configuration time via the method 'SetMetricEntryTypes'.
// These Type associated with an instance describes the shape of the instance
type Handler interface {
	adapter.Handler

	// HandleMetricEntry is called by Mixer at request time to deliver instances to
	// to an adapter.
	HandleMetricEntry(context.Context, []*Instance) error
}
