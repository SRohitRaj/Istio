package mts

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

// UpdateMediaPublishState invokes the mts.UpdateMediaPublishState API synchronously
// api document: https://help.aliyun.com/api/mts/updatemediapublishstate.html
func (client *Client) UpdateMediaPublishState(request *UpdateMediaPublishStateRequest) (response *UpdateMediaPublishStateResponse, err error) {
	response = CreateUpdateMediaPublishStateResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateMediaPublishStateWithChan invokes the mts.UpdateMediaPublishState API asynchronously
// api document: https://help.aliyun.com/api/mts/updatemediapublishstate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateMediaPublishStateWithChan(request *UpdateMediaPublishStateRequest) (<-chan *UpdateMediaPublishStateResponse, <-chan error) {
	responseChan := make(chan *UpdateMediaPublishStateResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateMediaPublishState(request)
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

// UpdateMediaPublishStateWithCallback invokes the mts.UpdateMediaPublishState API asynchronously
// api document: https://help.aliyun.com/api/mts/updatemediapublishstate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateMediaPublishStateWithCallback(request *UpdateMediaPublishStateRequest, callback func(response *UpdateMediaPublishStateResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateMediaPublishStateResponse
		var err error
		defer close(result)
		response, err = client.UpdateMediaPublishState(request)
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

// UpdateMediaPublishStateRequest is the request struct for api UpdateMediaPublishState
type UpdateMediaPublishStateRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	MediaId              string           `position:"Query" name:"MediaId"`
	Publish              requests.Boolean `position:"Query" name:"Publish"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// UpdateMediaPublishStateResponse is the response struct for api UpdateMediaPublishState
type UpdateMediaPublishStateResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpdateMediaPublishStateRequest creates a request to invoke UpdateMediaPublishState API
func CreateUpdateMediaPublishStateRequest() (request *UpdateMediaPublishStateRequest) {
	request = &UpdateMediaPublishStateRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "UpdateMediaPublishState", "mts", "openAPI")
	return
}

// CreateUpdateMediaPublishStateResponse creates a response to parse from UpdateMediaPublishState response
func CreateUpdateMediaPublishStateResponse() (response *UpdateMediaPublishStateResponse) {
	response = &UpdateMediaPublishStateResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
