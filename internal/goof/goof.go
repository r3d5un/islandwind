// Package goof contains functionality for structured error handling.
package goof

func New(err error) error {
	if err == nil {
		return nil
	}

	return newBuilder().New(err)
}

// Wrap is an alias for New, provided for semantic clarity when wrapping an existing error.
func Wrap(err error) error {
	return newBuilder().New(err)
}

func With(key string, value any) ErrorBuilder {
	return newBuilder().With(key, value)
}

func Code(code string) ErrorBuilder {
	return newBuilder().Code(code)
}

func Service(service string) ErrorBuilder {
	return newBuilder().Service(service)
}

func Message(message string) ErrorBuilder {
	return newBuilder().Message(message)
}
