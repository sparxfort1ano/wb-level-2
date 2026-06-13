package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/ctxutil"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/logger"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

// RequestID ensures every request has a unique identifier.
// It reads the X-Request-ID head from the client or generates a new UUID.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			w.Header().Set(requestIDHeader, requestID)

			ctx := ctxutil.WithRequestID(r.Context(), requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logger injects a context-aware zap.Logger into the request context.
// It binds the request_id and URL to the logger consistent structured logging.
func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := ctxutil.RequestID(r.Context())

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := ctxutil.WithLogger(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Trace logs the start and completion of an HTTP request handling.
// It prevents the server from crashing and returns a graceful 500 response.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := ctxutil.Logger(r.Context())
			rw := response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.StatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

// Panic recovers from unexpected panics during HTTP request handling.
// It prevents the server from crashing and returns a graceful 500 response.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					log := ctxutil.Logger(r.Context())
					responseHandler := response.NewHTTPResponseHandler(w, log)

					responseHandler.PanicResponse(
						p,
						"during handling HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
