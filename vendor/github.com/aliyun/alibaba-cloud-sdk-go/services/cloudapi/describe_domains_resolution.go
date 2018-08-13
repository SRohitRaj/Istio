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

// DescribeDomainsResolution invokes the cloudapi.DescribeDomainsResolution API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describedomainsresolution.html
func (client *Client) DescribeDomainsResolution(request *DescribeDomainsResolutionRequest) (response *DescribeDomainsResolutionResponse, err error) {
	response = CreateDescribeDomainsResolutionResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainsResolutionWithChan invokes the cloudapi.DescribeDomainsResolution API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describedomainsresolution.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainsResolutionWithChan(request *DescribeDomainsResolutionRequest) (<-chan *DescribeDomainsResolutionResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainsResolutionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainsResolution(request)
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

// DescribeDomainsResolutionWithCallback invokes the cloudapi.DescribeDomainsResolution API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describedomainsresolution.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainsResolutionWithCallback(request *DescribeDomainsResolutionRequest, callback func(response *DescribeDomainsResolutionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainsResolutionResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainsResolution(request)
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

// DescribeDomainsResolutionRequest is the request struct for api DescribeDomainsResolution
type DescribeDomainsResolutionRequest struct {
	*requests.RpcRequest
	GroupId     string `position:"Query" name:"GroupId"`
	DomainNames string `position:"Query" name:"DomainNames"`
}

// DescribeDomainsResolutionResponse is the response struct for api DescribeDomainsResolution
type DescribeDomainsResolutionResponse struct {
	*responses.BaseResponse
	RequestId         string            `json:"RequestId" xml:"RequestId"`
	GroupId           string            `json:"GroupId" xml:"GroupId"`
	DomainResolutions DomainResolutions `json:"DomainResolutions" xml:"DomainResolutions"`
}

// CreateDescribeDomainsResolutionRequest creates a request to invoke DescribeDomainsResolution API
func CreateDescribeDomainsResolutionRequest() (request *DescribeDomainsResolutionRequest) {
	request = &DescribeDomainsResolutionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeDomainsResolution", "apigateway", "openAPI")
	return
}

// CreateDescribeDomainsResolutionResponse creates a response to parse from DescribeDomainsResolution response
func CreateDescribeDomainsResolutionResponse() (response *DescribeDomainsResolutionResponse) {
	response = &DescribeDomainsResolutionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
