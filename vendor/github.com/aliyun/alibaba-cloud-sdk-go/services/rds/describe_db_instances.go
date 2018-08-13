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

// DescribeDBInstances invokes the rds.DescribeDBInstances API synchronously
// api document: https://help.aliyun.com/api/rds/describedbinstances.html
func (client *Client) DescribeDBInstances(request *DescribeDBInstancesRequest) (response *DescribeDBInstancesResponse, err error) {
	response = CreateDescribeDBInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDBInstancesWithChan invokes the rds.DescribeDBInstances API asynchronously
// api document: https://help.aliyun.com/api/rds/describedbinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDBInstancesWithChan(request *DescribeDBInstancesRequest) (<-chan *DescribeDBInstancesResponse, <-chan error) {
	responseChan := make(chan *DescribeDBInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDBInstances(request)
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

// DescribeDBInstancesWithCallback invokes the rds.DescribeDBInstances API asynchronously
// api document: https://help.aliyun.com/api/rds/describedbinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDBInstancesWithCallback(request *DescribeDBInstancesRequest, callback func(response *DescribeDBInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDBInstancesResponse
		var err error
		defer close(result)
		response, err = client.DescribeDBInstances(request)
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

// DescribeDBInstancesRequest is the request struct for api DescribeDBInstances
type DescribeDBInstancesRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ClientToken          string           `position:"Query" name:"ClientToken"`
	ProxyId              string           `position:"Query" name:"proxyId"`
	Engine               string           `position:"Query" name:"Engine"`
	DBInstanceStatus     string           `position:"Query" name:"DBInstanceStatus"`
	SearchKey            string           `position:"Query" name:"SearchKey"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	DBInstanceType       string           `position:"Query" name:"DBInstanceType"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	InstanceNetworkType  string           `position:"Query" name:"InstanceNetworkType"`
	VpcId                string           `position:"Query" name:"VpcId"`
	VSwitchId            string           `position:"Query" name:"VSwitchId"`
	ConnectionMode       string           `position:"Query" name:"ConnectionMode"`
	Tag1Key              string           `position:"Query" name:"Tag.1.key"`
	Tag2Key              string           `position:"Query" name:"Tag.2.key"`
	Tag3Key              string           `position:"Query" name:"Tag.3.key"`
	Tag4Key              string           `position:"Query" name:"Tag.4.key"`
	Tag5Key              string           `position:"Query" name:"Tag.5.key"`
	Tag1Value            string           `position:"Query" name:"Tag.1.value"`
	Tag2Value            string           `position:"Query" name:"Tag.2.value"`
	Tag3Value            string           `position:"Query" name:"Tag.3.value"`
	Tag4Value            string           `position:"Query" name:"Tag.4.value"`
	Tag5Value            string           `position:"Query" name:"Tag.5.value"`
	Tags                 string           `position:"Query" name:"Tags"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// DescribeDBInstancesResponse is the response struct for api DescribeDBInstances
type DescribeDBInstancesResponse struct {
	*responses.BaseResponse
	RequestId        string                     `json:"RequestId" xml:"RequestId"`
	PageNumber       int                        `json:"PageNumber" xml:"PageNumber"`
	TotalRecordCount int                        `json:"TotalRecordCount" xml:"TotalRecordCount"`
	PageRecordCount  int                        `json:"PageRecordCount" xml:"PageRecordCount"`
	Items            ItemsInDescribeDBInstances `json:"Items" xml:"Items"`
}

// CreateDescribeDBInstancesRequest creates a request to invoke DescribeDBInstances API
func CreateDescribeDBInstancesRequest() (request *DescribeDBInstancesRequest) {
	request = &DescribeDBInstancesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeDBInstances", "rds", "openAPI")
	return
}

// CreateDescribeDBInstancesResponse creates a response to parse from DescribeDBInstances response
func CreateDescribeDBInstancesResponse() (response *DescribeDBInstancesResponse) {
	response = &DescribeDBInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
