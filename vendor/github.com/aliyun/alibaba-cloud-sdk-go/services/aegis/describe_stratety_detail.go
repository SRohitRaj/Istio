package aegis

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

// DescribeStratetyDetail invokes the aegis.DescribeStratetyDetail API synchronously
// api document: https://help.aliyun.com/api/aegis/describestratetydetail.html
func (client *Client) DescribeStratetyDetail(request *DescribeStratetyDetailRequest) (response *DescribeStratetyDetailResponse, err error) {
	response = CreateDescribeStratetyDetailResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeStratetyDetailWithChan invokes the aegis.DescribeStratetyDetail API asynchronously
// api document: https://help.aliyun.com/api/aegis/describestratetydetail.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeStratetyDetailWithChan(request *DescribeStratetyDetailRequest) (<-chan *DescribeStratetyDetailResponse, <-chan error) {
	responseChan := make(chan *DescribeStratetyDetailResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeStratetyDetail(request)
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

// DescribeStratetyDetailWithCallback invokes the aegis.DescribeStratetyDetail API asynchronously
// api document: https://help.aliyun.com/api/aegis/describestratetydetail.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeStratetyDetailWithCallback(request *DescribeStratetyDetailRequest, callback func(response *DescribeStratetyDetailResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeStratetyDetailResponse
		var err error
		defer close(result)
		response, err = client.DescribeStratetyDetail(request)
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

// DescribeStratetyDetailRequest is the request struct for api DescribeStratetyDetail
type DescribeStratetyDetailRequest struct {
	*requests.RpcRequest
	SourceIp        string           `position:"Query" name:"SourceIp"`
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Id              string           `position:"Query" name:"Id"`
}

// DescribeStratetyDetailResponse is the response struct for api DescribeStratetyDetail
type DescribeStratetyDetailResponse struct {
	*responses.BaseResponse
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Strategy  Strategy `json:"Strategy" xml:"Strategy"`
}

// CreateDescribeStratetyDetailRequest creates a request to invoke DescribeStratetyDetail API
func CreateDescribeStratetyDetailRequest() (request *DescribeStratetyDetailRequest) {
	request = &DescribeStratetyDetailRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("aegis", "2016-11-11", "DescribeStratetyDetail", "vipaegis", "openAPI")
	return
}

// CreateDescribeStratetyDetailResponse creates a response to parse from DescribeStratetyDetail response
func CreateDescribeStratetyDetailResponse() (response *DescribeStratetyDetailResponse) {
	response = &DescribeStratetyDetailResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
