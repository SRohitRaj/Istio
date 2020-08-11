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

package mesh

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestValidateSetFlags(t *testing.T) {

	tests := []struct {
		name string
		args []string
		want error
	}{
		{
			name: "Test when no flag prams in sent",
			args: []string{},
			want: nil,
		},
		{
			name: "Test invalid flag format",
			args: []string{
				"values.global.sds.enabled",
			},
			want: fmt.Errorf("\n Invalid flag format %q", "values.global.sds.enabled"),
		},
		{
			name: "Test valid flag format",
			args: []string{
				"values.global.sds.enabled=true",
			},
			want: nil,
		},
		{
			name: "Test flag name not available",
			args: []string{
				"values.global.controlPlaneSecurity=true",
			},
			want: fmt.Errorf("\n Invalid flag: %q", "values.global.controlPlaneSecurity"),
		},
		{
			name: "Test flag name available",
			args: []string{
				"values.global.controlPlaneSecurityEnabled=true",
			},
			want: nil,
		},
		{
			name: "Test Unsupported values",
			args: []string{
				"values.global.imagePullPolicy=Occasionally",
			},
			want: fmt.Errorf("\n Unsupported value: %q, supported values for: %q is %q",
				"Occasionally", "values.global.imagePullPolicy", strings.Join(imagePullPolicy, ", ")),
		},
		{
			name: "Test supported values",
			args: []string{
				"values.global.imagePullPolicy=IfNotPresent",
			},
			want: nil,
		},
		{
			name: "Test supported traceSampling",
			args: []string{
				"values.pilot.traceSampling=10.5",
			},
			want: nil,
		},
		{
			name: "Test Unsupported traceSampling",
			args: []string{
				"values.pilot.traceSampling=100.5",
			},
			want: fmt.Errorf("\n Unsupported value: %q, supported values for: %q is between %.1f to %.1f",
				"100.5", "values.pilot.traceSampling", traceSamplingMin, traceSamplingMax),
		},
		{
			name: "Test valid namespace",
			args: []string{
				"namespace=istio-system",
			},
			want: nil,
		},
		{
			name: "Test invalid namespace",
			args: []string{
				"namespace=foo.bar",
			},
			want: fmt.Errorf("\n Unsupported format: %q for flag %q", "foo.bar", "namespace"),
		},
		{
			name: "Test valid revision",
			args: []string{
				"revision=canary",
			},
			want: nil,
		},
		{
			name: "Test invalid revision",
			args: []string{
				"revision=v1.2.3",
			},
			want: fmt.Errorf("\n Unsupported format: %q for flag %q", "v1.2.3", "revision"),
		},
		{
			name: "Test valid reportBatchMaxEntries",
			args: []string{
				"values.mixer.telemetry.reportBatchMaxEntries=100",
			},
			want: nil,
		},
		{
			name: "Test invalid reportBatchMaxEntries",
			args: []string{
				"values.mixer.telemetry.reportBatchMaxEntries=ten",
			},
			want: fmt.Errorf("\n Unsupported value: %q for flag %q, use valid value eg: 100",
				"ten", "values.mixer.telemetry.reportBatchMaxEntries"),
		},
		{
			name: "Test valid outboundTrafficPolicy mode",
			args: []string{
				"values.meshConfig.outboundTrafficPolicy.mode=ALLOW_ANY",
			},
			want: nil,
		},
		{
			name: "Test valid outboundClusterStatName",
			args: []string{
				"values.meshConfig.outboundClusterStatName=%SERVICE_FQDN%_%SERVICE_PORT%",
			},
			want: nil,
		},
		{
			name: "Test invalid outboundClusterStatName",
			args: []string{
				"values.meshConfig.outboundClusterStatName=%SERVICE_FQDN%_%SERVICE_POT%",
			},
			want: fmt.Errorf("\n Unsupported value: %q, supported values for: %q is %q",
				"%SERVICE_FQDN%_%SERVICE_POT%", "values.meshConfig.outboundClusterStatName",
				strings.Join(outboundClusterStatName, ", ")),
		},
		{
			name: "Test valid h2UpgradePolicy",
			args: []string{
				"values.meshConfig.h2UpgradePolicy=UPGRADE",
			},
			want: nil,
		},
		{
			name: "Test invalid h2UpgradePolicy",
			args: []string{
				"values.meshConfig.h2UpgradePolicy=DONOTUPGRADE",
			},
			want: fmt.Errorf("\n Unsupported value: %q, supported values for: %q is %q",
				"DONOTUPGRADE", "values.meshConfig.h2UpgradePolicy",
				strings.Join(h2UpgradePolicy, ", ")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateSetFlags(tt.args)
			if got != nil && fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestValidateDuration(t *testing.T) {

	tests := []struct {
		name     string
		flagName string
		duration string
		want     error
	}{
		{
			name:     "Test valid convertDuration",
			flagName: "values.global.proxy.dnsRefreshRate",
			duration: "10s",
			want:     nil,
		},
		{
			name:     "Test invalid convertDuration",
			flagName: "values.global.proxy.dnsRefreshRate",
			duration: "sam",
			want:     fmt.Errorf("invalid duration format %q", "sam"),
		},
		{
			name:     "Test valid protocolDetectionTimeout",
			flagName: "values.global.proxy.protocolDetectionTimeout",
			duration: "100ms",
			want:     nil,
		},
		{
			name:     "Test 0s is valid for protocolDetectionTimeout",
			flagName: "values.global.proxy.protocolDetectionTimeout",
			duration: "0s",
			want:     nil,
		},
		{
			name:     "Test invalid protocolDetectionTimeout",
			flagName: "values.global.proxy.protocolDetectionTimeout",
			duration: "-1s",
			want:     errors.New("only durations to ms precision are supported"),
		},
		{
			name:     "Test valid dnsRefreshRate",
			flagName: "values.global.proxy.dnsRefreshRate",
			duration: "10s",
			want:     nil,
		},
		{
			name:     "Test invalid dnsRefreshRate",
			flagName: "values.global.proxy.dnsRefreshRate",
			duration: "100ms",
			want:     errors.New("DNS refresh rate only supports durations to seconds precision"),
		},
		{
			name:     "Test valid connectTimeout",
			flagName: "values.global.connectTimeout",
			duration: "100ms",
			want:     nil,
		},
		{
			name:     "Test invalid connectTimeout",
			flagName: "values.global.connectTimeout",
			duration: "-1ms",
			want:     errors.New("duration must be greater than 1ms"),
		},
		{
			name:     "Test valid reportBatchMaxTime",
			flagName: "values.mixer.telemetry.reportBatchMaxTime",
			duration: "10s",
			want:     nil,
		},
		{
			name:     "Test invalid reportBatchMaxTime",
			flagName: "values.mixer.telemetry.reportBatchMaxTime",
			duration: "-1ms",
			want:     errors.New("duration must be greater than 1ms"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateDuration(tt.flagName, tt.duration)
			if got != nil && fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func TestValidateClusterStatName(t *testing.T) {

	tests := []struct {
		name      string
		flagName  string
		flagValue string
		want      bool
	}{
		{
			name:      "Test valid inboundClusterStatName",
			flagName:  "values.meshConfig.inboundClusterStatName",
			flagValue: "%SERVICE_FQDN%_%SERVICE_PORT%",
			want:      true,
		},
		{
			name:      "Test invalid inboundClusterStatName",
			flagName:  "values.meshConfig.inboundClusterStatName",
			flagValue: "%SERVICE_FQDN%_%SERVICE_POT%",
			want:      false,
		},
		{
			name:      "Test valid outboundClusterStatName",
			flagName:  "values.meshConfig.outboundClusterStatName",
			flagValue: "%SERVICE_FQDN%_%SERVICE_PORT%",
			want:      true,
		},
		{
			name:      "Test invalid outboundClusterStatName",
			flagName:  "values.meshConfig.outboundClusterStatName",
			flagValue: "%SERVICE_FQDN%_%SERVICE_POT%",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateClusterStatName(tt.flagName, tt.flagValue); got != tt.want {
				t.Errorf("Failed: got valid=%t but wanted valid=%t: %v for %v", got, tt.want, got, tt.flagValue)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {

	tests := []struct {
		name      string
		flagValue string
		want      error
	}{
		{
			name:      "Test valid port",
			flagValue: "3000",
			want:      nil,
		},
		{
			name:      "Test invalid port with string value",
			flagValue: "test",
			want: fmt.Errorf("\n Unsupported value: %q, port number must be in the range 1..65535",
				"test"),
		},
		{
			name:      "Test invalid port",
			flagValue: "0",
			want:      fmt.Errorf("port number %d must be in the range 1..65535", 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validatePort(tt.flagValue)
			if got != nil && fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
