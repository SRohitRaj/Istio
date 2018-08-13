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

// DescribeDeployedApi invokes the cloudapi.DescribeDeployedApi API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describedeployedapi.html
func (client *Client) DescribeDeployedApi(request *DescribeDeployedApiRequest) (response *DescribeDeployedApiResponse, err error) {
	response = CreateDescribeDeployedApiResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeDeployedApiWithChan invokes the cloudapi.DescribeDeployedApi API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describedeployedapi.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDeployedApiWithChan(request *DescribeDeployedApiRequest) (<-chan *DescribeDeployedApiResponse, <-chan error) {
	responseChan := make(chan *DescribeDeployedApiResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeDeployedApi(request)
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

// DescribeDeployedApiWithCallback invokes the cloudapi.DescribeDeployedApi API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describedeployedapi.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeDeployedApiWithCallback(request *DescribeDeployedApiRequest, callback func(response *DescribeDeployedApiResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeDeployedApiResponse
		var err error
		defer close(result)
		response, err = client.DescribeDeployedApi(request)
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

// DescribeDeployedApiRequest is the request struct for api DescribeDeployedApi
type DescribeDeployedApiRequest struct {
	*requests.RpcRequest
	GroupId   string `position:"Query" name:"GroupId"`
	ApiId     string `position:"Query" name:"ApiId"`
	StageName string `position:"Query" name:"StageName"`
}

// DescribeDeployedApiResponse is the response struct for api DescribeDeployedApi
type DescribeDeployedApiResponse struct {
	*responses.BaseResponse
	RequestId              string                                      `json:"RequestId" xml:"RequestId"`
	RegionId               string                                      `json:"RegionId" xml:"RegionId"`
	GroupId                string                                      `json:"GroupId" xml:"GroupId"`
	GroupName              string                                      `json:"GroupName" xml:"GroupName"`
	StageName              string                                      `json:"StageName" xml:"StageName"`
	ApiId                  string                                      `json:"ApiId" xml:"ApiId"`
	ApiName                string                                      `json:"ApiName" xml:"ApiName"`
	Description            string                                      `json:"Description" xml:"Description"`
	Visibility             string                                      `json:"Visibility" xml:"Visibility"`
	AuthType               string                                      `json:"AuthType" xml:"AuthType"`
	ResultType             string                                      `json:"ResultType" xml:"ResultType"`
	ResultSample           string                                      `json:"ResultSample" xml:"ResultSample"`
	FailResultSample       string                                      `json:"FailResultSample" xml:"FailResultSample"`
	DeployedTime           string                                      `json:"DeployedTime" xml:"DeployedTime"`
	AllowSignatureMethod   string                                      `json:"AllowSignatureMethod" xml:"AllowSignatureMethod"`
	RequestConfig          RequestConfig                               `json:"RequestConfig" xml:"RequestConfig"`
	ServiceConfig          ServiceConfig                               `json:"ServiceConfig" xml:"ServiceConfig"`
	OpenIdConnectConfig    OpenIdConnectConfig                         `json:"OpenIdConnectConfig" xml:"OpenIdConnectConfig"`
	ErrorCodeSamples       ErrorCodeSamplesInDescribeDeployedApi       `json:"ErrorCodeSamples" xml:"ErrorCodeSamples"`
	ResultDescriptions     ResultDescriptionsInDescribeDeployedApi     `json:"ResultDescriptions" xml:"ResultDescriptions"`
	SystemParameters       SystemParametersInDescribeDeployedApi       `json:"SystemParameters" xml:"SystemParameters"`
	CustomSystemParameters CustomSystemParametersInDescribeDeployedApi `json:"CustomSystemParameters" xml:"CustomSystemParameters"`
	ConstantParameters     ConstantParametersInDescribeDeployedApi     `json:"ConstantParameters" xml:"ConstantParameters"`
	RequestParameters      RequestParametersInDescribeDeployedApi      `json:"RequestParameters" xml:"RequestParameters"`
	ServiceParameters      ServiceParametersInDescribeDeployedApi      `json:"ServiceParameters" xml:"ServiceParameters"`
	ServiceParametersMap   ServiceParametersMapInDescribeDeployedApi   `json:"ServiceParametersMap" xml:"ServiceParametersMap"`
}

// CreateDescribeDeployedApiRequest creates a request to invoke DescribeDeployedApi API
func CreateDescribeDeployedApiRequest() (request *DescribeDeployedApiRequest) {
	request = &DescribeDeployedApiRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeDeployedApi", "apigateway", "openAPI")
	return
}

// CreateDescribeDeployedApiResponse creates a response to parse from DescribeDeployedApi response
func CreateDescribeDeployedApiResponse() (response *DescribeDeployedApiResponse) {
	response = &DescribeDeployedApiResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
