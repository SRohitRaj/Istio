// Copyright 2018 Istio Authors
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

package source

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gogo/protobuf/types"
	kube_meta "istio.io/istio/galley/pkg/metadata/kube"
	"istio.io/istio/pkg/features/galley"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic/fake"
	dtesting "k8s.io/client-go/testing"

	"istio.io/istio/galley/pkg/kube"
	"istio.io/istio/galley/pkg/kube/converter"
	"istio.io/istio/galley/pkg/meshconfig"
	"istio.io/istio/galley/pkg/runtime/resource"
	"istio.io/istio/galley/pkg/testing/mock"
)

var emptyInfo resource.Info

func init() {
	b := resource.NewSchemaBuilder()
	b.Register("type.googleapis.com/google.protobuf.Empty")
	schema := b.Build()
	emptyInfo, _ = schema.Lookup("type.googleapis.com/google.protobuf.Empty")
}

func mockKubeMeta() *mock.Kube {
	var k mock.Kube
	//len(kube_meta.Types.All())
	for i := 0; i < 100; i++ {
		cl := fake.NewSimpleDynamicClient(runtime.NewScheme())
		k.AddResponse(cl, nil)
	}
	return &k
}

func TestNewSource(t *testing.T) {
	k := mockKubeMeta()
	cfg := converter.Config{
		Mesh: meshconfig.NewInMemory(),
	}
	p, err := New(k, 0, &cfg)
	if err != nil {
		t.Fatalf("Unexpected error found: %v", err)
	}

	if p == nil {
		t.Fatal("expected non-nil source")
	}
}

func TestNewSource_Error(t *testing.T) {
	k := &mock.Kube{}
	k.AddResponse(nil, errors.New("newDynamicClient error"))

	cfg := converter.Config{
		Mesh: meshconfig.NewInMemory(),
	}
	_, err := New(k, 0, &cfg)
	if err == nil || err.Error() != "newDynamicClient error" {
		t.Fatalf("Expected error not found: %v", err)
	}
}

func TestNewSource_ServiceEntry(t *testing.T) {
	typeCount := len(kube_meta.Types.All())
	tests := []struct {
		name    string
		convert bool
		wantLen int
	}{
		{
			name:    "Unset",
			wantLen: typeCount - 1,
		},
		{
			name:    "Disabled",
			convert: false,
			wantLen: typeCount - 1,
		},
		{
			name:    "Enabled",
			convert: true,
			wantLen: typeCount,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			x := galley.ConvertK8SService
			defer func() {
				galley.ConvertK8SService = x
			}()
			galley.ConvertK8SService = test.convert

			s, err := New(mockKubeMeta(), 0, &converter.Config{
				Mesh: meshconfig.NewInMemory(),
			})
			if err != nil {
				t.Fatalf("New failed: %v", err)
			}
			if got := len(s.(*sourceImpl).listeners); got != test.wantLen {
				t.Errorf("Constructed with %d listeners, want %d", got, test.wantLen)
			}
		})
	}
}

func TestSource_BasicEvents(t *testing.T) {
	k := &mock.Kube{}
	cl := fake.NewSimpleDynamicClient(runtime.NewScheme())

	k.AddResponse(cl, nil)

	i1 := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "List",
			"metadata": map[string]interface{}{
				"name":            "f1",
				"namespace":       "ns",
				"resourceVersion": "rv1",
			},
			"spec": map[string]interface{}{},
		},
	}

	cl.PrependReactor("*", "foos", func(action dtesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &unstructured.UnstructuredList{Items: []unstructured.Unstructured{i1}}, nil
	})
	w := mock.NewWatch()
	cl.PrependWatchReactor("foos", func(action dtesting.Action) (handled bool, ret watch.Interface, err error) {
		return true, w, nil
	})

	entries := []kube.ResourceSpec{
		{
			Kind:      "List",
			Singular:  "List",
			Plural:    "foos",
			Target:    emptyInfo,
			Converter: converter.Get("identity"),
		},
	}

	cfg := converter.Config{
		Mesh: meshconfig.NewInMemory(),
	}
	s, err := newSource(k, 0, &cfg, entries)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if s == nil {
		t.Fatal("Expected non nil source")
	}

	ch, err := s.Start()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	log := logChannelOutput(ch, 2)
	expected := strings.TrimSpace(`
[Event](Added: [VKey](type.googleapis.com/google.protobuf.Empty:ns/f1 @rv1))
[Event](FullSync)
`)
	if log != expected {
		t.Fatalf("Event mismatch:\nActual:\n%s\nExpected:\n%s\n", log, expected)
	}

	w.Send(watch.Event{Type: watch.Deleted, Object: &i1})
	log = logChannelOutput(ch, 1)
	expected = strings.TrimSpace(`
[Event](Deleted: [VKey](type.googleapis.com/google.protobuf.Empty:ns/f1 @rv1))
		`)
	if log != expected {
		t.Fatalf("Event mismatch:\nActual:\n%s\nExpected:\n%s\n", log, expected)
	}

	s.Stop()
}

func TestSource_BasicEvents_NoConversion(t *testing.T) {

	k := &mock.Kube{}
	cl := fake.NewSimpleDynamicClient(runtime.NewScheme())

	k.AddResponse(cl, nil)

	i1 := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "List",
			"metadata": map[string]interface{}{
				"name":            "f1",
				"namespace":       "ns",
				"resourceVersion": "rv1",
			},
			"spec": map[string]interface{}{},
		},
	}

	cl.PrependReactor("*", "foos", func(action dtesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &unstructured.UnstructuredList{Items: []unstructured.Unstructured{i1}}, nil
	})
	cl.PrependWatchReactor("foos", func(action dtesting.Action) (handled bool, ret watch.Interface, err error) {
		return true, mock.NewWatch(), nil
	})

	entries := []kube.ResourceSpec{
		{
			Kind:      "List",
			Singular:  "List",
			Plural:    "foos",
			Target:    emptyInfo,
			Converter: converter.Get("nil"),
		},
	}

	cfg := converter.Config{
		Mesh: meshconfig.NewInMemory(),
	}
	s, err := newSource(k, 0, &cfg, entries)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if s == nil {
		t.Fatal("Expected non nil source")
	}

	ch, err := s.Start()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	log := logChannelOutput(ch, 1)
	expected := strings.TrimSpace(`
[Event](FullSync)
`)
	if log != expected {
		t.Fatalf("Event mismatch:\nActual:\n%s\nExpected:\n%s\n", log, expected)
	}

	s.Stop()
}

func TestSource_ProtoConversionError(t *testing.T) {
	k := &mock.Kube{}
	cl := fake.NewSimpleDynamicClient(runtime.NewScheme())
	k.AddResponse(cl, nil)

	i1 := unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":            "f1",
				"namespace":       "ns",
				"resourceVersion": "rv1",
			},
			"spec": map[string]interface{}{
				"zoo": "zar",
			},
		},
	}

	cl.PrependReactor("*", "foos", func(action dtesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &unstructured.UnstructuredList{Items: []unstructured.Unstructured{i1}}, nil
	})
	cl.PrependWatchReactor("foos", func(action dtesting.Action) (handled bool, ret watch.Interface, err error) {
		return true, mock.NewWatch(), nil
	})

	entries := []kube.ResourceSpec{
		{
			Kind:     "foo",
			Singular: "foo",
			Plural:   "foos",
			Target:   emptyInfo,
			Converter: func(_ *converter.Config, info resource.Info, name resource.FullName, kind string, u *unstructured.Unstructured) ([]converter.Entry, error) {
				return nil, fmt.Errorf("cant convert")
			},
		},
	}

	cfg := converter.Config{
		Mesh: meshconfig.NewInMemory(),
	}
	s, err := newSource(k, 0, &cfg, entries)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if s == nil {
		t.Fatal("Expected non nil source")
	}

	ch, err := s.Start()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// The add event should not appear.
	log := logChannelOutput(ch, 1)
	expected := strings.TrimSpace(`
[Event](FullSync)
`)
	if log != expected {
		t.Fatalf("Event mismatch:\nActual:\n%s\nExpected:\n%s\n", log, expected)
	}

	s.Stop()
}

func TestSource_MangledNames(t *testing.T) {
	k := &mock.Kube{}
	cl := fake.NewSimpleDynamicClient(runtime.NewScheme())
	k.AddResponse(cl, nil)

	i1 := unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name":            "f1",
				"namespace":       "ns",
				"resourceVersion": "rv1",
			},
			"spec": map[string]interface{}{
				"zoo": "zar",
			},
		},
	}

	cl.PrependReactor("*", "foos", func(action dtesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &unstructured.UnstructuredList{Items: []unstructured.Unstructured{i1}}, nil
	})
	cl.PrependWatchReactor("foos", func(action dtesting.Action) (handled bool, ret watch.Interface, err error) {
		return true, mock.NewWatch(), nil
	})

	entries := []kube.ResourceSpec{
		{
			Kind:     "foo",
			Singular: "foo",
			Plural:   "foos",
			Target:   emptyInfo,
			Converter: func(_ *converter.Config, info resource.Info, name resource.FullName, kind string, u *unstructured.Unstructured) ([]converter.Entry, error) {
				e := converter.Entry{
					Key:      resource.FullNameFromNamespaceAndName("foo", name.String()),
					Resource: &types.Struct{},
				}

				return []converter.Entry{e}, nil
			},
		},
	}

	cfg := converter.Config{
		Mesh: meshconfig.NewInMemory(),
	}
	s, err := newSource(k, 0, &cfg, entries)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if s == nil {
		t.Fatal("Expected non nil source")
	}

	ch, err := s.Start()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// The mangled name foo/ns/f1 should appear.
	log := logChannelOutput(ch, 1)
	expected := strings.TrimSpace(`
[Event](Added: [VKey](type.googleapis.com/google.protobuf.Empty:foo/ns/f1 @rv1))
`)
	if log != expected {
		t.Fatalf("Event mismatch:\nActual:\n%s\nExpected:\n%s\n", log, expected)
	}

	s.Stop()
}

func logChannelOutput(ch chan resource.Event, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		item := <-ch
		result += fmt.Sprintf("%v\n", item)
	}

	result = strings.TrimSpace(result)

	return result
}
