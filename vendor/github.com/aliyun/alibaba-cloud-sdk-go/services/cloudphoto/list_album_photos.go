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

// ListAlbumPhotos invokes the cloudphoto.ListAlbumPhotos API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/listalbumphotos.html
func (client *Client) ListAlbumPhotos(request *ListAlbumPhotosRequest) (response *ListAlbumPhotosResponse, err error) {
	response = CreateListAlbumPhotosResponse()
	err = client.DoAction(request, response)
	return
}

// ListAlbumPhotosWithChan invokes the cloudphoto.ListAlbumPhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/listalbumphotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListAlbumPhotosWithChan(request *ListAlbumPhotosRequest) (<-chan *ListAlbumPhotosResponse, <-chan error) {
	responseChan := make(chan *ListAlbumPhotosResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ListAlbumPhotos(request)
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

// ListAlbumPhotosWithCallback invokes the cloudphoto.ListAlbumPhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/listalbumphotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ListAlbumPhotosWithCallback(request *ListAlbumPhotosRequest, callback func(response *ListAlbumPhotosResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ListAlbumPhotosResponse
		var err error
		defer close(result)
		response, err = client.ListAlbumPhotos(request)
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

// ListAlbumPhotosRequest is the request struct for api ListAlbumPhotos
type ListAlbumPhotosRequest struct {
	*requests.RpcRequest
	AlbumId   requests.Integer `position:"Query" name:"AlbumId"`
	Direction string           `position:"Query" name:"Direction"`
	Size      requests.Integer `position:"Query" name:"Size"`
	Cursor    string           `position:"Query" name:"Cursor"`
	State     string           `position:"Query" name:"State"`
	StoreName string           `position:"Query" name:"StoreName"`
	LibraryId string           `position:"Query" name:"LibraryId"`
}

// ListAlbumPhotosResponse is the response struct for api ListAlbumPhotos
type ListAlbumPhotosResponse struct {
	*responses.BaseResponse
	Code       string   `json:"Code" xml:"Code"`
	Message    string   `json:"Message" xml:"Message"`
	NextCursor string   `json:"NextCursor" xml:"NextCursor"`
	TotalCount int      `json:"TotalCount" xml:"TotalCount"`
	RequestId  string   `json:"RequestId" xml:"RequestId"`
	Action     string   `json:"Action" xml:"Action"`
	Results    []Result `json:"Results" xml:"Results"`
}

// CreateListAlbumPhotosRequest creates a request to invoke ListAlbumPhotos API
func CreateListAlbumPhotosRequest() (request *ListAlbumPhotosRequest) {
	request = &ListAlbumPhotosRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "ListAlbumPhotos", "cloudphoto", "openAPI")
	return
}

// CreateListAlbumPhotosResponse creates a response to parse from ListAlbumPhotos response
func CreateListAlbumPhotosResponse() (response *ListAlbumPhotosResponse) {
	response = &ListAlbumPhotosResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
