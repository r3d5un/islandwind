package db

import (
	"database/sql"
	"maps"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/r3d5un/islandwind/internal/ensure"
)

type Order string

const (
	Asc  Order = "ASC"
	Desc Order = "DESC"
)

type OrderBy struct {
	Column string
	Order  Order
}

type QueryBuilder struct {
	orderBy          []OrderBy
	whereClauses     []string
	namedArgs        pgx.NamedArgs
	returning        bool
	returningColumns []string
	from             string
}

func (qb QueryBuilder) clone() QueryBuilder {
	return QueryBuilder{
		orderBy:      slices.Clone(qb.orderBy),
		whereClauses: slices.Clone(qb.whereClauses),
		namedArgs:    maps.Clone(qb.namedArgs),
		returning:    qb.returning,
		from:         qb.from,
	}
}

func newQueryBuilder() QueryBuilder {
	return QueryBuilder{
		orderBy:      make([]OrderBy, 0),
		whereClauses: make([]string, 0),
		namedArgs:    make(pgx.NamedArgs, 0),
		returning:    false,
		from:         "",
	}
}

func (qb QueryBuilder) OrderBy(column string, order Order) QueryBuilder {
	if column == "" {
		return qb
	}

	clone := qb.clone()
	return clone.OrderBy(column, order)
}

func (qb QueryBuilder) From(from string) QueryBuilder {
	clone := qb.clone()
	clone.from = from
	return clone
}

func (qb QueryBuilder) Where(condition string, arg pgx.NamedArgs) QueryBuilder {
	if condition == "" || len(arg) == 0 {
		return qb
	}

	clone := qb.clone()
	maps.Copy(clone.namedArgs, arg)
	return clone.Where(condition, arg)
}

func (qb QueryBuilder) Returning(returning bool) QueryBuilder {
	clone := qb.clone()
	clone.returning = returning
	return clone
}

func (qb QueryBuilder) Select(cols ...string) (string, pgx.NamedArgs) {
	qb.returningColumns = append(qb.returningColumns, cols...)
	return qb.selectExp()
}

func (qb QueryBuilder) selectExp() (string, pgx.NamedArgs) {
	var builder strings.Builder
	builder.WriteString("SELECT ")
	for i, column := range qb.returningColumns {
		if _, err := builder.WriteString(column); err != nil {
			ensure.NoError(err, "select builder failed to write column")
		}
		if i != len(qb.returningColumns)-1 {
			if _, err := builder.WriteString(", "); err != nil {
				ensure.NoError(err, "select builder failed to write comma")
			}
		}
	}
	builder.WriteString(" FROM ")
	builder.WriteString(qb.from)

	if len(qb.whereClauses) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(qb.whereClauses, " AND "))
	}

	if len(qb.orderBy) > 0 {
		for i, orderBy := range qb.orderBy {
			builder.WriteString(" ORDER BY ")
			builder.WriteString(orderBy.Column)
			builder.WriteString(" ")
			builder.WriteString(string(orderBy.Order))
			if i != len(qb.orderBy)-1 {
				builder.WriteString(", ")
			}
		}
	}

	builder.WriteString(";")

	return builder.String(), qb.namedArgs
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

func From(from string) QueryBuilder {
	return newQueryBuilder().From(from)
}
