package builder

import (
	"database/sql"
	"errors"
	"maps"
	"slices"
	"sort"
	"strconv"
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
	joins            []string
	namedArgs        pgx.NamedArgs
	returningColumns []string
	insertColumns    []string
	table            string
	limitSet         bool
	limit            int
	tuples           []Tuple
}

func (qb QueryBuilder) clone() QueryBuilder {
	return QueryBuilder{
		orderBy:          slices.Clone(qb.orderBy),
		whereClauses:     slices.Clone(qb.whereClauses),
		joins:            slices.Clone(qb.joins),
		namedArgs:        maps.Clone(qb.namedArgs),
		returningColumns: slices.Clone(qb.returningColumns),
		insertColumns:    slices.Clone(qb.insertColumns),
		table:            qb.table,
		tuples:           slices.Clone(qb.tuples),
	}
}

func newQueryBuilder() QueryBuilder {
	return QueryBuilder{
		orderBy:      make([]OrderBy, 0),
		whereClauses: make([]string, 0),
		joins:        make([]string, 0),
		namedArgs:    make(pgx.NamedArgs, 0),
		table:        "",
		limitSet:     false,
		limit:        0,
		tuples:       make([]Tuple, 0),
	}
}

func (qb QueryBuilder) OrderBy(order ...OrderBy) QueryBuilder {
	clone := qb.clone()
	clone.orderBy = append(clone.orderBy, order...)
	return clone
}

func (qb QueryBuilder) From(from string) QueryBuilder {
	clone := qb.clone()
	clone.table = from
	return clone
}

func (qb QueryBuilder) Join(join string) QueryBuilder {
	clone := qb.clone()
	clone.joins = append(clone.joins, join)
	return clone
}

func (qb QueryBuilder) Returning(cols ...string) QueryBuilder {
	clone := qb.clone()
	clone.returningColumns = append(clone.returningColumns, cols...)
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
	builder.WriteString(qb.table)

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

func (qb QueryBuilder) Delete() (string, pgx.NamedArgs) {
	var builder strings.Builder

	builder.WriteString("DELETE FROM ")
	builder.WriteString(qb.table)

	if len(qb.whereClauses) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(qb.whereClauses, " AND "))
	}

	colsLength := len(qb.returningColumns)
	if colsLength > 0 {
		builder.WriteString(" RETURNING ")
		for i, column := range qb.returningColumns {
			builder.WriteString(column)
			if i != colsLength-1 {
				builder.WriteString(", ")
			}
		}
	}

	builder.WriteString(";")

	return builder.String(), qb.namedArgs
}

var (
	ErrNoInsertTuples      = errors.New("no tuples provided for insertion")
	ErrTupleColumnNotFound = errors.New("tuple column not found")
)

type Tuple map[string]sql.Null[any]

func Insert(records ...Tuple) QueryBuilder {
	qb := newQueryBuilder()
	qb.tuples = append(qb.tuples, records...)
	if len(records) > 0 {
		for columnName := range records[0] {
			qb.insertColumns = append(qb.insertColumns, columnName)
		}
		sort.Strings(qb.insertColumns)
	}

	return qb
}

func (qb QueryBuilder) Into(into string) (string, pgx.NamedArgs, error) {
	var builder strings.Builder

	if len(qb.tuples) < 1 {
		return "", make(pgx.NamedArgs), ErrNoInsertTuples
	}

	builder.WriteString("INSERT INTO ")
	builder.WriteString(into)
	builder.WriteString(" (")
	builder.WriteString(strings.Join(qb.insertColumns, ", "))
	builder.WriteString(") ")
	builder.WriteString("VALUES ")

	for i, tuple := range qb.tuples {
		tupleParams := make([]string, 0)
		for _, columnName := range qb.insertColumns {
			val, ok := tuple[columnName]
			if !ok {
				return "", qb.namedArgs, ErrTupleColumnNotFound
			}
			param := columnName + "_" + strconv.Itoa(i)
			qb.namedArgs[param] = val
			tupleParams = append(tupleParams, "@"+param)
		}
		builder.WriteString("(")
		builder.WriteString(strings.Join(tupleParams, ", "))
		builder.WriteString(")")
		if i != len(qb.tuples)-1 {
			builder.WriteString(", ")
		}
	}

	colsLength := len(qb.returningColumns)
	if colsLength > 0 {
		builder.WriteString(" RETURNING ")
		for i, column := range qb.returningColumns {
			builder.WriteString(column)
			if i != colsLength-1 {
				builder.WriteString(", ")
			}
		}
	}

	builder.WriteString(";")

	return builder.String(), qb.namedArgs, nil
}

func Update(table string) QueryBuilder {
	qb := newQueryBuilder()
	qb.table = table
	return qb
}

func (qb QueryBuilder) Set(assignments ...Assignment) (string, pgx.NamedArgs, error) {
	var builder strings.Builder

	if qb.table == "" {
		return "", make(pgx.NamedArgs), errors.New("table not set")
	}
	builder.WriteString("UPDATE ")
	builder.WriteString(qb.table)
	builder.WriteString(" SET ")
	if len(assignments) < 1 {
		return "", make(pgx.NamedArgs), errors.New("no assignments set")
	}
	assignmentStrings := make([]string, len(assignments))
	for i, assignment := range assignments {
		assignmentStrings[i] = assignment.Text
		maps.Copy(qb.namedArgs, assignment.Args)
	}
	builder.WriteString(strings.Join(assignmentStrings, ", "))

	if len(qb.whereClauses) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(qb.whereClauses, " AND "))
	}

	colsLength := len(qb.returningColumns)
	if colsLength > 0 {
		builder.WriteString(" RETURNING ")
		for i, column := range qb.returningColumns {
			builder.WriteString(column)
			if i != colsLength-1 {
				builder.WriteString(", ")
			}
		}
	}

	builder.WriteString(";")

	return builder.String(), qb.namedArgs, nil
}
