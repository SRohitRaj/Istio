package ecs

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

// DescribeInstanceStatus invokes the ecs.DescribeInstanceStatus API synchronously
// api document: https://help.aliyun.com/api/ecs/describeinstancestatus.html
func (client *Client) DescribeInstanceStatus(request *DescribeInstanceStatusRequest) (response *DescribeInstanceStatusResponse, err error) {
	response = CreateDescribeInstanceStatusResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeInstanceStatusWithChan invokes the ecs.DescribeInstanceStatus API asynchronously
// api document: https://help.aliyun.com/api/ecs/describeinstancestatus.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeInstanceStatusWithChan(request *DescribeInstanceStatusRequest) (<-chan *DescribeInstanceStatusResponse, <-chan error) {
	responseChan := make(chan *DescribeInstanceStatusResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeInstanceStatus(request)
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

// DescribeInstanceStatusWithCallback invokes the ecs.DescribeInstanceStatus API asynchronously
// api document: https://help.aliyun.com/api/ecs/describeinstancestatus.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeInstanceStatusWithCallback(request *DescribeInstanceStatusRequest, callback func(response *DescribeInstanceStatusResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeInstanceStatusResponse
		var err error
		defer close(result)
		response, err = client.DescribeInstanceStatus(request)
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

// DescribeInstanceStatusRequest is the request struct for api DescribeInstanceStatus
type DescribeInstanceStatusRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ZoneId               string           `position:"Query" name:"ZoneId"`
	ClusterId            string           `position:"Query" name:"ClusterId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// DescribeInstanceStatusResponse is the response struct for api DescribeInstanceStatus
type DescribeInstanceStatusResponse struct {
	*responses.BaseResponse
	RequestId        string           `json:"RequestId" xml:"RequestId"`
	TotalCount       int              `json:"TotalCount" xml:"TotalCount"`
	PageNumber       int              `json:"PageNumber" xml:"PageNumber"`
	PageSize         int              `json:"PageSize" xml:"PageSize"`
	InstanceStatuses InstanceStatuses `json:"InstanceStatuses" xml:"InstanceStatuses"`
}

// CreateDescribeInstanceStatusRequest creates a request to invoke DescribeInstanceStatus API
func CreateDescribeInstanceStatusRequest() (request *DescribeInstanceStatusRequest) {
	request = &DescribeInstanceStatusRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DescribeInstanceStatus", "ecs", "openAPI")
	return
}

// CreateDescribeInstanceStatusResponse creates a response to parse from DescribeInstanceStatus response
func CreateDescribeInstanceStatusResponse() (response *DescribeInstanceStatusResponse) {
	response = &DescribeInstanceStatusResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
