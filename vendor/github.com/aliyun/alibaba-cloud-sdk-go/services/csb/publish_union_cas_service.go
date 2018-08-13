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

// PublishUnionCasService invokes the csb.PublishUnionCasService API synchronously
// api document: https://help.aliyun.com/api/csb/publishunioncasservice.html
func (client *Client) PublishUnionCasService(request *PublishUnionCasServiceRequest) (response *PublishUnionCasServiceResponse, err error) {
	response = CreatePublishUnionCasServiceResponse()
	err = client.DoAction(request, response)
	return
}

// PublishUnionCasServiceWithChan invokes the csb.PublishUnionCasService API asynchronously
// api document: https://help.aliyun.com/api/csb/publishunioncasservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) PublishUnionCasServiceWithChan(request *PublishUnionCasServiceRequest) (<-chan *PublishUnionCasServiceResponse, <-chan error) {
	responseChan := make(chan *PublishUnionCasServiceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.PublishUnionCasService(request)
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

// PublishUnionCasServiceWithCallback invokes the csb.PublishUnionCasService API asynchronously
// api document: https://help.aliyun.com/api/csb/publishunioncasservice.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) PublishUnionCasServiceWithCallback(request *PublishUnionCasServiceRequest, callback func(response *PublishUnionCasServiceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *PublishUnionCasServiceResponse
		var err error
		defer close(result)
		response, err = client.PublishUnionCasService(request)
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

// PublishUnionCasServiceRequest is the request struct for api PublishUnionCasService
type PublishUnionCasServiceRequest struct {
	*requests.RpcRequest
	CasCsbName string `position:"Query" name:"CasCsbName"`
	Data       string `position:"Body" name:"Data"`
}

// PublishUnionCasServiceResponse is the response struct for api PublishUnionCasService
type PublishUnionCasServiceResponse struct {
	*responses.BaseResponse
	Code      int    `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreatePublishUnionCasServiceRequest creates a request to invoke PublishUnionCasService API
func CreatePublishUnionCasServiceRequest() (request *PublishUnionCasServiceRequest) {
	request = &PublishUnionCasServiceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CSB", "2017-11-18", "PublishUnionCasService", "CSB", "openAPI")
	return
}

// CreatePublishUnionCasServiceResponse creates a response to parse from PublishUnionCasService response
func CreatePublishUnionCasServiceResponse() (response *PublishUnionCasServiceResponse) {
	response = &PublishUnionCasServiceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
