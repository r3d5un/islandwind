package builder

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/oapi-codegen/nullable"
)

const (
	predicatePrefix = "predicate_"
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

	parameter := predicatePrefix + column

	predicate.Text = column + " " + string(cond) + " @" + parameter
	if value.IsNull() {
		predicate.Arg = pgx.NamedArgs{parameter: nil}
		return predicate
	}

	v, err := value.Get()
	if err != nil {
		return predicate
	}
	predicate.Arg = pgx.NamedArgs{parameter: v}

	return predicate
}

func NewNullPredicate[T any](
	column string, cond PredicateCondition, value sql.Null[T],
) Predicate {
	predicate := newPredicate()

	if !value.Valid {
		return predicate
	}

	parameter := predicatePrefix + column
	predicate.Text = column + " " + string(cond) + " @" + parameter
	predicate.Arg = pgx.NamedArgs{parameter: value.V}

	return predicate
}

func NewGenericPredicate[T any](
	column string, cond PredicateCondition, value T,
) Predicate {
	parameter := predicatePrefix + column
	return Predicate{
		Text: column + " " + string(cond) + " @" + parameter,
		Arg:  pgx.NamedArgs{parameter: value},
	}
}

func NewPredicate(text string, args pgx.NamedArgs) Predicate {
	return Predicate{Text: text, Arg: args}
}
