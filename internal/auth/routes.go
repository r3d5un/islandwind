package auth

import (
	"context"
	"fmt"
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
		Path          string `json:"path"`
		handler       http.HandlerFunc
		Method        string `json:"method"`
		authType      authType
		CorsPreflight bool `json:"corsPreflight"`
	}{
		// healthcheck
		{
			"/api/v1/auth/healthcheck",
			m.healthcheckHandler,
			http.MethodGet,
			noAuth,
			false,
		},
		// login
		{
			"/api/v1/auth/login",
			handlers.LoginHandler(m.repo.Tokens),
			http.MethodPost,
			// Basic authentication should only be used for logging in. Other resources
			// should be accessible with access tokens.
			basicAuth,
			true,
		},
		// refresh
		{
			"/api/v1/auth/refresh",
			handlers.RefreshHandler(m.repo.Tokens),
			http.MethodPost,
			// The RefreshHandler authenticates and validates the request as part of the
			// refresh process. No extra auth required.
			noAuth,
			true,
		},
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
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
				switch route.authType {
				case noAuth:
					return next
				case basicAuth:
					handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						next.ServeHTTP(w, r)
					})
					return middleware.BasicAuthMiddleware(handlerFunc, m.cfg.BasicAuth)
				case accessToken:
					fallthrough
				default:
					return m.AccessTokenMiddleware(next)

				}
			},
		)

		if route.CorsPreflight {
			m.mux.Handle(
				fmt.Sprintf("%s %s", http.MethodOptions, route.Path),
				chain.Then(api.CorsPreflightHandler()),
			)
		}

		m.mux.Handle(
			fmt.Sprintf("%s %s", route.Method, route.Path),
			chain.Then(route.handler),
		)
	}
}
