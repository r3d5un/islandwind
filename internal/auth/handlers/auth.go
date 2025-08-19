package handlers

import (
	"net/http"

	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/auth/repo"
)

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func LoginHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := tokens.CreateAccessToken()
		if err != nil {
			api.ServerErrorResponse(w, r, err)
			return
		}
		refreshToken, err := tokens.CreateRefreshToken(r.Context())
		if err != nil {
			api.ServerErrorResponse(w, r, err)
			return
		}

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			LoginResponse{
				AccessToken:  *accessToken,
				RefreshToken: *refreshToken,
			},
			nil,
		)
	}
}
