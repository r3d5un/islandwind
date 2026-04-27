package db

// ExplicitNull represents a tri-state nullable value used in query filters.
//
// States:
// - unspecified: filter should be ignored
// - null: filter should target SQL NULL
// - value: filter should apply with a concrete value
type ExplicitNull interface {
	IsSpecified() bool
	IsNull() bool
}
