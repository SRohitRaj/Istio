package rds

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

// DescribeReplicaPerformance invokes the rds.DescribeReplicaPerformance API synchronously
// api document: https://help.aliyun.com/api/rds/describereplicaperformance.html
func (client *Client) DescribeReplicaPerformance(request *DescribeReplicaPerformanceRequest) (response *DescribeReplicaPerformanceResponse, err error) {
	response = CreateDescribeReplicaPerformanceResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeReplicaPerformanceWithChan invokes the rds.DescribeReplicaPerformance API asynchronously
// api document: https://help.aliyun.com/api/rds/describereplicaperformance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeReplicaPerformanceWithChan(request *DescribeReplicaPerformanceRequest) (<-chan *DescribeReplicaPerformanceResponse, <-chan error) {
	responseChan := make(chan *DescribeReplicaPerformanceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeReplicaPerformance(request)
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

// DescribeReplicaPerformanceWithCallback invokes the rds.DescribeReplicaPerformance API asynchronously
// api document: https://help.aliyun.com/api/rds/describereplicaperformance.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeReplicaPerformanceWithCallback(request *DescribeReplicaPerformanceRequest, callback func(response *DescribeReplicaPerformanceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeReplicaPerformanceResponse
		var err error
		defer close(result)
		response, err = client.DescribeReplicaPerformance(request)
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

// DescribeReplicaPerformanceRequest is the request struct for api DescribeReplicaPerformance
type DescribeReplicaPerformanceRequest struct {
	*requests.RpcRequest
	SecurityToken        string           `position:"Query" name:"SecurityToken"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	ReplicaId            string           `position:"Query" name:"ReplicaId"`
	SourceDBInstanceId   string           `position:"Query" name:"SourceDBInstanceId"`
	Key                  string           `position:"Query" name:"Key"`
	StartTime            string           `position:"Query" name:"StartTime"`
	EndTime              string           `position:"Query" name:"EndTime"`
}

// DescribeReplicaPerformanceResponse is the response struct for api DescribeReplicaPerformance
type DescribeReplicaPerformanceResponse struct {
	*responses.BaseResponse
	RequestId       string          `json:"RequestId" xml:"RequestId"`
	StartTime       string          `json:"StartTime" xml:"StartTime"`
	EndTime         string          `json:"EndTime" xml:"EndTime"`
	ReplicaId       string          `json:"ReplicaId" xml:"ReplicaId"`
	PerformanceKeys PerformanceKeys `json:"PerformanceKeys" xml:"PerformanceKeys"`
}

// CreateDescribeReplicaPerformanceRequest creates a request to invoke DescribeReplicaPerformance API
func CreateDescribeReplicaPerformanceRequest() (request *DescribeReplicaPerformanceRequest) {
	request = &DescribeReplicaPerformanceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeReplicaPerformance", "rds", "openAPI")
	return
}

// CreateDescribeReplicaPerformanceResponse creates a response to parse from DescribeReplicaPerformance response
func CreateDescribeReplicaPerformanceResponse() (response *DescribeReplicaPerformanceResponse) {
	response = &DescribeReplicaPerformanceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
