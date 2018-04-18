// Code generated by protoc-gen-validate
// source: envoy/config/accesslog/v2/als.proto
// DO NOT EDIT!!!

package v2

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

// Validate checks the field values on TcpGrpcAccessLogConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *TcpGrpcAccessLogConfig) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetCommonConfig() == nil {
		return TcpGrpcAccessLogConfigValidationError{
			Field:  "CommonConfig",
			Reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetCommonConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TcpGrpcAccessLogConfigValidationError{
				Field:  "CommonConfig",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// TcpGrpcAccessLogConfigValidationError is the validation error returned by
// TcpGrpcAccessLogConfig.Validate if the designated constraints aren't met.
type TcpGrpcAccessLogConfigValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e TcpGrpcAccessLogConfigValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTcpGrpcAccessLogConfig.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = TcpGrpcAccessLogConfigValidationError{}

// Validate checks the field values on HttpGrpcAccessLogConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *HttpGrpcAccessLogConfig) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetCommonConfig() == nil {
		return HttpGrpcAccessLogConfigValidationError{
			Field:  "CommonConfig",
			Reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetCommonConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HttpGrpcAccessLogConfigValidationError{
				Field:  "CommonConfig",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// HttpGrpcAccessLogConfigValidationError is the validation error returned by
// HttpGrpcAccessLogConfig.Validate if the designated constraints aren't met.
type HttpGrpcAccessLogConfigValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e HttpGrpcAccessLogConfigValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHttpGrpcAccessLogConfig.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = HttpGrpcAccessLogConfigValidationError{}

// Validate checks the field values on CommonGrpcAccessLogConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CommonGrpcAccessLogConfig) Validate() error {
	if m == nil {
		return nil
	}

	if len(m.GetLogName()) < 1 {
		return CommonGrpcAccessLogConfigValidationError{
			Field:  "LogName",
			Reason: "value length must be at least 1 bytes",
		}
	}

	if m.GetGrpcService() == nil {
		return CommonGrpcAccessLogConfigValidationError{
			Field:  "GrpcService",
			Reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetGrpcService()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommonGrpcAccessLogConfigValidationError{
				Field:  "GrpcService",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// CommonGrpcAccessLogConfigValidationError is the validation error returned by
// CommonGrpcAccessLogConfig.Validate if the designated constraints aren't met.
type CommonGrpcAccessLogConfigValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e CommonGrpcAccessLogConfigValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCommonGrpcAccessLogConfig.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = CommonGrpcAccessLogConfigValidationError{}
