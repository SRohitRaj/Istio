package mts

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

// GetPackage invokes the mts.GetPackage API synchronously
// api document: https://help.aliyun.com/api/mts/getpackage.html
func (client *Client) GetPackage(request *GetPackageRequest) (response *GetPackageResponse, err error) {
	response = CreateGetPackageResponse()
	err = client.DoAction(request, response)
	return
}

// GetPackageWithChan invokes the mts.GetPackage API asynchronously
// api document: https://help.aliyun.com/api/mts/getpackage.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetPackageWithChan(request *GetPackageRequest) (<-chan *GetPackageResponse, <-chan error) {
	responseChan := make(chan *GetPackageResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetPackage(request)
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

// GetPackageWithCallback invokes the mts.GetPackage API asynchronously
// api document: https://help.aliyun.com/api/mts/getpackage.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetPackageWithCallback(request *GetPackageRequest, callback func(response *GetPackageResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetPackageResponse
		var err error
		defer close(result)
		response, err = client.GetPackage(request)
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

// GetPackageRequest is the request struct for api GetPackage
type GetPackageRequest struct {
	*requests.RpcRequest
	OwnerId              string `position:"Query" name:"OwnerId"`
	ResourceOwnerId      string `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string `position:"Query" name:"OwnerAccount"`
	Data                 string `position:"Query" name:"Data"`
}

// GetPackageResponse is the response struct for api GetPackage
type GetPackageResponse struct {
	*responses.BaseResponse
	RequestId   string `json:"RequestId" xml:"RequestId"`
	CertPackage string `json:"CertPackage" xml:"CertPackage"`
}

// CreateGetPackageRequest creates a request to invoke GetPackage API
func CreateGetPackageRequest() (request *GetPackageRequest) {
	request = &GetPackageRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "GetPackage", "mts", "openAPI")
	return
}

// CreateGetPackageResponse creates a response to parse from GetPackage response
func CreateGetPackageResponse() (response *GetPackageResponse) {
	response = &GetPackageResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
