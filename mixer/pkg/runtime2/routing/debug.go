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

package routing

// tableDebugInfo contains debugging information for the table.
type tableDebugInfo struct {

	// mapping from the handler entry in the routing table to the name of the handler.
	handlerEntries map[uint32]string

	inputSets map[uint32]inputSetDebugInfo
}

type inputSetDebugInfo struct {
	// match condition text for the input set.
	match string

	// The name of instances used to create builders.
	instanceNames []string
}
