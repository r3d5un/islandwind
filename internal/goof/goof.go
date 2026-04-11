// Package goof contains functionality for structured error handling.
package goof

// New creates a new Error with default values.
//
// Example:
//
//	err := errors.New("original error")
//	ge := goof.New(err)
func New(err error) error {
	return newBuilder().New(err)
}

// Wrap is an alias for New, provided for semantic clarity when wrapping an existing error.
//
// Example:
//
//	err := errors.New("database connection failed")
//	ge := goof.Wrap(err)
func Wrap(err error) error {
	return newBuilder().New(err)
}

// With starts building an error with metadata.
//
// Example:
//
//	err := goof.With("field", "email").
//		With("reason", "invalid format").
//		New(errors.New("validation failed"))
func With(key string, value any) ErrorBuilder {
	return newBuilder().With(key, value)
}

// Code starts building an error with a specific code.
//
// Example:
//
//	err := goof.Code("ERR_UNAUTHORIZED").New(errors.New("unauthorized access"))
func Code(code string) ErrorBuilder {
	return newBuilder().Code(code)
}

// Service starts building an error with a specific service.
//
// Example:
//
//	err := goof.Service("payment-gateway").New(errors.New("payment failed"))
func Service(service string) ErrorBuilder {
	return newBuilder().Service(service)
}

// Message starts building an error with a specific message.
//
// Example:
//
//	err := goof.Message("could not save file").New(errors.New("low-level disk error"))
func Message(message string) ErrorBuilder {
	return newBuilder().Message(message)
}
