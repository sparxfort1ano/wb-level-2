package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/logger"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/middleware"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/server"
	eventsStorage "github.com/sparxfort1ano/wb-level-2/calendar/internal/features/events/repository/storage"
	eventsService "github.com/sparxfort1ano/wb-level-2/calendar/internal/features/events/service"
	eventsTransport "github.com/sparxfort1ano/wb-level-2/calendar/internal/features/events/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Printf("failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("init feature", zap.String("feature", "events"))
	eventsRepository := eventsStorage.NewEventsRepository()
	eventsService := eventsService.NewEventsService(eventsRepository)
	eventsHTTPHandler := eventsTransport.NewEventsHTTPHandler(eventsService)

	logger.Debug("init HTTP server")

	httpConfig := server.NewConfigMust()
	httpServer := server.NewHTTPServer(
		httpConfig,
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.Panic(),
	)

	apiVersionRouterV1 := server.NewAPIVersionRouter(
		server.APIVersion1,
	)
	apiVersionRouterV1.RegisterRoutes(eventsHTTPHandler.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error")
	}
}
