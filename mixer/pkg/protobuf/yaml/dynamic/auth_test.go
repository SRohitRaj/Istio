// Copyright 2018 Istio Authors.
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

package dynamic

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gogo/protobuf/types"

	"istio.io/api/mixer/adapter/model/v1beta1"
	attributeV1beta1 "istio.io/api/policy/v1beta1"
	policypb "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/template/metric"
	spy "istio.io/istio/mixer/test/spybackend"
)

type AuthMode int

const (
	PLAINTEXT AuthMode = iota
	TLS
	MTLS
)

func TestAuth(t *testing.T) {
	// prep for test.
	// write token file.
	tokenPath := path.Join(os.TempDir(), "token")
	writeTestToken(t, tokenPath)
	defer func() {
		_ = os.Remove(tokenPath)
	}()
	// start oauth server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		cn, cs, ok := r.BasicAuth()
		wn := "clientid"
		ws := "test-token"
		if !ok {
			t.Fatal("cannot parse oauth request basic auth header")
		}
		if cn != wn || cs != ws {
			t.Errorf("client name and client secret are not expected, want %v and %v, got %v and %v", wn, ws, cn, cs)
		}
		w.Write([]byte("access_token=test-token&token_type=bearer"))
	}))
	defer ts.Close()

	testcases := []struct {
		name                  string
		mode                  AuthMode
		headerKey             string
		headerToken           string
		configErrorMessage    string
		handshakeErrorMessage string
		adapterCrt            string
		adapterKey            string
		authCfg               *policypb.Authentication
	}{
		{
			name: "no auth",
			mode: PLAINTEXT,
		},
		{
			name: "mtls",
			mode: MTLS,
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Mutual{
					Mutual: &policypb.Mutual{
						PrivateKey:        "../testdata/auth/mixer.key",
						ClientCertificate: "../testdata/auth/mixer.crt",
						CaCertificates:    "../testdata/auth/ca.pem",
					},
				},
			},
			adapterKey: "../testdata/auth/adapter.key",
			adapterCrt: "../testdata/auth/adapter.crt",
		},
		{
			name: "mtls non mixer san",
			mode: MTLS,
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Mutual{
					Mutual: &policypb.Mutual{
						PrivateKey:        "../testdata/auth/mixer.key",
						ClientCertificate: "../testdata/auth/mixer.crt",
						CaCertificates:    "../testdata/auth/ca.pem",
					},
				},
			},
			adapterKey:            "../testdata/auth/bad.key",
			adapterCrt:            "../testdata/auth/bad.crt",
			handshakeErrorMessage: "cert SAN [spiffe://cluster.local/ns/istio-system/sa/bad-service-account] is not whitelisted",
		},
		{
			name: "mtls untrusted certs",
			mode: MTLS,
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Mutual{
					Mutual: &policypb.Mutual{
						PrivateKey:        "../testdata/auth/mixer.key",
						ClientCertificate: "../testdata/auth/mixer.crt",
						CaCertificates:    "../testdata/auth/ca.pem",
					},
				},
			},
			adapterKey:            "../testdata/auth/untrusted.key",
			adapterCrt:            "../testdata/auth/untrusted.crt",
			handshakeErrorMessage: "certificate signed by unknown authority",
		},
		{
			name:        "tls only token",
			mode:        TLS,
			headerKey:   "authorization",
			headerToken: "test-token",
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "../testdata/auth/ca.pem",
						TokenSource:    &policypb.Tls_TokenPath{TokenPath: tokenPath},
						TokenType:      &policypb.Tls_AuthHeader_{AuthHeader: policypb.PLAIN},
					},
				},
			},
		},
		{
			name:        "tls authorization header",
			mode:        TLS,
			headerKey:   "authorization",
			headerToken: "Bearer test-token",
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "../testdata/auth/ca.pem",
						TokenSource:    &policypb.Tls_TokenPath{TokenPath: tokenPath},
						TokenType:      &policypb.Tls_AuthHeader_{AuthHeader: policypb.BEARER},
					},
				},
			},
		},
		{
			name:        "tls authorization custom header",
			mode:        TLS,
			headerKey:   "x-api-key",
			headerToken: "test-token",
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "../testdata/auth/ca.pem",
						TokenSource:    &policypb.Tls_TokenPath{TokenPath: tokenPath},
						TokenType:      &policypb.Tls_CustomHeader{CustomHeader: "x-api-key"},
					},
				},
			},
		},
		{
			name: "tls authorization error no token type",
			mode: TLS,
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "../testdata/auth/ca.pem",
						TokenSource:    &policypb.Tls_TokenPath{TokenPath: tokenPath},
					},
				},
			},
			configErrorMessage: "cannot get grpc per rpc credentials token type should be specified",
		},
		{
			name: "tls authorization missing ca",
			mode: TLS,
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "some/ca/path",
						TokenSource:    &policypb.Tls_TokenPath{TokenPath: tokenPath},
						TokenType:      &policypb.Tls_AuthHeader_{AuthHeader: policypb.BEARER},
					},
				},
			},
			configErrorMessage: "ca cert cannot be load open some/ca/path: no such file or directory",
		},
		{
			name: "tls authorization no token source",
			mode: TLS,
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "../testdata/auth/ca.pem",
						TokenType:      &policypb.Tls_CustomHeader{CustomHeader: "x-api-key"},
					},
				},
			},
			configErrorMessage: "cannot get grpc per rpc credentials unexpected tls token source type",
		},
		{
			name:        "tls oauth token",
			mode:        TLS,
			headerKey:   "authorization",
			headerToken: "Bearer test-token",
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "../testdata/auth/ca.pem",
						TokenSource: &policypb.Tls_Oauth{
							Oauth: &policypb.OAuth{
								ClientId:     "clientid",
								ClientSecret: tokenPath,
								TokenUrl:     ts.URL,
							},
						},
					},
				},
			},
		},
		{
			name:        "tls oauth token no secret",
			mode:        TLS,
			headerKey:   "authorization",
			headerToken: "Bearer test-token",
			authCfg: &policypb.Authentication{
				AuthType: &policypb.Authentication_Tls{
					Tls: &policypb.Tls{
						CaCertificates: "../testdata/auth/ca.pem",
						TokenSource: &policypb.Tls_Oauth{
							Oauth: &policypb.OAuth{
								ClientId: "clientid",
								TokenUrl: ts.URL,
							},
						},
					},
				},
			},
			configErrorMessage: "cannot get grpc per rpc credentials oauth secret cannot be empty",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			var s spy.Server
			var err error
			args := getServerArgs(tt.mode, tt.headerKey, tt.headerToken, tt.adapterKey, tt.adapterCrt)
			if s, err = spy.NewNoSessionServer(args); err != nil {
				t.Fatalf("unable to start Spy %v", err)
			}
			s.Run()
			defer func() {
				_ = s.Close()
			}()
			t.Logf("Started server at: %v", s.Addr())

			adapterConfig := &types.Any{
				TypeUrl: "@abc",
				Value:   []byte("abcd"),
			}
			metricDi := loadInstance(t, "metric", "template/metric/template_handler_service.descriptor_set",
				v1beta1.TEMPLATE_VARIETY_REPORT)
			h, err := BuildHandler("spy",
				&attributeV1beta1.Connection{
					Address:        s.Addr().String(),
					Authentication: tt.authCfg,
				},
				false, adapterConfig,
				[]*TemplateConfig{metricDi}, true)
			if err != nil {
				if tt.configErrorMessage != "" {
					if !strings.Contains(err.Error(), tt.configErrorMessage) {
						t.Errorf("want %v in error message, got %v", tt.configErrorMessage, err.Error())
					}
					return
				}
				t.Fatalf("cannot connect to remote handler %v", err)
			}
			mi := buildMetricInst(t)
			if err := h.HandleRemoteReport(context.Background(), []*adapter.EncodedInstance{mi}); err != nil {
				if tt.handshakeErrorMessage != "" {
					if !strings.Contains(err.Error(), tt.handshakeErrorMessage) {
						t.Errorf("want %v in error message, got %v", tt.handshakeErrorMessage, err.Error())
					}
					return
				}
				t.Errorf("get rpc error %v", err)
			}
		})
	}
}

func getServerArgs(auth AuthMode, headerKey, headerToken, key, crt string) *spy.Args {
	args := spy.DefaultArgs()
	switch auth {
	case TLS:
		args.Behavior.RequireTLS = true
		args.Behavior.HeaderKey = headerKey
		args.Behavior.HeaderToken = headerToken
	case MTLS:
		args.Behavior.RequireMTls = true
		args.Behavior.InsecureSkipVerification = true
	}
	if key == "" {
		key = "../testdata/auth/adapter.key"
	}
	if crt == "" {
		crt = "../testdata/auth/adapter.crt"
	}
	args.Behavior.KeyPath = key
	args.Behavior.CredsPath = crt
	args.Behavior.CertPath = "../testdata/auth/ca.pem"
	return args
}

func writeTestToken(t *testing.T, tp string) {
	var err error
	f, err := os.Create(tp)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("test-token")
	f.Close()
}

func buildMetricInst(t *testing.T) *adapter.EncodedInstance {
	metricDi := loadInstance(t, "metric", "template/metric/template_handler_service.descriptor_set",
		v1beta1.TEMPLATE_VARIETY_REPORT)

	minst := &metric.InstanceMsg{
		Name: metricDi.Name,
		Value: &attributeV1beta1.Value{
			Value: &attributeV1beta1.Value_StringValue{
				StringValue: "aaaaaaaaaaaaaaaa",
			},
		},
	}
	minstBa, _ := minst.Marshal()
	return &adapter.EncodedInstance{
		Name: metricDi.Name,
		Data: minstBa,
	}
}
