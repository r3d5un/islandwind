package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/auth/config"
	"github.com/r3d5un/islandwind/internal/auth/repo"
	"github.com/r3d5un/islandwind/internal/logging"
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

		valid, err := tokens.Validate(r.Context(), repo.AccessTokenType, accessTokenString[1])
		if err != nil || !valid {
			api.UnauthorizedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func BasicAuthMiddleware(next http.Handler, cfg config.BasicAuthConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.LoggerFromContext(ctx)

		logger.LogAttrs(ctx, slog.LevelInfo, "authenticating request")
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(cfg.Username))
			expectedPasswordHash := sha256.Sum256([]byte(cfg.Password))
			usernameMatch := subtle.ConstantTimeCompare(
				usernameHash[:],
				expectedUsernameHash[:],
			) == 1
			passwordMatch := subtle.ConstantTimeCompare(
				passwordHash[:],
				expectedPasswordHash[:],
			) == 1

			if usernameMatch && passwordMatch {
				logger.LogAttrs(ctx, slog.LevelInfo, "request authenticated")
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authentication", `Basic realm="restricted", charset="UTF-8"`)
		api.UnauthorizedResponse(w, r)
	})
}
