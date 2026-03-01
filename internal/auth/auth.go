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
	"github.com/r3d5un/islandwind/internal/logging"
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

func NewModule(ctx context.Context, cfg *config.Config, db *pgxpool.Pool) (*Module, error) {
	ctx, logger := logging.ContextLogger(ctx, slog.Group("module", slog.String("name", moduleName)))

	logger.LogAttrs(ctx, slog.LevelInfo, "setting up module")
	timeout := time.Duration(cfg.DB.TimeoutSeconds) * time.Second
	module := Module{
		name:   moduleName,
		logger: logger,
		db:     db,
		cfg:    cfg,
		repo:   repo.NewRepository(db, &timeout, cfg.Auth),
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "module setup complete")

	return &module, nil
}

func (m *Module) Start(ctx context.Context, mux *http.ServeMux) {
	m.mux = mux
	m.addRoutes(ctx)
}

func (m *Module) Shutdown() {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "shutting down module")
}
