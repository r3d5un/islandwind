package builder_test

import (
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/oapi-codegen/nullable"
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
			Where("id = @id", pgx.NamedArgs{"id": "1"}).
			Select("col1")

		assert.Equal(t, "SELECT col1 FROM table WHERE id = @id;", stmt)
	})

	t.Run("MultipleColumns", func(t *testing.T) {
		stmt, _ := builder.From("table").
			Where("id = @id", pgx.NamedArgs{"id": "1"}).
			Select("col1", "col2")

		assert.Equal(t, "SELECT col1, col2 FROM table WHERE id = @id;", stmt)
	})

	t.Run("ComplexConditions", func(t *testing.T) {
		where, whereArgs := builder.And(
			builder.Equal(1, "id"),
			builder.Or(
				builder.Equal("active", "status"),
				builder.Equal("pending", "status"),
			),
		)()

		stmt, args := builder.From("table").
			Where(where, whereArgs).
			Select("col1")

		expectedSQL := "SELECT col1 FROM table WHERE (id = @id AND (status = @status OR status = @status_1));"
		assert.Equal(t, expectedSQL, stmt)
		assert.Equal(t, 1, args["id"])
		assert.Equal(t, "active", args["status"])
		assert.Equal(t, "pending", args["status_1"])
	})
}

func TestWhereNull(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		stmt, args := builder.From("table").
			WhereNull(sql.Null[string]{Valid: true, V: "abc"}, "id = @id", pgx.NamedArgs{"id": "abc"}).
			Select("id")

		assert.Equal(t, "SELECT id FROM table WHERE id = @id;", stmt)
		assert.Equal(t, "abc", args["id"])
	})

	t.Run("InvalidSkipped", func(t *testing.T) {
		stmt, args := builder.From("table").
			WhereNull(sql.Null[string]{Valid: false}, "id = @id", pgx.NamedArgs{"id": "abc"}).
			Select("id")

		assert.Equal(t, "SELECT id FROM table;", stmt)
		assert.Empty(t, args)
	})

	t.Run("ExplicitNullableValue", func(t *testing.T) {
		stmt, args := builder.From("table").
			WhereNull(nullable.NewNullableWithValue("abc"), "id = @id", pgx.NamedArgs{"id": "abc"}).
			Select("id")

		assert.Equal(t, "SELECT id FROM table WHERE id = @id;", stmt)
		assert.Equal(t, "abc", args["id"])
	})

	t.Run("ExplicitNullableNullSkipped", func(t *testing.T) {
		stmt, args := builder.From("table").
			WhereNull(nullable.NewNullNullable[string](), "id = @id", pgx.NamedArgs{"id": "abc"}).
			Select("id")

		assert.Equal(t, "SELECT id FROM table;", stmt)
		assert.Empty(t, args)
	})
}

func TestWhereExplicitNull(t *testing.T) {
	t.Run("UnspecifiedSkipped", func(t *testing.T) {
		stmt, args := builder.From("table").
			WhereExplicitNull(
				nullable.Nullable[string]{},
				"deleted_at = @deleted_at",
				pgx.NamedArgs{"deleted_at": "2026-01-01"},
				"deleted_at IS NULL",
			).
			Select("id")

		assert.Equal(t, "SELECT id FROM table;", stmt)
		assert.Empty(t, args)
	})

	t.Run("NullClause", func(t *testing.T) {
		stmt, args := builder.From("table").
			WhereExplicitNull(
				nullable.NewNullNullable[string](),
				"deleted_at = @deleted_at",
				pgx.NamedArgs{"deleted_at": "2026-01-01"},
				"deleted_at IS NULL",
			).
			Select("id")

		assert.Equal(t, "SELECT id FROM table WHERE deleted_at IS NULL;", stmt)
		assert.Empty(t, args)
	})

	t.Run("ValueClause", func(t *testing.T) {
		stmt, args := builder.From("table").
			WhereExplicitNull(
				nullable.NewNullableWithValue("2026-01-01"),
				"deleted_at = @deleted_at",
				pgx.NamedArgs{"deleted_at": "2026-01-01"},
				"deleted_at IS NULL",
			).
			Select("id")

		assert.Equal(t, "SELECT id FROM table WHERE deleted_at = @deleted_at;", stmt)
		assert.Equal(t, "2026-01-01", args["deleted_at"])
	})

	t.Run("TypedNilPointerSkipped", func(t *testing.T) {
		var n *nullable.Nullable[string]

		stmt, args := builder.From("table").
			WhereExplicitNull(
				n,
				"deleted_at = @deleted_at",
				pgx.NamedArgs{"deleted_at": "2026-01-01"},
				"deleted_at IS NULL",
			).
			Select("id")

		assert.Equal(t, "SELECT id FROM table;", stmt)
		assert.Empty(t, args)
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

func TestNullFunctions(t *testing.T) {
	t.Run("NullEqual Valid", func(t *testing.T) {
		clause, args := builder.NullEqual(sql.Null[int]{Valid: true, V: 10}, "col")()
		assert.Equal(t, "col = @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})
	t.Run("NullEqual Invalid", func(t *testing.T) {
		clause, args := builder.NullEqual(sql.Null[int]{Valid: false}, "col")()
		assert.Equal(t, "", clause)
		assert.Nil(t, args)
	})

	t.Run("NullNotEqual Valid", func(t *testing.T) {
		clause, args := builder.NullNotEqual(sql.Null[int]{Valid: true, V: 10}, "col")()
		assert.Equal(t, "col != @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullLike Valid", func(t *testing.T) {
		clause, args := builder.NullLike(sql.Null[string]{Valid: true, V: "test"}, "col")()
		assert.Equal(t, "col LIKE @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": "%test%"}, args)
	})

	t.Run("NullGreater Valid", func(t *testing.T) {
		clause, args := builder.NullGreater(sql.Null[int]{Valid: true, V: 10}, "col")()
		assert.Equal(t, "col > @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullGreaterOrEqual Valid", func(t *testing.T) {
		clause, args := builder.NullGreaterOrEqual(sql.Null[int]{Valid: true, V: 10}, "col")()
		assert.Equal(t, "col >= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullLess Valid", func(t *testing.T) {
		clause, args := builder.NullLess(sql.Null[int]{Valid: true, V: 10}, "col")()
		assert.Equal(t, "col < @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullLessOrEqual Valid", func(t *testing.T) {
		clause, args := builder.NullLessOrEqual(sql.Null[int]{Valid: true, V: 10}, "col")()
		assert.Equal(t, "col <= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})
}

func TestNotNullFunctions(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		clause, args := builder.Equal(10, "col")()
		assert.Equal(t, "col = @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NotEqual", func(t *testing.T) {
		clause, args := builder.NotEqual(10, "col")()
		assert.Equal(t, "col != @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("Greater", func(t *testing.T) {
		clause, args := builder.Greater(10, "col")()
		assert.Equal(t, "col > @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("GreaterOrEqual", func(t *testing.T) {
		clause, args := builder.GreaterOrEqual(10, "col")()
		assert.Equal(t, "col >= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("Less", func(t *testing.T) {
		clause, args := builder.Less(10, "col")()
		assert.Equal(t, "col < @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("LessOrEqual", func(t *testing.T) {
		clause, args := builder.LessOrEqual(10, "col")()
		assert.Equal(t, "col <= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("Like", func(t *testing.T) {
		clause, args := builder.Like("test", "col")()
		assert.Equal(t, "col LIKE @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": "%test%"}, args)
	})
}
