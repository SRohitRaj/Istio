// Copyright 2018 Istio Authors
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

package ingress

import (
	"flag"
	"testing"
)

var (
	tc = &testConfig{
		V1alpha1: true,  //implies envoyv1
		V1alpha3: false, //implies envoyv2
		Ingress:  true,
		Egress:   true,
	}
)

func init() {
	flag.BoolVar(&tc.V1alpha1, "v1alpha1", tc.V1alpha1, "Enable / disable v1alpha1 routing rules.")
	flag.BoolVar(&tc.V1alpha3, "v1alpha3", tc.V1alpha3, "Enable / disable v1alpha3 routing rules.")
	flag.BoolVar(&tc.Ingress, "ingress", tc.Ingress, "Enable / disable Ingress tests.")
	flag.BoolVar(&tc.Egress, "egress", tc.Egress, "Enable / disable Egress tests.")
}

func TestMain(m *testing.M) {
	flag.Parse()
	//cc, err := framework.NewCommonConfig("pilot_ingress_test")
	//if err != nil {
	//	os.Exit(-1)
	//}
	//tc.CommonConfig = cc
}

type testConfig struct {
	//	*framework.CommonConfig
	V1alpha1 bool
	V1alpha3 bool
	Ingress  bool
	Egress   bool
}
