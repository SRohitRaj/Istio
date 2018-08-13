package domain

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

// QueryBidRecords invokes the domain.QueryBidRecords API synchronously
// api document: https://help.aliyun.com/api/domain/querybidrecords.html
func (client *Client) QueryBidRecords(request *QueryBidRecordsRequest) (response *QueryBidRecordsResponse, err error) {
	response = CreateQueryBidRecordsResponse()
	err = client.DoAction(request, response)
	return
}

// QueryBidRecordsWithChan invokes the domain.QueryBidRecords API asynchronously
// api document: https://help.aliyun.com/api/domain/querybidrecords.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryBidRecordsWithChan(request *QueryBidRecordsRequest) (<-chan *QueryBidRecordsResponse, <-chan error) {
	responseChan := make(chan *QueryBidRecordsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryBidRecords(request)
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

// QueryBidRecordsWithCallback invokes the domain.QueryBidRecords API asynchronously
// api document: https://help.aliyun.com/api/domain/querybidrecords.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryBidRecordsWithCallback(request *QueryBidRecordsRequest, callback func(response *QueryBidRecordsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryBidRecordsResponse
		var err error
		defer close(result)
		response, err = client.QueryBidRecords(request)
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

// QueryBidRecordsRequest is the request struct for api QueryBidRecords
type QueryBidRecordsRequest struct {
	*requests.RpcRequest
	AuctionId   string           `position:"Body" name:"AuctionId"`
	CurrentPage requests.Integer `position:"Body" name:"CurrentPage"`
	PageSize    requests.Integer `position:"Body" name:"PageSize"`
}

// QueryBidRecordsResponse is the response struct for api QueryBidRecords
type QueryBidRecordsResponse struct {
	*responses.BaseResponse
	RequestId      string      `json:"RequestId" xml:"RequestId"`
	TotalItemNum   int         `json:"TotalItemNum" xml:"TotalItemNum"`
	CurrentPageNum int         `json:"CurrentPageNum" xml:"CurrentPageNum"`
	PageSize       int         `json:"PageSize" xml:"PageSize"`
	TotalPageNum   int         `json:"TotalPageNum" xml:"TotalPageNum"`
	Data           []BidRecord `json:"Data" xml:"Data"`
}

// CreateQueryBidRecordsRequest creates a request to invoke QueryBidRecords API
func CreateQueryBidRecordsRequest() (request *QueryBidRecordsRequest) {
	request = &QueryBidRecordsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-02-08", "QueryBidRecords", "", "")
	return
}

// CreateQueryBidRecordsResponse creates a response to parse from QueryBidRecords response
func CreateQueryBidRecordsResponse() (response *QueryBidRecordsResponse) {
	response = &QueryBidRecordsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
