package builder_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
)

func TestAllColumnsFrom(t *testing.T) {
	type MyStruct struct {
		ID    int    `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
		Other string // No db tag
	}

	t.Run("ExtractsAllDbTags", func(t *testing.T) {
		cols := builder.ColumnsFrom(MyStruct{})
		expected := []string{"id", "name", "email"}
		assert.ElementsMatch(t, expected, cols)
	})

	t.Run("QueryBuilderIntegration", func(t *testing.T) {
		stmt, _ := builder.From("users").Select(builder.ColumnsFrom(MyStruct{})...)
		// ElementsMatch because order might not be guaranteed by reflection (though usually it matches struct order)
		// But SELECT statement has fixed order in string.
		// Actually reflection order is deterministic (same as struct definition).
		assert.Contains(t, stmt, "SELECT")
		assert.Contains(t, stmt, "id")
		assert.Contains(t, stmt, "name")
		assert.Contains(t, stmt, "email")
		assert.NotContains(t, stmt, "Other")
		assert.Contains(t, stmt, "FROM users;")
	})

	t.Run("PointerToStruct", func(t *testing.T) {
		cols := builder.ColumnsFrom(&MyStruct{})
		expected := []string{"id", "name", "email"}
		assert.ElementsMatch(t, expected, cols)
	})
}
