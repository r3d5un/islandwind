package repo

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/auth/data"
)

type Repository struct {
	db     *pgxpool.Pool
	models data.Models
}

func NewRepository(db *pgxpool.Pool, timeout *time.Duration) Repository {
	return Repository{
		db:     db,
		models: data.NewModels(db, timeout),
	}
}
