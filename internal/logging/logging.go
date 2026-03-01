package logging

import (
	"context"
	"log/slog"
	"regexp"
	"strings"
)

type ContextKey string

const LoggerKey ContextKey = "logger"

// WithLogger embeds a logger in the given context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// LoggerFromContext attempts to extract an embedded logger from the
// given context. If no context is found, it returns the default logger
// registered for the application.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

func MinifySQL(query string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(query, " "))
}

// ContextLogger accepts a context.Context and a slog.Attr, and returns a new enriched context
// and a logger object.
//
// It's a convenience function wrapping the functionality of LoggerFromContext and WithLogger
// in a single call.
func ContextLogger(ctx context.Context, attr slog.Attr) (context.Context, *slog.Logger) {
	logger := LoggerFromContext(ctx).With(attr)
	return WithLogger(ctx, logger), logger
}
