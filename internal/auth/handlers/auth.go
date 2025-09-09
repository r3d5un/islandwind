package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/auth/data"
	"github.com/r3d5un/islandwind/internal/testsuite"
	"github.com/r3d5un/islandwind/internal/validator"

	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/auth/repo"
)

type Response struct {
	RequestID    uuid.UUID `json:"requestId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

type LogoutResponse struct {
	RequestID uuid.UUID `json:"requestId"`
}

type RefreshRequestBody struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenListResponse struct {
	RequestID     uuid.UUID            `json:"requestId"`
	Metadata      data.Metadata        `json:"metadata"`
	RefreshTokens []*repo.RefreshToken `json:"refreshTokens,omitzero"`
}

type RefreshTokenDeleteResponse struct {
	RequestID uuid.UUID `json:"requestId"`
	Data      struct {
		NumberDeleted int64 `json:"numberDeleted"`
	} `json:"data"`
}

func ListRefreshTokenHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v := validator.New()
		qs := r.URL.Query()
		filters := data.Filter{}

		addFilters(&filters, qs, v)
		if !v.Valid() {
			api.ValidationFailedResponse(ctx, w, r, v.Errors)
			return
		}

		refreshTokens, metadata, err := tokens.List(ctx, filters)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				api.TimeoutResponse(ctx, w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}
		testsuite.Assert(
			refreshTokens != nil, "refresh tokens should not be nil without errors", refreshTokens,
		)
		testsuite.Assert(
			metadata != nil, "refresh token metadata should not be nil without errors", metadata,
		)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			RefreshTokenListResponse{
				RequestID:     api.RequestIDFromContext(ctx),
				Metadata:      *metadata,
				RefreshTokens: refreshTokens,
			},
			nil,
		)
	}
}

func DeleteRefreshTokenHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v := validator.New()
		qs := r.URL.Query()
		filters := data.Filter{}

		addFilters(&filters, qs, v)
		if !v.Valid() {
			api.ValidationFailedResponse(ctx, w, r, v.Errors)
			return
		}

		numberDeleted, err := tokens.Delete(ctx, filters)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				api.TimeoutResponse(ctx, w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			RefreshTokenDeleteResponse{
				RequestID: api.RequestIDFromContext(ctx),
				Data: struct {
					NumberDeleted int64 `json:"numberDeleted"`
				}{
					NumberDeleted: numberDeleted,
				},
			},
			nil,
		)
	}
}

func LoginHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		accessToken, err := tokens.CreateAccessToken()
		if err != nil {
			api.ServerErrorResponse(w, r, err)
			return
		}
		testsuite.Assert(
			accessToken != nil, "accessToken should never be nil without errors", accessToken,
		)

		refreshToken, err := tokens.CreateRefreshToken(ctx)
		if err != nil {
			api.ServerErrorResponse(w, r, err)
			return
		}
		testsuite.Assert(
			refreshToken != nil, "refreshToken should never be nil without errors", refreshToken,
		)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			Response{
				RequestID:    api.RequestIDFromContext(ctx),
				AccessToken:  *accessToken,
				RefreshToken: *refreshToken,
			},
			nil,
		)
	}
}

func LogoutHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body RefreshRequestBody
		if err := api.ReadJSON(r, &body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		err := tokens.InvalidateRefreshToken(ctx, body.RefreshToken)
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
			LogoutResponse{RequestID: api.RequestIDFromContext(ctx)},
			nil,
		)
	}
}

func RefreshHandler(tokens repo.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body RefreshRequestBody
		if err := api.ReadJSON(r, &body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		accessToken, refreshToken, err := tokens.Refresh(ctx, body.RefreshToken)
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
		testsuite.Assert(
			accessToken != nil, "accessToken should never be nil without errors", accessToken,
		)
		testsuite.Assert(
			refreshToken != nil, "refreshToken should never be nil without errors", refreshToken,
		)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			Response{
				RequestID:    api.RequestIDFromContext(ctx),
				AccessToken:  *accessToken,
				RefreshToken: *refreshToken,
			},
			nil,
		)
	}
}

func addFilters(filters *data.Filter, qs url.Values, v *validator.Validator) {
	filters.PageSize = api.ReadRequiredQueryInt(qs, "page_size", 25, v)
	filters.ID = api.ReadOptionalQueryUUID(qs, "id", v)
	filters.Issuer = api.ReadOptionalQueryString(qs, "id")
	filters.IssuedAtFrom = api.ReadOptionalQueryDate(qs, "issued_at_from", v)
	filters.IssuedAtTo = api.ReadOptionalQueryDate(qs, "issued_at_to", v)
	filters.ExpirationFrom = api.ReadOptionalQueryDate(qs, "expiration_from", v)
	filters.ExpirationTo = api.ReadOptionalQueryDate(qs, "expiration_to", v)
	filters.Invalidated = api.ReadOptionalQueryBoolean(qs, "invalidated")
	filters.InvalidatedBy = api.ReadOptionalQueryUUID(qs, "invalidated_by", v)
}
