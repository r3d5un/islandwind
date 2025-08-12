package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/logging"
)

type Models struct {
	db    *pgxpool.Pool
	Posts PostModel
}

func NewModels(pool *pgxpool.Pool, timeout *time.Duration) Models {
	return Models{
		db:    pool,
		Posts: PostModel{DB: pool, Timeout: timeout},
	}
}

func (m *Models) BeginTx(ctx context.Context) (pgx.Tx, func(), error) {
	tx, err := m.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}

	rollbackFunc := func() {
		logger := logging.LoggerFromContext(ctx)
		logger.LogAttrs(ctx, slog.LevelInfo, "performing rollback")

		if err := tx.Rollback(ctx); err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				"error upon rollback",
				slog.String("error", err.Error()),
			)
			return
		}
	}

	return tx, rollbackFunc, nil
}
