package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/r3d5un/islandwind/internal/auth/handlers"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	var login handlers.Response

	t.Run("LoginHandler", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.LoginHandler(authRepo.Tokens)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		err = json.Unmarshal(rr.Body.Bytes(), &login)
		assert.NoError(t, err)
	})

	t.Run("LogoutHandler", func(t *testing.T) {
		token, err := authRepo.Tokens.CreateRefreshToken(ctx)
		assert.NoError(t, err)
		body, err := json.Marshal(handlers.RefreshRequestBody{RefreshToken: *token})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.LogoutHandler(authRepo.Tokens)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		err = json.Unmarshal(rr.Body.Bytes(), &login)
		assert.NoError(t, err)
	})

	t.Run("RefreshHandler", func(t *testing.T) {
		body, err := json.Marshal(handlers.RefreshRequestBody{
			RefreshToken: login.RefreshToken,
		})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.RefreshHandler(authRepo.Tokens)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotNil(t, rr.Body)

		err = json.Unmarshal(rr.Body.Bytes(), &login)
		assert.NoError(t, err)
	})

	t.Run("RefreshHandlerUnauthorized", func(t *testing.T) {
		body, err := json.Marshal(handlers.RefreshRequestBody{
			RefreshToken: login.AccessToken,
		})
		assert.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := handlers.RefreshHandler(authRepo.Tokens)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.NotNil(t, rr.Body)

		err = json.Unmarshal(rr.Body.Bytes(), &login)
		assert.NoError(t, err)
	})
}
