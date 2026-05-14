package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/blog/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("Create", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		require.NoError(t, err)
		require.NotNil(t, created)

		t.Cleanup(func() {
			require.NoError(t, blog.Posts.Purge(ctx, created.ID))
		})
	})

	t.Run("Read", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, blog.Posts.Purge(ctx, created.ID))
		})

		read, err := blog.Posts.Read(ctx, created.ID)
		require.NoError(t, err)
		require.NotNil(t, read)
		assert.Equal(t, *created, *read)
	})

	t.Run("List", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, blog.Posts.Purge(ctx, created.ID))
		})

		list, metadata, err := blog.Posts.List(
			ctx,
			data.PostFilter{
				PageSize: 1,
				ID:       uuid.NullUUID{UUID: created.ID, Valid: true},
			},
		)
		require.NoError(t, err)
		require.NotNil(t, list)
		require.NotEmpty(t, list)
		assert.NotEmpty(t, metadata)
		assert.Equal(t, list[len(list)-1].ID, metadata.LastSeen)
	})

	t.Run("Update", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, blog.Posts.Purge(ctx, created.ID))
		})

		updated, err := blog.Posts.Update(
			ctx,
			repo.PostPatch{ID: created.ID, Deleted: new(false)},
		)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, created.ID, updated.ID)
		assert.Empty(t, updated.DeletedAt)
		assert.False(t, updated.Deleted)
	})

	t.Run("Delete", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, blog.Posts.Purge(ctx, created.ID))
		})

		require.NoError(t, blog.Posts.Delete(ctx, created.ID))
	})

	t.Run("Restore", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, blog.Posts.Purge(ctx, created.ID))
		})

		require.NoError(t, blog.Posts.Delete(ctx, created.ID))

		restored, err := blog.Posts.Restore(ctx, created.ID)
		require.NoError(t, err)
		require.NotEmpty(t, restored)
		assert.Equal(t, created.ID, restored.ID)
		assert.False(t, restored.Deleted)
		assert.Nil(t, restored.DeletedAt)
		assert.Equal(t, created.ID, restored.ID)
	})

	t.Run("Purge", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		require.NoError(t, err)

		require.NoError(t, blog.Posts.Purge(ctx, created.ID))
	})
}
