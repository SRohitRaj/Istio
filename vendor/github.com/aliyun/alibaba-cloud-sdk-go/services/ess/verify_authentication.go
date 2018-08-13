package ess

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

// VerifyAuthentication invokes the ess.VerifyAuthentication API synchronously
// api document: https://help.aliyun.com/api/ess/verifyauthentication.html
func (client *Client) VerifyAuthentication(request *VerifyAuthenticationRequest) (response *VerifyAuthenticationResponse, err error) {
	response = CreateVerifyAuthenticationResponse()
	err = client.DoAction(request, response)
	return
}

// VerifyAuthenticationWithChan invokes the ess.VerifyAuthentication API asynchronously
// api document: https://help.aliyun.com/api/ess/verifyauthentication.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) VerifyAuthenticationWithChan(request *VerifyAuthenticationRequest) (<-chan *VerifyAuthenticationResponse, <-chan error) {
	responseChan := make(chan *VerifyAuthenticationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.VerifyAuthentication(request)
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

// VerifyAuthenticationWithCallback invokes the ess.VerifyAuthentication API asynchronously
// api document: https://help.aliyun.com/api/ess/verifyauthentication.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) VerifyAuthenticationWithCallback(request *VerifyAuthenticationRequest, callback func(response *VerifyAuthenticationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *VerifyAuthenticationResponse
		var err error
		defer close(result)
		response, err = client.VerifyAuthentication(request)
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

// VerifyAuthenticationRequest is the request struct for api VerifyAuthentication
type VerifyAuthenticationRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Uid                  requests.Integer `position:"Query" name:"Uid"`
}

// VerifyAuthenticationResponse is the response struct for api VerifyAuthentication
type VerifyAuthenticationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateVerifyAuthenticationRequest creates a request to invoke VerifyAuthentication API
func CreateVerifyAuthenticationRequest() (request *VerifyAuthenticationRequest) {
	request = &VerifyAuthenticationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ess", "2014-08-28", "VerifyAuthentication", "ess", "openAPI")
	return
}

// CreateVerifyAuthenticationResponse creates a response to parse from VerifyAuthentication response
func CreateVerifyAuthenticationResponse() (response *VerifyAuthenticationResponse) {
	response = &VerifyAuthenticationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
