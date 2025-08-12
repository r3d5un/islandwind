package data_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/stretchr/testify/assert"
)

func TestBlogModel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	var post data.Post

	t.Run("Insert", func(t *testing.T) {
		inserted, err := models.Posts.Insert(ctx, data.PostInput{
			Title:     "Test",
			Content:   "Some example content",
			Published: true,
		})
		assert.NoError(t, err)

		t.Logf("post post inserted: %v\n", inserted)

		post = *inserted
	})

	t.Run("Select", func(t *testing.T) {
		selected, err := models.Posts.SelectOne(ctx, post.ID)
		assert.NoError(t, err)
		assert.Equal(t, post, *selected)
	})

	t.Run("SelectMany", func(t *testing.T) {
		selected, metadata, err := models.Posts.SelectMany(ctx, data.Filter{
			PageSize: 1,
			ID:       &post.ID,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, selected)
		assert.NotEmpty(t, metadata)
		assert.Equal(t, selected[len(selected)-1].ID, metadata.LastSeen)
	})

	t.Run("Update", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Delete", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
