package cdn

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

// DescribeDomainRealTimeData invokes the cdn.DescribeDomainRealTimeData API synchronously
// api document: https://help.aliyun.com/api/cdn/describedomainrealtimedata.html
func (client *Client) DescribeDomainRealTimeData(request *DescribeDomainRealTimeDataRequest) (response *DescribeDomainRealTimeDataResponse, err error) {
	response = CreateDescribeDomainRealTimeDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDomainRealTimeDataWithChan invokes the cdn.DescribeDomainRealTimeData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainrealtimedata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainRealTimeDataWithChan(request *DescribeDomainRealTimeDataRequest) (<-chan *DescribeDomainRealTimeDataResponse, <-chan error) {
	responseChan := make(chan *DescribeDomainRealTimeDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDomainRealTimeData(request)
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

// DescribeDomainRealTimeDataWithCallback invokes the cdn.DescribeDomainRealTimeData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describedomainrealtimedata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDomainRealTimeDataWithCallback(request *DescribeDomainRealTimeDataRequest, callback func(response *DescribeDomainRealTimeDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDomainRealTimeDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeDomainRealTimeData(request)
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

// DescribeDomainRealTimeDataRequest is the request struct for api DescribeDomainRealTimeData
type DescribeDomainRealTimeDataRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	StartTime     string           `position:"Query" name:"StartTime"`
	EndTime       string           `position:"Query" name:"EndTime"`
	Field         string           `position:"Query" name:"Field"`
}

// DescribeDomainRealTimeDataResponse is the response struct for api DescribeDomainRealTimeData
type DescribeDomainRealTimeDataResponse struct {
	*responses.BaseResponse
	RequestId       string                                      `json:"RequestId" xml:"RequestId"`
	DomainName      string                                      `json:"DomainName" xml:"DomainName"`
	Field           string                                      `json:"Field" xml:"Field"`
	StartTime       string                                      `json:"StartTime" xml:"StartTime"`
	EndTime         string                                      `json:"EndTime" xml:"EndTime"`
	DataPerInterval DataPerIntervalInDescribeDomainRealTimeData `json:"DataPerInterval" xml:"DataPerInterval"`
}

// CreateDescribeDomainRealTimeDataRequest creates a request to invoke DescribeDomainRealTimeData API
func CreateDescribeDomainRealTimeDataRequest() (request *DescribeDomainRealTimeDataRequest) {
	request = &DescribeDomainRealTimeDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeDomainRealTimeData", "", "")
	return
}

// CreateDescribeDomainRealTimeDataResponse creates a response to parse from DescribeDomainRealTimeData response
func CreateDescribeDomainRealTimeDataResponse() (response *DescribeDomainRealTimeDataResponse) {
	response = &DescribeDomainRealTimeDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
