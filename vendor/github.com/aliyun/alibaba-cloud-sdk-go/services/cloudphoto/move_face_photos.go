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

// MoveFacePhotos invokes the cloudphoto.MoveFacePhotos API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/movefacephotos.html
func (client *Client) MoveFacePhotos(request *MoveFacePhotosRequest) (response *MoveFacePhotosResponse, err error) {
	response = CreateMoveFacePhotosResponse()
	err = client.DoAction(request, response)
	return
}

// MoveFacePhotosWithChan invokes the cloudphoto.MoveFacePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/movefacephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MoveFacePhotosWithChan(request *MoveFacePhotosRequest) (<-chan *MoveFacePhotosResponse, <-chan error) {
	responseChan := make(chan *MoveFacePhotosResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.MoveFacePhotos(request)
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

// MoveFacePhotosWithCallback invokes the cloudphoto.MoveFacePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/movefacephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MoveFacePhotosWithCallback(request *MoveFacePhotosRequest, callback func(response *MoveFacePhotosResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *MoveFacePhotosResponse
		var err error
		defer close(result)
		response, err = client.MoveFacePhotos(request)
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

// MoveFacePhotosRequest is the request struct for api MoveFacePhotos
type MoveFacePhotosRequest struct {
	*requests.RpcRequest
	SourceFaceId requests.Integer `position:"Query" name:"SourceFaceId"`
	PhotoId      *[]string        `position:"Query" name:"PhotoId"  type:"Repeated"`
	TargetFaceId requests.Integer `position:"Query" name:"TargetFaceId"`
	StoreName    string           `position:"Query" name:"StoreName"`
	LibraryId    string           `position:"Query" name:"LibraryId"`
}

// MoveFacePhotosResponse is the response struct for api MoveFacePhotos
type MoveFacePhotosResponse struct {
	*responses.BaseResponse
	Code      string   `json:"Code" xml:"Code"`
	Message   string   `json:"Message" xml:"Message"`
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Action    string   `json:"Action" xml:"Action"`
	Results   []Result `json:"Results" xml:"Results"`
}

// CreateMoveFacePhotosRequest creates a request to invoke MoveFacePhotos API
func CreateMoveFacePhotosRequest() (request *MoveFacePhotosRequest) {
	request = &MoveFacePhotosRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "MoveFacePhotos", "cloudphoto", "openAPI")
	return
}

// CreateMoveFacePhotosResponse creates a response to parse from MoveFacePhotos response
func CreateMoveFacePhotosResponse() (response *MoveFacePhotosResponse) {
	response = &MoveFacePhotosResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
