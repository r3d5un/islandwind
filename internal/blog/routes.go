package blog

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/r3d5un/islandwind/internal/blog/handlers"
)

func (m *Module) addRoutes(ctx context.Context) {
	routes := []struct {
		Path    string `json:"path"`
		Handler http.HandlerFunc
	}{
		// healthcheck
		{"GET /api/v1/blog/healthcheck", m.healthcheckHandler},
		// blog posts
		{"POST /api/v1/blog/post", handlers.PostBlogpostHandler(m.repo.Posts)},
		{"GET /api/v1/blog/post/{id}", handlers.GetBlogpostHandler(m.repo.Posts)},
		{"GET /api/v1/blog/post/", handlers.ListBlogpostHandler(m.repo.Posts)},
		{"PATCH /api/v1/blog/post", handlers.PatchBlogpostHandler(m.repo.Posts)},
		{"DELETE /api/v1/blog/post", handlers.DeleteBlogpostHandler(m.repo.Posts)},
	}

	m.logger.LogAttrs(ctx, slog.LevelInfo, "adding routes")
	for _, route := range routes {
		m.logger.LogAttrs(ctx, slog.LevelInfo, "adding route", slog.String("route", route.Path))
		m.mux.Handle(route.Path, route.Handler)
	}
}
