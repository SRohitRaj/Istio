package mts

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

// AddMediaTag invokes the mts.AddMediaTag API synchronously
// api document: https://help.aliyun.com/api/mts/addmediatag.html
func (client *Client) AddMediaTag(request *AddMediaTagRequest) (response *AddMediaTagResponse, err error) {
	response = CreateAddMediaTagResponse()
	err = client.DoAction(request, response)
	return
}

// AddMediaTagWithChan invokes the mts.AddMediaTag API asynchronously
// api document: https://help.aliyun.com/api/mts/addmediatag.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddMediaTagWithChan(request *AddMediaTagRequest) (<-chan *AddMediaTagResponse, <-chan error) {
	responseChan := make(chan *AddMediaTagResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddMediaTag(request)
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

// AddMediaTagWithCallback invokes the mts.AddMediaTag API asynchronously
// api document: https://help.aliyun.com/api/mts/addmediatag.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddMediaTagWithCallback(request *AddMediaTagRequest, callback func(response *AddMediaTagResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddMediaTagResponse
		var err error
		defer close(result)
		response, err = client.AddMediaTag(request)
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

// AddMediaTagRequest is the request struct for api AddMediaTag
type AddMediaTagRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	MediaId              string           `position:"Query" name:"MediaId"`
	Tag                  string           `position:"Query" name:"Tag"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// AddMediaTagResponse is the response struct for api AddMediaTag
type AddMediaTagResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateAddMediaTagRequest creates a request to invoke AddMediaTag API
func CreateAddMediaTagRequest() (request *AddMediaTagRequest) {
	request = &AddMediaTagRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "AddMediaTag", "mts", "openAPI")
	return
}

// CreateAddMediaTagResponse creates a response to parse from AddMediaTag response
func CreateAddMediaTagResponse() (response *AddMediaTagResponse) {
	response = &AddMediaTagResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
