package server

import (
	"net/http"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/middleware"
)

// Route binds an HTTP method and URI pattern to specific handler.
// It also contains route-level middlewares that are applied exclusively to this endpoint.
type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []middleware.Middleware
}

func (r *Route) withMiddleware() http.Handler {
	return middleware.ChainMiddleware(
		r.Handler, 
		r.Middleware...,
	)
}