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

// QueryRegistrantProfileRealNameVerificationInfo invokes the domain.QueryRegistrantProfileRealNameVerificationInfo API synchronously
// api document: https://help.aliyun.com/api/domain/queryregistrantprofilerealnameverificationinfo.html
func (client *Client) QueryRegistrantProfileRealNameVerificationInfo(request *QueryRegistrantProfileRealNameVerificationInfoRequest) (response *QueryRegistrantProfileRealNameVerificationInfoResponse, err error) {
	response = CreateQueryRegistrantProfileRealNameVerificationInfoResponse()
	err = client.DoAction(request, response)
	return
}

// QueryRegistrantProfileRealNameVerificationInfoWithChan invokes the domain.QueryRegistrantProfileRealNameVerificationInfo API asynchronously
// api document: https://help.aliyun.com/api/domain/queryregistrantprofilerealnameverificationinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryRegistrantProfileRealNameVerificationInfoWithChan(request *QueryRegistrantProfileRealNameVerificationInfoRequest) (<-chan *QueryRegistrantProfileRealNameVerificationInfoResponse, <-chan error) {
	responseChan := make(chan *QueryRegistrantProfileRealNameVerificationInfoResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryRegistrantProfileRealNameVerificationInfo(request)
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

// QueryRegistrantProfileRealNameVerificationInfoWithCallback invokes the domain.QueryRegistrantProfileRealNameVerificationInfo API asynchronously
// api document: https://help.aliyun.com/api/domain/queryregistrantprofilerealnameverificationinfo.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryRegistrantProfileRealNameVerificationInfoWithCallback(request *QueryRegistrantProfileRealNameVerificationInfoRequest, callback func(response *QueryRegistrantProfileRealNameVerificationInfoResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryRegistrantProfileRealNameVerificationInfoResponse
		var err error
		defer close(result)
		response, err = client.QueryRegistrantProfileRealNameVerificationInfo(request)
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

// QueryRegistrantProfileRealNameVerificationInfoRequest is the request struct for api QueryRegistrantProfileRealNameVerificationInfo
type QueryRegistrantProfileRealNameVerificationInfoRequest struct {
	*requests.RpcRequest
	UserClientIp        string           `position:"Query" name:"UserClientIp"`
	Lang                string           `position:"Query" name:"Lang"`
	RegistrantProfileId requests.Integer `position:"Query" name:"RegistrantProfileId"`
	FetchImage          requests.Boolean `position:"Query" name:"FetchImage"`
}

// QueryRegistrantProfileRealNameVerificationInfoResponse is the response struct for api QueryRegistrantProfileRealNameVerificationInfo
type QueryRegistrantProfileRealNameVerificationInfoResponse struct {
	*responses.BaseResponse
	RequestId              string `json:"RequestId" xml:"RequestId"`
	SubmissionDate         string `json:"SubmissionDate" xml:"SubmissionDate"`
	ModificationDate       string `json:"ModificationDate" xml:"ModificationDate"`
	IdentityCredential     string `json:"IdentityCredential" xml:"IdentityCredential"`
	RegistrantProfileId    int    `json:"RegistrantProfileId" xml:"RegistrantProfileId"`
	IdentityCredentialNo   string `json:"IdentityCredentialNo" xml:"IdentityCredentialNo"`
	IdentityCredentialType string `json:"IdentityCredentialType" xml:"IdentityCredentialType"`
}

// CreateQueryRegistrantProfileRealNameVerificationInfoRequest creates a request to invoke QueryRegistrantProfileRealNameVerificationInfo API
func CreateQueryRegistrantProfileRealNameVerificationInfoRequest() (request *QueryRegistrantProfileRealNameVerificationInfoRequest) {
	request = &QueryRegistrantProfileRealNameVerificationInfoRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "QueryRegistrantProfileRealNameVerificationInfo", "", "")
	return
}

// CreateQueryRegistrantProfileRealNameVerificationInfoResponse creates a response to parse from QueryRegistrantProfileRealNameVerificationInfo response
func CreateQueryRegistrantProfileRealNameVerificationInfoResponse() (response *QueryRegistrantProfileRealNameVerificationInfoResponse) {
	response = &QueryRegistrantProfileRealNameVerificationInfoResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
