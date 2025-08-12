package data

import (
	"reflect"

	"github.com/google/uuid"
)

// Metadata contains metadata about the results of a query.
type Metadata struct {
	// LastSeen is a [uuid.UUID] from the last item of the result the metadata object describes.
	LastSeen uuid.UUID `json:"lastSeen,omitzero"`
	// Next is true when there are more results in the dataset.
	Next bool `json:"next"`
	// ResponseLength is the number of results returned.
	ResponseLength int `json:"responseLength"`
}

// NewMetadata uses query results and filter to create a new [Metadata] instance.
//
// The following is required by the input:
//   - rows must be a slice of pointers.
//   - an individual row must have an ID field of [uuid.UUID].
func NewMetadata[T any](rows []T, filter Filter) Metadata {
	var metadata Metadata
	length := len(rows)

	if length >= filter.PageSize {
		metadata.Next = true
	}
	metadata.ResponseLength = length

	if length < 1 {
		return metadata
	}

	lastRow := rows[length-1]
	deref := reflect.Indirect(reflect.ValueOf(lastRow))
	if !deref.IsValid() {
		return metadata
	}

	field := deref.FieldByName("ID")
	if !field.IsValid() || !field.CanInterface() {
		return metadata
	}
	if id, ok := field.Interface().(uuid.UUID); ok {
		metadata.LastSeen = id
	}

	return metadata
}
