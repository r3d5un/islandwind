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

// Post is the database record for a blog post.
type Post struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Published bool         `json:"published"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Deleted   bool         `json:"deleted"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}

// PostInput is the input type used by the BlogModel for creating new blog post records.
type PostInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

// PostPatch is used for updating any existing blog post records. All fields except
// the ID is optional, but if populated will update the record when given to
// BlogModel.Update.
type PostPatch struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title"`
	Content   *string   `json:"content"`
	Published *bool     `json:"published"`
	Deleted   *bool     `json:"deleted"`
}

type PostModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *PostModel) insert(ctx context.Context, q db.Queryable, input PostInput) (*Post, error) {
	const stmt string = `
INSERT INTO blog.post (title,
                       content,
                       published)
VALUES ($1::VARCHAR(1024),
        $2::TEXT,
        $3::BOOLEAN)
RETURNING
    id,
    title,
    content,
    published,
    created_at,
    updated_at,
    deleted,
    deleted_at;
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
	var p Post
	err := q.QueryRow(
		ctx,
		stmt,
		input.Title,
		input.Content,
		input.Published,
	).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Published,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Deleted,
		&p.DeletedAt,
	)
	if err != nil {
		return nil, db.HandleError(err, logger)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "post inserted", slog.Any("post", p))

	return &p, nil
}

func (m *PostModel) Insert(ctx context.Context, input PostInput) (*Post, error) {
	return m.insert(ctx, m.DB, input)
}

func (m *PostModel) InsertTx(ctx context.Context, tx pgx.Tx, input PostInput) (*Post, error) {
	return m.insert(ctx, tx, input)
}

func (m *PostModel) selectOne(ctx context.Context, q db.Queryable, id uuid.UUID) (*Post, error) {
	const stmt string = `
SELECT id,
       title,
       content,
       published,
       created_at,
       updated_at,
       deleted,
       deleted_at
FROM blog.post
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
	var p Post
	err := q.QueryRow(
		ctx,
		stmt,
		id,
	).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Published,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Deleted,
		&p.DeletedAt,
	)
	if err != nil {
		return nil, db.HandleError(err, logger)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "post selected", slog.Any("post", p))

	return &p, nil
}

func (m *PostModel) SelectOne(ctx context.Context, id uuid.UUID) (*Post, error) {
	return m.selectOne(ctx, m.DB, id)
}

func (m *PostModel) SelectOneTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Post, error) {
	return m.selectOne(ctx, tx, id)
}

func (m *PostModel) selectMany(
	ctx context.Context,
	q db.Queryable,
	filter Filter,
) ([]*Post, *Metadata, error) {
	const stmt string = `
SELECT id,
       title,
       content,
       published,
       created_at,
       updated_at,
       deleted,
       deleted_at
FROM blog.post
WHERE ($2::UUID IS NULL OR id = $2::UUID)
  AND ($3::VARCHAR(1024) IS NULL OR title = $3::VARCHAR(1024))
  AND ($4::BOOLEAN IS NULL OR published = $4::BOOLEAN)
  AND ($5::TIMESTAMPTZ IS NULL OR created_at <= $5::TIMESTAMPTZ)
  AND ($6::TIMESTAMPTZ IS NULL OR created_at > $6::TIMESTAMPTZ)
  AND ($7::TIMESTAMPTZ IS NULL OR updated_at <= $7::TIMESTAMPTZ)
  AND ($8::TIMESTAMPTZ IS NULL OR updated_at > $8::TIMESTAMPTZ)
  AND ($9::BOOLEAN IS NULL OR deleted = $9::BOOLEAN)
  AND ($10::TIMESTAMPTZ IS NULL OR updated_at <= $10::TIMESTAMPTZ)
  AND ($11::TIMESTAMPTZ IS NULL OR updated_at > $11::TIMESTAMPTZ)
  AND id > $12::UUID
ORDER BY created_at, id
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
		filter.Title,
		filter.Published,
		filter.CreatedAtFrom,
		filter.CreatedAtTo,
		filter.UpdatedAtFrom,
		filter.UpdatedAtTo,
		filter.Deleted,
		filter.DeletedAtFrom,
		filter.DeletedAtTo,
		filter.LastSeen,
	)
	if err != nil {
		logger.LogAttrs(
			ctx, slog.LevelError, "unable to perform query", slog.String("error", err.Error()),
		)
		return nil, nil, err
	}

	posts := []*Post{}

	for rows.Next() {
		var b Post

		err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Content,
			&b.Published,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.Deleted,
			&b.DeletedAt,
		)
		if err != nil {
			return nil, nil, db.HandleError(err, logger)
		}
		posts = append(posts, &b)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, db.HandleError(err, logger)
	}
	metadata := NewMetadata(posts, filter)

	logger.LogAttrs(ctx, slog.LevelInfo, "posts selected", slog.Any("metadata", metadata))

	return posts, &metadata, nil
}

func (m *PostModel) SelectMany(
	ctx context.Context,
	filter Filter,
) ([]*Post, *Metadata, error) {
	return m.selectMany(ctx, m.DB, filter)
}

func (m *PostModel) SelectManyTx(
	ctx context.Context,
	tx pgx.Tx,
	filter Filter,
) ([]*Post, *Metadata, error) {
	return m.selectMany(ctx, tx, filter)
}

func (m *PostModel) update(ctx context.Context, q db.Queryable, patch PostPatch) (*Post, error) {
	const stmt string = `
UPDATE blog.post
SET title      = COALESCE($2::VARCHAR(1024), title),
    content    = COALESCE($3::TEXT, content),
    published  = COALESCE($4::BOOLEAN, published),
    deleted    = COALESCE($5::BOOLEAN, deleted),
    deleted_at = CASE
                     WHEN COALESCE($5::BOOLEAN, deleted) = TRUE AND deleted_at IS NULL THEN NOW()
                     WHEN COALESCE($5::BOOLEAN, deleted) = FALSE THEN NULL
                     ELSE deleted_at
        END,
    updated_at = NOW()
WHERE id = $1::UUID
RETURNING id,
    title,
    content,
    published,
    created_at,
    updated_at,
    deleted,
    deleted_at;
`

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(stmt)),
		slog.Any("id", patch),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	var p Post
	err := q.QueryRow(
		ctx,
		stmt,
		patch.ID,
		patch.Title,
		patch.Content,
		patch.Published,
		patch.Deleted,
	).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Published,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Deleted,
		&p.DeletedAt,
	)
	if err != nil {
		return nil, db.HandleError(err, logger)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "post updated", slog.Any("post", p))

	return &p, nil
}

func (m *PostModel) Update(ctx context.Context, patch PostPatch) (*Post, error) {
	return m.update(ctx, m.DB, patch)
}

func (m *PostModel) UpdateTx(ctx context.Context, tx pgx.Tx, patch PostPatch) (*Post, error) {
	return m.update(ctx, tx, patch)
}

func (m *PostModel) delete(ctx context.Context, q db.Queryable, id uuid.UUID) (*Post, error) {
	const stmt string = `
DELETE
FROM blog.post
WHERE id = $1::UUID
RETURNING id,
    title,
    content,
    published,
    created_at,
    updated_at,
    deleted,
    deleted_at;
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
	var p Post
	err := q.QueryRow(
		ctx,
		stmt,
		id,
	).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Published,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.Deleted,
		&p.DeletedAt,
	)
	if err != nil {
		return nil, db.HandleError(err, logger)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "post deleted", slog.Any("post", p))

	return &p, nil
}

func (m *PostModel) Delete(ctx context.Context, id uuid.UUID) (*Post, error) {
	return m.delete(ctx, m.DB, id)
}

func (m *PostModel) DeleteTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Post, error) {
	return m.delete(ctx, tx, id)
}
