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

// Blog is the database record for a blog post.
type Blog struct {
	ID        uuid.UUID    `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Published bool         `json:"published"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Deleted   bool         `json:"deleted"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}

// BlogInput is the input type used by the BlogModel for creating new blog post records.
type BlogInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

// BlogPatch is used for updating any existing blog post records. All fields except
// the ID is optional, but if populated will update the record when given to
// BlogModel.Update.
type BlogPatch struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title"`
	Content   *string   `json:"content"`
	Published *bool     `json:"published"`
}

type BlogModel struct {
	DB      *pgxpool.Pool
	Timeout *time.Duration
}

func (m *BlogModel) insert(ctx context.Context, q db.Queryable, input BlogInput) (*Blog, error) {
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
	var b Blog
	err := q.QueryRow(
		ctx,
		stmt,
		input.Title,
		input.Content,
		input.Published,
	).Scan(
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
		return nil, db.HandleError(err, logger)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog inserted", slog.Any("blog", b))

	return &b, nil
}

func (m *BlogModel) Insert(ctx context.Context, input BlogInput) (*Blog, error) {
	return m.insert(ctx, m.DB, input)
}

func (m *BlogModel) InsertTx(ctx context.Context, tx pgx.Tx, input BlogInput) (*Blog, error) {
	return m.insert(ctx, tx, input)
}

func (m *BlogModel) selectOne(ctx context.Context, q db.Queryable, id uuid.UUID) (*Blog, error) {
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
	var b Blog
	err := q.QueryRow(
		ctx,
		stmt,
		id,
	).Scan(
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
		return nil, db.HandleError(err, logger)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog selected", slog.Any("blog", b))

	return &b, nil
}

func (m *BlogModel) SelectOne(ctx context.Context, id uuid.UUID) (*Blog, error) {
	return m.selectOne(ctx, m.DB, id)
}

func (m *BlogModel) SelectOneTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Blog, error) {
	return m.selectOne(ctx, tx, id)
}

func (m *BlogModel) selectMany(
	ctx context.Context,
	q db.Queryable,
	filter Filter,
) ([]*Blog, *Metadata, error) {
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

	posts := []*Blog{}

	for rows.Next() {
		var b Blog

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

func (m *BlogModel) SelectMany(
	ctx context.Context,
	filter Filter,
) ([]*Blog, *Metadata, error) {
	return m.selectMany(ctx, m.DB, filter)
}

func (m *BlogModel) SelectManyTx(
	ctx context.Context,
	tx pgx.Tx,
	filter Filter,
) ([]*Blog, *Metadata, error) {
	return m.selectMany(ctx, tx, filter)
}

func (m *BlogModel) update(ctx context.Context, q db.Queryable, patch BlogPatch) (*Blog, error) {
	// TODO: Implement
	return nil, nil
}

func (m *BlogModel) Update(ctx context.Context, patch BlogPatch) (*Blog, error) {
	return m.update(ctx, m.DB, patch)
}

func (m *BlogModel) UpdateTx(ctx context.Context, tx pgx.Tx, patch BlogPatch) (*Blog, error) {
	return m.update(ctx, tx, patch)
}

func (m *BlogModel) delete(ctx context.Context, q db.Queryable, id uuid.UUID) (*Blog, error) {
	// TODO: Implement
	return nil, nil
}

func (m *BlogModel) Delete(ctx context.Context, id uuid.UUID) (*Blog, error) {
	return m.delete(ctx, m.DB, id)
}

func (m *BlogModel) DeleteTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Blog, error) {
	return m.delete(ctx, tx, id)
}
