package builder_test

import (
	"database/sql"
	"testing"

	"github.com/oapi-codegen/nullable"
	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
)

func TestNullableAssignment(t *testing.T) {
	column := "column1"
	value := "value"

	t.Run("WithValue", func(t *testing.T) {
		nullableValue := nullable.NewNullableWithValue(value)
		predicate := builder.NewNullableAssignment(column, nullableValue)
		assert.Equal(t, "column1 = @column1", predicate.Text)
		assert.NotEmpty(t, predicate.Args)
		assert.Equal(t, value, predicate.Args[column])
	})

	t.Run("ExplicitNull", func(t *testing.T) {
		explicitNull := nullable.NewNullNullable[string]()
		explicitNull.SetNull()
		predicate := builder.NewNullableAssignment(column, explicitNull)
		assert.Equal(t, "column1 = @column1", predicate.Text)
		assert.NotEmpty(t, predicate.Args)
		assert.Nil(t, predicate.Args[column])
	})

	t.Run("NotSpecified", func(t *testing.T) {
		notSpecifiedValue := nullable.Nullable[string]{}
		predicate := builder.NewNullableAssignment(column, notSpecifiedValue)
		assert.Empty(t, predicate.Text)
		assert.Empty(t, predicate.Args)
	})
}

func TestNullAssignment(t *testing.T) {
	column := "column1"

	t.Run("Valid", func(t *testing.T) {
		value := sql.Null[string]{V: "string", Valid: true}
		predicate := builder.NewNullAssignment(column, value)
		assert.Equal(t, "column1 = @column1", predicate.Text)
		assert.NotEmpty(t, predicate.Args)
		assert.Equal(t, value.V, predicate.Args[column])
	})

	t.Run("Invalid", func(t *testing.T) {
		value := sql.Null[string]{Valid: false}
		predicate := builder.NewNullAssignment(column, value)
		assert.Empty(t, predicate.Text)
		assert.Empty(t, predicate.Args)
	})
}

func TestNewGenericAssignment(t *testing.T) {
	column := "column1"
	value := "value"

	predicate := builder.NewGenericAssignment(column, value)
	assert.Equal(t, "column1 = @column1", predicate.Text)
	assert.NotEmpty(t, predicate.Args)
	assert.Equal(t, value, predicate.Args[column])
}
