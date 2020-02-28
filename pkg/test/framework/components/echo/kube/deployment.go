// Copyright 2019 Istio Authors
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

package kube

import (
	"fmt"
	"text/template"

	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/core/image"
	"istio.io/istio/pkg/test/util/tmpl"
)

const (
	deploymentYAML = `
{{- if .ServiceAccount }}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Service }}
---
{{- end }}
apiVersion: v1
kind: Service
metadata:
  name: {{ .Service }}
  labels:
    app: {{ .Service }}
{{- if .ServiceAnnotations }}
  annotations:
{{- range $name, $value := .ServiceAnnotations }}
    {{ $name }}: {{ printf "%q" $value }}
{{- end }}
{{- end }}
spec:
{{- if .Headless }}
  clusterIP: None
{{- end }}
  ports:
{{- range $i, $p := .Ports }}
  - name: {{ $p.Name }}
    port: {{ $p.ServicePort }}
    targetPort: {{ $p.InstancePort }}
{{- end }}
  selector:
    app: {{ .Service }}
---
{{$deployConfig := .Workloads }}
{{- range $i, $w := $deployConfig }}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ $.Service }}-{{ $w.Version }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ $.Service }}
      version: {{ $w.Version }}
{{- if ne $.Locality "" }}
      istio-locality: {{ $.Locality }}
{{- end }}
  template:
    metadata:
      labels:
        app: {{ $.Service }}
        version: {{ $w.Version }}
{{- if ne $.Locality "" }}
        istio-locality: {{ $.Locality }}
{{- end }}
      annotations:
        foo: bar
{{- range $name, $value := index $.WorkloadAnnotations $w.Version }}
        {{ $name }}: {{ printf "%q" $value }}
{{- end }}
{{- if $.IncludeInboundPorts }}
        traffic.sidecar.istio.io/includeInboundPorts: "{{ $.IncludeInboundPorts }}"
{{- end }}
    spec:
{{- if $.ServiceAccount }}
      serviceAccountName: {{ $.Service }}
{{- end }}
      containers:
      - name: app
        image: {{ $.Hub }}/app:{{ $.Tag }}
        imagePullPolicy: {{ $.PullPolicy }}
        securityContext:
          runAsUser: 1
        args:
{{- range $i, $p := $.ContainerPorts }}
{{- if eq .Protocol "GRPC" }}
          - --grpc
{{- else }}
          - --port
{{- end }}
          - "{{ $p.Port }}"
{{- end }}
{{- range $i, $p := $.WorkloadOnlyPorts }}
          - --port
          - "{{ $p }}"
{{- end }}
          - --version
          - "{{ $w.Version }}"
        ports:
{{- range $i, $p := $.ContainerPorts }}
        - containerPort: {{ $p.Port }} 
{{- if eq .Port 3333 }}
          name: tcp-health-port
{{- end }}
{{- end }}
        readinessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          failureThreshold: 10
        livenessProbe:
          tcpSocket:
            port: tcp-health-port
          initialDelaySeconds: 10
          periodSeconds: 10
          failureThreshold: 10
---
{{- end}}
apiVersion: v1
kind: Secret
metadata:
  name: sdstokensecret
type: Opaque
stringData:
  sdstoken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2\
VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Ii\
wia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InZhdWx0LWNpdGFkZWwtc2\
EtdG9rZW4tNzR0d3MiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC\
5uYW1lIjoidmF1bHQtY2l0YWRlbC1zYSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2Vydm\
ljZS1hY2NvdW50LnVpZCI6IjJhYzAzYmEyLTY5MTUtMTFlOS05NjkwLTQyMDEwYThhMDExNCIsInN1Yi\
I6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnZhdWx0LWNpdGFkZWwtc2EifQ.pZ8SiyNeO0p\
1p8HB9oXvXOAI1XCJZKk2wVHXBsTSzKWxlVD9HrHbAcSbO2dlhFpeCgknt6eZywvhShZJh2F6-iHP_Yo\
UVoCqQmzjPoB3c3JoYFpJo-9jTN1_mNRtZUcNvYl-tDlTmBlaKEvoC5P2WGVUF3AoLsES66u4FG9Wllm\
LV92LG1WNqx_ltkT1tahSy9WiHQgyzPqwtwE72T1jAGdgVIoJy1lfSaLam_bo9rqkRlgSg-au9BAjZiD\
Gtm9tf3lwrcgfbxccdlG4jAsTFa2aNs3dW4NLk7mFnWCJa-iWj-TgFxf9TW-9XPK0g3oYIQ0Id0CIW2S\
iFxKGPAjB-g"
`
)

var (
	deploymentTemplate *template.Template
)

func init() {
	deploymentTemplate = template.New("echo_deployment")
	if _, err := deploymentTemplate.Parse(deploymentYAML); err != nil {
		panic(fmt.Sprintf("unable to parse echo deployment template: %v", err))
	}
}

func generateYAML(cfg echo.Config) (string, error) {
	// Create the parameters for the YAML template.
	settings, err := image.SettingsFromCommandLine()
	if err != nil {
		return "", err
	}
	return generateYAMLWithSettings(cfg, settings)
}

func generateYAMLWithSettings(cfg echo.Config, settings *image.Settings) (string, error) {
	// Convert legacy config to workload oritended.
	if cfg.Workloads == nil {
		cfg.Workloads = []echo.WorkloadConfig{
			{
				Version: cfg.Version,
			},
		}
	}

	// Separate the annotations.
	serviceAnnotations := make(map[string]string)
	wlas := make(map[string]map[string]string)
	for i := range cfg.Workloads {
		w := &cfg.Workloads[i]
		if cfg.Workloads[i].Version == "" {
			cfg.Workloads[i].Version = "v1"
		}
		wlas[w.Version] = map[string]string{}
		if cfg.Workloads[i].Version == "" {
			cfg.Workloads[i].Version = cfg.Version
		}
		for k, val := range w.Annotations {
			wlas[w.Version][k.Name] = val.Value
		}
	}
	for k, v := range cfg.ServiceAnnotations {
		serviceAnnotations[k.Name] = v.Value
	}

	params := map[string]interface{}{
		"Hub":                 settings.Hub,
		"Tag":                 settings.Tag,
		"PullPolicy":          settings.PullPolicy,
		"Service":             cfg.Service,
		"Version":             cfg.Version,
		"Headless":            cfg.Headless,
		"Locality":            cfg.Locality,
		"ServiceAccount":      cfg.ServiceAccount,
		"Ports":               cfg.Ports,
		"WorkloadOnlyPorts":   cfg.WorkloadOnlyPorts,
		"ContainerPorts":      getContainerPorts(cfg.Ports),
		"ServiceAnnotations":  serviceAnnotations,
		"WorkloadAnnotations": wlas,
		"IncludeInboundPorts": cfg.IncludeInboundPorts,
		"Workloads":           cfg.Workloads,
	}

	// Generate the YAML content.
	return tmpl.Execute(deploymentTemplate, params)
}
