package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/auth/repo"
)

func AccessTokenMiddleware(next http.Handler, tokens repo.TokenService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.BadRequestResponse(
				w,
				r,
				errors.New("missing authorization header"),
				"missing authorization header",
			)
		}

		accessTokenString := strings.Split(authHeader, " ")
		if len(accessTokenString) != 2 {
			api.BadRequestResponse(
				w,
				r,
				errors.New("invalid authorization header format"),
				"invalid authorization header",
			)
			return
		}
		if accessTokenString[0] != "Bearer" {
			w.Header().Set("WWW-Authenticate", `Bearer error="invalid_request"`)
			http.Error(w, "invalid authorization header prefix", http.StatusUnauthorized)
			return
		}

		valid, err := tokens.Validate(r.Context(), accessTokenString[1])
		if err != nil || !valid {
			api.UnauthorizedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
