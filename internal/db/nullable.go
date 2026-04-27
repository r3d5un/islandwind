package db

import (
	"github.com/r3d5un/islandwind/internal/nullable"
)

var (
	_ ExplicitNull = nullable.Nullable[string]{}
	_ ExplicitNull = (*nullable.Nullable[string])(nil)
)

// ExplicitNull represents a tri-state nullable value used in query filters.
//
// States:
// - Unspecified: filter should be ignored
// - Null: filter should target SQL NULL
// - Value: filter should apply with a concrete value
type ExplicitNull interface {
	IsSpecified() bool
	IsNull() bool
}
