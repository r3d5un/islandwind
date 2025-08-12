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
	logger.LogAttrs(ctx, slog.LevelInfo, "blog selected", slog.Any("blog", b))

	return &b, nil
}

func (m *BlogModel) Insert(ctx context.Context, input BlogInput) (*Blog, error) {
	return m.insert(ctx, m.DB, input)
}

func (m *BlogModel) InsertTx(ctx context.Context, tx pgx.Tx, input BlogInput) (*Blog, error) {
	return m.insert(ctx, tx, input)
}

func (m *BlogModel) selectOne(ctx context.Context, q db.Queryable, id uuid.UUID) (*Blog, error) {
	// TODO: Implement
	return nil, nil
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
	// TODO: Implement
	return nil, nil, nil
}

func (m *BlogModel) SelectMany(
	ctx context.Context,
	q db.Queryable,
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
