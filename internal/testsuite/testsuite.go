package testsuite

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
)

func Assert(condition bool, msg string, object any) {
	if !condition {
		slog.LogAttrs(context.Background(), slog.LevelError, msg, slog.Any("object", object))
		panic(msg)
	}
}

func FindProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get current working directory: %w", err)
	}

	markerFile := ".git"

	for {
		if _, err := os.Stat(filepath.Join(cwd, markerFile)); err == nil {
			return cwd, nil
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "", fmt.Errorf("project root not found")
		}
		cwd = parent
	}
}

func NewMigrateClient(connStr string) (*migrate.Migrate, error) {
	projectRoot, err := FindProjectRoot()
	if err != nil {
		return nil, err
	}
	migrationURL := fmt.Sprintf("file://%s/migrations", projectRoot)

	m, err := migrate.New(migrationURL, connStr)
	if err != nil {
		return nil, err
	}

	return m, nil
}
