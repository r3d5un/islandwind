package cache

import (
	"errors"

	"github.com/google/uuid"
)

type Cache interface {
	Set(uuid.UUID, any)
	Get(uuid.UUID, any) error
	Delete(uuid.UUID)
}

var (
	ErrCacheMiss = errors.New("cache miss")
)
