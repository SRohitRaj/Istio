package vpc

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

// VirtualBorderRouterType is a nested struct in vpc response
type VirtualBorderRouterType struct {
	VbrId                            string                                       `json:"VbrId" xml:"VbrId"`
	CreationTime                     string                                       `json:"CreationTime" xml:"CreationTime"`
	ActivationTime                   string                                       `json:"ActivationTime" xml:"ActivationTime"`
	TerminationTime                  string                                       `json:"TerminationTime" xml:"TerminationTime"`
	RecoveryTime                     string                                       `json:"RecoveryTime" xml:"RecoveryTime"`
	Status                           string                                       `json:"Status" xml:"Status"`
	VlanId                           int                                          `json:"VlanId" xml:"VlanId"`
	CircuitCode                      string                                       `json:"CircuitCode" xml:"CircuitCode"`
	RouteTableId                     string                                       `json:"RouteTableId" xml:"RouteTableId"`
	VlanInterfaceId                  string                                       `json:"VlanInterfaceId" xml:"VlanInterfaceId"`
	LocalGatewayIp                   string                                       `json:"LocalGatewayIp" xml:"LocalGatewayIp"`
	PeerGatewayIp                    string                                       `json:"PeerGatewayIp" xml:"PeerGatewayIp"`
	PeeringSubnetMask                string                                       `json:"PeeringSubnetMask" xml:"PeeringSubnetMask"`
	PhysicalConnectionId             string                                       `json:"PhysicalConnectionId" xml:"PhysicalConnectionId"`
	PhysicalConnectionStatus         string                                       `json:"PhysicalConnectionStatus" xml:"PhysicalConnectionStatus"`
	PhysicalConnectionBusinessStatus string                                       `json:"PhysicalConnectionBusinessStatus" xml:"PhysicalConnectionBusinessStatus"`
	PhysicalConnectionOwnerUid       string                                       `json:"PhysicalConnectionOwnerUid" xml:"PhysicalConnectionOwnerUid"`
	AccessPointId                    string                                       `json:"AccessPointId" xml:"AccessPointId"`
	Name                             string                                       `json:"Name" xml:"Name"`
	Description                      string                                       `json:"Description" xml:"Description"`
	AssociatedPhysicalConnections    AssociatedPhysicalConnections                `json:"AssociatedPhysicalConnections" xml:"AssociatedPhysicalConnections"`
	AssociatedCens                   AssociatedCensInDescribeVirtualBorderRouters `json:"AssociatedCens" xml:"AssociatedCens"`
}
