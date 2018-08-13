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

// DescribeIpControlPolicyItems invokes the cloudapi.DescribeIpControlPolicyItems API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describeipcontrolpolicyitems.html
func (client *Client) DescribeIpControlPolicyItems(request *DescribeIpControlPolicyItemsRequest) (response *DescribeIpControlPolicyItemsResponse, err error) {
	response = CreateDescribeIpControlPolicyItemsResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeIpControlPolicyItemsWithChan invokes the cloudapi.DescribeIpControlPolicyItems API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeipcontrolpolicyitems.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeIpControlPolicyItemsWithChan(request *DescribeIpControlPolicyItemsRequest) (<-chan *DescribeIpControlPolicyItemsResponse, <-chan error) {
	responseChan := make(chan *DescribeIpControlPolicyItemsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeIpControlPolicyItems(request)
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

// DescribeIpControlPolicyItemsWithCallback invokes the cloudapi.DescribeIpControlPolicyItems API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeipcontrolpolicyitems.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeIpControlPolicyItemsWithCallback(request *DescribeIpControlPolicyItemsRequest, callback func(response *DescribeIpControlPolicyItemsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeIpControlPolicyItemsResponse
		var err error
		defer close(result)
		response, err = client.DescribeIpControlPolicyItems(request)
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

// DescribeIpControlPolicyItemsRequest is the request struct for api DescribeIpControlPolicyItems
type DescribeIpControlPolicyItemsRequest struct {
	*requests.RpcRequest
	IpControlId  string           `position:"Query" name:"IpControlId"`
	PolicyItemId string           `position:"Query" name:"PolicyItemId"`
	PageNumber   requests.Integer `position:"Query" name:"PageNumber"`
	PageSize     requests.Integer `position:"Query" name:"PageSize"`
}

// DescribeIpControlPolicyItemsResponse is the response struct for api DescribeIpControlPolicyItems
type DescribeIpControlPolicyItemsResponse struct {
	*responses.BaseResponse
	RequestId            string               `json:"RequestId" xml:"RequestId"`
	TotalCount           int                  `json:"TotalCount" xml:"TotalCount"`
	PageSize             int                  `json:"PageSize" xml:"PageSize"`
	PageNumber           int                  `json:"PageNumber" xml:"PageNumber"`
	IpControlPolicyItems IpControlPolicyItems `json:"IpControlPolicyItems" xml:"IpControlPolicyItems"`
}

// CreateDescribeIpControlPolicyItemsRequest creates a request to invoke DescribeIpControlPolicyItems API
func CreateDescribeIpControlPolicyItemsRequest() (request *DescribeIpControlPolicyItemsRequest) {
	request = &DescribeIpControlPolicyItemsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeIpControlPolicyItems", "apigateway", "openAPI")
	return
}

// CreateDescribeIpControlPolicyItemsResponse creates a response to parse from DescribeIpControlPolicyItems response
func CreateDescribeIpControlPolicyItemsResponse() (response *DescribeIpControlPolicyItemsResponse) {
	response = &DescribeIpControlPolicyItemsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
