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

package diag

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"istio.io/istio/pkg/config/resource"
)

// DocPrefix is the root URL for validation message docs
const DocPrefix = "https://istio.io/docs/reference/config/analysis"

// MessageType is a type of diagnostic message
type MessageType struct {
	// The level of the message.
	level Level

	// The error code of the message
	code string

	// TODO: Make this localizable
	template string
}

// Level returns the level of the MessageType
func (m *MessageType) Level() Level { return m.level }

// Code returns the code of the MessageType
func (m *MessageType) Code() string { return m.code }

// Template returns the message template used by the MessageType
func (m *MessageType) Template() string { return m.template }

// Message is a specific diagnostic message
// TODO: Implement using Analysis message API
type Message struct {
	Type *MessageType

	// The Parameters to the message
	Parameters []interface{}

	// Resource is the underlying resource instance associated with the
	// message, or nil if no resource is associated with it.
	Resource *resource.Instance

	// DocRef is an optional reference tracker for the documentation URL
	DocRef string

	line int
}

// Unstructured returns this message as a JSON-style unstructured map
func (m *Message) Unstructured(includeOrigin bool) map[string]interface{} {
	result := make(map[string]interface{})

	result["code"] = m.Type.Code()
	result["level"] = m.Type.Level().String()
	if includeOrigin && m.Resource != nil {
		result["origin"] = m.Resource.Origin.FriendlyName()
		if m.Resource.Origin.Reference() != nil {

			loc := m.Resource.Origin.Reference().String()
			if m.line != 0 {
				loc = m.ReplaceLine(loc)
			}
			result["reference"] = loc
		}
	}
	result["message"] = fmt.Sprintf(m.Type.Template(), m.Parameters...)

	docQueryString := ""
	if m.DocRef != "" {
		docQueryString = fmt.Sprintf("?ref=%s", m.DocRef)
	}
	result["documentation_url"] = fmt.Sprintf("%s/%s/%s", DocPrefix, strings.ToLower(m.Type.Code()), docQueryString)

	return result
}

// Origin returns the origin of the message
func (m *Message) Origin() string {
	origin := ""
	if m.Resource != nil {
		loc := ""
		if m.Resource.Origin.Reference() != nil {
			loc = " " + m.Resource.Origin.Reference().String()
			if m.GetLine() != 0 {
				loc = m.ReplaceLine(loc)
			}
		}
		origin = " (" + m.Resource.Origin.FriendlyName() + loc + ")"
	}
	return origin
}

// String implements io.Stringer
func (m *Message) String() string {
	return fmt.Sprintf("%v [%v]%s %s",
		m.Type.Level(), m.Type.Code(), m.Origin(),
		fmt.Sprintf(m.Type.Template(), m.Parameters...))
}

// MarshalJSON satisfies the Marshaler interface
func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Unstructured(true))
}

// NewMessageType returns a new MessageType instance.
func NewMessageType(level Level, code, template string) *MessageType {
	return &MessageType{
		level:    level,
		code:     code,
		template: template,
	}
}

// NewMessage returns a new Message instance from an existing type.
func NewMessage(mt *MessageType, r *resource.Instance, p ...interface{}) Message {
	return Message{
		Type:       mt,
		Resource:   r,
		Parameters: p,
	}
}

// SetLine sets the error line number of the message
func (m *Message) SetLine(l int) {
	m.line = l
}

// GetLine gets the error line number of the message
func (m *Message) GetLine() int {
	return m.line
}

// ReplaceLine replaces the line number from the input String method of the Reference object to the line from the Message object
func (m Message) ReplaceLine(l string) string {
	loc := l
	_, err := strconv.Atoi(string(loc[len(loc)-1]))
	for err == nil {
		loc = loc[:len(loc)-1]
		_, err = strconv.Atoi(string(loc[len(loc)-1]))
	}
	loc += fmt.Sprintf("%d", m.GetLine())
	return loc
}
