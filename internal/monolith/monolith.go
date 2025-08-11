package monolith

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justinas/alice"
	"github.com/r3d5un/islandwind/internal/config"
	database "github.com/r3d5un/islandwind/internal/db"
	"github.com/spf13/viper"
)

type Monolith struct {
	cfg    *config.Config
	mux    *http.ServeMux
	logger *slog.Logger
	db     *pgxpool.Pool
	id     uuid.UUID
}

func NewMonolith() (*Monolith, error) {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	instanceID := uuid.New()
	logGroup := slog.Group(
		"instance",
		slog.String("name", viper.GetString("app.name")),
		slog.String("environment", viper.GetString("app.environment")),
		slog.String("id", instanceID.String()),
	)
	var handler slog.Handler
	switch cfg.App.Environment {
	case "testing":
		fallthrough
	case "production":
		handler = slog.NewJSONHandler(os.Stderr, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}
	logger := slog.New(handler).With(logGroup)
	slog.SetDefault(logger)

	logger.LogAttrs(ctx, slog.LevelInfo, "creating database pool", slog.Any("cfg", cfg.DB))
	db, err := database.OpenPool(ctx, cfg.DB)
	if err != nil {
		return nil, err
	}

	return &Monolith{
		id:     instanceID,
		cfg:    cfg,
		mux:    http.NewServeMux(),
		logger: slog.Default(),
		db:     db,
	}, nil
}

func (m *Monolith) Serve() error {
	ctx := context.Background()
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", m.cfg.Server.Port),
		Handler:      m.routes(),
		IdleTimeout:  time.Duration(m.cfg.Server.IdleTimeout) * time.Second,
		ReadTimeout:  time.Duration(m.cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(m.cfg.Server.WriteTimeout) * time.Second,
		ErrorLog:     slog.NewLogLogger(m.logger.Handler(), slog.LevelError),
	}
	srvLogGroup := slog.Group(
		"serverSettings",
		slog.String("addr", srv.Addr),
		slog.Any("idleTimeout", srv.IdleTimeout.Seconds()),
		slog.Any("readTimeout", srv.ReadTimeout.Seconds()),
		slog.Any("writeTimeout", srv.WriteTimeout.Seconds()),
	)

	shutdownError := make(chan error)
	go func() {
		ctx := context.Background()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		slog.LogAttrs(
			ctx, slog.LevelInfo, "shutting down server", slog.String("signal", s.String()),
		)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	m.logger.LogAttrs(ctx, slog.LevelInfo, "starting server", srvLogGroup)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err = <-shutdownError; err != nil {
		m.logger.LogAttrs(
			ctx, slog.LevelError,
			"unable to shutdown server",
			srvLogGroup, slog.String("error", err.Error()),
		)
		return err
	}
	m.logger.LogAttrs(ctx, slog.LevelInfo, "stopped server", srvLogGroup)

	return nil
}

func (m *Monolith) routes() http.Handler {
	m.logger.Info("creating standard middleware chain")
	standard := alice.New(
		m.recoverPanic,
		m.enableCORS,
		m.logRequest,
	)

	// healthcheck
	m.mux.HandleFunc("GET /api/v1/mono/healthcheck", m.healthcheckHandler)

	// profiling
	m.mux.HandleFunc("GET /debug/pprof/", http.DefaultServeMux.ServeHTTP)
	m.mux.HandleFunc("GET /debug/pprof/profile", http.DefaultServeMux.ServeHTTP)
	m.mux.HandleFunc("GET /debug/pprof/heap", http.DefaultServeMux.ServeHTTP)

	handler := standard.Then(m.mux)
	return handler
}

type ContextKey string

const InstanceIDKey string = "instanceID"

func WithInstanceID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, InstanceIDKey, id)
}

func InstanceIDFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(InstanceIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}
	return id
}
