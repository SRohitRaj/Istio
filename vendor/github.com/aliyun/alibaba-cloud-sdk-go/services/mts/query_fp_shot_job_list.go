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

// QueryFpShotJobList invokes the mts.QueryFpShotJobList API synchronously
// api document: https://help.aliyun.com/api/mts/queryfpshotjoblist.html
func (client *Client) QueryFpShotJobList(request *QueryFpShotJobListRequest) (response *QueryFpShotJobListResponse, err error) {
	response = CreateQueryFpShotJobListResponse()
	err = client.DoAction(request, response)
	return
}

// QueryFpShotJobListWithChan invokes the mts.QueryFpShotJobList API asynchronously
// api document: https://help.aliyun.com/api/mts/queryfpshotjoblist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryFpShotJobListWithChan(request *QueryFpShotJobListRequest) (<-chan *QueryFpShotJobListResponse, <-chan error) {
	responseChan := make(chan *QueryFpShotJobListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryFpShotJobList(request)
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

// QueryFpShotJobListWithCallback invokes the mts.QueryFpShotJobList API asynchronously
// api document: https://help.aliyun.com/api/mts/queryfpshotjoblist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryFpShotJobListWithCallback(request *QueryFpShotJobListRequest, callback func(response *QueryFpShotJobListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryFpShotJobListResponse
		var err error
		defer close(result)
		response, err = client.QueryFpShotJobList(request)
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

// QueryFpShotJobListRequest is the request struct for api QueryFpShotJobList
type QueryFpShotJobListRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	JobIds               string           `position:"Query" name:"JobIds"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// QueryFpShotJobListResponse is the response struct for api QueryFpShotJobList
type QueryFpShotJobListResponse struct {
	*responses.BaseResponse
	RequestId     string                          `json:"RequestId" xml:"RequestId"`
	NonExistIds   NonExistIdsInQueryFpShotJobList `json:"NonExistIds" xml:"NonExistIds"`
	FpShotJobList FpShotJobList                   `json:"FpShotJobList" xml:"FpShotJobList"`
}

// CreateQueryFpShotJobListRequest creates a request to invoke QueryFpShotJobList API
func CreateQueryFpShotJobListRequest() (request *QueryFpShotJobListRequest) {
	request = &QueryFpShotJobListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "QueryFpShotJobList", "mts", "openAPI")
	return
}

// CreateQueryFpShotJobListResponse creates a response to parse from QueryFpShotJobList response
func CreateQueryFpShotJobListResponse() (response *QueryFpShotJobListResponse) {
	response = &QueryFpShotJobListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
