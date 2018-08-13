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

// ImportDataFromDatabase invokes the rds.ImportDataFromDatabase API synchronously
// api document: https://help.aliyun.com/api/rds/importdatafromdatabase.html
func (client *Client) ImportDataFromDatabase(request *ImportDataFromDatabaseRequest) (response *ImportDataFromDatabaseResponse, err error) {
	response = CreateImportDataFromDatabaseResponse()
	err = client.DoAction(request, response)
	return
}

// ImportDataFromDatabaseWithChan invokes the rds.ImportDataFromDatabase API asynchronously
// api document: https://help.aliyun.com/api/rds/importdatafromdatabase.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ImportDataFromDatabaseWithChan(request *ImportDataFromDatabaseRequest) (<-chan *ImportDataFromDatabaseResponse, <-chan error) {
	responseChan := make(chan *ImportDataFromDatabaseResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ImportDataFromDatabase(request)
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

// ImportDataFromDatabaseWithCallback invokes the rds.ImportDataFromDatabase API asynchronously
// api document: https://help.aliyun.com/api/rds/importdatafromdatabase.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ImportDataFromDatabaseWithCallback(request *ImportDataFromDatabaseRequest, callback func(response *ImportDataFromDatabaseResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ImportDataFromDatabaseResponse
		var err error
		defer close(result)
		response, err = client.ImportDataFromDatabase(request)
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

// ImportDataFromDatabaseRequest is the request struct for api ImportDataFromDatabase
type ImportDataFromDatabaseRequest struct {
	*requests.RpcRequest
	OwnerId                requests.Integer `position:"Query" name:"OwnerId"`
	ResourceOwnerAccount   string           `position:"Query" name:"ResourceOwnerAccount"`
	ResourceOwnerId        requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DBInstanceId           string           `position:"Query" name:"DBInstanceId"`
	SourceDatabaseIp       string           `position:"Query" name:"SourceDatabaseIp"`
	SourceDatabasePort     string           `position:"Query" name:"SourceDatabasePort"`
	SourceDatabaseUserName string           `position:"Query" name:"SourceDatabaseUserName"`
	SourceDatabasePassword string           `position:"Query" name:"SourceDatabasePassword"`
	ImportDataType         string           `position:"Query" name:"ImportDataType"`
	IsLockTable            string           `position:"Query" name:"IsLockTable"`
	SourceDatabaseDBNames  string           `position:"Query" name:"SourceDatabaseDBNames"`
	OwnerAccount           string           `position:"Query" name:"OwnerAccount"`
}

// ImportDataFromDatabaseResponse is the response struct for api ImportDataFromDatabase
type ImportDataFromDatabaseResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	ImportId  int    `json:"ImportId" xml:"ImportId"`
}

// CreateImportDataFromDatabaseRequest creates a request to invoke ImportDataFromDatabase API
func CreateImportDataFromDatabaseRequest() (request *ImportDataFromDatabaseRequest) {
	request = &ImportDataFromDatabaseRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "ImportDataFromDatabase", "rds", "openAPI")
	return
}

// CreateImportDataFromDatabaseResponse creates a response to parse from ImportDataFromDatabase response
func CreateImportDataFromDatabaseResponse() (response *ImportDataFromDatabaseResponse) {
	response = &ImportDataFromDatabaseResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
