// Package goof contains functionality for structured error handling.
package goof

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

func Code(code string) ErrorBuilder {
	return newBuilder().Code(code)
}

func Service(service string) ErrorBuilder {
	return newBuilder().Service(service)
}
