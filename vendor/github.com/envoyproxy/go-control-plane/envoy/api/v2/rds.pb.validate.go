// Code generated by protoc-gen-validate
// source: envoy/api/v2/rds.proto
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

// Validate checks the field values on RouteConfiguration with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RouteConfiguration) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Name

	for idx, item := range m.GetVirtualHosts() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RouteConfigurationValidationError{
					Field:  fmt.Sprintf("VirtualHosts[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetResponseHeadersToAdd() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RouteConfigurationValidationError{
					Field:  fmt.Sprintf("ResponseHeadersToAdd[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetRequestHeadersToAdd() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RouteConfigurationValidationError{
					Field:  fmt.Sprintf("RequestHeadersToAdd[%v]", idx),
					Reason: "embedded message failed validation",
					Cause:  err,
				}
			}
		}

	}

	if v, ok := interface{}(m.GetValidateClusters()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RouteConfigurationValidationError{
				Field:  "ValidateClusters",
				Reason: "embedded message failed validation",
				Cause:  err,
			}
		}
	}

	return nil
}

// RouteConfigurationValidationError is the validation error returned by
// RouteConfiguration.Validate if the designated constraints aren't met.
type RouteConfigurationValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e RouteConfigurationValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRouteConfiguration.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = RouteConfigurationValidationError{}
