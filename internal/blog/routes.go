package blog

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/blog/handlers"
	"github.com/rs/cors"
)

func (m *Module) addRoutes(ctx context.Context) {
	routes := []struct {
		Path         string `json:"path"`
		handler      http.HandlerFunc
		Method       string `json:"methods"`
		AuthRequried bool   `json:"authRequried"`
	}{
		// healthcheck
		{
			"/api/v1/blog/healthcheck",
			m.healthcheckHandler,
			http.MethodGet,
			false,
		},
		// blog posts
		{
			"/api/v1/blog/post",
			handlers.PostBlogpostHandler(m.repo.Posts),
			http.MethodPost,
			true,
		},
		{
			"/api/v1/blog/post/{id}",
			handlers.GetBlogpostHandler(m.repo.Posts),
			http.MethodGet,
			false,
		},
		{
			"/api/v1/blog/post",
			handlers.ListBlogpostHandler(m.repo.Posts),
			http.MethodGet,
			false,
		},
		{
			"/api/v1/blog/post",
			handlers.PatchBlogpostHandler(m.repo.Posts),
			http.MethodPatch,
			true,
		},
		{
			"/api/v1/blog/post",
			handlers.DeleteBlogpostHandler(m.repo.Posts),
			http.MethodDelete,
			true,
		},
		{
			"/api/v1/blog/post",
			api.CorsPreflightHandler(),
			http.MethodOptions,
			false,
		},
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodGet,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	m.logger.LogAttrs(ctx, slog.LevelInfo, "adding routes")
	for _, route := range routes {
		m.logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"adding route",
			slog.Group("route", slog.Any("route", route)),
		)

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
				return corsMiddleware.Handler(next)
			},
			// Require authentication for write requests
			func(next http.Handler) http.Handler {
				if !route.AuthRequried {
					return next
				}
				return m.auth.AccessTokenMiddleware(next)
			},
		)

		m.mux.Handle(
			fmt.Sprintf("%s %s", route.Method, route.Path),
			chain.Then(route.handler),
		)
	}
}
