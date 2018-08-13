package rds

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

// SwitchDBInstanceHA invokes the rds.SwitchDBInstanceHA API synchronously
// api document: https://help.aliyun.com/api/rds/switchdbinstanceha.html
func (client *Client) SwitchDBInstanceHA(request *SwitchDBInstanceHARequest) (response *SwitchDBInstanceHAResponse, err error) {
	response = CreateSwitchDBInstanceHAResponse()
	err = client.DoAction(request, response)
	return
}

// SwitchDBInstanceHAWithChan invokes the rds.SwitchDBInstanceHA API asynchronously
// api document: https://help.aliyun.com/api/rds/switchdbinstanceha.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SwitchDBInstanceHAWithChan(request *SwitchDBInstanceHARequest) (<-chan *SwitchDBInstanceHAResponse, <-chan error) {
	responseChan := make(chan *SwitchDBInstanceHAResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SwitchDBInstanceHA(request)
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

// SwitchDBInstanceHAWithCallback invokes the rds.SwitchDBInstanceHA API asynchronously
// api document: https://help.aliyun.com/api/rds/switchdbinstanceha.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SwitchDBInstanceHAWithCallback(request *SwitchDBInstanceHARequest, callback func(response *SwitchDBInstanceHAResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SwitchDBInstanceHAResponse
		var err error
		defer close(result)
		response, err = client.SwitchDBInstanceHA(request)
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

// SwitchDBInstanceHARequest is the request struct for api SwitchDBInstanceHA
type SwitchDBInstanceHARequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	NodeId               string           `position:"Query" name:"NodeId"`
	Operation            string           `position:"Query" name:"Operation"`
	Force                string           `position:"Query" name:"Force"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	EffectiveTime        string           `position:"Query" name:"EffectiveTime"`
}

// SwitchDBInstanceHAResponse is the response struct for api SwitchDBInstanceHA
type SwitchDBInstanceHAResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSwitchDBInstanceHARequest creates a request to invoke SwitchDBInstanceHA API
func CreateSwitchDBInstanceHARequest() (request *SwitchDBInstanceHARequest) {
	request = &SwitchDBInstanceHARequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "SwitchDBInstanceHA", "rds", "openAPI")
	return
}

// CreateSwitchDBInstanceHAResponse creates a response to parse from SwitchDBInstanceHA response
func CreateSwitchDBInstanceHAResponse() (response *SwitchDBInstanceHAResponse) {
	response = &SwitchDBInstanceHAResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
