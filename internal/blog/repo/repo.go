package repo

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/blog/data"
)

type Repository struct {
	db     *pgxpool.Pool
	models data.Models
	Posts  PostReaderWriter
}

func NewRepository(db *pgxpool.Pool, timeout *time.Duration) Repository {
	return Repository{
		db:     db,
		models: data.NewModels(db, timeout),
		Posts:  newPostRepository(db, timeout),
	}
}
