// Package ctxutil prevents context key collsions.
package ctxutil

import (
	"context"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/logger"
)

type ctxKey string

const (
	loggerKey    ctxKey = "logger"
	requestIDKey ctxKey = "request_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		panic("no request ID in context")
	}

	return requestID
}

func WithLogger(ctx context.Context, log *logger.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func Logger(ctx context.Context) *logger.Logger {
	logger, ok := ctx.Value(loggerKey).(*logger.Logger)
	if !ok {
		panic("no logger in context")
	}

	return logger
}
