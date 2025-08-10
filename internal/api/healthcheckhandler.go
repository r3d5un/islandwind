package api

import (
	"net/http"
)

type HealthCheckMessage struct {
	Status string `json:"status"`
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, r, http.StatusOK, HealthCheckMessage{Status: "available"}, nil)
}
