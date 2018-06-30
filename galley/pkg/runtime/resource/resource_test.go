//  Copyright 2018 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package resource

import (
	"reflect"
	"testing"

	"github.com/gogo/protobuf/types"
)

func TestKind_Equality_True(t *testing.T) {
	k1 := MessageName{"a"}
	k2 := MessageName{"a"}

	if k1 != k2 {
		t.Fatalf("Expected to be equal: %v == %v", k1, k2)
	}
}

func TestKind_Equality_False(t *testing.T) {
	k1 := MessageName{"a"}
	k2 := MessageName{"v"}

	if k1 == k2 {
		t.Fatalf("Expected to be not equal: %v == %v", k1, k2)
	}
}

func TestVersion_Equality_True(t *testing.T) {
	v1 := Version("a")
	v2 := Version("a")

	if v1 != v2 {
		t.Fatalf("Expected to be equal: %v == %v", v1, v2)
	}
}

func TestVersion_Equality_False(t *testing.T) {
	v1 := Version("a")
	v2 := Version("v")

	if v1 == v2 {
		t.Fatalf("Expected to be not equal: %v == %v", v1, v2)
	}
}
func TestKey_Equality_True(t *testing.T) {
	k1 := Key{MessageName: MessageName{"a"}, FullName: "ks"}
	k2 := Key{MessageName: MessageName{"a"}, FullName: "ks"}

	if k1 != k2 {
		t.Fatalf("Expected to be equal: %v == %v", k1, k2)
	}
}

func TestKey_Equality_False_DifferentKind(t *testing.T) {
	k1 := Key{MessageName: MessageName{"a"}, FullName: "ks"}
	k2 := Key{MessageName: MessageName{"b"}, FullName: "ks"}

	if k1 == k2 {
		t.Fatalf("Expected to be not equal: %v == %v", k1, k2)
	}
}

func TestKey_Equality_False_DifferentName(t *testing.T) {
	k1 := Key{MessageName: MessageName{"a"}, FullName: "ks"}
	k2 := Key{MessageName: MessageName{"a"}, FullName: "otherks"}

	if k1 == k2 {
		t.Fatalf("Expected to be not equal: %v == %v", k1, k2)
	}
}

func TestKey_String(t *testing.T) {
	k1 := Key{MessageName: MessageName{"a"}, FullName: "ks"}
	// Ensure that it doesn't crash
	_ = k1.String()
}

func TestVersionedKey_Equality_True(t *testing.T) {
	k1 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "ks"}, Version: Version("v1")}
	k2 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "ks"}, Version: Version("v1")}

	if k1 != k2 {
		t.Fatalf("Expected to be equal: %v == %v", k1, k2)
	}
}

func TestVersionedKey_Equality_False_DifferentKind(t *testing.T) {
	k1 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "ks"}, Version: Version("v1")}
	k2 := VersionedKey{
		Key: Key{MessageName: MessageName{"b"}, FullName: "ks"}, Version: Version("v1")}

	if k1 == k2 {
		t.Fatalf("Expected to be not equal: %v == %v", k1, k2)
	}
}

func TestVersionedKey_Equality_False_DifferentName(t *testing.T) {
	k1 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "ks"}, Version: Version("v1")}
	k2 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "otherks"}, Version: Version("v1")}

	if k1 == k2 {
		t.Fatalf("Expected to be not equal: %v == %v", k1, k2)
	}
}

func TestVersionedKey_Equality_False_DifferentVersion(t *testing.T) {
	k1 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "ks"}, Version: Version("v1")}
	k2 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "ks"}, Version: Version("v2")}

	if k1 == k2 {
		t.Fatalf("Expected to be not equal: %v == %v", k1, k2)
	}
}

func TestVersionedKey_String(t *testing.T) {
	k1 := VersionedKey{
		Key: Key{MessageName: MessageName{"a"}, FullName: "ks"}, Version: Version("v1")}
	// Ensure that it doesn't crash
	_ = k1.String()
}

func TestResource_IsEmpty(t *testing.T) {
	r := Entry{}
	if !r.IsEmpty() {
		t.Fatal("should have been empty")
	}

	r.Item = &types.Empty{}
	if r.IsEmpty() {
		t.Fatal("should have not been empty")
	}
}

func TestInfo_newProtoInstance_Success(t *testing.T) {
	i := Info{
		goType: reflect.TypeOf(types.Empty{}),
	}
	p := i.NewProtoInstance()

	if p == nil || reflect.TypeOf(p) != reflect.PtrTo(reflect.TypeOf(types.Empty{})) {
		t.Fatalf("Unexpected proto type returned: %v", p)
	}
}

func TestInfo_newProtoInstance_PanicAtNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Expected panic not found")
		}
	}()

	i := Info{
		goType: nil,
	}
	_ = i.NewProtoInstance()
}

func TestInfo_newProtoInstance_PanicAtNonProto(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Expected panic not found")
		}
	}()

	i := Info{
		goType: reflect.TypeOf(""),
	}
	_ = i.NewProtoInstance()
}

func TestInfo_String(t *testing.T) {
	i := Info{
		TypeURL:     "http://foo.bar.com",
		MessageName: MessageName{"foo"},
	}
	// Ensure that it doesn't crash
	_ = i.String()
}
