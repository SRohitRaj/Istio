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

package distribution

import (
	"gopkg.in/yaml.v2"
)

type DistributionReport struct {
	Reporter            string         `json:"reporter"`
	DataPlaneCount      int            `json:"dataPlaneCount"`
	InProgressResources map[string]int `json:"inProgressResources"`
}

func ReportFromYaml(content []byte) (DistributionReport, error) {
	out := DistributionReport{}
	err := yaml.Unmarshal(content, &out)
	return out, err
}

