//go:build integ
// +build integ

//  Copyright Istio Authors
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

package security

import (
	"os"
	"path"
	"testing"

	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/tests/integration/security/util"
)

func readCertFromFile(filename string) string {
	csrBytes, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(csrBytes)
}

var (
	ist        istio.Instance
	apps              = &util.EchoDeployments{}
	rootCACert string = readCertFromFile(path.Join(env.IstioSrc, "samples/jwt-server/testdata", "ca.crt"))
)

func TestMain(m *testing.M) {
	framework.
		NewSuite(m).
		Setup(istio.Setup(&ist, setupConfig)).
		Setup(func(ctx resource.Context) error {
			return util.SetupApps(ctx, ist, apps, true)
		}).
		Run()
}

func setupConfig(ctx resource.Context, cfg *istio.Config) {
	if cfg == nil {
		return
	}

	// command to generate certificate
	// use the generated ca.crt by following https://github.com/istio/istio/blob/master/samples/jwt-server/testdata/README.MD
	// TODO(garyan): enable the test for "PILOT_JWT_ENABLE_REMOTE_JWKS: true" as well.
	cfg.ControlPlaneValues = `
values:
  pilot: 
    jwksResolverExtraRootCA: |
      {{ rootCACert }}
    env: 
      PILOT_JWT_ENABLE_REMOTE_JWKS: false
meshConfig:
  accessLogFile: /dev/stdout`
}
