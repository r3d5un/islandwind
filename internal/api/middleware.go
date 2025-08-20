package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/r3d5un/islandwind/internal/logging"
)

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
