package error

import (
	"errors"
	"log/slog"
)

type Error struct {
	Code     string         `json:"code"`
	Message  string         `json:"message"`
	Internal error          `json:"internal"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Internal error `json:"-"`
}

func New(code, message string, internal error, metadata map[string]any) *Error {
	return &Error{
		Code:     code,
		Message:  message,
		Internal: errors.Join(internal),
		Metadata: metadata,
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
		slog.String("internal", e.Internal.Error()),
		slog.Any("metadata", e.Metadata),
	)
}
