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
		handler := handlers.PostBlogpostHandler(blogReaderWriter)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		err = json.Unmarshal(rr.Body.Bytes(), &post)
		assert.NoError(t, err)
	})

	t.Run("GetBlogpostHandler", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)
		req.SetPathValue("id", post.Data.ID.String())

		rr := httptest.NewRecorder()
		handler := handlers.GetBlogpostHandler(blogReaderWriter)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		var resp handlers.BlogpostResponse
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, post, resp)
	})

	t.Run("ListBlogpostHandler", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.ListBlogpostHandler(blogReaderWriter)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		var resp handlers.BlogpostListResponse
		err = json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)

		t.Log(resp)
	})

	t.Run("PatchBlogpostHandler", func(t *testing.T) {
		change := false
		body, err := json.Marshal(handlers.PatchRequestBody{
			Data: repo.PostPatch{
				ID:        post.Data.ID,
				Published: &change,
			},
		})
		assert.NoError(t, err)

		req, err := http.NewRequest(
			http.MethodPatch, "", strings.NewReader(string(body)),
		)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.PatchBlogpostHandler(blogReaderWriter)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		err = json.Unmarshal(rr.Body.Bytes(), &post)
		assert.NoError(t, err)
		assert.False(t, post.Data.Published)
	})

	t.Run("DeleteBlogpostHandlerSoftDelete", func(t *testing.T) {
		body, err := json.Marshal(handlers.DeleteRequestBody{
			Data: handlers.DeleteOptions{
				ID:    post.Data.ID,
				Purge: false,
			},
		})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodDelete, "", strings.NewReader(string(body)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.DeleteBlogpostHandler(blogReaderWriter)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)
	})

	t.Run("DeleteBlogpostHandlerPurge", func(t *testing.T) {
		body, err := json.Marshal(handlers.DeleteRequestBody{
			Data: handlers.DeleteOptions{
				ID:    post.Data.ID,
				Purge: true,
			},
		})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodDelete, "", strings.NewReader(string(body)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.DeleteBlogpostHandler(blogReaderWriter)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)
	})
}
