package ecs

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

// GetInstanceScreenshot invokes the ecs.GetInstanceScreenshot API synchronously
// api document: https://help.aliyun.com/api/ecs/getinstancescreenshot.html
func (client *Client) GetInstanceScreenshot(request *GetInstanceScreenshotRequest) (response *GetInstanceScreenshotResponse, err error) {
	response = CreateGetInstanceScreenshotResponse()
	err = client.DoAction(request, response)
	return
}

// GetInstanceScreenshotWithChan invokes the ecs.GetInstanceScreenshot API asynchronously
// api document: https://help.aliyun.com/api/ecs/getinstancescreenshot.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetInstanceScreenshotWithChan(request *GetInstanceScreenshotRequest) (<-chan *GetInstanceScreenshotResponse, <-chan error) {
	responseChan := make(chan *GetInstanceScreenshotResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetInstanceScreenshot(request)
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

// GetInstanceScreenshotWithCallback invokes the ecs.GetInstanceScreenshot API asynchronously
// api document: https://help.aliyun.com/api/ecs/getinstancescreenshot.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetInstanceScreenshotWithCallback(request *GetInstanceScreenshotRequest, callback func(response *GetInstanceScreenshotResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetInstanceScreenshotResponse
		var err error
		defer close(result)
		response, err = client.GetInstanceScreenshot(request)
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

// GetInstanceScreenshotRequest is the request struct for api GetInstanceScreenshot
type GetInstanceScreenshotRequest struct {
	*requests.RpcRequest
}

// GetInstanceScreenshotResponse is the response struct for api GetInstanceScreenshot
type GetInstanceScreenshotResponse struct {
	*responses.BaseResponse
	RequestId  string `json:"RequestId" xml:"RequestId"`
	InstanceId string `json:"InstanceId" xml:"InstanceId"`
	Screenshot string `json:"Screenshot" xml:"Screenshot"`
}

// CreateGetInstanceScreenshotRequest creates a request to invoke GetInstanceScreenshot API
func CreateGetInstanceScreenshotRequest() (request *GetInstanceScreenshotRequest) {
	request = &GetInstanceScreenshotRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "GetInstanceScreenshot", "ecs", "openAPI")
	return
}

// CreateGetInstanceScreenshotResponse creates a response to parse from GetInstanceScreenshot response
func CreateGetInstanceScreenshotResponse() (response *GetInstanceScreenshotResponse) {
	response = &GetInstanceScreenshotResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
