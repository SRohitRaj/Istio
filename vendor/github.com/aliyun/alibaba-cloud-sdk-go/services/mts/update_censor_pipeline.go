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

// UpdateCensorPipeline invokes the mts.UpdateCensorPipeline API synchronously
// api document: https://help.aliyun.com/api/mts/updatecensorpipeline.html
func (client *Client) UpdateCensorPipeline(request *UpdateCensorPipelineRequest) (response *UpdateCensorPipelineResponse, err error) {
	response = CreateUpdateCensorPipelineResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateCensorPipelineWithChan invokes the mts.UpdateCensorPipeline API asynchronously
// api document: https://help.aliyun.com/api/mts/updatecensorpipeline.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateCensorPipelineWithChan(request *UpdateCensorPipelineRequest) (<-chan *UpdateCensorPipelineResponse, <-chan error) {
	responseChan := make(chan *UpdateCensorPipelineResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateCensorPipeline(request)
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

// UpdateCensorPipelineWithCallback invokes the mts.UpdateCensorPipeline API asynchronously
// api document: https://help.aliyun.com/api/mts/updatecensorpipeline.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateCensorPipelineWithCallback(request *UpdateCensorPipelineRequest, callback func(response *UpdateCensorPipelineResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateCensorPipelineResponse
		var err error
		defer close(result)
		response, err = client.UpdateCensorPipeline(request)
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

// UpdateCensorPipelineRequest is the request struct for api UpdateCensorPipeline
type UpdateCensorPipelineRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PipelineId           string           `position:"Query" name:"PipelineId"`
	Name                 string           `position:"Query" name:"Name"`
	State                string           `position:"Query" name:"State"`
	Priority             requests.Integer `position:"Query" name:"Priority"`
	NotifyConfig         string           `position:"Query" name:"NotifyConfig"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// UpdateCensorPipelineResponse is the response struct for api UpdateCensorPipeline
type UpdateCensorPipelineResponse struct {
	*responses.BaseResponse
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Pipeline  Pipeline `json:"Pipeline" xml:"Pipeline"`
}

// CreateUpdateCensorPipelineRequest creates a request to invoke UpdateCensorPipeline API
func CreateUpdateCensorPipelineRequest() (request *UpdateCensorPipelineRequest) {
	request = &UpdateCensorPipelineRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "UpdateCensorPipeline", "mts", "openAPI")
	return
}

// CreateUpdateCensorPipelineResponse creates a response to parse from UpdateCensorPipeline response
func CreateUpdateCensorPipelineResponse() (response *UpdateCensorPipelineResponse) {
	response = &UpdateCensorPipelineResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
