package nullable

import (
	"bytes"
	"encoding/json"
	"errors"
)

var (
	ErrNullValue    = errors.New("value is nil")
	ErrNotSpecified = errors.New("value not specified")
)

// Nullable is a generic type that is able to represent three different states:
//
//   - Unspecified: `map[bool]{}`
//   - Null: `map[false]T`
//   - Value: `map[true]T`
//
// The intended use is for fields in JSON structs that can be explicitly set to null
// by the caller in a JSON request body. If the field is expected to be optional, add the
// `omitempty` JSON tag.
type Nullable[T any] map[bool]T

// NewNullableValue creates a new instance of Nullable with a given type and value.
//
//	`val := NewNullableValue[string]("example")`
func NewNullableValue[T any](t T) Nullable[T] {
	var nullable Nullable[T]
	nullable.Set(t)
	return nullable
}

// NewNullableNull creates a new instance of Nullable explicitly set to null.
func NewNullableNull[T any]() Nullable[T] {
	var nullable Nullable[T]
	nullable.Null()
	return nullable
}

// NewNullableUnspecified creates a new instance of Nullable that is deliberately set without any
// value.
func NewNullableUnspecified[T any]() Nullable[T] {
	var nullable Nullable[T]
	nullable.SetUnspecified()
	return nullable
}

// IsNull checks if the Nullable is null
func (t Nullable[T]) IsNull() bool {
	_, ok := t[false]
	return ok
}

// IsSpecified checks if the Nullable is explicitly set
func (t Nullable[T]) IsSpecified() bool {
	return len(t) != 0
}

// SetUnspecified explicitly sets the Nullable as unspecified. This would be equivalent to a
// missing field in a JSON object.
func (t *Nullable[T]) SetUnspecified() {
	*t = map[bool]T{}
}

// Set sets the value of the Nullable explicitly
func (t *Nullable[T]) Set(value T) {
	*t = map[bool]T{true: value}
}

// Null explicitly sets the Nullable to null. This would be the equivalent to a JSON object
// `{"field": null}`.
func (t *Nullable[T]) Null() {
	var null T
	*t = map[bool]T{false: null}
}

// Get returns the value if specified or set. If the Nullable is null or not specified a
// ErrNullValue or ErrNotSpecified error is returned respectively.
func (t Nullable[T]) Get() (T, error) {
	var null T
	if t.IsNull() {
		return null, ErrNullValue
	}
	if !t.IsSpecified() {
		return null, ErrNotSpecified
	}

	return t[true], nil
}

// MarshalJSON implements the json.Marshaler interface used by the standard library to marshal
// JSON objects.
func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return []byte("null"), nil
	}

	return json.Marshal(t[true])
}

// UnmarshalJSON implements the json.Unmarshaler interface used by the standard library to
// unmarshal JSON objects.
func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Null()
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.Set(v)

	return nil
}
