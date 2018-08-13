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

// QueryChangeLogList invokes the domain_intl.QueryChangeLogList API synchronously
// api document: https://help.aliyun.com/api/domain-intl/querychangeloglist.html
func (client *Client) QueryChangeLogList(request *QueryChangeLogListRequest) (response *QueryChangeLogListResponse, err error) {
	response = CreateQueryChangeLogListResponse()
	err = client.DoAction(request, response)
	return
}

// QueryChangeLogListWithChan invokes the domain_intl.QueryChangeLogList API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/querychangeloglist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryChangeLogListWithChan(request *QueryChangeLogListRequest) (<-chan *QueryChangeLogListResponse, <-chan error) {
	responseChan := make(chan *QueryChangeLogListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryChangeLogList(request)
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

// QueryChangeLogListWithCallback invokes the domain_intl.QueryChangeLogList API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/querychangeloglist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryChangeLogListWithCallback(request *QueryChangeLogListRequest, callback func(response *QueryChangeLogListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryChangeLogListResponse
		var err error
		defer close(result)
		response, err = client.QueryChangeLogList(request)
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

// QueryChangeLogListRequest is the request struct for api QueryChangeLogList
type QueryChangeLogListRequest struct {
	*requests.RpcRequest
	UserClientIp string           `position:"Query" name:"UserClientIp"`
	Lang         string           `position:"Query" name:"Lang"`
	DomainName   string           `position:"Query" name:"DomainName"`
	PageNum      requests.Integer `position:"Query" name:"PageNum"`
	PageSize     requests.Integer `position:"Query" name:"PageSize"`
	StartDate    requests.Integer `position:"Query" name:"StartDate"`
	EndDate      requests.Integer `position:"Query" name:"EndDate"`
}

// QueryChangeLogListResponse is the response struct for api QueryChangeLogList
type QueryChangeLogListResponse struct {
	*responses.BaseResponse
	RequestId      string                   `json:"RequestId" xml:"RequestId"`
	TotalItemNum   int                      `json:"TotalItemNum" xml:"TotalItemNum"`
	CurrentPageNum int                      `json:"CurrentPageNum" xml:"CurrentPageNum"`
	TotalPageNum   int                      `json:"TotalPageNum" xml:"TotalPageNum"`
	PageSize       int                      `json:"PageSize" xml:"PageSize"`
	PrePage        bool                     `json:"PrePage" xml:"PrePage"`
	NextPage       bool                     `json:"NextPage" xml:"NextPage"`
	ResultLimit    bool                     `json:"ResultLimit" xml:"ResultLimit"`
	Data           DataInQueryChangeLogList `json:"Data" xml:"Data"`
}

// CreateQueryChangeLogListRequest creates a request to invoke QueryChangeLogList API
func CreateQueryChangeLogListRequest() (request *QueryChangeLogListRequest) {
	request = &QueryChangeLogListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain-intl", "2017-12-18", "QueryChangeLogList", "domain", "openAPI")
	return
}

// CreateQueryChangeLogListResponse creates a response to parse from QueryChangeLogList response
func CreateQueryChangeLogListResponse() (response *QueryChangeLogListResponse) {
	response = &QueryChangeLogListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
