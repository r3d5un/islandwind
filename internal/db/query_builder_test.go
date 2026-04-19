package db_test

import (
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	stmt, _ := db.From("table").
		Select("col1", "col2")

	t.Log(stmt)
}

func TestWhere(t *testing.T) {
	stmt, _ := db.From("table").
		Where("@id = id", pgx.NamedArgs{"id": "1"}).
		Select("col1", "col2")

	assert.Equal(t, "SELECT col1, col2 FROM table WHERE @id = id;", stmt)
}

func TestOrderBy(t *testing.T) {
	stmt, _ := db.From("table").
		OrderBy(db.OrderBy{Column: "col1", Order: db.Asc}).
		Select("col1", "col2")

	assert.Equal(t, "SELECT col1, col2 FROM table ORDER BY col1 ASC;", stmt)
}

func TestNullFunctions(t *testing.T) {
	t.Run("NullEqual Valid", func(t *testing.T) {
		clause, args := db.NullEqual(sql.Null[int]{Valid: true, V: 10}, "col")
		assert.Equal(t, "col = @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})
	t.Run("NullEqual Invalid", func(t *testing.T) {
		clause, args := db.NullEqual(sql.Null[int]{Valid: false}, "col")
		assert.Equal(t, "", clause)
		assert.Nil(t, args)
	})

	t.Run("NullNotEqual Valid", func(t *testing.T) {
		clause, args := db.NullNotEqual(sql.Null[int]{Valid: true, V: 10}, "col")
		assert.Equal(t, "col != @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullLike Valid", func(t *testing.T) {
		clause, args := db.NullLike(sql.Null[string]{Valid: true, V: "test"}, "col")
		assert.Equal(t, "col LIKE @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": "%test%"}, args)
	})

	t.Run("NullGreater Valid", func(t *testing.T) {
		clause, args := db.NullGreater(sql.Null[int]{Valid: true, V: 10}, "col")
		assert.Equal(t, "col > @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullGreaterOrEqual Valid", func(t *testing.T) {
		clause, args := db.NullGreaterOrEqual(sql.Null[int]{Valid: true, V: 10}, "col")
		assert.Equal(t, "col >= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullLess Valid", func(t *testing.T) {
		clause, args := db.NullLess(sql.Null[int]{Valid: true, V: 10}, "col")
		assert.Equal(t, "col < @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NullLessOrEqual Valid", func(t *testing.T) {
		clause, args := db.NullLessOrEqual(sql.Null[int]{Valid: true, V: 10}, "col")
		assert.Equal(t, "col <= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})
}

func TestNotNullFunctions(t *testing.T) {
	t.Run("Equal", func(t *testing.T) {
		clause, args := db.Equal(10, "col")
		assert.Equal(t, "col = @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("NotEqual", func(t *testing.T) {
		clause, args := db.NotEqual(10, "col")
		assert.Equal(t, "col != @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("Greater", func(t *testing.T) {
		clause, args := db.Greater(10, "col")
		assert.Equal(t, "col > @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("GreaterOrEqual", func(t *testing.T) {
		clause, args := db.GreaterOrEqual(10, "col")
		assert.Equal(t, "col >= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("Less", func(t *testing.T) {
		clause, args := db.Less(10, "col")
		assert.Equal(t, "col < @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("LessOrEqual", func(t *testing.T) {
		clause, args := db.LessOrEqual(10, "col")
		assert.Equal(t, "col <= @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": 10}, args)
	})

	t.Run("Like", func(t *testing.T) {
		clause, args := db.Like("test", "col")
		assert.Equal(t, "col LIKE @col", clause)
		assert.Equal(t, pgx.NamedArgs{"col": "%test%"}, args)
	})
}
