package builder_test

import (
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/oapi-codegen/nullable"
	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
)

func TestNewNullablePredicate(t *testing.T) {
	column := "column1"
	value := "value"
	namedArgKey := "predicate_" + column

	t.Run("WithValue", func(t *testing.T) {
		nullableValue := nullable.NewNullableWithValue(value)
		predicate := builder.NewNullablePredicate(column, builder.Equal, nullableValue)
		assert.Equal(t, "column1 = @predicate_column1", predicate.Text)
		assert.NotEmpty(t, predicate.Arg)
		assert.Equal(t, value, predicate.Arg[namedArgKey])
	})

	t.Run("ExplicitNull", func(t *testing.T) {
		explicitNull := nullable.NewNullNullable[string]()
		explicitNull.SetNull()
		predicate := builder.NewNullablePredicate(column, builder.Equal, explicitNull)
		assert.Equal(t, "column1 = @predicate_column1", predicate.Text)
		assert.NotEmpty(t, predicate.Arg)
		assert.Nil(t, predicate.Arg[namedArgKey])
	})

	t.Run("NotSpecified", func(t *testing.T) {
		notSpecifiedValue := nullable.Nullable[string]{}
		predicate := builder.NewNullablePredicate(column, builder.Equal, notSpecifiedValue)
		assert.Empty(t, predicate.Text)
		assert.Empty(t, predicate.Arg)
	})
}

func TestNewNullPredicate(t *testing.T) {
	column := "column1"
	namedArgKey := "predicate_" + column

	t.Run("Valid", func(t *testing.T) {
		value := sql.Null[string]{V: "string", Valid: true}
		predicate := builder.NewNullPredicate(column, builder.Equal, value)
		assert.Equal(t, "column1 = @predicate_column1", predicate.Text)
		assert.NotEmpty(t, predicate.Arg)
		assert.Equal(t, value.V, predicate.Arg[namedArgKey])
	})

	t.Run("Invalid", func(t *testing.T) {
		value := sql.Null[string]{Valid: false}
		predicate := builder.NewNullPredicate(column, builder.Equal, value)
		assert.Empty(t, predicate.Text)
		assert.Empty(t, predicate.Arg)
	})
}

func TestNewGenericPredicate(t *testing.T) {
	column := "column1"
	value := "value"
	namedArgKey := "predicate_" + column

	predicate := builder.NewGenericPredicate(column, builder.Equal, value)
	assert.Equal(t, "column1 = @"+namedArgKey, predicate.Text)
	assert.NotEmpty(t, predicate.Arg)
	assert.Equal(t, value, predicate.Arg[namedArgKey])
}

func TestNewPredicate(t *testing.T) {
	predicate := builder.NewPredicate("column1 = @column1", pgx.NamedArgs{"column1": "value"})
	assert.Equal(t, "column1 = @column1", predicate.Text)
	assert.NotEmpty(t, predicate.Arg)
	assert.Equal(t, "value", predicate.Arg["column1"])
}
