package handlers

import (
	"errors"
	"net/http"

	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/auth/repo"
)

type Response struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequestBody struct {
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
			Response{
				AccessToken:  *accessToken,
				RefreshToken: *refreshToken,
			},
			nil,
		)
	}
}

func LogoutHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body RefreshRequestBody
		if err := api.ReadJSON(r, &body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		err := tokens.InvalidateRefreshToken(r.Context(), body.RefreshToken)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrVerifyingToken), errors.Is(err, repo.ErrParsingToken):
				api.BadRequestResponse(w, r, err, "unable to verify token")
			case errors.Is(err, repo.ErrTokenExpired):
				api.BadRequestResponse(w, r, err, "token expired")
			case errors.Is(err, repo.ErrInvalidIssuedAt):
				api.BadRequestResponse(w, r, err, "token iat invalid")
			case errors.Is(err, repo.ErrIssuerMismatch), errors.Is(err, repo.ErrUnauthorized):
				api.UnauthorizedResponse(w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			nil,
			nil,
		)
	}
}

func RefreshHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body RefreshRequestBody
		if err := api.ReadJSON(r, &body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		accessToken, refreshToken, err := tokens.Refresh(r.Context(), body.RefreshToken)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrVerifyingToken), errors.Is(err, repo.ErrParsingToken):
				api.BadRequestResponse(w, r, err, "unable to verify token")
			case errors.Is(err, repo.ErrTokenExpired):
				api.BadRequestResponse(w, r, err, "token expired")
			case errors.Is(err, repo.ErrInvalidIssuedAt):
				api.BadRequestResponse(w, r, err, "token iat invalid")
			case errors.Is(err, repo.ErrIssuerMismatch), errors.Is(err, repo.ErrUnauthorized):
				api.UnauthorizedResponse(w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			Response{
				AccessToken:  *accessToken,
				RefreshToken: *refreshToken,
			},
			nil,
		)
	}
}
