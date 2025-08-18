package api

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/r3d5un/islandwind/internal/logging"
)

// BasicAuthConfig contains the username and password used in the basic authentication
// for the HTTP server.
type BasicAuthConfig struct {
	// Username is the admin username.
	//
	// Set through the ISLANDWIND_SERVER_AUTHENTICATION_USERNAME environment variable.
	Username string `json:"username"`
	// Password is the admin password.
	//
	// Field is safe for logging as the [BasicAuthConfig] contains a custom [BasicAuthConfig.LogValue] method.
	//
	// Set through the ISLANDWIND_SERVER_AUTHENTICATION_PASSWORD environment variable.
	Password string `json:"password"`
}

func (c BasicAuthConfig) LogValue() slog.Value {
	return slog.GroupValue(slog.String("username", c.Username))
}

func BasicAuthMiddleware(next http.Handler, cfg BasicAuthConfig) http.Handler {
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
			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				logger.LogAttrs(ctx, slog.LevelInfo, "request authenticated")
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authentication", `Basic realm="restricted", charset="UTF-8"`)
		UnauthorizedResponse(w, r)
	})
}

// CORSMiddleware returns middleware enabeling Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler, allowedMethod string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", allowedMethod)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type requestUrlKey string

const RequestUrlKey requestUrlKey = "requestUrlKey"

// LogRequestMiddleware returns a middleware which adds a [slog.Logger] to the request context
func LogRequestMiddleware(next http.Handler, logger slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logger.With(
			slog.Group(
				"request",
				slog.String("method", r.Method),
				slog.String("protocol", r.Proto),
				slog.String("url", r.URL.Path),
			),
		)
		ctx = logging.WithLogger(ctx, logger)
		ctx = context.WithValue(ctx, RequestUrlKey, r.URL.Path)

		logger.LogAttrs(ctx, slog.LevelInfo, "received request")
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.LogAttrs(ctx, slog.LevelInfo, "request completed")
	})
}

// RecoverPanicMiddleware returns a middleware that recovers in case of a panic further down
// the chain
func RecoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
