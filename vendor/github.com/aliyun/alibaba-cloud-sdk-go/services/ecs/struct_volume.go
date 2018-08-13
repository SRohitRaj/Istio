package ecs

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

// Volume is a nested struct in ecs response
type Volume struct {
	VolumeId                      string                          `json:"VolumeId" xml:"VolumeId"`
	RegionId                      string                          `json:"RegionId" xml:"RegionId"`
	ZoneId                        string                          `json:"ZoneId" xml:"ZoneId"`
	VolumeName                    string                          `json:"VolumeName" xml:"VolumeName"`
	Description                   string                          `json:"Description" xml:"Description"`
	Category                      string                          `json:"Category" xml:"Category"`
	Size                          int                             `json:"Size" xml:"Size"`
	SourceSnapshotId              string                          `json:"SourceSnapshotId" xml:"SourceSnapshotId"`
	AutoSnapshotPolicyId          string                          `json:"AutoSnapshotPolicyId" xml:"AutoSnapshotPolicyId"`
	SnapshotLinkId                string                          `json:"SnapshotLinkId" xml:"SnapshotLinkId"`
	Status                        string                          `json:"Status" xml:"Status"`
	EnableAutomatedSnapshotPolicy bool                            `json:"EnableAutomatedSnapshotPolicy" xml:"EnableAutomatedSnapshotPolicy"`
	CreationTime                  string                          `json:"CreationTime" xml:"CreationTime"`
	VolumeChargeType              string                          `json:"VolumeChargeType" xml:"VolumeChargeType"`
	MountInstanceNum              int                             `json:"MountInstanceNum" xml:"MountInstanceNum"`
	Encrypted                     bool                            `json:"Encrypted" xml:"Encrypted"`
	OperationLocks                OperationLocksInDescribeVolumes `json:"OperationLocks" xml:"OperationLocks"`
	MountInstances                MountInstancesInDescribeVolumes `json:"MountInstances" xml:"MountInstances"`
	Tags                          TagsInDescribeVolumes           `json:"Tags" xml:"Tags"`
}
