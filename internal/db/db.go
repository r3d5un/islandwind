package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	ConnStr         string `json:"-"`
	MaxOpenConns    int32  `json:"maxOpenConns"`
	IdleTimeMinutes int    `json:"idleTimeMinutes"`
	TimeoutSeconds  int    `json:"timeoutSeconds"`
}

func (c *Config) TimeoutDuration() time.Duration {
	return time.Duration(c.TimeoutSeconds) * time.Second
}

func (c *Config) IdleTime() time.Duration {
	return time.Duration(c.IdleTimeMinutes) * time.Minute
}

func OpenPool(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	pgxCfg, err := pgxpool.ParseConfig(config.ConnStr)
	if err != nil {
		return nil, err
	}
	pgxCfg.MaxConnIdleTime = config.IdleTime()
	pgxCfg.MaxConns = config.MaxOpenConns

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
