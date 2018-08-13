package ddospro

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

// Data is a nested struct in ddospro response
type Data struct {
	Cycle             int          `json:"Cycle" xml:"Cycle"`
	CycleResidue      int          `json:"CycleResidue" xml:"CycleResidue"`
	LastBlockTime     int          `json:"LastBlockTime" xml:"LastBlockTime"`
	HoleThresholdMbps int          `json:"HoleThresholdMbps" xml:"HoleThresholdMbps"`
	CcMode            string       `json:"CcMode" xml:"CcMode"`
	WhiteList         string       `json:"WhiteList" xml:"WhiteList"`
	ExpectionOpenTime int          `json:"ExpectionOpenTime" xml:"ExpectionOpenTime"`
	Status            bool         `json:"Status" xml:"Status"`
	ResourceId        string       `json:"ResourceId" xml:"ResourceId"`
	OpDesc            string       `json:"OpDesc" xml:"OpDesc"`
	OpAction          int          `json:"OpAction" xml:"OpAction"`
	LastCloseTime     int          `json:"LastCloseTime" xml:"LastCloseTime"`
	RemainingTime     int          `json:"RemainingTime" xml:"RemainingTime"`
	BlockZone         string       `json:"BlockZone" xml:"BlockZone"`
	BlockTime         int          `json:"BlockTime" xml:"BlockTime"`
	TotalClose        int          `json:"TotalClose" xml:"TotalClose"`
	Region            string       `json:"Region" xml:"Region"`
	TotalTime         int          `json:"TotalTime" xml:"TotalTime"`
	CycleTime         int          `json:"CycleTime" xml:"CycleTime"`
	DdosStatus        int          `json:"DdosStatus" xml:"DdosStatus"`
	BlackList         string       `json:"BlackList" xml:"BlackList"`
	Result            int          `json:"Result" xml:"Result"`
	Vip               string       `json:"Vip" xml:"Vip"`
	OpDate            int          `json:"OpDate" xml:"OpDate"`
	UnblockTime       int          `json:"UnblockTime" xml:"UnblockTime"`
	BpsDrop           []string     `json:"BpsDrop" xml:"BpsDrop"`
	Attacks           []string     `json:"Attacks" xml:"Attacks"`
	PpsDrop           []string     `json:"PpsDrop" xml:"PpsDrop"`
	PpsTotal          []string     `json:"PpsTotal" xml:"PpsTotal"`
	BpsTotal          []string     `json:"BpsTotal" xml:"BpsTotal"`
	Total             []string     `json:"Total" xml:"Total"`
	TimeScope         TimeScope    `json:"TimeScope" xml:"TimeScope"`
	PageInfo          PageInfo     `json:"PageInfo" xml:"PageInfo"`
	List              []AttackInfo `json:"List" xml:"List"`
}
