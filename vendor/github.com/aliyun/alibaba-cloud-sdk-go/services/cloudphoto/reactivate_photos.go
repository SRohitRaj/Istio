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

// ReactivatePhotos invokes the cloudphoto.ReactivatePhotos API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/reactivatephotos.html
func (client *Client) ReactivatePhotos(request *ReactivatePhotosRequest) (response *ReactivatePhotosResponse, err error) {
	response = CreateReactivatePhotosResponse()
	err = client.DoAction(request, response)
	return
}

// ReactivatePhotosWithChan invokes the cloudphoto.ReactivatePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/reactivatephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ReactivatePhotosWithChan(request *ReactivatePhotosRequest) (<-chan *ReactivatePhotosResponse, <-chan error) {
	responseChan := make(chan *ReactivatePhotosResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ReactivatePhotos(request)
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

// ReactivatePhotosWithCallback invokes the cloudphoto.ReactivatePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/reactivatephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ReactivatePhotosWithCallback(request *ReactivatePhotosRequest, callback func(response *ReactivatePhotosResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ReactivatePhotosResponse
		var err error
		defer close(result)
		response, err = client.ReactivatePhotos(request)
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

// ReactivatePhotosRequest is the request struct for api ReactivatePhotos
type ReactivatePhotosRequest struct {
	*requests.RpcRequest
	PhotoId   *[]string `position:"Query" name:"PhotoId"  type:"Repeated"`
	StoreName string    `position:"Query" name:"StoreName"`
	LibraryId string    `position:"Query" name:"LibraryId"`
}

// ReactivatePhotosResponse is the response struct for api ReactivatePhotos
type ReactivatePhotosResponse struct {
	*responses.BaseResponse
	Code      string   `json:"Code" xml:"Code"`
	Message   string   `json:"Message" xml:"Message"`
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Action    string   `json:"Action" xml:"Action"`
	Results   []Result `json:"Results" xml:"Results"`
}

// CreateReactivatePhotosRequest creates a request to invoke ReactivatePhotos API
func CreateReactivatePhotosRequest() (request *ReactivatePhotosRequest) {
	request = &ReactivatePhotosRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "ReactivatePhotos", "cloudphoto", "openAPI")
	return
}

// CreateReactivatePhotosResponse creates a response to parse from ReactivatePhotos response
func CreateReactivatePhotosResponse() (response *ReactivatePhotosResponse) {
	response = &ReactivatePhotosResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
