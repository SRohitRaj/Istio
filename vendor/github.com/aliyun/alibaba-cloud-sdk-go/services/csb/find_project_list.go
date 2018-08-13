package csb

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

// FindProjectList invokes the csb.FindProjectList API synchronously
// api document: https://help.aliyun.com/api/csb/findprojectlist.html
func (client *Client) FindProjectList(request *FindProjectListRequest) (response *FindProjectListResponse, err error) {
	response = CreateFindProjectListResponse()
	err = client.DoAction(request, response)
	return
}

// FindProjectListWithChan invokes the csb.FindProjectList API asynchronously
// api document: https://help.aliyun.com/api/csb/findprojectlist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) FindProjectListWithChan(request *FindProjectListRequest) (<-chan *FindProjectListResponse, <-chan error) {
	responseChan := make(chan *FindProjectListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.FindProjectList(request)
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

// FindProjectListWithCallback invokes the csb.FindProjectList API asynchronously
// api document: https://help.aliyun.com/api/csb/findprojectlist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) FindProjectListWithCallback(request *FindProjectListRequest, callback func(response *FindProjectListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *FindProjectListResponse
		var err error
		defer close(result)
		response, err = client.FindProjectList(request)
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

// FindProjectListRequest is the request struct for api FindProjectList
type FindProjectListRequest struct {
	*requests.RpcRequest
	CsbId       requests.Integer `position:"Query" name:"CsbId"`
	PageNum     requests.Integer `position:"Query" name:"PageNum"`
	ProjectName string           `position:"Query" name:"ProjectName"`
}

// FindProjectListResponse is the response struct for api FindProjectList
type FindProjectListResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateFindProjectListRequest creates a request to invoke FindProjectList API
func CreateFindProjectListRequest() (request *FindProjectListRequest) {
	request = &FindProjectListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CSB", "2017-11-18", "FindProjectList", "CSB", "openAPI")
	return
}

// CreateFindProjectListResponse creates a response to parse from FindProjectList response
func CreateFindProjectListResponse() (response *FindProjectListResponse) {
	response = &FindProjectListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
