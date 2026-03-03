package blog

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/blog/repo"
	"github.com/r3d5un/islandwind/internal/cache"
	"github.com/r3d5un/islandwind/internal/config"
	"github.com/r3d5un/islandwind/internal/logging"
)

const moduleName string = "blog"

type Module struct {
	name   string
	logger *slog.Logger
	db     *pgxpool.Pool
	cache  cache.Cache
	cfg    *config.Config
	repo   repo.Repository
	mux    *http.ServeMux
	auth   AuthMiddlewareService
}

type AuthMiddlewareService interface {
	AccessTokenMiddleware(next http.Handler) http.Handler
}

func NewModule(
	ctx context.Context,
	cfg *config.Config,
	db *pgxpool.Pool,
	cache cache.Cache,
	authModule AuthMiddlewareService,
) (*Module, error) {
	ctx, logger := logging.ContextLogger(ctx, slog.Group("module", slog.String("name", moduleName)))

	timeout := time.Duration(cfg.DB.TimeoutSeconds) * time.Second
	module := Module{
		name:   moduleName,
		logger: logger,
		db:     db,
		cache:  cache,
		cfg:    cfg,
		repo:   repo.NewRepository(db, cache, &timeout),
		auth:   authModule,
	}

	return &module, nil
}

func (m *Module) Start(ctx context.Context, mux *http.ServeMux) {
	m.mux = mux
	m.addRoutes(ctx)
}

func (m *Module) Shutdown() {
	m.logger.LogAttrs(context.Background(), slog.LevelInfo, "shutting down module")
}
