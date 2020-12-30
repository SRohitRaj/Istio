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

package istio

import (
	"context"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/framework/components/istio/ingress"
	"istio.io/istio/pkg/test/framework/image"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/pkg/test/helm"
	kube2 "istio.io/istio/pkg/test/kube"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/istio/pkg/test/util/retry"
)

const (
	ReleasePrefix     = "istio-"
	BaseChart         = "base"
	DiscoveryChart    = "istio-discovery"
	BaseReleaseName   = ReleasePrefix + BaseChart
	IstiodReleaseName = "istiod"
	ControlChartsDir  = "istio-control"
	retryTimeOut      = 5 * time.Minute
	helmTimeout       = 2 * time.Minute
)

var _ io.Closer = &helmComponent{}
var _ Instance = &helmComponent{}
var _ resource.Dumper = &helmComponent{}

var (
	// chartBasePath is path of local Helm charts used for testing.
	chartBasePath = filepath.Join(env.IstioSrc, "manifests/charts")
	// versionedChartBasePath is path of checked-in charts for older Istio versions
	versionedChartBasePath = filepath.Join(env.IstioSrc, "tests/integration/helm/testdata")
)

type helmComponent struct {
	id         resource.ID
	settings   Config
	helmCmd    *helm.Helm
	cs         resource.Cluster
	deployTime time.Duration
}

func (h *helmComponent) Dump(ctx resource.Context) {
	scopes.Framework.Errorf("=== Dumping Istio Deployment State...")
	ns := h.settings.SystemNamespace
	d, err := ctx.CreateTmpDirectory("istio-state")
	if err != nil {
		scopes.Framework.Errorf("Unable to create directory for dumping Istio contents: %v", err)
		return
	}
	kube2.DumpPods(ctx, d, ns)
}

func (h *helmComponent) ID() resource.ID {
	return h.id
}

func (h *helmComponent) IngressFor(cluster resource.Cluster) ingress.Instance {
	panic("implement me")
}

func (h *helmComponent) CustomIngressFor(cluster resource.Cluster, serviceName, istioLabel string) ingress.Instance {
	panic("implement me")
}

func (h *helmComponent) RemoteDiscoveryAddressFor(cluster resource.Cluster) (net.TCPAddr, error) {
	panic("implement me")
}

func (h *helmComponent) Settings() Config {
	return h.settings
}

func (h *helmComponent) Close() error {
	scopes.Framework.Infof("cleaning up resources")
	// TODO remove ingress and egress charts
	istiodRelease := IstiodReleaseName + h.settings.Revision
	if err := h.helmCmd.DeleteChart(istiodRelease, h.settings.IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", istiodRelease)
	}
	if err := h.helmCmd.DeleteChart(BaseReleaseName, h.settings.IstioNamespace); err != nil {
		return fmt.Errorf("failed to delete %s release", BaseReleaseName)
	}
	if err := h.cs.CoreV1().Namespaces().Delete(context.Background(), h.settings.IstioNamespace, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete istio namespace: %v", err)
	}
	if err := kube2.WaitForNamespaceDeletion(h.cs, h.settings.IstioNamespace, retry.Timeout(retryTimeOut)); err != nil {
		return fmt.Errorf("wating for istio namespace to be deleted: %v", err)
	}

	return nil
}

func deployWithHelm(ctx resource.Context, env *kube.Environment, cfg Config) (Instance, error) {
	scopes.Framework.Infof("=== Istio Component Config ===")
	scopes.Framework.Infof("\n%s", cfg.String())
	scopes.Framework.Infof("================================")

	// install control plane clusters
	cluster := ctx.Clusters().Default().(*kube.Cluster)

	var helmCmd *helm.Helm
	if cfg.Version == "" {
		helmCmd = helm.New(cluster.Filename(), chartBasePath)
	} else {
		chartPath := filepath.Join(versionedChartBasePath, cfg.Version, "charts")
		helmCmd = helm.New(cluster.Filename(), chartPath)
	}

	h := &helmComponent{
		settings: cfg,
		helmCmd:  helmCmd,
		cs:       cluster,
	}

	t0 := time.Now()
	defer func() {
		h.deployTime = time.Since(t0)
	}()
	h.id = ctx.TrackResource(h)

	if !cfg.DeployIstio {
		scopes.Framework.Info("skipping helm deployment as specified in the config")
		return h, nil
	}

	if env.IsMulticluster() || len(ctx.Clusters()) > 1 {
		scopes.Framework.Error("multicluster support not implemented for helm deployments")
		return h, nil
	}

	err := helmInstall(h)
	if err != nil {
		scopes.Framework.Error("multicluster support not implemented for helm deployments")
		return h, err
	}

	return h, nil
}

func helmInstall(h *helmComponent) error {
	scopes.Framework.Infof("setting up %s as control-plane cluster", h.cs.Name())

	if !h.cs.IsConfig() {
		return fmt.Errorf("cluster is not config cluster")
	}

	if _, err := h.cs.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: h.settings.SystemNamespace,
		},
	}, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			if _, err := h.cs.CoreV1().Namespaces().Update(context.Background(), &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: h.settings.SystemNamespace,
				},
			}, metav1.UpdateOptions{}); err != nil {
				scopes.Framework.Errorf("failed updating namespace %s on cluster %s. This can happen when deploying "+
					"multiple control planes. Error: %v", h.settings.SystemNamespace, h.cs.Name(), err)
			}
		} else {
			scopes.Framework.Errorf("failed creating namespace %s on cluster %s. This can happen when deploying "+
				"multiple control planes. Error: %v", h.settings.SystemNamespace, h.cs.Name(), err)
		}
	}

	overridesArgs, err := generateCommonInstallSettings(h.settings)
	if err != nil {
		return fmt.Errorf("failed to install istio %s chart", BaseChart)
	}

	// Install base chart
	err = h.helmCmd.InstallChartWithValues(BaseReleaseName, BaseChart,
		h.settings.IstioNamespace, overridesArgs, helmTimeout)
	if err != nil {
		return fmt.Errorf("failed to install istio %s chart", BaseChart)
	}

	// Install discovery chart
	err = h.helmCmd.InstallChartWithValues(IstiodReleaseName+h.settings.Revision, filepath.Join(ControlChartsDir, DiscoveryChart),
		h.settings.IstioNamespace, overridesArgs, helmTimeout)
	if err != nil {
		return fmt.Errorf("failed to install istio %s chart", DiscoveryChart)
	}

	return nil
}

func generateCommonInstallSettings(cfg Config) ([]string, error) {
	s, err := image.SettingsFromCommandLine()
	if err != nil {
		return nil, err
	}

	installSettings := []string{
		"--set", fmt.Sprintf("%s=%s", image.ImagePullPolicyValuesKey, s.PullPolicy),
	}

	// Include all user-specified values.
	for k, v := range cfg.Values {
		if cfg.Version != "" {
			if k == image.TagValuesKey || k == image.HubValuesKey {
				continue
			}
		}
		installSettings = append(installSettings, "--set", fmt.Sprintf("%s=%s", k, v))
	}

	if cfg.Revision != "" {
		installSettings = append(installSettings,
			"--set", fmt.Sprintf("%s=%s", "revision", cfg.Revision))
	}

	// set the tag and hub to custom values if version specified
	if cfg.Version != "" {
		installSettings = append(installSettings,
			"--set", fmt.Sprintf("%s=%s", image.TagValuesKey, cfg.Version),
			"--set", fmt.Sprintf("%s=%s", image.HubValuesKey, "gcr.io/istio-release"))
	}

	return installSettings, nil
}
