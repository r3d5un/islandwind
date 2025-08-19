package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/r3d5un/islandwind/internal/auth/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	accessToken, err := authRepo.Tokens.CreateAccessToken()
	assert.NoError(t, err)

	mw := middleware.AccessTokenMiddleware(handler, authRepo.Tokens)

	t.Run("Authorize", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accessToken))

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
	})

	t.Run("NoToken", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusBadRequest)
	})

	t.Run("UnauthorizedNotBearer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("NotBearer %s", *accessToken))

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	})
}
