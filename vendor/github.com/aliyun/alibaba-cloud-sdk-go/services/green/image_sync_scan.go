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

// ImageSyncScan invokes the green.ImageSyncScan API synchronously
// api document: https://help.aliyun.com/api/green/imagesyncscan.html
func (client *Client) ImageSyncScan(request *ImageSyncScanRequest) (response *ImageSyncScanResponse, err error) {
	response = CreateImageSyncScanResponse()
	err = client.DoAction(request, response)
	return
}

// ImageSyncScanWithChan invokes the green.ImageSyncScan API asynchronously
// api document: https://help.aliyun.com/api/green/imagesyncscan.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ImageSyncScanWithChan(request *ImageSyncScanRequest) (<-chan *ImageSyncScanResponse, <-chan error) {
	responseChan := make(chan *ImageSyncScanResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ImageSyncScan(request)
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

// ImageSyncScanWithCallback invokes the green.ImageSyncScan API asynchronously
// api document: https://help.aliyun.com/api/green/imagesyncscan.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ImageSyncScanWithCallback(request *ImageSyncScanRequest, callback func(response *ImageSyncScanResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ImageSyncScanResponse
		var err error
		defer close(result)
		response, err = client.ImageSyncScan(request)
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

// ImageSyncScanRequest is the request struct for api ImageSyncScan
type ImageSyncScanRequest struct {
	*requests.RoaRequest
	ClientInfo string `position:"Query" name:"ClientInfo"`
}

// ImageSyncScanResponse is the response struct for api ImageSyncScan
type ImageSyncScanResponse struct {
	*responses.BaseResponse
}

// CreateImageSyncScanRequest creates a request to invoke ImageSyncScan API
func CreateImageSyncScanRequest() (request *ImageSyncScanRequest) {
	request = &ImageSyncScanRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Green", "2017-08-25", "ImageSyncScan", "/green/image/scan", "green", "openAPI")
	request.Method = requests.POST
	return
}

// CreateImageSyncScanResponse creates a response to parse from ImageSyncScan response
func CreateImageSyncScanResponse() (response *ImageSyncScanResponse) {
	response = &ImageSyncScanResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
