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

// ListFacePhotos invokes the cloudphoto.ListFacePhotos API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/listfacephotos.html
func (client *Client) ListFacePhotos(request *ListFacePhotosRequest) (response *ListFacePhotosResponse, err error) {
	response = CreateListFacePhotosResponse()
	err = client.DoAction(request, response)
	return
}

// ListFacePhotosWithChan invokes the cloudphoto.ListFacePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/listfacephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListFacePhotosWithChan(request *ListFacePhotosRequest) (<-chan *ListFacePhotosResponse, <-chan error) {
	responseChan := make(chan *ListFacePhotosResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListFacePhotos(request)
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

// ListFacePhotosWithCallback invokes the cloudphoto.ListFacePhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/listfacephotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListFacePhotosWithCallback(request *ListFacePhotosRequest, callback func(response *ListFacePhotosResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListFacePhotosResponse
		var err error
		defer close(result)
		response, err = client.ListFacePhotos(request)
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

// ListFacePhotosRequest is the request struct for api ListFacePhotos
type ListFacePhotosRequest struct {
	*requests.RpcRequest
	FaceId    requests.Integer `position:"Query" name:"FaceId"`
	Direction string           `position:"Query" name:"Direction"`
	Size      requests.Integer `position:"Query" name:"Size"`
	Cursor    string           `position:"Query" name:"Cursor"`
	State     string           `position:"Query" name:"State"`
	StoreName string           `position:"Query" name:"StoreName"`
	LibraryId string           `position:"Query" name:"LibraryId"`
}

// ListFacePhotosResponse is the response struct for api ListFacePhotos
type ListFacePhotosResponse struct {
	*responses.BaseResponse
	Code       string   `json:"Code" xml:"Code"`
	Message    string   `json:"Message" xml:"Message"`
	NextCursor string   `json:"NextCursor" xml:"NextCursor"`
	TotalCount int      `json:"TotalCount" xml:"TotalCount"`
	RequestId  string   `json:"RequestId" xml:"RequestId"`
	Action     string   `json:"Action" xml:"Action"`
	Results    []Result `json:"Results" xml:"Results"`
}

// CreateListFacePhotosRequest creates a request to invoke ListFacePhotos API
func CreateListFacePhotosRequest() (request *ListFacePhotosRequest) {
	request = &ListFacePhotosRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "ListFacePhotos", "cloudphoto", "openAPI")
	return
}

// CreateListFacePhotosResponse creates a response to parse from ListFacePhotos response
func CreateListFacePhotosResponse() (response *ListFacePhotosResponse) {
	response = &ListFacePhotosResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
