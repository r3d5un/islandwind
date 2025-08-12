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

	var blog data.Blog

	t.Run("Insert", func(t *testing.T) {
		inserted, err := models.Blogs.Insert(ctx, data.BlogInput{
			Title:     "Test",
			Content:   "Some example content",
			Published: true,
		})
		assert.NoError(t, err)

		t.Logf("blog post inserted: %v\n", inserted)

		blog = *inserted
	})

	t.Log(blog)

	t.Run("Select", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("SelectMany", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Update", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("Delete", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
