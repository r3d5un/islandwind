package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/auth/config"
	"github.com/r3d5un/islandwind/internal/auth/middleware"
	"github.com/stretchr/testify/assert"
)

func TestBasicAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	cfg := config.BasicAuthConfig{
		Username: "username",
		Password: "password",
	}

	mw := middleware.BasicAuthMiddleware(handler, cfg)

	t.Run("Authorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.SetBasicAuth(cfg.Username, cfg.Password)

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
	})

	t.Run("UnauthorizedIncorrectCredentials", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.SetBasicAuth("incorrect", "incorrect")

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	})

	t.Run("UnauthorizedNoCredential", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		rr := httptest.NewRecorder()
		mw.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	})
}

func TestRecoverPanicMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("the middleware should recover from this panic")
	})
	mw := api.RecoverPanicMiddleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		mw.ServeHTTP(rr, req)
	})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
