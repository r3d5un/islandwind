package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/r3d5un/islandwind/internal/monolith"
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
	mono, err := monolith.NewMonolith()
	if err != nil {
		return err
	}

	err = mono.Serve()
	if err != nil {
		return err
	}

	return nil
}
