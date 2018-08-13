package domain

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

// SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID invokes the domain.SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID API synchronously
// api document: https://help.aliyun.com/api/domain/savetaskforsubmittingdomainrealnameverificationbyregistrantprofileid.html
func (client *Client) SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID(request *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest) (response *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse, err error) {
	response = CreateSaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse()
	err = client.DoAction(request, response)
	return
}

// SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDWithChan invokes the domain.SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID API asynchronously
// api document: https://help.aliyun.com/api/domain/savetaskforsubmittingdomainrealnameverificationbyregistrantprofileid.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDWithChan(request *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest) (<-chan *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse, <-chan error) {
	responseChan := make(chan *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID(request)
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

// SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDWithCallback invokes the domain.SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID API asynchronously
// api document: https://help.aliyun.com/api/domain/savetaskforsubmittingdomainrealnameverificationbyregistrantprofileid.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDWithCallback(request *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest, callback func(response *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse
		var err error
		defer close(result)
		response, err = client.SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID(request)
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

// SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest is the request struct for api SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID
type SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest struct {
	*requests.RpcRequest
	UserClientIp        string           `position:"Query" name:"UserClientIp"`
	Lang                string           `position:"Query" name:"Lang"`
	DomainName          string           `position:"Query" name:"DomainName"`
	InstanceId          string           `position:"Query" name:"InstanceId"`
	RegistrantProfileId requests.Integer `position:"Query" name:"RegistrantProfileId"`
}

// SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse is the response struct for api SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID
type SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	TaskNo    string `json:"TaskNo" xml:"TaskNo"`
}

// CreateSaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest creates a request to invoke SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID API
func CreateSaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest() (request *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest) {
	request = &SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID", "", "")
	return
}

// CreateSaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse creates a response to parse from SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileID response
func CreateSaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse() (response *SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse) {
	response = &SaveTaskForSubmittingDomainRealNameVerificationByRegistrantProfileIDResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
