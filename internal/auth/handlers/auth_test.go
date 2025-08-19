package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/r3d5un/islandwind/internal/auth/handlers"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandlers(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Logf("cleaning up test: %s", t.Name())
		defer cancel()
	})

	var login handlers.LoginResponse

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

	t.Log(login)
}
