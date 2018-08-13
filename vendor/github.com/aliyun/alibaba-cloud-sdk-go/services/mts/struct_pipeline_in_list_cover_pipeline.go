package mts

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// PipelineInListCoverPipeline is a nested struct in mts response
type PipelineInListCoverPipeline struct {
	UserId       int    `json:"UserId" xml:"UserId"`
	PipelineId   string `json:"PipelineId" xml:"PipelineId"`
	Name         string `json:"Name" xml:"Name"`
	State        string `json:"State" xml:"State"`
	Priority     string `json:"Priority" xml:"Priority"`
	QuotaNum     int    `json:"quotaNum" xml:"quotaNum"`
	QuotaUsed    int    `json:"quotaUsed" xml:"quotaUsed"`
	NotifyConfig string `json:"NotifyConfig" xml:"NotifyConfig"`
	Role         string `json:"Role" xml:"Role"`
	ExtendConfig string `json:"ExtendConfig" xml:"ExtendConfig"`
}
