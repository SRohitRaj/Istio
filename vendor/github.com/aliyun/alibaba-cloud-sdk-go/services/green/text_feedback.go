package green

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

// TextFeedback invokes the green.TextFeedback API synchronously
// api document: https://help.aliyun.com/api/green/textfeedback.html
func (client *Client) TextFeedback(request *TextFeedbackRequest) (response *TextFeedbackResponse, err error) {
	response = CreateTextFeedbackResponse()
	err = client.DoAction(request, response)
	return
}

// TextFeedbackWithChan invokes the green.TextFeedback API asynchronously
// api document: https://help.aliyun.com/api/green/textfeedback.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TextFeedbackWithChan(request *TextFeedbackRequest) (<-chan *TextFeedbackResponse, <-chan error) {
	responseChan := make(chan *TextFeedbackResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.TextFeedback(request)
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

// TextFeedbackWithCallback invokes the green.TextFeedback API asynchronously
// api document: https://help.aliyun.com/api/green/textfeedback.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) TextFeedbackWithCallback(request *TextFeedbackRequest, callback func(response *TextFeedbackResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *TextFeedbackResponse
		var err error
		defer close(result)
		response, err = client.TextFeedback(request)
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

// TextFeedbackRequest is the request struct for api TextFeedback
type TextFeedbackRequest struct {
	*requests.RoaRequest
	ClientInfo string `position:"Query" name:"ClientInfo"`
}

// TextFeedbackResponse is the response struct for api TextFeedback
type TextFeedbackResponse struct {
	*responses.BaseResponse
}

// CreateTextFeedbackRequest creates a request to invoke TextFeedback API
func CreateTextFeedbackRequest() (request *TextFeedbackRequest) {
	request = &TextFeedbackRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("Green", "2017-08-25", "TextFeedback", "/green/text/feedback", "green", "openAPI")
	request.Method = requests.POST
	return
}

// CreateTextFeedbackResponse creates a response to parse from TextFeedback response
func CreateTextFeedbackResponse() (response *TextFeedbackResponse) {
	response = &TextFeedbackResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
