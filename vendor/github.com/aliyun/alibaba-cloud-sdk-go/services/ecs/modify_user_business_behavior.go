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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyUserBusinessBehavior invokes the ecs.ModifyUserBusinessBehavior API synchronously
// api document: https://help.aliyun.com/api/ecs/modifyuserbusinessbehavior.html
func (client *Client) ModifyUserBusinessBehavior(request *ModifyUserBusinessBehaviorRequest) (response *ModifyUserBusinessBehaviorResponse, err error) {
	response = CreateModifyUserBusinessBehaviorResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyUserBusinessBehaviorWithChan invokes the ecs.ModifyUserBusinessBehavior API asynchronously
// api document: https://help.aliyun.com/api/ecs/modifyuserbusinessbehavior.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyUserBusinessBehaviorWithChan(request *ModifyUserBusinessBehaviorRequest) (<-chan *ModifyUserBusinessBehaviorResponse, <-chan error) {
	responseChan := make(chan *ModifyUserBusinessBehaviorResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyUserBusinessBehavior(request)
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

// ModifyUserBusinessBehaviorWithCallback invokes the ecs.ModifyUserBusinessBehavior API asynchronously
// api document: https://help.aliyun.com/api/ecs/modifyuserbusinessbehavior.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyUserBusinessBehaviorWithCallback(request *ModifyUserBusinessBehaviorRequest, callback func(response *ModifyUserBusinessBehaviorResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyUserBusinessBehaviorResponse
		var err error
		defer close(result)
		response, err = client.ModifyUserBusinessBehavior(request)
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

// ModifyUserBusinessBehaviorRequest is the request struct for api ModifyUserBusinessBehavior
type ModifyUserBusinessBehaviorRequest struct {
	*requests.RpcRequest
}

// ModifyUserBusinessBehaviorResponse is the response struct for api ModifyUserBusinessBehavior
type ModifyUserBusinessBehaviorResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyUserBusinessBehaviorRequest creates a request to invoke ModifyUserBusinessBehavior API
func CreateModifyUserBusinessBehaviorRequest() (request *ModifyUserBusinessBehaviorRequest) {
	request = &ModifyUserBusinessBehaviorRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "ModifyUserBusinessBehavior", "ecs", "openAPI")
	return
}

// CreateModifyUserBusinessBehaviorResponse creates a response to parse from ModifyUserBusinessBehavior response
func CreateModifyUserBusinessBehaviorResponse() (response *ModifyUserBusinessBehaviorResponse) {
	response = &ModifyUserBusinessBehaviorResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
