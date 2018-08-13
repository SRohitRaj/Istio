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

// DeleteEmailVerification invokes the domain.DeleteEmailVerification API synchronously
// api document: https://help.aliyun.com/api/domain/deleteemailverification.html
func (client *Client) DeleteEmailVerification(request *DeleteEmailVerificationRequest) (response *DeleteEmailVerificationResponse, err error) {
	response = CreateDeleteEmailVerificationResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteEmailVerificationWithChan invokes the domain.DeleteEmailVerification API asynchronously
// api document: https://help.aliyun.com/api/domain/deleteemailverification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteEmailVerificationWithChan(request *DeleteEmailVerificationRequest) (<-chan *DeleteEmailVerificationResponse, <-chan error) {
	responseChan := make(chan *DeleteEmailVerificationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteEmailVerification(request)
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

// DeleteEmailVerificationWithCallback invokes the domain.DeleteEmailVerification API asynchronously
// api document: https://help.aliyun.com/api/domain/deleteemailverification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteEmailVerificationWithCallback(request *DeleteEmailVerificationRequest, callback func(response *DeleteEmailVerificationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteEmailVerificationResponse
		var err error
		defer close(result)
		response, err = client.DeleteEmailVerification(request)
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

// DeleteEmailVerificationRequest is the request struct for api DeleteEmailVerification
type DeleteEmailVerificationRequest struct {
	*requests.RpcRequest
	Lang         string `position:"Query" name:"Lang"`
	Email        string `position:"Query" name:"Email"`
	UserClientIp string `position:"Query" name:"UserClientIp"`
}

// DeleteEmailVerificationResponse is the response struct for api DeleteEmailVerification
type DeleteEmailVerificationResponse struct {
	*responses.BaseResponse
	RequestId   string       `json:"RequestId" xml:"RequestId"`
	SuccessList []SendResult `json:"SuccessList" xml:"SuccessList"`
	FailList    []SendResult `json:"FailList" xml:"FailList"`
}

// CreateDeleteEmailVerificationRequest creates a request to invoke DeleteEmailVerification API
func CreateDeleteEmailVerificationRequest() (request *DeleteEmailVerificationRequest) {
	request = &DeleteEmailVerificationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "DeleteEmailVerification", "", "")
	return
}

// CreateDeleteEmailVerificationResponse creates a response to parse from DeleteEmailVerification response
func CreateDeleteEmailVerificationResponse() (response *DeleteEmailVerificationResponse) {
	response = &DeleteEmailVerificationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
