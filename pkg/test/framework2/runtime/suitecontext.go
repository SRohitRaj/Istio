//  Copyright 2019 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package runtime

import (
	"fmt"
	"io/ioutil"
	"sync"

	"istio.io/istio/pkg/test/framework2/common"
	"istio.io/istio/pkg/test/framework2/components/environment"
	"istio.io/istio/pkg/test/framework2/resource"
	"istio.io/istio/pkg/test/scopes"
)

// SuiteContext contains suite-level items used during runtime.
type SuiteContext struct {
	settings    *common.Settings
	environment environment.Instance

	contextMu    sync.Mutex
	contextNames map[string]struct{}

	// context-level resources
	globalScope *scope
}

var _ resource.Context = &SuiteContext{}

func newSuiteContext(s *common.Settings, envFn environment.FactoryFn) (*SuiteContext, error) {
	scopeID := fmt.Sprintf("[suite(%s)]", s.TestID)

	c := &SuiteContext{
		settings:    s,
		globalScope: newScope(scopeID, nil),

		contextNames: make(map[string]struct{}),
	}

	env, err := envFn(s.Environment, c)
	if err != nil {
		return nil, err
	}
	c.environment = env
	c.globalScope.add(env)

	return c, nil
}

// allocateContextID allocates a unique context id for TestContexts. Useful for creating unique names to help with
// debugging
func (s *SuiteContext) allocateContextID(prefix string) string {
	s.contextMu.Lock()
	defer s.contextMu.Unlock()

	candidate := prefix
	discriminator := 0
	for {
		if _, found := s.contextNames[candidate]; !found {
			s.contextNames[candidate] = struct{}{}
			return candidate
		}

		candidate = fmt.Sprintf("%s-%d", prefix, discriminator)
		discriminator++
	}
}

// TrackResource adds a new resource to track to the context at this level.
func (s *SuiteContext) TrackResource(r resource.Instance) {
	s.globalScope.add(r)
}

// Environment implements ResourceContext
func (s *SuiteContext) Environment() environment.Instance {
	return s.environment
}

// Settings returns the current runtime.Settings.
func (s *SuiteContext) Settings() *common.Settings {
	return s.settings
}

func (s *SuiteContext) done() error {
	return s.globalScope.done(s.settings.NoCleanup)
}

// CreateTmpDirectory creates a new temporary directory with the given prefix.
func (s *SuiteContext) CreateTmpDirectory(prefix string) (string, error) {
	dir, err := ioutil.TempDir(s.settings.RunDir(), prefix)
	if err != nil {
		scopes.Framework.Errorf("Error creating temp dir: runID='%s', prefix='%s', workDir='%v', err='%v'",
			s.settings.RunID, prefix, s.settings.RunDir(), err)
	} else {
		scopes.Framework.Debugf("Created a temp dir: runID='%s', name='%s'", s.settings.RunID, dir)
	}

	return dir, err
}
