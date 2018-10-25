// Code generated by protoc-gen-validate
// source: envoy/config/retry/other_priority/other_priority_config.proto
// DO NOT EDIT!!!

package envoy_config_retry_other_priority

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

// Validate checks the field values on OtherPriorityConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *OtherPriorityConfig) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for UpdateFrequency

	return nil
}

// OtherPriorityConfigValidationError is the validation error returned by
// OtherPriorityConfig.Validate if the designated constraints aren't met.
type OtherPriorityConfigValidationError struct {
	Field  string
	Reason string
	Cause  error
	Key    bool
}

// Error satisfies the builtin error interface
func (e OtherPriorityConfigValidationError) Error() string {
	cause := ""
	if e.Cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.Cause)
	}

	key := ""
	if e.Key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOtherPriorityConfig.%s: %s%s",
		key,
		e.Field,
		e.Reason,
		cause)
}

var _ error = OtherPriorityConfigValidationError{}
