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

package testing

import (
	gt "testing"

	"istio.io/mixer/pkg/adapter"
)

// RegisterFunc is the function that registers adapters into the supplied registry
type RegisterFunc func(adapter.Registrar) error

type fakeRegistrar struct {
	registrations int
}

func (r *fakeRegistrar) RegisterCheckList(adapter.ListCheckerAdapter) error {
	r.registrations++
	return nil
}

func (r *fakeRegistrar) RegisterDeny(adapter.DenyCheckerAdapter) error {
	r.registrations++
	return nil
}

func (r *fakeRegistrar) RegisterLogger(adapter.LoggerAdapter) error {
	r.registrations++
	return nil
}

func (r *fakeRegistrar) RegisterQuota(adapter.QuotaAdapter) error {
	r.registrations++
	return nil
}

// TestAdapterInvariants ensures that adapters implement expected semantics.
func TestAdapterInvariants(a adapter.Adapter, r RegisterFunc, t *gt.T) {
	if a.Name() == "" {
		t.Error("Name() => all adapters need names")
	}

	if a.Description() == "" {
		t.Errorf("Description() => adapter '%s' doesn't provide a valid description", a.Name())
	}

	c := a.DefaultConfig()
	if err := a.ValidateConfig(c); err != nil {
		t.Errorf("ValidateConfig() => adapter '%s' can't validate its default configuration: %v", a.Name(), err)
	}

	if err := a.Close(); err != nil {
		t.Errorf("Close() => adapter '%s' fails to close when used with its default configuration: %v", a.Name(), err)
	}

	fr := &fakeRegistrar{}
	if err := r(fr); err != nil {
		t.Errorf("Register() => adapter '%s' didn't register properly: %v", a.Name(), err)
	}

	if fr.registrations < 1 {
		t.Errorf("Register() => adapter '%s' didn't register anything", a.Name())
	}
}
