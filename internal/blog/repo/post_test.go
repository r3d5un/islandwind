package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/blog/repo"
	"github.com/stretchr/testify/assert"
)

func TestPostRepository(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	var post repo.Post

	t.Run("Create", func(t *testing.T) {
		created, err := blog.Posts.Create(
			ctx,
			repo.PostInput{
				Title:     "Example Title",
				Content:   "Some placeholder content",
				Published: true,
			},
		)
		assert.NoError(t, err)
		assert.NotNil(t, created)

		post = *created
	})

	t.Run("Read", func(t *testing.T) {
		read, err := blog.Posts.Read(ctx, post.ID)
		assert.NoError(t, err)
		assert.NotNil(t, read)
		assert.Equal(t, post, *read)
	})

	t.Run("List", func(t *testing.T) {
		list, metadata, err := blog.Posts.List(
			ctx,
			data.Filter{
				PageSize: 1,
				ID:       &post.ID,
			},
		)
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.NotEmpty(t, list)
		assert.NotEmpty(t, metadata)
		assert.Equal(t, list[len(list)-1].ID, metadata.LastSeen)
	})

	t.Run("Update", func(t *testing.T) {
		del := false
		updated, err := blog.Posts.Update(
			ctx,
			repo.PostPatch{ID: post.ID, Deleted: &del},
		)
		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, post.ID, updated.ID)
		assert.Empty(t, updated.DeletedAt)
		assert.False(t, updated.Deleted)

		post = *updated
	})

	t.Run("SoftDelete", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Delete", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
