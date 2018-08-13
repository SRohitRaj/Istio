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

// AddPhoneNumber invokes the ccc.AddPhoneNumber API synchronously
// api document: https://help.aliyun.com/api/ccc/addphonenumber.html
func (client *Client) AddPhoneNumber(request *AddPhoneNumberRequest) (response *AddPhoneNumberResponse, err error) {
	response = CreateAddPhoneNumberResponse()
	err = client.DoAction(request, response)
	return
}

// AddPhoneNumberWithChan invokes the ccc.AddPhoneNumber API asynchronously
// api document: https://help.aliyun.com/api/ccc/addphonenumber.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddPhoneNumberWithChan(request *AddPhoneNumberRequest) (<-chan *AddPhoneNumberResponse, <-chan error) {
	responseChan := make(chan *AddPhoneNumberResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddPhoneNumber(request)
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

// AddPhoneNumberWithCallback invokes the ccc.AddPhoneNumber API asynchronously
// api document: https://help.aliyun.com/api/ccc/addphonenumber.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddPhoneNumberWithCallback(request *AddPhoneNumberRequest, callback func(response *AddPhoneNumberResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddPhoneNumberResponse
		var err error
		defer close(result)
		response, err = client.AddPhoneNumber(request)
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

// AddPhoneNumberRequest is the request struct for api AddPhoneNumber
type AddPhoneNumberRequest struct {
	*requests.RpcRequest
	InstanceId    string `position:"Query" name:"InstanceId"`
	PhoneNumber   string `position:"Query" name:"PhoneNumber"`
	Usage         string `position:"Query" name:"Usage"`
	ContactFlowId string `position:"Query" name:"ContactFlowId"`
}

// AddPhoneNumberResponse is the response struct for api AddPhoneNumber
type AddPhoneNumberResponse struct {
	*responses.BaseResponse
	RequestId      string      `json:"RequestId" xml:"RequestId"`
	Success        bool        `json:"Success" xml:"Success"`
	Code           string      `json:"Code" xml:"Code"`
	Message        string      `json:"Message" xml:"Message"`
	HttpStatusCode int         `json:"HttpStatusCode" xml:"HttpStatusCode"`
	PhoneNumber    PhoneNumber `json:"PhoneNumber" xml:"PhoneNumber"`
}

// CreateAddPhoneNumberRequest creates a request to invoke AddPhoneNumber API
func CreateAddPhoneNumberRequest() (request *AddPhoneNumberRequest) {
	request = &AddPhoneNumberRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CCC", "2017-07-05", "AddPhoneNumber", "", "")
	return
}

// CreateAddPhoneNumberResponse creates a response to parse from AddPhoneNumber response
func CreateAddPhoneNumberResponse() (response *AddPhoneNumberResponse) {
	response = &AddPhoneNumberResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
