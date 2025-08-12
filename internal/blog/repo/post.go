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
	// TODO: Implement
	return nil, nil, nil
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
	// TODO: Implement
	return nil, nil
}

func (r *PostRepository) SoftDelete(ctx context.Context, ID uuid.UUID) (*Post, error) {
	// TODO: Implement
	return nil, nil
}

func (r *PostRepository) Delete(ctx context.Context, ID uuid.UUID) (*Post, error) {
	// TODO: Implement
	return nil, nil
}
