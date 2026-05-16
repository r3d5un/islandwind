package builder_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		stmt, args, err := builder.Insert(
			builder.Tuple{"column1": {V: "test", Valid: true}},
		).Into("table1")
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		require.NoError(t, err)
		assert.Equal(t, "INSERT INTO table1 (column1) VALUES (@column1_0);", stmt)
	})

	t.Run("MultipleColumns", func(t *testing.T) {
		stmt, args, err := builder.Insert(
			builder.Tuple{"column1": {V: "test", Valid: true}, "column2": {V: "test", Valid: true}},
		).Into("table1")
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		require.NoError(t, err)

		expected := "INSERT INTO table1 (column1, column2) VALUES (@column1_0, @column2_0);"
		assert.Equal(t, expected, stmt)
	})

	t.Run("MultipleRows", func(t *testing.T) {
		stmt, args, err := builder.Insert(
			builder.Tuple{"column1": {V: "test", Valid: true}, "column2": {V: "test", Valid: true}},
			builder.Tuple{"column1": {V: "test", Valid: true}, "column2": {V: "test", Valid: true}},
		).Into("table1")
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		require.NoError(t, err)

		expected := "INSERT INTO table1 (column1, column2) VALUES (@column1_0, @column2_0), (@column1_1, @column2_1);"
		assert.Equal(t, expected, stmt)
	})

	t.Run("Returning", func(t *testing.T) {
		stmt, args, err := builder.
			Insert(
				builder.Tuple{"column1": {V: "test", Valid: true}},
			).
			Returning("column1").
			Into("table1")
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		require.NoError(t, err)
		assert.Equal(t, "INSERT INTO table1 (column1) VALUES (@column1_0) RETURNING column1;", stmt)
	})

	t.Run("NoTuples", func(t *testing.T) {
		stmt, args, err := builder.
			Insert().
			Into("table1")
		require.Empty(t, stmt)
		require.Empty(t, args)
		require.ErrorIs(t, err, builder.ErrNoInsertTuples)
	})
}
