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

package denyChecker

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/code"

	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/adapter/denyChecker"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/aspect/denyChecker/config"
	"istio.io/mixer/pkg/attribute"
	"istio.io/mixer/pkg/expr"
)

const (
	kind = "istio/denyChecker"
)

type (
	manager struct{}

	aspectWrapper struct {
		adapter denyChecker.Adapter
		aspect  denyChecker.Aspect
	}
)

// NewManager returns "this" aspect Manager
func NewManager() aspect.Manager {
	return &manager{}
}

// NewAspect creates a denyChecker aspect.
func (m *manager) NewAspect(cfg *aspect.CombinedConfig, ga adapter.Adapter, env adapter.Env) (aspect.Wrapper, error) {
	aa, ok := ga.(denyChecker.Adapter)
	if !ok {
		return nil, fmt.Errorf("adapter of incorrect type; expected denyChecker.Adapter got %#v %T", ga, ga)
	}

	// TODO: convert from proto Struct to Go struct here!
	adapterCfg := aa.DefaultConfig()
	// TODO: parse cfg.Adapter.Params (*ptypes.struct) into adapterCfg
	var asp denyChecker.Aspect
	var err error

	if asp, err = aa.NewDenyChecker(env, adapterCfg); err != nil {
		return nil, err
	}

	return &aspectWrapper{
		adapter: aa,
		aspect:  asp,
	}, nil
}

func (*manager) Kind() string {
	return kind
}

func (*manager) DefaultConfig() adapter.AspectConfig {
	return &config.Params{}
}

func (*manager) ValidateConfig(c adapter.AspectConfig) (ce *adapter.ConfigErrors) {
	return
}

func (a *aspectWrapper) AdapterName() string {
	return a.adapter.Name()
}

func (a *aspectWrapper) Execute(attrs attribute.Bag, mapper expr.Evaluator) (*aspect.Output, error) {
	status := a.aspect.Deny()
	return &aspect.Output{Code: code.Code(status.Code)}, nil
}
