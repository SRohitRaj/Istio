package vpc

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

// CreatePhysicalConnectionNew invokes the vpc.CreatePhysicalConnectionNew API synchronously
// api document: https://help.aliyun.com/api/vpc/createphysicalconnectionnew.html
func (client *Client) CreatePhysicalConnectionNew(request *CreatePhysicalConnectionNewRequest) (response *CreatePhysicalConnectionNewResponse, err error) {
	response = CreateCreatePhysicalConnectionNewResponse()
	err = client.DoAction(request, response)
	return
}

// CreatePhysicalConnectionNewWithChan invokes the vpc.CreatePhysicalConnectionNew API asynchronously
// api document: https://help.aliyun.com/api/vpc/createphysicalconnectionnew.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreatePhysicalConnectionNewWithChan(request *CreatePhysicalConnectionNewRequest) (<-chan *CreatePhysicalConnectionNewResponse, <-chan error) {
	responseChan := make(chan *CreatePhysicalConnectionNewResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreatePhysicalConnectionNew(request)
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

// CreatePhysicalConnectionNewWithCallback invokes the vpc.CreatePhysicalConnectionNew API asynchronously
// api document: https://help.aliyun.com/api/vpc/createphysicalconnectionnew.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreatePhysicalConnectionNewWithCallback(request *CreatePhysicalConnectionNewRequest, callback func(response *CreatePhysicalConnectionNewResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreatePhysicalConnectionNewResponse
		var err error
		defer close(result)
		response, err = client.CreatePhysicalConnectionNew(request)
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

// CreatePhysicalConnectionNewRequest is the request struct for api CreatePhysicalConnectionNew
type CreatePhysicalConnectionNewRequest struct {
	*requests.RpcRequest
	AccessPointId                 string           `position:"Query" name:"AccessPointId"`
	Type                          string           `position:"Query" name:"Type"`
	LineOperator                  string           `position:"Query" name:"LineOperator"`
	Bandwidth                     requests.Integer `position:"Query" name:"bandwidth"`
	PeerLocation                  string           `position:"Query" name:"PeerLocation"`
	PortType                      string           `position:"Query" name:"PortType"`
	RedundantPhysicalConnectionId string           `position:"Query" name:"RedundantPhysicalConnectionId"`
	Description                   string           `position:"Query" name:"Description"`
	Name                          string           `position:"Query" name:"Name"`
	CircuitCode                   string           `position:"Query" name:"CircuitCode"`
	ClientToken                   string           `position:"Query" name:"ClientToken"`
	DeviceName                    string           `position:"Query" name:"DeviceName"`
	InterfaceName                 string           `position:"Query" name:"InterfaceName"`
	OwnerId                       requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount          string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId               requests.Integer `position:"Query" name:"ResourceOwnerId"`
	OwnerAccount                  string           `position:"Query" name:"OwnerAccount"`
}

// CreatePhysicalConnectionNewResponse is the response struct for api CreatePhysicalConnectionNew
type CreatePhysicalConnectionNewResponse struct {
	*responses.BaseResponse
	RequestId            string `json:"RequestId" xml:"RequestId"`
	PhysicalConnectionId string `json:"PhysicalConnectionId" xml:"PhysicalConnectionId"`
}

// CreateCreatePhysicalConnectionNewRequest creates a request to invoke CreatePhysicalConnectionNew API
func CreateCreatePhysicalConnectionNewRequest() (request *CreatePhysicalConnectionNewRequest) {
	request = &CreatePhysicalConnectionNewRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Vpc", "2016-04-28", "CreatePhysicalConnectionNew", "vpc", "openAPI")
	return
}

// CreateCreatePhysicalConnectionNewResponse creates a response to parse from CreatePhysicalConnectionNew response
func CreateCreatePhysicalConnectionNewResponse() (response *CreatePhysicalConnectionNewResponse) {
	response = &CreatePhysicalConnectionNewResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
