package cloudapi

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

// DescribeApiIpControls invokes the cloudapi.DescribeApiIpControls API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapiipcontrols.html
func (client *Client) DescribeApiIpControls(request *DescribeApiIpControlsRequest) (response *DescribeApiIpControlsResponse, err error) {
	response = CreateDescribeApiIpControlsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeApiIpControlsWithChan invokes the cloudapi.DescribeApiIpControls API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapiipcontrols.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiIpControlsWithChan(request *DescribeApiIpControlsRequest) (<-chan *DescribeApiIpControlsResponse, <-chan error) {
	responseChan := make(chan *DescribeApiIpControlsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeApiIpControls(request)
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

// DescribeApiIpControlsWithCallback invokes the cloudapi.DescribeApiIpControls API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapiipcontrols.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiIpControlsWithCallback(request *DescribeApiIpControlsRequest, callback func(response *DescribeApiIpControlsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeApiIpControlsResponse
		var err error
		defer close(result)
		response, err = client.DescribeApiIpControls(request)
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

// DescribeApiIpControlsRequest is the request struct for api DescribeApiIpControls
type DescribeApiIpControlsRequest struct {
	*requests.RpcRequest
	StageName  string           `position:"Query" name:"StageName"`
	GroupId    string           `position:"Query" name:"GroupId"`
	ApiIds     string           `position:"Query" name:"ApiIds"`
	PageNumber requests.Integer `position:"Query" name:"PageNumber"`
	PageSize   requests.Integer `position:"Query" name:"PageSize"`
}

// DescribeApiIpControlsResponse is the response struct for api DescribeApiIpControls
type DescribeApiIpControlsResponse struct {
	*responses.BaseResponse
	RequestId     string        `json:"RequestId" xml:"RequestId"`
	TotalCount    int           `json:"TotalCount" xml:"TotalCount"`
	PageSize      int           `json:"PageSize" xml:"PageSize"`
	PageNumber    int           `json:"PageNumber" xml:"PageNumber"`
	ApiIpControls ApiIpControls `json:"ApiIpControls" xml:"ApiIpControls"`
}

// CreateDescribeApiIpControlsRequest creates a request to invoke DescribeApiIpControls API
func CreateDescribeApiIpControlsRequest() (request *DescribeApiIpControlsRequest) {
	request = &DescribeApiIpControlsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeApiIpControls", "apigateway", "openAPI")
	return
}

// CreateDescribeApiIpControlsResponse creates a response to parse from DescribeApiIpControls response
func CreateDescribeApiIpControlsResponse() (response *DescribeApiIpControlsResponse) {
	response = &DescribeApiIpControlsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
