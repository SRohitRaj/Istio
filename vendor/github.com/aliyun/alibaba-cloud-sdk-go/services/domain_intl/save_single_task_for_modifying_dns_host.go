package domain_intl

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

// SaveSingleTaskForModifyingDnsHost invokes the domain_intl.SaveSingleTaskForModifyingDnsHost API synchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskformodifyingdnshost.html
func (client *Client) SaveSingleTaskForModifyingDnsHost(request *SaveSingleTaskForModifyingDnsHostRequest) (response *SaveSingleTaskForModifyingDnsHostResponse, err error) {
	response = CreateSaveSingleTaskForModifyingDnsHostResponse()
	err = client.DoAction(request, response)
	return
}

// SaveSingleTaskForModifyingDnsHostWithChan invokes the domain_intl.SaveSingleTaskForModifyingDnsHost API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskformodifyingdnshost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForModifyingDnsHostWithChan(request *SaveSingleTaskForModifyingDnsHostRequest) (<-chan *SaveSingleTaskForModifyingDnsHostResponse, <-chan error) {
	responseChan := make(chan *SaveSingleTaskForModifyingDnsHostResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveSingleTaskForModifyingDnsHost(request)
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

// SaveSingleTaskForModifyingDnsHostWithCallback invokes the domain_intl.SaveSingleTaskForModifyingDnsHost API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/savesingletaskformodifyingdnshost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveSingleTaskForModifyingDnsHostWithCallback(request *SaveSingleTaskForModifyingDnsHostRequest, callback func(response *SaveSingleTaskForModifyingDnsHostResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveSingleTaskForModifyingDnsHostResponse
		var err error
		defer close(result)
		response, err = client.SaveSingleTaskForModifyingDnsHost(request)
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

// SaveSingleTaskForModifyingDnsHostRequest is the request struct for api SaveSingleTaskForModifyingDnsHost
type SaveSingleTaskForModifyingDnsHostRequest struct {
	*requests.RpcRequest
	InstanceId string    `position:"Query" name:"InstanceId"`
	Lang       string    `position:"Query" name:"Lang"`
	DnsName    string    `position:"Query" name:"DnsName"`
	Ip         *[]string `position:"Query" name:"Ip"  type:"Repeated"`
}

// SaveSingleTaskForModifyingDnsHostResponse is the response struct for api SaveSingleTaskForModifyingDnsHost
type SaveSingleTaskForModifyingDnsHostResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	TaskNo    string `json:"TaskNo" xml:"TaskNo"`
}

// CreateSaveSingleTaskForModifyingDnsHostRequest creates a request to invoke SaveSingleTaskForModifyingDnsHost API
func CreateSaveSingleTaskForModifyingDnsHostRequest() (request *SaveSingleTaskForModifyingDnsHostRequest) {
	request = &SaveSingleTaskForModifyingDnsHostRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain-intl", "2017-12-18", "SaveSingleTaskForModifyingDnsHost", "", "")
	return
}

// CreateSaveSingleTaskForModifyingDnsHostResponse creates a response to parse from SaveSingleTaskForModifyingDnsHost response
func CreateSaveSingleTaskForModifyingDnsHostResponse() (response *SaveSingleTaskForModifyingDnsHostResponse) {
	response = &SaveSingleTaskForModifyingDnsHostResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
