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

// DescribeDBInstancesByExpireTime invokes the rds.DescribeDBInstancesByExpireTime API synchronously
// api document: https://help.aliyun.com/api/rds/describedbinstancesbyexpiretime.html
func (client *Client) DescribeDBInstancesByExpireTime(request *DescribeDBInstancesByExpireTimeRequest) (response *DescribeDBInstancesByExpireTimeResponse, err error) {
	response = CreateDescribeDBInstancesByExpireTimeResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDBInstancesByExpireTimeWithChan invokes the rds.DescribeDBInstancesByExpireTime API asynchronously
// api document: https://help.aliyun.com/api/rds/describedbinstancesbyexpiretime.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDBInstancesByExpireTimeWithChan(request *DescribeDBInstancesByExpireTimeRequest) (<-chan *DescribeDBInstancesByExpireTimeResponse, <-chan error) {
	responseChan := make(chan *DescribeDBInstancesByExpireTimeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDBInstancesByExpireTime(request)
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

// DescribeDBInstancesByExpireTimeWithCallback invokes the rds.DescribeDBInstancesByExpireTime API asynchronously
// api document: https://help.aliyun.com/api/rds/describedbinstancesbyexpiretime.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDBInstancesByExpireTimeWithCallback(request *DescribeDBInstancesByExpireTimeRequest, callback func(response *DescribeDBInstancesByExpireTimeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDBInstancesByExpireTimeResponse
		var err error
		defer close(result)
		response, err = client.DescribeDBInstancesByExpireTime(request)
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

// DescribeDBInstancesByExpireTimeRequest is the request struct for api DescribeDBInstancesByExpireTime
type DescribeDBInstancesByExpireTimeRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ProxyId              string           `position:"Query" name:"proxyId"`
	ExpirePeriod         requests.Integer `position:"Query" name:"ExpirePeriod"`
	Expired              requests.Boolean `position:"Query" name:"Expired"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	Tags                 string           `position:"Query" name:"Tags"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// DescribeDBInstancesByExpireTimeResponse is the response struct for api DescribeDBInstancesByExpireTime
type DescribeDBInstancesByExpireTimeResponse struct {
	*responses.BaseResponse
	RequestId        string                                 `json:"RequestId" xml:"RequestId"`
	PageNumber       int                                    `json:"PageNumber" xml:"PageNumber"`
	TotalRecordCount int                                    `json:"TotalRecordCount" xml:"TotalRecordCount"`
	PageRecordCount  int                                    `json:"PageRecordCount" xml:"PageRecordCount"`
	Items            ItemsInDescribeDBInstancesByExpireTime `json:"Items" xml:"Items"`
}

// CreateDescribeDBInstancesByExpireTimeRequest creates a request to invoke DescribeDBInstancesByExpireTime API
func CreateDescribeDBInstancesByExpireTimeRequest() (request *DescribeDBInstancesByExpireTimeRequest) {
	request = &DescribeDBInstancesByExpireTimeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeDBInstancesByExpireTime", "rds", "openAPI")
	return
}

// CreateDescribeDBInstancesByExpireTimeResponse creates a response to parse from DescribeDBInstancesByExpireTime response
func CreateDescribeDBInstancesByExpireTimeResponse() (response *DescribeDBInstancesByExpireTimeResponse) {
	response = &DescribeDBInstancesByExpireTimeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
