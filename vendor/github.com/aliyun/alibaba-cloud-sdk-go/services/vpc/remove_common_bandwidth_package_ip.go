package vpc

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

// RemoveCommonBandwidthPackageIp invokes the vpc.RemoveCommonBandwidthPackageIp API synchronously
// api document: https://help.aliyun.com/api/vpc/removecommonbandwidthpackageip.html
func (client *Client) RemoveCommonBandwidthPackageIp(request *RemoveCommonBandwidthPackageIpRequest) (response *RemoveCommonBandwidthPackageIpResponse, err error) {
	response = CreateRemoveCommonBandwidthPackageIpResponse()
	err = client.DoAction(request, response)
	return
}

// RemoveCommonBandwidthPackageIpWithChan invokes the vpc.RemoveCommonBandwidthPackageIp API asynchronously
// api document: https://help.aliyun.com/api/vpc/removecommonbandwidthpackageip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveCommonBandwidthPackageIpWithChan(request *RemoveCommonBandwidthPackageIpRequest) (<-chan *RemoveCommonBandwidthPackageIpResponse, <-chan error) {
	responseChan := make(chan *RemoveCommonBandwidthPackageIpResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.RemoveCommonBandwidthPackageIp(request)
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

// RemoveCommonBandwidthPackageIpWithCallback invokes the vpc.RemoveCommonBandwidthPackageIp API asynchronously
// api document: https://help.aliyun.com/api/vpc/removecommonbandwidthpackageip.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) RemoveCommonBandwidthPackageIpWithCallback(request *RemoveCommonBandwidthPackageIpRequest, callback func(response *RemoveCommonBandwidthPackageIpResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *RemoveCommonBandwidthPackageIpResponse
		var err error
		defer close(result)
		response, err = client.RemoveCommonBandwidthPackageIp(request)
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

// RemoveCommonBandwidthPackageIpRequest is the request struct for api RemoveCommonBandwidthPackageIp
type RemoveCommonBandwidthPackageIpRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	BandwidthPackageId   string           `position:"Query" name:"BandwidthPackageId"`
	IpInstanceId         string           `position:"Query" name:"IpInstanceId"`
}

// RemoveCommonBandwidthPackageIpResponse is the response struct for api RemoveCommonBandwidthPackageIp
type RemoveCommonBandwidthPackageIpResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateRemoveCommonBandwidthPackageIpRequest creates a request to invoke RemoveCommonBandwidthPackageIp API
func CreateRemoveCommonBandwidthPackageIpRequest() (request *RemoveCommonBandwidthPackageIpRequest) {
	request = &RemoveCommonBandwidthPackageIpRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "RemoveCommonBandwidthPackageIp", "vpc", "openAPI")
	return
}

// CreateRemoveCommonBandwidthPackageIpResponse creates a response to parse from RemoveCommonBandwidthPackageIp response
func CreateRemoveCommonBandwidthPackageIpResponse() (response *RemoveCommonBandwidthPackageIpResponse) {
	response = &RemoveCommonBandwidthPackageIpResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
