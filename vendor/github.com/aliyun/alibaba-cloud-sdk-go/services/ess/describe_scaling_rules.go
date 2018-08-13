package ess

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

// DescribeScalingRules invokes the ess.DescribeScalingRules API synchronously
// api document: https://help.aliyun.com/api/ess/describescalingrules.html
func (client *Client) DescribeScalingRules(request *DescribeScalingRulesRequest) (response *DescribeScalingRulesResponse, err error) {
	response = CreateDescribeScalingRulesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeScalingRulesWithChan invokes the ess.DescribeScalingRules API asynchronously
// api document: https://help.aliyun.com/api/ess/describescalingrules.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeScalingRulesWithChan(request *DescribeScalingRulesRequest) (<-chan *DescribeScalingRulesResponse, <-chan error) {
	responseChan := make(chan *DescribeScalingRulesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeScalingRules(request)
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

// DescribeScalingRulesWithCallback invokes the ess.DescribeScalingRules API asynchronously
// api document: https://help.aliyun.com/api/ess/describescalingrules.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeScalingRulesWithCallback(request *DescribeScalingRulesRequest, callback func(response *DescribeScalingRulesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeScalingRulesResponse
		var err error
		defer close(result)
		response, err = client.DescribeScalingRules(request)
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

// DescribeScalingRulesRequest is the request struct for api DescribeScalingRules
type DescribeScalingRulesRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	ScalingGroupId       string           `position:"Query" name:"ScalingGroupId"`
	ScalingRuleId1       string           `position:"Query" name:"ScalingRuleId.1"`
	ScalingRuleId2       string           `position:"Query" name:"ScalingRuleId.2"`
	ScalingRuleId3       string           `position:"Query" name:"ScalingRuleId.3"`
	ScalingRuleId4       string           `position:"Query" name:"ScalingRuleId.4"`
	ScalingRuleId5       string           `position:"Query" name:"ScalingRuleId.5"`
	ScalingRuleId6       string           `position:"Query" name:"ScalingRuleId.6"`
	ScalingRuleId7       string           `position:"Query" name:"ScalingRuleId.7"`
	ScalingRuleId8       string           `position:"Query" name:"ScalingRuleId.8"`
	ScalingRuleId9       string           `position:"Query" name:"ScalingRuleId.9"`
	ScalingRuleId10      string           `position:"Query" name:"ScalingRuleId.10"`
	ScalingRuleName1     string           `position:"Query" name:"ScalingRuleName.1"`
	ScalingRuleName2     string           `position:"Query" name:"ScalingRuleName.2"`
	ScalingRuleName3     string           `position:"Query" name:"ScalingRuleName.3"`
	ScalingRuleName4     string           `position:"Query" name:"ScalingRuleName.4"`
	ScalingRuleName5     string           `position:"Query" name:"ScalingRuleName.5"`
	ScalingRuleName6     string           `position:"Query" name:"ScalingRuleName.6"`
	ScalingRuleName7     string           `position:"Query" name:"ScalingRuleName.7"`
	ScalingRuleName8     string           `position:"Query" name:"ScalingRuleName.8"`
	ScalingRuleName9     string           `position:"Query" name:"ScalingRuleName.9"`
	ScalingRuleName10    string           `position:"Query" name:"ScalingRuleName.10"`
	ScalingRuleAri1      string           `position:"Query" name:"ScalingRuleAri.1"`
	ScalingRuleAri2      string           `position:"Query" name:"ScalingRuleAri.2"`
	ScalingRuleAri3      string           `position:"Query" name:"ScalingRuleAri.3"`
	ScalingRuleAri4      string           `position:"Query" name:"ScalingRuleAri.4"`
	ScalingRuleAri5      string           `position:"Query" name:"ScalingRuleAri.5"`
	ScalingRuleAri6      string           `position:"Query" name:"ScalingRuleAri.6"`
	ScalingRuleAri7      string           `position:"Query" name:"ScalingRuleAri.7"`
	ScalingRuleAri8      string           `position:"Query" name:"ScalingRuleAri.8"`
	ScalingRuleAri9      string           `position:"Query" name:"ScalingRuleAri.9"`
	ScalingRuleAri10     string           `position:"Query" name:"ScalingRuleAri.10"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
}

// DescribeScalingRulesResponse is the response struct for api DescribeScalingRules
type DescribeScalingRulesResponse struct {
	*responses.BaseResponse
	TotalCount   int          `json:"TotalCount" xml:"TotalCount"`
	PageNumber   int          `json:"PageNumber" xml:"PageNumber"`
	PageSize     int          `json:"PageSize" xml:"PageSize"`
	RequestId    string       `json:"RequestId" xml:"RequestId"`
	ScalingRules ScalingRules `json:"ScalingRules" xml:"ScalingRules"`
}

// CreateDescribeScalingRulesRequest creates a request to invoke DescribeScalingRules API
func CreateDescribeScalingRulesRequest() (request *DescribeScalingRulesRequest) {
	request = &DescribeScalingRulesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Ess", "2014-08-28", "DescribeScalingRules", "ess", "openAPI")
	return
}

// CreateDescribeScalingRulesResponse creates a response to parse from DescribeScalingRules response
func CreateDescribeScalingRulesResponse() (response *DescribeScalingRulesResponse) {
	response = &DescribeScalingRulesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
