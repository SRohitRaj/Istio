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

// Package stdioLogger provides an implementation of the mixer logger aspect
// that writes logs (serialized as JSON) to a standard stream (stdout | stderr).
package stdioLogger

import (
	"encoding/json"
	"io"
	"os"

	"istio.io/mixer/adapter/stdioLogger/config"
	"istio.io/mixer/pkg/aspect"
	"istio.io/mixer/pkg/aspect/logger"
	"istio.io/mixer/pkg/registry"

	me "github.com/hashicorp/go-multierror"
)

type (
	adapter    struct{}
	aspectImpl struct {
		logStream io.Writer
	}
)

// Register adds the stdioLogger adapter to the list of logger.Aspects known to
// mixer.
func Register(r registry.Registrar) error { return r.RegisterLogger(&adapter{}) }

func (a *adapter) Name() string { return "istio/stdioLogger" }
func (a *adapter) Description() string {
	return "Writes structured log entries to a standard I/O stream"
}
func (a *adapter) DefaultConfig() aspect.Config                             { return &config.Params{} }
func (a *adapter) Close() error                                             { return nil }
func (a *adapter) ValidateConfig(c aspect.Config) (ce *aspect.ConfigErrors) { return nil }
func (a *adapter) NewAspect(env aspect.Env, cfg aspect.Config) (logger.Aspect, error) {
	c := cfg.(*config.Params)

	w := os.Stderr
	if c.LogStream == config.Params_STDOUT {
		w = os.Stdout
	}

	return &aspectImpl{w}, nil
}

func (a *aspectImpl) Close() error { return nil }
func (a *aspectImpl) Log(entries []logger.Entry) error {
	var errors *me.Error
	for _, entry := range entries {
		if err := writeJSON(a.logStream, entry); err != nil {
			errors = me.Append(errors, err)
		}
	}

	return errors.ErrorOrNil()
}

func writeJSON(w io.Writer, le logger.Entry) error {
	return json.NewEncoder(w).Encode(le)
}
