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

// QueryDomainList invokes the domain.QueryDomainList API synchronously
// api document: https://help.aliyun.com/api/domain/querydomainlist.html
func (client *Client) QueryDomainList(request *QueryDomainListRequest) (response *QueryDomainListResponse, err error) {
	response = CreateQueryDomainListResponse()
	err = client.DoAction(request, response)
	return
}

// QueryDomainListWithChan invokes the domain.QueryDomainList API asynchronously
// api document: https://help.aliyun.com/api/domain/querydomainlist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryDomainListWithChan(request *QueryDomainListRequest) (<-chan *QueryDomainListResponse, <-chan error) {
	responseChan := make(chan *QueryDomainListResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryDomainList(request)
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

// QueryDomainListWithCallback invokes the domain.QueryDomainList API asynchronously
// api document: https://help.aliyun.com/api/domain/querydomainlist.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryDomainListWithCallback(request *QueryDomainListRequest, callback func(response *QueryDomainListResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryDomainListResponse
		var err error
		defer close(result)
		response, err = client.QueryDomainList(request)
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

// QueryDomainListRequest is the request struct for api QueryDomainList
type QueryDomainListRequest struct {
	*requests.RpcRequest
	StartExpirationDate   requests.Integer `position:"Query" name:"StartExpirationDate"`
	UserClientIp          string           `position:"Query" name:"UserClientIp"`
	Lang                  string           `position:"Query" name:"Lang"`
	QueryType             string           `position:"Query" name:"QueryType"`
	EndExpirationDate     requests.Integer `position:"Query" name:"EndExpirationDate"`
	StartRegistrationDate requests.Integer `position:"Query" name:"StartRegistrationDate"`
	EndRegistrationDate   requests.Integer `position:"Query" name:"EndRegistrationDate"`
	DomainName            string           `position:"Query" name:"DomainName"`
	OrderByType           string           `position:"Query" name:"OrderByType"`
	OrderKeyType          string           `position:"Query" name:"OrderKeyType"`
	ProductDomainType     string           `position:"Query" name:"ProductDomainType"`
	PageNum               requests.Integer `position:"Query" name:"PageNum"`
	PageSize              requests.Integer `position:"Query" name:"PageSize"`
}

// QueryDomainListResponse is the response struct for api QueryDomainList
type QueryDomainListResponse struct {
	*responses.BaseResponse
	RequestId      string                `json:"RequestId" xml:"RequestId"`
	TotalItemNum   int                   `json:"TotalItemNum" xml:"TotalItemNum"`
	CurrentPageNum int                   `json:"CurrentPageNum" xml:"CurrentPageNum"`
	TotalPageNum   int                   `json:"TotalPageNum" xml:"TotalPageNum"`
	PageSize       int                   `json:"PageSize" xml:"PageSize"`
	PrePage        bool                  `json:"PrePage" xml:"PrePage"`
	NextPage       bool                  `json:"NextPage" xml:"NextPage"`
	Data           DataInQueryDomainList `json:"Data" xml:"Data"`
}

// CreateQueryDomainListRequest creates a request to invoke QueryDomainList API
func CreateQueryDomainListRequest() (request *QueryDomainListRequest) {
	request = &QueryDomainListRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "QueryDomainList", "", "")
	return
}

// CreateQueryDomainListResponse creates a response to parse from QueryDomainList response
func CreateQueryDomainListResponse() (response *QueryDomainListResponse) {
	response = &QueryDomainListResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
