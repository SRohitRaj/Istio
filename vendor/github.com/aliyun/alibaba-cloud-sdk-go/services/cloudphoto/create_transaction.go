package cloudphoto

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

// CreateTransaction invokes the cloudphoto.CreateTransaction API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/createtransaction.html
func (client *Client) CreateTransaction(request *CreateTransactionRequest) (response *CreateTransactionResponse, err error) {
	response = CreateCreateTransactionResponse()
	err = client.DoAction(request, response)
	return
}

// CreateTransactionWithChan invokes the cloudphoto.CreateTransaction API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/createtransaction.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateTransactionWithChan(request *CreateTransactionRequest) (<-chan *CreateTransactionResponse, <-chan error) {
	responseChan := make(chan *CreateTransactionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateTransaction(request)
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

// CreateTransactionWithCallback invokes the cloudphoto.CreateTransaction API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/createtransaction.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateTransactionWithCallback(request *CreateTransactionRequest, callback func(response *CreateTransactionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateTransactionResponse
		var err error
		defer close(result)
		response, err = client.CreateTransaction(request)
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

// CreateTransactionRequest is the request struct for api CreateTransaction
type CreateTransactionRequest struct {
	*requests.RpcRequest
	Size      requests.Integer `position:"Query" name:"Size"`
	Ext       string           `position:"Query" name:"Ext"`
	Force     string           `position:"Query" name:"Force"`
	Md5       string           `position:"Query" name:"Md5"`
	StoreName string           `position:"Query" name:"StoreName"`
	LibraryId string           `position:"Query" name:"LibraryId"`
}

// CreateTransactionResponse is the response struct for api CreateTransaction
type CreateTransactionResponse struct {
	*responses.BaseResponse
	Code        string      `json:"Code" xml:"Code"`
	Message     string      `json:"Message" xml:"Message"`
	RequestId   string      `json:"RequestId" xml:"RequestId"`
	Action      string      `json:"Action" xml:"Action"`
	Transaction Transaction `json:"Transaction" xml:"Transaction"`
}

// CreateCreateTransactionRequest creates a request to invoke CreateTransaction API
func CreateCreateTransactionRequest() (request *CreateTransactionRequest) {
	request = &CreateTransactionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "CreateTransaction", "cloudphoto", "openAPI")
	return
}

// CreateCreateTransactionResponse creates a response to parse from CreateTransaction response
func CreateCreateTransactionResponse() (response *CreateTransactionResponse) {
	response = &CreateTransactionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
