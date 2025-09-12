package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3d5un/islandwind/internal/api"
	"github.com/r3d5un/islandwind/internal/blog/data"
	"github.com/r3d5un/islandwind/internal/blog/repo"
	"github.com/r3d5un/islandwind/internal/db"
	"github.com/r3d5un/islandwind/internal/testsuite"
	"github.com/r3d5un/islandwind/internal/validator"
)

type BlogpostResponse struct {
	Data repo.Post `json:"data"`
}

type BlogpostListResponse struct {
	Metadata data.Metadata `json:"metadata"`
	Data     []*repo.Post  `json:"data"`
}

type PostRequestBody struct {
	Data repo.PostInput `json:"data"`
}

type PatchRequestBody struct {
	Data repo.PostPatch `json:"data"`
}

type DeleteRequestBody struct {
	Data DeleteOptions `json:"data"`
}

type DeleteOptions struct {
	ID    uuid.UUID `json:"id"`
	Purge bool      `json:"purge"`
}

func PostBlogpostHandler(
	blogposts repo.PostWriter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body PostRequestBody
		if err := api.ReadJSON(r, &body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		blogpost, err := blogposts.Create(ctx, body.Data)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrUniqueConstraintViolation):
				api.ConstraintViolationResponse(w, r, err, "blogpost ID already exists")
			case errors.Is(err, context.DeadlineExceeded):
				api.TimeoutResponse(ctx, w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}
		testsuite.Assert(blogpost != nil, "blogpost should not be nil without errors", blogpost)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{
				Data: *blogpost,
			},
			nil,
		)
	}
}

func GetBlogpostHandler(
	blogposts repo.PostReader,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := api.ReadPathParamID(ctx, "id", r)
		if err != nil {
			api.InvalidParameterResponse(ctx, w, r, "id", err)
			return
		}

		blogpost, err := blogposts.Read(ctx, *id)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrRecordNotFound):
				api.NotFoundResponse(ctx, w, r)
			case errors.Is(err, context.DeadlineExceeded):
				api.TimeoutResponse(ctx, w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}
		testsuite.Assert(blogpost != nil, "blogpost should not be nil without errors", blogpost)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{
				Data: *blogpost,
			},
			nil,
		)
	}
}

func ListBlogpostHandler(
	blogposts repo.PostReader,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v := validator.New()
		qs := r.URL.Query()
		filters := data.Filter{}

		filters.PageSize = api.ReadRequiredQueryInt(qs, "page_size", 25, v)
		filters.ID = api.ReadOptionalQueryUUID(qs, "id", v)
		filters.Title = api.ReadOptionalQueryString(qs, "title")
		filters.CreatedAtFrom = api.ReadOptionalQueryDate(qs, "created_at_from", v)
		filters.CreatedAtTo = api.ReadOptionalQueryDate(qs, "created_at_to", v)
		filters.UpdatedAtFrom = api.ReadOptionalQueryDate(qs, "updated_at_from", v)
		filters.UpdatedAtTo = api.ReadOptionalQueryDate(qs, "updated_at_to", v)
		filters.Deleted = api.ReadOptionalQueryBoolean(qs, "deleted")
		filters.DeletedAtFrom = api.ReadOptionalQueryDate(qs, "deleted_at_from", v)
		filters.DeletedAtTo = api.ReadOptionalQueryDate(qs, "deleted_at_to", v)
		filters.LastSeen = *api.ReadRequiredQueryUUID(
			qs, "last_seen", v, uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		)

		if !v.Valid() {
			api.ValidationFailedResponse(ctx, w, r, v.Errors)
			return
		}

		blogposts, metadata, err := blogposts.List(ctx, filters)
		if err != nil {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				api.TimeoutResponse(ctx, w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}
		testsuite.Assert(blogposts != nil, "blogpost should not be nil without errors", blogposts)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostListResponse{
				Data:     blogposts,
				Metadata: *metadata,
			},
			nil,
		)
	}
}

func PatchBlogpostHandler(
	blogposts repo.PostWriter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body PatchRequestBody
		if err := api.ReadJSON(r, &body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		blogpost, err := blogposts.Update(ctx, body.Data)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrUniqueConstraintViolation):
				api.ConstraintViolationResponse(w, r, err, "blogpost ID already exists")
			case errors.Is(err, context.DeadlineExceeded):
				api.TimeoutResponse(ctx, w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}
		testsuite.Assert(blogposts != nil, "blogpost should not be nil without errors", blogposts)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{
				Data: *blogpost,
			},
			nil,
		)
	}
}

func DeleteBlogpostHandler(
	blogposts repo.PostWriter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var body DeleteRequestBody
		if err := api.ReadJSON(r, &body); err != nil {
			api.BadRequestResponse(w, r, err, "unable to parse JSON request body")
			return
		}

		var blogpost *repo.Post
		var err error
		switch body.Data.Purge {
		case true:
			blogpost, err = blogposts.Delete(ctx, body.Data.ID)
		default:
			change := true
			blogpost, err = blogposts.Update(
				ctx,
				repo.PostPatch{ID: body.Data.ID, Deleted: &change},
			)
		}
		if err != nil {
			switch {
			case errors.Is(err, db.ErrUniqueConstraintViolation):
				api.ConstraintViolationResponse(w, r, err, "blogpost ID already exists")
			case errors.Is(err, db.ErrForeignKeyConstraintViolation):
				api.ConstraintViolationResponse(w, r, err, "blogpost referenced by other resources")
			case errors.Is(err, db.ErrRecordNotFound):
				api.NotFoundResponse(ctx, w, r)
			case errors.Is(err, context.DeadlineExceeded):
				api.TimeoutResponse(ctx, w, r)
			default:
				api.ServerErrorResponse(w, r, err)
			}
			return
		}
		testsuite.Assert(blogposts != nil, "blogpost should not be nil without errors", blogposts)

		api.RespondWithJSON(
			w,
			r,
			http.StatusOK,
			BlogpostResponse{
				Data: *blogpost,
			},
			nil,
		)
	}
}
