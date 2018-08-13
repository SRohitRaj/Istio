package green

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

// GetFaces invokes the green.GetFaces API synchronously
// api document: https://help.aliyun.com/api/green/getfaces.html
func (client *Client) GetFaces(request *GetFacesRequest) (response *GetFacesResponse, err error) {
	response = CreateGetFacesResponse()
	err = client.DoAction(request, response)
	return
}

// GetFacesWithChan invokes the green.GetFaces API asynchronously
// api document: https://help.aliyun.com/api/green/getfaces.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetFacesWithChan(request *GetFacesRequest) (<-chan *GetFacesResponse, <-chan error) {
	responseChan := make(chan *GetFacesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetFaces(request)
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

// GetFacesWithCallback invokes the green.GetFaces API asynchronously
// api document: https://help.aliyun.com/api/green/getfaces.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetFacesWithCallback(request *GetFacesRequest, callback func(response *GetFacesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetFacesResponse
		var err error
		defer close(result)
		response, err = client.GetFaces(request)
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

// GetFacesRequest is the request struct for api GetFaces
type GetFacesRequest struct {
	*requests.RoaRequest
	ClientInfo string `position:"Query" name:"ClientInfo"`
}

// GetFacesResponse is the response struct for api GetFaces
type GetFacesResponse struct {
	*responses.BaseResponse
}

// CreateGetFacesRequest creates a request to invoke GetFaces API
func CreateGetFacesRequest() (request *GetFacesRequest) {
	request = &GetFacesRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Green", "2017-08-25", "GetFaces", "/green/sface/getFaces", "green", "openAPI")
	request.Method = requests.POST
	return
}

// CreateGetFacesResponse creates a response to parse from GetFaces response
func CreateGetFacesResponse() (response *GetFacesResponse) {
	response = &GetFacesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
