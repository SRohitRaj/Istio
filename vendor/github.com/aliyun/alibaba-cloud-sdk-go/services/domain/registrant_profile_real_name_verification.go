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

// RegistrantProfileRealNameVerification invokes the domain.RegistrantProfileRealNameVerification API synchronously
// api document: https://help.aliyun.com/api/domain/registrantprofilerealnameverification.html
func (client *Client) RegistrantProfileRealNameVerification(request *RegistrantProfileRealNameVerificationRequest) (response *RegistrantProfileRealNameVerificationResponse, err error) {
	response = CreateRegistrantProfileRealNameVerificationResponse()
	err = client.DoAction(request, response)
	return
}

// RegistrantProfileRealNameVerificationWithChan invokes the domain.RegistrantProfileRealNameVerification API asynchronously
// api document: https://help.aliyun.com/api/domain/registrantprofilerealnameverification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RegistrantProfileRealNameVerificationWithChan(request *RegistrantProfileRealNameVerificationRequest) (<-chan *RegistrantProfileRealNameVerificationResponse, <-chan error) {
	responseChan := make(chan *RegistrantProfileRealNameVerificationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RegistrantProfileRealNameVerification(request)
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

// RegistrantProfileRealNameVerificationWithCallback invokes the domain.RegistrantProfileRealNameVerification API asynchronously
// api document: https://help.aliyun.com/api/domain/registrantprofilerealnameverification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RegistrantProfileRealNameVerificationWithCallback(request *RegistrantProfileRealNameVerificationRequest, callback func(response *RegistrantProfileRealNameVerificationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RegistrantProfileRealNameVerificationResponse
		var err error
		defer close(result)
		response, err = client.RegistrantProfileRealNameVerification(request)
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

// RegistrantProfileRealNameVerificationRequest is the request struct for api RegistrantProfileRealNameVerification
type RegistrantProfileRealNameVerificationRequest struct {
	*requests.RpcRequest
	UserClientIp           string           `position:"Query" name:"UserClientIp"`
	Lang                   string           `position:"Query" name:"Lang"`
	RegistrantProfileID    requests.Integer `position:"Query" name:"RegistrantProfileID"`
	IdentityCredential     string           `position:"Body" name:"IdentityCredential"`
	IdentityCredentialNo   string           `position:"Query" name:"IdentityCredentialNo"`
	IdentityCredentialType string           `position:"Query" name:"IdentityCredentialType"`
}

// RegistrantProfileRealNameVerificationResponse is the response struct for api RegistrantProfileRealNameVerification
type RegistrantProfileRealNameVerificationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateRegistrantProfileRealNameVerificationRequest creates a request to invoke RegistrantProfileRealNameVerification API
func CreateRegistrantProfileRealNameVerificationRequest() (request *RegistrantProfileRealNameVerificationRequest) {
	request = &RegistrantProfileRealNameVerificationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "RegistrantProfileRealNameVerification", "", "")
	return
}

// CreateRegistrantProfileRealNameVerificationResponse creates a response to parse from RegistrantProfileRealNameVerification response
func CreateRegistrantProfileRealNameVerificationResponse() (response *RegistrantProfileRealNameVerificationResponse) {
	response = &RegistrantProfileRealNameVerificationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
