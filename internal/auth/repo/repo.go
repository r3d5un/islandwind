package repo

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/auth/data"
)

type Repository struct {
	db     *pgxpool.Pool
	models data.Models
	Tokens TokenRepository
}

func NewRepository(
	db *pgxpool.Pool,
	timeout *time.Duration,
	secret []byte,
	issuer string,
) Repository {
	models := data.NewModels(db, timeout)
	return Repository{
		db:     db,
		models: models,
		Tokens: NewTokenRepository(secret, issuer, &models),
	}
}
