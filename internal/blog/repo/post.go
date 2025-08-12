package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/logging"
	"github.com/r3d5un/islandwind/internal/testsuite"
)

type Post struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Published bool       `json:"published"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Deleted   bool       `json:"deleted"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func newPostFromRow(row data.Post) *Post {
	return &Post{
		ID:        row.ID,
		Title:     row.Title,
		Content:   row.Content,
		Published: row.Published,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		Deleted:   row.Deleted,
		DeletedAt: db.NullTimeToPtr(row.DeletedAt),
	}
}

type PostInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

func (p *PostInput) postInputRow() data.PostInput {
	return data.PostInput{
		Title:     p.Title,
		Content:   p.Content,
		Published: p.Published,
	}
}

type PostPatch struct {
	ID        uuid.UUID `json:"id"`
	Title     *string   `json:"title"`
	Content   *string   `json:"content"`
	Published *bool     `json:"published"`
	Deleted   *bool     `json:"deleted"`
}

func (p *PostPatch) postPatchRow() data.PostPatch {
	return data.PostPatch{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		Published: p.Published,
		Deleted:   p.Deleted,
	}
}

type PostReader interface {
	Read(ctx context.Context, ID uuid.UUID) (*Post, error)
	List(ctx context.Context, filter data.Filter) ([]*Post, *data.Metadata, error)
}

type PostWriter interface {
	Create(ctx context.Context, input PostInput) (*Post, error)
	Update(ctx context.Context, patch PostPatch) (*Post, error)
	SoftDelete(ctx context.Context, ID uuid.UUID) (*Post, error)
	Restore(ctx context.Context, ID uuid.UUID) (*Post, error)
	Delete(ctx context.Context, ID uuid.UUID) (*Post, error)
}

type PostReaderWriter interface {
	PostReader
	PostWriter
}

type PostRepository struct {
	db     *pgxpool.Pool
	models data.Models
}

func newPostRepository(db *pgxpool.Pool, timeout *time.Duration) PostReaderWriter {
	return &PostRepository{
		db:     db,
		models: data.NewModels(db, timeout),
	}
}

func (r *PostRepository) Read(ctx context.Context, ID uuid.UUID) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"post",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "reading blog post")
	row, err := r.models.Posts.SelectOne(ctx, ID)
	if err != nil {
		return nil, err
	}
	testsuite.Assert(row != nil, "blog post database record is nil", nil)
	post := newPostFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post retrieved")

	return post, nil
}

func (r *PostRepository) List(
	ctx context.Context,
	filter data.Filter,
) ([]*Post, *data.Metadata, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"posts",
		slog.Any("filter", filter),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "reading blog posts")
	rows, metadata, err := r.models.Posts.SelectMany(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	testsuite.Assert(rows != nil, "blog post list are nil", nil)
	testsuite.Assert(metadata != nil, "blog post metadata are nil", nil)

	posts := make([]*Post, metadata.ResponseLength)
	for i, row := range rows {
		posts[i] = newPostFromRow(*row)
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog posts retrieved")

	return posts, metadata, nil
}

func (r *PostRepository) Create(ctx context.Context, input PostInput) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"newPost",
		slog.Any("input", input),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "creating blog post")
	row, err := r.models.Posts.Insert(ctx, input.postInputRow())
	if err != nil {
		return nil, err
	}
	testsuite.Assert(row != nil, "blog post database record is nil", nil)
	post := newPostFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post created")

	return post, nil
}

func (r *PostRepository) Update(ctx context.Context, patch PostPatch) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"changes",
		slog.Any("patch", patch),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "updating blog post")
	row, err := r.models.Posts.Update(ctx, patch.postPatchRow())
	if err != nil {
		return nil, err
	}
	testsuite.Assert(row != nil, "blog post database record is nil", nil)
	post := newPostFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post updated")

	return post, nil
}

func (r *PostRepository) SoftDelete(ctx context.Context, ID uuid.UUID) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"blogpost",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "soft deleting blog post")
	state := true
	row, err := r.models.Posts.Update(
		ctx,
		data.PostPatch{ID: ID, Deleted: &state},
	)
	if err != nil {
		return nil, err
	}
	testsuite.Assert(row != nil, "blog post database record is nil", nil)
	post := newPostFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post soft deleted")

	return post, nil
}

func (r *PostRepository) Restore(ctx context.Context, ID uuid.UUID) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"blogpost",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "restoring blog post")
	state := false
	row, err := r.models.Posts.Update(
		ctx,
		data.PostPatch{ID: ID, Deleted: &state},
	)
	if err != nil {
		return nil, err
	}
	testsuite.Assert(row != nil, "blog post database record is nil", nil)
	post := newPostFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post restored")

	return post, nil
}

func (r *PostRepository) Delete(ctx context.Context, ID uuid.UUID) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"blogpost",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "starting database transaction")
	tx, rollback, err := r.models.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer rollback()

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting blog post")
	state := true
	row, err := r.models.Posts.UpdateTx(
		ctx,
		tx,
		data.PostPatch{ID: ID, Deleted: &state},
	)
	if err != nil {
		return nil, err
	}
	row, err = r.models.Posts.DeleteTx(ctx, tx, row.ID)
	if err != nil {
		return nil, err
	}
	testsuite.Assert(row != nil, "blog post database record is nil", nil)

	logger.LogAttrs(ctx, slog.LevelInfo, "committing changes")
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	post := newPostFromRow(*row)
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post deleted")

	return post, nil
}
