// Package http acts as the transport layer for the Event feature.
// It is responsible for parsing HTTP requests, formatting responses and routing.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/server"
)

// EventsHTTPHandler handles HTTP requests related to events management.
// All of its methods delegate the logic to the service layer
// after decoding the payload. At the end they write an JSON response.
type EventsHTTPHandler struct {
	eventsService eventsService
}

// eventService defines the contract that decouples the HTTP transport layer
// from the underlying domain logic.
type eventsService interface {
	// CreateEvent enforces business rules using validation on the event domain.
	CreateEvent(
		ctx context.Context,
		domain domain.Event,
	) (domain.Event, error)

	DeleteEvent(
		ctx context.Context,
		id int,
	) error

	// UpdateEvent requests to get the event by the event ID.
	// Then enforces business rules on the patching Event object.
	UpdateEvent(
		ctx context.Context,
		id int,
		patch domain.EventPatch,
	) (domain.Event, error)

	// GetEventsForDay calculates from and to parameters to get one day time range.
	GetEventsForDay(
		ctx context.Context,
		userID *int,
		date *time.Time,
	) ([]domain.Event, error)

	// GetEventsForWeek calculates from and to parameters to get one week time range
	// using the date parameter's time.Week().
	GetEventsForWeek(
		ctx context.Context,
		userID *int,
		date *time.Time,
	) ([]domain.Event, error)

	// GetEventsForMonth calculates from and to parameters to get one week time range
	// using the date parameter's time.Month().
	GetEventsForMonth(
		ctx context.Context,
		userID *int,
		date *time.Time,
	) ([]domain.Event, error)
}

// NewEventsHTTPHandler creates a new instance of EventsHTTPHandler.
func NewEventsHTTPHandler(
	eventsService eventsService,
) *EventsHTTPHandler {
	return &EventsHTTPHandler{
		eventsService: eventsService,
	}
}

// Routes returns a list of HTTP routes to be registered in the server router.
func (h *EventsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/create_event",
			Handler: h.CreateEvent,
		},
		{
			Method:  http.MethodPost,
			Path:    "/delete_event",
			Handler: h.DeleteEvent,
		},
		{
			Method:  http.MethodPost,
			Path:    "/update_event",
			Handler: h.UpdateEvent,
		},
		{
			Method:  http.MethodGet,
			Path:    "/events_for_day",
			Handler: h.GetEventsForDay,
		},
		{
			Method:  http.MethodGet,
			Path:    "/events_for_week",
			Handler: h.GetEventsForWeek,
		},
		{
			Method:  http.MethodGet,
			Path:    "/events_for_month",
			Handler: h.GetEventsForMonth,
		},
	}
}
