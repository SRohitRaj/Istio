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

// QueryAsrPipelineList invokes the mts.QueryAsrPipelineList API synchronously
// api document: https://help.aliyun.com/api/mts/queryasrpipelinelist.html
func (client *Client) QueryAsrPipelineList(request *QueryAsrPipelineListRequest) (response *QueryAsrPipelineListResponse, err error) {
	response = CreateQueryAsrPipelineListResponse()
	err = client.DoAction(request, response)
	return
}

// QueryAsrPipelineListWithChan invokes the mts.QueryAsrPipelineList API asynchronously
// api document: https://help.aliyun.com/api/mts/queryasrpipelinelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryAsrPipelineListWithChan(request *QueryAsrPipelineListRequest) (<-chan *QueryAsrPipelineListResponse, <-chan error) {
	responseChan := make(chan *QueryAsrPipelineListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryAsrPipelineList(request)
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

// QueryAsrPipelineListWithCallback invokes the mts.QueryAsrPipelineList API asynchronously
// api document: https://help.aliyun.com/api/mts/queryasrpipelinelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryAsrPipelineListWithCallback(request *QueryAsrPipelineListRequest, callback func(response *QueryAsrPipelineListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryAsrPipelineListResponse
		var err error
		defer close(result)
		response, err = client.QueryAsrPipelineList(request)
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

// QueryAsrPipelineListRequest is the request struct for api QueryAsrPipelineList
type QueryAsrPipelineListRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PipelineIds          string           `position:"Query" name:"PipelineIds"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// QueryAsrPipelineListResponse is the response struct for api QueryAsrPipelineList
type QueryAsrPipelineListResponse struct {
	*responses.BaseResponse
	RequestId    string                             `json:"RequestId" xml:"RequestId"`
	NonExistIds  NonExistIdsInQueryAsrPipelineList  `json:"NonExistIds" xml:"NonExistIds"`
	PipelineList PipelineListInQueryAsrPipelineList `json:"PipelineList" xml:"PipelineList"`
}

// CreateQueryAsrPipelineListRequest creates a request to invoke QueryAsrPipelineList API
func CreateQueryAsrPipelineListRequest() (request *QueryAsrPipelineListRequest) {
	request = &QueryAsrPipelineListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "QueryAsrPipelineList", "mts", "openAPI")
	return
}

// CreateQueryAsrPipelineListResponse creates a response to parse from QueryAsrPipelineList response
func CreateQueryAsrPipelineListResponse() (response *QueryAsrPipelineListResponse) {
	response = &QueryAsrPipelineListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
