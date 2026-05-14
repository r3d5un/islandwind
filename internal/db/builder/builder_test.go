package builder_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	stmt, _ := builder.From("table").
		Select("col1", "col2")

	t.Log(stmt)
}

func TestWhere(t *testing.T) {
	t.Run("SingleColumn", func(t *testing.T) {
		stmt, _ := builder.From("table").
			Where(builder.NewGenericPredicate("id", builder.Equal, "1")).
			Select("col1")

		assert.Equal(t, "SELECT col1 FROM table WHERE id = @id;", stmt)
	})

	t.Run("MultipleColumns", func(t *testing.T) {
		stmt, _ := builder.From("table").
			Where(
				builder.NewGenericPredicate("id", builder.Equal, "1"),
				builder.NewGenericPredicate("col1", builder.Equal, "value"),
			).
			Select("col1", "col2")

		assert.Equal(t, "SELECT col1, col2 FROM table WHERE id = @id AND col1 = @col1;", stmt)
	})

	t.Run("And", func(t *testing.T) {
		t.Skip()
	})

	t.Run("Or", func(t *testing.T) {
		t.Skip()
	})

	t.Run("ComplexConditions", func(t *testing.T) {
		// SELECT col1 FROM table WHERE (id = @id AND (status = @status OR status = @status_1));
		t.Skip()
	})
}

func TestJoin(t *testing.T) {
	t.Run("SingleJoin", func(t *testing.T) {
		stmt, _ := builder.From("table1 a").
			Join("table2 b ON a.id = b.a_id").
			Select("a.col1", "b.col2")

		assert.Equal(t, "SELECT a.col1, b.col2 FROM table1 a JOIN table2 b ON a.id = b.a_id;", stmt)
	})

	t.Run("MultipleJoins", func(t *testing.T) {
		stmt, _ := builder.From("table1 a").
			Join("table2 b ON a.id = b.a_id").
			Join("table3 c ON b.id = c.b_id").
			Select("a.col1", "b.col2", "c.col3")

		assert.Equal(
			t,
			"SELECT a.col1, b.col2, c.col3 FROM table1 a JOIN table2 b ON a.id = b.a_id JOIN table3 c ON b.id = c.b_id;",
			stmt,
		)
	})
}

func TestOrderBy(t *testing.T) {
	t.Run("Asc", func(t *testing.T) {
		stmt, _ := builder.From("table").
			OrderBy(builder.OrderBy{Column: "col1", Order: builder.Asc}).
			Select("col1", "col2")

		assert.Equal(t, "SELECT col1, col2 FROM table ORDER BY col1 ASC;", stmt)
	})

	t.Run("Desc", func(t *testing.T) {
		stmt, _ := builder.From("table").
			OrderBy(builder.OrderBy{Column: "col1", Order: builder.Desc}).
			Select("col1", "col2")

		assert.Equal(t, "SELECT col1, col2 FROM table ORDER BY col1 DESC;", stmt)
	})
}
