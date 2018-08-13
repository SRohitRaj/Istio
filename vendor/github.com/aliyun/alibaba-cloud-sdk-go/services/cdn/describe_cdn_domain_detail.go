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

// DescribeCdnDomainDetail invokes the cdn.DescribeCdnDomainDetail API synchronously
// api document: https://help.aliyun.com/api/cdn/describecdndomaindetail.html
func (client *Client) DescribeCdnDomainDetail(request *DescribeCdnDomainDetailRequest) (response *DescribeCdnDomainDetailResponse, err error) {
	response = CreateDescribeCdnDomainDetailResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeCdnDomainDetailWithChan invokes the cdn.DescribeCdnDomainDetail API asynchronously
// api document: https://help.aliyun.com/api/cdn/describecdndomaindetail.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCdnDomainDetailWithChan(request *DescribeCdnDomainDetailRequest) (<-chan *DescribeCdnDomainDetailResponse, <-chan error) {
	responseChan := make(chan *DescribeCdnDomainDetailResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeCdnDomainDetail(request)
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

// DescribeCdnDomainDetailWithCallback invokes the cdn.DescribeCdnDomainDetail API asynchronously
// api document: https://help.aliyun.com/api/cdn/describecdndomaindetail.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeCdnDomainDetailWithCallback(request *DescribeCdnDomainDetailRequest, callback func(response *DescribeCdnDomainDetailResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeCdnDomainDetailResponse
		var err error
		defer close(result)
		response, err = client.DescribeCdnDomainDetail(request)
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

// DescribeCdnDomainDetailRequest is the request struct for api DescribeCdnDomainDetail
type DescribeCdnDomainDetailRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
}

// DescribeCdnDomainDetailResponse is the response struct for api DescribeCdnDomainDetail
type DescribeCdnDomainDetailResponse struct {
	*responses.BaseResponse
	RequestId            string               `json:"RequestId" xml:"RequestId"`
	GetDomainDetailModel GetDomainDetailModel `json:"GetDomainDetailModel" xml:"GetDomainDetailModel"`
}

// CreateDescribeCdnDomainDetailRequest creates a request to invoke DescribeCdnDomainDetail API
func CreateDescribeCdnDomainDetailRequest() (request *DescribeCdnDomainDetailRequest) {
	request = &DescribeCdnDomainDetailRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeCdnDomainDetail", "", "")
	return
}

// CreateDescribeCdnDomainDetailResponse creates a response to parse from DescribeCdnDomainDetail response
func CreateDescribeCdnDomainDetailResponse() (response *DescribeCdnDomainDetailResponse) {
	response = &DescribeCdnDomainDetailResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
