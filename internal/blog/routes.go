package blog

import (
	"context"
	"log/slog"
	"net/http"
)

func (m *Module) addRoutes(ctx context.Context) {
	routes := []struct {
		Path    string `json:"path"`
		Handler http.HandlerFunc
	}{
		{"GET /api/v1/blog/healthcheck", m.healthcheckHandler},
	}

	m.logger.LogAttrs(ctx, slog.LevelInfo, "adding routes")
	for _, route := range routes {
		m.logger.LogAttrs(ctx, slog.LevelInfo, "adding route", slog.String("route", route.Path))
		m.mux.Handle(route.Path, route.Handler)
	}
}
