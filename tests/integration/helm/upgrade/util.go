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

package helmupgrade

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/api/label"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/cluster"
	kubecluster "istio.io/istio/pkg/test/framework/components/cluster/kube"
	"istio.io/istio/pkg/test/helm"
	kubetest "istio.io/istio/pkg/test/kube"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/istio/pkg/test/shell"
	"istio.io/istio/pkg/test/util/retry"
	helmtest "istio.io/istio/tests/integration/helm"
	"istio.io/istio/tests/util/sanitycheck"
)

const (
	gcrHub            = "gcr.io/istio-release"
	prodTag           = "prod"
	canaryTag         = "canary"
	latestRevisionTag = "latest"
)

// upgradeCharts upgrades Istio using Helm charts with the provided
// override values file to the latest charts in $ISTIO_SRC/manifests
func upgradeCharts(ctx framework.TestContext, h *helm.Helm, overrideValuesFile string, isAmbient bool) {
	execCmd := fmt.Sprintf(
		"kubectl apply -n %v -f %v",
		helmtest.IstioNamespace,
		filepath.Join(helmtest.ManifestsChartPath, helmtest.BaseChart, helmtest.CRDsFolder))
	_, err := shell.Execute(false, execCmd)
	if err != nil {
		ctx.Fatalf("couldn't run kubectl apply on crds folder: %v", err)
	}

	// Upgrade base chart
	err = h.UpgradeChart(helmtest.BaseReleaseName, filepath.Join(helmtest.ManifestsChartPath, helmtest.BaseChart),
		helmtest.IstioNamespace, overrideValuesFile, helmtest.Timeout, "--skip-crds")
	if err != nil {
		ctx.Fatalf("failed to upgrade istio %s chart", helmtest.BaseReleaseName)
	}

	// Upgrade discovery chart
	err = h.UpgradeChart(helmtest.IstiodReleaseName, filepath.Join(helmtest.ManifestsChartPath, helmtest.ControlChartsDir, helmtest.DiscoveryChartsDir),
		helmtest.IstioNamespace, overrideValuesFile, helmtest.Timeout)
	if err != nil {
		ctx.Fatalf("failed to upgrade istio %s chart", helmtest.IstiodReleaseName)
	}

	if isAmbient {
		// Upgrade istio-cni chart
		err = h.UpgradeChart(helmtest.CniReleaseName, filepath.Join(helmtest.ManifestsChartPath, helmtest.CniChartsDir),
			helmtest.IstioNamespace, overrideValuesFile, helmtest.Timeout)
		if err != nil {
			ctx.Fatalf("failed to upgrade istio %s chart", helmtest.CniReleaseName)
		}
		// Upgrade ztunnel chart
		err = h.UpgradeChart(helmtest.ZtunnelReleaseName, filepath.Join(helmtest.ManifestsChartPath, helmtest.ZtunnelChartsDir),
			helmtest.IstioNamespace, overrideValuesFile, helmtest.Timeout)
		if err != nil {
			ctx.Fatalf("failed to upgrade istio %s chart", helmtest.ZtunnelReleaseName)
		}
	}

	// Upgrade ingress gateway chart
	err = h.UpgradeChart(helmtest.IngressReleaseName, filepath.Join(helmtest.ManifestsChartPath, helmtest.GatewayChartsDir),
		helmtest.IstioNamespace, overrideValuesFile, helmtest.Timeout)
	if err != nil {
		ctx.Fatalf("failed to upgrade istio %s chart", helmtest.IngressReleaseName)
	}
}

// deleteIstio deletes installed Istio Helm charts and resources
func deleteIstio(cs cluster.Cluster, h *helm.Helm, gatewayChartInstalled bool) error {
	scopes.Framework.Infof("cleaning up resources")
	if gatewayChartInstalled {
		if err := h.DeleteChart(helmtest.IngressReleaseName, helmtest.IstioNamespace); err != nil {
			return fmt.Errorf("failed to delete %s release", helmtest.IngressReleaseName)
		}
	}

	if err := h.DeleteChart(helmtest.IstiodReleaseName, helmtest.IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", helmtest.IstiodReleaseName)
	}

	return cleanupIstio(cs, h)
}

func cleanupIstio(cs cluster.Cluster, h *helm.Helm) error {
	if err := h.DeleteChart(helmtest.BaseReleaseName, helmtest.IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", helmtest.BaseReleaseName)
	}
	if err := cs.Kube().CoreV1().Namespaces().Delete(context.TODO(), helmtest.IstioNamespace, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete istio namespace: %v", err)
	}
	if err := kubetest.WaitForNamespaceDeletion(cs.Kube(), helmtest.IstioNamespace, retry.Timeout(helmtest.RetryTimeOut)); err != nil {
		return fmt.Errorf("waiting for istio namespace to be deleted: %v", err)
	}
	return nil
}

// deleteIstioCanary deletes installed Istio Helm charts and resources
func deleteIstioRevision(h *helm.Helm, revision string) error {
	scopes.Framework.Infof("cleaning up revision resources (%s)", revision)
	name := helmtest.IstiodReleaseName + "-" + strings.ReplaceAll(revision, ".", "-")
	if err := h.DeleteChart(name, helmtest.IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete revision (%s)", name)
	}

	return nil
}

// performInPlaceUpgradeFunc returns the provided function necessary to run inside an integration test
// for upgrade capability
func performInPlaceUpgradeFunc(previousVersion string, isAmbient bool) func(framework.TestContext) {
	return func(t framework.TestContext) {
		cs := t.Clusters().Default().(*kubecluster.Cluster)
		h := helm.New(cs.Filename())

		t.CleanupConditionally(func() {
			// only need to do call this once as helm doesn't need to remove
			// all versions
			helmtest.DeleteIstio(t, h, cs, isAmbient)
		})
		overrideValuesFile := helmtest.GetValuesOverrides(t, gcrHub, previousVersion, "", isAmbient)
		helmtest.InstallIstio(t, cs, h, overrideValuesFile, previousVersion, true, isAmbient)
		helmtest.VerifyInstallation(t, cs, true, isAmbient)

		_, oldClient, oldServer := sanitycheck.SetupTrafficTest(t, t, "")
		sanitycheck.RunTrafficTestClientServer(t, oldClient, oldServer)

		s := t.Settings()
		overrideValuesFile = helmtest.GetValuesOverrides(t, s.Image.Hub, s.Image.Tag, "", isAmbient)
		upgradeCharts(t, h, overrideValuesFile, isAmbient)
		helmtest.VerifyInstallation(t, cs, true, isAmbient)

		_, newClient, newServer := sanitycheck.SetupTrafficTest(t, t, "")
		sanitycheck.RunTrafficTestClientServer(t, newClient, newServer)

		// now check that we are compatible with N-1 proxy with N proxy
		sanitycheck.RunTrafficTestClientServer(t, oldClient, newServer)
	}
}

// performCanaryUpgradeFunc returns the provided function necessary to run inside an integration test
// for upgrade capability with revisions
func performCanaryUpgradeFunc(previousVersion string) func(framework.TestContext) {
	return func(t framework.TestContext) {
		cs := t.Clusters().Default().(*kubecluster.Cluster)
		h := helm.New(cs.Filename())
		t.CleanupConditionally(func() {
			err := deleteIstioRevision(h, canaryTag)
			if err != nil {
				t.Fatalf("could not delete istio: %v", err)
			}
			err = deleteIstio(cs, h, false)
			if err != nil {
				t.Fatalf("could not delete istio: %v", err)
			}
		})

		overrideValuesFile := helmtest.GetValuesOverrides(t, gcrHub, previousVersion, "", false)
		helmtest.InstallIstio(t, cs, h, overrideValuesFile, previousVersion, false, false)
		helmtest.VerifyInstallation(t, cs, false, false)

		_, oldClient, oldServer := sanitycheck.SetupTrafficTest(t, t, "")
		sanitycheck.RunTrafficTestClientServer(t, oldClient, oldServer)

		s := t.Settings()
		overrideValuesFile = helmtest.GetValuesOverrides(t, s.Image.Hub, s.Image.Tag, canaryTag, false)
		helmtest.InstallIstioWithRevision(t, cs, h, "", canaryTag, overrideValuesFile, true, false)
		helmtest.VerifyInstallation(t, cs, false, false)

		// now that we've installed with a revision we have a new mutating webhook
		helmtest.VerifyMutatingWebhookConfigurations(t, cs, []string{
			"istio-sidecar-injector",
			"istio-sidecar-injector-canary",
		})

		_, newClient, newServer := sanitycheck.SetupTrafficTest(t, t, canaryTag)
		sanitycheck.RunTrafficTestClientServer(t, newClient, newServer)

		// now check that we are compatible with N-1 proxy with N proxy
		sanitycheck.RunTrafficTestClientServer(t, oldClient, newServer)
	}
}

// performRevisionTagsUpgradeFunc returns the provided function necessary to run inside an integration test
// for upgrade capability with stable label revision upgrades
func performRevisionTagsUpgradeFunc(previousVersion string) func(framework.TestContext) {
	return func(t framework.TestContext) {
		cs := t.Clusters().Default().(*kubecluster.Cluster)
		h := helm.New(cs.Filename())
		t.CleanupConditionally(func() {
			err := deleteIstioRevision(h, latestRevisionTag)
			if err != nil {
				t.Fatalf("could not delete istio revision (%v): %v", latestRevisionTag, err)
			}
			err = deleteIstioRevision(h, previousVersion)
			if err != nil {
				t.Fatalf("could not delete istio revision (%v): %v", previousVersion, err)
			}

			err = cleanupIstio(cs, h)
			if err != nil {
				t.Fatalf("could not cleanup istio: %v", err)
			}
		})

		// install MAJOR.MINOR.PATCH charts with revision set to "MAJOR-MINOR-PATCH" name. For example,
		// helm install istio-base istio/base --version 1.15.0 --namespace istio-system -f values.yaml
		// helm install istiod-1-15 istio/istiod --version 1.15.0 -f values.yaml
		previousRevision := strings.ReplaceAll(previousVersion, ".", "-")
		overrideValuesFile := helmtest.GetValuesOverrides(t, gcrHub, previousVersion, previousRevision, false)
		helmtest.InstallIstioWithRevision(t, cs, h, previousVersion, previousRevision, overrideValuesFile, false, true)
		helmtest.VerifyInstallation(t, cs, false, false)

		// helm template istiod-1-15-0 istio/istiod --version 1.15.0 -s templates/revision-tags.yaml --set revision=1-15-0 --set revisionTags={prod}
		helmtest.SetRevisionTagWithVersion(t, h, previousRevision, prodTag, previousVersion)
		helmtest.VerifyMutatingWebhookConfigurations(t, cs, []string{
			"istio-revision-tag-prod",
			fmt.Sprintf("istio-sidecar-injector-%s", previousRevision),
		})

		// setup istio.io/rev=1-15-0 for the default-1 namespace
		oldNs, oldClient, oldServer := sanitycheck.SetupTrafficTest(t, t, previousRevision)
		sanitycheck.RunTrafficTestClientServer(t, oldClient, oldServer)

		// install the charts from this branch with revision set to "latest"
		// helm upgrade istio-base ../manifests/charts/base --namespace istio-system -f values.yaml
		// helm install istiod-latest ../manifests/charts/istio-control/istio-discovery -f values.yaml
		s := t.Settings()
		overrideValuesFile = helmtest.GetValuesOverrides(t, s.Image.Hub, s.Image.Tag, latestRevisionTag, false)
		helmtest.InstallIstioWithRevision(t, cs, h, "", latestRevisionTag, overrideValuesFile, true, false)
		helmtest.VerifyInstallation(t, cs, false, false)

		// helm template istiod-latest ../manifests/charts/istio-control/istio-discovery --namespace istio-system
		//    -s templates/revision-tags.yaml --set revision=latest --set revisionTags={canary}
		helmtest.SetRevisionTag(t, h, "", latestRevisionTag, canaryTag, helmtest.ManifestsChartPath, "")
		helmtest.VerifyMutatingWebhookConfigurations(t, cs, []string{
			"istio-revision-tag-prod",
			fmt.Sprintf("istio-sidecar-injector-%v", previousRevision),
			"istio-revision-tag-canary",
			"istio-sidecar-injector-latest",
		})

		// setup istio.io/rev=latest for the default-2 namespace
		_, newClient, newServer := sanitycheck.SetupTrafficTest(t, t, latestRevisionTag)
		sanitycheck.RunTrafficTestClientServer(t, newClient, newServer)

		// now check that we are compatible with N-1 proxy with N proxy between a client
		// in default-1 namespace and a server in the default-2 namespace, respectively
		sanitycheck.RunTrafficTestClientServer(t, oldClient, newServer)

		// change the mutating webhook configuration to use the latest revision (istiod-latest service in istio-system)
		// helm template istiod-latest ../manifests/charts/istio-control/istio-discovery --namespace istio-system
		//    -s templates/revision-tags.yaml --set revision=latest --set revisionTags={prod}
		helmtest.SetRevisionTag(t, h, "", latestRevisionTag, prodTag, helmtest.ManifestsChartPath, "")

		// change the old namespace that was pointing to the old prod (1-15-0) to point to the
		// 'latest' revision by setting the `istio.io/rev=prod` label on the namespace
		err := oldNs.SetLabel(label.IoIstioRev.Name, prodTag)
		if err != nil {
			t.Fatal("could not remove istio.io/rev from old namespace")
		}

		err = oldClient.Restart()
		if err != nil {
			t.Fatal("could not restart old client")
		}
		err = oldServer.Restart()
		if err != nil {
			t.Fatal("could not restart old server")
		}

		// make sure the restarted pods in default-1 namespace do not use
		// the previous version (check for the previousVersion in the image string)
		err = checkVersion(t, oldNs.Name(), previousVersion)
		if err != nil {
			t.Fatalf("found a pod in namespace (%s) with the previous version: %v", oldNs.Name(), err)
		}

		// now check traffic still works between the proxies
		sanitycheck.RunTrafficTestClientServer(t, oldClient, newServer)
	}
}

func checkVersion(t framework.TestContext, namespace, version string) error {
	// func NewPodFetch(a istioKube.CLIClient, namespace string, selectors ...string) PodFetchFunc {
	fetch := kubetest.NewPodFetch(t.Clusters().Default(), namespace)
	pods, err := kubetest.CheckPodsAreReady(fetch)
	if err != nil {
		return fmt.Errorf("failed to retrieve pods: %v", err)
	}
	for _, p := range pods {
		for _, c := range p.Spec.Containers {
			if strings.Contains(c.Image, version) {
				return fmt.Errorf("expected container image to not include version %q, got %q", version, c.Image)
			}
		}
	}

	return nil
}
