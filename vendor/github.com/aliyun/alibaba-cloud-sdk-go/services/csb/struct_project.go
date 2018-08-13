package csb

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

// Project is a nested struct in csb response
type Project struct {
	InterfaceJarLocation string `json:"InterfaceJarLocation" xml:"InterfaceJarLocation"`
	Status               int    `json:"Status" xml:"Status"`
	JarFileKey           string `json:"JarFileKey" xml:"JarFileKey"`
	ProjectOwnerEmail    string `json:"ProjectOwnerEmail" xml:"ProjectOwnerEmail"`
	Id                   int    `json:"Id" xml:"Id"`
	ProjectOwnerName     string `json:"ProjectOwnerName" xml:"ProjectOwnerName"`
	ProjectOwnerPhoneNum string `json:"ProjectOwnerPhoneNum" xml:"ProjectOwnerPhoneNum"`
	GmtCreate            int    `json:"GmtCreate" xml:"GmtCreate"`
	InterfaceJarName     string `json:"InterfaceJarName" xml:"InterfaceJarName"`
	DeleteFlag           int    `json:"DeleteFlag" xml:"DeleteFlag"`
	OwnerId              string `json:"OwnerId" xml:"OwnerId"`
	ProjectName          string `json:"ProjectName" xml:"ProjectName"`
	UserId               string `json:"UserId" xml:"UserId"`
	ApiNum               int    `json:"ApiNum" xml:"ApiNum"`
	GmtModified          int    `json:"GmtModified" xml:"GmtModified"`
	CsbId                int    `json:"CsbId" xml:"CsbId"`
	Description          string `json:"Description" xml:"Description"`
}
