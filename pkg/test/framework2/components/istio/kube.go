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

package istio

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"istio.io/istio/pkg/test/framework2/core"

	"istio.io/istio/pkg/test/deployment"
	"istio.io/istio/pkg/test/framework2/components/environment/kube"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/istio/pkg/test/util/retry"
)

type kubeComponent struct {
	id          core.ResourceID
	settings    *Config
	environment *kube.Environment
	deployment  *deployment.Instance
}

var _ io.Closer = &kubeComponent{}
var _ Instance = &kubeComponent{}

func deploy(ctx core.Context, env *kube.Environment, cfg *Config) (Instance, error) {
	scopes.CI.Infof("=== Istio Component Config ===")
	scopes.CI.Infof("\n%s", cfg.String())
	scopes.CI.Infof("HUB: %s", env.Settings().Hub)
	scopes.CI.Infof("TAG: %s", env.Settings().Hub)
	scopes.CI.Infof("================================")

	i := &kubeComponent{
		environment: env,
		settings:    cfg,
	}
	i.id = ctx.TrackResource(i)

	if !cfg.DeployIstio {
		scopes.Framework.Info("skipping deployment due to Config")
		return i, nil
	}

	helmDir, err := ctx.CreateTmpDirectory("istio")
	if err != nil {
		return nil, err
	}

	generatedYaml, err := generateIstioYaml(helmDir, cfg)
	if err != nil {
		return nil, err
	}

	// split installation & configuration into two distinct steps int
	installYaml, configureYaml := splitIstioYaml(generatedYaml)

	installYamlFilePath := path.Join(helmDir, "istio-install.yaml")
	if err = ioutil.WriteFile(installYamlFilePath, []byte(installYaml), os.ModePerm); err != nil {
		return nil, fmt.Errorf("unable to write helm generated yaml: %v", err)
	}

	configureYamlFilePath := path.Join(helmDir, "istio-configure.yaml")
	if err = ioutil.WriteFile(configureYamlFilePath, []byte(configureYaml), os.ModePerm); err != nil {
		return nil, fmt.Errorf("unable to write helm generated yaml: %v", err)
	}

	scopes.CI.Infof("Created Helm-generated Yaml file(s): %s, %s", installYamlFilePath, configureYamlFilePath)
	i.deployment = deployment.NewYamlDeployment(cfg.SystemNamespace, installYamlFilePath)

	if err = i.deployment.Deploy(env.Accessor, true, retry.Timeout(cfg.DeployTimeout)); err != nil {
		return nil, err
	}

	if err = env.Accessor.Apply(cfg.SystemNamespace, configureYamlFilePath); err != nil {
		return nil, err
	}

	return i, nil
}

// ID implements resource.Instance
func (i *kubeComponent) ID() core.ResourceID {
	return i.id
}

func (i *kubeComponent) Settings() *Config {
	s := *i.settings
	return &s
}

func (i *kubeComponent) Close() error {
	if i.settings.DeployIstio {
		// TODO: There is a problem with  orderly cleanup. Re-enable this once it is fixed. Delete the system namespace
		// instead
		//return i.deployment.Delete(i.environment.Accessor, true, retry.Timeout(s.DeployTimeout))
		return i.environment.Accessor.DeleteNamespace(i.settings.SystemNamespace)
	}

	return nil
}
