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

// DescribeZones invokes the vpc.DescribeZones API synchronously
// api document: https://help.aliyun.com/api/vpc/describezones.html
func (client *Client) DescribeZones(request *DescribeZonesRequest) (response *DescribeZonesResponse, err error) {
	response = CreateDescribeZonesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeZonesWithChan invokes the vpc.DescribeZones API asynchronously
// api document: https://help.aliyun.com/api/vpc/describezones.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeZonesWithChan(request *DescribeZonesRequest) (<-chan *DescribeZonesResponse, <-chan error) {
	responseChan := make(chan *DescribeZonesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeZones(request)
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

// DescribeZonesWithCallback invokes the vpc.DescribeZones API asynchronously
// api document: https://help.aliyun.com/api/vpc/describezones.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeZonesWithCallback(request *DescribeZonesRequest, callback func(response *DescribeZonesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeZonesResponse
		var err error
		defer close(result)
		response, err = client.DescribeZones(request)
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

// DescribeZonesRequest is the request struct for api DescribeZones
type DescribeZonesRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// DescribeZonesResponse is the response struct for api DescribeZones
type DescribeZonesResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Zones     Zones  `json:"Zones" xml:"Zones"`
}

// CreateDescribeZonesRequest creates a request to invoke DescribeZones API
func CreateDescribeZonesRequest() (request *DescribeZonesRequest) {
	request = &DescribeZonesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "DescribeZones", "vpc", "openAPI")
	return
}

// CreateDescribeZonesResponse creates a response to parse from DescribeZones response
func CreateDescribeZonesResponse() (response *DescribeZonesResponse) {
	response = &DescribeZonesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
