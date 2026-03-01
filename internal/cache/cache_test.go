package cache_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/r3d5un/islandwind/internal/cache"
	database "github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/testsuite"
)

var postgresCache cache.PostgresCache
var logger slog.Logger

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger = testsuite.NewTestLogger()

	logger.Info("creating PostgreSQL container")
	dbContainer, shutdown, err := database.NewPostgresTestcontainer(ctx)
	if err != nil {
		logger.Error("unable to start container", slog.String("error", err.Error()))
		return
	}
	defer shutdown()

	db, _, err := database.NewTestPool(ctx, dbContainer)
	if err != nil {
		logger.Error("unable to create database connection pool", slog.String("error", err.Error()))
		return
	}
	postgresCache = cache.NewPostgresCache(db, &logger)

	exitCode := m.Run()

	defer os.Exit(exitCode)
}
