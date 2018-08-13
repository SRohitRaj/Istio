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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DeleteSkillGroup invokes the ccc.DeleteSkillGroup API synchronously
// api document: https://help.aliyun.com/api/ccc/deleteskillgroup.html
func (client *Client) DeleteSkillGroup(request *DeleteSkillGroupRequest) (response *DeleteSkillGroupResponse, err error) {
	response = CreateDeleteSkillGroupResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteSkillGroupWithChan invokes the ccc.DeleteSkillGroup API asynchronously
// api document: https://help.aliyun.com/api/ccc/deleteskillgroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteSkillGroupWithChan(request *DeleteSkillGroupRequest) (<-chan *DeleteSkillGroupResponse, <-chan error) {
	responseChan := make(chan *DeleteSkillGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteSkillGroup(request)
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

// DeleteSkillGroupWithCallback invokes the ccc.DeleteSkillGroup API asynchronously
// api document: https://help.aliyun.com/api/ccc/deleteskillgroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteSkillGroupWithCallback(request *DeleteSkillGroupRequest, callback func(response *DeleteSkillGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteSkillGroupResponse
		var err error
		defer close(result)
		response, err = client.DeleteSkillGroup(request)
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

// DeleteSkillGroupRequest is the request struct for api DeleteSkillGroup
type DeleteSkillGroupRequest struct {
	*requests.RpcRequest
	InstanceId   string `position:"Query" name:"InstanceId"`
	SkillGroupId string `position:"Query" name:"SkillGroupId"`
}

// DeleteSkillGroupResponse is the response struct for api DeleteSkillGroup
type DeleteSkillGroupResponse struct {
	*responses.BaseResponse
	RequestId      string `json:"RequestId" xml:"RequestId"`
	Success        bool   `json:"Success" xml:"Success"`
	Code           string `json:"Code" xml:"Code"`
	Message        string `json:"Message" xml:"Message"`
	HttpStatusCode int    `json:"HttpStatusCode" xml:"HttpStatusCode"`
}

// CreateDeleteSkillGroupRequest creates a request to invoke DeleteSkillGroup API
func CreateDeleteSkillGroupRequest() (request *DeleteSkillGroupRequest) {
	request = &DeleteSkillGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CCC", "2017-07-05", "DeleteSkillGroup", "", "")
	return
}

// CreateDeleteSkillGroupResponse creates a response to parse from DeleteSkillGroup response
func CreateDeleteSkillGroupResponse() (response *DeleteSkillGroupResponse) {
	response = &DeleteSkillGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
