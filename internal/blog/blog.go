package blog

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/monolith/interfaces"
)

const moduleName string = "blog"

type Module struct {
	name   string
	logger *slog.Logger
	db     *pgxpool.Pool
}

func (m *Module) Setup(ctx context.Context, mono interfaces.Monolith) {
	logger := mono.Logger().With(slog.Group(
		"module",
		slog.String("name", moduleName),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "setting up module")
	m.name = moduleName
	m.logger = logger
	m.db = mono.DB()
	logger.LogAttrs(ctx, slog.LevelInfo, "module setup complete")
}

func (m *Module) Shutdown() {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "shutting down module")
}
