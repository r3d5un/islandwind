package builder_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		stmt, args := builder.From("table1").Delete()
		require.NotEmpty(t, stmt)
		require.Empty(t, args)
		assert.Equal(t, "DELETE FROM table1;", stmt)
	})

	t.Run("Where", func(t *testing.T) {
		stmt, args := builder.From("table1").
			Where(builder.NewGenericPredicate("col1", builder.Equal, "test")).
			Delete()
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		require.Equal(t, "test", args["col1"])
		assert.Equal(t, "DELETE FROM table1 WHERE col1 = @col1;", stmt)
	})

	t.Run("ReturningSingleColumn", func(t *testing.T) {
		stmt, args := builder.From("table1").
			Where(builder.NewGenericPredicate("col1", builder.Equal, "test")).
			Returning("col1").
			Delete()
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		require.Equal(t, "test", args["col1"])
		assert.Equal(t, "DELETE FROM table1 WHERE col1 = @col1 RETURNING col1;", stmt)
	})

	t.Run("ReturningMultipleColumns", func(t *testing.T) {
		stmt, args := builder.From("table1").
			Where(builder.NewGenericPredicate("col1", builder.Equal, "test")).
			Returning("col1, col2").
			Delete()
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		require.Equal(t, "test", args["col1"])
		assert.Equal(t, "DELETE FROM table1 WHERE col1 = @col1 RETURNING col1, col2;", stmt)
	})
}
