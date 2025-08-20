package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/auth/handlers"
	"github.com/r3d5un/islandwind/internal/auth/middleware"
	"github.com/rs/cors"
)

type authType int

const (
	basicAuth authType = iota
	accessToken
	noAuth
)

func (m *Module) addRoutes(ctx context.Context) {
	routes := []struct {
		Path     string `json:"path"`
		Handler  http.HandlerFunc
		Methods  []string `json:"methods"`
		authType authType
	}{
		// healthcheck
		{
			"GET /api/v1/auth/healthcheck",
			m.healthcheckHandler,
			[]string{http.MethodGet},
			noAuth,
		},
		// login
		{
			"POST /api/v1/auth/login",
			handlers.LoginHandler(m.repo.Tokens),
			[]string{http.MethodPost},
			// Basic authentication should only be used for logging in. Other resources
			// should be accessible with access tokens.
			basicAuth,
		},
		// refresh
		{
			"POST /api/v1/auth/refresh",
			handlers.RefreshHandler(m.repo.Tokens),
			[]string{http.MethodPost},
			// The RefreshHandler authenticates and validates the request as part of the
			// refresh process. No extra auth required.
			noAuth,
		},
	}

	m.logger.LogAttrs(ctx, slog.LevelInfo, "adding routes")
	for _, route := range routes {
		m.logger.LogAttrs(ctx, slog.LevelInfo, "adding route", slog.Group(
			"route",
			slog.Any("method", route.Methods),
			slog.String("path", route.Path),
			slog.Any("authRequired", route.authType),
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
				switch route.authType {
				case noAuth:
					return next
				case basicAuth:
					handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r)
					})
					return middleware.BasicAuthMiddleware(handlerFunc, m.cfg.Server.BasicAuth)
				case accessToken:
					fallthrough
				default:
					return m.AccessTokenMiddleware(next)

				}
			},
		)

		m.mux.Handle(route.Path, chain.Then(route.Handler))
	}
}
