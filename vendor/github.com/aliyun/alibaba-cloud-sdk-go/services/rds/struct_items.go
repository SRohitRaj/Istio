package rds

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

// Items is a nested struct in rds response
type Items struct {
	ReplicaMode        string  `json:"ReplicaMode" xml:"ReplicaMode"`
	Role               string  `json:"Role" xml:"Role"`
	ReadWriteType      string  `json:"ReadWriteType" xml:"ReadWriteType"`
	DomainMode         string  `json:"DomainMode" xml:"DomainMode"`
	ReplicaDescription string  `json:"ReplicaDescription" xml:"ReplicaDescription"`
	DBInstanceId       string  `json:"DBInstanceId" xml:"DBInstanceId"`
	ReplicaStatus      string  `json:"ReplicaStatus" xml:"ReplicaStatus"`
	ReplicaId          string  `json:"ReplicaId" xml:"ReplicaId"`
	DBInstances        []Items `json:"DBInstances" xml:"DBInstances"`
}
