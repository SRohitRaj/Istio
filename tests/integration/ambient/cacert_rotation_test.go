//go:build integ
// +build integ

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

package ambient

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"

	"istio.io/istio/istioctl/pkg/util/configdump"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/istioctl"
	"istio.io/istio/pkg/test/framework/components/namespace"
	kubetest "istio.io/istio/pkg/test/kube"
	"istio.io/istio/pkg/test/util/assert"
	"istio.io/istio/pkg/test/util/retry"
	"istio.io/istio/security/pkg/pki/util"
	"istio.io/istio/tests/integration/security/util/cert"
)

func TestIntermediateCertificateRefresh(t *testing.T) {
	framework.NewTest(t).
		Features("security.peer.cacert-rotation").
		Run(func(t framework.TestContext) {
			istioCfg := istio.DefaultConfigOrFail(t, t)
			istioCtl := istioctl.NewOrFail(t, t, istioctl.Config{})
			namespace.ClaimOrFail(t, t, istioCfg.SystemNamespace)
			newX509 := getX509FromFile(t, "ca-cert-alt.pem")

			sa := apps.Captured[0].ServiceAccountName()

			// we do not know which ztunnel instance is located on the node as the workload, so we need to check all of them initially
			ztunnelPods, err := kubetest.NewPodFetch(t.AllClusters()[0], istioCfg.SystemNamespace, "app=ztunnel")()
			assert.NoError(t, err)

			originalWorkloadSecret, ztunnelPod, err := getWorkloadSecret(t, ztunnelPods, sa, istioCtl)
			if err != nil {
				t.Errorf("failed to get initial workload cert: %v", err)
			}

			// Update CA with new intermediate cert
			if err := cert.CreateCustomCASecret(t,
				"ca-cert-alt.pem", "ca-key-alt.pem",
				"cert-chain-alt.pem", "root-cert-combined.pem"); err != nil {
				t.Errorf("failed to update CA secret: %v", err)
			}

			newWorkloadCert := waitForWorkloadCertUpdate(t, ztunnelPod, sa, istioCtl, originalWorkloadSecret)

			verifyWorkloadCert(t, newWorkloadCert, newX509)

			// reset CA to original values
			if err := cert.CreateCustomCASecret(t,
				"ca-cert.pem", "ca-key.pem",
				"cert-chain.pem", "root-cert.pem"); err != nil {
				t.Errorf("failed to update CA secret: %v", err)
			}

			_ = waitForWorkloadCertUpdate(t, ztunnelPod, sa, istioCtl, newWorkloadCert)
		})
}

func getWorkloadSecret(t framework.TestContext, zPods []v1.Pod, serviceAccount string, ctl istioctl.Instance) (*configdump.CertsDump, v1.Pod, error) {
	for _, ztunnel := range zPods {
		podName := fmt.Sprintf("%s.%s", ztunnel.Name, ztunnel.Namespace)
		out, errOut, err := ctl.Invoke([]string{"pc", "s", podName, "-o", "json"})
		if err != nil || errOut != "" {
			t.Errorf("failed to retrieve pod secrets from %s, err: %v errOut: %s", podName, err, errOut)
		}

		dump := []configdump.CertsDump{}
		if err := json.Unmarshal([]byte(out), &dump); err != nil {
			t.Errorf("failed to unmarshal secret dump: %v", err)
		}

		for _, s := range dump {
			if strings.Contains(s.Identity, serviceAccount) {
				if len(s.CaCert) == 0 {
					t.Errorf("ca cert missing in %v for identity: %v", ztunnel.Name, s.Identity)
				}
				if len(s.CertChain) == 0 {
					t.Errorf("cert chain missing in %v for identity: %v", ztunnel.Name, s.Identity)
				}
				return &s, ztunnel, nil
			}
		}
	}
	return nil, v1.Pod{}, errors.New("failed to find workload secret")
}

// Abstracted function to wait for workload cert to be updated
func waitForWorkloadCertUpdate(t framework.TestContext, ztunnelPod v1.Pod, serviceAccount string, istioCtl istioctl.Instance, originalCert *configdump.CertsDump) *configdump.CertsDump {
	var newSecret *configdump.CertsDump
	retry.UntilOrFail(t, func() bool {
		updatedCert, _, err := getWorkloadSecret(t, []v1.Pod{ztunnelPod}, serviceAccount, istioCtl)
		if err != nil {
			t.Logf("failed to get current workload secret: %v", err)
			return false
		}

		// retry when workload cert is not updated
		if originalCert.CaCert[0].ValidFrom != updatedCert.CaCert[0].ValidFrom {
			newSecret = updatedCert
			t.Logf("workload cert is updated")
			return true
		}

		return false
	}, retry.Timeout(5*time.Minute), retry.Delay(10*time.Second))
	return newSecret
}

func verifyWorkloadCert(t framework.TestContext, workloadSecret *configdump.CertsDump, caX590 *x509.Certificate) {
	intermediateCert, err := base64.StdEncoding.DecodeString(workloadSecret.CertChain[0].Pem)
	if err != nil {
		t.Errorf("failed to decode intermediate certificate: %v", err)
	}
	intermediateX509 := parseCert(t, intermediateCert)
	// verify the correct intermediate cert is in the certificate chain
	if intermediateX509.SerialNumber.String() != caX590.SerialNumber.String() {
		t.Errorf("intermediate certificate serial numbers do not match: got %v, wanted %v", intermediateX509.SerialNumber.String(), caX590.SerialNumber.String())
	}

	workloadCert, err := base64.StdEncoding.DecodeString(workloadSecret.CaCert[0].Pem)
	if err != nil {
		t.Errorf("failed to decode workload certificate: ", err)
	}
	workloadX509 := parseCert(t, workloadCert)

	// verify workload cert contains the correct intermediate cert
	if string(workloadX509.AuthorityKeyId) != string(caX590.SubjectKeyId) {
		t.Errorf("workload certificate did not have expected authority key id: got %v wanted %v", string(workloadX509.AuthorityKeyId), string(caX590.SubjectKeyId))
	}
}

func getX509FromFile(t framework.TestContext, caCertFile string) *x509.Certificate {
	certBytes, err := cert.ReadSampleCertFromFile(caCertFile)
	if err != nil {
		t.Errorf("failed to read %s file: %v", caCertFile, err)
	}
	return parseCert(t, certBytes)
}

func parseCert(t framework.TestContext, certBytes []byte) *x509.Certificate {
	parsedCert, err := util.ParsePemEncodedCertificate(certBytes)
	if err != nil {
		t.Errorf("failed to parse certificate pem file: %v", err)
	}
	return parsedCert
}
