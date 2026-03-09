package repo

import (
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/cache"
)

type tokenStore struct {
	models *data.Models
	cache  cache.Cache
}

func newTokenStore(models *data.Models, cache cache.Cache) tokenStore {
	return tokenStore{models: models, cache: cache}
}
