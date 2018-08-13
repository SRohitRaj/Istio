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

// FetchMomentPhotos invokes the cloudphoto.FetchMomentPhotos API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/fetchmomentphotos.html
func (client *Client) FetchMomentPhotos(request *FetchMomentPhotosRequest) (response *FetchMomentPhotosResponse, err error) {
	response = CreateFetchMomentPhotosResponse()
	err = client.DoAction(request, response)
	return
}

// FetchMomentPhotosWithChan invokes the cloudphoto.FetchMomentPhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/fetchmomentphotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) FetchMomentPhotosWithChan(request *FetchMomentPhotosRequest) (<-chan *FetchMomentPhotosResponse, <-chan error) {
	responseChan := make(chan *FetchMomentPhotosResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.FetchMomentPhotos(request)
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

// FetchMomentPhotosWithCallback invokes the cloudphoto.FetchMomentPhotos API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/fetchmomentphotos.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) FetchMomentPhotosWithCallback(request *FetchMomentPhotosRequest, callback func(response *FetchMomentPhotosResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *FetchMomentPhotosResponse
		var err error
		defer close(result)
		response, err = client.FetchMomentPhotos(request)
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

// FetchMomentPhotosRequest is the request struct for api FetchMomentPhotos
type FetchMomentPhotosRequest struct {
	*requests.RpcRequest
	MomentId  requests.Integer `position:"Query" name:"MomentId"`
	OrderBy   string           `position:"Query" name:"OrderBy"`
	Order     string           `position:"Query" name:"Order"`
	Size      requests.Integer `position:"Query" name:"Size"`
	Page      requests.Integer `position:"Query" name:"Page"`
	StoreName string           `position:"Query" name:"StoreName"`
	LibraryId string           `position:"Query" name:"LibraryId"`
}

// FetchMomentPhotosResponse is the response struct for api FetchMomentPhotos
type FetchMomentPhotosResponse struct {
	*responses.BaseResponse
	Code       string  `json:"Code" xml:"Code"`
	Message    string  `json:"Message" xml:"Message"`
	TotalCount int     `json:"TotalCount" xml:"TotalCount"`
	RequestId  string  `json:"RequestId" xml:"RequestId"`
	Action     string  `json:"Action" xml:"Action"`
	Photos     []Photo `json:"Photos" xml:"Photos"`
}

// CreateFetchMomentPhotosRequest creates a request to invoke FetchMomentPhotos API
func CreateFetchMomentPhotosRequest() (request *FetchMomentPhotosRequest) {
	request = &FetchMomentPhotosRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "FetchMomentPhotos", "cloudphoto", "openAPI")
	return
}

// CreateFetchMomentPhotosResponse creates a response to parse from FetchMomentPhotos response
func CreateFetchMomentPhotosResponse() (response *FetchMomentPhotosResponse) {
	response = &FetchMomentPhotosResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
