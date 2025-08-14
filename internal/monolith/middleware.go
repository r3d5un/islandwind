package monolith

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/logging"
)

type requestUrlKey string

const RequestUrlKey requestUrlKey = "requestUrlKey"

func (m *Monolith) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Monolith) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := m.logger.With(
			slog.Group(
				"request",
				slog.String("id", uuid.New().String()),
				slog.String("method", r.Method),
				slog.String("protocol", r.Proto),
				slog.String("url", r.URL.Path),
			),
			// TODO: Add group for module
		)
		ctx = logging.WithLogger(ctx, logger)
		ctx = context.WithValue(ctx, RequestUrlKey, r.URL.Path)

		logger.LogAttrs(ctx, slog.LevelInfo, "received request")
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.LogAttrs(ctx, slog.LevelInfo, "request completed")
	})
}
