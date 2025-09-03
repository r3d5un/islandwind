package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"

	"github.com/r3d5un/islandwind/internal/logging"
)

type requestIDKey string

const RequestIDKey requestIDKey = "requestIDKey"

// LogRequestMiddleware returns a middleware which adds a [slog.Logger] to the request context
func LogRequestMiddleware(next http.Handler, logger slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := uuid.New()
		logger := logger.With(
			slog.Group(
				"request",
				slog.String("method", r.Method),
				slog.String("protocol", r.Proto),
				slog.String("url", r.URL.Path),
			),
		)
		ctx = logging.WithLogger(ctx, logger)
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		logger.LogAttrs(ctx, slog.LevelInfo, "received request")
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.LogAttrs(ctx, slog.LevelInfo, "request completed")
	})
}

func RequestIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(RequestIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.UUID{}
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
