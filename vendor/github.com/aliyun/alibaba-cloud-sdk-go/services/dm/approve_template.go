package dm

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

// ApproveTemplate invokes the dm.ApproveTemplate API synchronously
// api document: https://help.aliyun.com/api/dm/approvetemplate.html
func (client *Client) ApproveTemplate(request *ApproveTemplateRequest) (response *ApproveTemplateResponse, err error) {
	response = CreateApproveTemplateResponse()
	err = client.DoAction(request, response)
	return
}

// ApproveTemplateWithChan invokes the dm.ApproveTemplate API asynchronously
// api document: https://help.aliyun.com/api/dm/approvetemplate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ApproveTemplateWithChan(request *ApproveTemplateRequest) (<-chan *ApproveTemplateResponse, <-chan error) {
	responseChan := make(chan *ApproveTemplateResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ApproveTemplate(request)
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

// ApproveTemplateWithCallback invokes the dm.ApproveTemplate API asynchronously
// api document: https://help.aliyun.com/api/dm/approvetemplate.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ApproveTemplateWithCallback(request *ApproveTemplateRequest, callback func(response *ApproveTemplateResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ApproveTemplateResponse
		var err error
		defer close(result)
		response, err = client.ApproveTemplate(request)
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

// ApproveTemplateRequest is the request struct for api ApproveTemplate
type ApproveTemplateRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	TemplateId           requests.Integer `position:"Query" name:"TemplateId"`
	FromType             requests.Integer `position:"Query" name:"FromType"`
}

// ApproveTemplateResponse is the response struct for api ApproveTemplate
type ApproveTemplateResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateApproveTemplateRequest creates a request to invoke ApproveTemplate API
func CreateApproveTemplateRequest() (request *ApproveTemplateRequest) {
	request = &ApproveTemplateRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dm", "2015-11-23", "ApproveTemplate", "", "")
	return
}

// CreateApproveTemplateResponse creates a response to parse from ApproveTemplate response
func CreateApproveTemplateResponse() (response *ApproveTemplateResponse) {
	response = &ApproveTemplateResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
