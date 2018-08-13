package cms

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

// Resource is a nested struct in cms response
type Resource struct {
	Category      string                      `json:"Category" xml:"Category"`
	Desc          string                      `json:"Desc" xml:"Desc"`
	InstanceName  string                      `json:"InstanceName" xml:"InstanceName"`
	GroupName     string                      `json:"GroupName" xml:"GroupName"`
	Id            int                         `json:"Id" xml:"Id"`
	BindUrls      string                      `json:"BindUrls" xml:"BindUrls"`
	ServiceId     string                      `json:"ServiceId" xml:"ServiceId"`
	RegionId      string                      `json:"RegionId" xml:"RegionId"`
	InstanceId    string                      `json:"InstanceId" xml:"InstanceId"`
	GroupId       int                         `json:"GroupId" xml:"GroupId"`
	AliUid        int                         `json:"AliUid" xml:"AliUid"`
	NetworkType   string                      `json:"NetworkType" xml:"NetworkType"`
	Type          string                      `json:"Type" xml:"Type"`
	Vpc           Vpc                         `json:"Vpc" xml:"Vpc"`
	Region        Region                      `json:"Region" xml:"Region"`
	Tags          Tags                        `json:"Tags" xml:"Tags"`
	ContactGroups ContactGroupsInListMyGroups `json:"ContactGroups" xml:"ContactGroups"`
}
