// Copyright 2018 Istio Authors
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

package server

import (
	"github.com/gogo/protobuf/proto"
)

// TODO - this will be replaced by the common envelope and metadata API from istoi/api
type Corev1Metadata struct {
	Name string
}

func (md *Corev1Metadata) GetName() string {
	if md != nil {
		return md.Name
	}
	return ""
}

// Resource defines the common protobuf interface for resources
// delivered through MCP.
type Resource interface {
	proto.Message
	GetMetadata() *Corev1Metadata
}
