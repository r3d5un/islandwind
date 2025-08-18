package data_test

import (
	"context"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/auth/data"
	database "github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/testsuite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     string = "postgres"
	dbUser     string = "postgres"
	dbPassword string = "postgres"
)

var models data.Models
var db *pgxpool.Pool

func TestMain(m *testing.M) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logger.Info("creating postgres container")
	dbContainer, err := postgres.Run(
		ctx,
		"postgres:17.4",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithLogger(log.Default()),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(dbContainer); err != nil {
			logger.Info("failed to terminate container", slog.String("error", err.Error()))
		}
	}()
	if err != nil {
		logger.Error("unable to start container", slog.String("error", err.Error()))
		return
	}

	connStr, err := dbContainer.ConnectionString(ctx, "sslmode=disable", "application_name=rosetta")
	if err != nil {
		logger.Error("unable to get database connection string", slog.String("error", err.Error()))
		return
	}

	migrate, err := testsuite.NewMigrateClient(connStr)
	if err != nil {
		logger.Error("unable to create migration client", slog.String("error", err.Error()))
		return
	}

	err = migrate.Up()
	if err != nil {
		logger.Error("up migrations failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		err = migrate.Down()
		if err != nil {
			logger.Error("up migrations failed", slog.String("error", err.Error()))
			return
		}
	}()

	dbConfig := database.Config{
		ConnStr:         connStr,
		MaxOpenConns:    20,
		IdleTimeMinutes: 1,
		TimeoutSeconds:  30,
	}
	db, err = database.OpenPool(ctx, dbConfig)
	if err != nil {
		logger.Error("unable to create database connection pool", slog.String("error", err.Error()))
		return
	}
	timeout := dbConfig.TimeoutDuration()
	models = data.NewModels(db, &timeout)

	exitCode := m.Run()

	err = migrate.Drop()
	if err != nil {
		logger.Error("unable to drop database schemas", slog.String("error", err.Error()))
	}

	defer os.Exit(exitCode)
}
