// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"encoding/binary"
	"testing"
	"unsafe"
)

const intWidth int = int(unsafe.Sizeof(0))

var byteOrder binary.ByteOrder

// Inspired by etcd cpuutil
func init() {
	i := int(0x1)
	if v := (*[intWidth]byte)(unsafe.Pointer(&i)); v[0] == 0 {
		byteOrder = binary.BigEndian
	} else {
		byteOrder = binary.LittleEndian
	}
}

func TestNtohs(t *testing.T) {
	hostValue := ntohs(0xbeef)
	expectValue := 0xbeef
	if byteOrder == binary.LittleEndian {
		expectValue = 0xefbe
	}
	if hostValue != uint16(expectValue) {
		t.Errorf("Expected evaluating ntohs(%v) is %v, actual %v", 0xbeef, expectValue, hostValue)
	}
}
