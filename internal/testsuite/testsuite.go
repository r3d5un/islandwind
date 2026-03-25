package testsuite

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
)

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

func NewTestLogger() slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return *logger
}
