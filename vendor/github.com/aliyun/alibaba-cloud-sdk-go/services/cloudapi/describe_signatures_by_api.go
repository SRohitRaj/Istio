package cloudapi

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

// DescribeSignaturesByApi invokes the cloudapi.DescribeSignaturesByApi API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describesignaturesbyapi.html
func (client *Client) DescribeSignaturesByApi(request *DescribeSignaturesByApiRequest) (response *DescribeSignaturesByApiResponse, err error) {
	response = CreateDescribeSignaturesByApiResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeSignaturesByApiWithChan invokes the cloudapi.DescribeSignaturesByApi API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describesignaturesbyapi.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeSignaturesByApiWithChan(request *DescribeSignaturesByApiRequest) (<-chan *DescribeSignaturesByApiResponse, <-chan error) {
	responseChan := make(chan *DescribeSignaturesByApiResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeSignaturesByApi(request)
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

// DescribeSignaturesByApiWithCallback invokes the cloudapi.DescribeSignaturesByApi API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describesignaturesbyapi.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeSignaturesByApiWithCallback(request *DescribeSignaturesByApiRequest, callback func(response *DescribeSignaturesByApiResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeSignaturesByApiResponse
		var err error
		defer close(result)
		response, err = client.DescribeSignaturesByApi(request)
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

// DescribeSignaturesByApiRequest is the request struct for api DescribeSignaturesByApi
type DescribeSignaturesByApiRequest struct {
	*requests.RpcRequest
	GroupId   string `position:"Query" name:"GroupId"`
	ApiId     string `position:"Query" name:"ApiId"`
	StageName string `position:"Query" name:"StageName"`
}

// DescribeSignaturesByApiResponse is the response struct for api DescribeSignaturesByApi
type DescribeSignaturesByApiResponse struct {
	*responses.BaseResponse
	RequestId  string     `json:"RequestId" xml:"RequestId"`
	Signatures Signatures `json:"Signatures" xml:"Signatures"`
}

// CreateDescribeSignaturesByApiRequest creates a request to invoke DescribeSignaturesByApi API
func CreateDescribeSignaturesByApiRequest() (request *DescribeSignaturesByApiRequest) {
	request = &DescribeSignaturesByApiRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeSignaturesByApi", "apigateway", "openAPI")
	return
}

// CreateDescribeSignaturesByApiResponse creates a response to parse from DescribeSignaturesByApi response
func CreateDescribeSignaturesByApiResponse() (response *DescribeSignaturesByApiResponse) {
	response = &DescribeSignaturesByApiResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
