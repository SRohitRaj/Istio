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

// DescribeLoadBalancers invokes the slb.DescribeLoadBalancers API synchronously
// api document: https://help.aliyun.com/api/slb/describeloadbalancers.html
func (client *Client) DescribeLoadBalancers(request *DescribeLoadBalancersRequest) (response *DescribeLoadBalancersResponse, err error) {
	response = CreateDescribeLoadBalancersResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLoadBalancersWithChan invokes the slb.DescribeLoadBalancers API asynchronously
// api document: https://help.aliyun.com/api/slb/describeloadbalancers.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLoadBalancersWithChan(request *DescribeLoadBalancersRequest) (<-chan *DescribeLoadBalancersResponse, <-chan error) {
	responseChan := make(chan *DescribeLoadBalancersResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLoadBalancers(request)
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

// DescribeLoadBalancersWithCallback invokes the slb.DescribeLoadBalancers API asynchronously
// api document: https://help.aliyun.com/api/slb/describeloadbalancers.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLoadBalancersWithCallback(request *DescribeLoadBalancersRequest, callback func(response *DescribeLoadBalancersResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLoadBalancersResponse
		var err error
		defer close(result)
		response, err = client.DescribeLoadBalancers(request)
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

// DescribeLoadBalancersRequest is the request struct for api DescribeLoadBalancers
type DescribeLoadBalancersRequest struct {
	*requests.RpcRequest
	OwnerId               requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount  string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId       requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ServerId              string           `position:"Query" name:"ServerId"`
	LoadBalancerId        string           `position:"Query" name:"LoadBalancerId"`
	LoadBalancerName      string           `position:"Query" name:"LoadBalancerName"`
	ServerIntranetAddress string           `position:"Query" name:"ServerIntranetAddress"`
	AddressType           string           `position:"Query" name:"AddressType"`
	InternetChargeType    string           `position:"Query" name:"InternetChargeType"`
	VpcId                 string           `position:"Query" name:"VpcId"`
	VSwitchId             string           `position:"Query" name:"VSwitchId"`
	NetworkType           string           `position:"Query" name:"NetworkType"`
	Address               string           `position:"Query" name:"Address"`
	MasterZoneId          string           `position:"Query" name:"MasterZoneId"`
	SlaveZoneId           string           `position:"Query" name:"SlaveZoneId"`
	OwnerAccount          string           `position:"Query" name:"OwnerAccount"`
	AccessKeyId           string           `position:"Query" name:"access_key_id"`
	Tags                  string           `position:"Query" name:"Tags"`
	PayType               string           `position:"Query" name:"PayType"`
	ResourceGroupId       string           `position:"Query" name:"ResourceGroupId"`
	PageNumber            requests.Integer `position:"Query" name:"PageNumber"`
	PageSize              requests.Integer `position:"Query" name:"PageSize"`
}

// DescribeLoadBalancersResponse is the response struct for api DescribeLoadBalancers
type DescribeLoadBalancersResponse struct {
	*responses.BaseResponse
	RequestId     string                               `json:"RequestId" xml:"RequestId"`
	PageNumber    int                                  `json:"PageNumber" xml:"PageNumber"`
	PageSize      int                                  `json:"PageSize" xml:"PageSize"`
	TotalCount    int                                  `json:"TotalCount" xml:"TotalCount"`
	LoadBalancers LoadBalancersInDescribeLoadBalancers `json:"LoadBalancers" xml:"LoadBalancers"`
}

// CreateDescribeLoadBalancersRequest creates a request to invoke DescribeLoadBalancers API
func CreateDescribeLoadBalancersRequest() (request *DescribeLoadBalancersRequest) {
	request = &DescribeLoadBalancersRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Slb", "2014-05-15", "DescribeLoadBalancers", "slb", "openAPI")
	return
}

// CreateDescribeLoadBalancersResponse creates a response to parse from DescribeLoadBalancers response
func CreateDescribeLoadBalancersResponse() (response *DescribeLoadBalancersResponse) {
	response = &DescribeLoadBalancersResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
