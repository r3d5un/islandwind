package auth

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/auth/repo"
	"github.com/r3d5un/islandwind/internal/config"
	"github.com/r3d5un/islandwind/internal/monolith/interfaces"
)

const moduleName string = "auth"

type Module struct {
	name       string
	logger     *slog.Logger
	db         *pgxpool.Pool
	cfg        *config.Config
	mux        *http.ServeMux
	instanceID uuid.UUID
	repo       repo.Repository
}

func (m *Module) Setup(ctx context.Context, mono interfaces.Monolith) {
	logger := mono.Logger().With(slog.Group(
		"module",
		slog.String("name", moduleName),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "setting up module")
	m.instanceID = mono.InstanceID()
	m.name = moduleName
	m.logger = logger
	m.db = mono.DB()
	m.cfg = mono.Config()
	timeout := time.Duration(m.cfg.DB.TimeoutSeconds) * time.Second
	m.repo = repo.NewRepository(m.db, &timeout, m.cfg.Auth)
	m.mux = mono.Mux()
	m.addRoutes(ctx)
	logger.LogAttrs(ctx, slog.LevelInfo, "module setup complete")
}

func (m *Module) Shutdown() {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "shutting down module")
}
