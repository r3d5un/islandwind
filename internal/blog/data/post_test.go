package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlogModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	t.Run("Insert", func(t *testing.T) {
		inserted, err := models.Posts.Insert(ctx, data.PostInput{
			Title:     "Test",
			Content:   "Some example content",
			Published: true,
		})
		require.NoError(t, err)
		require.NotNil(t, inserted)

		t.Cleanup(func() {
			_, err := models.Posts.Delete(ctx, inserted.ID)
			require.NoError(t, err)
		})
	})

	t.Run("Select", func(t *testing.T) {
		inserted, err := models.Posts.Insert(ctx, data.PostInput{
			Title:     "Test",
			Content:   "Some example content",
			Published: true,
		})
		require.NoError(t, err)
		require.NotNil(t, inserted)

		t.Cleanup(func() {
			_, err := models.Posts.Delete(ctx, inserted.ID)
			require.NoError(t, err)
		})

		selected, err := models.Posts.SelectOne(ctx, inserted.ID)
		require.NoError(t, err)
		require.NotNil(t, selected)
		assert.Equal(t, *inserted, *selected)
	})

	t.Run("SelectMany", func(t *testing.T) {
		inserted, err := models.Posts.Insert(ctx, data.PostInput{
			Title:     "Test",
			Content:   "Some example content",
			Published: true,
		})
		require.NoError(t, err)
		require.NotNil(t, inserted)

		t.Cleanup(func() {
			_, err := models.Posts.Delete(ctx, inserted.ID)
			require.NoError(t, err)
		})

		selected, metadata, err := models.Posts.SelectMany(ctx, data.PostFilter{
			PageSize: 1,
			ID:       uuid.NullUUID{UUID: inserted.ID, Valid: true},
		})
		require.NoError(t, err)
		require.NotNil(t, selected)
		assert.NotEmpty(t, selected)
		assert.NotEmpty(t, metadata)
		assert.Equal(t, selected[len(selected)-1].ID, metadata.LastSeen)
	})

	t.Run("Update", func(t *testing.T) {
		inserted, err := models.Posts.Insert(ctx, data.PostInput{
			Title:     "Test",
			Content:   "Some example content",
			Published: true,
		})
		require.NoError(t, err)
		require.NotNil(t, inserted)

		t.Cleanup(func() {
			_, err := models.Posts.Delete(ctx, inserted.ID)
			require.NoError(t, err)
		})

		updated, err := models.Posts.Update(
			ctx,
			data.PostPatch{ID: inserted.ID, Deleted: new(true)},
		)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, inserted.ID, updated.ID)
		assert.NotEmpty(t, updated.DeletedAt)
		assert.True(t, updated.Deleted)
	})

	t.Run("Purge", func(t *testing.T) {
		inserted, err := models.Posts.Insert(ctx, data.PostInput{
			Title:     "Test",
			Content:   "Some example content",
			Published: true,
		})
		require.NoError(t, err)
		require.NotNil(t, inserted)

		_, err = models.Posts.Delete(ctx, inserted.ID)
		require.NoError(t, err)
	})
}
