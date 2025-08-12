package db_test

import (
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/r3d5un/islandwind/internal/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseMigration(t *testing.T) {
	projectRoot, err := testsuite.FindProjectRoot()
	if err != nil {
		t.Logf("unable to find project root: %s\n", err.Error())
		return
	}
	migrationURL := fmt.Sprintf("file://%s/migrations", projectRoot)

	m, err := migrate.New(migrationURL, connectionString)
	assert.NoError(t, err)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())

		err = m.Drop()
		assert.NoError(t, err)
	})

	t.Run("UpMigrations", func(t *testing.T) {
		err = m.Up()
		assert.NoError(t, err)
	})

	t.Run("DownMigrations", func(t *testing.T) {
		err = m.Down()
		assert.NoError(t, err)
	})
}
