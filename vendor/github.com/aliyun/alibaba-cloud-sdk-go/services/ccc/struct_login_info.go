package ccc

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

// LoginInfo is a nested struct in ccc response
type LoginInfo struct {
	UserName       string `json:"UserName" xml:"UserName"`
	DisplayName    string `json:"DisplayName" xml:"DisplayName"`
	PhoneNumber    string `json:"PhoneNumber" xml:"PhoneNumber"`
	Region         string `json:"Region" xml:"Region"`
	WebRtcUrl      string `json:"WebRtcUrl" xml:"WebRtcUrl"`
	AgentServerUrl string `json:"AgentServerUrl" xml:"AgentServerUrl"`
	Extension      string `json:"Extension" xml:"Extension"`
	TenantId       string `json:"TenantId" xml:"TenantId"`
	Signature      string `json:"Signature" xml:"Signature"`
	SignData       string `json:"SignData" xml:"SignData"`
}
