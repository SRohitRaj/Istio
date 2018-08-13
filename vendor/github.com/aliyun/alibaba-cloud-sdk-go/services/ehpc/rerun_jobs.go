package ehpc

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

// RerunJobs invokes the ehpc.RerunJobs API synchronously
// api document: https://help.aliyun.com/api/ehpc/rerunjobs.html
func (client *Client) RerunJobs(request *RerunJobsRequest) (response *RerunJobsResponse, err error) {
	response = CreateRerunJobsResponse()
	err = client.DoAction(request, response)
	return
}

// RerunJobsWithChan invokes the ehpc.RerunJobs API asynchronously
// api document: https://help.aliyun.com/api/ehpc/rerunjobs.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RerunJobsWithChan(request *RerunJobsRequest) (<-chan *RerunJobsResponse, <-chan error) {
	responseChan := make(chan *RerunJobsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RerunJobs(request)
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

// RerunJobsWithCallback invokes the ehpc.RerunJobs API asynchronously
// api document: https://help.aliyun.com/api/ehpc/rerunjobs.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RerunJobsWithCallback(request *RerunJobsRequest, callback func(response *RerunJobsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RerunJobsResponse
		var err error
		defer close(result)
		response, err = client.RerunJobs(request)
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

// RerunJobsRequest is the request struct for api RerunJobs
type RerunJobsRequest struct {
	*requests.RpcRequest
	ClusterId string `position:"Query" name:"ClusterId"`
	Jobs      string `position:"Query" name:"Jobs"`
}

// RerunJobsResponse is the response struct for api RerunJobs
type RerunJobsResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateRerunJobsRequest creates a request to invoke RerunJobs API
func CreateRerunJobsRequest() (request *RerunJobsRequest) {
	request = &RerunJobsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("EHPC", "2017-07-14", "RerunJobs", "ehs", "openAPI")
	return
}

// CreateRerunJobsResponse creates a response to parse from RerunJobs response
func CreateRerunJobsResponse() (response *RerunJobsResponse) {
	response = &RerunJobsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
