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

package diag

import (
	"istio.io/istio/pkg/config/resource"
)

var _ resource.Origin = &testOrigin{}
var _ resource.Reference = &testReference{}

type testOrigin struct {
	name     string
	ref      resource.Reference
	fieldMap map[string]int
}

func (o testOrigin) FriendlyName() string {
	return o.name
}

func (o testOrigin) Namespace() resource.Namespace {
	return ""
}

func (o testOrigin) Reference() resource.Reference {
	return o.ref
}

func (o testOrigin) GetFieldMap() map[string]int {
	return o.fieldMap
}

type testReference struct {
	name string
}

func (r testReference) String() string {
	return r.name
}

func MockResource(name string) *resource.Instance {
	return &resource.Instance{
		Metadata: resource.Metadata{
			FullName: resource.NewShortOrFullName("default", name),
		},
		Origin: testOrigin{name: name},
	}
}
