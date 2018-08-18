// Copyright 2018 Istio Authors. All Rights Reserved.
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

package client_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"istio.io/istio/mixer/test/client/env"
)

// The Istio authn envoy config
// nolint
const authnConfig = `
- name: jwt-auth
  config: {
     "rules": [ 
       {
         "issuer": "issuer@foo.com",
         "local_jwks": {
           "inline_string": '{ "keys" : [ {"e":   "AQAB", "kid": "DHFbpoIUqrY8t2zpA2qXfCmr5VO5ZEr4RzHU_-envvQ", "kty": "RSA","n":   "xAE7eB6qugXyCAG3yhh7pkDkT65pHymX-P7KfIupjf59vsdo91bSP9C8H07pSAGQO1MV_xFj9VswgsCg4R6otmg5PV2He95lZdHtOcU5DXIg_pbhLdKXbi66GlVeK6ABZOUW3WYtnNHD-91gVuoeJT_DwtGGcp4ignkgXfkiEm4sw-4sfb4qdt5oLbyVpmW6x9cfa7vs2WTfURiCrBoUqgBo_-4WTiULmmHSGZHOjzwa8WtrtOQGsAFjIbno85jp6MnGGGZPYZbDAa_b3y5u-YpW7ypZrvD8BgtKVjgtQgZhLAGezMt0ua3DRrWnKqTZ0BJ_EyxOGuHJrLsn00fnMQ"}]}'
         },
         "audiences": ["aud1"],
         "forward_payload_header": "test-jwt-payload-output"
       }
     ]
  }
- name: istio_authn
  config: {
    "policy": {
      "peers": [
        {
          "jwt": {
            "issuer": "issuer@foo.com",
            "jwks_uri": "http://localhost:8081/"
          }
        }
      ],
      "principal_binding": 1
    }
  }
`

const secIstioAuthUserInfoHeaderKey = "sec-istio-auth-jwt-output"

const secIstioAuthUserinfoHeaderValue = `{"aud":"aud1","exp":20000000000,` +
	`"iat":1500000000,"iss":"issuer@foo.com","some-other-string-claims":"some-claims-kept",` +
	`"sub":"sub@foo.com"}`

// jwt is formed from the value in secIstioAuthUserinfoHeaderValue
const jwt = "eyJhbGciOiJSUzI1NiIsImtpZCI6IkRIRmJwb0lVcXJZOHQyenBBMnFYZ" +
	"kNtcjVWTzVaRXI0UnpIVV8tZW52dlEiLCJ0eXAiOiJKV1QifQ.eyJhdWQ" +
	"iOiJhdWQxIiwiZXhwIjoyMDAwMDAwMDAwMCwiaWF0IjoxNTAwMDAwMDAw" +
	"LCJpc3MiOiJpc3N1ZXJAZm9vLmNvbSIsInNvbWUtb3RoZXItc3RyaW5nL" +
	"WNsYWltcyI6InNvbWUtY2xhaW1zLWtlcHQiLCJzdWIiOiJzdWJAZm9vLm" +
	"NvbSJ9.VYQdAqzlzpVBoKQMkmwm4oCX-wgMieR7rEpJiOggYocEJbEINr" +
	"ZSMas9bJ0CQXdv5UWR6NiO-p1Ko1Zol1X5Ma93Aego18vygY1K1bZ5whX" +
	"qVtbkpDe5tUaPNP58uKWsh8g3EA2Mpr1jF7RgGCYmiW_LlWJnLlBMEvbb" +
	"pkBFy43Yfzn_wpLHNBTO8cUGHGMErBeBSe2jUYmdOda1s51rGmS-CuQDL" +
	"GMeJPmc2l50AOO0tnNbSp3S3KfeyF918uDFfDRLYp7j16cx71ETXfLsrX" +
	"UkcLOLthIYGpuD0RgvLi5soHDpV_uNO8FDiOPMs8y60EUQUcuSKZZHTS_" +
	"hzONkhg"

const respExpected = "Origin authentication failed."

func TestAuthnCheckReportAttributesPeerJwtBoundToOrigin(t *testing.T) {
	s := env.NewTestSetup(env.CheckReportIstioAuthnAttributesTestPeerJwtBoundToOrigin, t)
	// In the Envoy config, principal_binding binds to origin
	s.SetFiltersBeforeMixer(authnConfig)
	// Disable the HotRestart of Envoy
	s.SetDisableHotRestart(true)

	env.SetStatsUpdateInterval(s.MfConfig(), 1)
	if err := s.SetUp(); err != nil {
		t.Fatalf("Failed to setup test: %v", err)
	}
	defer s.TearDown()

	url := fmt.Sprintf("http://localhost:%d/echo", s.Ports().ClientProxyPort)

	// Issues a GET echo request with 0 size body
	tag := "OKGet"

	// Add jwt_auth header to be consumed by Istio authn filter
	headers := map[string]string{}
	headers[secIstioAuthUserInfoHeaderKey] =
		base64.StdEncoding.EncodeToString([]byte(secIstioAuthUserinfoHeaderValue))
	headers["Authorization"] = "Bearer " + jwt

	// Principal is binded to origin, but no method specified in origin policy.
	// The request will be rejected by Istio authn filter.
	code, resp, err := env.HTTPGetWithHeaders(url, headers)
	if err != nil {
		t.Errorf("Failed in request %s: %v", tag, err)
	}
	// Verify that the http request is rejected
	if code != 401 {
		t.Errorf("Status code 401 is expected, got %d.", code)
	}
	if resp != respExpected {
		t.Errorf("Expected response: %s, got %s.", respExpected, resp)
	}
}
