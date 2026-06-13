// Package response provides utilities for formatting HTTP responses.
// mapping domain errors to appropriate HTTP status codes and logging.
package response

import (
	"encoding/json"
	"errors"

	"fmt"
	"net/http"

	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/logger"
	"go.uber.org/zap"
)

// HTTPResponseHandler contains an http.ResponseWriter to provide standardized
// JSON formatting and automated error logging capabilities.
type HTTPResponseHandler struct {
	rw  http.ResponseWriter
	log *logger.Logger
}

// NewHTTPResponseHandler creates a new instance of HTTPResponseHandler.
func NewHTTPResponseHandler(
	w http.ResponseWriter,
	log *logger.Logger,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		rw:  w,
		log: log,
	}
}

// JSONResponse serializes the response body to JSON, sets the HTTP status code
// and content-type and logs an error if the encoding process fails.
func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.Header().Set("Content-Type", "application/json")

	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.JSONResponse(response, statusCode)
}

// ErrorResponse maps sentinel errors to the correct
// HTTP status codes and logging level, ensuring uniform error handling across the app.
func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, errs.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Debug
	case errors.Is(err, errs.ErrBadArgument):
		statusCode = http.StatusServiceUnavailable
		logFunc = h.log.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

// PanicResponse logs the recovered panic information
// and sends a 500 status code response to the client.
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)
	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		http.StatusInternalServerError,
		err,
		msg,
	)
}
