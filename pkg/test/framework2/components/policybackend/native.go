//  Copyright 2018 Istio Authors
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

package policybackend

import (
	"fmt"
	"io"
	"time"

	"istio.io/istio/pkg/test/framework2/core"

	"github.com/hashicorp/go-multierror"

	"istio.io/istio/pkg/test/framework2/components/environment/native"

	"istio.io/istio/pkg/test/fakes/policy"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/istio/pkg/test/util/retry"
)

var (
	retryDelay = retry.Delay(time.Second)

	_ Instance      = &nativeComponent{}
	_ io.Closer     = &nativeComponent{}
	_ core.Resetter = &nativeComponent{}
)

type nativeComponent struct {
	id  core.ResourceID
	ctx core.Context
	env *native.Environment

	*client
	namespace string
	backend   *policy.Backend
}

// NewNativeComponent factory function for the component
func newNative(ctx core.Context, env *native.Environment) (Instance, error) {
	c := &nativeComponent{
		ctx:    ctx,
		env:    env,
		client: &client{},
	}
	c.id = ctx.TrackResource(c)

	var err error
	scopes.CI.Infof("=== BEGIN: Start local PolicyBackend ===")
	defer func() {
		if err != nil {
			scopes.CI.Infof("=== FAILED: Start local PolicyBackend ===")
			_ = c.Close()
		} else {
			scopes.CI.Infof("=== SUCCEEDED: Start local PolicyBackend ===")
		}
	}()

	c.backend = policy.NewPolicyBackend(0) // auto-allocate port
	if err = c.backend.Start(); err != nil {
		return nil, err
	}

	var ctl interface{}
	ctl, err = retry.Do(func() (interface{}, bool, error) {
		c, err := policy.NewController(fmt.Sprintf(":%d", c.backend.Port()))
		if err != nil {
			scopes.Framework.Debugf("error while connecting to the PolicyBackend controller: %v", err)
			return nil, false, err
		}
		return c, true, nil
	}, retryDelay)
	if err != nil {
		return nil, err
	}
	c.client.controller = ctl.(*policy.Controller)

	return c, nil
}

func (c *nativeComponent) CreateConfigSnippet(name string) string {
	return fmt.Sprintf(
		`apiVersion: "config.istio.io/v1alpha2"
kind: bypass
metadata:
  name: %s
  namespace: %s
spec:
  backend_address: 127.0.0.1:%d
`, name, c.namespace, c.backend.Port())
}

func (c *nativeComponent) Reset() error {
	if c.client != nil {
		return c.client.Reset()
	}
	return nil
}

func (c *nativeComponent) ID() core.ResourceID {
	return c.id
}

func (c *nativeComponent) Close() (err error) {
	if c.backend != nil {
		err = multierror.Append(err, c.backend.Close()).ErrorOrNil()
	}
	return
}
