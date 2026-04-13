package db

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
)

type Query struct {
	q Queryable
}

type QueryBuilder struct {
	q Queryable
}

func New(q Queryable) QueryBuilder {
	return QueryBuilder{
		q: q,
	}
}

func NullEqual[T any](nullable sql.Null[T], column string) (string, pgx.NamedArgs) {
	if !nullable.Valid {
		return "", nil
	}

	return column + " = @" + column, pgx.NamedArgs{column: nullable.V}
}

func NullNotEqual[T any](nullable sql.Null[T], column string) (string, pgx.NamedArgs) {
	if !nullable.Valid {
		return "", nil
	}

	return column + " != @" + column, pgx.NamedArgs{column: nullable.V}
}

func NullLike(nullable sql.Null[string], column string) (string, pgx.NamedArgs) {
	if !nullable.Valid {
		return "", nil
	}

	return column + " LIKE @" + column, pgx.NamedArgs{column: "%" + nullable.V + "%"}
}

func NullGreater[T any](nullable sql.Null[T], column string) (string, pgx.NamedArgs) {
	if !nullable.Valid {
		return "", nil
	}

	return column + " > @" + column, pgx.NamedArgs{column: nullable.V}
}

func NullGreaterOrEqual[T any](nullable sql.Null[T], column string) (string, pgx.NamedArgs) {
	if !nullable.Valid {
		return "", nil
	}

	return column + " >= @" + column, pgx.NamedArgs{column: nullable.V}
}

func NullLess[T any](nullable sql.Null[T], column string) (string, pgx.NamedArgs) {
	if !nullable.Valid {
		return "", nil
	}

	return column + " < @" + column, pgx.NamedArgs{column: nullable.V}
}

func NullLessOrEqual[T any](nullable sql.Null[T], column string) (string, pgx.NamedArgs) {
	if !nullable.Valid {
		return "", nil
	}

	return column + " <= @" + column, pgx.NamedArgs{column: nullable.V}
}

func Equal[T any](val T, column string) (string, pgx.NamedArgs) {
	return column + " = @" + column, pgx.NamedArgs{column: val}
}

func NotEqual[T any](val T, column string) (string, pgx.NamedArgs) {
	return column + " != @" + column, pgx.NamedArgs{column: val}
}

func Greater[T any](val T, column string) (string, pgx.NamedArgs) {
	return column + " > @" + column, pgx.NamedArgs{column: val}
}

func GreaterOrEqual[T any](val T, column string) (string, pgx.NamedArgs) {
	return column + " >= @" + column, pgx.NamedArgs{column: val}
}

func Less[T any](val T, column string) (string, pgx.NamedArgs) {
	return column + " < @" + column, pgx.NamedArgs{column: val}
}

func LessOrEqual[T any](val T, column string) (string, pgx.NamedArgs) {
	return column + " <= @" + column, pgx.NamedArgs{column: val}
}

func Like(val string, column string) (string, pgx.NamedArgs) {
	return column + " LIKE @" + column, pgx.NamedArgs{column: "%" + val + "%"}
}
