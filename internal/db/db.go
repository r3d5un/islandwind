package db

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	ConnStr         string `json:"-"`
	MaxOpenConns    int32  `json:"maxOpenConns"`
	IdleTimeMinutes int    `json:"idleTimeMinutes"`
	TimeoutSeconds  int    `json:"timeoutSeconds"`
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("maxOpenConns", int(c.MaxOpenConns)),
		slog.Int("idleTimeMinutes", int(c.IdleTimeMinutes)),
		slog.Int("timeoutSeconds", int(c.TimeoutSeconds)),
	)
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

type Queryable interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
}

var (
	ErrUnsafeDeleteFilter = errors.New("filter is unsafe")
)

// deleteManyGuardrail accepts a slice of values and raises and error if all
// given values are nil. If all values are nil an ErrUnsafeDeleteFilter error
// is returned. If a value is not nil the function returns nil, and the
// filter is safe to use.
//
// WARNING: This function assumed any given value in the input slice is a
// pointer that can be checked for nil. A non-pointer value will return
// early, but the the filter may still be unsafe.
func deleteManyGuardrail(input ...any) error {
	for _, x := range input {
		if v := reflect.ValueOf(x); !v.IsNil() {
			return nil
		}
	}

	return ErrUnsafeDeleteFilter
}

// isEmpty checks if a slice is nil or empty
func isEmpty[T comparable](x []*T) bool {
	if len(x) < 1 {
		return true
	}

	return false
}
