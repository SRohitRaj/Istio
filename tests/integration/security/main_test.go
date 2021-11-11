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
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/tests/integration/security/util"
)

var (
	ist  istio.Instance
	apps = &util.EchoDeployments{}
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
	// use the server.crt generated by following https://github.com/istio/istio/blob/master/samples/jwt-server/testdata/README.MD
	cfg.ControlPlaneValues = `
values:
  pilot: 
    jwksResolverExtraRootCA: |
      MIIDNzCCAh+gAwIBAgIUHwGzTNIiabqPqcvi1zMOKnI58TQwDQYJKoZIhvcNAQEL
      BQAwRjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkFaMRMwEQYDVQQKDApBY21lLCBJ
      bmMuMRUwEwYDVQQDDAxBY21lIFJvb3QgQ0EwHhcNMjExMTExMjMzNzEwWhcNMzEx
      MTA5MjMzNzEwWjA/MQswCQYDVQQGEwJVUzELMAkGA1UECAwCQVoxEzARBgNVBAoM
      CkFjbWUsIEluYy4xDjAMBgNVBAMMBSouY29tMIIBIjANBgkqhkiG9w0BAQEFAAOC
      AQ8AMIIBCgKCAQEA88fGFszZuZxPJPS6PG0XQEwZfwFxizF3oUrlfFoBJL9UZprD
      XF3nS0W8PVd41mzfrTY6UaX8WT2Focj1/+L9W99zKJr8UOa2ev99zMH2C45rRHDP
      bdNpmDzVVb77t3kZH2CJoCmDkIJ87kJstKTnfScAAw4EXGjzNljY1eayRzZdSoMn
      PGmxOa0+deWa32hnrWA1RPNhnjQPOStm/yzqdTkTjF7GzCQL6KBF9GgiIdX1RhgQ
      eKJNEy/7FGeOFjEkVq/gTYOs5hbcAUsBhztv44Rq4wfXjNGFMMYAM8lKPOaGcWoH
      00iNf2XGTCx/FfBF+XfoSg1yzraYG/kkkJ0gpQIDAQABoyQwIjAgBgNVHREEGTAX
      ggpqd3Qtc2VydmVygglsb2NhbGhvc3QwDQYJKoZIhvcNAQELBQADggEBABjiQmNf
      LaIQCH7JXSZMBPAWJDBiTySSgNVPLMT1rB60NU/ycrpmaj7JJwH8jW2m1J1HCA27
      G8e+foatiHKnMAxqvlpxrtDheZKTliZerHNNxkQGfvUccmLoHshXFUkYGtnEFrIg
      hK3ZYQkGh3XvDaPswpMvMpg33RSWijbVZhrVmHArgM4EEwyKLz3sR5DZx7Z+cOLd
      IRq6eFA3V34ReKononkGZgKzDwmg21ZCKGYeHcujv5RHREPmm7vzS+VlPKtgAYeJ
      FO+KS+6eG86bQiLc9EfjYNHx3WYmtojcyRkO9XGj8AGaILE1SEziW8zsWdYlqK2H
      nJQcWI2cTD0qwoI=
    env: 
      PILOT_JWT_ENABLE_REMOTE_JWKS: true
meshConfig:
  accessLogEncoding: JSON
  accessLogFile: /dev/stdout
  defaultConfig:
    gatewayTopology:
      numTrustedProxies: 1`
}
