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

// DescribePurchasedApiGroup invokes the cloudapi.DescribePurchasedApiGroup API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describepurchasedapigroup.html
func (client *Client) DescribePurchasedApiGroup(request *DescribePurchasedApiGroupRequest) (response *DescribePurchasedApiGroupResponse, err error) {
	response = CreateDescribePurchasedApiGroupResponse()
	err = client.DoAction(request, response)
	return
}

// DescribePurchasedApiGroupWithChan invokes the cloudapi.DescribePurchasedApiGroup API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describepurchasedapigroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribePurchasedApiGroupWithChan(request *DescribePurchasedApiGroupRequest) (<-chan *DescribePurchasedApiGroupResponse, <-chan error) {
	responseChan := make(chan *DescribePurchasedApiGroupResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribePurchasedApiGroup(request)
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

// DescribePurchasedApiGroupWithCallback invokes the cloudapi.DescribePurchasedApiGroup API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describepurchasedapigroup.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribePurchasedApiGroupWithCallback(request *DescribePurchasedApiGroupRequest, callback func(response *DescribePurchasedApiGroupResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribePurchasedApiGroupResponse
		var err error
		defer close(result)
		response, err = client.DescribePurchasedApiGroup(request)
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

// DescribePurchasedApiGroupRequest is the request struct for api DescribePurchasedApiGroup
type DescribePurchasedApiGroupRequest struct {
	*requests.RpcRequest
	GroupId string `position:"Query" name:"GroupId"`
}

// DescribePurchasedApiGroupResponse is the response struct for api DescribePurchasedApiGroup
type DescribePurchasedApiGroupResponse struct {
	*responses.BaseResponse
	RequestId     string  `json:"RequestId" xml:"RequestId"`
	GroupId       string  `json:"GroupId" xml:"GroupId"`
	GroupName     string  `json:"GroupName" xml:"GroupName"`
	Description   string  `json:"Description" xml:"Description"`
	PurchasedTime string  `json:"PurchasedTime" xml:"PurchasedTime"`
	RegionId      string  `json:"RegionId" xml:"RegionId"`
	Status        string  `json:"Status" xml:"Status"`
	Domains       Domains `json:"Domains" xml:"Domains"`
}

// CreateDescribePurchasedApiGroupRequest creates a request to invoke DescribePurchasedApiGroup API
func CreateDescribePurchasedApiGroupRequest() (request *DescribePurchasedApiGroupRequest) {
	request = &DescribePurchasedApiGroupRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribePurchasedApiGroup", "apigateway", "openAPI")
	return
}

// CreateDescribePurchasedApiGroupResponse creates a response to parse from DescribePurchasedApiGroup response
func CreateDescribePurchasedApiGroupResponse() (response *DescribePurchasedApiGroupResponse) {
	response = &DescribePurchasedApiGroupResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
