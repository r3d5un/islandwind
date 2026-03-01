package blog

import (
	"net/http"

	"github.com/r3d5un/islandwind/internal/api"
)

type HealthCheckMessage struct {
	Module      string `json:"module"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}

func (m *Module) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	api.RespondWithJSON(
		w,
		r,
		http.StatusOK,
		HealthCheckMessage{
			Module:      moduleName,
			Environment: m.cfg.App.Environment,
			Status:      "available",
		},
		nil,
	)
}
