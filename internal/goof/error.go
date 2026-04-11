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
	// message is a human-readable description of the error.
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

// Unwrap implements the errors.Wrapper interface, returning the underlying internal error.
func (e Error) Unwrap() error {
	return e.internal
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.message
}

// LogValue implements the slog.LogValuer interface. It returns a slog.Value containing the details
// of the error for structured logging.
func (e Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("code", e.code),
		slog.String("message", e.message),
		slog.Any("internal", e.internal),
		slog.Any("metadata", e.metadata),
	)
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Metadata() map[string]any {
	return e.metadata
}

func (e Error) Time() time.Time {
	return e.time
}

func (e Error) Service() string {
	return e.service
}

func (e Error) Internal() error {
	return e.internal
}
