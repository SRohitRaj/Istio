package ccc

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

// RemovePhoneNumber invokes the ccc.RemovePhoneNumber API synchronously
// api document: https://help.aliyun.com/api/ccc/removephonenumber.html
func (client *Client) RemovePhoneNumber(request *RemovePhoneNumberRequest) (response *RemovePhoneNumberResponse, err error) {
	response = CreateRemovePhoneNumberResponse()
	err = client.DoAction(request, response)
	return
}

// RemovePhoneNumberWithChan invokes the ccc.RemovePhoneNumber API asynchronously
// api document: https://help.aliyun.com/api/ccc/removephonenumber.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemovePhoneNumberWithChan(request *RemovePhoneNumberRequest) (<-chan *RemovePhoneNumberResponse, <-chan error) {
	responseChan := make(chan *RemovePhoneNumberResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemovePhoneNumber(request)
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

// RemovePhoneNumberWithCallback invokes the ccc.RemovePhoneNumber API asynchronously
// api document: https://help.aliyun.com/api/ccc/removephonenumber.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemovePhoneNumberWithCallback(request *RemovePhoneNumberRequest, callback func(response *RemovePhoneNumberResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemovePhoneNumberResponse
		var err error
		defer close(result)
		response, err = client.RemovePhoneNumber(request)
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

// RemovePhoneNumberRequest is the request struct for api RemovePhoneNumber
type RemovePhoneNumberRequest struct {
	*requests.RpcRequest
	InstanceId    string `position:"Query" name:"InstanceId"`
	PhoneNumberId string `position:"Query" name:"PhoneNumberId"`
}

// RemovePhoneNumberResponse is the response struct for api RemovePhoneNumber
type RemovePhoneNumberResponse struct {
	*responses.BaseResponse
	RequestId      string `json:"RequestId" xml:"RequestId"`
	Success        bool   `json:"Success" xml:"Success"`
	Code           string `json:"Code" xml:"Code"`
	Message        string `json:"Message" xml:"Message"`
	HttpStatusCode int    `json:"HttpStatusCode" xml:"HttpStatusCode"`
}

// CreateRemovePhoneNumberRequest creates a request to invoke RemovePhoneNumber API
func CreateRemovePhoneNumberRequest() (request *RemovePhoneNumberRequest) {
	request = &RemovePhoneNumberRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CCC", "2017-07-05", "RemovePhoneNumber", "", "")
	return
}

// CreateRemovePhoneNumberResponse creates a response to parse from RemovePhoneNumber response
func CreateRemovePhoneNumberResponse() (response *RemovePhoneNumberResponse) {
	response = &RemovePhoneNumberResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
