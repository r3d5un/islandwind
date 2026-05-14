package builder_test

import (
	"testing"

	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseTag(t *testing.T) {
	type TestStruct struct {
		Column1     int64  `db:"column1"`
		Column2     string `db:"column2"`
		HiddenField string
	}

	columns := builder.ColumnsFrom(TestStruct{})
	assert.Contains(t, columns, "column1")
	assert.Contains(t, columns, "column2")
	assert.Len(t, columns, 2)
}
