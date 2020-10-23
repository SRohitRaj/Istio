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

package filemountedcerts

import (
	"context"
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/label"
	"istio.io/istio/pkg/test/framework/resource"

	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	inst istio.Instance
)

const (
	PilotCertsPath = "tests/testdata/certs/pilot"
	PilotSecretName = "test-istiod-server-cred"

	ProxyMetadataJson = `
	{
      "PILOT_CERT_PROVIDER": "mycopki",
      "FILE_MOUNTED_CERTS": "true",
      "ISTIO_META_TLS_CLIENT_CERT_CHAIN": "/client-certs/cert-chain.pem",
      "ISTIO_META_TLS_CLIENT_KEY": "/client-certs/key.pem",
      "ISTIO_META_TLS_CLIENT_ROOT_CERT": "/client-certs/root-cert.pem",
      "ISTIO_META_TLS_SERVER_CERT_CHAIN": "/server-certs/cert-chain.pem",
      "ISTIO_META_TLS_SERVER_KEY": "/server-certs/key.pem",
      "ISTIO_META_TLS_SERVER_ROOT_CERT": "/server-certs/root-cert.pem",
      "PROXY_XDS_VIA_AGENT": "false",
	}
`
)

func TestMain(m *testing.M) {
	framework.
		NewSuite(m).
		Label(label.CustomSetup).

		// SDS requires Kubernetes 1.13
		RequireEnvironmentVersion("1.13").
		Label("CustomSetup").
		Setup(istio.Setup(&inst, setupConfig, CreateCustomIstiodSecret)).
		Run()
}


func setupConfig(_ resource.Context, cfg *istio.Config) {
	if cfg == nil {
		return
	}

	cfg.ControlPlaneValues = `
components:
  egressGateways:
  - enabled: false
    name: istio-egressgateway
  ingressGateways:
  - enabled: false
    name: istio-ingressgateway
  pilot:
    enabled: true
    k8s:
      overlays:
        - kind: Deployment
          name: istiod
          patches:
            - path: spec.template.spec.containers.[name:discovery].args[1001]
              value: "--caCertFile=/server-certs/root-cert.pem"
            - path: spec.template.spec.containers.[name:discovery].args[1002]
              value: "--tlsCertFile=/server-certs/cert-chain.pem"
            - path: spec.template.spec.containers.[name:discovery].args[1003]
              value: "--tlsKeyFile=/server-certs/key.pem"
            - path: spec.template.spec.volumes[-1]
              value: |-
                name: server-certs
                secret:
                  secretName: ` + PilotSecretName + `
                  defaultMode: 420
            - path: spec.template.spec.containers[name:discovery].volumeMounts[-1]
              value: |-
                name: server-certs
                mountPath: /server-certs

meshConfig:
  defaultConfig:
    controlPlaneAuthPolicy: "MUTUAL_TLS"
    proxyMetadata: ` + strings.Replace(ProxyMetadataJson, "\n", "", -1) +
`
values:
  global:
    pilotCertProvider: "mycopki"
    logging:
      level: "default:debug"
    proxy:
      logLevel: debug
      componentLogLevel: "misc:debug"
  pilot:
    env:
      XDS_AUTH: "false"
      ENABLE_CA_SERVER: "false"
`
}

func CreateCustomIstiodSecret(ctx resource.Context) error {
	systemNs, err := istio.ClaimSystemNamespace(ctx)
	if err != nil {
		return err
	}


	err = CreateCustomSecret(ctx, PilotSecretName, systemNs, PilotCertsPath)
	if err != nil {
		return err
	}

	return nil
}

func CreateCustomSecret(ctx resource.Context, name string, namespace namespace.Instance, certsPath string) error {
	var privateKey, clientCert, caCert []byte
	var err error
	if privateKey, err = ReadCustomCertFromFile(certsPath,"key.pem"); err != nil {
		return err
	}
	if clientCert, err = ReadCustomCertFromFile(certsPath, "cert-chain.pem"); err != nil {
		return err
	}
	if caCert, err = ReadCustomCertFromFile(certsPath, "root-cert.pem"); err != nil {
		return err
	}

	kubeAccessor := ctx.Environment().(*kube.Environment).KubeClusters[0]
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace.Name(),
		},
		Data: map[string][]byte{
			"key.pem"       : privateKey,
			"cert-chain.pem": clientCert,
			"root-cert.pem" : caCert,
		},
	}

	_, err = kubeAccessor.CoreV1().Secrets(namespace.Name()).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func ReadCustomCertFromFile(certsPath string, f string) ([]byte, error) {
	b, err := ioutil.ReadFile(path.Join(env.IstioSrc, certsPath, f))
	if err != nil {
		return nil, err
	}
	return b, nil
}
