package data

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/logging"
)

type RefreshToken struct {
	ID            uuid.UUID     `json:"id"`
	Issuer        string        `json:"issuer"`
	Expiration    time.Time     `json:"expiration"`
	IssuedAt      time.Time     `json:"issuedAt"`
	Invalidated   bool          `json:"invalidated"`
	InvalidatedBy uuid.NullUUID `json:"invalidatedBy"`
}

type RefreshTokenInput struct {
	Issuer     string    `json:"issuer"`
	Expiration time.Time `json:"expiration"`
	IssuedAt   time.Time `json:"issuedAt"`
}

type RefreshTokenPatch struct {
	ID            uuid.UUID      `json:"id"`
	Issuer        sql.NullString `json:"issuer"`
	Invalidated   sql.NullBool   `json:"invalidated"`
	InvalidatedBy uuid.NullUUID  `json:"invalidatedBy"`
}

type RefreshTokenModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *RefreshTokenModel) insert(
	ctx context.Context,
	q db.Queryable,
	input RefreshTokenInput,
) (*RefreshToken, error) {
	const stmt string = `
INSERT INTO auth.refresh_token (issuer,
                                expiration,
                                issued_at)
VALUES ($1::VARCHAR(512),
        $2::TIMESTAMP,
        $3::TIMESTAMP)
RETURNING
    id,
    issuer,
    expiration,
    issued_at,
    invalidated,
    invalidated_by;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(stmt)),
		slog.Any("input", input),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	r, err := m.scan(q.QueryRow(ctx, stmt, input.Issuer, input.Expiration, input.IssuedAt))
	if err != nil {
		return nil, db.HandleError(ctx, err)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token inserted", slog.Any("refreshToken", r))

	return &r, nil
}

func (m *RefreshTokenModel) Insert(
	ctx context.Context,
	input RefreshTokenInput,
) (*RefreshToken, error) {
	return m.insert(ctx, m.DB, input)
}

func (m *RefreshTokenModel) InsertTx(
	ctx context.Context,
	tx pgx.Tx,
	input RefreshTokenInput,
) (*RefreshToken, error) {
	return m.insert(ctx, tx, input)
}

func (m *RefreshTokenModel) selectOne(
	ctx context.Context,
	q db.Queryable,
	id uuid.UUID,
) (*RefreshToken, error) {
	const stmt string = `
SELECT id,
       issuer,
       expiration,
       issued_at,
       invalidated,
       invalidated_by
FROM auth.refresh_token
WHERE id = $1::UUID;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(stmt)),
		slog.String("id", id.String()),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	r, err := m.scan(q.QueryRow(ctx, stmt, id))
	if err != nil {
		return nil, db.HandleError(ctx, err)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token selected", slog.Any("refreshToken", r))

	return &r, nil
}

func (m *RefreshTokenModel) SelectOne(
	ctx context.Context,
	id uuid.UUID,
) (*RefreshToken, error) {
	return m.selectOne(ctx, m.DB, id)
}

func (m *RefreshTokenModel) SelectOneTx(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
) (*RefreshToken, error) {
	return m.selectOne(ctx, tx, id)
}

func (m *RefreshTokenModel) selectMany(
	ctx context.Context,
	q db.Queryable,
	filter Filter,
) ([]*RefreshToken, *Metadata, error) {
	const stmt string = `
SELECT id,
       issuer,
       expiration,
       issued_at,
       invalidated,
       invalidated_by
FROM auth.refresh_token
WHERE ($2::UUID IS NULL OR id = $2::UUID)
  AND ($3::VARCHAR(512) IS NULL OR issuer = $3::VARCHAR(512))
  AND ($4::TIMESTAMP IS NULL OR expiration <= $4::TIMESTAMP)
  AND ($5::TIMESTAMP IS NULL OR expiration > $5::TIMESTAMP)
  AND ($6::TIMESTAMP IS NULL OR issued_at <= $6::TIMESTAMP)
  AND ($7::TIMESTAMP IS NULL OR issued_at > $7::TIMESTAMP)
  AND id > $8::UUID
ORDER BY expiration, id
LIMIT $1;
`

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("statement", logging.MinifySQL(stmt)),
		slog.Any("filter", filter),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	rows, err := q.Query(
		ctx,
		stmt,
		filter.PageSize,
		filter.ID,
		filter.Issuer,
		filter.ExpirationFrom,
		filter.ExpirationTo,
		filter.IssuedAtFrom,
		filter.IssuedAtTo,
		filter.LastSeen,
	)
	if err != nil {
		return nil, nil, db.HandleError(ctx, err)
	}

	var tokens []*RefreshToken

	for rows.Next() {
		var r RefreshToken
		r, err := m.scan(rows)
		if err != nil {
			return nil, nil, db.HandleError(ctx, err)
		}
		tokens = append(tokens, &r)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, db.HandleError(ctx, err)
	}
	metadata := NewMetadata(tokens, filter)

	logger.LogAttrs(ctx, slog.LevelInfo, "posts selected", slog.Any("metadata", metadata))

	return tokens, &metadata, nil
}

func (m *RefreshTokenModel) SelectMany(
	ctx context.Context,
	filter Filter,
) ([]*RefreshToken, *Metadata, error) {
	return m.selectMany(ctx, m.DB, filter)
}

func (m *RefreshTokenModel) SelectManyTx(
	ctx context.Context,
	tx pgx.Tx,
	filter Filter,
) ([]*RefreshToken, *Metadata, error) {
	return m.selectMany(ctx, tx, filter)
}

func (m *RefreshTokenModel) update(
	ctx context.Context,
	q db.Queryable,
	input RefreshTokenPatch,
) (*RefreshToken, error) {
	const stmt string = `
UPDATE auth.refresh_token
SET issuer = COALESCE($2::VARCHAR(512), issuer),
    invalidated = COALESCE($3::BOOLEAN, invalidated),
	invalidated_by = COALESCE($4::UUID, invalidated_by)
WHERE id = $1::UUID
RETURNING
    id,
    issuer,
    expiration,
    issued_at,
    invalidated,
    invalidated_by;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(stmt)),
		slog.Any("input", input),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	r, err := m.scan(q.QueryRow(
		ctx,
		stmt,
		input.ID,
		input.Issuer,
		input.Invalidated,
		input.InvalidatedBy,
	))
	if err != nil {
		return nil, db.HandleError(ctx, err)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token updated", slog.Any("refreshToken", r))

	return &r, nil
}

func (m *RefreshTokenModel) Update(
	ctx context.Context,
	input RefreshTokenPatch,
) (*RefreshToken, error) {
	return m.update(ctx, m.DB, input)
}

func (m *RefreshTokenModel) UpdateTx(
	ctx context.Context,
	tx pgx.Tx,
	input RefreshTokenPatch,
) (*RefreshToken, error) {
	return m.update(ctx, tx, input)
}

func (m *RefreshTokenModel) delete(
	ctx context.Context,
	q db.Queryable,
	id uuid.UUID,
) (*RefreshToken, error) {
	const stmt string = `
DELETE
FROM auth.refresh_token
WHERE id = $1::UUID
RETURNING
    id,
    issuer,
    expiration,
    issued_at,
    invalidated,
    invalidated_by;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(stmt)),
		slog.String("id", id.String()),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	r, err := m.scan(q.QueryRow(ctx, stmt, id))
	if err != nil {
		return nil, db.HandleError(ctx, err)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "refresh token selected", slog.Any("refreshToken", r))

	return &r, nil
}

func (m *RefreshTokenModel) Delete(ctx context.Context, id uuid.UUID) (*RefreshToken, error) {
	return m.delete(ctx, m.DB, id)
}

func (m *RefreshTokenModel) DeleteTx(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
) (*RefreshToken, error) {
	return m.delete(ctx, tx, id)
}

func (m *RefreshTokenModel) deleteMany(
	ctx context.Context,
	q db.Queryable,
	filter Filter,
) (*int64, error) {
	const stmt string = `
DELETE
FROM auth.refresh_token
WHERE ($1::UUID IS NULL OR id = $1::UUID)
  AND ($2::VARCHAR(512) IS NULL OR issuer = $2::VARCHAR(512))
  AND ($3::TIMESTAMP IS NULL OR expiration <= $3::TIMESTAMP)
  AND ($4::TIMESTAMP IS NULL OR expiration > $4::TIMESTAMP)
  AND ($5::TIMESTAMP IS NULL OR issued_at <= $5::TIMESTAMP)
  AND ($6::TIMESTAMP IS NULL OR issued_at > $6::TIMESTAMP);
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(stmt)),
		slog.Any("filter", filter),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	err := db.DeleteManyGuardrail(
		filter.ID,
		filter.Issuer,
		filter.ExpirationFrom,
		filter.ExpirationTo,
		filter.IssuedAtFrom,
		filter.IssuedAtTo,
	)
	if err != nil {
		return nil, db.HandleError(ctx, err)
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	res, err := q.Exec(
		ctx,
		stmt,
		filter.ID,
		filter.Issuer,
		filter.ExpirationFrom,
		filter.ExpirationTo,
		filter.IssuedAtFrom,
		filter.IssuedAtTo,
	)
	if err != nil {
		return nil, db.HandleError(ctx, err)
	}
	rowsAffected := res.RowsAffected()
	logger.LogAttrs(
		ctx, slog.LevelInfo, "refresh tokens deleted", slog.Int64("rowsAffected", rowsAffected),
	)

	return &rowsAffected, nil
}

func (m *RefreshTokenModel) DeleteMany(ctx context.Context, filter Filter) (*int64, error) {
	return m.deleteMany(ctx, m.DB, filter)
}

func (m *RefreshTokenModel) DeleteManyTx(
	ctx context.Context,
	tx pgx.Tx,
	filter Filter,
) (*int64, error) {
	return m.deleteMany(ctx, tx, filter)
}

func (m *RefreshTokenModel) scan(row pgx.Row) (RefreshToken, error) {
	var r RefreshToken
	err := row.Scan(
		&r.ID,
		&r.Issuer,
		&r.Expiration,
		&r.IssuedAt,
		&r.Invalidated,
		&r.InvalidatedBy,
	)
	if err != nil {
		return r, err
	}
	return r, nil
}
