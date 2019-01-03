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

package cel

import (
	"errors"
	"reflect"
	"strings"

	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/interpreter"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/pkg/attribute"
)

// Attribute provider resolves typing information by modeling attributes
// as fields in nested proto messages with presence.
//
// For example, an integer attribute "a.b.c" is modeled as the following proto3 definition:
//
// message A {
//   message B {
//     google.protobuf.Int64Value c = 1;
//   }
//   B b = 1;
// }
//
// message Root {
//   A a = 1;
// }
//
type attributeProvider struct {
	root    *node
	typeMap map[string]*node
}

// node corresponds to the message types holding other message types or scalars.
// leaf nodes have field types, inner nodes have children.
type node struct {
	typeName string

	children map[string]*node

	typ       *exprpb.Type
	valueType v1beta1.ValueType
}

func (n *node) HasTrait(trait int) bool {
	// only support getter trait
	return trait == traits.IndexerType
}
func (n *node) TypeName() string {
	return n.typeName
}

func (ap *attributeProvider) newNode(typeName string) *node {
	out := &node{typeName: typeName}
	ap.typeMap[typeName] = out
	return out
}

func (ap *attributeProvider) insert(n *node, words []string, valueType v1beta1.ValueType) {
	if len(words) == 0 {
		n.valueType = valueType
		n.typ = convertType(valueType)
		return
	}

	if n.children == nil {
		n.children = make(map[string]*node)
	}

	child, ok := n.children[words[0]]
	if !ok {
		child = ap.newNode(n.typeName + "." + words[0])
		n.children[words[0]] = child
	}

	ap.insert(child, words[1:], valueType)
}

func newAttributeProvider(attributes map[string]*v1beta1.AttributeManifest_AttributeInfo) *attributeProvider {
	out := &attributeProvider{typeMap: make(map[string]*node)}
	out.root = out.newNode("")
	for name, info := range attributes {
		out.insert(out.root, strings.Split(name, "."), info.ValueType)
	}
	return out
}

func (ap *attributeProvider) newEnvironment() *checker.Env {
	env := checker.NewStandardEnv(packages.DefaultPackage, ap)

	// populate with root-level identifiers
	for name, node := range ap.root.children {
		if node.typ != nil {
			env.Add(decls.NewIdent(name, node.typ, nil))
		} else {
			env.Add(decls.NewIdent(name, decls.NewObjectType(node.typeName), nil))
		}
	}

	// populate with standard overloads
	env.Add(standardFunctions()...)

	return env
}

func (ap *attributeProvider) newInterpreter() interpreter.Interpreter {
	return interpreter.NewStandardInterpreter(packages.DefaultPackage, ap)
}

func (ap *attributeProvider) newActivation(bag attribute.Bag) interpreter.Activation {
	return &attributeActivation{provider: ap, bag: bag}
}

func (ap *attributeProvider) EnumValue(enumName string) ref.Value {
	return types.NewErr("custom enum values not implemented")
}
func (ap *attributeProvider) FindIdent(identName string) (ref.Value, bool) {
	// note: not used with decls in the checker environment
	return nil, false
}
func (ap *attributeProvider) FindType(typeName string) (*exprpb.Type, bool) {
	if _, ok := ap.typeMap[typeName]; ok {
		return decls.NewObjectType(typeName), true
	}
	return nil, false
}
func (ap *attributeProvider) FindFieldType(t *exprpb.Type, fieldName string) (*ref.FieldType, bool) {
	switch v := t.TypeKind.(type) {
	case *exprpb.Type_MessageType:
		node, ok := ap.typeMap[v.MessageType]
		if !ok {
			break
		}

		child, ok := node.children[fieldName]
		if !ok {
			break
		}

		typ := child.typ
		if typ == nil {
			typ = decls.NewObjectType(child.typeName)
		}

		return &ref.FieldType{
				Type:             typ,
				SupportsPresence: true},
			true
	}
	return nil, false
}
func (ap *attributeProvider) NewValue(typeName string, fields map[string]ref.Value) ref.Value {
	return types.NewErr("value construction not implemented")
}
func (ap *attributeProvider) RegisterType(types ...ref.Type) error {
	return nil
}

// Attribute activation binds attribute values to the expression nodes
type attributeActivation struct {
	provider *attributeProvider
	bag      attribute.Bag
}

type value struct {
	node *node
	bag  attribute.Bag
}

func (v value) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	return nil, errors.New("cannot convert attribute message to native types")
}
func (v value) ConvertToType(typeValue ref.Type) ref.Value {
	return types.NewErr("cannot convert attribute message to CEL types")
}
func (v value) Equal(other ref.Value) ref.Value {
	return types.NewErr("attribute message does not support equality")
}
func (v value) Type() ref.Type {
	return v.node
}
func (v value) Value() interface{} {
	return v
}
func (v value) Get(index ref.Value) ref.Value {
	if index.Type() != types.StringType {
		return types.NewErr("select not implemented")
	}

	field := index.Value().(string)
	child, ok := v.node.children[field]
	if !ok {
		return types.NewErr("cannot evaluate select of %q from %s", field, v.node.typeName)
	}

	if child.typ == nil {
		return value{node: child, bag: v.bag}
	}

	value, found := v.bag.Get(child.typeName[1:])
	if found {
		return convertValue(child.valueType, value)
	}
	return defaultValue(child.valueType)
}

func (a *attributeActivation) ResolveName(name string) (ref.Value, bool) {
	if node, ok := a.provider.root.children[name]; ok {
		return value{node: node, bag: a.bag}, true
	}
	return nil, false
}

func (a *attributeActivation) ResolveReference(exprID int64) (ref.Value, bool) {
	return nil, false
}

func (a *attributeActivation) Parent() interpreter.Activation {
	return nil
}
