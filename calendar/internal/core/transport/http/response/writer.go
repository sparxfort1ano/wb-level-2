package response

import "net/http"

const (
	StatusCodeUninitialized = -1
)

// ResponseWriter is a custom decorator around the http.ResponseWriter.
// It intercepts and stores the HTTP status code so it can be read later.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new instance of ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

// WriteHeader decorates the underlying WriteHeader method to capture an store
// the status code in memory before sending it to the client.
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// StatusCode retrieves the captured HTTP status code.
// It status code is not set, it panics. Panic is allowed:
// it is not correct to send an HTTP response without status code set.
func (rw *ResponseWriter) StatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}

	return rw.statusCode
}
