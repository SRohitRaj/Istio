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

// ModifyEipAddressAttribute invokes the vpc.ModifyEipAddressAttribute API synchronously
// api document: https://help.aliyun.com/api/vpc/modifyeipaddressattribute.html
func (client *Client) ModifyEipAddressAttribute(request *ModifyEipAddressAttributeRequest) (response *ModifyEipAddressAttributeResponse, err error) {
	response = CreateModifyEipAddressAttributeResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyEipAddressAttributeWithChan invokes the vpc.ModifyEipAddressAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifyeipaddressattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyEipAddressAttributeWithChan(request *ModifyEipAddressAttributeRequest) (<-chan *ModifyEipAddressAttributeResponse, <-chan error) {
	responseChan := make(chan *ModifyEipAddressAttributeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyEipAddressAttribute(request)
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

// ModifyEipAddressAttributeWithCallback invokes the vpc.ModifyEipAddressAttribute API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifyeipaddressattribute.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyEipAddressAttributeWithCallback(request *ModifyEipAddressAttributeRequest, callback func(response *ModifyEipAddressAttributeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyEipAddressAttributeResponse
		var err error
		defer close(result)
		response, err = client.ModifyEipAddressAttribute(request)
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

// ModifyEipAddressAttributeRequest is the request struct for api ModifyEipAddressAttribute
type ModifyEipAddressAttributeRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	AllocationId         string           `position:"Query" name:"AllocationId"`
	Bandwidth            string           `position:"Query" name:"Bandwidth"`
	Name                 string           `position:"Query" name:"Name"`
	Description          string           `position:"Query" name:"Description"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// ModifyEipAddressAttributeResponse is the response struct for api ModifyEipAddressAttribute
type ModifyEipAddressAttributeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyEipAddressAttributeRequest creates a request to invoke ModifyEipAddressAttribute API
func CreateModifyEipAddressAttributeRequest() (request *ModifyEipAddressAttributeRequest) {
	request = &ModifyEipAddressAttributeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ModifyEipAddressAttribute", "vpc", "openAPI")
	return
}

// CreateModifyEipAddressAttributeResponse creates a response to parse from ModifyEipAddressAttribute response
func CreateModifyEipAddressAttributeResponse() (response *ModifyEipAddressAttributeResponse) {
	response = &ModifyEipAddressAttributeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
