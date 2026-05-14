package builder

import (
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/oapi-codegen/nullable"
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
	joins            []string
	namedArgs        pgx.NamedArgs
	returning        bool
	returningColumns []string
	from             string
	limitSet         bool
	limit            int
}

func (qb QueryBuilder) clone() QueryBuilder {
	return QueryBuilder{
		orderBy:          slices.Clone(qb.orderBy),
		whereClauses:     slices.Clone(qb.whereClauses),
		joins:            slices.Clone(qb.joins),
		namedArgs:        maps.Clone(qb.namedArgs),
		returning:        qb.returning,
		returningColumns: slices.Clone(qb.returningColumns),
		from:             qb.from,
	}
}

func newQueryBuilder() QueryBuilder {
	return QueryBuilder{
		orderBy:      make([]OrderBy, 0),
		whereClauses: make([]string, 0),
		joins:        make([]string, 0),
		namedArgs:    make(pgx.NamedArgs, 0),
		returning:    false,
		from:         "",
		limitSet:     false,
		limit:        0,
	}
}

func (qb QueryBuilder) OrderBy(order ...OrderBy) QueryBuilder {
	clone := qb.clone()
	clone.orderBy = append(clone.orderBy, order...)
	return clone
}

func (qb QueryBuilder) From(from string) QueryBuilder {
	clone := qb.clone()
	clone.from = from
	return clone
}

func (qb QueryBuilder) Join(join string) QueryBuilder {
	clone := qb.clone()
	clone.joins = append(clone.joins, join)
	return clone
}

func (qb QueryBuilder) Returning(returning bool) QueryBuilder {
	clone := qb.clone()
	clone.returning = returning
	return clone
}

func (qb QueryBuilder) Limit(limit int) QueryBuilder {
	clone := qb.clone()
	clone.limit = limit
	clone.limitSet = true
	return clone
}

func (qb QueryBuilder) Select(cols ...string) (string, pgx.NamedArgs) {
	clone := qb.clone()
	clone.returningColumns = append(clone.returningColumns, cols...)
	return clone.selectExp()
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

	if len(qb.joins) > 0 {
		for _, join := range qb.joins {
			builder.WriteString(" JOIN ")
			builder.WriteString(join)
		}
	}

	if len(qb.whereClauses) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(qb.whereClauses, " AND "))
	}

	if len(qb.orderBy) > 0 {
		builder.WriteString(" ORDER BY ")
		for i, orderBy := range qb.orderBy {
			builder.WriteString(orderBy.Column)
			builder.WriteString(" ")
			builder.WriteString(string(orderBy.Order))
			if i != len(qb.orderBy)-1 {
				builder.WriteString(", ")
			}
		}
	}

	if qb.limitSet {
		builder.WriteString("LIMIT ")
		builder.WriteString(strconv.Itoa(qb.limit))
	}

	builder.WriteString(";")

	return builder.String(), qb.namedArgs
}

func From(from string) QueryBuilder {
	return newQueryBuilder().From(from)
}

var (
	_ ExplicitNull = nullable.Nullable[string]{}
	_ ExplicitNull = (*nullable.Nullable[string])(nil)
)

// ExplicitNull represents a tri-state nullable value used in query filters.
//
// States:
// - Unspecified: filter should be ignored
// - Null: filter should target SQL NULL
// - Value: filter should apply with a concrete value
type ExplicitNull interface {
	IsSpecified() bool
	IsNull() bool
}

func (qb QueryBuilder) Where(predicates ...Predicate) QueryBuilder {
	clone := qb.clone()
	for i, p := range predicates {
		if strings.TrimSpace(p.Text) == "" {
			continue
		}

		text := p.Text
		args := make(pgx.NamedArgs, len(p.Arg))
		for k, v := range p.Arg {
			newK := k
			if _, exists := clone.namedArgs[k]; exists {
				newK = k + "_" + strconv.Itoa(i)
				text = strings.ReplaceAll(text, "@"+k, "@"+newK)
			}
			args[newK] = v
		}

		clone.whereClauses = append(clone.whereClauses, text)
		maps.Copy(clone.namedArgs, args)
	}
	return clone
}
