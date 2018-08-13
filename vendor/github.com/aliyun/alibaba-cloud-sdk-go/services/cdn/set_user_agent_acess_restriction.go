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

// SetUserAgentAcessRestriction invokes the cdn.SetUserAgentAcessRestriction API synchronously
// api document: https://help.aliyun.com/api/cdn/setuseragentacessrestriction.html
func (client *Client) SetUserAgentAcessRestriction(request *SetUserAgentAcessRestrictionRequest) (response *SetUserAgentAcessRestrictionResponse, err error) {
	response = CreateSetUserAgentAcessRestrictionResponse()
	err = client.DoAction(request, response)
	return
}

// SetUserAgentAcessRestrictionWithChan invokes the cdn.SetUserAgentAcessRestriction API asynchronously
// api document: https://help.aliyun.com/api/cdn/setuseragentacessrestriction.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetUserAgentAcessRestrictionWithChan(request *SetUserAgentAcessRestrictionRequest) (<-chan *SetUserAgentAcessRestrictionResponse, <-chan error) {
	responseChan := make(chan *SetUserAgentAcessRestrictionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetUserAgentAcessRestriction(request)
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

// SetUserAgentAcessRestrictionWithCallback invokes the cdn.SetUserAgentAcessRestriction API asynchronously
// api document: https://help.aliyun.com/api/cdn/setuseragentacessrestriction.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetUserAgentAcessRestrictionWithCallback(request *SetUserAgentAcessRestrictionRequest, callback func(response *SetUserAgentAcessRestrictionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetUserAgentAcessRestrictionResponse
		var err error
		defer close(result)
		response, err = client.SetUserAgentAcessRestriction(request)
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

// SetUserAgentAcessRestrictionRequest is the request struct for api SetUserAgentAcessRestriction
type SetUserAgentAcessRestrictionRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	UserAgent     string           `position:"Query" name:"UserAgent"`
	Type          string           `position:"Query" name:"Type"`
}

// SetUserAgentAcessRestrictionResponse is the response struct for api SetUserAgentAcessRestriction
type SetUserAgentAcessRestrictionResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetUserAgentAcessRestrictionRequest creates a request to invoke SetUserAgentAcessRestriction API
func CreateSetUserAgentAcessRestrictionRequest() (request *SetUserAgentAcessRestrictionRequest) {
	request = &SetUserAgentAcessRestrictionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "SetUserAgentAcessRestriction", "", "")
	return
}

// CreateSetUserAgentAcessRestrictionResponse creates a response to parse from SetUserAgentAcessRestriction response
func CreateSetUserAgentAcessRestrictionResponse() (response *SetUserAgentAcessRestrictionResponse) {
	response = &SetUserAgentAcessRestrictionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
