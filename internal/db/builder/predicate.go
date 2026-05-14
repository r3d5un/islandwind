package builder

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/oapi-codegen/nullable"
)

type PredicateCondition string

const (
	Equal          PredicateCondition = "="
	NotEqual       PredicateCondition = "!="
	Greater        PredicateCondition = ">"
	GreaterOrEqual PredicateCondition = ">="
	Less           PredicateCondition = "<"
	LessOrEqual    PredicateCondition = "<="

	// NOTE: LIKE have been deliberately omitted as these queries can carry heaby performance
	//  penalties, especially for prefixed wildcard quries. The user can opt-in to these queries
	//  using NewPredicate themselves as plain text.
)

type Predicate struct {
	Text string        `json:"text"`
	Arg  pgx.NamedArgs `json:"arg"`
}

func newPredicate() Predicate {
	return Predicate{
		Text: "",
		Arg:  make(pgx.NamedArgs),
	}
}

func NewNullablePredicate[T any](
	column string, cond PredicateCondition, value nullable.Nullable[T],
) Predicate {
	predicate := newPredicate()

	if !value.IsSpecified() {
		return predicate
	}

	predicate.Text = column + " " + string(cond) + " @" + column
	if value.IsNull() {
		predicate.Arg = pgx.NamedArgs{column: nil}
		return predicate
	}

	v, err := value.Get()
	if err != nil {
		return predicate
	}
	predicate.Arg = pgx.NamedArgs{column: v}

	return predicate
}

func NewNullPredicate[T any](
	column string, cond PredicateCondition, value sql.Null[T],
) Predicate {
	predicate := newPredicate()

	if !value.Valid {
		return predicate
	}

	predicate.Text = column + " " + string(cond) + " @" + column
	predicate.Arg = pgx.NamedArgs{column: value.V}

	return predicate
}

func NewGenericPredicate[T any](
	column string, cond PredicateCondition, value T,
) Predicate {
	return Predicate{
		Text: column + " " + string(cond) + " @" + column,
		Arg:  pgx.NamedArgs{column: value},
	}
}

func NewPredicate(text string, args pgx.NamedArgs) Predicate {
	return Predicate{Text: text, Arg: args}
}
