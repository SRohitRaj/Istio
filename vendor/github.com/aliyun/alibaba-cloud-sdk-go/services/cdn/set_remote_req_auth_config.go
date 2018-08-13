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

// SetRemoteReqAuthConfig invokes the cdn.SetRemoteReqAuthConfig API synchronously
// api document: https://help.aliyun.com/api/cdn/setremotereqauthconfig.html
func (client *Client) SetRemoteReqAuthConfig(request *SetRemoteReqAuthConfigRequest) (response *SetRemoteReqAuthConfigResponse, err error) {
	response = CreateSetRemoteReqAuthConfigResponse()
	err = client.DoAction(request, response)
	return
}

// SetRemoteReqAuthConfigWithChan invokes the cdn.SetRemoteReqAuthConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setremotereqauthconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetRemoteReqAuthConfigWithChan(request *SetRemoteReqAuthConfigRequest) (<-chan *SetRemoteReqAuthConfigResponse, <-chan error) {
	responseChan := make(chan *SetRemoteReqAuthConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetRemoteReqAuthConfig(request)
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

// SetRemoteReqAuthConfigWithCallback invokes the cdn.SetRemoteReqAuthConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setremotereqauthconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetRemoteReqAuthConfigWithCallback(request *SetRemoteReqAuthConfigRequest, callback func(response *SetRemoteReqAuthConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetRemoteReqAuthConfigResponse
		var err error
		defer close(result)
		response, err = client.SetRemoteReqAuthConfig(request)
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

// SetRemoteReqAuthConfigRequest is the request struct for api SetRemoteReqAuthConfig
type SetRemoteReqAuthConfigRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	AuthType      string           `position:"Query" name:"AuthType"`
	TimeOut       string           `position:"Query" name:"TimeOut"`
	AuthAddr      string           `position:"Query" name:"AuthAddr"`
	AuthCrash     string           `position:"Query" name:"AuthCrash"`
	AuthEnable    string           `position:"Query" name:"AuthEnable"`
	AuthProvider  string           `position:"Query" name:"AuthProvider"`
	AuthPath      string           `position:"Query" name:"AuthPath"`
	AuthFileType  string           `position:"Query" name:"AuthFileType"`
}

// SetRemoteReqAuthConfigResponse is the response struct for api SetRemoteReqAuthConfig
type SetRemoteReqAuthConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetRemoteReqAuthConfigRequest creates a request to invoke SetRemoteReqAuthConfig API
func CreateSetRemoteReqAuthConfigRequest() (request *SetRemoteReqAuthConfigRequest) {
	request = &SetRemoteReqAuthConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "SetRemoteReqAuthConfig", "", "")
	return
}

// CreateSetRemoteReqAuthConfigResponse creates a response to parse from SetRemoteReqAuthConfig response
func CreateSetRemoteReqAuthConfigResponse() (response *SetRemoteReqAuthConfigResponse) {
	response = &SetRemoteReqAuthConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
