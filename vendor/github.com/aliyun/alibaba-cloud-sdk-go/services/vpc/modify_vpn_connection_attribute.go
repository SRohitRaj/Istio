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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyVpnConnectionAttribute invokes the vpc.ModifyVpnConnectionAttribute API synchronously
// api document: https://help.aliyun.com/api/vpc/modifyvpnconnectionattribute.html
func (client *Client) ModifyVpnConnectionAttribute(request *ModifyVpnConnectionAttributeRequest) (response *ModifyVpnConnectionAttributeResponse, err error) {
	response = CreateModifyVpnConnectionAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyVpnConnectionAttributeWithChan invokes the vpc.ModifyVpnConnectionAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifyvpnconnectionattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyVpnConnectionAttributeWithChan(request *ModifyVpnConnectionAttributeRequest) (<-chan *ModifyVpnConnectionAttributeResponse, <-chan error) {
	responseChan := make(chan *ModifyVpnConnectionAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyVpnConnectionAttribute(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ModifyVpnConnectionAttributeWithCallback invokes the vpc.ModifyVpnConnectionAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifyvpnconnectionattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyVpnConnectionAttributeWithCallback(request *ModifyVpnConnectionAttributeRequest, callback func(response *ModifyVpnConnectionAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyVpnConnectionAttributeResponse
		var err error
		defer close(result)
		response, err = client.ModifyVpnConnectionAttribute(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ModifyVpnConnectionAttributeRequest is the request struct for api ModifyVpnConnectionAttribute
type ModifyVpnConnectionAttributeRequest struct {
	*requests.RpcRequest
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	VpnConnectionId      string           `position:"Query" name:"VpnConnectionId"`
	Name                 string           `position:"Query" name:"Name"`
	LocalSubnet          string           `position:"Query" name:"LocalSubnet"`
	RemoteSubnet         string           `position:"Query" name:"RemoteSubnet"`
	EffectImmediately    requests.Boolean `position:"Query" name:"EffectImmediately"`
	IkeConfig            string           `position:"Query" name:"IkeConfig"`
	IpsecConfig          string           `position:"Query" name:"IpsecConfig"`
}

// ModifyVpnConnectionAttributeResponse is the response struct for api ModifyVpnConnectionAttribute
type ModifyVpnConnectionAttributeResponse struct {
	*responses.BaseResponse
	RequestId         string      `json:"RequestId" xml:"RequestId"`
	VpnConnectionId   string      `json:"VpnConnectionId" xml:"VpnConnectionId"`
	CustomerGatewayId string      `json:"CustomerGatewayId" xml:"CustomerGatewayId"`
	VpnGatewayId      string      `json:"VpnGatewayId" xml:"VpnGatewayId"`
	Name              string      `json:"Name" xml:"Name"`
	Description       string      `json:"Description" xml:"Description"`
	LocalSubnet       string      `json:"LocalSubnet" xml:"LocalSubnet"`
	RemoteSubnet      string      `json:"RemoteSubnet" xml:"RemoteSubnet"`
	CreateTime        int         `json:"CreateTime" xml:"CreateTime"`
	EffectImmediately bool        `json:"EffectImmediately" xml:"EffectImmediately"`
	IkeConfig         IkeConfig   `json:"IkeConfig" xml:"IkeConfig"`
	IpsecConfig       IpsecConfig `json:"IpsecConfig" xml:"IpsecConfig"`
}

// CreateModifyVpnConnectionAttributeRequest creates a request to invoke ModifyVpnConnectionAttribute API
func CreateModifyVpnConnectionAttributeRequest() (request *ModifyVpnConnectionAttributeRequest) {
	request = &ModifyVpnConnectionAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ModifyVpnConnectionAttribute", "vpc", "openAPI")
	return
}

// CreateModifyVpnConnectionAttributeResponse creates a response to parse from ModifyVpnConnectionAttribute response
func CreateModifyVpnConnectionAttributeResponse() (response *ModifyVpnConnectionAttributeResponse) {
	response = &ModifyVpnConnectionAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
