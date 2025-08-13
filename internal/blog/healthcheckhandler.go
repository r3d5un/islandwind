package blog

import (
	"net/http"

	"github.com/r3d5un/islandwind/internal/api"
)

type HealthCheckMessage struct {
	Module      string `json:"module"`
	InstanceID  string `json:"instanceId"`
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
			InstanceID:  m.instanceID.String(),
			Environment: m.cfg.App.Environment,
			Status:      "available",
		},
		nil,
	)
}
