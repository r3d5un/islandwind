package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	database "github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/logging"
)

type PostgresCache struct {
	db         *pgxpool.Pool
	logger     *slog.Logger
	setChan    chan postgresSetCacheMessage
	deleteChan chan uuid.UUID
	done       chan struct{}
}

func NewPostgresCache(db *pgxpool.Pool, logger *slog.Logger) PostgresCache {
	return PostgresCache{
		db:         db,
		logger:     logger,
		setChan:    make(chan postgresSetCacheMessage, 64),
		deleteChan: make(chan uuid.UUID, 64),
		done:       make(chan struct{}),
	}
}

func (c *PostgresCache) Start() error {
	go func() {
		for {
			select {
			case <-c.done:
				return
			case msg, ok := <-c.setChan:
				if !ok {
					return
				}
				if err := c.set(msg); err != nil {
					c.logger.Error(
						"unable to cache data",
						slog.String("error", err.Error()),
						slog.Any("msg", msg),
					)
				}
			case ID, ok := <-c.deleteChan:
				if !ok {
					return
				}
				if err := c.delete(ID); err != nil {
					c.logger.Error(
						"unable to invalidate cache entry",
						slog.String("error", err.Error()),
						slog.Any("id", ID),
					)
				}
			}
		}
	}()

	go func() {
		for range time.Tick(30 * time.Second) {
			if err := c.deleteExpired(); err != nil {
				c.logger.Error(
					"unable to delete expired cache data",
					slog.String("error", err.Error()),
				)
			}
		}
	}()

	return nil
}

func (c *PostgresCache) Shutdown() {
	c.done <- struct{}{}
	return
}

type postgresSetCacheMessage struct {
	ID   uuid.UUID `json:"id"`
	Data any       `json:"data"`
}

func (c *PostgresCache) set(msg postgresSetCacheMessage) error {
	const stmt string = `
INSERT INTO cache.general (id, data)
VALUES ($1, $2)
ON CONFLICT (id)
    DO UPDATE SET data = $2;
`

	logger := c.logger.With(slog.Group(
		"query",
		slog.String("statement", logging.MinifySQL(stmt)),
		slog.Any("msg", msg),
	))

	marshalled, err := json.Marshal(msg.Data)
	if err != nil {
		logger.Error("unable to marshal data", slog.String("error", err.Error()))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("performing query")
	if _, err := c.db.Exec(ctx, stmt, msg.ID, marshalled); err != nil {
		logger.Error("unable to insert the cache data", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (c *PostgresCache) Set(ID uuid.UUID, data any) {
	c.setChan <- postgresSetCacheMessage{ID: ID, Data: data}
}

func (c *PostgresCache) Get(ID uuid.UUID, data any) error {
	const stmt string = `
SELECT g.data
FROM cache.general g
WHERE id = $1;
`

	logger := c.logger.With(slog.Group(
		"query",
		slog.String("statement", logging.MinifySQL(stmt)),
		slog.String("id", ID.String()),
	))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("performing query")
	var jsonb []byte
	row := c.db.QueryRow(ctx, stmt, ID)
	err := row.Scan(&jsonb)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrCacheMiss
		}
		return database.HandleError(ctx, err)
	}

	if err := json.Unmarshal(jsonb, data); err != nil {
		return err
	}

	return nil
}

func (c *PostgresCache) Delete(ID uuid.UUID) {
	c.deleteChan <- ID
}

func (c *PostgresCache) delete(ID uuid.UUID) error {
	const stmt string = `
DELETE
FROM cache.general
WHERE id = $1;
`

	logger := c.logger.With(slog.Group(
		"query",
		slog.String("statement", logging.MinifySQL(stmt)),
		slog.String("id", ID.String()),
	))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("performing query")
	if _, err := c.db.Exec(ctx, stmt, ID); err != nil {
		logger.Error("unable to insert the cache data", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (c *PostgresCache) deleteExpired() error {
	const stmt string = `
DELETE
FROM cache.general
WHERE expires_at < NOW();
`

	logger := c.logger.With(slog.Group("query", slog.String("statement", logging.MinifySQL(stmt))))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("performing query")
	if _, err := c.db.Exec(ctx, stmt); err != nil {
		logger.Error("unable to insert the cache data", slog.String("error", err.Error()))
		return err
	}

	return nil
}
