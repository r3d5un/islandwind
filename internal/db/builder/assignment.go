package builder

import (
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/oapi-codegen/nullable"
)

type Assignment struct {
	Text string        `json:"text"`
	Args pgx.NamedArgs `json:"args"`
}

func newAssignment() Assignment {
	return Assignment{Text: "", Args: make(pgx.NamedArgs)}
}

func NewNullableAssignment[T any](
	column string, value nullable.Nullable[T],
) Assignment {
	assignment := newAssignment()

	if !value.IsSpecified() {
		return assignment
	}

	assignment.Text = column + " = " + "@" + column
	if value.IsNull() {
		assignment.Args = pgx.NamedArgs{column: nil}
		return assignment
	}

	v, err := value.Get()
	if err != nil {
		return assignment
	}
	assignment.Args = pgx.NamedArgs{column: v}

	return assignment
}

func NewNullAssignment[T any](
	column string, value sql.Null[T],
) Assignment {
	assignment := newAssignment()

	if !value.Valid {
		return assignment
	}

	assignment.Text = column + " = " + "@" + column
	assignment.Args = pgx.NamedArgs{column: value.V}

	return assignment
}

func NewGenericAssignment[T any](column string, value T) Assignment {
	return Assignment{
		Text: column + " = " + "@" + column,
		Args: pgx.NamedArgs{column: value},
	}
}

func NewAssignment(text string, args pgx.NamedArgs) Assignment {
	return Assignment{Text: text, Args: args}
}
