package blog

import (
	"context"
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
		Handler      http.HandlerFunc
		Methods      []string `json:"methods"`
		AuthRequried bool
	}{
		// healthcheck
		{
			"GET /api/v1/blog/healthcheck",
			m.healthcheckHandler,
			[]string{http.MethodGet},
			false,
		},
		// blog posts
		{
			"POST /api/v1/blog/post",
			handlers.PostBlogpostHandler(m.repo.Posts),
			[]string{http.MethodPost},
			true,
		},
		{
			"GET /api/v1/blog/post/{id}",
			handlers.GetBlogpostHandler(m.repo.Posts),
			[]string{http.MethodGet},
			false,
		},
		{
			"GET /api/v1/blog/post",
			handlers.ListBlogpostHandler(m.repo.Posts),
			[]string{http.MethodGet},
			false,
		},
		{
			"PATCH /api/v1/blog/post",
			handlers.PatchBlogpostHandler(m.repo.Posts),
			[]string{http.MethodPatch},
			true,
		},
		{
			"DELETE /api/v1/blog/post",
			handlers.DeleteBlogpostHandler(m.repo.Posts),
			[]string{http.MethodDelete},
			true,
		},
	}

	m.logger.LogAttrs(ctx, slog.LevelInfo, "adding routes")
	for _, route := range routes {
		m.logger.LogAttrs(ctx, slog.LevelInfo, "adding route", slog.Group(
			"route",
			slog.Any("method", route.Methods),
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
				c := cors.New(cors.Options{
					AllowedOrigins:     []string{"*"},
					AllowedMethods:     route.Methods,
					AllowedHeaders:     []string{"Content-Type", "Authorization"},
					AllowCredentials:   false,
					Debug:              true,
					OptionsPassthrough: false,
				})
				return c.Handler(next)
			},
			// Require authentication for write requests
			func(next http.Handler) http.Handler {
				if !route.AuthRequried {
					return next
				}
				return m.auth.AccessTokenMiddleware(next)
			},
		)

		m.mux.Handle(route.Path, chain.Then(route.Handler))
	}
}
