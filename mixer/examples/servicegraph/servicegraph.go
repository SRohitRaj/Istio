// Copyright 2017 Google Inc.
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

// Package servicegraph defines the core model for the servicegraph service.
package servicegraph

import "io"

type (
	// Dynamic represents a service graph generated on the fly.
	Dynamic struct {
		Nodes map[string]struct{} `json:"nodes"`
		Edges []*Edge             `json:"edges"`
	}

	// Edge represents an edge in a dynamic service graph.
	Edge struct {
		Source string     `json:"source"`
		Target string     `json:"target"`
		Labels Attributes `json:"labels"`
	}

	// Attributes contain a set of annotations for an edge.
	Attributes map[string]string

	// SerializeFn provides a mechanism for writing out the service graph.
	SerializeFn func(w io.Writer, g *Dynamic) error
)

// AddEdge adds a new edge to an existing dynamic graph.
func (d *Dynamic) AddEdge(src, target string, lbls map[string]string) {
	d.Edges = append(d.Edges, &Edge{src, target, lbls})
	d.Nodes[src] = struct{}{}
	d.Nodes[target] = struct{}{}
}
