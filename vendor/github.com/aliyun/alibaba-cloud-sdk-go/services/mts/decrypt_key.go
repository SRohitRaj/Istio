package mts

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

// DecryptKey invokes the mts.DecryptKey API synchronously
// api document: https://help.aliyun.com/api/mts/decryptkey.html
func (client *Client) DecryptKey(request *DecryptKeyRequest) (response *DecryptKeyResponse, err error) {
	response = CreateDecryptKeyResponse()
	err = client.DoAction(request, response)
	return
}

// DecryptKeyWithChan invokes the mts.DecryptKey API asynchronously
// api document: https://help.aliyun.com/api/mts/decryptkey.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DecryptKeyWithChan(request *DecryptKeyRequest) (<-chan *DecryptKeyResponse, <-chan error) {
	responseChan := make(chan *DecryptKeyResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DecryptKey(request)
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

// DecryptKeyWithCallback invokes the mts.DecryptKey API asynchronously
// api document: https://help.aliyun.com/api/mts/decryptkey.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DecryptKeyWithCallback(request *DecryptKeyRequest, callback func(response *DecryptKeyResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DecryptKeyResponse
		var err error
		defer close(result)
		response, err = client.DecryptKey(request)
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

// DecryptKeyRequest is the request struct for api DecryptKey
type DecryptKeyRequest struct {
	*requests.RpcRequest
	OwnerId              string `position:"Query" name:"OwnerId"`
	ResourceOwnerId      string `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string `position:"Query" name:"OwnerAccount"`
	CiphertextBlob       string `position:"Query" name:"CiphertextBlob"`
	Rand                 string `position:"Query" name:"Rand"`
}

// DecryptKeyResponse is the response struct for api DecryptKey
type DecryptKeyResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Plaintext string `json:"Plaintext" xml:"Plaintext"`
	Rand      string `json:"Rand" xml:"Rand"`
}

// CreateDecryptKeyRequest creates a request to invoke DecryptKey API
func CreateDecryptKeyRequest() (request *DecryptKeyRequest) {
	request = &DecryptKeyRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Mts", "2014-06-18", "DecryptKey", "mts", "openAPI")
	return
}

// CreateDecryptKeyResponse creates a response to parse from DecryptKey response
func CreateDecryptKeyResponse() (response *DecryptKeyResponse) {
	response = &DecryptKeyResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
