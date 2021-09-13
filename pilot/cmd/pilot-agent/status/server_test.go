// Copyright Istio Authors
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

package status

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/prometheus/pkg/textparse"
	"k8s.io/apimachinery/pkg/util/intstr"

	"istio.io/istio/pilot/cmd/pilot-agent/status/ready"
	"istio.io/istio/pilot/cmd/pilot-agent/status/testserver"
	"istio.io/istio/pkg/kube/apimirror"
	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/util/retry"
	"istio.io/pkg/log"
)

type handler struct{}

const (
	testHeader      = "Some-Header"
	testHeaderValue = "some-value"
	testHostValue   = "host"
)

var liveServerStats = "cluster_manager.cds.update_success: 1\nlistener_manager.lds.update_success: 1\nserver.state: 0\nlistener_manager.workers_started: 1"

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path[1:], "/")
	switch segments[0] {
	case "header":
		if r.Host != testHostValue {
			log.Errorf("Missing expected host header, got %v", r.Host)
			w.WriteHeader(http.StatusBadRequest)
		}
		if r.Header.Get(testHeader) != testHeaderValue {
			log.Errorf("Missing expected Some-Header, got %v", r.Header)
			w.WriteHeader(http.StatusBadRequest)
		}
	case "redirect":
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	case "redirect-loop":
		http.Redirect(w, r, "/redirect-loop", http.StatusMovedPermanently)
	case "remote-redirect":
		http.Redirect(w, r, "http://example.com/foo", http.StatusMovedPermanently)
	case "", "hello/sunnyvale":
		w.Write([]byte("welcome, it works"))
	case "status":
		code, _ := strconv.Atoi(segments[1])
		w.Header().Set("Location", "/")
		w.WriteHeader(code)
	default:
		return
	}
}

func TestNewServer(t *testing.T) {
	testCases := []struct {
		probe string
		err   string
	}{
		// Json can't be parsed.
		{
			probe: "invalid-prober-json-encoding",
			err:   "failed to decode",
		},
		// map key is not well formed.
		{
			probe: `{"abc": {"path": "/app-foo/health"}}`,
			err:   "invalid path",
		},
		// invalid probe type
		{
			probe: `{"/app-health/hello-world/readyz": {"exec": {"command": [ "true" ]}}}`,
			err:   "invalid prober type",
		},
		// tcp probes are valid as well
		{
			probe: `{"/app-health/hello-world/readyz": {"tcpSocket": {"port": 8888}}}`,
		},
		// probes must be tcp or http, not both
		{
			probe: `{"/app-health/hello-world/readyz": {"tcpSocket": {"port": 8888}, "httpGet": {"path": "/", "port": 7777}}}`,
			err:   "must be either httpGet or tcpSocket",
		},
		// Port is not Int typed (tcpSocket).
		{
			probe: `{"/app-health/hello-world/readyz": {"tcpSocket": {"port": "tcp"}}}`,
			err:   "must be int type",
		},
		// Port is not Int typed (httpGet).
		{
			probe: `{"/app-health/hello-world/readyz": {"httpGet": {"path": "/hello/sunnyvale", "port": "container-port-dontknow"}}}`,
			err:   "must be int type",
		},
		// A valid input.
		{
			probe: `{"/app-health/hello-world/readyz": {"httpGet": {"path": "/hello/sunnyvale", "port": 8080}},` +
				`"/app-health/business/livez": {"httpGet": {"path": "/buisiness/live", "port": 9090}}}`,
		},
		// long request timeout
		{
			probe: `{"/app-health/hello-world/readyz": {"httpGet": {"path": "/hello/sunnyvale", "port": 8080},` +
				`"initialDelaySeconds": 120,"timeoutSeconds": 10,"periodSeconds": 20}}`,
		},
		// A valid input with empty probing path, which happens when HTTPGetAction.Path is not specified.
		{
			probe: `{"/app-health/hello-world/readyz": {"httpGet": {"path": "/hello/sunnyvale", "port": 8080}},
"/app-health/business/livez": {"httpGet": {"port": 9090}}}`,
		},
		// A valid input without any prober info.
		{
			probe: `{}`,
		},
		// A valid input with probing path not starting with /, which happens when HTTPGetAction.Path does not start with a /.
		{
			probe: `{"/app-health/hello-world/readyz": {"httpGet": {"path": "hello/sunnyvale", "port": 8080}},
"/app-health/business/livez": {"httpGet": {"port": 9090}}}`,
		},
	}
	for _, tc := range testCases {
		_, err := NewServer(Options{
			KubeAppProbers: tc.probe,
		})

		if err == nil {
			if tc.err != "" {
				t.Errorf("test case failed [%v], expect error %v", tc.probe, tc.err)
			}
			continue
		}
		if tc.err == "" {
			t.Errorf("test case failed [%v], expect no error, got %v", tc.probe, err)
		}
		// error case, error string should match.
		if !strings.Contains(err.Error(), tc.err) {
			t.Errorf("test case failed [%v], expect error %v, got %v", tc.probe, tc.err, err)
		}
	}
}

func TestPprof(t *testing.T) {
	pprofPath := "/debug/pprof/cmdline"
	// Starts the pilot agent status server.
	server, err := NewServer(Options{StatusPort: 0})
	if err != nil {
		t.Fatalf("failed to create status server %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.Run(ctx)

	var statusPort uint16
	for statusPort == 0 {
		server.mutex.RLock()
		statusPort = server.statusPort
		server.mutex.RUnlock()
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/%s", statusPort, pprofPath), nil)
	if err != nil {
		t.Fatalf("[%v] failed to create request", pprofPath)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("request failed: ", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("[%v] unexpected status code, want = %v, got = %v", pprofPath, http.StatusOK, resp.StatusCode)
	}
}

func TestStats(t *testing.T) {
	cases := []struct {
		name             string
		envoy            string
		app              string
		output           string
		expectParseError bool
	}{
		{
			name: "envoy metric only",
			envoy: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0
`,
			output: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0
`,
		},
		{
			name: "app metric only",
			app: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0
# EOF
`,
			output: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0
`,
		},
		{
			name: "multiple metric",
			envoy: `# TYPE my_metric counter
my_metric{} 0
`,
			app: `# TYPE my_other_metric counter
my_other_metric{} 0
# EOF
`,
			output: `# TYPE my_other_metric counter
my_other_metric{} 0
# TYPE my_metric counter
my_metric{} 0
`,
		},
		{
			name:  "agent metric",
			envoy: ``,
			app:   ``,
			// Agent metric is dynamic, so we just check a substring of it not the actual metric
			output: `
# TYPE istio_agent_scrapes_total counter
istio_agent_scrapes_total`,
		},
		// When the application and envoy share a metric, Prometheus will fail. This negative check validates this
		// assumption.
		{
			name: "conflict metric",
			envoy: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0
`,
			app: `# TYPE my_metric counter
my_metric{} 0
# EOF
`,
			output: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0
`,
			expectParseError: true,
		},
		{
			name: "conflict metric labeled",
			envoy: `# TYPE my_metric counter
my_metric{app="foo"} 0
`,
			app: `# TYPE my_metric counter
my_metric{app="bar"} 0
# EOF
`,
			output: `# TYPE my_metric counter
my_metric{app="bar"} 0
# TYPE my_metric counter
my_metric{app="foo"} 0
`,
			expectParseError: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			envoy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte(tt.envoy)); err != nil {
					t.Fatalf("write failed: %v", err)
				}
			}))
			defer envoy.Close()
			app := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte(tt.app)); err != nil {
					t.Fatalf("write failed: %v", err)
				}
			}))
			defer app.Close()
			envoyPort, err := strconv.Atoi(strings.Split(envoy.URL, ":")[2])
			if err != nil {
				t.Fatal(err)
			}
			server := &Server{
				prometheus: &PrometheusScrapeConfiguration{
					Port: strings.Split(app.URL, ":")[2],
				},
				envoyStatsPort: envoyPort,
			}
			req := &http.Request{}
			server.handleStats(rec, req)
			if rec.Code != 200 {
				t.Fatalf("handleStats() => %v; want 200", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), tt.output) {
				t.Fatalf("handleStats() => %v; want %v", rec.Body.String(), tt.output)
			}

			parser := expfmt.TextParser{}
			mfMap, err := parser.TextToMetricFamilies(strings.NewReader(rec.Body.String()))
			if err != nil && !tt.expectParseError {
				t.Fatalf("failed to parse metrics: %v", err)
			} else if err == nil && tt.expectParseError {
				t.Fatalf("expected a prse error, got %+v", mfMap)
			}
		})
	}
}

func TestStatsContentType(t *testing.T) {
	appOpenMetrics := `# TYPE jvm info
# HELP jvm VM version info
jvm_info{runtime="OpenJDK Runtime Environment",vendor="AdoptOpenJDK",version="16.0.1+9"} 1.0
# TYPE jmx_config_reload_success counter
# HELP jmx_config_reload_success Number of times configuration have successfully been reloaded.
jmx_config_reload_success_total 0.0
jmx_config_reload_success_created 1.623984612719E9
# EOF
`
	appText004 := `# HELP jvm_info VM version info
# TYPE jvm_info gauge
jvm_info{runtime="OpenJDK Runtime Environment",vendor="AdoptOpenJDK",version="16.0.1+9",} 1.0
# HELP jmx_config_reload_failure_created Number of times configuration have failed to be reloaded.
# TYPE jmx_config_reload_failure_created gauge
jmx_config_reload_failure_created 1.624025983489E9
`
	envoy := `# TYPE my_metric counter
my_metric{} 0

# TYPE my_other_metric counter
my_other_metric{} 0
`
	cases := []struct {
		name             string
		acceptHeader     string
		expectParseError bool
	}{
		{
			name:         "openmetrics accept header",
			acceptHeader: `application/openmetrics-text; version=0.0.1,text/plain;version=0.0.4;q=0.5,*/*;q=0.1`,
		},
		{
			name:         "plaintext accept header",
			acceptHeader: string(expfmt.FmtText),
		},
		{
			name:         "empty accept header",
			acceptHeader: "",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			envoy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte(envoy)); err != nil {
					t.Fatalf("write failed: %v", err)
				}
			}))
			defer envoy.Close()
			app := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				format := expfmt.NegotiateIncludingOpenMetrics(r.Header)
				var negotiatedMetrics string
				if format == expfmt.FmtOpenMetrics {
					negotiatedMetrics = appOpenMetrics
					w.Header().Set("Content-Type", string(expfmt.FmtOpenMetrics))
				} else {
					negotiatedMetrics = appText004
					w.Header().Set("Content-Type", string(expfmt.FmtText))
				}
				if _, err := w.Write([]byte(negotiatedMetrics)); err != nil {
					t.Fatalf("write failed: %v", err)
				}
			}))
			defer app.Close()
			envoyPort, err := strconv.Atoi(strings.Split(envoy.URL, ":")[2])
			if err != nil {
				t.Fatal(err)
			}
			server := &Server{
				prometheus: &PrometheusScrapeConfiguration{
					Port: strings.Split(app.URL, ":")[2],
				},
				envoyStatsPort: envoyPort,
			}
			req := &http.Request{}
			req.Header = make(http.Header)
			req.Header.Add("Accept", tt.acceptHeader)
			server.handleStats(rec, req)
			if rec.Code != 200 {
				t.Fatalf("handleStats() => %v; want 200", rec.Code)
			}

			var format expfmt.Format
			mediaType, _, err := mime.ParseMediaType(rec.Header().Get("Content-Type"))
			if err == nil && mediaType == "application/openmetrics-text" {
				format = expfmt.FmtOpenMetrics
			} else {
				format = expfmt.FmtText
			}

			if format == expfmt.FmtOpenMetrics {
				omParser := textparse.NewOpenMetricsParser(rec.Body.Bytes())
				for {
					_, err := omParser.Next()
					if err == io.EOF {
						break
					}
					if err != nil {
						t.Fatalf("failed to parse openmetrics: %v", err)
					}
				}
			} else {
				textParser := expfmt.TextParser{}
				str := rec.Body.String()
				_, err := textParser.TextToMetricFamilies(strings.NewReader(str))
				if err != nil {
					t.Fatalf("failed to parse text metrics: %v", err)
				}
			}
		})
	}
}

func TestStatsError(t *testing.T) {
	fail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer fail.Close()
	pass := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer pass.Close()
	failPort, err := strconv.Atoi(strings.Split(fail.URL, ":")[2])
	if err != nil {
		t.Fatal(err)
	}
	passPort, err := strconv.Atoi(strings.Split(pass.URL, ":")[2])
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name  string
		envoy int
		app   int
	}{
		{"both pass", passPort, passPort},
		{"envoy pass", passPort, failPort},
		{"app pass", failPort, passPort},
		{"both fail", failPort, failPort},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{
				prometheus: &PrometheusScrapeConfiguration{
					Port: strconv.Itoa(tt.app),
				},
				envoyStatsPort: tt.envoy,
			}
			req := &http.Request{}
			rec := httptest.NewRecorder()
			server.handleStats(rec, req)
			if rec.Code != 200 {
				t.Fatalf("handleStats() => %v; want 200", rec.Code)
			}
		})
	}
}

func TestAppProbe(t *testing.T) {
	// Starts the application first.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Errorf("failed to allocate unused port %v", err)
	}
	go http.Serve(listener, &handler{})
	appPort := listener.Addr().(*net.TCPAddr).Port

	simpleHTTPConfig := KubeAppProbers{
		"/app-health/hello-world/readyz": &Prober{
			HTTPGet: &apimirror.HTTPGetAction{
				Path: "/hello/sunnyvale",
				Port: intstr.IntOrString{IntVal: int32(appPort)},
			},
		},
		"/app-health/hello-world/livez": &Prober{
			HTTPGet: &apimirror.HTTPGetAction{
				Port: intstr.IntOrString{IntVal: int32(appPort)},
			},
		},
	}
	simpleTCPConfig := KubeAppProbers{
		"/app-health/hello-world/readyz": &Prober{
			TCPSocket: &apimirror.TCPSocketAction{
				Port: intstr.IntOrString{IntVal: int32(appPort)},
			},
		},
		"/app-health/hello-world/livez": &Prober{
			TCPSocket: &apimirror.TCPSocketAction{
				Port: intstr.IntOrString{IntVal: int32(appPort)},
			},
		},
	}

	type test struct {
		name       string
		probePath  string
		config     KubeAppProbers
		podIP      string
		ipv6       bool
		statusCode int
	}
	testCases := []test{
		{
			name:       "http-bad-path",
			probePath:  "bad-path-should-be-404",
			config:     simpleHTTPConfig,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "http-readyz",
			probePath:  "app-health/hello-world/readyz",
			config:     simpleHTTPConfig,
			statusCode: http.StatusOK,
		},
		{
			name:       "http-livez",
			probePath:  "app-health/hello-world/livez",
			config:     simpleHTTPConfig,
			statusCode: http.StatusOK,
		},
		{
			name:       "http-livez-localhost",
			probePath:  "app-health/hello-world/livez",
			config:     simpleHTTPConfig,
			statusCode: http.StatusOK,
			podIP:      "localhost",
		},
		{
			name:      "http-readyz-header",
			probePath: "app-health/header/readyz",
			config: KubeAppProbers{
				"/app-health/header/readyz": &Prober{
					HTTPGet: &apimirror.HTTPGetAction{
						Port: intstr.IntOrString{IntVal: int32(appPort)},
						Path: "/header",
						HTTPHeaders: []apimirror.HTTPHeader{
							{"Host", testHostValue},
							{testHeader, testHeaderValue},
						},
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name:      "http-readyz-path",
			probePath: "app-health/hello-world/readyz",
			config: KubeAppProbers{
				"/app-health/hello-world/readyz": &Prober{
					HTTPGet: &apimirror.HTTPGetAction{
						Path: "hello/texas",
						Port: intstr.IntOrString{IntVal: int32(appPort)},
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name:      "http-livez-path",
			probePath: "app-health/hello-world/livez",
			config: KubeAppProbers{
				"/app-health/hello-world/livez": &Prober{
					HTTPGet: &apimirror.HTTPGetAction{
						Path: "hello/texas",
						Port: intstr.IntOrString{IntVal: int32(appPort)},
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name:       "tcp-readyz",
			probePath:  "app-health/hello-world/readyz",
			config:     simpleTCPConfig,
			statusCode: http.StatusOK,
		},
		{
			name:       "tcp-livez",
			probePath:  "app-health/hello-world/livez",
			config:     simpleTCPConfig,
			statusCode: http.StatusOK,
		},
		{
			name:       "tcp-livez-ipv4",
			probePath:  "app-health/hello-world/livez",
			config:     simpleTCPConfig,
			statusCode: http.StatusOK,
			podIP:      "127.0.0.1",
		},
		{
			name:       "tcp-livez-ipv6",
			probePath:  "app-health/hello-world/livez",
			config:     simpleTCPConfig,
			statusCode: http.StatusOK,
			podIP:      "::1",
			ipv6:       true,
		},
		{
			name:       "tcp-livez-wrapped-ipv6",
			probePath:  "app-health/hello-world/livez",
			config:     simpleTCPConfig,
			statusCode: http.StatusOK,
			podIP:      "[::1]",
			ipv6:       true,
		},
		{
			name:       "tcp-livez-localhost",
			probePath:  "app-health/hello-world/livez",
			config:     simpleTCPConfig,
			statusCode: http.StatusOK,
			podIP:      "localhost",
		},
		{
			name:      "redirect",
			probePath: "app-health/redirect/livez",
			config: KubeAppProbers{
				"/app-health/redirect/livez": &Prober{
					HTTPGet: &apimirror.HTTPGetAction{
						Path: "redirect",
						Port: intstr.IntOrString{IntVal: int32(appPort)},
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name:      "redirect loop",
			probePath: "app-health/redirect-loop/livez",
			config: KubeAppProbers{
				"/app-health/redirect-loop/livez": &Prober{
					HTTPGet: &apimirror.HTTPGetAction{
						Path: "redirect-loop",
						Port: intstr.IntOrString{IntVal: int32(appPort)},
					},
				},
			},
			statusCode: http.StatusInternalServerError,
		},
		{
			name:      "remote redirect",
			probePath: "app-health/remote-redirect/livez",
			config: KubeAppProbers{
				"/app-health/remote-redirect/livez": &Prober{
					HTTPGet: &apimirror.HTTPGetAction{
						Path: "remote-redirect",
						Port: intstr.IntOrString{IntVal: int32(appPort)},
					},
				},
			},
			statusCode: http.StatusOK,
		},
	}
	testFn := func(t *testing.T, tc test) {
		appProber, err := json.Marshal(tc.config)
		if err != nil {
			t.Fatalf("invalid app probers")
		}
		config := Options{
			StatusPort:     0,
			KubeAppProbers: string(appProber),
			PodIP:          tc.podIP,
			IPv6:           tc.ipv6,
		}
		// Starts the pilot agent status server.
		server, err := NewServer(config)
		if err != nil {
			t.Fatalf("failed to create status server %v", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go server.Run(ctx)

		if tc.ipv6 {
			server.upstreamLocalAddress = &net.TCPAddr{IP: net.ParseIP("::1")} // required because ::6 is NOT a loopback address (IPv6 only has ::1)
		}

		var statusPort uint16
		for statusPort == 0 {
			server.mutex.RLock()
			statusPort = server.statusPort
			server.mutex.RUnlock()
		}

		client := http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/%s", statusPort, tc.probePath), nil)
		if err != nil {
			t.Fatalf("[%v] failed to create request", tc.probePath)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal("request failed: ", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != tc.statusCode {
			t.Errorf("[%v] unexpected status code, want = %v, got = %v", tc.probePath, tc.statusCode, resp.StatusCode)
		}
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) { testFn(t, tc) })
	}
	// Next we check ever
	t.Run("status codes", func(t *testing.T) {
		for code := http.StatusOK; code <= http.StatusNetworkAuthenticationRequired; code++ {
			if http.StatusText(code) == "" { // Not a valid HTTP code
				continue
			}
			expect := code
			if isRedirect(code) {
				expect = 200
			}
			t.Run(fmt.Sprint(code), func(t *testing.T) {
				testFn(t, test{
					probePath: "app-health/code/livez",
					config: KubeAppProbers{
						"/app-health/code/livez": &Prober{
							TimeoutSeconds: 1,
							HTTPGet: &apimirror.HTTPGetAction{
								Path: fmt.Sprintf("status/%d", code),
								Port: intstr.IntOrString{IntVal: int32(appPort)},
							},
						},
					},
					statusCode: expect,
				})
			})
		}
	})
}

func TestHttpsAppProbe(t *testing.T) {
	// Starts the application first.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Errorf("failed to allocate unused port %v", err)
	}
	keyFile := env.IstioSrc + "/pilot/cmd/pilot-agent/status/test-cert/cert.key"
	crtFile := env.IstioSrc + "/pilot/cmd/pilot-agent/status/test-cert/cert.crt"
	go http.ServeTLS(listener, &handler{}, crtFile, keyFile)
	appPort := listener.Addr().(*net.TCPAddr).Port

	// Starts the pilot agent status server.
	server, err := NewServer(Options{
		StatusPort: 0,
		KubeAppProbers: fmt.Sprintf(`{"/app-health/hello-world/readyz": {"httpGet": {"path": "/hello/sunnyvale", "port": %v, "scheme": "HTTPS"}},
"/app-health/hello-world/livez": {"httpGet": {"port": %v, "scheme": "HTTPS"}}}`, appPort, appPort),
	})
	if err != nil {
		t.Errorf("failed to create status server %v", err)
		return
	}
	go server.Run(context.Background())

	var statusPort uint16
	if err := retry.UntilSuccess(func() error {
		server.mutex.RLock()
		statusPort = server.statusPort
		server.mutex.RUnlock()
		if statusPort == 0 {
			return fmt.Errorf("no port allocated")
		}
		return nil
	}); err != nil {
		t.Fatalf("failed to getport: %v", err)
	}
	t.Logf("status server starts at port %v, app starts at port %v", statusPort, appPort)
	testCases := []struct {
		probePath  string
		statusCode int
	}{
		{
			probePath:  fmt.Sprintf(":%v/bad-path-should-be-disallowed", statusPort),
			statusCode: http.StatusNotFound,
		},
		{
			probePath:  fmt.Sprintf(":%v/app-health/hello-world/readyz", statusPort),
			statusCode: http.StatusOK,
		},
		{
			probePath:  fmt.Sprintf(":%v/app-health/hello-world/livez", statusPort),
			statusCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		client := http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost%s", tc.probePath), nil)
		if err != nil {
			t.Errorf("[%v] failed to create request", tc.probePath)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal("request failed")
		}
		defer resp.Body.Close()
		if resp.StatusCode != tc.statusCode {
			t.Errorf("[%v] unexpected status code, want = %v, got = %v", tc.probePath, tc.statusCode, resp.StatusCode)
		}
	}
}

func TestProbeHeader(t *testing.T) {
	headerChecker := func(t *testing.T, header http.Header) net.Listener {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			t.Fatalf("failed to allocate unused port %v", err)
		}
		go http.Serve(listener, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			r.Header.Del("User-Agent")
			r.Header.Del("Accept-Encoding")
			if !reflect.DeepEqual(r.Header, header) {
				t.Errorf("unexpected header, want = %v, got = %v", header, r.Header)
				http.Error(rw, "", http.StatusBadRequest)
				return
			}
			http.Error(rw, "", http.StatusOK)
		}))
		return listener
	}

	testCases := []struct {
		name          string
		originHeaders http.Header
		proxyHeaders  []apimirror.HTTPHeader
		want          http.Header
	}{
		{
			name: "Only Origin",
			originHeaders: http.Header{
				testHeader: []string{testHeaderValue},
			},
			proxyHeaders: []apimirror.HTTPHeader{},
			want: http.Header{
				testHeader: []string{testHeaderValue},
			},
		},
		{
			name: "Only Origin, has multiple values",
			originHeaders: http.Header{
				testHeader: []string{testHeaderValue, testHeaderValue},
			},
			proxyHeaders: []apimirror.HTTPHeader{},
			want: http.Header{
				testHeader: []string{testHeaderValue, testHeaderValue},
			},
		},
		{
			name:          "Only Proxy",
			originHeaders: http.Header{},
			proxyHeaders: []apimirror.HTTPHeader{
				{
					Name:  testHeader,
					Value: testHeaderValue,
				},
			},
			want: http.Header{
				testHeader: []string{testHeaderValue},
			},
		},
		{
			name:          "Only Proxy, has multiple values",
			originHeaders: http.Header{},
			proxyHeaders: []apimirror.HTTPHeader{
				{
					Name:  testHeader,
					Value: testHeaderValue,
				},
				{
					Name:  testHeader,
					Value: testHeaderValue,
				},
			},
			want: http.Header{
				testHeader: []string{testHeaderValue, testHeaderValue},
			},
		},
		{
			name: "Proxy overwrites Origin",
			originHeaders: http.Header{
				testHeader: []string{testHeaderValue, testHeaderValue},
			},
			proxyHeaders: []apimirror.HTTPHeader{
				{
					Name:  testHeader,
					Value: testHeaderValue + "Over",
				},
			},
			want: http.Header{
				testHeader: []string{testHeaderValue + "Over"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := headerChecker(t, tc.want)
			defer svc.Close()
			probePath := "/app-health/hello-world/livez"
			appAddress := svc.Addr().(*net.TCPAddr)
			appProber, err := json.Marshal(KubeAppProbers{
				probePath: &Prober{
					HTTPGet: &apimirror.HTTPGetAction{
						Port:        intstr.IntOrString{IntVal: int32(appAddress.Port)},
						Host:        appAddress.IP.String(),
						Path:        "/header",
						HTTPHeaders: tc.proxyHeaders,
					},
				},
			})
			if err != nil {
				t.Fatalf("invalid app probers")
			}
			config := Options{
				StatusPort:     0,
				KubeAppProbers: string(appProber),
			}
			// Starts the pilot agent status server.
			server, err := NewServer(config)
			if err != nil {
				t.Fatal("failed to create status server: ", err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go server.Run(ctx)

			var statusPort uint16
			for statusPort == 0 {
				server.mutex.RLock()
				statusPort = server.statusPort
				server.mutex.RUnlock()
			}

			client := http.Client{}
			req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v%s", statusPort, probePath), nil)
			if err != nil {
				t.Fatal("failed to create request: ", err)
			}
			req.Header = tc.originHeaders
			resp, err := client.Do(req)
			if err != nil {
				t.Fatal("request failed: ", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("unexpected status code, want = %v, got = %v", http.StatusOK, resp.StatusCode)
			}
		})
	}
}

func TestHandleQuit(t *testing.T) {
	statusPort := 15020
	s, err := NewServer(Options{StatusPort: uint16(statusPort)})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		method     string
		remoteAddr string
		expected   int
	}{
		{
			name:       "should send a sigterm for valid requests",
			method:     "POST",
			remoteAddr: "127.0.0.1",
			expected:   http.StatusOK,
		},
		{
			name:       "should send a sigterm for valid ipv6 requests",
			method:     "POST",
			remoteAddr: "[::1]",
			expected:   http.StatusOK,
		},
		{
			name:       "should require POST method",
			method:     "GET",
			remoteAddr: "127.0.0.1",
			expected:   http.StatusMethodNotAllowed,
		},
		{
			name:     "should require localhost",
			method:   "POST",
			expected: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Need to stop SIGTERM from killing the whole test run
			termChannel := make(chan os.Signal, 1)
			signal.Notify(termChannel, syscall.SIGTERM)
			defer signal.Reset(syscall.SIGTERM)

			req, err := http.NewRequest(tt.method, "/quitquitquit", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.remoteAddr != "" {
				req.RemoteAddr = tt.remoteAddr + ":15020"
			}

			resp := httptest.NewRecorder()
			s.handleQuit(resp, req)
			if resp.Code != tt.expected {
				t.Fatalf("Expected response code %v got %v", tt.expected, resp.Code)
			}

			if tt.expected == http.StatusOK {
				select {
				case <-termChannel:
				case <-time.After(time.Second):
					t.Fatalf("Failed to receive expected SIGTERM")
				}
			} else if len(termChannel) != 0 {
				t.Fatalf("A SIGTERM was sent when it should not have been")
			}
		})
	}
}

func TestAdditionalProbes(t *testing.T) {
	rp := readyProbe{}
	urp := unreadyProbe{}
	testCases := []struct {
		name   string
		probes []ready.Prober
		err    error
	}{
		{
			name:   "success probe",
			probes: []ready.Prober{rp},
			err:    nil,
		},
		{
			name:   "not ready probe",
			probes: []ready.Prober{urp},
			err:    errors.New("not ready"),
		},
		{
			name:   "both probes",
			probes: []ready.Prober{rp, urp},
			err:    errors.New("not ready"),
		},
	}
	testServer := testserver.CreateAndStartServer(liveServerStats)
	defer testServer.Close()
	for _, tc := range testCases {
		server, err := NewServer(Options{
			Probes:    tc.probes,
			AdminPort: uint16(testServer.Listener.Addr().(*net.TCPAddr).Port),
		})
		if err != nil {
			t.Errorf("failed to construct server")
		}
		err = server.isReady()
		if tc.err == nil {
			if err != nil {
				t.Errorf("Unexpected result, expected: %v got: %v", tc.err, err)
			}
		} else {
			if err.Error() != tc.err.Error() {
				t.Errorf("Unexpected result, expected: %v got: %v", tc.err, err)
			}
		}

	}
}

type readyProbe struct{}

func (s readyProbe) Check() error {
	return nil
}

type unreadyProbe struct{}

func (u unreadyProbe) Check() error {
	return errors.New("not ready")
}

func TestMetricReaderWriteTo(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantW   string
	}{
		{
			name: "should write everything",
			content: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0`,
			wantW: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0`,
		},
		{
			name: "should remove tail /n/n",
			content: `# TYPE my_metric counter
my_metric{} 0


`,
			wantW: `# TYPE my_metric counter
my_metric{} 0
`,
		},
		{
			name: "should remove head /n",
			content: `

# TYPE my_metric counter
my_metric{} 0
`,
			wantW: `# TYPE my_metric counter
my_metric{} 0
`,
		},
		{
			name: "should remove all /n/n",
			content: `


# TYPE my_metric counter


my_metric{} 0


# TYPE my_other_metric counter


my_other_metric{} 0


`,
			wantW: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &metricReader{
				buf: bufio.NewReader(bytes.NewReader([]byte(tt.content))),
			}
			w := &bytes.Buffer{}
			_, _ = r.WriteTo(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteTo() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMetricReaderWriteToError(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		errPostion int
		wantW      string
	}{
		{
			name: "remove last line",
			content: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0`,
			errPostion: 79,
			wantW: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
`,
		},
		{
			name: "remove nothing",
			content: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0`,
			errPostion: 89,
			wantW: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0`,
		},
		{
			name: "remove all",
			content: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0`,
			errPostion: 3,
			wantW:      ``,
		},
		{
			name: "remove half",
			content: `# TYPE my_metric counter
my_metric{} 0
# TYPE my_other_metric counter
my_other_metric{} 0`,
			errPostion: 50,
			wantW: `# TYPE my_metric counter
my_metric{} 0
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", strconv.Itoa(len(tt.content)))
				if _, err := w.Write([]byte(tt.content[:tt.errPostion])); err != nil {
					t.Fatalf("write failed: %v", err)
				}
			}))
			defer server.Close()

			resp, _ := server.Client().Get(server.URL)
			defer resp.Body.Close()
			r := &metricReader{
				buf: bufio.NewReader(resp.Body),
			}
			w := &bytes.Buffer{}
			_, _ = r.WriteTo(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("WriteTo() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestMetricReaderWriteToLongLine(t *testing.T) {
	tests := []struct {
		name        string
		blockLength int
		blockNums   int
		insertBytes []byte
		wantLength  int
	}{
		{name: "nothing left",
			blockLength: 0,
			blockNums:   4 * 4096,
			insertBytes: []byte("\n"),
			wantLength:  0,
		},
		{name: "normal",
			blockLength: 1022,
			blockNums:   4 * 2,
			insertBytes: []byte("\n\n"),
			wantLength:  2 * 4 * (1024 - 1),
		},
		// bufio reader default buffer size is 4096
		{name: "line length less than buffer size",
			blockLength: 4094,
			blockNums:   4,
			insertBytes: []byte("\n\n"),
			wantLength:  4 * 4095,
		},
		{name: "line length equal buffer size",
			blockLength: 4096,
			blockNums:   4,
			insertBytes: []byte("\n\n"),
			wantLength:  4 * 4097,
		},
		{name: "line length greater than buffer size",
			blockLength: 4097,
			blockNums:   4,
			insertBytes: []byte("\n\n"),
			wantLength:  4 * 4098,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipeR, pipeW := io.Pipe()
			go func() {
				block := make([]byte, tt.blockLength)
				for i := 0; i < tt.blockNums; i++ {
					_, _ = pipeW.Write(block)
					_, _ = pipeW.Write(tt.insertBytes)
				}
				_ = pipeW.Close()
			}()
			r := &metricReader{
				buf: bufio.NewReader(pipeR),
			}
			w := &bytes.Buffer{}
			_, _ = r.WriteTo(w)
			if w.Len() != tt.wantLength {
				t.Errorf("WriteTo() gotL = %v, want %v", w.Len(), tt.wantLength)
			}
		})
	}
}


func Test_metricReader_Read(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantW   string
	}{
		{
			name: "should write everything",
			content: `test
test`,
			wantW: `test
test`,
		},
		{
			name: "should get compact content",
			content: `test

test
test

test

`,
			wantW: `test
test
test
test
`,
		},
		{
			name: "should remove tail /n/n",
			content: `test

`,
			wantW: `test
`,
		},
		{
			name: "should remove head /n",
			content: `

test`,
			wantW: `test`,
		},
		{
			name: "should remove all /n/n",
			content: `test



test

`,
			wantW: `test
test
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &metricReader{
				buf: bufio.NewReader(bytes.NewReader([]byte(tt.content))),
			}
			b := make([]byte, len(tt.wantW))
			l, _ := r.Read(b)
			if gotW := string(b[:l]); gotW != tt.wantW {
				t.Errorf("Read() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Benchmark_metricReader_WriteTo(b *testing.B) {
	block := make([]byte, 1024)
	doubleNewline := []byte("\n\n")
	data := &bytes.Buffer{}
	for i := 0; i < 100*1024; i++ {
		_, _ = data.Write(block)
		_, _ = data.Write(doubleNewline)
	}
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data.Bytes())
		r := &metricReader{
			buf: bufio.NewReader(buf),
		}
		w := &bytes.Buffer{}
		_, _ = r.WriteTo(w)
	}
}

func Benchmark_not_metricReader(b *testing.B) {
	block := make([]byte, 1024)
	doubleNewline := []byte("\n\n")
	data := &bytes.Buffer{}
	for i := 0; i < 100*1024; i++ {
		_, _ = data.Write(block)
		_, _ = data.Write(doubleNewline)
	}
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(data.Bytes())
		raw, _ := io.ReadAll(buf)
		bytes.ReplaceAll(raw, []byte("\n\n"), []byte("\n"))
	}
}

func Benchmark_metricReader_Read(b *testing.B) {
	block := make([]byte, 1024)
	doubleNewline := []byte("\n\n")
	data := &bytes.Buffer{}
	for i := 0; i < 100*1024; i++ {
		_, _ = data.Write(block)
		_, _ = data.Write(doubleNewline)
	}
	for i := 0; i < b.N; i++ {
		bs := data.Bytes()
		r := &metricReader{
			buf: bufio.NewReader(bytes.NewReader(bs)),
		}
		res := make([]byte, len(bs))
		_, _ = r.Read(res)
	}
}