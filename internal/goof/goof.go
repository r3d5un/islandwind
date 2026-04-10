// Package goof contains functionality for structured error handling.
package goof

import (
	"errors"
	"log/slog"
	"maps"
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

type ErrorBuilder Error

func newBuilder() ErrorBuilder {
	return ErrorBuilder{
		code:     "",
		message:  "",
		internal: nil,
		metadata: make(map[string]any),
		time:     time.Now(),
		service:  "",
	}
}

func (e ErrorBuilder) clone() ErrorBuilder {
	clone := ErrorBuilder{
		code:     e.code,
		message:  e.message,
		internal: e.internal,
		metadata: maps.Clone(e.metadata),
		time:     e.time,
		service:  e.service,
	}

	if clone.metadata == nil {
		clone.metadata = make(map[string]any)
	}
	return clone
}

func (e ErrorBuilder) New(message string) error {
	clone := e.clone()
	clone.internal = errors.New(message)

	return Error(clone)
}

func (e ErrorBuilder) Wrap(err error) error {
	clone := e.clone()
	clone.internal = err

	return Error(clone)
}

// With adds a key-value pair to the error's metadata. It is intended to be used to provide
// additional context or details about the error.
func (e ErrorBuilder) With(key string, value any) ErrorBuilder {
	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}

	e.metadata[key] = value
	return e
}

func New(message string) error {
	return newBuilder().New(message)
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	return newBuilder().Wrap(err)
}

func With(key string, value any) ErrorBuilder {
	return newBuilder().With(key, value)
}
