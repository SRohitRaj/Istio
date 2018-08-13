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

// ModifyGlobalAccelerationInstanceSpec invokes the vpc.ModifyGlobalAccelerationInstanceSpec API synchronously
// api document: https://help.aliyun.com/api/vpc/modifyglobalaccelerationinstancespec.html
func (client *Client) ModifyGlobalAccelerationInstanceSpec(request *ModifyGlobalAccelerationInstanceSpecRequest) (response *ModifyGlobalAccelerationInstanceSpecResponse, err error) {
	response = CreateModifyGlobalAccelerationInstanceSpecResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyGlobalAccelerationInstanceSpecWithChan invokes the vpc.ModifyGlobalAccelerationInstanceSpec API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifyglobalaccelerationinstancespec.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyGlobalAccelerationInstanceSpecWithChan(request *ModifyGlobalAccelerationInstanceSpecRequest) (<-chan *ModifyGlobalAccelerationInstanceSpecResponse, <-chan error) {
	responseChan := make(chan *ModifyGlobalAccelerationInstanceSpecResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyGlobalAccelerationInstanceSpec(request)
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

// ModifyGlobalAccelerationInstanceSpecWithCallback invokes the vpc.ModifyGlobalAccelerationInstanceSpec API asynchronously
// api document: https://help.aliyun.com/api/vpc/modifyglobalaccelerationinstancespec.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyGlobalAccelerationInstanceSpecWithCallback(request *ModifyGlobalAccelerationInstanceSpecRequest, callback func(response *ModifyGlobalAccelerationInstanceSpecResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyGlobalAccelerationInstanceSpecResponse
		var err error
		defer close(result)
		response, err = client.ModifyGlobalAccelerationInstanceSpec(request)
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

// ModifyGlobalAccelerationInstanceSpecRequest is the request struct for api ModifyGlobalAccelerationInstanceSpec
type ModifyGlobalAccelerationInstanceSpecRequest struct {
	*requests.RpcRequest
	OwnerId                      requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount         string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId              requests.Integer `position:"Query" name:"ResourceOwnerId"`
	GlobalAccelerationInstanceId string           `position:"Query" name:"GlobalAccelerationInstanceId"`
	Bandwidth                    string           `position:"Query" name:"Bandwidth"`
	OwnerAccount                 string           `position:"Query" name:"OwnerAccount"`
}

// ModifyGlobalAccelerationInstanceSpecResponse is the response struct for api ModifyGlobalAccelerationInstanceSpec
type ModifyGlobalAccelerationInstanceSpecResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyGlobalAccelerationInstanceSpecRequest creates a request to invoke ModifyGlobalAccelerationInstanceSpec API
func CreateModifyGlobalAccelerationInstanceSpecRequest() (request *ModifyGlobalAccelerationInstanceSpecRequest) {
	request = &ModifyGlobalAccelerationInstanceSpecRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "ModifyGlobalAccelerationInstanceSpec", "vpc", "openAPI")
	return
}

// CreateModifyGlobalAccelerationInstanceSpecResponse creates a response to parse from ModifyGlobalAccelerationInstanceSpec response
func CreateModifyGlobalAccelerationInstanceSpecResponse() (response *ModifyGlobalAccelerationInstanceSpecResponse) {
	response = &ModifyGlobalAccelerationInstanceSpecResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
