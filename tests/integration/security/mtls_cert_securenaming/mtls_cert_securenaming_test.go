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

package mtlscertsecurenaming

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"istio.io/istio/pkg/config/constants"
	"istio.io/istio/pkg/test/echo/common/scheme"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/echoboot"
	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/resource/environment"
	"istio.io/istio/pkg/test/util/retry"
	"istio.io/istio/tests/integration/security/util"
	"istio.io/istio/tests/integration/security/util/cert"
	"istio.io/istio/tests/integration/security/util/connection"
)

const (
	// The length of the example certificate chain.
	exampleCertChainLength = 3

	defaultIdentityDR = `apiVersion: "networking.istio.io/v1alpha3"
kind: "DestinationRule"
metadata:
  name: "service-b-dr"
spec:
  host: "b.NS.svc.cluster.local"
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
      subjectAltNames:
      - "spiffe://cluster.local/ns/NS/sa/default"
`
	correctIdentityDR = `apiVersion: "networking.istio.io/v1alpha3"
kind: "DestinationRule"
metadata:
  name: "service-b-dr"
spec:
  host: "b.NS.svc.cluster.local"
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
      subjectAltNames:
      - "spiffe://cluster.local/ns/NS/sa/b"
`
	nonExistIdentityDR = `apiVersion: "networking.istio.io/v1alpha3"
kind: "DestinationRule"
metadata:
  name: "service-b-dr"
spec:
  host: "b.NS.svc.cluster.local"
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
      subjectAltNames:
      - "I-do-not-exist"
`
	identityListDR = `apiVersion: "networking.istio.io/v1alpha3"
kind: "DestinationRule"
metadata:
  name: "service-b-dr"
spec:
  host: "b.NS.svc.cluster.local"
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL
      subjectAltNames:
      - "spiffe://cluster.local/ns/NS/sa/a"
      - "spiffe://cluster.local/ns/NS/sa/b"
      - "spiffe://cluster.local/ns/default/sa/default"
      - "I-do-not-exist"
`
)

// TestMTLSCertSecureNaming verifies:
// - The certificate issued by CA to the sidecar is as expected and that strict mTLS works as expected.
// - The CA certificate in the configmap of each namespace is as expected, which
//   is used for data plane to control plane TLS authentication.
func TestMTLSCertSecureNaming(t *testing.T) {
	framework.NewTest(t).
		RequiresEnvironment(environment.Kube).
		Run(func(ctx framework.TestContext) {
			istioCfg := istio.DefaultConfigOrFail(t, ctx)
			testNamespace := namespace.NewOrFail(t, ctx, namespace.Config{
				Prefix: "test-mtls",
				Inject: true,
			})
			namespace.ClaimOrFail(t, ctx, istioCfg.SystemNamespace)
			// Check that the CA certificate in the configmap of each namespace is as expected, which
			// is used for data plane to control plane TLS authentication.
			retry.UntilSuccessOrFail(t, func() error {
				return checkCACert(ctx, t, testNamespace)
			}, retry.Delay(time.Second), retry.Timeout(10*time.Second))
			var a, b echo.Instance
			echoboot.NewBuilderOrFail(ctx, ctx).
				With(&a, util.EchoConfig("a", testNamespace, false, nil, g, p)).
				With(&b, util.EchoConfig("b", testNamespace, false, nil, g, p)).
				BuildOrFail(t)

			ctx.NewSubTest("mTLS cert validation").
				Run(func(ctx framework.TestContext) {
					// Verify that the certificate issued to the sidecar is as expected.
					connectTarget := fmt.Sprintf("b.%s:80", testNamespace.Name())
					out, err := cert.DumpCertFromSidecar(testNamespace, "app=a", "istio-proxy",
						connectTarget)
					if err != nil {
						t.Fatalf("failed to dump certificate: %v", err)
						return
					}
					verifyCertificates(t, out)

					// Verify mTLS works between a and b
					callOptions := echo.CallOptions{
						Target:   b,
						PortName: "http",
						Scheme:   scheme.HTTP,
					}
					checker := connection.Checker{
						From:          a,
						Options:       callOptions,
						ExpectSuccess: true,
					}
					checker.CheckOrFail(ctx)
				})

			secureNamingTestCases := []struct {
				name            string
				destinationRule string
				expectSuccess   bool
			}{
				{
					name:            "connection fails when DR doesn't match SA",
					destinationRule: defaultIdentityDR,
					expectSuccess:   false,
				},
				{
					name:            "connection succeeds when DR matches SA",
					destinationRule: correctIdentityDR,
					expectSuccess:   true,
				},
				{
					name:            "connection fails when DR contains non-matching, non-existing SA",
					destinationRule: nonExistIdentityDR,
					expectSuccess:   false,
				},
				{
					name:            "connection succeeds when SA is in the list of SANs",
					destinationRule: identityListDR,
					expectSuccess:   true,
				},
			}
			for _, tc := range secureNamingTestCases {
				ctx.NewSubTest(tc.name).
					Run(func(ctx framework.TestContext) {
						dr := strings.ReplaceAll(tc.destinationRule, "NS", testNamespace.Name())
						g.ApplyConfigOrFail(t, testNamespace, dr)
						// Verify mTLS works between a and b
						callOptions := echo.CallOptions{
							Target:   b,
							PortName: "http",
							Scheme:   scheme.HTTP,
						}
						checker := connection.Checker{
							From:          a,
							Options:       callOptions,
							ExpectSuccess: tc.expectSuccess,
						}
						if err := retry.UntilSuccess(
							checker.Check, retry.Delay(time.Second), retry.Timeout(15*time.Second), retry.Converge(5)); err != nil {
							ctx.Fatal(err)
						}
					})
			}
		})
}

func verifyCertificates(t *testing.T, dump string) {
	var certExp = regexp.MustCompile("(?sU)-----BEGIN CERTIFICATE-----(.+)-----END CERTIFICATE-----")
	certs := certExp.FindAll([]byte(dump), -1)
	// Verify that the certificate chain length is as expected
	if len(certs) != exampleCertChainLength {
		t.Errorf("expect %v certs in the cert chain but getting %v certs",
			exampleCertChainLength, len(certs))
		return
	}
	var rootCert []byte
	var err error
	if rootCert, err = cert.ReadSampleCertFromFile("root-cert.pem"); err != nil {
		t.Errorf("error when reading expected CA cert: %v", err)
		return
	}
	// Verify that the CA certificate is as expected
	if strings.TrimSpace(string(rootCert)) != strings.TrimSpace(string(certs[2])) {
		t.Errorf("the actual CA cert is different from the expected. expected: %v, actual: %v",
			strings.TrimSpace(string(rootCert)), strings.TrimSpace(string(certs[2])))
		return
	}
	t.Log("the CA certificate is as expected")
}

func checkCACert(testCtx framework.TestContext, t *testing.T, testNamespace namespace.Instance) error {
	configMapName := "istio-ca-root-cert"
	kEnv := testCtx.Environment().(*kube.Environment)
	cm, err := kEnv.KubeClusters[0].GetConfigMap(configMapName, testNamespace.Name())
	if err != nil {
		return err
	}
	var certData string
	var pluginCert []byte
	var ok bool
	if certData, ok = cm.Data[constants.CACertNamespaceConfigMapDataName]; !ok {
		return fmt.Errorf("CA certificate %v not found", constants.CACertNamespaceConfigMapDataName)
	}
	t.Logf("CA certificate %v found", constants.CACertNamespaceConfigMapDataName)
	if pluginCert, err = cert.ReadSampleCertFromFile("root-cert.pem"); err != nil {
		return err
	}
	if string(pluginCert) != certData {
		return fmt.Errorf("CA certificate (%v) not matching plugin cert (%v)", certData, string(pluginCert))
	}

	return nil
}
