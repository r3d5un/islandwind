package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/r3d5un/islandwind/internal/logging"
)

var (
	ErrPathParamID = errors.New("path parameter is invalid")
)

const (
	notFoundMsg string = "resource not found"
	timeoutMsg  string = "the server took to long to respond"
)

type ErrorMessage struct {
	Message any `json:"message"`
}

func ErrorResponse(
	w http.ResponseWriter, r *http.Request, status int, message any,
) {
	ctx := r.Context()
	logger := logging.LoggerFromContext(ctx)

	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"writing error response",
		slog.Int("status", status),
		slog.Any("message", message),
	)
	RespondWithJSON(w, r, status, ErrorMessage{Message: message}, nil)
}

func LogError(r *http.Request, err error) {
	logging.LoggerFromContext(r.Context()).Error(
		"an error occurred",
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String()),
		slog.String("error", err.Error()),
	)
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	LogError(r, err)
	const serverErrorMsg string = "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, serverErrorMsg)
}

func RespondWithJSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	data any,
	headers http.Header,
) {
	logger := logging.LoggerFromContext(r.Context())

	js, err := json.Marshal(data)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	logger.LogAttrs(r.Context(), slog.LevelInfo, "writing response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(js); err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error, msg string) {
	logger := logging.LoggerFromContext(r.Context())

	logger.Info("bad request", slog.String("error", err.Error()), slog.String("message", msg))
	ErrorResponse(w, r, http.StatusBadRequest, msg)
}

func ConstraintViolationResponse(w http.ResponseWriter, r *http.Request, err error, msg string) {
	logger := logging.LoggerFromContext(r.Context())

	logger.Info(
		"a constraint violation occurred",
		slog.String("error", err.Error()),
		slog.String("message", msg),
	)
	ErrorResponse(w, r, http.StatusConflict, msg)
}

func TimeoutResponse(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(r.Context())
	logger.LogAttrs(ctx, slog.LevelInfo, timeoutMsg)
	ErrorResponse(w, r, http.StatusRequestTimeout, timeoutMsg)
}

func NotFoundResponse(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	logger := logging.LoggerFromContext(ctx)
	logger.LogAttrs(ctx, slog.LevelInfo, notFoundMsg)
	ErrorResponse(w, r, http.StatusNotFound, notFoundMsg)
}
