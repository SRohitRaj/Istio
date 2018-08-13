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

// DescribeLiveStreamBitRateData invokes the cdn.DescribeLiveStreamBitRateData API synchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreambitratedata.html
func (client *Client) DescribeLiveStreamBitRateData(request *DescribeLiveStreamBitRateDataRequest) (response *DescribeLiveStreamBitRateDataResponse, err error) {
	response = CreateDescribeLiveStreamBitRateDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLiveStreamBitRateDataWithChan invokes the cdn.DescribeLiveStreamBitRateData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreambitratedata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamBitRateDataWithChan(request *DescribeLiveStreamBitRateDataRequest) (<-chan *DescribeLiveStreamBitRateDataResponse, <-chan error) {
	responseChan := make(chan *DescribeLiveStreamBitRateDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLiveStreamBitRateData(request)
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

// DescribeLiveStreamBitRateDataWithCallback invokes the cdn.DescribeLiveStreamBitRateData API asynchronously
// api document: https://help.aliyun.com/api/cdn/describelivestreambitratedata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamBitRateDataWithCallback(request *DescribeLiveStreamBitRateDataRequest, callback func(response *DescribeLiveStreamBitRateDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLiveStreamBitRateDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeLiveStreamBitRateData(request)
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

// DescribeLiveStreamBitRateDataRequest is the request struct for api DescribeLiveStreamBitRateData
type DescribeLiveStreamBitRateDataRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	AppName       string           `position:"Query" name:"AppName"`
	StreamName    string           `position:"Query" name:"StreamName"`
	StartTime     string           `position:"Query" name:"StartTime"`
	EndTime       string           `position:"Query" name:"EndTime"`
}

// DescribeLiveStreamBitRateDataResponse is the response struct for api DescribeLiveStreamBitRateData
type DescribeLiveStreamBitRateDataResponse struct {
	*responses.BaseResponse
	RequestId                string                                                  `json:"RequestId" xml:"RequestId"`
	FrameRateAndBitRateInfos FrameRateAndBitRateInfosInDescribeLiveStreamBitRateData `json:"FrameRateAndBitRateInfos" xml:"FrameRateAndBitRateInfos"`
}

// CreateDescribeLiveStreamBitRateDataRequest creates a request to invoke DescribeLiveStreamBitRateData API
func CreateDescribeLiveStreamBitRateDataRequest() (request *DescribeLiveStreamBitRateDataRequest) {
	request = &DescribeLiveStreamBitRateDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DescribeLiveStreamBitRateData", "", "")
	return
}

// CreateDescribeLiveStreamBitRateDataResponse creates a response to parse from DescribeLiveStreamBitRateData response
func CreateDescribeLiveStreamBitRateDataResponse() (response *DescribeLiveStreamBitRateDataResponse) {
	response = &DescribeLiveStreamBitRateDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
