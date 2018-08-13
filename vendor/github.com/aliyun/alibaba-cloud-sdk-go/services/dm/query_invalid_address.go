package dm

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

// QueryInvalidAddress invokes the dm.QueryInvalidAddress API synchronously
// api document: https://help.aliyun.com/api/dm/queryinvalidaddress.html
func (client *Client) QueryInvalidAddress(request *QueryInvalidAddressRequest) (response *QueryInvalidAddressResponse, err error) {
	response = CreateQueryInvalidAddressResponse()
	err = client.DoAction(request, response)
	return
}

// QueryInvalidAddressWithChan invokes the dm.QueryInvalidAddress API asynchronously
// api document: https://help.aliyun.com/api/dm/queryinvalidaddress.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryInvalidAddressWithChan(request *QueryInvalidAddressRequest) (<-chan *QueryInvalidAddressResponse, <-chan error) {
	responseChan := make(chan *QueryInvalidAddressResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryInvalidAddress(request)
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

// QueryInvalidAddressWithCallback invokes the dm.QueryInvalidAddress API asynchronously
// api document: https://help.aliyun.com/api/dm/queryinvalidaddress.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryInvalidAddressWithCallback(request *QueryInvalidAddressRequest, callback func(response *QueryInvalidAddressResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryInvalidAddressResponse
		var err error
		defer close(result)
		response, err = client.QueryInvalidAddress(request)
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

// QueryInvalidAddressRequest is the request struct for api QueryInvalidAddress
type QueryInvalidAddressRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	StartTime            string           `position:"Query" name:"StartTime"`
	EndTime              string           `position:"Query" name:"EndTime"`
	KeyWord              string           `position:"Query" name:"KeyWord"`
	Length               requests.Integer `position:"Query" name:"Length"`
	NextStart            string           `position:"Query" name:"NextStart"`
}

// QueryInvalidAddressResponse is the response struct for api QueryInvalidAddress
type QueryInvalidAddressResponse struct {
	*responses.BaseResponse
	RequestId  string                    `json:"RequestId" xml:"RequestId"`
	NextStart  int                       `json:"NextStart" xml:"NextStart"`
	TotalCount int                       `json:"TotalCount" xml:"TotalCount"`
	Data       DataInQueryInvalidAddress `json:"data" xml:"data"`
}

// CreateQueryInvalidAddressRequest creates a request to invoke QueryInvalidAddress API
func CreateQueryInvalidAddressRequest() (request *QueryInvalidAddressRequest) {
	request = &QueryInvalidAddressRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dm", "2015-11-23", "QueryInvalidAddress", "", "")
	return
}

// CreateQueryInvalidAddressResponse creates a response to parse from QueryInvalidAddress response
func CreateQueryInvalidAddressResponse() (response *QueryInvalidAddressResponse) {
	response = &QueryInvalidAddressResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
