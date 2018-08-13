package cloudapi

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

// ModifySignature invokes the cloudapi.ModifySignature API synchronously
// api document: https://help.aliyun.com/api/cloudapi/modifysignature.html
func (client *Client) ModifySignature(request *ModifySignatureRequest) (response *ModifySignatureResponse, err error) {
	response = CreateModifySignatureResponse()
	err = client.DoAction(request, response)
	return
}

// ModifySignatureWithChan invokes the cloudapi.ModifySignature API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/modifysignature.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifySignatureWithChan(request *ModifySignatureRequest) (<-chan *ModifySignatureResponse, <-chan error) {
	responseChan := make(chan *ModifySignatureResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifySignature(request)
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

// ModifySignatureWithCallback invokes the cloudapi.ModifySignature API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/modifysignature.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifySignatureWithCallback(request *ModifySignatureRequest, callback func(response *ModifySignatureResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifySignatureResponse
		var err error
		defer close(result)
		response, err = client.ModifySignature(request)
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

// ModifySignatureRequest is the request struct for api ModifySignature
type ModifySignatureRequest struct {
	*requests.RpcRequest
	SignatureId     string `position:"Query" name:"SignatureId"`
	SignatureName   string `position:"Query" name:"SignatureName"`
	SignatureKey    string `position:"Query" name:"SignatureKey"`
	SignatureSecret string `position:"Query" name:"SignatureSecret"`
}

// ModifySignatureResponse is the response struct for api ModifySignature
type ModifySignatureResponse struct {
	*responses.BaseResponse
	RequestId     string `json:"RequestId" xml:"RequestId"`
	SignatureId   string `json:"SignatureId" xml:"SignatureId"`
	SignatureName string `json:"SignatureName" xml:"SignatureName"`
}

// CreateModifySignatureRequest creates a request to invoke ModifySignature API
func CreateModifySignatureRequest() (request *ModifySignatureRequest) {
	request = &ModifySignatureRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "ModifySignature", "apigateway", "openAPI")
	return
}

// CreateModifySignatureResponse creates a response to parse from ModifySignature response
func CreateModifySignatureResponse() (response *ModifySignatureResponse) {
	response = &ModifySignatureResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
