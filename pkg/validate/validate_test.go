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

package validate

import (
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"

	"istio.io/operator/pkg/apis/istio/v1alpha2"
	"istio.io/operator/pkg/util"
)

func TestUnmarshalKubernetes(t *testing.T) {
	tests := []struct {
		desc    string
		yamlStr string
		want    string
	}{
		{
			desc:    "nil success",
			yamlStr: "",
			want:    "{}",
		},
		{
			desc: "hpaSpec",
			yamlStr: `
hpaSpec:
  maxReplicas: 10
  minReplicas: 1
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: istio-pilot
  metrics:
    - type: Resource
      resource:
        name: cpu
        targetAverageUtilization: 80
`,
		},
		{
			desc: "resources",
			yamlStr: `
resources:
  limits:
    cpu: 444m
    memory: 333Mi
  requests:
    cpu: 222m
    memory: 111Mi
`,
		},
		{
			desc: "podDisruptionBudget",
			yamlStr: `
podDisruptionBudget:
  maxUnavailable: 1
  selector:
    matchLabels:
      app: pilot
`,
		},
		{
			desc: "readinessProbeWithPortNumber",
			yamlStr: `
readinessProbe:
  failureThreshold: 44
  initialDelaySeconds: 11
  periodSeconds: 22
  successThreshold: 33
  httpGet:
    path: /ready
    port: 8080
`,
		},
		{
			desc: "readinessProbeWithPortName",
			yamlStr: `
readinessProbe:
  failureThreshold: 44
  initialDelaySeconds: 11
  periodSeconds: 22
  successThreshold: 33
  httpGet:
    path: /ready
    port: http
`,
		},
		{
			desc: "affinity",
			yamlStr: `
affinity:
  podAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - labelSelector:
        matchExpressions:
        - key: security
          operator: In
          values:
          - S1
      topologyKey: failure-domain.beta.kubernetes.io/zone
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: security
            operator: In
            values:
            - S2
        topologyKey: failure-domain.beta.kubernetes.io/zone
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tk := &v1alpha2.TestKube{}
			err := util.UnmarshalWithJSONPB(tt.yamlStr, tk)
			if err != nil {
				t.Fatalf("unmarshalWithJSONPB(%s): got error %s", tt.desc, err)
			}
			s, err := util.MarshalWithJSONPB(tk)
			if err != nil {
				t.Fatalf("unmarshalWithJSONPB(%s): got error %s", tt.desc, err)
			}
			got, want := stripNL(s), stripNL(tt.want)
			if want == "" {
				want = stripNL(tt.yamlStr)
			}
			if !util.IsYAMLEqual(got, want) {
				t.Errorf("%s: got:\n%s\nwant:\n%s\n(-got, +want)\n%s\n", tt.desc, got, want, diff.Diff(got, want))
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		desc     string
		yamlStr  string
		wantErrs util.Errors
	}{
		{
			desc: "nil success",
		},
		{
			desc: "TrafficManagement",
			yamlStr: `
trafficManagement:
  enabled: true
  components:
    namespace: istio-system-traffic
`,
		},
		{
			desc: "SidecarInjectorConfig",
			yamlStr: `
autoInjection:
  components:
    namespace: istio-control
    injector:
      enabled: true
`,
		},
		{
			desc: "CommonConfig",
			// TODO:        debug: INFO
			yamlStr: `
hub: docker.io/istio
tag: v1.2.3
trafficManagement:
  components:
    proxy:
      enabled: true
      namespace: istio-control-system
      k8s:
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 11
          periodSeconds: 22
          successThreshold: 33
          failureThreshold: 44
        hpaSpec:
          scaleTargetRef:
            apiVersion: apps/v1
            kind: Deployment
            name: php-apache
          minReplicas: 1
          maxReplicas: 10
          metrics:
            - type: Resource
              resource:
                name: cpu
                targetAverageUtilization: 80
        nodeSelector:
          disktype: ssd
values:
  global:
    proxy:
      includeIpRanges: "1.1.0.0/16,2.2.0.0/16"
      excludeIpRanges: "3.3.0.0/16,4.4.0.0/16"

`,
		},
		{
			desc: "BadTag",
			yamlStr: `
hub: ?illegal-tag!
`,
			wantErrs: makeErrors([]string{`invalid value Hub: ?illegal-tag!`}),
		},
		{
			desc: "BadHub",
			yamlStr: `
hub: docker.io:tag/istio
`,
			wantErrs: makeErrors([]string{`invalid value Hub: docker.io:tag/istio`}),
		},
		{
			desc: "GoodURL",
			yamlStr: `
installPackagePath: /local/file/path
`,
		},
		{
			desc: "BadValuesIP",
			yamlStr: `
values:
  global:
    proxy:
      includeIpRanges: "1.1.0.300/16,2.2.0.0/16"
`,
			wantErrs: makeErrors([]string{`global.proxy.includeIpRanges invalid CIDR address: 1.1.0.300/16`}),
		},
		{
			desc: "EmptyValuesIP",
			yamlStr: `
values:
  global:
    proxy:
      includeIpRanges: ""
      includeInboundPorts: "*"
`,
		},
	}

	for _, tt := range tests {
		if tt.desc != "EmptyValuesIP" {
			continue
		}
		t.Run(tt.desc, func(t *testing.T) {
			ispec := &v1alpha2.IstioControlPlaneSpec{}
			err := util.UnmarshalWithJSONPB(tt.yamlStr, ispec)
			if err != nil {
				t.Fatalf("unmarshalWithJSONPB(%s): got error %s", tt.desc, err)
			}
			errs := CheckIstioControlPlaneSpec(ispec, false)
			if gotErrs, wantErrs := errs, tt.wantErrs; !util.EqualErrors(gotErrs, wantErrs) {
				t.Errorf("ProtoToValues(%s)(%v): gotErrs:%s, wantErrs:%s", tt.desc, tt.yamlStr, gotErrs, wantErrs)
			}
		})
	}
}

func stripNL(s string) string {
	return strings.Trim(s, "\n")
}
