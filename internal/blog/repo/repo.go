package repo

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/cache"
)

type Repository struct {
	db     *pgxpool.Pool
	cache  cache.Cache
	models data.Models
	Posts  PostReaderWriter
}

func NewRepository(db *pgxpool.Pool, c cache.Cache, timeout *time.Duration) Repository {
	return Repository{
		db:     db,
		cache:  c,
		models: data.NewModels(db, timeout),
		Posts:  newPostRepository(db, c, timeout),
	}
}
