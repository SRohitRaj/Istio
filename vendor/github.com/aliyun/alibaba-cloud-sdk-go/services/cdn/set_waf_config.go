package cdn

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

// SetWafConfig invokes the cdn.SetWafConfig API synchronously
// api document: https://help.aliyun.com/api/cdn/setwafconfig.html
func (client *Client) SetWafConfig(request *SetWafConfigRequest) (response *SetWafConfigResponse, err error) {
	response = CreateSetWafConfigResponse()
	err = client.DoAction(request, response)
	return
}

// SetWafConfigWithChan invokes the cdn.SetWafConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setwafconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetWafConfigWithChan(request *SetWafConfigRequest) (<-chan *SetWafConfigResponse, <-chan error) {
	responseChan := make(chan *SetWafConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetWafConfig(request)
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

// SetWafConfigWithCallback invokes the cdn.SetWafConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setwafconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetWafConfigWithCallback(request *SetWafConfigRequest, callback func(response *SetWafConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetWafConfigResponse
		var err error
		defer close(result)
		response, err = client.SetWafConfig(request)
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

// SetWafConfigRequest is the request struct for api SetWafConfig
type SetWafConfigRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	Enable        string           `position:"Query" name:"Enable"`
}

// SetWafConfigResponse is the response struct for api SetWafConfig
type SetWafConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetWafConfigRequest creates a request to invoke SetWafConfig API
func CreateSetWafConfigRequest() (request *SetWafConfigRequest) {
	request = &SetWafConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "SetWafConfig", "", "")
	return
}

// CreateSetWafConfigResponse creates a response to parse from SetWafConfig response
func CreateSetWafConfigResponse() (response *SetWafConfigResponse) {
	response = &SetWafConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
