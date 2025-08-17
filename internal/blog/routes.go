package blog

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/blog/handlers"
)

func (m *Module) addRoutes(ctx context.Context) {
	routes := []struct {
		Path         string `json:"path"`
		Handler      http.HandlerFunc
		Method       string
		AuthRequried bool
	}{
		// healthcheck
		{
			"GET /api/v1/blog/healthcheck",
			m.healthcheckHandler,
			http.MethodGet,
			false,
		},
		// blog posts
		{
			"POST /api/v1/blog/post",
			handlers.PostBlogpostHandler(m.repo.Posts),
			http.MethodPost,
			true,
		},
		{
			"GET /api/v1/blog/post/{id}",
			handlers.GetBlogpostHandler(m.repo.Posts),
			http.MethodGet,
			false,
		},
		{
			"GET /api/v1/blog/post",
			handlers.ListBlogpostHandler(m.repo.Posts),
			http.MethodGet,
			false,
		},
		{
			"PATCH /api/v1/blog/post",
			handlers.PatchBlogpostHandler(m.repo.Posts),
			http.MethodPatch,
			true,
		},
		{
			"DELETE /api/v1/blog/post",
			handlers.DeleteBlogpostHandler(m.repo.Posts),
			http.MethodDelete,
			true,
		},
		{
			// Route for testing basic auth credentials
			"GET /api/v1/auth/login",
			api.EmptyHandler(),
			http.MethodGet,
			true,
		},
	}

	m.logger.LogAttrs(ctx, slog.LevelInfo, "adding routes")
	for _, route := range routes {
		m.logger.LogAttrs(ctx, slog.LevelInfo, "adding route", slog.Group(
			"route",
			slog.String("method", route.Method),
			slog.String("path", route.Path),
			slog.Bool("authRequired", route.AuthRequried),
		))

		chain := alice.New(
			// Add logging middleware for all requests
			func(next http.Handler) http.Handler {
				handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
				return api.LogRequestMiddleware(handlerFunc, *m.logger)
			},
			// Enable CORS for all requests
			func(next http.Handler) http.Handler {
				return api.CORSMiddleware(next, route.Method)
			},
			// Require authentication for write requests
			func(next http.Handler) http.Handler {
				if !route.AuthRequried {
					return next
				}
				handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
				return api.BasicAuthMiddleware(handlerFunc, m.cfg.Server.Authentication)
			},
		)

		m.mux.Handle(
			fmt.Sprintf("%s %s", route.Method, route.Path),
			chain.Then(route.Handler),
		)
	}
}
