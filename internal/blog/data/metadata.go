package data

import (
	"reflect"

	"github.com/google/uuid"
)

type Metadata struct {
	LastSeen       uuid.UUID `json:"lastSeen,omitzero"`
	Next           bool      `json:"next"`
	ResponseLength int       `json:"responseLength"`
}

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
