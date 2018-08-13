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

// QueryWaterMarkTemplateList invokes the mts.QueryWaterMarkTemplateList API synchronously
// api document: https://help.aliyun.com/api/mts/querywatermarktemplatelist.html
func (client *Client) QueryWaterMarkTemplateList(request *QueryWaterMarkTemplateListRequest) (response *QueryWaterMarkTemplateListResponse, err error) {
	response = CreateQueryWaterMarkTemplateListResponse()
	err = client.DoAction(request, response)
	return
}

// QueryWaterMarkTemplateListWithChan invokes the mts.QueryWaterMarkTemplateList API asynchronously
// api document: https://help.aliyun.com/api/mts/querywatermarktemplatelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryWaterMarkTemplateListWithChan(request *QueryWaterMarkTemplateListRequest) (<-chan *QueryWaterMarkTemplateListResponse, <-chan error) {
	responseChan := make(chan *QueryWaterMarkTemplateListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryWaterMarkTemplateList(request)
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

// QueryWaterMarkTemplateListWithCallback invokes the mts.QueryWaterMarkTemplateList API asynchronously
// api document: https://help.aliyun.com/api/mts/querywatermarktemplatelist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryWaterMarkTemplateListWithCallback(request *QueryWaterMarkTemplateListRequest, callback func(response *QueryWaterMarkTemplateListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryWaterMarkTemplateListResponse
		var err error
		defer close(result)
		response, err = client.QueryWaterMarkTemplateList(request)
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

// QueryWaterMarkTemplateListRequest is the request struct for api QueryWaterMarkTemplateList
type QueryWaterMarkTemplateListRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	WaterMarkTemplateIds string           `position:"Query" name:"WaterMarkTemplateIds"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// QueryWaterMarkTemplateListResponse is the response struct for api QueryWaterMarkTemplateList
type QueryWaterMarkTemplateListResponse struct {
	*responses.BaseResponse
	RequestId             string                                            `json:"RequestId" xml:"RequestId"`
	NonExistWids          NonExistWids                                      `json:"NonExistWids" xml:"NonExistWids"`
	WaterMarkTemplateList WaterMarkTemplateListInQueryWaterMarkTemplateList `json:"WaterMarkTemplateList" xml:"WaterMarkTemplateList"`
}

// CreateQueryWaterMarkTemplateListRequest creates a request to invoke QueryWaterMarkTemplateList API
func CreateQueryWaterMarkTemplateListRequest() (request *QueryWaterMarkTemplateListRequest) {
	request = &QueryWaterMarkTemplateListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "QueryWaterMarkTemplateList", "mts", "openAPI")
	return
}

// CreateQueryWaterMarkTemplateListResponse creates a response to parse from QueryWaterMarkTemplateList response
func CreateQueryWaterMarkTemplateListResponse() (response *QueryWaterMarkTemplateListResponse) {
	response = &QueryWaterMarkTemplateListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
