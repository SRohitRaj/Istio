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

// FindServiceList invokes the csb.FindServiceList API synchronously
// api document: https://help.aliyun.com/api/csb/findservicelist.html
func (client *Client) FindServiceList(request *FindServiceListRequest) (response *FindServiceListResponse, err error) {
	response = CreateFindServiceListResponse()
	err = client.DoAction(request, response)
	return
}

// FindServiceListWithChan invokes the csb.FindServiceList API asynchronously
// api document: https://help.aliyun.com/api/csb/findservicelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) FindServiceListWithChan(request *FindServiceListRequest) (<-chan *FindServiceListResponse, <-chan error) {
	responseChan := make(chan *FindServiceListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.FindServiceList(request)
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

// FindServiceListWithCallback invokes the csb.FindServiceList API asynchronously
// api document: https://help.aliyun.com/api/csb/findservicelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) FindServiceListWithCallback(request *FindServiceListRequest, callback func(response *FindServiceListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *FindServiceListResponse
		var err error
		defer close(result)
		response, err = client.FindServiceList(request)
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

// FindServiceListRequest is the request struct for api FindServiceList
type FindServiceListRequest struct {
	*requests.RpcRequest
	PageNum        requests.Integer `position:"Query" name:"PageNum"`
	CasShowType    requests.Integer `position:"Query" name:"CasShowType"`
	ServiceName    string           `position:"Query" name:"ServiceName"`
	Alias          string           `position:"Query" name:"Alias"`
	ProjectName    string           `position:"Query" name:"ProjectName"`
	CsbId          requests.Integer `position:"Query" name:"CsbId"`
	ShowDelService requests.Boolean `position:"Query" name:"ShowDelService"`
}

// FindServiceListResponse is the response struct for api FindServiceList
type FindServiceListResponse struct {
	*responses.BaseResponse
	Message   string `json:"Message" xml:"Message"`
	Code      int    `json:"Code" xml:"Code"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateFindServiceListRequest creates a request to invoke FindServiceList API
func CreateFindServiceListRequest() (request *FindServiceListRequest) {
	request = &FindServiceListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CSB", "2017-11-18", "FindServiceList", "CSB", "openAPI")
	return
}

// CreateFindServiceListResponse creates a response to parse from FindServiceList response
func CreateFindServiceListResponse() (response *FindServiceListResponse) {
	response = &FindServiceListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
