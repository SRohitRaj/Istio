package cbn

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

// CenBandwidthPackage is a nested struct in cbn response
type CenBandwidthPackage struct {
	CenBandwidthPackageId      string `json:"CenBandwidthPackageId" xml:"CenBandwidthPackageId"`
	Name                       string `json:"Name" xml:"Name"`
	Description                string `json:"Description" xml:"Description"`
	Bandwidth                  int    `json:"Bandwidth" xml:"Bandwidth"`
	BandwidthPackageChargeType string `json:"BandwidthPackageChargeType" xml:"BandwidthPackageChargeType"`
	GeographicRegionAId        string `json:"GeographicRegionAId" xml:"GeographicRegionAId"`
	GeographicRegionBId        string `json:"GeographicRegionBId" xml:"GeographicRegionBId"`
	BusinessStatus             string `json:"BusinessStatus" xml:"BusinessStatus"`
	CreationTime               string `json:"CreationTime" xml:"CreationTime"`
	ExpiredTime                string `json:"ExpiredTime" xml:"ExpiredTime"`
	Status                     string `json:"Status" xml:"Status"`
	CenIds                     CenIds `json:"CenIds" xml:"CenIds"`
}
