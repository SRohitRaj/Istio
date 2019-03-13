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

package bookinfo

import (
	"fmt"
	"path"

	"istio.io/istio/pkg/test/framework2/core"

	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework2/components/environment/kube"

	"istio.io/istio/pkg/test/scopes"
)

type bookInfoConfig string

const (
	// BookInfoConfig uses "bookinfo.yaml"
	BookInfoConfig    bookInfoConfig = "bookinfo.yaml"
	BookinfoRatingsv2 bookInfoConfig = "bookinfo-ratings-v2.yaml"
	BookinfoDb        bookInfoConfig = "bookinfo-db.yaml"
)

var _ Instance = &kubeComponent{}

// NewKubeComponent factory function for the component
func newKube(ctx core.Context) (Instance, error) {
	c := &kubeComponent{}
	c.id = ctx.TrackResource(c)

	if err := deployBookInfoService(ctx, string(BookInfoConfig)); err != nil {
		return nil, err
	}

	return c, nil
}

type kubeComponent struct {
	id core.ResourceID
}

func (c *kubeComponent) ID() core.ResourceID {
	return c.id
}

// DeployRatingsV2 deploys ratings v2 service
func (c *kubeComponent) DeployRatingsV2(ctx core.Context) (err error) {
	return deployBookInfoService(ctx, string(BookinfoRatingsv2))
}

// DeployMongoDb deploys mongodb service
func (c *kubeComponent) DeployMongoDb(ctx core.Context) (err error) {
	return deployBookInfoService(ctx, string(BookinfoDb))
}

func deployBookInfoService(ctx core.Context, bookinfoYamlFile string) (err error) {
	e := ctx.Environment().(*kube.Environment)

	scopes.CI.Infof("=== BEGIN: Deploy BookInfoConfig (via Yaml File - %s) ===", bookinfoYamlFile)
	defer func() {
		if err != nil {
			err = fmt.Errorf("BookInfoConfig %s deployment failed: %v", bookinfoYamlFile, err) // nolint:golint
			scopes.CI.Infof("=== FAILED: Deploy BookInfoConfig %s ===", bookinfoYamlFile)
		} else {
			scopes.CI.Infof("=== SUCCEEDED: Deploy BookInfoConfig %s ===", bookinfoYamlFile)
		}
	}()

	ns, err := e.AllocateNamespace("bookinfo", true)
	if err != nil {
		return err
	}

	yamlFile := path.Join(env.BookInfoKube, bookinfoYamlFile)
	_, err = e.DeployYaml(ns.Name(), yamlFile)
	return
}
