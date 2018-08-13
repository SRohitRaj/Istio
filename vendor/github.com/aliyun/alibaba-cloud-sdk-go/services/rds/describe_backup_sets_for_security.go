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

// DescribeBackupSetsForSecurity invokes the rds.DescribeBackupSetsForSecurity API synchronously
// api document: https://help.aliyun.com/api/rds/describebackupsetsforsecurity.html
func (client *Client) DescribeBackupSetsForSecurity(request *DescribeBackupSetsForSecurityRequest) (response *DescribeBackupSetsForSecurityResponse, err error) {
	response = CreateDescribeBackupSetsForSecurityResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeBackupSetsForSecurityWithChan invokes the rds.DescribeBackupSetsForSecurity API asynchronously
// api document: https://help.aliyun.com/api/rds/describebackupsetsforsecurity.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeBackupSetsForSecurityWithChan(request *DescribeBackupSetsForSecurityRequest) (<-chan *DescribeBackupSetsForSecurityResponse, <-chan error) {
	responseChan := make(chan *DescribeBackupSetsForSecurityResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeBackupSetsForSecurity(request)
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

// DescribeBackupSetsForSecurityWithCallback invokes the rds.DescribeBackupSetsForSecurity API asynchronously
// api document: https://help.aliyun.com/api/rds/describebackupsetsforsecurity.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeBackupSetsForSecurityWithCallback(request *DescribeBackupSetsForSecurityRequest, callback func(response *DescribeBackupSetsForSecurityResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeBackupSetsForSecurityResponse
		var err error
		defer close(result)
		response, err = client.DescribeBackupSetsForSecurity(request)
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

// DescribeBackupSetsForSecurityRequest is the request struct for api DescribeBackupSetsForSecurity
type DescribeBackupSetsForSecurityRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	TargetAliUid         string           `position:"Query" name:"TargetAliUid"`
	TargetAliBid         string           `position:"Query" name:"TargetAliBid"`
	BackupId             string           `position:"Query" name:"BackupId"`
	BackupLocation       string           `position:"Query" name:"BackupLocation"`
	BackupStatus         string           `position:"Query" name:"BackupStatus"`
	BackupMode           string           `position:"Query" name:"BackupMode"`
	StartTime            string           `position:"Query" name:"StartTime"`
	EndTime              string           `position:"Query" name:"EndTime"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// DescribeBackupSetsForSecurityResponse is the response struct for api DescribeBackupSetsForSecurity
type DescribeBackupSetsForSecurityResponse struct {
	*responses.BaseResponse
	RequestId        string                               `json:"RequestId" xml:"RequestId"`
	TotalRecordCount string                               `json:"TotalRecordCount" xml:"TotalRecordCount"`
	PageNumber       string                               `json:"PageNumber" xml:"PageNumber"`
	PageRecordCount  string                               `json:"PageRecordCount" xml:"PageRecordCount"`
	TotalBackupSize  int                                  `json:"TotalBackupSize" xml:"TotalBackupSize"`
	Items            ItemsInDescribeBackupSetsForSecurity `json:"Items" xml:"Items"`
}

// CreateDescribeBackupSetsForSecurityRequest creates a request to invoke DescribeBackupSetsForSecurity API
func CreateDescribeBackupSetsForSecurityRequest() (request *DescribeBackupSetsForSecurityRequest) {
	request = &DescribeBackupSetsForSecurityRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeBackupSetsForSecurity", "rds", "openAPI")
	return
}

// CreateDescribeBackupSetsForSecurityResponse creates a response to parse from DescribeBackupSetsForSecurity response
func CreateDescribeBackupSetsForSecurityResponse() (response *DescribeBackupSetsForSecurityResponse) {
	response = &DescribeBackupSetsForSecurityResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
