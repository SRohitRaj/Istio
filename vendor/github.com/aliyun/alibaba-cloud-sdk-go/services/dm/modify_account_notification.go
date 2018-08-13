package dm

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

// ModifyAccountNotification invokes the dm.ModifyAccountNotification API synchronously
// api document: https://help.aliyun.com/api/dm/modifyaccountnotification.html
func (client *Client) ModifyAccountNotification(request *ModifyAccountNotificationRequest) (response *ModifyAccountNotificationResponse, err error) {
	response = CreateModifyAccountNotificationResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyAccountNotificationWithChan invokes the dm.ModifyAccountNotification API asynchronously
// api document: https://help.aliyun.com/api/dm/modifyaccountnotification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyAccountNotificationWithChan(request *ModifyAccountNotificationRequest) (<-chan *ModifyAccountNotificationResponse, <-chan error) {
	responseChan := make(chan *ModifyAccountNotificationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyAccountNotification(request)
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

// ModifyAccountNotificationWithCallback invokes the dm.ModifyAccountNotification API asynchronously
// api document: https://help.aliyun.com/api/dm/modifyaccountnotification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyAccountNotificationWithCallback(request *ModifyAccountNotificationRequest, callback func(response *ModifyAccountNotificationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyAccountNotificationResponse
		var err error
		defer close(result)
		response, err = client.ModifyAccountNotification(request)
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

// ModifyAccountNotificationRequest is the request struct for api ModifyAccountNotification
type ModifyAccountNotificationRequest struct {
	*requests.RpcRequest
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	Region               string           `position:"Query" name:"Region"`
	Status               string           `position:"Query" name:"Status"`
}

// ModifyAccountNotificationResponse is the response struct for api ModifyAccountNotification
type ModifyAccountNotificationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyAccountNotificationRequest creates a request to invoke ModifyAccountNotification API
func CreateModifyAccountNotificationRequest() (request *ModifyAccountNotificationRequest) {
	request = &ModifyAccountNotificationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Dm", "2015-11-23", "ModifyAccountNotification", "", "")
	return
}

// CreateModifyAccountNotificationResponse creates a response to parse from ModifyAccountNotification response
func CreateModifyAccountNotificationResponse() (response *ModifyAccountNotificationResponse) {
	response = &ModifyAccountNotificationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
