package csb

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

// GetOrder invokes the csb.GetOrder API synchronously
// api document: https://help.aliyun.com/api/csb/getorder.html
func (client *Client) GetOrder(request *GetOrderRequest) (response *GetOrderResponse, err error) {
	response = CreateGetOrderResponse()
	err = client.DoAction(request, response)
	return
}

// GetOrderWithChan invokes the csb.GetOrder API asynchronously
// api document: https://help.aliyun.com/api/csb/getorder.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetOrderWithChan(request *GetOrderRequest) (<-chan *GetOrderResponse, <-chan error) {
	responseChan := make(chan *GetOrderResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetOrder(request)
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

// GetOrderWithCallback invokes the csb.GetOrder API asynchronously
// api document: https://help.aliyun.com/api/csb/getorder.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetOrderWithCallback(request *GetOrderRequest, callback func(response *GetOrderResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetOrderResponse
		var err error
		defer close(result)
		response, err = client.GetOrder(request)
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

// GetOrderRequest is the request struct for api GetOrder
type GetOrderRequest struct {
	*requests.RpcRequest
	OrderId     requests.Integer `position:"Query" name:"OrderId"`
	ServiceName string           `position:"Query" name:"ServiceName"`
}

// GetOrderResponse is the response struct for api GetOrder
type GetOrderResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateGetOrderRequest creates a request to invoke GetOrder API
func CreateGetOrderRequest() (request *GetOrderRequest) {
	request = &GetOrderRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CSB", "2017-11-18", "GetOrder", "CSB", "openAPI")
	return
}

// CreateGetOrderResponse creates a response to parse from GetOrder response
func CreateGetOrderResponse() (response *GetOrderResponse) {
	response = &GetOrderResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
