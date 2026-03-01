package monolith

import (
	"context"
	"net/http"
)

type Module interface {
	Start(ctx context.Context, mux *http.ServeMux)
	Shutdown()
}
