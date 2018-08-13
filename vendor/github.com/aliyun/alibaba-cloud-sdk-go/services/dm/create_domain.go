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

// CreateDomain invokes the dm.CreateDomain API synchronously
// api document: https://help.aliyun.com/api/dm/createdomain.html
func (client *Client) CreateDomain(request *CreateDomainRequest) (response *CreateDomainResponse, err error) {
	response = CreateCreateDomainResponse()
	err = client.DoAction(request, response)
	return
}

// CreateDomainWithChan invokes the dm.CreateDomain API asynchronously
// api document: https://help.aliyun.com/api/dm/createdomain.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateDomainWithChan(request *CreateDomainRequest) (<-chan *CreateDomainResponse, <-chan error) {
	responseChan := make(chan *CreateDomainResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateDomain(request)
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

// CreateDomainWithCallback invokes the dm.CreateDomain API asynchronously
// api document: https://help.aliyun.com/api/dm/createdomain.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateDomainWithCallback(request *CreateDomainRequest, callback func(response *CreateDomainResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateDomainResponse
		var err error
		defer close(result)
		response, err = client.CreateDomain(request)
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

// CreateDomainRequest is the request struct for api CreateDomain
type CreateDomainRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DomainName           string           `position:"Query" name:"DomainName"`
}

// CreateDomainResponse is the response struct for api CreateDomain
type CreateDomainResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateCreateDomainRequest creates a request to invoke CreateDomain API
func CreateCreateDomainRequest() (request *CreateDomainRequest) {
	request = &CreateDomainRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dm", "2015-11-23", "CreateDomain", "", "")
	return
}

// CreateCreateDomainResponse creates a response to parse from CreateDomain response
func CreateCreateDomainResponse() (response *CreateDomainResponse) {
	response = &CreateDomainResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
