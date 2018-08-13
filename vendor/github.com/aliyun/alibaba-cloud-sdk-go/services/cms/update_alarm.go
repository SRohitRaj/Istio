package cms

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

// UpdateAlarm invokes the cms.UpdateAlarm API synchronously
// api document: https://help.aliyun.com/api/cms/updatealarm.html
func (client *Client) UpdateAlarm(request *UpdateAlarmRequest) (response *UpdateAlarmResponse, err error) {
	response = CreateUpdateAlarmResponse()
	err = client.DoAction(request, response)
	return
}

// UpdateAlarmWithChan invokes the cms.UpdateAlarm API asynchronously
// api document: https://help.aliyun.com/api/cms/updatealarm.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateAlarmWithChan(request *UpdateAlarmRequest) (<-chan *UpdateAlarmResponse, <-chan error) {
	responseChan := make(chan *UpdateAlarmResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpdateAlarm(request)
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

// UpdateAlarmWithCallback invokes the cms.UpdateAlarm API asynchronously
// api document: https://help.aliyun.com/api/cms/updatealarm.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) UpdateAlarmWithCallback(request *UpdateAlarmRequest, callback func(response *UpdateAlarmResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpdateAlarmResponse
		var err error
		defer close(result)
		response, err = client.UpdateAlarm(request)
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

// UpdateAlarmRequest is the request struct for api UpdateAlarm
type UpdateAlarmRequest struct {
	*requests.RpcRequest
	CallbyCmsOwner     string           `position:"Query" name:"callby_cms_owner"`
	Id                 string           `position:"Query" name:"Id"`
	Name               string           `position:"Query" name:"Name"`
	Period             requests.Integer `position:"Query" name:"Period"`
	Statistics         string           `position:"Query" name:"Statistics"`
	ComparisonOperator string           `position:"Query" name:"ComparisonOperator"`
	Threshold          string           `position:"Query" name:"Threshold"`
	ContactGroups      string           `position:"Query" name:"ContactGroups"`
	StartTime          requests.Integer `position:"Query" name:"StartTime"`
	EndTime            requests.Integer `position:"Query" name:"EndTime"`
	SilenceTime        requests.Integer `position:"Query" name:"SilenceTime"`
	EvaluationCount    requests.Integer `position:"Query" name:"EvaluationCount"`
	NotifyType         requests.Integer `position:"Query" name:"NotifyType"`
	Webhook            string           `position:"Query" name:"Webhook"`
}

// UpdateAlarmResponse is the response struct for api UpdateAlarm
type UpdateAlarmResponse struct {
	*responses.BaseResponse
	Success   bool   `json:"Success" xml:"Success"`
	Code      string `json:"Code" xml:"Code"`
	Message   string `json:"Message" xml:"Message"`
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpdateAlarmRequest creates a request to invoke UpdateAlarm API
func CreateUpdateAlarmRequest() (request *UpdateAlarmRequest) {
	request = &UpdateAlarmRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2018-03-08", "UpdateAlarm", "cms", "openAPI")
	return
}

// CreateUpdateAlarmResponse creates a response to parse from UpdateAlarm response
func CreateUpdateAlarmResponse() (response *UpdateAlarmResponse) {
	response = &UpdateAlarmResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
