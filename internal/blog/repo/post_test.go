package repo_test

import (
	"context"
	"testing"
	"time"

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

	t.Log(post)

	t.Run("Read", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("List", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Update", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("SoftDelete", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Delete", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
