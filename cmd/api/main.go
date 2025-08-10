package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/config"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		slog.LogAttrs(
			context.Background(),
			slog.LevelError,
			"unrecoverable error occurred",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	defer slog.Info("shutting down")
	ctx := context.Background()
	instanceID := uuid.New()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	logGroup := slog.Group(
		"instance",
		slog.String("name", viper.GetString("app.name")),
		slog.String("environment", viper.GetString("app.environment")),
		slog.String("id", instanceID.String()),
	)
	var handler slog.Handler

	switch cfg.App.Environment {
	case "testing":
		fallthrough
	case "production":
		handler = slog.NewJSONHandler(os.Stderr, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}
	logger := slog.New(handler).With(logGroup)
	slog.SetDefault(logger)

	logger.LogAttrs(ctx, slog.LevelInfo, "starting up")

	return nil
}
