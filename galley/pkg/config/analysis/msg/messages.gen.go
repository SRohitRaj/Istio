// GENERATED FILE -- DO NOT EDIT
//

package msg

import (
	"istio.io/istio/galley/pkg/config/analysis/diag"
	"istio.io/istio/galley/pkg/config/resource"
)

// InternalError returns a new diag.Message for message "Internal Error".
//
// There was an internal error in the toolchain. This is almost always a bug in the implementation.
func InternalError(entry *resource.Entry, detail string) diag.Message {
	return diag.NewMessage(
		diag.Error,
		diag.Code(1),
		entry.Origin,
		"Internal error: %v",
		detail,
	)
}

// NotYetImplemented returns a new diag.Message for message "Not Yet Implemented".
//
// A feature that the configuration is depending on is not implemented yet.
func NotYetImplemented(entry *resource.Entry, detail string) diag.Message {
	return diag.NewMessage(
		diag.Error,
		diag.Code(2),
		entry.Origin,
		"Not yet implemented: %s",
		detail,
	)
}

// ParseError returns a new diag.Message for message "Parse Error".
//
// There was a parse error during the parsing of the configuration text
func ParseError(entry *resource.Entry, detail string) diag.Message {
	return diag.NewMessage(
		diag.Warning,
		diag.Code(3),
		entry.Origin,
		"Parse error: %s",
		detail,
	)
}

// Deprecated returns a new diag.Message for message "Deprecated".
//
// A feature that the configuration is depending on is now deprecated.
func Deprecated(entry *resource.Entry, detail string) diag.Message {
	return diag.NewMessage(
		diag.Warning,
		diag.Code(4),
		entry.Origin,
		"Deprecated: %s",
		detail,
	)
}

// GatewayNotFound returns a new diag.Message for message "Gateway Not Found".
//
// The Gateway resource that the Virtual Service is referencing does not exist.
func GatewayNotFound(entry *resource.Entry, gateway string) diag.Message {
	return diag.NewMessage(
		diag.Error,
		diag.Code(20),
		entry.Origin,
		"Referenced Gateway not found: %q",
		gateway,
	)
}
