package auth

import (
	"net/http"

	"github.com/r3d5un/islandwind/internal/auth/middleware"
)

func (m *Module) AccessTokenMiddleware(next http.Handler) http.Handler {
	return middleware.AccessTokenMiddleware(next, m.repo.Tokens)
}
