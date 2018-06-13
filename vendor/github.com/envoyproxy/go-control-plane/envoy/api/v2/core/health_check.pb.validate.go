// Code generated by protoc-gen-validate
// source: envoy/api/v2/core/health_check.proto
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

// Validate checks the field values on HealthCheck with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *HealthCheck) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetTimeout() == nil {
		return HealthCheckValidationError{
			Field:  "Timeout",
			Reason: "value is required",
		}
	}

	if m.GetInterval() == nil {
		return HealthCheckValidationError{
			Field:  "Interval",
			Reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetIntervalJitter()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "IntervalJitter",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetUnhealthyThreshold()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "UnhealthyThreshold",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetHealthyThreshold()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "HealthyThreshold",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetAltPort()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "AltPort",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetReuseConnection()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "ReuseConnection",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetNoTrafficInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "NoTrafficInterval",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetUnhealthyInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "UnhealthyInterval",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetUnhealthyEdgeInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "UnhealthyEdgeInterval",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetHealthyEdgeInterval()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheckValidationError{
				Field:  "HealthyEdgeInterval",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	switch m.HealthChecker.(type) {

	case *HealthCheck_HttpHealthCheck_:

		if v, ok := interface{}(m.GetHttpHealthCheck()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HealthCheckValidationError{
					Field:  "HttpHealthCheck",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *HealthCheck_TcpHealthCheck_:

		if v, ok := interface{}(m.GetTcpHealthCheck()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HealthCheckValidationError{
					Field:  "TcpHealthCheck",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *HealthCheck_RedisHealthCheck_:

		if v, ok := interface{}(m.GetRedisHealthCheck()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HealthCheckValidationError{
					Field:  "RedisHealthCheck",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *HealthCheck_GrpcHealthCheck_:

		if v, ok := interface{}(m.GetGrpcHealthCheck()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HealthCheckValidationError{
					Field:  "GrpcHealthCheck",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	case *HealthCheck_CustomHealthCheck_:

		if v, ok := interface{}(m.GetCustomHealthCheck()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HealthCheckValidationError{
					Field:  "CustomHealthCheck",
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	default:
		return HealthCheckValidationError{
			Field:  "HealthChecker",
			Reason: "value is required",
		}

	}

	return nil
}

// HealthCheckValidationError is the validation error returned by
// HealthCheck.Validate if the designated constraints aren't met.
type HealthCheckValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HealthCheckValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHealthCheck.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HealthCheckValidationError{}

// Validate checks the field values on HealthCheck_Payload with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HealthCheck_Payload) Validate() error {
	if m == nil {
		return nil
	}

	switch m.Payload.(type) {

	case *HealthCheck_Payload_Text:

		if len(m.GetText()) < 1 {
			return HealthCheck_PayloadValidationError{
				Field:  "Text",
				Reason: "value length must be at least 1 bytes",
			}
		}

	case *HealthCheck_Payload_Binary:
		// no validation rules for Binary

	default:
		return HealthCheck_PayloadValidationError{
			Field:  "Payload",
			Reason: "value is required",
		}

	}

	return nil
}

// HealthCheck_PayloadValidationError is the validation error returned by
// HealthCheck_Payload.Validate if the designated constraints aren't met.
type HealthCheck_PayloadValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HealthCheck_PayloadValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHealthCheck_Payload.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HealthCheck_PayloadValidationError{}

// Validate checks the field values on HealthCheck_HttpHealthCheck with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HealthCheck_HttpHealthCheck) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Host

	if len(m.GetPath()) < 1 {
		return HealthCheck_HttpHealthCheckValidationError{
			Field:  "Path",
			Reason: "value length must be at least 1 bytes",
		}
	}

	if v, ok := interface{}(m.GetSend()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheck_HttpHealthCheckValidationError{
				Field:  "Send",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetReceive()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheck_HttpHealthCheckValidationError{
				Field:  "Receive",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	// no validation rules for ServiceName

	for idx, item := range m.GetRequestHeadersToAdd() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HealthCheck_HttpHealthCheckValidationError{
					Field:  fmt.Sprintf("RequestHeadersToAdd[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	// no validation rules for UseHttp2

	return nil
}

// HealthCheck_HttpHealthCheckValidationError is the validation error returned
// by HealthCheck_HttpHealthCheck.Validate if the designated constraints
// aren't met.
type HealthCheck_HttpHealthCheckValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HealthCheck_HttpHealthCheckValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHealthCheck_HttpHealthCheck.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HealthCheck_HttpHealthCheckValidationError{}

// Validate checks the field values on HealthCheck_TcpHealthCheck with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HealthCheck_TcpHealthCheck) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetSend()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheck_TcpHealthCheckValidationError{
				Field:  "Send",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	for idx, item := range m.GetReceive() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return HealthCheck_TcpHealthCheckValidationError{
					Field:  fmt.Sprintf("Receive[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	return nil
}

// HealthCheck_TcpHealthCheckValidationError is the validation error returned
// by HealthCheck_TcpHealthCheck.Validate if the designated constraints aren't met.
type HealthCheck_TcpHealthCheckValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HealthCheck_TcpHealthCheckValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHealthCheck_TcpHealthCheck.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HealthCheck_TcpHealthCheckValidationError{}

// Validate checks the field values on HealthCheck_RedisHealthCheck with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HealthCheck_RedisHealthCheck) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Key

	return nil
}

// HealthCheck_RedisHealthCheckValidationError is the validation error returned
// by HealthCheck_RedisHealthCheck.Validate if the designated constraints
// aren't met.
type HealthCheck_RedisHealthCheckValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HealthCheck_RedisHealthCheckValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHealthCheck_RedisHealthCheck.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HealthCheck_RedisHealthCheckValidationError{}

// Validate checks the field values on HealthCheck_GrpcHealthCheck with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HealthCheck_GrpcHealthCheck) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ServiceName

	return nil
}

// HealthCheck_GrpcHealthCheckValidationError is the validation error returned
// by HealthCheck_GrpcHealthCheck.Validate if the designated constraints
// aren't met.
type HealthCheck_GrpcHealthCheckValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HealthCheck_GrpcHealthCheckValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHealthCheck_GrpcHealthCheck.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HealthCheck_GrpcHealthCheckValidationError{}

// Validate checks the field values on HealthCheck_CustomHealthCheck with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HealthCheck_CustomHealthCheck) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetName()) < 1 {
		return HealthCheck_CustomHealthCheckValidationError{
			Field:  "Name",
			Reason: "value length must be at least 1 bytes",
		}
	}

	if v, ok := interface{}(m.GetConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HealthCheck_CustomHealthCheckValidationError{
				Field:  "Config",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// HealthCheck_CustomHealthCheckValidationError is the validation error
// returned by HealthCheck_CustomHealthCheck.Validate if the designated
// constraints aren't met.
type HealthCheck_CustomHealthCheckValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HealthCheck_CustomHealthCheckValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHealthCheck_CustomHealthCheck.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HealthCheck_CustomHealthCheckValidationError{}
