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

// DeletePhotoStore invokes the cloudphoto.DeletePhotoStore API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/deletephotostore.html
func (client *Client) DeletePhotoStore(request *DeletePhotoStoreRequest) (response *DeletePhotoStoreResponse, err error) {
	response = CreateDeletePhotoStoreResponse()
	err = client.DoAction(request, response)
	return
}

// DeletePhotoStoreWithChan invokes the cloudphoto.DeletePhotoStore API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/deletephotostore.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeletePhotoStoreWithChan(request *DeletePhotoStoreRequest) (<-chan *DeletePhotoStoreResponse, <-chan error) {
	responseChan := make(chan *DeletePhotoStoreResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeletePhotoStore(request)
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

// DeletePhotoStoreWithCallback invokes the cloudphoto.DeletePhotoStore API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/deletephotostore.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeletePhotoStoreWithCallback(request *DeletePhotoStoreRequest, callback func(response *DeletePhotoStoreResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeletePhotoStoreResponse
		var err error
		defer close(result)
		response, err = client.DeletePhotoStore(request)
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

// DeletePhotoStoreRequest is the request struct for api DeletePhotoStore
type DeletePhotoStoreRequest struct {
	*requests.RpcRequest
	StoreName string `position:"Query" name:"StoreName"`
}

// DeletePhotoStoreResponse is the response struct for api DeletePhotoStore
type DeletePhotoStoreResponse struct {
	*responses.BaseResponse
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
	Action    string `json:"Action" xml:"Action"`
}

// CreateDeletePhotoStoreRequest creates a request to invoke DeletePhotoStore API
func CreateDeletePhotoStoreRequest() (request *DeletePhotoStoreRequest) {
	request = &DeletePhotoStoreRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "DeletePhotoStore", "cloudphoto", "openAPI")
	return
}

// CreateDeletePhotoStoreResponse creates a response to parse from DeletePhotoStore response
func CreateDeletePhotoStoreResponse() (response *DeletePhotoStoreResponse) {
	response = &DeletePhotoStoreResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
