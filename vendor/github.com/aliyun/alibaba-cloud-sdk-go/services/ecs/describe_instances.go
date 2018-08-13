package ecs

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

// DescribeInstances invokes the ecs.DescribeInstances API synchronously
// api document: https://help.aliyun.com/api/ecs/describeinstances.html
func (client *Client) DescribeInstances(request *DescribeInstancesRequest) (response *DescribeInstancesResponse, err error) {
	response = CreateDescribeInstancesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeInstancesWithChan invokes the ecs.DescribeInstances API asynchronously
// api document: https://help.aliyun.com/api/ecs/describeinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeInstancesWithChan(request *DescribeInstancesRequest) (<-chan *DescribeInstancesResponse, <-chan error) {
	responseChan := make(chan *DescribeInstancesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeInstances(request)
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

// DescribeInstancesWithCallback invokes the ecs.DescribeInstances API asynchronously
// api document: https://help.aliyun.com/api/ecs/describeinstances.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeInstancesWithCallback(request *DescribeInstancesRequest, callback func(response *DescribeInstancesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeInstancesResponse
		var err error
		defer close(result)
		response, err = client.DescribeInstances(request)
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

// DescribeInstancesRequest is the request struct for api DescribeInstances
type DescribeInstancesRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	VpcId                string           `position:"Query" name:"VpcId"`
	VSwitchId            string           `position:"Query" name:"VSwitchId"`
	ZoneId               string           `position:"Query" name:"ZoneId"`
	InstanceNetworkType  string           `position:"Query" name:"InstanceNetworkType"`
	SecurityGroupId      string           `position:"Query" name:"SecurityGroupId"`
	InstanceIds          string           `position:"Query" name:"InstanceIds"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	InnerIpAddresses     string           `position:"Query" name:"InnerIpAddresses"`
	PrivateIpAddresses   string           `position:"Query" name:"PrivateIpAddresses"`
	PublicIpAddresses    string           `position:"Query" name:"PublicIpAddresses"`
	EipAddresses         string           `position:"Query" name:"EipAddresses"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	InstanceChargeType   string           `position:"Query" name:"InstanceChargeType"`
	InternetChargeType   string           `position:"Query" name:"InternetChargeType"`
	InstanceName         string           `position:"Query" name:"InstanceName"`
	ImageId              string           `position:"Query" name:"ImageId"`
	Status               string           `position:"Query" name:"Status"`
	LockReason           string           `position:"Query" name:"LockReason"`
	Filter1Key           string           `position:"Query" name:"Filter.1.Key"`
	Filter2Key           string           `position:"Query" name:"Filter.2.Key"`
	Filter3Key           string           `position:"Query" name:"Filter.3.Key"`
	Filter4Key           string           `position:"Query" name:"Filter.4.Key"`
	Filter1Value         string           `position:"Query" name:"Filter.1.Value"`
	Filter2Value         string           `position:"Query" name:"Filter.2.Value"`
	Filter3Value         string           `position:"Query" name:"Filter.3.Value"`
	Filter4Value         string           `position:"Query" name:"Filter.4.Value"`
	DeviceAvailable      requests.Boolean `position:"Query" name:"DeviceAvailable"`
	IoOptimized          requests.Boolean `position:"Query" name:"IoOptimized"`
	Tag1Key              string           `position:"Query" name:"Tag.1.Key"`
	Tag2Key              string           `position:"Query" name:"Tag.2.Key"`
	Tag3Key              string           `position:"Query" name:"Tag.3.Key"`
	Tag4Key              string           `position:"Query" name:"Tag.4.Key"`
	Tag5Key              string           `position:"Query" name:"Tag.5.Key"`
	Tag1Value            string           `position:"Query" name:"Tag.1.Value"`
	Tag2Value            string           `position:"Query" name:"Tag.2.Value"`
	Tag3Value            string           `position:"Query" name:"Tag.3.Value"`
	Tag4Value            string           `position:"Query" name:"Tag.4.Value"`
	Tag5Value            string           `position:"Query" name:"Tag.5.Value"`
	InstanceType         string           `position:"Query" name:"InstanceType"`
	InstanceTypeFamily   string           `position:"Query" name:"InstanceTypeFamily"`
	KeyPairName          string           `position:"Query" name:"KeyPairName"`
	ResourceGroupId      string           `position:"Query" name:"ResourceGroupId"`
	HpcClusterId         string           `position:"Query" name:"HpcClusterId"`
	RdmaIpAddresses      string           `position:"Query" name:"RdmaIpAddresses"`
	DryRun               requests.Boolean `position:"Query" name:"DryRun"`
}

// DescribeInstancesResponse is the response struct for api DescribeInstances
type DescribeInstancesResponse struct {
	*responses.BaseResponse
	RequestId  string    `json:"RequestId" xml:"RequestId"`
	TotalCount int       `json:"TotalCount" xml:"TotalCount"`
	PageNumber int       `json:"PageNumber" xml:"PageNumber"`
	PageSize   int       `json:"PageSize" xml:"PageSize"`
	Instances  Instances `json:"Instances" xml:"Instances"`
}

// CreateDescribeInstancesRequest creates a request to invoke DescribeInstances API
func CreateDescribeInstancesRequest() (request *DescribeInstancesRequest) {
	request = &DescribeInstancesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ecs", "2014-05-26", "DescribeInstances", "ecs", "openAPI")
	return
}

// CreateDescribeInstancesResponse creates a response to parse from DescribeInstances response
func CreateDescribeInstancesResponse() (response *DescribeInstancesResponse) {
	response = &DescribeInstancesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
