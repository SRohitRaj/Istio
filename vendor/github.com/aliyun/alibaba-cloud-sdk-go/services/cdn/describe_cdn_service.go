package cdn

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

// DescribeCdnService invokes the cdn.DescribeCdnService API synchronously
// api document: https://help.aliyun.com/api/cdn/describecdnservice.html
func (client *Client) DescribeCdnService(request *DescribeCdnServiceRequest) (response *DescribeCdnServiceResponse, err error) {
	response = CreateDescribeCdnServiceResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeCdnServiceWithChan invokes the cdn.DescribeCdnService API asynchronously
// api document: https://help.aliyun.com/api/cdn/describecdnservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCdnServiceWithChan(request *DescribeCdnServiceRequest) (<-chan *DescribeCdnServiceResponse, <-chan error) {
	responseChan := make(chan *DescribeCdnServiceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeCdnService(request)
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

// DescribeCdnServiceWithCallback invokes the cdn.DescribeCdnService API asynchronously
// api document: https://help.aliyun.com/api/cdn/describecdnservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCdnServiceWithCallback(request *DescribeCdnServiceRequest, callback func(response *DescribeCdnServiceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeCdnServiceResponse
		var err error
		defer close(result)
		response, err = client.DescribeCdnService(request)
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

// DescribeCdnServiceRequest is the request struct for api DescribeCdnService
type DescribeCdnServiceRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
}

// DescribeCdnServiceResponse is the response struct for api DescribeCdnService
type DescribeCdnServiceResponse struct {
	*responses.BaseResponse
	RequestId          string         `json:"RequestId" xml:"RequestId"`
	InstanceId         string         `json:"InstanceId" xml:"InstanceId"`
	InternetChargeType string         `json:"InternetChargeType" xml:"InternetChargeType"`
	OpeningTime        string         `json:"OpeningTime" xml:"OpeningTime"`
	ChangingChargeType string         `json:"ChangingChargeType" xml:"ChangingChargeType"`
	ChangingAffectTime string         `json:"ChangingAffectTime" xml:"ChangingAffectTime"`
	OperationLocks     OperationLocks `json:"OperationLocks" xml:"OperationLocks"`
}

// CreateDescribeCdnServiceRequest creates a request to invoke DescribeCdnService API
func CreateDescribeCdnServiceRequest() (request *DescribeCdnServiceRequest) {
	request = &DescribeCdnServiceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeCdnService", "", "")
	return
}

// CreateDescribeCdnServiceResponse creates a response to parse from DescribeCdnService response
func CreateDescribeCdnServiceResponse() (response *DescribeCdnServiceResponse) {
	response = &DescribeCdnServiceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
