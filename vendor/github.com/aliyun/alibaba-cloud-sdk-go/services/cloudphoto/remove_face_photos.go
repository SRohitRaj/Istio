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

// RemoveFacePhotos invokes the cloudphoto.RemoveFacePhotos API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/removefacephotos.html
func (client *Client) RemoveFacePhotos(request *RemoveFacePhotosRequest) (response *RemoveFacePhotosResponse, err error) {
	response = CreateRemoveFacePhotosResponse()
	err = client.DoAction(request, response)
	return
}

// RemoveFacePhotosWithChan invokes the cloudphoto.RemoveFacePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/removefacephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveFacePhotosWithChan(request *RemoveFacePhotosRequest) (<-chan *RemoveFacePhotosResponse, <-chan error) {
	responseChan := make(chan *RemoveFacePhotosResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemoveFacePhotos(request)
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

// RemoveFacePhotosWithCallback invokes the cloudphoto.RemoveFacePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/removefacephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveFacePhotosWithCallback(request *RemoveFacePhotosRequest, callback func(response *RemoveFacePhotosResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemoveFacePhotosResponse
		var err error
		defer close(result)
		response, err = client.RemoveFacePhotos(request)
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

// RemoveFacePhotosRequest is the request struct for api RemoveFacePhotos
type RemoveFacePhotosRequest struct {
	*requests.RpcRequest
	FaceId    requests.Integer `position:"Query" name:"FaceId"`
	PhotoId   *[]string        `position:"Query" name:"PhotoId"  type:"Repeated"`
	StoreName string           `position:"Query" name:"StoreName"`
	LibraryId string           `position:"Query" name:"LibraryId"`
}

// RemoveFacePhotosResponse is the response struct for api RemoveFacePhotos
type RemoveFacePhotosResponse struct {
	*responses.BaseResponse
	Code      string   `json:"Code" xml:"Code"`
	Message   string   `json:"Message" xml:"Message"`
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Action    string   `json:"Action" xml:"Action"`
	Results   []Result `json:"Results" xml:"Results"`
}

// CreateRemoveFacePhotosRequest creates a request to invoke RemoveFacePhotos API
func CreateRemoveFacePhotosRequest() (request *RemoveFacePhotosRequest) {
	request = &RemoveFacePhotosRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "RemoveFacePhotos", "cloudphoto", "openAPI")
	return
}

// CreateRemoveFacePhotosResponse creates a response to parse from RemoveFacePhotos response
func CreateRemoveFacePhotosResponse() (response *RemoveFacePhotosResponse) {
	response = &RemoveFacePhotosResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
