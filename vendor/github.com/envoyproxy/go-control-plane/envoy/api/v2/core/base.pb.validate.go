// Code generated by protoc-gen-validate
// source: envoy/api/v2/core/base.proto
// DO NOT EDIT!!!

package core

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = types.DynamicAny{}
)

// Validate checks the field values on Locality with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Locality) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Region

	// no validation rules for Zone

	// no validation rules for SubZone

	return nil
}

// LocalityValidationError is the validation error returned by
// Locality.Validate if the designated constraints aren't met.
type LocalityValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e LocalityValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLocality.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = LocalityValidationError{}

// Validate checks the field values on Node with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Node) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for Cluster

	if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return NodeValidationError{
				Field:  "Metadata",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetLocality()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return NodeValidationError{
				Field:  "Locality",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	// no validation rules for BuildVersion

	return nil
}

// NodeValidationError is the validation error returned by Node.Validate if the
// designated constraints aren't met.
type NodeValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e NodeValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNode.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = NodeValidationError{}

// Validate checks the field values on Metadata with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Metadata) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for FilterMetadata

	return nil
}

// MetadataValidationError is the validation error returned by
// Metadata.Validate if the designated constraints aren't met.
type MetadataValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e MetadataValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetadata.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = MetadataValidationError{}

// Validate checks the field values on RuntimeUInt32 with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *RuntimeUInt32) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for DefaultValue

	if len(m.GetRuntimeKey()) < 1 {
		return RuntimeUInt32ValidationError{
			Field:  "RuntimeKey",
			Reason: "value length must be at least 1 bytes",
		}
	}

	return nil
}

// RuntimeUInt32ValidationError is the validation error returned by
// RuntimeUInt32.Validate if the designated constraints aren't met.
type RuntimeUInt32ValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e RuntimeUInt32ValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRuntimeUInt32.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = RuntimeUInt32ValidationError{}

// Validate checks the field values on HeaderValue with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *HeaderValue) Validate() error {
	if m == nil {
		return nil
	}

	if l := len(m.GetKey()); l < 1 || l > 16384 {
		return HeaderValueValidationError{
			Field:  "Key",
			Reason: "value length must be between 1 and 16384 bytes, inclusive",
		}
	}

	if len(m.GetValue()) > 16384 {
		return HeaderValueValidationError{
			Field:  "Value",
			Reason: "value length must be at most 16384 bytes",
		}
	}

	return nil
}

// HeaderValueValidationError is the validation error returned by
// HeaderValue.Validate if the designated constraints aren't met.
type HeaderValueValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HeaderValueValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHeaderValue.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HeaderValueValidationError{}

// Validate checks the field values on HeaderValueOption with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *HeaderValueOption) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetHeader() == nil {
		return HeaderValueOptionValidationError{
			Field:  "Header",
			Reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetHeader()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HeaderValueOptionValidationError{
				Field:  "Header",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetAppend()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HeaderValueOptionValidationError{
				Field:  "Append",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// HeaderValueOptionValidationError is the validation error returned by
// HeaderValueOption.Validate if the designated constraints aren't met.
type HeaderValueOptionValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HeaderValueOptionValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHeaderValueOption.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HeaderValueOptionValidationError{}

// Validate checks the field values on HeaderMap with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *HeaderMap) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetHeaders() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HeaderMapValidationError{
					Field:  fmt.Sprintf("Headers[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// HeaderMapValidationError is the validation error returned by
// HeaderMap.Validate if the designated constraints aren't met.
type HeaderMapValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HeaderMapValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHeaderMap.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HeaderMapValidationError{}

// Validate checks the field values on DataSource with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *DataSource) Validate() error {
	if m == nil {
		return nil
	}

	switch m.Specifier.(type) {

	case *DataSource_Filename:

		if len(m.GetFilename()) < 1 {
			return DataSourceValidationError{
				Field:  "Filename",
				Reason: "value length must be at least 1 bytes",
			}
		}

	case *DataSource_InlineBytes:

		if len(m.GetInlineBytes()) < 1 {
			return DataSourceValidationError{
				Field:  "InlineBytes",
				Reason: "value length must be at least 1 bytes",
			}
		}

	case *DataSource_InlineString:

		if len(m.GetInlineString()) < 1 {
			return DataSourceValidationError{
				Field:  "InlineString",
				Reason: "value length must be at least 1 bytes",
			}
		}

	default:
		return DataSourceValidationError{
			Field:  "Specifier",
			Reason: "value is required",
		}

	}

	return nil
}

// DataSourceValidationError is the validation error returned by
// DataSource.Validate if the designated constraints aren't met.
type DataSourceValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e DataSourceValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDataSource.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = DataSourceValidationError{}

// Validate checks the field values on TransportSocket with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *TransportSocket) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetName()) < 1 {
		return TransportSocketValidationError{
			Field:  "Name",
			Reason: "value length must be at least 1 bytes",
		}
	}

	switch m.ConfigType.(type) {

	case *TransportSocket_Config:

		if v, ok := interface{}(m.GetConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TransportSocketValidationError{
					Field:  "Config",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *TransportSocket_TypedConfig:

		if v, ok := interface{}(m.GetTypedConfig()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return TransportSocketValidationError{
					Field:  "TypedConfig",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// TransportSocketValidationError is the validation error returned by
// TransportSocket.Validate if the designated constraints aren't met.
type TransportSocketValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e TransportSocketValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTransportSocket.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = TransportSocketValidationError{}

// Validate checks the field values on SocketOption with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *SocketOption) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Description

	// no validation rules for Level

	// no validation rules for Name

	if _, ok := SocketOption_SocketState_name[int32(m.GetState())]; !ok {
		return SocketOptionValidationError{
			Field:  "State",
			Reason: "value must be one of the defined enum values",
		}
	}

	switch m.Value.(type) {

	case *SocketOption_IntValue:
		// no validation rules for IntValue

	case *SocketOption_BufValue:
		// no validation rules for BufValue

	default:
		return SocketOptionValidationError{
			Field:  "Value",
			Reason: "value is required",
		}

	}

	return nil
}

// SocketOptionValidationError is the validation error returned by
// SocketOption.Validate if the designated constraints aren't met.
type SocketOptionValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e SocketOptionValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSocketOption.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = SocketOptionValidationError{}

// Validate checks the field values on RuntimeFractionalPercent with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RuntimeFractionalPercent) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetDefaultValue() == nil {
		return RuntimeFractionalPercentValidationError{
			Field:  "DefaultValue",
			Reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetDefaultValue()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RuntimeFractionalPercentValidationError{
				Field:  "DefaultValue",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	// no validation rules for RuntimeKey

	return nil
}

// RuntimeFractionalPercentValidationError is the validation error returned by
// RuntimeFractionalPercent.Validate if the designated constraints aren't met.
type RuntimeFractionalPercentValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e RuntimeFractionalPercentValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRuntimeFractionalPercent.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = RuntimeFractionalPercentValidationError{}

// Validate checks the field values on ControlPlane with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ControlPlane) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Identifier

	return nil
}

// ControlPlaneValidationError is the validation error returned by
// ControlPlane.Validate if the designated constraints aren't met.
type ControlPlaneValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e ControlPlaneValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sControlPlane.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = ControlPlaneValidationError{}
