package domain_intl

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

// SaveSingleTaskForCreatingOrderRedeem invokes the domain_intl.SaveSingleTaskForCreatingOrderRedeem API synchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskforcreatingorderredeem.html
func (client *Client) SaveSingleTaskForCreatingOrderRedeem(request *SaveSingleTaskForCreatingOrderRedeemRequest) (response *SaveSingleTaskForCreatingOrderRedeemResponse, err error) {
	response = CreateSaveSingleTaskForCreatingOrderRedeemResponse()
	err = client.DoAction(request, response)
	return
}

// SaveSingleTaskForCreatingOrderRedeemWithChan invokes the domain_intl.SaveSingleTaskForCreatingOrderRedeem API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskforcreatingorderredeem.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForCreatingOrderRedeemWithChan(request *SaveSingleTaskForCreatingOrderRedeemRequest) (<-chan *SaveSingleTaskForCreatingOrderRedeemResponse, <-chan error) {
	responseChan := make(chan *SaveSingleTaskForCreatingOrderRedeemResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveSingleTaskForCreatingOrderRedeem(request)
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

// SaveSingleTaskForCreatingOrderRedeemWithCallback invokes the domain_intl.SaveSingleTaskForCreatingOrderRedeem API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskforcreatingorderredeem.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForCreatingOrderRedeemWithCallback(request *SaveSingleTaskForCreatingOrderRedeemRequest, callback func(response *SaveSingleTaskForCreatingOrderRedeemResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveSingleTaskForCreatingOrderRedeemResponse
		var err error
		defer close(result)
		response, err = client.SaveSingleTaskForCreatingOrderRedeem(request)
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

// SaveSingleTaskForCreatingOrderRedeemRequest is the request struct for api SaveSingleTaskForCreatingOrderRedeem
type SaveSingleTaskForCreatingOrderRedeemRequest struct {
	*requests.RpcRequest
	UserClientIp          string           `position:"Query" name:"UserClientIp"`
	Lang                  string           `position:"Query" name:"Lang"`
	DomainName            string           `position:"Query" name:"DomainName"`
	CurrentExpirationDate requests.Integer `position:"Query" name:"CurrentExpirationDate"`
}

// SaveSingleTaskForCreatingOrderRedeemResponse is the response struct for api SaveSingleTaskForCreatingOrderRedeem
type SaveSingleTaskForCreatingOrderRedeemResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	TaskNo    string `json:"TaskNo" xml:"TaskNo"`
}

// CreateSaveSingleTaskForCreatingOrderRedeemRequest creates a request to invoke SaveSingleTaskForCreatingOrderRedeem API
func CreateSaveSingleTaskForCreatingOrderRedeemRequest() (request *SaveSingleTaskForCreatingOrderRedeemRequest) {
	request = &SaveSingleTaskForCreatingOrderRedeemRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain-intl", "2017-12-18", "SaveSingleTaskForCreatingOrderRedeem", "", "")
	return
}

// CreateSaveSingleTaskForCreatingOrderRedeemResponse creates a response to parse from SaveSingleTaskForCreatingOrderRedeem response
func CreateSaveSingleTaskForCreatingOrderRedeemResponse() (response *SaveSingleTaskForCreatingOrderRedeemResponse) {
	response = &SaveSingleTaskForCreatingOrderRedeemResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
