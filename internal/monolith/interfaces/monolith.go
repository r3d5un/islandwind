package interfaces

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/r3d5un/islandwind/internal/config"
)

type Monolith interface {
	DB() *pgxpool.Pool
	Mux() *http.ServeMux
	Logger() *slog.Logger
	Config() *config.Config
	Modules() *Modules
}

type Modules struct {
	Blog BlogService
}
