package cms

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

// TaskConfigDelete invokes the cms.TaskConfigDelete API synchronously
// api document: https://help.aliyun.com/api/cms/taskconfigdelete.html
func (client *Client) TaskConfigDelete(request *TaskConfigDeleteRequest) (response *TaskConfigDeleteResponse, err error) {
	response = CreateTaskConfigDeleteResponse()
	err = client.DoAction(request, response)
	return
}

// TaskConfigDeleteWithChan invokes the cms.TaskConfigDelete API asynchronously
// api document: https://help.aliyun.com/api/cms/taskconfigdelete.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TaskConfigDeleteWithChan(request *TaskConfigDeleteRequest) (<-chan *TaskConfigDeleteResponse, <-chan error) {
	responseChan := make(chan *TaskConfigDeleteResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.TaskConfigDelete(request)
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

// TaskConfigDeleteWithCallback invokes the cms.TaskConfigDelete API asynchronously
// api document: https://help.aliyun.com/api/cms/taskconfigdelete.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TaskConfigDeleteWithCallback(request *TaskConfigDeleteRequest, callback func(response *TaskConfigDeleteResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *TaskConfigDeleteResponse
		var err error
		defer close(result)
		response, err = client.TaskConfigDelete(request)
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

// TaskConfigDeleteRequest is the request struct for api TaskConfigDelete
type TaskConfigDeleteRequest struct {
	*requests.RpcRequest
	IdList *[]string `position:"Query" name:"IdList"  type:"Repeated"`
}

// TaskConfigDeleteResponse is the response struct for api TaskConfigDelete
type TaskConfigDeleteResponse struct {
	*responses.BaseResponse
	ErrorCode    int    `json:"ErrorCode" xml:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage" xml:"ErrorMessage"`
	Success      bool   `json:"Success" xml:"Success"`
	RequestId    string `json:"RequestId" xml:"RequestId"`
}

// CreateTaskConfigDeleteRequest creates a request to invoke TaskConfigDelete API
func CreateTaskConfigDeleteRequest() (request *TaskConfigDeleteRequest) {
	request = &TaskConfigDeleteRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2018-03-08", "TaskConfigDelete", "cms", "openAPI")
	return
}

// CreateTaskConfigDeleteResponse creates a response to parse from TaskConfigDelete response
func CreateTaskConfigDeleteResponse() (response *TaskConfigDeleteResponse) {
	response = &TaskConfigDeleteResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
