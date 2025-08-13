package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/r3d5un/islandwind/internal/blog/handlers"
	"github.com/r3d5un/islandwind/internal/blog/repo"
	"github.com/stretchr/testify/assert"
)

func TestBlogpostHandlers(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	var post handlers.BlogpostResponse

	t.Run("PostBlogpostHandler", func(t *testing.T) {
		body, err := json.Marshal(handlers.PostRequestBody{
			Data: repo.PostInput{
				Title:     "Example Title",
				Content:   "Some sample content",
				Published: true,
			},
		})
		assert.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPost, "", strings.NewReader(string(body)),
		)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		postHandler := handlers.PostBlogpostHandler(blogReaderWriter)
		postHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		err = json.Unmarshal(rr.Body.Bytes(), &post)
		assert.NoError(t, err)
	})

	t.Run("GetBlogpostHandler", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("ListBlogpostHandler", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("PatchBlogpostHandler", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("DeleteBlogpostHandler", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
