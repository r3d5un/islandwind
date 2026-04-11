package goof

import (
	"maps"
	"time"
)

// ErrorBuilder is used to build a new Error with specific properties. It is immutable,
// meaning that each method call returns a new ErrorBuilder with the updated property.
type ErrorBuilder struct {
	// code is a short machine-readable identifier for the error.
	code string
	// message is a human-readable description of the error. The message is meant to be safe to
	// expose to external services.
	message string
	// metadata is a map of additional information about the error.
	metadata map[string]any
	// service refers to the service, module, domain, or category where the error occurred.
	service string
}

// newBuilder creates a new ErrorBuilder with default values.
func newBuilder() ErrorBuilder {
	return ErrorBuilder{
		code:     "ERROR",
		message:  "an error occurred",
		metadata: make(map[string]any),
		service:  "",
	}
}

// clone returns a shallow copy of the ErrorBuilder.
func (b ErrorBuilder) clone() ErrorBuilder {
	return ErrorBuilder{
		code:     b.code,
		message:  b.message,
		metadata: maps.Clone(b.metadata),
		service:  b.service,
	}
}

// New creates a new Error from the ErrorBuilder and the provided error.
// If the provided error is nil, it returns nil.
func (b ErrorBuilder) New(err error) error {
	if err == nil {
		return nil
	}

	return &Error{
		code:     b.code,
		message:  b.message,
		internal: err,
		metadata: maps.Clone(b.metadata),
		time:     time.Now(),
		service:  b.service,
	}
}

// With adds a key-value pair to the error's metadata. It is intended to be used to provide
// additional context or details about the error.
func (b ErrorBuilder) With(key string, value any) ErrorBuilder {
	clone := b.clone()
	if clone.metadata == nil {
		clone.metadata = make(map[string]any)
	}
	clone.metadata[key] = value
	return clone
}

// Code sets the error code for the ErrorBuilder.
func (b ErrorBuilder) Code(code string) ErrorBuilder {
	clone := b.clone()
	clone.code = code
	return clone
}

// Service sets the service for the ErrorBuilder.
func (b ErrorBuilder) Service(service string) ErrorBuilder {
	clone := b.clone()
	clone.service = service
	return clone
}

// Message sets the message for the ErrorBuilder.
func (b ErrorBuilder) Message(message string) ErrorBuilder {
	clone := b.clone()
	clone.message = message
	return clone
}
