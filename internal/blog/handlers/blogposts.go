package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/blog/repo"
	"github.com/r3d5un/islandwind/internal/db"
)

type BlogpostResponse struct {
	Data repo.Post `json:"data"`
}

type PostRequestBody struct {
	Data repo.PostInput `json:"data"`
}

type PatchRequestBody struct {
	Data repo.PostPatch `json:"data"`
}

type DeleteRequestBody struct {
	Data deleteOptions `json:"data"`
}

type deleteOptions struct {
	ID    uuid.UUID `json:"id"`
	Purge bool      `json:"purge"`
}

func PostBlogpostHandler(
	blogposts repo.PostReaderWriter,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body PostRequestBody
		if err := api.ReadJSON(r, body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		blogpost, err := blogposts.Create(ctx, body.Data)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrUniqueConstraintViolation):
				api.ConstraintViolationResponse(w, r, err, "blostpost ID already exists")
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
			BlogpostResponse{Data: *blogpost},
			nil,
		)
	})
}

func GetBlogpostHandler(
	blogposts repo.PostReaderWriter,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{},
			nil,
		)
	})
}

func ListBlogpostHandler(
	blogposts repo.PostReaderWriter,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{},
			nil,
		)
	})
}

func PatchBlogpostHandler(
	blogposts repo.PostReaderWriter,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{},
			nil,
		)
	})
}

func DeleteBlogpostHandler(
	blogposts repo.PostReaderWriter,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{},
			nil,
		)
	})
}
