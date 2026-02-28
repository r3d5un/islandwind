package handlers_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/r3d5un/islandwind/internal/blog/repo"
	database "github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/testsuite"
)

const (
	dbName     string = "postgres"
	dbUser     string = "postgres"
	dbPassword string = "postgres"
)

var blogReaderWriter repo.PostReaderWriter

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logger := testsuite.NewTestLogger()

	logger.Info("creating PostgreSQL container")
	dbContainer, shutdown, err := database.NewPostgresTestcontainer(ctx)
	if err != nil {
		logger.Error("unable to start container", slog.String("error", err.Error()))
		return
	}
	defer shutdown()

	db, cfg, err := database.NewTestPool(ctx, dbContainer)
	if err != nil {
		logger.Error("unable to create database connection pool", slog.String("error", err.Error()))
		return
	}
	timeout := cfg.TimeoutDuration()
	repository := repo.NewRepository(db, &timeout)
	blogReaderWriter = repository.Posts

	exitCode := m.Run()

	defer os.Exit(exitCode)
}
