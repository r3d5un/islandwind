package api

import "net/http"

func EmptyHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, r, http.StatusOK, nil, nil)
	})
}
