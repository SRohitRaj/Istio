package cdn

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

// DeleteLiveStreamTranscode invokes the cdn.DeleteLiveStreamTranscode API synchronously
// api document: https://help.aliyun.com/api/cdn/deletelivestreamtranscode.html
func (client *Client) DeleteLiveStreamTranscode(request *DeleteLiveStreamTranscodeRequest) (response *DeleteLiveStreamTranscodeResponse, err error) {
	response = CreateDeleteLiveStreamTranscodeResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteLiveStreamTranscodeWithChan invokes the cdn.DeleteLiveStreamTranscode API asynchronously
// api document: https://help.aliyun.com/api/cdn/deletelivestreamtranscode.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteLiveStreamTranscodeWithChan(request *DeleteLiveStreamTranscodeRequest) (<-chan *DeleteLiveStreamTranscodeResponse, <-chan error) {
	responseChan := make(chan *DeleteLiveStreamTranscodeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteLiveStreamTranscode(request)
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

// DeleteLiveStreamTranscodeWithCallback invokes the cdn.DeleteLiveStreamTranscode API asynchronously
// api document: https://help.aliyun.com/api/cdn/deletelivestreamtranscode.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteLiveStreamTranscodeWithCallback(request *DeleteLiveStreamTranscodeRequest, callback func(response *DeleteLiveStreamTranscodeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteLiveStreamTranscodeResponse
		var err error
		defer close(result)
		response, err = client.DeleteLiveStreamTranscode(request)
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

// DeleteLiveStreamTranscodeRequest is the request struct for api DeleteLiveStreamTranscode
type DeleteLiveStreamTranscodeRequest struct {
	*requests.RpcRequest
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
	OwnerAccount  string           `position:"Query" name:"OwnerAccount"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	Domain        string           `position:"Query" name:"Domain"`
	Template      string           `position:"Query" name:"Template"`
	App           string           `position:"Query" name:"App"`
}

// DeleteLiveStreamTranscodeResponse is the response struct for api DeleteLiveStreamTranscode
type DeleteLiveStreamTranscodeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateDeleteLiveStreamTranscodeRequest creates a request to invoke DeleteLiveStreamTranscode API
func CreateDeleteLiveStreamTranscodeRequest() (request *DeleteLiveStreamTranscodeRequest) {
	request = &DeleteLiveStreamTranscodeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "DeleteLiveStreamTranscode", "", "")
	return
}

// CreateDeleteLiveStreamTranscodeResponse creates a response to parse from DeleteLiveStreamTranscode response
func CreateDeleteLiveStreamTranscodeResponse() (response *DeleteLiveStreamTranscodeResponse) {
	response = &DeleteLiveStreamTranscodeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
