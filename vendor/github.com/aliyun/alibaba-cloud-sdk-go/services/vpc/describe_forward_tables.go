package vpc

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

// DescribeForwardTables invokes the vpc.DescribeForwardTables API synchronously
// api document: https://help.aliyun.com/api/vpc/describeforwardtables.html
func (client *Client) DescribeForwardTables(request *DescribeForwardTablesRequest) (response *DescribeForwardTablesResponse, err error) {
	response = CreateDescribeForwardTablesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeForwardTablesWithChan invokes the vpc.DescribeForwardTables API asynchronously
// api document: https://help.aliyun.com/api/vpc/describeforwardtables.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeForwardTablesWithChan(request *DescribeForwardTablesRequest) (<-chan *DescribeForwardTablesResponse, <-chan error) {
	responseChan := make(chan *DescribeForwardTablesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeForwardTables(request)
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

// DescribeForwardTablesWithCallback invokes the vpc.DescribeForwardTables API asynchronously
// api document: https://help.aliyun.com/api/vpc/describeforwardtables.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeForwardTablesWithCallback(request *DescribeForwardTablesRequest, callback func(response *DescribeForwardTablesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeForwardTablesResponse
		var err error
		defer close(result)
		response, err = client.DescribeForwardTables(request)
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

// DescribeForwardTablesRequest is the request struct for api DescribeForwardTables
type DescribeForwardTablesRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	ForwardTableId       string           `position:"Query" name:"ForwardTableId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
}

// DescribeForwardTablesResponse is the response struct for api DescribeForwardTables
type DescribeForwardTablesResponse struct {
	*responses.BaseResponse
	RequestId     string        `json:"RequestId" xml:"RequestId"`
	TotalCount    int           `json:"TotalCount" xml:"TotalCount"`
	PageNumber    int           `json:"PageNumber" xml:"PageNumber"`
	PageSize      int           `json:"PageSize" xml:"PageSize"`
	ForwardTables ForwardTables `json:"ForwardTables" xml:"ForwardTables"`
}

// CreateDescribeForwardTablesRequest creates a request to invoke DescribeForwardTables API
func CreateDescribeForwardTablesRequest() (request *DescribeForwardTablesRequest) {
	request = &DescribeForwardTablesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "DescribeForwardTables", "vpc", "openAPI")
	return
}

// CreateDescribeForwardTablesResponse creates a response to parse from DescribeForwardTables response
func CreateDescribeForwardTablesResponse() (response *DescribeForwardTablesResponse) {
	response = &DescribeForwardTablesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
