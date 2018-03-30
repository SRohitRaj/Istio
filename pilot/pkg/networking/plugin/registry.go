// Copyright 2018 Istio Authors.
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

package plugin

import (
	"istio.io/istio/pilot/pkg/networking/plugin/authn"
)

// New returns a list of plugin instance handles. Each plugin implements the plugin.Callbacks interfaces
func New() []Callbacks {
	plugins := make([]Callbacks, 0)
	plugins = append(plugins, authn.NewPlugin())
	// plugins = append(plugins, mixer.NewPlugin())
	// plugins = append(plugins, apim.NewPlugin())
	return plugins
}
