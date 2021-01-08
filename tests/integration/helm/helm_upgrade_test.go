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

package helm

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/framework/image"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/pkg/test/helm"
	kubetest "istio.io/istio/pkg/test/kube"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/istio/pkg/test/util/retry"
	"istio.io/istio/tests/util/sanitycheck"
)

var (
	// previousChartPath is path of Helm charts for previous Istio deployments.
	previousChartPath = filepath.Join(env.IstioSrc, "tests/integration/helm/testdata/")
)

const (
	gcrHub                   = "gcr.io/istio-release"
	previousSupportedVersion = "1.8.1"

	defaultValues = `
global:
  hub: %s
  tag: %s
`

	firstPartyJwtValues = `
global:
  hub: %s
  tag: %s
  jwtPolicy: first-party-jwt
`
)

// TestDefaultInPlaceUpgrades tests Istio installation using Helm with default options
func TestDefaultInPlaceUpgrades(t *testing.T) {
	framework.
		NewTest(t).
		Features("installation.helm.default.upgrade").
		Run(func(ctx framework.TestContext) {
			cs := ctx.Clusters().Default().(*kube.Cluster)
			h := helm.New(cs.Filename(), filepath.Join(previousChartPath, previousSupportedVersion))

			ctx.WhenDone(func() error {
				// only need to do call this once as helm doesn't need to remove
				// all versions
				return deleteIstio(t, cs, h)
			})

			overrideValuesFile := getValuesOverrides(ctx, defaultValues, gcrHub, previousSupportedVersion)
			installIstio(t, cs, h, overrideValuesFile)
			verifyInstallation(ctx, cs)

			oldClient, oldServer := sanitycheck.SetupTrafficTest(t, ctx)
			sanitycheck.RunTrafficTestClientServer(t, oldClient, oldServer)

			// now upgrade istio to the latest version found in this branch
			// use the command line or environmental vars from the user to set
			// the hub/tag
			s, err := image.SettingsFromCommandLine()
			if err != nil {
				ctx.Fatal(err)
			}

			overrideValuesFile = getValuesOverrides(ctx, defaultValues, s.Hub, s.Tag)
			upgradeCharts(ctx, h, overrideValuesFile)
			verifyInstallation(ctx, cs)

			newClient, newServer := sanitycheck.SetupTrafficTest(t, ctx)
			sanitycheck.RunTrafficTestClientServer(t, newClient, newServer)

			// now check that we are compatible with N-1 proxy with N proxy
			sanitycheck.RunTrafficTestClientServer(t, oldClient, newServer)
		})
}

// upgradeCharts upgrades Istio using Helm charts with the provided
// override values file to the latest charts in $ISTIO_SRC/manifests
func upgradeCharts(ctx framework.TestContext, h *helm.Helm, overrideValuesFile string) {

	// Upgrade base chart
	err := h.UpgradeChart(BaseReleaseName, filepath.Join(ChartPath, BaseChart),
		IstioNamespace, overrideValuesFile, helmTimeout)
	if err != nil {
		ctx.Fatalf("failed to upgrade istio %s chart", BaseChart)
	}

	// Upgrade discovery chart
	err = h.UpgradeChart(IstiodReleaseName, filepath.Join(ChartPath, ControlChartsDir, DiscoveryChart),
		IstioNamespace, overrideValuesFile, helmTimeout)
	if err != nil {
		ctx.Fatalf("failed to upgrade istio %s chart", DiscoveryChart)
	}

	// Upgrade ingress gateway chart
	err = h.UpgradeChart(IngressReleaseName, filepath.Join(ChartPath, GatewayChartsDir, IngressGatewayChart),
		IstioNamespace, overrideValuesFile, helmTimeout)
	if err != nil {
		ctx.Fatalf("failed to upgrade istio %s chart", IngressGatewayChart)
	}

	// Upgrade egress gateway chart
	err = h.UpgradeChart(EgressReleaseName, filepath.Join(ChartPath, GatewayChartsDir, EgressGatewayChart),
		IstioNamespace, overrideValuesFile, helmTimeout)
	if err != nil {
		ctx.Fatalf("failed to upgrade istio %s chart", EgressGatewayChart)
	}
}

// installIstio install Istio using Helm charts with the provided
// override values file and fails the tests on any failures.
func installIstio(t *testing.T, cs resource.Cluster,
	h *helm.Helm, overrideValuesFile string) {
	createIstioSystemNamespace(t, cs)

	// Install base chart
	err := h.InstallChart(BaseReleaseName, BaseChart,
		IstioNamespace, overrideValuesFile, helmTimeout)
	if err != nil {
		t.Errorf("failed to install istio %s chart", BaseChart)
	}

	// Install discovery chart
	err = h.InstallChart(IstiodReleaseName, filepath.Join(ControlChartsDir, DiscoveryChart),
		IstioNamespace, overrideValuesFile, helmTimeout)
	if err != nil {
		t.Errorf("failed to install istio %s chart", DiscoveryChart)
	}

	installGatewaysCharts(t, cs, h, overrideValuesFile)
}

// deleteIstio deletes installed Istio Helm charts and resources
func deleteIstio(t *testing.T, cs resource.Cluster, h *helm.Helm) error {
	scopes.Framework.Infof("cleaning up resources")
	if err := h.DeleteChart(EgressReleaseName, IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", EgressReleaseName)
	}
	if err := h.DeleteChart(IngressReleaseName, IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", IngressReleaseName)
	}
	if err := h.DeleteChart(IstiodReleaseName, IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", IngressReleaseName)
	}
	if err := h.DeleteChart(BaseReleaseName, IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", BaseReleaseName)
	}
	if err := cs.CoreV1().Namespaces().Delete(context.TODO(), IstioNamespace, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete istio namespace: %v", err)
	}
	if err := kubetest.WaitForNamespaceDeletion(cs, IstioNamespace, retry.Timeout(retryTimeOut)); err != nil {
		return fmt.Errorf("wating for istio namespace to be deleted: %v", err)
	}

	return nil
}
