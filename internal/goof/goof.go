// Package goof contains functionality for structured error handling.
package goof

// New creates a new Error with default values.
func New(err error) error {
	return newBuilder().New(err)
}

// Wrap is an alias for New, provided for semantic clarity when wrapping an existing error.
func Wrap(err error) error {
	return newBuilder().New(err)
}

// With starts building an error with metadata.
func With(key string, value any) ErrorBuilder {
	return newBuilder().With(key, value)
}

// Code starts building an error with a specific code.
func Code(code string) ErrorBuilder {
	return newBuilder().Code(code)
}

// Service starts building an error with a specific service.
func Service(service string) ErrorBuilder {
	return newBuilder().Service(service)
}

// Message starts building an error with a specific message.
func Message(message string) ErrorBuilder {
	return newBuilder().Message(message)
}
