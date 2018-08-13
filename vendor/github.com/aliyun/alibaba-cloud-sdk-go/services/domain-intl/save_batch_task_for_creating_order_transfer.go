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

// SaveBatchTaskForCreatingOrderTransfer invokes the domain_intl.SaveBatchTaskForCreatingOrderTransfer API synchronously
// api document: https://help.aliyun.com/api/domain-intl/savebatchtaskforcreatingordertransfer.html
func (client *Client) SaveBatchTaskForCreatingOrderTransfer(request *SaveBatchTaskForCreatingOrderTransferRequest) (response *SaveBatchTaskForCreatingOrderTransferResponse, err error) {
	response = CreateSaveBatchTaskForCreatingOrderTransferResponse()
	err = client.DoAction(request, response)
	return
}

// SaveBatchTaskForCreatingOrderTransferWithChan invokes the domain_intl.SaveBatchTaskForCreatingOrderTransfer API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savebatchtaskforcreatingordertransfer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveBatchTaskForCreatingOrderTransferWithChan(request *SaveBatchTaskForCreatingOrderTransferRequest) (<-chan *SaveBatchTaskForCreatingOrderTransferResponse, <-chan error) {
	responseChan := make(chan *SaveBatchTaskForCreatingOrderTransferResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveBatchTaskForCreatingOrderTransfer(request)
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

// SaveBatchTaskForCreatingOrderTransferWithCallback invokes the domain_intl.SaveBatchTaskForCreatingOrderTransfer API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savebatchtaskforcreatingordertransfer.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveBatchTaskForCreatingOrderTransferWithCallback(request *SaveBatchTaskForCreatingOrderTransferRequest, callback func(response *SaveBatchTaskForCreatingOrderTransferResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveBatchTaskForCreatingOrderTransferResponse
		var err error
		defer close(result)
		response, err = client.SaveBatchTaskForCreatingOrderTransfer(request)
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

// SaveBatchTaskForCreatingOrderTransferRequest is the request struct for api SaveBatchTaskForCreatingOrderTransfer
type SaveBatchTaskForCreatingOrderTransferRequest struct {
	*requests.RpcRequest
	UserClientIp       string                                                     `position:"Query" name:"UserClientIp"`
	Lang               string                                                     `position:"Query" name:"Lang"`
	OrderTransferParam *[]SaveBatchTaskForCreatingOrderTransferOrderTransferParam `position:"Query" name:"OrderTransferParam"  type:"Repeated"`
}

// SaveBatchTaskForCreatingOrderTransferOrderTransferParam is a repeated param struct in SaveBatchTaskForCreatingOrderTransferRequest
type SaveBatchTaskForCreatingOrderTransferOrderTransferParam struct {
	DomainName            string `name:"DomainName"`
	AuthorizationCode     string `name:"AuthorizationCode"`
	RegistrantProfileId   string `name:"RegistrantProfileId"`
	PermitPremiumTransfer string `name:"PermitPremiumTransfer"`
}

// SaveBatchTaskForCreatingOrderTransferResponse is the response struct for api SaveBatchTaskForCreatingOrderTransfer
type SaveBatchTaskForCreatingOrderTransferResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	TaskNo    string `json:"TaskNo" xml:"TaskNo"`
}

// CreateSaveBatchTaskForCreatingOrderTransferRequest creates a request to invoke SaveBatchTaskForCreatingOrderTransfer API
func CreateSaveBatchTaskForCreatingOrderTransferRequest() (request *SaveBatchTaskForCreatingOrderTransferRequest) {
	request = &SaveBatchTaskForCreatingOrderTransferRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain-intl", "2017-12-18", "SaveBatchTaskForCreatingOrderTransfer", "domain", "openAPI")
	return
}

// CreateSaveBatchTaskForCreatingOrderTransferResponse creates a response to parse from SaveBatchTaskForCreatingOrderTransfer response
func CreateSaveBatchTaskForCreatingOrderTransferResponse() (response *SaveBatchTaskForCreatingOrderTransferResponse) {
	response = &SaveBatchTaskForCreatingOrderTransferResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
