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

// UpdateProjectListStatus invokes the csb.UpdateProjectListStatus API synchronously
// api document: https://help.aliyun.com/api/csb/updateprojectliststatus.html
func (client *Client) UpdateProjectListStatus(request *UpdateProjectListStatusRequest) (response *UpdateProjectListStatusResponse, err error) {
	response = CreateUpdateProjectListStatusResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateProjectListStatusWithChan invokes the csb.UpdateProjectListStatus API asynchronously
// api document: https://help.aliyun.com/api/csb/updateprojectliststatus.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateProjectListStatusWithChan(request *UpdateProjectListStatusRequest) (<-chan *UpdateProjectListStatusResponse, <-chan error) {
	responseChan := make(chan *UpdateProjectListStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateProjectListStatus(request)
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

// UpdateProjectListStatusWithCallback invokes the csb.UpdateProjectListStatus API asynchronously
// api document: https://help.aliyun.com/api/csb/updateprojectliststatus.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateProjectListStatusWithCallback(request *UpdateProjectListStatusRequest, callback func(response *UpdateProjectListStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateProjectListStatusResponse
		var err error
		defer close(result)
		response, err = client.UpdateProjectListStatus(request)
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

// UpdateProjectListStatusRequest is the request struct for api UpdateProjectListStatus
type UpdateProjectListStatusRequest struct {
	*requests.RpcRequest
	CsbId requests.Integer `position:"Query" name:"CsbId"`
	Data  string           `position:"Body" name:"Data"`
}

// UpdateProjectListStatusResponse is the response struct for api UpdateProjectListStatus
type UpdateProjectListStatusResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpdateProjectListStatusRequest creates a request to invoke UpdateProjectListStatus API
func CreateUpdateProjectListStatusRequest() (request *UpdateProjectListStatusRequest) {
	request = &UpdateProjectListStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CSB", "2017-11-18", "UpdateProjectListStatus", "CSB", "openAPI")
	return
}

// CreateUpdateProjectListStatusResponse creates a response to parse from UpdateProjectListStatus response
func CreateUpdateProjectListStatusResponse() (response *UpdateProjectListStatusResponse) {
	response = &UpdateProjectListStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
