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
	"github.com/r3d5un/islandwind/internal/db/builder"
	"github.com/r3d5un/islandwind/internal/logging"
)

// Post is the database record for a blog post.
type Post struct {
	ID        uuid.UUID    `json:"id"        db:"id"`
	Title     string       `json:"title"     db:"title"`
	Content   string       `json:"content"   db:"content"`
	Published bool         `json:"published" db:"published"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt" db:"updated_at"`
	Deleted   bool         `json:"deleted"   db:"deleted"`
	DeletedAt sql.NullTime `json:"deletedAt" db:"deleted_at"`
}

var postColumns = builder.ColumnsFrom(Post{})

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
	p, err := m.scan(q.QueryRow(
		ctx,
		stmt,
		input.Title,
		input.Content,
		input.Published,
	))
	if err != nil {
		return nil, db.HandleError(ctx, err)
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
	stmt, args := builder.From("blog.post").
		Where(builder.NewGenericPredicate("id", builder.Equal, id)).
		Select(postColumns...)

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("query", logging.MinifySQL(stmt)),
		slog.Any("args", args),
		slog.Duration("timeout", *m.Timeout),
	))

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	p, err := m.scan(q.QueryRow(ctx, stmt, args))
	if err != nil {
		return nil, db.HandleError(ctx, err)
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

type PostFilter struct {
	ID            sql.Null[uuid.UUID] `json:"id"`
	Title         sql.Null[string]    `json:"title"`
	Published     sql.Null[bool]      `json:"published"`
	CreatedAtFrom sql.Null[time.Time] `json:"createdAtFrom"`
	CreatedAtTo   sql.Null[time.Time] `json:"createdAtTo"`
	UpdatedAtFrom sql.Null[time.Time] `json:"updatedAtFrom"`
	UpdatedAtTo   sql.Null[time.Time] `json:"updatedAtTo"`
	Deleted       sql.Null[bool]      `json:"deleted"`
	DeletedAtFrom sql.Null[time.Time] `json:"deletedAtFrom"`
	DeletedAtTo   sql.Null[time.Time] `json:"deletedAtTo"`

	LastSeen uuid.UUID `json:"lastSeen"`
	PageSize int       `json:"pageSize"`
}

func (m *PostModel) selectMany(
	ctx context.Context,
	q db.Queryable,
	filter PostFilter,
) ([]*Post, *Metadata, error) {
	stmt, args := builder.
		From("blog.post").
		Where(
			builder.NewNullPredicate("id", builder.Equal, filter.ID),
			builder.NewNullPredicate("title", builder.Equal, filter.Title),
			builder.NewNullPredicate("published", builder.Equal, filter.Published),
			builder.NewNullPredicate("deleted", builder.Equal, filter.Deleted),
			builder.NewNullPredicate("deleted_at", builder.GreaterOrEqual, filter.DeletedAtFrom),
			builder.NewNullPredicate("deleted_to", builder.Less, filter.DeletedAtTo),
			builder.NewNullPredicate("created_at", builder.GreaterOrEqual, filter.CreatedAtFrom),
			builder.NewNullPredicate("created_to", builder.Less, filter.CreatedAtFrom),
			builder.NewNullPredicate("updated_at", builder.GreaterOrEqual, filter.UpdatedAtFrom),
			builder.NewNullPredicate("updated_at", builder.Less, filter.UpdatedAtFrom),
			builder.NewGenericPredicate("id", builder.Greater, filter.LastSeen),
		).
		OrderBy(
			builder.OrderBy{Column: "created_at", Order: builder.Asc},
			builder.OrderBy{Column: "id", Order: builder.Asc},
		).
		Limit(filter.PageSize).
		Select(postColumns...)

	ctx, cancel := context.WithTimeout(ctx, *m.Timeout)
	defer cancel()

	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"query",
		slog.String("statement", logging.MinifySQL(stmt)),
		slog.Any("filter", filter),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "performing query")
	rows, err := q.Query(ctx, stmt, args)
	if err != nil {
		return nil, nil, db.HandleError(ctx, err)
	}

	posts := make([]*Post, filter.PageSize)
	i := 0
	for rows.Next() {
		p, err := m.scan(rows)
		if err != nil {
			return nil, nil, db.HandleError(ctx, err)
		}
		posts[i] = &p
		i++
	}
	posts = posts[:i]
	if err = rows.Err(); err != nil {
		return nil, nil, db.HandleError(ctx, err)
	}
	metadata := Metadata{
		Next:           false,
		ResponseLength: len(posts),
	}
	if len(posts) > 0 {
		metadata.LastSeen = posts[metadata.ResponseLength-1].ID
		metadata.Next = true
	}

	return posts, &metadata, nil
}

func (m *PostModel) SelectMany(
	ctx context.Context,
	filter PostFilter,
) ([]*Post, *Metadata, error) {
	return m.selectMany(ctx, m.DB, filter)
}

func (m *PostModel) SelectManyTx(
	ctx context.Context,
	tx pgx.Tx,
	filter PostFilter,
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
	p, err := m.scan(q.QueryRow(
		ctx,
		stmt,
		patch.ID,
		patch.Title,
		patch.Content,
		patch.Published,
		patch.Deleted,
	))
	if err != nil {
		return nil, db.HandleError(ctx, err)
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

func (m *PostModel) delete(ctx context.Context, q db.Queryable, id uuid.UUID) error {
	const stmt string = `
DELETE
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
	_, err := q.Exec(ctx, stmt, id)
	if err != nil {
		return db.HandleError(ctx, err)
	}

	return nil
}

func (m *PostModel) Delete(ctx context.Context, id uuid.UUID) error {
	return m.delete(ctx, m.DB, id)
}

func (m *PostModel) DeleteTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	return m.delete(ctx, tx, id)
}

func (m *PostModel) scan(row pgx.Row) (Post, error) {
	var p Post
	err := row.Scan(
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
		return p, err
	}
	return p, err
}
