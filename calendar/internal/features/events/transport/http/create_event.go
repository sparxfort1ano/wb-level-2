package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/ctxutil"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/request"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/response"
)

// CreateEventRequest represents the incoming URL body for creating an event.
type CreateEventRequest struct {
	Event  string `form:"event" validate:"required,min=1"`
	Date   string `form:"date" validate:"required"`
	UserID int    `form:"user_id" validate:"gte=0"`
}

// CreateEventResponse represents the outgoing JSON body for creating a event.
type CreateEventResponse resultEventResponse

// CreateEvent processes the HTTP to register a new event.
func (h *EventsHTTPHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctxutil.Logger(ctx)
	responseHandler := response.NewHTTPResponseHandler(w, log)

	var req CreateEventRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	eventDomain, err := domainFromDTO(req)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'date'",
		)
		return
	}

	eventDomain, err = h.eventsService.CreateEvent(ctx, eventDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create event",
		)
		return
	}

	response := CreateEventResponse(resultEventDTO(eventDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func domainFromDTO(dto CreateEventRequest) (domain.Event, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, dto.Date)
	if err != nil {
		return domain.Event{}, fmt.Errorf(
			"value %s is not a valid date: %v: %w",
			dto.Date,
			err,
			errs.ErrInvalidArgument,
		)
	}

	return domain.NewEventUnitialized(
		dto.Event,
		date,
		dto.UserID,
	), nil
}