package data

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/logging"
)

type RefreshToken struct {
	ID         uuid.UUID `json:"id"`
	Issuer     string    `json:"issuer"`
	Audience   string    `json:"audience"`
	Expiration time.Time `json:"expiration"`
	IssuedAt   time.Time `json:"issuedAt"`
	NotBefore  time.Time `json:"notBefore"`
}

type RefreshTokenInput struct {
	Issuer     string    `json:"issuer"`
	Audience   string    `json:"audience"`
	Expiration time.Time `json:"expiration"`
	IssuedAt   time.Time `json:"issuedAt"`
	NotBefore  time.Time `json:"notBefore"`
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
                                audience,
                                expiration,
                                issued_at,
                                not_before)
VALUES ($1::VARCHAR(512),
        $2::VARCHAR(512),
        $3::TIMESTAMP,
        $4::TIMESTAMP,
        $5::TIMESTAMP)
RETURNING
    id,
    issuer,
    audience,
    expiration,
    issued_at,
    not_before;
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
	var r RefreshToken
	err := q.QueryRow(
		ctx,
		stmt,
		input.Issuer,
		input.Audience,
		input.Expiration,
		input.IssuedAt,
		input.NotBefore,
	).Scan(
		&r.ID,
		&r.Issuer,
		&r.Audience,
		&r.Expiration,
		&r.IssuedAt,
		&r.NotBefore,
	)
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
	// TODO: Implement
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
       audience,
       expiration,
       issued_at,
       not_before
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
	var r RefreshToken
	err := q.QueryRow(
		ctx,
		stmt,
		id,
	).Scan(
		&r.ID,
		&r.Issuer,
		&r.Audience,
		&r.Expiration,
		&r.IssuedAt,
		&r.NotBefore,
	)
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
	return nil, nil, nil
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

func (m *RefreshTokenModel) delete(
	ctx context.Context,
	q db.Queryable,
	id uuid.UUID,
) (*RefreshToken, error) {
	return nil, nil
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
