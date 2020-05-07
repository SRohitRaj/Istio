// Copyright 2020 Istio Authors
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

package bootstrap

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"istio.io/pkg/log"

	"istio.io/istio/pilot/pkg/features"
)

const (
	HTTPSHandlerReadyPath = "/httpsReady"
)

func (s *Server) initHTTPSWebhookServer(args *PilotArgs) {
	if features.IstiodService.Get() == "" {
		log.Info("Not starting HTTPS webhook server: istiod address not set")
		return
	}

	log.Info("Setting up HTTPS webhook server for istiod webhooks")

	// create the https server for hosting the k8s injectionWebhook handlers.
	s.httpsMux = http.NewServeMux()
	s.httpsServer = &http.Server{
		Addr:    args.DiscoveryOptions.HTTPSAddr,
		Handler: s.httpsMux,
		TLSConfig: &tls.Config{
			GetCertificate: s.getIstiodCertificate,
		},
	}

	// setup our readiness handler and the corresponding client we'll use later to check it with.
	s.httpsMux.HandleFunc(HTTPSHandlerReadyPath, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s.httpsReadyClient = &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (s *Server) checkHTTPSWebhookServerReadiness() int {
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "https",
			Host:   s.httpsServer.Addr,
			Path:   HTTPSHandlerReadyPath,
		},
	}

	response, err := s.httpsReadyClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer response.Body.Close()

	return response.StatusCode
}
