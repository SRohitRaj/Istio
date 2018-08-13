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

// QueryContactInfo invokes the domain_intl.QueryContactInfo API synchronously
// api document: https://help.aliyun.com/api/domain-intl/querycontactinfo.html
func (client *Client) QueryContactInfo(request *QueryContactInfoRequest) (response *QueryContactInfoResponse, err error) {
	response = CreateQueryContactInfoResponse()
	err = client.DoAction(request, response)
	return
}

// QueryContactInfoWithChan invokes the domain_intl.QueryContactInfo API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/querycontactinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryContactInfoWithChan(request *QueryContactInfoRequest) (<-chan *QueryContactInfoResponse, <-chan error) {
	responseChan := make(chan *QueryContactInfoResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryContactInfo(request)
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

// QueryContactInfoWithCallback invokes the domain_intl.QueryContactInfo API asynchronously
// api document: https://help.aliyun.com/api/domain-intl/querycontactinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryContactInfoWithCallback(request *QueryContactInfoRequest, callback func(response *QueryContactInfoResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryContactInfoResponse
		var err error
		defer close(result)
		response, err = client.QueryContactInfo(request)
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

// QueryContactInfoRequest is the request struct for api QueryContactInfo
type QueryContactInfoRequest struct {
	*requests.RpcRequest
	UserClientIp string `position:"Query" name:"UserClientIp"`
	Lang         string `position:"Query" name:"Lang"`
	DomainName   string `position:"Query" name:"DomainName"`
	ContactType  string `position:"Query" name:"ContactType"`
}

// QueryContactInfoResponse is the response struct for api QueryContactInfo
type QueryContactInfoResponse struct {
	*responses.BaseResponse
	RequestId              string `json:"RequestId" xml:"RequestId"`
	CreateDate             string `json:"CreateDate" xml:"CreateDate"`
	RegistrantName         string `json:"RegistrantName" xml:"RegistrantName"`
	RegistrantOrganization string `json:"RegistrantOrganization" xml:"RegistrantOrganization"`
	Country                string `json:"Country" xml:"Country"`
	Province               string `json:"Province" xml:"Province"`
	City                   string `json:"City" xml:"City"`
	Address                string `json:"Address" xml:"Address"`
	Email                  string `json:"Email" xml:"Email"`
	PostalCode             string `json:"PostalCode" xml:"PostalCode"`
	TelArea                string `json:"TelArea" xml:"TelArea"`
	Telephone              string `json:"Telephone" xml:"Telephone"`
	TelExt                 string `json:"TelExt" xml:"TelExt"`
}

// CreateQueryContactInfoRequest creates a request to invoke QueryContactInfo API
func CreateQueryContactInfoRequest() (request *QueryContactInfoRequest) {
	request = &QueryContactInfoRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain-intl", "2017-12-18", "QueryContactInfo", "", "")
	return
}

// CreateQueryContactInfoResponse creates a response to parse from QueryContactInfo response
func CreateQueryContactInfoResponse() (response *QueryContactInfoResponse) {
	response = &QueryContactInfoResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
