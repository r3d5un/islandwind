// Package goof contains functionality for structured error handling.
package goof

import (
	"log/slog"
)

// Error is a type that enabled structured error handling while implementing the error interface
// from the standard library.
type Error struct {
	// Code is a short machine-readable identifier for the error.
	Code string `json:"code"`
	// Message is a human-readable description of the error.
	Message string `json:"message"`
	// Internal is the raw error. This field should not be exposed to clients.
	Internal error `json:"-"`
	// Metadata is a map of additional information about the error.
	Metadata map[string]any `json:"-"`
}

// New creates a new Error instance.
func New(code, message string, internal error) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		Internal: internal,
	}
}

// Unwrap implements the error interface.
func (e *Error) Unwrap() error {
	return e.Internal
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.Message
}

// LogValue implements the slog.LogValuer interface. It returns a slog.Value containing the details
// of the error for structured logging.
func (e *Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("code", e.Code),
		slog.String("message", e.Message),
		slog.Any("internal", e.Internal),
		slog.Any("metadata", e.Metadata),
	)
}

func (e *Error) With(key string, value any) *Error {
	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}

	e.Metadata[key] = value
	return e
}
