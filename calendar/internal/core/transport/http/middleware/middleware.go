// Package middleware provides HTTP interceptors that contain or wrap
// standard handlers with common cross-cutting logic (e.g. logging).
package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// ChainMiddleware builds a single http.Handler from a chain of middleware functions.
// It applies the middleware in reverse order so they execute in the exact order provided.
func ChainMiddleware(h http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
