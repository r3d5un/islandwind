package builder_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Run("Update", func(t *testing.T) {
		stmt, args, err := builder.
			Update("table1").
			Set(
				builder.NewGenericAssignment("column1", "value"),
			)
		require.NoError(t, err)
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		assert.Equal(t, "UPDATE table1 SET column1 = @assignment_column1;", stmt)
	})

	t.Run("UpdateMultipleColumns", func(t *testing.T) {
		stmt, args, err := builder.
			Update("table1").
			Set(
				builder.NewGenericAssignment("column1", "value"),
				builder.NewGenericAssignment("column2", "value"),
			)
		require.NoError(t, err)
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		assert.Equal(t, "UPDATE table1 SET column1 = @assignment_column1, column2 = @assignment_column2;", stmt)
	})

	t.Run("TableNotSet", func(t *testing.T) {
		stmt, args, err := builder.
			Update("").
			Set(
				builder.NewGenericAssignment("column1", "value"),
			)
		require.Error(t, err)
		require.Empty(t, stmt)
		require.Empty(t, args)
	})

	t.Run("NoAssignmentsSet", func(t *testing.T) {
		stmt, args, err := builder.
			Update("table1").
			Set()
		require.Error(t, err)
		require.Empty(t, stmt)
		require.Empty(t, args)
	})

	t.Run("Where", func(t *testing.T) {
		stmt, args, err := builder.
			Update("table1").
			Where(builder.NewGenericPredicate("column1", builder.Equal, "test")).
			Set(
				builder.NewGenericAssignment("column1", "value"),
			)
		require.NoError(t, err)
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		assert.Equal(t, "UPDATE table1 SET column1 = @assignment_column1 WHERE column1 = @predicate_column1;", stmt)
	})

	t.Run("Returning", func(t *testing.T) {
		stmt, args, err := builder.
			Update("table1").
			Returning("column1").
			Set(
				builder.NewGenericAssignment("column1", "value"),
			)
		require.NoError(t, err)
		require.NotEmpty(t, stmt)
		require.NotEmpty(t, args)
		assert.Equal(t, "UPDATE table1 SET column1 = @assignment_column1 RETURNING column1;", stmt)
	})
}
