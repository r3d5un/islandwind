package monolith

import (
	"net/http"

	"github.com/r3d5un/islandwind/internal/api"
)

type HealthCheckMessage struct {
	InstanceID  string `json:"instanceId"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}

func (m *Monolith) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	api.RespondWithJSON(
		w,
		r,
		http.StatusOK,
		HealthCheckMessage{
			InstanceID:  m.id.String(),
			Environment: m.cfg.App.Environment,
			Status:      "available",
		},
		nil,
	)
}
