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

// SaveSingleTaskForCreatingOrderActivate invokes the domain_intl.SaveSingleTaskForCreatingOrderActivate API synchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskforcreatingorderactivate.html
func (client *Client) SaveSingleTaskForCreatingOrderActivate(request *SaveSingleTaskForCreatingOrderActivateRequest) (response *SaveSingleTaskForCreatingOrderActivateResponse, err error) {
	response = CreateSaveSingleTaskForCreatingOrderActivateResponse()
	err = client.DoAction(request, response)
	return
}

// SaveSingleTaskForCreatingOrderActivateWithChan invokes the domain_intl.SaveSingleTaskForCreatingOrderActivate API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskforcreatingorderactivate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForCreatingOrderActivateWithChan(request *SaveSingleTaskForCreatingOrderActivateRequest) (<-chan *SaveSingleTaskForCreatingOrderActivateResponse, <-chan error) {
	responseChan := make(chan *SaveSingleTaskForCreatingOrderActivateResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveSingleTaskForCreatingOrderActivate(request)
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

// SaveSingleTaskForCreatingOrderActivateWithCallback invokes the domain_intl.SaveSingleTaskForCreatingOrderActivate API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskforcreatingorderactivate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForCreatingOrderActivateWithCallback(request *SaveSingleTaskForCreatingOrderActivateRequest, callback func(response *SaveSingleTaskForCreatingOrderActivateResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveSingleTaskForCreatingOrderActivateResponse
		var err error
		defer close(result)
		response, err = client.SaveSingleTaskForCreatingOrderActivate(request)
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

// SaveSingleTaskForCreatingOrderActivateRequest is the request struct for api SaveSingleTaskForCreatingOrderActivate
type SaveSingleTaskForCreatingOrderActivateRequest struct {
	*requests.RpcRequest
	Dns2                    string           `position:"Query" name:"Dns2"`
	UserClientIp            string           `position:"Query" name:"UserClientIp"`
	Lang                    string           `position:"Query" name:"Lang"`
	DomainName              string           `position:"Query" name:"DomainName"`
	SubscriptionDuration    requests.Integer `position:"Query" name:"SubscriptionDuration"`
	RegistrantProfileId     requests.Integer `position:"Query" name:"RegistrantProfileId"`
	EnableDomainProxy       requests.Boolean `position:"Query" name:"EnableDomainProxy"`
	PermitPremiumActivation requests.Boolean `position:"Query" name:"PermitPremiumActivation"`
	AliyunDns               requests.Boolean `position:"Query" name:"AliyunDns"`
	Dns1                    string           `position:"Query" name:"Dns1"`
	Country                 string           `position:"Query" name:"Country"`
	City                    string           `position:"Query" name:"City"`
	RegistrantOrganization  string           `position:"Query" name:"RegistrantOrganization"`
	RegistrantName          string           `position:"Query" name:"RegistrantName"`
	Province                string           `position:"Query" name:"Province"`
	Address                 string           `position:"Query" name:"Address"`
	Email                   string           `position:"Query" name:"Email"`
	PostalCode              string           `position:"Query" name:"PostalCode"`
	TelArea                 string           `position:"Query" name:"TelArea"`
	Telephone               string           `position:"Query" name:"Telephone"`
	TelExt                  string           `position:"Query" name:"TelExt"`
}

// SaveSingleTaskForCreatingOrderActivateResponse is the response struct for api SaveSingleTaskForCreatingOrderActivate
type SaveSingleTaskForCreatingOrderActivateResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	TaskNo    string `json:"TaskNo" xml:"TaskNo"`
}

// CreateSaveSingleTaskForCreatingOrderActivateRequest creates a request to invoke SaveSingleTaskForCreatingOrderActivate API
func CreateSaveSingleTaskForCreatingOrderActivateRequest() (request *SaveSingleTaskForCreatingOrderActivateRequest) {
	request = &SaveSingleTaskForCreatingOrderActivateRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain-intl", "2017-12-18", "SaveSingleTaskForCreatingOrderActivate", "domain", "openAPI")
	return
}

// CreateSaveSingleTaskForCreatingOrderActivateResponse creates a response to parse from SaveSingleTaskForCreatingOrderActivate response
func CreateSaveSingleTaskForCreatingOrderActivateResponse() (response *SaveSingleTaskForCreatingOrderActivateResponse) {
	response = &SaveSingleTaskForCreatingOrderActivateResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
