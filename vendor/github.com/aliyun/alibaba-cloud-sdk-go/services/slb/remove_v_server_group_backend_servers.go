package slb

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

// RemoveVServerGroupBackendServers invokes the slb.RemoveVServerGroupBackendServers API synchronously
// api document: https://help.aliyun.com/api/slb/removevservergroupbackendservers.html
func (client *Client) RemoveVServerGroupBackendServers(request *RemoveVServerGroupBackendServersRequest) (response *RemoveVServerGroupBackendServersResponse, err error) {
	response = CreateRemoveVServerGroupBackendServersResponse()
	err = client.DoAction(request, response)
	return
}

// RemoveVServerGroupBackendServersWithChan invokes the slb.RemoveVServerGroupBackendServers API asynchronously
// api document: https://help.aliyun.com/api/slb/removevservergroupbackendservers.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveVServerGroupBackendServersWithChan(request *RemoveVServerGroupBackendServersRequest) (<-chan *RemoveVServerGroupBackendServersResponse, <-chan error) {
	responseChan := make(chan *RemoveVServerGroupBackendServersResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemoveVServerGroupBackendServers(request)
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

// RemoveVServerGroupBackendServersWithCallback invokes the slb.RemoveVServerGroupBackendServers API asynchronously
// api document: https://help.aliyun.com/api/slb/removevservergroupbackendservers.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveVServerGroupBackendServersWithCallback(request *RemoveVServerGroupBackendServersRequest, callback func(response *RemoveVServerGroupBackendServersResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemoveVServerGroupBackendServersResponse
		var err error
		defer close(result)
		response, err = client.RemoveVServerGroupBackendServers(request)
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

// RemoveVServerGroupBackendServersRequest is the request struct for api RemoveVServerGroupBackendServers
type RemoveVServerGroupBackendServersRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	AccessKeyId          string           `position:"Query" name:"access_key_id"`
	Tags                 string           `position:"Query" name:"Tags"`
	VServerGroupId       string           `position:"Query" name:"VServerGroupId"`
	BackendServers       string           `position:"Query" name:"BackendServers"`
}

// RemoveVServerGroupBackendServersResponse is the response struct for api RemoveVServerGroupBackendServers
type RemoveVServerGroupBackendServersResponse struct {
	*responses.BaseResponse
	RequestId      string                                           `json:"RequestId" xml:"RequestId"`
	VServerGroupId string                                           `json:"VServerGroupId" xml:"VServerGroupId"`
	BackendServers BackendServersInRemoveVServerGroupBackendServers `json:"BackendServers" xml:"BackendServers"`
}

// CreateRemoveVServerGroupBackendServersRequest creates a request to invoke RemoveVServerGroupBackendServers API
func CreateRemoveVServerGroupBackendServersRequest() (request *RemoveVServerGroupBackendServersRequest) {
	request = &RemoveVServerGroupBackendServersRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Slb", "2014-05-15", "RemoveVServerGroupBackendServers", "slb", "openAPI")
	return
}

// CreateRemoveVServerGroupBackendServersResponse creates a response to parse from RemoveVServerGroupBackendServers response
func CreateRemoveVServerGroupBackendServersResponse() (response *RemoveVServerGroupBackendServersResponse) {
	response = &RemoveVServerGroupBackendServersResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
