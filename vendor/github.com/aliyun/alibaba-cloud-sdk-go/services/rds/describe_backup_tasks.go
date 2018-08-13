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

// DescribeBackupTasks invokes the rds.DescribeBackupTasks API synchronously
// api document: https://help.aliyun.com/api/rds/describebackuptasks.html
func (client *Client) DescribeBackupTasks(request *DescribeBackupTasksRequest) (response *DescribeBackupTasksResponse, err error) {
	response = CreateDescribeBackupTasksResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeBackupTasksWithChan invokes the rds.DescribeBackupTasks API asynchronously
// api document: https://help.aliyun.com/api/rds/describebackuptasks.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeBackupTasksWithChan(request *DescribeBackupTasksRequest) (<-chan *DescribeBackupTasksResponse, <-chan error) {
	responseChan := make(chan *DescribeBackupTasksResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeBackupTasks(request)
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

// DescribeBackupTasksWithCallback invokes the rds.DescribeBackupTasks API asynchronously
// api document: https://help.aliyun.com/api/rds/describebackuptasks.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeBackupTasksWithCallback(request *DescribeBackupTasksRequest, callback func(response *DescribeBackupTasksResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeBackupTasksResponse
		var err error
		defer close(result)
		response, err = client.DescribeBackupTasks(request)
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

// DescribeBackupTasksRequest is the request struct for api DescribeBackupTasks
type DescribeBackupTasksRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	Flag                 string           `position:"Query" name:"Flag"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	BackupJobId          string           `position:"Query" name:"BackupJobId"`
	BackupMode           string           `position:"Query" name:"BackupMode"`
	BackupJobStatus      string           `position:"Query" name:"BackupJobStatus"`
}

// DescribeBackupTasksResponse is the response struct for api DescribeBackupTasks
type DescribeBackupTasksResponse struct {
	*responses.BaseResponse
	RequestId string                     `json:"RequestId" xml:"RequestId"`
	Items     ItemsInDescribeBackupTasks `json:"Items" xml:"Items"`
}

// CreateDescribeBackupTasksRequest creates a request to invoke DescribeBackupTasks API
func CreateDescribeBackupTasksRequest() (request *DescribeBackupTasksRequest) {
	request = &DescribeBackupTasksRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeBackupTasks", "rds", "openAPI")
	return
}

// CreateDescribeBackupTasksResponse creates a response to parse from DescribeBackupTasks response
func CreateDescribeBackupTasksResponse() (response *DescribeBackupTasksResponse) {
	response = &DescribeBackupTasksResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
