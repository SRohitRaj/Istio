// Copyright 2020 Istio Authors
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

package caclient

import (
	"encoding/json"
	"os"
	"testing"
)

func TestKeyfactorConfigFromFile(t *testing.T) {

	testCases := map[string]struct {
		caName          string
		authToken       string
		caTemplate      string
		appKey          string
		customMetadatas []FieldAlias
		expectedErr     string
	}{
		"Missing caName in ENV": {
			caName:          "",
			caTemplate:      "Istio",
			authToken:       "Basic FakeAuthToken",
			appKey:          "FAKE_APP_KEY",
			customMetadatas: []FieldAlias{{Name: "Cluster", Alias: "Fake_Alias_Cluster"}},
			expectedErr:     "Missing caName (KEYFATOR_CA) in ENV",
		},
		"Missing appKey in ENV": {
			caName:          "FakeCAName",
			caTemplate:      "Istio",
			authToken:       "Basic FakeAuthToken",
			appKey:          "",
			customMetadatas: []FieldAlias{{Name: "Cluster", Alias: "Fake_Alias_Cluster"}},
			expectedErr:     "Missing appKey (KEYFATOR_APPKEY) in ENV",
		},
		"Missing authToken in ENV": {
			caName:          "FakeCAName",
			caTemplate:      "Istio",
			authToken:       "",
			appKey:          "FakeAppKey",
			customMetadatas: []FieldAlias{{Name: "Cluster", Alias: "Fake_Alias_Cluster"}},
			expectedErr:     "Missing authToken (KEYFATOR_AUTH_TOKEN) in ENV",
		},
		"Missing caTemplate in ENV": {
			caName:          "FakeCAName",
			caTemplate:      "",
			authToken:       "Basic FakeAuth",
			appKey:          "FakeAppKey",
			customMetadatas: []FieldAlias{{Name: "Cluster", Alias: "Fake_Alias_Cluster"}},
			expectedErr:     "Missing caTemplate (KEYFATOR_CA_TEMPLATE) in ENV",
		},
		"Do not supported new metadata field": {
			caName:          "FakeCAName",
			caTemplate:      "FakeCATemplate",
			authToken:       "Basic FakeAuth",
			appKey:          "FakeAppKey",
			customMetadatas: []FieldAlias{{Name: "ClusterInvalid", Alias: "Fake_Alias_Cluster"}},
			expectedErr:     "Do not support Metadata field name: ClusterInvalid",
		},
		"Invalid metadata configuration": {
			caName:          "FakeCAName",
			caTemplate:      "FakeCATemplate",
			authToken:       "Basic FakeAuth",
			appKey:          "FakeAppKey",
			customMetadatas: []FieldAlias{{Name: "Cluster", Alias: ""}},
			expectedErr:     "Invalid alias name for Metadata: Cluster",
		},
		"Empty metadata configuration": {
			caName:          "FakeCAName",
			caTemplate:      "FakeCATemplate",
			authToken:       "Basic FakeAuth",
			appKey:          "FakeAppKey",
			customMetadatas: nil,
			expectedErr:     "",
		},
	}

	for testID, tc := range testCases {
		t.Run(testID, func(tsub *testing.T) {
			os.Setenv("KEYFACTOR_CA", tc.caName)
			os.Setenv("KEYFACTOR_AUTH_TOKEN", tc.authToken)
			os.Setenv("KEYFACTOR_APPKEY", tc.appKey)
			os.Setenv("KEYFACTOR_CA_TEMPLATE", tc.caTemplate)

			metadataJSON, _ := json.Marshal(tc.customMetadatas)
			os.Setenv("KEYFACTOR_METADATA_JSON", string(metadataJSON))

			defer func() {
				os.Unsetenv("KEYFACTOR_CA")
				os.Unsetenv("KEYFACTOR_AUTH_TOKEN")
				os.Unsetenv("KEYFACTOR_APPKEY")
				os.Unsetenv("KEYFACTOR_CA_TEMPLATE")
				os.Unsetenv("KEYFACTOR_METADATA_JSON")
			}()

			_, err := LoadKeyfactorConfigFromENV()

			if err != nil && err.Error() != tc.expectedErr {
				tsub.Errorf("Failed testcase: %s - Expect (%v), but got (%v)", testID, tc.expectedErr, err)
			}
		})
	}
}
