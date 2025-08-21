package repo

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/auth/config"
	"github.com/r3d5un/islandwind/internal/auth/data"
)

type Repository struct {
	db     *pgxpool.Pool
	models data.Models
	cfg    config.Config
	Tokens TokenService
}

func NewRepository(
	db *pgxpool.Pool,
	timeout *time.Duration,
	cfg config.Config,
) Repository {
	models := data.NewModels(db, timeout)
	return Repository{
		db:     db,
		models: models,
		cfg:    cfg,
		Tokens: NewTokenRepository(
			[]byte(cfg.SigningSecret),
			[]byte(cfg.RefreshSigningSecret),
			cfg.TokenIssuer,
			&models,
		),
	}
}
