package ddospro

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

// DescribeUnBlockholeCount invokes the ddospro.DescribeUnBlockholeCount API synchronously
// api document: https://help.aliyun.com/api/ddospro/describeunblockholecount.html
func (client *Client) DescribeUnBlockholeCount(request *DescribeUnBlockholeCountRequest) (response *DescribeUnBlockholeCountResponse, err error) {
	response = CreateDescribeUnBlockholeCountResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeUnBlockholeCountWithChan invokes the ddospro.DescribeUnBlockholeCount API asynchronously
// api document: https://help.aliyun.com/api/ddospro/describeunblockholecount.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeUnBlockholeCountWithChan(request *DescribeUnBlockholeCountRequest) (<-chan *DescribeUnBlockholeCountResponse, <-chan error) {
	responseChan := make(chan *DescribeUnBlockholeCountResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeUnBlockholeCount(request)
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

// DescribeUnBlockholeCountWithCallback invokes the ddospro.DescribeUnBlockholeCount API asynchronously
// api document: https://help.aliyun.com/api/ddospro/describeunblockholecount.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeUnBlockholeCountWithCallback(request *DescribeUnBlockholeCountRequest, callback func(response *DescribeUnBlockholeCountResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeUnBlockholeCountResponse
		var err error
		defer close(result)
		response, err = client.DescribeUnBlockholeCount(request)
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

// DescribeUnBlockholeCountRequest is the request struct for api DescribeUnBlockholeCount
type DescribeUnBlockholeCountRequest struct {
	*requests.RpcRequest
	SourceIp        string           `position:"Query" name:"SourceIp"`
	ResourceOwnerId requests.Integer `position:"Query" name:"ResourceOwnerId"`
}

// DescribeUnBlockholeCountResponse is the response struct for api DescribeUnBlockholeCount
type DescribeUnBlockholeCountResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateDescribeUnBlockholeCountRequest creates a request to invoke DescribeUnBlockholeCount API
func CreateDescribeUnBlockholeCountRequest() (request *DescribeUnBlockholeCountRequest) {
	request = &DescribeUnBlockholeCountRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("DDoSPro", "2017-07-25", "DescribeUnBlockholeCount", "", "")
	return
}

// CreateDescribeUnBlockholeCountResponse creates a response to parse from DescribeUnBlockholeCount response
func CreateDescribeUnBlockholeCountResponse() (response *DescribeUnBlockholeCountResponse) {
	response = &DescribeUnBlockholeCountResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
