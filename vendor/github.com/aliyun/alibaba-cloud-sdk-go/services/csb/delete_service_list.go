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

// DeleteServiceList invokes the csb.DeleteServiceList API synchronously
// api document: https://help.aliyun.com/api/csb/deleteservicelist.html
func (client *Client) DeleteServiceList(request *DeleteServiceListRequest) (response *DeleteServiceListResponse, err error) {
	response = CreateDeleteServiceListResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteServiceListWithChan invokes the csb.DeleteServiceList API asynchronously
// api document: https://help.aliyun.com/api/csb/deleteservicelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteServiceListWithChan(request *DeleteServiceListRequest) (<-chan *DeleteServiceListResponse, <-chan error) {
	responseChan := make(chan *DeleteServiceListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteServiceList(request)
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

// DeleteServiceListWithCallback invokes the csb.DeleteServiceList API asynchronously
// api document: https://help.aliyun.com/api/csb/deleteservicelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteServiceListWithCallback(request *DeleteServiceListRequest, callback func(response *DeleteServiceListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteServiceListResponse
		var err error
		defer close(result)
		response, err = client.DeleteServiceList(request)
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

// DeleteServiceListRequest is the request struct for api DeleteServiceList
type DeleteServiceListRequest struct {
	*requests.RpcRequest
	CsbId requests.Integer `position:"Query" name:"CsbId"`
	Data  string           `position:"Body" name:"Data"`
}

// DeleteServiceListResponse is the response struct for api DeleteServiceList
type DeleteServiceListResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteServiceListRequest creates a request to invoke DeleteServiceList API
func CreateDeleteServiceListRequest() (request *DeleteServiceListRequest) {
	request = &DeleteServiceListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CSB", "2017-11-18", "DeleteServiceList", "CSB", "openAPI")
	return
}

// CreateDeleteServiceListResponse creates a response to parse from DeleteServiceList response
func CreateDeleteServiceListResponse() (response *DeleteServiceListResponse) {
	response = &DeleteServiceListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
