//  Copyright 2019 Istio Authors
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

package policy

import (
	"path"
	"testing"

	"istio.io/istio/pkg/test/scopes"

	"istio.io/istio/pkg/test/framework/components/environment/kube"
)

const (
	// The directory that contains the test data (e.g., test policies).
	// When using this util, please place the test data under this directory.
	testDataDir = "testdata"
)

type TestPolicy struct {
	t         *testing.T
	env       *kube.Environment
	namespace string
	fileName  string
}

func (p TestPolicy) TearDown() {
	scopes.CI.Infof("Tearing down policy %q.", p.fileName)
	if err := p.env.Delete(p.namespace, path.Join(testDataDir, p.fileName)); err != nil {
		p.t.Fatalf("Cannot delete %q from namespace %q: %v", p.fileName, p.namespace, err)
	}
}

func ApplyPolicyFile(t *testing.T, env *kube.Environment, namespace string, fileName string) *TestPolicy {
	scopes.CI.Infof("Applying policy file %v", fileName)
	if err := env.Apply(namespace, path.Join(testDataDir, fileName)); err != nil {
		t.Fatalf("Cannot apply %q to namespace %q: %v", fileName, namespace, err)
		return nil
	}
	return &TestPolicy{
		t:         t,
		env:       env,
		namespace: namespace,
		fileName:  fileName,
	}
}
