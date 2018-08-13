package cs

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

// ScaleInCluster invokes the cs.ScaleInCluster API synchronously
// api document: https://help.aliyun.com/api/cs/scaleincluster.html
func (client *Client) ScaleInCluster(request *ScaleInClusterRequest) (response *ScaleInClusterResponse, err error) {
	response = CreateScaleInClusterResponse()
	err = client.DoAction(request, response)
	return
}

// ScaleInClusterWithChan invokes the cs.ScaleInCluster API asynchronously
// api document: https://help.aliyun.com/api/cs/scaleincluster.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ScaleInClusterWithChan(request *ScaleInClusterRequest) (<-chan *ScaleInClusterResponse, <-chan error) {
	responseChan := make(chan *ScaleInClusterResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ScaleInCluster(request)
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

// ScaleInClusterWithCallback invokes the cs.ScaleInCluster API asynchronously
// api document: https://help.aliyun.com/api/cs/scaleincluster.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ScaleInClusterWithCallback(request *ScaleInClusterRequest, callback func(response *ScaleInClusterResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ScaleInClusterResponse
		var err error
		defer close(result)
		response, err = client.ScaleInCluster(request)
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

// ScaleInClusterRequest is the request struct for api ScaleInCluster
type ScaleInClusterRequest struct {
	*requests.RoaRequest
	ClusterId string `position:"Path" name:"ClusterId"`
}

// ScaleInClusterResponse is the response struct for api ScaleInCluster
type ScaleInClusterResponse struct {
	*responses.BaseResponse
}

// CreateScaleInClusterRequest creates a request to invoke ScaleInCluster API
func CreateScaleInClusterRequest() (request *ScaleInClusterRequest) {
	request = &ScaleInClusterRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("CS", "2015-12-15", "ScaleInCluster", "/clusters/[ClusterId]/scalein", "", "")
	request.Method = requests.POST
	return
}

// CreateScaleInClusterResponse creates a response to parse from ScaleInCluster response
func CreateScaleInClusterResponse() (response *ScaleInClusterResponse) {
	response = &ScaleInClusterResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
