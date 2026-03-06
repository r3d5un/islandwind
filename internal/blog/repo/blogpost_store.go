package repo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/cache"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/logging"
)

type blogpostStore struct {
	models *data.Models
	cache  cache.Cache
}

func newBlogpostStore(models *data.Models, cache cache.Cache) blogpostStore {
	return blogpostStore{models: models, cache: cache}
}

func (s *blogpostStore) Create(ctx context.Context, input PostInput) (*Post, error) {
	row, err := s.models.Posts.Insert(ctx, input.row())
	if err != nil {
		return nil, err
	}
	blogpost := s.newPostFromRow(*row)
	s.cache.Set(blogpost.ID, blogpost)

	return blogpost, nil
}

func (s *blogpostStore) Read(ctx context.Context, ID uuid.UUID) (*Post, error) {
	var blogpost Post
	var err error
	logger := logging.LoggerFromContext(ctx)

	if err := s.cache.Get(ID, blogpost); err == nil {
		return &blogpost, nil
	}
	switch {
	case errors.Is(err, cache.ErrCacheMiss):
		logger.LogAttrs(ctx, slog.LevelInfo, "cache miss")
	default:
		logger.LogAttrs(ctx, slog.LevelError, "unable to use cache")
	}

	row, err := s.models.Posts.SelectOne(ctx, ID)
	if err != nil {
		return nil, err
	}
	s.newPost(*row, &blogpost)

	return &blogpost, nil
}

func (s *blogpostStore) List(
	ctx context.Context,
	filter data.Filter,
) ([]*Post, *data.Metadata, error) {
	rows, metadata, err := s.models.Posts.SelectMany(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	blogposts := s.newBlogpostListFromRows(rows)

	return blogposts, metadata, err
}

func (s *blogpostStore) Update(ctx context.Context, patch PostPatch) (*Post, error) {
	row, err := s.models.Posts.Update(ctx, patch.row())
	if err != nil {
		return nil, err
	}
	blogpost := s.newPostFromRow(*row)

	if err := s.cache.Delete(blogpost.ID); err != nil {
		logging.LoggerFromContext(ctx).
			Error("unable to invalidate cache", slog.String("error", err.Error()))
	}

	return blogpost, nil
}

func (s *blogpostStore) Delete(ctx context.Context, ID uuid.UUID) error {
	_, err := s.models.Posts.Update(
		ctx,
		data.PostPatch{ID: ID, Deleted: new(false)},
	)
	if err != nil {
		return err
	}
	if err := s.cache.Delete(ID); err != nil {
		logging.LoggerFromContext(ctx).
			Error("unable to invalidate cache", slog.String("error", err.Error()))
	}

	return nil
}

func (s *blogpostStore) Restore(ctx context.Context, ID uuid.UUID) (*Post, error) {
	row, err := s.models.Posts.Update(
		ctx,
		data.PostPatch{ID: ID, Deleted: new(false)},
	)
	if err != nil {
		return nil, err
	}
	return s.newPostFromRow(*row), nil
}

func (s *blogpostStore) Purge(ctx context.Context, tx pgx.Tx, ID uuid.UUID) error {
	_, err := s.models.Posts.DeleteTx(ctx, tx, ID)
	if err != nil {
		return err
	}
	if err := s.cache.Delete(ID); err != nil {
		logging.LoggerFromContext(ctx).
			Error("unable to invalidate cache", slog.String("error", err.Error()))
	}

	return nil
}

func (s *blogpostStore) newPostFromRow(row data.Post) *Post {
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

func (s *blogpostStore) newBlogpostListFromRows(rows []*data.Post) []*Post {
	blogposts := make([]*Post, len(rows))
	for i, row := range rows {
		blogposts[i] = s.newPostFromRow(*row)
	}
	return blogposts
}

func (s *blogpostStore) newPost(row data.Post, post *Post) {
	post.ID = row.ID
	post.Title = row.Title
	post.Content = row.Content
	post.Published = row.Published
	post.CreatedAt = row.CreatedAt
	post.UpdatedAt = row.UpdatedAt
	post.Deleted = row.Deleted
	post.DeletedAt = db.NullTimeToPtr(row.DeletedAt)
}
