package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/cache"
	"github.com/r3d5un/islandwind/internal/logging"
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

type PostInput struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

func (p *PostInput) row() data.PostInput {
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

func (p *PostPatch) row() data.PostPatch {
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
	Delete(ctx context.Context, ID uuid.UUID) error
	Restore(ctx context.Context, ID uuid.UUID) (*Post, error)
	Purge(ctx context.Context, ID uuid.UUID) error
}

type PostReaderWriter interface {
	PostReader
	PostWriter
}

type BlogpostService struct {
	db            *pgxpool.Pool
	cache         cache.Cache
	models        data.Models
	blogpostStore blogpostStore
}

func newPostRepository(
	db *pgxpool.Pool,
	cache cache.Cache,
	timeout *time.Duration,
) PostReaderWriter {
	models := data.NewModels(db, timeout)
	return &BlogpostService{
		db:            db,
		cache:         cache,
		models:        models,
		blogpostStore: newBlogpostStore(&models, cache),
	}
}

func (svc *BlogpostService) Read(ctx context.Context, ID uuid.UUID) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"post",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "reading blog post")
	blogpost, err := svc.blogpostStore.Read(ctx, ID)
	if err != nil {
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post retrieved")

	return blogpost, nil
}

func (svc *BlogpostService) List(
	ctx context.Context,
	filter data.Filter,
) ([]*Post, *data.Metadata, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"posts",
		slog.Any("filter", filter),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "reading blog posts")
	posts, metadata, err := svc.blogpostStore.List(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog posts retrieved")

	return posts, metadata, nil
}

func (svc *BlogpostService) Create(ctx context.Context, input PostInput) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"newPost",
		slog.Any("input", input),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "creating blog post")
	blogpost, err := svc.blogpostStore.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post created")

	return blogpost, nil
}

func (svc *BlogpostService) Update(ctx context.Context, patch PostPatch) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"changes",
		slog.Any("patch", patch),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "updating blog post")
	post, err := svc.blogpostStore.Update(ctx, patch)
	if err != nil {
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post updated")

	return post, nil
}

func (svc *BlogpostService) Delete(ctx context.Context, ID uuid.UUID) error {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"blogpost",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "deleting blog post")
	err := svc.blogpostStore.Delete(ctx, ID)
	if err != nil {
		return err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post soft deleted")

	return nil
}

func (svc *BlogpostService) Restore(ctx context.Context, ID uuid.UUID) (*Post, error) {
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"blogpost",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "restoring blog post")
	blogpost, err := svc.blogpostStore.Restore(ctx, ID)
	if err != nil {
		return nil, err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post restored")

	return blogpost, nil
}

func (svc *BlogpostService) Purge(ctx context.Context, ID uuid.UUID) error {
	// TODO: Refactor the Purge method
	//  - Add a X-Purge header to the DELETE endpoint.
	//  - Add a softDelete helper function
	//  - Add a purge helper function
	logger := logging.LoggerFromContext(ctx).With(slog.Group(
		"blogpost",
		slog.String("id", ID.String()),
	))

	logger.LogAttrs(ctx, slog.LevelInfo, "purging blog post")
	tx, rollback, err := svc.models.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer rollback()

	err = svc.blogpostStore.Purge(ctx, tx, ID)
	if err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelInfo, "committing changes")
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "blog post deleted")

	return nil
}
