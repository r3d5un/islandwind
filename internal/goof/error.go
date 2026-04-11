package goof

import (
	"log/slog"
	"time"
)

var (
	_ error          = (*Error)(nil)
	_ slog.LogValuer = (*Error)(nil)
)

// Error is a type that enabled structured error handling while implementing the error interface
// from the standard library.
type Error struct {
	// code is a short machine-readable identifier for the error.
	code string
	// message is a human-readable description of the error. The message is meant to be safe to
	// expose to external services.
	message string
	// internal is the raw error. This field should not be exposed to clients.
	internal error
	// metadata is a map of additional information about the error.
	metadata map[string]any
	// time is the time at which the error occurred.
	time time.Time
	// service refers to the service, module, domain, or category where the error occurred.
	service string
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	return e.internal
}

// Error returns the error message from the underlying error.
func (e *Error) Error() string {
	return e.internal.Error()
}

// LogValue returns the error's properties for structured logging.
func (e *Error) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("code", e.code),
		slog.String("message", e.message),
		slog.Time("time", e.time),
	}

	if e.service != "" {
		attrs = append(attrs, slog.String("service", e.service))
	}

	if e.internal != nil {
		attrs = append(attrs, slog.Any("internal", e.internal))
	}

	if len(e.metadata) > 0 {
		attrs = append(attrs, slog.Any("metadata", e.metadata))
	}

	return slog.GroupValue(attrs...)
}

// Code returns the error code.
func (e *Error) Code() string {
	return e.code
}

// Message returns the human-readable message.
func (e *Error) Message() string {
	return e.message
}

// Metadata returns the error's metadata.
func (e *Error) Metadata() map[string]any {
	return e.metadata
}

// Time returns the time the error occurred.
func (e *Error) Time() time.Time {
	return e.time
}

// Service returns the service name associated with the error.
func (e *Error) Service() string {
	return e.service
}

// Internal returns the underlying error.
func (e *Error) Internal() error {
	return e.internal
}
