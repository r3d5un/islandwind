package cache

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// DatabaseCache expands on the Cache interface with a transaction aware DeleteTx method.
type DatabaseCache interface {
	Cache
	// DeleteTx is used to invalidate cache entries as part of a transaction.
	DeleteTx(pgx.Tx, uuid.UUID) error
}

// Cache is an interface for setting, reading and invalidating cache entries.
type Cache interface {
	Set(uuid.UUID, any)
	Get(uuid.UUID, any) error
	Delete(uuid.UUID) error
}

var (
	ErrCacheMiss = errors.New("cache miss")
)
