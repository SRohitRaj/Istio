// Copyright 2016 Google Inc.
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

package attribute

import (
	mixerpb "istio.io/api/mixer/api/v1"
)

// Tracker is responsible for tracking a set of live attributes over time.
//
// An instance of this type is created for every gRPC stream incoming to the
// mixer. The instance tracks a current dictionary along with a set of
// attribute contexts.
type Tracker interface {
	// Update refreshes the set of attributes tracked based on an incoming proto.
	//
	// This returns the Context that can be used to query the current
	// set of attributes.
	//
	// If this returns a non-nil error, it indicates there was a problem in the
	// supplied Attributes struct. When this happens, the state of the
	// Context is left unchanged.
	Update(attrs *mixerpb.Attributes) (MutableBag, error)
}

type tracker struct {
	dictionaries *dictionaries

	// all active attribute contexts
	contexts map[int32]*bag

	// the current live dictionary
	currentDictionary dictionary
}

func newTracker(dictionaries *dictionaries) Tracker {
	return &tracker{dictionaries, make(map[int32]*bag), nil}
}

func (at *tracker) Update(attrs *mixerpb.Attributes) (MutableBag, error) {
	// replace the dictionary if requested
	if len(attrs.Dictionary) > 0 {
		at.currentDictionary = at.dictionaries.Intern(attrs.Dictionary)
	}

	// find the context or create it if needed
	ab := at.contexts[attrs.AttributeContext]
	if ab == nil {
		ab = &bag{}
		at.contexts[attrs.AttributeContext] = ab
	}

	if err := ab.update(at.currentDictionary, attrs); err != nil {
		return nil, err
	}

	// to maintain the integrity of the API-level attribute protocol,
	// we need to ensure that the rest of the mixer's processing pipeline
	// doesn't mutate the set of protocol-level attributes. We do this
	// by returning a child bag.
	return ab.Child(), nil
}
