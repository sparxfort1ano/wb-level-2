// Package server provides utilities to configure, run and gracefully shutdown
// the main HTTP server, along with API versioning and routing mechanisms.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/logger"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

// HTTPServer contains the main HTTP multiplexer, global middleware chain,
// server configurastion and the application logger.
type HTTPServer struct {
	mux        *http.ServeMux
	cfg        config
	log        *logger.Logger
	middleware []middleware.Middleware
}

// NewHTTPServer creates a new instance of HTTPServer.
func NewHTTPServer(
	cfg config,
	log *logger.Logger,
	middleware ...middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		cfg:        cfg,
		log:        log,
		middleware: middleware,
	}
}

// RegisterAPIRouters mounts version-specific sub-routers onto the main HTTP server.
// It automatically wraps each router with version's specific middleware.
func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(prefix+"/", http.StripPrefix(prefix, router.withMiddleware()))
	}
}

// RegisterRoutes mount individual, top-level HTTP routes onto the main server multiplexer.
// It automatically wraps each handler with the route's specific middleware.
func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		s.mux.Handle(pattern, route.withMiddleware())
	}
}

// Run starts the HTTP server, supporting its graceful shutdown.
func (s *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("HTTP server start", zap.String("addr", server.Addr))

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutDownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutDownCtx); err != nil {
			server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
