package goof

import (
	"errors"
	"maps"
	"time"
)

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
	clone.message = message

	return Error(clone)
}

func (e ErrorBuilder) Wrap(err error) error {
	clone := e.clone()
	clone.internal = err
	if err != nil {
		clone.message = err.Error()
	}

	return Error(clone)
}

// With adds a key-value pair to the error's metadata. It is intended to be used to provide
// additional context or details about the error.
func (e ErrorBuilder) With(key string, value any) ErrorBuilder {
	clone := e.clone()
	clone.metadata[key] = value
	return clone
}

func (e ErrorBuilder) Code(code string) ErrorBuilder {
	clone := e.clone()
	clone.code = code

	return clone
}

func (e ErrorBuilder) Service(service string) ErrorBuilder {
	clone := e.clone()
	clone.service = service

	return clone
}
